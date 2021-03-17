package button

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/web/auth"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type ButtonCreate struct {
	BaseField
}

// 验证器语法，参见 Register.go文件，有详细说明

func (c ButtonCreate) CheckParams(context *gin.Context) {
	//1.基本的验证规则没有通过
	if err := context.ShouldBind(&c); err != nil {
		errs := gin.H{
			"tips": "ButtonStore参数校验失败，参数校验失败，请检查cn_name(>=2)、en_name(>=2)、color(>=2)、allow_method(>=2)、status>=0",
			"err":  err.Error(),
		}
		response.ErrorParam(context, errs)
		return
	}

	//  该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(c, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "ButtonCreate表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&auth.Button{}).Create(extraAddBindDataContext)
	}
}
