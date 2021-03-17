package model

//  按照健康状况统计老年人在不同的街道分布数据

func CreateHealthyStaticsFactory(sqlType string) *HealthyStatics {
	return &HealthyStatics{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type HealthyStatics struct {
	BaseModel                  `json:"-"`
	HealthyName                string `gorm:"column:healthy_name"`
	CityAreaComunityStreetName string `gorm:"city_area_comunity_street_name"`
	Num                        int
	OrderNo                    int
}

// 表名
func (h *HealthyStatics) TableName() string {
	return ""
}

// 根据城市id、年龄段最低值查询老年人数据分布
func (h *HealthyStatics) GetHealthyStatics(cityId float64) []HealthyStatics {
	sql := " CALL pro_get_aged_healthy_statictis(?)"
	var tmp []HealthyStatics
	if result := h.Raw(sql, cityId).Find(&tmp); result.RowsAffected > 0 {
		return tmp
	} else {
		return nil
	}
}
