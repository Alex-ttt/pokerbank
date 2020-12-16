package handlers

import (
	"github.com/Alex-ttt/pokerbank/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginPage(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", nil)
}

func Login(c *gin.Context) {
	var user UserLoginDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined:
	//if user.Username != u.Username || user.Password != u.Password {
	//	c.JSON(http.StatusUnauthorized, "Please provide valid login details")
	//	return
	//}
	tokenDetails, err := services.CreateToken(1)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := services.CreateAuth(1, tokenDetails)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	tokens := map[string]string{
		services.AccessTokenKey:  tokenDetails.AccessToken,
		services.RefreshTokenKey: tokenDetails.RefreshToken,
	}

	services.SetTokensToResponseCookie(&c.Writer, &tokens)
	c.JSON(http.StatusOK, tokens)
}

func Logout(c *gin.Context) {
	au, err := services.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleted, delErr := services.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

type UserLoginDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
