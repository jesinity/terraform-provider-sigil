# sigil_mark Data Source

Generates a resource name and its components using provider configuration, optional overrides, and style preferences.

## Example Usage

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

The data source selects the first valid style from `style_priority` (request-specific) or the provider `style_priority` when none is supplied. If `resource_style_overrides` defines an allowed style list for the current `what`, only those styles are considered. If no style matches, the provider falls back to `dashed`.

By default, `s3` and `s3_bucket` are restricted to `dashed` and `straight` to align with S3 naming rules.

Valid styles and their output shapes:
- `dashed` Lowercase words joined by `-`.
- `underscore` Lowercase words joined by `_`.
- `straight` Lowercase words concatenated.
- `pascal` Words in `PascalCase`.
- `camel` Words in `camelCase`.

Words are extracted from each component using the pattern `[A-Za-z0-9]+`, so punctuation or separators become word boundaries.

## Resource Constraints

Some resources enforce naming constraints after formatting. The constraint name is the `what` input (case-insensitive). If the computed name violates a constraint, the data source returns an error.

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
