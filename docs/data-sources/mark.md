# sigil_mark Data Source

Generates a resource name and its components using provider configuration, optional overrides, and style preferences.

Behavior is cloud-aware:
- `cloud = "aws"` uses built-in AWS acronyms and constraints.
- `cloud = "azure"` uses Azure CAF resource definitions (acronyms, style rules, and regex constraints).
- `cloud = "gcp"` uses starter GCP defaults with strict constraints for bucket/network/subnetwork, Pub/Sub, service account, BigQuery dataset, and Cloud Run service resources.
- CAF resource catalog JSON: https://github.com/aztfmod/terraform-provider-azurecaf/blob/main/resourceDefinition.json
- Azure CAF examples: `azurerm_resource_group -> rg`, `azurerm_storage_account -> st`.

## Example Usage

### Azure Bootstrap Snippet

Use this as a starter for a new Azure project that uses Sigil for naming:

```hcl
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }

    sigil = {
      source  = "jesinity/sigil"
      version = "~> 1.3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "sigil" {
  cloud      = "azure"
  org_prefix = "acme"
  project    = "payments"
  env        = "dev"
  region     = "westeurope"
  # Uses built-in Azure region short code mapping: westeurope -> weu.
}

data "sigil_mark" "rg" {
  what      = "azurerm_resource_group"
  qualifier = "core"
}

data "sigil_mark" "storage" {
  what      = "azurerm_storage_account"
  qualifier = "raw"

  # Storage accounts are lowercase + no separators in CAF defaults.
  recipe = ["org", "proj", "env", "resource", "qualifier"]
}

resource "azurerm_resource_group" "this" {
  name     = data.sigil_mark.rg.name
  location = "West Europe"
}

resource "azurerm_storage_account" "this" {
  name                     = data.sigil_mark.storage.name
  resource_group_name      = azurerm_resource_group.this.name
  location                 = azurerm_resource_group.this.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
```

These examples assume `ignore_region_for_regional_resources = false`. If you keep the default `true` and the resource is marked `regional`, the region segment and `region_code` will be omitted.

```hcl
data "sigil_mark" "bucket" {
  what      = "s3"
  qualifier = "mydata"
}

output "bucket_name" {
  value = data.sigil_mark.bucket.name
  # Example: "acme-iac-dev-apse2-s3b-mydata"
}

output "bucket_style" {
  value = data.sigil_mark.bucket.style
  # Example: "dashed"
}

output "bucket_region_code" {
  value = data.sigil_mark.bucket.region_code
  # Example: "apse2"
}

output "bucket_resource_acronym" {
  value = data.sigil_mark.bucket.resource_acronym
  # Example: "s3b"
}

output "bucket_parts" {
  value = data.sigil_mark.bucket.parts
  # Example: ["acme", "iac", "dev", "apse2", "s3b", "mydata"]
}

output "bucket_components" {
  value = data.sigil_mark.bucket.components
  # Example:
  # {
  #   org       = "acme"
  #   proj      = "iac"
  #   env       = "dev"
  #   region    = "apse2"
  #   resource  = "s3b"
  #   qualifier = "mydata"
  # }
}
```

### More Examples

```hcl
data "sigil_mark" "iam_role" {
  what           = "iam_role"
  qualifier      = "app"
  style_priority = ["pascal", "camel", "dashed"]
}

output "iam_role_name" {
  value = data.sigil_mark.iam_role.name
  # Example: "AcmeIacDevApse2RoleApp"
}

output "iam_role_style" {
  value = data.sigil_mark.iam_role.style
  # Example: "pascal"
}
```

```hcl
data "sigil_mark" "lambda" {
  what           = "lambda"
  qualifier      = "ingest"
  style_priority = ["underscore", "dashed"]
}

output "lambda_name" {
  value = data.sigil_mark.lambda.name
  # Example: "acme_iac_dev_apse2_lmbd_ingest"
}

output "lambda_resource_acronym" {
  value = data.sigil_mark.lambda.resource_acronym
  # Example: "lmbd"
}
```

```hcl
data "sigil_mark" "queue" {
  what      = "sqs"
  qualifier = "jobs"
}

output "queue_name" {
  value = data.sigil_mark.queue.name
  # Example: "acme-iac-dev-apse2-sqs-jobs"
}

output "queue_style" {
  value = data.sigil_mark.queue.style
  # Example: "dashed"
}
```

```hcl
data "sigil_mark" "azure_storage_account" {
  what      = "azurerm_storage_account"
  qualifier = "raw"

  recipe = ["org", "proj", "env", "resource", "qualifier"]
}

output "azure_storage_account_name" {
  value = data.sigil_mark.azure_storage_account.name
  # Example: "acmepaymentsprodstraw"
}

output "azure_storage_account_style" {
  value = data.sigil_mark.azure_storage_account.style
  # Example: "straight"
}
```

## Argument Reference

- `what` (Required) Resource identifier, such as `s3` or `iam_role`.
- `resource` (Deprecated) Alias for `what`. The `components` output still uses the `resource` key.
- `qualifier` (Optional) Additional name segment to distinguish similar resources.
- `overrides` (Optional) Map of component overrides, such as `org`, `proj`, `env`, `region`, `resource` (or `what`), or `qualifier`.
- `recipe` (Optional) Ordered list of components used to build the name for this request.
- `style_priority` (Optional) Preferred naming styles in order of precedence for this request.

## Attributes Reference

- `name` The final computed name.
- `style` The style used to format the name.
- `region_code` The resolved short region code.
- `resource_acronym` The resolved resource acronym.
- `components` Map of computed component values.
- `parts` Ordered list of name parts used to construct `name`.

## Style Priority Resolution

The data source selects the first valid style from `style_priority` (request-specific) or the provider `style_priority` when none is supplied. If `resource_style_overrides` defines an allowed style list for the current `what`, only those styles are considered. If no style matches, Sigil falls back to the first allowed style for that resource, or `dashed` when no style override exists.

Cloud-specific style overrides are applied automatically:
- `aws`: `s3` and `s3_bucket` are restricted to `dashed` and `straight`.
- `azure`: each CAF resource inherits style limits from CAF dash/lowercase metadata.
- `gcp`: starter resources include style restrictions for Tier-A compatibility, including bucket/network, service account, BigQuery dataset, Pub/Sub, and Cloud Run resources.

Valid styles and their output shapes:
- `dashed` Lowercase words joined by `-`.
- `underscore` Lowercase words joined by `_`.
- `straight` Lowercase words concatenated.
- `pascal` Words in `PascalCase`.
- `pascaldashed` Words in `Pascal-Case` joined by `-`.
- `camel` Words in `camelCase`.

Words are extracted from each component using the pattern `[A-Za-z0-9]+`, so punctuation or separators become word boundaries.

## Resource Constraints

Some resources enforce naming constraints after formatting. The constraint name is the `what` input (case-insensitive). If the computed name violates a constraint, the data source returns an error.

The table below lists built-in `aws` constraints. Azure constraints are listed in `../azure-caf-resources.md` and sourced from Azure naming rules plus Azure CAF definitions. GCP starter constraints currently cover `google_storage_bucket`, `google_compute_network`, `google_compute_subnetwork`, `google_pubsub_topic`, `google_pubsub_subscription`, `google_service_account`, `google_bigquery_dataset`, and `google_cloud_run_v2_service` (including aliases).

| Resource | Min | Max | Pattern | Notes |
| --- | --- | --- | --- | --- |
| `s3` | 3 | 63 | lowercase letters, numbers, dots, and hyphens; must start and end with a letter or number | Forbidden prefixes: `xn--`, `sthree-`, `amzn-s3-demo-`; forbidden suffixes: `-s3alias`, `--ol-s3`; forbidden substrings: `..`; disallow IPv4 |
| `s3_bucket` | 3 | 63 | lowercase letters, numbers, dots, and hyphens; must start and end with a letter or number | Forbidden prefixes: `xn--`, `sthree-`, `amzn-s3-demo-`; forbidden suffixes: `-s3alias`, `--ol-s3`; forbidden substrings: `..`; disallow IPv4 |
| `role` | 1 | 64 | alphanumeric and the following: `+=,.@_-` | none |
| `iam_role` | 1 | 64 | alphanumeric and the following: `+=,.@_-` | none |
| `iam_user` | 1 | 64 | alphanumeric and the following: `+=,.@_-` | none |
| `iam_group` | 1 | 128 | alphanumeric and the following: `+=,.@_-` | none |
| `iam_policy` | 1 | 128 | alphanumeric and the following: `+=,.@_-` | none |
| `role_policy` | 1 | 128 | alphanumeric and the following: `+=,.@_-` | none |
| `sns` | 1 | 256 | letters, numbers, underscores, and hyphens; FIFO topics must end with `.fifo` | none |
| `sns_topic` | 1 | 256 | letters, numbers, underscores, and hyphens; FIFO topics must end with `.fifo` | none |
| `sqs` | 1 | 80 | letters, numbers, underscores, and hyphens; FIFO queues must end with `.fifo` | none |
| `sqs_queue` | 1 | 80 | letters, numbers, underscores, and hyphens; FIFO queues must end with `.fifo` | none |
| `lambda` | 1 | 64 | letters, numbers, hyphens, and underscores | none |
| `kms_alias` | 1 | 256 | must begin with `alias/` and contain only letters, numbers, slashes, underscores, and hyphens | Forbidden prefix: `alias/aws/` |
| `log_group` | 1 | 512 | letters, numbers, underscore, hyphen, slash, period, and `#` | Forbidden prefix: `aws/` |
| `cloudwatch_log_group` | 1 | 512 | letters, numbers, underscore, hyphen, slash, period, and `#` | Forbidden prefix: `aws/` |
| `sec_group` | 1 | 255 | letters, numbers, spaces, and `._-:/()#,@[]+=&;{}!$*` | Forbidden prefix: `sg-` (case-insensitive) |
| `security_group` | 1 | 255 | letters, numbers, spaces, and `._-:/()#,@[]+=&;{}!$*` | Forbidden prefix: `sg-` (case-insensitive) |

Constraint types include minimum or maximum length, required pattern, forbidden prefixes or suffixes, forbidden substrings, and checks that the name is not formatted as an IPv4 address.
