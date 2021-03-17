package button

type BaseField struct {
	EnName        string   `form:"en_name" json:"en_name" binding:"required,min=2"`
	CnName        string   `form:"cn_name" json:"cn_name" binding:"required,min=2"`
	Color         string   `form:"color" json:"color" binding:"required,min=2"`
	RequestMethod string   `form:"allow_method" json:"allow_method" binding:"required,min=1"`
	Status        *float64 `form:"status" json:"status" binding:"required,min=0"`
	Remark        string   `form:"remark" json:"remark" `
}

type Id struct {
	Id float64 `form:"id" json:"id" binding:"required,min=1"` // 注意：gin框架数字的存储形式都是 float64
}
