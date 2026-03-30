package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/jesinity/terraform-provider-sigil/internal/naming"
)

type SigilProvider struct {
	version string
}

type ProviderData struct {
	Cloud                            string
	OrgPrefix                        string
	Project                          string
	Env                              string
	Region                           string
	RegionShortCode                  string
	RegionMap                        map[string]string
	Recipe                           []string
	StylePriority                    []string
	ResourceAcronyms                 map[string]string
	ResourceStyleOverrides           map[string][]string
	ResourceConstraints              map[string]naming.ResourceConstraint
	IgnoreRegionForRegionalResources bool
	RegionalResources                map[string]bool
}

type providerModel struct {
	Config                           types.Object `tfsdk:"config"`
	Overrides                        types.Object `tfsdk:"overrides"`
	Cloud                            types.String `tfsdk:"cloud"`
	OrgPrefix                        types.String `tfsdk:"org_prefix"`
	Project                          types.String `tfsdk:"project"`
	Env                              types.String `tfsdk:"env"`
	Region                           types.String `tfsdk:"region"`
	RegionShortCode                  types.String `tfsdk:"region_short_code"`
	Recipe                           types.List   `tfsdk:"recipe"`
	StylePriority                    types.List   `tfsdk:"style_priority"`
	RegionMap                        types.Map    `tfsdk:"region_map"`
	RegionOverrides                  types.Map    `tfsdk:"region_overrides"`
	ResourceAcronyms                 types.Map    `tfsdk:"resource_acronyms"`
	ResourceStyleOverrides           types.Map    `tfsdk:"resource_style_overrides"`
	IgnoreRegionForRegionalResources types.Bool   `tfsdk:"ignore_region_for_regional_resources"`
}

type providerConfigModel struct {
	Cloud                            types.String `tfsdk:"cloud"`
	OrgPrefix                        types.String `tfsdk:"org_prefix"`
	Project                          types.String `tfsdk:"project"`
	Env                              types.String `tfsdk:"env"`
	Region                           types.String `tfsdk:"region"`
	RegionShortCode                  types.String `tfsdk:"region_short_code"`
	Recipe                           types.List   `tfsdk:"recipe"`
	StylePriority                    types.List   `tfsdk:"style_priority"`
	RegionMap                        types.Map    `tfsdk:"region_map"`
	RegionOverrides                  types.Map    `tfsdk:"region_overrides"`
	ResourceAcronyms                 types.Map    `tfsdk:"resource_acronyms"`
	ResourceStyleOverrides           types.Map    `tfsdk:"resource_style_overrides"`
	IgnoreRegionForRegionalResources types.Bool   `tfsdk:"ignore_region_for_regional_resources"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SigilProvider{version: version}
	}
}

func (p *SigilProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sigil"
	resp.Version = p.version
}

func (p *SigilProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	providerConfigAttributes := providerConfigSchemaAttributes()
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"config": schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: providerConfigSchemaAttributes(),
			},
			"overrides": schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: providerConfigSchemaAttributes(),
			},
		},
	}
	for key, attr := range providerConfigAttributes {
		resp.Schema.Attributes[key] = attr
	}
}

func (p *SigilProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	baseConfig, hasBaseConfig := decodeProviderConfigObject(ctx, config.Config, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	overrideConfig, hasOverrideConfig := decodeProviderConfigObject(ctx, config.Overrides, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	cloud := resolveCloud(config.Cloud, baseConfig, hasBaseConfig, overrideConfig, hasOverrideConfig)
	if !naming.IsSupportedCloud(cloud) {
		resp.Diagnostics.AddError("Invalid cloud", fmt.Sprintf("Unsupported cloud %q. Valid values are %q, %q, and %q.", cloud, naming.CloudAWS, naming.CloudAzure, naming.CloudGCP))
		return
	}

	cloudDefaults, err := naming.DefaultCloudDefaults(cloud)
	if err != nil {
		resp.Diagnostics.AddError("Cloud defaults error", err.Error())
		return
	}

	data := &ProviderData{
		Cloud:                            cloud,
		OrgPrefix:                        "",
		Project:                          "",
		Env:                              "",
		Region:                           "",
		RegionShortCode:                  "",
		RegionMap:                        cloudDefaults.RegionMap,
		Recipe:                           naming.DefaultRecipe(),
		StylePriority:                    naming.DefaultStylePriority(),
		ResourceAcronyms:                 cloudDefaults.ResourceAcronyms,
		ResourceStyleOverrides:           cloudDefaults.ResourceStyleOverrides,
		ResourceConstraints:              cloudDefaults.ResourceConstraints,
		IgnoreRegionForRegionalResources: true,
		RegionalResources:                cloudDefaults.RegionalResources,
	}

	if hasBaseConfig {
		applyProviderConfig(ctx, resp, data, baseConfig)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	applyProviderConfig(ctx, resp, data, providerConfigFromModel(config))
	if resp.Diagnostics.HasError() {
		return
	}

	if hasOverrideConfig {
		applyProviderConfig(ctx, resp, data, overrideConfig)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if strings.TrimSpace(data.OrgPrefix) == "" {
		resp.Diagnostics.AddError("Missing org_prefix", "Set org_prefix at the top level, inside config, or inside overrides.")
	}
	if strings.TrimSpace(data.Env) == "" {
		resp.Diagnostics.AddError("Missing env", "Set env at the top level, inside config, or inside overrides.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = data
	resp.ResourceData = data
}

func (p *SigilProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMarkDataSource,
	}
}

func (p *SigilProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func providerConfigSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cloud": schema.StringAttribute{
			Optional: true,
		},
		"org_prefix": schema.StringAttribute{
			Optional: true,
		},
		"project": schema.StringAttribute{
			Optional: true,
		},
		"env": schema.StringAttribute{
			Optional: true,
		},
		"region": schema.StringAttribute{
			Optional: true,
		},
		"region_short_code": schema.StringAttribute{
			Optional: true,
		},
		"recipe": schema.ListAttribute{
			Optional:    true,
			ElementType: types.StringType,
		},
		"style_priority": schema.ListAttribute{
			Optional:    true,
			ElementType: types.StringType,
		},
		"region_map": schema.MapAttribute{
			Optional:    true,
			ElementType: types.StringType,
		},
		"region_overrides": schema.MapAttribute{
			Optional:    true,
			ElementType: types.StringType,
		},
		"resource_acronyms": schema.MapAttribute{
			Optional:    true,
			ElementType: types.StringType,
		},
		"resource_style_overrides": schema.MapAttribute{
			Optional:    true,
			ElementType: types.ListType{ElemType: types.StringType},
		},
		"ignore_region_for_regional_resources": schema.BoolAttribute{
			Optional: true,
		},
	}
}

func providerConfigFromModel(config providerModel) providerConfigModel {
	return providerConfigModel{
		Cloud:                            config.Cloud,
		OrgPrefix:                        config.OrgPrefix,
		Project:                          config.Project,
		Env:                              config.Env,
		Region:                           config.Region,
		RegionShortCode:                  config.RegionShortCode,
		Recipe:                           config.Recipe,
		StylePriority:                    config.StylePriority,
		RegionMap:                        config.RegionMap,
		RegionOverrides:                  config.RegionOverrides,
		ResourceAcronyms:                 config.ResourceAcronyms,
		ResourceStyleOverrides:           config.ResourceStyleOverrides,
		IgnoreRegionForRegionalResources: config.IgnoreRegionForRegionalResources,
	}
}

func decodeProviderConfigObject(ctx context.Context, obj types.Object, resp *provider.ConfigureResponse) (providerConfigModel, bool) {
	var decoded providerConfigModel
	if obj.IsNull() || obj.IsUnknown() {
		return decoded, false
	}
	resp.Diagnostics.Append(obj.As(ctx, &decoded, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return decoded, false
	}
	return decoded, true
}

func resolveCloud(topLevelCloud types.String, baseConfig providerConfigModel, hasBaseConfig bool, overrideConfig providerConfigModel, hasOverrideConfig bool) string {
	cloud := naming.DefaultCloud()
	if hasBaseConfig && !baseConfig.Cloud.IsNull() && !baseConfig.Cloud.IsUnknown() {
		cloud = naming.NormalizeCloud(baseConfig.Cloud.ValueString())
	}
	if !topLevelCloud.IsNull() && !topLevelCloud.IsUnknown() {
		cloud = naming.NormalizeCloud(topLevelCloud.ValueString())
	}
	if hasOverrideConfig && !overrideConfig.Cloud.IsNull() && !overrideConfig.Cloud.IsUnknown() {
		cloud = naming.NormalizeCloud(overrideConfig.Cloud.ValueString())
	}
	return cloud
}

func applyProviderConfig(ctx context.Context, resp *provider.ConfigureResponse, data *ProviderData, config providerConfigModel) {
	if !config.Cloud.IsNull() && !config.Cloud.IsUnknown() {
		data.Cloud = naming.NormalizeCloud(config.Cloud.ValueString())
	}
	if !config.OrgPrefix.IsNull() && !config.OrgPrefix.IsUnknown() {
		data.OrgPrefix = config.OrgPrefix.ValueString()
	}
	if !config.Project.IsNull() && !config.Project.IsUnknown() {
		data.Project = config.Project.ValueString()
	}
	if !config.Env.IsNull() && !config.Env.IsUnknown() {
		data.Env = config.Env.ValueString()
	}
	if !config.Region.IsNull() && !config.Region.IsUnknown() {
		data.Region = config.Region.ValueString()
	}
	if !config.RegionShortCode.IsNull() && !config.RegionShortCode.IsUnknown() {
		data.RegionShortCode = config.RegionShortCode.ValueString()
	}
	if !config.IgnoreRegionForRegionalResources.IsNull() && !config.IgnoreRegionForRegionalResources.IsUnknown() {
		data.IgnoreRegionForRegionalResources = config.IgnoreRegionForRegionalResources.ValueBool()
	}
	if !config.RegionMap.IsNull() && !config.RegionMap.IsUnknown() {
		regionMap := map[string]string{}
		resp.Diagnostics.Append(config.RegionMap.ElementsAs(ctx, &regionMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(regionMap) > 0 {
			data.RegionMap = regionMap
		}
	}
	if !config.RegionOverrides.IsNull() && !config.RegionOverrides.IsUnknown() {
		overrides := map[string]string{}
		resp.Diagnostics.Append(config.RegionOverrides.ElementsAs(ctx, &overrides, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for key, val := range overrides {
			data.RegionMap[key] = val
		}
	}
	if !config.Recipe.IsNull() && !config.Recipe.IsUnknown() {
		recipe := []string{}
		resp.Diagnostics.Append(config.Recipe.ElementsAs(ctx, &recipe, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(recipe) > 0 {
			data.Recipe = recipe
		}
	}
	if !config.StylePriority.IsNull() && !config.StylePriority.IsUnknown() {
		styles := []string{}
		resp.Diagnostics.Append(config.StylePriority.ElementsAs(ctx, &styles, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(styles) > 0 {
			data.StylePriority = styles
		}
	}
	if !config.ResourceAcronyms.IsNull() && !config.ResourceAcronyms.IsUnknown() {
		acronyms := map[string]string{}
		resp.Diagnostics.Append(config.ResourceAcronyms.ElementsAs(ctx, &acronyms, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for key, val := range acronyms {
			data.ResourceAcronyms[strings.ToLower(key)] = val
		}
	}
	if !config.ResourceStyleOverrides.IsNull() && !config.ResourceStyleOverrides.IsUnknown() {
		overrides := map[string][]string{}
		for key, value := range config.ResourceStyleOverrides.Elements() {
			list, ok := value.(types.List)
			if !ok {
				continue
			}
			styles := []string{}
			resp.Diagnostics.Append(list.ElementsAs(ctx, &styles, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			overrides[strings.ToLower(key)] = styles
		}
		for key, styles := range overrides {
			data.ResourceStyleOverrides[key] = styles
		}
	}
}
