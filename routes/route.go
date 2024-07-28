package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hetvivaghani/yourblog/controller"
	"github.com/hetvivaghani/yourblog/middleware"
)

func Setup(c *fiber.App) {
	c.Post("/api/register", controller.Register)
	c.Post("/api/login", controller.Login)

	c.Use(middleware.IsAuthenticate)
	c.Post("/api/post", controller.CreatePost)
	c.Get("/api/allpost", controller.AllPost)
	c.Get("/api/allpost/:id", controller.DetailPost)
	c.Put("/api/updatepost/:id", controller.UpdatePost)
	c.Get("/api/uniquepost", controller.UniquePost)
	c.Delete("/api/deletepost/:id", controller.DeletePost)
	c.Post("/api/upload-image", controller.Upload)
	c.Static("/api/uploads", "./uploads")
}
