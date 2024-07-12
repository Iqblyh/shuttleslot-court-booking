package middleware

import (
	"net/http"
	"strings"
	"team2/shuttleslot/model/dto"
	"team2/shuttleslot/service"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	CheckToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	service service.AuthService
}

// CheckToken implements AuthMiddleware.
func (a *authMiddleware) CheckToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		token := strings.Replace(authHeader, "Bearer ", "", -1)
		claims, err := a.service.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": dto.Status{
					Code: http.StatusUnauthorized,
					Message: "Unauthorized",
				},
			})
			return
		}
		ctx.Set("userId", claims["userId"])
		ctx.Set("role", claims["role"])
		var validRole bool
		for _, r := range roles {
			if r == claims["role"] {
				validRole = true
				break
			}
		}
		if !validRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status": dto.Status{
					Code: http.StatusForbidden,
					Message: "Forbidden Access",
				},
			})
		}
		ctx.Next()
	}
}

func NewAuthMiddleware(authService service.AuthService) AuthMiddleware {
	return &authMiddleware{service: authService}
}
