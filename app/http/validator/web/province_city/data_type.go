package province_city

type BaseField struct {
	Name   string   `form:"name" json:"name" binding:"required,min=1"`
	Status *float64 `form:"status" json:"status" binding:"required,min=0"`
	Sort   *float64 `form:"sort" json:"sort" binding:"required,min=0"`
	Remark string   `form:"remark" json:"remark" `
}

type Id struct {
	Id float64 `form:"id" json:"id" binding:"required,min=1"` // 注意：gin数字的存储形式以 float64 接受
}

type Ids struct {
	Ids string `form:"ids" json:"ids" binding:"required,min=1"` // 注意：gin数字的存储形式以 float64 接受
}

type Fid struct {
	Fid *float64 `form:"fid"  json:"fid"  binding:"required,min=0" ` // 注意：gin数字的存储形式以 float64 接受
}

type Name struct {
	Name string `form:"name" json:"name" `
}
