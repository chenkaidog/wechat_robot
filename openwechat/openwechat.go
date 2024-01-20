package openwechat

import (
	"fmt"
	"runtime/debug"
	"wechat_robot/logrus"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

func Start() {
	for {
		newRobot()
	}
}

func newRobot() {
	defer func() {
		if rec := recover(); rec != nil {
			logrus.GetLogger().Errorf("openwechat panic: %v\n, stack: %s", rec, debug.Stack())
			return
		}
	}()

	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	bot.UUIDCallback = consoleQrCode

	reloadStorage := openwechat.NewFileHotReloadStorage("storage/storage.json")
	defer reloadStorage.Close()

	// 第一次进行热登录的时候热存储容器是空的，这时候会发生错误。失败后执行扫码登录
	if err := bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		logrus.GetLogger().Errorf("login err: %v", err)
		panic(err)
	}

	bot.MessageHandler = newDispatcher().AsMessageHandler()

	bot.Block()
}

func consoleQrCode(uuid string) {
	url := fmt.Sprintf("https://login.weixin.qq.com/l/" + uuid)
	logrus.GetLogger().Infof("login url: %s", url)

	q, err := qrcode.New(url, qrcode.Low)
	if err != nil {
		logrus.GetLogger().Errorf("new qrcode err: %v", err)
		panic(err)
	}

	fmt.Print(q.ToSmallString(true))
}
