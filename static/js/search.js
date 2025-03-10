window.addEventListener('load', () => {
    let loadFlag = false
    let dataObj = []
    const $searchMask = document.getElementById('search-mask')

    const openSearch = () => {
        const bodyStyle = document.body.style
        bodyStyle.width = '100%'
        bodyStyle.overflow = 'hidden'
        btf.animateIn($searchMask, 'to_show 0.5s')
        btf.animateIn(document.querySelector('#local-search .search-dialog'), 'titleScale 0.5s')
        initsearch()
        setTimeout(() => {
            document.querySelector('#local-search-input input').focus()
        }, 100)
        if (!loadFlag) {
            search()
            loadFlag = true
        }
        // shortcut: ESC
        document.addEventListener('keydown', function f(event) {
            if (event.code === 'Escape') {
                closeSearch()
                document.removeEventListener('keydown', f)
            }
        })
    }
    const initsearch = () => {
        document.querySelector('#local-search-input input').value = ''
        document.getElementById('local-search-results').innerHTML = ''
    }

    const closeSearch = () => {
        const bodyStyle = document.body.style
        bodyStyle.width = ''
        bodyStyle.overflow = ''
        btf.animateOut(document.querySelector('#local-search .search-dialog'), 'search_close .5s')
        btf.animateOut($searchMask, 'to_hide 0.5s')
    }

    const searchClickFn = () => {
        document.querySelector('#search-button > .search').addEventListener('click', openSearch)
    }

    const searchClickFnOnce = () => {
        document.querySelector('#local-search .search-close-button').addEventListener('click', closeSearch)
        $searchMask.addEventListener('click', closeSearch)
    }

    const search = () => {
        const $input = document.querySelector('#local-search-input input')
        const $resultContent = document.getElementById('local-search-results')
        const $loadingStatus = document.getElementById('loading-status')
        $input.addEventListener('keydown', function (event) {
            const keywords = this.value.trim().toLowerCase().split(/[\s]+/)
            if (keywords[0] !== '') $loadingStatus.innerHTML = '<i class="fas fa-spinner fa-pulse"></i>'
            $resultContent.innerHTML = ''
            let str = '<div class="search-result-list">'
            if (keywords.length <= 0) return
            if (keywords.length >= 30) return;
            let count = 0
            if (event.key === "Enter") {
                // 在这里执行您的操作

                fetch(window.location.origin+'/search/' + keywords)
                    .then(function (response) {
                        console.log(response);
                        return response.json();
                    })
                    .then(function (res) {
                        if (res.code === 200) {
                            dataObj = res.payload
                        }
                    });

                // perform local searching
                dataObj.forEach(data => {
                    let dataTitle = data.title ? data.title.trim().toLowerCase() : ''
                    const dataUrl = "/post/" + data.postSlug
                    // show search results
                    // highlight all keywords
                    keywords.forEach(keyword => {
                        const regS = new RegExp(keyword, 'gi')
                        dataTitle = dataTitle.replace(regS, '<span class="search-keyword">' + keyword + '</span>')
                    })
                    str += '<div class="local-search__hit-item"><a href="' + dataUrl + '" class="search-result-title">' + dataTitle + '</a>'
                    count += 1
                    str += '</div>'
                })
                if (count === 0) {
                    str += '<div id="local-search__hits-empty">' + GLOBAL_CONFIG.localSearch.languages.hits_empty.replace(/\$\{query}/, this.value.trim()) + '</div>'
                }
                str += '</div>'
                $resultContent.innerHTML = str
            }
            if (keywords[0] !== '') $loadingStatus.innerHTML = ''
        })
    }

    searchClickFn()
    searchClickFnOnce()

    // pjax
    window.addEventListener('pjax:complete', () => {
        !btf.isHidden($searchMask) && closeSearch()
        searchClickFn()
    })
})
