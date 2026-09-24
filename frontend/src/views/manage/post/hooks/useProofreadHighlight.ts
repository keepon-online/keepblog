import { ref, type Ref } from "vue";
import { RangeSetBuilder, StateEffect, StateField } from "@codemirror/state";
import { Decoration, type DecorationSet, EditorView } from "@codemirror/view";
import type { ProofreadItem } from "@/utils/proofread";
import { message } from "@/utils/message";

export interface ProofreadLocation {
  item: ProofreadItem;
  from: number;
  to: number;
}

export const setProofreadEffect = StateEffect.define<ProofreadLocation[]>();
export const clearProofreadEffect = StateEffect.define<void>();

export const proofreadStateField = StateField.define<DecorationSet>({
  create() {
    return Decoration.none;
  },
  update(underlines, tr) {
    underlines = underlines.map(tr.changes);

    for (const effect of tr.effects) {
      if (effect.is(clearProofreadEffect)) {
        return Decoration.none;
      }
      if (effect.is(setProofreadEffect)) {
        const locations = effect.value;
        const builder = new RangeSetBuilder<Decoration>();

        // 必须按起始位置升序排序
        const sorted = [...locations].sort((a, b) => a.from - b.from);

        for (const loc of sorted) {
          if (loc.from >= loc.to || loc.to > tr.newDoc.length) continue;
          if (loc.item.status && loc.item.status !== "pending") continue;

          let className = "cm-proofread-typo";
          if (loc.item.type.includes("语法") || loc.item.type.includes("语病") || loc.item.type.includes("标点")) {
            className = "cm-proofread-grammar";
          } else if (loc.item.type.includes("风格") || loc.item.type.includes("表达") || loc.item.type.includes("冗余")) {
            className = "cm-proofread-style";
          }

          const mark = Decoration.mark({
            class: className,
            attributes: {
              "data-proofread-id": loc.item.id,
              "data-proofread-original": loc.item.original,
              "data-proofread-suggestion": loc.item.suggestion,
              "data-proofread-reason": loc.item.reason || "",
              "data-proofread-type": loc.item.type
            }
          });

          builder.add(loc.from, loc.to, mark);
        }

        return builder.finish();
      }
    }

    return underlines;
  },
  provide: f => EditorView.decorations.from(f)
});

export function useProofreadHighlight(
  editorRef: Ref<any>,
  proofreadItems: Ref<ProofreadItem[]>,
  onContentUpdated?: (newContent: string) => void
) {
  const hoverCardVisible = ref(false);
  const hoverCardItem = ref<ProofreadItem | null>(null);
  const hoverCardPos = ref({ top: 0, left: 0 });
  const activeLocation = ref<ProofreadLocation | null>(null);

  // 在文档中检索所有待处理校对项的精确位置
  const computeLocations = (view: EditorView, items: ProofreadItem[]): ProofreadLocation[] => {
    const docText = view.state.doc.toString();
    const locations: ProofreadLocation[] = [];
    const pendingItems = items.filter(i => !i.status || i.status === "pending");

    for (const item of pendingItems) {
      if (!item.original) continue;
      let startIdx = 0;
      let foundIdx = docText.indexOf(item.original, startIdx);

      // 仅标记第一个匹配，或匹配所有
      while (foundIdx !== -1) {
        locations.push({
          item,
          from: foundIdx,
          to: foundIdx + item.original.length
        });
        startIdx = foundIdx + item.original.length;
        // 为防重叠只匹配第一个有效位置
        break;
      }

      if (foundIdx === -1) {
        item.status = "not_found";
      }
    }

    return locations;
  };

  // 刷新波浪线高亮
  const refreshDecorations = () => {
    const view = editorRef.value?.getEditorView?.();
    if (!view) return;

    const locations = computeLocations(view, proofreadItems.value);
    view.dispatch({
      effects: setProofreadEffect.of(locations)
    });
  };

  // 清除所有高亮
  const clearDecorations = () => {
    const view = editorRef.value?.getEditorView?.();
    if (!view) return;
    view.dispatch({
      effects: clearProofreadEffect.of()
    });
    hoverCardVisible.value = false;
  };

  // 定位某一项到编辑器视口中央
  const locateProofreadItem = (item: ProofreadItem) => {
    const view = editorRef.value?.getEditorView?.();
    if (!view) return;

    const docText = view.state.doc.toString();
    const index = docText.indexOf(item.original);
    if (index === -1) {
      message(`在正文中未找到「${item.original}」，可能已被修改`, {
        type: "warning"
      });
      item.status = "not_found";
      return;
    }

    const from = index;
    const to = index + item.original.length;
    view.dispatch({
      selection: { anchor: from, head: to },
      scrollIntoView: true
    });

    const coords = view.coordsAtPos(from);
    if (coords) {
      showHoverCard(item, {
        top: coords.bottom + 8,
        left: Math.max(10, coords.left - 40)
      }, { item, from, to });
    }

    const lineNo = docText.slice(0, from).split("\n").length;
    message(`已在正文第 ${lineNo} 行定位并聚焦`, { type: "info" });
  };

  // 采纳单个修改
  const applyProofreadItem = (item: ProofreadItem) => {
    const view = editorRef.value?.getEditorView?.();
    if (!view) return;

    const docText = view.state.doc.toString();
    const index = docText.indexOf(item.original);
    if (index === -1) {
      message(`在正文中未找到「${item.original}」，无法采纳`, {
        type: "warning"
      });
      item.status = "not_found";
      hoverCardVisible.value = false;
      return;
    }

    const from = index;
    const to = index + item.original.length;
    view.dispatch({
      changes: { from, to, insert: item.suggestion },
      selection: { anchor: from + item.suggestion.length },
      scrollIntoView: true
    });

    item.status = "applied";
    hoverCardVisible.value = false;

    const newContent = view.state.doc.toString();
    onContentUpdated?.(newContent);

    message(`已采纳：「${item.original}」→「${item.suggestion}」`, {
      type: "success"
    });

    refreshDecorations();
  };

  // 忽略单个修改
  const ignoreProofreadItem = (item: ProofreadItem) => {
    item.status = "ignored";
    hoverCardVisible.value = false;
    refreshDecorations();
  };

  // 一键采纳全部待处理项
  const applyAllProofreadItems = () => {
    const view = editorRef.value?.getEditorView?.();
    if (!view) return;

    const pending = proofreadItems.value.filter(i => i.status === "pending" || !i.status);
    if (pending.length === 0) return;

    const docText = view.state.doc.toString();
    const replacements: {
      item: ProofreadItem;
      from: number;
      to: number;
      insert: string;
    }[] = [];

    for (const item of pending) {
      const idx = docText.indexOf(item.original);
      if (idx !== -1) {
        replacements.push({
          item,
          from: idx,
          to: idx + item.original.length,
          insert: item.suggestion
        });
      } else {
        item.status = "not_found";
      }
    }

    if (replacements.length === 0) {
      message("未在正文中找到可采纳的内容", { type: "warning" });
      return;
    }

    // 从后往前替换，保证前面的 offset 不失效
    replacements.sort((a, b) => b.from - a.from);

    const changes = replacements.map(r => ({
      from: r.from,
      to: r.to,
      insert: r.insert
    }));

    view.dispatch({
      changes,
      scrollIntoView: true
    });

    for (const r of replacements) {
      r.item.status = "applied";
    }

    const newContent = view.state.doc.toString();
    onContentUpdated?.(newContent);

    message(`已成功一键采纳 ${replacements.length} 处修改`, {
      type: "success"
    });

    clearDecorations();
  };

  const showHoverCard = (
    item: ProofreadItem,
    pos: { top: number; left: number },
    loc?: ProofreadLocation
  ) => {
    hoverCardItem.value = item;
    hoverCardPos.value = pos;
    activeLocation.value = loc || null;
    hoverCardVisible.value = true;
  };

  const hideHoverCard = () => {
    hoverCardVisible.value = false;
  };

  return {
    hoverCardVisible,
    hoverCardItem,
    hoverCardPos,
    activeLocation,
    refreshDecorations,
    clearDecorations,
    locateProofreadItem,
    applyProofreadItem,
    ignoreProofreadItem,
    applyAllProofreadItems,
    showHoverCard,
    hideHoverCard
  };
}
