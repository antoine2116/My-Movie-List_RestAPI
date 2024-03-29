package users

import (
	"github.com/gin-gonic/gin"
)

// --------------- Register ----------------
type RegisterValidator struct {
	UserRegister struct {
		Email                string `json:"email" binding:"required"`
		Password             string `json:"password" binding:"required"`
		PasswordConfirmation string `json:"passwordConfirmation" binding:"required"`
	} `json:"user"`
}

func (v *RegisterValidator) Bind(c *gin.Context) error {
	// Bind
	if err := c.ShouldBindJSON(v); err != nil {
		return err
	}

	// Validate
	// Check password confirmation
	if v.UserRegister.Password != v.UserRegister.PasswordConfirmation {
		return ErrMismatchedPasswords
	}

	return nil
}

// --------------- Login ----------------
type LoginValidator struct {
	UserLogin struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	} `json:"user"`
}

func (v *LoginValidator) Bind(c *gin.Context) error {
	// Bind
	if err := c.ShouldBindJSON(v); err != nil {
		return err
	}

	// No validations yet

	return nil
}
