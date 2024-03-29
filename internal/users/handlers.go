package users

import (
	"mml-rest-api/internal/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	cookie_maxAge = 60 * 60 * 24
)

func RegisterHandlers(c *gin.RouterGroup, s Service, clientURI string) {
	c.POST("/register", UserRegister(s))
	c.POST("/login", UserLogin(s))
	c.GET("/google/callback", GoogleLogin(s, clientURI))
	c.GET("/github/callback", GitHubLogin(s, clientURI))
}

func RegisterAuthenticatedHandlers(c *gin.RouterGroup) {
	c.GET("/profile", UserProfile())
}

func UserRegister(s Service) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		// Bind and validate
		v := RegisterValidator{}

		if err := v.Bind(ctx); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.NewValidationError(err)})
			return
		}

		// Register
		token, err := s.Register(ctx, v.UserRegister.Email, v.UserRegister.Password)

		if err != nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"token": token})
	}

	return gin.HandlerFunc(fn)
}

func UserLogin(s Service) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		// Bind
		v := LoginValidator{}

		if err := v.Bind(ctx); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.NewValidationError(err)})
			return
		}

		// Authenticate
		token, err := s.Login(ctx, v.UserLogin.Email, v.UserLogin.Password)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.NewCommonError(err)})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"token": token})
	}

	return gin.HandlerFunc(fn)
}

func GoogleLogin(s Service, clientURI string) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		// Extract code from query
		code := ctx.Query("code")

		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.NewCommonError(ErrInvalidToken)})
			return
		}

		// Login
		token, err := s.GoogleLogin(ctx, code)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.NewCommonError(err)})
			return
		}

		// Set cookie
		ctx.SetCookie("token", token, cookie_maxAge, "/", "", true, false)

		// Redirect to client
		ctx.Redirect(http.StatusPermanentRedirect, clientURI)
	}

	return gin.HandlerFunc(fn)
}

func GitHubLogin(s Service, clientURI string) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		code := ctx.Query("code")

		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.NewCommonError(ErrInvalidToken)})
			return
		}

		// Login
		token, err := s.GitHubLogin(ctx, code)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.NewCommonError(err)})
			return
		}

		// Set cookie
		ctx.SetCookie("token", token, cookie_maxAge, "/", "", true, false)

		// Redirect to client
		ctx.Redirect(http.StatusPermanentRedirect, clientURI)
	}

	return gin.HandlerFunc(fn)
}

func UserProfile() gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		user := GetCurrentUser(ctx)
		ctx.JSON(http.StatusOK, user)
	}

	return gin.HandlerFunc(fn)
}
