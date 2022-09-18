package login_policy

import (
	"errors"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/redis_factory"
	"time"
)

// 登陆安全策略：
// 账号登陆连续失败N次，会限制账号在M分钟内禁止登陆
// N 、M 通过配置项设置

func CreateUsersLoginPolicyFactory() *UserLoginPolicy {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		variable.ZapLog.Error("redis 初始化出错-请在配置项检查redis配置参数")
		return nil
	}
	return &UserLoginPolicy{redisClient: redCli, keyPre: "login_fail_times:"}
}

type UserLoginPolicy struct {
	redisClient *redis_factory.RedisClient
	keyPre      string
}

// CheckAccountIsForbidden 检查账号是否被禁止登陆
// @account 账号
// @return false 表示账号没有禁用；true 表示账号被禁用
func (u *UserLoginPolicy) CheckAccountIsForbidden(account string) (bool, error) {
	if !u.accountIsExists(account) {
		return false, nil
	}

	if val, err := u.getFailTotalTimes(account); err == nil {
		maxLoginFailTimes := variable.ConfigYml.GetInt64("LoginPolicy.MaxLoginFailTimes")
		if val >= maxLoginFailTimes {
			return true, nil
		}
	} else {
		variable.ZapLog.Error("redis命令执行出错", zap.Error(err))
	}
	return false, nil
}

// AccountIsExists 判断账号是否存在
// @account 账号
// @return false 表示账号不存在；true 表示账号已经存在 于redis
func (u *UserLoginPolicy) accountIsExists(account string) bool {
	if val, err := u.redisClient.Int(u.redisClient.Execute("exists", u.keyPre+account)); err == nil && val > 0 {
		return true
	}
	return false
}

// SetAccountLoginCache 每次登陆设置账号登陆缓存
// @account 登陆账号
// @isFail 登陆是否失败
// @return 返回登录失败的累计次数
func (u *UserLoginPolicy) SetAccountLoginCache(account string, isFail bool) (int64, error) {
	// 如果是登陆失败，对应账号的等次失败次数+1
	if isFail {
		// 只要登录失败，失败次数就+1，（相关key不存在会自动创建）
		if _, err := u.redisClient.Execute("incr", u.keyPre+account); err != nil {
			variable.ZapLog.Error("redis命令执行出错", zap.Error(err))
			return -1, err
		}
		// 返回已经累计失败的次数
		if failTimes, err := u.getFailTotalTimes(account); err == nil {
			// 累计失败次数超限，自动禁止登陆
			maxLoginFailTimes := variable.ConfigYml.GetInt64("LoginPolicy.MaxLoginFailTimes")
			if failTimes >= maxLoginFailTimes {
				if err = u.setCountDown(account); err == nil {
					return failTimes, nil
				} else {
					return failTimes, err
				}
			} else {
				return failTimes, nil
			}
		} else {
			variable.ZapLog.Error("redis命令执行出错", zap.Error(err))
			return failTimes, err
		}

	} else {
		// 凡是登陆成功，累计登陆失败次数全部设置为0
		if _, err := u.redisClient.Execute("set", u.keyPre+account, 0); err == nil {
			return 0, nil
		} else {
			variable.ZapLog.Error("redis命令执行出错", zap.Error(err))
			return 0, err
		}
	}
}

// getFailTotalTimes 获取登录失败的累计次数
func (u *UserLoginPolicy) getFailTotalTimes(account string) (int64, error) {
	// 继续获取已经累计失败的次数
	if totalFilaTimes, err := u.redisClient.Int64(u.redisClient.Execute("get", u.keyPre+account)); err == nil {
		return totalFilaTimes, nil
	} else {
		variable.ZapLog.Error("redis命令执行出错", zap.Error(err))
		return -1, err
	}
}

// setCountDown 当超过登陆最大失败次数时，键设置倒计时的过期时间
// @account 过期的账号
// @second 过期的时间戳(时间点)
func (u *UserLoginPolicy) setCountDown(account string) error {
	expireSeconds := variable.ConfigYml.GetInt64("LoginPolicy.LoginFailCountDown")
	if expireSeconds < 1 {
		variable.ZapLog.Error("登陆策略配置错 - LoginPolicy.MaxLoginFailTimes 必须设置一个有效的禁止登陆倒计时(秒)")
		return errors.New("LoginPolicy.MaxLoginFailTimes 设置错误")
	}

	if _, err := u.redisClient.Execute("expireAt", u.keyPre+account, time.Now().Unix()+expireSeconds); err != nil {
		variable.ZapLog.Error("redis命令执行出错", zap.Error(err))
		return err
	}
	return nil
}

// ReleaseRedisConn 释放redis
func (u *UserLoginPolicy) ReleaseRedisConn() {
	u.redisClient.ReleaseOneRedisClient()
}
