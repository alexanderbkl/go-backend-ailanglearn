package helpers

import (
	"encoding/json"
	"go-backend-ailanglearn/configs"
)

const COMPLETIONS_URL = "https://api.openai.com/v1/chat/completions"

var openAIKey string = configs.EnvOpenAIKey()

var client = NewClient(openAIKey, "")


type CreateCompletionsRequest struct {
	Model            string            `json:"model,omitempty"`
	Messages         []Message         `json:"messages,omitempty"`
	Prompt           StrArray          `json:"prompt,omitempty"`
	Suffix           string            `json:"suffix,omitempty"`
	MaxTokens        int               `json:"max_tokens,omitempty"`
	Temperature      float64           `json:"temperature,omitempty"`
	TopP             float64           `json:"top_p,omitempty"`
	N                int               `json:"n,omitempty"`
	Stream           bool              `json:"stream,omitempty"`
	LogProbs         int               `json:"logprobs,omitempty"`
	Echo             bool              `json:"echo,omitempty"`
	Stop             StrArray          `json:"stop,omitempty"`
	PresencePenalty  float64           `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64           `json:"frequency_penalty,omitempty"`
	BestOf           int               `json:"best_of,omitempty"`
	LogitBias        map[string]string `json:"logit_bias,omitempty"`
	User             string            `json:"user,omitempty"`
}

func CreateCompletionMessage(content string) (response CreateCompletionsResponse, err error){
	
	r := CreateCompletionsRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a language learning assistant, you will correct, give tips, or simulate a conversation accordingly, using maximum 100 words per answer. Answer really briefly.",
			},
			{
				Role:    "user",
				Content: "(Spanish learning Catalan) yo tingo vinte ans",
			},
			{
				Role:    "assistant",
				Content: "En catalán se dice \"tinc vint anys\". La conjugación correcta del verbo \"tener\" en catalán es \"tenir\".",
			},
			{
				Role:    "user",
				Content: content,
			},
		},
		Temperature: 0.7,
		MaxTokens:   100,
	}

	completions, err := client.CreateCompletions(r)
	if err != nil {
		panic(err)
	}

	return completions, err
}


func (c *Client) CreateCompletionsRaw(r CreateCompletionsRequest) ([]byte, error) {
	return c.Post(COMPLETIONS_URL, r)
}

func (c *Client) CreateCompletions(r CreateCompletionsRequest) (response CreateCompletionsResponse, err error) {
	raw, err := c.CreateCompletionsRaw(r)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(raw, &response)
	return response, err
}

type CreateCompletionsResponse struct {
	ID      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	Created int    `json:"created,omitempty"`
	Model   string `json:"model,omitempty"`
	Choices []struct {
		Message struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		}
		Index        int         `json:"index,omitempty"`
		FinishReason string      `json:"finish_reason,omitempty"`
	} `json:"choices,omitempty"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens,omitempty"`
		CompletionTokens int `json:"completion_tokens,omitempty"`
		TotalTokens      int `json:"total_tokens,omitempty"`
	} `json:"usage,omitempty"`

	Error Error `json:"error,omitempty"`
}

// Error is the error standard response from the API
type Error struct {
	Message string      `json:"message,omitempty"`
	Type    string      `json:"type,omitempty"`
	Param   interface{} `json:"param,omitempty"`
	Code    interface{} `json:"code,omitempty"`
}

type Message struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}