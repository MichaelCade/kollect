package cost

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/michaelcade/kollect/pkg/azure"
)

// CalculateAzureDiskSnapshotCosts calculates costs for Azure disk snapshots
func CalculateAzureDiskSnapshotCosts(snapshots []map[string]string) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(snapshots))

	for _, snapshot := range snapshots {
		// Parse size
		var sizeGB float64 = 0
		if sizeStr, ok := snapshot["DiskSizeGB"]; ok && sizeStr != "" {
			// Clean up the string to extract just the number
			numStr := strings.TrimSpace(sizeStr)
			numStr = strings.TrimSuffix(numStr, " GB") // Remove unit if present
			numStr = strings.Split(numStr, " ")[0]     // Take the first part before any space

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
		log.Printf("Azure disk snapshot pricing for region %s: $%.4f per GB/month (Source: %s, Last verified: %s)",
			region, pricePerGB, priceSource, priceInfo.LastVerified.Format("2006-01-02"))

		// Create result with cost info
		result := map[string]interface{}{
			"Name":            snapshot["Name"],
			"ResourceGroup":   snapshot["ResourceGroup"],
			"SizeGB":          sizeGB,
			"Location":        region,
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

// ConvertAzureDataForCostAnalysis converts the structured AzureData into a generic map for cost analysis
func ConvertAzureDataForCostAnalysis(ctx context.Context) (map[string]interface{}, error) {
	// Use the existing Azure inventory collection function
	azureData, err := azure.CollectAzureData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect Azure data: %v", err)
	}

	// Create a generic map to hold all resource types
	inventory := make(map[string]interface{})

	// Convert Azure VMs to generic map entries
	if len(azureData.AzureVMs) > 0 {
		vms := make([]map[string]interface{}, len(azureData.AzureVMs))
		for i, vm := range azureData.AzureVMs {
			// Extract resource group from ID
			resourceGroup := "unknown"
			if vm.ID != nil {
				parts := strings.Split(*vm.ID, "/")
				for j := 0; j < len(parts); j++ {
					if parts[j] == "resourceGroups" && j+1 < len(parts) {
						resourceGroup = parts[j+1]
						break
					}
				}
			}

			vmSize := "unknown"
			if vm.Properties != nil && vm.Properties.HardwareProfile != nil && vm.Properties.HardwareProfile.VMSize != nil {
				vmSize = fmt.Sprint(*vm.Properties.HardwareProfile.VMSize)
			}

			name := "unknown"
			if vm.Name != nil {
				name = *vm.Name
			}
			location := "unknown"
			if vm.Location != nil {
				location = *vm.Location
			}
			vms[i] = map[string]interface{}{
				"Name":          name,
				"ResourceGroup": resourceGroup,
				"Location":      location,
				"VMSize":        vmSize,
			}
		}
		inventory["VirtualMachines"] = vms
	}

	// Convert Storage Accounts to generic map entries
	if len(azureData.AzureStorageAccounts) > 0 {
		accounts := make([]map[string]interface{}, len(azureData.AzureStorageAccounts))
		for i, account := range azureData.AzureStorageAccounts {
			// Extract resource group from ID
			resourceGroup := "unknown"
			if account.ID != nil {
				parts := strings.Split(*account.ID, "/")
				for j := 0; j < len(parts); j++ {
					if parts[j] == "resourceGroups" && j+1 < len(parts) {
						resourceGroup = parts[j+1]
						break
					}
				}
			}

			name := "unknown"
			if account.Name != nil {
				name = *account.Name
			}
			location := "unknown"
			if account.Location != nil {
				location = *account.Location
			}
			accounts[i] = map[string]interface{}{
				"Name":           name,
				"ResourceGroup":  resourceGroup,
				"Location":       location,
				"UsedCapacityGB": 100.0, // Default estimate
			}
		}
		inventory["StorageAccounts"] = accounts
	}

	// Convert SQL Databases to generic map entries
	if len(azureData.AzureSQLDatabases) > 0 {
		databases := make([]map[string]interface{}, len(azureData.AzureSQLDatabases))
		for i, db := range azureData.AzureSQLDatabases {
			// Extract resource group from ID
			resourceGroup := "unknown"
			if db.ID != nil {
				parts := strings.Split(*db.ID, "/")
				for j := 0; j < len(parts); j++ {
					if parts[j] == "resourceGroups" && j+1 < len(parts) {
						resourceGroup = parts[j+1]
						break
					}
				}
			}

			name := "unknown"
			if db.Name != nil {
				name = *db.Name
			}
			location := "unknown"
			if db.Location != nil {
				location = *db.Location
			}
			databases[i] = map[string]interface{}{
				"Name":          name,
				"ResourceGroup": resourceGroup,
				"Location":      location,
			}
		}
		inventory["SQLDatabases"] = databases
	}

	// Include snapshot data
	snapshotData, err := azure.CollectSnapshotData(ctx)
	if err != nil {
		log.Printf("Warning: Failed to collect snapshot data: %v", err)
	} else {
		for k, v := range snapshotData {
			inventory[k] = v
		}
	}

	return inventory, nil
}
