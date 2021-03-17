package model

//  查询养老机构的具体信息

func CreateAgedOrgFactory(sqlType string) *AgedOrg {
	return &AgedOrg{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type AgedOrg struct {
	BaseModel            `json:"-"`
	Name                 string `json:"name"`
	Logo                 string `json:"logo"`
	Addr                 string `json:"addr"`
	ConcatWay            string `gorm:"column:concat_way" json:"concat_way"`
	WorkTimeRange        string `gorm:"column:work_time_range" json:"work_time_range"`
	Employers            int    `json:"employers"`
	ProvidedServerPerson int    `gorm:"column:provided_server_person" json:"provided_server_person"`
	Longitude            string `json:"longitude"`
	Latitude             string `json:"latitude"`
	ServerItems          string `json:"server_items"`
}

// 表名
func (a *AgedOrg) TableName() string {
	return ""
}

// 根据城市id、年龄段最低值查询老年人数据分布
func (a *AgedOrg) GetAgedOrg(cityId, orgPropertyId float64) []AgedOrg {
	sql := " CALL  pro_get_aged_org_info(?,?) "
	var tmp []AgedOrg
	if result := a.Raw(sql, cityId, orgPropertyId).Find(&tmp); result.RowsAffected > 0 {
		return tmp
	} else {
		return nil
	}
}
