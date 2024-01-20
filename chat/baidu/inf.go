package baidu

import (
	"context"
	"wechat_robot/chat/common"
)

func NewBaiduChatGen(model string) *BaiduCahtGen {
	return &BaiduCahtGen{}
}

type BaiduCahtGen struct{}

func (b *BaiduCahtGen) Reply2Chat(ctx context.Context, content string, history []*common.Message) (string, error) {
	return SendMsg(ctx, content, history...)
}
