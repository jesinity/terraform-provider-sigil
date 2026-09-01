// Command generate-databricks-docs renders the coverage table from the naming manifest.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type manifest struct {
	Provider        string     `json:"provider"`
	ProviderVersion string     `json:"provider_version"`
	Resources       []resource `json:"resources"`
}

type resource struct {
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
	SupportedClouds   []string `json:"supported_clouds"`
}

func main() {
	manifestPath := flag.String("manifest", "internal/naming/databricks_resource_definition.json", "manifest path")
	outPath := flag.String("out", "docs/databricks-resources.md", "output path")
	flag.Parse()

	contents, err := os.ReadFile(*manifestPath)
	if err != nil {
		fail(err)
	}
	var data manifest
	if err := json.Unmarshal(contents, &data); err != nil {
		fail(err)
	}
	if err := os.WriteFile(*outPath, render(data), 0o644); err != nil {
		fail(err)
	}
}

func render(data manifest) []byte {
	resources := append([]resource(nil), data.Resources...)
	sort.Slice(resources, func(i, j int) bool { return resources[i].Name < resources[j].Name })
	supported := 0
	for _, resource := range resources {
		if resource.SupportStatus == "supported" {
			supported++
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Databricks platform resources\n\n")
	fmt.Fprintf(&b, "Generated from [`internal/naming/databricks_resource_definition.json`](../internal/naming/databricks_resource_definition.json), audited against `%s` provider version **%s**. This first milestone covers %d name-bearing resources; it does not claim full provider coverage.\n\n", data.Provider, data.ProviderVersion, supported)
	b.WriteString("Supported entries require a canonical Terraform name, acronym, naming attribute, scope, styles, and upstream documentation link. `not_documented` means Sigil intentionally does not invent an API constraint. Path-component scopes are not complete workspace paths. Sigil does not generate identities, email addresses, application IDs, tokens, or secret values.\n\n")
	b.WriteString("| Resource type | Acronym | Category | Naming attribute | Scope | Supported clouds | Regional | Allowed styles | Known constraints | Support status | Upstream documentation |\n| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |\n")
	for _, resource := range resources {
		constraints := resource.ConstraintStatus
		if resource.MinLength != 0 || resource.MaxLength != 0 {
			constraints += fmt.Sprintf("; length %d–%d", resource.MinLength, resource.MaxLength)
		}
		if resource.ValidationRegex != "" {
			constraints += "; documented regex"
		}
		if len(resource.ForbiddenPrefixes) > 0 {
			constraints += "; forbidden prefixes: " + strings.Join(resource.ForbiddenPrefixes, ", ")
		}
		clouds := strings.Join(resource.SupportedClouds, ", ")
		if clouds == "" {
			clouds = "workspace-level"
		}
		fmt.Fprintf(&b, "| `%s` | `%s` | %s | `%s` | %s | %s | %t | %s | %s | %s | [source](%s) |\n", resource.Name, resource.Acronym, resource.Category, resource.NamingAttribute, resource.Scope, clouds, resource.Regional, strings.Join(resource.Styles, ", "), constraints, resource.SupportStatus, resource.DocumentationURL)
	}
	return []byte(b.String())
}

func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
