package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gin-notes/database"
	"gin-notes/configs"
	"gin-notes/handlers"

)

func main(){
	r := gin.Default()
	db, err := configs.InitDB()
	if err != nil {
		panic(err)
	}

	database.Migrate(db)
	
	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/login", handlers.Login(db))
	r.POST("/register", handlers.Register(db))

	r.Run()
}