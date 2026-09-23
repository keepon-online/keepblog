package ai

import (
	"encoding/json"
	"net/http"

	"gitee.com/jieepre/keepblog/internal/pkg/core"
	aisvc "gitee.com/jieepre/keepblog/internal/service/ai"
	aiclient "gitee.com/jieepre/keepblog/pkg/ai"
	"gitee.com/jieepre/keepblog/pkg/result"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*core.Context
}

// Status 探测 AI 是否可用（未配置 apiKey 时前端据此隐藏所有 AI 入口）。
func (h *Handler) Status(c *gin.Context) {
	result.Ok(c, h.Service.AIService.Status())
}

// sseEvent 输出一条 SSE 事件并立即刷出。
func sseEvent(c *gin.Context, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	_, _ = c.Writer.Write(append([]byte("data: "), append(data, '\n', '\n')...))
	c.Writer.Flush()
}

// Edit AI 编辑任务，SSE 流式返回。
// 事件格式：{"type":"delta","content":"..."} / {"type":"done","usage":{...}} /
// {"type":"error","message":"..."}。客户端断开时 ctx 取消，上游请求随之中止。
func (h *Handler) Edit(c *gin.Context) {
	var req aisvc.EditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, "参数错误: "+err.Error())
		return
	}

	msgs, err := aisvc.BuildMessages(req,
		h.Service.AIService.TemplateOverride(aisvc.TemplateKey(req.Task, req.Mode)),
		h.Service.AIService.EffectiveStyleHint())
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	if err := h.Service.AIService.AllowCall(); err != nil {
		result.Error(c, err.Error())
		return
	}
	client, err := h.Service.AIService.ClientFor(req.Model)
	if err != nil {
		result.Error(c, err.Error())
		return
	}

	// SSE 响应头：no-gzip 防中间件压缩缓冲，flush 由 sseEvent 显式触发
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Flush()

	ctx := c.Request.Context()
	cfgMax := h.Service.AIService.MaxTokensFor(req.Model)

	usage, err := client.StreamChat(ctx, aiclient.StreamRequest{
		Messages:  msgs,
		MaxTokens: cfgMax,
		OnDelta: func(delta string) {
			sseEvent(c, gin.H{"type": "delta", "content": delta})
		},
		// 推理模型的思考过程单独透传，前端展示为"思考中"，
		// 避免正文输出前的长时间静默被当成卡死
		OnReasoning: func(delta string) {
			sseEvent(c, gin.H{"type": "reasoning", "content": delta})
		},
	})
	if err != nil {
		// 客户端已断开时写事件无人消费，直接返回即可
		sseEvent(c, gin.H{"type": "error", "message": err.Error()})
		return
	}
	h.Service.AIService.RecordUsage(req.Task, h.Service.AIService.ActualModel(req.Model), usage)
	sseEvent(c, gin.H{"type": "done", "usage": usage})
}

// Templates 模板管理：内置默认与用户覆盖同屏返回。
func (h *Handler) Templates(c *gin.Context) {
	result.Ok(c, h.Service.AIService.TemplateItems())
}

// SaveTemplate 保存模板覆盖；content 为空即恢复默认。
func (h *Handler) SaveTemplate(c *gin.Context) {
	var body struct {
		Key     string `json:"key" binding:"required"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		result.Error(c, "参数错误")
		return
	}
	if err := h.Service.AIService.SaveTemplate(body.Key, body.Content); err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}

// AIConfig 读当前生效的 AI 配置（apiKey 只回掩码）。
func (h *Handler) AIConfig(c *gin.Context) {
	result.Ok(c, h.Service.AIService.ConfigInfo())
}

// SaveAIConfig 保存配置覆盖层（字段留空回落 config.yaml 基线）。
func (h *Handler) SaveAIConfig(c *gin.Context) {
	var req aisvc.SaveAIConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, "参数错误: "+err.Error())
		return
	}
	if err := h.Service.AIService.SaveAIConfig(req); err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
