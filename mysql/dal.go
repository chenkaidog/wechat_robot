package mysql

import (
	"context"
	"time"
)

func InsertUserMsg(ctx context.Context, msg *UserMsg) error {
	return gormDB.WithContext(ctx).Create(msg).Error
}

func InsertRepliedMsg(ctx context.Context, msg *RepliedMsg) error {
	return gormDB.WithContext(ctx).Create(msg).Error
}

type Result struct {
	UserContent    string `gorm:"column:user_content"`
	RepliedContent string `gorm:"column:replied_content"`
}

func SelectHistoryMsg(ctx context.Context, chatId string, limit int, traceback time.Duration) ([]*Result, error) {
	var results []*Result

	err := gormDB.WithContext(ctx).
		Model(&UserMsg{}).
		Select("user_msg.content AS user_content, replied_msg.content AS replied_content").
		Joins("JOIN replied_msg ON user_msg.msg_id=replied_msg.user_msg_id").
		Where("user_msg.chat_id", chatId).
		Where("user_msg.interact", true).
		Where("user_msg.created_at >= ?", time.Now().Add(traceback)).
		Order("user_msg.created_at DESC").
		Limit(limit).
		Scan(&results).
		Error
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(results)/2; i++ {
		results[i], results[len(results)-i-1] = results[len(results)-i-1], results[i]
	}

	return results, nil
}
