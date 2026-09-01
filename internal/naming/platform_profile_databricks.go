package naming

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

//go:generate go run ../../cmd/generate-databricks-docs -manifest databricks_resource_definition.json -out ../../docs/databricks-resources.md
//go:embed databricks_resource_definition.json
var databricksResourceDefinitionJSON []byte

type databricksManifest struct {
	Provider        string                         `json:"provider"`
	ProviderVersion string                         `json:"provider_version"`
	Resources       []databricksResourceDefinition `json:"resources"`
}

type databricksResourceDefinition struct {
	Name              string   `json:"name"`
	Category          string   `json:"category"`
	Acronym           string   `json:"acronym"`
	NamingAttribute   string   `json:"naming_attribute"`
	Scope             string   `json:"scope"`
	Styles            []string `json:"styles"`
	MinLength         int      `json:"min_length"`
	MaxLength         int      `json:"max_length"`
	ValidationRegex   string   `json:"validation_regex"`
	ForbiddenPrefixes []string `json:"forbidden_prefixes"`
	Regional          bool     `json:"regional"`
	SupportStatus     string   `json:"support_status"`
	ConstraintStatus  string   `json:"constraint_status"`
	DocumentationURL  string   `json:"documentation_url"`
	Aliases           []string `json:"aliases"`
	SupportedClouds   []string `json:"supported_clouds"`
}

type databricksPlatformProfile struct {
	once     sync.Once
	defaults CloudDefaults
	aliases  map[string]string
	err      error
}

func newDatabricksPlatformProfile() PlatformProfile { return &databricksPlatformProfile{} }
func (*databricksPlatformProfile) Platform() string { return PlatformDatabricks }
func (p *databricksPlatformProfile) Defaults() (CloudDefaults, error) {
	p.once.Do(p.load)
	if p.err != nil {
		return CloudDefaults{}, p.err
	}
	return copyCloudDefaults(p.defaults), nil
}
func (p *databricksPlatformProfile) Aliases() map[string]string {
	p.once.Do(p.load)
	return copyStringMap(p.aliases)
}
func (p *databricksPlatformProfile) load() {
	p.defaults, p.aliases, p.err = loadDatabricksPlatformDefaults(databricksResourceDefinitionJSON)
}

func loadDatabricksPlatformDefaults(contents []byte) (CloudDefaults, map[string]string, error) {
	var manifest databricksManifest
	if err := json.Unmarshal(contents, &manifest); err != nil {
		return CloudDefaults{}, nil, fmt.Errorf("decode Databricks resource manifest: %w", err)
	}
	if manifest.Provider != "databricks/databricks" || strings.TrimSpace(manifest.ProviderVersion) == "" {
		return CloudDefaults{}, nil, fmt.Errorf("Databricks manifest must declare provider databricks/databricks and a provider version")
	}
	defaults := CloudDefaults{RegionMap: map[string]string{}, ResourceAcronyms: map[string]string{}, ResourceStyleOverrides: map[string][]string{}, ResourceConstraints: map[string]ResourceConstraint{}, RegionalResources: map[string]bool{}, ResourceClouds: map[string][]string{}}
	aliases := map[string]string{}
	acronyms := map[string]bool{}
	for _, resource := range manifest.Resources {
		name := strings.ToLower(strings.TrimSpace(resource.Name))
		status := strings.TrimSpace(resource.SupportStatus)
		if name == "" || (!strings.HasPrefix(name, "databricks_") && name != "azurerm_databricks_workspace") {
			return CloudDefaults{}, nil, fmt.Errorf("Databricks resource name %q must use its canonical Terraform provider prefix", resource.Name)
		}
		if _, exists := defaults.ResourceAcronyms[name]; exists {
			return CloudDefaults{}, nil, fmt.Errorf("duplicate Databricks canonical resource %q", name)
		}
		if status != "supported" && status != "not_name_bearing" && status != "deferred" && status != "deprecated" {
			return CloudDefaults{}, nil, fmt.Errorf("Databricks resource %q has unsupported support_status %q", name, status)
		}
		for _, cloud := range resource.SupportedClouds {
			if !IsSupportedCloud(cloud) {
				return CloudDefaults{}, nil, fmt.Errorf("Databricks resource %q has unsupported cloud %q", name, cloud)
			}
		}
		if resource.ConstraintStatus != "audited" && resource.ConstraintStatus != "partially_documented" && resource.ConstraintStatus != "not_documented" && resource.ConstraintStatus != "not_applicable" {
			return CloudDefaults{}, nil, fmt.Errorf("Databricks resource %q has unsupported constraint_status %q", name, resource.ConstraintStatus)
		}
		if status != "supported" {
			continue
		}
		if strings.TrimSpace(resource.Acronym) == "" || strings.TrimSpace(resource.NamingAttribute) == "" || strings.TrimSpace(resource.Scope) == "" {
			return CloudDefaults{}, nil, fmt.Errorf("supported Databricks resource %q requires acronym, naming_attribute, and scope", name)
		}
		if strings.TrimSpace(resource.DocumentationURL) == "" {
			return CloudDefaults{}, nil, fmt.Errorf("supported Databricks resource %q requires documentation_url", name)
		}
		if acronyms[resource.Acronym] {
			return CloudDefaults{}, nil, fmt.Errorf("duplicate Databricks acronym %q", resource.Acronym)
		}
		acronyms[resource.Acronym] = true
		styles := normalizeStyles(resource.Styles)
		if len(styles) == 0 {
			return CloudDefaults{}, nil, fmt.Errorf("supported Databricks resource %q requires supported styles", name)
		}
		if resource.MinLength < 0 || resource.MaxLength < 0 || resource.MaxLength > 0 && resource.MinLength > resource.MaxLength {
			return CloudDefaults{}, nil, fmt.Errorf("Databricks resource %q has invalid length constraints", name)
		}
		constraint := ResourceConstraint{MinLen: resource.MinLength, MaxLen: resource.MaxLength, ForbiddenPrefixes: append([]string(nil), resource.ForbiddenPrefixes...)}
		if patternText := strings.TrimSpace(resource.ValidationRegex); patternText != "" {
			pattern, err := regexp.Compile(patternText)
			if err != nil {
				return CloudDefaults{}, nil, fmt.Errorf("Databricks resource %q has invalid validation_regex: %w", name, err)
			}
			constraint.Pattern = pattern
			constraint.PatternDescription = "must match Databricks documented pattern"
		}
		defaults.ResourceAcronyms[name] = resource.Acronym
		defaults.ResourceStyleOverrides[name] = styles
		defaults.ResourceConstraints[name] = constraint
		defaults.RegionalResources[name] = resource.Regional
		if len(resource.SupportedClouds) > 0 {
			defaults.ResourceClouds[name] = append([]string(nil), resource.SupportedClouds...)
		}
		for _, alias := range resource.Aliases {
			alias = strings.ToLower(strings.TrimSpace(alias))
			if alias == "" || strings.HasPrefix(alias, "databricks_") {
				return CloudDefaults{}, nil, fmt.Errorf("Databricks resource %q has invalid compatibility alias %q", name, alias)
			}
			if existing := aliases[alias]; existing != "" && existing != name {
				return CloudDefaults{}, nil, fmt.Errorf("ambiguous Databricks alias %q", alias)
			}
			aliases[alias] = name
		}
	}
	return defaults, aliases, nil
}
