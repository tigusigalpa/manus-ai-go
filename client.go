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

const (
	DefaultBaseURL       = "https://api.manus.ai"
	DefaultTimeout       = 30 * time.Second
	DefaultConnectTimeout = 10 * time.Second
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type ClientOption func(*Client)

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, &AuthenticationError{Message: "API key cannot be empty"}
	}

	client := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

func (c *Client) CreateTask(prompt string, options *TaskOptions) (*TaskResponse, error) {
	if strings.TrimSpace(prompt) == "" {
		return nil, &ValidationError{Message: "Task prompt cannot be empty"}
	}

	message := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": prompt,
			},
		},
	}

	payload := map[string]interface{}{
		"message": message,
	}

	if options != nil {
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
		if options.EnableAskUser != nil {
			payload["enable_ask_user"] = *options.EnableAskUser
		}
		if options.Connectors != nil && len(options.Connectors) > 0 {
			message["connectors"] = options.Connectors
		}
		if options.EnableSkills != nil && len(options.EnableSkills) > 0 {
			message["enable_skills"] = options.EnableSkills
		}
		if options.ForceSkills != nil && len(options.ForceSkills) > 0 {
			message["force_skills"] = options.ForceSkills
		}
		if options.Attachments != nil && len(options.Attachments) > 0 {
			for _, att := range options.Attachments {
				if attMap, ok := att.(map[string]interface{}); ok {
					content := message["content"].([]map[string]interface{})
					content = append(content, attMap)
					message["content"] = content
				}
			}
		}
	}

	var result TaskResponse
	err := c.request("POST", "/v2/task.create", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

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
	}

	var result TaskListResponse
	err := c.request("GET", "/v2/task.list", nil, query, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetTask(taskID string) (*TaskDetail, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: "Task ID cannot be empty"}
	}

	query := url.Values{}
	query.Set("task_id", taskID)

	var result TaskDetail
	err := c.request("GET", "/v2/task.detail", nil, query, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) UpdateTask(taskID string, updates *TaskUpdate) (*TaskDetail, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: "Task ID cannot be empty"}
	}

	if updates == nil {
		return nil, &ValidationError{Message: "Updates cannot be nil"}
	}

	payload := map[string]interface{}{
		"task_id": taskID,
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
	if updates.HideInTaskList != nil {
		payload["hide_in_task_list"] = *updates.HideInTaskList
		hasUpdates = true
	}

	if !hasUpdates {
		return nil, &ValidationError{Message: "No valid update fields provided"}
	}

	var result TaskDetail
	err := c.request("POST", "/v2/task.update", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeleteTask(taskID string) (*DeleteResponse, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: "Task ID cannot be empty"}
	}

	payload := map[string]interface{}{
		"task_id": taskID,
	}

	var result DeleteResponse
	err := c.request("POST", "/v2/task.delete", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) CreateFile(filename string) (*FileResponse, error) {
	if strings.TrimSpace(filename) == "" {
		return nil, &ValidationError{Message: "Filename cannot be empty"}
	}

	payload := map[string]string{
		"filename": filename,
	}

	var result FileResponse
	err := c.request("POST", "/v2/file.upload", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) UploadFileContent(uploadURL string, fileContent []byte, contentType string) error {
	if strings.TrimSpace(uploadURL) == "" {
		return &ValidationError{Message: "Upload URL cannot be empty"}
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(fileContent))
	if err != nil {
		return &ManusAIError{Message: fmt.Sprintf("Failed to create upload request: %v", err)}
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &ManusAIError{Message: fmt.Sprintf("Failed to upload file content: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return &ManusAIError{
			Message:    fmt.Sprintf("Upload failed with status %d: %s", resp.StatusCode, string(body)),
			StatusCode: resp.StatusCode,
		}
	}

	return nil
}

func (c *Client) ListFiles(limit int, cursor string) (*FileListResponse, error) {
	query := url.Values{}
	if limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", limit))
	}
	if cursor != "" {
		query.Set("cursor", cursor)
	}

	var result FileListResponse
	err := c.request("GET", "/v2/file.list", nil, query, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetFile(fileID string) (*FileDetail, error) {
	if strings.TrimSpace(fileID) == "" {
		return nil, &ValidationError{Message: "File ID cannot be empty"}
	}

	query := url.Values{}
	query.Set("file_id", fileID)

	var result FileDetail
	err := c.request("GET", "/v2/file.detail", nil, query, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeleteFile(fileID string) (*DeleteResponse, error) {
	if strings.TrimSpace(fileID) == "" {
		return nil, &ValidationError{Message: "File ID cannot be empty"}
	}

	payload := map[string]interface{}{
		"file_id": fileID,
	}

	var result DeleteResponse
	err := c.request("POST", "/v2/file.delete", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) CreateWebhook(webhook *WebhookConfig) (*WebhookResponse, error) {
	if webhook == nil {
		return nil, &ValidationError{Message: "Webhook configuration cannot be nil"}
	}

	if webhook.URL == "" {
		return nil, &ValidationError{Message: "Webhook URL is required"}
	}

	payload := map[string]interface{}{
		"url":    webhook.URL,
		"events": webhook.Events,
	}

	var result WebhookResponse
	err := c.request("POST", "/v2/webhook.create", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

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

func (c *Client) ListMessages(taskID string, limit int, cursor string, order string, verbose bool) (*TaskMessagesResponse, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: "Task ID cannot be empty"}
	}

	query := url.Values{}
	query.Set("task_id", taskID)
	
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

func (c *Client) SendMessage(taskID string, message string, attachments []interface{}) (*SendMessageResponse, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: "Task ID cannot be empty"}
	}
	if strings.TrimSpace(message) == "" {
		return nil, &ValidationError{Message: "Message cannot be empty"}
	}

	content := []map[string]interface{}{
		{
			"type": "text",
			"text": message,
		},
	}

	if attachments != nil && len(attachments) > 0 {
		for _, att := range attachments {
			if attMap, ok := att.(map[string]interface{}); ok {
				content = append(content, attMap)
			}
		}
	}

	payload := map[string]interface{}{
		"task_id": taskID,
		"message": map[string]interface{}{
			"content": content,
		},
	}

	var result SendMessageResponse
	err := c.request("POST", "/v2/task.sendMessage", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) StopTask(taskID string) (*StopTaskResponse, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: "Task ID cannot be empty"}
	}

	payload := map[string]interface{}{
		"task_id": taskID,
	}

	var result StopTaskResponse
	err := c.request("POST", "/v2/task.stop", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) ConfirmAction(taskID string, eventID string, input map[string]interface{}) (*ConfirmActionResponse, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: "Task ID cannot be empty"}
	}
	if strings.TrimSpace(eventID) == "" {
		return nil, &ValidationError{Message: "Event ID cannot be empty"}
	}

	payload := map[string]interface{}{
		"task_id":  taskID,
		"event_id": eventID,
		"input":    input,
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
	if query != nil && len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return &ManusAIError{Message: fmt.Sprintf("Failed to marshal request body: %v", err)}
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return &ManusAIError{Message: fmt.Sprintf("Failed to create request: %v", err)}
	}

	req.Header.Set("x-manus-api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &ManusAIError{Message: fmt.Sprintf("Request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ManusAIError{Message: fmt.Sprintf("Failed to read response body: %v", err)}
	}

	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(resp.StatusCode, respBody)
	}

	if len(respBody) == 0 {
		return nil
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return &ManusAIError{Message: fmt.Sprintf("Failed to decode response: %v", err)}
		}
	}

	return nil
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
		return &ManusAIError{
			Message:    fmt.Sprintf("API request failed: %s", message),
			StatusCode: statusCode,
		}
	}
}
