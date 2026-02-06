# cloudomen Provider

Terraform provider for generating consistent AWS resource names based on a shared recipe and style priorities.

## Example Usage

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

  region_overrides = {
    "us-east-1" = "ueue1"
  }

  # recipe = ["org", "proj", "env", "region", "resource", "qualifier"]
  # style_priority = ["dashed", "pascal", "camel", "straight", "underscore"]
}
```

## Argument Reference

- `org_prefix` (Required) Short organization identifier.
- `project` (Optional) Project or workload identifier.
- `env` (Required) Environment identifier, such as `dev`, `staging`, or `prod`.
- `region` (Optional) AWS region name, used to derive a short region code.
- `region_short_code` (Optional) Explicit short region code to use instead of mapping.
- `region_map` (Optional) Full region map; when set, replaces the default map.
- `region_overrides` (Optional) Map of region overrides applied on top of the default map.
- `recipe` (Optional) Ordered list of components used to build the name.
- `style_priority` (Optional) Preferred naming styles in order of precedence.
- `resource_acronyms` (Optional) Map of resource identifiers to acronyms.
- `resource_style_overrides` (Optional) Map of resource identifiers to allowed styles.

## Style Priority

The provider formats names by applying the first usable style in `style_priority`. Styles are evaluated in order and normalized to lowercase. If a resource has an allowed style list (from `resource_style_overrides`), only styles in that list are considered. If none match, the provider falls back to `dashed`.

By default, `s3` and `s3_bucket` are restricted to `dashed` and `straight` to align with S3 naming rules.

Valid styles and how they format parts:
- `dashed` Lowercase words joined by `-` (e.g., `acme-iac-dev-apse2-s3b-mydata`).
- `underscore` Lowercase words joined by `_` (e.g., `acme_iac_dev_apse2_s3b_mydata`).
- `straight` Lowercase words concatenated (e.g., `acmeiacdevapse2s3bmydata`).
- `pascal` Words in `PascalCase` (e.g., `AcmeIacDevApse2S3bMydata`).
- `camel` Words in `camelCase` (e.g., `acmeIacDevApse2S3bMydata`).

Words are extracted from each component using the pattern `[A-Za-z0-9]+`. Any separators or punctuation are treated as boundaries.

## Resource Constraints

Some resource identifiers have naming constraints enforced after formatting. If the generated name violates a constraint, the data source will return an error. The constraint names match the `resource` input (case-insensitive). Resources with built-in constraints include:
- `s3`
- `s3_bucket`
- `role`
- `iam_role`
- `iam_user`
- `iam_group`
- `iam_policy`
- `role_policy`
- `sns`
- `sns_topic`
- `sqs`
- `sqs_queue`
- `lambda`
- `kms_alias`
- `log_group`
- `cloudwatch_log_group`
- `sec_group`
- `security_group`

Constraint types include minimum or maximum length, required pattern, forbidden prefixes or suffixes, forbidden substrings, and checks that the name is not formatted as an IPv4 address.

## Notes

Default recipe components are `org`, `proj`, `env`, `region`, `resource`, and `qualifier`. If both `region_map` and `region_overrides` are set, overrides are applied to the map.
