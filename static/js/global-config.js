// 全站常量配置。原内联于 layout/head.html，外置为静态文件以命中长缓存，
// 同时减少每个 HTML 页面的体积。此处为纯数据常量，不含任何模板变量。
// 注意：依赖此配置的引导脚本（暗色模式检测等）仍内联在 head.html，
// 必须在本文件加载后再执行。本文件需在引导脚本之前引入。
const GLOBAL_CONFIG = {
    root: '/',
    algolia: undefined,
    localSearch: {
        path: "/search.xml",
        preload: !1,
        top_n_per_article: 1,
        unescape: !1,
        languages: {
            hits_empty: "找不到您查詢的內容：${query}",
            hits_stats: "共找到 ${hits} 篇文章"
        }
    },
    translate: undefined,
    noticeOutdate: undefined,
    highlight: {
        "plugin": "highlighjs",
        "highlightCopy": true,
        "highlightLang": true,
        "highlightHeightLimit": true
    },
    copy: {
        success: '复制成功',
        error: '复制错误',
        noSupport: '浏览器不支持'
    },
    relativeDate: {
        homepage: false,
        post: false
    },
    runtime: '天',
    date_suffix: {
        just: '刚刚',
        min: '分钟前',
        hour: '小时前',
        day: '天前',
        month: '个月前'
    },
    copyright: {
        "limitCount": 50,
        "languages": {
            "author": "作者: jieepre",
            "link": "链接: ",
            "source": "来源: 小助理",
            "info": "著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。"
        }
    },
    Snackbar: {
        "chs_to_cht": "你已切换为繁体",
        "cht_to_chs": "你已切换为简体",
        "day_to_night": "你已切换为深色模式",
        "night_to_day": "你已切换为浅色模式",
        "bgLight": "#49b1f5",
        "bgDark": "#1f1f1f",
        "position": "top-right"
    },
    isPhotoFigcaption: false,
    islazyload: false,
    isAnchor: false
}
