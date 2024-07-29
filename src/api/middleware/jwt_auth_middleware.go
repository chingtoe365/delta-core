package middleware

import (
	"net/http"
	"strings"

	"delta-core/bootstrap"
	"delta-core/domain"
	"delta-core/internal/tokenutil"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(env *bootstrap.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, err := tokenutil.IsAuthorized(authToken, env.AccessTokenSecret)
			if authorized {
				userID, err := tokenutil.ExtractIDFromToken(authToken, env.AccessTokenSecret)
				if err != nil {
					c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
					c.Abort()
					return
				}
				c.Set("x-user-id", userID)
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
			c.Abort()
			return
		}
		if env.AppEnv == "development" {
			c.Set("x-user-id", env.AnonUserId)
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Not authorized"})
		c.Abort()
	}
}
