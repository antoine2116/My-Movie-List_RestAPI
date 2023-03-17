package users

import (
	"apous-films-rest-api/cookies"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(c *gin.RouterGroup, s Service, clientURI string) {
	c.POST("/register", UserRegister(s))
	c.POST("/login", UserLogin(s))
	c.GET("/google/callback", GoogleLogin(s, clientURI))
	c.GET("/github/callback", GitHubLogin(s, clientURI))
	c.GET("/profile", UserProfile(s))
}

func UserRegister(s Service) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		// Bind and validate
		v := RegisterValidator{}

		if err := v.BindAndValidate(ctx); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

		if err := v.BindAndValidate(ctx); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Register
		token, err := s.Login(ctx, v.UserLogin.Email, v.UserLogin.Password)

		if err != nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"token": token})
	}

	return gin.HandlerFunc(fn)
}

func GoogleLogin(s Service, clientURI string) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		// Extract code from query
		code := ctx.Query("code")

		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
			return
		}

		// Login
		token, err := s.GoogleLogin(ctx, code)

		if err != nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		// Set cookie
		cookies.SetToken(ctx, token)

		// Redirect to client
		ctx.Redirect(http.StatusPermanentRedirect, clientURI)
	}

	return gin.HandlerFunc(fn)
}

func GitHubLogin(s Service, clientURI string) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		code := ctx.Query("code")

		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
			return
		}

		// Login
		token, err := s.GitHubLogin(ctx, code)

		if err != nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		// Set cookie
		cookies.SetToken(ctx, token)

		// Redirect to client
		ctx.Redirect(http.StatusPermanentRedirect, clientURI)
	}

	return gin.HandlerFunc(fn)
}

func UserProfile(s Service) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		// Get the current user id
		userId, ok := ctx.Get("user_id")

		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		}

		// Get the user
		user, err := s.GetById(ctx, userId.(string))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, gin.H{"user": user})
	}

	return gin.HandlerFunc(fn)
}
