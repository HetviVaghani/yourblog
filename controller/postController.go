package controller

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hetvivaghani/yourblog/database"
	"github.com/hetvivaghani/yourblog/models"
	"github.com/hetvivaghani/yourblog/util"
	"gorm.io/gorm"
)

func CreatePost(c *fiber.Ctx) error {
	var blogpost models.Blog
	if err := c.BodyParser(&blogpost); err != nil {
		fmt.Println("Unable to parse body")
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid payload",
		})
	}
	if err := database.DB.Create(&blogpost).Error; err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Unable to create post",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Congratulations! Your post is live",
	})
}

func AllPost(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit := 5
	offset := (page - 1) * limit
	var total int64
	var getblog []models.Blog
	database.DB.Preload("User").Offset(offset).Limit(limit).Find(&getblog)
	database.DB.Model(&models.Blog{}).Count(&total)
	return c.JSON(fiber.Map{
		"data": getblog,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": math.Ceil(float64(total) / float64(limit)),
		},
	})
}

func DetailPost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var blogpost models.Blog
	database.DB.Where("id = ?", id).Preload("User").First(&blogpost)
	return c.JSON(fiber.Map{
		"data": blogpost,
	})
}

func UpdatePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var blog models.Blog
	if err := database.DB.First(&blog, id).Error; err != nil {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Blog post not found",
		})
	}

	if err := c.BodyParser(&blog); err != nil {
		fmt.Println("Unable to parse body")
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid payload",
		})
	}

	if err := database.DB.Model(&blog).Updates(blog).Error; err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Unable to update post",
		})
	}
	return c.JSON(fiber.Map{
		"message":"post updated successfully",
	})
}

func UniquePost(c *fiber.Ctx) error{
	cookie := c.Cookies("jwt")
	id, _:= util.ParseJwt(cookie)
	var blog []models.Blog
	database.DB.Model(&blog).Where("user_id=?",id).Preload("User").Find(&blog)

	return c.JSON(blog)
}

func DeletePost(c *fiber.Ctx) error{
	id,_ := strconv.Atoi(c.Params("id"))
	blog := models.Blog{
		Id: uint(id),
	}
	deleteQuery := database.DB.Delete(&blog)
	if errors.Is(deleteQuery.Error, gorm.ErrRecordNotFound){
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Opps!, record Not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "post deleted Succesfully",
	})
}