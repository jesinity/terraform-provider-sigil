package naming

type awsCloudProfile struct{}

func (awsCloudProfile) Cloud() string {
	return CloudAWS
}

func (awsCloudProfile) Defaults() (CloudDefaults, error) {
	return CloudDefaults{
		RegionMap:              DefaultRegionMap(),
		ResourceAcronyms:       DefaultResourceAcronyms(),
		ResourceStyleOverrides: DefaultResourceStyleOverrides(),
		ResourceConstraints:    DefaultResourceConstraints(),
		RegionalResources:      DefaultRegionalResources(),
	}, nil
}
