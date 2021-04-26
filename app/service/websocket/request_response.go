package websocket

//  这里可以定义业务相关的逻辑供ws.go文件调用

type Request struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
