package auth

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const prefix = "Bearer "

func JWTMiddleware(verifier *oidc.IDTokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, prefix) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		rawToken := authHeader[len(prefix):]
		_, err := verifier.Verify(c.Request.Context(), rawToken)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
	}
}
