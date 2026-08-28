package manusai

// TaskOptions configures a task created with Client.CreateTask.
type TaskOptions struct {
	AgentProfile    string        `json:"agent_profile,omitempty"`
	Locale          string        `json:"locale,omitempty"`
	HideInTaskList  *bool         `json:"hide_in_task_list,omitempty"`
	ShareVisibility string        `json:"share_visibility,omitempty"`
	Title           string        `json:"title,omitempty"`
	ProjectID       string        `json:"project_id,omitempty"`
	EnableAskUser   *bool         `json:"enable_ask_user,omitempty"`
	Connectors      []string      `json:"connectors,omitempty"`
	EnableSkills    []string      `json:"enable_skills,omitempty"`
	ForceSkills     []string      `json:"force_skills,omitempty"`
	Attachments     []interface{} `json:"attachments,omitempty"`
}

// TaskResponse is the API response returned after creating a task.
type TaskResponse struct {
	OK              bool   `json:"ok"`
	RequestID       string `json:"request_id"`
	TaskID          string `json:"task_id"`
	TaskTitle       string `json:"task_title"`
	TaskURL         string `json:"task_url"`
	ShareURL        string `json:"share_url,omitempty"`
	ShareVisibility string `json:"share_visibility,omitempty"`
}

// TaskFilters narrows and paginates a task list request.
type TaskFilters struct {
	Cursor    string `json:"cursor,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Order     string `json:"order,omitempty"`
	Scope     string `json:"scope,omitempty"`
	AgentID   string `json:"agent_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

// TaskListResponse is a paginated list of task summaries.
type TaskListResponse struct {
	OK         bool          `json:"ok"`
	RequestID  string        `json:"request_id"`
	Tasks      []TaskSummary `json:"tasks"`
	HasMore    bool          `json:"has_more"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

// TaskSummary contains the list-view metadata for a task.
type TaskSummary struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	AgentStatus     string `json:"agent_status"`
	ShareVisibility string `json:"share_visibility"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

// TaskDetail contains detailed metadata for a task.
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

// TaskMessage represents one event or message in a task conversation.
type TaskMessage struct {
	ID               string                 `json:"id"`
	Type             string                 `json:"type"`
	Timestamp        int64                  `json:"timestamp"`
	UserMessage      map[string]interface{} `json:"user_message,omitempty"`
	AssistantMessage map[string]interface{} `json:"assistant_message,omitempty"`
	ErrorMessage     map[string]interface{} `json:"error_message,omitempty"`
	StatusUpdate     map[string]interface{} `json:"status_update,omitempty"`
}

// TaskUpdate specifies task properties to change.
type TaskUpdate struct {
	Title           *string `json:"title,omitempty"`
	ShareVisibility *string `json:"share_visibility,omitempty"`
	HideInTaskList  *bool   `json:"hide_in_task_list,omitempty"`
}

// DeleteResponse is returned after deleting a task or file.
type DeleteResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
	Deleted   bool   `json:"deleted"`
}

// FileResponse contains the upload information for a newly created file.
type FileResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
	FileID    string `json:"file_id"`
	Filename  string `json:"filename"`
	UploadURL string `json:"upload_url"`
	ExpiresAt int64  `json:"expires_at"`
}

// FileListResponse is a paginated list of files.
type FileListResponse struct {
	OK         bool         `json:"ok"`
	RequestID  string       `json:"request_id"`
	Files      []FileDetail `json:"files"`
	HasMore    bool         `json:"has_more"`
	NextCursor string       `json:"next_cursor,omitempty"`
}

// FileDetail contains metadata for a Manus file.
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

// WebhookConfig specifies the endpoint and events for a webhook.
type WebhookConfig struct {
	URL    string   `json:"url"`
	Events []string `json:"events,omitempty"`
}

// WebhookResponse is returned after registering a webhook.
type WebhookResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
	WebhookID string `json:"webhook_id"`
}

// WebhookPayload is the decoded payload sent to a webhook endpoint.
type WebhookPayload struct {
	EventType  string                 `json:"event_type"`
	TaskDetail map[string]interface{} `json:"task_detail,omitempty"`
}

// TaskAttachment represents a file attachment accepted by Manus task messages.
type TaskAttachment struct {
	Type     string `json:"type"`
	FileID   string `json:"file_id,omitempty"`
	FileURL  string `json:"file_url,omitempty"`
	FileData string `json:"file_data,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}

// TaskMessagesResponse is a paginated list of messages for a task.
type TaskMessagesResponse struct {
	OK         bool          `json:"ok"`
	RequestID  string        `json:"request_id"`
	TaskID     string        `json:"task_id"`
	Messages   []TaskMessage `json:"messages"`
	HasMore    bool          `json:"has_more"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

// SendMessageResponse is returned after sending a task message.
type SendMessageResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
}

// StopTaskResponse is returned after requesting that a task stop.
type StopTaskResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
}

// ConfirmActionResponse is returned after confirming a pending action.
type ConfirmActionResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
}
