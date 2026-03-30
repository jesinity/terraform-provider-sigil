package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sigil": providerserver.NewProtocol6WithError(New("test")()),
}

func TestMarkDataSource_gcpRegionalResourceIncludesRegion(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud                                = "gcp"
  org_prefix                           = "acme"
  env                                  = "dev"
  region                               = "northamerica-south1"
  ignore_region_for_regional_resources = false
`, `
data "sigil_mark" "subnet" {
  what      = "google_compute_subnetwork"
  qualifier = "app"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "name", "acme-dev-nas1-snet-app"),
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "region_code", "nas1"),
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "style", "dashed"),
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "resource_acronym", "snet"),
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "components.region", "nas1"),
				),
			},
		},
	})
}

func TestMarkDataSource_gcpRegionalResourceOmitsRegionByDefault(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud      = "gcp"
  org_prefix = "acme"
  env        = "dev"
  region     = "northamerica-south1"
`, `
data "sigil_mark" "subnet" {
  what      = "google_compute_subnetwork"
  qualifier = "app"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "name", "acme-dev-snet-app"),
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "region_code", ""),
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "components.region", ""),
				),
			},
		},
	})
}

func TestMarkDataSource_gcpRegionOverrideWins(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud      = "gcp"
  org_prefix = "acme"
  env        = "dev"
  region     = "northamerica-south1"
  region_overrides = {
    "northamerica-south1" = "mxs1"
  }
  ignore_region_for_regional_resources = false
`, `
data "sigil_mark" "subnet" {
  what      = "google_compute_subnetwork"
  qualifier = "app"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "name", "acme-dev-mxs1-snet-app"),
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "region_code", "mxs1"),
				),
			},
		},
	})
}

func TestMarkDataSource_conflictingWhatAndResource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud      = "gcp"
  org_prefix = "acme"
  env        = "dev"
`, `
data "sigil_mark" "conflict" {
  what     = "google_compute_network"
  resource = "google_storage_bucket"
}
`),
				ExpectError: regexp.MustCompile("cannot both be set to different values"),
			},
		},
	})
}

func TestMarkDataSource_bucketConstraintFailure(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud      = "gcp"
  org_prefix = "acme"
  env        = "dev"
`, `
data "sigil_mark" "bucket" {
  what      = "google_storage_bucket"
  qualifier = "google-data"
}
`),
				ExpectError: regexp.MustCompile(`must not\s+contain "google"`),
			},
		},
	})
}

func TestProvider_missingEnv(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud      = "gcp"
  org_prefix = "acme"
`, `
data "sigil_mark" "network" {
  what = "google_compute_network"
}
`),
				ExpectError: regexp.MustCompile("Missing env"),
			},
		},
	})
}

func TestMarkDataSource_awsRegionalResourceIncludesRegion(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud                                = "aws"
  org_prefix                           = "acme"
  env                                  = "dev"
  region                               = "ap-southeast-2"
  ignore_region_for_regional_resources = false
`, `
data "sigil_mark" "vpc" {
  what      = "vpc"
  qualifier = "core"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.vpc", "name", "acme-dev-apse2-vpcn-core"),
					resource.TestCheckResourceAttr("data.sigil_mark.vpc", "region_code", "apse2"),
					resource.TestCheckResourceAttr("data.sigil_mark.vpc", "style", "dashed"),
					resource.TestCheckResourceAttr("data.sigil_mark.vpc", "resource_acronym", "vpcn"),
					resource.TestCheckResourceAttr("data.sigil_mark.vpc", "components.region", "apse2"),
				),
			},
		},
	})
}

func TestMarkDataSource_awsRegionalResourceOmitsRegionByDefault(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud      = "aws"
  org_prefix = "acme"
  env        = "dev"
  region     = "ap-southeast-2"
`, `
data "sigil_mark" "subnet" {
  what      = "subnet"
  qualifier = "app"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "name", "acme-dev-subn-app"),
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "region_code", ""),
					resource.TestCheckResourceAttr("data.sigil_mark.subnet", "components.region", ""),
				),
			},
		},
	})
}

func TestMarkDataSource_awsGlobalResourceKeepsRegionByDefault(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud      = "aws"
  org_prefix = "acme"
  env        = "dev"
  region     = "us-east-1"
`, `
data "sigil_mark" "role" {
  what      = "iam_role"
  qualifier = "app"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.role", "name", "acme-dev-use1-role-app"),
					resource.TestCheckResourceAttr("data.sigil_mark.role", "region_code", "use1"),
					resource.TestCheckResourceAttr("data.sigil_mark.role", "resource_acronym", "role"),
					resource.TestCheckResourceAttr("data.sigil_mark.role", "components.region", "use1"),
				),
			},
		},
	})
}

func TestMarkDataSource_awsRegionOverrideWins(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud                                = "aws"
  org_prefix                           = "acme"
  env                                  = "dev"
  region                               = "us-east-1"
  ignore_region_for_regional_resources = false
  region_overrides = {
    "us-east-1" = "usex"
  }
`, `
data "sigil_mark" "vpc" {
  what      = "vpc"
  qualifier = "core"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.vpc", "name", "acme-dev-usex-vpcn-core"),
					resource.TestCheckResourceAttr("data.sigil_mark.vpc", "region_code", "usex"),
				),
			},
		},
	})
}

func TestMarkDataSource_awsBucketFallsBackToAllowedStyle(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud          = "aws"
  org_prefix     = "acme"
  env            = "dev"
  style_priority = ["underscore"]
`, `
data "sigil_mark" "bucket" {
  what      = "s3_bucket"
  qualifier = "logs"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.bucket", "name", "acme-dev-s3bk-logs"),
					resource.TestCheckResourceAttr("data.sigil_mark.bucket", "style", "dashed"),
					resource.TestCheckResourceAttr("data.sigil_mark.bucket", "resource_acronym", "s3bk"),
				),
			},
		},
	})
}

func TestMarkDataSource_awsBucketConstraintFailure(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud      = "aws"
  org_prefix = "acme"
  env        = "dev"
`, `
data "sigil_mark" "bucket" {
  what      = "s3_bucket"
  qualifier = "logs-s3alias"
  recipe    = ["qualifier"]
}
`),
				ExpectError: regexp.MustCompile(`must not\s+end with suffix "-s3alias"`),
			},
		},
	})
}

func TestMarkDataSource_azureRegionalResourceIncludesNormalizedRegion(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud                                = "azure"
  org_prefix                           = "acme"
  env                                  = "dev"
  region                               = "west-europe"
  ignore_region_for_regional_resources = false
`, `
data "sigil_mark" "vnet" {
  what      = "azurerm_virtual_network"
  qualifier = "core"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.vnet", "name", "acme-dev-weu-vnet-core"),
					resource.TestCheckResourceAttr("data.sigil_mark.vnet", "region_code", "weu"),
					resource.TestCheckResourceAttr("data.sigil_mark.vnet", "style", "dashed"),
					resource.TestCheckResourceAttr("data.sigil_mark.vnet", "resource_acronym", "vnet"),
					resource.TestCheckResourceAttr("data.sigil_mark.vnet", "components.region", "weu"),
				),
			},
		},
	})
}

func TestMarkDataSource_azureRegionalResourceOmitsRegionByDefault(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud      = "azure"
  org_prefix = "acme"
  env        = "dev"
  region     = "westeurope"
`, `
data "sigil_mark" "vnet" {
  what      = "azurerm_virtual_network"
  qualifier = "core"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.vnet", "name", "acme-dev-vnet-core"),
					resource.TestCheckResourceAttr("data.sigil_mark.vnet", "region_code", ""),
					resource.TestCheckResourceAttr("data.sigil_mark.vnet", "components.region", ""),
				),
			},
		},
	})
}

func TestMarkDataSource_azureGlobalStorageAccountKeepsRegionAndUsesStraightStyle(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud          = "azure"
  org_prefix     = "acme"
  env            = "dev"
  region         = "westeurope"
  style_priority = ["pascal"]
`, `
data "sigil_mark" "storage" {
  what      = "azurerm_storage_account"
  qualifier = "data"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.storage", "name", "acmedevweustdata"),
					resource.TestCheckResourceAttr("data.sigil_mark.storage", "region_code", "weu"),
					resource.TestCheckResourceAttr("data.sigil_mark.storage", "style", "straight"),
					resource.TestCheckResourceAttr("data.sigil_mark.storage", "resource_acronym", "st"),
					resource.TestCheckResourceAttr("data.sigil_mark.storage", "components.region", "weu"),
				),
			},
		},
	})
}

func TestMarkDataSource_azureStorageAccountFallsBackFromDashedToStraight(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud          = "azure"
  org_prefix     = "acme"
  env            = "dev"
  region         = "westeurope"
  style_priority = ["dashed", "pascaldashed", "camel"]
`, `
data "sigil_mark" "storage" {
  what      = "azurerm_storage_account"
  qualifier = "data-lake"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.storage", "name", "acmedevweustdatalake"),
					resource.TestCheckResourceAttr("data.sigil_mark.storage", "region_code", "weu"),
					resource.TestCheckResourceAttr("data.sigil_mark.storage", "style", "straight"),
					resource.TestCheckResourceAttr("data.sigil_mark.storage", "resource_acronym", "st"),
				),
			},
		},
	})
}

func TestMarkDataSource_azureRegionOverrideWins(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud                                = "azure"
  org_prefix                           = "acme"
  env                                  = "dev"
  region                               = "westeurope"
  ignore_region_for_regional_resources = false
  region_overrides = {
    "westeurope" = "euwx"
  }
`, `
data "sigil_mark" "vnet" {
  what      = "azurerm_virtual_network"
  qualifier = "core"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.vnet", "name", "acme-dev-euwx-vnet-core"),
					resource.TestCheckResourceAttr("data.sigil_mark.vnet", "region_code", "euwx"),
				),
			},
		},
	})
}

func TestMarkDataSource_azureRegionMapReplacementWins(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud                                = "azure"
  org_prefix                           = "acme"
  env                                  = "dev"
  region                               = "westeurope"
  ignore_region_for_regional_resources = false
  region_map = {
    "westeurope" = "cust"
  }
`, `
data "sigil_mark" "vnet" {
  what      = "azurerm_virtual_network"
  qualifier = "core"
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sigil_mark.vnet", "name", "acme-dev-cust-vnet-core"),
					resource.TestCheckResourceAttr("data.sigil_mark.vnet", "region_code", "cust"),
				),
			},
		},
	})
}

func TestMarkDataSource_azureStorageAccountConstraintFailure(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMarkDataSourceConfig(`
  cloud      = "azure"
  org_prefix = "verylongorganization"
  env        = "production"
  region     = "westeurope"
`, `
data "sigil_mark" "storage" {
  what      = "azurerm_storage_account"
  qualifier = "analytics"
}
`),
				ExpectError: regexp.MustCompile(`exceeds 24 characters`),
			},
		},
	})
}

func testAccMarkDataSourceConfig(providerBody, dataBody string) string {
	return fmt.Sprintf(`
%s

%s
`, testAccProviderConfig(providerBody), dataBody)
}

func testAccProviderConfig(providerBody string) string {
	return fmt.Sprintf(`
provider "sigil" {
%s
}
`, providerBody)
}
