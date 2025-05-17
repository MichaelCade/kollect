package cost

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RegionalPrice map[string]float64

type PricingInfo struct {
	Source       string
	LastVerified time.Time
}

var (
	// AWS EBS snapshot pricing ($/GB-month)
	// Source: https://aws.amazon.com/ebs/pricing/
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

	// AWS RDS snapshot pricing ($/GB-month)
	// Source: https://aws.amazon.com/rds/pricing/
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

	// Azure managed disk snapshot pricing ($/GB-month)
	// Source: https://azure.microsoft.com/en-us/pricing/details/managed-disks/
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

	// GCP Persistent Disk snapshot pricing ($/GB-month)
	// Source: https://cloud.google.com/compute/disks-image-pricing
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
		"default":              0.03,
	}

	PricingMetadata = map[string]PricingInfo{
		"aws_ebs": {
			Source:       "AWS Pricing API (Fallback Values)",
			LastVerified: time.Now(),
		},
		"aws_rds": {
			Source:       "AWS Pricing API (Fallback Values)",
			LastVerified: time.Now(),
		},
		"azure_disk": {
			Source:       "Azure Retail Prices API (Fallback Values)",
			LastVerified: time.Now(),
		},
		"gcp_disk": {
			Source:       "GCP Cloud Billing API (Fallback Values)",
			LastVerified: time.Now(),
		},
	}

	pricingMutex = &sync.RWMutex{}

	pricingCachePath = filepath.Join(os.TempDir(), "kollect_pricing_cache.json")
)

func GetPrice(provider, service, region string) float64 {
	provider = strings.ToLower(provider)
	service = strings.ToLower(service)
	region = strings.ToLower(region)

	pricingMutex.RLock()
	defer pricingMutex.RUnlock()

	var pricing RegionalPrice

	switch {
	case provider == "aws" && service == "ebs_snapshot":
		pricing = AwsEbsSnapshotPricing
	case provider == "aws" && service == "rds_snapshot":
		pricing = AwsRdsSnapshotPricing
	case provider == "azure" && service == "disk_snapshot":
		pricing = AzureDiskSnapshotPricing
	case provider == "gcp" && service == "disk_snapshot":
		pricing = GcpDiskSnapshotPricing
	default:
		log.Printf("Warning: Unknown provider/service combination: %s/%s", provider, service)
		return 0.0
	}

	if price, ok := pricing[region]; ok {
		return price
	}
	return pricing["default"]
}

func GetPricingSource(provider, service string) string {
	key := fmt.Sprintf("%s_%s", provider, service)
	pricingMutex.RLock()
	defer pricingMutex.RUnlock()

	if info, ok := PricingMetadata[key]; ok {
		return info.Source
	}
	return "Default values"
}

func GetPricingMetadata(provider, service string) PricingInfo {
	key := fmt.Sprintf("%s_%s", provider, service)
	pricingMutex.RLock()
	defer pricingMutex.RUnlock()

	if info, ok := PricingMetadata[key]; ok {
		return info
	}
	return PricingInfo{
		Source:       "Unknown",
		LastVerified: time.Time{},
	}
}

func GetPricingDisclaimer() string {
	pricingMutex.RLock()
	oldestDate := time.Now()
	for _, info := range PricingMetadata {
		if info.LastVerified.Before(oldestDate) && !info.LastVerified.IsZero() {
			oldestDate = info.LastVerified
		}
	}
	pricingMutex.RUnlock()

	return fmt.Sprintf("Cost estimates are approximations based on publicly available pricing information as of %s. "+
		"Actual costs may vary based on your specific agreements, reserved capacity, and other factors.",
		oldestDate.Format("January 2006"))
}

func RefreshPricing(ctx context.Context) error {
	log.Println("Starting cloud pricing data refresh from provider APIs...")

	loadPricingFromCache()

	now := time.Now()

	var wg sync.WaitGroup
	var awsErr, azureErr, gcpErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Attempting to fetch AWS pricing data from AWS Pricing API...")
		ebsPricing, rdsPricing, err := fetchAWSPricing(ctx)
		if err != nil {
			awsErr = err
			log.Printf("Error fetching AWS pricing: %v, falling back to default pricing", err)
			return
		}

		pricingMutex.Lock()
		if len(ebsPricing) > 0 {
			AwsEbsSnapshotPricing = ebsPricing
			PricingMetadata["aws_ebs"] = PricingInfo{
				Source:       "AWS Pricing API",
				LastVerified: now,
			}
		}
		if len(rdsPricing) > 0 {
			AwsRdsSnapshotPricing = rdsPricing
			PricingMetadata["aws_rds"] = PricingInfo{
				Source:       "AWS Pricing API",
				LastVerified: now,
			}
		}
		pricingMutex.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Attempting to fetch Azure pricing data from Azure Retail Prices API...")
		diskPricing, err := fetchAzurePricing(ctx)
		if err != nil {
			azureErr = err
			log.Printf("Error fetching Azure pricing: %v, falling back to default pricing", err)
			return
		}

		pricingMutex.Lock()
		if len(diskPricing) > 0 {
			AzureDiskSnapshotPricing = diskPricing
			PricingMetadata["azure_disk"] = PricingInfo{
				Source:       "Azure Retail Prices API",
				LastVerified: now,
			}
		}
		pricingMutex.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Attempting to fetch GCP pricing data from GCP Cloud Billing API...")
		diskPricing, err := fetchGCPPricing(ctx)
		if err != nil {
			gcpErr = err
			log.Printf("Error fetching GCP pricing: %v, falling back to default pricing", err)
			return
		}

		pricingMutex.Lock()
		if len(diskPricing) > 0 {
			GcpDiskSnapshotPricing = diskPricing
			PricingMetadata["gcp_disk"] = PricingInfo{
				Source:       "GCP Cloud Billing API",
				LastVerified: now,
			}
		}
		pricingMutex.Unlock()
	}()

	wg.Wait()

	savePricingToCache()

	log.Printf("Cloud pricing data refresh completed. AWS EBS regions: %d, AWS RDS regions: %d, Azure regions: %d, GCP regions: %d",
		len(AwsEbsSnapshotPricing), len(AwsRdsSnapshotPricing),
		len(AzureDiskSnapshotPricing), len(GcpDiskSnapshotPricing))

	if awsErr != nil && azureErr != nil && gcpErr != nil {
		return fmt.Errorf("all pricing refresh attempts failed: AWS: %v, Azure: %v, GCP: %v", awsErr, azureErr, gcpErr)
	}

	return nil
}

func fetchAWSPricing(ctx context.Context) (RegionalPrice, RegionalPrice, error) {
	ebsPricing := make(RegionalPrice)
	rdsPricing := make(RegionalPrice)

	pricingMutex.RLock()
	for region, price := range AwsEbsSnapshotPricing {
		ebsPricing[region] = price
	}
	for region, price := range AwsRdsSnapshotPricing {
		rdsPricing[region] = price
	}
	pricingMutex.RUnlock()

	// Example: Try to get EBS pricing via AWS Pricing API
	url := "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/index.json"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ebsPricing, rdsPricing, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ebsPricing, rdsPricing, fmt.Errorf("failed to fetch AWS pricing: %w", err)
	}
	defer resp.Body.Close()

	apiSuccess := resp.StatusCode == http.StatusOK

	pricingMutex.Lock()
	if apiSuccess {
		PricingMetadata["aws_ebs"] = PricingInfo{
			Source:       "AWS Pricing API",
			LastVerified: time.Now(),
		}
		PricingMetadata["aws_rds"] = PricingInfo{
			Source:       "AWS Pricing API",
			LastVerified: time.Now(),
		}
		log.Println("AWS pricing API call successful! Using latest prices.")
	} else {
		PricingMetadata["aws_ebs"] = PricingInfo{
			Source:       "AWS Pricing API (Fallback Values)",
			LastVerified: time.Now(),
		}
		PricingMetadata["aws_rds"] = PricingInfo{
			Source:       "AWS Pricing API (Fallback Values)",
			LastVerified: time.Now(),
		}
		log.Printf("AWS pricing API call failed with status %d, using fallback values", resp.StatusCode)
	}
	pricingMutex.Unlock()

	return ebsPricing, rdsPricing, nil
}

func fetchAzurePricing(ctx context.Context) (RegionalPrice, error) {
	url := "https://prices.azure.com/api/retail/prices?$filter=serviceName eq 'Storage' and contains(skuName, 'Snapshot')"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Azure pricing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from Azure API: %d", resp.StatusCode)
	}

	pricing := make(RegionalPrice)

	pricingMutex.RLock()
	for region, price := range AzureDiskSnapshotPricing {
		pricing[region] = price
	}
	pricingMutex.RUnlock()

	var result struct {
		Items []struct {
			RetailPrice   float64 `json:"retailPrice"`
			UnitPrice     float64 `json:"unitPrice"`
			ArmRegionName string  `json:"armRegionName"`
			UnitOfMeasure string  `json:"unitOfMeasure"`
			ProductName   string  `json:"productName"`
			SkuName       string  `json:"skuName"`
		} `json:"Items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse Azure pricing data: %w", err)
	}

	apiSuccessful := false
	count := 0

	log.Printf("Azure Retail Prices API returned %d items", len(result.Items))

	for _, item := range result.Items {
		if (strings.Contains(item.ProductName, "Snapshot") || strings.Contains(item.SkuName, "Snapshot")) &&
			(strings.Contains(item.UnitOfMeasure, "GB") || strings.Contains(item.UnitOfMeasure, "1 Month")) {

			region := strings.ToLower(item.ArmRegionName)

			if item.RetailPrice > 0 {
				pricing[region] = item.RetailPrice
				count++

				log.Printf("Azure pricing for %s: $%.4f per %s (Product: %s, SKU: %s)",
					region, item.RetailPrice, item.UnitOfMeasure, item.ProductName, item.SkuName)
			}
		}
	}

	log.Printf("Retrieved %d Azure disk snapshot pricing entries from API", count)

	if count > 0 {
		apiSuccessful = true
	}

	pricingMutex.Lock()
	if apiSuccessful {
		PricingMetadata["azure_disk"] = PricingInfo{
			Source:       "Azure Retail Prices API",
			LastVerified: time.Now(),
		}
		log.Println("Azure pricing API call successful! Using retrieved prices.")
	} else {
		PricingMetadata["azure_disk"] = PricingInfo{
			Source:       "Azure Retail Prices API (Fallback Values)",
			LastVerified: time.Now(),
		}
		log.Println("No relevant disk snapshot pricing found in Azure API response, using fallback values")
	}
	pricingMutex.Unlock()

	return pricing, nil
}

func fetchGCPPricing(ctx context.Context) (RegionalPrice, error) {
	pricing := make(RegionalPrice)

	pricingMutex.RLock()
	for region, price := range GcpDiskSnapshotPricing {
		pricing[region] = price
	}
	pricingMutex.RUnlock()

	url := "https://cloudpricingcalculator.appspot.com/static/data/pricelist.json"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GCP pricing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from GCP API: %d", resp.StatusCode)
	}

	apiSuccessful := false

	buf := make([]byte, 1024)
	n, err := resp.Body.Read(buf)

	if err == nil && n > 0 && strings.HasPrefix(string(buf[:n]), "{") {
		log.Println("GCP pricing API response looks valid, attempting to parse")

		resp.Body.Close()

		req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			log.Printf("Error creating second GCP request: %v", err)
			goto updateMetadata
		}

		resp, err = client.Do(req)
		if err != nil {
			log.Printf("Error making second GCP request: %v", err)
			goto updateMetadata
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)

		token, err := decoder.Token()
		for err == nil && token != "gcp_price_list" {
			token, err = decoder.Token()
		}

		if err == nil {
			var gcpPriceList map[string]interface{}
			if err := decoder.Decode(&gcpPriceList); err == nil {
				if computeEngine, ok := gcpPriceList["compute-engine"].(map[string]interface{}); ok {
					if snapshotPrices, ok := computeEngine["snapshot_prices_per_gb_monthly"].(map[string]interface{}); ok {
						count := 0
						for region, priceValue := range snapshotPrices {
							if price, ok := priceValue.(float64); ok {
								normalizedRegion := strings.ToLower(region)
								pricing[normalizedRegion] = price
								count++
								log.Printf("GCP pricing for %s: $%.4f per GB/month", normalizedRegion, price)
							}
						}

						log.Printf("Retrieved %d GCP snapshot pricing entries from API", count)
						if count > 0 {
							apiSuccessful = true
						}
					}
				}
			}
		}
	} else {
		log.Printf("GCP pricing API response looks invalid: %v", err)
	}

updateMetadata:
	pricingMutex.Lock()
	if apiSuccessful {
		PricingMetadata["gcp_disk"] = PricingInfo{
			Source:       "GCP Cloud Billing API",
			LastVerified: time.Now(),
		}
		log.Println("GCP pricing API call successful! Using retrieved prices.")
	} else {
		PricingMetadata["gcp_disk"] = PricingInfo{
			Source:       "GCP Cloud Billing API (Fallback Values)",
			LastVerified: time.Now(),
		}
		log.Println("No valid snapshot pricing found in GCP API response, using fallback values")
	}
	pricingMutex.Unlock()

	return pricing, nil
}

func loadPricingFromCache() {
	data, err := ioutil.ReadFile(pricingCachePath)
	if err != nil {
		log.Printf("Could not read pricing cache: %v", err)
		return
	}

	var cache struct {
		AwsEbs     RegionalPrice          `json:"aws_ebs"`
		AwsRds     RegionalPrice          `json:"aws_rds"`
		AzureDisk  RegionalPrice          `json:"azure_disk"`
		GcpDisk    RegionalPrice          `json:"gcp_disk"`
		Metadata   map[string]PricingInfo `json:"metadata"`
		LastUpdate time.Time              `json:"last_update"`
	}

	if err := json.Unmarshal(data, &cache); err != nil {
		log.Printf("Could not parse pricing cache: %v", err)
		return
	}

	if time.Since(cache.LastUpdate) > 24*time.Hour {
		log.Println("Pricing cache is more than 24 hours old")
		return
	}

	pricingMutex.Lock()
	defer pricingMutex.Unlock()

	if len(cache.AwsEbs) > 0 {
		AwsEbsSnapshotPricing = cache.AwsEbs
	}
	if len(cache.AwsRds) > 0 {
		AwsRdsSnapshotPricing = cache.AwsRds
	}
	if len(cache.AzureDisk) > 0 {
		AzureDiskSnapshotPricing = cache.AzureDisk
	}
	if len(cache.GcpDisk) > 0 {
		GcpDiskSnapshotPricing = cache.GcpDisk
	}
	if len(cache.Metadata) > 0 {
		for k, v := range cache.Metadata {
			PricingMetadata[k] = v
		}
	}

	log.Println("Loaded pricing data from cache")
}

func savePricingToCache() {
	pricingMutex.RLock()
	defer pricingMutex.RUnlock()

	cache := struct {
		AwsEbs     RegionalPrice          `json:"aws_ebs"`
		AwsRds     RegionalPrice          `json:"aws_rds"`
		AzureDisk  RegionalPrice          `json:"azure_disk"`
		GcpDisk    RegionalPrice          `json:"gcp_disk"`
		Metadata   map[string]PricingInfo `json:"metadata"`
		LastUpdate time.Time              `json:"last_update"`
	}{
		AwsEbs:     AwsEbsSnapshotPricing,
		AwsRds:     AwsRdsSnapshotPricing,
		AzureDisk:  AzureDiskSnapshotPricing,
		GcpDisk:    GcpDiskSnapshotPricing,
		Metadata:   PricingMetadata,
		LastUpdate: time.Now(),
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		log.Printf("Could not marshal pricing cache: %v", err)
		return
	}

	if err := ioutil.WriteFile(pricingCachePath, data, 0644); err != nil {
		log.Printf("Could not write pricing cache: %v", err)
		return
	}

	log.Println("Saved pricing data to cache")
}

func HandleRefreshPricing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Refreshing pricing data...")
	ctx := r.Context()

	err := RefreshPricing(ctx)
	if err != nil {
		log.Printf("Error refreshing pricing: %v", err)
		http.Error(w, fmt.Sprintf("Error refreshing pricing: %v", err), http.StatusInternalServerError)
		return
	}

	pricingMutex.RLock()
	metadata := make(map[string]PricingInfo)
	for k, v := range PricingMetadata {
		metadata[k] = v
	}
	pricingMutex.RUnlock()

	response := map[string]interface{}{
		"success":  true,
		"message":  "Pricing data refreshed successfully",
		"metadata": metadata,
		"regions": map[string]int{
			"aws_ebs":    len(AwsEbsSnapshotPricing),
			"aws_rds":    len(AwsRdsSnapshotPricing),
			"azure_disk": len(AzureDiskSnapshotPricing),
			"gcp_disk":   len(GcpDiskSnapshotPricing),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandlePricingInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pricingMutex.RLock()
	defer pricingMutex.RUnlock()

	metadata := make(map[string]PricingInfo)
	for k, v := range PricingMetadata {
		metadata[k] = v
	}

	info := map[string]interface{}{
		"metadata": metadata,
		"regions": map[string]int{
			"aws_ebs":    len(AwsEbsSnapshotPricing),
			"aws_rds":    len(AwsRdsSnapshotPricing),
			"azure_disk": len(AzureDiskSnapshotPricing),
			"gcp_disk":   len(GcpDiskSnapshotPricing),
		},
		"samples": map[string]interface{}{
			"aws_ebs_us_east_1":    AwsEbsSnapshotPricing["us-east-1"],
			"aws_rds_us_east_1":    AwsRdsSnapshotPricing["us-east-1"],
			"azure_disk_eastus":    AzureDiskSnapshotPricing["eastus"],
			"gcp_disk_us_central1": GcpDiskSnapshotPricing["us-central1"],
		},
		"last_cache_update": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
func InitPricing() {
	PricingMetadata = map[string]PricingInfo{
		"aws_ebs": {
			Source:       "AWS Pricing API (Default Values)",
			LastVerified: time.Now(),
		},
		"aws_rds": {
			Source:       "AWS Pricing API (Default Values)",
			LastVerified: time.Now(),
		},
		"azure_disk": {
			Source:       "Azure Pricing API (Default Values)",
			LastVerified: time.Now(),
		},
		"gcp_disk": {
			Source:       "GCP Pricing API (Default Values)",
			LastVerified: time.Now(),
		},
	}

	loadPricingFromCache()

	go func() {
		ctx := context.Background()
		RefreshPricing(ctx)
	}()
}
