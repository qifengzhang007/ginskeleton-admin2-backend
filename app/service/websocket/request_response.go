package websocket

import (
	"encoding/json"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/redis_factory"
)

//Code 常量定义
// 1000   进出人员列表
// 1050   养老服务

// redis 中的常量键
const PeopleRela = "person_mobility"
const ServiceRela = "service_list"

// 错误定义
const (
	RedisErrNoRecord     = "redis中没有人员进出记录信息"
	RequestErrNoCode     = "请求的Code码未定义"
	ResponseErrNotDecode = "服务端无法解析的json请求格式"
)

type Request struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

// 请求正确，完整响应
type Response struct {
	Code        int64  `json:"code"`
	Msg         string `json:"msg"`
	RequestCode string `json:"request_code"`
	Data        struct {
		InTotalPeople int64 `json:"in_total_people"`
		InList        []struct {
			ResidentName string `json:"resident_name"`
			Property     string `json:"property"`
			ImgPath      string `json:"img_path"`
			EventTime    string `json:"event_time"`
			CameraName   string `json:"camera_name"`
			Gender       string `json:"gender"`
		} `json:"in_list"`
		OutList []struct {
			ResidentName string `json:"resident_name"`
			Property     string `json:"property"`
			ImgPath      string `json:"img_path"`
			EventTime    string `json:"event_time"`
			CameraName   string `json:"camera_name"`
			Gender       string `json:"gender"`
		} `json:"out_list"`
	} `json:"data"`
}

// 请求的code 码未定义，简短响应错误码

type SimpleResponse struct {
	Code        int64  `json:"code"`
	Msg         string `json:"msg"`
	RequestCode int64  `json:"request_code"`
	Data        string `json:"data"`
}

// 请求相关的数据处理函数
func (r Request) DecodeJson(jsonStr string) string {
	// 定义一个简单的错误响应
	var simpleRes SimpleResponse
	if err := json.Unmarshal([]byte(jsonStr), &r); err == nil {
		simpleRes.RequestCode = r.Code
		if r.Code == 1000 {
			//进入人员清单
			redisClient := redis_factory.GetOneRedisClient()
			if res, err := redisClient.String(redisClient.Execute("Get", PeopleRela)); err == nil {
				redisClient.ReleaseOneRedisClient()
				return res
			} else {
				simpleRes.Msg = RedisErrNoRecord
				variable.ZapLog.Error(simpleRes.Msg, zap.Error(err))
			}
		} else if r.Code == 1050 {
			redisClient := redis_factory.GetOneRedisClient()
			if res, err := redisClient.String(redisClient.Execute("Get", ServiceRela)); err == nil {
				redisClient.ReleaseOneRedisClient()
				return res
			} else {
				simpleRes.Msg = RedisErrNoRecord
				variable.ZapLog.Error(simpleRes.Msg, zap.Error(err))
			}

		} else {
			simpleRes.Msg = RequestErrNoCode
			variable.ZapLog.Warn(simpleRes.Msg)
		}
	} else {
		simpleRes.Msg = ResponseErrNotDecode
		variable.ZapLog.Error(simpleRes.Msg, zap.Error(err))
	}
	res, _ := json.Marshal(simpleRes)
	return string(res)
}

// 服务器定时检测数据是否有变化，主动向客户端发送消息
func (r Request) GetCacheMsg(key string) string {
	// 定义一个简单的错误响应
	redisClient := redis_factory.GetOneRedisClient()
	if res, err := redisClient.String(redisClient.Execute("Get", key)); err == nil {
		redisClient.ReleaseOneRedisClient()
		return res
	} else {
		return ""
	}
}
