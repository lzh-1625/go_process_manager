package middle

import (
	"github.com/lzh-1625/go_process_manager/internal/app/constants"

	"github.com/gin-gonic/gin"
)

func PanicMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err == 0 {
				if err, ok := c.Get(constants.CTXFLG_ERR); ok {
					rErr(c, -1, err.(string), nil)
				} else {
					rErr(c, -1, "Internal error!", nil)
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
