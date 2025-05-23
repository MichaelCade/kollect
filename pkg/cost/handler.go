package cost

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/michaelcade/kollect/pkg/aws"
	"github.com/michaelcade/kollect/pkg/azure"
	"github.com/michaelcade/kollect/pkg/gcp"
)

func HandleCostRequest(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	useMock := r.URL.Query().Get("mock") == "true"
	resourceType := r.URL.Query().Get("type")

	if resourceType == "" {
		resourceType = "snapshots"
	}

	log.Printf("Cost request received for platform: %s, resource type: %s (mock: %v)",
		platform, resourceType, useMock)

	ctx := r.Context()
	var costs interface{}
	var err error

	switch platform {

	case "aws":
		var awsData map[string]interface{}
		if useMock {
			awsData = GenerateMockResourceData("aws", resourceType)
			log.Printf("Using mock AWS data for %s", resourceType)
		} else {
			if resourceType == "snapshots" {
				awsData, err = aws.CollectSnapshotData(ctx)
			} else {
				awsData, err = ConvertAwsDataForCostAnalysis(ctx)
			}

			if err != nil {
				log.Printf("Error collecting AWS resource data: %v", err)
				http.Error(w, fmt.Sprintf("Error collecting AWS resource data: %v", err), http.StatusInternalServerError)
				return
			}

			log.Printf("AWS resource data collected with %d EBS snapshots, %d RDS snapshots, %d EC2 instances, %d RDS instances, %d S3 buckets",
				countResources(awsData, "EBSSnapshots"),
				countResources(awsData, "RDSSnapshots"),
				countResources(awsData, "EC2Instances"),
				countResources(awsData, "RDSInstances"),
				countResources(awsData, "S3Buckets"))

			if isEmpty(awsData) {
				log.Printf("No AWS resources found")
				costs = map[string]interface{}{
					"Summary": map[string]interface{}{
						"TotalSnapshotStorage": 0.0,
						"TotalMonthlyCost":     0.0,
						"Currency":             "USD",
					},
					"Message": "No AWS resources found. Real data is being shown.",
				}
				costs = map[string]interface{}{"aws": costs}
				break
			}
		}

		if resourceType == "snapshots" {
			costs, err = EstimateAwsResourceCosts(awsData)
		} else {
			estimator := NewResourceCostEstimator()
			costs, err = estimator.EstimateAwsResourcesCost(awsData)
		}
		costs = map[string]interface{}{"aws": costs}

	case "azure":
		var azureData map[string]interface{}

		if useMock {
			azureData = GenerateMockResourceData("azure", resourceType)
			log.Printf("Using mock Azure data for %s", resourceType)
		} else {
			if resourceType == "snapshots" {
				azureData, err = azure.CollectSnapshotData(ctx)
			} else {
				azureData, err = ConvertAzureDataForCostAnalysis(ctx)
			}

			if err != nil {
				log.Printf("Error collecting Azure resource data: %v", err)
				http.Error(w, fmt.Sprintf("Error collecting Azure resource data: %v", err), http.StatusInternalServerError)
				return
			}

			log.Printf("Azure resource data collected with %d disk snapshots, %d VMs, %d storage accounts, %d SQL databases",
				countResources(azureData, "DiskSnapshots"),
				countResources(azureData, "VirtualMachines"),
				countResources(azureData, "StorageAccounts"),
				countResources(azureData, "SQLDatabases"))

			if isEmpty(azureData) {
				log.Printf("No Azure resources found")
				costs = map[string]interface{}{
					"Summary": map[string]interface{}{
						"TotalSnapshotStorage": 0.0,
						"TotalMonthlyCost":     0.0,
						"Currency":             "USD",
					},
					"Message": "No Azure resources found. Real data is being shown.",
				}
				costs = map[string]interface{}{"azure": costs}
				break
			}
		}

		if resourceType == "snapshots" {
			costs, err = EstimateAzureResourceCosts(azureData)
		} else {
			estimator := NewResourceCostEstimator()
			costs, err = estimator.EstimateAzureResourcesCost(azureData)
		}
		costs = map[string]interface{}{"azure": costs}

	case "gcp":
		var gcpData map[string]interface{}

		if useMock {
			gcpData = GenerateMockResourceData("gcp", resourceType)
			log.Printf("Using mock GCP data for %s", resourceType)
		} else {
			if resourceType == "snapshots" {
				gcpData, err = gcp.CollectSnapshotData(ctx)
			} else {
				gcpData, err = ConvertGcpDataForCostAnalysis(ctx)
			}

			if err != nil {
				log.Printf("Error collecting GCP resource data: %v", err)
				http.Error(w, fmt.Sprintf("Error collecting GCP resource data: %v", err), http.StatusInternalServerError)
				return
			}

			log.Printf("GCP resource data collected with %d disk snapshots, %d compute instances, %d Cloud SQL instances, %d GCS buckets",
				countResources(gcpData, "DiskSnapshots"),
				countResources(gcpData, "ComputeInstances"),
				countResources(gcpData, "CloudSQLInstances"),
				countResources(gcpData, "GCSBuckets"))

			if isEmpty(gcpData) {
				log.Printf("No GCP resources found")
				costs = map[string]interface{}{
					"Summary": map[string]interface{}{
						"TotalSnapshotStorage": 0.0,
						"TotalMonthlyCost":     0.0,
						"Currency":             "USD",
					},
					"Message": "No GCP resources found. Real data is being shown.",
				}
				costs = map[string]interface{}{"gcp": costs}
				break
			}
		}

		if resourceType == "snapshots" {
			costs, err = EstimateGcpResourceCosts(gcpData)
		} else {
			estimator := NewResourceCostEstimator()
			costs, err = estimator.EstimateGcpResourcesCost(gcpData)
			log.Printf("Using ResourceCostEstimator for comprehensive GCP resource cost analysis")
		}

		if err != nil {
			log.Printf("Error estimating GCP costs: %v", err)
			http.Error(w, fmt.Sprintf("Error estimating GCP costs: %v", err), http.StatusInternalServerError)
			return
		}

		costs = map[string]interface{}{"gcp": costs}

	case "all":
		allCosts := make(map[string]interface{})

		var awsData map[string]interface{}
		if useMock {
			awsData = GenerateMockResourceData("aws", resourceType)
			log.Printf("Using mock AWS data for 'all' option (%s)", resourceType)

			if resourceType == "snapshots" {
				awsCosts, _ := EstimateAwsResourceCosts(awsData)
				allCosts["aws"] = awsCosts
			} else {
				estimator := NewResourceCostEstimator()
				awsCosts, _ := estimator.EstimateAwsResourcesCost(awsData)
				allCosts["aws"] = awsCosts
			}
		} else {
			if resourceType == "snapshots" {
				awsData, err = aws.CollectSnapshotData(ctx)
			} else {
				awsData, err = ConvertAwsDataForCostAnalysis(ctx)
			}

			if err != nil {
				log.Printf("Warning: Error collecting AWS resources: %v", err)
			} else {
				log.Printf("AWS resource data collected for 'all' option")

				if !isEmpty(awsData) {
					if resourceType == "snapshots" {
						awsCosts, _ := EstimateAwsResourceCosts(awsData)
						allCosts["aws"] = awsCosts
					} else {
						estimator := NewResourceCostEstimator()
						awsCosts, _ := estimator.EstimateAwsResourcesCost(awsData)
						allCosts["aws"] = awsCosts
					}
				} else {
					log.Printf("No AWS resources found")
					allCosts["aws"] = map[string]interface{}{
						"Summary": map[string]interface{}{
							"TotalSnapshotStorage": 0.0,
							"TotalMonthlyCost":     0.0,
							"Currency":             "USD",
						},
						"Message": "No AWS resources found. Real data is being shown.",
					}
				}
			}
		}

		var azureData map[string]interface{}
		if useMock {
			azureData = GenerateMockResourceData("azure", resourceType)
			log.Printf("Using mock Azure data for 'all' option (%s)", resourceType)

			if resourceType == "snapshots" {
				azureCosts, _ := EstimateAzureResourceCosts(azureData)
				allCosts["azure"] = azureCosts
			} else {
				estimator := NewResourceCostEstimator()
				azureCosts, _ := estimator.EstimateAzureResourcesCost(azureData)
				allCosts["azure"] = azureCosts
			}
		} else {
			if resourceType == "snapshots" {
				azureData, err = azure.CollectSnapshotData(ctx)
			} else {
				azureData, err = ConvertAzureDataForCostAnalysis(ctx)
			}

			if err != nil {
				log.Printf("Warning: Error collecting Azure resources: %v", err)
			} else {
				log.Printf("Azure resource data collected for 'all' option")

				if !isEmpty(azureData) {
					if resourceType == "snapshots" {
						azureCosts, _ := EstimateAzureResourceCosts(azureData)
						allCosts["azure"] = azureCosts
					} else {
						estimator := NewResourceCostEstimator()
						azureCosts, _ := estimator.EstimateAzureResourcesCost(azureData)
						allCosts["azure"] = azureCosts
					}
				} else {
					log.Printf("No Azure resources found")
					allCosts["azure"] = map[string]interface{}{
						"Summary": map[string]interface{}{
							"TotalSnapshotStorage": 0.0,
							"TotalMonthlyCost":     0.0,
							"Currency":             "USD",
						},
						"Message": "No Azure resources found. Real data is being shown.",
					}
				}
			}
		}

		var gcpData map[string]interface{}
		if useMock {
			gcpData = GenerateMockResourceData("gcp", resourceType)
			log.Printf("Using mock GCP data for 'all' option (%s)", resourceType)

			if resourceType == "snapshots" {
				gcpCosts, _ := EstimateGcpResourceCosts(gcpData)
				allCosts["gcp"] = gcpCosts
			} else {
				estimator := NewResourceCostEstimator()
				gcpCosts, _ := estimator.EstimateGcpResourcesCost(gcpData)
				allCosts["gcp"] = gcpCosts
			}
		} else {
			if resourceType == "snapshots" {
				gcpData, err = gcp.CollectSnapshotData(ctx)
			} else {
				gcpData, err = ConvertGcpDataForCostAnalysis(ctx)
			}

			if err != nil {
				log.Printf("Warning: Error collecting GCP resources: %v", err)
			} else {
				log.Printf("GCP resource data collected for 'all' option")

				if !isEmpty(gcpData) {
					if resourceType == "snapshots" {
						gcpCosts, _ := EstimateGcpResourceCosts(gcpData)
						allCosts["gcp"] = gcpCosts
					} else {
						estimator := NewResourceCostEstimator()
						gcpCosts, _ := estimator.EstimateGcpResourcesCost(gcpData)
						allCosts["gcp"] = gcpCosts
					}
				} else {
					log.Printf("No GCP resources found")
					allCosts["gcp"] = map[string]interface{}{
						"Summary": map[string]interface{}{
							"TotalSnapshotStorage": 0.0,
							"TotalMonthlyCost":     0.0,
							"Currency":             "USD",
						},
						"Message": "No GCP resources found. Real data is being shown.",
					}
				}
			}
		}

		var totalStorage float64
		var totalCost float64

		for _, platformCosts := range allCosts {
			if costMap, ok := platformCosts.(map[string]interface{}); ok {
				if summary, ok := costMap["Summary"].(map[string]interface{}); ok {
					if storage, ok := summary["TotalSnapshotStorage"].(float64); ok {
						totalStorage += storage
					}

					if cost, ok := summary["TotalMonthlyCost"].(float64); ok {
						totalCost += cost
					}
				}
			}
		}

		allCosts["GlobalSummary"] = map[string]interface{}{
			"TotalSnapshotStorage": totalStorage,
			"TotalMonthlyCost":     totalCost,
			"Currency":             "USD",
		}

		costs = allCosts

	default:
		http.Error(w, "Invalid platform specified", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Error estimating costs: %v", err), http.StatusInternalServerError)
		return
	}

	disclaimer := GetPricingDisclaimer()

	result := map[string]interface{}{
		"costs":      costs,
		"disclaimer": disclaimer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func countResources(data map[string]interface{}, key string) int {
	if resources, ok := data[key].([]interface{}); ok {
		return len(resources)
	}
	if resources, ok := data[key].([]map[string]string); ok {
		return len(resources)
	}
	if resources, ok := data[key].([]map[string]interface{}); ok {
		return len(resources)
	}
	return 0
}

func isEmpty(data map[string]interface{}) bool {
	if data == nil {
		return true
	}

	for _, key := range []string{
		"EBSSnapshots", "RDSSnapshots", "EC2Instances", "RDSInstances", "S3Buckets", "DynamoDBTables", "VPCs",
		"DiskSnapshots", "VirtualMachines", "StorageAccounts", "SQLDatabases",
		"DiskSnapshots", "ComputeInstances", "GCSBuckets", "CloudSQLInstances", "CloudRunServices", "CloudFunctions"} {

		if resources := countResources(data, key); resources > 0 {
			return false
		}
	}

	return true
}

func GenerateMockResourceData(platform, resourceType string) map[string]interface{} {
	if resourceType == "snapshots" {
		return GenerateMockSnapshotData(platform)
	}

	data := make(map[string]interface{})

	switch platform {
	case "aws":
		data["EC2Instances"] = []map[string]interface{}{
			{
				"InstanceId":   "i-0123456789abcdef0",
				"InstanceType": "t3.medium",
				"Region":       "us-east-1",
				"State":        "running",
				"LaunchTime":   "2023-01-15T10:00:00Z",
			},
			{
				"InstanceId":   "i-0123456789abcdef1",
				"InstanceType": "m5.large",
				"Region":       "us-west-2",
				"State":        "running",
				"LaunchTime":   "2023-02-20T14:30:00Z",
			},
		}

		data["RDSInstances"] = []map[string]interface{}{
			{
				"DBInstanceIdentifier": "database-1",
				"DBInstanceClass":      "db.t3.medium",
				"Engine":               "mysql",
				"Region":               "us-east-1",
				"Status":               "available",
				"AllocatedStorage":     100,
			},
		}

		data["S3Buckets"] = []map[string]interface{}{
			{
				"Name":         "my-important-bucket",
				"Region":       "us-east-1",
				"SizeGB":       1024,
				"StorageClass": "STANDARD",
				"CreationDate": "2022-10-05T08:40:00Z",
			},
			{
				"Name":         "backup-bucket",
				"Region":       "us-west-2",
				"SizeGB":       2048,
				"StorageClass": "STANDARD_IA",
				"CreationDate": "2022-12-12T11:20:00Z",
			},
		}

		snapshots := GenerateMockSnapshotData(platform)
		for k, v := range snapshots {
			data[k] = v
		}

	case "azure":
		data["VirtualMachines"] = []map[string]interface{}{
			{
				"Name":          "vm-prod-app1",
				"ResourceGroup": "production-rg",
				"Location":      "eastus",
				"VMSize":        "Standard_D2s_v3",
				"PowerState":    "running",
			},
			{
				"Name":          "vm-dev-app2",
				"ResourceGroup": "development-rg",
				"Location":      "westeurope",
				"VMSize":        "Standard_B2s",
				"PowerState":    "running",
			},
		}

		data["StorageAccounts"] = []map[string]interface{}{
			{
				"Name":            "prodstorageacct",
				"ResourceGroup":   "production-rg",
				"Location":        "eastus",
				"AccountTier":     "Standard",
				"ReplicationType": "LRS",
				"UsedCapacityGB":  256,
			},
		}

		data["SQLDatabases"] = []map[string]interface{}{
			{
				"Name":          "prod-db",
				"ResourceGroup": "production-rg",
				"ServerName":    "prod-sqlserver",
				"Location":      "eastus",
				"Edition":       "Standard",
				"DTU":           100,
				"MaxSizeBytes":  107.374,
				"Status":        "Online",
			},
		}

		snapshots := GenerateMockSnapshotData(platform)
		for k, v := range snapshots {
			data[k] = v
		}

	case "gcp":
		data["ComputeInstances"] = []map[string]interface{}{
			{
				"Name":        "instance-1",
				"Project":     "my-project",
				"Zone":        "us-central1-a",
				"MachineType": "n1-standard-2",
				"Status":      "RUNNING",
			},
			{
				"Name":        "instance-2",
				"Project":     "my-project",
				"Zone":        "us-east1-b",
				"MachineType": "e2-medium",
				"Status":      "RUNNING",
			},
		}

		data["CloudSQLInstances"] = []map[string]interface{}{
			{
				"Name":            "sql-instance-1",
				"Project":         "my-project",
				"Region":          "us-central1",
				"Tier":            "db-n1-standard-1",
				"DatabaseVersion": "MYSQL_5_7",
				"DiskSizeGB":      100,
				"State":           "RUNNABLE",
			},
		}

		data["GCSBuckets"] = []map[string]interface{}{
			{
				"Name":         "my-important-bucket",
				"Project":      "my-project",
				"Location":     "US",
				"StorageClass": "STANDARD",
				"SizeGB":       1024,
			},
			{
				"Name":         "archival-bucket",
				"Project":      "my-project",
				"Location":     "US-CENTRAL1",
				"StorageClass": "NEARLINE",
				"SizeGB":       2048,
			},
		}

		snapshots := GenerateMockSnapshotData(platform)
		for k, v := range snapshots {
			data[k] = v
		}
	}

	return data
}
