package common

const (
	PlatformBaidu      = "baidu"
	PlatformChatGPT    = "chatgpt"
	ModelErnie         = "ernie"
	ModelChatgpt3turbo = "gpt-3.5-turbo"
	ModelChatgpt4      = "gpt-4"
)

type Message struct {
	RoleMsg      string
	AssistantMsg string
}
