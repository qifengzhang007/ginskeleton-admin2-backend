package users

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/web"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type Destroy struct {
	Id
}

// 验证器语法，参见 Register.go文件，有详细说明

func (d Destroy) CheckParams(context *gin.Context) {

	if err := context.ShouldBind(&d); err != nil {
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(d, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "Destroy 表单参数验证器json化失败", "")
		return
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&web.Users{}).Destroy(extraAddBindDataContext)

		// 以下代码为模拟 前置、后置函数的回调代码
		/*
			func(before_callback_fn func(context *gin.Context) bool, after_callback_fn func(context *gin.Context)) {
				if before_callback_fn(extraAddBindDataContext) {
					defer after_callback_fn(extraAddBindDataContext)
					(&Web.Users{}).Destroy(extraAddBindDataContext)
				} else {
					// 这里编写前置函数验证不通过的相关返回提示逻辑...

				}
			}((&Users.DestroyBefore{}).Before, (&Users.DestroyAfter{}).After)
		*/
	}
}
