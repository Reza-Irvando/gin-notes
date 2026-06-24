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

	// Seeder database
	database.Seed(db)
	
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
		notes.POST("", handlers.CreateNote(db))
		notes.GET("", handlers.GetAllNotes(db))
		notes.GET("/:id", handlers.GetNoteDetail(db))
		notes.PUT("/:id", handlers.UpdateNote(db))
		notes.DELETE("/:id", handlers.DeleteNote(db))

		// Favorite routes
		notes.POST("/:id/favorite", handlers.AddToFavorite(db))
		notes.GET("/favorites", handlers.GetFavoriteNotes(db))
		notes.DELETE("/:id/favorite", handlers.RemoveFromFavorite(db))
		// notes.GET("/favorite/check/:id", handlers.IsFavorite(db))
	}

	// Protected routes - Categories
	categories := r.Group("/categories")
	categories.Use(middlewares.AuthMiddleware())
	{
		categories.POST("", handlers.CreateCategory(db))
		categories.GET("", handlers.GetAllCategories(db))
		categories.PUT("/:id", handlers.UpdateCategory(db))
		categories.DELETE("/:id", handlers.DeleteCategory(db))
	}

	// Protected routes - Tags
	tags := r.Group("/tags")
	tags.Use(middlewares.AuthMiddleware())
	{
		tags.POST("", handlers.CreateTag(db))
		tags.GET("", handlers.GetAllTags(db))
		tags.PUT("/:id", handlers.UpdateTag(db))
		tags.DELETE("/:id", handlers.DeleteTag(db))

		// Note-Tag relationship
		tags.POST("/add", handlers.AddTagToNote(db))
		tags.DELETE("/remove", handlers.RemoveTagFromNote(db))
	}

	// Protected routes - Activity Log
	activities := r.Group("/activities")
	activities.Use(middlewares.AuthMiddleware())
	{
		activities.GET("", handlers.GetActivityLog(db))
		activities.GET("/:entity/:entityId", handlers.GetEntityActivityLog(db))
	}

	r.Run()
}