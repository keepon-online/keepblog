/**
 * 文章正文图片平滑灯箱 (Lightbox)
 *
 * 在文章详情页（#article-container 内），读者点击任意插图即可半透明平滑放大居中预览。
 * 支持 ESC 键、点击背景或点击放大图退出；自动提取 alt/title 作为底部图注；
 * 自动识别并跳过小徽标（badges）、外部普通超链接等无需放大的元素。
 *
 * 时机：通过 document 全局事件委托处理，天然兼容 PJAX 无刷新换页与图片懒加载。
 * 用 var/IIFE 编写，防止重复声明抛 SyntaxError。
 */
;(function () {
    var OVERLAY_ID = 'image-lightbox';
    var STYLE_ID = 'image-lightbox-style';

    // 动态注入自包含样式
    if (!document.getElementById(STYLE_ID)) {
        var style = document.createElement('style');
        style.id = STYLE_ID;
        style.textContent = [
            '#article-container img:not(.no-lightbox) {',
            '  cursor: zoom-in;',
            '  transition: filter .2s ease;',
            '}',
            '#' + OVERLAY_ID + ' {',
            '  position: fixed; top: 0; left: 0; right: 0; bottom: 0;',
            '  z-index: 99999;',
            '  display: flex; align-items: center; justify-content: center; flex-direction: column;',
            '  background: rgba(15, 23, 42, 0.88);',
            '  backdrop-filter: blur(8px); -webkit-backdrop-filter: blur(8px);',
            '  opacity: 0; visibility: hidden; pointer-events: none;',
            '  transition: opacity .25s ease, visibility .25s ease;',
            '  cursor: zoom-out; user-select: none; padding: 24px; box-sizing: border-box;',
            '}',
            '#' + OVERLAY_ID + '.is-open {',
            '  opacity: 1; visibility: visible; pointer-events: auto;',
            '}',
            '#' + OVERLAY_ID + ' img {',
            '  max-width: 92vw; max-height: 85vh; object-fit: contain;',
            '  border-radius: 8px;',
            '  box-shadow: 0 20px 45px rgba(0, 0, 0, 0.55);',
            '  transform: scale(0.94);',
            '  transition: transform .25s cubic-bezier(0.16, 1, 0.3, 1);',
            '  cursor: zoom-out;',
            '}',
            '#' + OVERLAY_ID + '.is-open img {',
            '  transform: scale(1);',
            '}',
            '#' + OVERLAY_ID + ' .lightbox-caption {',
            '  margin-top: 14px; color: #f8fafc; font-size: 14px;',
            '  padding: 4px 16px; border-radius: 20px;',
            '  background: rgba(0, 0, 0, 0.55); border: 1px solid rgba(255, 255, 255, 0.1);',
            '  max-width: 80vw; text-align: center; text-overflow: ellipsis; overflow: hidden; white-space: nowrap;',
            '}',
            '#' + OVERLAY_ID + ' .lightbox-close {',
            '  position: absolute; top: 20px; right: 24px;',
            '  width: 38px; height: 38px; border-radius: 50%;',
            '  background: rgba(255, 255, 255, 0.15); color: #ffffff;',
            '  border: 1px solid rgba(255, 255, 255, 0.2); font-size: 16px;',
            '  display: flex; align-items: center; justify-content: center;',
            '  cursor: pointer; transition: all .2s ease;',
            '}',
            '#' + OVERLAY_ID + ' .lightbox-close:hover {',
            '  background: rgba(255, 255, 255, 0.3); transform: scale(1.08);',
            '}'
        ].join('\n');
        document.head.appendChild(style);
    }

    // 获取或创建单例 DOM
    function getOverlay() {
        var overlay = document.getElementById(OVERLAY_ID);
        if (!overlay) {
            overlay = document.createElement('div');
            overlay.id = OVERLAY_ID;
            overlay.innerHTML = '<button class="lightbox-close" type="button" title="关闭 (ESC)"><i class="fas fa-times"></i></button>' +
                                '<img alt="lightbox preview">' +
                                '<div class="lightbox-caption" style="display:none;"></div>';
            document.body.appendChild(overlay);

            // 点击遮罩、图片或关闭按钮均关闭
            overlay.addEventListener('click', closeLightbox);
        }
        return overlay;
    }

    function openLightbox(img) {
        var overlay = getOverlay();
        var previewImg = overlay.querySelector('img');
        var caption = overlay.querySelector('.lightbox-caption');

        previewImg.src = img.currentSrc || img.src;
        previewImg.alt = img.alt || '';

        var captionText = (img.alt || img.title || '').trim();
        if (captionText) {
            caption.textContent = captionText;
            caption.style.display = 'block';
        } else {
            caption.style.display = 'none';
        }

        overlay.classList.add('is-open');
        document.body.style.overflow = 'hidden';
    }

    function closeLightbox() {
        var overlay = document.getElementById(OVERLAY_ID);
        if (!overlay || !overlay.classList.contains('is-open')) return;
        overlay.classList.remove('is-open');
        document.body.style.overflow = '';
    }

    // 全局事件委托：监听 #article-container 内的图片点击
    document.addEventListener('click', function (e) {
        var img = e.target;
        if (!img || img.tagName !== 'IMG') return;
        if (!img.closest('#article-container')) return;
        if (img.classList.contains('no-lightbox')) return;

        // 跳过宽度或高度过小的徽章/表情（如 < 40px）
        if (img.naturalWidth && img.naturalWidth < 40 && img.naturalHeight < 40) return;

        // 若图片在 <a> 标签内，如果是普通页面导航链接则不拦截；若指向图片本身则拦截并放大
        var link = img.closest('a');
        if (link) {
            var href = link.getAttribute('href') || '';
            var isImgLink = /\.(png|jpe?g|gif|webp|svg|bmp)(\?.*)?$/i.test(href) || href === img.getAttribute('src');
            if (!isImgLink) return;
            e.preventDefault();
        }

        openLightbox(img);
    });

    // ESC 键关闭
    window.addEventListener('keydown', function (e) {
        if (e.key === 'Escape' || e.keyCode === 27) {
            closeLightbox();
        }
    });

    // PJAX 换页时清理关闭
    document.addEventListener('pjax:send', closeLightbox);
})();
