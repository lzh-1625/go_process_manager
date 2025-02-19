package middle

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

// 演示模式
func DemoMiddle() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		whiteListUri := []string{
			"/api/user/login",
			"/api/log",
		}
		if ctx.Request.Method == http.MethodGet || slices.Contains(whiteListUri, ctx.Request.URL.String()) {
			ctx.Next()
		} else {
			ctx.Abort()
		}
	}
}
