package manusai

import "encoding/json"

// UnmarshalJSON decodes task timestamps returned as numbers, strings, or RFC3339 dates.
func (t *TaskSummary) UnmarshalJSON(data []byte) error {
	type taskSummaryWire struct {
		ID              string          `json:"id"`
		Title           string          `json:"title"`
		AgentStatus     string          `json:"status"`
		ShareVisibility string          `json:"share_visibility"`
		CreditUsage     float64         `json:"credit_usage"`
		TaskURL         string          `json:"task_url"`
		TaskType        string          `json:"task_type"`
		AgentProfile    string          `json:"agent_profile"`
		CreatedAt       json.RawMessage `json:"created_at"`
		UpdatedAt       json.RawMessage `json:"updated_at"`
	}
	var wire taskSummaryWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	createdAt, err := parseTimestamp(wire.CreatedAt, "task.created_at")
	if err != nil {
		return err
	}
	updatedAt, err := parseTimestamp(wire.UpdatedAt, "task.updated_at")
	if err != nil {
		return err
	}
	*t = TaskSummary{
		ID: wire.ID, Title: wire.Title, AgentStatus: wire.AgentStatus, ShareVisibility: wire.ShareVisibility,
		CreditUsage: wire.CreditUsage, TaskURL: wire.TaskURL, TaskType: wire.TaskType, AgentProfile: wire.AgentProfile,
		CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
	return nil
}

// UnmarshalJSON decodes project creation timestamps returned by Manus API v2.
func (p *Project) UnmarshalJSON(data []byte) error {
	var wire struct {
		ID          string          `json:"id"`
		Name        string          `json:"name"`
		Instruction string          `json:"instruction"`
		CreatedAt   json.RawMessage `json:"created_at"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	createdAt, err := parseTimestamp(wire.CreatedAt, "project.created_at")
	if err != nil {
		return err
	}
	*p = Project{ID: wire.ID, Name: wire.Name, Instruction: wire.Instruction, CreatedAt: createdAt}
	return nil
}

// UnmarshalJSON decodes skill timestamps returned by Manus API v2.
func (s *Skill) UnmarshalJSON(data []byte) error {
	var wire struct {
		ID          string                 `json:"id"`
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		OwnerType   string                 `json:"owner_type"`
		CreatorInfo map[string]interface{} `json:"creator_info"`
		CreatedAt   json.RawMessage        `json:"created_at"`
		UpdatedAt   json.RawMessage        `json:"updated_at"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	createdAt, err := parseTimestamp(wire.CreatedAt, "skill.created_at")
	if err != nil {
		return err
	}
	updatedAt, err := parseTimestamp(wire.UpdatedAt, "skill.updated_at")
	if err != nil {
		return err
	}
	*s = Skill{ID: wire.ID, Name: wire.Name, Description: wire.Description, OwnerType: wire.OwnerType, CreatorInfo: wire.CreatorInfo, CreatedAt: createdAt, UpdatedAt: updatedAt}
	return nil
}

// UnmarshalJSON decodes webhook creation timestamps returned by Manus API v2.
func (w *Webhook) UnmarshalJSON(data []byte) error {
	var wire struct {
		ID        string          `json:"id"`
		URL       string          `json:"url"`
		Status    string          `json:"status"`
		CreatedAt json.RawMessage `json:"created_at"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	createdAt, err := parseTimestamp(wire.CreatedAt, "webhook.created_at")
	if err != nil {
		return err
	}
	*w = Webhook{ID: wire.ID, URL: wire.URL, Status: wire.Status, CreatedAt: createdAt}
	return nil
}

// UnmarshalJSON decodes usage timestamps returned by Manus API v2.
func (u *UsageRecord) UnmarshalJSON(data []byte) error {
	var wire struct {
		TaskID           string          `json:"task_id"`
		Title            string          `json:"title"`
		Credits          float64         `json:"credits"`
		CreatedAt        json.RawMessage `json:"created_at"`
		Type             string          `json:"type"`
		CollaborateInfos []interface{}   `json:"collaborate_infos"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	createdAt, err := parseTimestamp(wire.CreatedAt, "usage.created_at")
	if err != nil {
		return err
	}
	*u = UsageRecord{TaskID: wire.TaskID, Title: wire.Title, Credits: wire.Credits, CreatedAt: createdAt, Type: wire.Type, CollaborateInfos: wire.CollaborateInfos}
	return nil
}

// UnmarshalJSON decodes credit-refresh timestamps returned by Manus API v2.
func (c *AvailableCredits) UnmarshalJSON(data []byte) error {
	type creditsWire struct {
		TotalCredits      float64         `json:"total_credits"`
		FreeCredits       float64         `json:"free_credits"`
		PeriodicCredits   float64         `json:"periodic_credits"`
		AddonCredits      float64         `json:"addon_credits"`
		ProMonthlyCredits float64         `json:"pro_monthly_credits"`
		EventCredits      float64         `json:"event_credits"`
		RefreshCredits    float64         `json:"refresh_credits"`
		MaxRefreshCredits float64         `json:"max_refresh_credits"`
		NextRefreshTime   json.RawMessage `json:"next_refresh_time"`
		RefreshInterval   json.RawMessage `json:"refresh_interval"`
		CurrentPeriodEnd  json.RawMessage `json:"current_period_end"`
	}
	var wire creditsWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	nextRefreshTime, err := parseTimestamp(wire.NextRefreshTime, "available_credits.next_refresh_time")
	if err != nil {
		return err
	}
	refreshInterval, err := parseInteger(wire.RefreshInterval, "available_credits.refresh_interval")
	if err != nil {
		return err
	}
	currentPeriodEnd, err := parseTimestamp(wire.CurrentPeriodEnd, "available_credits.current_period_end")
	if err != nil {
		return err
	}
	*c = AvailableCredits{TotalCredits: wire.TotalCredits, FreeCredits: wire.FreeCredits, PeriodicCredits: wire.PeriodicCredits, AddonCredits: wire.AddonCredits, ProMonthlyCredits: wire.ProMonthlyCredits, EventCredits: wire.EventCredits, RefreshCredits: wire.RefreshCredits, MaxRefreshCredits: wire.MaxRefreshCredits, NextRefreshTime: nextRefreshTime, RefreshInterval: refreshInterval, CurrentPeriodEnd: currentPeriodEnd}
	return nil
}

// UnmarshalJSON decodes website status timestamps returned by Manus API v2.
func (w *WebsiteStatusResponse) UnmarshalJSON(data []byte) error {
	var wire struct {
		OK              bool            `json:"ok"`
		RequestID       string          `json:"request_id"`
		WebsiteID       string          `json:"website_id"`
		PublishStatus   string          `json:"publish_status"`
		SiteURLs        []string        `json:"site_urls"`
		VersionID       string          `json:"version_id"`
		StatusUpdatedAt json.RawMessage `json:"status_updated_at"`
		Visibility      string          `json:"visibility"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	statusUpdatedAt, err := parseTimestamp(wire.StatusUpdatedAt, "website.status_updated_at")
	if err != nil {
		return err
	}
	*w = WebsiteStatusResponse{OK: wire.OK, RequestID: wire.RequestID, WebsiteID: wire.WebsiteID, PublishStatus: wire.PublishStatus, SiteURLs: wire.SiteURLs, VersionID: wire.VersionID, StatusUpdatedAt: statusUpdatedAt, Visibility: wire.Visibility}
	return nil
}

// UnmarshalJSON decodes website checkpoint creation timestamps returned by Manus API v2.
func (w *WebsiteCheckpoint) UnmarshalJSON(data []byte) error {
	var wire struct {
		VersionID string          `json:"version_id"`
		Message   string          `json:"message"`
		Status    string          `json:"status"`
		CreatedAt json.RawMessage `json:"created_at"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	createdAt, err := parseTimestamp(wire.CreatedAt, "website_checkpoint.created_at")
	if err != nil {
		return err
	}
	*w = WebsiteCheckpoint{VersionID: wire.VersionID, Message: wire.Message, Status: wire.Status, CreatedAt: createdAt}
	return nil
}
