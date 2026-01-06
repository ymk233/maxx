package converter

// Claude API types

type ClaudeRequest struct {
	Model         string            `json:"model"`
	Messages      []ClaudeMessage   `json:"messages"`
	System        interface{}       `json:"system,omitempty"` // string or []SystemBlock
	MaxTokens     int               `json:"max_tokens,omitempty"`
	Temperature   *float64          `json:"temperature,omitempty"`
	TopP          *float64          `json:"top_p,omitempty"`
	TopK          *int              `json:"top_k,omitempty"`
	StopSequences []string          `json:"stop_sequences,omitempty"`
	Stream        bool              `json:"stream,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	Tools         []ClaudeTool      `json:"tools,omitempty"`
	ToolChoice    interface{}       `json:"tool_choice,omitempty"`
}

type ClaudeMessage struct {
	Role    string               `json:"role"`
	Content interface{}          `json:"content"` // string or []ContentBlock
}

type ClaudeContentBlock struct {
	Type      string      `json:"type"`
	Text      string      `json:"text,omitempty"`
	ID        string      `json:"id,omitempty"`
	Name      string      `json:"name,omitempty"`
	Input     interface{} `json:"input,omitempty"`
	ToolUseID string      `json:"tool_use_id,omitempty"`
	Content   string      `json:"content,omitempty"`
}

type ClaudeTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	InputSchema interface{} `json:"input_schema"`
}

type ClaudeResponse struct {
	ID           string               `json:"id"`
	Type         string               `json:"type"`
	Role         string               `json:"role"`
	Content      []ClaudeContentBlock `json:"content"`
	Model        string               `json:"model"`
	StopReason   string               `json:"stop_reason"`
	StopSequence string               `json:"stop_sequence,omitempty"`
	Usage        ClaudeUsage          `json:"usage"`
}

type ClaudeUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens,omitempty"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
}

// Claude streaming events
type ClaudeStreamEvent struct {
	Type         string               `json:"type"`
	Message      *ClaudeResponse      `json:"message,omitempty"`
	Index        int                  `json:"index,omitempty"`
	ContentBlock *ClaudeContentBlock  `json:"content_block,omitempty"`
	Delta        *ClaudeStreamDelta   `json:"delta,omitempty"`
	Usage        *ClaudeUsage         `json:"usage,omitempty"`
}

type ClaudeStreamDelta struct {
	Type         string `json:"type,omitempty"`
	Text         string `json:"text,omitempty"`
	PartialJSON  string `json:"partial_json,omitempty"`
	StopReason   string `json:"stop_reason,omitempty"`
	StopSequence string `json:"stop_sequence,omitempty"`
}
