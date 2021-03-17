package auth

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	modeAuth "goskeleton/app/model/auth"
	"goskeleton/app/utils/response"
)

type PostMembersController struct {
}

// 1.岗位成员列表
func (p *PostMembersController) List(context *gin.Context) {
	var orgPostId = context.GetFloat64(consts.ValidatorPrefix + "org_post_id")
	var userName = context.GetString(consts.ValidatorPrefix + "user_name")
	var limit = context.GetFloat64(consts.ValidatorPrefix + "limit")
	var limitStart = (context.GetFloat64(consts.ValidatorPrefix+"page") - 1) * limit

	postMemberFac := modeAuth.CreateAuthPostMembersModelFactory("")
	counts := postMemberFac.GetCount(orgPostId, userName)
	if counts > 0 {
		res := postMemberFac.List(orgPostId, limitStart, limit, userName)
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"count": counts, "data": res})
		return
	}
	response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
}

//2.新增
func (p *PostMembersController) Create(context *gin.Context) {
	if modeAuth.CreateAuthPostMembersModelFactory("").InsertData(context) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg+"请注意不允许重复数据新增", "")
	}
}

//修改
func (p *PostMembersController) Edit(context *gin.Context) {
	if modeAuth.CreateAuthPostMembersModelFactory("").UpdateData(context) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
	}
}

// 删除

func (p *PostMembersController) Destroy(context *gin.Context) {
	var id = context.GetFloat64(consts.ValidatorPrefix + "id")
	if modeAuth.CreateAuthPostMembersModelFactory("").DeleteData(id) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg, "")
	}
}
