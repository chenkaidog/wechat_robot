package openwechat

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"wechat_robot/chat"
	"wechat_robot/chat/baidu"
	"wechat_robot/chat/chatgpt"
	"wechat_robot/chat/common"
	"wechat_robot/logrus"
	"wechat_robot/mysql"
	"wechat_robot/redis"

	"github.com/eatmoreapple/openwechat"
	"github.com/google/uuid"
)

func getRobotOption(ctx context.Context) *robotOption {
	result, err := redis.GetRobotOption(ctx)
	if err != nil || result == nil {
		return &robotOption{
			platform: common.PlatformBaidu,
			model:    common.ModelErnie,
		}
	}

	return &robotOption{
		platform: result.Platform,
		model:    result.Model,
	}
}

func handleText(ctx context.Context, msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		logrus.GetLogger().CtxErrorf(ctx, "get sender err: %v", err)
		return
	}

	data, _ := json.Marshal(sender)
	logrus.GetLogger().CtxInfof(ctx, "sender info: %s", data)

	var userMsg *userMsg
	switch {
	case sender.IsFriend():
		friend, _ := sender.AsFriend()
		userMsg = newUserMsgFromFriend(friend, msg)
	case sender.IsGroup():
		group, _ := sender.AsGroup()
		userMsg, err = newUserMsgFromGroup(group, msg)
		if err != nil {
			logrus.GetLogger().CtxErrorf(ctx, "newUserMsgFromGroup err: %v", err)
			return
		}
	case sender.IsMP():
		
	default:
		return
	}

	if err := userMsg.store(ctx); err != nil {
		return
	}
	if !userMsg.IsInteract {
		return
	}

	replyContent, err := userMsg.generateReply(ctx)
	if err != nil {
		return
	}

	if _, err = msg.ReplyText(replyContent); err != nil {
		logrus.GetLogger().CtxErrorf(ctx, "reply text err: %v", err)
		return
	}
}

func newUserMsgFromFriend(friend *openwechat.Friend, msg *openwechat.Message) *userMsg {
	return &userMsg{
		MsgId:      msg.MsgId,
		MsgType:    msg.MsgType.String(),
		Content:    msg.Content,
		ChatId:     friend.UserName,
		ChatName:   friend.String(),
		SenderId:   friend.UserName,
		SenderName: friend.String(),
		IsInteract: true,
		IsStar:     friend.StarFriend == 1,
	}
}

func newUserMsgFromGroup(group *openwechat.Group, msg *openwechat.Message) (*userMsg, error) {
	friend, err := msg.SenderInGroup()
	if err != nil {
		return nil, err
	}

	content := msg.Content
	for k, v := range redis.GetWordMask(msg.Context()) {
		content = strings.ReplaceAll(content, k, v)
	}

	content = fmt.Sprintf("%s\n asked by: %s", content, friend.DisplayName)

	return &userMsg{
		MsgId:      msg.MsgId,
		MsgType:    msg.MsgType.String(),
		Content:    content,
		ChatId:     group.UserName,
		ChatName:   group.String(),
		SenderId:   friend.UserName,
		SenderName: friend.String(),
		IsInteract: msg.IsAt(),
		IsStar:     false,
	}, nil
}

type robotOption struct {
	platform string
	model    string
}

type userMsg struct {
	MsgId      string
	ChatId     string
	ChatName   string
	SenderId   string
	SenderName string
	MsgType    string
	Content    string
	IsInteract bool
	IsStar     bool
}

func (u *userMsg) generateReply(ctx context.Context) (string, error) {
	handleResp, isCmd := u.handleCommand(ctx)
	if isCmd {
		return handleResp, nil
	}

	historyMsgs, err := u.loadHistoryMsg(ctx)
	if err != nil {
		return "", err
	}

	robotOpt := getRobotOption(ctx)
	replyContent, err := chat.NewChatGenerator(robotOpt.platform, robotOpt.model).
		Reply2Chat(ctx, u.Content, historyMsgs)
	if err != nil {
		replyContent = "sth wrong"
	}

	_ = u.storeReply(ctx, robotOpt.platform, robotOpt.model, replyContent)

	return replyContent, nil
}

func (u *userMsg) handleCommand(ctx context.Context) (string, bool) {
	content := u.Content

	if strings.EqualFold(content, "ping") {
		return "pong", true
	}

	if !u.IsStar {
		return "", false
	}

	logrus.GetLogger().CtxInfof(ctx, "execute cmd: [%s]", content)

	if strings.HasPrefix(content, "set robot ") {
		if cmds := strings.Split(content, " "); len(cmds) == 4 {
			platform, model := cmds[2], cmds[3]
			if err := redis.SetRobotOption(ctx, &redis.RobotOption{
				Platform: platform,
				Model:    model,
			}); err != nil {
				return err.Error(), true
			} else {
				return fmt.Sprintf("set %s %s successfully", platform, model), true
			}
		}
	} else if strings.HasPrefix(content, "set words ") {
		if cmds := strings.Split(content, " "); len(cmds) == 4 {
			word, mask := cmds[2], cmds[3]
			if err := redis.SetWordMask(ctx, word, mask); err != nil {
				return err.Error(), true
			} else {
				return fmt.Sprintf("set word %s => %s successfully", word, mask), true
			}
		}
	} else if strings.HasPrefix(content, "del words ") {
		if cmds := strings.Split(content, " "); len(cmds) == 3 {
			word := cmds[2]
			if err := redis.DeleteWordMask(ctx, word); err != nil {
				return err.Error(), true
			} else {
				return fmt.Sprintf("delete %s successfully", word), true
			}
		}
	} else if strings.HasPrefix(content, "ernie set app_key ") {
		if cmds := strings.Split(content, " "); len(cmds) == 4 {
			if err := baidu.SetErnieAppKey(ctx, cmds[3]); err != nil {
				return err.Error(), true
			} else {
				return "set app_key successfully", true
			}
		}
	} else if strings.HasPrefix(content, "ernie set app_secret ") {
		if cmds := strings.Split(content, " "); len(cmds) == 4 {
			if err := baidu.SetErnieAppSecret(ctx, cmds[3]); err != nil {
				return err.Error(), true
			} else {
				return "set app_secret successfully", true
			}
		}
	} else if strings.EqualFold(content, "ernie refresh token") {
		if err := baidu.RefreshToken(ctx); err != nil {
			return err.Error(), true
		} else {
			return "refresh successfully", true
		}
	} else if strings.HasPrefix(content, "chatgpt set api_key ") {
		if cmds := strings.Split(content, " "); len(cmds) == 4 {
			if err := chatgpt.SetChatgptApiKey(ctx, cmds[3]); err != nil {
				return err.Error(), true
			} else {
				return "set app_secret successfully", true
			}
		}
	}

	return "", false
}

func (u *userMsg) loadHistoryMsg(ctx context.Context) ([]*common.Message, error) {
	results, err := mysql.SelectHistoryMsg(ctx, u.ChatId, 5, -time.Hour)
	if err != nil {
		return nil, err
	}

	var msgs []*common.Message
	for _, result := range results {
		msgs = append(
			msgs,
			&common.Message{
				RoleMsg:      result.UserContent,
				AssistantMsg: result.RepliedContent,
			})
	}

	return msgs, nil
}

func (u *userMsg) store(ctx context.Context) error {
	if u.MsgId == "" {
		u.MsgId = uuid.NewString()
	}

	return mysql.InsertUserMsg(
		ctx,
		&mysql.UserMsg{
			MsgID:      u.MsgId,
			ChatID:     u.ChatId,
			ChatName:   u.ChatName,
			SenderID:   u.SenderId,
			SenderName: u.SenderName,
			MsgType:    u.MsgType,
			Content:    u.Content,
			Interact:   u.IsInteract,
		})
}

func (u *userMsg) storeReply(ctx context.Context, platform, model, content string) error {
	return mysql.InsertRepliedMsg(
		ctx,
		&mysql.RepliedMsg{
			MsgID:     uuid.NewString(),
			UserMsgID: u.MsgId,
			Platform:  platform,
			Model:     model,
			Content:   content,
		},
	)
}
