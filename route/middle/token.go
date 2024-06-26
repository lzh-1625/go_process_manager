package middle

import (
	"errors"
	"msm/consts/ctxflag"
	"msm/dao"
	"msm/log"
	"msm/utils"
	"slices"

	"github.com/gin-gonic/gin"
)

// code -1为失败,-2为token失效
func rErr(ctx *gin.Context, code int, message string, err error) {
	var statusCode int
	switch code {
	case -1:
		statusCode = 500
	case -2:
		statusCode = 401
	default:
		statusCode = 200
	}
	log.Logger.Warn(err)
	ctx.JSON(statusCode, map[string]any{
		"code": code,
		"msg":  message,
	})
	ctx.Abort()
}

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		whiteList := []string{
			"/api/user/login",
			"/api/user/register/admin",
		}
		if !slices.Contains(whiteList, c.Request.URL.Path) {
			var token string
			if c.Request.Header.Get("token") != "" {
				token = c.Request.Header.Get("token")
			} else {
				token = c.Query("token")
			}
			if _, err := utils.ParseToken(token); err != nil {
				rErr(c, -2, "token校验失败", err)
				return
			}
			if username, err := getUser(c); err != nil {
				rErr(c, -1, "无法获取user信息", err)
			} else {
				c.Set(ctxflag.USER_NAME, username)
				c.Set(ctxflag.ROLE, dao.UserDao.GetUserByName(username).Role)
			}
		}
		c.Next()
	}
}

func getUser(ctx *gin.Context) (string, error) {
	var token string
	if ctx.Request.Header.Get("token") != "" {
		token = ctx.Request.Header.Get("token")
	} else {
		token = ctx.Query("token")
	}
	if mc, err := utils.ParseToken(token); err == nil && mc != nil {
		return mc.UserName, nil
	} else {
		return "", errors.Join(errors.New("用户信息获取失败"), err)
	}
}
