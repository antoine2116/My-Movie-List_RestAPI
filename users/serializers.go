package users

import (
	"apous-films-rest-api/utils"

	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type UserSerializer struct {
	c *gin.Context
}

func (serializer *UserSerializer) Response() UserResponse {
	userModel := serializer.c.MustGet("user").(User)

	return UserResponse{
		Email: userModel.Email,
		Token: utils.GenerateJWT(userModel.ID.Hex()),
	}
}
