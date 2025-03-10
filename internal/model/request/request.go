package request

type PageRequest struct {
	PageNum  int `form:"pageNum" ,json:"pageNum"`
	PageSize int `form:"pageSize" ,json:"pageSize"`
}

type PostRequest struct {
	PageRequest
	Title      string  `form:"title" ,json:"title"`
	CategoryId *uint32 `form:"categoryId" ,json:"categoryId"`
	Published  *uint8  `form:"published" ,json:"published"`
}

type TopRequest struct {
	PostId int64 `form:"postId"  json:"postId"`
	Top    uint8 `form:"top" json:"top"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginLogQuery struct {
	PageRequest
	Success *uint8 `form:"success"`
}
type AccessLogQuery struct {
	PageRequest
	IP     string `form:"ip"`
	Status *int   `form:"status"`
	Start  string `form:"start"`
	End    string `form:"end"`
}

type ChangePasswordRequest struct {
	OldPassword string `form:"oldPassword" json:"oldPassword"`
	NewPassword string `form:"newPassword" json:"newPassword"`
	Username    string `json:"-"`
}
