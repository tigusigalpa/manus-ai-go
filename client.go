package manusai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Default client settings.
const (
	DefaultBaseURL        = "https://api.manus.ai"
	DefaultTimeout        = 30 * time.Second
	DefaultConnectTimeout = 10 * time.Second
	defaultContentType    = "application/octet-stream"
	maxResponseBodySize   = 4 << 20 // 4 MiB
	attachmentTypeKey     = "type"
	attachmentTypeFile    = "file"
	fileIDField           = "file_id"
	messageContentKey     = "content"
	messageKey            = "message"
	messageTextKey        = "text"
	messageTypeText       = "text"
	taskIDField           = "task_id"
	emptyTaskIDMessage    = "Task ID cannot be empty"
	emptyFileIDMessage    = "File ID cannot be empty"
	attachmentsError      = "Attachments must be created with an attachment helper or be map[string]interface{}"
)

// Client is a Manus API v2 client.
type Client struct {
	apiKey      string
	bearerToken string
	baseURL     string
	httpClient  *http.Client
}

type taskDetailResponse struct {
	OK        bool        `json:"ok"`
	RequestID string      `json:"request_id"`
	Task      TaskSummary `json:"task"`
}

// ClientOption configures a Client created by NewClient.
type ClientOption func(*Client)

// WithBaseURL overrides the Manus API base URL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithHTTPClient uses httpClient to make API and file-upload requests.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBearerToken authenticates requests with an OAuth access token instead of an API key.
func WithBearerToken(token string) ClientOption {
	return func(c *Client) {
		c.bearerToken = strings.TrimSpace(token)
	}
}

// WithTimeout sets the overall timeout for requests made by the SDK client.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates a Client using apiKey and optional configuration.
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	client := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		if opt == nil {
			return nil, &ValidationError{Message: "Client option cannot be nil"}
		}
		opt(client)
	}

	if client.httpClient == nil {
		return nil, &ValidationError{Message: "HTTP client cannot be nil"}
	}
	if strings.TrimSpace(client.apiKey) == "" && client.bearerToken == "" {
		return nil, &AuthenticationError{Message: "API key or bearer token cannot be empty"}
	}
	if strings.TrimSpace(client.apiKey) != "" && client.bearerToken != "" {
		return nil, &ValidationError{Message: "API key and bearer token cannot be used together"}
	}

	baseURL, err := url.ParseRequestURI(client.baseURL)
	if err != nil || baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, &ValidationError{Message: "Base URL must be an absolute HTTP URL"}
	}
	if baseURL.Scheme != "http" && baseURL.Scheme != "https" {
		return nil, &ValidationError{Message: "Base URL must use HTTP or HTTPS"}
	}

	return client, nil
}

// CreateTask creates a Manus task from prompt and optional task configuration.
func (c *Client) CreateTask(prompt string, options *TaskOptions) (*TaskResponse, error) {
	if strings.TrimSpace(prompt) == "" {
		return nil, &ValidationError{Message: "Task prompt cannot be empty"}
	}

	payload, err := newTaskPayload(prompt, options)
	if err != nil {
		return nil, err
	}

	var result TaskResponse
	if err := c.request("POST", "/v2/task.create", payload, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func newTaskPayload(prompt string, options *TaskOptions) (map[string]interface{}, error) {
	message := map[string]interface{}{
		messageContentKey: []map[string]interface{}{
			{
				attachmentTypeKey: messageTypeText,
				messageTextKey:    prompt,
			},
		},
	}

	payload := map[string]interface{}{
		messageKey: message,
	}

	if options == nil {
		return payload, nil
	}

	if err := applyTaskOptions(payload, message, options); err != nil {
		return nil, err
	}

	return payload, nil
}

func applyTaskOptions(payload, message map[string]interface{}, options *TaskOptions) error {
	if options.AgentProfile != "" {
		payload["agent_profile"] = options.AgentProfile
	}
	if options.Locale != "" {
		payload["locale"] = options.Locale
	}
	if options.HideInTaskList != nil {
		payload["hide_in_task_list"] = *options.HideInTaskList
	}
	if options.ShareVisibility != "" {
		payload["share_visibility"] = options.ShareVisibility
	}
	if options.Title != "" {
		payload["title"] = options.Title
	}
	if options.ProjectID != "" {
		payload["project_id"] = options.ProjectID
	}
	interactiveMode := options.InteractiveMode
	if interactiveMode == nil {
		interactiveMode = options.EnableAskUser
	}
	if interactiveMode != nil {
		payload["interactive_mode"] = *interactiveMode
	}
	if options.StructuredOutputSchema != nil {
		payload["structured_output_schema"] = options.StructuredOutputSchema
	}
	if len(options.Connectors) > 0 {
		message["connectors"] = options.Connectors
	}
	if len(options.EnableSkills) > 0 {
		message["enable_skills"] = options.EnableSkills
	}
	if len(options.ForceSkills) > 0 {
		message["force_skills"] = options.ForceSkills
	}

	attachments, err := appendAttachments(message[messageContentKey].([]map[string]interface{}), options.Attachments)
	if err != nil {
		return err
	}
	message[messageContentKey] = attachments

	return nil
}

func appendAttachments(content []map[string]interface{}, attachments []interface{}) ([]map[string]interface{}, error) {
	for _, attachment := range attachments {
		attachmentMap, ok := attachment.(map[string]interface{})
		if !ok {
			return nil, &ValidationError{Message: attachmentsError}
		}
		content = append(content, attachmentMap)
	}

	return content, nil
}

// GetTasks returns tasks that match filters. A nil filter returns the default list.
func (c *Client) GetTasks(filters *TaskFilters) (*TaskListResponse, error) {
	query := url.Values{}

	if filters != nil {
		if filters.Cursor != "" {
			query.Set("cursor", filters.Cursor)
		}
		if filters.Limit > 0 {
			query.Set("limit", fmt.Sprintf("%d", filters.Limit))
		}
		if filters.Order != "" {
			query.Set("order", filters.Order)
		}
		if filters.Scope != "" {
			query.Set("scope", filters.Scope)
		}
		if filters.AgentID != "" {
			query.Set("agent_id", filters.AgentID)
		}
		if filters.ProjectID != "" {
			query.Set("project_id", filters.ProjectID)
		}
		if filters.OAuthClientID != "" {
			query.Set("oauth_client_id", filters.OAuthClientID)
		}
		if filters.APIKeyID != "" {
			query.Set("api_key_id", filters.APIKeyID)
		}
	}

	var result TaskListResponse
	err := c.request("GET", "/v2/task.list", nil, query, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTask returns task details for taskID.
func (c *Client) GetTask(taskID string) (*TaskDetail, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: emptyTaskIDMessage}
	}

	query := url.Values{}
	query.Set(taskIDField, taskID)

	var response taskDetailResponse
	err := c.request("GET", "/v2/task.detail", nil, query, &response)
	if err != nil {
		return nil, err
	}
	return &TaskDetail{
		OK: response.OK, RequestID: response.RequestID, ID: response.Task.ID, Title: response.Task.Title,
		AgentStatus: response.Task.AgentStatus, ShareVisibility: response.Task.ShareVisibility,
		CreditUsage: response.Task.CreditUsage, TaskURL: response.Task.TaskURL, TaskType: response.Task.TaskType,
		AgentProfile: response.Task.AgentProfile, CreatedAt: response.Task.CreatedAt, UpdatedAt: response.Task.UpdatedAt,
	}, nil
}

// UpdateTask updates the supported fields of taskID.
func (c *Client) UpdateTask(taskID string, updates *TaskUpdate) (*TaskDetail, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: emptyTaskIDMessage}
	}

	if updates == nil {
		return nil, &ValidationError{Message: "Updates cannot be nil"}
	}

	payload := map[string]interface{}{
		taskIDField: taskID,
	}
	hasUpdates := false

	if updates.Title != nil {
		payload["title"] = *updates.Title
		hasUpdates = true
	}
	if updates.ShareVisibility != nil {
		payload["share_visibility"] = *updates.ShareVisibility
		hasUpdates = true
	}
	visibleInTaskList := updates.EnableVisibleInTaskList
	if visibleInTaskList == nil && updates.HideInTaskList != nil {
		visible := !*updates.HideInTaskList
		visibleInTaskList = &visible
	}
	if visibleInTaskList != nil {
		payload["enable_visible_in_task_list"] = *visibleInTaskList
		hasUpdates = true
	}

	if !hasUpdates {
		return nil, &ValidationError{Message: "No valid update fields provided"}
	}

	var response struct {
		OK              bool   `json:"ok"`
		RequestID       string `json:"request_id"`
		TaskID          string `json:"task_id"`
		TaskTitle       string `json:"task_title"`
		TaskURL         string `json:"task_url"`
		ShareVisibility string `json:"share_visibility"`
	}
	err := c.request("POST", "/v2/task.update", payload, nil, &response)
	if err != nil {
		return nil, err
	}
	return &TaskDetail{OK: response.OK, RequestID: response.RequestID, ID: response.TaskID, Title: response.TaskTitle, TaskURL: response.TaskURL, ShareVisibility: response.ShareVisibility}, nil
}

// DeleteTask deletes taskID and returns the API deletion result.
func (c *Client) DeleteTask(taskID string) (*DeleteResponse, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: emptyTaskIDMessage}
	}

	payload := map[string]interface{}{
		taskIDField: taskID,
	}

	var result DeleteResponse
	err := c.request("POST", "/v2/task.delete", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateFile creates a Manus file record and returns its upload URL.
func (c *Client) CreateFile(filename string) (*FileResponse, error) {
	if strings.TrimSpace(filename) == "" {
		return nil, &ValidationError{Message: "Filename cannot be empty"}
	}

	payload := map[string]string{
		"filename": filename,
	}

	var response struct {
		OK        bool   `json:"ok"`
		RequestID string `json:"request_id"`
		File      struct {
			ID        string          `json:"id"`
			Filename  string          `json:"filename"`
			Status    string          `json:"status"`
			CreatedAt json.RawMessage `json:"created_at"`
		} `json:"file"`
		UploadURL       string          `json:"upload_url"`
		UploadExpiresAt json.RawMessage `json:"upload_expires_at"`
	}
	err := c.request("POST", "/v2/file.upload", payload, nil, &response)
	if err != nil {
		return nil, err
	}
	uploadExpiresAt, err := parseTimestamp(response.UploadExpiresAt, "file.upload_expires_at")
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Failed to decode response: %v", err)}
	}
	return &FileResponse{OK: response.OK, RequestID: response.RequestID, FileID: response.File.ID, Filename: response.File.Filename, UploadURL: response.UploadURL, ExpiresAt: uploadExpiresAt}, nil
}

// UploadFileContent uploads fileContent to a URL returned by CreateFile.
func (c *Client) UploadFileContent(uploadURL string, fileContent []byte, contentType string) error {
	if strings.TrimSpace(uploadURL) == "" {
		return &ValidationError{Message: "Upload URL cannot be empty"}
	}

	if contentType == "" {
		contentType = defaultContentType
	}

	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(fileContent))
	if err != nil {
		return &Error{Message: fmt.Sprintf("Failed to create upload request: %v", err)}
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &Error{Message: fmt.Sprintf("Failed to upload file content: %v", err)}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := readResponseBody(resp.Body)
		return &Error{
			Message:    fmt.Sprintf("Upload failed with status %d: %s", resp.StatusCode, string(body)),
			StatusCode: resp.StatusCode,
		}
	}

	return nil
}

// ListFiles is deprecated because Manus API v2 has no file-list operation.
func (c *Client) ListFiles(limit int, cursor string) (*FileListResponse, error) {
	return nil, &ValidationError{Message: "Manus API v2 does not provide a file-list endpoint"}
}

// GetFile returns metadata for fileID.
func (c *Client) GetFile(fileID string) (*FileDetail, error) {
	if strings.TrimSpace(fileID) == "" {
		return nil, &ValidationError{Message: emptyFileIDMessage}
	}

	query := url.Values{}
	query.Set(fileIDField, fileID)

	var response struct {
		OK        bool   `json:"ok"`
		RequestID string `json:"request_id"`
		File      struct {
			ID        string          `json:"id"`
			Filename  string          `json:"filename"`
			Status    string          `json:"status"`
			Bytes     int64           `json:"bytes"`
			CreatedAt json.RawMessage `json:"created_at"`
			ExpiresAt json.RawMessage `json:"expires_at"`
		} `json:"file"`
	}
	err := c.request("GET", "/v2/file.detail", nil, query, &response)
	if err != nil {
		return nil, err
	}
	createdAt, err := parseTimestamp(response.File.CreatedAt, "file.created_at")
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Failed to decode response: %v", err)}
	}
	expiresAt, err := parseTimestamp(response.File.ExpiresAt, "file.expires_at")
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Failed to decode response: %v", err)}
	}
	return &FileDetail{OK: response.OK, RequestID: response.RequestID, FileID: response.File.ID, Filename: response.File.Filename, Status: response.File.Status, SizeBytes: response.File.Bytes, CreatedAt: createdAt, ExpiresAt: expiresAt}, nil
}

// DeleteFile deletes fileID and returns the API deletion result.
func (c *Client) DeleteFile(fileID string) (*DeleteResponse, error) {
	if strings.TrimSpace(fileID) == "" {
		return nil, &ValidationError{Message: emptyFileIDMessage}
	}

	payload := map[string]interface{}{
		fileIDField: fileID,
	}

	var result DeleteResponse
	err := c.request("POST", "/v2/file.delete", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateWebhook registers a webhook using webhook configuration.
func (c *Client) CreateWebhook(webhook *WebhookConfig) (*WebhookResponse, error) {
	if webhook == nil {
		return nil, &ValidationError{Message: "Webhook configuration cannot be nil"}
	}

	if strings.TrimSpace(webhook.URL) == "" {
		return nil, &ValidationError{Message: "Webhook URL is required"}
	}

	payload := map[string]interface{}{
		"url": webhook.URL,
	}

	var response struct {
		OK        bool   `json:"ok"`
		RequestID string `json:"request_id"`
		Webhook   struct {
			ID string `json:"id"`
		} `json:"webhook"`
	}
	err := c.request("POST", "/v2/webhook.create", payload, nil, &response)
	if err != nil {
		return nil, err
	}
	return &WebhookResponse{OK: response.OK, RequestID: response.RequestID, WebhookID: response.Webhook.ID}, nil
}

// DeleteWebhook deletes the webhook identified by webhookID.
func (c *Client) DeleteWebhook(webhookID string) error {
	if strings.TrimSpace(webhookID) == "" {
		return &ValidationError{Message: "Webhook ID cannot be empty"}
	}

	payload := map[string]interface{}{
		"webhook_id": webhookID,
	}

	err := c.request("POST", "/v2/webhook.delete", payload, nil, nil)
	return err
}

// ListMessages returns messages for taskID using optional cursor pagination.
func (c *Client) ListMessages(taskID string, limit int, cursor string, order string, verbose bool) (*TaskMessagesResponse, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: emptyTaskIDMessage}
	}

	query := url.Values{}
	query.Set(taskIDField, taskID)

	if limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", limit))
	}
	if cursor != "" {
		query.Set("cursor", cursor)
	}
	if order != "" {
		query.Set("order", order)
	}
	if verbose {
		query.Set("verbose", "true")
	}

	var result TaskMessagesResponse
	err := c.request("GET", "/v2/task.listMessages", nil, query, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// SendMessage sends message and optional attachments to an existing task.
func (c *Client) SendMessage(taskID string, message string, attachments []interface{}) (*SendMessageResponse, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: emptyTaskIDMessage}
	}
	if strings.TrimSpace(message) == "" {
		return nil, &ValidationError{Message: "Message cannot be empty"}
	}

	content := []map[string]interface{}{
		{
			attachmentTypeKey: messageTypeText,
			messageTextKey:    message,
		},
	}

	if len(attachments) > 0 {
		for _, att := range attachments {
			attMap, ok := att.(map[string]interface{})
			if !ok {
				return nil, &ValidationError{Message: attachmentsError}
			}
			content = append(content, attMap)
		}
	}

	payload := map[string]interface{}{
		taskIDField: taskID,
		messageKey: map[string]interface{}{
			messageContentKey: content,
		},
	}

	var result SendMessageResponse
	err := c.request("POST", "/v2/task.sendMessage", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// StopTask requests that Manus stop the running task identified by taskID.
func (c *Client) StopTask(taskID string) (*StopTaskResponse, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: emptyTaskIDMessage}
	}

	payload := map[string]interface{}{
		taskIDField: taskID,
	}

	var result StopTaskResponse
	err := c.request("POST", "/v2/task.stop", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ConfirmAction submits input for a pending task action event.
func (c *Client) ConfirmAction(taskID string, eventID string, input map[string]interface{}) (*ConfirmActionResponse, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: emptyTaskIDMessage}
	}
	if strings.TrimSpace(eventID) == "" {
		return nil, &ValidationError{Message: "Event ID cannot be empty"}
	}

	payload := map[string]interface{}{
		taskIDField: taskID,
		"event_id":  eventID,
		"input":     input,
	}

	var result ConfirmActionResponse
	err := c.request("POST", "/v2/task.confirmAction", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) request(method, endpoint string, body interface{}, query url.Values, result interface{}) error {
	fullURL := c.baseURL + endpoint
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Failed to marshal request body: %v", err)}
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return &Error{Message: fmt.Sprintf("Failed to create request: %v", err)}
	}

	if c.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	} else {
		req.Header.Set("x-manus-api-key", c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &Error{Message: fmt.Sprintf("Request failed: %v", err)}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	respBody, err := readResponseBody(resp.Body)
	if err != nil {
		return &Error{Message: fmt.Sprintf("Failed to read response body: %v", err)}
	}

	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(resp.StatusCode, respBody)
	}

	if len(respBody) == 0 {
		return nil
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return &Error{Message: fmt.Sprintf("Failed to decode response: %v", err)}
		}
	}

	return nil
}

func readResponseBody(body io.Reader) ([]byte, error) {
	response, err := io.ReadAll(io.LimitReader(body, maxResponseBodySize+1))
	if err != nil {
		return nil, err
	}
	if len(response) > maxResponseBodySize {
		return nil, fmt.Errorf("response body exceeds %d bytes", maxResponseBodySize)
	}

	return response, nil
}

func (c *Client) handleErrorResponse(statusCode int, body []byte) error {
	message := string(body)

	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return &AuthenticationError{
			Message:    fmt.Sprintf("Authentication failed: %s", message),
			StatusCode: statusCode,
		}
	case http.StatusBadRequest:
		return &ValidationError{
			Message:    fmt.Sprintf("Validation error: %s", message),
			StatusCode: statusCode,
		}
	default:
		return &Error{
			Message:    fmt.Sprintf("API request failed: %s", message),
			StatusCode: statusCode,
		}
	}
}
