package services

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"gopkg.in/yaml.v3"
)

// VersionInfo represents the version information structure from version.yml
type VersionInfo struct {
	Version struct {
		Major      int    `yaml:"major"`
		Minor      int    `yaml:"minor"`
		Patch      int    `yaml:"patch"`
		PreRelease string `yaml:"prerelease"`
	} `yaml:"version"`
	SDK struct {
		Name        string `yaml:"name"`
		FullVersion string `yaml:"full_version"`
		UserAgent   string `yaml:"user_agent"`
	} `yaml:"sdk"`
	Release struct {
		Date   string `yaml:"date"`
		Branch string `yaml:"branch"`
	} `yaml:"release"`
}

var (
	// versionCache stores the parsed version information
	versionCache     *VersionInfo
	versionCacheLock sync.RWMutex
	versionCacheOnce sync.Once

	// Fallback constants if version.yml cannot be read
	fallbackSDKName = "Go-Zalo-Bot-SDK"
	fallbackVersion = "0.0.1"
	fallbackMajor   = 0
	fallbackMinor   = 0
	fallbackPatch   = 1
)

// loadVersionInfo loads version information from version.yml
func loadVersionInfo() (*VersionInfo, error) {
	// Try to find version.yml in the module root
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to determine current file path")
	}

	// Get the directory containing version.go (module root)
	moduleRoot := filepath.Dir(currentFile)
	versionFile := filepath.Join(moduleRoot, "version.yml")

	// Read the version.yml file
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read version.yml: %w", err)
	}

	// Parse YAML
	var versionInfo VersionInfo
	if err := yaml.Unmarshal(data, &versionInfo); err != nil {
		return nil, fmt.Errorf("failed to parse version.yml: %w", err)
	}

	return &versionInfo, nil
}

// getVersionInfo returns the cached version information, loading it if necessary
func getVersionInfo() *VersionInfo {
	versionCacheOnce.Do(func() {
		versionCacheLock.Lock()
		defer versionCacheLock.Unlock()

		info, err := loadVersionInfo()
		if err != nil {
			// If we can't load version.yml, create fallback version info
			versionCache = &VersionInfo{}
			versionCache.Version.Major = fallbackMajor
			versionCache.Version.Minor = fallbackMinor
			versionCache.Version.Patch = fallbackPatch
			versionCache.SDK.Name = fallbackSDKName
			versionCache.SDK.FullVersion = fallbackVersion
			versionCache.SDK.UserAgent = fmt.Sprintf("%s/%s", fallbackSDKName, fallbackVersion)
		} else {
			versionCache = info
		}
	})

	versionCacheLock.RLock()
	defer versionCacheLock.RUnlock()
	return versionCache
}

// ReloadVersion forces a reload of version information from version.yml
// This is useful for testing or when version.yml is updated at runtime
func ReloadVersion() error {
	info, err := loadVersionInfo()
	if err != nil {
		return err
	}

	versionCacheLock.Lock()
	versionCache = info
	versionCacheLock.Unlock()

	return nil
}

// Version returns the full semantic version string (e.g., "0.0.1" or "0.1.0-beta")
func Version() string {
	info := getVersionInfo()

	// Use pre-computed full_version from YAML if available
	if info.SDK.FullVersion != "" {
		return info.SDK.FullVersion
	}

	// Otherwise construct it from components
	version := fmt.Sprintf("%d.%d.%d",
		info.Version.Major,
		info.Version.Minor,
		info.Version.Patch)

	if info.Version.PreRelease != "" {
		version += "-" + info.Version.PreRelease
	}

	return version
}

// UserAgent returns the User-Agent string used in HTTP requests
// Format: Go-Zalo-Bot-SDK/0.0.1
func UserAgent() string {
	info := getVersionInfo()

	// Use pre-computed user_agent from YAML if available
	if info.SDK.UserAgent != "" {
		return info.SDK.UserAgent
	}

	// Otherwise construct it
	return fmt.Sprintf("%s/%s", info.SDK.Name, Version())
}

// SDKName returns the SDK name
func SDKName() string {
	info := getVersionInfo()
	if info.SDK.Name != "" {
		return info.SDK.Name
	}
	return fallbackSDKName
}

// VersionMajor returns the major version number
func VersionMajor() int {
	return getVersionInfo().Version.Major
}

// VersionMinor returns the minor version number
func VersionMinor() int {
	return getVersionInfo().Version.Minor
}

// VersionPatch returns the patch version number
func VersionPatch() int {
	return getVersionInfo().Version.Patch
}

// VersionPreRelease returns the pre-release identifier
func VersionPreRelease() string {
	return getVersionInfo().Version.PreRelease
}

// ReleaseDate returns the release date
func ReleaseDate() string {
	return getVersionInfo().Release.Date
}

// ReleaseBranch returns the release branch
func ReleaseBranch() string {
	return getVersionInfo().Release.Branch
}

// VersionDetails returns detailed version information
func VersionDetails() map[string]interface{} {
	info := getVersionInfo()
	return map[string]interface{}{
		"version":        Version(),
		"major":          info.Version.Major,
		"minor":          info.Version.Minor,
		"patch":          info.Version.Patch,
		"prerelease":     info.Version.PreRelease,
		"sdk_name":       info.SDK.Name,
		"user_agent":     UserAgent(),
		"release_date":   info.Release.Date,
		"release_branch": info.Release.Branch,
	}
}
