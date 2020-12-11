package handlers

import (
	"github.com/Alex-ttt/pokerbank/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined:
	//if user.Username != u.Username || user.Password != u.Password {
	//	c.JSON(http.StatusUnauthorized, "Please provide valid login details")
	//	return
	//}
	ts, err := services.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := services.CreateAuth(user.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     services.AccessTokenKey,
		Value:    url.QueryEscape(ts.AccessToken),
		HttpOnly: true,
		Secure:   false,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     services.RefreshTokenKey,
		Value:    url.QueryEscape(ts.RefreshToken),
		HttpOnly: true,
		Secure:   false,
	})
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
