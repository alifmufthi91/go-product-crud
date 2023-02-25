package middlewares

import (
	"fmt"
	"product-crud/config"
	"product-crud/dto/app"
	errorUtil "product-crud/util/error"
	"product-crud/util/logger"
	responseUtil "product-crud/util/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.Contains(authHeader, "Bearer") {
		err := errorUtil.Unauthorized("Internal Error")
		responseUtil.Fail(c, app.ErrorHttpResponse{
			Message:    err.Error(),
			HttpStatus: err.HttpStatus,
			ErrorName:  err.ErrorName,
		})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", -1)
	token, err := jwt.ParseWithClaims(tokenString, &app.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Env.JWTSECRET), nil
	})

	if token != nil && err == nil {
		claims, _ := token.Claims.(*app.UserClaims)
		c.Set("user", claims)
		logger.Info(`token verified, claims = %+v`, claims)
	} else {
		err := errorUtil.Unauthorized("Not Authorized")
		responseUtil.Fail(c, app.ErrorHttpResponse{
			Message:    err.Error(),
			HttpStatus: err.HttpStatus,
			ErrorName:  err.ErrorName,
		})
	}

}
