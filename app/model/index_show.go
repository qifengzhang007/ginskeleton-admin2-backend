package model

import (
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
)

func CreateIndexShowFactory(sqlType string) *IndexShow {
	return &IndexShow{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type IndexShow struct {
	BaseModel `json:"-"`
	ApiValue  string `gorm:"column:api_value"  json:"api_value"`
}

// 表名
func (i *IndexShow) TableName() string {
	return "tb_index_show"
}

// 根据接口查询数据
func (i *IndexShow) GetApiValueByUrl(requestUrl string) string {
	sql := " SELECT   api_value   FROM `tb_index_show`  WHERE     api_uri_key=?  "
	result := i.Raw(sql, requestUrl).First(i)
	if result.Error == nil {
		return i.ApiValue
	} else {
		variable.ZapLog.Error("查询出错", zap.Error(result.Error))
		return ""
	}
}
