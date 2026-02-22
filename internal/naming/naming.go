package naming

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	StyleDashed       = "dashed"
	StyleUnderscore   = "underscore"
	StyleStraight     = "straight"
	StylePascal       = "pascal"
	StylePascalDashed = "pascaldashed"
	StyleCamel        = "camel"
)

var (
	wordRe = regexp.MustCompile(`[A-Za-z0-9]+`)
)

type Config struct {
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
	ResourceConstraints              map[string]ResourceConstraint
	IgnoreRegionForRegionalResources bool
	RegionalResources                map[string]bool
}

type BuildInput struct {
	Resource      string
	Qualifier     string
	Overrides     map[string]string
	Recipe        []string
	StylePriority []string
}

type BuildResult struct {
	Name            string
	Style           string
	Components      map[string]string
	Parts           []string
	RegionCode      string
	ResourceAcronym string
}

type ResourceConstraint struct {
	MinLen              int
	MaxLen              int
	Pattern             *regexp.Regexp
	PatternDescription  string
	ForbiddenPrefixes   []string
	ForbiddenSuffixes   []string
	ForbiddenSubstrings []string
	DisallowIPAddress   bool
	CaseInsensitive     bool
}

func DefaultRecipe() []string {
	return []string{"org", "proj", "env", "region", "resource", "qualifier"}
}

func DefaultStylePriority() []string {
	return []string{StyleDashed, StylePascal, StylePascalDashed, StyleCamel, StyleStraight, StyleUnderscore}
}

func BuildName(cfg Config, in BuildInput) (BuildResult, error) {
	effective := cfg
	if len(effective.RegionMap) == 0 || len(effective.ResourceAcronyms) == 0 || len(effective.ResourceStyleOverrides) == 0 || len(effective.ResourceConstraints) == 0 || len(effective.RegionalResources) == 0 {
		defaults, err := DefaultCloudDefaults(effective.Cloud)
		if err != nil {
			return BuildResult{}, err
		}
		if len(effective.RegionMap) == 0 {
			effective.RegionMap = defaults.RegionMap
		}
		if len(effective.ResourceAcronyms) == 0 {
			effective.ResourceAcronyms = defaults.ResourceAcronyms
		}
		if len(effective.ResourceStyleOverrides) == 0 {
			effective.ResourceStyleOverrides = defaults.ResourceStyleOverrides
		}
		if len(effective.ResourceConstraints) == 0 {
			effective.ResourceConstraints = defaults.ResourceConstraints
		}
		if len(effective.RegionalResources) == 0 {
			effective.RegionalResources = defaults.RegionalResources
		}
	}

	regionCode := strings.TrimSpace(effective.RegionShortCode)
	region := strings.TrimSpace(effective.Region)
	if regionCode == "" && region != "" {
		regionCode = strings.TrimSpace(effective.RegionMap[region])
		if regionCode == "" {
			regionCode = region
		}
	}

	resourceKey := strings.ToLower(strings.TrimSpace(in.Resource))
	resourceAcronym := strings.TrimSpace(in.Resource)
	if resourceKey != "" {
		if v, ok := effective.ResourceAcronyms[resourceKey]; ok && v != "" {
			resourceAcronym = v
		}
	}

	components := map[string]string{
		"org":       strings.TrimSpace(effective.OrgPrefix),
		"proj":      strings.TrimSpace(effective.Project),
		"env":       strings.TrimSpace(effective.Env),
		"region":    strings.TrimSpace(regionCode),
		"resource":  strings.TrimSpace(resourceAcronym),
		"qualifier": strings.TrimSpace(in.Qualifier),
	}

	if effective.IgnoreRegionForRegionalResources && isRegionalResource(resourceKey, effective.RegionalResources) {
		components["region"] = ""
	}

	overrides := in.Overrides
	if overrides == nil {
		overrides = map[string]string{}
	}

	for key, val := range overrides {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		canonical := canonicalComponentKey(key)
		if _, ok := components[canonical]; ok {
			components[canonical] = strings.TrimSpace(val)
		} else {
			components[key] = strings.TrimSpace(val)
		}
	}
	regionCode = components["region"]

	recipe := effective.Recipe
	if len(in.Recipe) > 0 {
		recipe = in.Recipe
	}
	if len(recipe) == 0 {
		recipe = DefaultRecipe()
	}

	parts := make([]string, 0, len(recipe))
	for _, item := range recipe {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		canonical := canonicalComponentKey(item)
		val := ""
		if v, ok := components[canonical]; ok {
			val = v
		} else if v, ok := components[item]; ok {
			val = v
		}
		if strings.TrimSpace(val) == "" {
			continue
		}
		parts = append(parts, val)
	}

	stylePriority := effective.StylePriority
	if len(in.StylePriority) > 0 {
		stylePriority = in.StylePriority
	}
	if len(stylePriority) == 0 {
		stylePriority = DefaultStylePriority()
	}

	allowedStyles := []string{}
	if len(effective.ResourceStyleOverrides) > 0 && resourceKey != "" {
		if v, ok := effective.ResourceStyleOverrides[resourceKey]; ok {
			allowedStyles = normalizeStyles(v)
		}
	}

	chosenStyle := ""
	for _, style := range stylePriority {
		style = normalizeStyle(style)
		if !isValidStyle(style) {
			continue
		}
		if len(allowedStyles) > 0 && !containsString(allowedStyles, style) {
			continue
		}
		chosenStyle = style
		break
	}
	if chosenStyle == "" {
		chosenStyle = StyleDashed
	}

	name, err := formatName(chosenStyle, parts)
	if err != nil {
		return BuildResult{}, err
	}
	if err := validateResourceConstraints(resourceKey, name, effective.ResourceConstraints); err != nil {
		return BuildResult{}, err
	}

	return BuildResult{
		Name:            name,
		Style:           chosenStyle,
		Components:      components,
		Parts:           parts,
		RegionCode:      regionCode,
		ResourceAcronym: components["resource"],
	}, nil
}

func canonicalComponentKey(key string) string {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "org_prefix", "org":
		return "org"
	case "project", "proj":
		return "proj"
	case "environment", "env":
		return "env"
	case "region", "region_code", "region_short_code":
		return "region"
	case "resource", "resource_type":
		return "resource"
	case "what":
		return "resource"
	case "qualifier", "qual":
		return "qualifier"
	default:
		return strings.ToLower(strings.TrimSpace(key))
	}
}

func normalizeStyle(style string) string {
	return strings.ToLower(strings.TrimSpace(style))
}

func isRegionalResource(resourceKey string, regionalResources map[string]bool) bool {
	if resourceKey == "" {
		return false
	}
	return regionalResources[resourceKey]
}

func validateResourceConstraints(resourceKey, name string, constraints map[string]ResourceConstraint) error {
	if resourceKey == "" || len(name) == 0 {
		return nil
	}
	c, ok := constraints[resourceKey]
	if !ok {
		return nil
	}
	if c.MinLen > 0 && len(name) < c.MinLen {
		return fmt.Errorf("resource %q name %q is shorter than %d characters", resourceKey, name, c.MinLen)
	}
	if c.MaxLen > 0 && len(name) > c.MaxLen {
		return fmt.Errorf("resource %q name %q exceeds %d characters", resourceKey, name, c.MaxLen)
	}
	if c.Pattern != nil && !c.Pattern.MatchString(name) {
		desc := c.PatternDescription
		if desc == "" {
			desc = c.Pattern.String()
		}
		return fmt.Errorf("resource %q name %q must match: %s", resourceKey, name, desc)
	}
	comparisonName := name
	if c.CaseInsensitive {
		comparisonName = strings.ToLower(name)
	}
	if len(c.ForbiddenPrefixes) > 0 {
		for _, prefix := range c.ForbiddenPrefixes {
			if prefix == "" {
				continue
			}
			candidate := prefix
			if c.CaseInsensitive {
				candidate = strings.ToLower(prefix)
			}
			if strings.HasPrefix(comparisonName, candidate) {
				return fmt.Errorf("resource %q name %q must not start with prefix %q", resourceKey, name, prefix)
			}
		}
	}
	if len(c.ForbiddenSuffixes) > 0 {
		for _, suffix := range c.ForbiddenSuffixes {
			if suffix == "" {
				continue
			}
			candidate := suffix
			if c.CaseInsensitive {
				candidate = strings.ToLower(suffix)
			}
			if strings.HasSuffix(comparisonName, candidate) {
				return fmt.Errorf("resource %q name %q must not end with suffix %q", resourceKey, name, suffix)
			}
		}
	}
	if len(c.ForbiddenSubstrings) > 0 {
		for _, sub := range c.ForbiddenSubstrings {
			if sub == "" {
				continue
			}
			candidate := sub
			if c.CaseInsensitive {
				candidate = strings.ToLower(sub)
			}
			if strings.Contains(comparisonName, candidate) {
				return fmt.Errorf("resource %q name %q must not contain %q", resourceKey, name, sub)
			}
		}
	}
	if c.DisallowIPAddress && isIPv4Address(name) {
		return fmt.Errorf("resource %q name %q must not be formatted as an IP address", resourceKey, name)
	}
	return nil
}

func isIPv4Address(value string) bool {
	ip := net.ParseIP(value)
	return ip != nil && ip.To4() != nil
}

func normalizeStyles(styles []string) []string {
	out := make([]string, 0, len(styles))
	for _, s := range styles {
		n := normalizeStyle(s)
		if n == "" {
			continue
		}
		out = append(out, n)
	}
	return out
}

func isValidStyle(style string) bool {
	switch style {
	case StyleDashed, StyleUnderscore, StyleStraight, StylePascal, StylePascalDashed, StyleCamel:
		return true
	default:
		return false
	}
}

func formatName(style string, parts []string) (string, error) {
	switch style {
	case StyleDashed:
		return strings.Join(normalizeParts(parts, "-", false), "-"), nil
	case StyleUnderscore:
		return strings.Join(normalizeParts(parts, "_", false), "_"), nil
	case StyleStraight:
		return strings.Join(normalizeParts(parts, "", false), ""), nil
	case StylePascal:
		return strings.Join(normalizeParts(parts, "", true), ""), nil
	case StylePascalDashed:
		return pascalDashedize(parts), nil
	case StyleCamel:
		return camelize(parts), nil
	default:
		return "", errors.New("unsupported style")
	}
}

func normalizeParts(parts []string, sep string, pascal bool) []string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if pascal {
			out = append(out, pascalize(p))
			continue
		}
		words := splitWords(p)
		if len(words) == 0 {
			continue
		}
		if sep == "" {
			out = append(out, strings.ToLower(strings.Join(words, "")))
		} else {
			out = append(out, strings.ToLower(strings.Join(words, sep)))
		}
	}
	return out
}

func camelize(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	firstWords := splitWords(parts[0])
	first := ""
	if len(firstWords) > 0 {
		first = strings.ToLower(strings.Join(firstWords, ""))
	}

	rest := make([]string, 0, len(parts))
	for _, p := range parts[1:] {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		rest = append(rest, pascalize(p))
	}

	return first + strings.Join(rest, "")
}

func pascalDashedize(parts []string) string {
	words := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		for _, w := range splitWords(p) {
			if w == "" {
				continue
			}
			words = append(words, titleWord(w))
		}
	}
	return strings.Join(words, "-")
}

func pascalize(value string) string {
	words := splitWords(value)
	if len(words) == 0 {
		return ""
	}
	out := make([]string, 0, len(words))
	for _, w := range words {
		out = append(out, titleWord(w))
	}
	return strings.Join(out, "")
}

func splitWords(value string) []string {
	return wordRe.FindAllString(value, -1)
}

func titleWord(word string) string {
	if word == "" {
		return ""
	}
	if len(word) == 1 {
		return strings.ToUpper(word)
	}
	return strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
}

func containsString(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
