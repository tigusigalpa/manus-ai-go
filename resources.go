package manusai

import (
	"net/url"
	"strconv"
	"strings"
)

const (
	agentIDField        = "agent_id"
	emptyAgentIDMessage = "Agent ID cannot be empty"
	websiteIDField      = "website_id"
)

// CreateProject creates a project with optional instructions applied to its tasks.
func (c *Client) CreateProject(name, instruction string) (*ProjectResponse, error) {
	if strings.TrimSpace(name) == "" {
		return nil, &ValidationError{Message: "Project name cannot be empty"}
	}

	payload := map[string]string{"name": name}
	if instruction != "" {
		payload["instruction"] = instruction
	}
	var result ProjectResponse
	if err := c.request("POST", "/v2/project.create", payload, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListProjects returns projects available to the authenticated user.
func (c *Client) ListProjects() (*ProjectListResponse, error) {
	var result ProjectListResponse
	if err := c.request("GET", "/v2/project.list", nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListSkills returns global skills and, when projectID is set, project-specific skills.
func (c *Client) ListSkills(projectID string) (*SkillListResponse, error) {
	query := url.Values{}
	if projectID != "" {
		query.Set("project_id", projectID)
	}
	var result SkillListResponse
	if err := c.request("GET", "/v2/skill.list", nil, query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListAgents returns custom agents in the authenticated account.
func (c *Client) ListAgents() (*AgentListResponse, error) {
	var result AgentListResponse
	if err := c.request("GET", "/v2/agent.list", nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAgent returns details for agentID.
func (c *Client) GetAgent(agentID string) (*AgentResponse, error) {
	if strings.TrimSpace(agentID) == "" {
		return nil, &ValidationError{Message: emptyAgentIDMessage}
	}
	query := url.Values{agentIDField: {agentID}}
	var result AgentResponse
	if err := c.request("GET", "/v2/agent.detail", nil, query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateAgent updates an agent's nickname and about text. At least one field is required.
func (c *Client) UpdateAgent(agentID, nickname, about string) (*AgentResponse, error) {
	if strings.TrimSpace(agentID) == "" {
		return nil, &ValidationError{Message: emptyAgentIDMessage}
	}
	if nickname == "" && about == "" {
		return nil, &ValidationError{Message: "At least one agent update field is required"}
	}
	payload := map[string]string{agentIDField: agentID}
	if nickname != "" {
		payload["nickname"] = nickname
	}
	if about != "" {
		payload["about"] = about
	}
	var result AgentResponse
	if err := c.request("POST", "/v2/agent.update", payload, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListWebhooks returns registered webhooks.
func (c *Client) ListWebhooks() (*WebhookListResponse, error) {
	var result WebhookListResponse
	if err := c.request("GET", "/v2/webhook.list", nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListOnlineBrowserClients returns browser clients currently connected to Manus.
func (c *Client) ListOnlineBrowserClients() (*BrowserClientListResponse, error) {
	var result BrowserClientListResponse
	if err := c.request("GET", "/v2/browser.onlineList", nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetWebhookPublicKey returns the public key for verifying webhook signatures.
func (c *Client) GetWebhookPublicKey() (*WebhookPublicKeyResponse, error) {
	var result WebhookPublicKeyResponse
	if err := c.request("GET", "/v2/webhook.publicKey", nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListConnectors returns connectors installed in the authenticated account.
func (c *Client) ListConnectors() (*ConnectorListResponse, error) {
	var result ConnectorListResponse
	if err := c.request("GET", "/v2/connector.list", nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListUsage returns the caller's credit changes using optional cursor pagination.
func (c *Client) ListUsage(limit int, cursor string) (*UsageListResponse, error) {
	query := paginationQuery(limit, cursor)
	var result UsageListResponse
	if err := c.request("GET", "/v2/usage.list", nil, query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTeamUsageStatistic returns daily team credit usage for the optional date range.
func (c *Client) GetTeamUsageStatistic(startDate, endDate string) (*TeamUsageStatisticResponse, error) {
	query := dateRangeQuery(startDate, endDate)
	var result TeamUsageStatisticResponse
	if err := c.request("GET", "/v2/usage.teamStatistic", nil, query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListTeamUsageLog returns aggregated team usage with optional filters.
func (c *Client) ListTeamUsageLog(limit int, cursor, startDate, endDate, sortBy string, isAscending *bool) (*TeamUsageLogResponse, error) {
	query := paginationQuery(limit, cursor)
	for key, value := range dateRangeQuery(startDate, endDate) {
		query.Set(key, value[0])
	}
	if sortBy != "" {
		query.Set("sort_by", sortBy)
	}
	if isAscending != nil {
		if *isAscending {
			query.Set("is_asc", "true")
		} else {
			query.Set("is_asc", "false")
		}
	}
	var result TeamUsageLogResponse
	if err := c.request("GET", "/v2/usage.teamLog", nil, query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAvailableCredits returns the caller's available credit balance.
func (c *Client) GetAvailableCredits() (*AvailableCreditsResponse, error) {
	var result AvailableCreditsResponse
	if err := c.request("GET", "/v2/usage.availableCredits", nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetWebsiteStatus returns the publication status for a task website.
func (c *Client) GetWebsiteStatus(taskID, websiteID string) (*WebsiteStatusResponse, error) {
	query, err := websiteQuery(taskID, websiteID)
	if err != nil {
		return nil, err
	}
	var result WebsiteStatusResponse
	if err := c.request("GET", "/v2/website.status", nil, query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListWebsiteCheckpoints returns version history for a task website.
func (c *Client) ListWebsiteCheckpoints(taskID, websiteID string) (*WebsiteCheckpointListResponse, error) {
	query, err := websiteQuery(taskID, websiteID)
	if err != nil {
		return nil, err
	}
	var result WebsiteCheckpointListResponse
	if err := c.request("GET", "/v2/website.listCheckpoints", nil, query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PublishWebsite deploys a task website with the requested visibility.
func (c *Client) PublishWebsite(taskID, websiteID, visibility string) (*WebsitePublishResponse, error) {
	payload, err := websitePayload(taskID, websiteID)
	if err != nil {
		return nil, err
	}
	if visibility != "" {
		payload["visibility"] = visibility
	}
	var result WebsitePublishResponse
	if err := c.request("POST", "/v2/website.publish", payload, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateWebsite changes a task website's title or visibility. At least one field is required.
func (c *Client) UpdateWebsite(taskID, websiteID, title, visibility string) error {
	payload, err := websitePayload(taskID, websiteID)
	if err != nil {
		return err
	}
	if title == "" && visibility == "" {
		return &ValidationError{Message: "At least one website update field is required"}
	}
	if title != "" {
		payload["title"] = title
	}
	if visibility != "" {
		payload["visibility"] = visibility
	}
	return c.request("POST", "/v2/website.update", payload, nil, nil)
}

func paginationQuery(limit int, cursor string) url.Values {
	query := url.Values{}
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		query.Set("cursor", cursor)
	}
	return query
}

func dateRangeQuery(startDate, endDate string) url.Values {
	query := url.Values{}
	if startDate != "" {
		query.Set("start_date", startDate)
	}
	if endDate != "" {
		query.Set("end_date", endDate)
	}
	return query
}

func websiteQuery(taskID, websiteID string) (url.Values, error) {
	payload, err := websitePayload(taskID, websiteID)
	if err != nil {
		return nil, err
	}
	return url.Values{taskIDField: {payload[taskIDField]}, websiteIDField: {payload[websiteIDField]}}, nil
}

func websitePayload(taskID, websiteID string) (map[string]string, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, &ValidationError{Message: "Task ID cannot be empty"}
	}
	if strings.TrimSpace(websiteID) == "" {
		return nil, &ValidationError{Message: "Website ID cannot be empty"}
	}
	return map[string]string{taskIDField: taskID, websiteIDField: websiteID}, nil
}
