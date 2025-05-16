package cost

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/michaelcade/kollect/pkg/aws"
	"github.com/michaelcade/kollect/pkg/azure"
	"github.com/michaelcade/kollect/pkg/gcp"
)

// HandleCostRequest processes cost calculation requests for cloud snapshots
func HandleCostRequest(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	useMock := r.URL.Query().Get("mock") == "true"

	log.Printf("Cost request received for platform: %s (mock: %v)", platform, useMock)

	ctx := r.Context()
	var costs interface{}
	var err error

	switch platform {
	case "aws":
		var awsData map[string]interface{}

		if useMock {
			awsData = GenerateMockSnapshotData("aws")
			log.Printf("Using mock AWS data")
		} else {
			awsData, err = aws.CollectSnapshotData(ctx)
			if err != nil {
				log.Printf("Error collecting AWS snapshot data: %v", err)
				http.Error(w, fmt.Sprintf("Error collecting AWS snapshot data: %v", err), http.StatusInternalServerError)
				return
			}

			log.Printf("AWS snapshot data collected with %d EBS snapshots and %d RDS snapshots",
				countSnapshots(awsData, "EBSSnapshots"),
				countSnapshots(awsData, "RDSSnapshots"))

			// If no snapshots found, return empty data with a message instead of using mock data
			if isEmpty(awsData) {
				log.Printf("No AWS snapshots found")
				costs = map[string]interface{}{
					"Summary": map[string]interface{}{
						"TotalSnapshotStorage": 0.0,
						"TotalMonthlyCost":     0.0,
						"Currency":             "USD",
					},
					"Message": "No AWS snapshots found. Real data is being shown.",
				}
				costs = map[string]interface{}{"aws": costs}
				break
			}
		}

		costs, err = EstimateAwsResourceCosts(awsData)
		costs = map[string]interface{}{"aws": costs}

	case "azure":
		var azureData map[string]interface{}

		if useMock {
			azureData = GenerateMockSnapshotData("azure")
			log.Printf("Using mock Azure data")
		} else {
			azureData, err = azure.CollectSnapshotData(ctx)
			if err != nil {
				log.Printf("Error collecting Azure snapshot data: %v", err)
				http.Error(w, fmt.Sprintf("Error collecting Azure snapshot data: %v", err), http.StatusInternalServerError)
				return
			}

			log.Printf("Azure snapshot data collected with %d disk snapshots",
				countSnapshots(azureData, "DiskSnapshots"))

			// If no snapshots found, return empty data with a message instead of using mock data
			if isEmpty(azureData) {
				log.Printf("No Azure snapshots found")
				costs = map[string]interface{}{
					"Summary": map[string]interface{}{
						"TotalSnapshotStorage": 0.0,
						"TotalMonthlyCost":     0.0,
						"Currency":             "USD",
					},
					"Message": "No Azure snapshots found. Real data is being shown.",
				}
				costs = map[string]interface{}{"azure": costs}
				break
			}
		}

		costs, err = EstimateAzureResourceCosts(azureData)
		costs = map[string]interface{}{"azure": costs}

	case "gcp":
		var gcpData map[string]interface{}

		if useMock {
			gcpData = GenerateMockSnapshotData("gcp")
			log.Printf("Using mock GCP data")
		} else {
			gcpData, err = gcp.CollectSnapshotData(ctx)
			if err != nil {
				log.Printf("Error collecting GCP snapshot data: %v", err)
				http.Error(w, fmt.Sprintf("Error collecting GCP snapshot data: %v", err), http.StatusInternalServerError)
				return
			}

			log.Printf("GCP snapshot data collected with %d disk snapshots",
				countSnapshots(gcpData, "DiskSnapshots"))

			// If no snapshots found, return empty data with a message instead of using mock data
			if isEmpty(gcpData) {
				log.Printf("No GCP snapshots found")
				costs = map[string]interface{}{
					"Summary": map[string]interface{}{
						"TotalSnapshotStorage": 0.0,
						"TotalMonthlyCost":     0.0,
						"Currency":             "USD",
					},
					"Message": "No GCP snapshots found. Real data is being shown.",
				}
				costs = map[string]interface{}{"gcp": costs}
				break
			}
		}

		costs, err = EstimateGcpResourceCosts(gcpData)
		costs = map[string]interface{}{"gcp": costs}

	case "all":
		allCosts := make(map[string]interface{})

		// AWS
		var awsData map[string]interface{}
		if useMock {
			awsData = GenerateMockSnapshotData("aws")
			log.Printf("Using mock AWS data for 'all' option")
			awsCosts, _ := EstimateAwsResourceCosts(awsData)
			allCosts["aws"] = awsCosts
		} else {
			awsData, err = aws.CollectSnapshotData(ctx)
			if err != nil {
				log.Printf("Warning: Error collecting AWS snapshots: %v", err)
			} else {
				log.Printf("AWS snapshot data collected with %d EBS snapshots and %d RDS snapshots",
					countSnapshots(awsData, "EBSSnapshots"),
					countSnapshots(awsData, "RDSSnapshots"))

				if !isEmpty(awsData) {
					awsCosts, _ := EstimateAwsResourceCosts(awsData)
					allCosts["aws"] = awsCosts
				} else {
					log.Printf("No AWS snapshots found")
					allCosts["aws"] = map[string]interface{}{
						"Summary": map[string]interface{}{
							"TotalSnapshotStorage": 0.0,
							"TotalMonthlyCost":     0.0,
							"Currency":             "USD",
						},
						"Message": "No AWS snapshots found. Real data is being shown.",
					}
				}
			}
		}

		// Azure
		var azureData map[string]interface{}
		if useMock {
			azureData = GenerateMockSnapshotData("azure")
			log.Printf("Using mock Azure data for 'all' option")
			azureCosts, _ := EstimateAzureResourceCosts(azureData)
			allCosts["azure"] = azureCosts
		} else {
			azureData, err = azure.CollectSnapshotData(ctx)
			if err != nil {
				log.Printf("Warning: Error collecting Azure snapshots: %v", err)
			} else {
				log.Printf("Azure snapshot data collected with %d disk snapshots",
					countSnapshots(azureData, "DiskSnapshots"))

				if !isEmpty(azureData) {
					azureCosts, _ := EstimateAzureResourceCosts(azureData)
					allCosts["azure"] = azureCosts
				} else {
					log.Printf("No Azure snapshots found")
					allCosts["azure"] = map[string]interface{}{
						"Summary": map[string]interface{}{
							"TotalSnapshotStorage": 0.0,
							"TotalMonthlyCost":     0.0,
							"Currency":             "USD",
						},
						"Message": "No Azure snapshots found. Real data is being shown.",
					}
				}
			}
		}

		// GCP
		var gcpData map[string]interface{}
		if useMock {
			gcpData = GenerateMockSnapshotData("gcp")
			log.Printf("Using mock GCP data for 'all' option")
			gcpCosts, _ := EstimateGcpResourceCosts(gcpData)
			allCosts["gcp"] = gcpCosts
		} else {
			gcpData, err = gcp.CollectSnapshotData(ctx)
			if err != nil {
				log.Printf("Warning: Error collecting GCP snapshots: %v", err)
			} else {
				log.Printf("GCP snapshot data collected with %d disk snapshots",
					countSnapshots(gcpData, "DiskSnapshots"))

				if !isEmpty(gcpData) {
					gcpCosts, _ := EstimateGcpResourceCosts(gcpData)
					allCosts["gcp"] = gcpCosts
				} else {
					log.Printf("No GCP snapshots found")
					allCosts["gcp"] = map[string]interface{}{
						"Summary": map[string]interface{}{
							"TotalSnapshotStorage": 0.0,
							"TotalMonthlyCost":     0.0,
							"Currency":             "USD",
						},
						"Message": "No GCP snapshots found. Real data is being shown.",
					}
				}
			}
		}

		// Calculate global summary
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

// InitPricing initializes pricing data at application startup
func InitPricing() {
	// Attempt to initialize pricing from APIs at startup
	go func() {
		log.Println("Initializing pricing data...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := RefreshPricing(ctx); err != nil {
			log.Printf("Warning: Error initializing pricing data: %v", err)
		} else {
			log.Println("Pricing data initialized successfully")
		}
	}()
}

// Helper function to count snapshots
func countSnapshots(data map[string]interface{}, key string) int {
	if snapshots, ok := data[key].([]interface{}); ok {
		return len(snapshots)
	}
	if snapshots, ok := data[key].([]map[string]string); ok {
		return len(snapshots)
	}
	return 0
}

// Helper function to check if snapshot data is empty
func isEmpty(data map[string]interface{}) bool {
	if len(data) == 0 {
		return true
	}

	for _, v := range data {
		if arr, ok := v.([]interface{}); ok && len(arr) > 0 {
			return false
		}
		if arr, ok := v.([]map[string]string); ok && len(arr) > 0 {
			return false
		}
		if arr, ok := v.([]map[string]interface{}); ok && len(arr) > 0 {
			return false
		}
	}

	return true
}
