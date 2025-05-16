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
			log.Printf("Using mock AWS data: %+v", awsData)
		} else {
			awsData, err = aws.CollectSnapshotData(ctx)
			if err != nil {
				log.Printf("Error collecting AWS snapshot data: %v", err)
				http.Error(w, fmt.Sprintf("Error collecting AWS snapshot data: %v", err), http.StatusInternalServerError)
				return
			}
		}

		log.Printf("AWS snapshot data collected with %d EBS snapshots and %d RDS snapshots",
			countSnapshots(awsData, "EBSSnapshots"),
			countSnapshots(awsData, "RDSSnapshots"))

		if isEmpty(awsData) {
			log.Printf("No AWS snapshots found, using mock data")
			awsData = GenerateMockSnapshotData("aws")
		}

		costs, err = EstimateAwsResourceCosts(awsData)
		log.Printf("AWS costs calculated: %+v", costs)

		// Wrap the costs in a platform map for the client
		costs = map[string]interface{}{"aws": costs}
		log.Printf("Final AWS costs object: %+v", costs)

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
		}

		log.Printf("Azure snapshot data collected with %d disk snapshots",
			countSnapshots(azureData, "DiskSnapshots"))

		if isEmpty(azureData) {
			log.Printf("No Azure snapshots found, using mock data")
			azureData = GenerateMockSnapshotData("azure")
		}

		costs, err = EstimateAzureResourceCosts(azureData)
		// Wrap the costs in a platform map for the client
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
		}

		log.Printf("GCP snapshot data collected with %d disk snapshots",
			countSnapshots(gcpData, "DiskSnapshots"))

		if isEmpty(gcpData) {
			log.Printf("No GCP snapshots found, using mock data")
			gcpData = GenerateMockSnapshotData("gcp")
		}

		costs, err = EstimateGcpResourceCosts(gcpData)
		// Wrap the costs in a platform map for the client
		costs = map[string]interface{}{"gcp": costs}

	case "all":
		allCosts := make(map[string]interface{})

		// AWS
		var awsData map[string]interface{}
		if useMock {
			awsData = GenerateMockSnapshotData("aws")
			log.Printf("Using mock AWS data for 'all' option")
		} else {
			awsData, err = aws.CollectSnapshotData(ctx)
			if err != nil {
				log.Printf("Warning: Error collecting AWS snapshots: %v", err)
			} else if isEmpty(awsData) {
				log.Printf("No AWS snapshots found, using mock data")
				awsData = GenerateMockSnapshotData("aws")
			}
		}

		if len(awsData) > 0 {
			awsCosts, _ := EstimateAwsResourceCosts(awsData)
			allCosts["aws"] = awsCosts
		}

		// Azure
		var azureData map[string]interface{}
		if useMock {
			azureData = GenerateMockSnapshotData("azure")
			log.Printf("Using mock Azure data for 'all' option")
		} else {
			azureData, err = azure.CollectSnapshotData(ctx)
			if err != nil {
				log.Printf("Warning: Error collecting Azure snapshots: %v", err)
			} else if isEmpty(azureData) {
				log.Printf("No Azure snapshots found, using mock data")
				azureData = GenerateMockSnapshotData("azure")
			}
		}

		if len(azureData) > 0 {
			azureCosts, _ := EstimateAzureResourceCosts(azureData)
			allCosts["azure"] = azureCosts
		}

		// GCP
		var gcpData map[string]interface{}
		if useMock {
			gcpData = GenerateMockSnapshotData("gcp")
			log.Printf("Using mock GCP data for 'all' option")
		} else {
			gcpData, err = gcp.CollectSnapshotData(ctx)
			if err != nil {
				log.Printf("Warning: Error collecting GCP snapshots: %v", err)
			} else if isEmpty(gcpData) {
				log.Printf("No GCP snapshots found, using mock data")
				gcpData = GenerateMockSnapshotData("gcp")
			}
		}

		if len(gcpData) > 0 {
			gcpCosts, _ := EstimateGcpResourceCosts(gcpData)
			allCosts["gcp"] = gcpCosts
		}

		// Calculate global summary
		var totalStorage float64
		var totalCost float64

		if awsCosts, ok := allCosts["aws"].(map[string]interface{}); ok {
			if summary, ok := awsCosts["Summary"].(map[string]interface{}); ok {
				if storage, ok := summary["TotalSnapshotStorage"].(float64); ok {
					totalStorage += storage
				}
				if cost, ok := summary["TotalMonthlyCost"].(float64); ok {
					totalCost += cost
				}
			}
		}

		if azureCosts, ok := allCosts["azure"].(map[string]interface{}); ok {
			if summary, ok := azureCosts["Summary"].(map[string]interface{}); ok {
				if storage, ok := summary["TotalSnapshotStorage"].(float64); ok {
					totalStorage += storage
				}
				if cost, ok := summary["TotalMonthlyCost"].(float64); ok {
					totalCost += cost
				}
			}
		}

		if gcpCosts, ok := allCosts["gcp"].(map[string]interface{}); ok {
			if summary, ok := gcpCosts["Summary"].(map[string]interface{}); ok {
				if storage, ok := summary["TotalSnapshotStorage"].(float64); ok {
					totalStorage += storage
				}
				if cost, ok := summary["TotalMonthlyCost"].(float64); ok {
					totalCost += cost
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

	// Add disclaimer
	disclaimer := "Cost estimates are approximations based on publicly available pricing information. " +
		"Actual costs may vary based on your specific agreements, reserved capacity, and other factors."

	result := map[string]interface{}{
		"costs":      costs,
		"disclaimer": disclaimer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
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

	hasSnapshots := false
	for _, key := range []string{"EBSSnapshots", "RDSSnapshots", "DiskSnapshots"} {
		if snapshots, ok := data[key].([]interface{}); ok && len(snapshots) > 0 {
			hasSnapshots = true
			break
		}
		if snapshots, ok := data[key].([]map[string]string); ok && len(snapshots) > 0 {
			hasSnapshots = true
			break
		}
	}

	return !hasSnapshots
}
