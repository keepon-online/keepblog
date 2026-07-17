/**
 * 搜索功能优化版
 * 特性：
 * - 实时搜索（防抖 300ms）
 * - 异步请求正确处理
 * - 动态搜索结果统计
 * - 改进的加载状态
 * - 键盘导航支持
 * - 搜索历史记录
 */
window.addEventListener('load', () => {
    let loadFlag = false;
    let dataObj = [];
    let debounceTimer = null;
    let currentIndex = -1; // 用于键盘导航
    const DEBOUNCE_DELAY = 300;
    const HISTORY_KEY = 'search_history';
    const MAX_HISTORY = 5;

    const $searchMask = document.getElementById('search-mask');
    const $searchDialog = document.querySelector('#local-search .search-dialog');
    const $input = document.querySelector('#local-search-input input');
    const $resultContent = document.getElementById('local-search-results');
    const $loadingStatus = document.getElementById('loading-status');
    const $statsWrap = document.getElementById('local-search-stats-wrap');

    // 获取搜索历史
    const getSearchHistory = () => {
        try {
            return JSON.parse(localStorage.getItem(HISTORY_KEY)) || [];
        } catch {
            return [];
        }
    };

    // 保存搜索历史
    const saveSearchHistory = (keyword) => {
        if (!keyword.trim()) return;
        let history = getSearchHistory();
        history = history.filter(h => h !== keyword);
        history.unshift(keyword);
        history = history.slice(0, MAX_HISTORY);
        localStorage.setItem(HISTORY_KEY, JSON.stringify(history));
    };

    // 显示搜索历史
    const showSearchHistory = () => {
        const history = getSearchHistory();
        if (history.length === 0) {
            $resultContent.innerHTML = '<div class="search-hint">输入关键词开始搜索</div>';
            $statsWrap.style.display = 'none';
            return;
        }

        let html = '<div class="search-history">';
        html += '<div class="search-history-title"><i class="fas fa-history"></i> 搜索历史</div>';
        history.forEach(keyword => {
            html += `<div class="search-history-item" data-keyword="${keyword}">
                <span>${keyword}</span>
                <i class="fas fa-times search-history-delete" data-keyword="${keyword}"></i>
            </div>`;
        });
        html += '</div>';
        $resultContent.innerHTML = html;
        $statsWrap.style.display = 'none';

        // 绑定历史点击事件
        document.querySelectorAll('.search-history-item span').forEach(item => {
            item.addEventListener('click', () => {
                $input.value = item.textContent;
                performSearch(item.textContent);
            });
        });

        // 绑定删除历史事件
        document.querySelectorAll('.search-history-delete').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                const keyword = btn.dataset.keyword;
                let history = getSearchHistory().filter(h => h !== keyword);
                localStorage.setItem(HISTORY_KEY, JSON.stringify(history));
                showSearchHistory();
            });
        });
    };

    // 显示加载状态
    const showLoading = () => {
        $loadingStatus.innerHTML = '<i class="fas fa-spinner fa-pulse"></i>';
        $resultContent.innerHTML = '<div class="search-loading"><i class="fas fa-spinner fa-pulse"></i> 搜索中...</div>';
    };

    // 隐藏加载状态
    const hideLoading = () => {
        $loadingStatus.innerHTML = '';
    };

    // 执行搜索
    const performSearch = async (keyword) => {
        const keywords = keyword.trim().toLowerCase().split(/[\s]+/);

        if (!keywords[0]) {
            showSearchHistory();
            return;
        }

        if (keywords.length > 30) {
            $resultContent.innerHTML = '<div class="search-error">关键词过多，请精简搜索条件</div>';
            return;
        }

        showLoading();

        try {
            const response = await fetch(window.location.origin + '/search/' + encodeURIComponent(keyword));
            const res = await response.json();

            if (res.code === 200) {
                dataObj = res.payload || [];
                renderResults(keywords);
                saveSearchHistory(keyword);
            } else {
                $resultContent.innerHTML = '<div class="search-error">搜索失败，请稍后重试</div>';
            }
        } catch (error) {
            console.error('Search error:', error);
            $resultContent.innerHTML = '<div class="search-error">网络错误，请检查连接</div>';
        } finally {
            hideLoading();
        }
    };

    // 渲染搜索结果
    const renderResults = (keywords) => {
        currentIndex = -1;

        if (dataObj.length === 0) {
            const emptyMsg = GLOBAL_CONFIG.localSearch?.languages?.hits_empty?.replace(/\$\{query}/, $input.value.trim())
                || `没有找到与 "${$input.value.trim()}" 相关的文章`;
            $resultContent.innerHTML = `<div id="local-search__hits-empty">${emptyMsg}</div>`;
            $statsWrap.style.display = 'none';
            return;
        }

        let html = '<div class="search-result-list">';
        dataObj.forEach((data, index) => {
            let dataTitle = data.title ? data.title.trim() : '无标题';
            let displayTitle = dataTitle;
            const dataUrl = "/post/" + data.postSlug;

            // 高亮关键词
            keywords.forEach(keyword => {
                const regS = new RegExp(`(${keyword})`, 'gi');
                displayTitle = displayTitle.replace(regS, '<span class="search-keyword">$1</span>');
            });

            // 摘要处理
            let summary = data.summary || '';
            if (summary.length > 80) {
                summary = summary.substring(0, 80) + '...';
            }
            keywords.forEach(keyword => {
                const regS = new RegExp(`(${keyword})`, 'gi');
                summary = summary.replace(regS, '<span class="search-keyword">$1</span>');
            });

            html += `<div class="local-search__hit-item" data-index="${index}">
                <a href="${dataUrl}" class="search-result-title">${displayTitle}</a>
                ${summary ? `<p class="search-result-summary">${summary}</p>` : ''}
            </div>`;
        });
        html += '</div>';

        $resultContent.innerHTML = html;

        // 让动态生成的结果链接被 PJAX 接管，避免整页刷新中断音乐
        if (window.PjaxMusic?.PjaxManager?.refreshLinks) {
            window.PjaxMusic.PjaxManager.refreshLinks($resultContent);
        }

        // 更新统计
        $statsWrap.style.display = 'block';
        $statsWrap.querySelector('.search-result-stats').textContent = `共找到 ${dataObj.length} 篇文章`;
    };

    // 键盘导航
    const handleKeyNavigation = (e) => {
        const items = document.querySelectorAll('.local-search__hit-item');
        if (items.length === 0) return;

        if (e.key === 'ArrowDown') {
            e.preventDefault();
            currentIndex = Math.min(currentIndex + 1, items.length - 1);
            updateActiveItem(items);
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            currentIndex = Math.max(currentIndex - 1, 0);
            updateActiveItem(items);
        } else if (e.key === 'Enter' && currentIndex >= 0) {
            e.preventDefault();
            const link = items[currentIndex].querySelector('a');
            if (link) {
                // 优先走 PJAX，与鼠标点击行为一致，避免整页刷新中断音乐
                const pjax = window.PjaxMusic?.PjaxManager?.instance;
                if (pjax && typeof pjax.loadUrl === 'function') {
                    pjax.loadUrl(link.href);
                } else {
                    window.location.href = link.href;
                }
            }
        }
    };

    const updateActiveItem = (items) => {
        items.forEach((item, index) => {
            if (index === currentIndex) {
                item.classList.add('active');
                item.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
            } else {
                item.classList.remove('active');
            }
        });
    };

    // 防抖搜索
    const debouncedSearch = (keyword) => {
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
            performSearch(keyword);
        }, DEBOUNCE_DELAY);
    };

    // 打开搜索
    const openSearch = () => {
        const bodyStyle = document.body.style;
        bodyStyle.width = '100%';
        bodyStyle.overflow = 'hidden';
        btf.animateIn($searchMask, 'to_show 0.5s');
        btf.animateIn($searchDialog, 'titleScale 0.5s');

        // 清空并显示历史
        $input.value = '';
        showSearchHistory();

        setTimeout(() => {
            $input.focus();
        }, 100);

        // 快捷键: ESC 关闭
        document.addEventListener('keydown', escHandler);
    };

    const escHandler = (event) => {
        if (event.code === 'Escape') {
            closeSearch();
            document.removeEventListener('keydown', escHandler);
        }
    };

    // 关闭搜索
    const closeSearch = () => {
        const bodyStyle = document.body.style;
        bodyStyle.width = '';
        bodyStyle.overflow = '';
        btf.animateOut($searchDialog, 'search_close .5s');
        btf.animateOut($searchMask, 'to_hide 0.5s');
        document.removeEventListener('keydown', escHandler);
    };

    // 初始化事件绑定
    const initSearchEvents = () => {
        // 搜索按钮点击
        const searchBtn = document.querySelector('#search-button > .search');
        if (searchBtn) {
            searchBtn.addEventListener('click', openSearch);
        }

        // 关闭按钮
        document.querySelector('#local-search .search-close-button')?.addEventListener('click', closeSearch);
        $searchMask?.addEventListener('click', closeSearch);

        // 输入事件 - 实时搜索
        $input?.addEventListener('input', (e) => {
            const value = e.target.value.trim();
            if (value) {
                debouncedSearch(value);
            } else {
                showSearchHistory();
            }
        });

        // 键盘导航
        $input?.addEventListener('keydown', (e) => {
            if (['ArrowDown', 'ArrowUp'].includes(e.key)) {
                handleKeyNavigation(e);
            } else if (e.key === 'Enter') {
                const value = $input.value.trim();
                if (value && currentIndex < 0) {
                    clearTimeout(debounceTimer);
                    performSearch(value);
                } else if (currentIndex >= 0) {
                    handleKeyNavigation(e);
                }
            }
        });
    };

    // 初始化
    initSearchEvents();

    // PJAX 支持
    window.addEventListener('pjax:complete', () => {
        if ($searchMask && !btf.isHidden($searchMask)) {
            closeSearch();
        }
        // 重新绑定搜索按钮事件
        const searchBtn = document.querySelector('#search-button > .search');
        if (searchBtn) {
            searchBtn.removeEventListener('click', openSearch);
            searchBtn.addEventListener('click', openSearch);
        }
    });
});
