package middlewares

import (
	"net/http"
	"product-crud/app"
	"product-crud/util/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth(c *gin.Context) {
	bearer := c.Request.Header.Get("Authorization")
	tokenString := strings.Fields(bearer)[1]
	// var token *string
	// var err error
	at(time.Unix(0, 0), func() {
		token, err := jwt.ParseWithClaims(tokenString, &app.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if claims, ok := token.Claims.(*app.MyCustomClaims); ok && token.Valid {
			logger.Info("%v %v %v", claims.FullName, claims.StandardClaims.ExpiresAt, claims.IssuedAt)
			logger.Info("%v", token.Valid)
		} else {
			logger.Error(err.Error())
		}

		if token != nil && err == nil {
			logger.Info("token verified")
		} else {
			result := gin.H{
				"message": "not authorized",
				"error":   err.Error(),
			}
			c.JSON(http.StatusUnauthorized, result)
			c.Abort()
		}
	})
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	if jwt.GetSigningMethod("HS256") != token.Method {
	// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	// 	}

	// 	return []byte("secret"), nil
	// })

}

func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
}
