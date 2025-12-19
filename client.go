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

	payload := map[string]interface{}{
		"prompt":       prompt,
		"agentProfile": "manus-1.6",
	}

	if options != nil {
		if options.AgentProfile != "" {
			payload["agentProfile"] = options.AgentProfile
		}
		if options.TaskMode != "" {
			payload["taskMode"] = options.TaskMode
		}
		if options.Locale != "" {
			payload["locale"] = options.Locale
		}
		if options.HideInTaskList != nil {
			payload["hideInTaskList"] = *options.HideInTaskList
		}
		if options.CreateShareableLink != nil {
			payload["createShareableLink"] = *options.CreateShareableLink
		}
		if options.Attachments != nil && len(options.Attachments) > 0 {
			payload["attachments"] = options.Attachments
		}
	}

	var result TaskResponse
	err := c.request("POST", "/v1/tasks", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetTasks(filters *TaskFilters) (*TaskListResponse, error) {
	query := url.Values{}

	if filters != nil {
		if filters.After != "" {
			query.Set("after", filters.After)
		}
		if filters.Limit > 0 {
			query.Set("limit", fmt.Sprintf("%d", filters.Limit))
		}
		if filters.Order != "" {
			query.Set("order", filters.Order)
		}
		if filters.OrderBy != "" {
			query.Set("orderBy", filters.OrderBy)
		}
		if filters.Query != "" {
			query.Set("query", filters.Query)
		}
		if filters.Status != nil && len(filters.Status) > 0 {
			for _, status := range filters.Status {
				query.Add("status", status)
			}
		}
		if filters.CreatedAfter != "" {
			query.Set("createdAfter", filters.CreatedAfter)
		}
		if filters.CreatedBefore != "" {
			query.Set("createdBefore", filters.CreatedBefore)
		}
	}

	var result TaskListResponse
	err := c.request("GET", "/v1/tasks", nil, query, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetTask(taskID string) (*TaskDetail, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: "Task ID cannot be empty"}
	}

	var result TaskDetail
	err := c.request("GET", fmt.Sprintf("/v1/tasks/%s", taskID), nil, nil, &result)
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

	payload := make(map[string]interface{})
	hasUpdates := false

	if updates.Title != nil {
		payload["title"] = *updates.Title
		hasUpdates = true
	}
	if updates.EnableShared != nil {
		payload["enableShared"] = *updates.EnableShared
		hasUpdates = true
	}
	if updates.EnableVisibleInTaskList != nil {
		payload["enableVisibleInTaskList"] = *updates.EnableVisibleInTaskList
		hasUpdates = true
	}

	if !hasUpdates {
		return nil, &ValidationError{Message: "No valid update fields provided"}
	}

	var result TaskDetail
	err := c.request("PATCH", fmt.Sprintf("/v1/tasks/%s", taskID), payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeleteTask(taskID string) (*DeleteResponse, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: "Task ID cannot be empty"}
	}

	var result DeleteResponse
	err := c.request("DELETE", fmt.Sprintf("/v1/tasks/%s", taskID), nil, nil, &result)
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
	err := c.request("POST", "/v1/files", payload, nil, &result)
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

func (c *Client) ListFiles() (*FileListResponse, error) {
	var result FileListResponse
	err := c.request("GET", "/v1/files", nil, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetFile(fileID string) (*FileDetail, error) {
	if strings.TrimSpace(fileID) == "" {
		return nil, &ValidationError{Message: "File ID cannot be empty"}
	}

	var result FileDetail
	err := c.request("GET", fmt.Sprintf("/v1/files/%s", fileID), nil, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeleteFile(fileID string) (*DeleteResponse, error) {
	if strings.TrimSpace(fileID) == "" {
		return nil, &ValidationError{Message: "File ID cannot be empty"}
	}

	var result DeleteResponse
	err := c.request("DELETE", fmt.Sprintf("/v1/files/%s", fileID), nil, nil, &result)
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
		"webhook": webhook,
	}

	var result WebhookResponse
	err := c.request("POST", "/v1/webhooks", payload, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeleteWebhook(webhookID string) error {
	if strings.TrimSpace(webhookID) == "" {
		return &ValidationError{Message: "Webhook ID cannot be empty"}
	}

	err := c.request("DELETE", fmt.Sprintf("/v1/webhooks/%s", webhookID), nil, nil, nil)
	return err
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

	req.Header.Set("Authorization", c.apiKey)
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
