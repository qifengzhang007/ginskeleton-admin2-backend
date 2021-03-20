package province_city

//根据fid查询子级树形化数据结构
type ProvinceCityTree struct {
	Id         int64              `json:"id" primaryKey:"yes"`
	CreatedAt  string             `json:"created_at"`
	UpdatedAt  string             `json:"updated_at"`
	Fid        int64              `json:"fid"  fid:"Id"`
	Name       string             `json:"title"` // iview  框架树形组件的标题字段必须是  title
	Status     int                `json:"status"`
	Sort       int                `json:"sort"`
	Selected   int                `json:"selected"`
	Loading    bool               `json:"loading"`
	NodeLevel  int                `json:"node_level"`
	Remark     string             `json:"remark"`
	HasSubNode int                `json:"has_sub_node"`
	Children   []ProvinceCityTree `gorm:"-" json:"children"`
}
