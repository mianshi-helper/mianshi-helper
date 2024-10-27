package tools

import (
	"context"
	"mianshi-helper/engine"
)

func VerifyLogin(redisHelper *engine.RedisHelper, token string) string {
	// 去除Bearer
	tokenResult := token[7:]
	redisResult := redisHelper.Get(context.Background(), tokenResult).Val()
	claims, err := ParseToken(string(tokenResult))
	if err != nil {
		return "出现异常,请重新登陆"
	}
	if redisResult == "" {
		return "token已过期"
	} else if claims.Username != redisResult {
		return "token错误"
	} else {
		return "success"
	}
}
