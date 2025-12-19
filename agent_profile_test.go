package manusai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllAgentProfiles(t *testing.T) {
	profiles := AllAgentProfiles()
	assert.Len(t, profiles, 5)
	assert.Contains(t, profiles, AgentProfileManus16)
	assert.Contains(t, profiles, AgentProfileManus16Lite)
	assert.Contains(t, profiles, AgentProfileManus16Max)
	assert.Contains(t, profiles, AgentProfileSpeed)
	assert.Contains(t, profiles, AgentProfileQuality)
}

func TestRecommendedAgentProfiles(t *testing.T) {
	profiles := RecommendedAgentProfiles()
	assert.Len(t, profiles, 3)
	assert.Contains(t, profiles, AgentProfileManus16)
	assert.Contains(t, profiles, AgentProfileManus16Lite)
	assert.Contains(t, profiles, AgentProfileManus16Max)
	assert.NotContains(t, profiles, AgentProfileSpeed)
	assert.NotContains(t, profiles, AgentProfileQuality)
}

func TestIsValidAgentProfile(t *testing.T) {
	tests := []struct {
		profile string
		valid   bool
	}{
		{AgentProfileManus16, true},
		{AgentProfileManus16Lite, true},
		{AgentProfileManus16Max, true},
		{AgentProfileSpeed, true},
		{AgentProfileQuality, true},
		{"invalid-profile", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.profile, func(t *testing.T) {
			assert.Equal(t, tt.valid, IsValidAgentProfile(tt.profile))
		})
	}
}

func TestIsDeprecatedAgentProfile(t *testing.T) {
	tests := []struct {
		profile    string
		deprecated bool
	}{
		{AgentProfileManus16, false},
		{AgentProfileManus16Lite, false},
		{AgentProfileManus16Max, false},
		{AgentProfileSpeed, true},
		{AgentProfileQuality, true},
		{"invalid-profile", false},
	}

	for _, tt := range tests {
		t.Run(tt.profile, func(t *testing.T) {
			assert.Equal(t, tt.deprecated, IsDeprecatedAgentProfile(tt.profile))
		})
	}
}
