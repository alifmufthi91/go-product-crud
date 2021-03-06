package middlewares

import (
	"fmt"
	"net/http"
	"product-crud/app"
	"product-crud/config"
	"product-crud/util/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.Contains(authHeader, "Bearer") {
		result := gin.H{
			"message": "not authorized",
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", -1)
	token, err := jwt.ParseWithClaims(tokenString, &app.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Env.JWTSECRET), nil
	})

	if token != nil && err == nil {
		claims, _ := token.Claims.(*app.MyCustomClaims)
		c.Set("user", claims)
		logger.Info(`token verified, claims = %+v`, claims)
	} else {
		result := gin.H{
			"message": "not authorized",
			"error":   err.Error(),
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
	}

}
