# sigil

Terraform provider for AWS naming conventions and consistent resource naming. Today it's AWS-first, with room to expand to other clouds in the future.

## Provider Configuration

```hcl
terraform {
  required_providers {
    sigil = {
      source  = "jesinity/sigil"
      version = "0.2.0"
    }
  }
}

provider "sigil" {
  org_prefix = "acme"
  project    = "iac"
  env        = "dev"
  region     = "ap-southeast-2"
  # Optional: omit region for regional resources (default true)
  ignore_region_for_regional_resources = false

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

## Data Source `sigil_mark`

`what` identifies the resource type (formerly `resource`) and drives acronyms, style overrides, and constraints. The `resource` argument is still accepted but deprecated. In recipes and outputs, the component key remains `resource` (alias `what`).

### Basic name

```hcl
data "sigil_mark" "bucket" {
  what      = "s3"
  qualifier = "mydata"
}

output "bucket_name" {
  value = data.sigil_mark.bucket.name
}
```

### Override a component

```hcl
data "sigil_mark" "bucket" {
  what      = "s3"
  qualifier = "mydata"

  overrides = {
    env    = "prod"
    region = "use1"
  }
}
```

### Custom recipe and style priority

```hcl
data "sigil_mark" "bucket" {
  what      = "s3"
  qualifier = "mydata"

  recipe         = ["org", "proj", "env", "resource", "qualifier"]
  style_priority = ["pascal", "camel", "straight", "dashed", "underscore"]
}
```

### More examples

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

## Recipe and Optional Components

The recipe is an ordered list of components. Components are only included when they have a non-empty value, so you can omit any component by removing it from the recipe or leaving it empty. Configure a default recipe at the provider level and override per data source with `recipe`.

Example: omit `env` from the name:

```hcl
provider "sigil" {
  org_prefix = "acme"
  project    = "iac"
  env        = "dev"
  region     = "ap-southeast-2"

  recipe = ["org", "proj", "region", "resource", "qualifier"]
}
```

## Region Handling

When `ignore_region_for_regional_resources` is `true` (default), the `region` component is omitted for resources marked as `regional` in the table below. Resources marked `global` keep the region component even when the flag is enabled. Set it to `false` to always include the region in names. You can still force a region per name via `overrides`. When the region is omitted, `region_code` will be empty unless overridden.

## Resource Acronyms and Scope

Default resource acronyms and scope. Scope is used by `ignore_region_for_regional_resources`. You can override acronyms with `resource_acronyms`.

| Resource | Acronym | Scope |
| --- | --- | --- |
| `acm_cert` | `acmc` | `regional` |
| `alb` | `albl` | `regional` |
| `api_gateway_model` | `agmd` | `regional` |
| `api_gateway_rest_api` | `agra` | `regional` |
| `api_gateway_v2` | `agv2` | `regional` |
| `appsync` | `apsy` | `regional` |
| `athena` | `athn` | `regional` |
| `aurora_cluster` | `arcl` | `regional` |
| `autoscaling_group` | `asgr` | `regional` |
| `cloudformation_stack` | `cfst` | `regional` |
| `cloudfront` | `clfr` | `global` |
| `cloudtrail` | `ctra` | `regional` |
| `cloudwatch_alarm` | `cwal` | `regional` |
| `cloudwatch_log_group` | `cwlg` | `regional` |
| `codebuild` | `cdbd` | `regional` |
| `codedeploy` | `cddp` | `regional` |
| `codepipeline` | `cdpl` | `regional` |
| `config_rule` | `cfrl` | `regional` |
| `dynamodb` | `dydb` | `regional` |
| `dynamodb_table` | `dydb` | `regional` |
| `ebs` | `ebs` | `regional` |
| `ec2_instance` | `ec2i` | `regional` |
| `ecr` | `ecr` | `regional` |
| `ecs` | `ecs` | `regional` |
| `ecs_cluster` | `ecsc` | `regional` |
| `ecs_service` | `ecss` | `regional` |
| `ecs_task` | `ecst` | `regional` |
| `efs` | `efs` | `regional` |
| `eks` | `eks` | `regional` |
| `eks_cluster` | `eksc` | `regional` |
| `eks_node_group` | `ekng` | `regional` |
| `elastic_ip` | `elip` | `regional` |
| `elasticache` | `elch` | `regional` |
| `elasticsearch` | `elsr` | `regional` |
| `elb` | `elbl` | `regional` |
| `eventbridge_bus` | `evbb` | `regional` |
| `eventbridge_rule` | `evbr` | `regional` |
| `glue` | `glue` | `regional` |
| `guardduty` | `gdty` | `regional` |
| `iam_group` | `iamg` | `global` |
| `iam_policy` | `iamp` | `global` |
| `iam_role` | `role` | `global` |
| `iam_user` | `iamu` | `global` |
| `igw` | `igtw` | `regional` |
| `kms_key` | `kmsk` | `regional` |
| `lambda` | `lmbd` | `regional` |
| `launch_template` | `lcht` | `regional` |
| `log_group` | `logg` | `regional` |
| `msk_cluster` | `mskc` | `regional` |
| `nacl` | `nacl` | `regional` |
| `nat_gw` | `ngtw` | `regional` |
| `nlb` | `nlbl` | `regional` |
| `opensearch` | `opsr` | `regional` |
| `rds` | `rds` | `regional` |
| `rds_cluster` | `rdsc` | `regional` |
| `redshift` | `rdsh` | `regional` |
| `role` | `role` | `global` |
| `role_policy` | `rlpl` | `global` |
| `route53_record` | `r53r` | `global` |
| `route53_zone` | `rt53` | `global` |
| `route_table` | `rttb` | `regional` |
| `s3` | `s3b` | `regional` |
| `s3_access_point` | `s3ap` | `regional` |
| `s3_bucket` | `s3bk` | `regional` |
| `s3_dir` | `s3dr` | `regional` |
| `s3_object` | `s3ob` | `regional` |
| `s3_table` | `s3tb` | `regional` |
| `sagemaker` | `sgmk` | `regional` |
| `sec_group` | `scgp` | `regional` |
| `secretsmanager_secret` | `smse` | `regional` |
| `sfn` | `stfn` | `regional` |
| `snow_notification_integration` | `snti` | `regional` |
| `sns` | `sns` | `regional` |
| `sqs` | `sqs` | `regional` |
| `ssm_parameter` | `ssmp` | `regional` |
| `step_function` | `stfn` | `regional` |
| `subnet` | `subn` | `regional` |
| `target_group` | `tgpt` | `regional` |
| `vpc` | `vpcn` | `regional` |
| `wafv2_ip_set` | `wfis` | `regional` |
| `wafv2_web_acl` | `wfac` | `regional` |
| `wafv2_web_acl_rule` | `wfar` | `regional` |

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

Some resources have naming constraints enforced after formatting. The constraint name matches the `what` input (case-insensitive). The table below lists built-in constraints and their limits.

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

## Argument Reference

- `org_prefix` (Required) Short organization identifier.
- `project` (Optional) Project or workload identifier.
- `env` (Required) Environment identifier, such as `dev`, `staging`, or `prod`.
- `region` (Optional) AWS region name, used to derive a short region code.
- `region_short_code` (Optional) Explicit short region code to use instead of mapping.
- `region_map` (Optional) Full region map; when set, replaces the default map.
- `region_overrides` (Optional) Map of region overrides applied on top of the default map.
- `ignore_region_for_regional_resources` (Optional) When `true` (default), omit the region component for resources marked as `regional` in the acronyms table.
- `recipe` (Optional) Ordered list of components used to build the name.
- `style_priority` (Optional) Preferred naming styles in order of precedence.
- `resource_acronyms` (Optional) Map of resource identifiers to acronyms.
- `resource_style_overrides` (Optional) Map of resource identifiers to allowed styles.

## Notes

Default recipe components are `org`, `proj`, `env`, `region`, `resource`, and `qualifier`. Components are only included when non-empty, and you can omit them by removing items from the recipe. If both `region_map` and `region_overrides` are set, overrides are applied to the map. When `ignore_region_for_regional_resources` is `true`, the region component is omitted for regional resources unless explicitly overridden.
