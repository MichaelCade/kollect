// veeam.js

registerDataHandler('veeam', 
    function(data) {
        return data.ServerInfo || data.BackupJobs || data.Repositories || 
               data.Credentials || data.CloudCredentials;
    },
    function(data) {
        console.log("Processing Veeam data");
        
        const chartsContainer = document.getElementById('charts-container');
        if (chartsContainer) {
            chartsContainer.style.display = 'grid';
        }
        
        if (data.ServerInfo) {
            createTable('Server Info', [data.ServerInfo], serverInfoRowTemplate, 
                ['Name', 'Build Version', 'Database Vendor', 'SQL Server Version', 'VBR ID']);
        }
        
        if (data.Credentials) {
            createTable('Credentials', data.Credentials, credentialsRowTemplate, 
                ['Username', 'Description', 'Type']);
        }
        
        if (data.CloudCredentials) {
            createTable('Cloud Credentials', data.CloudCredentials, cloudCredentialsRowTemplate, 
                ['Account', 'Type', 'Description']);
        }
        
        if (data.KMSServers) {
            createTable('KMS Servers', data.KMSServers, kmsServersRowTemplate, 
                ['ID', 'Name']);
        }
        
        if (data.ManagedServers) {
            createTable('Managed Servers', data.ManagedServers, managedServersRowTemplate, 
                ['Name', 'Type', 'Status', 'Description']);
        }
        
        if (data.Repositories) {
            createTable('Repositories', data.Repositories, repositoriesRowTemplate, 
                ['Name', 'Type', 'Description', 'Bucket Name', 'Folder Name', 'Region ID', 'Infrequent Access Storage', 'Immutability Status', 'Immutable Period']);
        }
        
        if (data.ScaleOutRepositories) {
            const tableId = 'Scale-Out-Repositories'.replace(/\s+/g, '-').toLowerCase();
            
            const tableContainer = document.createElement('div');
            tableContainer.className = 'collapsible-table';
            tableContainer.id = `table-container-${tableId}`;
            
            const tableHeader = document.createElement('div');
            tableHeader.className = 'table-header collapsed';
            tableHeader.innerHTML = `
                <span>Scale-Out Repositories</span>
                <div>
                    <span class="table-counter">${data.ScaleOutRepositories.length}</span>
                    <span class="icon">▼</span>
                </div>
            `;
            
            const tableContent = document.createElement('div');
            tableContent.className = 'table-content collapsed';
            
            const table = document.createElement('table');
            
            const thead = document.createElement('thead');
            const headerRow = document.createElement('tr');
            const headers = ['Name', 'Description', 'Details'];
            
            headers.forEach(header => {
                const th = document.createElement('th');
                th.textContent = header;
                headerRow.appendChild(th);
            });
            
            thead.appendChild(headerRow);
            table.appendChild(thead);
            
            const tbody = document.createElement('tbody');
            data.ScaleOutRepositories.forEach(item => {
                const row = document.createElement('tr');
                row.innerHTML = scaleOutRepositoriesRowTemplate(item, data.Repositories);
                tbody.appendChild(row);
            });
            
            table.appendChild(tbody);
            
            tableHeader.addEventListener('click', function() {
                tableHeader.classList.toggle('collapsed');
                tableContent.classList.toggle('collapsed');
            });
            
            tableContent.appendChild(table);
            tableContainer.appendChild(tableHeader);
            tableContainer.appendChild(tableContent);
            
            document.getElementById('content').appendChild(tableContainer);
        }
        
        if (data.Proxies) {
            createTable('Proxies', data.Proxies, proxiesRowTemplate, 
                ['Name', 'Type', 'Description', 'Max Task Count', 'Transport Mode']);
        }
        
        if (data.BackupJobs) {
            createTable('Backup Jobs', data.BackupJobs, backupJobsRowTemplate, 
                ['Job Name', 'ID', 'Description', 'Type', 'Is Disabled', 'Is High Priority', 'Job Details']);
        }
        
        setTimeout(() => {
            console.log(`Created Veeam tables`);
            generateCharts(data);
        }, 100);
    }
);


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
            <div id="details-${item.id}" style="display:none;" class="details-panel">
                <h4>Performance Tier: ${performanceTier}</h4>
                <ul>${performanceExtents}</ul>
                
                <h4>Capacity Tier: ${capacityTier}</h4>
                <ul>${capacityExtents}</ul>
                <p>Operational Restore Period Days: ${operationalRestorePeriodDays}</p>
                
                <h4>Archive Tier: ${archiveTier}</h4>
                <ul>${archiveExtents}</ul>
                <p>Archive Period Days: ${archivePeriodDays}</p>
                
                <p>Copy Policy: ${copyPolicy}</p>
                <p>Move Policy: ${movePolicy}</p>
            </div>
        </td>
    `;
}

function toggleDetails(id) {
    const details = document.getElementById(`details-${id}`);
    if (details) {
        details.style.display = details.style.display === 'none' ? 'block' : 'none';
    }
}

function proxiesRowTemplate(item) {
    return `<td>${item.name}</td><td>${item.type}</td><td>${item.description}</td><td>${item.server && item.server.maxTaskCount || 'N/A'}</td><td>${item.server && item.server.transportMode || 'N/A'}</td>`;
}

function backupJobsRowTemplate(item) {
    let vms = '';
    if (item.virtualMachines && Array.isArray(item.virtualMachines.includes)) {
        vms = item.virtualMachines.includes.map(vm => 
            `<li>Name: ${vm.name || 'Unknown'}, Host: ${vm.hostName || 'Unknown'}, Size: ${vm.size || 'Unknown'}</li>`
        ).join('');
    } else if (Array.isArray(item.sourceObjects)) {
        vms = item.sourceObjects.map(obj => 
            `<li>Name: ${obj.name || 'Unknown'}, Type: ${obj.type || 'Unknown'}</li>`
        ).join('');
    }
    
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
            <div id="details-${item.id}" style="display:none;" class="details-panel">
                <h4>Included VMs</h4>
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
        const currentTheme = document.documentElement.getAttribute('data-theme') || 'dark';
        
        const textColor = getComputedStyle(document.documentElement).getPropertyValue('--text-color').trim() || 
                         (currentTheme === 'dark' ? '#ffffff' : '#333333');
        
        const fontFamily = getComputedStyle(document.documentElement).getPropertyValue('--chart-font-family').trim() || 'Arial, sans-serif';
        const fontSize = parseInt(getComputedStyle(document.documentElement).getPropertyValue('--chart-font-size').trim(), 10) || 12;
        const titleFontSize = parseInt(getComputedStyle(document.documentElement).getPropertyValue('--chart-title-font-size').trim(), 10) || 16;

        Chart.defaults.color = textColor;
        Chart.defaults.font.family = fontFamily;
        Chart.defaults.font.size = fontSize;
        
        Chart.defaults.plugins.legend.labels.color = textColor;

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
                        if (job.virtualMachines && job.virtualMachines.includes) {
                            return job.virtualMachines.includes.length;
                        } else if (job.sourceObjects) {
                            return job.sourceObjects.length;
                        } else {
                            return 0; 
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
                    indexAxis: 'y', 
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
                                autoSkip: false 
                            }
                        },
                        x: {
                            display: false
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
                                            fontColor: textColor
                                        },
                                        {
                                            text: 'Replica',
                                            fillStyle: '#FF6384',
                                            hidden: false,
                                            index: 1,
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

            const maxValue = Math.max(
                ...scaleOutReposData.datasets.map(dataset => Math.max(...dataset.data))
            );

            new Chart(scaleOutReposCtx2d, {
                type: 'bar',
                data: scaleOutReposData,
                options: {
                    indexAxis: 'y',
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
                                autoSkip: false 
                            }
                        },
                        x: {
                            display: false,
                            max: maxValue + 0.1 
                        }
                    },
                    layout: {
                        padding: {
                            right: 20 
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
                            displayColors: false, 
                            callbacks: {
                                label: function(context) {
                                    const repo = data.ScaleOutRepositories[context.dataIndex];
                                    const repositories = data.Repositories; 

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
                                    return label.split('\n'); 
                                }
                            },
                            bodyFont: {
                                family: fontFamily,
                                size: fontSize
                            },
                            boxWidth: 0 
                        }
                    }
                }
            });
        } else {
            console.warn("No Scale-Out Repositories data available");
        }

        if (data.Credentials || data.CloudCredentials) {
            const credentialsCtx = document.getElementById('credentialsChart');
            if (!credentialsCtx) {
                console.error("Canvas element 'credentialsChart' not found.");
                return;
            }
            const credentialsCtx2d = credentialsCtx.getContext('2d');
            const credentialTypes = {};
            
            if (Array.isArray(data.Credentials)) {
                data.Credentials.forEach(cred => {
                    if (cred && cred.type && cred.type !== 'Standard') {
                        credentialTypes[cred.type] = (credentialTypes[cred.type] || 0) + 1;
                    }
                });
            }
            
            if (Array.isArray(data.CloudCredentials)) {
                data.CloudCredentials.forEach(cred => {
                    if (cred && cred.type && cred.type !== 'Standard') {
                        credentialTypes[cred.type] = (credentialTypes[cred.type] || 0) + 1;
                    }
                });
            }

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

document.getElementById('veeam-button')?.addEventListener('click', () => {
    const button = document.getElementById('veeam-button');
    if (button && button.classList.contains('connected')) {
        showLoadingIndicator();
        fetch('/api/switch?type=veeam')
            .then(response => response.json())
            .then(data => {
                location.reload();
            })
            .catch(error => console.error('Error switching to Veeam:', error))
            .finally(() => hideLoadingIndicator());
    } else {
        showVeeamConnectionModal();
    }
});

function showVeeamConnectionModal() {
    const modal = document.createElement('div');
    modal.className = 'modal';
    modal.style.display = 'flex';
    modal.style.position = 'fixed';
    modal.style.zIndex = 1000;
    modal.style.left = 0;
    modal.style.top = 0;
    modal.style.width = '100%';
    modal.style.height = '100%';
    modal.style.backgroundColor = 'rgba(0,0,0,0.7)'; 
    modal.style.alignItems = 'center';
    modal.style.justifyContent = 'center';
    
    const modalContent = document.createElement('div');
    modalContent.className = 'modal-content veeam-modal';
    modalContent.style.backgroundColor = 'var(--card-bg)';
    modalContent.style.color = 'var(--text-color)';
    modalContent.style.padding = '25px';
    modalContent.style.borderRadius = '8px';
    modalContent.style.maxWidth = '450px';
    modalContent.style.width = '90%';
    modalContent.style.boxShadow = '0 5px 20px rgba(0,0,0,0.4)';
    modalContent.style.border = '1px solid var(--border-color)';
    
    modalContent.innerHTML = `
        <h3 style="margin-top: 0; color: var(--accent-color); font-size: 1.5em; border-bottom: 1px solid var(--border-color); padding-bottom: 10px;">
            <i class="fas fa-server"></i> Connect to Veeam Backup & Replication
        </h3>
        
        <div class="veeam-connection-form" style="margin-top: 20px;">
            <div class="form-group">
                <label for="veeam-server" style="font-weight: bold; margin-bottom: 5px;">Server IP/Hostname:</label>
                <div style="display: flex; align-items: center;">
                    <span style="background: var(--secondary-bg-color); padding: 8px; border-radius: 4px 0 0 4px; border: 1px solid var(--border-color); border-right: none;">https://</span>
                    <input type="text" id="veeam-server" placeholder="vbr-server.example.com" style="flex-grow: 1; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 0;">
                    <span style="background: var(--secondary-bg-color); padding: 8px; border-radius: 0 4px 4px 0; border: 1px solid var(--border-color); border-left: none;">:9419</span>
                </div>
                <small style="color: var(--text-color); opacity: 0.7; font-size: 0.8em; margin-top: 3px;">Default port is 9419 for Veeam REST API</small>
            </div>
            
            <div class="form-group" style="margin-top: 15px;">
                <label for="veeam-username" style="font-weight: bold; margin-bottom: 5px;">Username:</label>
                <input type="text" id="veeam-username" placeholder="Username with admin rights" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;">
            </div>
            
            <div class="form-group" style="margin-top: 15px;">
                <label for="veeam-password" style="font-weight: bold; margin-bottom: 5px;">Password:</label>
                <input type="password" id="veeam-password" placeholder="Enter password" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;">
            </div>
            
            <div class="form-group" style="margin-top: 10px;">
                <div style="display: flex; align-items: center;">
                    <input type="checkbox" id="veeam-ignore-ssl" checked style="margin-right: 8px;">
                    <label for="veeam-ignore-ssl">Ignore SSL certificate errors</label>
                </div>
            </div>
            
            <div class="modal-buttons" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 25px; border-top: 1px solid var(--border-color); padding-top: 15px;">
                <button id="veeam-cancel-btn" class="btn" style="padding: 10px 20px; background-color: var(--button-bg-color); color: var(--button-text-color); border: none; border-radius: 4px; cursor: pointer;">Cancel</button>
                <button id="veeam-connect-btn" class="btn btn-primary" style="padding: 10px 20px; background-color: var(--accent-color); color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">
                    <i class="fas fa-plug"></i> Connect
                </button>
            </div>
        </div>
    `;
    
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    document.getElementById('veeam-cancel-btn').addEventListener('click', () => {
        modal.remove();
    });
    
document.getElementById('veeam-connect-btn').addEventListener('click', () => {
    const server = document.getElementById('veeam-server').value.trim();
    const username = document.getElementById('veeam-username').value.trim();
    const password = document.getElementById('veeam-password').value;
    const ignoreSSL = document.getElementById('veeam-ignore-ssl').checked;
    
    if (!server || !username || !password) {
        alert('Please provide all required fields');
        return;
    }
    
    showLoadingIndicator();
    
    let serverUrl = server;  

    if (!serverUrl.startsWith('http')) {
        serverUrl = `https://${serverUrl}`;
    }
    
    if (!serverUrl.includes(':9419') && !server.includes(':')) {
        serverUrl += ':9419';
    }
    
    console.log("Connecting to Veeam server:", serverUrl);
    
    fetch('/api/veeam/connect', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            baseUrl: serverUrl,
            username: username,
            password: password,
            ignoreSSL: ignoreSSL
        })
    })
    .then(response => {
        if (!response.ok) {
            return response.text().then(text => {
                throw new Error(`HTTP error ${response.status}: ${text}`);
            });
        }
        return response.json();
    })
    .then(data => {
    if (data.status === 'success') {
        const button = document.getElementById('veeam-button');
        if (button) {
            const existingBadges = button.querySelectorAll('.connection-badge');
            existingBadges.forEach(badge => badge.remove());
            
            button.classList.add('connected');
            button.classList.remove('not-connected');
            
            const badge = document.createElement('span');
            badge.className = 'connection-badge connected';
            button.appendChild(badge);
            
            button.title = 'Veeam (Connected)';
        }
        
        modal.remove();
        
        setTimeout(() => {
            location.reload();
        }, 300);
    } else {
        throw new Error(data.message || 'Failed to connect to Veeam server');
    }
})
    .catch(error => {
        console.error("Veeam connection error:", error);
        alert(`Error connecting to Veeam server: ${error.message}`);
    })
    .finally(() => {
        hideLoadingIndicator();
    });
});
}