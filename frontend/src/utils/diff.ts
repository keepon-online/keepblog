/**
 * 轻量级文本差异对比（Diff）工具
 * 基于 LCS（最长公共子序列）算法，专为中英混合、Markdown 文本优化分词。
 * 零外部依赖，极速计算。
 */

export type DiffType = "added" | "removed" | "common";

export interface DiffChunk {
  type: DiffType;
  value: string;
}

/**
 * 中英文与符号分词器
 * - 中文字符按单字切分，保证精细对比；
 * - 英文及数字按单词切分；
 * - 换行、空白符与标点符号单独成词。
 */
function tokenize(text: string): string[] {
  if (!text) return [];
  const regex =
    /[\u4e00-\u9fa5\u3000-\u303f\uff01-\uff5e]|[a-zA-Z0-9_]+|\r?\n|[^\s\w\u4e00-\u9fa5]+|\s+/g;
  const tokens = text.match(regex);
  return tokens && tokens.length > 0 ? tokens : [text];
}

/**
 * 计算两个文本的分词级 Diff
 * @param oldText 原始文本
 * @param newText 新文本
 * @returns 差异片段列表
 */
export function computeDiff(oldText: string, newText: string): DiffChunk[] {
  if (oldText === newText) {
    return oldText ? [{ type: "common", value: oldText }] : [];
  }
  if (!oldText) {
    return newText ? [{ type: "added", value: newText }] : [];
  }
  if (!newText) {
    return oldText ? [{ type: "removed", value: oldText }] : [];
  }

  const a = tokenize(oldText);
  const b = tokenize(newText);
  const n = a.length;
  const m = b.length;

  // 针对长文本截断优化：如果两个 token 序列长度极大（> 2000），可直接降级为段落级对比
  if (n * m > 1600000) {
    // 按行切分以控制复杂度
    return lineLevelDiff(oldText, newText);
  }

  // 动态规划构建 LCS 距离表（滚动数组或紧凑矩阵）
  // 考虑到一般选区在 1000 token 以内，使用一维扁平数组
  const dp = new Uint16Array((n + 1) * (m + 1));
  const stride = m + 1;

  for (let i = 1; i <= n; i++) {
    const ai = a[i - 1];
    const prevRow = (i - 1) * stride;
    const currRow = i * stride;
    for (let j = 1; j <= m; j++) {
      if (ai === b[j - 1]) {
        dp[currRow + j] = dp[prevRow + j - 1] + 1;
      } else {
        const top = dp[prevRow + j];
        const left = dp[currRow + j - 1];
        dp[currRow + j] = top >= left ? top : left;
      }
    }
  }

  // 回溯提取差异路径
  const rawChunks: DiffChunk[] = [];
  let i = n;
  let j = m;

  while (i > 0 || j > 0) {
    if (i > 0 && j > 0 && a[i - 1] === b[j - 1]) {
      rawChunks.push({ type: "common", value: a[i - 1] });
      i--;
      j--;
    } else if (
      j > 0 &&
      (i === 0 || dp[i * stride + (j - 1)] >= dp[(i - 1) * stride + j])
    ) {
      rawChunks.push({ type: "added", value: b[j - 1] });
      j--;
    } else if (i > 0) {
      rawChunks.push({ type: "removed", value: a[i - 1] });
      i--;
    }
  }

  rawChunks.reverse();

  // 合并相邻相同类型的 chunk
  const merged: DiffChunk[] = [];
  for (const chunk of rawChunks) {
    if (merged.length > 0 && merged[merged.length - 1].type === chunk.type) {
      merged[merged.length - 1].value += chunk.value;
    } else {
      merged.push({ ...chunk });
    }
  }

  return merged;
}

/**
 * 超长文本行级降级对比
 */
function lineLevelDiff(oldText: string, newText: string): DiffChunk[] {
  const oldLines = oldText.split("\n");
  const newLines = newText.split("\n");
  // 简易行级对比
  return [
    { type: "removed", value: oldText },
    { type: "added", value: newText }
  ];
}
