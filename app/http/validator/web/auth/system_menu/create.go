package system_menu

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/controller/web/auth"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/http/validator/web/auth/system_menu/data_type"
	"goskeleton/app/utils/response"
)

type SystemMenuCreate struct {
	data_type.MenuCreate
}

// 验证器语法，参见 Register.go文件，有详细说明
func (s SystemMenuCreate) CheckParams(context *gin.Context) {
	//1.基本的验证规则没有通过
	if err := context.ShouldBindJSON(&s); err != nil {
		errs := gin.H{
			"tips": "SystemMenuCreate  参数校验失败，参数不符合规定,name ≥ 0、path ≥ 0、component ≥ 0、 status ≥ 0、 fid  ≥ 0 、title ≥ 1",
			"err":  err.Error(),
		}
		response.ErrorParam(context, errs)
		return
	}
	context.Set(variable.SystemCreateKey, s.MenuCreate)
	//  该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(s, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "SystemMenuCreate 表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&auth.SystemMenuController{}).Create(extraAddBindDataContext)
	}
}
