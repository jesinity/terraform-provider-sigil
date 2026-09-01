# Changelog

## v1.5.0 - 2026-08-28

### Added

- Optional `platform = "databricks"` configuration with composable Databricks defaults for AWS, Azure, and GCP.
- Audited, source-controlled Databricks naming manifest and canonical Terraform resource lookup with platform-only aliases.
- Databricks provisioning coverage, complete v1.129.0 catalog classification, and an upstream drift comparison tool.

### Changed

- Hardened manifest validation for unsupported styles and duplicate canonical resources.
- Deep-copied constraint slices and completed category/scope classification for all 170 upstream resources.

### Compatibility

- This is an additive minor release. Existing cloud-only configurations preserve their prior behavior when `platform` is omitted.

## v1.4.0 - 2026-08-04

### Added

- Audited AWS naming coverage for Athena, Glue, EMR, Kinesis, Lake Formation, Redshift, SageMaker, Bedrock, Bedrock Agents, Bedrock AgentCore, DataZone, and QuickSight.
- Added 153 constrained, nameable Terraform resources and eight service aliases.
- Added canonical `aws_` Terraform resource lookup and compatibility aliases for established Sigil keys.
- Added Bedrock evaluation job coverage.

### Changed

- Enforced AWS-compatible styles and length/pattern constraints for the new data engineering, analytics, and ML resources.
- Classified and documented 48 configuration, association, policy, registration, and generated-version resources that cannot consume an independent Sigil name.
- Updated GoReleaser configuration for the current v2 schema.

### Compatibility

- All 82 AWS acronym mappings from v1.3.0 are unchanged.
- No legacy AWS acronym was removed or renamed.
- The only repeated acronym values remain the intentional `role`/`iam_role` and `sfn`/`step_function` aliases.
