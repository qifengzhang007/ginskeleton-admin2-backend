package users

type BaseField struct {
	UserName string   `form:"user_name" binding:"required,min=2"`
	Pass     string   `form:"pass" binding:"required,min=6"`
	RealName string   `form:"real_name" binding:"required,min=2"`
	Avatar   string   `form:"avatar"`
	Phone    string   `form:"phone"`
	Status   *float64 `form:"status" binding:"required,min=0"`
	Remark   string   `form:"remark" `
}

type Id struct {
	Id float64 `form:"id" binding:"required,min=1"` // 注意：gin框架数字的存储形式都是 float64
}

type MenuId struct {
	MenuId float64 `form:"menu_id" binding:"required,min=1"` // 注意：gin框架数字的存储形式都是 float64
}

type UserName struct {
	UserName string `form:"user_name"`
}

type Pass struct {
	Pass string `form:"pass" binding:"required,min=6,max=20"` //  密码为 必填，长度>=6
}
