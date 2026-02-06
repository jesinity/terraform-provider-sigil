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

## Notes

Default recipe components are `org`, `proj`, `env`, `region`, `resource`, and `qualifier`. If both `region_map` and `region_overrides` are set, overrides are applied to the map.
