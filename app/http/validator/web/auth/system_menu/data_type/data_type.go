package data_type

type Id struct {
	Id float64 `form:"id" json:"id" binding:"min=1"`
}

type Fid struct {
	Fid float64 `form:"fid"  json:"fid" binding:"min=0"`
}

// 系统菜单以及子表数据结构
type MenuCreate struct {
	Title        string   `form:"title" json:"title" binding:"required,min=1"`
	Icon         string   `form:"icon" json:"icon"`
	Fid          *float64 `form:"fid" json:"fid" binding:"required,min=0"`
	Status       *float64 `form:"status" json:"status" binding:"required,min=0"`
	Sort         *float64 `form:"sort" json:"sort" binding:"required,min=0"`
	Name         string   `form:"name" json:"name" binding:"required,min=1"`
	Path         string   `form:"path" json:"path" `
	Component    string   `form:"component" json:"component" binding:"min=1"`
	Remark       string   `form:"remark" json:"remark"`
	ButtonDelete string   `json:"button_delete"`
	ButtonArray  `json:"button_array"`
}

// 数据类型被使用时，shouldbindjson 对于数字是可以接受  int  int64   float64 ,shouldbind 函数对于数字只能接受  float64
type ButtonArray []struct {
	Id                 int64  `gorm:"primarykey" json:"id"`
	FrAuthSystemMenuId int64  `json:"fr_auth_system_menu_id"`
	FrAuthButtonCnEnId int64  `json:"fr_auth_button_cn_en_id"`
	RequestUrl         string `json:"request_url"`
	RequestMethod      string `json:"request_method"`
	Remark             string `json:"remark"`
	Status             int64  `json:"status"`
	CreatedAt          string
	UpdatedAt          string
}

// 菜单主表以及子表修改的数据结构
type MenuEdit struct {
	Id int64 `json:"id"`
	MenuCreate
}
