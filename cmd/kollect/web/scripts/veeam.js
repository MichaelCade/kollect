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
            console.log("Fetched Data:", data);

            // Show charts container if we have Veeam data
            const chartsContainer = document.getElementById('charts-container');
            if (data.ServerInfo || data.BackupJobs || data.Repositories) {
                chartsContainer.style.display = 'grid';
            } else {
                chartsContainer.style.display = 'none';
            }

            const content = document.getElementById('content');
            const template = document.getElementById('table-template').content;

            function createTable(headerText, data, rowTemplate, headers, repositories) {
                if (!data || data.length === 0) return;

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

            // Generate tables for each dataset
            if (data.ServerInfo) {
                createTable('Server Info', [data.ServerInfo], serverInfoRowTemplate, ['Name', 'Build Version', 'Database Vendor', 'SQL Server Version', 'VBR ID']);
            }
            if (data.Credentials) {
                createTable('Credentials', data.Credentials, credentialsRowTemplate, ['Username', 'Description', 'Type']);
            }
            if (data.CloudCredentials) {
                createTable('Cloud Credentials', data.CloudCredentials, cloudCredentialsRowTemplate, ['Account', 'Type', 'Description']);
            }
            if (data.KMSServers) {
                createTable('KMS Servers', data.KMSServers, kmsServersRowTemplate, ['ID', 'Name']);
            }
            if (data.ManagedServers) {
                createTable('Managed Servers', data.ManagedServers, managedServersRowTemplate, ['Name', 'Type', 'Status', 'Description']);
            }
            if (data.Repositories) {
                createTable('Repositories', data.Repositories, repositoriesRowTemplate, ['Name', 'Type', 'Description', 'Bucket Name', 'Folder Name', 'Region ID', 'Infrequent Access Storage', 'Immutability Status', 'Immutable Period']);
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

            // Generate charts after creating tables
            generateCharts(data);
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
    let immutabilityStatus = false;
    let immutablePeriod = 'N/A';
    let bucketName = 'N/A';
    let regionId = 'N/A';
    let infrequentAccessStorage = 'N/A';

    if (item.repository) {
        if (item.repository.makeRecentBackupsImmutableDays) {
            immutabilityStatus = true;
            immutablePeriod = item.repository.makeRecentBackupsImmutableDays;
        }
    } else if (item.bucket) {
        bucketName = item.bucket.bucketName || 'N/A';
        regionId = item.bucket.regionId || 'N/A';
        infrequentAccessStorage = item.bucket.infrequentAccessStorage && item.bucket.infrequentAccessStorage.isEnabled ? 'Enabled' : 'Disabled';

        if (item.bucket.immutability && item.bucket.immutability.isEnabled) {
            immutabilityStatus = true;
            immutablePeriod = item.bucket.immutability.daysCount;
        } else if (item.bucket.immutabilityEnabled) {
            immutabilityStatus = true;
            immutablePeriod = 'N/A';
        }
    }

    return `<td>${item.name}</td><td>${item.type}</td><td>${item.description}</td><td>${bucketName}</td><td>${item.bucket ? item.bucket.folderName : 'N/A'}</td><td>${regionId}</td><td>${infrequentAccessStorage}</td><td>${immutabilityStatus}</td><td>${immutablePeriod}</td>`;
}

function scaleOutRepositoriesRowTemplate(item, repositories) {
    if (!repositories || !Array.isArray(repositories)) {
        repositories = [];
        console.warn("Repositories data is undefined or not an array");
    }

    // Make properties safe to access with fallbacks
    const performanceTier = (item.performanceTier && 
                            (item.performanceTier.performanceExtents || 
                             item.performanceTier.type)) ? 'Enabled' : 'Disabled';
                             
    const capacityTier = (item.capacityTier && 
                         (item.capacityTier.isEnabled || 
                          item.capacityTier.extentId)) ? 'Enabled' : 'Disabled';
                          
    const archiveTier = (item.archiveTier && 
                        (item.archiveTier.isEnabled || 
                         item.archiveTier.extentId)) ? 'Enabled' : 'Disabled';
                         
    const copyPolicy = item.capacityTier && item.capacityTier.copyPolicyEnabled ? 'Enabled' : 'Disabled';
    const movePolicy = item.capacityTier && item.capacityTier.movePolicyEnabled ? 'Enabled' : 'Disabled';
    const operationalRestorePeriodDays = item.capacityTier && item.capacityTier.operationalRestorePeriodDays ? 
                                         item.capacityTier.operationalRestorePeriodDays : 'N/A';
    const archivePeriodDays = item.archiveTier && item.archiveTier.archivePeriodDays ? 
                              item.archiveTier.archivePeriodDays : 'N/A';

    // Safely handle different possible data structures
    let performanceExtents = 'N/A';
    if (item.performanceTier && Array.isArray(item.performanceTier.performanceExtents)) {
        performanceExtents = item.performanceTier.performanceExtents.map(extent => {
            const repo = repositories.find(repo => repo.id === extent.id);
            return `<li>${repo ? repo.name : extent.id}</li>`;
        }).join('');
    } else if (Array.isArray(item.extentIds)) {
        performanceExtents = item.extentIds.map(id => {
            const repo = repositories.find(repo => repo.id === id);
            return `<li>${repo ? repo.name : id}</li>`;
        }).join('');
    }

    let capacityExtents = 'N/A';
    if (item.capacityTier && Array.isArray(item.capacityTier.extents)) {
        capacityExtents = item.capacityTier.extents.map(extent => {
            const repo = repositories.find(repo => repo.id === extent.id);
            return `<li>${repo ? repo.name : extent.id}</li>`;
        }).join('');
    } else if (item.capacityTier && item.capacityTier.extentId) {
        const repo = repositories.find(repo => repo.id === item.capacityTier.extentId);
        capacityExtents = `<li>${repo ? repo.name : item.capacityTier.extentId}</li>`;
    }

    let archiveExtents = 'N/A';
    if (item.archiveTier && item.archiveTier.extentId) {
        const repo = repositories.find(repo => repo.id === item.archiveTier.extentId);
        archiveExtents = `<li>${repo ? repo.name : item.archiveTier.extentId}</li>`;
    }

    return `
        <td>${item.name || 'Unnamed'}</td>
        <td>${item.description || 'No description'}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${item.id}')"><i class="fas fa-info-circle"></i> Details</button>
            <div id="details-${item.id}" style="display:none;">
                <p>Performance Tier: ${performanceTier}</p>
                <ul>${performanceExtents}</ul>
                <p>Capacity Tier: ${capacityTier}</p>
                <ul>${capacityExtents}</ul>
                <p>Operational Restore Period Days: ${operationalRestorePeriodDays}</p>
                <p>Archive Tier: ${archiveTier}</p>
                <ul>${archiveExtents}</ul>
                <p>Archive Period Days: ${archivePeriodDays}</ul>
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
    let vms = '';
    // Check for VMs in various possible locations in the data
    if (item.virtualMachines && Array.isArray(item.virtualMachines.includes)) {
        vms = item.virtualMachines.includes.map(vm => 
            `<li>Name: ${vm.name || 'Unknown'}, Host: ${vm.hostName || 'Unknown'}, Size: ${vm.size || 'Unknown'}</li>`
        ).join('');
    } else if (Array.isArray(item.sourceObjects)) {
        vms = item.sourceObjects.map(obj => 
            `<li>Name: ${obj.name || 'Unknown'}, Type: ${obj.type || 'Unknown'}</li>`
        ).join('');
    }
    
    // Extract retention policy from different possible structures
    let retentionPolicy = 'N/A';
    if (item.storage && item.storage.retentionPolicy) {
        retentionPolicy = `${item.storage.retentionPolicy.type} for ${item.storage.retentionPolicy.quantity} days`;
    } else if (item.retentionPolicy) {
        if (item.retentionPolicy.type === 'GFS') {
            retentionPolicy = `GFS (Daily: ${item.retentionPolicy.dailyBackups}, Weekly: ${item.retentionPolicy.weeklyBackups}, Monthly: ${item.retentionPolicy.monthlyBackups}, Yearly: ${item.retentionPolicy.yearlyBackups})`;
        } else {
            retentionPolicy = `${item.retentionPolicy.type} for ${item.retentionPolicy.count} days`;
        }
    }
    
    // Extract schedule information from different possible structures
    let dailySchedule = 'N/A';
    if (item.schedule && item.schedule.daily) {
        dailySchedule = `${item.schedule.daily.dailyKind} at ${item.schedule.daily.localTime}`;
    } else if (item.schedule && item.schedule.type) {
        dailySchedule = `${item.schedule.type} at ${item.schedule.time || 'Unknown time'}`;
    }
    
    return `
        <td>${item.name || 'Unnamed'}</td>
        <td>${item.id || 'No ID'}</td>
        <td>${item.description || 'No description'}</td>
        <td>${item.type || 'Unknown'}</td>
        <td>${item.isDisabled || false}</td>
        <td>${item.isHighPriority || false}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${item.id}')"><i class="fas fa-info-circle"></i> Details</button>
            <div id="details-${item.id}" style="display:none;">
                <p>Included VMs:</p>
                <ul>${vms || '<li>No VM data available</li>'}</ul>
                <p>Backup Repository ID: ${(item.storage && item.storage.backupRepositoryId) || item.storageId || 'N/A'}</p>
                <p>Retention Policy: ${retentionPolicy}</p>
                <p>Run Automatically: ${item.schedule ? (item.schedule.runAutomatically || 'Yes') : 'N/A'}</p>
                <p>Daily Schedule: ${dailySchedule}</p>
            </div>
        </td>
    `;
}

function generateCharts(data) {
    try {
        console.log("Data for charts:", data);

        // Get the current theme
        const currentTheme = document.documentElement.getAttribute('data-theme') || 'dark';
        
        // Get the current text color and font properties from CSS variables
        const textColor = getComputedStyle(document.documentElement).getPropertyValue('--text-color').trim() || 
                         (currentTheme === 'dark' ? '#ffffff' : '#333333');
        
        const fontFamily = getComputedStyle(document.documentElement).getPropertyValue('--chart-font-family').trim() || 'Arial, sans-serif';
        const fontSize = parseInt(getComputedStyle(document.documentElement).getPropertyValue('--chart-font-size').trim(), 10) || 12;
        const titleFontSize = parseInt(getComputedStyle(document.documentElement).getPropertyValue('--chart-title-font-size').trim(), 10) || 16;

        // Set Chart.js defaults for this page - explicitly handle dark/light mode
        Chart.defaults.color = textColor;
        Chart.defaults.font.family = fontFamily;
        Chart.defaults.font.size = fontSize;
        
        // Override default legend behavior for better theme compatibility
        Chart.defaults.plugins.legend.labels.color = textColor;

        // Check for Backup Jobs data
        if (data.BackupJobs && data.BackupJobs.length) {
            const backupJobsCtx = document.getElementById('backupJobsChart');
            if (!backupJobsCtx) {
                console.error("Canvas element 'backupJobsChart' not found.");
                return;
            }
            const backupJobsCtx2d = backupJobsCtx.getContext('2d');
            
            const backupJobsData = {
                labels: data.BackupJobs.map(job => job.name || 'Unnamed Job'),
                datasets: [{
                    label: 'Number of Protected VMs',
                    data: data.BackupJobs.map(job => {
                        // Check all possible paths for VM count data
                        if (job.virtualMachines && job.virtualMachines.includes) {
                            return job.virtualMachines.includes.length;
                        } else if (job.sourceObjects) {
                            return job.sourceObjects.length;
                        } else {
                            return 0; // Default to 0 if no data found
                        }
                    }),
                    backgroundColor: data.BackupJobs.map(job => 
                        job.type === 'Backup' ? '#36A2EB' : 
                        job.type === 'VSphereReplica' || job.type === 'ReplicationJob' ? '#FF6384' : 
                        '#CCCCCC'
                    ),
                    borderWidth: 1
                }]
            };
            console.log("Backup Jobs Data:", backupJobsData);
            new Chart(backupJobsCtx2d, {
                type: 'bar',
                data: backupJobsData,
                options: {
                    indexAxis: 'y', // Use horizontal bar chart
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        y: {
                            beginAtZero: true,
                            ticks: {
                                color: textColor,
                                font: {
                                    family: fontFamily,
                                    size: fontSize
                                },
                                autoSkip: false // Ensure all labels are shown
                            }
                        },
                        x: {
                            display: false // Remove the numbers on the bottom axis
                        }
                    },
                    plugins: {
                        title: {
                            display: true,
                            text: 'Number of Protected VMs',
                            color: textColor,
                            font: {
                                family: fontFamily,
                                size: titleFontSize
                            }
                        },
                        legend: {
                            position: 'bottom',
                            labels: {
                                color: textColor,
                                font: {
                                    family: fontFamily,
                                    size: fontSize
                                },
                                generateLabels: function(chart) {
                                    return [
                                        {
                                            text: 'Backup',
                                            fillStyle: '#36A2EB',
                                            hidden: false,
                                            index: 0,
                                            // Force the text color to match the theme
                                            fontColor: textColor
                                        },
                                        {
                                            text: 'Replica',
                                            fillStyle: '#FF6384',
                                            hidden: false,
                                            index: 1,
                                            // Force the text color to match the theme
                                            fontColor: textColor
                                        }
                                    ];
                                },
                                usePointStyle: true,
                                boxWidth: 10
                            },
                        }
                    }
                }
            });
        } else {
            console.warn("No Backup Jobs data available");
        }

        // Check for Scale-Out Repositories data
// Similarly, update the ScaleOutRepositories chart code
        if (data.ScaleOutRepositories && data.ScaleOutRepositories.length) {
            const scaleOutReposCtx = document.getElementById('scaleOutReposChart');
            if (!scaleOutReposCtx) {
                console.error("Canvas element 'scaleOutReposChart' not found.");
                return;
            }
            const scaleOutReposCtx2d = scaleOutReposCtx.getContext('2d');
            const scaleOutReposData = {
                labels: data.ScaleOutRepositories.map(repo => repo.name || 'Unnamed Repository'),
                datasets: [{
                    label: 'Performance Tier',
                    data: data.ScaleOutRepositories.map(repo => {
                        if (repo.performanceTier && Array.isArray(repo.performanceTier.performanceExtents)) {
                            return repo.performanceTier.performanceExtents.length;
                        } else if (repo.extentIds) {
                            // Fallback for test data
                            return Array.isArray(repo.extentIds) ? repo.extentIds.length : 0;
                        }
                        return 0;
                    }),
                    backgroundColor: '#FF6384',
                    borderWidth: 1
                }, {
                    label: 'Capacity Tier',
                    data: data.ScaleOutRepositories.map(repo => {
                        if (repo.capacityTier) {
                            if (Array.isArray(repo.capacityTier.extents)) {
                                return repo.capacityTier.extents.length;
                            } else if (repo.capacityTier.extentId) {
                                return 1;
                            } else if (repo.capacityTier.isEnabled) {
                                return 1;
                            }
                        }
                        return 0;
                    }),
                    backgroundColor: '#FFCE56',
                    borderWidth: 1
                }, {
                    label: 'Archive Tier',
                    data: data.ScaleOutRepositories.map(repo => 
                        (repo.archiveTier && (repo.archiveTier.isEnabled || repo.archiveTier.extentId)) ? 1 : 0
                    ),
                    backgroundColor: '#4BC0C0',
                    borderWidth: 1
                }]
            };
            console.log("Scale-Out Repositories Data:", scaleOutReposData);

            // Calculate the maximum value for the x-axis
            const maxValue = Math.max(
                ...scaleOutReposData.datasets.map(dataset => Math.max(...dataset.data))
            );

            new Chart(scaleOutReposCtx2d, {
                type: 'bar',
                data: scaleOutReposData,
                options: {
                    indexAxis: 'y', // Use horizontal bar chart
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        y: {
                            beginAtZero: true,
                            ticks: {
                                color: textColor,
                                font: {
                                    family: fontFamily,
                                    size: fontSize
                                },
                                autoSkip: false // Ensure all labels are shown
                            }
                        },
                        x: {
                            display: false, // Remove the numbers on the bottom axis
                            max: maxValue + 0.1 // Set the maximum value for the x-axis dynamically with a buffer
                        }
                    },
                    layout: {
                        padding: {
                            right: 20 // Add padding to the right side of the chart
                        }
                    },
                    plugins: {
                        title: {
                            display: true,
                            text: 'Scale-Out Repositories',
                            color: textColor,
                            font: {
                                family: fontFamily,
                                size: titleFontSize
                            }
                        },
                        legend: {
                            position: 'bottom',
                            labels: {
                                color: textColor,
                                font: {
                                    family: fontFamily,
                                    size: fontSize
                                },
                                usePointStyle: true,
                                boxWidth: 10
                            }
                        },
                        tooltip: {
                            displayColors: false, // Remove color box
                            callbacks: {
                                label: function(context) {
                                    const repo = data.ScaleOutRepositories[context.dataIndex];
                                    const repositories = data.Repositories; // Assuming repositories data is available in data.Repositories

                                    const performanceTier = repo.performanceTier && repo.performanceTier.performanceExtents ? 'Enabled' : 'Disabled';
                                    const capacityTier = repo.capacityTier && repo.capacityTier.isEnabled ? 'Enabled' : 'Disabled';
                                    const archiveTier = repo.archiveTier && repo.archiveTier.isEnabled ? 'Enabled' : 'Disabled';
                                    const copyPolicy = repo.capacityTier && repo.capacityTier.copyPolicyEnabled ? 'Enabled' : 'Disabled';
                                    const movePolicy = repo.capacityTier && repo.capacityTier.movePolicyEnabled ? 'Enabled' : 'Disabled';
                                    const operationalRestorePeriodDays = repo.capacityTier && repo.capacityTier.operationalRestorePeriodDays ? repo.capacityTier.operationalRestorePeriodDays : 'N/A';
                                    const archivePeriodDays = repo.archiveTier && repo.archiveTier.archivePeriodDays ? repo.archiveTier.archivePeriodDays : 'N/A';

                                    const performanceExtents = repo.performanceTier && repo.performanceTier.performanceExtents ? repo.performanceTier.performanceExtents.map(extent => {
                                        const repository = repositories.find(r => r.id === extent.id);
                                        return repository ? repository.name : extent.id;
                                    }).join(', ') : 'N/A';

                                    const capacityExtents = repo.capacityTier && repo.capacityTier.extents ? repo.capacityTier.extents.map(extent => {
                                        const repository = repositories.find(r => r.id === extent.id);
                                        return repository ? repository.name : extent.id;
                                    }).join(', ') : 'N/A';

                                    const archiveExtents = repo.archiveTier && repo.archiveTier.extentId ? (() => {
                                        const repository = repositories.find(r => r.id === repo.archiveTier.extentId);
                                        return repository ? repository.name : repo.archiveTier.extentId;
                                    })() : 'N/A';

                                    let label = `${context.dataset.label}: ${context.raw}\n`;
                                    if (context.dataset.label === 'Performance Tier') {
                                        label += `Performance Tier: ${performanceTier}\nPerformance Extents: ${performanceExtents}`;
                                    } else if (context.dataset.label === 'Capacity Tier') {
                                        label += `Capacity Tier: ${capacityTier}\nCapacity Extents: ${capacityExtents}\nOperational Restore Period Days: ${operationalRestorePeriodDays}\nCopy Policy: ${copyPolicy}\nMove Policy: ${movePolicy}`;
                                    } else if (context.dataset.label === 'Archive Tier') {
                                        label += `Archive Tier: ${archiveTier}\nArchive Extents: ${archiveExtents}\nArchive Period Days: ${archivePeriodDays}`;
                                    }
                                    return label.split('\n'); // Add line breaks
                                }
                            },
                            bodyFont: {
                                family: fontFamily,
                                size: fontSize
                            },
                            boxWidth: 0 // Remove the box width to allow more space for text
                        }
                    }
                }
            });
        } else {
            console.warn("No Scale-Out Repositories data available");
        }

        // Create a polar chart for Credential Types
        // Update the chart for credential types to be more resilient
        if (data.Credentials || data.CloudCredentials) {
            const credentialsCtx = document.getElementById('credentialsChart');
            if (!credentialsCtx) {
                console.error("Canvas element 'credentialsChart' not found.");
                return;
            }
            const credentialsCtx2d = credentialsCtx.getContext('2d');
        
                        // Initialize credential types object
            const credentialTypes = {};
            
            // Safely add credentials if they exist
            if (Array.isArray(data.Credentials)) {
                data.Credentials.forEach(cred => {
                    if (cred && cred.type && cred.type !== 'Standard') {
                        credentialTypes[cred.type] = (credentialTypes[cred.type] || 0) + 1;
                    }
                });
            }
            
            // Safely add cloud credentials if they exist
            if (Array.isArray(data.CloudCredentials)) {
                data.CloudCredentials.forEach(cred => {
                    if (cred && cred.type && cred.type !== 'Standard') {
                        credentialTypes[cred.type] = (credentialTypes[cred.type] || 0) + 1;
                    }
                });
            }

            // Create the data structure for the chart
            const credentialLabels = Object.keys(credentialTypes);
            const credentialData = Object.values(credentialTypes);

            const credentialsData = {
                labels: credentialLabels,
                datasets: [{
                    label: 'Credential Types',
                    data: credentialData,
                    backgroundColor: ['#36A2EB', '#FF6384', '#FFCE56', '#4BC0C0', '#9966FF', '#FF9F40'],
                    borderWidth: 1
                }]
            };
            
            console.log("Credentials Data:", credentialsData);
            new Chart(credentialsCtx2d, {
                type: 'polarArea',
                data: credentialsData,
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        r: {
                            ticks: {
                                color: textColor,
                                font: {
                                    family: fontFamily,
                                    size: fontSize
                                }
                            }
                        }
                    },
                    plugins: {
                        title: {
                            display: true,
                            text: 'Credential Types',
                            color: textColor,
                            font: {
                                family: fontFamily,
                                size: titleFontSize
                            }
                        },
                        legend: {
                            position: 'bottom',
                            labels: {
                                color: textColor,
                                font: {
                                    family: fontFamily,
                                    size: fontSize
                                },
                                usePointStyle: true,
                                boxWidth: 10
                            }
                        }
                    }
                }
            });
        } else {
            console.warn("No Credentials or Cloud Credentials data available");
        }

        console.log("Charts created successfully.");
    } catch (error) {
        console.error("Error generating charts:", error);
    }
}