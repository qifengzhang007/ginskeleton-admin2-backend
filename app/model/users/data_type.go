package users

// 定义一个部门、岗位、用户名关键词查询返回所需结构体
type OrgPostList struct {
	Id          int64  `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	UserName    string `gorm:"column:user_name" json:"user_name"`
	Phone       string `json:"phone"`
	RealName    string `gorm:"column:real_name" json:"real_name"`
	Status      int    `json:"status"`
	Remark      string `json:"remark"`
	LastLoginIp string `gorm:"column:last_login_ip" json:"last_login_ip"`
	OrgPostName string `json:"org_post_name"`
}

// 用户权限分析界面带岗位数据查询
type AnalysisiUserList struct {
	Id       int    `json:"id"`
	UserName string `gorm:"column:user_name" json:"user_name"`
	RealName string `gorm:"column:real_name" json:"real_name"`
	PostName string `json:"post_name"`
}

// 用户在指定页面已分配的按钮列表
type UserHasButtons struct {
	Id     int    `json:"id"`
	CnName string `json:"cn_name"`
	EnName string `json:"en_name"`
}

// 待缓存到 redis的有效 token数据
type TokenToRedis struct {
	Id        int    `json:"id"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}
