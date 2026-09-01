// Command compare-databricks-catalog reports upstream resource drift; it never modifies the catalog.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type catalog struct {
	Resources []struct {
		Name string `json:"name"`
	} `json:"resources"`
}

func main() {
	catalogPath := flag.String("catalog", "internal/naming/databricks_resource_catalog.json", "catalog path")
	docsDir := flag.String("docs-dir", "", "upstream docs/resources directory")
	flag.Parse()
	if *docsDir == "" {
		fail("-docs-dir is required")
	}
	contents, err := os.ReadFile(*catalogPath)
	if err != nil {
		fail(err.Error())
	}
	var c catalog
	if err := json.Unmarshal(contents, &c); err != nil {
		fail(err.Error())
	}
	catalogNames := map[string]bool{}
	for _, resource := range c.Resources {
		catalogNames[resource.Name] = true
	}
	entries, err := os.ReadDir(*docsDir)
	if err != nil {
		fail(err.Error())
	}
	upstream := map[string]bool{}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			upstream["databricks_"+strings.TrimSuffix(filepath.Base(entry.Name()), ".md")] = true
		}
	}
	added, removed := difference(upstream, catalogNames), difference(catalogNames, upstream)
	fmt.Printf("catalog: %d; upstream: %d\n", len(catalogNames), len(upstream))
	printList("added upstream resources (require audit)", added)
	printList("removed upstream resources (retain/review classification)", removed)
	if len(added)+len(removed) > 0 {
		os.Exit(1)
	}
}

func difference(left, right map[string]bool) []string {
	var out []string
	for name := range left {
		if !right[name] {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}
func printList(label string, names []string) {
	fmt.Printf("%s: %d\n", label, len(names))
	for _, name := range names {
		fmt.Printf("  %s\n", name)
	}
}
func fail(message string) { fmt.Fprintln(os.Stderr, message); os.Exit(2) }
