package manusai

type TaskOptions struct {
	AgentProfile        string        `json:"agentProfile,omitempty"`
	TaskMode            string        `json:"taskMode,omitempty"`
	Locale              string        `json:"locale,omitempty"`
	HideInTaskList      *bool         `json:"hideInTaskList,omitempty"`
	CreateShareableLink *bool         `json:"createShareableLink,omitempty"`
	Attachments         []interface{} `json:"attachments,omitempty"`
}

type TaskResponse struct {
	TaskID    string `json:"task_id"`
	TaskTitle string `json:"task_title"`
	TaskURL   string `json:"task_url"`
}

type TaskFilters struct {
	After         string   `json:"after,omitempty"`
	Limit         int      `json:"limit,omitempty"`
	Order         string   `json:"order,omitempty"`
	OrderBy       string   `json:"orderBy,omitempty"`
	Query         string   `json:"query,omitempty"`
	Status        []string `json:"status,omitempty"`
	CreatedAfter  string   `json:"createdAfter,omitempty"`
	CreatedBefore string   `json:"createdBefore,omitempty"`
}

type TaskListResponse struct {
	Data    []TaskSummary `json:"data"`
	HasMore bool          `json:"has_more"`
}

type TaskSummary struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TaskDetail struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Status      string        `json:"status"`
	CreditUsage float64       `json:"credit_usage"`
	Output      []TaskMessage `json:"output"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

type TaskMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TaskUpdate struct {
	Title                    *string `json:"title,omitempty"`
	EnableShared             *bool   `json:"enableShared,omitempty"`
	EnableVisibleInTaskList  *bool   `json:"enableVisibleInTaskList,omitempty"`
}

type DeleteResponse struct {
	Deleted bool `json:"deleted"`
}

type FileResponse struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	UploadURL string `json:"upload_url"`
	Status    string `json:"status"`
}

type FileListResponse struct {
	Data []FileDetail `json:"data"`
}

type FileDetail struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	Status    string `json:"status"`
	SizeBytes int64  `json:"size_bytes,omitempty"`
	CreatedAt string `json:"created_at"`
}

type WebhookConfig struct {
	URL    string   `json:"url"`
	Events []string `json:"events,omitempty"`
}

type WebhookResponse struct {
	WebhookID string `json:"webhook_id"`
}

type WebhookPayload struct {
	EventType  string                 `json:"event_type"`
	TaskDetail map[string]interface{} `json:"task_detail,omitempty"`
}

type TaskAttachment struct {
	Type     string `json:"type"`
	FileID   string `json:"file_id,omitempty"`
	URL      string `json:"url,omitempty"`
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}
