// Package ai 提供 OpenAI 兼容 Chat Completions 协议的流式客户端。
// GLM（open.bigmodel.cn）、DeepSeek、通义、本地 Ollama 等均兼容该协议，
// 一个实现通吃，提供商由 config.ai 的 baseURL/model 决定。
package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Message 对话消息（OpenAI 兼容格式）。
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Usage token 用量（部分兼容服务在流式末块返回，缺失时为 0）。
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Client OpenAI 兼容客户端。零值不可用，经 New 创建。
type Client struct {
	baseURL string
	apiKey  string
	model   string
	http    *http.Client
}

// New 创建客户端。timeout 为单次生成请求的总超时（流式含等待时间，需宽于普通接口）。
func New(baseURL, apiKey, model string, timeoutSec int) *Client {
	if timeoutSec <= 0 {
		timeoutSec = 120
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		model:   model,
		http:    &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
	}
}

// StreamRequest 流式请求参数。
type StreamRequest struct {
	Messages  []Message
	MaxTokens int
	// Temperature 采样温度，0 使用服务端默认
	Temperature float64
	// OnDelta 正文增量回调（nil 忽略）
	OnDelta func(delta string)
	// OnReasoning 思考增量回调（推理模型正文前的 reasoning_content，nil 忽略）
	OnReasoning func(delta string)
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream"`
}

// chatChunk SSE data 块的反序列化目标。兼容服务在末块可能带 usage。
// 推理模型（如 GLM 思考系列）在正文前先流式输出 reasoning_content，
// 单独回调给上层展示"思考中"，不计入正文。
type chatChunk struct {
	Choices []struct {
		Delta struct {
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
	Usage *Usage `json:"usage"`
}

// StreamChat 发起流式对话。正文/思考增量经 req.OnDelta / req.OnReasoning
// 回调（内容顺序保证），返回最终 usage（缺失时为零值）。ctx 取消时中止
// 上游请求并返回 ctx.Err()。
func (c *Client) StreamChat(ctx context.Context, req StreamRequest) (Usage, error) {
	var zero Usage
	body, err := json.Marshal(chatRequest{
		Model:       c.model,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      true,
	})
	if err != nil {
		return zero, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return zero, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return zero, fmt.Errorf("request upstream: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return zero, fmt.Errorf("上游返回 %d: %s", resp.StatusCode, strings.TrimSpace(string(errBody)))
	}

	// SSE 解析：事件以空行分隔，data: 前缀携带 JSON；[DONE] 结束。
	usage := Usage{}
	scanner := bufio.NewScanner(resp.Body)
	// 单行兜底上限：内容块远小于此，超长行按错误处理而不是静默丢弃
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "" || payload == "[DONE]" {
			continue
		}
		var chunk chatChunk
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			// 单块损坏不致命：兼容服务可能插入非标准行，跳过继续
			continue
		}
		if chunk.Usage != nil {
			usage = *chunk.Usage
		}
		if len(chunk.Choices) > 0 {
			if d := chunk.Choices[0].Delta.Content; d != "" && req.OnDelta != nil {
				req.OnDelta(d)
			}
			if d := chunk.Choices[0].Delta.ReasoningContent; d != "" && req.OnReasoning != nil {
				req.OnReasoning(d)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return usage, fmt.Errorf("read stream: %w", err)
	}
	return usage, nil
}
