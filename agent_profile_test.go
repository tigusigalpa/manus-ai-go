package manusai

import "testing"

const changedProfile = "changed"

func TestAgentProfiles(t *testing.T) {
	profiles := AllAgentProfiles()
	if len(profiles) != 5 || !IsValidAgentProfile(AgentProfileManus16) || !IsValidAgentProfile(AgentProfileQuality) {
		t.Fatalf("unexpected profiles: %#v", profiles)
	}
	if IsValidAgentProfile("unknown") || IsValidAgentProfile("") {
		t.Fatal("unknown profiles must not be valid")
	}

	recommended := RecommendedAgentProfiles()
	if len(recommended) != 3 || IsDeprecatedAgentProfile(AgentProfileManus16) {
		t.Fatalf("unexpected recommended profiles: %#v", recommended)
	}
	if !IsDeprecatedAgentProfile(AgentProfileSpeed) || !IsDeprecatedAgentProfile(AgentProfileQuality) {
		t.Fatal("legacy profiles must be marked deprecated")
	}

	profiles[0] = changedProfile
	if AllAgentProfiles()[0] == changedProfile {
		t.Fatal("AllAgentProfiles() exposed its internal slice")
	}
}
