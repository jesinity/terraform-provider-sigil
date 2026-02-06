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

## Outputs

The data source returns:
- `name`
- `style`
- `region_code`
- `resource_acronym`
- `components`
- `parts`

## Naming Styles

Style priority determines how names are formatted. If a resource has style constraints, the provider selects the first allowed style in the priority list.

Valid styles:
- `dashed`
- `underscore`
- `straight`
- `pascal`
- `camel`
