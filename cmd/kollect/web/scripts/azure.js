// azure.js

document.getElementById('azure-button').addEventListener('click', () => {
    showLoadingIndicator();
    fetch('/api/switch?type=azure')
        .then(response => response.json())
        .then(data => {
            location.reload();
        })
        .catch(error => console.error('Error switching to Azure:', error))
        .finally(() => hideLoadingIndicator());
});

document.addEventListener('htmx:afterSwap', (event) => {
    if (event.detail.target.id === 'hidden-content') {
        try {
            const data = JSON.parse(event.detail.xhr.responseText);
            console.log("Fetched Data:", data); // Log fetched data
            const content = document.getElementById('content');
            const template = document.getElementById('table-template').content;
            function createTable(headerText, data, rowTemplate, headers) {
                if (!data || data.length === 0) return; // Ensure data is not null or empty
                const table = template.cloneNode(true);
                table.querySelector('th').textContent = headerText;
                const thead = table.querySelector('thead');
                const headerRow = document.createElement('tr');
                headers.forEach(header => {
                    const th = document.createElement('th');
                    th.textContent = header;
                    headerRow.appendChild(th);
                });
                thead.appendChild(headerRow);
                const tbody = table.querySelector('tbody');
                data.forEach(item => {
                    const row = document.createElement('tr');
                    row.innerHTML = rowTemplate(item);
                    tbody.appendChild(row);
                });
                content.appendChild(table);
            }
            if (data.AzureResourceGroups) {
                createTable('Azure Resource Groups', data.AzureResourceGroups, azureResourceGroupRowTemplate, 
                    ['Name', 'Location', 'Tags', 'Provisioning State']);
            }
            if (data.AzureVMs) {
                createTable('Azure VMs', data.AzureVMs, azureVMRowTemplate, ['Name', 'Location', 'VM Size']);
            }
            if (data.AzureStorageAccounts) {
                createTable('Azure Storage Accounts', data.AzureStorageAccounts, azureStorageAccountRowTemplate, ['Name', 'Location', 'Kind']);
            }
            if (data.AzureBlobContainers) {
                createTable('Azure Blob Containers', data.AzureBlobContainers, azureBlobContainerRowTemplate, ['Name', 'Immutable', 'ID']);
            }
            if (data.AzureVirtualNetworks) {
                createTable('Azure Virtual Networks', data.AzureVirtualNetworks, azureVirtualNetworkRowTemplate, ['Name', 'Location']);
            }
            if (data.AzureSQLDatabases) {
                createTable('Azure SQL Databases', data.AzureSQLDatabases, azureSQLDatabaseRowTemplate, ['Name', 'Location']);
            }
            if (data.AzureCosmosDBs) {
                createTable('Azure CosmosDB Accounts', data.AzureCosmosDBs, azureCosmosDBRowTemplate, ['Name', 'Location']);
            }
        } catch (error) {
            console.error("Error processing data:", error);
        }
    }
});

function azureVMRowTemplate(item) {
    return `<td>${item.name}</td><td>${item.location}</td><td>${item.properties.hardwareProfile.vmSize}</td>`;
}

function azureStorageAccountRowTemplate(item) {
    return `<td>${item.name}</td><td>${item.location}</td><td>${item.kind}</td>`;
}

function azureBlobContainerRowTemplate(item) {
    return `<td>${item.name}</td><td>${item.properties.immutableStorageWithVersioning.enabled}</td><td>${item.id}</td>`;
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