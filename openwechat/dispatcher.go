package openwechat

import (
	"context"
	"wechat_robot/logrus"

	"github.com/eatmoreapple/openwechat"
)

func newDispatcher() *openwechat.MessageMatchDispatcher {
	dispatcher := openwechat.NewMessageMatchDispatcher()

	dispatcher.RegisterHandler(
		func(m *openwechat.Message) bool {
			return true
		},
		func(msgCtx *openwechat.MessageContext) {
			_ = msgCtx.AsRead()
			if msgCtx.IsSendBySelf() {
				msgCtx.Abort()
				return
			}
		},
	)

	dispatcher.OnText(func(msgCtx *openwechat.MessageContext) {
		ctx := logrus.Newcontext(context.Background(), msgCtx.MsgId)

		msgCtx.WithContext(ctx)
		logrus.GetLogger().CtxInfof(ctx, "receive text: %s", msgCtx.Content)
		handleText(ctx, msgCtx.Message)
	})

	return dispatcher
}
