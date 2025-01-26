package api

import (
	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"
	"github.com/lzh-1625/go_process_manager/internal/app/service"

	"github.com/gin-gonic/gin"
)

type logApi struct{}

var LogApi = new(logApi)

func (a *logApi) GetLog(ctx *gin.Context, req model.GetLogReq) {
	filterName := make([]string, 0, len(req.FilterName))
	processNameList := repository.PermissionRepository.GetProcessNameByPermission(getUserName(ctx), constants.OPERATION_LOG)
	if len(filterName) != 0 {
		for _, v := range processNameList {
			for _, m := range req.FilterName {
				if v == m {
					filterName = append(filterName, m)
					break
				}
			}
		}
	} else {
		filterName = append(filterName, processNameList...)
	}
	errCheck(ctx, !isAdmin(ctx) && len(filterName) == 0, "No information found!")
	rOk(ctx, "Query successful!", service.LogServiceImpl.Search(req, filterName...))
}

func (a *logApi) GetRunningLog(ctx *gin.Context) {
	rOk(ctx, "Query successful!", service.Loghandler.GetRunning())
}
