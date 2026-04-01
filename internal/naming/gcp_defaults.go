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
		"northamerica-south1":     "nas1",
		"southamerica-east1":      "sae1",
		"southamerica-west1":      "saw1",
		"europe-central2":         "euc2",
		"europe-north1":           "eun1",
		"europe-north2":           "eun2",
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
		"asia-southeast3":         "asse3",
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
		"storage_bucket":                        "gcs",
		"gcs_bucket":                            "gcs",
		"gcs":                                   "gcs",
		"compute_network":                       "vpc",
		"compute_subnetwork":                    "snet",
		"vpc":                                   "vpc",
		"subnet":                                "snet",
		"pubsub_topic":                          "pst",
		"pubsub_subscription":                   "pss",
		"service_account":                       "gsa",
		"bigquery_dataset":                      "bqd",
		"artifact_registry_repository":          "arr",
		"cloud_run_v2_service":                  "crs",
		"cloud_run_service":                     "crs",
		"compute_router":                        "rtr",
		"compute_firewall":                      "fwl",
		"compute_address":                       "addr",
		"compute_global_address":                "gadr",
		"compute_route":                         "rte",
		"compute_router_nat":                    "rnat",
		"compute_vpn_gateway":                   "vpng",
		"compute_vpn_tunnel":                    "vpnt",
		"compute_ha_vpn_gateway":                "hvgw",
		"compute_url_map":                       "umap",
		"compute_target_http_proxy":             "thp",
		"compute_target_https_proxy":            "thps",
		"compute_backend_service":               "bksv",
		"compute_region_backend_service":        "rbksv",
		"compute_instance_template":             "itpl",
		"compute_instance_group_manager":        "igm",
		"compute_region_instance_group_manager": "rigm",
		"compute_disk":                          "pdsk",
		"compute_image":                         "img",
		"compute_snapshot":                      "snap",
		"dns_managed_zone":                      "dnsz",
		"secret_manager_secret":                 "sms",
		"kms_key_ring":                          "kmr",
		"kms_crypto_key":                        "kmk",
		"sql_database_instance":                 "sqli",
		"sql_instance":                          "sqli",
		"container_cluster":                     "gke",
		"gke_cluster":                           "gke",
		"container_node_pool":                   "npol",
		"gke_node_pool":                         "npol",
		"vpc_access_connector":                  "vpac",
		"redis_instance":                        "rdis",
		"memcache_instance":                     "memc",
		"filestore_instance":                    "fils",
		"spanner_instance":                      "spni",
		"spanner_database":                      "spdb",
		"cloudbuild_trigger":                    "cbt",
		"eventarc_trigger":                      "evtr",
		"cloud_scheduler_job":                   "schd",
		"cloud_tasks_queue":                     "ctsk",
		"workflows_workflow":                    "wflw",
		"monitoring_notification_channel":       "mnc",
		"logging_metric":                        "lgmt",
		"logging_project_sink":                  "lpsk",
		"pubsub_schema":                         "psch",
		"pubsub_snapshot":                       "psnp",
	}
}

func DefaultGCPGlobalResources() map[string]bool {
	return map[string]bool{
		"storage_bucket":                  true,
		"gcs_bucket":                      true,
		"gcs":                             true,
		"compute_network":                 true,
		"vpc":                             true,
		"pubsub_topic":                    true,
		"pubsub_subscription":             true,
		"service_account":                 true,
		"compute_firewall":                true,
		"compute_global_address":          true,
		"compute_route":                   true,
		"compute_url_map":                 true,
		"compute_target_http_proxy":       true,
		"compute_target_https_proxy":      true,
		"compute_backend_service":         true,
		"compute_image":                   true,
		"dns_managed_zone":                true,
		"secret_manager_secret":           true,
		"cloudbuild_trigger":              true,
		"monitoring_notification_channel": true,
		"logging_metric":                  true,
		"logging_project_sink":            true,
		"pubsub_schema":                   true,
		"pubsub_snapshot":                 true,
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
		"storage_bucket":       {StyleDashed, StyleUnderscore, StyleStraight},
		"gcs_bucket":           {StyleDashed, StyleUnderscore, StyleStraight},
		"gcs":                  {StyleDashed, StyleUnderscore, StyleStraight},
		"compute_network":      {StyleDashed, StyleStraight},
		"compute_subnetwork":   {StyleDashed, StyleStraight},
		"vpc":                  {StyleDashed, StyleStraight},
		"subnet":               {StyleDashed, StyleStraight},
		"pubsub_topic":         {StyleDashed, StyleUnderscore, StyleStraight},
		"pubsub_subscription":  {StyleDashed, StyleUnderscore, StyleStraight},
		"service_account":      {StyleDashed, StyleStraight},
		"bigquery_dataset":     {StyleStraight, StyleUnderscore},
		"cloud_run_v2_service": {StyleDashed, StyleStraight},
		"cloud_run_service":    {StyleDashed, StyleStraight},
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

	pubsubConstraint := ResourceConstraint{
		MinLen:             3,
		MaxLen:             255,
		Pattern:            regexp.MustCompile(`^[A-Za-z][A-Za-z0-9._~+%-]*$`),
		PatternDescription: "must start with a letter and contain only letters, numbers, hyphens, underscores, periods, tildes, plus signs, or percent signs",
		ForbiddenPrefixes:  []string{"goog"},
	}

	serviceAccountConstraint := ResourceConstraint{
		MinLen:             6,
		MaxLen:             30,
		Pattern:            regexp.MustCompile(`^[a-z]([-a-z0-9]*[a-z0-9])$`),
		PatternDescription: "must be 6-30 characters, start with a lowercase letter, and contain only lowercase letters, numbers, or hyphens",
	}

	bigQueryDatasetConstraint := ResourceConstraint{
		MinLen:             1,
		MaxLen:             1024,
		Pattern:            regexp.MustCompile(`^[A-Za-z0-9_]+$`),
		PatternDescription: "must contain only letters, numbers, or underscores",
	}

	// Cloud Run documents start/end/length rules for service IDs. The RFC1035-like
	// character class here is the compatible subset we generate for dashed/straight styles.
	cloudRunServiceConstraint := ResourceConstraint{
		MinLen:             1,
		MaxLen:             49,
		Pattern:            regexp.MustCompile(`^[a-z]([-a-z0-9]*[a-z0-9])?$`),
		PatternDescription: "must start with a lowercase letter, end with a lowercase letter or number, and contain only lowercase letters, numbers, or hyphens",
	}

	return map[string]ResourceConstraint{
		"storage_bucket":       bucketConstraint,
		"gcs_bucket":           bucketConstraint,
		"gcs":                  bucketConstraint,
		"compute_network":      rfc1035Constraint,
		"compute_subnetwork":   rfc1035Constraint,
		"vpc":                  rfc1035Constraint,
		"subnet":               rfc1035Constraint,
		"pubsub_topic":         pubsubConstraint,
		"pubsub_subscription":  pubsubConstraint,
		"service_account":      serviceAccountConstraint,
		"bigquery_dataset":     bigQueryDatasetConstraint,
		"cloud_run_v2_service": cloudRunServiceConstraint,
		"cloud_run_service":    cloudRunServiceConstraint,
	}
}
