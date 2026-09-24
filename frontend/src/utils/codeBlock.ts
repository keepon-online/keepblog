/**
 * Markdown 代码块识别与辅助工具
 */

export interface CodeBlockInfo {
  isCodeBlock: boolean;
  lang: string;
  code: string;
  fullBlock: string;
  from: number;
  to: number;
}

/**
 * 探测光标或选区是否位于 Markdown 代码块 (```...```) 内部，
 * 或者当前选区本身就是一个代码块。
 *
 * @param doc 全文文档内容
 * @param from 选区起始偏移（或者光标位置）
 * @param to 选区结束偏移（如果是光标则 from === to）
 */
export function detectCodeBlock(
  doc: string,
  from: number,
  to: number
): CodeBlockInfo | null {
  if (!doc) return null;

  // 1. 如果用户明确选中了一段文本，先看选中的文本本身是否就是一个或多个完整代码块
  const selectedText = doc.slice(from, to).trim();
  const directMatch = selectedText.match(
    /^```([a-zA-Z0-9_\-\+\.]*)\s*\n([\s\S]*?)\n```$/
  );
  if (directMatch) {
    return {
      isCodeBlock: true,
      lang: directMatch[1] || "",
      code: directMatch[2],
      fullBlock: selectedText,
      from,
      to
    };
  }

  // 2. 在全文中搜索所有 ``` 代码块，检查当前 [from, to] 是否落在某一个代码块范围内
  const blockRegex = /(?:^|\n)```([a-zA-Z0-9_\-\+\.]*)[ \t]*\r?\n([\s\S]*?)(?:\r?\n```(?:\r?\n|$)|$)/g;
  let match: RegExpExecArray | null;

  while ((match = blockRegex.exec(doc)) !== null) {
    const rawMatch = match[0];
    const leadingNewline = rawMatch.startsWith("\n") ? 1 : 0;
    const blockStart = match.index + leadingNewline;
    const blockEnd = match.index + rawMatch.length;

    // 检查光标或选区与该代码块是否有交集
    const cursorInside =
      (from >= blockStart && from <= blockEnd) ||
      (to >= blockStart && to <= blockEnd) ||
      (from <= blockStart && to >= blockEnd);

    if (cursorInside) {
      return {
        isCodeBlock: true,
        lang: match[1]?.trim() || "",
        code: match[2] ?? "",
        fullBlock: rawMatch.slice(leadingNewline).trim(),
        from: blockStart,
        to: blockEnd
      };
    }
  }

  // 3. 如果选中的是普通多行代码（例如有多行缩进或典型编程关键字，且超过 1 行）
  if (selectedText.includes("\n") && selectedText.length > 10) {
    const codeKeywords =
      /\b(func|function|const|let|var|class|interface|type|import|export|def|return|if|for|while|select|package|namespace)\b/;
    if (codeKeywords.test(selectedText)) {
      return {
        isCodeBlock: false,
        lang: guessLanguage(selectedText),
        code: selectedText,
        fullBlock: selectedText,
        from,
        to
      };
    }
  }

  return null;
}

/**
 * 依据代码特征猜测常见编程语言
 */
export function guessLanguage(code: string): string {
  if (/\bpackage\s+[a-z0-9_]+|\bfunc\s+[A-Za-z0-9_]+\s*\(/.test(code))
    return "go";
  if (/\bdef\s+[a-z0-9_]+\s*\(|\bimport\s+[a-z0-9_]+|\bself\b/.test(code))
    return "python";
  if (/\binterface\s+[A-Z]|\btype\s+[A-Z][a-zA-Z0-9_]*\s*=|:\s*(string|number|boolean|any)\b/.test(code))
    return "typescript";
  if (/\bconst\s+[a-zA-Z0-9_]+\s*=|\bconsole\.log\(/.test(code))
    return "javascript";
  if (/\b(SELECT|INSERT|UPDATE|DELETE|FROM|WHERE|JOIN)\b/i.test(code))
    return "sql";
  if (/\b(docker|kubectl|npm|pnpm|yarn|cargo|git|go run|curl)\b/.test(code))
    return "bash";
  if (/<template>|<script\s+lang=/.test(code))
    return "vue";
  if (/fn\s+main\(\)|let\s+mut\s+/.test(code))
    return "rust";
  return "";
}

/**
 * 确保代码被对应语言的反引号围栏包裹
 */
export function wrapCodeFence(code: string, lang = ""): string {
  const trimmed = code.trim();
  if (trimmed.startsWith("```") && trimmed.endsWith("```")) {
    return trimmed;
  }
  return `\`\`\`${lang}\n${trimmed}\n\`\`\``;
}
