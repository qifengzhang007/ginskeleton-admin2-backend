package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/data_bind"
)

func CreateAuthOrganizationFactory(sqlType string) *AuthOrganizationPostModel {
	return &AuthOrganizationPostModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type AuthOrganizationPostModel struct {
	model.BaseModel
	Fid      int    `json:"fid"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	PathInfo string `json:"path_info"`
	Remark   string `json:"remark"`
}

// 表名
func (a *AuthOrganizationPostModel) TableName() string {
	return "tb_auth_organization_post"
}

func (a *AuthOrganizationPostModel) GetCount(fid int, title string) (count int64) {
	a.Model(a).Select("id").Where("fid = ? AND title like ?", fid, "%"+title+"%").Count(&count)
	return
}
func (a *AuthOrganizationPostModel) List(limitStart int, limit int, fid int, title string) (data []AuthOrganizationPostModel) {
	a.Model(a).Select("id", "fid", "title", "status", "path_info", "remark", "created_at", "updated_at").Where("fid = ? AND title like ?", fid, "%"+title+"%").Offset(limitStart).Limit(limit).Find(&data)
	return
}

func (a *AuthOrganizationPostModel) InsertData(c *gin.Context) bool {
	var tmp AuthOrganizationPostModel
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		var counts int64
		if res := a.Model(a).Where("fid=? and title=?", tmp.Fid, tmp.Title).Count(&counts); res.Error == nil && counts > 0 {
			return false
		} else {
			if res := a.Create(&tmp); res.Error == nil {
				_ = a.updatePathInfoNodeLevel(int(tmp.Id))
				return true
			} else {
				variable.ZapLog.Error("AuthOrganizationPostModel 数据新增出错：", zap.Error(res.Error))
			}
		}
	} else {
		variable.ZapLog.Error("AuthOrganizationPostModel 数据从验证器绑定到model出错：", zap.Error(err))
	}
	return false
}

// 更新path_info 、node_level 字段
func (a *AuthOrganizationPostModel) updatePathInfoNodeLevel(curItemid int) bool {
	sql := `
		UPDATE tb_auth_organization_post a  LEFT JOIN tb_auth_organization_post  b
		ON  a.fid=b.id
		SET  a.node_level=b.node_level+1,  a.path_info=CONCAT(b.path_info,',',a.id)
		WHERE  a.id=?
		`
	if res := a.Exec(sql, curItemid); res.Error == nil && res.RowsAffected >= 0 {
		return true
	} else {
		variable.ZapLog.Error("tb_auth_organization_post 更新 path_info 失败", zap.Error(res.Error))
	}
	return false
}

func (a *AuthOrganizationPostModel) GetByFid(fid int, data *[]AuthOrganizationPostTree) (err error) {
	sql := `
		SELECT  
		id,fid,title, STATUS,path_info,remark ,
		(SELECT  CASE  WHEN  COUNT(*) >0 THEN 1 ELSE  0 END  FROM tb_auth_organization_post  WHERE  fid=a.id ) AS  has_sub_node
		FROM   tb_auth_organization_post  a  WHERE  fid=?
	`
	err = a.Raw(sql, fid).Scan(data).Error
	return
}

//根据ID查询单挑数据
func (a *AuthOrganizationPostModel) GetById(id int, data *AuthOrganizationPostModel) (err error) {
	err = a.Model(a).Where("id = ?", id).Find(data).Error
	return
}

func (a *AuthOrganizationPostModel) UpdateData(c *gin.Context) bool {
	var tmp AuthOrganizationPostModel

	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		// Omit 表示忽略指定字段(CreatedAt)，其他字段全量更新
		if res := a.Omit("CreatedAt").Save(&tmp); res.Error == nil {
			_ = a.updatePathInfoNodeLevel(int(tmp.Id))
		}
		return true
	} else {
		variable.ZapLog.Error("AuthOrganizationPostModel 数据更新失败，错误详情：", zap.Error(err))
	}
	return false

}

func (a *AuthOrganizationPostModel) DeleteData(id int) bool {
	if res := a.Delete(a, id); res.Error == nil {
		return true
	} else {
		variable.ZapLog.Error("AuthOrganizationPostModel 删除数据出错：", zap.Error(res.Error))
	}
	return false
}
func (a *AuthOrganizationPostModel) HasSubList(id int) (count int64) {
	a.Model(a).Select("id").Where("fid=?", id).Count(&count)
	return
}

func (a *AuthOrganizationPostModel) GetByIds(ids []int) (AuthOrganizationPostModel []AuthOrganizationPostModel) {
	a.Where("id IN ?", ids).Find(&AuthOrganizationPostModel)
	return
}

func (a *AuthOrganizationPostModel) GetByIdsScan(ids []int) (AllAuth []AllAuth) {
	a.Model(a).Select("id", "title", "fid").Where("id IN ?", ids).Scan(&AllAuth)
	return
}
