package naming

import "regexp"

func DefaultGCPRegionMap() map[string]string {
	return map[string]string{
		"us-central1":             "usc1",
		"us-east1":                "use1",
		"us-east4":                "use4",
		"us-east5":                "use5",
		"us-south1":               "uss1",
		"us-west1":                "usw1",
		"us-west2":                "usw2",
		"us-west3":                "usw3",
		"us-west4":                "usw4",
		"northamerica-northeast1": "nan1",
		"northamerica-northeast2": "nan2",
		"southamerica-east1":      "sae1",
		"southamerica-west1":      "saw1",
		"europe-central2":         "euc2",
		"europe-north1":           "eun1",
		"europe-southwest1":       "eusw1",
		"europe-west1":            "euw1",
		"europe-west2":            "euw2",
		"europe-west3":            "euw3",
		"europe-west4":            "euw4",
		"europe-west6":            "euw6",
		"europe-west8":            "euw8",
		"europe-west9":            "euw9",
		"europe-west10":           "euw10",
		"europe-west12":           "euw12",
		"asia-east1":              "ase1",
		"asia-east2":              "ase2",
		"asia-northeast1":         "asne1",
		"asia-northeast2":         "asne2",
		"asia-northeast3":         "asne3",
		"asia-south1":             "ass1",
		"asia-south2":             "ass2",
		"asia-southeast1":         "asse1",
		"asia-southeast2":         "asse2",
		"australia-southeast1":    "ause1",
		"australia-southeast2":    "ause2",
		"me-central1":             "mec1",
		"me-central2":             "mec2",
		"me-west1":                "mew1",
		"africa-south1":           "afs1",
	}
}

func DefaultGCPResourceAcronyms() map[string]string {
	return map[string]string{
		"google_storage_bucket":               "gcs",
		"gcs_bucket":                          "gcs",
		"gcs":                                 "gcs",
		"google_compute_network":              "vpc",
		"google_compute_subnetwork":           "snet",
		"vpc":                                 "vpc",
		"subnet":                              "snet",
		"google_pubsub_topic":                 "pst",
		"google_pubsub_subscription":          "pss",
		"pubsub_topic":                        "pst",
		"pubsub_subscription":                 "pss",
		"google_service_account":              "gsa",
		"service_account":                     "gsa",
		"google_bigquery_dataset":             "bqd",
		"bigquery_dataset":                    "bqd",
		"google_artifact_registry_repository": "arr",
		"artifact_registry_repository":        "arr",
		"google_cloud_run_v2_service":         "crs",
		"cloud_run_service":                   "crs",
	}
}

func DefaultGCPGlobalResources() map[string]bool {
	return map[string]bool{
		"google_storage_bucket":      true,
		"gcs_bucket":                 true,
		"gcs":                        true,
		"google_compute_network":     true,
		"vpc":                        true,
		"google_pubsub_topic":        true,
		"pubsub_topic":               true,
		"google_pubsub_subscription": true,
		"pubsub_subscription":        true,
		"google_service_account":     true,
		"service_account":            true,
	}
}

func DefaultGCPRegionalResources() map[string]bool {
	regional := map[string]bool{}
	for key := range DefaultGCPResourceAcronyms() {
		regional[key] = true
	}
	for key := range DefaultGCPGlobalResources() {
		delete(regional, key)
	}
	return regional
}

func DefaultGCPResourceStyleOverrides() map[string][]string {
	return map[string][]string{
		"google_storage_bucket":      {StyleDashed, StyleUnderscore, StyleStraight},
		"gcs_bucket":                 {StyleDashed, StyleUnderscore, StyleStraight},
		"gcs":                        {StyleDashed, StyleUnderscore, StyleStraight},
		"google_compute_network":     {StyleDashed, StyleStraight},
		"google_compute_subnetwork":  {StyleDashed, StyleStraight},
		"vpc":                        {StyleDashed, StyleStraight},
		"subnet":                     {StyleDashed, StyleStraight},
		"google_pubsub_topic":        {StyleDashed, StyleUnderscore, StyleStraight},
		"google_pubsub_subscription": {StyleDashed, StyleUnderscore, StyleStraight},
		"pubsub_topic":               {StyleDashed, StyleUnderscore, StyleStraight},
		"pubsub_subscription":        {StyleDashed, StyleUnderscore, StyleStraight},
	}
}

func DefaultGCPResourceConstraints() map[string]ResourceConstraint {
	bucketConstraint := ResourceConstraint{
		MinLen:              3,
		MaxLen:              63,
		Pattern:             regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*[a-z0-9]$`),
		PatternDescription:  "lowercase letters, numbers, dots, underscores, and hyphens; must start and end with a letter or number",
		ForbiddenPrefixes:   []string{"goog"},
		ForbiddenSubstrings: []string{"google"},
		DisallowIPAddress:   true,
	}

	rfc1035Constraint := ResourceConstraint{
		MinLen:             1,
		MaxLen:             63,
		Pattern:            regexp.MustCompile(`^[a-z]([-a-z0-9]*[a-z0-9])?$`),
		PatternDescription: "must start with a lowercase letter and contain only lowercase letters, numbers, or hyphens",
	}

	return map[string]ResourceConstraint{
		"google_storage_bucket":     bucketConstraint,
		"gcs_bucket":                bucketConstraint,
		"gcs":                       bucketConstraint,
		"google_compute_network":    rfc1035Constraint,
		"google_compute_subnetwork": rfc1035Constraint,
		"vpc":                       rfc1035Constraint,
		"subnet":                    rfc1035Constraint,
	}
}
