package baidu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"wechat_robot/chat/common"
	"wechat_robot/logrus"
	"wechat_robot/redis"
)

func init() {
	ctx := context.Background()

	accessInfo, err := refreshToken(ctx)
	if err != nil {
		panic(err)
	}

	setAccassInfo(accessInfo)

	go func() {
		timer := time.NewTimer(time.Duration(accessInfo.ExpiresIn) * time.Second)
		for range timer.C {
			accessInfo, err := refreshToken(ctx)
			if err != nil {
				timer.Reset(time.Second)
				continue
			}

			setAccassInfo(accessInfo)
			timer.Reset(time.Duration(accessInfo.ExpiresIn) * time.Second)
			logrus.GetLogger().Info("refresh token successfully")
		}
	}()
}

var accessInfo *AppAccessInfo

func setAccassInfo(info *AppAccessInfo) {
	accessInfo = info
}

func getAccessToken() string {
	if accessInfo == nil {
		return ""
	}

	return accessInfo.AccessToken
}

func refreshToken(ctx context.Context) (*AppAccessInfo, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, accessUrl, nil)
	if err != nil {
		panic(err)
	}

	param := req.URL.Query()
	param.Set("client_id", redis.GetErnieAppKey(ctx))
	param.Set("client_secret", redis.GetErnieAppSecret(ctx))
	req.URL.RawQuery = param.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.GetLogger().Errorf("http request err: %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.GetLogger().Errorf("read resp body err: %v", err)
		return nil, err
	}

	var accessInfo *AppAccessInfo
	if err = json.Unmarshal(data, &accessInfo); err != nil {
		logrus.GetLogger().Errorf("data invalid: %s", data)
		return nil, err
	}

	return accessInfo, nil
}

func SendMsg(ctx context.Context, msg string, histories ...*common.Message) (string, error) {
	var messages []*Message
	for _, his := range histories {
		messages = append(messages,
			&Message{
				Role:    roleUser,
				Content: his.RoleMsg,
			},
			&Message{
				Role:    roleAssistant,
				Content: his.AssistantMsg,
			})
	}

	messages = append(messages,
		&Message{
			Role:    roleUser,
			Content: msg,
		})

	chatReq := &ChatCreateReq{Messages: messages}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		logrus.GetLogger().CtxErrorf(ctx, "marshal err: %v", err)
		return "", nil
	}

	logrus.GetLogger().CtxDebugf(ctx, "ernie request: %s", reqBody)

	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, chatUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.GetLogger().Errorf("new request err: %v", err)
		return "", err
	}

	param := req.URL.Query()
	param.Set("access_token", getAccessToken())
	req.URL.RawQuery = param.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.GetLogger().CtxErrorf(ctx, "http client err: %v", err)
		return "", err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.GetLogger().CtxErrorf(ctx, "read resp body err: %v", err)
		return "", err
	}

	var chatResp ChatCreateResp
	if err = json.Unmarshal(respBody, &chatResp); err != nil {
		logrus.GetLogger().CtxErrorf(ctx, "unmarshal err: %v", err)
		return "", nil
	}

	if chatResp.ErrorCode != 0 {
		logrus.GetLogger().CtxErrorf(ctx, "chat request fails, :%s", respBody)
		return "", fmt.Errorf("%d:%s", chatResp.ErrorCode, chatResp.ErrorMsg)
	}
	if chatResp.Error != "" {
		logrus.GetLogger().CtxErrorf(ctx, "chat request fails, :%s", respBody)
		return "", fmt.Errorf("%s:%s", chatResp.Error, chatResp.ErrorDescription)
	}

	return chatResp.Result, nil
}
