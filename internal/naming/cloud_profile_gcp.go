package naming

type gcpCloudProfile struct{}

func (gcpCloudProfile) Cloud() string {
	return CloudGCP
}

func (gcpCloudProfile) Defaults() (CloudDefaults, error) {
	return CloudDefaults{
		RegionMap:              DefaultGCPRegionMap(),
		ResourceAcronyms:       DefaultGCPResourceAcronyms(),
		ResourceStyleOverrides: DefaultGCPResourceStyleOverrides(),
		ResourceConstraints:    DefaultGCPResourceConstraints(),
		RegionalResources:      DefaultGCPRegionalResources(),
	}, nil
}
