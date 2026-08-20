package jwttoken

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
)

// MySecret JWT 签名密钥，由应用启动时通过 Configure 注入
// （来源：config.yaml 的 jwt.secret 或环境变量 JWT_SECRET）。
// 此前版本在 init() 里读配置并回退硬编码默认值，现为显式注入。
var MySecret []byte

// Configure 注入签名密钥，应用启动阶段在 config.Load 之后调用。
// 密钥为空返回错误，与 config.ValidateConfig 的 fail-fast 保持一致。
func Configure(secret []byte) error {
	if len(secret) == 0 {
		return errors.New("JWT 密钥为空：请配置 jwt.secret 或环境变量 JWT_SECRET")
	}
	if len(secret) < 32 {
		slog.Warn("JWT 密钥长度建议至少 32 字符")
	}
	MySecret = secret
	return nil
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
