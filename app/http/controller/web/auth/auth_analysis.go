package auth

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model/users"
	"goskeleton/app/service/auth_post_members"
	"goskeleton/app/utils/response"
)

// 分析用户权限来源

type AuthAnalysis struct {
}

//  查询用户信息(带岗位)
func (a *AuthAnalysis) ListWithPost(context *gin.Context) {
	userName := context.GetString(consts.ValidatorPrefix + "user_name")
	page := context.GetFloat64(consts.ValidatorPrefix + "page")
	limit := context.GetFloat64(consts.ValidatorPrefix + "limit")
	limitStart := (page - 1) * limit

	totalCounts, showList := users.CreateUserFactory("").ListWithPost(userName, int(limitStart), int(limit))
	if totalCounts > 0 && showList != nil {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"count": totalCounts, "data": showList})
	} else {
		response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}

}

//根据用户ID获取所有权限的来源
func (a *AuthAnalysis) GetAuthByUserId(context *gin.Context) {
	id := context.GetFloat64(consts.ValidatorPrefix + "id")
	//根据用户ID,查询隶属哪些组织机构
	if data := (&auth_post_members.AuthPostMembersService{}).FindOrgs(int64(id)); data != nil {
		response.Success(context, consts.CurdStatusOkMsg, data)
	} else {
		response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}

}
