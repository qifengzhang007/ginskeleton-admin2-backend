### GinSkeleton-Admin2 (后端部分)
> 基于 GinSkeleton v1.5.xx 开发的全新后台管理系统.


###  [在线文档](https://www.yuque.com/xiaofensinixidaouxiang/qmanaq/qmucb4)
> 文档包含了最主要的使用功能说明、界面效果图、演示地址等.


### 更新日志
**v2.0.09  2022-09-19**
 - 更新 
主要是将主线版本最近的更新同步到admin系统：
 - 1.增加账号登陆安全策略,账号密码连续出错超过配置项允许的最大次数时,自动禁止2分钟.
 - 2.以上功能依赖于redis，默认没有开启,建议开发者正确配置redis后开启此功能,详情参见配置项.
