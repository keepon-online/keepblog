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
  // 使用前瞻断言 (?=\r?\n|$) 确保 blockEnd 精确截止在闭合反引号处，
  // 避免吞掉后续换行符导致替换回写后正文段落粘连
  const blockRegex =
    /(?:^|\n)(```([a-zA-Z0-9_\-\+\.]*)[ \t]*\r?\n([\s\S]*?)(?:\r?\n```|$))(?=\r?\n|$)/g;
  let match: RegExpExecArray | null;

  while ((match = blockRegex.exec(doc)) !== null) {
    const rawMatch = match[0];
    const leadingNewline = rawMatch.startsWith("\n") ? 1 : 0;
    const blockStart = match.index + leadingNewline;
    const blockContent = match[1];
    const blockEnd = blockStart + blockContent.length;

    // 检查光标或选区与该代码块是否有交集
    const cursorInside =
      (from >= blockStart && from <= blockEnd) ||
      (to >= blockStart && to <= blockEnd) ||
      (from <= blockStart && to >= blockEnd);

    if (cursorInside) {
      return {
        isCodeBlock: true,
        lang: match[2]?.trim() || "",
        code: match[3] ?? "",
        fullBlock: blockContent,
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
  const trimmed = code.trim();
  if (!trimmed) return "";

  // 1. Mermaid 图表
  if (
    /^(graph\s+(TD|LR|TB|BT)|sequenceDiagram|classDiagram|erDiagram|gantt|pie|flowchart\s+(TD|LR|TB|BT)|stateDiagram)/m.test(
      trimmed
    )
  ) {
    return "mermaid";
  }

  // 2. JSON
  if (
    (/^\{[\s\S]*\}$/.test(trimmed) || /^\[[\s\S]*\]$/.test(trimmed)) &&
    (/"[a-zA-Z0-9_\-]+"\s*:\s*/.test(trimmed) ||
      /:\s*("[^"]*"|\d+|true|false|null)/.test(trimmed))
  ) {
    return "json";
  }

  // 3. YAML
  if (
    /^[a-zA-Z0-9_\-]+:\s*([^\n]*|$)/m.test(trimmed) &&
    !trimmed.includes("{") &&
    !trimmed.includes("}")
  ) {
    if (/\n\s+[a-zA-Z0-9_\-]+:\s*/.test(trimmed)) {
      return "yaml";
    }
  }

  // 4. HTML / Vue 模板
  if (/<template>|<script\s+lang=/.test(trimmed)) {
    return "vue";
  }
  if (
    /<!DOCTYPE html>|<html[\s>]|<div[\s>]|<head[\s>]|<body[\s>]|<span[\s>]|<p[\s>]/i.test(
      trimmed
    )
  ) {
    return "html";
  }

  // 5. CSS / SCSS
  if (
    /(@media|@keyframes|\b(display|margin|padding|border-radius|font-size|box-shadow)\s*:)/i.test(
      trimmed
    ) &&
    trimmed.includes("{") &&
    trimmed.includes("}")
  ) {
    return "css";
  }

  // 6. Go
  if (
    /\bpackage\s+[a-z0-9_]+|\bfunc\s+[A-Za-z0-9_]+\s*\(|:=|\bchan\s+[a-zA-Z0-9_]+|\bdefer\s+/.test(
      trimmed
    )
  ) {
    return "go";
  }

  // 7. Rust
  if (
    /fn\s+main\(\)|\blet\s+mut\s+|\bimpl\b|\bmatch\s+[a-zA-Z0-9_]+|\bpub\s+fn\b/.test(
      trimmed
    )
  ) {
    return "rust";
  }

  // 8. Python
  if (
    /\bdef\s+[a-z0-9_]+\s*\(|\bimport\s+[a-z0-9_]+|\bself\b|\belif\b|\bif\s+__name__\s*==\s*['"]__main__['"]/.test(
      trimmed
    )
  ) {
    return "python";
  }

  // 9. TypeScript / JavaScript
  if (
    /\binterface\s+[A-Z]|\btype\s+[A-Z][a-zA-Z0-9_]*\s*=|:\s*(string|number|boolean|any)\b/.test(
      trimmed
    )
  ) {
    return "typescript";
  }
  if (
    /\bconst\s+[a-zA-Z0-9_]+\s*=|\bconsole\.log\(|\bfunction\s+[a-zA-Z0-9_]+\s*\(|=>\s*\{/.test(
      trimmed
    )
  ) {
    return "javascript";
  }

  // 10. Java
  if (
    /\b(public|private|protected)\s+(class|interface|void|int|String|boolean)\b|System\.out\.println\(|@Override/.test(
      trimmed
    )
  ) {
    return "java";
  }

  // 11. C / C++
  if (
    /#include\s+[<"][a-z0-9_\.]+[>"]|\bstd::cout\b|\bprintf\s*\(|\bint\s+main\s*\(/.test(
      trimmed
    )
  ) {
    return "cpp";
  }

  // 12. SQL
  if (
    /\b(SELECT|INSERT\s+INTO|UPDATE|DELETE\s+FROM|FROM|WHERE|JOIN|GROUP\s+BY|ORDER\s+BY)\b/i.test(
      trimmed
    )
  ) {
    return "sql";
  }

  // 13. Bash / Shell
  if (
    /\b(docker|kubectl|npm|pnpm|yarn|cargo|git|go run|go build|curl|chmod|chown|systemctl|grep|awk|sed)\b/.test(
      trimmed
    ) ||
    /^#!\/bin\/(bash|sh)/m.test(trimmed)
  ) {
    return "bash";
  }

  return "";
}

/**
 * 确保代码被对应语言的反引号围栏包裹
 */
export function wrapCodeFence(code: string, lang = ""): string {
  const trimmed = code.trim();
  const match = trimmed.match(/^```([a-zA-Z0-9_\-\+\.]*)\s*\n([\s\S]*?)\n```$/);
  if (match) {
    if (match[1] || !lang) {
      return trimmed;
    }
    return `\`\`\`${lang}\n${match[2]}\n\`\`\``;
  }
  return `\`\`\`${lang}\n${trimmed}\n\`\`\``;
}
