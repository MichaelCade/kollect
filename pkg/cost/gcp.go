package cost

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/michaelcade/kollect/pkg/gcp"
)

// CalculateGcpDiskSnapshotCosts calculates costs for GCP disk snapshots
func CalculateGcpDiskSnapshotCosts(snapshots []map[string]string) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(snapshots))

	for _, snapshot := range snapshots {
		// Parse size
		var sizeGB float64 = 0
		if sizeStr, ok := snapshot["DiskSizeGB"]; ok && sizeStr != "" {
			// Clean up the string to extract just the number
			numStr := strings.TrimSpace(sizeStr)
			numStr = strings.Split(numStr, " ")[0] // Take the first part before any space

			// Try to parse it as a float
			if val, err := strconv.ParseFloat(numStr, 64); err == nil {
				sizeGB = val
			}
		}

		// If size is still 0, use a default size
		if sizeGB == 0 {
			sizeGB = 100.0 // Default to 100 GB
		}

		// Determine region or use default
		region := "us-central1" // Default region
		if location, ok := snapshot["Location"]; ok && location != "" {
			region = strings.ToLower(location)
		}

		// Calculate monthly cost
		pricePerGB := GetPrice("gcp", "disk_snapshot", region)
		monthlyCost := sizeGB * pricePerGB

		// Get pricing source and metadata
		priceSource := GetPricingSource("gcp", "disk_snapshot")
		priceInfo := GetPricingMetadata("gcp", "disk_snapshot")

		// Log price information for debugging
		log.Printf("GCP disk pricing for region %s: $%.4f per GB/month (Source: %s, Last verified: %s)",
			region, pricePerGB, priceSource, priceInfo.LastVerified.Format("2006-01-02"))

		// Create result with cost info
		result := map[string]interface{}{
			"Name":            snapshot["Name"],
			"SizeGB":          sizeGB,
			"Region":          region,
			"Status":          snapshot["Status"],
			"CreationTime":    snapshot["CreationTimestamp"],
			"PricePerGBMonth": pricePerGB,
			"PriceSource":     priceSource,
			"LastVerified":    priceInfo.LastVerified.Format("2006-01-02"),
			"MonthlyCost":     monthlyCost,
			"MonthlyCostUSD":  fmt.Sprintf("$%.2f", monthlyCost),
		}

		results = append(results, result)
	}

	return results
}

// EstimateGcpResourceCosts calculates costs for GCP resources
func EstimateGcpResourceCosts(resourceData map[string]interface{}) (map[string]interface{}, error) {
	costData := make(map[string]interface{})

	// Disk Snapshots
	if diskSnapshotsRaw, ok := resourceData["DiskSnapshots"]; ok {
		if diskSnapshots, ok := convertToSnapshotList(diskSnapshotsRaw); ok {
			costData["DiskSnapshotCosts"] = CalculateGcpDiskSnapshotCosts(diskSnapshots)
		}
	}

	// Calculate summary
	var totalSnapshotStorage float64
	var totalMonthlyCost float64

	if diskCosts, ok := costData["DiskSnapshotCosts"].([]map[string]interface{}); ok {
		for _, cost := range diskCosts {
			if storage, ok := cost["SizeGB"].(float64); ok {
				totalSnapshotStorage += storage
			}
			if monthlyCost, ok := cost["MonthlyCost"].(float64); ok {
				totalMonthlyCost += monthlyCost
			}
		}
	}

	// Include pricing source information in the summary
	priceSource := GetPricingSource("gcp", "disk_snapshot")
	lastVerified := GetPricingMetadata("gcp", "disk_snapshot").LastVerified.Format("2006-01-02")

	costData["Summary"] = map[string]interface{}{
		"TotalSnapshotStorage": totalSnapshotStorage,
		"TotalMonthlyCost":     totalMonthlyCost,
		"Currency":             "USD",
		"PriceSource":          priceSource,
		"LastVerified":         lastVerified,
	}

	return costData, nil
}

func getGCPProjectID() (string, error) {
	projectID := gcp.GetProjectID()
	if projectID == "" {
		return "", fmt.Errorf("could not determine GCP project ID")
	}
	return projectID, nil
}

func ConvertGcpDataForCostAnalysis(ctx context.Context) (map[string]interface{}, error) {
	// Use the existing GCP inventory collection function
	gcpData, err := gcp.CollectGCPData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect GCP data: %v", err)
	}

	// Create a generic map to hold all resource types
	inventory := make(map[string]interface{})

	// Convert Compute Instances to generic map entries
	if len(gcpData.ComputeInstances) > 0 {
		instances := make([]map[string]interface{}, len(gcpData.ComputeInstances))
		for i, instance := range gcpData.ComputeInstances {
			instances[i] = map[string]interface{}{
				"Name":        instance.Name,
				"Zone":        instance.Zone,
				"MachineType": instance.MachineType,
				"Status":      instance.Status,
				"Project":     instance.Project,
			}
		}
		inventory["ComputeInstances"] = instances
	}

	// Convert GCS Buckets to generic map entries
	if len(gcpData.GCSBuckets) > 0 {
		buckets := make([]map[string]interface{}, len(gcpData.GCSBuckets))
		for i, bucket := range gcpData.GCSBuckets {
			buckets[i] = map[string]interface{}{
				"Name":         bucket.Name,
				"Location":     bucket.Location,
				"StorageClass": bucket.StorageClass,
				"Project":      bucket.Project,
				"SizeGB":       100.0, // Default size estimate, would need to be calculated
			}
		}
		inventory["GCSBuckets"] = buckets
	}

	// Convert Cloud SQL Instances to generic map entries
	if len(gcpData.CloudSQLInstances) > 0 {
		instances := make([]map[string]interface{}, len(gcpData.CloudSQLInstances))
		for i, instance := range gcpData.CloudSQLInstances {
			instances[i] = map[string]interface{}{
				"Name":            instance.Name,
				"DatabaseVersion": instance.DatabaseVersion,
				"Region":          instance.Region,
				"Tier":            instance.Tier,
				"Status":          instance.Status,
				"Project":         instance.Project,
				"DiskSizeGB":      100.0, // Default size estimate, would need to be calculated
			}
		}
		inventory["CloudSQLInstances"] = instances
	}

	// Convert Cloud Run Services to generic map entries
	if len(gcpData.CloudRunServices) > 0 {
		services := make([]map[string]interface{}, len(gcpData.CloudRunServices))
		for i, service := range gcpData.CloudRunServices {
			services[i] = map[string]interface{}{
				"Name":      service.Name,
				"Region":    service.Region,
				"URL":       service.URL,
				"Project":   service.Project,
				"Replicas":  service.Replicas,
				"Container": service.Container,
			}
		}
		inventory["CloudRunServices"] = services
	}

	// Convert Cloud Functions to generic map entries
	if len(gcpData.CloudFunctions) > 0 {
		functions := make([]map[string]interface{}, len(gcpData.CloudFunctions))
		for i, function := range gcpData.CloudFunctions {
			functions[i] = map[string]interface{}{
				"Name":            function.Name,
				"Region":          function.Region,
				"Runtime":         function.Runtime,
				"Status":          function.Status,
				"EntryPoint":      function.EntryPoint,
				"AvailableMemory": function.AvailableMemory,
				"Project":         function.Project,
			}
		}
		inventory["CloudFunctions"] = functions
	}

	// Include snapshot data
	snapshotData, err := gcp.CollectSnapshotData(ctx)
	if err != nil {
		log.Printf("Warning: Failed to collect snapshot data: %v", err)
	} else {
		for k, v := range snapshotData {
			inventory[k] = v
		}
	}

	return inventory, nil
}
