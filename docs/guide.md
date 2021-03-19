###   1.按钮管理表(tb_auth_button_cn_en) 的增删改查全功能开发  
> 本次业务介绍顺序将按照我个人的开发习惯展开：数据库创建表 -> 编写model -> 编写控制器 -> 编写表单参数验证器 -> 编写路由 
####  1.1 按钮管理表结构 数据库创建代码  
```code  
// 以下代码仅用于展示该表的结构  
CREATE TABLE tb_auth_button_cn_en (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  en_name char(60) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '英文字母',
  cn_name char(60) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '中文名',
  color varchar(100) DEFAULT 'default' COMMENT '菜单显示待分配的权限使用',
  allow_method varchar(30) DEFAULT '*' COMMENT '允许的请求方式',
  status tinyint(4) DEFAULT 1,
  remark varchar(300) CHARACTER SET utf8mb4 DEFAULT '',
  created_at datetime DEFAULT current_timestamp(),
  updated_at datetime DEFAULT current_timestamp(),
  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8 COMMENT='按钮代码统一管理'
```
####  2 为该表编写 model 代码
>   由于该表已经被本系统集成，详情代码请点击查看[tb_auth_button_cn_en model 详情](../app/model/button_cn_en.go)
```code   
// 我们现将最核心的增、删、改、查列举如下

func CreateButtonCnEnFactory(sqlType string) *ButtonCnEnModel {
return &ButtonCnEnModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type ButtonCnEnModel struct {
BaseModel
CnName      string `json:"cn_name"`
EnName      string `json:"en_name"`
AllowMethod string `json:"allow_method"`
Color       string `json:"color"`
Status      int    `json:"status"`
Remark      string `json:"remark"`
}

// 设置表名
func (b *ButtonCnEnModel) TableName() string {
return "tb_auth_button_cn_en"
}

//新增
func (b *ButtonCnEnModel) InsertData(c *gin.Context) bool {
    // 使用下面的绑定函数时必须定义一个独立变量
	var tmp ButtonCnEnModel
	
	// data_bind.ShouldBindFormDataToModel(c, &tmp) 是本系统对 gin.ShouldBind 函数的精简版
	// 核心作用是将表单参数验证器已经验证通过的表单参数直接绑定在 model 模型上
	// 绑定的依据是 model 表设置的 json 标签值、对应的字段类型与表单参数验证器设置的json值、数据类型一致
	// 此外，tmp 还会被自动绑定 created_at 、updated_at 二个字段值
	// tmp 是绑定后的值，实际开发时可自行打印查看绑定结果
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		if res := b.Create(&tmp); res.Error == nil {
			return true
		} else {
			variable.ZapLog.Error("ButtonModel 数据新增出错", zap.Error(res.Error))
		}
	} else {
		variable.ZapLog.Error("ButtonModel 数据绑定出错", zap.Error(err))
	}
	return false
}

//更新
func (b *ButtonCnEnModel) UpdateData(c *gin.Context) bool {
	var tmp ButtonCnEnModel
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
	
	// updates 函数对于零值字段则跳过
	// save 函数则是所有字段全量更新
		tmp.CreatedAt = ""  // 我们不需要更新 created_at 字段，那么需要手动设置该字段值为空即可
		if res := b.Updates(tmp); res.Error == nil {
			return true
		} else {
			variable.ZapLog.Error("ButtonModel 数据修改出错", zap.Error(res.Error))
		}
	} else {
		variable.ZapLog.Error("ButtonModel 数据绑定出错", zap.Error(err))
	}
	return false
}

//删除
func (b *ButtonCnEnModel) DeleteData(id int) bool {
	if b.Delete(b, id).Error == nil {
		return true
	}
	return false
}

// 查询
func (a *ButtonCnEnModel) List(cnName string, limitStart, limit float64) (counts int64, data []ButtonCnEnModel) {
	
	counts = a.getCountsByButtonName(cnName) // 该函数的具体代码暂时忽略即可
	if counts > 0 {
	// 查询不管是使用 gorm 语法还是写原生 sql , 按照数据库性能最优原则，都必须指定查询的字段
	// 日期时间字段在本系统全部按照字符串处理即可 
		if err := a.Model(a).
			Select("id", "en_name", "cn_name", "allow_method", "color", "status", "cn_name", "remark", "DATE_FORMAT(created_at,'%Y-%m-%d %H:%i:%s')  as created_at", "DATE_FORMAT(updated_at,'%Y-%m-%d %H:%i:%s')  as updated_at").
			Where("cn_name like ?", "%"+cnName+"%").Offset(int(limitStart)).Limit(int(limit)).Find(&data); err.Error == nil {
			return
		}
	}

	return 0, nil
}
```

####  3 编写 controller 代码
> 该表的控制器代码非常简洁、简单,没有任何陌生语法, 请点击查看即可：[进入详情](../app/http/controller/web/auth/button.go)


####  4 编写表单参数验证器(必看)  
> 表单参数验证器代码也是非常简洁, [点击查看详情](../app/http/validator/web/auth/button)  
```code  
// 这里特别说明一下表单参数验证器使用时注意事项
type ButtonCreate struct {
	BaseField
}

// 禁忌点：CheckParams 函数没有绑定在 指针上 ，错误的绑定形式 func (c *ButtonCreate)CheckParams {}
// 因为表单参数验证器在本系统是在容器中进行管理的，该函数绑定在指针，当第一请求过后，容器的代码段会被污染，所以这里绝对不能绑定在指针上 
func (c ButtonCreate) CheckParams(context *gin.Context) {
    // 中间代码省略
}

// 关于验证表单字段的语法补充介绍
// 数字字段请务必设置为  float64 系列，需要 int型，请使用  int()  int64()  等函数进行二次转换
// *float64 则表示接受该字段为 0 的用户提交 ； 
//  float64  无法接受用户将该子弹设置为 0 提交
Status        *float64 `form:"status" json:"status" binding:"required,min=0"` 

```
####  5 将以上表单参数验证器在容器进行注册  
```code 
	//按钮部分
	{
	    // 省略其他代码...
		key = consts.ValidatorPrefix + "ButtonCreate"
		containers.Set(key, button.ButtonCreate{})
		
	}

```

####  6 编写路由
```code 
    // 按钮模块
    button := backend.Group("button/")
    {
        // 新增
        button.POST("create", validatorFactory.Create(consts.ValidatorPrefix+"ButtonCreate"))
    }
    
```

#### 7 将按钮管理表(tb_auth_button_cn_en) 绑定在某个菜单
> 每一个 model 都可以绑定给一个菜单, 然后该菜单就可以分配给任何岗位、部门、公司，而最终用户则通过他所挂接的岗位，逐级向上继承一系列权限.  
![将后台开发完成的模块绑定给菜单](https://www.ginskeleton.com/images/menu_set.jpg)  

#### 8 将按钮管理表(tb_auth_button_cn_en)已绑定的菜单分配给组织机构(岗位、部门、公司等) 
![本次我们分配给超级管理员岗位](https://www.ginskeleton.com/images/auth_assgined.jpg)  

#### 9 最后,前端开发相关的页面即可  
> 相关代码请打开前端项目，路径：src/view/button/index.vue , 参考编写即可.  
> 由于前后端是一个整体，前端开发也有注意事项，具体参见前端开发指南，后端部分至此已经开发完成.  