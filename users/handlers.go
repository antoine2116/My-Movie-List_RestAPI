package users

import (
	"apous-films-rest-api/oauth"
	"apous-films-rest-api/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddUserAuthentication(c *gin.RouterGroup) {
	c.POST("/register", UserRegister)
	c.POST("/login", UserLogin)
}

func AddUserProfile(c *gin.RouterGroup) {
	c.GET("/profile", UserProfile)
}

func AddGoogleOAuth(c *gin.RouterGroup) {
	c.GET("/google/callback", GoogleLogin)
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
	c.Set("user", validator.userModel)

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
	c.Set("user", user)

	c.JSON(http.StatusOK, serializer.Response())
}

func GoogleLogin(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
		return
	}

	// Get google user
	var googleUser oauth.GoogleUser
	gp := oauth.NewGoogleProvider()

	if err := gp.GetGoogleUser(code, &googleUser); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to perform google authentication"})
		return
	}

	user, err := FindUserByEmail(googleUser.Email)

	// If user does not exists, insert the new user
	if err == mongo.ErrNoDocuments {
		user.ID = primitive.NewObjectID()
		user.Email = googleUser.Email
		user.Provider = "google"
		user.PasswordHash = ""

		if err := CreateUser(&user); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
	}

	// Generate token
	token := utils.GenerateJWT(user.ID.Hex())
	utils.SetCookieToken(c, token)

	c.Redirect(http.StatusPermanentRedirect, "http://localhost:3000")
}

func UserProfile(c *gin.Context) {
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, serializer.Response())
}
