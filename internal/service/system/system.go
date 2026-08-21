package system

import (
	"fmt"
	"time"

	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
	"gorm.io/gorm"

	"gitee.com/jieepre/go-site/internal/model"
	"gitee.com/jieepre/go-site/internal/model/request"
	"gitee.com/jieepre/go-site/internal/model/response"
	"gitee.com/jieepre/go-site/internal/model/system"
	"gitee.com/jieepre/go-site/pkg"
	"gitee.com/jieepre/go-site/pkg/area"
	"gitee.com/jieepre/go-site/pkg/jwttoken"
	"gitee.com/jieepre/go-site/pkg/page"
)

type Service struct {
	db *gorm.DB
}

func NewSystemService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (service *Service) Login(request request.LoginRequest, c *gin.Context) (*response.LoginResponse, error) {
	user := model.User{}
	ipL := pkg.Ip2long(c.ClientIP())

	log := system.LoginLog{
		UA:      c.Request.UserAgent(),
		Referer: c.Request.Referer(),
		Ip:      &ipL,
		Area:    area.Area(ipL),
	}
	if err := service.db.Where("username=?", request.Username).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		log.Note = "用户名不存在"
		errSuccess := 2
		log.Success = &errSuccess
		go service.LoginLog(log)
		slog.Errorf("登录失败: %s", err.Error())
		return nil, errors.New("用户名或密码错误")
	}
	password := user.Password
	hash := pkg.CheckPasswordHash(request.Password, password)
	if !hash {
		log.Note = "用户名或密码错误"
		errSuccess := 2
		log.Success = &errSuccess
		go service.LoginLog(log)
		return nil, errors.New("用户名或密码错误")
	}
	accessToken, err := jwttoken.AccessToken(user.Username)
	if err != nil {
		log.Note = "系统错误" + err.Error()
		errSuccess := 2
		log.Success = &errSuccess
		go service.LoginLog(log)
		slog.Errorf("登录失败: %s", err.Error())
		return nil, errors.New("系统错误，请联系管理员")
	}

	refreshToken, err := jwttoken.RefreshToken(user.Username)
	if err != nil {
		log.Note = "系统错误" + err.Error()
		errSuccess := 2
		log.Success = &errSuccess
		go service.LoginLog(log)
		slog.Errorf("登录失败: %s", err.Error())
		return nil, errors.New("系统错误，请联系管理员")
	}
	loginResponse := response.LoginResponse{
		Username:     user.Username,
		Roles:        []string{"admin"},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expires:      time.Now().Add(30 * time.Minute).Format("2006/01/02 15:04:05"),
	}
	log.Note = "登录成功"
	go service.LoginLog(log)
	return &loginResponse, nil
}

func (service *Service) RefreshToken(refreshToken string) (*response.RefreshToken, error) {
	parseToken, err := jwttoken.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("请重新登录")
	}

	// 检查令牌是否已被撤销（密码修改后的令牌失效检查）
	if parseToken.IssuedAt != nil {
		if !jwttoken.IsUserTokenValid(parseToken.Username, parseToken.IssuedAt.Unix()) {
			return nil, errors.New("令牌已失效，请重新登录")
		}
	}

	accessToken, err := jwttoken.AccessToken(parseToken.Username)
	if err != nil {
		return nil, errors.New("系统错误，请联系管理员")
	}

	// 令牌轮换：生成新的refreshToken替换旧的
	newRefreshToken, err := jwttoken.RefreshToken(parseToken.Username)
	if err != nil {
		return nil, errors.New("系统错误，请联系管理员")
	}

	return &response.RefreshToken{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		Expires:      time.Now().Add(30 * time.Minute).Format("2006/01/02 15:04:05"),
	}, nil
}

// ChangePassword 修改密码
func (service *Service) ChangePassword(req request.ChangePasswordRequest) error {
	user := model.User{}
	if err := service.db.Where("username=?", req.Username).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Errorf("用户名%s不存在", err.Error())
		return fmt.Errorf("用户名%s不存在", req.Username)
	}
	password := user.Password
	hash := pkg.CheckPasswordHash(req.OldPassword, password)
	if !hash {
		return errors.New("旧密码错误")
	}

	newPassword, _ := pkg.HashPassword(req.NewPassword)

	if err := service.db.Where("username=?", req.Username).Updates(model.User{Password: newPassword}).Error; err != nil {
		slog.Errorf("修改密码失败", err.Error())
		return fmt.Errorf("密码修改失败")
	}

	// 使该用户所有旧令牌失效（密码修改后的安全措施）
	jwttoken.InvalidateUserTokens(req.Username, time.Now().Unix())
	slog.Infof("用户 %s 密码已修改，旧令牌已失效", req.Username)

	return nil
}

// GetUserInfo 获取用户个人信息
func (service *Service) GetUserInfo(username string) (*response.UserInfoResponse, error) {
	user := model.User{}
	if err := service.db.Where("username", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		slog.Errorf("获取用户信息失败: %v", err)
		return nil, errors.New("获取用户信息失败")
	}

	// 获取最近登录信息
	var loginLog system.LoginLog
	var loginCount int64
	service.db.Model(&system.LoginLog{}).Where("success = ?", 1).Count(&loginCount)
	service.db.Where("success = ?", 1).Order("create_at desc").First(&loginLog)

	loginIP := ""
	loginTime := ""
	if loginLog.Ip != nil {
		loginIP = pkg.Long2ip(*loginLog.Ip)
	}
	if loginLog.CreatedAt > 0 {
		loginTime = time.Unix(int64(loginLog.CreatedAt), 0).Format("2006-01-02 15:04:05")
	}

	return &response.UserInfoResponse{
		UserId:       user.UserId,
		Username:     user.Username,
		NickName:     user.NickName,
		Email:        user.Email,
		Phonenumber:  user.Phonenumber,
		Sex:          user.Sex,
		Avatar:       user.Avatar,
		Role:         "超级管理员",
		LoginIP:      loginIP,
		LoginTime:    loginTime,
		LoginCount:   int(loginCount),
		RegisterTime: "2021-01-01 00:00:00", // TODO: 添加用户注册时间字段
	}, nil
}

// UpdateProfile 更新用户资料
func (service *Service) UpdateProfile(username string, req request.UpdateProfileRequest) error {
	user := model.User{}
	if err := service.db.Where("username", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return errors.New("更新用户资料失败")
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.NickName != "" {
		updates["nick_name"] = req.NickName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phonenumber != "" {
		updates["phonenumber"] = req.Phonenumber
	}
	if req.Sex != nil {
		updates["sex"] = *req.Sex
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) == 0 {
		return nil
	}

	if err := service.db.Model(&user).Updates(updates).Error; err != nil {
		slog.Errorf("更新用户资料失败: %v", err)
		return errors.New("更新用户资料失败")
	}

	return nil
}

func (service *Service) LoginLog(log system.LoginLog) {
	service.db.Save(&log)
}
func (service *Service) LoginLogList(req request.LoginLogQuery) (*page.Info, error) {
	logs := make([]system.LoginLog, 0)
	total := int64(0)
	service.db.Model(system.LoginLog{}).Scopes(success(req.Success)).Count(&total)
	service.db.Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Order("create_at desc").Scopes(success(req.Success)).Find(&logs)
	bInfo := page.PaginationInfo(logs, req.PageNum, req.PageSize, int(total))
	return bInfo, nil
}

func (service *Service) AccessLogList(req request.AccessLogQuery) (*page.Info, error) {
	logs := make([]system.AccessLog, 0)
	total := int64(0)
	service.db.Model(system.AccessLog{}).Scopes(daterange(req.Start, req.End), ip(req.IP), status(req.Status)).Count(&total)
	service.db.Offset((req.PageNum-1)*req.PageSize).
		Limit(req.PageSize).
		Order("create_at desc").Scopes(daterange(req.Start, req.End), ip(req.IP), status(req.Status)).Find(&logs)
	bInfo := page.PaginationInfo(logs, req.PageNum, req.PageSize, int(total))
	return bInfo, nil
}

func success(success *uint8) func(db *gorm.DB) *gorm.DB {
	if success == nil || *success == 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("success", success)
	}
}

func daterange(start, end string) func(db *gorm.DB) *gorm.DB {
	if start == "" || end == "" {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(" strftime('%Y-%m-%d', create_at, 'unixepoch' ) between ? and ? ", start, end)
	}
}

func ip(ip string) func(db *gorm.DB) *gorm.DB {
	if ip == "" {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(" ip = ?", pkg.Ip2long(ip))
	}
}

func status(code *int) func(db *gorm.DB) *gorm.DB {
	if code == nil || *code == 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status", code)
	}
}
