package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/qifengzhang007/sql_res_to_tree"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/data_bind"
)

func CreateAuthSystemMenuFactory(sqlType string) *AuthSystemMenuModel {
	return &AuthSystemMenuModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type AuthSystemMenuModel struct {
	model.BaseModel
	Fid       int    `json:"fid"`
	Icon      string `json:"icon"`
	Title     string `json:"title"`
	Name      string `json:"name"`
	Loading   bool   `gorm:"-" json:"loading"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Status    int    `json:"status"`
	Sort      int    `json:"sort"`
	Remark    string `json:"remark"`
}

// 表名
func (a *AuthSystemMenuModel) TableName() string {
	return "tb_auth_system_menu"
}

func (a *AuthSystemMenuModel) getCounts(fid int, title string) (count int) {
	sql := "SELECT COUNT(*) AS counts FROM   `tb_auth_system_menu`   WHERE fid=? AND title  LIKE  ?"
	a.Raw(sql, fid, "%"+title+"%").First(&count)
	return
}

// 查询
func (a *AuthSystemMenuModel) List(limitStart int, limit int, fid int, title string) (counts int, data []AuthSystemMenuButtonListTree) {
	counts = a.getCounts(fid, title)
	if counts > 0 {
		sql := `
			SELECT   a.id,  a.fid,  a.icon,  a.title,  a.name,  a.path,  a.component, a.status,a.remark ,a.sort,
			IFNULL(b.fr_auth_system_menu_id ,0) AS fr_auth_system_menu_id,
			IFNULL(c.id,0)  AS button_id,IFNULL(c.cn_name,'')  AS   button_name, IFNULL(c.color,'') AS  button_color 
			FROM   tb_auth_system_menu a  LEFT JOIN  tb_auth_system_menu_button  b ON  a.id =b.fr_auth_system_menu_id
			LEFT  JOIN  tb_auth_button_cn_en  c  ON  c.id=b.fr_auth_button_cn_en_id
			WHERE a.id  IN(
				SELECT id   FROM (SELECT id FROM   tb_auth_system_menu   WHERE fid=? AND title  LIKE  ?  LIMIT ?,?) AS tb_tmp 
			)
			ORDER   BY   a.sort  DESC,a.fid  ASC,button_id ASC
		`
		var sqlSlice []AuthSystemMenuButtonList
		if res := a.Raw(sql, fid, "%"+title+"%", limitStart, limit).Find(&sqlSlice); res.Error == nil {
			var dest = make([]AuthSystemMenuButtonListTree, 0)
			if err := sql_res_to_tree.CreateSqlResFormatFactory().ScanToTreeData(sqlSlice, &dest); err == nil {
				return counts, dest
			} else {
				variable.ZapLog.Error("AuthSystemMenuModel 树形化出错:" + err.Error())

			}
		}
	}
	return 0, nil
}

// 通过fid查询子节点数据
func (a *AuthSystemMenuModel) GetByFid(fid int) (data []AuthSystemMenuTree, err error) {
	sql := `
		SELECT  
		id,fid,title,name, icon,path,component,remark ,
		(SELECT  CASE  WHEN  COUNT(*) >0 THEN 1 ELSE  0 END  FROM tb_auth_system_menu  WHERE  fid=a.id ) AS  has_sub_node,
		(SELECT  CASE  WHEN  COUNT(*) =0 THEN 1 ELSE  0 END  FROM tb_auth_system_menu  WHERE  fid=a.id ) AS  is_leaf
		FROM   tb_auth_system_menu  a  WHERE  fid=?
	`
	err = a.Raw(sql, fid).Scan(&data).Error
	return
}

//新增
func (a *AuthSystemMenuModel) InsertData(c *gin.Context) (bool, AuthSystemMenuModel) {
	var tmp AuthSystemMenuModel
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		var counts int64
		if res := a.Model(a).Where("fid=? AND  title=?", tmp.Fid, tmp.Title).Count(&counts); res.Error == nil && counts == 0 {
			if res := a.Create(&tmp); res.Error == nil {
				//新增菜单后,处理按钮
				go a.updatePathInfoNodeLevel(tmp.Id)
				return true, tmp
			} else {
				variable.ZapLog.Error("AuthSystemMenuModel 新增失败", zap.Error(res.Error))
			}
		} else {
			variable.ZapLog.Warn("AuthSystemMenuModel 不允许重复新增")
		}
	} else {
		variable.ZapLog.Warn("AuthSystemMenuModel 数据绑定出错", zap.Error(err))
	}
	return false, tmp
}

// 更新
func (a *AuthSystemMenuModel) UpdateData(c *gin.Context) bool {
	var tmp AuthSystemMenuModel
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		// Omit 表示忽略指定字段(CreatedAt)，其他字段全量更新
		if res := a.Omit("CreatedAt").Save(&tmp); res.Error == nil {
			go a.updatePathInfoNodeLevel(tmp.Id)
			return true
		} else {
			variable.ZapLog.Error("AuthSystemMenuModel 数据更新出错：", zap.Error(res.Error))
		}
	} else {
		variable.ZapLog.Error("AuthSystemMenuModel 更新失败（或无数据被更新）", zap.Error(err))
	}
	return false
}

// 新增、更新继续hook，更新path_info 、node_level 字段
func (a *AuthSystemMenuModel) updatePathInfoNodeLevel(curItemid int64) bool {
	sql := `
		UPDATE tb_auth_system_menu a  LEFT JOIN tb_auth_system_menu  b
		ON  a.fid=b.id
		SET  a.node_level=IFNULL(b.node_level,0)+1,  a.path_info=CONCAT(IFNULL(b.path_info,0),',',a.id)
		WHERE  a.id=?
		`
	if res := a.Exec(sql, curItemid); res.Error == nil && res.RowsAffected >= 0 {
		return true
	} else {
		variable.ZapLog.Error("tb_auth_system_menu 更新path_info失败", zap.Error(res.Error))
	}
	return false
}

//根据id查询是否有子节点数据
func (a *AuthSystemMenuModel) GetSubNodeCount(id int) (count int64) {
	if res := a.Model(a).Where("fid = ?", id).Count(&count); res.Error != nil {
		variable.ZapLog.Error("AuthSystemMenuModel 查询子节点是否有数据出错：", zap.Error(res.Error))
	}
	return count
}

// 删除数据
func (a *AuthSystemMenuModel) DeleteData(id int) bool {
	if res := a.Delete(a, id); res.Error == nil {
		go a.DeleteDataHook(id) // 删除菜单关联的所有数据
		return true
	} else {
		variable.ZapLog.Error("AuthSystemMenuModel 数据删除失败", zap.Error(res.Error))
	}
	return false
}

//根据IDS获取菜单信息
func (a *AuthSystemMenuModel) GetByIds(ids []int) (AuthSystemMenuTree []AuthSystemMenuTree) {
	sql := `
			SELECT a.id ,a.fid, a.title, a.name, a.icon, a.name as  path, a.node_level,a.component ,
			IFNULL((SELECT 1 FROM tb_auth_system_menu b WHERE  b.fid=a.id  LIMIT 1),0) as has_sub_node
			FROM tb_auth_system_menu a  WHERE id IN (?) AND a.status=1 
			ORDER BY a.sort desc
		`
	a.Raw(sql, ids).Scan(&AuthSystemMenuTree)
	return
}

// 菜单主表数据删除，菜单关联的业务数据表同步删除
func (a *AuthSystemMenuModel) DeleteDataHook(menuId int) {

	//1.菜单可能被分配给  tb_auth_casbin_rules 的权限
	sql := `
		DELETE    FROM  tb_auth_casbin_rule  WHERE   fr_auth_post_mount_has_menu_button_id  IN(
			SELECT   id   FROM  tb_auth_post_mount_has_menu_button  WHERE  fr_auth_post_mount_has_menu_id  IN(
				SELECT   id   FROM  tb_auth_post_mount_has_menu  WHERE   fr_auth_system_menu_id=?
			)
		)
		`
	if res := a.Exec(sql, menuId); res.Error != nil {
		variable.ZapLog.Error("AuthSystemMenuModel 删除 tb_auth_casbin_rule 失败", zap.Error(res.Error))
	}

	//2. 菜单可能被分配给组织机构的权限关联数据
	sql = `
		DELETE FROM  tb_auth_post_mount_has_menu_button  WHERE  fr_auth_post_mount_has_menu_id  IN(
			SELECT   id   FROM  tb_auth_post_mount_has_menu  WHERE   fr_auth_system_menu_id=?
		)
		`
	if res := a.Exec(sql, menuId); res.Error != nil {
		variable.ZapLog.Error("AuthSystemMenuModel 删除 tb_auth_post_mount_has_menu_button 失败", zap.Error(res.Error))
	}
	//3. 菜单可能被分配给组织机构的权限按钮数据
	sql = `
		DELETE FROM tb_auth_post_mount_has_menu  WHERE  fr_auth_system_menu_id=?
		`
	if res := a.Exec(sql, menuId); res.Error != nil {
		variable.ZapLog.Error("AuthSystemMenuModel 删除 tb_auth_post_mount_has_menu 失败", zap.Error(res.Error))
	}

	//4. 删除菜单关联的待分配按钮子表
	sql = `DELETE  FROM  tb_auth_system_menu_button  WHERE   fr_auth_system_menu_id  = ? `
	if res := a.Exec(sql, menuId); res.Error != nil {
		variable.ZapLog.Error("AuthSystemMenuModel 删除 tb_auth_system_menu_button 失败", zap.Error(res.Error))
	}

}
