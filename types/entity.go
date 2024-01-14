package types

const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)

type ChatModelMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type ChatCompletionRequest struct {
	Model       string             `json:"model"`
	Messages    []ChatModelMessage `json:"messages"`
	MaxTokens   int                `json:"max_tokens,omitempty"`
	Temperature float32            `json:"temperature,omitempty"`
	Stream      bool               `json:"stream,omitempty"`
	User        string             `json:"user,omitempty"`
}

// Chat Completion
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionChoice struct {
	Index        int              `json:"index"`
	Message      ChatModelMessage `json:"message"`
	FinishReason string           `json:"finish_reason"`
}

type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   Usage                  `json:"usage"`
}

// Stream Chat Completion
type ChatCompletionStreamChoiceDelta struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`
}

type ChatCompletionStreamChoice struct {
	Index        int                             `json:"index"`
	Delta        ChatCompletionStreamChoiceDelta `json:"delta"`
	FinishReason string                          `json:"finish_reason"`
}

type ChatCompletionStreamResponse struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int64                        `json:"created"`
	Model   string                       `json:"model"`
	Choices []ChatCompletionStreamChoice `json:"choices"`
	Usage   *Usage                       `json:"usage"`
}

// Embedding
type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type Embedding []float32

// Green check for Chinese mainland
type WordCheckData struct {
	Keyword  string `json:"keyword"`
	Category string `json:"category"`
	Level    string `json:"level"`
	Position string `json:"position"`
}

type WordCheckResponse struct {
	Code      string
	Msg       string
	ReturnStr string          `json:"return_str"`
	WordList  []WordCheckData `json:"word_list"`
}

// Document
type Metadata struct {
	Corpus string `json:"corpus,omitempty"`
	Url    string `json:"url,omitempty"`
	Source string `json:"source,omitempty"`
}

type Document struct {
	ID       string   `json:"id,omitempty"`
	Text     string   `json:"text,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

type Chunk struct {
	ID         string    `json:"id,omitempty"`
	Text       string    `json:"text,omitempty"`
	DocumentID string    `json:"document_id,omitempty"`
	Metadata   Metadata  `json:"metadata,omitempty"`
	Embedding  Embedding `json:"embedding,omitempty"`
}
