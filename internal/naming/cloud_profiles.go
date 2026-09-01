package naming

import "fmt"

const (
	CloudAWS   = "aws"
	CloudAzure = "azure"
	CloudGCP   = "gcp"
)

type CloudDefaults struct {
	RegionMap              map[string]string
	ResourceAcronyms       map[string]string
	ResourceStyleOverrides map[string][]string
	ResourceConstraints    map[string]ResourceConstraint
	RegionalResources      map[string]bool
	ResourceClouds         map[string][]string
}

// composeCloudDefaults overlays platform defaults on cloud defaults.  Every map is
// copied so defaults returned to one provider instance can never affect another.
func composeCloudDefaults(cloud, platform string) (CloudDefaults, error) {
	base, err := DefaultCloudDefaults(cloud)
	if err != nil {
		return CloudDefaults{}, err
	}
	platformDefaults, err := DefaultPlatformDefaults(platform)
	if err != nil {
		return CloudDefaults{}, err
	}
	for key, value := range platformDefaults.RegionMap {
		base.RegionMap[key] = value
	}
	for key, value := range platformDefaults.ResourceAcronyms {
		base.ResourceAcronyms[key] = value
	}
	for key, value := range platformDefaults.ResourceStyleOverrides {
		base.ResourceStyleOverrides[key] = append([]string(nil), value...)
	}
	for key, value := range platformDefaults.ResourceConstraints {
		base.ResourceConstraints[key] = value
	}
	for key, value := range platformDefaults.RegionalResources {
		base.RegionalResources[key] = value
	}
	for key, value := range platformDefaults.ResourceClouds {
		base.ResourceClouds[key] = append([]string(nil), value...)
	}
	return base, nil
}

type CloudProfile interface {
	Cloud() string
	Defaults() (CloudDefaults, error)
}

var cloudProfiles = map[string]CloudProfile{
	CloudAWS:   awsCloudProfile{},
	CloudAzure: newAzureCloudProfile(),
	CloudGCP:   gcpCloudProfile{},
}

func DefaultCloud() string {
	return CloudAWS
}

func NormalizeCloud(cloud string) string {
	normalized := normalizeStyle(cloud)
	if normalized == "" {
		return DefaultCloud()
	}
	return normalized
}

func IsSupportedCloud(cloud string) bool {
	_, ok := cloudProfiles[NormalizeCloud(cloud)]
	return ok
}

func DefaultCloudDefaults(cloud string) (CloudDefaults, error) {
	normalized := NormalizeCloud(cloud)
	profile, ok := cloudProfiles[normalized]
	if !ok {
		return CloudDefaults{}, fmt.Errorf("unsupported cloud %q", cloud)
	}
	return profile.Defaults()
}
