package auth

import (
	"net/http"
	"nong/jwt-api-login/orm"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"
)

// Binding from JSON
var hmacSampleSecret []byte

type RegisterBody struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
	Fullname string `json:"fullname"  binding:"required"`
	Avatar   string `json:"avatar"  binding:"required"`
}

func Register(c *gin.Context) {
	var json RegisterBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check user exist
	// Get first matched record
	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExist)

	if userExist.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Exists"})
		return
	}
	//create user
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	user := &orm.User{Username: json.Username, Password: string(encryptedPassword), Fullname: json.Fullname, Avatar: json.Avatar}
	orm.Db.Create(&user)
	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User registered success", "userId": user.ID})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User registered failed", "userId": user.ID})
	}
}

// Binding from JSON

type LoginBody struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

func Login(c *gin.Context) {
	var json LoginBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check user

	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExist)

	if userExist.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Does Not Exist"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(json.Password))
	if err == nil {

		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": userExist.ID,
			//"nbf":    time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
			"exp": time.Now().Add(time.Minute * 1).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, _ := token.SignedString(hmacSampleSecret)

		//println(tokenString)

		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Login success", "token": tokenString})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Login failed"})
		return
	}
}
