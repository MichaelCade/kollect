package cost

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/michaelcade/kollect/pkg/azure"
)

func CalculateAzureDiskSnapshotCosts(snapshots []map[string]string) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(snapshots))

	for _, snapshot := range snapshots {
		var sizeGB float64 = 0
		if sizeStr, ok := snapshot["DiskSizeGB"]; ok && sizeStr != "" {
			numStr := strings.TrimSpace(sizeStr)
			numStr = strings.TrimSuffix(numStr, " GB")
			numStr = strings.Split(numStr, " ")[0]

			if val, err := strconv.ParseFloat(numStr, 64); err == nil {
				sizeGB = val
			}
		}

		if sizeGB == 0 {
			sizeGB = 100.0
		}

		region := "eastus"
		if location, ok := snapshot["Location"]; ok && location != "" {
			region = strings.ToLower(location)
		}

		pricePerGB := GetPrice("azure", "disk_snapshot", region)
		monthlyCost := sizeGB * pricePerGB

		priceSource := GetPricingSource("azure", "disk_snapshot")
		priceInfo := GetPricingMetadata("azure", "disk_snapshot")

		log.Printf("Azure disk snapshot pricing for region %s: $%.4f per GB/month (Source: %s, Last verified: %s)",
			region, pricePerGB, priceSource, priceInfo.LastVerified.Format("2006-01-02"))

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

func EstimateAzureResourceCosts(resourceData map[string]interface{}) (map[string]interface{}, error) {
	costData := make(map[string]interface{})

	if diskSnapshotsRaw, ok := resourceData["DiskSnapshots"]; ok {
		if diskSnapshots, ok := convertToSnapshotList(diskSnapshotsRaw); ok {
			costData["DiskSnapshotCosts"] = CalculateAzureDiskSnapshotCosts(diskSnapshots)
		}
	}

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

func ConvertAzureDataForCostAnalysis(ctx context.Context) (map[string]interface{}, error) {
	azureData, err := azure.CollectAzureData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect Azure data: %v", err)
	}

	inventory := make(map[string]interface{})

	if len(azureData.AzureVMs) > 0 {
		vms := make([]map[string]interface{}, len(azureData.AzureVMs))
		for i, vm := range azureData.AzureVMs {
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

	if len(azureData.AzureStorageAccounts) > 0 {
		accounts := make([]map[string]interface{}, len(azureData.AzureStorageAccounts))
		for i, account := range azureData.AzureStorageAccounts {
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
				"UsedCapacityGB": 100.0,
			}
		}
		inventory["StorageAccounts"] = accounts
	}

	if len(azureData.AzureSQLDatabases) > 0 {
		databases := make([]map[string]interface{}, len(azureData.AzureSQLDatabases))
		for i, db := range azureData.AzureSQLDatabases {
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
