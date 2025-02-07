package api

import (
	"slices"

	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"
	"github.com/lzh-1625/go_process_manager/internal/app/service"

	"github.com/gin-gonic/gin"
)

type logApi struct{}

var LogApi = new(logApi)

func (a *logApi) GetLog(ctx *gin.Context, req model.GetLogReq) {
	if isAdmin(ctx) {
		rOk(ctx, "Query successful!", service.LogServiceImpl.Search(req, req.FilterName...))
	} else {
		processNameList := repository.PermissionRepository.GetProcessNameByPermission(getUserName(ctx), constants.OPERATION_LOG)
		filterName := slices.DeleteFunc(req.FilterName, func(s string) bool {
			return !slices.Contains(processNameList, s)
		})
		if len(filterName) == 0 {
			filterName = processNameList
		}
		errCheck(ctx, len(filterName) == 0, "No information found!")
		rOk(ctx, "Query successful!", service.LogServiceImpl.Search(req, filterName...))
	}
}

func (a *logApi) GetRunningLog(ctx *gin.Context) {
	rOk(ctx, "Query successful!", service.Loghandler.GetRunning())
}
