package auth_system_menu

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/qifengzhang007/sql_res_to_tree"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/validator/web/auth/system_menu/data_type"
	"goskeleton/app/model/auth"
	"time"
)

type AuthSystemMenuService struct {
}

// 待分配系统菜单、mmodel、按钮树形化
func (a *AuthSystemMenuService) SystemMenuButtonToTree2(sqlRes []auth.AuthSystemMenuButton) []MenuListTree {
	var dest = make([]MenuListTree, 0)
	if err := sql_res_to_tree.CreateSqlResFormatFactory().ScanToTreeData(sqlRes, &dest); err == nil {
		return dest
	} else {
		variable.ZapLog.Error("sql结果数据树形化失败，错误明细：", zap.Error(err))
	}
	return nil
}

// 已分配给系统菜单、按钮树形化
func (a *AuthSystemMenuService) AssginedMenuButtonToTree(sqlRes []auth.AssignedSystemMenuButton) []AssignedSystemMenuButton {
	var dest = make([]AssignedSystemMenuButton, 0)
	if err := sql_res_to_tree.CreateSqlResFormatFactory().ScanToTreeData(sqlRes, &dest); err == nil {
		return dest
	} else {
		variable.ZapLog.Error("sql结果数据树形化失败，错误明细：", zap.Error(err))
	}
	return nil
}

// 处理按钮字符串
func (a *AuthSystemMenuService) ButtonStringToArray(jsonStr string) []map[string]interface{} {
	mSlice := make([]map[string]interface{}, 0)
	_ = json.Unmarshal([]byte(jsonStr), &mSlice)
	return mSlice
}

// 将按钮循环加入表中
func (a *AuthSystemMenuService) InsertButton(context *gin.Context, menuId int64) bool {
	if menuButtonList, isOk := context.MustGet(variable.SystemCreateKey).(data_type.MenuCreate); isOk {
		for index, item := range menuButtonList.ButtonArray {
			item.FrAuthSystemMenuId = menuId
			item.Status = 1
			item.CreatedAt = time.Now().Format(variable.DateFormat)
			item.UpdatedAt = item.CreatedAt
			menuButtonList.ButtonArray[index] = item
		}
		if len(menuButtonList.ButtonArray) > 0 && auth.CreateAuthSystemMenuButtonFactory("").InsertData(menuButtonList.ButtonArray) {
			return true
		} else if len(menuButtonList.ButtonArray) == 0 {
			return true
		}
	}
	return false
}

// 更新按钮数据
func (a *AuthSystemMenuService) UpdateButton(context *gin.Context, menuId int64) bool {
	if menuButtonList, isOk := context.MustGet(variable.SystemEditKey).(data_type.MenuEdit); isOk {
		//修改数据过程中可能存在单条数据被删除的情况，首先删除已标记的数据
		if len(menuButtonList.ButtonDelete) > 0 {
			auth.CreateAuthSystemMenuButtonFactory("").BatchDeleteData(menuButtonList.ButtonDelete)
		}
		for index, item := range menuButtonList.ButtonArray {
			item.FrAuthSystemMenuId = menuId
			item.Status = 1
			item.CreatedAt = time.Now().Format(variable.DateFormat)
			item.UpdatedAt = item.CreatedAt
			menuButtonList.ButtonArray[index] = item
		}
		if auth.CreateAuthSystemMenuButtonFactory("").UpdateData(menuButtonList) {
			go a.UpdateHook(menuId)
			return true
		}
	}
	return false
}

// 菜单挂接的待分配权限按钮数据被更新后，需要自动更新tb_auth_casbin_rule表数据
func (a *AuthSystemMenuService) UpdateHook(menuId int64) {
	auth.CreateAuthSystemMenuButtonFactory("").UpdateHook(menuId)
}
