package chatgpt

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatGPTGen_createChat(t *testing.T) {
	gen := NewChatGPTgen("gpt-3.5-turbo")
	reply, err := gen.Reply2Chat(context.Background(), "hello!", nil)
	assert.Nil(t, err)
	t.Log(reply)
}
