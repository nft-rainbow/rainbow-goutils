package jwtmiddleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

const (
	USER_TYPE = "user_type"
)

func GetDefaultIdFromClaim(c *gin.Context) uint {
	// return c.GetUint(JwtDefaultIdentityKey)
	claims := jwt.ExtractClaims(c)
	idVal := claims[JwtDefaultIdentityKey]
	return uint(idVal.(float64))
}

func GetUserTypeFromClaim(c *gin.Context) int {
	claims := jwt.ExtractClaims(c)
	idVal := claims[USER_TYPE]
	return int(idVal.(float64))
}

func ExtractUserTypeFromToken(mw *jwt.GinJWTMiddleware, token string) (int, error) {
	t, err := mw.ParseTokenString(token)
	if err != nil {
		return 0, err
	}
	userTypeFloat := (t.Claims.(jwtv4.MapClaims))[USER_TYPE]
	return int(userTypeFloat.(float64)), nil
}
