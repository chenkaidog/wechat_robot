package openwechat

import (
	"context"
	"wechat_robot/logrus"

	"github.com/eatmoreapple/openwechat"
)

func newDispatcher() *openwechat.MessageMatchDispatcher {
	dispatcher := openwechat.NewMessageMatchDispatcher()

	dispatcher.OnText(func(msgCtx *openwechat.MessageContext) {
		msg := msgCtx.Message

		if msg.IsSendBySelf() {
			return
		}

		ctx := logrus.Newcontext(context.Background(), msg.MsgId)

		msg.WithContext(ctx)
		logrus.GetLogger().CtxInfof(ctx, "receive msg: %s", msg.Content)
		handleText(ctx, msg)
	})

	return dispatcher
}
