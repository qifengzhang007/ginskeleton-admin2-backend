package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/data_bind"
)

func CreateAuthPostMembersModelFactory(sqlType string) *AuthPostMembersModel {
	return &AuthPostMembersModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type AuthPostMembersModel struct {
	model.BaseModel
	FrAuthOrganizationPostId int    `gorm:"column:fr_auth_organization_post_id" json:"org_post_id"`
	FrUserId                 int    `gorm:"column:fr_user_id" json:"user_id"`
	Status                   int    `json:"status"`
	Remark                   string `json:"remark"`
}

// 表名
func (a *AuthPostMembersModel) TableName() string {
	return "tb_auth_post_members"
}

// 查询类
func (a *AuthPostMembersModel) GetCount(postId float64, userName string) (count int64) {
	sql := `SELECT  count(*) as  counts  
			FROM    tb_auth_post_members  a,tb_users  b   
			WHERE a.fr_user_id=b.id  AND   ( a.fr_auth_organization_post_id=?  or 0=?)
			AND ( b.user_name  LIKE ?  or  b.real_name  like  ?)
			`
	a.Raw(sql, postId, postId, "%"+userName+"%", "%"+userName+"%").First(&count)
	return
}
func (a *AuthPostMembersModel) List(postId, limitStart, limits int, userName string) (data []MemberList) {
	sql := `SELECT  a.id, a.fr_auth_organization_post_id AS org_post_id, a.fr_user_id AS user_id,b.user_name,b.real_name,a.status,c.title AS post_name, a.remark, 
			DATE_FORMAT(a.created_at,'%Y-%m-%d %h:%i:%s')  AS created_at, DATE_FORMAT(a.updated_at,'%Y-%m-%d %h:%i:%s')  AS updated_at  
			FROM    tb_auth_post_members  a,tb_users  b   ,tb_auth_organization_post  c
			WHERE  a.fr_user_id=b.id  AND c.id=a.fr_auth_organization_post_id
			AND    (a.fr_auth_organization_post_id=? or 0=?)
			AND  ( b.user_name  LIKE ?  or  b.real_name  like  ?)
			LIMIT   ?,?`
	a.Raw(sql, postId, postId, "%"+userName+"%", "%"+userName+"%", limitStart, limits).Find(&data)
	return
}

// 新增
func (a *AuthPostMembersModel) InsertData(c *gin.Context) bool {
	var tmp AuthPostMembersModel
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		var counts int64
		if res := a.Model(a).Where("fr_auth_organization_post_id=? AND  fr_user_id=?", tmp.FrAuthOrganizationPostId, tmp.FrUserId).Count(&counts); res.Error == nil && counts == 0 {
			if res := a.Create(&tmp); res.Error == nil {
				return true
			} else {
				variable.ZapLog.Error("AuthPostMembersModel 新增失败", zap.Error(res.Error))
			}
		} else {
			variable.ZapLog.Warn("AuthPostMembersModel 不允许重复新增")
		}
	} else {
		variable.ZapLog.Warn("AuthPostMembersModel 数据绑定出错", zap.Error(err))
	}
	return false
}

//修改
func (a *AuthPostMembersModel) UpdateData(c *gin.Context) bool {
	var tmp AuthPostMembersModel
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		// Omit 表示忽略指定字段(CreatedAt)，其他字段全量更新
		if res := a.Omit("CreatedAt").Save(&tmp); res.Error == nil {
			return true
		} else {
			variable.ZapLog.Error("AuthPostMembersModel 数据更新出错：", zap.Error(res.Error))
		}
	} else {
		variable.ZapLog.Error("AuthPostMembersModel 数据绑定出错：", zap.Error(err))
	}
	return false
}

//删除
func (a *AuthPostMembersModel) DeleteData(id float64) bool {
	// 只能删除除了 admin 之外的用户
	var count int64
	a.Model(a).Select("fr_user_id").Where("id=?", id).First(&count)
	if count == 1 {
		return false
	}
	if res := a.Delete(a, id); res.Error == nil {
		return true
	} else {
		variable.ZapLog.Error("AuthPostMembersModel 删除数据出错：", zap.Error(res.Error))
	}
	return false
}

//修改
func (a *AuthPostMembersModel) GetByUserId(user_id int64) (AuthPostMembersModel []AuthPostMembersModel) {
	a.Where("fr_user_id = ?", user_id).Find(&AuthPostMembersModel)
	return
}
