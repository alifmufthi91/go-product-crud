package middlewares

import (
	"fmt"
	"product-crud/config"
	"product-crud/dto/app"
	"product-crud/util/apiresponse"
	"product-crud/util/errorhandler"
	"product-crud/util/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.Contains(authHeader, "Bearer") {
		err := errorhandler.Unauthorized("Internal Error")
		apiresponse.Fail(c, app.ErrorHttpResponse{
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

		return []byte(config.GetEnv().JWTSECRET), nil
	})

	if token != nil && err == nil {
		claims, _ := token.Claims.(*app.UserClaims)
		c.Set("user", claims)
		logger.Info(`token verified, claims = %+v`, claims)
	} else {
		err := errorhandler.Unauthorized("Not Authorized")
		apiresponse.Fail(c, app.ErrorHttpResponse{
			Message:    err.Error(),
			HttpStatus: err.HttpStatus,
			ErrorName:  err.ErrorName,
		})
	}

}
