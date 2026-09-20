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

    function initCodeCopy() {
        var pres = document.querySelectorAll(SELECTOR);
        if (!pres || !pres.length) return;

        pres.forEach(function (pre) {
            // 跳过已处理、或嵌套在另一个 pre 内的、或位于 figure.highlight 内的
            if (pre.dataset.codecopy === '1') return;
            if (pre.parentElement && pre.parentElement.closest(SELECTOR)) return;
            if (pre.closest('figure.highlight')) {
                pre.dataset.codecopy = '1';
                return;
            }
            pre.dataset.codecopy = '1';

            // 若未包裹，创建 code-block-wrap 容器
            var parent = pre.parentElement;
            var wrap;
            if (parent && parent.classList.contains('code-block-wrap')) {
                wrap = parent;
            } else {
                wrap = document.createElement('div');
                wrap.className = 'code-block-wrap';
                parent.insertBefore(wrap, pre);
                wrap.appendChild(pre);
            }

            // 识别代码块语言
            var codeEl = pre.querySelector('code');
            var lang = '';
            var dataLang = (codeEl && codeEl.getAttribute('data-lang')) || pre.getAttribute('data-lang');
            if (dataLang) {
                lang = dataLang.trim();
            } else {
                var langClass = (codeEl && codeEl.className) || pre.className || '';
                var match = langClass.match(/language-([a-zA-Z0-9_+-]+)/);
                if (match && match[1]) {
                    lang = match[1];
                }
            }

            // 创建或获取 Mac 顶栏
            var tools = wrap.querySelector('.code-block-tools');
            if (!tools) {
                tools = document.createElement('div');
                tools.className = 'code-block-tools';

                // Mac 三色圆点
                var dots = document.createElement('div');
                dots.className = 'code-mac-dots';
                dots.innerHTML = '<span></span><span></span><span></span>';
                tools.appendChild(dots);

                // 语言标签（若有）
                if (lang) {
                    var langSpan = document.createElement('span');
                    langSpan.className = 'code-lang-label';
                    langSpan.textContent = lang.toUpperCase();
                    tools.appendChild(langSpan);
                }

                wrap.insertBefore(tools, pre);
            }

            // 复制按钮
            var btn = document.createElement('button');
            btn.type = 'button';
            btn.className = 'code-copy-btn';
            btn.title = '复制代码';
            btn.setAttribute('aria-label', '复制代码');
            btn.innerHTML = '<i class="fas fa-paste"></i>';

            btn.addEventListener('click', function () {
                var text = extractCodeText(pre);
                copyText(text).then(function (ok) {
                    flash(btn, ok ? 'fas fa-check' : 'fas fa-times');
                });
            });

            tools.appendChild(btn);
        });
    }

    // 提取代码文本：克隆节点后移除 chroma 行号（带 user-select:none 的 span），
    // 再取 textContent。用 DOM 操作而非正则，避免行号格式（" 1"/"1"/"01"等）误判。
    function extractCodeText(pre) {
        var clone = pre.cloneNode(true);
        // 移除内部按钮（防御性）
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
