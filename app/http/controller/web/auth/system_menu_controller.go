package auth

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	modeAuth "goskeleton/app/model/auth"
	"goskeleton/app/service/auth_system_menu"
	"goskeleton/app/utils/response"
)

type SystemMenuController struct {
}

// 1.系统菜单列表
func (a *SystemMenuController) List(context *gin.Context) {
	var fid = context.GetFloat64(consts.ValidatorPrefix + "fid")
	var title = context.GetString(consts.ValidatorPrefix + "title")
	var limit = context.GetFloat64(consts.ValidatorPrefix + "limit")
	var limitStart = (context.GetFloat64(consts.ValidatorPrefix+"page") - 1) * limit

	systemMenuFac := modeAuth.CreateAuthSystemMenuFactory("")
	if counts, res := systemMenuFac.List(int(limitStart), int(limit), int(fid), title); counts > 0 {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"count": counts, "data": res})
		return
	}
	response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
}

// 根据ID获取子节点
func (a *SystemMenuController) GetByFid(c *gin.Context) {
	fid := c.GetFloat64(consts.ValidatorPrefix + "fid")
	data, err := modeAuth.CreateAuthSystemMenuFactory("").GetByFid(int(fid))
	if err == nil {
		response.Success(c, consts.CurdStatusOkMsg, data)
	} else {
		response.Fail(c, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, err)
	}
}

// 根据系统菜单的id获取挂接在全部可分配的按钮
// 获取菜单下的Button信息
func (a *SystemMenuController) GetMountButtonList(c *gin.Context) {
	menuId := c.GetFloat64(consts.ValidatorPrefix + "fr_auth_system_menu_id")
	data := modeAuth.CreateAuthSystemMenuButtonFactory("").MenuButton(int(menuId))
	response.Success(c, consts.CurdStatusOkMsg, data)
}

// 1.新增
func (a *SystemMenuController) Create(c *gin.Context) {
	// 限制菜单最大深度为3级，超过此值不允许添加
	fid := c.GetFloat64(consts.ValidatorPrefix + "fid")
	sysMenuFac := modeAuth.CreateAuthSystemMenuFactory("")
	if sysMenuFac.GetMenuLevel(int(fid)) >= 3 {
		response.Fail(c, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+",节点最大深度为3级", struct{}{})
		return
	}
	isOk, data := sysMenuFac.InsertData(c)
	if isOk {
		if (&auth_system_menu.AuthSystemMenuService{}).InsertButton(c, data.Id) {
			response.Success(c, consts.CurdStatusOkMsg, consts.CurdStatusOkCode)
			return
		}
	}
	response.Fail(c, consts.CurdCreatFailCode, consts.CurdCreatFailMsg, "请注意不要添加重复数据")
}

// 1.修改
func (a *SystemMenuController) Edit(c *gin.Context) {
	res := modeAuth.CreateAuthSystemMenuFactory("").UpdateData(c)
	if res {
		menuId := c.GetFloat64(consts.ValidatorPrefix + "id")
		if (&auth_system_menu.AuthSystemMenuService{}).UpdateButton(c, int64(menuId)) {
			response.Success(c, consts.CurdStatusOkMsg, consts.CurdStatusOkCode)
			return
		}
	}
	response.Fail(c, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "无数据被修改")
}

// 1.删除
func (a *SystemMenuController) Destroy(c *gin.Context) {
	id := c.GetFloat64(consts.ValidatorPrefix + "id")

	sysMenuFac := modeAuth.CreateAuthSystemMenuFactory("")
	//判断是否有子节点,如果有,则禁止删除
	if sysMenuFac.GetSubNodeCount(int(id)) > 0 {
		response.Fail(c, consts.CurdDeleteFailCode, "删除失败,存在子节点无法删除！", "")
		return
	}
	if sysMenuFac.DeleteData(int(id)) {
		response.Success(c, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(c, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg, "")
	}
}
