package middle

import (
	"msm/consts/ctxflag"

	"github.com/gin-gonic/gin"
)

func PanicMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err == 0 {
				if err, ok := c.Get(ctxflag.ERR); ok {
					rErr(c, -1, err.(error).Error(), err.(error))
				} else {
					rErr(c, -1, "内部错误", nil)
				}
			} else {
				if err != nil {
					panic(err)
				}
			}
		}()
		c.Next()
	}
}
