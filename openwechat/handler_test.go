package openwechat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_userMsg_handleCommand(t *testing.T) {
	t.Run("ping", func(t *testing.T) {
		userMs := &userMsg{
			Content: "ping",
		}

		result, ok := userMs.handleCommand(context.Background())
		assert.True(t, ok)
		assert.Equal(t, "pong", result)
	})

	t.Run("cmd", func(t *testing.T) {
		userMs := &userMsg{
			IsStar: true,
			Content: "set robot chatgpt gpt-3-turbo",
		}

		result, ok := userMs.handleCommand(context.Background())
		assert.True(t, ok)
		t.Log(result)
	})
}
