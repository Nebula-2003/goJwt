package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Nebula-2003/goJwt/initializers"
	"github.com/Nebula-2003/goJwt/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var body struct {
		Email    string
		Password string
		Name     string
		Role     string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse body",
		})
		return
	}

	byteArrHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error hashing password",
		})
		return
	}

	user := models.User{Name: body.Name, Email: body.Email, Password: string(byteArrHash)}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		fmt.Print(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating user",
		})
		return
	}

	// Print inserted data
	fmt.Printf("Inserted user: %+v\n", user)

	// Return the created user (excluding the password)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse body",
		})
		return
	}
	var user models.User
	initializers.DB.Find(&user, "email = ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error finding User",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Incorrect Password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"role": "admin",
	})
	mySigningKey := []byte("AllYourBase")

	ss, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Creating Auth token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": user,
		"jwt":  ss,
	})
}
