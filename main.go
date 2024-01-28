package main

import (
	AuthController "nong/jwt-api-login/controller/auth"
	UserController "nong/jwt-api-login/controller/user"
	"nong/jwt-api-login/middleware"
	"nong/jwt-api-login/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "golang.org/x/crypto/bcrypt"
)

// Binding from JSON

type Register struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
	Fullname string `json:"fullname"  binding:"required"`
	Avatar   string `json:"avatar"  binding:"required"`
}

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		println("Error loading .env file")
	}

	orm.InitDB()

	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)
	authorization := r.Group("/users", middleware.JWTAuthen())
	authorization.GET("/readall", UserController.ReadAll)
	authorization.GET("/profile", UserController.Profile)

	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}
