package middle

import (
	"reflect"
	"strconv"

	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func RolePermission(needPermission constants.Role) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if v, ok := ctx.Get(constants.CTXFLG_ROLE); !ok || v.(constants.Role) > needPermission {
			rErr(ctx, -1, "Insufficient permissions; please check your access rights!", nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func OprPermission(op constants.OprPermission) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		uuid, err := strconv.Atoi(ctx.Query("uuid"))
		if err != nil {
			rErr(ctx, -1, "Invalid parameters!", nil)
			ctx.Abort()
			return
		}
		if v, ok := ctx.Get(constants.CTXFLG_ROLE); !ok || v.(constants.Role) <= constants.ROLE_ADMIN {
			ctx.Next()
			return
		}
		if !reflect.ValueOf(repository.PermissionRepository.GetPermission(ctx.GetString(constants.CTXFLG_USER_NAME), uuid)).FieldByName(string(op)).Bool() {
			rErr(ctx, -1, "Insufficient permissions; please check your access rights!", nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
