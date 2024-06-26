package middle

import (
	"msm/consts/ctxflag"
	"msm/consts/permission"
	"msm/consts/role"
	"msm/dao"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RolePermission(needPermission role.Role) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if r := ctx.GetInt(ctxflag.ROLE); r > int(needPermission) {
			rErr(ctx, -1, "角色权限不足", nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func OprPermission(op permission.OprPermission) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		uuid, err := strconv.Atoi(ctx.Query("uuid"))
		if err != nil {
			rErr(ctx, -1, "参数有误", nil)
			ctx.Abort()
			return
		}
		if ctx.GetInt(ctxflag.ROLE) < int(role.USER) {
			ctx.Next()
			return
		}
		if !reflect.ValueOf(dao.PermissionDao.GetPermission(ctx.GetString(ctxflag.USER_NAME), uuid)).FieldByName(string(op)).Bool() {
			rErr(ctx, -1, "操作权限不足", nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
