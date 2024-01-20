package baidu_test

import (
	"context"
	"testing"
	"wechat_robot/chat/baidu"

	"github.com/stretchr/testify/assert"
)

func TestSendMsg(t *testing.T) {
	reply, err := baidu.SendMsg(context.Background(), "hello!")
	assert.Nil(t, err)
	t.Log(reply)
}
