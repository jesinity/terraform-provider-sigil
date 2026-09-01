package naming

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestDatabricksPlatformDefaultsAndAliases(t *testing.T) {
	defaults, err := DefaultDefaults(CloudAWS, PlatformDatabricks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if defaults.ResourceAcronyms["databricks_cluster"] != "dbc" {
		t.Fatalf("expected Databricks cluster acronym, got %q", defaults.ResourceAcronyms["databricks_cluster"])
	}
	if _, ok := defaults.ResourceAcronyms["s3_bucket"]; !ok {
		t.Fatal("expected cloud defaults to remain available")
	}

	canonical, err := BuildName(Config{Cloud: CloudAWS, Platform: PlatformDatabricks, OrgPrefix: "acme", Env: "prod"}, BuildInput{Resource: "databricks_cluster", Qualifier: "etl"})
	if err != nil {
		t.Fatal(err)
	}
	alias, err := BuildName(Config{Cloud: CloudAWS, Platform: PlatformDatabricks, OrgPrefix: "acme", Env: "prod"}, BuildInput{Resource: "cluster", Qualifier: "etl"})
	if err != nil {
		t.Fatal(err)
	}
	if canonical.ResourceAcronym != "dbc" || alias.ResourceAcronym != canonical.ResourceAcronym {
		t.Fatalf("canonical and alias should resolve to dbc: %#v %#v", canonical, alias)
	}

	withoutPlatform, err := BuildName(Config{Cloud: CloudAWS, OrgPrefix: "acme", Env: "prod"}, BuildInput{Resource: "cluster", Qualifier: "etl"})
	if err != nil {
		t.Fatal(err)
	}
	if withoutPlatform.ResourceAcronym == "dbc" {
		t.Fatal("Databricks alias must not resolve without the platform")
	}
}

func TestDatabricksOverridesAzureAndDefaultsAreIsolated(t *testing.T) {
	azure, err := DefaultCloudDefaults(CloudAzure)
	if err != nil {
		t.Fatal(err)
	}
	overlay, err := DefaultDefaults(CloudAzure, PlatformDatabricks)
	if err != nil {
		t.Fatal(err)
	}
	if overlay.ResourceAcronyms["databricks_cluster"] != "dbc" {
		t.Fatalf("platform must override Azure Databricks defaults, got %q", overlay.ResourceAcronyms["databricks_cluster"])
	}
	overlay.ResourceAcronyms["databricks_cluster"] = "changed"
	overlay.ResourceStyleOverrides["databricks_cluster"][0] = "changed"
	again, err := DefaultDefaults(CloudAzure, PlatformDatabricks)
	if err != nil {
		t.Fatal(err)
	}
	if again.ResourceAcronyms["databricks_cluster"] != "dbc" || again.ResourceStyleOverrides["databricks_cluster"][0] != StyleDashed {
		t.Fatal("default maps or slices leaked between provider instances")
	}
	if len(azure.ResourceStyleOverrides["databricks_cluster"]) < 2 || azure.ResourceStyleOverrides["databricks_cluster"][1] != StylePascalDashed {
		t.Fatal("Azure-only Databricks defaults must remain unchanged")
	}
}

func TestDatabricksManifestValidation(t *testing.T) {
	invalid := []byte(`{"provider":"databricks/databricks","provider_version":"1.129.0","resources":[{"name":"databricks_one","acronym":"same","naming_attribute":"name","scope":"workspace","styles":["dashed"],"support_status":"supported","constraint_status":"audited","documentation_url":"https://example.test"},{"name":"databricks_two","acronym":"same","naming_attribute":"name","scope":"workspace","styles":["dashed"],"support_status":"supported","constraint_status":"audited","documentation_url":"https://example.test"}]}`)
	if _, _, err := loadDatabricksPlatformDefaults(invalid); err == nil {
		t.Fatal("expected duplicate acronym error")
	}
}

func TestDatabricksSupportedResources(t *testing.T) {
	var manifest databricksManifest
	if err := json.Unmarshal(databricksResourceDefinitionJSON, &manifest); err != nil {
		t.Fatalf("decode embedded manifest: %v", err)
	}

	for _, resource := range manifest.Resources {
		if resource.SupportStatus != "supported" {
			continue
		}
		t.Run(resource.Name, func(t *testing.T) {
			for _, cloud := range []string{CloudAWS, CloudAzure, CloudGCP} {
				if len(resource.SupportedClouds) > 0 && !containsString(resource.SupportedClouds, cloud) {
					continue
				}
				result, err := BuildName(Config{
					Cloud:                            cloud,
					Platform:                         PlatformDatabricks,
					OrgPrefix:                        "acme",
					Project:                          "lake",
					Env:                              "prod",
					RegionShortCode:                  "eu1",
					IgnoreRegionForRegionalResources: true,
				}, BuildInput{Resource: resource.Name, Qualifier: "test", StylePriority: resource.Styles})
				if err != nil {
					t.Fatalf("%s representative name failed: %v", cloud, err)
				}
				if result.ResourceAcronym != resource.Acronym {
					t.Fatalf("%s: expected acronym %q, got %q", cloud, resource.Acronym, result.ResourceAcronym)
				}
				if result.Style != resource.Styles[0] {
					t.Fatalf("%s: expected preferred style %q, got %q", cloud, resource.Styles[0], result.Style)
				}
				if resource.Regional && strings.Contains(result.Name, "eu1") {
					t.Fatalf("%s: regional resource name should omit region: %q", cloud, result.Name)
				}
				if !resource.Regional && !strings.Contains(result.Name, "eu1") {
					t.Fatalf("%s: global resource name should include region: %q", cloud, result.Name)
				}
			}

			for _, alias := range resource.Aliases {
				aliased, err := BuildName(Config{Cloud: CloudAWS, Platform: PlatformDatabricks, OrgPrefix: "acme", Env: "prod"}, BuildInput{Resource: alias})
				if err != nil || aliased.ResourceAcronym != resource.Acronym {
					t.Fatalf("alias %q did not resolve to %q: %#v, %v", alias, resource.Acronym, aliased, err)
				}
			}
		})
	}
}

func TestDatabricksProvisioningCloudCompatibility(t *testing.T) {
	for _, test := range []struct {
		cloud, resource string
		wantError       bool
	}{
		{CloudAWS, "databricks_mws_workspaces", false}, {CloudGCP, "databricks_mws_workspaces", false},
		{CloudAzure, "azurerm_databricks_workspace", false}, {CloudAzure, "databricks_mws_workspaces", true},
		{CloudGCP, "databricks_mws_credentials", true}, {CloudAWS, "azurerm_databricks_workspace", true},
	} {
		_, err := BuildName(Config{Cloud: test.cloud, Platform: PlatformDatabricks, OrgPrefix: "acme", Env: "prod"}, BuildInput{Resource: test.resource})
		if (err != nil) != test.wantError {
			t.Errorf("%s on %s: err=%v, wantError=%t", test.resource, test.cloud, err, test.wantError)
		}
	}
}

func TestDatabricksProvisioningRegionMaps(t *testing.T) {
	for _, test := range []struct{ cloud, region, want string }{
		{CloudAWS, "us-east-1", "use1"}, {CloudAzure, "westeurope", "weu"}, {CloudGCP, "us-central1", "usc1"},
	} {
		resource := "databricks_mws_workspaces"
		if test.cloud == CloudAzure {
			resource = "azurerm_databricks_workspace"
		}
		result, err := BuildName(Config{Cloud: test.cloud, Platform: PlatformDatabricks, OrgPrefix: "acme", Env: "prod", Region: test.region, IgnoreRegionForRegionalResources: false}, BuildInput{Resource: resource})
		if err != nil {
			t.Fatalf("%s: %v", test.cloud, err)
		}
		if result.RegionCode != test.want || !strings.Contains(result.Name, test.want) {
			t.Errorf("%s: expected region code %q in %q", test.cloud, test.want, result.Name)
		}
	}
}

func TestDatabricksAppConstraint(t *testing.T) {
	_, err := BuildName(Config{Cloud: CloudAWS, Platform: PlatformDatabricks, OrgPrefix: "9invalid", Env: "prod"}, BuildInput{Resource: "databricks_app", Qualifier: "app"})
	if err == nil {
		t.Fatal("expected lowercase dashed Databricks App constraint error")
	}
}

func TestDatabricksDocumentationCoversManifest(t *testing.T) {
	contents, err := os.ReadFile("../../docs/databricks-resources.md")
	if err != nil {
		t.Fatal(err)
	}
	var manifest databricksManifest
	if err := json.Unmarshal(databricksResourceDefinitionJSON, &manifest); err != nil {
		t.Fatal(err)
	}
	for _, resource := range manifest.Resources {
		if resource.SupportStatus == "supported" && !strings.Contains(string(contents), "`"+resource.Name+"`") {
			t.Fatalf("documentation is missing %q", resource.Name)
		}
	}
}

func TestDatabricksCatalogClassifiesV11290Resources(t *testing.T) {
	contents, err := os.ReadFile("databricks_resource_catalog.json")
	if err != nil {
		t.Fatal(err)
	}
	var catalog struct {
		Provider        string `json:"provider"`
		ProviderVersion string `json:"provider_version"`
		Resources       []struct {
			Name             string `json:"name"`
			Category         string `json:"category"`
			SupportStatus    string `json:"support_status"`
			Reason           string `json:"reason"`
			DocumentationURL string `json:"documentation_url"`
		} `json:"resources"`
	}
	if err := json.Unmarshal(contents, &catalog); err != nil {
		t.Fatal(err)
	}
	if catalog.Provider != "databricks/databricks" || catalog.ProviderVersion != "1.129.0" {
		t.Fatalf("unexpected catalog source: %#v", catalog)
	}
	if len(catalog.Resources) != 170 {
		t.Fatalf("expected 170 v1.129.0 resources, got %d", len(catalog.Resources))
	}
	seen := map[string]bool{}
	for _, resource := range catalog.Resources {
		if seen[resource.Name] || resource.Category == "" || resource.SupportStatus == "" || resource.Reason == "" || resource.DocumentationURL == "" {
			t.Fatalf("incomplete or duplicate catalog entry: %#v", resource)
		}
		seen[resource.Name] = true
	}
}

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

func TestDefaultCloudDefaultsAWSDataScienceAndEngineering(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudAWS)
	if err != nil {
		t.Fatalf("unexpected error loading AWS defaults: %v", err)
	}

	expected := map[string]string{
		"athena_workgroup":                 "athwg",
		"glue_catalog_database":            "glcdb",
		"glue_job":                         "gljob",
		"sagemaker_domain":                 "sgdmn",
		"sagemaker_feature_group":          "sgfgr",
		"bedrock_guardrail":                "brgr",
		"bedrockagent_agent":               "braga",
		"bedrockagent_knowledge_base":      "brakb",
		"bedrockagentcore_agent_runtime":   "bracr",
		"datazone_domain":                  "dzdmn",
		"emrserverless_application":        "emrsa",
		"kinesis_firehose_delivery_stream": "knsfh",
		"lakeformation_lf_tag":             "lftag",
		"quicksight_dashboard":             "qsdsh",
		"redshiftserverless_workgroup":     "rsswg",
	}

	for resource, acronym := range expected {
		if got := defaults.ResourceAcronyms[resource]; got != acronym {
			t.Fatalf("expected AWS resource %q acronym %q, got %q", resource, acronym, got)
		}
		if !defaults.RegionalResources[resource] {
			t.Fatalf("expected AWS resource %q to be regional", resource)
		}
	}
}

func TestDefaultAWSLegacyResourceAcronymsRemainStable(t *testing.T) {
	expected := map[string]string{
		"role": "role", "role_policy": "rlpl", "iam_role": "role", "iam_policy": "iamp",
		"iam_user": "iamu", "iam_group": "iamg", "s3": "s3b", "s3_bucket": "s3bk",
		"s3_object": "s3ob", "s3_access_point": "s3ap", "s3_table": "s3tb", "s3_dir": "s3dr",
		"sns": "sns", "sqs": "sqs", "ecs_cluster": "ecsc", "ecs_service": "ecss",
		"ecs_task": "ecst", "eks": "eks", "eks_cluster": "eksc", "eks_node_group": "ekng",
		"msk_cluster": "mskc", "vpc": "vpcn", "subnet": "subn", "igw": "igtw",
		"nat_gw": "ngtw", "sec_group": "scgp", "nacl": "nacl", "route_table": "rttb",
		"elastic_ip": "elip", "wafv2_web_acl": "wfac", "wafv2_web_acl_rule": "wfar", "wafv2_ip_set": "wfis",
		"lambda": "lmbd", "api_gateway_rest_api": "agra", "api_gateway_model": "agmd", "api_gateway_v2": "agv2",
		"log_group": "logg", "cloudwatch_log_group": "cwlg", "cloudwatch_alarm": "cwal", "eventbridge_bus": "evbb",
		"eventbridge_rule": "evbr", "step_function": "stfn", "sfn": "stfn", "dynamodb": "dydb",
		"dynamodb_table": "dybt", "rds": "rds", "rds_cluster": "rdsc", "aurora_cluster": "arcl",
		"redshift": "rdsh", "elasticache": "elch", "opensearch": "opsr", "elasticsearch": "elsr",
		"ecr": "ecr", "ecs": "ecs", "ec2_instance": "ec2i", "launch_template": "lcht",
		"autoscaling_group": "asgr", "alb": "albl", "nlb": "nlbl", "elb": "elbl",
		"target_group": "tgpt", "cloudfront": "clfr", "route53_zone": "rt53", "route53_record": "r53r",
		"acm_cert": "acmc", "kms_key": "kmsk", "secretsmanager_secret": "smse", "ssm_parameter": "ssmp",
		"cloudtrail": "ctra", "guardduty": "gdty", "config_rule": "cfrl", "efs": "efs",
		"ebs": "ebs", "athena": "athn", "glue": "glue", "sagemaker": "sgmk",
		"codebuild": "cdbd", "codepipeline": "cdpl", "codedeploy": "cddp", "cloudformation_stack": "cfst",
		"appsync": "apsy", "snow_notification_integration": "snti",
	}

	actual := DefaultResourceAcronyms()
	if len(expected) != 82 {
		t.Fatalf("legacy test fixture must contain 82 entries, got %d", len(expected))
	}
	for resource, acronym := range expected {
		if got := actual[resource]; got != acronym {
			t.Fatalf("legacy AWS resource %q changed from %q to %q", resource, acronym, got)
		}
	}
}

func TestDefaultAWSOpaqueResourcesAreNotNameable(t *testing.T) {
	acronyms := DefaultResourceAcronyms()
	opaque := DefaultAWSOpaqueResources()
	if len(opaque) != 48 {
		t.Fatalf("expected 48 explicitly classified opaque AWS resources, got %d", len(opaque))
	}
	for resource := range opaque {
		if acronym, exists := acronyms[resource]; exists {
			t.Fatalf("opaque AWS resource %q must not expose acronym %q", resource, acronym)
		}
	}
}

func TestDefaultAWSCanonicalTerraformAliasesResolve(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudAWS)
	if err != nil {
		t.Fatalf("unexpected error loading AWS defaults: %v", err)
	}

	for canonical, legacy := range defaultAWSResourceAliases {
		expected := defaults.ResourceAcronyms[legacy]
		result, err := BuildName(Config{
			Cloud:                            CloudAWS,
			OrgPrefix:                        "acme",
			Project:                          "core",
			Env:                              "dev",
			Region:                           "us-east-1",
			RegionMap:                        defaults.RegionMap,
			IgnoreRegionForRegionalResources: false,
			ResourceAcronyms:                 defaults.ResourceAcronyms,
			ResourceStyleOverrides:           defaults.ResourceStyleOverrides,
			ResourceConstraints:              defaults.ResourceConstraints,
			RegionalResources:                defaults.RegionalResources,
		}, BuildInput{Resource: "aws_" + canonical})
		if err != nil {
			t.Fatalf("canonical AWS resource %q failed through legacy alias %q: %v", canonical, legacy, err)
		}
		if result.ResourceAcronym != expected {
			t.Fatalf("canonical AWS resource %q resolved to %q, expected legacy acronym %q", canonical, result.ResourceAcronym, expected)
		}
	}
}

func TestDefaultAWSDataResourcesBuildWithinConstraints(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudAWS)
	if err != nil {
		t.Fatalf("unexpected error loading AWS defaults: %v", err)
	}

	checked := 0
	for resource := range defaults.ResourceAcronyms {
		if _, isDataResource := defaultAWSDataResourceConstraint(resource); !isDataResource {
			continue
		}
		checked++
		if _, exists := defaults.ResourceConstraints[resource]; !exists {
			t.Fatalf("AWS data resource %q has no enforced naming constraint", resource)
		}
		if _, err := BuildName(Config{
			Cloud:                            CloudAWS,
			OrgPrefix:                        "acme",
			Project:                          "ml",
			Env:                              "dev",
			Region:                           "us-east-1",
			RegionMap:                        defaults.RegionMap,
			IgnoreRegionForRegionalResources: false,
			ResourceAcronyms:                 defaults.ResourceAcronyms,
			ResourceStyleOverrides:           defaults.ResourceStyleOverrides,
			ResourceConstraints:              defaults.ResourceConstraints,
			RegionalResources:                defaults.RegionalResources,
		}, BuildInput{Resource: "aws_" + resource}); err != nil {
			t.Fatalf("AWS data resource %q generated an invalid default name: %v", resource, err)
		}
	}
	if checked != 153 {
		t.Fatalf("expected 153 constrained AWS data resources, checked %d", checked)
	}
}

func TestBuildNameAWSStripsTerraformResourcePrefix(t *testing.T) {
	defaults, err := DefaultCloudDefaults(CloudAWS)
	if err != nil {
		t.Fatalf("unexpected error loading AWS defaults: %v", err)
	}

	result, err := BuildName(Config{
		Cloud:                            CloudAWS,
		OrgPrefix:                        "acme",
		Project:                          "ml",
		Env:                              "prod",
		Region:                           "us-east-1",
		RegionMap:                        defaults.RegionMap,
		IgnoreRegionForRegionalResources: false,
		ResourceAcronyms:                 defaults.ResourceAcronyms,
		ResourceStyleOverrides:           defaults.ResourceStyleOverrides,
		ResourceConstraints:              defaults.ResourceConstraints,
		RegionalResources:                defaults.RegionalResources,
	}, BuildInput{
		Resource:  "aws_bedrockagent_knowledge_base",
		Qualifier: "retrieval",
	})
	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}

	if result.ResourceAcronym != "brakb" {
		t.Fatalf("expected resource acronym %q, got %q", "brakb", result.ResourceAcronym)
	}
	if result.Name != "acme-ml-prod-use1-brakb-retrieval" {
		t.Fatalf("expected generated name %q, got %q", "acme-ml-prod-use1-brakb-retrieval", result.Name)
	}
}

func TestDefaultAWSPrimaryResourceAcronymsAreUnique(t *testing.T) {
	acronyms := DefaultResourceAcronyms()
	owners := map[string]string{}

	for resource, acronym := range acronyms {
		if previous, ok := owners[acronym]; ok {
			if !isAWSAllowedAcronymAlias(acronym, previous, resource) {
				t.Fatalf("duplicate AWS acronym %q for %q and %q", acronym, previous, resource)
			}
		}
		owners[acronym] = resource
	}
}

func isAWSAllowedAcronymAlias(acronym, a, b string) bool {
	if a > b {
		a, b = b, a
	}
	switch acronym {
	case "role":
		return a == "iam_role" && b == "role"
	case "stfn":
		return a == "sfn" && b == "step_function"
	default:
		return false
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
