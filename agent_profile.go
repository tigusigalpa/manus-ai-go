package manusai

const (
	AgentProfileManus16     = "manus-1.6"
	AgentProfileManus16Lite = "manus-1.6-lite"
	AgentProfileManus16Max  = "manus-1.6-max"
	AgentProfileSpeed       = "speed"
	AgentProfileQuality     = "quality"
)

var (
	allProfiles = []string{
		AgentProfileManus16,
		AgentProfileManus16Lite,
		AgentProfileManus16Max,
		AgentProfileSpeed,
		AgentProfileQuality,
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

func AllAgentProfiles() []string {
	result := make([]string, len(allProfiles))
	copy(result, allProfiles)
	return result
}

func RecommendedAgentProfiles() []string {
	result := make([]string, len(recommendedProfiles))
	copy(result, recommendedProfiles)
	return result
}

func IsValidAgentProfile(profile string) bool {
	for _, p := range allProfiles {
		if p == profile {
			return true
		}
	}
	return false
}

func IsDeprecatedAgentProfile(profile string) bool {
	for _, p := range deprecatedProfiles {
		if p == profile {
			return true
		}
	}
	return false
}
