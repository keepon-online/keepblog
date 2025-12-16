/**
 * PJAX + 音乐播放器管理模块
 * 实现无刷新页面切换，同时保持音乐播放不中断
 */
(function () {
    'use strict';

    // ========================================
    // 音乐状态管理
    // ========================================
    const MusicState = {
        STORAGE_KEY: 'musicPlayerState',

        /**
         * 保存当前播放状态到 localStorage
         */
        save() {
            if (!window.ap || !window.ap.audio) return;

            try {
                const state = {
                    isPlaying: !window.ap.audio.paused,
                    currentTime: window.ap.audio.currentTime || 0,
                    currentIndex: window.ap.list.index || 0,
                    volume: window.ap.audio.volume || 0.5,
                    timestamp: Date.now()
                };
                localStorage.setItem(this.STORAGE_KEY, JSON.stringify(state));
            } catch (e) {
                console.warn('[MusicState] 保存状态失败:', e);
            }
        },

        /**
         * 从 localStorage 恢复播放状态
         */
        restore() {
            if (!window.ap) return;

            try {
                const saved = localStorage.getItem(this.STORAGE_KEY);
                if (!saved) return;

                const state = JSON.parse(saved);

                // 检查状态是否过期（超过24小时）
                if (Date.now() - state.timestamp > 24 * 60 * 60 * 1000) {
                    this.clear();
                    return;
                }

                // 恢复播放列表索引
                if (state.currentIndex !== undefined && window.ap.list.audios.length > state.currentIndex) {
                    window.ap.list.switch(state.currentIndex);
                }

                // 恢复音量
                if (state.volume !== undefined) {
                    window.ap.volume(state.volume, true);
                }

                // 恢复播放进度（延迟执行以确保音频加载）
                if (state.currentTime > 0) {
                    setTimeout(() => {
                        window.ap.seek(state.currentTime);
                    }, 500);
                }

                // 恢复播放状态
                if (state.isPlaying) {
                    setTimeout(() => {
                        window.ap.play();
                    }, 600);
                }
            } catch (e) {
                console.warn('[MusicState] 恢复状态失败:', e);
            }
        },

        /**
         * 清除保存的状态
         */
        clear() {
            localStorage.removeItem(this.STORAGE_KEY);
        }
    };

    // ========================================
    // NProgress 进度条管理
    // ========================================
    const ProgressBar = {
        /**
         * 开始显示进度条
         */
        start() {
            if (window.NProgress) {
                NProgress.start();
            }
        },

        /**
         * 完成进度条
         */
        done() {
            if (window.NProgress) {
                NProgress.done();
            }
        },

        /**
         * 初始化 NProgress 配置
         */
        init() {
            if (window.NProgress) {
                NProgress.configure({
                    showSpinner: false,
                    speed: 400,
                    minimum: 0.2,
                    trickleSpeed: 200
                });
            }
        }
    };

    // ========================================
    // PJAX 管理器
    // ========================================
    const PjaxManager = {
        instance: null,

        /**
         * 初始化 PJAX
         */
        init() {
            if (typeof Pjax === 'undefined') {
                console.warn('[PjaxManager] Pjax 库未加载');
                return;
            }

            this.instance = new Pjax({
                selectors: [
                    "head > title",
                    "#config-diff",
                    "#body-wrap",
                    "#rightside-config-hide",
                    "#rightside-config-show",
                    ".js-pjax"
                ],
                cacheBust: false,
                timeout: 8000,
                scrollTo: false,
                scrollRestoration: false
            });

            this.bindEvents();
            console.log('[PjaxManager] PJAX 初始化完成');
        },

        /**
         * 绑定 PJAX 事件
         */
        bindEvents() {
            // 开始请求
            document.addEventListener('pjax:send', this.onSend.bind(this));

            // 请求完成
            document.addEventListener('pjax:complete', this.onComplete.bind(this));

            // 请求成功
            document.addEventListener('pjax:success', this.onSuccess.bind(this));

            // 请求失败
            document.addEventListener('pjax:error', this.onError.bind(this));
        },

        /**
         * PJAX 开始请求时
         */
        onSend() {
            // 保存音乐状态
            MusicState.save();

            // 显示进度条
            ProgressBar.start();

            // 添加加载中类名
            document.body.classList.add('pjax-loading');
        },

        /**
         * PJAX 请求完成时
         */
        onComplete() {
            // 隐藏进度条
            ProgressBar.done();

            // 移除加载中类名
            document.body.classList.remove('pjax-loading');

            // 重新初始化页面功能
            if (typeof window.refreshFn === 'function') {
                window.refreshFn();
            }

            // 执行 PJAX 回调（如果有）
            this.executePjaxCallbacks();

            // 平滑滚动到顶部
            if (typeof btf !== 'undefined' && typeof btf.scrollToDest === 'function') {
                btf.scrollToDest(0, 300);
            } else {
                window.scrollTo({ top: 0, behavior: 'smooth' });
            }

            // 重新加载不雅刷新的脚本
            this.reloadScripts();
        },

        /**
         * PJAX 请求成功时
         */
        onSuccess() {
            // 确保音乐继续播放
            if (window.ap && window.ap.audio.paused) {
                const savedState = localStorage.getItem(MusicState.STORAGE_KEY);
                if (savedState) {
                    const state = JSON.parse(savedState);
                    if (state.isPlaying) {
                        window.ap.play();
                    }
                }
            }
        },

        /**
         * PJAX 请求失败时
         */
        onError(e) {
            console.warn('[PjaxManager] PJAX 请求失败，降级为传统导航');
            ProgressBar.done();
            document.body.classList.remove('pjax-loading');

            // 降级为传统导航
            if (e && e.request && e.request.responseURL) {
                window.location.href = e.request.responseURL;
            }
        },

        /**
         * 执行 PJAX 回调函数
         */
        executePjaxCallbacks() {
            // 执行注册的回调
            if (window.pjaxCallbacks && Array.isArray(window.pjaxCallbacks)) {
                window.pjaxCallbacks.forEach(callback => {
                    if (typeof callback === 'function') {
                        try {
                            callback();
                        } catch (e) {
                            console.warn('[PjaxManager] 执行回调失败:', e);
                        }
                    }
                });
            }
        },

        /**
         * 重新加载需要重新执行的脚本
         */
        reloadScripts() {
            // 重新初始化代码高亮
            if (typeof Prism !== 'undefined') {
                Prism.highlightAll();
            }

            // 重新初始化 busuanzi 统计
            if (typeof bszCaller !== 'undefined') {
                bszCaller.fetch();
            }
        }
    };

    // ========================================
    // 音乐播放器管理器
    // ========================================
    const MusicPlayer = {
        /**
         * 初始化音乐播放器
         */
        init(options = {}) {
            // 如果播放器已存在，不重复初始化
            if (window.ap) {
                console.log('[MusicPlayer] 播放器已存在，跳过初始化');
                return;
            }

            const container = document.getElementById('player');
            if (!container) {
                console.warn('[MusicPlayer] 找不到播放器容器 #player');
                return;
            }

            // 默认音乐列表
            const defaultAudio = [{
                name: '平凡之路',
                artist: '朴树',
                url: 'http://music.163.com/song/media/outer/url?id=1497917991.mp3',
                cover: 'http://p2.music.126.net/IwEI0tFPh4w9OjY6RM2IJQ==/109951163009071893.jpg?param=90y90',
                type: 'auto'
            }];

            // 合并配置
            const config = Object.assign({
                container: container,
                fixed: true,
                autoplay: false,
                theme: '#b7daff',
                loop: 'all',
                order: 'random',
                preload: 'auto',
                lrcType: 0,
                volume: 0.5,
                mutex: true,
                listFolded: true,
                storageName: 'aplayer-setting',
                audio: defaultAudio
            }, options);

            try {
                window.ap = new APlayer(config);

                // 绑定播放器事件
                this.bindEvents();

                // 恢复之前的播放状态
                MusicState.restore();

                console.log('[MusicPlayer] 播放器初始化完成');
            } catch (e) {
                console.error('[MusicPlayer] 初始化失败:', e);
            }
        },

        /**
         * 绑定播放器事件
         */
        bindEvents() {
            if (!window.ap) return;

            // 播放/暂停时保存状态
            window.ap.on('play', () => {
                MusicState.save();
            });

            window.ap.on('pause', () => {
                MusicState.save();
            });

            // 切换歌曲时保存状态
            window.ap.on('listswitch', () => {
                MusicState.save();
            });

            // 定期保存播放进度（每10秒）
            setInterval(() => {
                if (window.ap && !window.ap.audio.paused) {
                    MusicState.save();
                }
            }, 10000);

            // 页面关闭前保存状态
            window.addEventListener('beforeunload', () => {
                MusicState.save();
            });
        }
    };

    // ========================================
    // 注册 PJAX 回调
    // ========================================
    window.pjaxCallbacks = window.pjaxCallbacks || [];

    /**
     * 注册 PJAX 完成后的回调函数
     * @param {Function} callback 回调函数
     */
    window.registerPjaxCallback = function (callback) {
        if (typeof callback === 'function') {
            window.pjaxCallbacks.push(callback);
        }
    };

    // ========================================
    // 初始化入口
    // ========================================
    window.initPjaxMusic = function (musicOptions) {
        // 初始化进度条
        ProgressBar.init();

        // 初始化音乐播放器（仅首次）
        MusicPlayer.init(musicOptions);

        // 初始化 PJAX
        PjaxManager.init();
    };

    // 导出模块（可选）
    window.PjaxMusic = {
        MusicState,
        ProgressBar,
        PjaxManager,
        MusicPlayer
    };

})();
