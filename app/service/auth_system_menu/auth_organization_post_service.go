package auth_system_menu

import (
	"encoding/json"
	"github.com/qifengzhang007/sql_res_to_tree"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/auth"
	"strconv"
	"strings"
)

type AuthSystemMenuService struct {
}

func (a *AuthSystemMenuService) GetOrgByFid(fid int) (err error, data []auth.AuthSystemMenuTree) {
	models := auth.CreateAuthSystemMenuFactory("")
	err = models.GetByFid(fid, &data)
	for key, value := range data {
		var hasSubNode []auth.AuthSystemMenuTree
		_ = models.GetByFid(value.Id, &hasSubNode)
		if len(hasSubNode) > 0 {
			value.Children = []auth.AuthSystemMenuTree{}
			data[key] = value
		} else {
			value.Children = nil
			data[key] = value
		}
	}
	return
}

// 待分配系统菜单、mmodel、按钮树形化
func (a *AuthSystemMenuService) SystemMenuButtonToTree(sqlRes []auth.AuthSystemMenuButton) []MenuListTree {
	var dest = make([]MenuListTree, 0)
	if err := sql_res_to_tree.CreateSqlResFormatFactory().ScanToTreeData(sqlRes, &dest); err == nil {
		return dest
	} else {
		variable.ZapLog.Error("sql结果数据树形化失败，错误明细：", zap.Error(err))
	}
	return nil
}

// 已分配给系统菜单、按钮树形化
func (a *AuthSystemMenuService) AssginedMenuButtonToTree(sqlRes []auth.AssignedSystemMenuButton) []AssignedMenuListTree {
	var dest = make([]AssignedMenuListTree, 0)
	if err := sql_res_to_tree.CreateSqlResFormatFactory().ScanToTreeData(sqlRes, &dest); err == nil {
		return dest
	} else {
		variable.ZapLog.Error("sql结果数据树形化失败，错误明细：", zap.Error(err))
	}
	return nil
}

//处理按钮字符串
func (a *AuthSystemMenuService) ButtonStringToArray(jsonStr string) []map[string]interface{} {
	mSlice := make([]map[string]interface{}, 0)
	_ = json.Unmarshal([]byte(jsonStr), &mSlice)
	return mSlice
}

//讲按钮循环加入表中
//处理按钮字符串
func (a *AuthSystemMenuService) InsertButton(jsonStr string, menuId int64) bool {
	buttonArray := a.ButtonStringToArray(jsonStr)
	for _, v := range buttonArray {
		data := auth.AuthSystemMenuButtonModel{}
		data.FrAuthButtonCnEnId = int(v["fr_auth_button_cn_en_id"].(float64))
		data.FrAuthSystemMenuId = int(menuId)
		data.Remark = v["remark"].(string)
		data.RequestMethod = v["request_method"].(string)
		data.RequestUrl = v["request_url"].(string)
		auth.CreateAuthSystemMenuButtonFactory("").InsertData(data)
	}
	return true
}

//讲按钮循环加入表中
//处理按钮字符串
func (a *AuthSystemMenuService) UpdateButton(jsonStr string, buttonDelete string, menuId int64) bool {
	buttonArray := a.ButtonStringToArray(jsonStr)
	//取出需要删除的项目
	buttonDeleteArr := strings.Split(buttonDelete, ",")
	//循环删除
	if len(buttonDeleteArr) != 0 {
		for _, v := range buttonDeleteArr {
			id, _ := strconv.Atoi(v)
			auth.CreateAuthSystemMenuButtonFactory("").DeleteData(id)
		}
	}
	for _, v := range buttonArray {
		data := auth.AuthSystemMenuButtonModel{}
		data.FrAuthButtonCnEnId = int(v["fr_auth_button_cn_en_id"].(float64))
		data.FrAuthSystemMenuId = int(menuId)
		data.Remark = v["remark"].(string)
		data.Id = int64(v["id"].(float64))
		data.RequestMethod = v["request_method"].(string)
		data.RequestUrl = v["request_url"].(string)
		//开始逻辑判断,如果button_id = 0 ,则新增
		if data.Id == 0 {
			auth.CreateAuthSystemMenuButtonFactory("").InsertData(data)
		} else {
			auth.CreateAuthSystemMenuButtonFactory("").UpdateData(data)
		}

	}
	go a.UpdateHook(menuId)
	return true
}

// 菜单挂接的待分配权限按钮数据被更新后，需要自动更新tb_auth_casbin_rule表数据
func (a *AuthSystemMenuService) UpdateHook(menuId int64) {
	auth.CreateAuthSystemMenuButtonFactory("").UpdateHook(menuId)
}
