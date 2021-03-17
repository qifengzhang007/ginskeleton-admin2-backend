package auth

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/utils/response"
)

type Button struct {
}

//1按钮列表
func (s *Button) List(context *gin.Context) {
	buttonName := context.GetString(consts.ValidatorPrefix + "button_name")
	page := context.GetFloat64(consts.ValidatorPrefix + "page")
	limit := context.GetFloat64(consts.ValidatorPrefix + "limit")
	limitStart := (page - 1) * limit

	totalCounts, showList := model.CreateButtonCnEnFactory("").List(buttonName, limitStart, limit)
	if totalCounts > 0 && showList != nil {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"count": totalCounts, "data": showList})
	} else {
		response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

//2.按钮新增(store)
func (s *Button) Create(context *gin.Context) {
	if model.CreateButtonCnEnFactory("").InsertData(context) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+",新增错误", "")
	}
}

//5.按钮更新(update)
func (s *Button) Edit(context *gin.Context) {
	//注意：这里没有实现权限控制逻辑，例如：超级管理管理员可以更新全部用户数据，普通用户只能修改自己的数据。目前只是验证了token有效、合法之后就可以进行后续操作
	// 实际使用请根据真是业务实现权限控制逻辑、再进行数据库操作
	if model.CreateButtonCnEnFactory("").UpdateData(context) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
	}

}

//6.删除记录
func (u *Button) Destroy(context *gin.Context) {
	userId := context.GetFloat64(consts.ValidatorPrefix + "id")
	if model.CreateButtonCnEnFactory("").DeleteData(int(userId)) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg, "")
	}
}
