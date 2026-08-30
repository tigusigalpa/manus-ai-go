package manusai

// TaskOptions configures a task created with Client.CreateTask.
type TaskOptions struct {
	AgentProfile    string `json:"agent_profile,omitempty"`
	Locale          string `json:"locale,omitempty"`
	HideInTaskList  *bool  `json:"hide_in_task_list,omitempty"`
	ShareVisibility string `json:"share_visibility,omitempty"`
	Title           string `json:"title,omitempty"`
	ProjectID       string `json:"project_id,omitempty"`
	// InteractiveMode lets the agent ask follow-up questions when needed.
	InteractiveMode        *bool                  `json:"interactive_mode,omitempty"`
	StructuredOutputSchema map[string]interface{} `json:"structured_output_schema,omitempty"`
	// EnableAskUser is retained for compatibility. Prefer InteractiveMode.
	EnableAskUser *bool         `json:"-"`
	Connectors    []string      `json:"connectors,omitempty"`
	EnableSkills  []string      `json:"enable_skills,omitempty"`
	ForceSkills   []string      `json:"force_skills,omitempty"`
	Attachments   []interface{} `json:"attachments,omitempty"`
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
	Cursor        string `json:"cursor,omitempty"`
	Limit         int    `json:"limit,omitempty"`
	Order         string `json:"order,omitempty"`
	Scope         string `json:"scope,omitempty"`
	AgentID       string `json:"agent_id,omitempty"`
	ProjectID     string `json:"project_id,omitempty"`
	OAuthClientID string `json:"oauth_client_id,omitempty"`
	APIKeyID      string `json:"api_key_id,omitempty"`
}

// TaskListResponse is a paginated list of task summaries.
type TaskListResponse struct {
	OK         bool          `json:"ok"`
	RequestID  string        `json:"request_id"`
	Tasks      []TaskSummary `json:"data"`
	HasMore    bool          `json:"has_more"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

// TaskSummary contains the list-view metadata for a task.
type TaskSummary struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	AgentStatus     string  `json:"status"`
	ShareVisibility string  `json:"share_visibility"`
	CreditUsage     float64 `json:"credit_usage"`
	TaskURL         string  `json:"task_url,omitempty"`
	TaskType        string  `json:"task_type,omitempty"`
	AgentProfile    string  `json:"agent_profile,omitempty"`
	CreatedAt       int64   `json:"created_at"`
	UpdatedAt       int64   `json:"updated_at"`
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
	TaskURL         string  `json:"task_url,omitempty"`
	TaskType        string  `json:"task_type,omitempty"`
	AgentProfile    string  `json:"agent_profile,omitempty"`
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
	// EnableVisibleInTaskList controls whether the task appears in the Manus task list.
	EnableVisibleInTaskList *bool `json:"enable_visible_in_task_list,omitempty"`
	// HideInTaskList is retained for compatibility. Prefer EnableVisibleInTaskList.
	HideInTaskList *bool `json:"-"`
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

// WebhookConfig specifies the endpoint for a webhook.
type WebhookConfig struct {
	URL string `json:"url"`
	// Events is retained for compatibility and is ignored: the v2 API does not accept it.
	Events []string `json:"-"`
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
	TaskID    string `json:"task_id"`
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
	TaskID    string `json:"task_id"`
	Confirmed bool   `json:"confirmed"`
}

// Project represents a Manus project.
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Instruction string `json:"instruction,omitempty"`
	CreatedAt   int64  `json:"created_at"`
}

// Skill represents an available Manus skill.
type Skill struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	OwnerType   string         `json:"owner_type,omitempty"`
	CreatorInfo map[string]any `json:"creator_info,omitempty"`
	CreatedAt   int64          `json:"created_at"`
	UpdatedAt   int64          `json:"updated_at"`
}

// Agent represents a custom Manus agent.
type Agent struct {
	ID       string `json:"id"`
	TaskID   string `json:"task_id"`
	Nickname string `json:"nickname,omitempty"`
	About    string `json:"about,omitempty"`
}

// Webhook contains registered webhook metadata.
type Webhook struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Status    string `json:"status,omitempty"`
	CreatedAt int64  `json:"created_at"`
}

// BrowserClient is an online Manus browser client.
type BrowserClient struct {
	ClientID   string `json:"client_id"`
	ClientName string `json:"client_name,omitempty"`
	UserAgent  string `json:"ua,omitempty"`
}

// Connector describes an installed connector.
type Connector struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
}

// ProjectResponse is returned after creating a project.
type ProjectResponse struct {
	OK        bool    `json:"ok"`
	RequestID string  `json:"request_id"`
	Project   Project `json:"project"`
}

// ProjectListResponse contains projects available to the authenticated user.
type ProjectListResponse struct {
	OK        bool      `json:"ok"`
	RequestID string    `json:"request_id"`
	Projects  []Project `json:"data"`
}

// SkillListResponse contains available skills.
type SkillListResponse struct {
	OK        bool    `json:"ok"`
	RequestID string  `json:"request_id"`
	Skills    []Skill `json:"data"`
}

// AgentListResponse contains custom agents.
type AgentListResponse struct {
	OK        bool    `json:"ok"`
	RequestID string  `json:"request_id"`
	Agents    []Agent `json:"data"`
}

// AgentResponse contains one agent.
type AgentResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
	Agent     Agent  `json:"agent"`
}

// WebhookListResponse contains registered webhooks.
type WebhookListResponse struct {
	OK        bool      `json:"ok"`
	RequestID string    `json:"request_id"`
	Webhooks  []Webhook `json:"data"`
}

// BrowserClientListResponse contains online browser clients.
type BrowserClientListResponse struct {
	OK        bool            `json:"ok"`
	RequestID string          `json:"request_id"`
	Clients   []BrowserClient `json:"data"`
}

// WebhookPublicKeyResponse contains the public key used for webhook signature verification.
type WebhookPublicKeyResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
	PublicKey string `json:"public_key"`
	Algorithm string `json:"algorithm"`
}

// ConnectorListResponse contains installed connectors.
type ConnectorListResponse struct {
	OK         bool        `json:"ok"`
	RequestID  string      `json:"request_id"`
	Connectors []Connector `json:"data"`
}

// UsageRecord describes credit usage for one task or session.
type UsageRecord struct {
	TaskID           string  `json:"task_id"`
	Title            string  `json:"title,omitempty"`
	Credits          float64 `json:"credits"`
	CreatedAt        int64   `json:"created_at"`
	Type             string  `json:"type,omitempty"`
	CollaborateInfos []any   `json:"collaborate_infos,omitempty"`
}

// UsageListResponse contains paginated credit changes.
type UsageListResponse struct {
	OK         bool          `json:"ok"`
	RequestID  string        `json:"request_id"`
	Records    []UsageRecord `json:"data"`
	HasMore    bool          `json:"has_more"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

// DailyUsageStatistic contains a team's credit usage for one date.
type DailyUsageStatistic struct {
	Date    string  `json:"date"`
	Credits float64 `json:"credits"`
}

// TeamUsageStatisticResponse contains daily team usage data.
type TeamUsageStatisticResponse struct {
	OK        bool                  `json:"ok"`
	RequestID string                `json:"request_id"`
	Records   []DailyUsageStatistic `json:"data"`
}

// TeamUsageLog contains one user's aggregated team credit usage.
type TeamUsageLog struct {
	UserID    string  `json:"user_id"`
	UserName  string  `json:"user_name,omitempty"`
	Email     string  `json:"email,omitempty"`
	TaskCount int     `json:"task_count"`
	Credits   float64 `json:"credits"`
}

// TeamUsageLogResponse contains paginated team usage records.
type TeamUsageLogResponse struct {
	OK         bool           `json:"ok"`
	RequestID  string         `json:"request_id"`
	Records    []TeamUsageLog `json:"data"`
	HasMore    bool           `json:"has_more"`
	NextCursor string         `json:"next_cursor,omitempty"`
}

// AvailableCredits describes available credits and their refresh schedule.
type AvailableCredits struct {
	TotalCredits      float64 `json:"total_credits"`
	FreeCredits       float64 `json:"free_credits"`
	PeriodicCredits   float64 `json:"periodic_credits"`
	AddonCredits      float64 `json:"addon_credits"`
	ProMonthlyCredits float64 `json:"pro_monthly_credits"`
	EventCredits      float64 `json:"event_credits"`
	RefreshCredits    float64 `json:"refresh_credits"`
	MaxRefreshCredits float64 `json:"max_refresh_credits"`
	NextRefreshTime   int64   `json:"next_refresh_time"`
	RefreshInterval   int64   `json:"refresh_interval"`
	CurrentPeriodEnd  int64   `json:"current_period_end"`
}

// AvailableCreditsResponse contains the caller's credit balance.
type AvailableCreditsResponse struct {
	OK        bool             `json:"ok"`
	RequestID string           `json:"request_id"`
	Credits   AvailableCredits `json:"data"`
}

// WebsiteStatusResponse contains a website's publication status.
type WebsiteStatusResponse struct {
	OK              bool     `json:"ok"`
	RequestID       string   `json:"request_id"`
	WebsiteID       string   `json:"website_id"`
	PublishStatus   string   `json:"publish_status"`
	SiteURLs        []string `json:"site_urls"`
	VersionID       string   `json:"version_id,omitempty"`
	StatusUpdatedAt int64    `json:"status_updated_at"`
	Visibility      string   `json:"visibility,omitempty"`
}

// WebsiteCheckpoint describes one saved website version.
type WebsiteCheckpoint struct {
	VersionID string `json:"version_id"`
	Message   string `json:"message,omitempty"`
	Status    string `json:"status,omitempty"`
	CreatedAt int64  `json:"created_at"`
}

// WebsiteCheckpointListResponse contains website version history.
type WebsiteCheckpointListResponse struct {
	OK                 bool                `json:"ok"`
	RequestID          string              `json:"request_id"`
	WebsiteID          string              `json:"website_id"`
	Checkpoints        []WebsiteCheckpoint `json:"data"`
	PublishedVersionID string              `json:"published_version_id,omitempty"`
}

// WebsitePublishResponse confirms a website publication request.
type WebsitePublishResponse struct {
	OK        bool   `json:"ok"`
	RequestID string `json:"request_id"`
	WebsiteID string `json:"website_id"`
	VersionID string `json:"version_id"`
}
