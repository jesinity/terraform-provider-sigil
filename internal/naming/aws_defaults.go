package naming

import (
	"regexp"
	"strings"
)

func DefaultRegionMap() map[string]string {
	return map[string]string{
		"us-east-1":      "use1",
		"us-east-2":      "use2",
		"us-west-1":      "usw1",
		"us-west-2":      "usw2",
		"af-south-1":     "afs1",
		"ap-east-1":      "ape1",
		"ap-south-1":     "aps1",
		"ap-south-2":     "aps2",
		"ap-southeast-1": "apse1",
		"ap-southeast-2": "apse2",
		"ap-southeast-3": "apse3",
		"ap-southeast-4": "apse4",
		"ap-northeast-1": "apne1",
		"ap-northeast-2": "apne2",
		"ap-northeast-3": "apne3",
		"ca-central-1":   "cac1",
		"ca-west-1":      "caw1",
		"cn-north-1":     "cnn1",
		"cn-northwest-1": "cnnw1",
		"eu-central-1":   "euc1",
		"eu-central-2":   "euc2",
		"eu-west-1":      "euw1",
		"eu-west-2":      "euw2",
		"eu-west-3":      "euw3",
		"eu-west-4":      "euw4",
		"eu-north-1":     "eun1",
		"eu-south-1":     "eus1",
		"eu-south-2":     "eus2",
		"il-central-1":   "ilc1",
		"me-south-1":     "mes1",
		"me-central-1":   "mec1",
		"sa-east-1":      "sae1",
		"us-gov-west-1":  "usgw1",
		"us-gov-east-1":  "usge1",
	}
}

func DefaultResourceAcronyms() map[string]string {
	return map[string]string{
		"role":                                  "role",
		"role_policy":                           "rlpl",
		"iam_role":                              "role",
		"iam_policy":                            "iamp",
		"iam_user":                              "iamu",
		"iam_group":                             "iamg",
		"s3":                                    "s3b",
		"s3_bucket":                             "s3bk",
		"s3_object":                             "s3ob",
		"s3_access_point":                       "s3ap",
		"s3_table":                              "s3tb",
		"s3_dir":                                "s3dr",
		"sns":                                   "sns",
		"sqs":                                   "sqs",
		"ecs_cluster":                           "ecsc",
		"ecs_service":                           "ecss",
		"ecs_task":                              "ecst",
		"eks":                                   "eks",
		"eks_cluster":                           "eksc",
		"eks_node_group":                        "ekng",
		"msk_cluster":                           "mskc",
		"vpc":                                   "vpcn",
		"subnet":                                "subn",
		"igw":                                   "igtw",
		"nat_gw":                                "ngtw",
		"sec_group":                             "scgp",
		"nacl":                                  "nacl",
		"route_table":                           "rttb",
		"elastic_ip":                            "elip",
		"wafv2_web_acl":                         "wfac",
		"wafv2_web_acl_rule":                    "wfar",
		"wafv2_ip_set":                          "wfis",
		"lambda":                                "lmbd",
		"api_gateway_rest_api":                  "agra",
		"api_gateway_model":                     "agmd",
		"api_gateway_v2":                        "agv2",
		"log_group":                             "logg",
		"cloudwatch_log_group":                  "cwlg",
		"cloudwatch_alarm":                      "cwal",
		"eventbridge_bus":                       "evbb",
		"eventbridge_rule":                      "evbr",
		"step_function":                         "stfn",
		"sfn":                                   "stfn",
		"dynamodb":                              "dydb",
		"dynamodb_table":                        "dybt",
		"rds":                                   "rds",
		"rds_cluster":                           "rdsc",
		"aurora_cluster":                        "arcl",
		"redshift":                              "rdsh",
		"redshift_authentication_profile":       "rsapf",
		"redshift_cluster":                      "rscl",
		"redshift_cluster_snapshot":             "rscss",
		"redshift_endpoint_access":              "rsepa",
		"redshift_event_subscription":           "rsesb",
		"redshift_hsm_client_certificate":       "rshcc",
		"redshift_hsm_configuration":            "rshcf",
		"redshift_idc_application":              "rsidc",
		"redshift_integration":                  "rsint",
		"redshift_parameter_group":              "rspg",
		"redshift_scheduled_action":             "rssca",
		"redshift_snapshot_copy_grant":          "rsscg",
		"redshift_snapshot_schedule":            "rsssh",
		"redshift_subnet_group":                 "rssng",
		"redshift_usage_limit":                  "rsusg",
		"redshiftdata_statement":                "rsdst",
		"redshiftserverless_endpoint_access":    "rssea",
		"redshiftserverless_namespace":          "rssns",
		"redshiftserverless_snapshot":           "rsssn",
		"redshiftserverless_workgroup":          "rsswg",
		"elasticache":                           "elch",
		"opensearch":                            "opsr",
		"elasticsearch":                         "elsr",
		"ecr":                                   "ecr",
		"ecs":                                   "ecs",
		"ec2_instance":                          "ec2i",
		"launch_template":                       "lcht",
		"autoscaling_group":                     "asgr",
		"alb":                                   "albl",
		"nlb":                                   "nlbl",
		"elb":                                   "elbl",
		"target_group":                          "tgpt",
		"cloudfront":                            "clfr",
		"route53_zone":                          "rt53",
		"route53_record":                        "r53r",
		"acm_cert":                              "acmc",
		"kms_key":                               "kmsk",
		"secretsmanager_secret":                 "smse",
		"ssm_parameter":                         "ssmp",
		"cloudtrail":                            "ctra",
		"guardduty":                             "gdty",
		"config_rule":                           "cfrl",
		"efs":                                   "efs",
		"ebs":                                   "ebs",
		"athena":                                "athn",
		"athena_capacity_reservation":           "athcr",
		"athena_data_catalog":                   "athdc",
		"athena_database":                       "athdb",
		"athena_named_query":                    "athnq",
		"athena_prepared_statement":             "athps",
		"athena_workgroup":                      "athwg",
		"glue":                                  "glue",
		"glue_catalog":                          "glcat",
		"glue_catalog_database":                 "glcdb",
		"glue_catalog_table":                    "glctb",
		"glue_classifier":                       "glclf",
		"glue_connection":                       "glcon",
		"glue_crawler":                          "glcrw",
		"glue_data_quality_ruleset":             "gldqr",
		"glue_dev_endpoint":                     "gldev",
		"glue_job":                              "gljob",
		"glue_ml_transform":                     "glmlt",
		"glue_partition_index":                  "glpix",
		"glue_registry":                         "glreg",
		"glue_schema":                           "glsch",
		"glue_security_configuration":           "glsec",
		"glue_trigger":                          "gltrg",
		"glue_user_defined_function":            "gludf",
		"glue_workflow":                         "glwfl",
		"sagemaker":                             "sgmk",
		"sagemaker_algorithm":                   "sgalg",
		"sagemaker_app":                         "sgapp",
		"sagemaker_app_image_config":            "sgimc",
		"sagemaker_code_repository":             "sgcdr",
		"sagemaker_data_quality_job_definition": "sgdqj",
		"sagemaker_device":                      "sgdev",
		"sagemaker_device_fleet":                "sgdfl",
		"sagemaker_domain":                      "sgdmn",
		"sagemaker_endpoint":                    "sgend",
		"sagemaker_endpoint_configuration":      "sgecf",
		"sagemaker_feature_group":               "sgfgr",
		"sagemaker_flow_definition":             "sgfld",
		"sagemaker_hub":                         "sghub",
		"sagemaker_hub_content_reference":       "sghcr",
		"sagemaker_human_task_ui":               "sghtu",
		"sagemaker_hyper_parameter_tuning_job":  "sghtj",
		"sagemaker_image":                       "sgimg",
		"sagemaker_labeling_job":                "sglbj",
		"sagemaker_mlflow_app":                  "sgmfa",
		"sagemaker_mlflow_tracking_server":      "sgmft",
		"sagemaker_model":                       "sgmdl",
		"sagemaker_model_card":                  "sgmcd",
		"sagemaker_model_card_export_job":       "sgmce",
		"sagemaker_model_package_group":         "sgmpg",
		"sagemaker_monitoring_schedule":         "sgmon",
		"sagemaker_notebook_instance":           "sgnbi",
		"sagemaker_notebook_instance_lifecycle_configuration": "sgnlc",
		"sagemaker_pipeline":                           "sgppl",
		"sagemaker_project":                            "sgprj",
		"sagemaker_space":                              "sgspc",
		"sagemaker_studio_lifecycle_config":            "sgslc",
		"sagemaker_training_job":                       "sgtrn",
		"sagemaker_user_profile":                       "sgusr",
		"sagemaker_workforce":                          "sgwkf",
		"sagemaker_workteam":                           "sgwkt",
		"bedrock":                                      "bdrk",
		"bedrock_custom_model":                         "brcm",
		"bedrock_evaluation_job":                       "brevj",
		"bedrock_guardrail":                            "brgr",
		"bedrock_inference_profile":                    "brip",
		"bedrock_provisioned_model_throughput":         "brpmt",
		"bedrockagent":                                 "brag",
		"bedrockagent_agent":                           "braga",
		"bedrockagent_agent_action_group":              "bragg",
		"bedrockagent_agent_alias":                     "braal",
		"bedrockagent_agent_collaborator":              "bracl",
		"bedrockagent_data_source":                     "brads",
		"bedrockagent_flow":                            "brafl",
		"bedrockagent_knowledge_base":                  "brakb",
		"bedrockagent_prompt":                          "brapt",
		"bedrockagentcore":                             "brac",
		"bedrockagentcore_agent_runtime":               "bracr",
		"bedrockagentcore_agent_runtime_endpoint":      "brace",
		"bedrockagentcore_api_key_credential_provider": "brcak",
		"bedrockagentcore_browser":                     "brcbr",
		"bedrockagentcore_browser_profile":             "brcbp",
		"bedrockagentcore_code_interpreter":            "brcci",
		"bedrockagentcore_evaluator":                   "brcev",
		"bedrockagentcore_gateway":                     "brcgw",
		"bedrockagentcore_gateway_target":              "brcgt",
		"bedrockagentcore_harness":                     "brchr",
		"bedrockagentcore_memory":                      "brcme",
		"bedrockagentcore_memory_strategy":             "brcms",
		"bedrockagentcore_oauth2_credential_provider":  "brco2",
		"bedrockagentcore_online_evaluation_config":    "brcoe",
		"bedrockagentcore_policy":                      "brcpl",
		"bedrockagentcore_policy_engine":               "brcpe",
		"bedrockagentcore_registry":                    "brcrg",
		"bedrockagentcore_workload_identity":           "brcwi",
		"datazone":                                     "dtzn",
		"datazone_asset_type":                          "dzaty",
		"datazone_domain":                              "dzdmn",
		"datazone_environment":                         "dzenv",
		"datazone_environment_profile":                 "dzepf",
		"datazone_form_type":                           "dzfty",
		"datazone_glossary":                            "dzglo",
		"datazone_glossary_term":                       "dzglt",
		"datazone_project":                             "dzprj",
		"datazone_user_profile":                        "dzusr",
		"emr":                                          "emr",
		"emr_cluster":                                  "emrc",
		"emr_instance_fleet":                           "emrif",
		"emr_instance_group":                           "emrig",
		"emr_security_configuration":                   "emrsc",
		"emr_studio":                                   "emrst",
		"emrcontainers_job_template":                   "emrcj",
		"emrcontainers_virtual_cluster":                "emrcv",
		"emrserverless_application":                    "emrsa",
		"kinesis":                                      "knss",
		"kinesis_analytics_application":                "knsaa",
		"kinesis_firehose_delivery_stream":             "knsfh",
		"kinesis_stream":                               "knsst",
		"kinesis_stream_consumer":                      "knssc",
		"kinesis_video_stream":                         "knsvs",
		"kinesisanalyticsv2_application":               "knsv2",
		"kinesisanalyticsv2_application_snapshot":      "kns2s",
		"lakeformation":                                "lkfm",
		"lakeformation_data_cells_filter":              "lfdcf",
		"lakeformation_lf_tag":                         "lftag",
		"lakeformation_lf_tag_expression":              "lftge",
		"quicksight":                                   "qkst",
		"quicksight_account_subscription":              "qsasu",
		"quicksight_analysis":                          "qsana",
		"quicksight_custom_permissions":                "qscpm",
		"quicksight_dashboard":                         "qsdsh",
		"quicksight_data_set":                          "qsds",
		"quicksight_data_source":                       "qsdsr",
		"quicksight_folder":                            "qsfld",
		"quicksight_group":                             "qsgrp",
		"quicksight_iam_policy_assignment":             "qsipa",
		"quicksight_ingestion":                         "qsing",
		"quicksight_namespace":                         "qsns",
		"quicksight_refresh_schedule":                  "qsref",
		"quicksight_template":                          "qstpl",
		"quicksight_template_alias":                    "qstal",
		"quicksight_theme":                             "qsthm",
		"quicksight_user":                              "qsusr",
		"quicksight_vpc_connection":                    "qsvpc",
		"codebuild":                                    "cdbd",
		"codepipeline":                                 "cdpl",
		"codedeploy":                                   "cddp",
		"cloudformation_stack":                         "cfst",
		"appsync":                                      "apsy",
		"snow_notification_integration":                "snti",
	}
}

// DefaultAWSOpaqueResources lists Terraform resources in the covered data and
// ML services that do not expose an independent name or standard AWS tags.
func DefaultAWSOpaqueResources() map[string]string {
	return map[string]string{
		"bedrock_foundation_model_agreement":             "service configuration",
		"bedrock_guardrail_version":                      "generated version",
		"bedrock_model_invocation_logging_configuration": "service configuration",
		"bedrock_use_case_for_model_access":              "service configuration",
		"bedrockagent_agent_knowledge_base_association":  "association",
		"bedrockagentcore_resource_policy":               "resource policy",
		"bedrockagentcore_token_vault_cmk":               "service configuration",
		"datazone_environment_blueprint_configuration":   "service configuration",
		"emr_block_public_access_configuration":          "service configuration",
		"emr_managed_scaling_policy":                     "attached policy",
		"emr_studio_session_mapping":                     "association",
		"glue_catalog_table_optimizer":                   "service configuration",
		"glue_data_catalog_encryption_settings":          "service configuration",
		"glue_partition":                                 "generated identifier",
		"glue_resource_policy":                           "resource policy",
		"kinesis_account_settings":                       "service configuration",
		"kinesis_resource_policy":                        "resource policy",
		"lakeformation_data_lake_settings":               "service configuration",
		"lakeformation_identity_center_configuration":    "service configuration",
		"lakeformation_opt_in":                           "association",
		"lakeformation_permissions":                      "permission grant",
		"lakeformation_resource":                         "resource registration",
		"lakeformation_resource_lf_tag":                  "association",
		"lakeformation_resource_lf_tags":                 "association",
		"quicksight_account_settings":                    "service configuration",
		"quicksight_folder_membership":                   "association",
		"quicksight_group_membership":                    "association",
		"quicksight_ip_restriction":                      "service configuration",
		"quicksight_key_registration":                    "resource registration",
		"quicksight_role_custom_permission":              "association",
		"quicksight_role_membership":                     "association",
		"quicksight_user_custom_permission":              "association",
		"redshift_cluster_iam_roles":                     "association",
		"redshift_data_share_authorization":              "authorization",
		"redshift_data_share_consumer_association":       "association",
		"redshift_endpoint_authorization":                "authorization",
		"redshift_logging":                               "service configuration",
		"redshift_namespace_registration":                "resource registration",
		"redshift_partner":                               "association",
		"redshift_resource_policy":                       "resource policy",
		"redshift_snapshot_copy":                         "service configuration",
		"redshift_snapshot_schedule_association":         "association",
		"redshiftserverless_custom_domain_association":   "association",
		"redshiftserverless_resource_policy":             "resource policy",
		"redshiftserverless_usage_limit":                 "service configuration",
		"sagemaker_image_version":                        "generated version",
		"sagemaker_model_package_group_policy":           "resource policy",
		"sagemaker_servicecatalog_portfolio_status":      "service configuration",
	}
}

// defaultAWSResourceAliases maps canonical Terraform resource keys to legacy
// Sigil keys. Exact keys still win, so this lookup is backwards compatible.
var defaultAWSResourceAliases = map[string]string{
	"acm_certificate":         "acm_cert",
	"apigatewayv2_api":        "api_gateway_v2",
	"appsync_graphql_api":     "appsync",
	"cloudfront_distribution": "cloudfront",
	"cloudwatch_event_bus":    "eventbridge_bus",
	"cloudwatch_event_rule":   "eventbridge_rule",
	"cloudwatch_metric_alarm": "cloudwatch_alarm",
	"codebuild_project":       "codebuild",
	"codedeploy_app":          "codedeploy",
	"config_config_rule":      "config_rule",
	"ebs_volume":              "ebs",
	"ecr_repository":          "ecr",
	"ecs_task_definition":     "ecs_task",
	"efs_file_system":         "efs",
	"eip":                     "elastic_ip",
	"elasticache_cluster":     "elasticache",
	"elasticsearch_domain":    "elasticsearch",
	"iam_role_policy":         "role_policy",
	"instance":                "ec2_instance",
	"internet_gateway":        "igw",
	"lambda_function":         "lambda",
	"lb_target_group":         "target_group",
	"nat_gateway":             "nat_gw",
	"network_acl":             "nacl",
	"opensearch_domain":       "opensearch",
	"s3_directory_bucket":     "s3_dir",
	"s3tables_table":          "s3_table",
	"security_group":          "sec_group",
	"sfn_state_machine":       "sfn",
	"sns_topic":               "sns",
	"sqs_queue":               "sqs",
}

func DefaultGlobalResources() map[string]bool {
	return map[string]bool{
		"role":           true,
		"role_policy":    true,
		"iam_role":       true,
		"iam_policy":     true,
		"iam_user":       true,
		"iam_group":      true,
		"cloudfront":     true,
		"route53_zone":   true,
		"route53_record": true,
	}
}

func DefaultRegionalResources() map[string]bool {
	regional := map[string]bool{}
	for key := range DefaultResourceAcronyms() {
		regional[key] = true
	}
	for key := range DefaultGlobalResources() {
		delete(regional, key)
	}
	return regional
}

func DefaultResourceStyleOverrides() map[string][]string {
	return map[string][]string{
		"athena_database":                           {StyleStraight},
		"athena_prepared_statement":                 {StyleStraight},
		"bedrockagentcore_agent_runtime":            {StyleStraight},
		"bedrockagentcore_agent_runtime_endpoint":   {StyleStraight},
		"bedrockagentcore_browser":                  {StyleStraight},
		"bedrockagentcore_browser_profile":          {StyleStraight},
		"bedrockagentcore_code_interpreter":         {StyleStraight},
		"bedrockagentcore_evaluator":                {StyleStraight},
		"bedrockagentcore_harness":                  {StyleStraight},
		"bedrockagentcore_memory":                   {StyleStraight},
		"bedrockagentcore_memory_strategy":          {StyleStraight},
		"bedrockagentcore_online_evaluation_config": {StyleStraight},
		"bedrockagentcore_policy":                   {StyleStraight},
		"bedrockagentcore_policy_engine":            {StyleStraight},
		"datazone_form_type":                        {StyleStraight},
		"s3":                                        {StyleDashed, StyleStraight},
		"s3_bucket":                                 {StyleDashed, StyleStraight},
	}
}

func DefaultResourceConstraints() map[string]ResourceConstraint {
	constraints := map[string]ResourceConstraint{
		"s3": {
			MinLen:              3,
			MaxLen:              63,
			Pattern:             regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`),
			PatternDescription:  "lowercase letters, numbers, dots, and hyphens; must start and end with a letter or number",
			ForbiddenPrefixes:   []string{"xn--", "sthree-", "amzn-s3-demo-"},
			ForbiddenSuffixes:   []string{"-s3alias", "--ol-s3"},
			ForbiddenSubstrings: []string{".."},
			DisallowIPAddress:   true,
		},
		"s3_bucket": {
			MinLen:              3,
			MaxLen:              63,
			Pattern:             regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`),
			PatternDescription:  "lowercase letters, numbers, dots, and hyphens; must start and end with a letter or number",
			ForbiddenPrefixes:   []string{"xn--", "sthree-", "amzn-s3-demo-"},
			ForbiddenSuffixes:   []string{"-s3alias", "--ol-s3"},
			ForbiddenSubstrings: []string{".."},
			DisallowIPAddress:   true,
		},
		"role": {
			MinLen:             1,
			MaxLen:             64,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`),
			PatternDescription: "alphanumeric and the following: +=,.@_-",
		},
		"iam_role": {
			MinLen:             1,
			MaxLen:             64,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`),
			PatternDescription: "alphanumeric and the following: +=,.@_-",
		},
		"iam_user": {
			MinLen:             1,
			MaxLen:             64,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`),
			PatternDescription: "alphanumeric and the following: +=,.@_-",
		},
		"iam_group": {
			MinLen:             1,
			MaxLen:             128,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`),
			PatternDescription: "alphanumeric and the following: +=,.@_-",
		},
		"iam_policy": {
			MinLen:             1,
			MaxLen:             128,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`),
			PatternDescription: "alphanumeric and the following: +=,.@_-",
		},
		"role_policy": {
			MinLen:             1,
			MaxLen:             128,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`),
			PatternDescription: "alphanumeric and the following: +=,.@_-",
		},
		"sns": {
			MinLen:             1,
			MaxLen:             256,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9_-]+(\.fifo)?$`),
			PatternDescription: "letters, numbers, underscores, and hyphens; FIFO topics must end with .fifo",
		},
		"sns_topic": {
			MinLen:             1,
			MaxLen:             256,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9_-]+(\.fifo)?$`),
			PatternDescription: "letters, numbers, underscores, and hyphens; FIFO topics must end with .fifo",
		},
		"sqs": {
			MinLen:             1,
			MaxLen:             80,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9_-]+(\.fifo)?$`),
			PatternDescription: "letters, numbers, underscores, and hyphens; FIFO queues must end with .fifo",
		},
		"sqs_queue": {
			MinLen:             1,
			MaxLen:             80,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9_-]+(\.fifo)?$`),
			PatternDescription: "letters, numbers, underscores, and hyphens; FIFO queues must end with .fifo",
		},
		"lambda": {
			MinLen:             1,
			MaxLen:             64,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9-_]+$`),
			PatternDescription: "letters, numbers, hyphens, and underscores",
		},
		"kms_alias": {
			MinLen:             1,
			MaxLen:             256,
			Pattern:            regexp.MustCompile(`^alias/[a-zA-Z0-9/_-]+$`),
			PatternDescription: "must begin with alias/ and contain only letters, numbers, slashes, underscores, and hyphens",
			ForbiddenPrefixes:  []string{"alias/aws/"},
		},
		"log_group": {
			MinLen:             1,
			MaxLen:             512,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9_\-/.#]+$`),
			PatternDescription: "letters, numbers, underscore, hyphen, slash, period, and #",
			ForbiddenPrefixes:  []string{"aws/"},
		},
		"cloudwatch_log_group": {
			MinLen:             1,
			MaxLen:             512,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9_\-/.#]+$`),
			PatternDescription: "letters, numbers, underscore, hyphen, slash, period, and #",
			ForbiddenPrefixes:  []string{"aws/"},
		},
		"sec_group": {
			MinLen:             1,
			MaxLen:             255,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9 ._\-:/()#,@\[\]+=&;{}!$*]+$`),
			PatternDescription: "letters, numbers, spaces, and ._-:/()#,@[]+=&;{}!$*",
			ForbiddenPrefixes:  []string{"sg-"},
			CaseInsensitive:    true,
		},
		"security_group": {
			MinLen:             1,
			MaxLen:             255,
			Pattern:            regexp.MustCompile(`^[a-zA-Z0-9 ._\-:/()#,@\[\]+=&;{}!$*]+$`),
			PatternDescription: "letters, numbers, spaces, and ._-:/()#,@[]+=&;{}!$*",
			ForbiddenPrefixes:  []string{"sg-"},
			CaseInsensitive:    true,
		},
	}

	for key := range DefaultResourceAcronyms() {
		if _, exists := constraints[key]; exists {
			continue
		}
		if constraint, ok := defaultAWSDataResourceConstraint(key); ok {
			constraints[key] = constraint
		}
	}

	return constraints
}

func defaultAWSDataResourceConstraint(key string) (ResourceConstraint, bool) {
	maxLen := 0
	switch {
	case strings.HasPrefix(key, "athena_"):
		maxLen = 128
	case strings.HasPrefix(key, "bedrockagentcore_"):
		maxLen = 48
	case strings.HasPrefix(key, "bedrockagent_"):
		maxLen = 100
	case strings.HasPrefix(key, "bedrock_"):
		maxLen = 63
	case strings.HasPrefix(key, "datazone_"):
		maxLen = 64
	case strings.HasPrefix(key, "emrcontainers_") || strings.HasPrefix(key, "emrserverless_"):
		maxLen = 64
	case strings.HasPrefix(key, "emr_"):
		maxLen = 256
	case strings.HasPrefix(key, "glue_"):
		maxLen = 255
	case strings.HasPrefix(key, "kinesis_firehose_"):
		maxLen = 64
	case strings.HasPrefix(key, "kinesis_") || strings.HasPrefix(key, "kinesisanalyticsv2_"):
		maxLen = 128
	case strings.HasPrefix(key, "lakeformation_"):
		maxLen = 128
	case strings.HasPrefix(key, "quicksight_"):
		maxLen = 64
	case strings.HasPrefix(key, "redshiftserverless_") || strings.HasPrefix(key, "redshiftdata_") || strings.HasPrefix(key, "redshift_"):
		maxLen = 63
	case strings.HasPrefix(key, "sagemaker_"):
		maxLen = 63
	default:
		return ResourceConstraint{}, false
	}

	constraint := ResourceConstraint{
		MinLen:             1,
		MaxLen:             maxLen,
		Pattern:            regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9_.:/+=,@-]*[a-zA-Z0-9])?$`),
		PatternDescription: "letters, numbers, and AWS-safe separators; must start and end with a letter or number",
	}

	switch key {
	case "athena_database":
		constraint.MaxLen = 255
		constraint.Pattern = regexp.MustCompile(`^[0-9a-z_]+$`)
		constraint.PatternDescription = "lowercase letters, numbers, and underscores"
	case "athena_prepared_statement":
		constraint.MaxLen = 256
		constraint.Pattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_@:]{1,255}$`)
		constraint.PatternDescription = "2-256 characters; must start with a letter or underscore and contain only letters, numbers, underscores, @, and :"
	case "bedrock_guardrail":
		constraint.MaxLen = 50
	case "bedrockagentcore_harness":
		constraint.MaxLen = 40
		constraint.Pattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
		constraint.PatternDescription = "must start with a letter and contain only letters, numbers, and underscores"
	case "bedrockagentcore_agent_runtime",
		"bedrockagentcore_agent_runtime_endpoint",
		"bedrockagentcore_browser",
		"bedrockagentcore_browser_profile",
		"bedrockagentcore_code_interpreter",
		"bedrockagentcore_evaluator",
		"bedrockagentcore_memory",
		"bedrockagentcore_memory_strategy",
		"bedrockagentcore_online_evaluation_config",
		"bedrockagentcore_policy",
		"bedrockagentcore_policy_engine":
		constraint.Pattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
		constraint.PatternDescription = "must start with a letter and contain only letters, numbers, and underscores"
	case "datazone_form_type":
		constraint.MaxLen = 36
		constraint.Pattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
		constraint.PatternDescription = "a Smithy identifier containing only letters, numbers, and underscores"
	case "redshift_cluster", "redshift_parameter_group", "redshift_subnet_group":
		constraint.Pattern = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)
		constraint.PatternDescription = "lowercase letters, numbers, and hyphens; must start with a letter and not end with a hyphen"
		constraint.ForbiddenSubstrings = []string{"--"}
	case "redshift_endpoint_access", "redshiftserverless_endpoint_access":
		constraint.MaxLen = 30
		constraint.Pattern = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)
		constraint.PatternDescription = "lowercase letters, numbers, and hyphens; must start with a letter and not end with a hyphen"
		constraint.ForbiddenSubstrings = []string{"--"}
	case "sagemaker_flow_definition", "sagemaker_human_task_ui":
		constraint.Pattern = regexp.MustCompile(`^[0-9a-z](?:[0-9a-z-]*[0-9a-z])?$`)
		constraint.PatternDescription = "lowercase letters, numbers, and hyphens; must start and end with a letter or number"
	case "sagemaker_hyper_parameter_tuning_job", "sagemaker_project":
		constraint.MaxLen = 32
	}

	return constraint, true
}
