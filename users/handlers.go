package users

import (
	"apous-films-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRegestration(c *gin.Context) {
	var regUser *RegisterUser

	// Model validation
	if err := c.ShouldBindJSON(&regUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Check password confirmation
	if regUser.Password != regUser.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	// Insert
	user, err := CreateUser(regUser)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	serializer := UserSerializer{c}
	c.Set("user_model", user)

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": serializer.Response()})
}

func UserLogin(c *gin.Context) {
	var login *LoginUser

	// Model validation
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Find user
	user, err := FindUserByEmail(login.Email)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
			return
		}
	}

	if err := utils.CompareHashAndPassword(user.PasswordHash, login.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
		return
	}

	serializer := UserSerializer{c}
	c.Set("user_model", user)

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": serializer.Response()})
}
