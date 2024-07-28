package controller

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/hetvivaghani/yourblog/database"
	"github.com/hetvivaghani/yourblog/models"
	"github.com/hetvivaghani/yourblog/util"
)

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return Re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
		return c.Status(400).JSON(fiber.Map{
			"message": "Unable to parse body",
		})
	}

	// Check if password is less than 6 characters
	if len(data["password"].(string)) <= 6 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Password must be greater than 6 characters",
		})
	}

	// Validate email
	if !validateEmail(strings.TrimSpace(data["email"].(string))) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid Email Address",
		})
	}

	// Check if email already exists
	database.DB.Where("email = ?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.Id != 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// Create user
	user := models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Phone:     data["phone"].(string),
		Email:     strings.TrimSpace(data["email"].(string)),
	}

	user.SetPassword(data["password"].(string))
	if err := database.DB.Create(&user).Error; err != nil {
		log.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not create user",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"user":    user,
		"message": "Account created successfully",
	})
}

func Login(c *fiber.Ctx)error {
	var data map[string] string

	if err:=c.BodyParser(&data); err!= nil {
		fmt.Println("Unable to parse body")
	}
	var user models.User
	database.DB.Where("email=?", data["email"]).First(&user)
	if user.Id ==  0{
		c.Status(404)
		return c.JSON(fiber.Map{
			"message":"Email Address doesn't exist,kindly create an account",
		})
	}
	if err:= user.ComparePassword(data["password"]); err!=nil{
		c.Status(400)
		return c.JSON(fiber.Map{
			"message":"incorrect password",
		})
	}

	token,err:=util.GenerateJwt(strconv.Itoa(int(user.Id)),)
	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	cookie := fiber.Cookie{
		Name:"jwt",
		Value:token,
		Expires: time.Now().Add(time.Hour *24),
		HTTPOnly: true,
	}		

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message":"you have successfully login",
		"user":user,
	})
}

type Claims struct{
	jwt.StandardClaims
}
