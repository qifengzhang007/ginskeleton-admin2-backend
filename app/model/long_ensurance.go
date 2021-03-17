package model

//  查询长护险信息

func CreateLongEnsuranceFactory(sqlType string) *LongEnsurance {
	return &LongEnsurance{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type LongEnsurance struct {
	BaseModel       `json:"-"`
	DateCategory    string `gorm:"column:date_category"`
	OrderNo         int    `gorm:"column:order_no"`
	PassAppDatetime string `gorm:"column:pass_app_datetime"`
	Num             int
}

// 表名
func (a *LongEnsurance) TableName() string {
	return ""
}

// 根据城市id、年龄段最低值查询老年人数据分布
func (a *LongEnsurance) GetLongEnsuranceList(cityId float64) []LongEnsurance {
	sql := "CALL  pro_get_long_insurance_detail(?) "
	var tmp []LongEnsurance
	if result := a.Raw(sql, cityId).Find(&tmp); result.RowsAffected > 0 {
		return tmp
	} else {
		return nil
	}
}
