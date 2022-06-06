package auth_system_menu

// MenuListTree2 待分配的所有菜单，树形列表
type MenuListTree2 struct {
	SystemMenuId  int             `primaryKey:"yes" json:"system_menu_id"`
	SystemMenuFid int             `fid:"SystemMenuId"  json:"system_menu_fid"`
	Title         string          `json:"title"`
	NodeType      string          `json:"node_type" default:"menu"`
	Expand        int8            `json:"expand"`
	Sort          int             `json:"sort"`
	Children      []MenuListTree2 `json:"children"`
}

// 待分配系统菜单、model、按钮数据树形化
//type MenuListTree struct {
//	SystemMenuId  int    `primaryKey:"yes" json:"system_menu_id"`
//	SystemMenuFid int    `json:"system_menu_fid"`
//	Title         string `json:"title"`
//	MenuNodeType  string `json:"node_type" default:"menu"`
//	Expand        int8   `json:"expand"`
//	Children      []struct {
//		SystemMenuId  int    `primaryKey:"yes" json:"system_menu_id"`
//		SystemMenuFid int    `fid:"SystemMenuId" json:"system_menu_fid"`
//		Title         string `json:"title"`
//		MenuNodeType  string `json:"node_type" default:"menu"`
//		Expand        int8   `json:"expand"`
//		Children      []struct {
//			FrAuthSystemMenuId int    `fid:"SystemMenuId" json:"fr_auth_system_menu_id"`
//			SystemMenuId       int    `json:"system_menu_id"`
//			SystemMenuFid      int    `json:"system_menu_fid"`
//			ButtonId           int    `primaryKey:"yes" json:"button_id"`
//			ButtonName         string `json:"title"` // 这里json化的时候 设置为  title
//			NodeType           string `json:"node_type"`
//			Expand             int8   `json:"expand"`
//		} `json:"children"`
//	} `json:"children"`
//}

// 已分配给部门、岗位的系统菜单、按钮结构体
type AssignedMenuListTree struct {
	OrgPostId          int    `json:"org_post_id"`
	SystemMenuId       int8   `primaryKey:"yes" json:"system_menu_id"`
	SystemMenuFid      int8   `json:"system_menu_fid"`
	Title              string `json:"title"`
	PostMountHasMenuId int    `json:"post_mount_has_menu_id"`
	MenuNodeType       string `json:"node_type" default:"menu"`
	Expand             int8   `json:"expand"`
	Children           []struct {
		OrgPostId          int    `json:"org_post_id"`
		SystemMenuId       int8   `json:"system_menu_id" primaryKey:"yes"`
		SystemMenuFid      int8   `json:"system_menu_fid" fid:"SystemMenuId"`
		Title              string `json:"title"`
		PostMountHasMenuId int    `json:"post_mount_has_menu_id"`
		MenuNodeType       string `json:"node_type" default:"menu"`
		Children           []struct {
			FrMountHasMenuId         int    `fid:"PostMountHasMenuId" json:"post_mount_has_menu_id"`
			PostMountHasMenuButtonId int    `primaryKey:"yes" json:"post_mount_has_menu_button_id"`
			ButtonName               string `json:"title"` // 这里json化的时候 设置为  title
			NodeType                 string `json:"node_type"`
			Checked                  int8   `json:"checked"`
		} `json:"children"`
	} `json:"children"`
}
