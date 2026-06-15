package manusai

type TaskOptions struct {
	AgentProfile     string        `json:"agent_profile,omitempty"`
	Locale           string        `json:"locale,omitempty"`
	HideInTaskList   *bool         `json:"hide_in_task_list,omitempty"`
	ShareVisibility  string        `json:"share_visibility,omitempty"`
	Title            string        `json:"title,omitempty"`
	ProjectID        string        `json:"project_id,omitempty"`
	EnableAskUser    *bool         `json:"enable_ask_user,omitempty"`
	Connectors       []string      `json:"connectors,omitempty"`
	EnableSkills     []string      `json:"enable_skills,omitempty"`
	ForceSkills      []string      `json:"force_skills,omitempty"`
	Attachments      []interface{} `json:"attachments,omitempty"`
}

type TaskResponse struct {
	OK              bool   `json:"ok"`
	RequestID       string `json:"request_id"`
	TaskID          string `json:"task_id"`
	TaskTitle       string `json:"task_title"`
	TaskURL         string `json:"task_url"`
	ShareURL        string `json:"share_url,omitempty"`
	ShareVisibility string `json:"share_visibility,omitempty"`
}

type TaskFilters struct {
	Cursor    string `json:"cursor,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Order     string `json:"order,omitempty"`
	Scope     string `json:"scope,omitempty"`
	AgentID   string `json:"agent_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

type TaskListResponse struct {
	OK         bool          `json:"ok"`
	RequestID  string        `json:"request_id"`
	Tasks      []TaskSummary `json:"tasks"`
	HasMore    bool          `json:"has_more"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

type TaskSummary struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	AgentStatus     string `json:"agent_status"`
	ShareVisibility string `json:"share_visibility"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

type TaskDetail struct {
	OK              bool    `json:"ok"`
	RequestID       string  `json:"request_id"`
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	AgentStatus     string  `json:"agent_status"`
	ShareVisibility string  `json:"share_visibility"`
	CreditUsage     float64 `json:"credit_usage"`
	CreatedAt       int64   `json:"created_at"`
	UpdatedAt       int64   `json:"updated_at"`
}

type TaskMessage struct {
	ID               string                 `json:"id"`
	Type             string                 `json:"type"`
	Timestamp        int64                  `json:"timestamp"`
	UserMessage      map[string]interface{} `json:"user_message,omitempty"`
	AssistantMessage map[string]interface{} `json:"assistant_message,omitempty"`
	ErrorMessage     map[string]interface{} `json:"error_message,omitempty"`
	StatusUpdate     map[string]interface{} `json:"status_update,omitempty"`
}

type TaskUpdate struct {
	Title           *string `json:"title,omitempty"`
	ShareVisibility *string `json:"share_visibility,omitempty"`
	HideInTaskList  *bool   `json:"hide_in_task_list,omitempty"`
}

type DeleteResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
	Deleted   bool   `json:"deleted"`
}

type FileResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
	FileID    string `json:"file_id"`
	Filename  string `json:"filename"`
	UploadURL string `json:"upload_url"`
	ExpiresAt int64  `json:"expires_at"`
}

type FileListResponse struct {
	OK         bool         `json:"ok"`
	RequestID  string       `json:"request_id"`
	Files      []FileDetail `json:"files"`
	HasMore    bool         `json:"has_more"`
	NextCursor string       `json:"next_cursor,omitempty"`
}

type FileDetail struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id,omitempty"`
	FileID    string `json:"file_id"`
	Filename  string `json:"filename"`
	Status    string `json:"status"`
	SizeBytes int64  `json:"size_bytes,omitempty"`
	CreatedAt int64  `json:"created_at"`
	ExpiresAt int64  `json:"expires_at"`
}

type WebhookConfig struct {
	URL    string   `json:"url"`
	Events []string `json:"events,omitempty"`
}

type WebhookResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
	WebhookID string `json:"webhook_id"`
}

type WebhookPayload struct {
	EventType  string                 `json:"event_type"`
	TaskDetail map[string]interface{} `json:"task_detail,omitempty"`
}

type TaskAttachment struct {
	Type     string `json:"type"`
	FileID   string `json:"file_id,omitempty"`
	FileURL  string `json:"file_url,omitempty"`
	FileData string `json:"file_data,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}

type TaskMessagesResponse struct {
	OK         bool          `json:"ok"`
	RequestID  string        `json:"request_id"`
	TaskID     string        `json:"task_id"`
	Messages   []TaskMessage `json:"messages"`
	HasMore    bool          `json:"has_more"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

type SendMessageResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
}

type StopTaskResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
}

type ConfirmActionResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
}
