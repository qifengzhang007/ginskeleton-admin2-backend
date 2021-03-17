package common

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/utils/response"
)

// 公共模块： 按钮

type ButtonCnEn struct {
}

func (b *ButtonCnEn) List(context *gin.Context) {
	keyWord := context.GetString(consts.ValidatorPrefix + "key_word")
	page := context.GetFloat64(consts.ValidatorPrefix + "page")
	limits := context.GetFloat64(consts.ValidatorPrefix + "limit")
	limitStart := (page - 1) * limits

	totalCOunts, showList := model.CreateButtonCnEnFactory("").Show(keyWord, int(limitStart), int(limits))
	if totalCOunts > 0 && showList != nil {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"count": totalCOunts, "data": showList})
	} else {
		response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}
