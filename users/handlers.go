package users

import (
	"apous-films-rest-api/config"
	"apous-films-rest-api/oauth"
	"apous-films-rest-api/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterHandlers(c *gin.RouterGroup, s Service) {
	c.POST("/register", UserRegister(s))
	c.POST("/login", UserLogin(s))
	c.GET("/google/callback", GoogleLogin(s))
	c.GET("/github/callback", GitHubLogin(s))
	c.GET("/profile", UserProfile(s))
}

func UserRegister(s Service) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Bind and validate
		v := RegisterValidator{}

		if err := v.BindAndValidate(c); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Register
		token, err := s.Register(v.UserRegister.Email, v.UserRegister.Password)

		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, token)
	}

	return gin.HandlerFunc(fn)
}

func UserLogin(s Service) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

	return gin.HandlerFunc(fn)
}

func GoogleLogin(s Service) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

		c.Redirect(http.StatusPermanentRedirect, config.Config.Client.URI)
	}

	return gin.HandlerFunc(fn)
}

func GitHubLogin(s Service) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		code := c.Query("code")

		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
			return
		}

		// Get google user
		var gitHubUser oauth.GitHubUser
		gp := oauth.NewGitHubProvider()

		if err := gp.GetGitHubUser(code, &gitHubUser); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to perform google authentication"})
			return
		}

		user, err := FindUserByEmail(gitHubUser.Email)

		// If user does not exists, insert the new user
		if err == mongo.ErrNoDocuments {
			user.ID = primitive.NewObjectID()
			user.Email = gitHubUser.Email
			user.Provider = "github"
			user.PasswordHash = ""

			if err := CreateUser(&user); err != nil {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
		}

		// Generate token
		token := utils.GenerateJWT(user.ID.Hex())
		utils.SetCookieToken(c, token)

		c.Redirect(http.StatusPermanentRedirect, config.Config.Client.URI)
	}

	return gin.HandlerFunc(fn)
}

func UserProfile(s Service) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		serializer := UserSerializer{c}
		c.JSON(http.StatusOK, serializer.Response())
	}

	return gin.HandlerFunc(fn)
}
