/**
 * 全文校对结构化解析工具
 * 支持兼容主流大模型在校对任务下的多种输出形态（直角引号、双引号、箭头、序号、标点等）
 */

export interface ProofreadItem {
  id: string;
  original: string; // 待修改原文
  suggestion: string; // 建议修改后内容
  type: string; // 问题分类：错别字 | 标点符号 | 语病 | 格式规范 | 优化建议
  reason?: string; // 错误原因或详细说明
  status: "pending" | "applied" | "ignored" | "not_found";
}

/**
 * 智能分类提取
 */
function classifyType(rawType: string): { type: string; reason?: string } {
  if (!rawType) return { type: "优化建议" };
  const trimmed = rawType.trim();

  let type = "优化建议";
  let reason: string | undefined;

  // 检查是否包含冒号分隔的详情，例如 "错别字：应为再次的再"
  const colonIdx = trimmed.indexOf("：");
  const colonAsciiIdx = trimmed.indexOf(":");
  const splitIdx = colonIdx > -1 ? colonIdx : colonAsciiIdx;

  let mainType = trimmed;
  if (splitIdx > -1) {
    mainType = trimmed.slice(0, splitIdx).trim();
    reason = trimmed.slice(splitIdx + 1).trim();
  }

  if (/错[别]?字|别字|错词/.test(mainType)) {
    type = "错别字";
  } else if (/标点|符号|顿号|逗号|句号|引号/.test(mainType)) {
    type = "标点符号";
  } else if (/语病|语意|搭配|主谓|语序|通顺/.test(mainType)) {
    type = "语病问题";
  } else if (/格式|规范|大小写|空格|排版/.test(mainType)) {
    type = "格式规范";
  } else {
    type = mainType || "优化建议";
  }

  return { type, reason: reason || (mainType !== type ? mainType : undefined) };
}

/**
 * 解析大模型返回的校对文本为结构化列表
 */
export function parseProofreadOutput(text: string): ProofreadItem[] {
  if (!text || text.trim() === "未发现问题") return [];

  const items: ProofreadItem[] = [];
  const lines = text.split("\n");

  // 正则1: 原文「...」→ 建议「...」（类型）
  // 兼容各种括号「」、""、“”、箭头 →、->、=>、建议为
  const regexStandard =
    /(?:^|\s)(?:原文|原句|待改|原)?\s*[「"“](.+?)[」"”]\s*(?:→|->|=>|建议为|改成|应改为|建议)\s*(?:建议)?\s*[「"“](.+?)[」"”](?:\s*[（(](.+?)[）)])?/;

  // 正则2: 1. 「...」 -> 「...」 (类型)
  const regexArrow =
    /(?:^|\s|\d+[\.、\s]|-|\*)\s*[「"“](.+?)[」"”]\s*(?:→|->|=>)\s*[「"“](.+?)[」"”](?:\s*[（(](.+?)[）)])?/;

  // 正则3: 原文：...，建议：...（类型）
  const regexColon =
    /(?:原文|原句)[：:]\s*(.+?)\s*[,，;；]\s*(?:建议|改为|应为)[：:]\s*(.+?)(?:\s*[（(](.+?)[）)]|$)/;

  let idCounter = 0;

  for (const rawLine of lines) {
    const line = rawLine.trim();
    if (!line) continue;
    if (line.includes("未发现问题") && lines.length === 1) return [];

    let match = line.match(regexStandard);
    if (!match) match = line.match(regexArrow);
    if (!match) match = line.match(regexColon);

    if (match) {
      const original = (match[1] || "").trim();
      const suggestion = (match[2] || "").trim();
      const rawType = (match[3] || "").trim();

      if (original && suggestion && original !== suggestion) {
        const { type, reason } = classifyType(rawType);
        items.push({
          id: `pr-${idCounter++}`,
          original,
          suggestion,
          type,
          reason,
          status: "pending"
        });
      }
    }
  }

  return items;
}
