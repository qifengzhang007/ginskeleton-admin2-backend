package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/users"
	"goskeleton/app/service/users/curd"
	"goskeleton/app/service/users/login_policy"
	userstoken "goskeleton/app/service/users/token"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/response"
	"time"
)

type Users struct {
}

// 1.用户注册
func (u *Users) Register(context *gin.Context) {
	//  由于本项目骨架已经将表单验证器的字段(成员)绑定在上下文，因此可以按照 GetString()、GetBool()、GetFloat64（）等快捷获取需要的数据类型，注意：相关键名规则：  前缀+验证器结构体中的 json 标签
	// 注意：在 ginskeleton 中获取表单参数验证器中的数字键（字段）,请统一使用 GetFloat64(),其它获取数字键（字段）的函数无效，例如：GetInt()、GetInt64()等
	// 当然也可以通过gin框架的上下文原始方法获取，例如： context.PostForm("user_name") 获取，这样获取的数据格式为文本，需要自己继续转换
	userName := context.GetString(consts.ValidatorPrefix + "user_name")
	pass := context.GetString(consts.ValidatorPrefix + "pass")
	userIp := context.ClientIP()
	if curd.CreateUserCurdFactory().Register(userName, pass, userIp) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdRegisterFailCode, consts.CurdRegisterFailMsg, "")
	}
}

// 2.用户登录
func (u *Users) Login(context *gin.Context) {
	userName := context.GetString(consts.ValidatorPrefix + "user_name")
	pass := context.GetString(consts.ValidatorPrefix + "pass")
	phone := context.GetString(consts.ValidatorPrefix + "phone")

	// 登陆安全策略是否开启
	loginPolicyIsOpen := variable.ConfigYml.GetInt64("LoginPolicy.IsOpenPolicy")
	// 登陆之前检查账号是否被登陆策略禁用
	var lpf *login_policy.UserLoginPolicy
	if loginPolicyIsOpen == 1 {
		lpf = login_policy.CreateUsersLoginPolicyFactory()
		if lpf != nil {
			defer lpf.ReleaseRedisConn()
			if forbidden, err := lpf.CheckAccountIsForbidden(userName); err == nil && forbidden {
				response.Fail(context, consts.CurdLoginFailCode, "该账号已被登陆策略自动禁止登陆, 请至少等待2分钟后再试", "")
				return
			}
		}
	}

	userModel := users.CreateUserFactory("").Login(userName, pass)
	if userModel != nil {
		// 成功登陆调用登陆策略，自动会将之前的失败次数清0
		if loginPolicyIsOpen == 1 && lpf != nil {
			_, _ = lpf.SetAccountLoginCache(userName, false)
		}

		userTokenFactory := userstoken.CreateUserFactory()
		if userToken, err := userTokenFactory.GenerateToken(userModel.Id, userModel.UserName, userModel.Phone, variable.ConfigYml.GetInt64("Token.JwtTokenCreatedExpireAt")); err == nil {
			if userTokenFactory.RecordLoginToken(userToken, context.ClientIP()) {
				data := gin.H{
					"id":         userModel.Id,
					"user_name":  userName,
					"real_name":  userModel.RealName,
					"phone":      phone,
					"token":      userToken,
					"updated_at": time.Now().Format(variable.DateFormat),
				}
				response.Success(context, consts.CurdStatusOkMsg, data)
				return
			}
		} else {
			fmt.Println("生成token出错：", err.Error())
		}
	}
	if loginPolicyIsOpen == 1 && lpf != nil {
		// 登陆失败 - 调用登陆策略，自动会将之前的失败次数+1
		_, _ = lpf.SetAccountLoginCache(userName, true)
	}

	response.Fail(context, consts.CurdLoginFailCode, "账号或密码错误", "")
}

// 刷新用户token
func (u *Users) RefreshToken(context *gin.Context) {
	oldToken := context.GetString(consts.ValidatorPrefix + "token")
	if newToken, ok := userstoken.CreateUserFactory().RefreshToken(oldToken, context.ClientIP()); ok {
		res := gin.H{
			"token": newToken,
		}
		response.Success(context, consts.CurdStatusOkMsg, res)
	} else {
		response.Fail(context, consts.CurdRefreshTokenFailCode, consts.CurdRefreshTokenFailMsg, "")
	}
}

// 3.用户查询（show）
func (u *Users) List(context *gin.Context) {
	userName := context.GetString(consts.ValidatorPrefix + "user_name")
	page := context.GetFloat64(consts.ValidatorPrefix + "page")
	limit := context.GetFloat64(consts.ValidatorPrefix + "limit")
	limitStart := (page - 1) * limit

	totalCounts, showList := users.CreateUserFactory("").List(userName, int(limitStart), int(limit))
	if totalCounts > 0 && showList != nil {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"count": totalCounts, "data": showList})
	} else {
		response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

// 3.用户查询（PostList），根据部门、岗位id，用户名关键词查询数据
func (u *Users) PostList(context *gin.Context) {
	orgPostName := context.GetString(consts.ValidatorPrefix + "org_post_name")
	userName := context.GetString(consts.ValidatorPrefix + "user_name")
	page := context.GetFloat64(consts.ValidatorPrefix + "page")
	limit := context.GetFloat64(consts.ValidatorPrefix + "limit")
	limitStart := (page - 1) * limit

	totalCounts, postList := users.CreateUserFactory("").PostList(userName, orgPostName, int(limitStart), int(limit))
	if totalCounts > 0 && postList != nil {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"count": totalCounts, "data": postList})
	} else {
		response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

// 4.用户新增(store)
func (u *Users) Create(context *gin.Context) {
	if users.CreateUserFactory("").InsertData(context) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+",用户名不能重复", "")
	}
}

// 5.用户更新(update)
func (u *Users) Edit(context *gin.Context) {
	userId := context.GetFloat64(consts.ValidatorPrefix + "id")
	userName := context.GetString(consts.ValidatorPrefix + "user_name")
	// 检查正在修改的用户名是否被其他人使用
	if users.CreateUserFactory("").UpdateDataCheckUserNameIsUsed(int(userId), userName) > 0 {
		response.Fail(context, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg+", "+userName+" 已经被其他人使用", "")
		return
	}

	//注意：这里没有实现权限控制逻辑，例如：超级管理管理员可以更新全部用户数据，普通用户只能修改自己的数据。目前只是验证了token有效、合法之后就可以进行后续操作
	// 实际使用请根据真是业务实现权限控制逻辑、再进行数据库操作
	if users.CreateUserFactory("").UpdateData(context) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
	}

}

// 6.删除记录
func (u *Users) Destroy(context *gin.Context) {
	userId := context.GetFloat64(consts.ValidatorPrefix + "id")
	if users.CreateUserFactory("").DeleteData(int(userId)) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg, "注意：admin 不能删除")
	}
}

// 6.获取用户token信息+动态菜单
func (u *Users) UserInfo(context *gin.Context) {
	UserId, exist := cur_userinfo.GetCurrentUserId(context)
	if !exist {
		response.Fail(context, consts.CurdTokenFailCode, consts.CurdTokenFailMsg, "")
	} else {
		userService := curd.CreateUserCurdFactory()
		if data := userService.FindUserInfo(UserId); data != nil {
			response.Success(context, consts.CurdStatusOkMsg, data)
		} else {
			response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
		}
	}
}

// 查询用户当前打开的页面允许显示的按钮（查询指定页面拥有的按钮权限）
func (u *Users) GetButtonListByMenuId(context *gin.Context) {
	menuId := context.GetFloat64(consts.ValidatorPrefix + "menu_id")
	UserId, exist := cur_userinfo.GetCurrentUserId(context)
	if !exist {
		response.Fail(context, consts.CurdTokenFailCode, consts.CurdTokenFailMsg, "")
	} else {
		data := curd.CreateUserCurdFactory().GetButtonListByMenuId(UserId, int64(menuId))
		response.Success(context, consts.CurdStatusOkMsg, data)
	}

}

// GetPersonalInfo 每个用户查询自己的个人信息
func (u *Users) GetPersonalInfo(context *gin.Context) {
	userId, exists := cur_userinfo.GetCurrentUserId(context)
	if !exists {
		response.Fail(context, consts.CurdTokenFailCode, consts.CurdTokenFailMsg, "")
	} else {
		user, _ := users.CreateUserFactory("").ShowOneItem(userId)
		response.Success(context, consts.CurdStatusOkMsg, user)
	}
}

// EditPersonalInfo 编辑自己的信息
func (u *Users) EditPersonalInfo(context *gin.Context) {
	// 获取当前请求用户id
	userId, exists := cur_userinfo.GetCurrentUserId(context)
	if !exists {
		response.Fail(context, consts.CurdTokenFailCode, consts.CurdTokenFailMsg, "")
	} else {

		userName := context.GetString(consts.ValidatorPrefix + "user_name")
		usersModel := users.CreateUserFactory("")

		// 检查正在修改的用户名是否被其他站使用
		if usersModel.UpdateDataCheckUserNameIsUsed(int(userId), userName) > 0 {
			response.Fail(context, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg+",该用户名: "+userName+" 已经被其他人占用", "")
			return
		}
		// 这里使用token解析的id更新表单参数里面的id，加固安全
		context.Set(consts.ValidatorPrefix+"id", float64(userId))

		if usersModel.UpdateData(context) {
			response.Success(context, consts.CurdStatusOkMsg, "")
		} else {
			response.Fail(context, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
		}
	}
}
