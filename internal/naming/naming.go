package naming

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	StyleDashed    = "dashed"
	StyleUnderscore = "underscore"
	StyleStraight  = "straight"
	StylePascal    = "pascal"
	StyleCamel     = "camel"
)

var (
	wordRe = regexp.MustCompile(`[A-Za-z0-9]+`)
)

type Config struct {
	OrgPrefix              string
	Project                string
	Env                    string
	Region                 string
	RegionShortCode         string
	RegionMap              map[string]string
	Recipe                 []string
	StylePriority          []string
	ResourceAcronyms        map[string]string
	ResourceStyleOverrides map[string][]string
	ResourceConstraints    map[string]ResourceConstraint
}

type BuildInput struct {
	Resource      string
	Qualifier     string
	Overrides     map[string]string
	Recipe        []string
	StylePriority []string
}

type BuildResult struct {
	Name           string
	Style          string
	Components     map[string]string
	Parts          []string
	RegionCode     string
	ResourceAcronym string
}

type ResourceConstraint struct {
	MinLen             int
	MaxLen             int
	Pattern            *regexp.Regexp
	PatternDescription string
	ForbiddenPrefixes  []string
	ForbiddenSuffixes  []string
	ForbiddenSubstrings []string
	DisallowIPAddress  bool
	CaseInsensitive    bool
}

func DefaultRecipe() []string {
	return []string{"org", "proj", "env", "region", "resource", "qualifier"}
}

func DefaultStylePriority() []string {
	return []string{StyleDashed, StylePascal, StyleCamel, StyleStraight, StyleUnderscore}
}

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
		"role":                     "role",
		"role_policy":              "rlpl",
		"iam_role":                 "role",
		"iam_policy":               "iamp",
		"iam_user":                 "iamu",
		"iam_group":                "iamg",
		"s3":                       "s3b",
		"s3_bucket":                "s3bk",
		"s3_object":                "s3ob",
		"s3_access_point":          "s3ap",
		"s3_table":                 "s3tb",
		"s3_dir":                   "s3dr",
		"sns":                      "sns",
		"sqs":                      "sqs",
		"ecs_cluster":              "ecsc",
		"ecs_service":              "ecss",
		"ecs_task":                 "ecst",
		"eks":                      "eks",
		"eks_cluster":              "eksc",
		"eks_node_group":           "ekng",
		"msk_cluster":              "mskc",
		"vpc":                      "vpcn",
		"subnet":                   "subn",
		"igw":                      "igtw",
		"nat_gw":                   "ngtw",
		"sec_group":                "scgp",
		"nacl":                     "nacl",
		"route_table":              "rttb",
		"elastic_ip":               "elip",
		"wafv2_web_acl":            "wfac",
		"wafv2_web_acl_rule":       "wfar",
		"wafv2_ip_set":             "wfis",
		"lambda":                   "lmbd",
		"api_gateway_rest_api":     "agra",
		"api_gateway_model":        "agmd",
		"api_gateway_v2":           "agv2",
		"log_group":                "logg",
		"cloudwatch_log_group":     "cwlg",
		"cloudwatch_alarm":         "cwal",
		"eventbridge_bus":          "evbb",
		"eventbridge_rule":         "evbr",
		"step_function":            "stfn",
		"sfn":                      "stfn",
		"dynamodb":                 "dydb",
		"dynamodb_table":           "dydb",
		"rds":                      "rds",
		"rds_cluster":              "rdsc",
		"aurora_cluster":           "arcl",
		"redshift":                 "rdsh",
		"elasticache":              "elch",
		"opensearch":               "opsr",
		"elasticsearch":            "elsr",
		"ecr":                      "ecr",
		"ecs":                      "ecs",
		"ec2_instance":             "ec2i",
		"launch_template":          "lcht",
		"autoscaling_group":        "asgr",
		"alb":                      "albl",
		"nlb":                      "nlbl",
		"elb":                      "elbl",
		"target_group":             "tgpt",
		"cloudfront":               "clfr",
		"route53_zone":             "rt53",
		"route53_record":           "r53r",
		"acm_cert":                 "acmc",
		"kms_key":                  "kmsk",
		"secretsmanager_secret":    "smse",
		"ssm_parameter":            "ssmp",
		"cloudtrail":               "ctra",
		"guardduty":                "gdty",
		"config_rule":              "cfrl",
		"efs":                      "efs",
		"ebs":                      "ebs",
		"athena":                   "athn",
		"glue":                     "glue",
		"sagemaker":                "sgmk",
		"codebuild":                "cdbd",
		"codepipeline":             "cdpl",
		"codedeploy":               "cddp",
		"cloudformation_stack":     "cfst",
		"appsync":                  "apsy",
		"snow_notification_integration": "snti",
	}
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
			MinLen:             3,
			MaxLen:             63,
			Pattern:            regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`),
			PatternDescription: "lowercase letters, numbers, dots, and hyphens; must start and end with a letter or number",
			ForbiddenPrefixes:  []string{"xn--", "sthree-", "amzn-s3-demo-"},
			ForbiddenSuffixes:  []string{"-s3alias", "--ol-s3"},
			ForbiddenSubstrings: []string{".."},
			DisallowIPAddress:  true,
		},
		"s3_bucket": {
			MinLen:             3,
			MaxLen:             63,
			Pattern:            regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`),
			PatternDescription: "lowercase letters, numbers, dots, and hyphens; must start and end with a letter or number",
			ForbiddenPrefixes:  []string{"xn--", "sthree-", "amzn-s3-demo-"},
			ForbiddenSuffixes:  []string{"-s3alias", "--ol-s3"},
			ForbiddenSubstrings: []string{".."},
			DisallowIPAddress:  true,
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

func BuildName(cfg Config, in BuildInput) (BuildResult, error) {
	regionCode := strings.TrimSpace(cfg.RegionShortCode)
	if regionCode == "" && strings.TrimSpace(cfg.Region) != "" {
		regionCode = cfg.RegionMap[strings.TrimSpace(cfg.Region)]
	}

	resourceKey := strings.ToLower(strings.TrimSpace(in.Resource))
	resourceAcronym := strings.TrimSpace(in.Resource)
	if resourceKey != "" {
		if v, ok := cfg.ResourceAcronyms[resourceKey]; ok && v != "" {
			resourceAcronym = v
		}
	}

	components := map[string]string{
		"org":       strings.TrimSpace(cfg.OrgPrefix),
		"proj":      strings.TrimSpace(cfg.Project),
		"env":       strings.TrimSpace(cfg.Env),
		"region":    strings.TrimSpace(regionCode),
		"resource":  strings.TrimSpace(resourceAcronym),
		"qualifier": strings.TrimSpace(in.Qualifier),
	}

	overrides := in.Overrides
	if overrides == nil {
		overrides = map[string]string{}
	}

	for key, val := range overrides {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		canonical := canonicalComponentKey(key)
		if _, ok := components[canonical]; ok {
			components[canonical] = strings.TrimSpace(val)
		} else {
			components[key] = strings.TrimSpace(val)
		}
	}
	regionCode = components["region"]

	recipe := cfg.Recipe
	if len(in.Recipe) > 0 {
		recipe = in.Recipe
	}
	if len(recipe) == 0 {
		recipe = DefaultRecipe()
	}

	parts := make([]string, 0, len(recipe))
	for _, item := range recipe {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		canonical := canonicalComponentKey(item)
		val := ""
		if v, ok := components[canonical]; ok {
			val = v
		} else if v, ok := components[item]; ok {
			val = v
		}
		if strings.TrimSpace(val) == "" {
			continue
		}
		parts = append(parts, val)
	}

	stylePriority := cfg.StylePriority
	if len(in.StylePriority) > 0 {
		stylePriority = in.StylePriority
	}
	if len(stylePriority) == 0 {
		stylePriority = DefaultStylePriority()
	}

	allowedStyles := []string{}
	if len(cfg.ResourceStyleOverrides) > 0 && resourceKey != "" {
		if v, ok := cfg.ResourceStyleOverrides[resourceKey]; ok {
			allowedStyles = normalizeStyles(v)
		}
	}

	chosenStyle := ""
	for _, style := range stylePriority {
		style = normalizeStyle(style)
		if !isValidStyle(style) {
			continue
		}
		if len(allowedStyles) > 0 && !containsString(allowedStyles, style) {
			continue
		}
		chosenStyle = style
		break
	}
	if chosenStyle == "" {
		chosenStyle = StyleDashed
	}

	name, err := formatName(chosenStyle, parts)
	if err != nil {
		return BuildResult{}, err
	}
	if err := validateResourceConstraints(resourceKey, name, cfg.ResourceConstraints); err != nil {
		return BuildResult{}, err
	}

	return BuildResult{
		Name:            name,
		Style:           chosenStyle,
		Components:      components,
		Parts:           parts,
		RegionCode:      regionCode,
		ResourceAcronym: components["resource"],
	}, nil
}

func canonicalComponentKey(key string) string {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "org_prefix", "org":
		return "org"
	case "project", "proj":
		return "proj"
	case "environment", "env":
		return "env"
	case "region", "region_code", "region_short_code":
		return "region"
	case "resource", "resource_type":
		return "resource"
	case "qualifier", "qual":
		return "qualifier"
	default:
		return strings.ToLower(strings.TrimSpace(key))
	}
}

func normalizeStyle(style string) string {
	return strings.ToLower(strings.TrimSpace(style))
}

func validateResourceConstraints(resourceKey, name string, constraints map[string]ResourceConstraint) error {
	if resourceKey == "" || len(name) == 0 {
		return nil
	}
	if len(constraints) == 0 {
		constraints = DefaultResourceConstraints()
	}
	c, ok := constraints[resourceKey]
	if !ok {
		return nil
	}
	if c.MinLen > 0 && len(name) < c.MinLen {
		return fmt.Errorf("resource %q name %q is shorter than %d characters", resourceKey, name, c.MinLen)
	}
	if c.MaxLen > 0 && len(name) > c.MaxLen {
		return fmt.Errorf("resource %q name %q exceeds %d characters", resourceKey, name, c.MaxLen)
	}
	if c.Pattern != nil && !c.Pattern.MatchString(name) {
		desc := c.PatternDescription
		if desc == "" {
			desc = c.Pattern.String()
		}
		return fmt.Errorf("resource %q name %q must match: %s", resourceKey, name, desc)
	}
	comparisonName := name
	if c.CaseInsensitive {
		comparisonName = strings.ToLower(name)
	}
	if len(c.ForbiddenPrefixes) > 0 {
		for _, prefix := range c.ForbiddenPrefixes {
			if prefix == "" {
				continue
			}
			candidate := prefix
			if c.CaseInsensitive {
				candidate = strings.ToLower(prefix)
			}
			if strings.HasPrefix(comparisonName, candidate) {
				return fmt.Errorf("resource %q name %q must not start with prefix %q", resourceKey, name, prefix)
			}
		}
	}
	if len(c.ForbiddenSuffixes) > 0 {
		for _, suffix := range c.ForbiddenSuffixes {
			if suffix == "" {
				continue
			}
			candidate := suffix
			if c.CaseInsensitive {
				candidate = strings.ToLower(suffix)
			}
			if strings.HasSuffix(comparisonName, candidate) {
				return fmt.Errorf("resource %q name %q must not end with suffix %q", resourceKey, name, suffix)
			}
		}
	}
	if len(c.ForbiddenSubstrings) > 0 {
		for _, sub := range c.ForbiddenSubstrings {
			if sub == "" {
				continue
			}
			candidate := sub
			if c.CaseInsensitive {
				candidate = strings.ToLower(sub)
			}
			if strings.Contains(comparisonName, candidate) {
				return fmt.Errorf("resource %q name %q must not contain %q", resourceKey, name, sub)
			}
		}
	}
	if c.DisallowIPAddress && isIPv4Address(name) {
		return fmt.Errorf("resource %q name %q must not be formatted as an IP address", resourceKey, name)
	}
	return nil
}

func isIPv4Address(value string) bool {
	ip := net.ParseIP(value)
	return ip != nil && ip.To4() != nil
}

func normalizeStyles(styles []string) []string {
	out := make([]string, 0, len(styles))
	for _, s := range styles {
		n := normalizeStyle(s)
		if n == "" {
			continue
		}
		out = append(out, n)
	}
	return out
}

func isValidStyle(style string) bool {
	switch style {
	case StyleDashed, StyleUnderscore, StyleStraight, StylePascal, StyleCamel:
		return true
	default:
		return false
	}
}

func formatName(style string, parts []string) (string, error) {
	switch style {
	case StyleDashed:
		return strings.Join(normalizeParts(parts, "-", false), "-"), nil
	case StyleUnderscore:
		return strings.Join(normalizeParts(parts, "_", false), "_"), nil
	case StyleStraight:
		return strings.Join(normalizeParts(parts, "", false), ""), nil
	case StylePascal:
		return strings.Join(normalizeParts(parts, "", true), ""), nil
	case StyleCamel:
		return camelize(parts), nil
	default:
		return "", errors.New("unsupported style")
	}
}

func normalizeParts(parts []string, sep string, pascal bool) []string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if pascal {
			out = append(out, pascalize(p))
			continue
		}
		words := splitWords(p)
		if len(words) == 0 {
			continue
		}
		if sep == "" {
			out = append(out, strings.ToLower(strings.Join(words, "")))
		} else {
			out = append(out, strings.ToLower(strings.Join(words, sep)))
		}
	}
	return out
}

func camelize(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	firstWords := splitWords(parts[0])
	first := ""
	if len(firstWords) > 0 {
		first = strings.ToLower(strings.Join(firstWords, ""))
	}

	rest := make([]string, 0, len(parts))
	for _, p := range parts[1:] {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		rest = append(rest, pascalize(p))
	}

	return first + strings.Join(rest, "")
}

func pascalize(value string) string {
	words := splitWords(value)
	if len(words) == 0 {
		return ""
	}
	out := make([]string, 0, len(words))
	for _, w := range words {
		out = append(out, titleWord(w))
	}
	return strings.Join(out, "")
}

func splitWords(value string) []string {
	return wordRe.FindAllString(value, -1)
}

func titleWord(word string) string {
	if word == "" {
		return ""
	}
	if len(word) == 1 {
		return strings.ToUpper(word)
	}
	return strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
}

func containsString(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
