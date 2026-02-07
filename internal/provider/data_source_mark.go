package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jesinity/terraform-provider-sigil/internal/naming"
)

type MarkDataSource struct {
	providerData *ProviderData
}

type markDataSourceModel struct {
	Resource        types.String `tfsdk:"resource"`
	What            types.String `tfsdk:"what"`
	Qualifier       types.String `tfsdk:"qualifier"`
	Overrides       types.Map    `tfsdk:"overrides"`
	Recipe          types.List   `tfsdk:"recipe"`
	StylePriority   types.List   `tfsdk:"style_priority"`
	Name            types.String `tfsdk:"name"`
	Style           types.String `tfsdk:"style"`
	RegionCode      types.String `tfsdk:"region_code"`
	ResourceAcronym types.String `tfsdk:"resource_acronym"`
	Components      types.Map    `tfsdk:"components"`
	Parts           types.List   `tfsdk:"parts"`
}

func NewMarkDataSource() datasource.DataSource {
	return &MarkDataSource{}
}

func (d *MarkDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mark"
}

func (d *MarkDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource": schema.StringAttribute{
				Optional:           true,
				DeprecationMessage: "Use `what` instead.",
			},
			"what": schema.StringAttribute{
				Optional: true,
			},
			"qualifier": schema.StringAttribute{
				Optional: true,
			},
			"overrides": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"recipe": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"style_priority": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"style": schema.StringAttribute{
				Computed: true,
			},
			"region_code": schema.StringAttribute{
				Computed: true,
			},
			"resource_acronym": schema.StringAttribute{
				Computed: true,
			},
			"components": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"parts": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *MarkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*ProviderData)
	if !ok {
		return
	}

	d.providerData = providerData
}

func (d *MarkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.providerData == nil {
		resp.Diagnostics.AddError("Provider not configured", "The provider has not been configured yet.")
		return
	}

	var data markDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	overrides := map[string]string{}
	if !data.Overrides.IsNull() && !data.Overrides.IsUnknown() {
		resp.Diagnostics.Append(data.Overrides.ElementsAs(ctx, &overrides, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	recipe := []string{}
	if !data.Recipe.IsNull() && !data.Recipe.IsUnknown() {
		resp.Diagnostics.Append(data.Recipe.ElementsAs(ctx, &recipe, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	stylePriority := []string{}
	if !data.StylePriority.IsNull() && !data.StylePriority.IsUnknown() {
		resp.Diagnostics.Append(data.StylePriority.ElementsAs(ctx, &stylePriority, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	what := strings.TrimSpace(data.What.ValueString())
	resource := strings.TrimSpace(data.Resource.ValueString())
	if what == "" && resource == "" {
		resp.Diagnostics.AddError("Missing required attribute", "Either `what` (preferred) or `resource` must be set.")
		return
	}
	if what != "" && resource != "" && what != resource {
		resp.Diagnostics.AddError("Conflicting attributes", "`what` and `resource` cannot both be set to different values.")
		return
	}
	if what == "" {
		what = resource
	}

	result, err := naming.BuildName(naming.Config{
		OrgPrefix:                        d.providerData.OrgPrefix,
		Project:                          d.providerData.Project,
		Env:                              d.providerData.Env,
		Region:                           d.providerData.Region,
		RegionShortCode:                  d.providerData.RegionShortCode,
		RegionMap:                        d.providerData.RegionMap,
		Recipe:                           d.providerData.Recipe,
		StylePriority:                    d.providerData.StylePriority,
		ResourceAcronyms:                 d.providerData.ResourceAcronyms,
		ResourceStyleOverrides:           d.providerData.ResourceStyleOverrides,
		IgnoreRegionForRegionalResources: d.providerData.IgnoreRegionForRegionalResources,
		RegionalResources:                d.providerData.RegionalResources,
	}, naming.BuildInput{
		Resource:      what,
		Qualifier:     data.Qualifier.ValueString(),
		Overrides:     overrides,
		Recipe:        recipe,
		StylePriority: stylePriority,
	})
	if err != nil {
		resp.Diagnostics.AddError("Name build failed", err.Error())
		return
	}

	data.Name = types.StringValue(result.Name)
	data.Style = types.StringValue(result.Style)
	data.RegionCode = types.StringValue(result.RegionCode)
	data.ResourceAcronym = types.StringValue(result.ResourceAcronym)

	componentsValue, diags := types.MapValueFrom(ctx, types.StringType, result.Components)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Components = componentsValue

	partsValue, diags := types.ListValueFrom(ctx, types.StringType, result.Parts)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Parts = partsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
