package redis

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
	"wechat_robot/logrus"

	"github.com/redis/go-redis/v9"
)


type loggerHook struct {}

func NewRedisHook() redis.Hook {
	return new(loggerHook)
}

func (*loggerHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (*loggerHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		startTime := time.Now()

		err := next(ctx, cmd)

		costTime := float64(time.Since(startTime).Microseconds()) / 1000.0

		var args []string
		for _, arg := range cmd.Args() {
			args = append(args, fmt.Sprintf("%v", arg))
		}
		commands := strings.Join(args, " ")

		if err != nil && err != redis.Nil {
			logrus.GetLogger().CtxErrorf(ctx, "go-redis command fail: %s, err: %s, cost: %.3f", commands, err.Error(), costTime)
		} else {
			logrus.GetLogger().CtxInfof(ctx, "redis command: %s, cost: %.3f", commands, costTime)
		}

		return err
	}
}

func (*loggerHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		startTime := time.Now()

		err := next(ctx, cmds)

		costTime := float64(time.Since(startTime).Microseconds()) / 1000.0

		var cmdAggregation []string
		for _, cmd := range cmds {
			cmdAggregation = append(cmdAggregation, cmd.String())
		}

		if err != nil && err != redis.Nil {
			logrus.GetLogger().CtxErrorf(ctx, "pipeline fail: \n%s\n, err: %s, cost: %.3f", strings.Join(cmdAggregation, "\n"), err.Error(), costTime)
		} else {
			logrus.GetLogger().CtxInfof(ctx, "pipeline success: \n%s\n, cost: %.3f", strings.Join(cmdAggregation, "\n"), costTime)
		}

		return err
	}
}
