package baidu

import (
	"context"
	"wechat_robot/redis"
)

func SetErnieAppKey(ctx context.Context, appKey string) error {
	return redis.SetErnieAppKey(ctx, appKey)
}

func SetErnieAppSecret(ctx context.Context, appSecret string) error {
	return redis.SetErnieAppSecret(ctx, appSecret)
}

func RefreshToken(ctx context.Context) error {
	accessInfo, err := refreshToken(ctx)
	if err != nil {
		return err
	}

	setAccassInfo(accessInfo)

	return nil
}

