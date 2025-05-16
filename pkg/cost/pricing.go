package cost

import (
	"time"
)

// RegionalPrice stores price information per region
type RegionalPrice map[string]float64

var (
	// AWS EBS snapshot pricing ($/GB-month) as of 2023
	AwsEbsSnapshotPricing = RegionalPrice{
		"us-east-1":      0.05,
		"us-east-2":      0.05,
		"us-west-1":      0.05,
		"us-west-2":      0.05,
		"eu-west-1":      0.05,
		"eu-central-1":   0.05,
		"ap-northeast-1": 0.05,
		"ap-southeast-1": 0.05,
		"ap-southeast-2": 0.05,
		"sa-east-1":      0.065,
		"ap-south-1":     0.055,
		"ca-central-1":   0.055,
		// Default fallback
		"default": 0.05,
	}

	// AWS RDS snapshot pricing (slightly higher than EBS)
	AwsRdsSnapshotPricing = RegionalPrice{
		"us-east-1":      0.095,
		"us-east-2":      0.095,
		"us-west-1":      0.095,
		"us-west-2":      0.095,
		"eu-west-1":      0.095,
		"eu-central-1":   0.105,
		"ap-northeast-1": 0.10,
		"ap-southeast-1": 0.10,
		"ap-southeast-2": 0.105,
		"sa-east-1":      0.115,
		"ap-south-1":     0.105,
		"ca-central-1":   0.105,
		// Default fallback
		"default": 0.095,
	}

	// Azure managed disk snapshot pricing ($/GB-month) as of 2023
	AzureDiskSnapshotPricing = RegionalPrice{
		"eastus":             0.05,
		"eastus2":            0.05,
		"westus":             0.05,
		"westus2":            0.05,
		"centralus":          0.05,
		"northeurope":        0.05,
		"westeurope":         0.05,
		"southeastasia":      0.05,
		"eastasia":           0.05,
		"australiaeast":      0.07,
		"australiasoutheast": 0.07,
		// Default fallback
		"default": 0.05,
	}

	// GCP Persistent Disk snapshot pricing ($/GB-month) as of 2023
	GcpDiskSnapshotPricing = RegionalPrice{
		"us-central1":          0.026,
		"us-east1":             0.026,
		"us-west1":             0.026,
		"europe-west1":         0.026,
		"europe-west2":         0.031,
		"europe-west3":         0.031,
		"asia-east1":           0.031,
		"asia-southeast1":      0.031,
		"australia-southeast1": 0.036,
		// Default fallback
		"default": 0.03,
	}

	// Last update time
	lastPricingUpdate = time.Now()
)

// GetPrice returns the price for a given region or the default price if region not found
func GetPrice(pricing RegionalPrice, region string) float64 {
	if price, ok := pricing[region]; ok {
		return price
	}
	return pricing["default"]
}
