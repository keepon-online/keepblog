/**
 * 代码块复制按钮
 *
 * 背景：主题 main.js 的 addHighlightTool 选择器是 figure.highlight（Hexo 渲染
 * 格式），但后端 goldmark+chroma 输出的是 <pre><code>，选择器匹配不到，
 * 导致主题原生的复制/折叠功能从未启用。这里用独立补丁给文章正文里的每个
 * 代码块加复制按钮，不改动 main.js 与后端渲染，与主题逻辑解耦。
 *
 * 时机：DOMContentLoaded 首次执行；PJAX 切换页面后通过 pjax:complete 重新
 * 扫描（避免重复加按钮：已处理的 pre 打 data-codecopy 标记）。
 *
 * 注：用 var 与 function 声明，避免 PJAX 重新执行外部脚本时 let/const
 * 重复声明抛 SyntaxError（与 global-config.js、head.html 内联脚本同策略）。
 */
;(function () {
    var SELECTOR = '#article-container pre';

    // 注入按钮样式（一次性，自包含，不污染 index.css）
    if (!document.getElementById('code-copy-style')) {
        var style = document.createElement('style');
        style.id = 'code-copy-style';
        style.textContent = [
            '#article-container pre { position: relative; }',
            '.code-copy-btn {',
            '  position: absolute; top: 6px; right: 6px;',
            '  display: flex; align-items: center; justify-content: center;',
            '  width: 28px; height: 28px; padding: 0;',
            '  border: none; border-radius: 4px;',
            '  background: rgba(255,255,255,0.12);',
            '  color: #fff; cursor: pointer; opacity: 0;',
            '  transition: opacity .2s, background .2s;',
            '  font-size: 13px; line-height: 1; z-index: 2;',
            '}',
            '#article-container pre:hover .code-copy-btn,',
            '.code-copy-btn:focus-visible { opacity: 1; }',
            '.code-copy-btn:hover { background: rgba(255,255,255,0.25); }',
            '.code-copy-btn--done { opacity: 1; color: #4ade80; }',
            // 浅色背景的代码块（非 monokai）按钮也要可见
            '#article-container pre:not([style*="background-color"]) .code-copy-btn,',
            '#article-container pre[style*="#ffffff"] .code-copy-btn {',
            '  background: rgba(0,0,0,0.06); color: #555;',
            '}',
            '#article-container pre:not([style*="background-color"]) .code-copy-btn:hover,',
            '#article-container pre[style*="#ffffff"] .code-copy-btn:hover {',
            '  background: rgba(0,0,0,0.12);',
            '}'
        ].join('\n');
        document.head.appendChild(style);
    }

    function initCodeCopy() {
        var pres = document.querySelectorAll(SELECTOR);
        if (!pres || !pres.length) return;

        pres.forEach(function (pre) {
            // 跳过已处理、或嵌套在另一个 pre 内的（防御性，一般不会出现）
            if (pre.dataset.codecopy === '1') return;
            if (pre.parentElement && pre.parentElement.closest(SELECTOR)) return;
            pre.dataset.codecopy = '1';

            // 若 pre 本身是 static 定位，按钮绝对定位会锚到更上层；先确保定位上下文
            var computed = getComputedStyle(pre);
            if (computed.position === 'static') {
                pre.style.position = 'relative';
            }

            var btn = document.createElement('button');
            btn.type = 'button';
            btn.className = 'code-copy-btn';
            btn.title = '复制代码';
            btn.setAttribute('aria-label', '复制代码');
            btn.innerHTML = '<i class="fas fa-paste"></i>';

            btn.addEventListener('click', function () {
                // 优先取纯文本，剔除 chroma 行号锚点（<a> 内的数字）残留
                var text = extractCodeText(pre);
                copyText(text).then(function (ok) {
                    flash(btn, ok ? 'fas fa-check' : 'fas fa-times');
                });
            });

            pre.appendChild(btn);
        });
    }

    // 提取代码文本：克隆节点后移除 chroma 行号（带 user-select:none 的 span），
    // 再取 textContent。用 DOM 操作而非正则，避免行号格式（" 1"/"1"/"01"等）误判。
    function extractCodeText(pre) {
        var clone = pre.cloneNode(true);
        // 移除复制按钮自身（避免按钮文字混入）
        var btns = clone.querySelectorAll('.code-copy-btn');
        btns.forEach(function (b) { b.remove(); });
        // chroma 行号 span 带 user-select:none，且通常含 <a> 锚点；逐个检查 inline style
        var spans = clone.querySelectorAll('span');
        spans.forEach(function (s) {
            if (s.style.userSelect === 'none' || s.style.webkitUserSelect === 'none') {
                s.remove();
            }
        });
        var code = clone.querySelector('code');
        var text = code ? code.textContent : clone.textContent;
        // 去掉每行尾部多余换行差异（chroma 每行末尾有 \n，整体末尾可能多一个）
        return text.replace(/\n$/, '');
    }

    // 复制：优先 Clipboard API（HTTPS/localhost），降级 execCommand
    function copyText(text) {
        return new Promise(function (resolve) {
            if (navigator.clipboard && window.isSecureContext) {
                navigator.clipboard.writeText(text).then(
                    function () { resolve(true); },
                    function () { resolve(execCopy(text)); }
                );
            } else {
                resolve(execCopy(text));
            }
        });
    }

    function execCopy(text) {
        try {
            var ta = document.createElement('textarea');
            ta.value = text;
            ta.style.position = 'fixed';
            ta.style.opacity = '0';
            document.body.appendChild(ta);
            ta.select();
            var ok = document.execCommand('copy');
            document.body.removeChild(ta);
            return ok;
        } catch (e) {
            return false;
        }
    }

    // 复制后图标切换反馈：粘贴图标 → 对勾/叉 → 回弹
    function flash(btn, iconClass) {
        var icon = btn.querySelector('i');
        if (!icon) return;
        var origin = icon.className;
        icon.className = iconClass;
        btn.classList.add('code-copy-btn--done');
        setTimeout(function () {
            icon.className = origin;
            btn.classList.remove('code-copy-btn--done');
        }, 1500);
    }

    // 首次加载
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initCodeCopy);
    } else {
        initCodeCopy();
    }
    // PJAX 切换页面后重新扫描
    document.addEventListener('pjax:complete', initCodeCopy);
})();
