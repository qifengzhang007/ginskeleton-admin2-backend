package web

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model/province_city"
	"goskeleton/app/utils/response"
)

type ProvinceCityController struct {
}

// 1.省份城市列表
func (a *ProvinceCityController) List(context *gin.Context) {
	var fid = context.GetFloat64(consts.ValidatorPrefix + "fid")
	var name = context.GetString(consts.ValidatorPrefix + "name")
	var limit = context.GetFloat64(consts.ValidatorPrefix + "limit")
	var limitStart = (context.GetFloat64(consts.ValidatorPrefix+"page") - 1) * limit

	cityFac := province_city.CreateProvinceCityFactory("")
	if counts := cityFac.GetCount(int(fid), name); counts > 0 {
		res := cityFac.List(name, int(fid), int(limitStart), int(limit))
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"count": counts, "data": res})
		return
	}
	response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
}

// 1.根据fid查询子节点列表
func (a *ProvinceCityController) SubList(context *gin.Context) {
	var fid = context.GetFloat64(consts.ValidatorPrefix + "fid")

	if subList := province_city.CreateProvinceCityFactory("").GetSubListByfid(int(fid)); len(subList) > 0 {
		response.Success(context, consts.CurdStatusOkMsg, subList)
	} else {
		response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

// 新增
func (a *ProvinceCityController) Create(c *gin.Context) {
	if province_city.CreateProvinceCityFactory("").InsertData(c) {
		response.Success(c, consts.CurdStatusOkMsg, consts.CurdStatusOkCode)
	} else {
		response.Fail(c, consts.CurdCreatFailCode, consts.CurdCreatFailMsg+"请注意不要添加重复数据", "")
	}
}

// 1.根据ID获取子节点
func (a *ProvinceCityController) GetSubList(c *gin.Context) {
	id := c.GetFloat64(consts.ValidatorPrefix + "id")
	data := province_city.CreateProvinceCityFactory("").GetSubListByfid(int(id))
	if len(data) == 0 {
		response.Fail(c, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	} else {
		response.Success(c, consts.CurdStatusOkMsg, data)
	}
}

// 修改
func (a *ProvinceCityController) Edit(c *gin.Context) {
	if res := province_city.CreateProvinceCityFactory("").UpdateData(c); res {
		response.Success(c, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(c, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
	}
}

// 删除
func (a *ProvinceCityController) Destroy(c *gin.Context) {
	id := c.GetFloat64(consts.ValidatorPrefix + "id")
	cityFac := province_city.CreateProvinceCityFactory("")

	if cityFac.HasSubNode(int(id)) > 0 {
		response.Fail(c, consts.CurdDeleteFailCode, "该节点下有子节点,禁止删除", "")
	} else {
		if cityFac.DeleteData(int(id)) {
			response.Success(c, consts.CurdStatusOkMsg, "")
		} else {
			response.Fail(c, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg, "")
		}
	}
}
