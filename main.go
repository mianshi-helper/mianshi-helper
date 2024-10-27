package main

import (
	Router "mianshi-helper/router"
)

// 测试版不做登录注册，主要做面试助手本身模型能力的测试
func main() {
	// 初始化路由
	router := Router.SetupRouter()
	router.Run(":7912")

}
