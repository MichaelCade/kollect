package cost

import (
	"fmt"
	"strconv"
)

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

	costData["Summary"] = map[string]interface{}{
		"TotalSnapshotStorage": totalSnapshotStorage,
		"TotalMonthlyCost":     totalMonthlyCost,
		"Currency":             "USD",
	}

	return costData, nil
}

// CalculateAzureDiskSnapshotCosts calculates costs for Azure disk snapshots
func CalculateAzureDiskSnapshotCosts(snapshots []map[string]string) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(snapshots))

	for _, snapshot := range snapshots {
		// Parse size
		var sizeGB float64 = 0
		if sizeStr, ok := snapshot["DiskSizeGB"]; ok && sizeStr != "" {
			if val, err := strconv.ParseFloat(sizeStr, 64); err == nil {
				sizeGB = val
			}
		}

		// Get region/location
		region := "eastus" // Default region
		if location, ok := snapshot["Location"]; ok && location != "" {
			region = location
		}

		// Calculate monthly cost
		pricePerGB := GetPrice(AzureDiskSnapshotPricing, region)
		monthlyCost := sizeGB * pricePerGB

		// Create result with cost info
		result := map[string]interface{}{
			"Name":            snapshot["Name"],
			"SizeGB":          sizeGB,
			"Region":          region,
			"State":           snapshot["ProvisioningState"],
			"CreationTime":    snapshot["TimeCreated"],
			"PricePerGBMonth": pricePerGB,
			"MonthlyCost":     monthlyCost,
			"MonthlyCostUSD":  fmt.Sprintf("$%.2f", monthlyCost),
		}

		results = append(results, result)
	}

	return results
}
