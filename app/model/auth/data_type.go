package auth

//菜单分配文件相关的数据类型

// 待分配系统的菜单、model、按钮返回结构体
type AuthSystemMenuButton struct {
	SystemMenuFid      int
	SystemMenuId       int
	FrAuthSystemMenuId int
	Title              string
	NodeType           string
	ButtonId           int
	ButtonName         string
	Expand             int8
}

// AssignedSystemMenuButton2 已分配系统给部门、岗位的的菜单、按钮返回结构体
type AssignedSystemMenuButton2 struct {
	SystemMenuId           int
	SystemMenuFid          int
	Title                  string
	NodeType               string
	Expand                 int8
	OrgPostId              int
	AuthPostMountHasMenuId int
	sort1                  int
	sort2                  int
}

// 已分配系统给部门、岗位的的菜单、按钮返回结构体
type AssignedSystemMenuButton struct {
	OrgPostId                int
	SystemMenuId             int8
	SystemMenuFid            int8
	PostMountHasMenuId       int
	FrMountHasMenuId         int
	PostMountHasMenuButtonId int
	Title                    string
	NodeType                 string
	ButtonId                 int
	ButtonName               string
	Checked                  int8
	Expand                   int8
}

// 根据用户id查询已经分配的权限树形结构
type OrgTree struct {
	Id       int       ` json:"id" primaryKey:"yes"`
	OrgTitle string    `json:"title"`
	OrgFid   int       `json:"org_fid" fid:"Id"`
	NodeType string    `json:"node_type"`
	Expand   bool      `json:"expand"`
	Children []OrgTree `gorm:"-" json:"children"`
}

// 组织机构数据结构

type AuthOrganizationPostTree struct {
	Id       int                        `json:"id"`
	Fid      string                     `json:"fid"`
	Title    string                     `json:"title"`
	Status   string                     `json:"status"`
	PathInfo string                     `json:"path_info"`
	Remark   string                     `json:"remark"`
	IsLeaf   bool                       `json:"is_leaf"` // 是否为叶子节点
	Children []AuthOrganizationPostTree `gorm:"-" json:"children"`
}

// 根据ids查询数据列表
type AllAuth struct {
	Id    int    `json:"id"`
	Title string `json:"string"`
	Fid   int    `json:"fid"`
}

//岗位成员
//  定义不同的查询结果返回的数据结构体
type MemberList struct {
	Id        int    `json:"id"`
	OrgPostId int    `json:"org_post_id"`
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	RealName  string `json:"real_name"`
	Status    int    `json:"status"`
	PostName  string `json:"post_name"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// 系统菜单

//查询显示的菜单列表，结合了待分配的按钮作为子级数据显示
type AuthSystemMenuButtonList struct {
	Id                 int
	Fid                int
	Icon               string
	Title              string
	Name               string
	Loading            bool
	Path               string
	Component          string
	Status             int
	Sort               int
	Remark             string
	FrAuthSystemMenuId int
	ButtonId           int
	ButtonName         string
	ButtonColor        string
}

// 系统菜单以及待分配的按钮树形化结构体
type AuthSystemMenuButtonListTree struct {
	Id        int    `json:"id" primaryKey:"yes"`
	Fid       int    `json:"fid"`
	Icon      string `json:"icon"`
	Title     string `json:"title"`
	Name      string `json:"name"`
	Loading   bool   `json:"loading"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Status    int    `json:"status"`
	Sort      int    `json:"sort"`
	Remark    string `json:"remark"`
	Children  []struct {
		FrAuthSystemMenuId int    `fid:"Id" json:"fr_auth_system_menu_id" `
		ButtonId           int    `json:"button_id" primaryKey:"yes"`
		ButtonName         string `json:"button_name"`
		ButtonColor        string `json:"button_color"`
	} `json:"button_list" gorm:"-"`
}

//接口返回数据类型结构体（左侧菜单树）
type AuthSystemMenuTree struct {
	Id         int                  `json:"id" primaryKey:"yes" `
	Fid        int                  `json:"fid" fid:"Id"`
	Title      string               `json:"title"`
	Name       string               `json:"name"`
	Icon       string               `gorm:"icon" json:"icon"`
	Path       string               `gorm:"path" json:"path"`
	NodeLevel  int                  `json:"node_level"`
	Component  string               `json:"component"`
	HasSubNode int                  `json:"has_sub_node"`
	IsLeaf     bool                 `json:"is_leaf"`
	Children   []AuthSystemMenuTree `gorm:"-" json:"children"`
}

// 系统菜单挂接的按钮列表

type SystemMenuButtonList struct {
	Id                 int    `json:"id"`
	FrAuthSystemMenuId int    `json:"fr_auth_system_menu_id"`
	FrAuthButtonCnEnId int    `json:"fr_auth_button_cn_en_id"`
	ButtonName         string `json:"button_name"`
	RequestUrl         string `json:"request_url"`
	RequestMethod      string `json:"request_method"`
	Status             int    `json:"status"`
	Remark             string `json:"remark"`
}
