package authorization

import (
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/users/curd"
	userstoken "goskeleton/app/service/users/token"
	"goskeleton/app/utils/response"
	"strconv"
	"strings"
)

type HeaderParams struct {
	Authorization string `header:"Authorization" binding:"required,min=20"`
}

// CheckTokenAuth 检查token完整性、有效性中间件
func CheckTokenAuth() gin.HandlerFunc {
	return func(context *gin.Context) {

		headerParams := HeaderParams{}

		//  推荐使用 ShouldBindHeader 方式获取头参数
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.ErrorParam(context, consts.JwtTokenMustValid+err.Error())
			return
		}
		token := strings.Split(headerParams.Authorization, " ")
		if len(token) == 2 && len(token[1]) >= 20 {
			tokenIsEffective := userstoken.CreateUserFactory().IsEffective(token[1])
			if tokenIsEffective {
				if customeToken, err := userstoken.CreateUserFactory().ParseToken(token[1]); err == nil {
					key := variable.ConfigYml.GetString("Token.BindContextKeyName")
					// token验证通过，同时绑定在请求上下文
					context.Set(key, customeToken)
				}
				context.Next()
			} else {
				response.ErrorTokenAuthFail(context)
			}
		} else {
			response.ErrorTokenBaseInfo(context)
		}
	}
}

// RefreshTokenConditionCheck 刷新token条件检查中间件，针对已经过期的token，要求是token格式以及携带的信息满足配置参数即可
func RefreshTokenConditionCheck() gin.HandlerFunc {
	return func(context *gin.Context) {

		headerParams := HeaderParams{}
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.ErrorParam(context, consts.JwtTokenMustValid+err.Error())
			return
		}
		token := strings.Split(headerParams.Authorization, " ")
		if len(token) == 2 && len(token[1]) >= 20 {
			// 判断token是否满足刷新条件
			if userstoken.CreateUserFactory().TokenIsMeetRefreshCondition(token[1]) {
				context.Next()
			} else {
				response.ErrorTokenRefreshFail(context)
			}
		} else {
			response.ErrorTokenBaseInfo(context)
		}
	}
}

// casbin检查用户对应的角色权限是否允许访问接口
func CheckCasbinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		requstUrl := c.Request.URL.Path
		method := c.Request.Method
		// 用户角色id需要存储在缓存，加快接口验证的效率(2021-03-11  后续实现)
		orgIds := curd.CreateUserCurdFactory().GetUserOrgIdsByRedis(c)
		var roleId int
		var isPass bool
		var err error
		for i := 0; i < len(orgIds); i++ {
			roleId = orgIds[i]
			isPass, err = variable.Enforcer.Enforce(strconv.Itoa(roleId), requstUrl, method)
			//fmt.Printf("Casbin权限校验参数打印：isPass:%v,角色ID：%d ,url：%s ,method: %s\n", isPass,roleId, requstUrl, method)
			if isPass == true {
				break
			}
		}

		if err != nil {
			response.ErrorCasbinAuthFail(c, err.Error())
			return
		} else if !isPass {
			response.ErrorCasbinAuthFail(c, "")
		} else {
			c.Next()
		}
	}
}

// CheckCaptchaAuth 验证码中间件
func CheckCaptchaAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		captchaIdKey := variable.ConfigYml.GetString("Captcha.captchaId")
		captchaValueKey := variable.ConfigYml.GetString("Captcha.captchaValue")
		captchaId := c.PostForm(captchaIdKey)
		value := c.PostForm(captchaValueKey)
		if captchaId == "" || value == "" {
			response.Fail(c, consts.CaptchaCheckParamsInvalidCode, consts.CaptchaCheckParamsInvalidMsg, "")
			return
		}
		if captcha.VerifyString(captchaId, value) {
			c.Next()
		} else {
			response.Fail(c, consts.CaptchaCheckFailCode, consts.CaptchaCheckFailMsg, "")
		}
	}
}
