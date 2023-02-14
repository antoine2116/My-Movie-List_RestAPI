package handlers

import (
	"apous-films-rest-api/models"
	"apous-films-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRegestration(c *gin.Context) {
	var regUser *models.RegisterUser

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
	userId, err := models.CreateUser(regUser)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	userResp := models.UserResponse{
		Username: regUser.Username,
		Email:    regUser.Email,
		Token:    utils.GenerateJWT(userId),
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": userResp})
}

func UserLogin(c *gin.Context) {
	var login *models.LoginUser

	// Model validation
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Find user
	user, err := models.FindUserByEmail(login.Email)

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

	userResp := models.UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Token:    utils.GenerateJWT(user.ID.Hex()),
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": userResp})
}
