package post_members

type Create struct {
	OrgPostId float64  `form:"org_post_id" json:"org_post_id" binding:"required,min=1"`
	UserId    float64  `form:"user_id" json:"user_id" binding:"required,min=1"`
	Status    *float64 `form:"status" json:"status" binding:"required,min=0"`
	Remark    string   `form:"remark" json:"remark"`
}

type Id struct {
	Id float64 `form:"id" json:"id" binding:"required,min=1"`
}
