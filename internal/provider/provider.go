package provider

import (
	"context"
	"strings"

	"github.com/awsnomicon/terraform-provider-awsnomicon/internal/naming"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AwsnomiconProvider struct {
	version string
}

type ProviderData struct {
	OrgPrefix              string
	Project                string
	Env                    string
	Region                 string
	RegionShortCode         string
	RegionMap              map[string]string
	Recipe                 []string
	StylePriority          []string
	ResourceAcronyms        map[string]string
	ResourceStyleOverrides map[string][]string
}

type providerModel struct {
	OrgPrefix              types.String `tfsdk:"org_prefix"`
	Project                types.String `tfsdk:"project"`
	Env                    types.String `tfsdk:"env"`
	Region                 types.String `tfsdk:"region"`
	RegionShortCode         types.String `tfsdk:"region_short_code"`
	Recipe                 types.List   `tfsdk:"recipe"`
	StylePriority          types.List   `tfsdk:"style_priority"`
	RegionMap              types.Map    `tfsdk:"region_map"`
	RegionOverrides        types.Map    `tfsdk:"region_overrides"`
	ResourceAcronyms        types.Map    `tfsdk:"resource_acronyms"`
	ResourceStyleOverrides types.Map    `tfsdk:"resource_style_overrides"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AwsnomiconProvider{version: version}
	}
}

func (p *AwsnomiconProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "awsnomicon"
	resp.Version = p.version
}

func (p *AwsnomiconProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"org_prefix": schema.StringAttribute{
				Required: true,
			},
			"project": schema.StringAttribute{
				Optional: true,
			},
			"env": schema.StringAttribute{
				Required: true,
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
		},
	}
}

func (p *AwsnomiconProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data := &ProviderData{
		OrgPrefix:              config.OrgPrefix.ValueString(),
		Project:                config.Project.ValueString(),
		Env:                    config.Env.ValueString(),
		Region:                 config.Region.ValueString(),
		RegionShortCode:         config.RegionShortCode.ValueString(),
		RegionMap:              naming.DefaultRegionMap(),
		Recipe:                 naming.DefaultRecipe(),
		StylePriority:          naming.DefaultStylePriority(),
		ResourceAcronyms:        naming.DefaultResourceAcronyms(),
		ResourceStyleOverrides: naming.DefaultResourceStyleOverrides(),
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

	resp.DataSourceData = data
	resp.ResourceData = data
}

func (p *AwsnomiconProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewBrewDataSource,
	}
}

func (p *AwsnomiconProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
