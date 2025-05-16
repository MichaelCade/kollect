package cost

import (
	"fmt"
	"strconv"
	"strings"
)

// First, let's improve the CalculateEbsSnapshotCosts function
func CalculateEbsSnapshotCosts(snapshots []map[string]string) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(snapshots))

	for _, snapshot := range snapshots {
		// Parse size from format like "100 GiB" or just "100"
		var sizeGB float64 = 0
		if sizeStr, ok := snapshot["VolumeSize"]; ok && sizeStr != "" {
			// First, clean up the string to extract just the number
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
		region := "us-east-1" // Default region
		if snapshotId, ok := snapshot["SnapshotId"]; ok && strings.Contains(snapshotId, ":") {
			parts := strings.Split(snapshotId, ":")
			if len(parts) >= 4 {
				region = parts[3]
			}
		}

		// Calculate monthly cost
		pricePerGB := GetPrice(AwsEbsSnapshotPricing, region)
		monthlyCost := sizeGB * pricePerGB

		// Create result with cost info
		result := map[string]interface{}{
			"SnapshotId":      snapshot["SnapshotId"],
			"VolumeId":        snapshot["VolumeId"],
			"SizeGB":          sizeGB,
			"Region":          region,
			"State":           snapshot["State"],
			"CreationTime":    snapshot["StartTime"],
			"PricePerGBMonth": pricePerGB,
			"MonthlyCost":     monthlyCost,
			"MonthlyCostUSD":  fmt.Sprintf("$%.2f", monthlyCost),
		}

		results = append(results, result)
	}

	return results
}

// Now, let's do the same for RDS snapshots
func CalculateRdsSnapshotCosts(snapshots []map[string]string) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(snapshots))

	for _, snapshot := range snapshots {
		// Parse size
		var sizeGB float64 = 0
		if sizeStr, ok := snapshot["AllocatedStorage"]; ok && sizeStr != "" {
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
		region := "us-east-1" // Default region
		if arn, ok := snapshot["SnapshotId"]; ok && strings.Contains(arn, ":") {
			parts := strings.Split(arn, ":")
			if len(parts) >= 4 {
				region = parts[3]
			}
		}

		// Calculate monthly cost - using RDS-specific pricing
		pricePerGB := GetPrice(AwsRdsSnapshotPricing, region)
		monthlyCost := sizeGB * pricePerGB

		// Create result with cost info
		result := map[string]interface{}{
			"SnapshotId":      snapshot["SnapshotId"],
			"Engine":          snapshot["Engine"],
			"SizeGB":          sizeGB,
			"Region":          region,
			"Status":          snapshot["Status"],
			"CreationTime":    snapshot["SnapshotCreateTime"],
			"PricePerGBMonth": pricePerGB,
			"MonthlyCost":     monthlyCost,
			"MonthlyCostUSD":  fmt.Sprintf("$%.2f", monthlyCost),
		}

		results = append(results, result)
	}

	return results
}

// Now, let's fix the main function
func EstimateAwsResourceCosts(resourceData map[string]interface{}) (map[string]interface{}, error) {
	costData := make(map[string]interface{})

	// EBS Snapshots
	if ebsSnapshotsRaw, ok := resourceData["EBSSnapshots"]; ok {
		if ebsSnapshots, ok := convertToSnapshotList(ebsSnapshotsRaw); ok {
			costData["EBSSnapshotCosts"] = CalculateEbsSnapshotCosts(ebsSnapshots)
		}
	}

	// RDS Snapshots
	if rdsSnapshotsRaw, ok := resourceData["RDSSnapshots"]; ok {
		if rdsSnapshots, ok := convertToSnapshotList(rdsSnapshotsRaw); ok {
			costData["RDSSnapshotCosts"] = CalculateRdsSnapshotCosts(rdsSnapshots)
		}
	}

	// Calculate summary
	var totalSnapshotStorage float64
	var totalMonthlyCost float64

	if ebsCosts, ok := costData["EBSSnapshotCosts"].([]map[string]interface{}); ok {
		for _, cost := range ebsCosts {
			if storage, ok := cost["SizeGB"].(float64); ok {
				totalSnapshotStorage += storage
			}
			if monthlyCost, ok := cost["MonthlyCost"].(float64); ok {
				totalMonthlyCost += monthlyCost
			}
		}
	}

	if rdsCosts, ok := costData["RDSSnapshotCosts"].([]map[string]interface{}); ok {
		for _, cost := range rdsCosts {
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

// Also ensure the convertToSnapshotList function is working properly
func convertToSnapshotList(data interface{}) ([]map[string]string, bool) {
	// First, try direct type assertion
	if snapshots, ok := data.([]map[string]string); ok {
		return snapshots, true
	}

	// Next, try to handle it as an array of interfaces
	if snapshotsArray, ok := data.([]interface{}); ok {
		result := make([]map[string]string, 0, len(snapshotsArray))
		for _, item := range snapshotsArray {
			if mapItem, ok := item.(map[string]interface{}); ok {
				strMap := make(map[string]string)
				for k, v := range mapItem {
					strMap[k] = fmt.Sprintf("%v", v)
				}
				result = append(result, strMap)
			} else if stringMap, ok := item.(map[string]string); ok {
				result = append(result, stringMap)
			}
		}
		return result, len(result) > 0
	}

	return nil, false
}
