package jwtmiddleware

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewJwtMiddleware(
	name string, key string, timeout time.Duration,
	authenticator func(c *gin.Context) (interface{}, error),
) (*jwt.GinJWTMiddleware, error) {
	jwtAuthMw, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       name,
		Key:         []byte(key),
		Timeout:     timeout,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: JwtDefaultIdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					JwtDefaultIdentityKey: v.Id,
					"user_type":           v.UserType,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			id := claims[JwtDefaultIdentityKey]
			return uint(id.(float64))
		},
		Authenticator: authenticator,
		Unauthorized: func(c *gin.Context, code int, message string) {
			errCode := 100
			if message == "Token is expired" {
				errCode = 198
			}
			if message == "auth header is empty" {
				errCode = 199
			}
			_errCode, errCodeExist := c.Get("code")
			_message, messageExist := c.Get("message")
			if errCodeExist {
				errCode = _errCode.(int)
			}
			if messageExist {
				message = _message.(string)
			}
			c.JSON(code, UnauthorizedResp{errCode, message})
		},
		TokenLookup:   "header: Authorization", // cookie: jwt, query: token
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	setLoginResponse(jwtAuthMw)

	if err != nil {
		return nil, errors.WithMessage(err, "failed to initialize jwt middleware")
	}

	return jwtAuthMw, nil
}

func setLoginResponse(mw *jwt.GinJWTMiddleware) {
	loginResponseFunc := func(c *gin.Context, code int, token string, time time.Time) {
		userType, err := ExtractUserTypeFromToken(mw, token)
		if err != nil {
			c.JSON(http.StatusBadRequest, UnauthorizedResp{100, err.Error()})
			return
		}
		c.JSON(code, LoginResp{time, token, userType})
	}
	mw.LoginResponse = loginResponseFunc
	mw.RefreshResponse = loginResponseFunc
}
