document.getElementById('veeam-button').addEventListener('click', () => {
    showLoadingIndicator();
    fetch('/api/switch?type=veeam')
        .then(response => response.json())
        .then(data => {
            location.reload();
        })
        .catch(error => console.error('Error switching to Veeam:', error))
        .finally(() => hideLoadingIndicator());
});

document.addEventListener('htmx:afterSwap', (event) => {
    if (event.detail.target.id === 'hidden-content') {
        try {
            const data = JSON.parse(event.detail.xhr.responseText);
            console.log("Fetched Data:", data); // Log fetched data
            const content = document.getElementById('content');
            const template = document.getElementById('table-template').content;
            function createTable(headerText, data, rowTemplate, headers, repositories) {
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
                    row.innerHTML = rowTemplate(item, repositories);
                    tbody.appendChild(row);
                });
                content.appendChild(table);
            }
            if (data.ServerInfo) {
                createTable('Server Info', [data.ServerInfo], serverInfoRowTemplate, ['Name', 'Build Version', 'Database Vendor', 'SQL Server Version', 'VBR ID']);
            }
            if (data.Credentials) {
                createTable('Credentials', data.Credentials, credentialsRowTemplate, ['Username', 'Description', 'Type']);
            }
            if (data.CloudCredentials) {
                createTable('Cloud Credentials', data.CloudCredentials, cloudCredentialsRowTemplate, ['Account', 'Description', 'Type']);
            }
            if (data.KMSServers) {
                createTable('KMS Servers', data.KMSServers, kmsServersRowTemplate, ['ID', 'Name']);
            }
            if (data.ManagedServers) {
                createTable('Managed Servers', data.ManagedServers, managedServersRowTemplate, ['Name', 'Type', 'Status', 'Description']);
            }
            if (data.Repositories) {
                createTable('Repositories', data.Repositories, repositoriesRowTemplate, ['Name', 'Type', 'Description', 'Folder Name', 'Immutability Status']);
            }
            if (data.ScaleOutRepositories) {
                createTable('Scale-Out Repositories', data.ScaleOutRepositories, scaleOutRepositoriesRowTemplate, ['Name', 'Description', 'Details'], data.Repositories);
            }
            if (data.Proxies) {
                createTable('Proxies', data.Proxies, proxiesRowTemplate, ['Name', 'Type', 'Description', 'Max Task Count', 'Transport Mode']);
            }
            if (data.BackupJobs) {
                createTable('Backup Jobs', data.BackupJobs, backupJobsRowTemplate, ['Job Name', 'ID', 'Description', 'Type', 'Is Disabled', 'Is High Priority', 'Job Details']);
            }
        } catch (error) {
            console.error("Error processing data:", error);
        }
    }
});

function serverInfoRowTemplate(item) {
    return `<td>${item.name}</td><td>${item.buildVersion}</td><td>${item.databaseVendor}</td><td>${item.sqlServerVersion}</td><td>${item.vbrId}</td>`;
}

function credentialsRowTemplate(item) {
    return `<td>${item.username}</td><td>${item.description}</td><td>${item.type}</td>`;
}

function cloudCredentialsRowTemplate(item) {
    let details = '';
    switch (item.type) {
        case 'Amazon':
            details = `${item.accessKey}`;
            break;
        case 'AzureCompute':
            details = `${item.connectionName}`;
            break;
        case 'Google':
            details = `${item.accessKey}`;
            break;
        case 'AzureStorage':
            details = `${item.account}`;
            break;
        case 'S3Compatible':
            details = `${item.accessKey}`;
            break;
        default:
            details = 'N/A';
    }
    return `<td>${details}</td><td>${item.type}</td><td>${item.description}</td>`;
}

function kmsServersRowTemplate(item) {
    return `<td>${item.id}</td><td>${item.name}</td>`;
}

function managedServersRowTemplate(item) {
    return `<td>${item.name}</td><td>${item.type}</td><td>${item.status}</td><td>${item.description}</td>`;
}

function repositoriesRowTemplate(item) {
    const immutabilityStatus = item.bucket && item.bucket.immutability ? item.bucket.immutability.isEnabled : false;
    return `<td>${item.name}</td><td>${item.type}</td><td>${item.description}</td><td>${item.bucket ? item.bucket.folderName : 'N/A'}</td><td>${immutabilityStatus}</td>`;
}

function scaleOutRepositoriesRowTemplate(item, repositories) {
    if (!repositories) {
        console.error("Repositories data is undefined");
        return '';
    }

    const performanceTier = item.performanceTier && item.performanceTier.performanceExtents ? 'Enabled' : 'Disabled';
    const capacityTier = item.capacityTier && item.capacityTier.isEnabled ? 'Enabled' : 'Disabled';
    const archiveTier = item.archiveTier && item.archiveTier.isEnabled ? 'Enabled' : 'Disabled';
    const copyPolicy = item.capacityTier && item.capacityTier.copyPolicyEnabled ? 'Enabled' : 'Disabled';
    const movePolicy = item.capacityTier && item.capacityTier.movePolicyEnabled ? 'Enabled' : 'Disabled';

    const performanceExtents = item.performanceTier && item.performanceTier.performanceExtents ? item.performanceTier.performanceExtents.map(extent => {
        const repo = repositories.find(repo => repo.id === extent.id);
        return `<li>${repo ? repo.name : extent.id}</li>`;
    }).join('') : 'N/A';

    const capacityExtents = item.capacityTier && item.capacityTier.extents ? item.capacityTier.extents.map(extent => {
        const repo = repositories.find(repo => repo.id === extent.id);
        return `<li>${repo ? repo.name : extent.id}</li>`;
    }).join('') : 'N/A';

    const archiveExtents = item.archiveTier && item.archiveTier.extentId ? (() => {
        const repo = repositories.find(repo => repo.id === item.archiveTier.extentId);
        return `<li>${repo ? repo.name : item.archiveTier.extentId}</li>`;
    })() : 'N/A';

    return `
        <td>${item.name}</td>
        <td>${item.description}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${item.id}')"><i class="fas fa-info-circle"></i> Details</button>
            <div id="details-${item.id}" style="display:none;">
                <p>Performance Tier: ${performanceTier}</p>
                <ul>${performanceExtents}</ul>
                <p>Capacity Tier: ${capacityTier}</p>
                <ul>${capacityExtents}</ul>
                <p>Archive Tier: ${archiveTier}</p>
                <ul>${archiveExtents}</ul>
                <p>Copy Policy: ${copyPolicy}</p>
                <p>Move Policy: ${movePolicy}</p>
            </div>
        </td>
    `;
}

function toggleDetails(id) {
    const details = document.getElementById(`details-${id}`);
    details.style.display = details.style.display === 'none' ? 'block' : 'none';
}

function proxiesRowTemplate(item) {
    return `<td>${item.name}</td><td>${item.type}</td><td>${item.description}</td><td>${item.server.maxTaskCount}</td><td>${item.server.transportMode}</td>`;
}

function backupJobsRowTemplate(item) {
    const vms = item.virtualMachines && item.virtualMachines.includes 
        ? item.virtualMachines.includes.map(vm => `<li>Name: ${vm.name}, Host: ${vm.hostName}, Size: ${vm.size}</li>`).join('') 
        : '';
    const retentionPolicy = item.storage && item.storage.retentionPolicy 
        ? `${item.storage.retentionPolicy.type} for ${item.storage.retentionPolicy.quantity} days` 
        : 'N/A';
    const dailySchedule = item.schedule && item.schedule.daily 
        ? `${item.schedule.daily.dailyKind} at ${item.schedule.daily.localTime}` 
        : 'N/A';
    return `
        <td>${item.name}</td>
        <td>${item.id}</td>
        <td>${item.description}</td>
        <td>${item.type}</td>
        <td>${item.isDisabled}</td>
        <td>${item.isHighPriority}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${item.id}')"><i class="fas fa-info-circle"></i> Details</button>
            <div id="details-${item.id}" style="display:none;">
                <p>Included VMs:</p>
                <ul>${vms}</ul>
                <p>Backup Repository ID: ${item.storage ? item.storage.backupRepositoryId : 'N/A'}</p>
                <p>Retention Policy: ${retentionPolicy}</p>
                <p>Run Automatically: ${item.schedule ? item.schedule.runAutomatically : 'N/A'}</p>
                <p>Daily Schedule: ${dailySchedule}</p>
            </div>
        </td>
    `;
}

function toggleDetails(id) {
    const details = document.getElementById(`details-${id}`);
    details.style.display = details.style.display === 'none' ? 'block' : 'none';
}