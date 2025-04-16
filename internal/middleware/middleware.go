package middleware

import (
	"book-management/internal/route"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
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
		route.MiddleWare(Cors()),
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin) // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
