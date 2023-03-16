package users

import (
	"apous-films-rest-api/utils"
	"errors"

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

func (v *RegisterValidator) BindAndValidate(c *gin.Context) error {
	// Bind
	if err := c.ShouldBindJSON(v); err != nil {
		return utils.NewValidationError(err)
	}

	// Validate
	// Check password confirmation
	if v.UserRegister.Password != v.UserRegister.PasswordConfirmation {
		return errors.New("passwords do not match")
	}

	return nil
}

// --------------- Login ----------------
type LoginValidator struct {
	UserLogin struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	} `json:"user"`

	userModel User `json:"-"`
}

func (v *LoginValidator) BindAndValidate(c *gin.Context) error {
	// Bind
	if err := c.ShouldBind(v); err != nil {
		return utils.NewValidationError(err)
	}

	// No validations yet

	// Map
	v.userModel.Email = v.UserLogin.Email

	return nil
}
