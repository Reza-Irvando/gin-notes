package main

import (
	"gin-notes/configs"
	"gin-notes/database"
	"gin-notes/handlers"
	"gin-notes/middlewares"

	"github.com/gin-gonic/gin"
)

func main(){
	r := gin.Default()
	db, err := configs.InitDB()
	if err != nil {
		panic(err)
	}

	database.Migrate(db)
	
	// Public routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Authentication routes
	r.POST("/register", handlers.Register(db))
	r.POST("/login", handlers.Login(db))
	r.POST("/logout", middlewares.AuthMiddleware(), handlers.Logout)

	// Protected routes - Notes
	notes := r.Group("/notes")
	notes.Use(middlewares.AuthMiddleware())
	{
		notes.POST("add", handlers.CreateNote(db))
		notes.GET("", handlers.GetAllNotes(db))
		notes.GET("/:id", handlers.GetNoteDetail(db))
		notes.PUT("/update/:id", handlers.UpdateNote(db))
		notes.DELETE("/:id", handlers.DeleteNote(db))
	}

	r.Run()
}