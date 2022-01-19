package token_cache_redis

import (
	"goskeleton/app/utils/redis_factory"
	"strconv"
	"time"
)

func CreateUsersTokenCacheFactory() *userTokenCacheRedis {
	return &userTokenCacheRedis{redisClient: redis_factory.GetOneRedisClient()}
}

type userTokenCacheRedis struct {
	redisClient *redis_factory.RedisClient
}

// SetCache 设置缓存
func (u *userTokenCacheRedis) SetCache(userId, tokenId int64, token string) bool {
	keyName := "user_id:" + strconv.FormatInt(userId, 10)
	if _, err := u.redisClient.Int(u.redisClient.Execute("hSet", keyName, tokenId, token)); err == nil {
		return true
	}
	return false
}

// DelCache 删除缓存
func (u *userTokenCacheRedis) DelCache(keyName string) {

}

// CacheIsExists 查询token是否在redis存在
func (u *userTokenCacheRedis) TokenCacheIsExists(userId int64) bool {
	keyName := "user_id:" + strconv.FormatInt(userId, 10)
	// 打印全部
	if _, err := u.redisClient.Int(u.redisClient.Execute("hGetAll", keyName)); err == nil {
		return true
	}
	return false

}

// GetCache 查询缓存
func (u *userTokenCacheRedis) GetCache(keyName string) {

}

// 设置键的过期时间
func (u *userTokenCacheRedis) SetKeyExpire(keyName string, expireSec int64) bool {

	if _, err := u.redisClient.Int(u.redisClient.Execute("expireAt", keyName, time.Now().Unix()+expireSec)); err == nil {
		return true
	}
	return false
}

// releaseRedisConn 释放redis
func (u *userTokenCacheRedis) releaseRedisConn() {
	u.redisClient.ReleaseOneRedisClient()
}
