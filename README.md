### GinSkeleton-Admin2 (后端部分)
> 基于 GinSkeleton v1.5.xx 开发的全新后台管理系统.


###  [在线文档](https://www.yuque.com/xiaofensinixidaouxiang/qmanaq/qmucb4)
> 文档包含了最主要的使用功能说明、界面效果图、演示地址等.


### 更新日志
#### v2.0.19  2024-08-03
**更新**
- 1.`websocket` 修复断电、直接拔网线导致服务端检测的终端在线状态不准确的bug, 因为直接断电、拔网线客户端的回调事件(onClose、onError)根本无法传递出去,服务端对应的socket文件状态无法及时变化.
- 3.项目依赖包全部更新至最新版.

