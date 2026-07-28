package model

const TPostsTable = "post"

type Post struct {
	PostId           uint64 `json:"postId" gorm:"primaryKey;autoIncrement"`
	Title            string `json:"title" gorm:"default:''"`
	PostSlug         string `json:"postSlug" gorm:"default:''"`
	Author           string `json:"author" gorm:"default:''"`
	CoverImage       string `json:"coverImage" gorm:"default:''"`
	PostContent      string `json:"postContent" gorm:"default:''"`
	PostContentHtml  string `json:"postContentHtml" gorm:"default:''"`
	Summary          string `json:"summary" gorm:"default:''"`
	Type             uint8  `json:"type" gorm:"default:1"`
	Top              uint8  `json:"top" gorm:"default:0"`
	ReadCount        uint32 `json:"readCount" gorm:"default:0"`
	WordCount        uint32 `json:"wordCount" gorm:"default:0"`
	IsPublished      uint8  `json:"published" gorm:"default:0"`
	IsDeleted        uint8  `json:"isDeleted" gorm:"default:0"`
	Status           uint8  `json:"status" gorm:"default:0"`
	CreateTime       uint64 `json:"createTime" gorm:"autoCreateTime"`
	PubTime          uint64 `json:"pubTime" gorm:"autoUpdateTime"`
	LastModifiedTime uint64 `json:"lastModifiedTime" gorm:"autoUpdateTime"`
	CategoryId       uint32 `json:"categoryId" gorm:"default:0"`
	// 分类名称 只读不需要写入数据库
	CategoryName string   `json:"categoryName" gorm:"->;-:migration"`
	TagName      string   `json:"tagName" gorm:"->;-:migration"`
	Tags         []string `json:"tags" gorm:"-"`
}

const TCategoryTable = "category"

type Category struct {
	CategoryId   uint32 `json:"categoryId" gorm:"primaryKey;autoIncrement"`
	CreateTime   uint64 `json:"createTime" gorm:"autoCreateTime"`
	CategoryName string `json:"categoryName" gorm:"default:''"`
	Note         string `json:"note" gorm:"default:''"`
	State        uint8  `json:"state" gorm:"default:1"`
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
	LastUpdateTime string `json:"lastUpdateTime"` // 最后更新时间
	SiteStartDate  string `json:"siteStartDate"`  // 网站创建日期
}
