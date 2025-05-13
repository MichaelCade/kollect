package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
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

	// Collect VMs
	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		return data, err
	}
	vmPager := vmClient.NewListAllPager(nil)
	for vmPager.More() {
		page, err := vmPager.NextPage(ctx)
		if err != nil {
			log.Fatalf("Failed to get VMs: %v", err)
		}
		for _, vm := range page.Value {
			data.AzureVMs = append(data.AzureVMs, *vm)
		}
	}

	// Collect VMSS
	vmssClient, err := armcompute.NewVirtualMachineScaleSetsClient(subscriptionID, cred, nil)
	if err != nil {
		return data, err
	}
	vmssPager := vmssClient.NewListAllPager(nil)
	for vmssPager.More() {
		page, err := vmssPager.NextPage(ctx)
		if err != nil {
			log.Fatalf("Failed to get VMSS: %v", err)
		}
		for _, vmss := range page.Value {
			data.AzureVMSS = append(data.AzureVMSS, *vmss)
		}
	}

	// Collect AKS Clusters
	aksClient, err := armcontainerservice.NewManagedClustersClient(subscriptionID, cred, nil)
	if err != nil {
		return data, err
	}
	aksPager := aksClient.NewListPager(nil)
	for aksPager.More() {
		page, err := aksPager.NextPage(ctx)
		if err != nil {
			log.Fatalf("Failed to get AKS Clusters: %v", err)
		}
		for _, aks := range page.Value {
			data.AzureAKSClusters = append(data.AzureAKSClusters, *aks)
		}
	}

	// Collect Storage Accounts
	storageClient, err := armstorage.NewAccountsClient(subscriptionID, cred, nil)
	if err != nil {
		return data, err
	}
	storagePager := storageClient.NewListPager(nil)
	for storagePager.More() {
		page, err := storagePager.NextPage(ctx)
		if err != nil {
			log.Fatalf("Failed to get Storage Accounts: %v", err)
		}
		for _, account := range page.Value {
			data.AzureStorageAccounts = append(data.AzureStorageAccounts, *account)
		}
	}

	// Collect Blob Containers
	for _, account := range data.AzureStorageAccounts {
		resourceGroup := getResourceGroupFromID(*account.ID)
		blobClient, err := armstorage.NewBlobContainersClient(subscriptionID, cred, nil)
		if err != nil {
			return data, err
		}
		blobPager := blobClient.NewListPager(resourceGroup, *account.Name, nil)
		for blobPager.More() {
			page, err := blobPager.NextPage(ctx)
			if err != nil {
				log.Fatalf("Failed to get Blob Containers: %v", err)
			}
			for _, container := range page.Value {
				data.AzureBlobContainers = append(data.AzureBlobContainers, *container)
			}
		}
	}

	// Collect Virtual Networks
	vnetClient, err := armnetwork.NewVirtualNetworksClient(subscriptionID, cred, nil)
	if err != nil {
		return data, err
	}
	vnetPager := vnetClient.NewListAllPager(nil)
	for vnetPager.More() {
		page, err := vnetPager.NextPage(ctx)
		if err != nil {
			log.Fatalf("Failed to get Virtual Networks: %v", err)
		}
		for _, vnet := range page.Value {
			data.AzureVirtualNetworks = append(data.AzureVirtualNetworks, *vnet)
		}
	}

	// Collect SQL Databases
	sqlClient, err := armsql.NewDatabasesClient(subscriptionID, cred, nil)
	if err != nil {
		return data, err
	}
	// List SQL servers first to get the server names
	sqlServerClient, err := armsql.NewServersClient(subscriptionID, cred, nil)
	if err != nil {
		return data, err
	}
	sqlServerPager := sqlServerClient.NewListPager(nil)
	for sqlServerPager.More() {
		page, err := sqlServerPager.NextPage(ctx)
		if err != nil {
			log.Fatalf("Failed to get SQL Servers: %v", err)
		}
		for _, server := range page.Value {
			resourceGroup := getResourceGroupFromID(*server.ID)
			dbPager := sqlClient.NewListByServerPager(resourceGroup, *server.Name, nil)
			for dbPager.More() {
				dbPage, err := dbPager.NextPage(ctx)
				if err != nil {
					log.Fatalf("Failed to get SQL Databases: %v", err)
				}
				for _, db := range dbPage.Value {
					data.AzureSQLDatabases = append(data.AzureSQLDatabases, *db)
				}
			}
		}
	}

	// Collect CosmosDB Accounts
	cosmosClient, err := armcosmos.NewDatabaseAccountsClient(subscriptionID, cred, nil)
	if err != nil {
		return data, err
	}
	cosmosPager := cosmosClient.NewListPager(nil)
	for cosmosPager.More() {
		page, err := cosmosPager.NextPage(ctx)
		if err != nil {
			log.Fatalf("Failed to get CosmosDB Accounts: %v", err)
		}
		for _, db := range page.Value {
			data.AzureCosmosDBs = append(data.AzureCosmosDBs, *db)
		}
	}

	return data, nil
}

// Helper function to extract resource group from resource ID
func getResourceGroupFromID(resourceID string) string {
	parts := strings.Split(resourceID, "/")
	for i, part := range parts {
		if part == "resourceGroups" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// Helper function to get the Azure subscription ID using Azure CLI
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
