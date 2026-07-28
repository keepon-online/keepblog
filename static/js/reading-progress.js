/**
 * 文章阅读进度条
 *
 * 在文章详情页（存在 #article-container）顶部显示一条固定进度条，随滚动反映
 * 阅读进度（按整页滚动百分比，非按请求）。
 *
 * 时机：DOMContentLoaded 首次执行；PJAX 切换页面后通过 pjax:complete 重新
 * 初始化。仅在文章详情页创建进度条元素，其他页面隐藏/移除。
 *
 * 注：用 var 与 function 声明，避免 PJAX 重新执行外部脚本时 let/const
 * 重复声明抛 SyntaxError（与 code-copy.js、global-config.js 同策略）。
 */
;(function () {
    var BAR_ID = 'reading-progress-bar';
    var bar = null;        // 当前进度条 DOM 引用
    var scrollHandler = null; // 当前绑定的 scroll 监听器，便于 PJAX 切换前解绑

    // 注入样式（一次性，自包含，不污染 index.css）
    if (!document.getElementById('reading-progress-style')) {
        var style = document.createElement('style');
        style.id = 'reading-progress-style';
        style.textContent = [
            '#' + BAR_ID + ' {',
            '  position: fixed; top: 0; left: 0;',
            '  width: 0; height: 3px; z-index: 9999;',
            '  background: var(--lm-a1, #49b1f5);',
            '  transition: width .1s ease-out; opacity: 0;',
            '  pointer-events: none;',
            '}',
            '#' + BAR_ID + '.reading-progress--show { opacity: 1; }'
        ].join('\n');
        document.head.appendChild(style);
    }

    // 计算并更新进度条宽度
    function updateProgress() {
        if (!bar) return;
        var scrollTop = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop || 0;
        var scrollHeight = document.documentElement.scrollHeight || document.body.scrollHeight || 0;
        var clientHeight = window.innerHeight || document.documentElement.clientHeight || 0;
        var max = scrollHeight - clientHeight;
        var ratio = max > 0 ? scrollTop / max : 0;
        if (ratio < 0) ratio = 0;
        if (ratio > 1) ratio = 1;
        bar.style.width = (ratio * 100) + '%';
    }

    // 清理上一次绑定的监听器与 DOM（PJAX 切换页面前调用）
    function teardown() {
        if (scrollHandler) {
            window.removeEventListener('scroll', scrollHandler, { passive: true });
            window.removeEventListener('resize', scrollHandler);
            scrollHandler = null;
        }
        if (bar) {
            bar.remove();
            bar = null;
        }
    }

    // 初始化：仅当文章正文容器存在时创建进度条
    function initReadingProgress() {
        teardown();

        var article = document.getElementById('article-container');
        // 非文章页（首页/归档/分类等）不显示进度条
        if (!article) return;

        bar = document.createElement('div');
        bar.id = BAR_ID;
        document.body.appendChild(bar);
        // 下一帧再显示，触发透明度过渡
        requestAnimationFrame(function () {
            if (bar) bar.classList.add('reading-progress--show');
        });

        scrollHandler = function () {
            // 用 rAF 节流，避免 scroll 高频回调
            if (scrollHandler._ticking) return;
            scrollHandler._ticking = true;
            requestAnimationFrame(function () {
                updateProgress();
                scrollHandler._ticking = false;
            });
        };

        window.addEventListener('scroll', scrollHandler, { passive: true });
        window.addEventListener('resize', scrollHandler);
        updateProgress();
    }

    // 首次加载
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initReadingProgress);
    } else {
        initReadingProgress();
    }
    // PJAX 切换页面后重新初始化
    document.addEventListener('pjax:complete', initReadingProgress);
})();
