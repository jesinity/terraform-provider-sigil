# sigil

Terraform provider for consistent resource naming across multiple clouds. `aws` is the default cloud profile, including built-in data engineering and data science coverage for Athena, Glue, EMR, Kinesis, Lake Formation, Redshift, SageMaker, Bedrock, Bedrock Agents, and Bedrock AgentCore. `azure` uses Azure CAF resource coverage, and `gcp` includes built-in resource coverage with strict constraints for supported resource families.

## Provider Configuration

```hcl
terraform {
  required_providers {
    sigil = {
      source  = "jesinity/sigil"
      version = "~> 1.4.0"
    }
  }
}

provider "sigil" {
  # Optional, defaults to "aws"
  cloud = "aws"

  org_prefix = "acme"
  project    = "iac"
  env        = "dev"
  region     = "ap-southeast-2"
  # Optional: omit region for resources marked regional (default true)
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
  # style_priority = ["dashed", "pascal", "pascaldashed", "camel", "straight", "underscore"]
}
```

For reuse across multiple provider aliases, you can supply a base `config` object and apply `overrides`. Precedence is: `config` -> top-level attributes -> `overrides`. Top-level attributes are a shorthand for the common case.

```hcl
locals {
  sigil_config = {
    cloud      = "aws"
    org_prefix = "acme"
    project    = "iac"
    env        = "dev"
    region     = "ap-southeast-2"

    ignore_region_for_regional_resources = false
    region_overrides = {
      "us-east-1" = "ueue1"
    }
  }
}

provider "sigil" {
  config = local.sigil_config
}

provider "sigil" {
  alias  = "secondary"
  config = local.sigil_config

  overrides = {
    region = "us-east-1"
  }
}
```

Azure example (`cloud = "azure"`):

```hcl
provider "sigil" {
  cloud      = "azure"
  org_prefix = "acme"
  project    = "payments"
  env        = "prod"
  region     = "westeurope"

  # Azure defaults include a built-in region short code map.
  # Example: westeurope -> weu, eastus2 -> eus2.
  # Optional: override if your org uses different codes.
  # region_overrides = {
  #   westeurope = "weu"
  #   eastus2    = "eus2"
  # }

  # Optional Azure-specific overrides
  # resource_acronyms = {
  #   azurerm_storage_account = "st" # CAF default shown here as an explicit override example.
  # }
}
```

GCP example (`cloud = "gcp"`):

```hcl
provider "sigil" {
  cloud      = "gcp"
  org_prefix = "acme"
  project    = "payments"
  env        = "prod"
  region     = "us-central1"

  # GCP built-in constraints currently cover:
  # - google_storage_bucket
  # - named Compute Engine resources in the default GCP map
  # - google_pubsub_topic / google_pubsub_subscription
  # - google_service_account
  # - google_bigquery_dataset
  # - google_cloud_run_v2_service
}
```

**Why This Design**
Top-level attributes keep the provider fast to configure for the common single-provider case. The optional `config` + `overrides` pattern reduces repetition when you need multiple provider aliases with small differences (like region), without forcing everyone into extra nesting. The merge order is explicit so it is easy to reason about which values win.

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

```hcl
data "sigil_mark" "azure_storage_account" {
  what      = "azurerm_storage_account"
  qualifier = "raw"

  # Azure storage accounts are lowercase/no-dash constrained.
  # The Azure cloud defaults select an allowed style automatically.
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

## AWS Coverage and Resource Acronyms

Sigil accepts canonical Terraform identifiers with or without the `aws_` prefix. Exact keys take precedence, and canonical aliases fall back to their legacy Sigil key without changing existing output.

The AWS catalog contains **243 supported keys**:

- **82 legacy keys**, preserved byte-for-byte from v1.3.0.
- **153 exact Terraform resources** for the expanded data engineering, analytics, and ML services, each with an enforced naming constraint.
- **8 service aliases** for Bedrock, Bedrock Agents, Bedrock AgentCore, DataZone, EMR, Kinesis, Lake Formation, and QuickSight.
- **48 explicitly excluded Terraform resources** that have no independent name or standard tags.

The audit is pinned to HashiCorp AWS provider commit [`63681ad6`](https://github.com/hashicorp/terraform-provider-aws/commit/63681ad684441377ec2f220729f21c3d115778ca). The only repeated acronym values are the intentional legacy aliases `role`/`iam_role` and `sfn`/`step_function`.

`Surface` indicates where the generated value can be applied: a direct resource identifier, standard AWS tags, or an alias that resolves to another supported key.

| Resource | Acronym | Scope | Surface |
| --- | --- | --- | --- |
| `acm_cert` | `acmc` | `regional` | Legacy alias |
| `alb` | `albl` | `regional` | Legacy alias |
| `api_gateway_model` | `agmd` | `regional` | Direct |
| `api_gateway_rest_api` | `agra` | `regional` | Direct + tags |
| `api_gateway_v2` | `agv2` | `regional` | Legacy alias |
| `appsync` | `apsy` | `regional` | Legacy alias |
| `athena` | `athn` | `regional` | Legacy alias |
| `athena_capacity_reservation` | `athcr` | `regional` | Direct + tags |
| `athena_data_catalog` | `athdc` | `regional` | Direct + tags |
| `athena_database` | `athdb` | `regional` | Direct |
| `athena_named_query` | `athnq` | `regional` | Direct |
| `athena_prepared_statement` | `athps` | `regional` | Direct |
| `athena_workgroup` | `athwg` | `regional` | Direct + tags |
| `aurora_cluster` | `arcl` | `regional` | Legacy alias |
| `autoscaling_group` | `asgr` | `regional` | Direct |
| `bedrock` | `bdrk` | `regional` | Service alias |
| `bedrock_custom_model` | `brcm` | `regional` | Direct + tags |
| `bedrock_evaluation_job` | `brevj` | `regional` | Direct + tags |
| `bedrock_guardrail` | `brgr` | `regional` | Direct + tags |
| `bedrock_inference_profile` | `brip` | `regional` | Direct + tags |
| `bedrock_provisioned_model_throughput` | `brpmt` | `regional` | Direct + tags |
| `bedrockagent` | `brag` | `regional` | Service alias |
| `bedrockagent_agent` | `braga` | `regional` | Direct + tags |
| `bedrockagent_agent_action_group` | `bragg` | `regional` | Direct |
| `bedrockagent_agent_alias` | `braal` | `regional` | Direct + tags |
| `bedrockagent_agent_collaborator` | `bracl` | `regional` | Direct |
| `bedrockagent_data_source` | `brads` | `regional` | Direct |
| `bedrockagent_flow` | `brafl` | `regional` | Direct + tags |
| `bedrockagent_knowledge_base` | `brakb` | `regional` | Direct + tags |
| `bedrockagent_prompt` | `brapt` | `regional` | Direct + tags |
| `bedrockagentcore` | `brac` | `regional` | Service alias |
| `bedrockagentcore_agent_runtime` | `bracr` | `regional` | Direct + tags |
| `bedrockagentcore_agent_runtime_endpoint` | `brace` | `regional` | Direct + tags |
| `bedrockagentcore_api_key_credential_provider` | `brcak` | `regional` | Direct + tags |
| `bedrockagentcore_browser` | `brcbr` | `regional` | Direct + tags |
| `bedrockagentcore_browser_profile` | `brcbp` | `regional` | Direct + tags |
| `bedrockagentcore_code_interpreter` | `brcci` | `regional` | Direct + tags |
| `bedrockagentcore_evaluator` | `brcev` | `regional` | Direct + tags |
| `bedrockagentcore_gateway` | `brcgw` | `regional` | Direct + tags |
| `bedrockagentcore_gateway_target` | `brcgt` | `regional` | Direct |
| `bedrockagentcore_harness` | `brchr` | `regional` | Direct + tags |
| `bedrockagentcore_memory` | `brcme` | `regional` | Direct + tags |
| `bedrockagentcore_memory_strategy` | `brcms` | `regional` | Direct |
| `bedrockagentcore_oauth2_credential_provider` | `brco2` | `regional` | Direct + tags |
| `bedrockagentcore_online_evaluation_config` | `brcoe` | `regional` | Direct + tags |
| `bedrockagentcore_policy` | `brcpl` | `regional` | Direct |
| `bedrockagentcore_policy_engine` | `brcpe` | `regional` | Direct + tags |
| `bedrockagentcore_registry` | `brcrg` | `regional` | Direct |
| `bedrockagentcore_workload_identity` | `brcwi` | `regional` | Direct |
| `cloudformation_stack` | `cfst` | `regional` | Direct + tags |
| `cloudfront` | `clfr` | `regional` | Legacy alias |
| `cloudtrail` | `ctra` | `regional` | Direct + tags |
| `cloudwatch_alarm` | `cwal` | `regional` | Legacy alias |
| `cloudwatch_log_group` | `cwlg` | `regional` | Direct + tags |
| `codebuild` | `cdbd` | `regional` | Legacy alias |
| `codedeploy` | `cddp` | `regional` | Legacy alias |
| `codepipeline` | `cdpl` | `regional` | Direct + tags |
| `config_rule` | `cfrl` | `regional` | Legacy alias |
| `datazone` | `dtzn` | `regional` | Service alias |
| `datazone_asset_type` | `dzaty` | `regional` | Direct |
| `datazone_domain` | `dzdmn` | `regional` | Direct |
| `datazone_environment` | `dzenv` | `regional` | Direct |
| `datazone_environment_profile` | `dzepf` | `regional` | Direct |
| `datazone_form_type` | `dzfty` | `regional` | Direct |
| `datazone_glossary` | `dzglo` | `regional` | Direct |
| `datazone_glossary_term` | `dzglt` | `regional` | Direct |
| `datazone_project` | `dzprj` | `regional` | Direct |
| `datazone_user_profile` | `dzusr` | `regional` | Direct |
| `dynamodb` | `dydb` | `regional` | Legacy alias |
| `dynamodb_table` | `dybt` | `regional` | Direct + tags |
| `ebs` | `ebs` | `regional` | Legacy alias |
| `ec2_instance` | `ec2i` | `regional` | Legacy alias |
| `ecr` | `ecr` | `regional` | Legacy alias |
| `ecs` | `ecs` | `regional` | Legacy alias |
| `ecs_cluster` | `ecsc` | `regional` | Direct + tags |
| `ecs_service` | `ecss` | `regional` | Direct + tags |
| `ecs_task` | `ecst` | `regional` | Legacy alias |
| `efs` | `efs` | `regional` | Legacy alias |
| `eks` | `eks` | `regional` | Legacy alias |
| `eks_cluster` | `eksc` | `regional` | Direct + tags |
| `eks_node_group` | `ekng` | `regional` | Direct + tags |
| `elastic_ip` | `elip` | `regional` | Legacy alias |
| `elasticache` | `elch` | `regional` | Legacy alias |
| `elasticsearch` | `elsr` | `regional` | Legacy alias |
| `elb` | `elbl` | `regional` | Direct + tags |
| `emr` | `emr` | `regional` | Service alias |
| `emr_cluster` | `emrc` | `regional` | Direct + tags |
| `emr_instance_fleet` | `emrif` | `regional` | Direct |
| `emr_instance_group` | `emrig` | `regional` | Direct |
| `emr_security_configuration` | `emrsc` | `regional` | Direct |
| `emr_studio` | `emrst` | `regional` | Direct + tags |
| `emrcontainers_job_template` | `emrcj` | `regional` | Direct + tags |
| `emrcontainers_virtual_cluster` | `emrcv` | `regional` | Direct + tags |
| `emrserverless_application` | `emrsa` | `regional` | Direct + tags |
| `eventbridge_bus` | `evbb` | `regional` | Legacy alias |
| `eventbridge_rule` | `evbr` | `regional` | Legacy alias |
| `glue` | `glue` | `regional` | Legacy alias |
| `glue_catalog` | `glcat` | `regional` | Direct + tags |
| `glue_catalog_database` | `glcdb` | `regional` | Direct + tags |
| `glue_catalog_table` | `glctb` | `regional` | Direct |
| `glue_classifier` | `glclf` | `regional` | Direct |
| `glue_connection` | `glcon` | `regional` | Direct + tags |
| `glue_crawler` | `glcrw` | `regional` | Direct + tags |
| `glue_data_quality_ruleset` | `gldqr` | `regional` | Direct + tags |
| `glue_dev_endpoint` | `gldev` | `regional` | Direct + tags |
| `glue_job` | `gljob` | `regional` | Direct + tags |
| `glue_ml_transform` | `glmlt` | `regional` | Direct + tags |
| `glue_partition_index` | `glpix` | `regional` | Direct |
| `glue_registry` | `glreg` | `regional` | Direct + tags |
| `glue_schema` | `glsch` | `regional` | Direct + tags |
| `glue_security_configuration` | `glsec` | `regional` | Direct |
| `glue_trigger` | `gltrg` | `regional` | Direct + tags |
| `glue_user_defined_function` | `gludf` | `regional` | Direct |
| `glue_workflow` | `glwfl` | `regional` | Direct + tags |
| `guardduty` | `gdty` | `regional` | Legacy alias |
| `iam_group` | `iamg` | `regional` | Direct |
| `iam_policy` | `iamp` | `regional` | Direct + tags |
| `iam_role` | `role` | `regional` | Direct + tags |
| `iam_user` | `iamu` | `regional` | Direct + tags |
| `igw` | `igtw` | `regional` | Legacy alias |
| `kinesis` | `knss` | `regional` | Service alias |
| `kinesis_analytics_application` | `knsaa` | `regional` | Direct + tags |
| `kinesis_firehose_delivery_stream` | `knsfh` | `regional` | Direct + tags |
| `kinesis_stream` | `knsst` | `regional` | Direct + tags |
| `kinesis_stream_consumer` | `knssc` | `regional` | Direct |
| `kinesis_video_stream` | `knsvs` | `regional` | Direct + tags |
| `kinesisanalyticsv2_application` | `knsv2` | `regional` | Direct + tags |
| `kinesisanalyticsv2_application_snapshot` | `kns2s` | `regional` | Direct |
| `kms_key` | `kmsk` | `regional` | Tags |
| `lakeformation` | `lkfm` | `regional` | Service alias |
| `lakeformation_data_cells_filter` | `lfdcf` | `regional` | Direct |
| `lakeformation_lf_tag` | `lftag` | `regional` | Direct |
| `lakeformation_lf_tag_expression` | `lftge` | `regional` | Direct |
| `lambda` | `lmbd` | `regional` | Legacy alias |
| `launch_template` | `lcht` | `regional` | Direct + tags |
| `log_group` | `logg` | `regional` | Legacy alias |
| `msk_cluster` | `mskc` | `regional` | Direct + tags |
| `nacl` | `nacl` | `regional` | Legacy alias |
| `nat_gw` | `ngtw` | `regional` | Legacy alias |
| `nlb` | `nlbl` | `regional` | Legacy alias |
| `opensearch` | `opsr` | `regional` | Legacy alias |
| `quicksight` | `qkst` | `regional` | Service alias |
| `quicksight_account_subscription` | `qsasu` | `regional` | Direct |
| `quicksight_analysis` | `qsana` | `regional` | Direct + tags |
| `quicksight_custom_permissions` | `qscpm` | `regional` | Direct + tags |
| `quicksight_dashboard` | `qsdsh` | `regional` | Direct + tags |
| `quicksight_data_set` | `qsds` | `regional` | Direct + tags |
| `quicksight_data_source` | `qsdsr` | `regional` | Direct + tags |
| `quicksight_folder` | `qsfld` | `regional` | Direct + tags |
| `quicksight_group` | `qsgrp` | `regional` | Direct |
| `quicksight_iam_policy_assignment` | `qsipa` | `regional` | Direct |
| `quicksight_ingestion` | `qsing` | `regional` | Direct |
| `quicksight_namespace` | `qsns` | `regional` | Tags |
| `quicksight_refresh_schedule` | `qsref` | `regional` | Direct |
| `quicksight_template` | `qstpl` | `regional` | Direct + tags |
| `quicksight_template_alias` | `qstal` | `regional` | Direct |
| `quicksight_theme` | `qsthm` | `regional` | Direct + tags |
| `quicksight_user` | `qsusr` | `regional` | Direct |
| `quicksight_vpc_connection` | `qsvpc` | `regional` | Direct + tags |
| `rds` | `rds` | `regional` | Legacy alias |
| `rds_cluster` | `rdsc` | `regional` | Direct + tags |
| `redshift` | `rdsh` | `regional` | Legacy alias |
| `redshift_authentication_profile` | `rsapf` | `regional` | Direct |
| `redshift_cluster` | `rscl` | `regional` | Direct + tags |
| `redshift_cluster_snapshot` | `rscss` | `regional` | Direct + tags |
| `redshift_endpoint_access` | `rsepa` | `regional` | Direct |
| `redshift_event_subscription` | `rsesb` | `regional` | Direct + tags |
| `redshift_hsm_client_certificate` | `rshcc` | `regional` | Direct + tags |
| `redshift_hsm_configuration` | `rshcf` | `regional` | Direct + tags |
| `redshift_idc_application` | `rsidc` | `regional` | Direct |
| `redshift_integration` | `rsint` | `regional` | Direct + tags |
| `redshift_parameter_group` | `rspg` | `regional` | Direct + tags |
| `redshift_scheduled_action` | `rssca` | `regional` | Direct |
| `redshift_snapshot_copy_grant` | `rsscg` | `regional` | Direct + tags |
| `redshift_snapshot_schedule` | `rsssh` | `regional` | Tags |
| `redshift_subnet_group` | `rssng` | `regional` | Direct + tags |
| `redshift_usage_limit` | `rsusg` | `regional` | Direct + tags |
| `redshiftdata_statement` | `rsdst` | `regional` | Direct |
| `redshiftserverless_endpoint_access` | `rssea` | `regional` | Direct |
| `redshiftserverless_namespace` | `rssns` | `regional` | Direct + tags |
| `redshiftserverless_snapshot` | `rsssn` | `regional` | Direct |
| `redshiftserverless_workgroup` | `rsswg` | `regional` | Direct + tags |
| `role` | `role` | `regional` | Legacy alias |
| `role_policy` | `rlpl` | `regional` | Legacy alias |
| `route53_record` | `r53r` | `regional` | Direct |
| `route53_zone` | `rt53` | `regional` | Direct + tags |
| `route_table` | `rttb` | `regional` | Tags |
| `s3` | `s3b` | `regional` | Legacy alias |
| `s3_access_point` | `s3ap` | `regional` | Direct + tags |
| `s3_bucket` | `s3bk` | `regional` | Tags |
| `s3_dir` | `s3dr` | `regional` | Legacy alias |
| `s3_object` | `s3ob` | `regional` | Direct + tags |
| `s3_table` | `s3tb` | `regional` | Legacy alias |
| `sagemaker` | `sgmk` | `regional` | Legacy alias |
| `sagemaker_algorithm` | `sgalg` | `regional` | Direct + tags |
| `sagemaker_app` | `sgapp` | `regional` | Direct + tags |
| `sagemaker_app_image_config` | `sgimc` | `regional` | Direct + tags |
| `sagemaker_code_repository` | `sgcdr` | `regional` | Direct + tags |
| `sagemaker_data_quality_job_definition` | `sgdqj` | `regional` | Direct + tags |
| `sagemaker_device` | `sgdev` | `regional` | Direct |
| `sagemaker_device_fleet` | `sgdfl` | `regional` | Direct + tags |
| `sagemaker_domain` | `sgdmn` | `regional` | Direct + tags |
| `sagemaker_endpoint` | `sgend` | `regional` | Direct + tags |
| `sagemaker_endpoint_configuration` | `sgecf` | `regional` | Direct + tags |
| `sagemaker_feature_group` | `sgfgr` | `regional` | Direct + tags |
| `sagemaker_flow_definition` | `sgfld` | `regional` | Direct + tags |
| `sagemaker_hub` | `sghub` | `regional` | Direct + tags |
| `sagemaker_hub_content_reference` | `sghcr` | `regional` | Direct + tags |
| `sagemaker_human_task_ui` | `sghtu` | `regional` | Direct + tags |
| `sagemaker_hyper_parameter_tuning_job` | `sghtj` | `regional` | Direct + tags |
| `sagemaker_image` | `sgimg` | `regional` | Direct + tags |
| `sagemaker_labeling_job` | `sglbj` | `regional` | Direct + tags |
| `sagemaker_mlflow_app` | `sgmfa` | `regional` | Direct + tags |
| `sagemaker_mlflow_tracking_server` | `sgmft` | `regional` | Direct + tags |
| `sagemaker_model` | `sgmdl` | `regional` | Direct + tags |
| `sagemaker_model_card` | `sgmcd` | `regional` | Direct + tags |
| `sagemaker_model_card_export_job` | `sgmce` | `regional` | Direct |
| `sagemaker_model_package_group` | `sgmpg` | `regional` | Direct + tags |
| `sagemaker_monitoring_schedule` | `sgmon` | `regional` | Direct + tags |
| `sagemaker_notebook_instance` | `sgnbi` | `regional` | Direct + tags |
| `sagemaker_notebook_instance_lifecycle_configuration` | `sgnlc` | `regional` | Direct + tags |
| `sagemaker_pipeline` | `sgppl` | `regional` | Direct + tags |
| `sagemaker_project` | `sgprj` | `regional` | Direct + tags |
| `sagemaker_space` | `sgspc` | `regional` | Direct + tags |
| `sagemaker_studio_lifecycle_config` | `sgslc` | `regional` | Direct + tags |
| `sagemaker_training_job` | `sgtrn` | `regional` | Direct + tags |
| `sagemaker_user_profile` | `sgusr` | `regional` | Direct + tags |
| `sagemaker_workforce` | `sgwkf` | `regional` | Direct |
| `sagemaker_workteam` | `sgwkt` | `regional` | Direct + tags |
| `sec_group` | `scgp` | `regional` | Legacy alias |
| `secretsmanager_secret` | `smse` | `regional` | Direct + tags |
| `sfn` | `stfn` | `regional` | Legacy alias |
| `snow_notification_integration` | `snti` | `regional` | Legacy alias |
| `sns` | `sns` | `regional` | Legacy alias |
| `sqs` | `sqs` | `regional` | Legacy alias |
| `ssm_parameter` | `ssmp` | `regional` | Direct + tags |
| `step_function` | `stfn` | `regional` | Legacy alias |
| `subnet` | `subn` | `regional` | Tags |
| `target_group` | `tgpt` | `regional` | Legacy alias |
| `vpc` | `vpcn` | `regional` | Tags |
| `wafv2_ip_set` | `wfis` | `regional` | Direct + tags |
| `wafv2_web_acl` | `wfac` | `regional` | Direct + tags |
| `wafv2_web_acl_rule` | `wfar` | `regional` | Direct |

### Canonical Terraform Aliases

These canonical Terraform keys resolve to established legacy acronyms. The legacy keys remain valid and unchanged.

| Canonical key | Legacy key | Acronym |
| --- | --- | --- |
| `aws_acm_certificate` | `acm_cert` | `acmc` |
| `aws_apigatewayv2_api` | `api_gateway_v2` | `agv2` |
| `aws_appsync_graphql_api` | `appsync` | `apsy` |
| `aws_cloudfront_distribution` | `cloudfront` | `clfr` |
| `aws_cloudwatch_event_bus` | `eventbridge_bus` | `evbb` |
| `aws_cloudwatch_event_rule` | `eventbridge_rule` | `evbr` |
| `aws_cloudwatch_metric_alarm` | `cloudwatch_alarm` | `cwal` |
| `aws_codebuild_project` | `codebuild` | `cdbd` |
| `aws_codedeploy_app` | `codedeploy` | `cddp` |
| `aws_config_config_rule` | `config_rule` | `cfrl` |
| `aws_ebs_volume` | `ebs` | `ebs` |
| `aws_ecr_repository` | `ecr` | `ecr` |
| `aws_ecs_task_definition` | `ecs_task` | `ecst` |
| `aws_efs_file_system` | `efs` | `efs` |
| `aws_eip` | `elastic_ip` | `elip` |
| `aws_elasticache_cluster` | `elasticache` | `elch` |
| `aws_elasticsearch_domain` | `elasticsearch` | `elsr` |
| `aws_iam_role_policy` | `role_policy` | `rlpl` |
| `aws_instance` | `ec2_instance` | `ec2i` |
| `aws_internet_gateway` | `igw` | `igtw` |
| `aws_lambda_function` | `lambda` | `lmbd` |
| `aws_lb_target_group` | `target_group` | `tgpt` |
| `aws_nat_gateway` | `nat_gw` | `ngtw` |
| `aws_network_acl` | `nacl` | `nacl` |
| `aws_opensearch_domain` | `opensearch` | `opsr` |
| `aws_s3_directory_bucket` | `s3_dir` | `s3dr` |
| `aws_s3tables_table` | `s3_table` | `s3tb` |
| `aws_security_group` | `sec_group` | `scgp` |
| `aws_sfn_state_machine` | `sfn` | `stfn` |
| `aws_sns_topic` | `sns` | `sns` |
| `aws_sqs_queue` | `sqs` | `sqs` |

### Excluded Non-Nameable Resources

The following Terraform resources are covered by the audit but intentionally do not receive acronyms. They configure, attach, authorize, register, or version another resource and expose neither an independent name nor standard AWS tags.

| Terraform resource | Reason |
| --- | --- |
| `aws_bedrock_foundation_model_agreement` | service configuration |
| `aws_bedrock_guardrail_version` | generated version |
| `aws_bedrock_model_invocation_logging_configuration` | service configuration |
| `aws_bedrock_use_case_for_model_access` | service configuration |
| `aws_bedrockagent_agent_knowledge_base_association` | association |
| `aws_bedrockagentcore_resource_policy` | resource policy |
| `aws_bedrockagentcore_token_vault_cmk` | service configuration |
| `aws_datazone_environment_blueprint_configuration` | service configuration |
| `aws_emr_block_public_access_configuration` | service configuration |
| `aws_emr_managed_scaling_policy` | attached policy |
| `aws_emr_studio_session_mapping` | association |
| `aws_glue_catalog_table_optimizer` | service configuration |
| `aws_glue_data_catalog_encryption_settings` | service configuration |
| `aws_glue_partition` | generated identifier |
| `aws_glue_resource_policy` | resource policy |
| `aws_kinesis_account_settings` | service configuration |
| `aws_kinesis_resource_policy` | resource policy |
| `aws_lakeformation_data_lake_settings` | service configuration |
| `aws_lakeformation_identity_center_configuration` | service configuration |
| `aws_lakeformation_opt_in` | association |
| `aws_lakeformation_permissions` | permission grant |
| `aws_lakeformation_resource` | resource registration |
| `aws_lakeformation_resource_lf_tag` | association |
| `aws_lakeformation_resource_lf_tags` | association |
| `aws_quicksight_account_settings` | service configuration |
| `aws_quicksight_folder_membership` | association |
| `aws_quicksight_group_membership` | association |
| `aws_quicksight_ip_restriction` | service configuration |
| `aws_quicksight_key_registration` | resource registration |
| `aws_quicksight_role_custom_permission` | association |
| `aws_quicksight_role_membership` | association |
| `aws_quicksight_user_custom_permission` | association |
| `aws_redshift_cluster_iam_roles` | association |
| `aws_redshift_data_share_authorization` | authorization |
| `aws_redshift_data_share_consumer_association` | association |
| `aws_redshift_endpoint_authorization` | authorization |
| `aws_redshift_logging` | service configuration |
| `aws_redshift_namespace_registration` | resource registration |
| `aws_redshift_partner` | association |
| `aws_redshift_resource_policy` | resource policy |
| `aws_redshift_snapshot_copy` | service configuration |
| `aws_redshift_snapshot_schedule_association` | association |
| `aws_redshiftserverless_custom_domain_association` | association |
| `aws_redshiftserverless_resource_policy` | resource policy |
| `aws_redshiftserverless_usage_limit` | service configuration |
| `aws_sagemaker_image_version` | generated version |
| `aws_sagemaker_model_package_group_policy` | resource policy |
| `aws_sagemaker_servicecatalog_portfolio_status` | service configuration |

### AWS Constraint Policy

All 153 newly supported exact resources have minimum length, maximum length, and character-pattern validation. Service-family defaults are conservative subsets of the upstream AWS/provider rules. Resource-specific overrides enforce stricter limits where needed:

- Athena databases and prepared statements use `straight` style because their identifiers do not accept hyphens.
- Bedrock AgentCore identifiers with letter/alphanumeric/underscore rules use `straight` style and 40-48 character limits.
- DataZone form types use Smithy-compatible `straight` identifiers with a 36-character limit.
- Redshift endpoint names are limited to 30 lowercase characters; cluster-related identifiers reject trailing or consecutive hyphens.
- SageMaker projects and hyperparameter tuning jobs are limited to 32 characters; other SageMaker names use a conservative 63-character limit.

## Azure CAF Acronyms and Constraints

For `cloud = "azure"`, Sigil loads **all Azure CAF resource types** from `resourceDefinition.json` and applies:
- CAF acronyms from the Azure CAF resource catalog.
- Per-resource min/max/regex constraints.
- Per-resource style allowances derived from CAF dash/lowercase metadata.

Comprehensive reference (395 resource types):
- `azure-caf-resources.md`
- CAF resource catalog JSON: https://github.com/aztfmod/terraform-provider-azurecaf/blob/main/resourceDefinition.json
- Azure naming rules: https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules
- CAF abbreviations: https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations

### Supported Azure Resources and Acronyms

Supported Azure `what` values are the Azure CAF resource identifiers listed in `azure-caf-resources.md`. The `Acronym` column in that table is the value returned by `resource_acronym`.

Quick reference:

| Azure Resource (`what`) | Acronym |
| --- | --- |
| `azurerm_resource_group` | `rg` |
| `azurerm_storage_account` | `st` |
| `azurerm_virtual_network` | `vnet` |
| `azurerm_subnet` | `snet` |
| `azurerm_kubernetes_cluster` | `aks` |
| `azurerm_container_registry` | `cr` |
| `azurerm_key_vault` | `kv` |
| `azurerm_linux_virtual_machine` | `vm` |

Sigil uses CAF acronyms directly by default. Use `resource_acronyms` only when you need explicit overrides.

For the complete list of all 395 supported Azure resources and acronyms, see `azure-caf-resources.md`.

## GCP Coverage and Strategy

`cloud = "gcp"` has broad built-in coverage. Unlike Azure CAF, Google Cloud does not provide a single official catalog that includes all Terraform resource identifiers, acronyms, scopes, and naming regex rules in one place.

Sigil accepts either Terraform-style GCP resource identifiers such as `google_compute_network` or normalized keys such as `compute_network`. The built-in GCP defaults are stored in normalized form and the optional `google_` prefix is resolved at lookup time.

### Why GCP Needs a Different Approach

- GCP naming rules are mostly per-service, not centralized.
- Many resources are identified by fully-qualified paths (`projects/.../locations/.../...`) or server-generated IDs.
- Some resources have both a technical identifier and a user-facing `display_name`, which need different handling.

### Nameability Tiers

Classify each `what` resource into one of these tiers:

| Tier | Meaning | Sigil Behavior |
| --- | --- | --- |
| `tier_a_named` | Resource has a real user-controlled identifier (`name`, `bucket`, `project_id`, `account_id`, etc.) with documented constraints. | Full acronym + style + strict constraints (min/max/regex/forbidden patterns). |
| `tier_b_display` | Primary identity is path-like or composite, but resource exposes `display_name`/labels for human naming. | Acronym + style only by default; no hard validation unless an explicit documented constraint exists. |
| `tier_c_opaque` | No stable user-defined name (provider/API generated IDs, bindings/memberships, attachment resources). | No strict naming profile; resource should not be targeted for canonical Sigil naming. |

### Constraint Policy

Apply strict validation only where deterministic and well-documented:

1. `strict` for `tier_a_named` resources with authoritative naming rules.
2. `best_effort` for `tier_b_display` resources (formatting consistency, usually no hard fail).
3. `none` for `tier_c_opaque` resources.

This avoids false failures on resources that are not truly user-nameable.

### Current Coverage

Strict constraints are enabled for:

- `google_storage_bucket` (plus aliases `gcs_bucket`, `gcs`)
- named Compute Engine resources in the default GCP map, including:
  `google_compute_network` / `vpc`,
  `google_compute_subnetwork` / `subnet`,
  `google_compute_router`,
  `google_compute_firewall`,
  `google_compute_address`,
  `google_compute_global_address`,
  `google_compute_route`,
  `google_compute_router_nat`,
  `google_compute_vpn_gateway`,
  `google_compute_vpn_tunnel`,
  `google_compute_ha_vpn_gateway`,
  `google_compute_url_map`,
  `google_compute_target_http_proxy`,
  `google_compute_target_https_proxy`,
  `google_compute_backend_service`,
  `google_compute_region_backend_service`,
  `google_compute_instance_template`,
  `google_compute_instance_group_manager`,
  `google_compute_region_instance_group_manager`,
  `google_compute_disk`,
  `google_compute_image`,
  `google_compute_snapshot`
- `google_pubsub_topic`, `google_pubsub_subscription`
- `google_service_account`
- `google_bigquery_dataset`
- `google_cloud_run_v2_service`

Built-in acronyms also cover additional common GCP resources such as:

- `google_artifact_registry_repository`
- `google_compute_router`, `google_compute_firewall`, `google_compute_address`, `google_compute_global_address`
- `google_compute_route`, `google_compute_router_nat`
- `google_compute_vpn_gateway`, `google_compute_vpn_tunnel`, `google_compute_ha_vpn_gateway`
- `google_compute_url_map`, `google_compute_target_http_proxy`, `google_compute_target_https_proxy`
- `google_compute_backend_service`, `google_compute_region_backend_service`
- `google_compute_instance_template`, `google_compute_instance_group_manager`, `google_compute_region_instance_group_manager`
- `google_compute_disk`, `google_compute_image`, `google_compute_snapshot`
- `google_dns_managed_zone`
- `google_secret_manager_secret`
- `google_kms_key_ring`, `google_kms_crypto_key`
- `google_sql_database_instance`
- `google_container_cluster`, `google_container_node_pool`
- `google_vpc_access_connector`
- `google_redis_instance`, `google_memcache_instance`, `google_filestore_instance`
- `google_spanner_instance`, `google_spanner_database`
- `google_cloudbuild_trigger`, `google_eventarc_trigger`
- `google_cloud_scheduler_job`, `google_cloud_tasks_queue`, `google_workflows_workflow`
- `google_monitoring_notification_channel`
- `google_logging_metric`, `google_logging_project_sink`
- `google_pubsub_schema`, `google_pubsub_snapshot`

GCP resources outside the built-in map remain permissive by default (resource key fallback, default style handling, and no hard constraints).

### Expansion Plan

1. Add constraints resource-family by resource-family, only when naming rules are explicit and stable.
2. Keep path/ID-based resources in `tier_b_display` or `tier_c_opaque` mode by default.
3. Add tests for each new constrained resource before adding it to defaults.

### Remaining Tier-A Review Set

The current open Tier-A review item is:

- `google_artifact_registry_repository`

It already has a stable acronym, but strict validation should only be added once we pin an authoritative Google naming rule source for `repository_id`.

### Data Sources for Coverage

Use multiple inputs, because no GCP equivalent to Azure CAF exists:

- Cloud Asset Inventory asset type list for broad resource inventory.
- Terraform Google provider resource schemas for argument names and shape.
- Service-specific Google Cloud documentation for authoritative naming constraints.

### Definition of Done for GCP Support

- Every supported GCP `what` is tagged with a tier.
- Only `tier_a_named` resources enforce hard constraints.
- Docs list supported GCP resources and constraint source for each.
- Tests cover acronym resolution, style filtering, and constraint behavior for representative resources in each tier.

## Naming Styles

Style priority determines how names are formatted. If a resource has style constraints, the provider selects the first allowed style in the priority list.

Valid styles:
- `dashed`
- `underscore`
- `straight`
- `pascal`
- `pascaldashed`
- `camel`

Style behaviors:
- `dashed` Lowercase words joined by `-`.
- `underscore` Lowercase words joined by `_`.
- `straight` Lowercase words concatenated.
- `pascal` Words in `PascalCase`.
- `pascaldashed` Words in `Pascal-Case` joined by `-`.
- `camel` Words in `camelCase`.

Words are extracted from each component using the pattern `[A-Za-z0-9]+`, so punctuation or separators become word boundaries. If no valid style matches, Sigil falls back to the first allowed style from `resource_style_overrides` for that resource, or `dashed` when no style override exists.

Cloud-specific style overrides are applied automatically:
- `aws`: S3 uses compatible dashed/straight styles; Athena, Bedrock AgentCore, and DataZone resources with identifier-only character sets automatically use `straight`.
- `azure`: each CAF resource inherits style limits from CAF dash/lowercase metadata.
- `gcp`: built-in constrained resources include style restrictions for compatibility with Google Cloud naming rules, including bucket, Compute Engine, service account, BigQuery dataset, Pub/Sub, and Cloud Run resources.

## Resource Constraints

Some resources have naming constraints enforced after formatting. The constraint name matches the `what` input (case-insensitive).

The table below lists built-in `aws` constraints. Azure constraints are listed in `azure-caf-resources.md`. GCP built-in constraints currently cover `google_storage_bucket`, the named Compute Engine resources included in the default GCP acronym map, `google_pubsub_topic`, `google_pubsub_subscription`, `google_service_account`, `google_bigquery_dataset`, and `google_cloud_run_v2_service` (including their aliases).

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

Constraint types include minimum or maximum length, required pattern, forbidden prefixes or suffixes, forbidden substrings, forbidden regex patterns, and checks that the name is not formatted as an IPv4 address.

## Argument Reference

- `config` (Optional) Base configuration object; accepts the same keys as the top-level attributes.
- `overrides` (Optional) Overrides applied after top-level attributes; accepts the same keys as the top-level attributes.
- `cloud` (Optional) Cloud naming profile. Supported values are `aws` (default), `azure`, and `gcp`.
- `org_prefix` (Required unless set in `config` or `overrides`) Short organization identifier.
- `project` (Optional) Project or workload identifier.
- `env` (Required unless set in `config` or `overrides`) Environment identifier, such as `dev`, `staging`, or `prod`.
- `region` (Optional) Cloud region name, used to derive a short region code. If no `region_map` entry exists, the raw region value is used.
- `region_short_code` (Optional) Explicit short region code to use instead of mapping.
- `region_map` (Optional) Full region map; when set, replaces the default map.
- `region_overrides` (Optional) Map of region overrides applied on top of the default map.
- `ignore_region_for_regional_resources` (Optional) When `true` (default), omit the region component for resources marked as `regional` in the acronyms tables.
- `recipe` (Optional) Ordered list of components used to build the name.
- `style_priority` (Optional) Preferred naming styles in order of precedence.
- `resource_acronyms` (Optional) Map of resource identifiers to acronyms.
- `resource_style_overrides` (Optional) Map of resource identifiers to allowed styles.

## Notes

Default recipe components are `org`, `proj`, `env`, `region`, `resource`, and `qualifier`. Components are only included when non-empty, and you can omit them by removing items from the recipe.

If both `region_map` and `region_overrides` are set, overrides are applied to the map.

When `ignore_region_for_regional_resources` is `true`, the region component is omitted for resources marked as regional unless explicitly overridden.

`cloud = "azure"` loads Azure CAF resource defaults (acronyms, style rules, and regex constraints) from `resourceDefinition.json`, plus a built-in Azure region short code map.

`cloud = "gcp"` loads built-in GCP defaults (region short codes, resource acronyms, style rules, and strict constraints for storage bucket, named Compute Engine resources in the default acronym map, Pub/Sub, service account, BigQuery dataset, and Cloud Run service resources).
