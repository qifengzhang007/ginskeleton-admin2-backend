package token_cache_redis

import (
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/redis_factory"
	"strconv"
	"strings"
	"time"
)

func CreateUsersTokenCacheFactory(userId int64) *userTokenCacheRedis {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &userTokenCacheRedis{redisClient: redCli, userTokenKey: "token_userid_" + strconv.FormatInt(userId, 10)}
}

type userTokenCacheRedis struct {
	redisClient  *redis_factory.RedisClient
	userTokenKey string
}

// SetTokenCache 设置缓存
func (u *userTokenCacheRedis) SetTokenCache(tokenExpire int64, token string) bool {
	// 存储用户token时转为MD5，下一步比较的时候可以更加快速地比较是否一致
	if _, err := u.redisClient.Int(u.redisClient.Execute("zAdd", u.userTokenKey, tokenExpire, md5_encrypt.MD5(token))); err == nil {
		return true
	}
	return false
}

// DelOverMaxOnlineCache 删除缓存,删除超过系统允许最大在线数量之外的用户
func (u *userTokenCacheRedis) DelOverMaxOnlineCache() bool {
	// 首先先删除过期的token
	_, _ = u.redisClient.Execute("zRemRangeByScore", u.userTokenKey, 0, time.Now().Unix()-1)

	onlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	alreadyCacheNum, err := u.redisClient.Int(u.redisClient.Execute("zCard", u.userTokenKey))
	if err == nil && alreadyCacheNum > onlineUsers {
		// 删除超过最大在线数量之外的token
		if alreadyCacheNum, err = u.redisClient.Int(u.redisClient.Execute("zRemRangeByRank", u.userTokenKey, 0, alreadyCacheNum-onlineUsers-1)); err == nil {
			return true
		} else {
			variable.ZapLog.Error("删除超过系统允许之外的token出错：", zap.Error(err))
		}
	}
	return false
}

// TokenCacheIsExists 查询token是否在redis存在
func (u *userTokenCacheRedis) TokenCacheIsExists(token string) (exists bool) {
	token = md5_encrypt.MD5(token)
	curTimestamp := time.Now().Unix()
	onlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	if strSlice, err := u.redisClient.Strings(u.redisClient.Execute("zRevRange", u.userTokenKey, 0, onlineUsers-1)); err == nil {
		for _, val := range strSlice {
			if score, err := u.redisClient.Int64(u.redisClient.Execute("zScore", u.userTokenKey, token)); err == nil {
				if score > curTimestamp {
					if strings.Compare(val, token) == 0 {
						exists = true
						break
					}
				}
			}
		}
	} else {
		variable.ZapLog.Error("获取用户在redis缓存的 token 值出错：", zap.Error(err))
	}
	return
}

// SetUserTokenExpire 设置用户的 usertoken 键过期时间
// 参数： 时间戳
func (u *userTokenCacheRedis) SetUserTokenExpire(ts int64) bool {
	if _, err := u.redisClient.Execute("expireAt", u.userTokenKey, ts); err == nil {
		return true
	}
	return false
}

// ClearUserToken 清除某个用户的全部缓存，当用户更改密码或者用户被禁用则删除该用户的全部缓存
func (u *userTokenCacheRedis) ClearUserToken() bool {
	if _, err := u.redisClient.Execute("del", u.userTokenKey); err == nil {
		return true
	}
	return false
}

// ReleaseRedisConn 释放redis
func (u *userTokenCacheRedis) ReleaseRedisConn() {
	u.redisClient.ReleaseOneRedisClient()
}
