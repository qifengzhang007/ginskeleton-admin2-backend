package model

//  查询长护险信息

func CreateseSverModuleFactory(sqlType string) *SverModule {
	return &SverModule{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type SverModule struct {
	BaseModel `json:"-"`
}

// 助浴服务list
type ZhuYuLists struct {
	RealName   string `json:"real_name"`
	Age        int    `json:"age"`
	CardNo     string `json:"card_no"`
	Telephone  string `json:"telephone"`
	StreetName string `json:"street_name"`
}

// 老吾老list
type LaoWuLaoLists struct {
	StreetName string

	ProvidedServerTitle      string
	ProvidedServerPersonNum  int
	ProvidedServerPersonUnit string

	AssistedAgedTitle string
	AssistedAgedNum   int
	AssistedAgedUnit  string

	JoinOrgTitle string
	JoinOrgNum   int
	JoinOrgUnit  string

	ProvidedServerNumTitle string
	ProvidedServerNum      int
	ProvidedServerUnit     string

	ServerItemTitle string
	ServerItemNum   int
	ServerItemUnit  string
}

// 表名
func (s *SverModule) TableName() string {
	return ""
}

// 享受服务的人员list
// 参数：   城市id、服务条目id、page、limit
func (s *SverModule) GetUseServicePeopleList(cityId, limitStart, limit float64) (int64, []ZhuYuLists) {
	// 查询数据总条数
	sql := "CALL  pro_get_useServiceAged_lists(?,?,?)"
	var tmp []ZhuYuLists
	res1 := s.Raw(sql, cityId, 0, 1000).Find(&tmp)
	if res1.RowsAffected > 0 {
		sql = "CALL  pro_get_useServiceAged_lists(?,?,?) "
		tmp = make([]ZhuYuLists, 0)
		if result := s.Raw(sql, cityId, limitStart, limit).Find(&tmp); result.RowsAffected > 0 {
			return res1.RowsAffected, tmp
		}
	}
	return 0, nil
}

// 老伙伴list
// 参数：   城市id
func (s *SverModule) GetServerModuleLists(cityId, itemId, limitStart, limit float64) []LaoWuLaoLists {
	sql := "CALL pro_get_server_module(?,?,?,?) "
	var tmp []LaoWuLaoLists
	if result := s.Raw(sql, cityId, itemId, limitStart, limit).Find(&tmp); result.RowsAffected > 0 {
		return tmp
	} else {
		return nil
	}
}
