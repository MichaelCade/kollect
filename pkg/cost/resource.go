package cost

import (
	"log"
	"strings"
)

type ResourceCostEstimator struct {
}

func NewResourceCostEstimator() *ResourceCostEstimator {
	return &ResourceCostEstimator{}
}

func (r *ResourceCostEstimator) EstimateAwsResourcesCost(resourceData map[string]interface{}) (map[string]interface{}, error) {
	costData, err := EstimateAwsResourceCosts(resourceData)
	if err != nil {
		return nil, err
	}

	var totalResourceCost float64

	if ec2InstancesRaw, ok := resourceData["EC2Instances"]; ok {
		if ec2Instances, ok := ec2InstancesRaw.([]map[string]interface{}); ok {
			var ec2Costs []map[string]interface{}

			log.Printf("Calculating costs for %d EC2 instances", len(ec2Instances))

			for _, instance := range ec2Instances {
				hourlyCost := 0.05

				if instanceType, ok := instance["InstanceType"].(string); ok {
					switch {
					case strings.HasPrefix(instanceType, "t2."):
						hourlyCost = 0.02
					case strings.HasPrefix(instanceType, "t3."):
						hourlyCost = 0.03
					case strings.HasPrefix(instanceType, "m5."):
						hourlyCost = 0.1
					case strings.HasPrefix(instanceType, "c5."):
						hourlyCost = 0.08
					}
				}

				monthlyCost := hourlyCost * 24 * 30

				cost := map[string]interface{}{
					"InstanceId":   instance["InstanceId"],
					"InstanceType": instance["InstanceType"],
					"Region":       instance["Region"],
					"HourlyCost":   hourlyCost,
					"MonthlyCost":  monthlyCost,
				}

				ec2Costs = append(ec2Costs, cost)
				totalResourceCost += monthlyCost
			}

			costData["EC2Costs"] = ec2Costs
			log.Printf("Added %d EC2 instance costs", len(ec2Costs))
		}
	}

	if s3BucketsRaw, ok := resourceData["S3Buckets"]; ok {
		if s3Buckets, ok := s3BucketsRaw.([]map[string]interface{}); ok {
			var s3Costs []map[string]interface{}

			log.Printf("Calculating costs for %d S3 buckets", len(s3Buckets))

			for _, bucket := range s3Buckets {
				sizeGB := 100.0
				if size, ok := bucket["SizeGB"].(float64); ok && size > 0 {
					sizeGB = size
				}

				storageClass := "STANDARD"
				if sc, ok := bucket["StorageClass"].(string); ok && sc != "" {
					storageClass = sc
				}

				pricePerGB := 0.023

				switch storageClass {
				case "STANDARD_IA":
					pricePerGB = 0.0125
				case "ONEZONE_IA":
					pricePerGB = 0.01
				case "GLACIER":
					pricePerGB = 0.004
				case "DEEP_ARCHIVE":
					pricePerGB = 0.00099
				}

				monthlyCost := sizeGB * pricePerGB

				cost := map[string]interface{}{
					"Name":         bucket["Name"],
					"Region":       bucket["Region"],
					"SizeGB":       sizeGB,
					"StorageClass": storageClass,
					"PricePerGB":   pricePerGB,
					"MonthlyCost":  monthlyCost,
				}

				s3Costs = append(s3Costs, cost)
				totalResourceCost += monthlyCost
			}

			costData["S3Costs"] = s3Costs
			log.Printf("Added %d S3 bucket costs", len(s3Costs))
		}
	}

	if summary, ok := costData["Summary"].(map[string]interface{}); ok {
		if totalCost, ok := summary["TotalMonthlyCost"].(float64); ok {
			summary["TotalMonthlyCost"] = totalCost + totalResourceCost
			summary["TotalComputeCost"] = totalResourceCost
			log.Printf("Updated summary with total AWS resource cost of $%.2f", totalResourceCost)
		}
	}

	return costData, nil
}

func (r *ResourceCostEstimator) EstimateGcpResourcesCost(resourceData map[string]interface{}) (map[string]interface{}, error) {
	costData, err := EstimateGcpResourceCosts(resourceData)
	if err != nil {
		return nil, err
	}

	var totalResourceCost float64

	if computeInstancesRaw, ok := resourceData["ComputeInstances"]; ok {
		if computeInstances, ok := computeInstancesRaw.([]map[string]interface{}); ok {
			var computeCosts []map[string]interface{}

			log.Printf("Calculating costs for %d GCP Compute instances", len(computeInstances))

			for _, instance := range computeInstances {
				hourlyCost := 0.05

				if machineType, ok := instance["MachineType"].(string); ok {
					switch {
					case strings.Contains(machineType, "f1-micro"):
						hourlyCost = 0.01
					case strings.Contains(machineType, "g1-small"):
						hourlyCost = 0.02
					case strings.Contains(machineType, "e2-"):
						hourlyCost = 0.03
					case strings.Contains(machineType, "n1-standard"):
						hourlyCost = 0.05
					case strings.Contains(machineType, "n2-standard"):
						hourlyCost = 0.06
					case strings.Contains(machineType, "c2-"):
						hourlyCost = 0.09
					}
				}

				zone := "us-central1-a"
				if z, ok := instance["Zone"].(string); ok && z != "" {
					zone = z
				}

				if strings.HasPrefix(zone, "australia-") ||
					strings.HasPrefix(zone, "europe-") ||
					strings.HasPrefix(zone, "asia-") {
					hourlyCost *= 1.2
				}

				monthlyCost := hourlyCost * 24 * 30

				cost := map[string]interface{}{
					"Name":        instance["Name"],
					"MachineType": instance["MachineType"],
					"Zone":        zone,
					"HourlyCost":  hourlyCost,
					"MonthlyCost": monthlyCost,
				}

				computeCosts = append(computeCosts, cost)
				totalResourceCost += monthlyCost
			}

			costData["ComputeCosts"] = computeCosts
			log.Printf("Added %d GCP Compute instance costs", len(computeCosts))
		}
	}

	if gcsBucketsRaw, ok := resourceData["GCSBuckets"]; ok {
		if gcsBuckets, ok := gcsBucketsRaw.([]map[string]interface{}); ok {
			var gcsCosts []map[string]interface{}

			log.Printf("Calculating costs for %d GCS buckets", len(gcsBuckets))

			for _, bucket := range gcsBuckets {
				sizeGB := 100.0
				if size, ok := bucket["SizeGB"].(float64); ok && size > 0 {
					sizeGB = size
				}

				storageClass := "STANDARD"
				if sc, ok := bucket["StorageClass"].(string); ok && sc != "" {
					storageClass = sc
				}

				pricePerGB := 0.02

				switch storageClass {
				case "NEARLINE":
					pricePerGB = 0.01
				case "COLDLINE":
					pricePerGB = 0.007
				case "ARCHIVE":
					pricePerGB = 0.004
				}

				monthlyCost := sizeGB * pricePerGB

				cost := map[string]interface{}{
					"Name":         bucket["Name"],
					"Location":     bucket["Location"],
					"SizeGB":       sizeGB,
					"StorageClass": storageClass,
					"PricePerGB":   pricePerGB,
					"MonthlyCost":  monthlyCost,
				}

				gcsCosts = append(gcsCosts, cost)
				totalResourceCost += monthlyCost
			}

			costData["GCSCosts"] = gcsCosts
			log.Printf("Added %d GCS bucket costs", len(gcsCosts))
		}
	}

	if cloudSQLInstancesRaw, ok := resourceData["CloudSQLInstances"]; ok {
		if cloudSQLInstances, ok := cloudSQLInstancesRaw.([]map[string]interface{}); ok {
			var sqlCosts []map[string]interface{}

			log.Printf("Calculating costs for %d Cloud SQL instances", len(cloudSQLInstances))

			for _, instance := range cloudSQLInstances {
				hourlyCost := 0.1

				if tier, ok := instance["Tier"].(string); ok {
					switch {
					case strings.Contains(tier, "db-f1-micro"):
						hourlyCost = 0.025
					case strings.Contains(tier, "db-g1-small"):
						hourlyCost = 0.05
					case strings.Contains(tier, "standard"):
						hourlyCost = 0.1
					case strings.Contains(tier, "highmem"):
						hourlyCost = 0.15
					case strings.Contains(tier, "highcpu"):
						hourlyCost = 0.12
					}
				}

				diskSizeGB := 100.0
				if size, ok := instance["DiskSizeGB"].(float64); ok && size > 0 {
					diskSizeGB = size
				}

				diskCostPerMonth := diskSizeGB * 0.17

				monthlyCost := (hourlyCost * 24 * 30) + diskCostPerMonth

				cost := map[string]interface{}{
					"Name":            instance["Name"],
					"DatabaseVersion": instance["DatabaseVersion"],
					"Region":          instance["Region"],
					"Tier":            instance["Tier"],
					"DiskSizeGB":      diskSizeGB,
					"HourlyCost":      hourlyCost,
					"MonthlyCost":     monthlyCost,
				}

				sqlCosts = append(sqlCosts, cost)
				totalResourceCost += monthlyCost
			}

			costData["CloudSQLCosts"] = sqlCosts
			log.Printf("Added %d Cloud SQL instance costs", len(sqlCosts))
		}
	}

	if summary, ok := costData["Summary"].(map[string]interface{}); ok {
		if totalCost, ok := summary["TotalMonthlyCost"].(float64); ok {
			summary["TotalMonthlyCost"] = totalCost + totalResourceCost
			summary["TotalComputeCost"] = totalResourceCost
			log.Printf("Updated summary with total GCP resource cost of $%.2f", totalResourceCost)
		}
	}

	return costData, nil
}

func (r *ResourceCostEstimator) EstimateAzureResourcesCost(resourceData map[string]interface{}) (map[string]interface{}, error) {
	costData, err := EstimateAzureResourceCosts(resourceData)
	if err != nil {
		return nil, err
	}

	var totalResourceCost float64

	if vmsRaw, ok := resourceData["VirtualMachines"]; ok {
		if vms, ok := vmsRaw.([]map[string]interface{}); ok {
			var vmCosts []map[string]interface{}

			for _, vm := range vms {
				hourlyCost := 0.05

				if vmSize, ok := vm["VMSize"].(string); ok {
					switch {
					case strings.Contains(vmSize, "Standard_B"):
						hourlyCost = 0.03
					case strings.Contains(vmSize, "Standard_D"):
						hourlyCost = 0.1
					case strings.Contains(vmSize, "Standard_E"):
						hourlyCost = 0.15
					}
				}

				monthlyCost := hourlyCost * 24 * 30

				cost := map[string]interface{}{
					"Name":          vm["Name"],
					"ResourceGroup": vm["ResourceGroup"],
					"Location":      vm["Location"],
					"VMSize":        vm["VMSize"],
					"HourlyCost":    hourlyCost,
					"MonthlyCost":   monthlyCost,
				}

				vmCosts = append(vmCosts, cost)
				totalResourceCost += monthlyCost
			}

			costData["VMCosts"] = vmCosts
		}
	}

	if summary, ok := costData["Summary"].(map[string]interface{}); ok {
		if totalCost, ok := summary["TotalMonthlyCost"].(float64); ok {
			summary["TotalMonthlyCost"] = totalCost + totalResourceCost
			summary["TotalComputeCost"] = totalResourceCost
		}
	}

	return costData, nil
}
