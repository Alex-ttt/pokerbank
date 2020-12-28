package handlers

import (
	"github.com/Alex-ttt/pokerbank/repository"
	"github.com/Alex-ttt/pokerbank/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func LoginPage(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", nil)
}

func Login(c *gin.Context) {
	var userCredentials CredentialsDto
	if err := c.ShouldBindJSON(&userCredentials); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	userCredentials.Login = strings.ToLower(strings.TrimSpace(userCredentials.Login))
	if doesLoginExist := repository.CheckLoginExists(services.Db, userCredentials.Login); !doesLoginExist {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	savedPassword := repository.GetPassword(services.Db, userCredentials.Login)
	if len(savedPassword) == 0 {
		// Register password for user
		if len(userCredentials.Password) < 5 {
			c.JSON(http.StatusBadRequest, "Short password is not required")
			return
		}

		newPassword, err := services.EncryptPassword(userCredentials.Password)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}

		repository.SetPassword(services.Db, userCredentials.Login, string(newPassword))
	} else {
		if !services.IsPasswordsEqual(savedPassword, userCredentials.Password) {
			c.JSON(http.StatusUnauthorized, "Please provide valid login details")
			return
		}
	}

	tokenDetails, err := services.CreateToken(userCredentials.Login)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := services.CreateAuth(userCredentials.Login, tokenDetails)
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
	accessTokenUid, err := services.ExtractTokenMetadata(c.Request, true)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	refreshTokenUid, err := services.ExtractTokenMetadata(c.Request, false)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleted, delErr := services.DeleteAuth(accessTokenUid.Uuid)
	if delErr != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleted, delErr = services.DeleteAuth(refreshTokenUid.Uuid)
	if delErr != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.Redirect(http.StatusSeeOther, repository.LoginRoute)
}

type CredentialsDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
