###    快速上手
>   1.安装的go语言版本最好>=1.15,只为更好的支持 `go module` 包管理.  
>   2.配置go包的代理，参见`https://goproxy.cn`,有详细设置教程.    
>   3.使用 `goland(>=2019.3版本)` 打开本项目，找到`database/db_demo_mysql.sql`导入数据库，自行配置账号、密码、端口等。    
>   4.双击`cmd/(web|api|cli)/main.go`，进入代码界面，鼠标右键`run`运行本项目，首次会自动下载依赖， 片刻后即可启动.    
>![业务主线图](https://www.ginskeleton.com/GinSkeleton.jpg)  

###  交叉编译(windows直接编译出linux可执行文件)    
>   1 `goland` 终端底栏打开`terminal`, 依次执行 `set GOARCH=amd64` 、`set GOOS=linux` 、`set CGO_ENABLED=0` , 特别说明：以上命令执行时后面不要有空格，否则报错!    
>   2 进入根目录（GinSkeleton所在目录）：`go build -o demo_goskeleton cmd/(web|api|cli)/main.go` 可交叉编译出（web|api|cli）对应的二进制文件。     

###    <font color="red">项目骨架主线、核心逻辑</font>  
>   1.这部分主要介绍了`项目初始化流程`、`路由`、`表单参数验证器`、`控制器`、`model`、`service` 以及 `websocket` 为核心的主线逻辑.   
[进入主线逻辑文档](docs/document.md)  

###    测试用例路由  
[进入Api接口测试用例文档](docs/api_doc.md)      

###    开发常用模块  
>   随着项目不断完善以下列表模块会陆续增加, 各个模块被贯穿在本项目骨架的主线中, 因此只要掌握主线核心逻辑, 其余在为主线提供服务.  

序号|功能模块 | 文档地址  
---|---|---
1| 全局变量(日志、gorm、配置模块)|  [清单一览](docs/global_variable.md)  
2 | 消息队列| [rabbitmq文档](docs/rabbitmq.md)   
3 | cli命令| [cobra文档](docs/cobra.md) 
4 | goCurl、httpClient|[httpClient客户端](https://gitee.com/daitougege/goCurl) 
5|[websocket js客户端](docs/ws_js_client.md)| [websocket服务端](app/service/websocket/ws.go)  
6|aop切面编程| [Aop切面编程](docs/aop.md) 
7|redis| [redis使用示例](test/redis_test.go) 
8|gorm_v2操作(mysql、sqlserver、postgreSql)| [gorm v2 测试用例](test/gormv2_test.go) 
9|日志记录|  [zap高性能日志](docs/zap_log.md) 
10|项目日志对接到 elk 服务器|  [elk 日志顶级解决方案](docs/elk_log.md) 
11| 验证码|  [验证码](docs/captcha.md)
12| nginx配置(https、负载均衡)|[nginx配置详情](docs/nginx.md) 
13|supervisor| [supervisor进程守护](docs/supervisor.md)   


###    项目上线后，运维方案(基于docker)    
序号|运维模块 | 文档地址  
---|---|---
1 | linux服务器| [详情](docs/deploy_linux.md)   
2 | mysql| [详情](docs/deploy_mysql.md)  
3 | redis| [详情](docs/deploy_redis.md)    
4 | nginx| [详情](docs/deploy_nginx.md)   
5 | go应用程序| [详情](docs/deploy_go.md)  

### 并发测试
[点击查看详情](docs/bench_cpu_memory.md)

### 性能分析报告  
> 1.开发之初，我们的目标就是追求极致的高性能,因此在项目整体功能越来越趋于完善之时，我们现将进行一次全面的性能分析评测.    
> 2.通过执行相关代码, 跟踪 cpu 耗时 和 内存占用 来分析各个部分的性能,CPU耗时越短性、内存占用越低能越优秀,反之就比较垃圾.        

###  通过CPU的耗时来分析相关代码段的性能  
序号|分析对象 | 文档地址  
---|---|---
1| 项目骨架主线逻辑| [主线分析报告](./docs/project_analysis_1.md)
2| 操作数据库代码段| [操作数据库代码段分析报告](./docs/project_analysis_2.md)

###  通过内存占用来分析相关代码段的性能 
序号|分析对象 | 文档地址  
---|---|---
1| 操作数据库代码段| [操作数据库代码段](./docs/project_analysis_3.md) 
 
 #### V1.0.00  2020-11-05  
 * 智慧养老项目初始化  