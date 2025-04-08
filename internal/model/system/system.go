package system

type User struct {
	UserId int64 `json:"userId" gorm:"primaryKey;autoIncrement"`

	/** 用户账号 */

	UserName string `json:"userName" gorm:"default:''"`

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

// AccessLog 访问日志统计
type AccessLog struct {
	Id        *int    `json:"id" gorm:"autoIncrement:id;primaryKey;comment:''"`
	Ip        *uint32 `json:"ip" gorm:"default:0"`
	URL       string  `json:"url" gorm:"default:'';column:url"`
	PV        *int    `json:"pv" gorm:"default:0;column:pv"`
	UV        *int    `json:"uv" gorm:"default:0;column:uv"`
	UA        string  `json:"ua" gorm:"default:'';column:ua"`
	Status    int     `json:"status" gorm:"default:0"`
	Referer   string  `json:"referer" gorm:"default:'';column:referer"`
	Area      string  `json:"area" gorm:"default:''"`
	CreatedAt int     `json:"createTime" gorm:"column:create_at;autoCreateTime"`
	UpdatedAt int     `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (AccessLog) TableName() string {

	return "system_access_log"
}

// LoginLog 登录日志记录
type LoginLog struct {
	Id        *int    `json:"id" gorm:"autoIncrement:id;primaryKey;comment:''"`
	Ip        *uint32 `json:"ip" gorm:"default:0"`
	UA        string  `json:"ua" gorm:"default:'';column:ua"`
	Note      string  `json:"note"`
	Referer   string  `json:"referer" gorm:"default:'';column:referer"`
	Area      string  `json:"area" gorm:"default:''"`
	Success   *int    `json:"success" gorm:"comment:1登录成功2登录失败;default:1"`
	CreatedAt int     `json:"createTime" gorm:"column:create_at;autoCreateTime"`
	UpdatedAt int     `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (LoginLog) TableName() string {

	return "system_login_log"
}

// WebSite 网站配置
type WebSite struct {
	Id          int    `json:"id" gorm:"autoIncrement:id;primaryKey;comment:''"`
	Icp         string `json:"icp"`         //icp备案号
	Notice      string `json:"notice"`      //公告
	Title       string `json:"title"`       //页面标题
	Description string `json:"description"` //页面简述
	URL         string `json:"url"`         //页面地址
	Keywords    string `json:"keywords"`    //关键词
	Copyright   string `json:"copyright"`   //版权
	BaiduStat   string `json:"stat"`        //百度统计
	BaiduSite   string `json:"site"`        //百度收录
	Github      string `json:"github"`      //GitHub
	Gitee       string `json:"gitee"`       //码云
	CreatedAt   int    `json:"createTime" gorm:"column:create_at;autoCreateTime"`
	UpdatedAt   int    `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (WebSite) TableName() string {

	return "web_site"
}

type PushSite struct {
	Remain      int    `json:"remain"`                                            //剩余可用的push数量
	CreatedAt   int    `json:"createTime" gorm:"column:create_at;autoCreateTime"` //创建时间
	UpdatedAt   int    `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"` //更新时间
	Success     int    `json:"success"`                                           //成功推送的url条数
	NotSameSite string `json:"notSameSite" gorm:"not_same_site"`                  //由于不是本站url而未处理的url列表
	NotValid    string `json:"notValid" gorm:"not_valid"`                         //不合法的url列表
	Error       int    `json:"error"`                                             //错误码
	Message     string `json:"message"`                                           //错误消息
}

func (PushSite) TableName() string {

	return "system_push_site"
}
