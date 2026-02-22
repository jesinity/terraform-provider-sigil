package naming

import "regexp"

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
		"role":                          "role",
		"role_policy":                   "rlpl",
		"iam_role":                      "role",
		"iam_policy":                    "iamp",
		"iam_user":                      "iamu",
		"iam_group":                     "iamg",
		"s3":                            "s3b",
		"s3_bucket":                     "s3bk",
		"s3_object":                     "s3ob",
		"s3_access_point":               "s3ap",
		"s3_table":                      "s3tb",
		"s3_dir":                        "s3dr",
		"sns":                           "sns",
		"sqs":                           "sqs",
		"ecs_cluster":                   "ecsc",
		"ecs_service":                   "ecss",
		"ecs_task":                      "ecst",
		"eks":                           "eks",
		"eks_cluster":                   "eksc",
		"eks_node_group":                "ekng",
		"msk_cluster":                   "mskc",
		"vpc":                           "vpcn",
		"subnet":                        "subn",
		"igw":                           "igtw",
		"nat_gw":                        "ngtw",
		"sec_group":                     "scgp",
		"nacl":                          "nacl",
		"route_table":                   "rttb",
		"elastic_ip":                    "elip",
		"wafv2_web_acl":                 "wfac",
		"wafv2_web_acl_rule":            "wfar",
		"wafv2_ip_set":                  "wfis",
		"lambda":                        "lmbd",
		"api_gateway_rest_api":          "agra",
		"api_gateway_model":             "agmd",
		"api_gateway_v2":                "agv2",
		"log_group":                     "logg",
		"cloudwatch_log_group":          "cwlg",
		"cloudwatch_alarm":              "cwal",
		"eventbridge_bus":               "evbb",
		"eventbridge_rule":              "evbr",
		"step_function":                 "stfn",
		"sfn":                           "stfn",
		"dynamodb":                      "dydb",
		"dynamodb_table":                "dybt",
		"rds":                           "rds",
		"rds_cluster":                   "rdsc",
		"aurora_cluster":                "arcl",
		"redshift":                      "rdsh",
		"elasticache":                   "elch",
		"opensearch":                    "opsr",
		"elasticsearch":                 "elsr",
		"ecr":                           "ecr",
		"ecs":                           "ecs",
		"ec2_instance":                  "ec2i",
		"launch_template":               "lcht",
		"autoscaling_group":             "asgr",
		"alb":                           "albl",
		"nlb":                           "nlbl",
		"elb":                           "elbl",
		"target_group":                  "tgpt",
		"cloudfront":                    "clfr",
		"route53_zone":                  "rt53",
		"route53_record":                "r53r",
		"acm_cert":                      "acmc",
		"kms_key":                       "kmsk",
		"secretsmanager_secret":         "smse",
		"ssm_parameter":                 "ssmp",
		"cloudtrail":                    "ctra",
		"guardduty":                     "gdty",
		"config_rule":                   "cfrl",
		"efs":                           "efs",
		"ebs":                           "ebs",
		"athena":                        "athn",
		"glue":                          "glue",
		"sagemaker":                     "sgmk",
		"codebuild":                     "cdbd",
		"codepipeline":                  "cdpl",
		"codedeploy":                    "cddp",
		"cloudformation_stack":          "cfst",
		"appsync":                       "apsy",
		"snow_notification_integration": "snti",
	}
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
		"s3":        {StyleDashed, StyleStraight},
		"s3_bucket": {StyleDashed, StyleStraight},
	}
}

func DefaultResourceConstraints() map[string]ResourceConstraint {
	return map[string]ResourceConstraint{
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
}
