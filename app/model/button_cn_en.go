package model

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/data_bind"
	"strings"
)

func CreateButtonCnEnFactory(sqlType string) *ButtonCnEnModel {
	return &ButtonCnEnModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type ButtonCnEnModel struct {
	BaseModel
	CnName      string `json:"cn_name"`
	EnName      string `json:"en_name"`
	AllowMethod string `json:"allow_method"`
	Color       string `json:"color"`
	Status      int    `json:"status"`
	Remark      string `json:"remark"`
}

// 表名
func (b *ButtonCnEnModel) TableName() string {
	return "tb_auth_button_cn_en"
}

// 根据关键词查询用户表的条数
func (b *ButtonCnEnModel) getCounts(keyWords string) (counts int64) {
	sql := "select  count(*) as counts from tb_auth_button_cn_en WHERE  ( en_name like ? or cn_name like  ?) "
	if _ = b.Raw(sql, "%"+keyWords+"%", "%"+keyWords+"%").First(&counts); counts > 0 {
		return counts
	} else {
		return 0
	}
}

// 查询（根据关键词模糊查询）
func (b *ButtonCnEnModel) Show(keyWords string, limitStart int, limitItems int) (totalCounts int64, temp []ButtonCnEnModel) {
	totalCounts = b.getCounts(keyWords)
	if totalCounts > 0 {
		sql := "SELECT  `id`, `cn_name`,`allow_method`, `en_name`, `remark`,DATE_FORMAT(created_at,'%Y-%m-%d %h:%i:%s')  AS created_at," +
			" DATE_FORMAT(updated_at,'%Y-%m-%d %h:%i:%s')  AS updated_at   FROM  `tb_auth_button_cn_en`  WHERE  ( cn_name like ? or en_name like  ?) LIMIT ?,?"
		if res := b.Raw(sql, "%"+keyWords+"%", "%"+keyWords+"%", limitStart, limitItems).Find(&temp); res.RowsAffected > 0 {
			return totalCounts, temp
		} else {
			return totalCounts, nil
		}
	}
	return 0, nil
}

//按钮编辑页的列表展示

func (a *ButtonCnEnModel) getCountsByButtonName(cnName string) (count int64) {
	if res := a.Model(a).Where("cn_name like ?", "%"+cnName+"%").Count(&count); res.Error == nil {
		return count
	}
	return 0
}

func (a *ButtonCnEnModel) List(cnName string, limitStart, limit float64) (counts int64, data []ButtonCnEnModel) {
	counts = a.getCountsByButtonName(cnName)
	if counts > 0 {
		if err := a.Model(a).
			Select("id", "en_name", "cn_name", "allow_method", "color", "status", "cn_name", "remark", "DATE_FORMAT(created_at,'%Y-%m-%d %H:%i:%s')  as created_at", "DATE_FORMAT(updated_at,'%Y-%m-%d %H:%i:%s')  as updated_at").
			Where("cn_name like ?", "%"+cnName+"%").Offset(int(limitStart)).Limit(int(limit)).Find(&data); err.Error == nil {
			return
		}
	}

	return 0, nil
}

//新增
func (b *ButtonCnEnModel) InsertData(c *gin.Context) bool {
	var tmp ButtonCnEnModel
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		tmp.AllowMethod = strings.ToUpper(tmp.AllowMethod)
		if res := b.Create(&tmp); res.Error == nil {
			return true
		} else {
			variable.ZapLog.Error("ButtonModel 数据新增出错", zap.Error(res.Error))
		}
	} else {
		variable.ZapLog.Error("ButtonModel 数据绑定出错", zap.Error(err))
	}
	return false
}

//更新
func (b *ButtonCnEnModel) UpdateData(c *gin.Context) bool {
	var tmp ButtonCnEnModel
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		tmp.AllowMethod = strings.ToUpper(tmp.AllowMethod)
		// Omit 表示忽略指定字段(CreatedAt)，其他字段全量更新
		if res := b.Omit("CreatedAt").Save(tmp); res.Error != nil {
			variable.ZapLog.Error("ButtonModel 数据修改出错", zap.Error(res.Error))
		} else {
			return true
		}
	} else {
		variable.ZapLog.Error("ButtonModel 数据绑定出错", zap.Error(err))
	}
	return false
}

//删除
func (b *ButtonCnEnModel) DeleteData(id int) bool {
	if b.Delete(b, id).Error == nil {
		return true
	}
	return false
}
