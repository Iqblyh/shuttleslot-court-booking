package AuthMiddlewareMock

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type AuthMiddlewareMock struct {
	mock.Mock
}

func (a *AuthMiddlewareMock) CheckToken(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {}

}
