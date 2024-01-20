package chatgpt


const (
	msgUrl = "https://api.openai.com/v1/chat/completions"
)

const (
	roleAssistant = "assistant"
	roleUser      = "user"
	roleSystem    = "system"
)

type MessageCreateReq struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MessageCreateResp struct {
	// ID      string `json:"id"`
	// Object  string `json:"object"`
	// Created int    `json:"created"`
	// Model   string `json:'model"`
	Choices []*Choice `json:"choices"`
}

type Choice struct {
	// Index int `json:"index"`
	Message *Message `json:"message"`
}
