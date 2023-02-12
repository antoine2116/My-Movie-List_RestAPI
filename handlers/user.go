package handlers

import (
	"apous-films-rest-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRegestration(c *gin.Context) {
	regUser := models.RegisterUser{}

	// Model validation
	if err := c.ShouldBindJSON(&regUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// TODO : generate JWT Token

	// Insert into database
	if err := models.CreateUser(regUser); err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": "User created"})
}
