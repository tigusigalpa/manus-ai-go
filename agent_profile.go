package manusai

// Agent profile identifiers accepted by Manus when creating a task.
const (
	AgentProfileManus16     = "manus-1.6"
	AgentProfileManus16Lite = "manus-1.6-lite"
	AgentProfileManus16Max  = "manus-1.6-max"
	// AgentProfileSpeed is retained for source compatibility but is not accepted by API v2.
	AgentProfileSpeed = "speed"
	// AgentProfileQuality is retained for source compatibility but is not accepted by API v2.
	AgentProfileQuality = "quality"
)

var (
	allProfiles = []string{
		AgentProfileManus16,
		AgentProfileManus16Lite,
		AgentProfileManus16Max,
	}

	recommendedProfiles = []string{
		AgentProfileManus16,
		AgentProfileManus16Lite,
		AgentProfileManus16Max,
	}

	deprecatedProfiles = []string{
		AgentProfileSpeed,
		AgentProfileQuality,
	}
)

// AllAgentProfiles returns a copy of every profile known to this SDK.
func AllAgentProfiles() []string {
	result := make([]string, len(allProfiles))
	copy(result, allProfiles)
	return result
}

// RecommendedAgentProfiles returns a copy of the current recommended profiles.
func RecommendedAgentProfiles() []string {
	result := make([]string, len(recommendedProfiles))
	copy(result, recommendedProfiles)
	return result
}

// IsValidAgentProfile reports whether profile is known to this SDK.
func IsValidAgentProfile(profile string) bool {
	for _, p := range allProfiles {
		if p == profile {
			return true
		}
	}
	return false
}

// IsDeprecatedAgentProfile reports whether profile is retained only for compatibility and is not accepted by API v2.
func IsDeprecatedAgentProfile(profile string) bool {
	for _, p := range deprecatedProfiles {
		if p == profile {
			return true
		}
	}
	return false
}
