package auth

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	modeAuth "goskeleton/app/model/auth"
	"goskeleton/app/service/auth_post_members"
	"goskeleton/app/utils/response"
)

type OrganizationPostController struct {
}

// 1.组织机构列表
func (a *OrganizationPostController) List(context *gin.Context) {
	var fid = context.GetFloat64(consts.ValidatorPrefix + "fid")
	var title = context.GetString(consts.ValidatorPrefix + "title")
	var limit = context.GetFloat64(consts.ValidatorPrefix + "limit")
	var limitStart = (context.GetFloat64(consts.ValidatorPrefix+"page") - 1) * limit

	orgPostFac := modeAuth.CreateAuthOrganizationFactory("")
	if counts := orgPostFac.GetCount(int(fid), title); counts > 0 {
		res := orgPostFac.List(int(limitStart), int(limit), int(fid), title)
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"count": counts, "data": res})
		return
	}
	response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
}

// 新增
func (a *OrganizationPostController) Create(c *gin.Context) {
	if modeAuth.CreateAuthOrganizationFactory("").InsertData(c) {
		response.Success(c, consts.CurdStatusOkMsg, consts.CurdStatusOkCode)
	} else {
		response.Fail(c, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+"请注意不要添加重复数据", "")
	}
}

// 1.根据ID获取子节点
func (a *OrganizationPostController) GetByFid(c *gin.Context) {
	fid := c.GetFloat64(consts.ValidatorPrefix + "fid")
	data, err := modeAuth.CreateAuthOrganizationFactory("").GetByFid(int(fid))
	if err != nil {
		response.Fail(c, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, err)
	} else {
		response.Success(c, consts.CurdStatusOkMsg, data)
	}
}

// 修改
func (a *OrganizationPostController) Edit(c *gin.Context) {
	res := modeAuth.CreateAuthOrganizationFactory("").UpdateData(c)
	if res {
		response.Success(c, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(c, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
	}
}

// 删除
func (a *OrganizationPostController) Destroy(c *gin.Context) {
	id := c.GetFloat64(consts.ValidatorPrefix + "id")

	models := modeAuth.CreateAuthOrganizationFactory("")
	//判断是否有子节点,如果有,则禁止删除
	if models.HasSubList(int(id)) > 0 {
		response.Fail(c, consts.CurdDeleteFailCode, "该节点下有子节点,禁止删除", "")
	} else {
		if models.DeleteData(int(id)) {
			response.Success(c, consts.CurdStatusOkMsg, "")
		} else {
			response.Fail(c, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg, "")
		}
	}
}

//根据用户ID获取所有的组织机构和权限
func (a *OrganizationPostController) GetAuthByUserId(c *gin.Context) {
	//var err error
	//models := model.CreateAuthOrganizationFactory("")
	id := c.GetFloat64(consts.ValidatorPrefix + "id")
	//根据用户ID,查询隶属哪些组织机构
	data := (&auth_post_members.AuthPostMembersService{}).FindOrgs(int64(id))
	//model.CreateAuthPostMembersModelFactory("")
	response.Success(c, consts.CurdStatusOkMsg, data)
}
