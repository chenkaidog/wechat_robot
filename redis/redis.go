package redis

import (
	"context"
	"fmt"
	"wechat_robot/config"
	"wechat_robot/logrus"

	"github.com/redis/go-redis/v9"
)

func init() {
	rdbClient = redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.GetRedis().IP, config.GetRedis().Port),
			Password: config.GetRedis().Password,
			DB:       config.GetRedis().DB,
		})

	rdbClient.AddHook(new(loggerHook))
}

var rdbClient *redis.Client

const (
	keyErnieAppKey    = "openwechat_ernie_app_key"
	keyErnieAppSecret = "openwechat_ernie_app_secret"
	keyChatgptApiKey  = "openwechat_chatgpt_api_kei"
	keyRobotOption    = "openwechat_robot_option"
	keyWordMask       = "openwechat_word_mask"
)

func Get(ctx context.Context, key string) (string, error) {
	result, err := rdbClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", redis.Nil
		}

		logrus.GetLogger().CtxErrorf(ctx, "redis get err: %v", err)
		return "", err
	}

	return result, nil
}

func GetErnieAppKey(ctx context.Context) string {
	result, _ := rdbClient.Get(ctx, keyErnieAppKey).Result()
	if result != "" {
		return result
	}

	appKey := config.GetBaidu().AppKey

	_, _ = rdbClient.SetNX(ctx, keyErnieAppKey, appKey, 0).Result()

	return appKey
}

func SetErnieAppKey(ctx context.Context, appKey string) error {
	_, err := rdbClient.Set(ctx, keyErnieAppKey, appKey, 0).Result()
	return err
}

func GetErnieAppSecret(ctx context.Context) string {
	result, _ := rdbClient.Get(ctx, keyErnieAppSecret).Result()
	if result != "" {
		return result
	}

	appSecret := config.GetBaidu().AppSecret

	_, _ = rdbClient.SetNX(ctx, keyErnieAppSecret, appSecret, 0).Result()

	return appSecret
}

func SetErnieAppSecret(ctx context.Context, appSecret string) error {
	_, err := rdbClient.Set(ctx, keyErnieAppSecret, appSecret, 0).Result()
	return err
}

func GetChatgptApiKey(ctx context.Context) string {
	result, _ := rdbClient.Get(ctx, keyChatgptApiKey).Result()
	if result != "" {
		return result
	}

	apiKey := config.GetChatGPT().ApiKey

	_, _ = rdbClient.SetNX(ctx, keyChatgptApiKey, apiKey, 0).Result()

	return apiKey
}

func SetChatgptApiKey(ctx context.Context, apiKey string) error {
	_, err := rdbClient.Set(ctx, keyChatgptApiKey, apiKey, 0).Result()
	return err
}

type RobotOption struct {
	Platform string `redis:"platform"`
	Model    string `redis:"model"`
}

func SetRobotOption(ctx context.Context, opt *RobotOption) error {
	_, err := rdbClient.HSet(ctx, keyRobotOption, opt).Result()
	return err
}

func GetRobotOption(ctx context.Context) (*RobotOption, error) {
	result, err := rdbClient.HGetAll(ctx, keyRobotOption).Result()
	if err != nil {
		logrus.GetLogger().CtxErrorf(ctx, "hgetall err: %v ", err)
		return nil, err
	}

	return &RobotOption{
		Platform: result["platform"],
		Model:    result["model"],
	}, nil
}

func SetWordMask(ctx context.Context, word, mask string) error {
	_, err := rdbClient.HSet(ctx, keyWordMask, word, mask).Result()
	return err
}

func DeleteWordMask(ctx context.Context, word string) error {
	_, err := rdbClient.HDel(ctx, keyWordMask).Result()
	return err
}

func GetWordMask(ctx context.Context) map[string]string {
	result, err := rdbClient.HGetAll(ctx, keyWordMask).Result()
	if err != nil {
		logrus.GetLogger().CtxErrorf(ctx, "hgetall err: %v", err)
	}

	// return anyway
	return result
}
