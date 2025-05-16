package cost

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// CalculateAzureDiskSnapshotCosts calculates costs for Azure disk snapshots
func CalculateAzureDiskSnapshotCosts(snapshots []map[string]string) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(snapshots))

	for _, snapshot := range snapshots {
		// Parse size from format like "128" or "128 GB"
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
		region := "eastus" // Default region
		if location, ok := snapshot["Location"]; ok && location != "" {
			region = strings.ToLower(location)
		}

		// Calculate monthly cost
		pricePerGB := GetPrice("azure", "disk_snapshot", region)
		monthlyCost := sizeGB * pricePerGB

		// Get pricing source and metadata
		priceSource := GetPricingSource("azure", "disk_snapshot")
		priceInfo := GetPricingMetadata("azure", "disk_snapshot")

		// Log price information for debugging
		log.Printf("Azure disk pricing for region %s: $%.4f per GB/month (Source: %s, Last verified: %s)",
			region, pricePerGB, priceSource, priceInfo.LastVerified.Format("2006-01-02"))

		// Create result with cost info
		result := map[string]interface{}{
			"Name":            snapshot["Name"],
			"SizeGB":          sizeGB,
			"Region":          region,
			"State":           snapshot["ProvisioningState"],
			"CreationTime":    snapshot["TimeCreated"],
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

// EstimateAzureResourceCosts calculates costs for Azure resources
func EstimateAzureResourceCosts(resourceData map[string]interface{}) (map[string]interface{}, error) {
	costData := make(map[string]interface{})

	// Disk Snapshots
	if diskSnapshotsRaw, ok := resourceData["DiskSnapshots"]; ok {
		if diskSnapshots, ok := convertToSnapshotList(diskSnapshotsRaw); ok {
			costData["DiskSnapshotCosts"] = CalculateAzureDiskSnapshotCosts(diskSnapshots)
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
	priceSource := GetPricingSource("azure", "disk_snapshot")
	lastVerified := GetPricingMetadata("azure", "disk_snapshot").LastVerified.Format("2006-01-02")

	costData["Summary"] = map[string]interface{}{
		"TotalSnapshotStorage": totalSnapshotStorage,
		"TotalMonthlyCost":     totalMonthlyCost,
		"Currency":             "USD",
		"PriceSource":          priceSource,
		"LastVerified":         lastVerified,
	}

	return costData, nil
}
