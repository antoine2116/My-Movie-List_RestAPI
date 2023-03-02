package users

import (
	"apous-films-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddRoutes(c *gin.RouterGroup) {
	c.POST("/register", UserRegister)
	c.POST("/login", UserLogin)
}

func UserRegister(c *gin.Context) {
	// Bind and validate
	validator := RegisterValidator{}

	if err := validator.BindAndValidate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert
	if err := CreateUser(&validator.userModel); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// Send response
	serializer := UserSerializer{c}
	c.Set("user_model", validator.userModel)

	c.JSON(http.StatusCreated, serializer.Response())
}

func UserLogin(c *gin.Context) {
	// Bind
	validator := LoginValidator{}

	if err := validator.BindAndValidate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	user, err := FindUserByEmail(validator.userModel.Email)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
			return
		}
	}

	if err := utils.CompareHashAndPassword(user.PasswordHash, validator.UserLogin.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	serializer := UserSerializer{c}
	c.Set("user_model", user)

	c.JSON(http.StatusOK, gin.H{"data": serializer.Response()})
}
