package cost

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/michaelcade/kollect/pkg/gcp"
)

func CalculateGcpDiskSnapshotCosts(snapshots []map[string]string) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(snapshots))

	for _, snapshot := range snapshots {
		var sizeGB float64 = 0
		if sizeStr, ok := snapshot["DiskSizeGB"]; ok && sizeStr != "" {
			numStr := strings.TrimSpace(sizeStr)
			numStr = strings.Split(numStr, " ")[0]

			if val, err := strconv.ParseFloat(numStr, 64); err == nil {
				sizeGB = val
			}
		}

		if sizeGB == 0 {
			sizeGB = 100.0
		}

		region := "us-central1"
		if location, ok := snapshot["Location"]; ok && location != "" {
			region = strings.ToLower(location)
		}

		pricePerGB := GetPrice("gcp", "disk_snapshot", region)
		monthlyCost := sizeGB * pricePerGB

		priceSource := GetPricingSource("gcp", "disk_snapshot")
		priceInfo := GetPricingMetadata("gcp", "disk_snapshot")

		log.Printf("GCP disk pricing for region %s: $%.4f per GB/month (Source: %s, Last verified: %s)",
			region, pricePerGB, priceSource, priceInfo.LastVerified.Format("2006-01-02"))

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

func EstimateGcpResourceCosts(resourceData map[string]interface{}) (map[string]interface{}, error) {
	costData := make(map[string]interface{})

	if diskSnapshotsRaw, ok := resourceData["DiskSnapshots"]; ok {
		if diskSnapshots, ok := convertToSnapshotList(diskSnapshotsRaw); ok {
			costData["DiskSnapshotCosts"] = CalculateGcpDiskSnapshotCosts(diskSnapshots)
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
	log.Println("Starting conversion of GCP data for cost analysis...")

	gcpData, err := gcp.CollectGCPData(ctx)
	if err != nil {
		log.Printf("Failed to collect GCP data: %v", err)
		return nil, fmt.Errorf("failed to collect GCP data: %v", err)
	}

	log.Printf("GCP data collected: %d compute instances, %d GCS buckets, %d Cloud SQL instances, %d Cloud Run services, %d Cloud Functions",
		len(gcpData.ComputeInstances), len(gcpData.GCSBuckets), len(gcpData.CloudSQLInstances),
		len(gcpData.CloudRunServices), len(gcpData.CloudFunctions))

	inventory := make(map[string]interface{})

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
		log.Printf("Added %d GCP Compute instances to cost inventory", len(instances))
	}

	if len(gcpData.GCSBuckets) > 0 {
		buckets := make([]map[string]interface{}, len(gcpData.GCSBuckets))
		for i, bucket := range gcpData.GCSBuckets {
			buckets[i] = map[string]interface{}{
				"Name":         bucket.Name,
				"Location":     bucket.Location,
				"StorageClass": bucket.StorageClass,
				"Project":      bucket.Project,
				"SizeGB":       100.0,
			}
		}
		inventory["GCSBuckets"] = buckets
		log.Printf("Added %d GCP Storage buckets to cost inventory", len(buckets))
	}

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
				"DiskSizeGB":      100.0,
			}
		}
		inventory["CloudSQLInstances"] = instances
		log.Printf("Added %d GCP Cloud SQL instances to cost inventory", len(instances))
	}

	snapshotData, err := gcp.CollectSnapshotData(ctx)
	if err != nil {
		log.Printf("Warning: Failed to collect GCP snapshot data: %v", err)
	} else {
		for k, v := range snapshotData {
			inventory[k] = v
			log.Printf("Added GCP snapshot data type %s to cost inventory", k)
		}
	}

	return inventory, nil
}
