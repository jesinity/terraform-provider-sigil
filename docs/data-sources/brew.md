# awsnomicon_brew Data Source

Generates a resource name and its components using provider configuration, optional overrides, and style preferences.

## Example Usage

```hcl
data "awsnomicon_brew" "bucket" {
  resource  = "s3"
  qualifier = "mydata"
}

output "bucket_name" {
  value = data.awsnomicon_brew.bucket.name
}
```

## Argument Reference

- `resource` (Required) Resource identifier, such as `s3` or `iam_role`.
- `qualifier` (Optional) Additional name segment to distinguish similar resources.
- `overrides` (Optional) Map of component overrides, such as `org`, `proj`, `env`, `region`, `resource`, or `qualifier`.
- `recipe` (Optional) Ordered list of components used to build the name for this request.
- `style_priority` (Optional) Preferred naming styles in order of precedence for this request.

## Attributes Reference

- `name` The final computed name.
- `style` The style used to format the name.
- `region_code` The resolved short region code.
- `resource_acronym` The resolved resource acronym.
- `components` Map of computed component values.
- `parts` Ordered list of name parts used to construct `name`.
