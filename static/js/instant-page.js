/**
 * 链接悬停智能预加载 (Instant Page Prefetching)
 *
 * 原理：当用户鼠标悬停在站内链接超过 65ms（或移动端 touchstart 瞬间），
 * 在用户实际点击之前提前在后台低优先级预取目标页面的 HTML 并放入内存缓存。
 * 当用户随后点击触发 PJAX 跳转时，直接从内存读取内容完成 0ms 瞬时换页。
 *
 * 特性：
 * 1. 深度整合 PJAX：劫持 Pjax.prototype.doRequest，命中缓存直接同步回调；
 * 2. 流量节省感知：若用户开启 Save-Data 模式或处于 2G 慢速网络，自动停用；
 * 3. 严格边界防御：仅预取站内 GET 页面，排除 /admin、/api、锚点、下载附件及带 target="_blank" 的外链；
 * 4. 内存控制：LRU 策略限制最多缓存 15 个最近页面，避免占用内存；
 * 5. PJAX 幂等性：全部采用 var/IIFE 声明，多次加载不报错。
 */
;(function () {
    // 数据节省模式或慢网络不预取
    if (navigator.connection && (navigator.connection.saveData || /2g/.test(navigator.connection.effectiveType))) {
        return;
    }

    var MAX_CACHE_SIZE = 15;
    var HOVER_DELAY = 65; // ms
    var pageCache = new Map(); // url -> html text
    var inflightRequests = new Map(); // url -> Promise<string>
    var hoverTimer = null;
    var currentHoverUrl = null;

    // 获取标准化后的站内绝对路径（pathname + search）
    function getCleanUrl(url) {
        try {
            var parsed = new URL(url, window.location.origin);
            return parsed.pathname + parsed.search;
        } catch (e) {
            return null;
        }
    }

    // 检查链接是否值得预加载
    function shouldPrefetch(anchor) {
        if (!anchor || anchor.tagName !== 'A') return false;
        var href = anchor.getAttribute('href');
        if (!href || href.startsWith('#') || href.startsWith('javascript:') || href.startsWith('mailto:') || href.startsWith('tel:')) {
            return false;
        }
        if (anchor.target && anchor.target !== '_self') return false;
        if (anchor.hasAttribute('download') || anchor.getAttribute('data-no-instant') || anchor.getAttribute('data-no-prefetch')) {
            return false;
        }

        // 协议与域名必须相同
        if (anchor.origin !== window.location.origin) return false;

        var clean = getCleanUrl(anchor.href);
        if (!clean) return false;

        // 不预取当前所在页面
        var currentClean = getCleanUrl(window.location.href);
        if (clean === currentClean) return false;

        // 排除管理后台、API 接口和常见非 HTML 文件
        if (/^\/(admin|api)(\/|$)/i.test(clean)) return false;
        if (/\.(zip|tar|gz|rar|7z|pdf|png|jpe?g|gif|svg|webp|mp3|mp4|exe|dmg)(\?.*)?$/i.test(clean)) {
            return false;
        }

        return true;
    }

    // 执行预取
    function prefetchUrl(url) {
        var clean = getCleanUrl(url);
        if (!clean || pageCache.has(clean) || inflightRequests.has(clean)) {
            return;
        }

        // 构造 fetch 请求（带 PJAX 标识头，使服务端支持按需响应）
        var req = fetch(clean, {
            headers: {
                'X-PJAX': 'true',
                'X-Requested-With': 'XMLHttpRequest'
            },
            credentials: 'same-origin',
            priority: 'low'
        })
        .then(function (res) {
            if (!res.ok) throw new Error('HTTP ' + res.status);
            var ct = res.headers.get('content-type') || '';
            if (!ct.includes('text/html')) throw new Error('Not HTML');
            return res.text();
        })
        .then(function (html) {
            // 写入缓存并维护 LRU 容量
            if (pageCache.size >= MAX_CACHE_SIZE) {
                var firstKey = pageCache.keys().next().value;
                pageCache.delete(firstKey);
            }
            pageCache.set(clean, html);
            inflightRequests.delete(clean);
            return html;
        })
        .catch(function () {
            inflightRequests.delete(clean);
        });

        inflightRequests.set(clean, req);
    }

    // 劫持 Pjax 的请求通道
    function hookPjax() {
        if (typeof Pjax === 'undefined' || !Pjax.prototype || Pjax.prototype._instantHooked) {
            return;
        }
        Pjax.prototype._instantHooked = true;
        var origDoRequest = Pjax.prototype.doRequest;

        Pjax.prototype.doRequest = function (location, options, callback) {
            var clean = getCleanUrl(location);
            if (clean && pageCache.has(clean)) {
                var cachedHtml = pageCache.get(clean);
                // 模拟标准 XHR 响应，直接在微任务中完成 PJAX DOM 切换
                var mockXhr = {
                    readyState: 4,
                    status: 200,
                    responseText: cachedHtml,
                    responseURL: location,
                    getResponseHeader: function (h) {
                        if (h && h.toLowerCase() === 'x-pjax-url') return location;
                        return null;
                    }
                };
                setTimeout(function () {
                    callback(cachedHtml, mockXhr, location, options);
                }, 0);
                return { readyState: 4, abort: function () {} };
            }

            // 若正处于预取 inflight，等待其完成以复用数据
            if (clean && inflightRequests.has(clean)) {
                var inflight = inflightRequests.get(clean);
                var aborted = false;
                var pendingReq = {
                    readyState: 1,
                    abort: function () { aborted = true; }
                };
                inflight.then(function (html) {
                    if (aborted) return;
                    if (html) {
                        var mockXhr = {
                            readyState: 4,
                            status: 200,
                            responseText: html,
                            responseURL: location,
                            getResponseHeader: function () { return null; }
                        };
                        callback(html, mockXhr, location, options);
                    } else {
                        origDoRequest.call(this, location, options, callback);
                    }
                }.bind(this)).catch(function () {
                    if (!aborted) origDoRequest.call(this, location, options, callback);
                }.bind(this));
                return pendingReq;
            }

            return origDoRequest.call(this, location, options, callback);
        };
    }

    // 事件监听：指针悬停 65ms 后发起预取
    document.addEventListener('mouseover', function (e) {
        var anchor = e.target.closest('a');
        if (!shouldPrefetch(anchor)) return;

        var url = anchor.href;
        currentHoverUrl = url;
        clearTimeout(hoverTimer);
        hoverTimer = setTimeout(function () {
            if (currentHoverUrl === url) {
                prefetchUrl(url);
            }
        }, HOVER_DELAY);
    }, { passive: true });

    document.addEventListener('mouseout', function (e) {
        var anchor = e.target.closest('a');
        if (anchor && anchor.href === currentHoverUrl) {
            clearTimeout(hoverTimer);
            currentHoverUrl = null;
        }
    }, { passive: true });

    // 移动端：touchstart 瞬间直接发起预加载（触控到触发 click 约有 100~300ms 间隔）
    document.addEventListener('touchstart', function (e) {
        var anchor = e.target.closest('a');
        if (shouldPrefetch(anchor)) {
            prefetchUrl(anchor.href);
        }
    }, { passive: true });

    // 监听 PJAX 初始化或就绪
    if (typeof Pjax !== 'undefined') {
        hookPjax();
    } else {
        document.addEventListener('DOMContentLoaded', hookPjax);
    }
    document.addEventListener('pjax:complete', hookPjax);
})();
