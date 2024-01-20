package main

import (
	"wechat_robot/config"
	"wechat_robot/openwechat"
)

func main() {
	config.Init()
	openwechat.Init()
}
