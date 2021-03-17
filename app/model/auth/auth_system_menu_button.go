package auth

import (
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

func CreateAuthSystemMenuButtonFactory(sqlType string) *AuthSystemMenuButtonModel {
	return &AuthSystemMenuButtonModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type AuthSystemMenuButtonModel struct {
	model.BaseModel
	FrAuthSystemMenuId int    `json:"fr_auth_system_menu_id"`
	FrAuthButtonCnEnId int    `json:"fr_auth_button_cn_en_id"`
	RequestMethod      string `json:"request_method"`
	RequestUrl         string `json:"request_url"`
	Status             int    `gorm:"-" json:"status"`
	Remark             string `json:"remark"`
}

// 表名
func (a *AuthSystemMenuButtonModel) TableName() string {
	return "tb_auth_system_menu_button"
}

func (a *AuthSystemMenuButtonModel) getCounts(sysMenuId int) (count int64) {
	if res := a.Debug().Model(a).Where("fr_auth_system_menu_id=?", sysMenuId).Count(&count); res.Error == nil {
		return count
	}
	return 0
}

// 查询
func (a *AuthSystemMenuButtonModel) List(sysMenuId int) (counts int, data []SystemMenuButtonList) {
	if a.getCounts(sysMenuId) > 0 {
		sql := `
		SELECT
		a.id,a.fr_auth_system_menu_id,
		a.fr_auth_button_cn_en_id,
		b.cn_name AS button_name,
		a.request_method,  a.request_url,
		a.status,a.remark
		FROM
		tb_auth_system_menu_button  a LEFT JOIN  tb_auth_button_cn_en  b  ON  a.fr_auth_button_cn_en_id=b.id
		WHERE   a.fr_auth_system_menu_id=?
		`
		if res := a.Raw(sql, sysMenuId).Find(&data); res.Error != nil {
			variable.ZapLog.Error("AuthSystemMenuButtonModel 查询出错：" + res.Error.Error())
		}
	} else {
		return 0, nil
	}
	return
}

//新增
func (a *AuthSystemMenuButtonModel) InsertData(data AuthSystemMenuButtonModel) bool {
	a.Create(&data)
	return true
}

// 更新
func (a *AuthSystemMenuButtonModel) UpdateData(data AuthSystemMenuButtonModel) bool {
	a.Updates(&data)
	return true
}

// 删除数据
func (a *AuthSystemMenuButtonModel) DeleteData(id int) bool {
	if res := a.Delete(a, id); res.Error == nil {
		return true
	} else {
		variable.ZapLog.Error("AuthSystemMenuButtonModel 数据删除失败", zap.Error(res.Error))
	}
	return false
}

//新增
func (a *AuthSystemMenuButtonModel) InsertMap(data map[string]interface{}) bool {
	a.Model(a.TableName()).Create(&data)
	return true
}

//根据菜单ID获取按钮信息
func (a *AuthSystemMenuButtonModel) MenuButton(menuId float64) (data []SystemMenuButtonList) {
	m := a.Table("tb_auth_system_menu_button a")
	m.Joins("left join tb_auth_button_cn_en b on a.fr_auth_button_cn_en_id = b.id").
		Where("a.fr_auth_system_menu_id = ?", menuId).
		Select("a.*,b.cn_name as button_name").Scan(&data)
	return
}

// 数据更新hook函数，负责更新菜单被引用的地方，同步更新
func (a *AuthSystemMenuButtonModel) UpdateHook(menuId int64) {
	// 更新菜单挂接的按钮之后，可能存在按钮被删除，因此需要删除的数据主要有：1. tb_auth_casbin_rules 表被应用的按钮数据
	sql := `
		DELETE   FROM  tb_auth_casbin_rule WHERE   fr_auth_post_mount_has_menu_button_id  IN(
			SELECT  b.id  FROM tb_auth_post_mount_has_menu  a ,tb_auth_post_mount_has_menu_button  b
			WHERE  a.id=b.fr_auth_post_mount_has_menu_id
			AND a.fr_auth_system_menu_id=?
			AND b.fr_auth_button_cn_en_id  NOT   IN(
				SELECT d.fr_auth_button_cn_en_id  FROM  
				tb_auth_system_menu  c,tb_auth_system_menu_button  d   
				WHERE c.id=d.fr_auth_system_menu_id
				AND c.id=?
			)
		)	
	`
	if res := a.Exec(sql, menuId, menuId); res.Error != nil {
		variable.ZapLog.Error("AuthSystemMenuButtonModel UpdateHook 删 除tb_auth_casbin_rule 关联按钮数据出错", zap.Error(res.Error))
	}

	sql = `
		DELETE   FROM tb_auth_post_mount_has_menu_button   WHERE    id   IN(
		SELECT id  FROM (	SELECT  b.id  FROM tb_auth_post_mount_has_menu  a ,tb_auth_post_mount_has_menu_button  b
			WHERE  a.id=b.fr_auth_post_mount_has_menu_id
			AND a.fr_auth_system_menu_id=?
			AND b.fr_auth_button_cn_en_id  NOT   IN(
				SELECT d.fr_auth_button_cn_en_id  FROM  
				tb_auth_system_menu  c,tb_auth_system_menu_button  d   
				WHERE c.id=d.fr_auth_system_menu_id
				AND c.id=?
			)) AS  tmp
		)	
	`
	if res := a.Exec(sql, menuId, menuId); res.Error != nil {
		variable.ZapLog.Error("AuthSystemMenuButtonModel UpdateHook 删 tb_auth_post_mount_has_menu_button 关联按钮数据出错", zap.Error(res.Error))
	}

	// 批量更新菜单被引用的所有地方
	sql = `
		UPDATE  tb_auth_casbin_rule e  LEFT  JOIN  (
		SELECT  
		DISTINCT
		a.id,b.fr_auth_button_cn_en_id,b.request_method,b.request_url ,d.id AS   auth_post_mount_has_menu_button_id
		FROM  tb_auth_system_menu  a,tb_auth_system_menu_button  b   ,
		tb_auth_post_mount_has_menu  c,tb_auth_post_mount_has_menu_button  d
		WHERE   a.id=b.fr_auth_system_menu_id
		AND c.id=d.fr_auth_post_mount_has_menu_id
		AND c.fr_auth_system_menu_id = a.id  AND   d.fr_auth_button_cn_en_id=b.fr_auth_button_cn_en_id
		AND a.id=?  
		)  AS   f  ON  e.fr_auth_post_mount_has_menu_button_id=f.auth_post_mount_has_menu_button_id
		SET  e.v1=f.request_url  ,  e.v2=f.request_method 
		WHERE   f.auth_post_mount_has_menu_button_id  IS  NOT  NULL 
		AND  LENGTH(IFNULL(f.request_url,''))>0 
		AND   LENGTH(IFNULL(f.request_method,''))>0 
		`
	if res := a.Exec(sql, menuId); res.Error != nil {
		variable.ZapLog.Error("AuthSystemMenuButtonModel UpdateHook 更新 tb_auth_casbin_rule 出错", zap.Error(res.Error))
	}
}

func (a *AuthSystemMenuButtonModel) GetByButtonId(butonId int) bool {
	data := []AuthSystemMenuButtonModel{}
	a.Where("fr_auth_button_cn_en_id = ?", butonId).Find(&data)
	if len(data) != 0 {
		return false
	}
	return true
}
