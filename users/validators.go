package users

import (
	"apous-films-rest-api/utils"
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --------------- Register ----------------
type RegisterValidator struct {
	UserRegister struct {
		Email                string `json:"email" binding:"required"`
		Password             string `json:"password" binding:"required"`
		PasswordConfirmation string `json:"passwordConfirmation" binding:"required"`
	} `json:"user"`

	userModel User `json:"-"`
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

	// Map
	v.userModel.ID = primitive.NewObjectID()
	v.userModel.Email = v.UserRegister.Email
	v.userModel.PasswordHash = utils.HashPassword(v.UserRegister.Password)
	v.userModel.Provider = "local"

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
