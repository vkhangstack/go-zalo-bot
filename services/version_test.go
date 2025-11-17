package services

import (
	"testing"
)

func TestVersion(t *testing.T) {
	version := Version()
	if version == "" {
		t.Error("Version() returned empty string")
	}

	// Should match the version in version.yml
	expectedVersion := "0.0.1"
	if version != expectedVersion {
		t.Errorf("Expected version %s, got %s", expectedVersion, version)
	}
}

func TestUserAgent(t *testing.T) {
	userAgent := UserAgent()
	if userAgent == "" {
		t.Error("UserAgent() returned empty string")
	}

	// Should match the user_agent in version.yml
	expectedUserAgent := "Go-Zalo-Bot-SDK/0.0.1"
	if userAgent != expectedUserAgent {
		t.Errorf("Expected user agent %s, got %s", expectedUserAgent, userAgent)
	}
}

func TestSDKName(t *testing.T) {
	name := SDKName()
	if name == "" {
		t.Error("SDKName() returned empty string")
	}

	expectedName := "Go-Zalo-Bot-SDK"
	if name != expectedName {
		t.Errorf("Expected SDK name %s, got %s", expectedName, name)
	}
}

func TestVersionComponents(t *testing.T) {
	tests := []struct {
		name     string
		function func() int
		expected int
	}{
		{"Major", VersionMajor, 0},
		{"Minor", VersionMinor, 0},
		{"Patch", VersionPatch, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.function()
			if got != tt.expected {
				t.Errorf("Version%s() = %d, want %d", tt.name, got, tt.expected)
			}
		})
	}
}

func TestVersionPreRelease(t *testing.T) {
	prerelease := VersionPreRelease()
	// Should be empty for stable release
	if prerelease != "" {
		t.Errorf("Expected empty prerelease, got %s", prerelease)
	}
}

func TestReleaseDate(t *testing.T) {
	date := ReleaseDate()
	if date == "" {
		t.Error("ReleaseDate() returned empty string")
	}

	expectedDate := "2025-11-17"
	if date != expectedDate {
		t.Errorf("Expected release date %s, got %s", expectedDate, date)
	}
}

func TestReleaseBranch(t *testing.T) {
	branch := ReleaseBranch()
	if branch == "" {
		t.Error("ReleaseBranch() returned empty string")
	}

	expectedBranch := "master"
	if branch != expectedBranch {
		t.Errorf("Expected release branch %s, got %s", expectedBranch, branch)
	}
}

func TestVersionDetails(t *testing.T) {
	details := VersionDetails()

	requiredFields := []string{
		"version",
		"major",
		"minor",
		"patch",
		"prerelease",
		"sdk_name",
		"user_agent",
		"release_date",
		"release_branch",
	}

	for _, field := range requiredFields {
		if _, ok := details[field]; !ok {
			t.Errorf("VersionDetails() missing required field: %s", field)
		}
	}

	// Verify some specific values
	if details["version"] != "0.0.1" {
		t.Errorf("Expected version 0.0.1, got %v", details["version"])
	}

	if details["sdk_name"] != "Go-Zalo-Bot-SDK" {
		t.Errorf("Expected SDK name Go-Zalo-Bot-SDK, got %v", details["sdk_name"])
	}

	if details["user_agent"] != "Go-Zalo-Bot-SDK/0.0.1" {
		t.Errorf("Expected user agent Go-Zalo-Bot-SDK/0.0.1, got %v", details["user_agent"])
	}
}

func TestReloadVersion(t *testing.T) {
	// This test verifies that ReloadVersion doesn't panic
	// and returns without error when version.yml is valid
	err := ReloadVersion()
	if err != nil {
		t.Errorf("ReloadVersion() returned error: %v", err)
	}

	// Verify version is still correct after reload
	if Version() != "0.0.1" {
		t.Errorf("Version mismatch after reload: got %s, want 0.0.1", Version())
	}
}

func TestVersionFallback(t *testing.T) {
	// This test verifies that the fallback mechanism works
	// We can't easily test this without moving version.yml,
	// but we can verify the fallback constants are defined
	if fallbackSDKName == "" {
		t.Error("fallbackSDKName should not be empty")
	}

	if fallbackVersion == "" {
		t.Error("fallbackVersion should not be empty")
	}

	// Verify fallback values are reasonable
	if fallbackMajor < 0 || fallbackMinor < 0 || fallbackPatch < 0 {
		t.Error("Fallback version components should not be negative")
	}
}
