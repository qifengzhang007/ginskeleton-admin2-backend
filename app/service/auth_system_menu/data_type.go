package auth_system_menu

// MenuListTree2 待分配的所有菜单，树形列表
type MenuListTree struct {
	SystemMenuButtonId int            `primaryKey:"yes" json:"system_menu_button_id"`
	SystemMenuFid      int            `fid:"SystemMenuButtonId"  json:"system_menu_fid"`
	Title              string         `json:"title"`
	NodeType           string         `json:"node_type" default:"menu"`
	Expand             int8           `json:"expand"`
	Sort               int            `json:"sort"`
	Children           []MenuListTree `json:"children"`
}

// AssignedSystemMenuButton 已分配给部门、岗位的系统菜单、按钮结构体
// 已分配系统给部门、岗位的的菜单、按钮返回结构体
type AssignedSystemMenuButton struct {
	SystemMenuButtonId     int                        `primaryKey:"yes" json:"system_menu_button_id"`
	SystemMenuFid          int                        `fid:"SystemMenuButtonId" json:"system_menu_fid"`
	Title                  string                     `json:"title"`
	NodeType               string                     `json:"node_type"`
	Expand                 int8                       `json:"expand"`
	OrgPostId              int                        `json:"org_post_id"`
	AuthPostMountHasMenuId int                        `json:"auth_post_mount_has_menu_id"`
	Children               []AssignedSystemMenuButton `gorm:"-" json:"children"`
}
