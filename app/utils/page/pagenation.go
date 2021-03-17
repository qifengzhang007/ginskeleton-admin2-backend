package page

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"math"
)

type pageJson struct {
	Page  int `form:"page" json:"page" binding:"number"`   // 必填，页面值>0
	Limit int `form:"limit" json:"limit" binding:"number"` // 必填，每页条数值>0
}

var pageStruct pageJson

func PageInfo(c *gin.Context, total int64) (page int, limit int) {
	pageStruct = pageJson{}
	c.ShouldBind(&pageStruct)
	limit = pageSize(c)
	totalPage := totalPage(total, limit)
	page = getPage(totalPage)

	return (page - 1) * limit, limit
}

func totalPage(total int64, limit int) int {
	totalPage := math.Ceil(float64(total) / float64(limit))
	return int(totalPage)
}

func getPage(totalPage int) int {
	page := pageStruct.Page
	if page != 0 {
		if page <= 0 {
			page = 1
		} else if page > totalPage {
			page = totalPage
		}
	}

	return page
}

func pageSize(c *gin.Context) int {
	limit := pageStruct.Limit
	if limit < 0 || limit > 200 {
		limit = variable.ConfigYml.GetInt("PageSize")
	}
	return limit
}
