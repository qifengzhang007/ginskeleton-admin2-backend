package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/controller/captcha"
	"goskeleton/app/http/controller/web"
	"goskeleton/app/http/middleware/authorization"
	"goskeleton/app/http/middleware/cors"
	validatorFactory "goskeleton/app/http/validator/core/factory"
	"goskeleton/app/utils/gin_release"
	"net/http"
)

// 该路由主要设置 后台管理系统等后端应用路由

func InitWebRouter() *gin.Engine {
	var router *gin.Engine
	// 非调试模式（生产模式） 日志写到日志文件
	if variable.ConfigYml.GetBool("AppDebug") == false {
		//1.gin自行记录接口访问日志，不需要nginx，如果开启以下3行，那么请屏蔽第 34 行代码
		//gin.DisableConsoleColor()
		//f, _ := os.Create(variable.BasePath + variable.ConfigYml.GetString("Logs.GinLogName"))
		//gin.DefaultWriter = io.MultiWriter(f)

		//【生产模式】
		// 根据 gin 官方的说明：[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
		// 如果部署到生产环境，请使用以下模式：
		// 1.生产模式(release) 和开发模式的变化主要是禁用 gin 记录接口访问日志，
		// 2.go服务就必须使用nginx作为前置代理服务，这样也方便实现负载均衡
		// 3.如果程序发生 panic 等异常使用自定义的 panic 恢复中间件拦截、记录到日志
		router = gin_release.ReleaseRouter()
	} else {
		// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
		router = gin.Default()
		pprof.Register(router)
	}

	// 设置可信任的代理服务器列表,gin (2021-11-24发布的v1.7.7版本之后出的新功能)
	if variable.ConfigYml.GetInt("HttpServer.TrustProxies.IsOpen") == 1 {
		if err := router.SetTrustedProxies(variable.ConfigYml.GetStringSlice("HttpServer.TrustProxies.ProxyServerList")); err != nil {
			variable.ZapLog.Error(consts.GinSetTrustProxyError, zap.Error(err))
		}
	} else {
		_ = router.SetTrustedProxies(nil)
	}

	//根据配置进行设置跨域
	if variable.ConfigYml.GetBool("HttpServer.AllowCrossDomain") {
		router.Use(cors.Next())
	}

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "GinSkeleton-Admin-Backend")
	})

	//处理静态资源
	router.Static("/public", "./public") //  定义静态资源路由与实际目录映射关系

	// 创建一个验证码路由
	verifyCode := router.Group("captcha")
	{
		// 验证码业务，该业务无需专门校验参数，所以可以直接调用控制器
		verifyCode.GET("/", (&captcha.Captcha{}).GenerateId)                          //  获取验证码ID
		verifyCode.GET("/:captcha_id", (&captcha.Captcha{}).GetImg)                   // 获取图像地址
		verifyCode.GET("/:captcha_id/:captcha_value", (&captcha.Captcha{}).CheckCode) // 校验验证码
	}
	//  创建一个后端接口路由组
	backend := router.Group("/admin/")
	{
		// 创建一个websocket,如果ws需要账号密码登录才能使用，就写在需要鉴权的分组，这里暂定是开放式的，不需要严格鉴权，我们简单验证一下token值
		backend.GET("ws", validatorFactory.Create(consts.ValidatorPrefix+"WebsocketConnect"))

		//  【不需要token】中间件验证的路由  用户注册、登录
		noAuth := backend.Group("users/") //authorization.CheckCasbinAuth()
		{
			// 关于路由的第二个参数用法说明
			// 1.编写一个表单参数验证器结构体，参见代码：   app/http/validator/web/users/register.go
			// 2.将以上表单参数验证器注册，遵守 键 =》值 格式注册即可 ，app/http/validator/common/register_validator/register_validator.go  20行就是注册时候的键 consts.ValidatorPrefix+"UsersRegister"
			// 3.按照注册时的键，直接从容器调用即可 ：validatorFactory.Create(consts.ValidatorPrefix+"UsersRegister")
			//noAuth.POST("register", validatorFactory.Create(consts.ValidatorPrefix+"UsersRegister")) // 将公开注册渠道关闭
			noAuth.Use(authorization.CheckCaptchaAuth()).POST("login", validatorFactory.Create(consts.ValidatorPrefix+"UsersLogin"))
		}

		// 刷新token
		refreshToken := backend.Group("users/")
		{
			// 刷新token，当过期的token在允许失效的延长时间范围内，用旧token换取新token
			refreshToken.Use(authorization.RefreshTokenConditionCheck()).POST("refreshtoken", validatorFactory.Create(consts.ValidatorPrefix+"RefreshToken"))
		}

		// 【需要token】中间件验证的路由
		backend.Use(authorization.CheckTokenAuth(), authorization.CheckCasbinAuth())
		{
			// 用户组路由
			users := backend.Group("users/")
			{
				// 查询 ，这里的验证器直接从容器获取，是因为程序启动时，将验证器注册在了容器，具体代码位置：App\Http\Validator\Web\Users\xxx
				users.GET("list", validatorFactory.Create(consts.ValidatorPrefix+"UserList"))
				// 新增
				users.POST("create", validatorFactory.Create(consts.ValidatorPrefix+"UserCreate"))
				// 更新
				users.POST("edit", validatorFactory.Create(consts.ValidatorPrefix+"UserEdit"))
				// 删除
				users.POST("destroy", validatorFactory.Create(consts.ValidatorPrefix+"UserDestroy"))
				// 用户获取动态菜单
				users.GET("info", (&web.Users{}).UserInfo)
				users.GET("personal_info", (&web.Users{}).GetPersonalInfo)                              // 每个用户获取属于自己的账号信息
				users.POST("personal_edit", validatorFactory.Create(consts.ValidatorPrefix+"UserEdit")) // 每个用户编辑属于自己的账号信息
				// 用户获取视图页面拥有的权限按钮
				users.GET("has_view_button_list", validatorFactory.Create(consts.ValidatorPrefix+"ViewButtonList"))
			}
			//文件上传公共路由
			uploadFiles := backend.Group("upload/")
			{
				uploadFiles.POST("files", validatorFactory.Create(consts.ValidatorPrefix+"UploadFiles"))
			}
			// 组织机构、岗位
			authOrganizationPost := backend.Group("organization/")
			{
				authOrganizationPost.GET("get_by_fid", validatorFactory.Create(consts.ValidatorPrefix+"OrgPostGetByFid"))
				authOrganizationPost.GET("list", validatorFactory.Create(consts.ValidatorPrefix+"OrgPostList"))
				authOrganizationPost.POST("create", validatorFactory.Create(consts.ValidatorPrefix+"OrgPostCreate"))
				authOrganizationPost.POST("edit", validatorFactory.Create(consts.ValidatorPrefix+"OrgPostEdit"))
				authOrganizationPost.POST("destroy", validatorFactory.Create(consts.ValidatorPrefix+"OrgPostDestroy"))
			}
			// 岗位成员
			postMembers := backend.Group("post_members/")
			{
				postMembers.GET("list", validatorFactory.Create(consts.ValidatorPrefix+"PostMembersList"))
				postMembers.POST("create", validatorFactory.Create(consts.ValidatorPrefix+"PostMembersCreate"))
				postMembers.POST("edit", validatorFactory.Create(consts.ValidatorPrefix+"PostMembersEdit"))
				postMembers.POST("destroy", validatorFactory.Create(consts.ValidatorPrefix+"PostMembersDestroy"))
			}

			// 按钮公共模块
			buttonCnEn := backend.Group("button_cn_en/")
			{
				buttonCnEn.GET("list", validatorFactory.Create(consts.ValidatorPrefix+"ButtonCnEnList"))
			}

			// 系统菜单
			systemMenuList := backend.Group("system_menu/")
			{
				systemMenuList.GET("get_by_fid", validatorFactory.Create(consts.ValidatorPrefix+"SystemMenuGetByFid"))
				systemMenuList.GET("list", validatorFactory.Create(consts.ValidatorPrefix+"SystemMenuList"))
				systemMenuList.POST("create", validatorFactory.Create(consts.ValidatorPrefix+"SystemMenuCreate"))
				systemMenuList.POST("edit", validatorFactory.Create(consts.ValidatorPrefix+"SystemMenuEdit"))
				systemMenuList.POST("destroy", validatorFactory.Create(consts.ValidatorPrefix+"SystemMenuDestroy"))
				//系统菜单拥有的待分配所有权限按钮
				systemMenuList.GET("mount_auth_button", validatorFactory.Create(consts.ValidatorPrefix+"SysMenuMountButton"))

				// 为部门、岗位分配菜单、按钮权限
				systemMenuList.POST("assgin_to_org", validatorFactory.Create(consts.ValidatorPrefix+"AssginSystemMenuToOrg"))
				// 删除已分配给部门岗位的系统菜单、按钮
				systemMenuList.POST("del_auth_from_org", validatorFactory.Create(consts.ValidatorPrefix+"DelAuthFromOrg"))

				//待分配系统菜单、按钮属性列表
				systemMenuList.GET("all_list", validatorFactory.Create(consts.ValidatorPrefix+"SystemMenuListAllList"))
				// 已分配给部门、岗位的菜单、按钮
				systemMenuList.GET("assgined_list", validatorFactory.Create(consts.ValidatorPrefix+"AssginedSystemMenuList"))
			}
			// 权限分析
			authAnalysis := backend.Group("auth_analysis/")
			{
				authAnalysis.GET("user_list_with_post", validatorFactory.Create(consts.ValidatorPrefix+"UserListWithPost"))
				authAnalysis.GET("has_auth_list", validatorFactory.Create(consts.ValidatorPrefix+"OrgPostGetByUserId"))
			}
			// 按钮模块
			button := backend.Group("button/")
			{
				// 查询
				button.GET("list", validatorFactory.Create(consts.ValidatorPrefix+"ButtonList"))
				// 新增
				button.POST("create", validatorFactory.Create(consts.ValidatorPrefix+"ButtonCreate"))
				// 更新
				button.POST("edit", validatorFactory.Create(consts.ValidatorPrefix+"ButtonEdit"))
				// 删除
				button.POST("destroy", validatorFactory.Create(consts.ValidatorPrefix+"ButtonDestroy"))
			}
			// 省份城市
			provinceCity := backend.Group("province_city/")
			{
				// 查询
				provinceCity.GET("list", validatorFactory.Create(consts.ValidatorPrefix+"ProvinceCityList"))
				// 新增
				provinceCity.POST("create", validatorFactory.Create(consts.ValidatorPrefix+"ProvinceCityCreate"))
				// 更新
				provinceCity.POST("edit", validatorFactory.Create(consts.ValidatorPrefix+"ProvinceCityEdit"))
				// 删除
				provinceCity.POST("destroy", validatorFactory.Create(consts.ValidatorPrefix+"ProvinceCityDestroy"))
				// 根据id查询子级数据
				provinceCity.GET("get_sublist", validatorFactory.Create(consts.ValidatorPrefix+"ProvinceCitySubList"))
			}

		}
	}
	return router
}
