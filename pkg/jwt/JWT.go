package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// 一些常量
var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "xieang"
)

type LoginResult struct {
	User  interface{}
	Token string
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// CustomClaims 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

// NewJWT 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(getSignKey()),
	}
}

// GetSignKey 获取signKey
func getSignKey() string {
	return SignKey
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// GenerateToken 生成令牌  创建jwt风格的token
func GenerateToken(id string, TokenTime int64) (LoginResult, error) {
	j := NewJWT()
	claims := CustomClaims{
		id,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 600), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + TokenTime), // 过期时间 一小时
			Issuer:    "Liang",                         //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		return LoginResult{}, err
	}

	data := LoginResult{
		User:  id,
		Token: token,
	}
	return data, nil
}

// ParseToken 解析Tokne
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

// JWTAuth 中间件，检查token
func JWTAuth(token string) error {
	if token == "" {
		return errors.New("请求未携带token，无权限访问")
	}
	j := NewJWT()
	// parseToken 解析token包含的信息
	_, err := j.ParseToken(token)
	//fmt.Println("claims", claims)
	if err != nil {
		if err == TokenExpired {
			return errors.New("授权已过期,请重新登陆")
		}
		return err
	}
	return nil
}
