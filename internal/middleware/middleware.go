package middleware

import (
	"book-management/internal/route"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewMyMiddleware, NewMiddlewares)

type VerifyToken interface {
	VerifyToken(token string) bool
}

type MyMiddleware struct {
	verifyToken VerifyToken
}

func NewMyMiddleware(verifyToken VerifyToken) *MyMiddleware {
	return &MyMiddleware{verifyToken: verifyToken}
}

func NewMiddlewares(m *MyMiddleware) []route.MiddleWare {
	var middlewares = []route.MiddleWare{
		route.MiddleWare(AuthMiddleware(m)),
	}
	return middlewares
}

func AuthMiddleware(m *MyMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		switch path {
		case "/api/v1/auth/get_verification_code", "/api/v1/auth/login":
			c.Next()
			return
		}
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		if !m.verifyToken.VerifyToken(token) {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
