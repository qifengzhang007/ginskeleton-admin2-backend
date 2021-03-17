package org_post

type Create struct {
	Fid    *float64 `form:"fid" json:"fid" binding:"required,min=0"`
	Title  string   `form:"title" json:"title" binding:"required,min=1,max=120"`
	Status *float64 `form:"status" json:"status" binding:"required,min=0,max=1"`
	Remark string   `form:"remark" json:"remark" `
}

type Id struct {
	Id float64 `form:"id" json:"id" binding:"required,min=1"`
}

type Fid struct {
	Fid *float64 `form:"fid" json:"fid" binding:"min=0"` // 必填，页面值>0
}
