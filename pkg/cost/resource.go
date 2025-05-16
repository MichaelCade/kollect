package cost

import (
	"log"
	"strings"
)

// ResourceCostEstimator calculates costs for cloud resources
type ResourceCostEstimator struct {
	// You can add configuration fields here if needed
}

// NewResourceCostEstimator creates a new resource cost estimator
func NewResourceCostEstimator() *ResourceCostEstimator {
	return &ResourceCostEstimator{}
}

// EstimateAwsResourcesCost estimates costs for various AWS resources
func (r *ResourceCostEstimator) EstimateAwsResourcesCost(resourceData map[string]interface{}) (map[string]interface{}, error) {
	// Start with the snapshot costs
	costData, err := EstimateAwsResourceCosts(resourceData)
	if err != nil {
		return nil, err
	}

	// Add cost calculations for EC2 instances, EBS volumes, S3 buckets, etc.
	var totalResourceCost float64

	// EC2 Instances - calculate approximate costs
	if ec2InstancesRaw, ok := resourceData["EC2Instances"]; ok {
		if ec2Instances, ok := ec2InstancesRaw.([]map[string]interface{}); ok {
			var ec2Costs []map[string]interface{}

			log.Printf("Calculating costs for %d EC2 instances", len(ec2Instances))

			for _, instance := range ec2Instances {
				// Default cost for small instance
				hourlyCost := 0.05

				// Adjust based on instance type
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

	// S3 Buckets - calculate approximate costs
	if s3BucketsRaw, ok := resourceData["S3Buckets"]; ok {
		if s3Buckets, ok := s3BucketsRaw.([]map[string]interface{}); ok {
			var s3Costs []map[string]interface{}

			log.Printf("Calculating costs for %d S3 buckets", len(s3Buckets))

			for _, bucket := range s3Buckets {
				sizeGB := 100.0 // Default size
				if size, ok := bucket["SizeGB"].(float64); ok && size > 0 {
					sizeGB = size
				}

				storageClass := "STANDARD"
				if sc, ok := bucket["StorageClass"].(string); ok && sc != "" {
					storageClass = sc
				}

				// Price per GB per month
				pricePerGB := 0.023 // Standard storage

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

	// Update the summary to include all resource costs
	if summary, ok := costData["Summary"].(map[string]interface{}); ok {
		if totalCost, ok := summary["TotalMonthlyCost"].(float64); ok {
			summary["TotalMonthlyCost"] = totalCost + totalResourceCost
			summary["TotalComputeCost"] = totalResourceCost
			log.Printf("Updated summary with total resource cost of $%.2f", totalResourceCost)
		}
	}

	return costData, nil
}

// EstimateAzureResourcesCost estimates costs for various Azure resources
func (r *ResourceCostEstimator) EstimateAzureResourcesCost(resourceData map[string]interface{}) (map[string]interface{}, error) {
	// Start with the snapshot costs
	costData, err := EstimateAzureResourceCosts(resourceData)
	if err != nil {
		return nil, err
	}

	// Add cost calculations for VMs, storage accounts, etc.
	var totalResourceCost float64

	// Virtual Machines - calculate approximate costs
	if vmsRaw, ok := resourceData["VirtualMachines"]; ok {
		if vms, ok := vmsRaw.([]map[string]interface{}); ok {
			var vmCosts []map[string]interface{}

			for _, vm := range vms {
				// Default cost for small VM
				hourlyCost := 0.05

				// Adjust based on VM size
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

	// Update the summary to include all resource costs
	if summary, ok := costData["Summary"].(map[string]interface{}); ok {
		if totalCost, ok := summary["TotalMonthlyCost"].(float64); ok {
			summary["TotalMonthlyCost"] = totalCost + totalResourceCost
			summary["TotalComputeCost"] = totalResourceCost
		}
	}

	return costData, nil
}

// EstimateGcpResourcesCost estimates costs for various GCP resources
func (r *ResourceCostEstimator) EstimateGcpResourcesCost(resourceData map[string]interface{}) (map[string]interface{}, error) {
	// Start with the snapshot costs
	costData, err := EstimateGcpResourceCosts(resourceData)
	if err != nil {
		return nil, err
	}

	// Add cost calculations for compute instances, Cloud SQL, etc.
	var totalResourceCost float64

	// Compute Instances - calculate approximate costs
	if instancesRaw, ok := resourceData["ComputeInstances"]; ok {
		if instances, ok := instancesRaw.([]map[string]interface{}); ok {
			var instanceCosts []map[string]interface{}

			for _, instance := range instances {
				// Default cost for small instance
				hourlyCost := 0.05

				// Adjust based on machine type
				if machineType, ok := instance["MachineType"].(string); ok {
					switch {
					case strings.Contains(machineType, "n1-standard"):
						hourlyCost = 0.05
					case strings.Contains(machineType, "n2-standard"):
						hourlyCost = 0.07
					case strings.Contains(machineType, "e2-"):
						hourlyCost = 0.03
					}
				}

				monthlyCost := hourlyCost * 24 * 30

				cost := map[string]interface{}{
					"Name":        instance["Name"],
					"Zone":        instance["Zone"],
					"MachineType": instance["MachineType"],
					"HourlyCost":  hourlyCost,
					"MonthlyCost": monthlyCost,
				}

				instanceCosts = append(instanceCosts, cost)
				totalResourceCost += monthlyCost
			}

			costData["ComputeCosts"] = instanceCosts
		}
	}

	// Update the summary to include all resource costs
	if summary, ok := costData["Summary"].(map[string]interface{}); ok {
		if totalCost, ok := summary["TotalMonthlyCost"].(float64); ok {
			summary["TotalMonthlyCost"] = totalCost + totalResourceCost
			summary["TotalComputeCost"] = totalResourceCost
		}
	}

	return costData, nil
}
