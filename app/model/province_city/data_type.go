package province_city

//根据fid查询子级树形化数据结构
type ProvinceCityTree struct {
	Id        int64              `json:"id" primaryKey:"yes"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
	Fid       int64              `json:"fid"  fid:"Id"`
	Name      string             `json:"title"`
	Status    int                `json:"status"`
	Sort      int                `json:"sort"`
	Selected  int                `json:"selected"`
	Loading   bool               `json:"loading"`
	NodeLevel int                `json:"node_level"`
	Remark    string             `json:"remark"`
	IsLeaf    bool               `json:"is_leaf"` // 是否为叶子节点
	Children  []ProvinceCityTree `gorm:"-" json:"children"`
}

//type ProvinceCityModelList struct {
//	ProvinceCityModel
//	Ftitle string `json:"ftitle"`
//}
