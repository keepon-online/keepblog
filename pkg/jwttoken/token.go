package jwttoken

import (
	"os"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
)

// MySecret JWT 签名密钥（从环境变量读取，必须至少32字节）
var MySecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// 开发环境默认值，生产环境必须设置环境变量
		secret = "sq44jvTJfbBlsZSvNm440So77O9J9TA"
		slog.Warn("JWT_SECRET not set, using default value. Set JWT_SECRET environment variable in production!")
	}
	if len(secret) < 32 {
		slog.Warn("JWT_SECRET should be at least 32 characters for security")
	}
	MySecret = []byte(secret)
}

// MyClaims JWT 自定义声明
type MyClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// CreateToken username
func CreateToken(username string) (tokenString string, err error) {
	claim := MyClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)), // 过期时间3小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                     // 生效时间
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
	tokenString, err = token.SignedString(MySecret)
	return tokenString, err
}

// AccessToken username
func AccessToken(username string) (tokenString string, err error) {
	claim := MyClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // 过期时间2小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                    // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                    // 生效时间
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
	tokenString, err = token.SignedString(MySecret)
	return tokenString, err
}

// RefreshToken 这里传入的是手机号，因为我项目登陆用的是手机号和密码
func RefreshToken(username string) (tokenString string, err error) {
	claim := MyClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 过期时间7天
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                         // 生效时间
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
	tokenString, err = token.SignedString(MySecret)
	return tokenString, err
}

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		return MySecret, nil
	}
}

func ParseToken(tokensStr string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokensStr, &MyClaims{}, Secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			}

			return nil, errors.New("couldn't handle this token")
		}

		return nil, errors.New("couldn't handle this token")
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}
