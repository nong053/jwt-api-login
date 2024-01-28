package user

import (
	"net/http"
	"nong/jwt-api-login/orm"

	"github.com/gin-gonic/gin"
)

func ReadAll(c *gin.Context) {

	var users []orm.User
	orm.Db.Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "User read success",
		"users":   users,
	})

}
func Profile(c *gin.Context) {

	userId := c.MustGet("userId").(float64)

	var users orm.User
	orm.Db.First(&users, userId)
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "User read success",
		"user":    users,
	})

}
