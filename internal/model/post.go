package model

const TPostsTable = "post"

type Post struct {
	PostId          uint64 `json:"postId" gorm:"primaryKey;autoIncrement"`
	Title           string `json:"title" gorm:"default:''"`
	PostSlug        string `json:"postSlug" gorm:"default:''"`
	Author          string `json:"author" gorm:"default:''"`
	CoverImage      string `json:"coverImage" gorm:"default:''"`
	PostContent     string `json:"postContent" gorm:"default:''"`
	PostContentHtml string `json:"postContentHtml" gorm:"default:''"`
	Summary         string `json:"summary" gorm:"default:''"`
	Type            uint8  `json:"type" gorm:"default:1"`
	Top             uint8  `json:"top" gorm:"default:0"`
	ReadCount       uint32 `json:"readCount" gorm:"default:0"`
	WordCount       uint32 `json:"wordCount" gorm:"default:0"`
	IsPublished     uint8  `json:"published" gorm:"default:0"`
	IsDeleted       uint8  `json:"isDeleted" gorm:"default:0"`
	Status          uint8  `json:"status" gorm:"default:0"`
	CreateTime      uint64 `json:"createTime" gorm:"autoCreateTime"`
	// PubTime 发布时间。不可标 autoUpdateTime：GORM 在 Save 时会把非零值覆盖为
	// 当前时间，定时发布的未来时间会被冲掉，普通编辑也会悄悄重置发布时间。
	// 由 SavePost/UpdatePost 显式维护：新建缺省为当前时间，更新缺省保留原值。
	PubTime          uint64 `json:"pubTime"`
	LastModifiedTime uint64 `json:"lastModifiedTime" gorm:"autoUpdateTime"`
	CategoryId       uint32 `json:"categoryId" gorm:"default:0"`
	// Series 所属系列名（如"Spring Boot 系列"）。空表示不属于任何系列；
	// 同系列文章在文章页按发布时间列出，形成连载导航。
	Series string `json:"series" gorm:"default:''"`
	// 分类名称 只读不需要写入数据库
	CategoryName string   `json:"categoryName" gorm:"->;-:migration"`
	TagName      string   `json:"tagName" gorm:"->;-:migration"`
	Tags         []string `json:"tags" gorm:"-"`
}

func (Post) TableName() string {
	return TPostsTable
}

const TCategoryTable = "category"

type Category struct {
	CategoryId   uint32 `json:"categoryId" gorm:"primaryKey;autoIncrement"`
	CreateTime   uint64 `json:"createTime" gorm:"autoCreateTime"`
	CategoryName string `json:"categoryName" gorm:"default:''"`
	Note         string `json:"note" gorm:"default:''"`
	State        uint8  `json:"state" gorm:"default:1"`
	PostCount    uint32 `json:"postCount" gorm:"->"`
}

type CategoryCount struct {
	CategoryName string `json:"categoryName" gorm:"default:''"`
	Total        uint32 `json:"total"`
}

const TTagTable = "tag"

type Tag struct {
	TagId        uint32 `json:"tagId" gorm:"primaryKey;autoIncrement"`
	CreateTime   uint64 `json:"createTime" gorm:"autoCreateTime"`
	TagName      string `json:"tagName"`
	TagStyle     string `json:"tagStyle" gorm:"->;-:migration"`
	CardTagStyle string `json:"cardTagStyle" gorm:"->;-:migration"`
	PostCount    uint32 `json:"postCount" gorm:"->"`
}

const TAboutTable = "about"

type About struct {
	Id               uint32 `json:"Id" gorm:"primaryKey;autoIncrement"`
	Title            string `json:"title"  gorm:"default:''"`
	Note             string `json:"note" gorm:"default:''"`
	CreateTime       uint64 `json:"createTime" gorm:"autoCreateTime"`
	LastModifiedTime uint64 `json:"lastModifiedTime" gorm:"autoUpdateTime"`
}

type TagCount struct {
	TagName string `json:"tagName"`
	Ratio   int    `json:"ratio"`
}

const TPostTagTable = "post_tag"

type PostTag struct {
	PostId uint64 `json:"postId"`
	TagId  uint32 `json:"tagId"`
}

const TFriendLinkTable = "friend_link"

type FriendLink struct {
	Id               uint32 `json:"id" gorm:"primaryKey;autoIncrement"`
	Title            string `json:"title" gorm:"default:''"`
	LinkUrl          string `json:"linkUrl" gorm:"default:''"`
	LinkIcon         string `json:"linkIcon" gorm:"default:''"`
	State            uint8  `json:"state" gorm:"default:0"`
	Type             uint8  `json:"type"`
	LinkDesc         string `json:"linkDesc" gorm:"default:''"`
	CreateTime       uint64 `json:"createTime" gorm:"autoCreateTime"`
	LastModifiedTime uint64 `json:"LastModifiedTime" gorm:"autoUpdateTime"`
}

type User struct {
	UserId int64 `json:"userId" gorm:"primaryKey;autoIncrement"`

	/** 用户账号 */

	Username string `json:"username" gorm:"default:''"`

	/**用户密码*/
	Password string `json:"password" gorm:"default:''"`

	/** 用户昵称 */

	NickName string `json:"nickName" gorm:"default:''"`

	/** 用户邮箱 */
	Email string `json:"email" gorm:"default:''"`

	/** 手机号码 */
	Phonenumber string `json:"phonenumber" gorm:"default:''"`

	/** 用户性别 */

	Sex uint8 `json:"sex" gorm:"default:0"`

	/** 用户头像 */
	Avatar string `json:"avatar" gorm:"default:''"`
}

/**********页面显示*********/

type ArchivePosts struct {
	Title      string `json:"title"`
	PostSlug   string `json:"postSlug"`
	PubTime    uint64 `json:"pubTime"`
	CoverImage string `json:"coverImage"`
}

type TagCategoryPosts struct {
	Title      string `json:"title"`
	CoverImage string `json:"coverImage"`
	PostSlug   string `json:"postSlug"`
	PubTime    uint64 `json:"pubTime"`
}

type LatestPosts struct {
	Title           string `json:"title"`
	PostSlug        string `json:"postSlug"`
	Author          string `json:"author"`
	ContentAbstract string `json:"contentAbstract"`
	PubTime         uint64 `json:"pubTime"`
	CategoryId      uint32 `json:"categoryId"`
	CoverImage      string `json:"coverImage"`
	CategoryName    string `json:"categoryName"`
	TagName         string `json:"tagName"`
	Top             uint8  `json:"top"`
}

type SearchPost struct {
	Id         uint32 `json:"id"`
	Title      string `json:"title"`
	PostSlug   string `json:"postSlug"`
	Summary    string `json:"summary"`
	Highlight  string `json:"highlight"` // 匹配片段高亮
	CoverImage string `json:"coverImage"`
	CreateAt   int64  `json:"createAt"`
	Relevance  int    `json:"relevance"` // 相关度分数
}

// SeriesPost 系列导航条目：文章页尾部"本系列"列表用，按发布时间升序。
type SeriesPost struct {
	PostId   uint64 `json:"postId"`
	Title    string `json:"title"`
	PostSlug string `json:"postSlug"`
	PubTime  uint64 `json:"pubTime"`
}

// SearchResult 搜索结果包装
type SearchResult struct {
	Posts    []SearchPost `json:"posts"`
	Total    int64        `json:"total"`
	Keyword  string       `json:"keyword"`
	PageNum  int          `json:"pageNum"`
	PageSize int          `json:"pageSize"`
}

type ArchivesPosts struct {
	Archives map[string][]ArchivePosts `json:"archives"`
}

type CardInfo struct {
	Post     int64 `json:"post"`
	Tag      int64 `json:"tag"`
	Category int64 `json:"category"`
}
type Sidebar struct {
	CardInfo        CardInfo          `json:"cardInfo"`
	LatestPosts     []LatestPosts     `json:"latestPosts"`
	Category        []CategoryCount   `json:"categories"`
	Tag             []Tag             `json:"tags"`
	SidebarArchives []SidebarArchives `json:"sidebarArchives"`
	WebInfo         WebInfo           `json:"webInfo"`
}

type SidebarArchives struct {
	Year  string `json:"year"`
	Total int64  `json:"total"`
}

// WebInfo 网站资讯数据
type WebInfo struct {
	PostCount      int64  `json:"postCount"`      // 文章数目
	TotalWordCount int64  `json:"totalWordCount"` // 总字数
	RuntimeDays    int64  `json:"runtimeDays"`    // 运行天数
	LastUpdateTime string `json:"lastUpdateTime"` // 最后更新时间（展示格式：2006年1月2日 / 暂无文章）
	LastUpdateISO  string `json:"lastUpdateISO"`  // 最后更新时间（RFC3339，供前端 data-lastPushDate 相对时间计算）
	SiteStartDate  string `json:"siteStartDate"`  // 网站创建日期
	SitePV         int64  `json:"sitePV"`         // 本站总访问量（system_access_log 的 SUM(pv)）
	SiteUV         int64  `json:"siteUV"`         // 本站访客数（system_access_log 的 COUNT(DISTINCT ip)）
}
