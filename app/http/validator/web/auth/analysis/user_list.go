package analysis

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	webAuth "goskeleton/app/http/controller/web/auth"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/http/validator/web/auth/common_data_type"
	"goskeleton/app/utils/response"
)

type UserListWithPost struct {
	UserName string `form:"real_name" json:"user_name" `
	common_data_type.Page
}

func (u UserListWithPost) CheckParams(context *gin.Context) {
	//1.基本的验证规则没有通过
	if err := context.ShouldBind(&u); err != nil {
		errs := gin.H{
			"tips": "UserListWithPost 参数校验失败，参数不符合规定，user_name（可空）、page的值(>0)、limit 的值（>0)",
			"err":  err.Error(),
		}
		response.ErrorParam(context, errs)
		return
	}

	//  该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(u, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "UserListWithPost 表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&webAuth.AuthAnalysis{}).ListWithPost(extraAddBindDataContext)
	}
}
