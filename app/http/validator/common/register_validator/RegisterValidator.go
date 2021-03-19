package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/validator/api/home"
	"goskeleton/app/http/validator/common/upload_files"
	"goskeleton/app/http/validator/common/websocket"
	"goskeleton/app/http/validator/web/auth/analysis"
	"goskeleton/app/http/validator/web/auth/button"
	"goskeleton/app/http/validator/web/auth/button_en_cn"
	"goskeleton/app/http/validator/web/auth/org_post"
	"goskeleton/app/http/validator/web/auth/post_members"
	"goskeleton/app/http/validator/web/auth/system_menu"
	"goskeleton/app/http/validator/web/users"
)

// 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func RegisterValidator() {
	//创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string

	// 注册门户类表单参数验证器
	key = consts.ValidatorPrefix + "HomeNews"
	containers.Set(key, home.News{})

	// Users 模块表单验证器按照 key => value 形式注册在容器，方便路由模块中调用

	//key = consts.ValidatorPrefix + "UsersRegister"  // 关闭通过公共接口注册用户
	//containers.Set(key, users.Register{})

	key = consts.ValidatorPrefix + "UsersLogin"
	containers.Set(key, users.Login{})
	key = consts.ValidatorPrefix + "RefreshToken"
	containers.Set(key, users.RefreshToken{})

	// Users基本操作（CURD）
	{
		key = consts.ValidatorPrefix + "UserList"
		containers.Set(key, users.List{})
		key = consts.ValidatorPrefix + "UserCreate"
		containers.Set(key, users.Create{})
		key = consts.ValidatorPrefix + "UserEdit"
		containers.Set(key, users.Edit{})
		key = consts.ValidatorPrefix + "UserDestroy"
		containers.Set(key, users.Destroy{})

		// 用户打开一个菜单对应的页面地址时，拥有权限的按钮列表
		key = consts.ValidatorPrefix + "ViewButtonList"
		containers.Set(key, users.ViewButtonList{})
	}
	// 文件上传
	key = consts.ValidatorPrefix + "UploadFiles"
	containers.Set(key, upload_files.UpFiles{})

	// Websocket 连接验证器
	key = consts.ValidatorPrefix + "WebsocketConnect"
	containers.Set(key, websocket.Connect{})

	//组织结构部分
	{
		key = consts.ValidatorPrefix + "OrgPostList"
		containers.Set(key, org_post.OrgPostList{})

		key = consts.ValidatorPrefix + "OrgPostCreate"
		containers.Set(key, org_post.OrgPostCreate{})

		key = consts.ValidatorPrefix + "OrgPostGetByFid"
		containers.Set(key, org_post.OrgPostGetByFid{})

		key = consts.ValidatorPrefix + "OrgPostEdit"
		containers.Set(key, org_post.OrgPostEdit{})

		key = consts.ValidatorPrefix + "OrgPostDestroy"
		containers.Set(key, org_post.OrgPostDestroy{})

	}
	// 按钮公共接口
	{
		key = consts.ValidatorPrefix + "ButtonCnEnList"
		containers.Set(key, button_en_cn.ButtonCnEnList{})
	}

	// 系统菜单相关的全部接口验证器
	{
		key = consts.ValidatorPrefix + "SystemMenuList"
		containers.Set(key, system_menu.SystemMenuList{})

		key = consts.ValidatorPrefix + "SystemMenuCreate"
		containers.Set(key, system_menu.SystemMenuCreate{})

		key = consts.ValidatorPrefix + "SystemMenuEdit"
		containers.Set(key, system_menu.SystemMenuEdit{})

		key = consts.ValidatorPrefix + "SystemMenuDestroy"
		containers.Set(key, system_menu.SystemMenuDestroy{})

		key = consts.ValidatorPrefix + "SystemMenuGetByFid"
		containers.Set(key, system_menu.SystemMenuGetByFid{})
		// 系统菜单拥有可被分配的全部按钮列表
		key = consts.ValidatorPrefix + "SysMenuMountButton"
		containers.Set(key, system_menu.SysMenuMountButton{})
	}

	// 权限分配
	{
		key = consts.ValidatorPrefix + "SystemMenuListAllList"
		containers.Set(key, system_menu.SystemMenuListAllList{})

		key = consts.ValidatorPrefix + "AssginedSystemMenuList"
		containers.Set(key, system_menu.AssginedSystemMenuList{})

		key = consts.ValidatorPrefix + "AssginSystemMenuToOrg"
		containers.Set(key, system_menu.AssginSystemMenuToOrg{})

		key = consts.ValidatorPrefix + "DelAuthFromOrg"
		containers.Set(key, system_menu.DelAuthFromOrg{})
	}

	// 权限分析（左侧带岗位信息的用户列表+右侧用户权限来源）
	{
		key = consts.ValidatorPrefix + "UserListWithPost"
		containers.Set(key, analysis.UserListWithPost{})
		key = consts.ValidatorPrefix + "OrgPostGetByUserId"
		containers.Set(key, analysis.OrgPostGetByUserId{})
	}
	//6.岗位成员表
	{
		key = consts.ValidatorPrefix + "PostMembersList"
		containers.Set(key, post_members.PostMembersList{})
		key = consts.ValidatorPrefix + "PostMembersCreate"
		containers.Set(key, post_members.PostMembersCreate{})
		key = consts.ValidatorPrefix + "PostMembersEdit"
		containers.Set(key, post_members.PostMembersEdit{})
		key = consts.ValidatorPrefix + "PostMembersDestroy"
		containers.Set(key, post_members.PostMembersDestroy{})
	}

	//按钮部分
	{
		key = consts.ValidatorPrefix + "ButtonList"
		containers.Set(key, button.ButtonList{})
		key = consts.ValidatorPrefix + "ButtonCreate"
		containers.Set(key, button.ButtonCreate{})
		key = consts.ValidatorPrefix + "ButtonEdit"
		containers.Set(key, button.ButtonEdit{})
		key = consts.ValidatorPrefix + "ButtonDestroy"
		containers.Set(key, button.ButtonDestroy{})
	}
}
