package api

import (
	"github.com/lzh-1625/go_process_manager/config"
	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"
	"github.com/lzh-1625/go_process_manager/utils"

	"github.com/gin-gonic/gin"
)

type userApi struct{}

var UserApi = new(userApi)

const DEFAULT_ROOT_PASSWORD = "root"

func (u *userApi) LoginHandler(ctx *gin.Context, req map[string]string) {
	account := req["account"]
	password := req["password"]
	errCheck(ctx, !u.checkLoginInfo(account, password), "Incorrect username or password!")
	token, err := utils.GenToken(account)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", gin.H{
		"token":    token,
		"username": account,
		"role":     repository.UserRepository.GetUserByName(account).Role,
	})
}

func (u *userApi) CreateUser(ctx *gin.Context, req model.User) {
	errCheck(ctx, req.Role == constants.ROLE_ROOT, "Creation of root accounts is forbidden!")
	errCheck(ctx, req.Account == constants.CONSOLE, "Operation failed!")
	errCheck(ctx, len(req.Password) < config.CF.UserPassWordMinLength, "Password is too short")
	err := repository.UserRepository.CreateUser(req)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (u *userApi) ChangePassword(ctx *gin.Context, req model.User) {
	reqUser := getUserName(ctx)
	errCheck(ctx, getRole(ctx) != constants.ROLE_ROOT && req.Account != "", "Invalid parameters!")
	var userName string
	if req.Account != "" {
		userName = req.Account
	} else {
		userName = reqUser
	}
	errCheck(ctx, len(req.Password) < config.CF.UserPassWordMinLength, "Password is too short")
	err := repository.UserRepository.UpdatePassword(userName, req.Password)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)

}

func (u *userApi) DeleteUser(ctx *gin.Context) {
	account := getQueryString(ctx, "account")
	errCheck(ctx, account == "root", "Deletion of root accounts is forbidden!")
	err := repository.UserRepository.DeleteUser(account)
	errCheck(ctx, err != nil, "Deletion of root accounts failed!")
	rOk(ctx, "Operation successful!", nil)
}

func (u *userApi) GetUserList(ctx *gin.Context) {
	rOk(ctx, "Query successful!", repository.UserRepository.GetUserList())
}

func (u *userApi) checkLoginInfo(account, password string) bool {
	user := repository.UserRepository.GetUserByName(account)
	if account == "root" && user.Account == "" { // 如果root用户不存在，则创建一个root用户
		repository.UserRepository.CreateUser(model.User{
			Account:  "root",
			Password: DEFAULT_ROOT_PASSWORD,
			Role:     constants.ROLE_ROOT,
		})
		return password == DEFAULT_ROOT_PASSWORD
	}
	return user.Password == utils.Md5(password)
}
