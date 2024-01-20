package logrus

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type logrusHook struct{}

type LogID struct{}
type MsgID struct{}

func Newcontext(ctx context.Context, msgId string) context.Context {
	ctx = context.WithValue(ctx, LogID{}, uuid.NewString())
	ctx = context.WithValue(ctx, MsgID{}, msgId)

	return ctx
}

func newLogrusHook() *logrusHook {
	return new(logrusHook)
}

// Fire implements logrus.Hook.
func (h *logrusHook) Fire(entry *logrus.Entry) error {
	if entry != nil && entry.Context != nil {
		logId, ok := entry.Context.Value(LogID{}).(string)
		if ok {
			entry.Data["log_id"] = logId
		}

		msgId, ok := entry.Context.Value(MsgID{}).(string)
		if ok {
			entry.Data["msg_id"] = msgId
		}
	}

	return nil
}

// Levels implements logrus.Hook.
func (h *logrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
