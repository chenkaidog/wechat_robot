package chat

import (
	"context"
	"wechat_robot/chat/baidu"
	"wechat_robot/chat/chatgpt"
	"wechat_robot/chat/common"
)

type ChatGenerator interface {
	Reply2Chat(ctx context.Context, content string, history []*common.Message) (string, error)
}

func NewChatGenerator(platform, model string) ChatGenerator {
	switch platform {
	case common.PlatformChatGPT:
		return chatgpt.NewChatGPTgen(model)
	}

	return baidu.NewBaiduChatGen(model)
}
