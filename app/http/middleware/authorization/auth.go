package authorization

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/users/curd"
	userstoken "goskeleton/app/service/users/token"
	"goskeleton/app/utils/response"
	"strconv"
	"strings"
)

type HeaderParams struct {
	Authorization string `header:"Authorization"`
}

// 检查token权限
func CheckTokenAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		//  模拟验证token
		headerParams := HeaderParams{}

		//  推荐使用 ShouldBindHeader 方式获取头参数
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			variable.ZapLog.Error(my_errors.ErrorsValidatorBindParamsFail, zap.Error(err))
			context.Abort()
		}

		if len(headerParams.Authorization) >= 20 {
			token := strings.Split(headerParams.Authorization, " ")
			if len(token) == 2 && len(token[1]) >= 20 {
				tokenIsEffective := userstoken.CreateUserFactory().IsEffective(token[1])
				if tokenIsEffective {
					if customeToken, err := userstoken.CreateUserFactory().ParseToken(token[1]); err == nil {
						context.Set("customeToken", customeToken)
					}
					context.Next()
				} else {
					response.ErrorTokenAuthFail(context)
				}
			}
		} else {
			response.ErrorTokenAuthFail(context)
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
