package api

import (
	"msm/consts/ctxflag"
	"msm/consts/role"
	"msm/dao"
	"msm/model"
	"msm/utils"

	"github.com/gin-gonic/gin"
)

type userApi struct{}

var UserApi = new(userApi)

const DEFAULT_ROOT_PASSWORD = "root"

func (u *userApi) LoginHandler(ctx *gin.Context) {
	info := map[string]string{}
	ctx.ShouldBindJSON(&info)
	account := info["account"]
	password := info["password"]
	errCheck(ctx, !u.checkLoginInfo(account, password), "登入失败,账号或密码错误")
	token, err := utils.GenToken(account)
	errCheck(ctx, err != nil, err)
	ctx.JSON(200, gin.H{
		"code":     0,
		"msg":      "登入成功！",
		"token":    token,
		"username": account,
		"role":     dao.UserDao.GetUserByName(account).Role,
	})
}

func (u *userApi) CreateUser(ctx *gin.Context) {
	user := model.User{}
	err := ctx.ShouldBindJSON(&user)
	errCheck(ctx, err != nil, err)
	errCheck(ctx, user.Role == int(role.ROOT), "不能添加root账号")
	err = dao.UserDao.CreateUser(user)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "注册成功", nil)
}

func (u *userApi) ChangePassword(ctx *gin.Context) {
	user := model.User{}
	err := ctx.ShouldBindJSON(&user)
	errCheck(ctx, err != nil, err)
	reqUser := ctx.GetString(ctxflag.USER_NAME)
	errCheck(ctx, ctx.GetInt(ctxflag.ROLE) != int(role.ROOT) && user.Account != "", "参数错误")
	var userName string
	if user.Account != "" {
		userName = user.Account
	} else {
		userName = reqUser
	}
	err = dao.UserDao.UpdatePassword(userName, user.Password)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "修改密码成功", nil)

}

func (u *userApi) DeleteUser(ctx *gin.Context) {
	errCheck(ctx, ctx.Query("account") == "root", "无法删除root账户")
	err := dao.UserDao.DeleteUser(ctx.Query("account"))
	errCheck(ctx, err != nil, "无法删除root账户")
	rOk(ctx, "删除成功", nil)
}

func (u *userApi) GetUserList(ctx *gin.Context) {
	rOk(ctx, "查询成功", dao.UserDao.GetUserList())
}

func (u *userApi) checkLoginInfo(account, password string) bool {
	user := dao.UserDao.GetUserByName(account)
	if account == "root" && user.Account == "" {
		dao.UserDao.CreateUser(model.User{
			Account:  "root",
			Password: DEFAULT_ROOT_PASSWORD,
			Role:     int(role.ROOT),
		})
		return password == DEFAULT_ROOT_PASSWORD
	}
	return user.Password == utils.Md5(password)
}
