package button

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/web/auth"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/http/validator/web/auth/common_data_type"
	"goskeleton/app/utils/response"
)

type ButtonList struct {
	ButtonName string `form:"button_name" json:"button_name" binding:"min=0"` // 必填，页面值>0
	common_data_type.Page
}

// 验证器语法，参见 Register.go文件，有详细说明
func (b ButtonList) CheckParams(context *gin.Context) {
	//1.基本的验证规则没有通过
	if err := context.ShouldBind(&b); err != nil {
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(b, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "ButtonList 表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&auth.Button{}).List(extraAddBindDataContext)
	}
}
