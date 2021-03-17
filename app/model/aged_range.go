package model

//  按照年龄段统计老年人在不同的街道分布数据

func CreateAgedRangeFactory(sqlType string) *AgedRange {
	return &AgedRange{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type AgedRange struct {
	BaseModel                  `json:"-"`
	AgedRanged                 string `gorm:"column:aged_ranged"  json:"aged_ranged"`
	CityAreaComunityStreetName string `gorm:"column:city_area_comunity_street_name"  json:"city_area_comunity_street_name"`
	Num                        int    `json:"num"`
	OrderNo                    int    `json:"order_no"`
}

// 表名
func (i *AgedRange) TableName() string {
	return ""
}

// 根据城市id、年龄段最低值查询老年人数据分布
func (i *AgedRange) GetAgedDataRange(cityId, age float64) []AgedRange {
	sql := " CALL  pro_get_age_range_data(?,?) "
	var tmp []AgedRange
	if result := i.Raw(sql, cityId, age).Find(&tmp); result.RowsAffected > 0 {
		return tmp
	} else {
		return nil
	}
}
