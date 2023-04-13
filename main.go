package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nickemma/go-jwt/controllers"
	"github.com/nickemma/go-jwt/initializers"
)

func init() {
  initializers.LoadEnvVariables()
		initializers.ConnectToDb()
		initializers.SyncDatabase()
}

func main() {

	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	r.Run()
}