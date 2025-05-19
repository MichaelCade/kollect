package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

type AzureData struct {
	AzureVMs             []armcompute.VirtualMachine
	AzureVMSS            []armcompute.VirtualMachineScaleSet
	AzureAKSClusters     []armcontainerservice.ManagedCluster
	AzureStorageAccounts []armstorage.Account
	AzureBlobContainers  []armstorage.ListContainerItem
	AzureVirtualNetworks []armnetwork.VirtualNetwork
	AzureSQLDatabases    []armsql.Database
	AzureCosmosDBs       []armcosmos.DatabaseAccountGetResults
	AzureResourceGroups  []armresources.ResourceGroup
}

func CheckCredentials(ctx context.Context) (bool, error) {
	err := error(nil)

	_, err = getAzureSubscriptionID()

	return err == nil, err
}

func CollectAzureData(ctx context.Context) (AzureData, error) {
	var data AzureData
	subscriptionID, err := getAzureSubscriptionID()
	if err != nil {
		return data, fmt.Errorf("failed to get Azure subscription ID: %v", err)
	}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return data, err
	}
	rgClient, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
	if err != nil {
		return data, err
	}
	rgPager := rgClient.NewListPager(nil)
	for rgPager.More() {
		page, err := rgPager.NextPage(ctx)
		if err != nil {
			log.Printf("Warning: Failed to get Resource Groups: %v", err)
			break
		}
		for _, rg := range page.Value {
			data.AzureResourceGroups = append(data.AzureResourceGroups, *rg)
		}
	}

	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		return data, err
	}
	vmPager := vmClient.NewListAllPager(nil)
	for vmPager.More() {
		page, err := vmPager.NextPage(ctx)
		if err != nil {
			log.Printf("Warning: Failed to get VMs: %v", err)
			break
		}
		for _, vm := range page.Value {
			data.AzureVMs = append(data.AzureVMs, *vm)
		}
	}

	vmssClient, err := armcompute.NewVirtualMachineScaleSetsClient(subscriptionID, cred, nil)
	if err != nil {
		log.Printf("Warning: Failed to create VMSS client: %v", err)
	} else {
		vmssPager := vmssClient.NewListAllPager(nil)
		for vmssPager.More() {
			page, err := vmssPager.NextPage(ctx)
			if err != nil {
				log.Printf("Warning: Failed to get VMSS: %v", err)
				break
			}
			for _, vmss := range page.Value {
				data.AzureVMSS = append(data.AzureVMSS, *vmss)
			}
		}
	}

	aksClient, err := armcontainerservice.NewManagedClustersClient(subscriptionID, cred, nil)
	if err != nil {
		log.Printf("Warning: Failed to create AKS client: %v", err)
	} else {
		aksPager := aksClient.NewListPager(nil)
		for aksPager.More() {
			page, err := aksPager.NextPage(ctx)
			if err != nil {
				log.Printf("Warning: Failed to get AKS Clusters: %v", err)
				break
			}
			for _, aks := range page.Value {
				data.AzureAKSClusters = append(data.AzureAKSClusters, *aks)
			}
		}
	}

	storageClient, err := armstorage.NewAccountsClient(subscriptionID, cred, nil)
	if err != nil {
		log.Printf("Warning: Failed to create Storage client: %v", err)
	} else {
		storagePager := storageClient.NewListPager(nil)
		for storagePager.More() {
			page, err := storagePager.NextPage(ctx)
			if err != nil {
				log.Printf("Warning: Failed to get Storage Accounts: %v", err)
				break
			}
			for _, account := range page.Value {
				data.AzureStorageAccounts = append(data.AzureStorageAccounts, *account)
			}
		}
	}

	for _, account := range data.AzureStorageAccounts {
		resourceGroup := getResourceGroupFromID(*account.ID)
		blobClient, err := armstorage.NewBlobContainersClient(subscriptionID, cred, nil)
		if err != nil {
			log.Printf("Warning: Failed to create Blob client: %v", err)
			continue
		}
		blobPager := blobClient.NewListPager(resourceGroup, *account.Name, nil)
		for blobPager.More() {
			page, err := blobPager.NextPage(ctx)
			if err != nil {
				log.Printf("Warning: Failed to get Blob Containers: %v", err)
				break
			}
			for _, container := range page.Value {
				data.AzureBlobContainers = append(data.AzureBlobContainers, *container)
			}
		}
	}

	vnetClient, err := armnetwork.NewVirtualNetworksClient(subscriptionID, cred, nil)
	if err != nil {
		log.Printf("Warning: Failed to create VNet client: %v", err)
	} else {
		vnetPager := vnetClient.NewListAllPager(nil)
		for vnetPager.More() {
			page, err := vnetPager.NextPage(ctx)
			if err != nil {
				log.Printf("Warning: Failed to get Virtual Networks: %v", err)
				break
			}
			for _, vnet := range page.Value {
				data.AzureVirtualNetworks = append(data.AzureVirtualNetworks, *vnet)
			}
		}
	}

	sqlClient, err := armsql.NewDatabasesClient(subscriptionID, cred, nil)
	if err != nil {
		log.Printf("Warning: Failed to create SQL client: %v", err)
	} else {
		sqlServerClient, err := armsql.NewServersClient(subscriptionID, cred, nil)
		if err != nil {
			log.Printf("Warning: Failed to create SQL server client: %v", err)
		} else {
			sqlServerPager := sqlServerClient.NewListPager(nil)
			for sqlServerPager.More() {
				page, err := sqlServerPager.NextPage(ctx)
				if err != nil {
					log.Printf("Warning: Failed to get SQL Servers: %v", err)
					break
				}
				for _, server := range page.Value {
					resourceGroup := getResourceGroupFromID(*server.ID)
					dbPager := sqlClient.NewListByServerPager(resourceGroup, *server.Name, nil)
					for dbPager.More() {
						dbPage, err := dbPager.NextPage(ctx)
						if err != nil {
							log.Printf("Warning: Failed to get SQL Databases: %v", err)
							break
						}
						for _, db := range dbPage.Value {
							data.AzureSQLDatabases = append(data.AzureSQLDatabases, *db)
						}
					}
				}
			}
		}
	}

	cosmosClient, err := armcosmos.NewDatabaseAccountsClient(subscriptionID, cred, nil)
	if err != nil {
		log.Printf("Warning: Failed to create CosmosDB client: %v", err)
	} else {
		cosmosPager := cosmosClient.NewListPager(nil)
		for cosmosPager.More() {
			page, err := cosmosPager.NextPage(ctx)
			if err != nil {
				log.Printf("Warning: Failed to get CosmosDB Accounts: %v", err)
				break
			}
			for _, db := range page.Value {
				data.AzureCosmosDBs = append(data.AzureCosmosDBs, *db)
			}
		}
	}

	return data, nil
}

func getResourceGroupFromID(resourceID string) string {
	parts := strings.Split(resourceID, "/")
	for i, part := range parts {
		if part == "resourceGroups" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

func getAzureSubscriptionID() (string, error) {
	cmd := exec.Command("az", "account", "show", "--query", "id", "--output", "json")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var subscriptionID string
	if err := json.Unmarshal(output, &subscriptionID); err != nil {
		return "", err
	}

	return subscriptionID, nil
}

func CollectSnapshotData(ctx context.Context) (map[string]interface{}, error) {
	snapshots := map[string]interface{}{}

	subscriptionID, err := getAzureSubscriptionID()
	if err != nil {
		return nil, fmt.Errorf("failed to get Azure subscription ID: %v", err)
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %v", err)
	}

	resourceClient, err := armresources.NewProvidersClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource providers client: %v", err)
	}

	provider, err := resourceClient.Get(ctx, "Microsoft.Compute", nil)
	if err != nil {
		log.Printf("Warning: Failed to get compute provider, using default locations: %v", err)
		return collectSnapshotsFromDefaultRegions(ctx, subscriptionID, cred)
	}

	var locations []string
	for _, resourceType := range provider.ResourceTypes {
		if *resourceType.ResourceType == "snapshots" {
			for _, location := range resourceType.Locations {
				locations = append(locations, *location)
			}
			break
		}
	}

	if len(locations) == 0 {
		log.Printf("No snapshot locations found in provider, using default locations")
		return collectSnapshotsFromDefaultRegions(ctx, subscriptionID, cred)
	}

	var allDiskSnapshots []map[string]string

	for _, location := range locations {
		log.Printf("Collecting disk snapshots from location: %s", location)
		regionSnapshots, err := collectDiskSnapshotsFromRegion(ctx, subscriptionID, cred, location)
		if err != nil {
			log.Printf("Warning: Failed to collect disk snapshots from location %s: %v", location, err)
		} else if len(regionSnapshots) > 0 {
			for i := range regionSnapshots {
				if _, exists := regionSnapshots[i]["Location"]; !exists {
					regionSnapshots[i]["Location"] = location
				}
			}
			allDiskSnapshots = append(allDiskSnapshots, regionSnapshots...)
		}
	}

	if len(allDiskSnapshots) > 0 {
		snapshots["DiskSnapshots"] = allDiskSnapshots
	}

	return snapshots, nil
}
func collectSnapshotsFromDefaultRegions(ctx context.Context, subscriptionID string, cred *azidentity.DefaultAzureCredential) (map[string]interface{}, error) {
	snapshots := map[string]interface{}{}
	defaultRegions := []string{
		"eastus", "eastus2", "westus", "westus2", "centralus",
		"northeurope", "westeurope", "uksouth", "southeastasia",
		"eastasia", "australiaeast", "japaneast",
	}

	var allDiskSnapshots []map[string]string

	for _, location := range defaultRegions {
		log.Printf("Collecting disk snapshots from default location: %s", location)
		regionSnapshots, err := collectDiskSnapshotsFromRegion(ctx, subscriptionID, cred, location)
		if err != nil {
			log.Printf("Warning: Failed to collect disk snapshots from location %s: %v", location, err)
		} else if len(regionSnapshots) > 0 {
			for i := range regionSnapshots {
				if _, exists := regionSnapshots[i]["Location"]; !exists {
					regionSnapshots[i]["Location"] = location
				}
			}
			allDiskSnapshots = append(allDiskSnapshots, regionSnapshots...)
		}
	}

	if len(allDiskSnapshots) > 0 {
		snapshots["DiskSnapshots"] = allDiskSnapshots
	}

	return snapshots, nil
}

func collectDiskSnapshotsFromRegion(ctx context.Context, subscriptionID string, cred *azidentity.DefaultAzureCredential, location string) ([]map[string]string, error) {
	snapshotClient, err := armcompute.NewSnapshotsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshots client: %v", err)
	}

	pager := snapshotClient.NewListPager(nil)

	var snapshots []map[string]string
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of snapshots: %v", err)
		}

		for _, snapshot := range page.Value {
			if snapshot.Location != nil && strings.EqualFold(*snapshot.Location, location) {
				snapshotInfo := map[string]string{
					"Name":              *snapshot.Name,
					"Location":          location,
					"ProvisioningState": string(*snapshot.Properties.ProvisioningState),
				}

				if snapshot.ID != nil {
					snapshotInfo["ID"] = *snapshot.ID
				}

				if snapshot.Properties.TimeCreated != nil {
					snapshotInfo["TimeCreated"] = snapshot.Properties.TimeCreated.Format(time.RFC3339)
				}

				if snapshot.Properties.DiskSizeGB != nil {
					snapshotInfo["DiskSizeGB"] = fmt.Sprintf("%d", *snapshot.Properties.DiskSizeGB)
				}

				snapshots = append(snapshots, snapshotInfo)
			}
		}
	}

	return snapshots, nil
}

func collectDiskSnapshots(ctx context.Context, subscriptionID string, cred *azidentity.DefaultAzureCredential) ([]map[string]string, error) {
	snapshotsClient, err := armcompute.NewSnapshotsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshots client: %v", err)
	}

	pager := snapshotsClient.NewListPager(nil)
	snapshots := []map[string]string{}

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %v", err)
		}

		for _, snapshot := range page.Value {
			snapshotInfo := map[string]string{
				"Name":     *snapshot.Name,
				"ID":       *snapshot.ID,
				"Location": *snapshot.Location,
			}

			if snapshot.Properties != nil {
				if snapshot.Properties.TimeCreated != nil {
					snapshotInfo["CreationTime"] = snapshot.Properties.TimeCreated.Format(time.RFC3339)
				}

				if snapshot.Properties.DiskSizeGB != nil {
					snapshotInfo["SizeGB"] = fmt.Sprintf("%d", *snapshot.Properties.DiskSizeGB)
				}

				if snapshot.Properties.ProvisioningState != nil {
					snapshotInfo["ProvisioningState"] = *snapshot.Properties.ProvisioningState
				}

				if snapshot.Properties.DiskState != nil {
					snapshotInfo["State"] = string(*snapshot.Properties.DiskState)
				}
			}

			snapshots = append(snapshots, snapshotInfo)
		}
	}

	return snapshots, nil
}
