package test

import (
	"goskeleton/app/global/variable"
	_ "goskeleton/bootstrap"
	"testing"
	"time"
)

// 容器大面积异步注册时测试

func TestContainer(t *testing.T) {

	//cc:=container.CreateContainersFactory()

	for i := 0; i <= 200; i++ {

		go variable.ConfigYml.GetInt("Token.IsCacheToRedis")

	}

	time.Sleep(10 * time.Second)

}
