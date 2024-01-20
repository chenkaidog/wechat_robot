package chatgpt

import (
	"context"
	"wechat_robot/redis"
)

func SetChatgptApiKey(ctx context.Context, apiKey string) error {
	return redis.SetChatgptApiKey(ctx, apiKey)
}
