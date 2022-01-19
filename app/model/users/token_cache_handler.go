package users

import (
	"fmt"
	"goskeleton/app/global/variable"
	"time"
)

// 本文件专门处理 token 缓存到 redis 的相关逻辑

func (u *UsersModel) ValidTokenCacheToRedis(userId int) {

	sql := "SELECT   token  FROM  `tb_oauth_access_tokens`  WHERE   fr_user_id=?  AND  revoked=0  AND  expires_at>NOW() ORDER  BY  expires_at  DESC , updated_at  DESC  LIMIT ?"
	maxOnlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	var tokens = make([]TokenToRedis, maxOnlineUsers)
	if u.Raw(sql, userId, maxOnlineUsers).Find(&tokens).Error == nil {

		for _, item := range tokens {

			stamp, _ := time.ParseInLocation(variable.DateFormat, item.ExpiresAt, time.Local)
			fmt.Println(item, stamp)
		}

	}

}
