package system_menu

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/web/auth"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type SystemMenuListAllList struct {
	Title string `form:"title" json:"title"`
}

func (s SystemMenuListAllList) CheckParams(context *gin.Context) {
	//1.基本的验证规则没有通过
	if err := context.ShouldBind(&s); err != nil {
		errs := gin.H{
			"tips": "SystemMenuListAllList  参数校验失败，参数不符合规定 title 可空",
			"err":  err.Error(),
		}
		response.ErrorParam(context, errs)
		return
	}
	//  该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(s, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "SystemMenuList 表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&auth.SystemMenuAssignController{}).SystemMenuButtonAllList(extraAddBindDataContext)
	}
}
