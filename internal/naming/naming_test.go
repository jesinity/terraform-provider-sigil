package naming

import "testing"

func TestDefaultCloudDefaultsAzure(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudAzure)
	if err != nil {
		t.Fatalf("unexpected error loading Azure defaults: %v", err)
	}

	if len(defaults.RegionMap) == 0 {
		t.Fatal("expected Azure region map defaults to be populated")
	}

	if len(defaults.ResourceAcronyms) < 300 {
		t.Fatalf("expected Azure CAF acronyms to be populated, got %d entries", len(defaults.ResourceAcronyms))
	}

	for resource, acronym := range defaults.ResourceAcronyms {
		if acronym == "" {
			t.Fatalf("resource %q has empty acronym", resource)
		}
	}

	if got := defaults.ResourceAcronyms["azurerm_storage_account"]; got != "st" {
		t.Fatalf("expected CAF storage account acronym %q, got %q", "st", got)
	}

	if got := defaults.ResourceAcronyms["azurerm_resource_group"]; got != "rg" {
		t.Fatalf("expected CAF resource group acronym %q, got %q", "rg", got)
	}

	if got := defaults.ResourceAcronyms["azurerm_virtual_machine"]; got != "vm" {
		t.Fatalf("expected CAF virtual machine acronym %q, got %q", "vm", got)
	}

	if got := defaults.ResourceAcronyms["azurerm_linux_virtual_machine"]; got != "vm" {
		t.Fatalf("expected CAF linux virtual machine acronym %q, got %q", "vm", got)
	}

	if got := defaults.ResourceAcronyms["azurerm_windows_virtual_machine"]; got != "vm" {
		t.Fatalf("expected CAF windows virtual machine acronym %q, got %q", "vm", got)
	}

	if got := defaults.ResourceAcronyms["azurerm_api_management"]; got != "apim" {
		t.Fatalf("expected CAF API management acronym %q, got %q", "apim", got)
	}

	if got := defaults.ResourceAcronyms["azurerm_api_management_group"]; got != "apimgr" {
		t.Fatalf("expected CAF API management group acronym %q, got %q", "apimgr", got)
	}

	if got := defaults.ResourceAcronyms["azurerm_api_management_logger"]; got != "apimlg" {
		t.Fatalf("expected CAF API management logger acronym %q, got %q", "apimlg", got)
	}

	if got := defaults.ResourceAcronyms["azurerm_api_management_service"]; got != "apim" {
		t.Fatalf("expected CAF API management service acronym %q, got %q", "apim", got)
	}

	if _, ok := defaults.ResourceAcronyms["general"]; ok {
		t.Fatalf("expected no acronym entry for %q because CAF slug is empty", "general")
	}

	if got := defaults.RegionMap["westeurope"]; got != "weu" {
		t.Fatalf("expected Azure region short code %q for westeurope, got %q", "weu", got)
	}
	if got := defaults.RegionMap["eastus2"]; got != "eus2" {
		t.Fatalf("expected Azure region short code %q for eastus2, got %q", "eus2", got)
	}
}

func TestBuildNameAzureStorageAccountSelectsStraightStyle(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudAzure)
	if err != nil {
		t.Fatalf("unexpected error loading Azure defaults: %v", err)
	}

	result, err := BuildName(Config{
		StylePriority:          DefaultStylePriority(),
		ResourceAcronyms:       defaults.ResourceAcronyms,
		ResourceStyleOverrides: defaults.ResourceStyleOverrides,
		ResourceConstraints:    defaults.ResourceConstraints,
		RegionalResources:      defaults.RegionalResources,
	}, BuildInput{
		Resource:  "azurerm_storage_account",
		Qualifier: "data",
		Recipe:    []string{"resource", "qualifier"},
	})
	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}

	if result.Style != StyleStraight {
		t.Fatalf("expected style %q, got %q", StyleStraight, result.Style)
	}

	if result.Name == "" {
		t.Fatal("expected a generated name")
	}
}

func TestBuildNameAzureRegionLookupNormalizesSeparators(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudAzure)
	if err != nil {
		t.Fatalf("unexpected error loading Azure defaults: %v", err)
	}

	result, err := BuildName(Config{
		Cloud:                            CloudAzure,
		Region:                           "west-europe",
		RegionMap:                        defaults.RegionMap,
		IgnoreRegionForRegionalResources: false,
		ResourceAcronyms:                 defaults.ResourceAcronyms,
		ResourceStyleOverrides:           defaults.ResourceStyleOverrides,
		ResourceConstraints:              defaults.ResourceConstraints,
		RegionalResources:                defaults.RegionalResources,
	}, BuildInput{
		Resource: "azurerm_storage_account",
		Recipe:   []string{"region"},
	})
	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}

	if result.RegionCode != "weu" {
		t.Fatalf("expected region code %q, got %q", "weu", result.RegionCode)
	}
	if result.Name != "weu" {
		t.Fatalf("expected generated name %q, got %q", "weu", result.Name)
	}
}

func TestBuildNameAzureFallsBackToAllowedStyleWhenPriorityDoesNotMatch(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudAzure)
	if err != nil {
		t.Fatalf("unexpected error loading Azure defaults: %v", err)
	}

	result, err := BuildName(Config{
		Cloud:                            CloudAzure,
		OrgPrefix:                        "acme",
		Project:                          "payments",
		Env:                              "prod",
		Region:                           "westeurope",
		RegionMap:                        defaults.RegionMap,
		IgnoreRegionForRegionalResources: false,
		StylePriority:                    []string{StylePascal}, // not allowed for storage accounts
		ResourceAcronyms:                 defaults.ResourceAcronyms,
		ResourceStyleOverrides:           defaults.ResourceStyleOverrides,
		ResourceConstraints:              defaults.ResourceConstraints,
		RegionalResources:                defaults.RegionalResources,
	}, BuildInput{
		Resource:  "azurerm_storage_account",
		Qualifier: "raw",
		Recipe:    []string{"org", "proj", "env", "resource", "qualifier"},
	})
	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}

	if result.Style != StyleStraight {
		t.Fatalf("expected fallback style %q, got %q", StyleStraight, result.Style)
	}
}

func TestDefaultCloudDefaultsAzureRestrictiveStylesFromCAF(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudAzure)
	if err != nil {
		t.Fatalf("unexpected error loading Azure defaults: %v", err)
	}

	if got := defaults.ResourceStyleOverrides["azurerm_storage_account"]; len(got) != 1 || got[0] != StyleStraight {
		t.Fatalf("expected storage account styles to be [%q], got %#v", StyleStraight, got)
	}

	if got := defaults.ResourceStyleOverrides["azurerm_analysis_services_server"]; len(got) != 1 || got[0] != StyleStraight {
		t.Fatalf("expected analysis services server styles to be [%q], got %#v", StyleStraight, got)
	}

	if got := defaults.ResourceStyleOverrides["azurerm_cdn_frontdoor_rule"]; len(got) != 3 || got[0] != StylePascal || got[1] != StyleCamel || got[2] != StyleStraight {
		t.Fatalf("expected CDN Front Door rule styles to be [%q %q %q], got %#v", StylePascal, StyleCamel, StyleStraight, got)
	}
}

func TestBuildNameAzureStorageAccountFallsBackFromDashedToStraightAndLowercase(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudAzure)
	if err != nil {
		t.Fatalf("unexpected error loading Azure defaults: %v", err)
	}

	result, err := BuildName(Config{
		Cloud:                            CloudAzure,
		StylePriority:                    []string{StyleDashed, StylePascal, StyleCamel},
		ResourceAcronyms:                 defaults.ResourceAcronyms,
		ResourceStyleOverrides:           defaults.ResourceStyleOverrides,
		ResourceConstraints:              defaults.ResourceConstraints,
		RegionalResources:                defaults.RegionalResources,
		IgnoreRegionForRegionalResources: false,
	}, BuildInput{
		Resource:  "azurerm_storage_account",
		Qualifier: "data-lake",
		Recipe:    []string{"resource", "qualifier"},
	})
	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}

	if result.Style != StyleStraight {
		t.Fatalf("expected fallback style %q, got %q", StyleStraight, result.Style)
	}
	if result.Name != "stdatalake" {
		t.Fatalf("expected generated name %q, got %q", "stdatalake", result.Name)
	}
}

func TestBuildNameAzureCdnFrontdoorRuleFallsBackFromDashedToPascal(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudAzure)
	if err != nil {
		t.Fatalf("unexpected error loading Azure defaults: %v", err)
	}

	result, err := BuildName(Config{
		Cloud:                            CloudAzure,
		StylePriority:                    []string{StyleDashed, StylePascal, StyleCamel},
		ResourceAcronyms:                 defaults.ResourceAcronyms,
		ResourceStyleOverrides:           defaults.ResourceStyleOverrides,
		ResourceConstraints:              defaults.ResourceConstraints,
		RegionalResources:                defaults.RegionalResources,
		IgnoreRegionForRegionalResources: false,
	}, BuildInput{
		Resource:  "azurerm_cdn_frontdoor_rule",
		Qualifier: "edge-prod",
		Recipe:    []string{"resource", "qualifier"},
	})
	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}

	if result.Style != StylePascal {
		t.Fatalf("expected fallback style %q, got %q", StylePascal, result.Style)
	}
	if result.Name != "CfdrEdgeProd" {
		t.Fatalf("expected generated name %q, got %q", "CfdrEdgeProd", result.Name)
	}
}

func TestDefaultCloudDefaultsGCP(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudGCP)
	if err != nil {
		t.Fatalf("unexpected error loading GCP defaults: %v", err)
	}

	if len(defaults.RegionMap) == 0 {
		t.Fatal("expected GCP region map defaults to be populated")
	}

	if got := defaults.RegionMap["us-central1"]; got != "usc1" {
		t.Fatalf("expected GCP region short code %q for us-central1, got %q", "usc1", got)
	}
	if got := defaults.RegionMap["northamerica-south1"]; got != "nas1" {
		t.Fatalf("expected GCP region short code %q for northamerica-south1, got %q", "nas1", got)
	}
	if got := defaults.RegionMap["europe-north2"]; got != "eun2" {
		t.Fatalf("expected GCP region short code %q for europe-north2, got %q", "eun2", got)
	}
	if got := defaults.RegionMap["asia-southeast3"]; got != "asse3" {
		t.Fatalf("expected GCP region short code %q for asia-southeast3, got %q", "asse3", got)
	}

	if got := defaults.ResourceAcronyms["storage_bucket"]; got != "gcs" {
		t.Fatalf("expected GCP bucket acronym %q, got %q", "gcs", got)
	}
	if got := defaults.ResourceAcronyms["compute_network"]; got != "vpc" {
		t.Fatalf("expected GCP network acronym %q, got %q", "vpc", got)
	}
	if got := defaults.ResourceAcronyms["compute_subnetwork"]; got != "snet" {
		t.Fatalf("expected GCP subnetwork acronym %q, got %q", "snet", got)
	}
	if got := defaults.ResourceAcronyms["compute_router"]; got != "crtr" {
		t.Fatalf("expected GCP router acronym %q, got %q", "crtr", got)
	}
	if got := defaults.ResourceAcronyms["compute_firewall"]; got != "cfwl" {
		t.Fatalf("expected GCP firewall acronym %q, got %q", "cfwl", got)
	}
	if got := defaults.ResourceAcronyms["compute_global_address"]; got != "gaddr" {
		t.Fatalf("expected GCP global address acronym %q, got %q", "gaddr", got)
	}
	if got := defaults.ResourceAcronyms["compute_target_https_proxy"]; got != "cthps" {
		t.Fatalf("expected GCP target HTTPS proxy acronym %q, got %q", "cthps", got)
	}
	if got := defaults.ResourceAcronyms["compute_region_backend_service"]; got != "crbs" {
		t.Fatalf("expected GCP regional backend service acronym %q, got %q", "crbs", got)
	}
	if got := defaults.ResourceAcronyms["dns_managed_zone"]; got != "dnsz" {
		t.Fatalf("expected GCP DNS managed zone acronym %q, got %q", "dnsz", got)
	}
	if got := defaults.ResourceAcronyms["sql_database_instance"]; got != "sqli" {
		t.Fatalf("expected GCP SQL instance acronym %q, got %q", "sqli", got)
	}
	if got := defaults.ResourceAcronyms["container_cluster"]; got != "gkec" {
		t.Fatalf("expected GCP GKE cluster acronym %q, got %q", "gkec", got)
	}
	if got := defaults.ResourceAcronyms["gke_cluster"]; got != "gkec" {
		t.Fatalf("expected GCP GKE cluster alias acronym %q, got %q", "gkec", got)
	}
	if got := defaults.ResourceAcronyms["gke_node_pool"]; got != "gkenp" {
		t.Fatalf("expected GCP GKE node pool alias acronym %q, got %q", "gkenp", got)
	}
	if got := defaults.ResourceAcronyms["workflows_workflow"]; got != "wflw" {
		t.Fatalf("expected GCP workflow acronym %q, got %q", "wflw", got)
	}
	if _, ok := defaults.ResourceAcronyms["google_compute_network"]; ok {
		t.Fatal("expected GCP defaults to normalize rather than store google_* resource keys")
	}

	styles := defaults.ResourceStyleOverrides["storage_bucket"]
	if len(styles) == 0 {
		t.Fatal("expected storage bucket style overrides to be populated")
	}
	if !containsString(styles, StyleDashed) || !containsString(styles, StyleUnderscore) || !containsString(styles, StyleStraight) {
		t.Fatalf("expected bucket styles to include dashed, underscore, and straight; got %#v", styles)
	}

	if !defaults.RegionalResources["compute_subnetwork"] {
		t.Fatal("expected compute_subnetwork to be marked regional")
	}
	if defaults.RegionalResources["storage_bucket"] {
		t.Fatal("expected storage_bucket to be marked non-regional")
	}
	if defaults.RegionalResources["compute_firewall"] {
		t.Fatal("expected compute_firewall to be marked non-regional")
	}
	if defaults.RegionalResources["compute_global_address"] {
		t.Fatal("expected compute_global_address to be marked non-regional")
	}
	if !defaults.RegionalResources["compute_router_nat"] {
		t.Fatal("expected compute_router_nat to be marked regional")
	}
}

func TestBuildNameGCPBucketFallsBackToAllowedStyle(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudGCP)
	if err != nil {
		t.Fatalf("unexpected error loading GCP defaults: %v", err)
	}

	result, err := BuildName(Config{
		Cloud:                  CloudGCP,
		OrgPrefix:              "acme",
		Project:                "payments",
		Env:                    "prod",
		StylePriority:          []string{StylePascal},
		ResourceAcronyms:       defaults.ResourceAcronyms,
		ResourceStyleOverrides: defaults.ResourceStyleOverrides,
		ResourceConstraints:    defaults.ResourceConstraints,
		RegionalResources:      defaults.RegionalResources,
	}, BuildInput{
		Resource:  "google_storage_bucket",
		Qualifier: "raw",
		Recipe:    []string{"org", "proj", "env", "resource", "qualifier"},
	})
	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}

	if result.Style != StyleDashed {
		t.Fatalf("expected fallback style %q, got %q", StyleDashed, result.Style)
	}
}

func TestBuildNameGCPBucketRejectsReservedGoogleSubstring(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudGCP)
	if err != nil {
		t.Fatalf("unexpected error loading GCP defaults: %v", err)
	}

	_, err = BuildName(Config{
		Cloud:                  CloudGCP,
		OrgPrefix:              "acme",
		ResourceAcronyms:       defaults.ResourceAcronyms,
		ResourceStyleOverrides: defaults.ResourceStyleOverrides,
		ResourceConstraints:    defaults.ResourceConstraints,
		RegionalResources:      defaults.RegionalResources,
	}, BuildInput{
		Resource:  "google_storage_bucket",
		Qualifier: "google-data",
		Recipe:    []string{"org", "qualifier"},
	})
	if err == nil {
		t.Fatal("expected bucket constraint error, got nil")
	}
}

func TestBuildNameGCPBucketRejectsGoogleLookalike(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudGCP)
	if err != nil {
		t.Fatalf("unexpected error loading GCP defaults: %v", err)
	}

	_, err = BuildName(Config{
		Cloud:                  CloudGCP,
		OrgPrefix:              "acme",
		ResourceAcronyms:       defaults.ResourceAcronyms,
		ResourceStyleOverrides: defaults.ResourceStyleOverrides,
		ResourceConstraints:    defaults.ResourceConstraints,
		RegionalResources:      defaults.RegionalResources,
	}, BuildInput{
		Resource:  "google_storage_bucket",
		Qualifier: "g00gle-data",
		Recipe:    []string{"org", "qualifier"},
	})
	if err == nil {
		t.Fatal("expected bucket lookalike constraint error, got nil")
	}
}

func TestBuildNameGCPPubSubTopicRejectsGoogPrefix(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudGCP)
	if err != nil {
		t.Fatalf("unexpected error loading GCP defaults: %v", err)
	}

	_, err = BuildName(Config{
		Cloud:                  CloudGCP,
		ResourceAcronyms:       defaults.ResourceAcronyms,
		ResourceStyleOverrides: defaults.ResourceStyleOverrides,
		ResourceConstraints:    defaults.ResourceConstraints,
		RegionalResources:      defaults.RegionalResources,
	}, BuildInput{
		Resource:  "google_pubsub_topic",
		Qualifier: "goog-events",
		Recipe:    []string{"qualifier"},
	})
	if err == nil {
		t.Fatal("expected pubsub constraint error, got nil")
	}
}

func TestBuildNameGCPServiceAccountRejectsShortNames(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudGCP)
	if err != nil {
		t.Fatalf("unexpected error loading GCP defaults: %v", err)
	}

	_, err = BuildName(Config{
		Cloud:                  CloudGCP,
		ResourceAcronyms:       defaults.ResourceAcronyms,
		ResourceStyleOverrides: defaults.ResourceStyleOverrides,
		ResourceConstraints:    defaults.ResourceConstraints,
		RegionalResources:      defaults.RegionalResources,
	}, BuildInput{
		Resource:  "google_service_account",
		Qualifier: "svc1",
		Recipe:    []string{"qualifier"},
	})
	if err == nil {
		t.Fatal("expected service account constraint error, got nil")
	}
}

func TestBuildNameGCPBigQueryDatasetFallsBackToStraightStyle(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudGCP)
	if err != nil {
		t.Fatalf("unexpected error loading GCP defaults: %v", err)
	}

	result, err := BuildName(Config{
		Cloud:                  CloudGCP,
		OrgPrefix:              "acme",
		Project:                "analytics",
		StylePriority:          []string{StyleDashed},
		ResourceAcronyms:       defaults.ResourceAcronyms,
		ResourceStyleOverrides: defaults.ResourceStyleOverrides,
		ResourceConstraints:    defaults.ResourceConstraints,
		RegionalResources:      defaults.RegionalResources,
	}, BuildInput{
		Resource: "google_bigquery_dataset",
		Recipe:   []string{"org", "proj"},
	})
	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}

	if result.Style != StyleStraight {
		t.Fatalf("expected fallback style %q, got %q", StyleStraight, result.Style)
	}
	if result.Name != "acmeanalytics" {
		t.Fatalf("expected generated name %q, got %q", "acmeanalytics", result.Name)
	}
}

func TestBuildNameGCPCloudRunRejectsLeadingDigit(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudGCP)
	if err != nil {
		t.Fatalf("unexpected error loading GCP defaults: %v", err)
	}

	_, err = BuildName(Config{
		Cloud:                  CloudGCP,
		ResourceAcronyms:       defaults.ResourceAcronyms,
		ResourceStyleOverrides: defaults.ResourceStyleOverrides,
		ResourceConstraints:    defaults.ResourceConstraints,
		RegionalResources:      defaults.RegionalResources,
	}, BuildInput{
		Resource:  "google_cloud_run_v2_service",
		Qualifier: "9api",
		Recipe:    []string{"qualifier"},
	})
	if err == nil {
		t.Fatal("expected cloud run constraint error, got nil")
	}
}

func TestBuildNameGCPComputeRouterRejectsUnderscoreStyle(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudGCP)
	if err != nil {
		t.Fatalf("unexpected error loading GCP defaults: %v", err)
	}

	_, err = BuildName(Config{
		Cloud:                  CloudGCP,
		OrgPrefix:              "acme",
		Env:                    "dev",
		ResourceAcronyms:       defaults.ResourceAcronyms,
		ResourceStyleOverrides: defaults.ResourceStyleOverrides,
		ResourceConstraints:    defaults.ResourceConstraints,
		RegionalResources:      defaults.RegionalResources,
	}, BuildInput{
		Resource:      "google_compute_router",
		Qualifier:     "edge",
		Recipe:        []string{"org", "env", "resource", "qualifier"},
		StylePriority: []string{StyleUnderscore},
	})
	if err == nil {
		t.Fatal("expected compute router constraint error, got nil")
	}
}

func TestDefaultGCPPrimaryResourceAcronymsAreUnique(t *testing.T) {
	acronyms := DefaultGCPResourceAcronyms()
	owners := map[string]string{}
	aliases := map[string]bool{
		"gcs_bucket":        true,
		"gcs":               true,
		"vpc":               true,
		"subnet":            true,
		"cloud_run_service": true,
		"sql_instance":      true,
		"gke_cluster":       true,
		"gke_node_pool":     true,
	}

	for resource, acronym := range acronyms {
		if aliases[resource] {
			continue
		}
		if previous, ok := owners[acronym]; ok {
			t.Fatalf("duplicate GCP acronym %q for %q and %q", acronym, previous, resource)
		}
		owners[acronym] = resource
	}
}

func TestDefaultGCPPrimaryResourceAcronymsStayCompact(t *testing.T) {
	acronyms := DefaultGCPResourceAcronyms()
	aliases := map[string]bool{
		"gcs_bucket":        true,
		"gcs":               true,
		"vpc":               true,
		"subnet":            true,
		"cloud_run_service": true,
		"sql_instance":      true,
		"gke_cluster":       true,
		"gke_node_pool":     true,
	}

	for resource, acronym := range acronyms {
		if aliases[resource] {
			continue
		}
		if len(acronym) < 3 || len(acronym) > 6 {
			t.Fatalf("expected compact GCP acronym for %q, got %q", resource, acronym)
		}
	}
}
