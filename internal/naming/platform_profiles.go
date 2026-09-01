package naming

import (
	"fmt"
	"strings"
)

const PlatformDatabricks = "databricks"

type PlatformProfile interface {
	Platform() string
	Defaults() (CloudDefaults, error)
	Aliases() map[string]string
}

var platformProfiles = map[string]PlatformProfile{
	PlatformDatabricks: newDatabricksPlatformProfile(),
}

func NormalizePlatform(platform string) string { return strings.ToLower(strings.TrimSpace(platform)) }

func IsSupportedPlatform(platform string) bool {
	platform = NormalizePlatform(platform)
	if platform == "" {
		return true
	}
	_, ok := platformProfiles[platform]
	return ok
}

func DefaultPlatformDefaults(platform string) (CloudDefaults, error) {
	platform = NormalizePlatform(platform)
	if platform == "" {
		return CloudDefaults{RegionMap: map[string]string{}, ResourceAcronyms: map[string]string{}, ResourceStyleOverrides: map[string][]string{}, ResourceConstraints: map[string]ResourceConstraint{}, RegionalResources: map[string]bool{}, ResourceClouds: map[string][]string{}}, nil
	}
	profile, ok := platformProfiles[platform]
	if !ok {
		return CloudDefaults{}, fmt.Errorf("unsupported platform %q", platform)
	}
	return profile.Defaults()
}

// DefaultDefaults returns isolated, composed defaults for a cloud/platform pair.
func DefaultDefaults(cloud, platform string) (CloudDefaults, error) {
	return composeCloudDefaults(cloud, platform)
}

func platformAliases(platform string) map[string]string {
	profile, ok := platformProfiles[NormalizePlatform(platform)]
	if !ok {
		return nil
	}
	return profile.Aliases()
}
