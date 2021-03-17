package auth

import "goskeleton/app/model"

func CreateAuthPostMountHasMenuModelFactory(sqlType string) *AuthPostMountHasMenuModel {
	return &AuthPostMountHasMenuModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type AuthPostMountHasMenuModel struct {
	model.BaseModel         `gorm:"-"`
	FrAuthOrgnizationPostId int    `form:"fr_auth_orgnization_post_id" gorm:"column:fr_auth_orgnization_post_id" json:"fr_auth_orgnization_post_id"`
	FrAuthSystemMenuId      int    `form:"fr_auth_system_menu_id" gorm:"column:fr_auth_system_menu_id" json:"fr_auth_system_menu_id"`
	Status                  int    `form:"status" json:"status"`
	Remark                  string `form:"remark" json:"remark"`
}

// 表名
func (a *AuthPostMountHasMenuModel) TableName() string {
	return "tb_auth_post_mount_has_menu"
}

// 根据id获取菜单id
func (a *AuthPostMountHasMenuModel) GetByIds(ids []int) (AuthPostMountHasMenuModel []AuthPostMountHasMenuModel) {
	a.Where("fr_auth_orgnization_post_id IN ?", ids).Select("distinct fr_auth_system_menu_id").Find(&AuthPostMountHasMenuModel)
	return
}

// 根据postID获取菜单ID
func (a *AuthPostMountHasMenuModel) GetByPostId(id int) (AuthPostMountHasMenuModel []AuthPostMountHasMenuModel) {
	a.Where("fr_auth_orgnization_post_id = ?", id).Select("distinct fr_auth_system_menu_id").Find(&AuthPostMountHasMenuModel)
	return
}
