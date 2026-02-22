package naming

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

//go:embed azure_caf_resource_definition.json
var azureCAFResourceDefinitionJSON []byte

type azureCAFResourceDefinition struct {
	Name            string `json:"name"`
	MinLength       int    `json:"min_length"`
	MaxLength       int    `json:"max_length"`
	ValidationRegex string `json:"validation_regex"`
	Scope           string `json:"scope"`
	Slug            string `json:"slug"`
	Dashes          bool   `json:"dashes"`
	Lowercase       bool   `json:"lowercase"`
}

type azureCloudProfile struct {
	defaultsOnce sync.Once
	defaults     CloudDefaults
	defaultsErr  error
}

func newAzureCloudProfile() CloudProfile {
	return &azureCloudProfile{}
}

func (*azureCloudProfile) Cloud() string {
	return CloudAzure
}

func (p *azureCloudProfile) Defaults() (CloudDefaults, error) {
	p.defaultsOnce.Do(func() {
		var definitions []azureCAFResourceDefinition
		if err := json.Unmarshal(azureCAFResourceDefinitionJSON, &definitions); err != nil {
			p.defaultsErr = fmt.Errorf("decode Azure CAF resource definitions: %w", err)
			return
		}

		acronyms := make(map[string]string, len(definitions))
		styleOverrides := make(map[string][]string, len(definitions))
		constraints := make(map[string]ResourceConstraint, len(definitions))
		regionalResources := make(map[string]bool, len(definitions))

		for _, definition := range definitions {
			name := strings.ToLower(strings.TrimSpace(definition.Name))
			if name == "" {
				continue
			}

			acronyms[name] = azureCAFResourceAcronym(definition.Slug, definition.Name)
			styleOverrides[name] = azureCAFStyleOverrides(definition.Lowercase, definition.Dashes)
			regionalResources[name] = azureCAFIsRegionalScope(definition.Scope)
			constraints[name] = azureCAFConstraint(definition)
		}

		p.defaults = CloudDefaults{
			RegionMap:              map[string]string{},
			ResourceAcronyms:       acronyms,
			ResourceStyleOverrides: styleOverrides,
			ResourceConstraints:    constraints,
			RegionalResources:      regionalResources,
		}
	})

	if p.defaultsErr != nil {
		return CloudDefaults{}, p.defaultsErr
	}
	return copyCloudDefaults(p.defaults), nil
}

func azureCAFConstraint(definition azureCAFResourceDefinition) ResourceConstraint {
	constraint := ResourceConstraint{
		MinLen: definition.MinLength,
		MaxLen: definition.MaxLength,
	}

	regexValue := strings.TrimSpace(definition.ValidationRegex)
	regexValue = strings.Trim(regexValue, `"`)
	if regexValue == "" {
		return constraint
	}

	constraint.PatternDescription = fmt.Sprintf("must match Azure CAF regex %q", regexValue)
	pattern, err := regexp.Compile(regexValue)
	if err != nil {
		return constraint
	}
	constraint.Pattern = pattern
	return constraint
}

func azureCAFStyleOverrides(lowercase, dashes bool) []string {
	styles := []string{}
	if lowercase {
		if dashes {
			styles = append(styles, StyleDashed)
		}
		styles = append(styles, StyleStraight)
		return styles
	}

	if dashes {
		styles = append(styles, StyleDashed, StylePascalDashed)
	}
	styles = append(styles, StylePascal, StyleCamel, StyleStraight)
	return styles
}

func azureCAFResourceAcronym(slug, name string) string {
	base := toLowerAlnum(slug)
	fallback := toLowerAlnum(name)
	for len(base) < 4 && fallback != "" {
		base += fallback[:1]
		fallback = fallback[1:]
	}
	for len(base) < 4 {
		base += "x"
	}
	return base[:4]
}

func azureCAFIsRegionalScope(scope string) bool {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "resourcegroup", "region", "location", "parent":
		return true
	default:
		return false
	}
}

func toLowerAlnum(value string) string {
	var b strings.Builder
	b.Grow(len(value))
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

func copyCloudDefaults(in CloudDefaults) CloudDefaults {
	return CloudDefaults{
		RegionMap:              copyStringMap(in.RegionMap),
		ResourceAcronyms:       copyStringMap(in.ResourceAcronyms),
		ResourceStyleOverrides: copyStringSliceMap(in.ResourceStyleOverrides),
		ResourceConstraints:    copyConstraintMap(in.ResourceConstraints),
		RegionalResources:      copyBoolMap(in.RegionalResources),
	}
}

func copyStringMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func copyBoolMap(in map[string]bool) map[string]bool {
	out := make(map[string]bool, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func copyStringSliceMap(in map[string][]string) map[string][]string {
	out := make(map[string][]string, len(in))
	for key, value := range in {
		cloned := make([]string, len(value))
		copy(cloned, value)
		out[key] = cloned
	}
	return out
}

func copyConstraintMap(in map[string]ResourceConstraint) map[string]ResourceConstraint {
	out := make(map[string]ResourceConstraint, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}
