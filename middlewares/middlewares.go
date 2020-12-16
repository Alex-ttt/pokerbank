package middlewares

import (
	"github.com/Alex-ttt/pokerbank/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenAuthWithRedirectToIndexMiddleware(c *gin.Context) {
	if isAuthorized, _ := services.IsRequestAuthorized(c.Request); !isAuthorized {
		if isRefreshSucceed, _ := services.Refresh(c); !isRefreshSucceed {
			c.Next()
			return
		}
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func TokenAuthWithRedirectToLoginMiddleware(c *gin.Context) {
	if isAuthorized, _ := services.IsRequestAuthorized(c.Request); !isAuthorized {
		if isRefreshSucceed, _ := services.Refresh(c); !isRefreshSucceed {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
	}

	c.Next()
}

func TokenAuthMiddleware(c *gin.Context) {
	if isAuthorized, _ := services.IsRequestAuthorized(c.Request); !isAuthorized {
		if isRefreshSucceed, _ := services.Refresh(c); !isRefreshSucceed {
			c.JSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			c.Abort()
			return
		}
	}

	c.Next()
}
