package api

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"
	"github.com/lzh-1625/go_process_manager/log"

	"github.com/gin-gonic/gin"
)

func rOk(ctx *gin.Context, message string, data any) {
	jsonData := map[string]any{
		"code": 0,
		"msg":  message,
	}
	if data != nil {
		jsonData["data"] = data
	}
	ctx.JSON(http.StatusOK, jsonData)
}

func errCheck(ctx *gin.Context, isErr bool, errData any) {
	if !isErr {
		return
	}
	if err, ok := errData.(error); ok {
		log.Logger.Warn(errData)
		ctx.Set(constants.CTXFLG_ERR, err.Error())
	}
	if err, ok := errData.(string); ok {
		ctx.Set(constants.CTXFLG_ERR, err)
	}
	panic(0)
}

func getRole(ctx *gin.Context) constants.Role {
	if v, ok := ctx.Get(constants.CTXFLG_ROLE); ok {
		return v.(constants.Role)
	}
	return constants.ROLE_GUEST
}

func getUserName(ctx *gin.Context) string {
	return ctx.GetString(constants.CTXFLG_USER_NAME)
}

func isAdmin(ctx *gin.Context) bool {
	return getRole(ctx) <= constants.ROLE_ADMIN
}

func hasOprPermission(ctx *gin.Context, uuid int, op constants.OprPermission) bool {
	return isAdmin(ctx) || reflect.ValueOf(repository.PermissionRepository.GetPermission(getUserName(ctx), uuid)).FieldByName(string(op)).Bool()
}

func getQueryInt(ctx *gin.Context, query string) (i int) {
	i, err := strconv.Atoi(ctx.Query(query))
	errCheck(ctx, err != nil, "Invalid parameters!")
	return
}

func getQueryString(ctx *gin.Context, query string) (s string) {
	s = ctx.Query(query)
	errCheck(ctx, s == "", "Invalid parameters!")
	return
}

func bind[T any](ctx *gin.Context) T {
	var data T
	errCheck(ctx, ctx.ShouldBind(&data) != nil, "Invalid parameters!")
	return data
}
