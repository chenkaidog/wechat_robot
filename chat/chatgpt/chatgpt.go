package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"wechat_robot/chat/common"
	"wechat_robot/config"
	"wechat_robot/logrus"
)

func NewChatGPTgen(model string) *ChatGPTGen {
	return &ChatGPTGen{
		model: model,
	}
}

type ChatGPTGen struct {
	model string
}

func (c *ChatGPTGen) Reply2Chat(ctx context.Context, content string, history []*common.Message) (string, error) {
	return c.createChat(ctx, content, "", history...)
}

func (c *ChatGPTGen) createChat(ctx context.Context, msg, systemContent string, histories ...*common.Message) (string, error) {
	var messages []*Message

	if systemContent != "" {
		messages = append(
			messages,
			&Message{
				Role:    roleSystem,
				Content: systemContent,
			},
		)
	}

	for _, his := range histories {
		messages = append(
			messages,
			&Message{
				Role:    roleUser,
				Content: his.RoleMsg,
			},
			&Message{
				Role:    roleAssistant,
				Content: his.AssistantMsg,
			},
		)
	}

	messages = append(
		messages,
		&Message{
			Role:    roleUser,
			Content: msg,
		},
	)

	reqBody, err := json.Marshal(
		&MessageCreateReq{
			Model:    c.model,
			Messages: messages,
		})
	if err != nil {
		logrus.GetLogger().CtxErrorf(ctx, "marshal err: %v", err)
		return "", err
	}

	logrus.GetLogger().CtxDebugf(ctx, "req body: %s", reqBody)

	proxyUrl, err := url.Parse("http://localhost:7890")
	if err != nil {
		logrus.GetLogger().CtxErrorf(ctx, "parse proxy err: %v", err)
		return "", err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	req, err := http.NewRequest(http.MethodPost, msgUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.GetLogger().CtxErrorf(ctx, "new request err: %v", err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.GetChatGPT().ApiKey)

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

	var msgResp MessageCreateResp
	if err = json.Unmarshal(respBody, &msgResp); err != nil {
		logrus.GetLogger().CtxErrorf(ctx, "unmarshal err: %v", err)
		return "", nil
	}

	if len(msgResp.Choices) > 0 &&
		msgResp.Choices[0] != nil &&
		msgResp.Choices[0].Message != nil {
		return msgResp.Choices[0].Message.Content, nil
	}

	logrus.GetLogger().CtxErrorf(ctx, "resp body invalid: %s", respBody)
	return "sth wrong", nil
}
