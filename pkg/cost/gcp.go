package cost

import (
	"fmt"
	"strconv"
	"strings"
)

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

	costData["Summary"] = map[string]interface{}{
		"TotalSnapshotStorage": totalSnapshotStorage,
		"TotalMonthlyCost":     totalMonthlyCost,
		"Currency":             "USD",
	}

	return costData, nil
}

// CalculateGcpDiskSnapshotCosts calculates costs for GCP disk snapshots
func CalculateGcpDiskSnapshotCosts(snapshots []map[string]string) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(snapshots))

	for _, snapshot := range snapshots {
		// Parse size from DiskSizeGB
		var sizeGB float64 = 0
		if sizeStr, ok := snapshot["DiskSizeGB"]; ok && sizeStr != "" {
			if val, err := strconv.ParseFloat(sizeStr, 64); err == nil {
				sizeGB = val
			}
		}

		// Determine region by extracting from zone or location
		region := "us-central1" // Default region

		// First try to get it from Location
		if location, ok := snapshot["Location"]; ok && location != "" {
			region = strings.ToLower(location)
		}

		// If not found, try to extract from zone
		if zone, ok := snapshot["Zone"]; ok && zone != "" {
			parts := strings.Split(zone, "-")
			if len(parts) >= 2 {
				region = strings.Join(parts[0:2], "-")
			}
		}

		// Calculate monthly cost
		pricePerGB := GetPrice(GcpDiskSnapshotPricing, region)
		monthlyCost := sizeGB * pricePerGB

		// Create result with cost info
		result := map[string]interface{}{
			"Name":            snapshot["Name"],
			"SizeGB":          sizeGB,
			"Region":          region,
			"Status":          snapshot["Status"],
			"CreationTime":    snapshot["CreationTimestamp"],
			"PricePerGBMonth": pricePerGB,
			"MonthlyCost":     monthlyCost,
			"MonthlyCostUSD":  fmt.Sprintf("$%.2f", monthlyCost),
		}

		results = append(results, result)
	}

	return results
}
