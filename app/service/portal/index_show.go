package portal

import (
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/redis_factory"
)

// 门户网站，首页接口数据缓存逻辑
func CreateIndexShowFctory() *IndexShowApiCache {
	return &IndexShowApiCache{
		redis: redis_factory.GetOneRedisClient(),
	}
}

type IndexShowApiCache struct {
	redis *redis_factory.RedisClient
}

func (i *IndexShowApiCache) ApiUriIsCache(keyUri string) bool {
	if res, err := i.redis.Int(i.redis.Execute("exists", "IndexShow_"+md5_encrypt.MD5(keyUri))); err == nil && res == 1 {
		return true
	} else {
		return false
	}
}

func (i *IndexShowApiCache) SetApiUriCache(keyUri, valuesStr string) bool {
	if res, err := i.redis.String(i.redis.Execute("setex", "IndexShow_"+md5_encrypt.MD5(keyUri), 3600, valuesStr)); err == nil && res == "OK" {
		return true
	} else {
		return false
	}
}

func (i *IndexShowApiCache) GetApiUriCache(keyUri string) string {
	if res, err := i.redis.String(i.redis.Execute("get", "IndexShow_"+md5_encrypt.MD5(keyUri))); err == nil {
		return res
	} else {
		return ""
	}
}

// 释放redis连接至连接池

func (i *IndexShowApiCache) ReleaseRedisConn() {
	i.redis.ReleaseOneRedisClient()
}
