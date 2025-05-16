package cost

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/michaelcade/kollect/pkg/aws"
)

// CalculateEbsSnapshotCosts calculates costs for EBS snapshots
func CalculateEbsSnapshotCosts(snapshots []map[string]string) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(snapshots))

	for _, snapshot := range snapshots {
		// Parse size
		var sizeGB float64 = 0
		if sizeStr, ok := snapshot["VolumeSize"]; ok && sizeStr != "" {
			// Clean up the string to extract just the number
			numStr := strings.TrimSpace(sizeStr)
			numStr = strings.TrimSuffix(numStr, " GiB") // Remove unit if present
			numStr = strings.Split(numStr, " ")[0]      // Take the first part before any space

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
		if regionVal, ok := snapshot["Region"]; ok && regionVal != "" {
			region = strings.ToLower(regionVal)
		}

		// Calculate monthly cost
		pricePerGB := GetPrice("aws", "ebs_snapshot", region)
		monthlyCost := sizeGB * pricePerGB

		// Get pricing source and metadata
		priceSource := GetPricingSource("aws", "ebs_snapshot")
		priceInfo := GetPricingMetadata("aws", "ebs_snapshot")

		// Log price information for debugging
		log.Printf("AWS EBS snapshot pricing for region %s: $%.4f per GB/month (Source: %s, Last verified: %s)",
			region, pricePerGB, priceSource, priceInfo.LastVerified.Format("2006-01-02"))

		// Create result with cost info
		result := map[string]interface{}{
			"SnapshotId":      snapshot["SnapshotId"],
			"Description":     snapshot["Description"],
			"VolumeId":        snapshot["VolumeId"],
			"SizeGB":          sizeGB,
			"Region":          region,
			"State":           snapshot["State"],
			"StartTime":       snapshot["StartTime"],
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

// CalculateRdsSnapshotCosts calculates costs for RDS snapshots
func CalculateRdsSnapshotCosts(snapshots []map[string]string) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(snapshots))

	for _, snapshot := range snapshots {
		// Parse size
		var sizeGB float64 = 0
		if sizeStr, ok := snapshot["AllocatedStorage"]; ok && sizeStr != "" {
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
		region := "us-east-1" // Default region
		if regionVal, ok := snapshot["Region"]; ok && regionVal != "" {
			region = strings.ToLower(regionVal)
		}

		// Calculate monthly cost
		pricePerGB := GetPrice("aws", "rds_snapshot", region)
		monthlyCost := sizeGB * pricePerGB

		// Get pricing source and metadata
		priceSource := GetPricingSource("aws", "rds_snapshot")
		priceInfo := GetPricingMetadata("aws", "rds_snapshot")

		// Log price information for debugging
		log.Printf("AWS RDS snapshot pricing for region %s: $%.4f per GB/month (Source: %s, Last verified: %s)",
			region, pricePerGB, priceSource, priceInfo.LastVerified.Format("2006-01-02"))

		// Create result with cost info
		result := map[string]interface{}{
			"DBSnapshotIdentifier": snapshot["DBSnapshotIdentifier"],
			"DBInstanceIdentifier": snapshot["DBInstanceIdentifier"],
			"SnapshotType":         snapshot["SnapshotType"],
			"Engine":               snapshot["Engine"],
			"SizeGB":               sizeGB,
			"Region":               region,
			"Status":               snapshot["Status"],
			"SnapshotCreateTime":   snapshot["SnapshotCreateTime"],
			"PricePerGBMonth":      pricePerGB,
			"PriceSource":          priceSource,
			"LastVerified":         priceInfo.LastVerified.Format("2006-01-02"),
			"MonthlyCost":          monthlyCost,
			"MonthlyCostUSD":       fmt.Sprintf("$%.2f", monthlyCost),
		}

		results = append(results, result)
	}

	return results
}

// EstimateAwsResourceCosts calculates costs for AWS resources
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

	// Include pricing source information in the summary
	priceSource := GetPricingSource("aws", "ebs_snapshot")
	lastVerified := GetPricingMetadata("aws", "ebs_snapshot").LastVerified.Format("2006-01-02")

	costData["Summary"] = map[string]interface{}{
		"TotalSnapshotStorage": totalSnapshotStorage,
		"TotalMonthlyCost":     totalMonthlyCost,
		"Currency":             "USD",
		"PriceSource":          priceSource,
		"LastVerified":         lastVerified,
	}

	return costData, nil
}

// ConvertAwsDataForCostAnalysis converts the structured AWSData into a generic map for cost analysis
func ConvertAwsDataForCostAnalysis(ctx context.Context) (map[string]interface{}, error) {
	// Use the existing AWS inventory collection function
	awsData, err := aws.CollectAWSData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect AWS data: %v", err)
	}

	// Create a generic map to hold all resource types
	inventory := make(map[string]interface{})

	// Convert EC2 instances to generic map entries
	if len(awsData.EC2Instances) > 0 {
		instances := make([]map[string]interface{}, len(awsData.EC2Instances))
		for i, instance := range awsData.EC2Instances {
			instances[i] = map[string]interface{}{
				"InstanceId":   instance.InstanceID,
				"Name":         instance.Name,
				"InstanceType": instance.Type,
				"State":        instance.State,
				"Region":       instance.Region,
			}
		}
		inventory["EC2Instances"] = instances
		log.Printf("Added %d EC2 instances to cost inventory", len(instances))
	}

	// Convert S3 buckets to generic map entries
	if len(awsData.S3Buckets) > 0 {
		buckets := make([]map[string]interface{}, len(awsData.S3Buckets))
		for i, bucket := range awsData.S3Buckets {
			buckets[i] = map[string]interface{}{
				"Name":         bucket.Name,
				"Region":       bucket.Region,
				"Immutable":    bucket.Immutable,
				"SizeGB":       100.0, // Default size estimate
				"StorageClass": "STANDARD",
			}
		}
		inventory["S3Buckets"] = buckets
		log.Printf("Added %d S3 buckets to cost inventory", len(buckets))
	}

	// Convert RDS instances to generic map entries
	if len(awsData.RDSInstances) > 0 {
		instances := make([]map[string]interface{}, len(awsData.RDSInstances))
		for i, instance := range awsData.RDSInstances {
			instances[i] = map[string]interface{}{
				"DBInstanceIdentifier": instance.InstanceID,
				"Engine":               instance.Engine,
				"Status":               instance.Status,
				"Region":               instance.Region,
				"DBInstanceClass":      "db.t3.medium", // Default value
				"AllocatedStorage":     20.0,           // Default value
			}
		}
		inventory["RDSInstances"] = instances
		log.Printf("Added %d RDS instances to cost inventory", len(instances))
	}

	// Convert DynamoDB tables to generic map entries
	if len(awsData.DynamoDBTables) > 0 {
		tables := make([]map[string]interface{}, len(awsData.DynamoDBTables))
		for i, table := range awsData.DynamoDBTables {
			tables[i] = map[string]interface{}{
				"TableName": table.TableName,
				"Status":    table.Status,
				"Region":    table.Region,
				"SizeGB":    1.0, // Default value
			}
		}
		inventory["DynamoDBTables"] = tables
		log.Printf("Added %d DynamoDB tables to cost inventory", len(tables))
	}

	// Convert VPCs to generic map entries
	if len(awsData.VPCs) > 0 {
		vpcs := make([]map[string]interface{}, len(awsData.VPCs))
		for i, vpc := range awsData.VPCs {
			vpcs[i] = map[string]interface{}{
				"VPCID":  vpc.VPCID,
				"State":  vpc.State,
				"Region": vpc.Region,
			}
		}
		inventory["VPCs"] = vpcs
		log.Printf("Added %d VPCs to cost inventory", len(vpcs))
	}

	// Include snapshot data
	snapshotData, err := aws.CollectSnapshotData(ctx)
	if err != nil {
		log.Printf("Warning: Failed to collect snapshot data: %v", err)
	} else {
		for k, v := range snapshotData {
			inventory[k] = v
			log.Printf("Added snapshot data type %s to cost inventory", k)
		}
	}

	return inventory, nil
}
