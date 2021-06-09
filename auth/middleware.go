package auth

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const prefix = "Bearer "

func JWTMiddleware(verifier *oidc.IDTokenVerifier, claims map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, prefix) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		rawToken := authHeader[len(prefix):]
		token, err := verifier.Verify(c.Request.Context(), rawToken)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		var actualClaims map[string]interface{}
		if err := token.Claims(&actualClaims); err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		for k, v := range claims {
			if actualClaims[k] != v {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
	}
}
