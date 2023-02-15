package util

import (
	"errors"
	"product-crud/app"

	"github.com/gin-gonic/gin"
)

func GetUserClaims(c *gin.Context) (*app.UserClaims, error) {
	userClaims, _ := c.Get("user")
	user, ok := userClaims.(*app.UserClaims)
	if !ok {
		return nil, errors.New("userClaims type is not correct")
	}
	return user, nil
}
