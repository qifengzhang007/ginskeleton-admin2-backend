package auth

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	modeAuth "goskeleton/app/model/auth"
	"goskeleton/app/service/auth_system_menu"
	"goskeleton/app/utils/response"
)

type SystemMenuAssignController struct {
}

// 待分配的系统菜单以及挂接的按钮
func (a *SystemMenuAssignController) SystemMenuButtonAllList(context *gin.Context) {

	totalCOunts, showList := modeAuth.CreateAuthMenuAssignFactory("").GetSystemMenuButtonList()
	if totalCOunts > 0 && showList != nil {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"counts": totalCOunts, "data": (&auth_system_menu.AuthSystemMenuService{}).SystemMenuButtonToTree(showList)})
	} else {
		response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

// 已分配给部门、岗位的系统菜单，以及菜单挂接的按钮
func (a *SystemMenuAssignController) AssignedToOrgPostMenuButton(context *gin.Context) {
	orgPostId := context.GetFloat64(consts.ValidatorPrefix + "org_post_id")

	totalCOunts, showList := modeAuth.CreateAuthMenuAssignFactory("").GetAssignedMenuButtonList(int(orgPostId))
	if totalCOunts > 0 && showList != nil {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"counts": totalCOunts, "data": (&auth_system_menu.AuthSystemMenuService{}).AssginedMenuButtonToTree(showList)})
	} else {
		response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

// 为组织机构（部门、岗位）分配权限
func (a *SystemMenuAssignController) AssignAuthToOrg(context *gin.Context) {
	orgPostId := context.GetFloat64(consts.ValidatorPrefix + "org_post_id")
	systemMenuId := context.GetFloat64(consts.ValidatorPrefix + "system_menu_id")
	systemMenuFid := context.GetFloat64(consts.ValidatorPrefix + "system_menu_fid")
	buttonId := context.GetFloat64(consts.ValidatorPrefix + "button_id")
	nodeType := context.GetString(consts.ValidatorPrefix + "node_type")

	menuAssignFac := modeAuth.CreateAuthMenuAssignFactory("")
	res := menuAssignFac.AssginAuthForOrg(int(orgPostId), int(systemMenuId), int(systemMenuFid), int(buttonId), nodeType)
	if res {
		response.Success(context, consts.AuthAssginOkMsg, "")
		return
	}
	response.Fail(context, consts.AuthAssginFailCode, consts.AuthAssginFailMsg, "")
}

// 删除已经分配给组织机构（部门、岗位）的权限
func (a *SystemMenuAssignController) DeleteAuthFromOrg(context *gin.Context) {
	postMountHasMenuId := context.GetFloat64(consts.ValidatorPrefix + "post_mount_has_menu_id")
	postMountHasMenuButtonId := context.GetFloat64(consts.ValidatorPrefix + "post_mount_has_menu_button_id")
	nodeType := context.GetString(consts.ValidatorPrefix + "node_type")
	res := modeAuth.CreateAuthMenuAssignFactory("").DeleteAuthFromOrg(int(postMountHasMenuId), int(postMountHasMenuButtonId), nodeType)
	if res {
		response.Success(context, consts.AuthDelOkMsg, "")
	} else {
		response.Fail(context, consts.AuthDelFailCode, consts.AuthDelFailMsg, "")
	}
}
