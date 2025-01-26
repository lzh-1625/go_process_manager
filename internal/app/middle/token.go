package middle

import (
	"errors"
	"slices"
	"strings"

	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"
	"github.com/lzh-1625/go_process_manager/log"
	"github.com/lzh-1625/go_process_manager/utils"

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
	if err != nil {
		log.Logger.Warn(err)
	}
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
			"/api/task/api-key/",
		}
		if !slices.ContainsFunc(whiteList, func(s string) bool {
			return strings.HasPrefix(c.Request.URL.Path, s)
		}) {
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
				c.Set(constants.CTXFLG_USER_NAME, username)
				c.Set(constants.CTXFLG_ROLE, repository.UserRepository.GetUserByName(username).Role)
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
