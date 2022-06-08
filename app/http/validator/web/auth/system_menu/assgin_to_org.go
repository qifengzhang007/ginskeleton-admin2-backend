package system_menu

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/web/auth"
	"goskeleton/app/http/validator/web/auth/system_menu/data_type"
	"goskeleton/app/utils/response"
)

// 为组织机构（部门岗位）分配菜单、按钮权限 参数校验
type AssginSystemMenuToOrg struct {
	MenuList   data_type.MenuButtonList `json:"menu_list" binding:"dive"`
	ButtonList data_type.MenuButtonList `json:"button_list"  binding:"dive"`
}

// 验证器语法，参见 Register.go文件，有详细说明
func (a AssginSystemMenuToOrg) CheckParams(context *gin.Context) {
	//1.基本的验证规则没有通过
	if err := context.ShouldBindJSON(&a); err != nil {
		response.ValidatorError(context, err)
		return
	}

	context.Set("auth_assign_menu_list", a.MenuList)
	context.Set("auth_assign_button_list", a.ButtonList)
	// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
	(&auth.SystemMenuAssignController{}).AssignAuthToOrg(context)

}
