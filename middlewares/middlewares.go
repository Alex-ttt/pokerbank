package middlewares

import (
	"github.com/Alex-ttt/pokerbank/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
