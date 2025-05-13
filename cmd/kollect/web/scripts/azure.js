// azure.js


registerDataHandler('azure', 
    function(data) {
        return data.AzureVMs || data.AzureResourceGroups || data.AzureStorageAccounts ||
               data.AzureBlobContainers || data.AzureVirtualNetworks || 
               data.AzureSQLDatabases || data.AzureCosmosDBs;
    },
    function(data) {
        console.log("Processing Azure data");
        
        if (data.AzureResourceGroups) {
            createTable('Azure Resource Groups', data.AzureResourceGroups, azureResourceGroupRowTemplate, 
                ['Name', 'Location', 'Tags', 'Provisioning State']);
        }
        
        
        if (data.AzureVMs) {
            createTable('Azure VMs', data.AzureVMs, azureVMRowTemplate, 
                ['Name', 'Location', 'VM Size']);
        }
        
        if (data.AzureVMSS) {
            createTable('Azure VM Scale Sets', data.AzureVMSS, azureVMSSRowTemplate, 
                ['Name', 'Location', 'Capacity']);
        }
        
        if (data.AzureAKSClusters) {
            createTable('Azure AKS Clusters', data.AzureAKSClusters, azureAKSClusterRowTemplate, 
                ['Name', 'Location', 'K8s Version', 'Node Count']);
        }
        
        if (data.AzureStorageAccounts) {
            createTable('Azure Storage Accounts', data.AzureStorageAccounts, azureStorageAccountRowTemplate, 
                ['Name', 'Location', 'Kind']);
        }
        
        if (data.AzureBlobContainers) {
            createTable('Azure Blob Containers', data.AzureBlobContainers, azureBlobContainerRowTemplate, 
                ['Name', 'Immutable', 'ID']);
        }
        
        if (data.AzureVirtualNetworks) {
            createTable('Azure Virtual Networks', data.AzureVirtualNetworks, azureVirtualNetworkRowTemplate, 
                ['Name', 'Location']);
        }
        
        if (data.AzureSQLDatabases) {
            createTable('Azure SQL Databases', data.AzureSQLDatabases, azureSQLDatabaseRowTemplate, 
                ['Name', 'Location']);
        }
        
        if (data.AzureCosmosDBs) {
            createTable('Azure CosmosDB Accounts', data.AzureCosmosDBs, azureCosmosDBRowTemplate, 
                ['Name', 'Location']);
        }
        
        setTimeout(() => {
            console.log(`Created Azure tables`);
        }, 100);
    }
);

function azureResourceGroupRowTemplate(item) {
    const tags = item.tags ? Object.entries(item.tags).map(([key, value]) => 
        `${key}: ${value}`).join(', ') : 'No tags';
        
    return `
        <td>${item.name}</td>
        <td>${item.location}</td>
        <td>${tags}</td>
        <td>${item.properties.provisioningState}</td>
    `;
}

function azureVMRowTemplate(item) {
    let vmSize = 'N/A';
    if (item.properties && item.properties.hardwareProfile) {
        vmSize = item.properties.hardwareProfile.vmSize;
    }
    return `<td>${item.name}</td><td>${item.location}</td><td>${vmSize}</td>`;
}

function azureVMSSRowTemplate(item) {
    let capacity = 'N/A';
    if (item.sku && item.sku.capacity) {
        capacity = item.sku.capacity;
    }
    return `<td>${item.name}</td><td>${item.location}</td><td>${capacity}</td>`;
}

function azureAKSClusterRowTemplate(item) {
    let version = 'N/A';
    let nodeCount = 'N/A';
    
    if (item.properties && item.properties.kubernetesVersion) {
        version = item.properties.kubernetesVersion;
    }
    
    if (item.properties && item.properties.agentPoolProfiles && item.properties.agentPoolProfiles.length > 0) {
        nodeCount = item.properties.agentPoolProfiles[0].count || 'N/A';
    }
    
    return `<td>${item.name}</td><td>${item.location}</td><td>${version}</td><td>${nodeCount}</td>`;
}

function azureStorageAccountRowTemplate(item) {
    return `<td>${item.name}</td><td>${item.location}</td><td>${item.kind}</td>`;
}

function azureBlobContainerRowTemplate(item) {
    let immutable = false;
    if (item.properties && item.properties.immutableStorageWithVersioning) {
        immutable = item.properties.immutableStorageWithVersioning.enabled;
    }
    return `<td>${item.name}</td><td>${immutable}</td><td>${item.id}</td>`;
}

function azureVirtualNetworkRowTemplate(item) {
    return `<td>${item.name}</td><td>${item.location}</td>`;
}

function azureSQLDatabaseRowTemplate(item) {
    return `<td>${item.name}</td><td>${item.location}</td>`;
}

function azureCosmosDBRowTemplate(item) {
    return `<td>${item.name}</td><td>${item.location}</td>`;
}