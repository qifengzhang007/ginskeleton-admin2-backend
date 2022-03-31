
### 版本更新历史日志

#### V1.0.34  (2022-03-22)
###### 搭配的前端版本：>= [V1.0.25](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
* 细节完善
- 1.gorm 相关的回调函数增强条件判断, 加强程序严谨性。
- 2.自带数据库文件更新, 给字段设置了默认值.

V1.0.33 (2022-02-27)
###### 搭配的前端版本：>= [V1.0.25](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.token 缓存到 redis 逻辑加强严谨性.

#### V1.0.32 (2022-02-10)
###### 搭配的前端版本：>= [V1.0.25](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
* bug 修复
- 1.大批量分配权限时会出现部分失败的情况。
- 2.老版本增量升级：`app/model/users.go  139 行 对应的函数 OauthCheckTokenIsOk ` 覆盖老版本即可.


#### V1.0.31 (2022-02-05)
###### 搭配的前端版本：>= [V1.0.25](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
* bug 修复
- 1.大批量分配权限时会出现部分失败的情况，老版本 admin 升级最快捷方式：直接使用最新版本 `app/model/auth/auth_menu_assign.go` 文件覆盖同名文件即可.
* 更新
- 1.自带数据库文件更新, 为大批量操作的相关表字段创建索引，老版本 admin 系统可忽略，与上一条 bug 修复无依赖关系.

#### V1.0.29 (2022-01-25)
###### 搭配的前端版本：>= [V1.0.25](http://gitee.com/daitougege/gin-skeleton-admin-frontend)

* 新增
- 1.用户 `token` 缓存到 `redis` 功能,如果项目使用了 `redis` , 请直接在 config/config.yml 文件中设置 `Token.IsCacheToRedis = 1`
- 2.项目初始化时增加设置信任代理服务器ip列表，gin(v1.7.7)新增功能,详情参见相关配置项说明.
- 3.新增个人信息自主编辑功能，方便小权限账号编辑自己的账号密码.

* 更新
- 1.配置文件缓存时加锁,避免开发者频繁注册时,程序出现提示。
- 2.用户token鉴权时,如果开启了redis缓存功能，优先查询redis.
- 3.自带数据库文件更新，以便支持新增的个人信息编辑功能.
- 4.所有底层依赖包更新至最新版.

#### V1.0.28 (2022-02-27)
###### 搭配的前端版本：>= [V1.0.22](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.`token` 缓存到 `redis` 逻辑加强严谨性.

#### V1.0.27 (2021-12-20)
###### 搭配的前端版本：>= [V1.0.22](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.错误日志记录时同时记录调用链信息。
- 2.rabbitmq 消息队列增加消息延迟发送功能.
- 3.关于 rabbitmq 消息延迟使用请参考新版在线文档.

#### V1.0.26 (2021-11-28)
###### 搭配的前端版本：>= [V1.0.22](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.引入表单参数验证器全局翻译器,简化代码书写,提升开发效率.

#### V1.0.25 (2021-11-26)
###### 搭配的前端版本：>= [V1.0.22](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.将主线版本更新的内容合并至 `admin-backend` 版本.
- 2.详细更新日志参见主线版本(v1.5.29)的更新日志.

#### V1.0.24 (2021-10-24)
###### 搭配的前端版本：>= [V1.0.20](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.修复bug：添加系统菜单时，接口参数有空值时,会导致存储在数据库的个别字段为NULL.

#### V1.0.23 (2021-10-11)
###### 搭配的前端版本：>= [V1.0.14](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.更新token刷新接口逻辑，解决认证中间件被加载2次的小bug.

#### V1.0.22 (2021-09-13)
###### 搭配的前端版本：>= [V1.0.14](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.更新token刷新接口逻辑，支持在过期24小时内使用旧token换取新token.
- 2.验证码自定义验证逻辑部分完善代码严谨性.

#### V1.0.21 (2021-08-02)
###### 搭配的前端版本：>= [V1.0.14](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.升级项目依赖版本至最新版
- 2.修正一处单词拼写错误.

#### V1.0.20 (2021-07-09)
###### 搭配的前端版本：>= [V1.0.14](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.删除一处调试语句.

#### V1.0.19 (2021-06-29)
###### 搭配的前端版本：>= [V1.0.14](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.修正自带数据库(`./database/db_ginskeleton_20210629.7z`)——系统菜单表 `tb_auth_system_menu` 菜单排序字段 `sort` 为  `string` 的 `bug`, 导致菜单排序没有按照数字大小排序.
- 2.系统菜单表 `tb_auth_system_menu` , 菜单排序字段 `sort` 如果手动指定值，则父级值必须 > 子级菜单值，子级之间则没有任何限制.
- 3.本次 `bug` 也可以手动快速修复，自行修改已导入数据库表  `tb_auth_system_menu` 的 `sort` 字段类型为  `int` 即可.


#### V1.0.18 (2021-06-18)
###### 搭配的前端版本：>= [V1.0.14](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.修正常量定义处日期格式单词书写错误问题.
- 2.一个用户同时允许最大在线的token, 查询时优先按照 expires_at 倒序排列,便于不同系统间对接时,那种长久有效的token不会被"踢"下线.

#### V1.0.17 (2021-05-29)
###### 搭配的前端版本：>= [V1.0.14](http://gitee.com/daitougege/gin-skeleton-admin-frontend)
- 1.账号密码使用验证码中间件代替原来的异步逻辑,避免异步验证验证码时用户使用postman请求接口时绕过了验证码校验机制.

#### V1.0.16 (2021-05-28)
###### 搭配的前端版本：V1.0.13
- 1.针对小权限账号启动本项目时，涉及到文件上传，自动创建目录时，调整文件夹权限为 os.ModePerm.

#### V1.0.15 (2021-05-12)
- 1.修复数据查询时间格式化书写错误的bug,mysql日期时间格式化由 %Y-%m-%d %h:%i:%s 修复为 ：%Y-%m-%d %H:%i:%s ，由于 %h 小写导致了时间比实际滞后8小时.

#### V1.0.14 (2021-05-04)
- 1.cobra增加创建子目录的示例代码.

#### V1.0.12 (2021-04-26)
- 1.表单参数验证器独立为：api、web,进一步简化项目代码比较多时，程序的简洁性.
- 2.token生成时计算时间戳mysql调整为go函数计算,避免mysql函数 FROM_UNIXTIME 参数最大支持21亿的局限性.
- 3.ws删除了原有的部分业务代码,使程序保持最简洁性.
- 4.项目骨架依赖的核心包升级至最新版.

#### V1.0.11 (2021-04-16)
- 1.修复bug:在mysql8高版本系列，limit 后面的参数是浮点型查询不到数据。
- 2.更新自带数据库，修改oauth_token表的token字段长度为600.
- 3.更新sql结果树形化包版本，以便支持更多复杂的sql查询结果树形化.

#### V1.0.10 (2021-04-02)
- 1.系统菜单主表与子表数据接受方式以及后续的总体逻辑更新.

#### V1.0.01 (2021-03-27)
- Bug修复:
- 1.系统菜单添加按钮失败的错误.  
  1.1 涉及到的文件：app/model/auth/auth_system_menu_button.go  
  1.2 app/service/auth_system_menu/auth_organization_post_service.go  
  2.快速升级、更新的办法：可直接使用官方仓库最新代码覆盖自带代码即可.

####    V1.0.00 (2021-03-20)
> 1.GinSkeleton-Admin 系统 v1.0.0 版本发布. 