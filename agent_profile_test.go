package manusai

import "testing"

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

	profiles[0] = "changed"
	if AllAgentProfiles()[0] == "changed" {
		t.Fatal("AllAgentProfiles() exposed its internal slice")
	}
}
