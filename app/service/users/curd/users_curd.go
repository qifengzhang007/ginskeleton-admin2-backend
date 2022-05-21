package curd

import (
	"github.com/gin-gonic/gin"
	"github.com/qifengzhang007/sql_res_to_tree"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/middleware/my_jwt"
	modeAuth "goskeleton/app/model/auth"
	"goskeleton/app/model/users"
	"goskeleton/app/utils/md5_encrypt"
	"strconv"
	"strings"
)

func CreateUserCurdFactory() *UsersCurd {
	return &UsersCurd{users.CreateUserFactory("")}
}

type UsersCurd struct {
	userModel *users.UsersModel
}

func (u *UsersCurd) Register(userName, pass, userIp string) bool {
	pass = md5_encrypt.Base64Md5(pass) // 预先处理密码加密，然后存储在数据库
	return u.userModel.Register(userName, pass, userIp)
}

type userWithMenus struct {
	users.UsersModel
	Menus []modeAuth.AuthSystemMenuTree `json:"menus"`
}

//获取用户信息并处理业务逻辑
func (u *UsersCurd) FindUserInfo(userId int64) *userWithMenus {
	var data = &userWithMenus{}
	var user *users.UsersModel = nil
	user, _ = u.userModel.ShowOneItem(userId)
	data.UsersModel = *user
	orgIds := u.getUserAllOrgIds(user.Id)
	//根据岗位ID获取拥有的菜单ID,去重
	menuIdModel := modeAuth.CreateAuthPostMountHasMenuModelFactory("").GetByIds(orgIds)
	menuIdArray := []int{}
	for k, _ := range menuIdModel {
		menuIdArray = append(menuIdArray, menuIdModel[k].FrAuthSystemMenuId)
	}

	//根据菜单 Ids数组 获取菜单信息

	menus := modeAuth.CreateAuthSystemMenuFactory("").GetByIds(menuIdArray)
	var dest = make([]modeAuth.AuthSystemMenuTree, 0)
	if err := sql_res_to_tree.CreateSqlResFormatFactory().ScanToTreeData(menus, &dest); err != nil {
		variable.ZapLog.Error("根据用户id查询权限范围内的菜单数据树形化出错", zap.Error(err))
		return nil
	}
	data.Menus = dest
	return data
}

// 根据用户id查询所有可能的岗位节点id
func (u *UsersCurd) getUserAllOrgIds(userId int64) []int {
	//获取用户的所有岗位id
	postMember := modeAuth.CreateAuthPostMembersModelFactory("").GetByUserId(userId)

	postMemberIdArr := []int{}
	for k, _ := range postMember {
		postMemberIdArr = append(postMemberIdArr, postMember[k].FrAuthOrganizationPostId)
	}
	//根据岗位ID获取所有的岗位ID,父子级(需要去重)
	organization := modeAuth.CreateAuthOrganizationFactory("").GetByIds(postMemberIdArr)
	organizationIdArr := []int{}
	for _, v := range organization {
		idArr := strings.Split(v.PathInfo, ",")
		for _, vv := range idArr {
			id, _ := strconv.Atoi(vv)
			if id > 0 {
				organizationIdArr = append(organizationIdArr, id)
			}
		}
	}
	return organizationIdArr
}

// 查询用户打开指定的页面所拥有的按钮列表
func (u *UsersCurd) GetButtonListByMenuId(userId, menuId int64) (r []users.UserHasButtons) {
	orgIds := u.getUserAllOrgIds(userId)
	if list := u.userModel.GetButtonListByMenuId(orgIds, menuId); len(list) > 0 {
		return list
	}
	return nil
}

// casbin 控制接口权限使用
// 获取用户挂接组织机构的节点id(casbin对应表的roleId)，系统支持一人多岗位
func (u *UsersCurd) GetUserOrgIdsByRedis(context *gin.Context) []int {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	currentUser, exist := context.MustGet(tokenKey).(my_jwt.CustomClaims)
	if exist {
		return u.getUserAllOrgIds(currentUser.UserId)
	}
	return nil
}
