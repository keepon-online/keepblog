package request

type PageRequest struct {
	PageNum  int `form:"pageNum" json:"pageNum"`
	PageSize int `form:"pageSize" json:"pageSize" binding:"max=100"`
}

type PostRequest struct {
	PageRequest
	Title      string  `form:"title" json:"title" binding:"max=200"`
	CategoryId *uint32 `form:"categoryId" json:"categoryId"`
	Published  *uint8  `form:"published" json:"published"`
}

type TopRequest struct {
	PostId int64 `form:"postId" json:"postId" binding:"required"`
	Top    uint8 `form:"top" json:"top"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

type LoginLogQuery struct {
	PageRequest
	Success *uint8 `form:"success"`
}
type AccessLogQuery struct {
	PageRequest
	IP     string `form:"ip" binding:"max=45"`
	Status *int   `form:"status"`
	Start  string `form:"start"`
	End    string `form:"end"`
}

type ChangePasswordRequest struct {
	OldPassword string `form:"oldPassword" json:"oldPassword" binding:"required,min=6,max=128"`
	NewPassword string `form:"newPassword" json:"newPassword" binding:"required,min=6,max=128"`
	Username    string `json:"-"`
}
