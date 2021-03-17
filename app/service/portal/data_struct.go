package portal

type StaticsRange struct {
	Title   string           `json:"title"`
	OrderNo int              `json:"order_no"`
	Deatil  []ComunityStreet `json:"detail"`
}

//返回的数据结构
type ComunityStreet struct {
	StreetName string `json:"street_name"`
	Num        int    `json:"num"`
	Unit       string `json:"unit"`
}

// 长护险数据结构（主、从构成层次结构数据格式）
type LongEnsurance struct {
	DateCategory string   `json:"date_category"`
	OrderNo      int      `json:"order_no"`
	Deatil       []Detail `json:"detail"`
}

//返回的数据结构
type Detail struct {
	PassAppDate string `json:"pass_app_date"`
	Num         int    `json:"num"`
	Unit        string `json:"unit"`
}

// 服务项模块

type LaoWuLaoLists struct {
	StreetName string         `json:"street_name"`
	Detail     []LaoWuLaoItem `json:"detail"`
}

type LaoWuLaoItem struct {
	ItemName string `json:"item_name"`
	ItemNum  int    `json:"item_num"`
	ItemUnit string `json:"item_unit"`
}
