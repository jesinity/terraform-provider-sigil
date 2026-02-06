# cloudomen

Terraform provider for AWS naming conventions and consistent resource naming.

## Provider Configuration

```hcl
terraform {
  required_providers {
    cloudomen = {
      source  = "jesinity/cloudomen"
      version = "0.1.0"
    }
  }
}

provider "cloudomen" {
  org_prefix = "acme"
  project    = "iac"
  env        = "dev"
  region     = "ap-southeast-2"

  # Optional: override just one region
  region_overrides = {
    "us-east-1" = "ueue1"
  }

  # Optional: override the full region map
  # region_map = {
  #   "us-east-1" = "use1"
  # }

  # Optional: override the default recipe
  # recipe = ["org", "proj", "env", "region", "resource", "qualifier"]

  # Optional: override style priority
  # style_priority = ["dashed", "pascal", "camel", "straight", "underscore"]
}
```

## Data Source `cloudomen_nomen`

### Basic name

```hcl
data "cloudomen_nomen" "bucket" {
  resource  = "s3"
  qualifier = "mydata"
}

output "bucket_name" {
  value = data.cloudomen_nomen.bucket.name
}
```

### Override a component

```hcl
data "cloudomen_nomen" "bucket" {
  resource  = "s3"
  qualifier = "mydata"

  overrides = {
    env    = "prod"
    region = "use1"
  }
}
```

### Custom recipe and style priority

```hcl
data "cloudomen_nomen" "bucket" {
  resource  = "s3"
  qualifier = "mydata"

  recipe         = ["org", "proj", "env", "resource", "qualifier"]
  style_priority = ["pascal", "camel", "straight", "dashed", "underscore"]
}
```

### More examples

```hcl
data "cloudomen_nomen" "iam_role" {
  resource       = "iam_role"
  qualifier      = "app"
  style_priority = ["pascal", "camel", "dashed"]
}

output "iam_role_name" {
  value = data.cloudomen_nomen.iam_role.name
  # Example: "AcmeIacDevApse2RoleApp"
}

output "iam_role_style" {
  value = data.cloudomen_nomen.iam_role.style
  # Example: "pascal"
}
```

```hcl
data "cloudomen_nomen" "lambda" {
  resource       = "lambda"
  qualifier      = "ingest"
  style_priority = ["underscore", "dashed"]
}

output "lambda_name" {
  value = data.cloudomen_nomen.lambda.name
  # Example: "acme_iac_dev_apse2_lmbd_ingest"
}

output "lambda_resource_acronym" {
  value = data.cloudomen_nomen.lambda.resource_acronym
  # Example: "lmbd"
}
```

```hcl
data "cloudomen_nomen" "queue" {
  resource  = "sqs"
  qualifier = "jobs"
}

output "queue_name" {
  value = data.cloudomen_nomen.queue.name
  # Example: "acme-iac-dev-apse2-sqs-jobs"
}

output "queue_style" {
  value = data.cloudomen_nomen.queue.style
  # Example: "dashed"
}
```

## Outputs

The data source returns:
- `name`
- `style`
- `region_code`
- `resource_acronym`
- `components`
- `parts`

### Output Examples

```hcl
output "bucket_name" {
  value = data.cloudomen_nomen.bucket.name
  # Example: "acme-iac-dev-apse2-s3b-mydata"
}

output "bucket_style" {
  value = data.cloudomen_nomen.bucket.style
  # Example: "dashed"
}

output "bucket_region_code" {
  value = data.cloudomen_nomen.bucket.region_code
  # Example: "apse2"
}

output "bucket_resource_acronym" {
  value = data.cloudomen_nomen.bucket.resource_acronym
  # Example: "s3b"
}

output "bucket_parts" {
  value = data.cloudomen_nomen.bucket.parts
  # Example: ["acme", "iac", "dev", "apse2", "s3b", "mydata"]
}

output "bucket_components" {
  value = data.cloudomen_nomen.bucket.components
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

## Naming Styles

Style priority determines how names are formatted. If a resource has style constraints, the provider selects the first allowed style in the priority list.

Valid styles:
- `dashed`
- `underscore`
- `straight`
- `pascal`
- `camel`

Style behaviors:
- `dashed` Lowercase words joined by `-`.
- `underscore` Lowercase words joined by `_`.
- `straight` Lowercase words concatenated.
- `pascal` Words in `PascalCase`.
- `camel` Words in `camelCase`.

Words are extracted from each component using the pattern `[A-Za-z0-9]+`, so punctuation or separators become word boundaries. If no valid style matches the priority list and any resource overrides, the provider falls back to `dashed`.

By default, `s3` and `s3_bucket` are restricted to `dashed` and `straight` to align with S3 naming rules.

## Resource Constraints

Some resources have naming constraints enforced after formatting. The constraint name matches the `resource` input (case-insensitive). The table below lists built-in constraints and their limits.

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
