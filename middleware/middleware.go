package middleware

import (
	"net/http"

	"github.com/aca/permit/utils"
	"github.com/gin-gonic/gin"
)

func MainMiddleware(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	if _, err := utils.ExtractTokenID(token); err != nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	c.Next()
}