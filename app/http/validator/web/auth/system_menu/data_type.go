package system_menu

type Create struct {
	Title        string   `form:"title" json:"title" binding:"required,min=1"`
	Icon         string   `form:"icon" json:"icon"`
	Fid          *float64 `form:"fid" json:"fid" binding:"required,min=0"`
	Status       *float64 `form:"status" json:"status" binding:"required,min=0"`
	Sort         *float64 `form:"sort" json:"sort" binding:"required,min=0"`
	Name         string   `form:"name" json:"name" binding:"required,min=1"`
	Path         string   `form:"path" json:"path" binding:"required,min=1"`
	Component    string   `form:"component" json:"component" binding:"min=1"`
	ButtonString string   `form:"button_string" json:"button_string" binding:"min=0"`
	ButtonDelete string   `form:"button_delete" json:"button_delete" binding:"min=0"`
	Remark       string   `form:"remark" json:"remark"`
}

type Id struct {
	Id float64 `form:"id" json:"id" binding:"min=1"`
}

type Fid struct {
	Fid float64 `form:"fid"  json:"fid" binding:"min=0"`
}
