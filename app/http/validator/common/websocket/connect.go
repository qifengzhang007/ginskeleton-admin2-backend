package websocket

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	controllerWs "goskeleton/app/http/controller/websocket"
	"goskeleton/app/http/validator/core/data_transfer"
)

type Connect struct {
	Token string `form:"token" json:"token" binding:"required,min=1"`
}

// 验证器语法，参见 Register.go文件，有详细说明
// 注意：websocket 连接建立之前如果有错误，只能在服务端同构日志输出方式记录（因为使用response.Fail等函数，客户端是收不到任何信息的）

func (c Connect) CheckParams(context *gin.Context) {

	// 1. 首先检查是否开启websocket服务配置（在配置项中开启）
	if variable.ConfigYml.GetInt("Websocket.Start") != 1 {
		variable.ZapLog.Error(consts.WsServerNotStartMsg)
		return
	}
	//2.基本的验证规则没有通过
	if err := context.ShouldBind(&c); err != nil {
		variable.ZapLog.Error("客户端上线参数不合格", zap.Error(err))
		return
	}
	extraAddBindDataContext := data_transfer.DataAddContext(c, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		variable.ZapLog.Error("websocket-Connect 表单验证器json化失败")
		context.Abort()
		return
	} else {
		if serviceWs, ok := (&controllerWs.Ws{}).OnOpen(context); ok == false {
			variable.ZapLog.Error(consts.WsOpenFailMsg)
		} else {
			(&controllerWs.Ws{}).OnMessage(serviceWs, context) // 注意这里传递的service_ws必须是调用open返回的，必须保证的ws对象的一致性
		}
	}

}
