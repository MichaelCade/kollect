// cost.js - Cloud Cost Explorer

registerDataHandler('cost', 
    function(data) {
        return data.costs || data.aws || data.azure || data.gcp;
    },
    function(data) {
        console.log("Processing cost data:", data);
        logAllResourceTypes(data);
        
        // Add disclaimer
        const disclaimer = data.disclaimer || "Cost estimates are approximations based on publicly available pricing information. Actual costs may vary based on your specific agreements, reserved capacity, and other factors.";
        
        const disclaimerDiv = document.createElement('div');
        disclaimerDiv.className = 'cost-disclaimer';
        disclaimerDiv.innerHTML = `
            <div style="background-color: rgba(255, 165, 0, 0.1); border-left: 4px solid #FFA500; padding: 10px; margin-bottom: 20px;">
                <p style="margin: 0; font-size: 0.9em; color: var(--text-color);"><i class="fas fa-info-circle"></i> <strong>Disclaimer:</strong> ${disclaimer}</p>
            </div>
        `;
        document.getElementById('content').prepend(disclaimerDiv);
        
        // Check if we have costs data wrapper or direct data
        let costsData = data.costs || data;
        
        // Process cost data by platform
        if (costsData.aws) {
            processPlatformCosts('AWS', costsData.aws);
        }
        
        if (costsData.azure) {
            processPlatformCosts('Azure', costsData.azure);
        }
        
        if (costsData.gcp) {
            processPlatformCosts('GCP', costsData.gcp);
        }
        
        // If there's a global summary, create a summary section
        if (costsData.GlobalSummary) {
            createGlobalSummary(costsData.GlobalSummary);
        }
        
        // Create cost charts
        createCostCharts(costsData);
    }
);

function debugCostData(data) {
    console.log("--- COST DATA DEBUG ---");
    console.log("Raw data:", data);
    console.log("Data structure:", JSON.stringify(data, null, 2));
    
    // Check different possible structures
    if (data.aws) {
        console.log("Found AWS data at data.aws");
    } else if (data.costs && data.costs.aws) {
        console.log("Found AWS data at data.costs.aws");
    } else {
        console.log("No AWS data found");
    }
    
    if (data.azure) {
        console.log("Azure data found:");
        console.log("Disk Snapshots:", data.azure.DiskSnapshotCosts ? data.azure.DiskSnapshotCosts.length : "none");
        console.log("Azure Summary:", data.azure.Summary);
    } else {
        console.log("No Azure data found");
    }
    
    if (data.gcp) {
        console.log("GCP data found:");
        console.log("Disk Snapshots:", data.gcp.DiskSnapshotCosts ? data.gcp.DiskSnapshotCosts.length : "none");
        console.log("GCP Summary:", data.gcp.Summary);
    } else {
        console.log("No GCP data found");
    }
    
    console.log("--- END DEBUG ---");
}

function fetchCostData(platform) {
    showLoadingIndicator();
    
    console.log(`Fetching cost data for platform: ${platform}`);
    
    const useMock = document.getElementById('use-mock-data')?.checked || false;
    // Add type=all to get all resource types, not just snapshots
    const url = `/api/costs?platform=${platform}&type=all${useMock ? '&mock=true' : ''}`;
    
    console.log(`API URL: ${url}`);
    fetch(url)
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => {
                    throw new Error(`HTTP error ${response.status}: ${text}`);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("Cost data received:", data);
            console.log("Cost data keys:", Object.keys(data));
            console.log("Data costs:", data.costs);
            if (data.costs) {
                console.log("Costs keys:", Object.keys(data.costs));
            }
            document.getElementById('content').innerHTML = '';
            
            debugCostData(data.costs || data);
            
            processWithHandler({
                type: 'cost',
                ...data
            });
            
            updateResourceNav();
        })
        .catch(error => {
            console.error("Error fetching cost data:", error);
            document.getElementById('content').innerHTML = `
                <div class="error-message">
                    <h2>Error Collecting Cost Data</h2>
                    <p>${error.message}</p>
                </div>
            `;
        })
        .finally(() => {
            hideLoadingIndicator();
        });
}

function processPlatformCosts(platform, costData) {
    console.log(`Processing ${platform} cost data:`, costData);
    
    // Check if we have the expected cost data
    if (platform === 'AWS') {
        console.log("AWS cost data details:", {
            hasEBSSnapshots: costData.EBSSnapshotCosts ? costData.EBSSnapshotCosts.length : 'none',
            hasRDSSnapshots: costData.RDSSnapshotCosts ? costData.RDSSnapshotCosts.length : 'none',
            hasEC2Instances: costData.EC2Costs ? costData.EC2Costs.length : 'none',
            summary: costData.Summary
        });
    }

    // Check if we received just a Summary or a message
    if (costData.Message) {
        const messageDiv = document.createElement('div');
        messageDiv.className = 'cost-message';
        messageDiv.innerHTML = `
            <div style="background-color: rgba(255, 165, 0, 0.1); border-left: 4px solid #FFA500; padding: 15px; margin-bottom: 20px; border-radius: 4px;">
                <h3 style="margin-top: 0; margin-bottom: 10px;">${platform} Cost Analysis</h3>
                <p style="margin: 0;">${costData.Message}</p>
            </div>
        `;
        document.getElementById('content').appendChild(messageDiv);
        return;
    }
    
    // Display EBS Snapshot costs for AWS
    if (platform === 'AWS') {
        if (costData.EBSSnapshotCosts && costData.EBSSnapshotCosts.length > 0) {
            createTable(`${platform} EBS Snapshot Costs`, costData.EBSSnapshotCosts, 
                item => `<td>${item.SnapshotId}</td><td>${item.VolumeId || 'N/A'}</td><td>${item.SizeGB} GB</td><td>${item.Region}</td><td>$${item.PricePerGBMonth.toFixed(3)}</td><td>${item.MonthlyCostUSD}</td>`,
                ['Snapshot ID', 'Volume ID', 'Size', 'Region', 'Price per GB/Month', 'Monthly Cost']);
        } else {
            console.log("No EBS snapshot costs found");
        }
        
        if (costData.RDSSnapshotCosts && costData.RDSSnapshotCosts.length > 0) {
            createTable(`${platform} RDS Snapshot Costs`, costData.RDSSnapshotCosts, 
                item => `<td>${item.SnapshotId}</td><td>${item.Engine || 'N/A'}</td><td>${item.SizeGB} GB</td><td>${item.Region}</td><td>$${item.PricePerGBMonth.toFixed(3)}</td><td>${item.MonthlyCostUSD}</td>`,
                ['Snapshot ID', 'Engine', 'Size', 'Region', 'Price per GB/Month', 'Monthly Cost']);
        } else {
            console.log("No RDS snapshot costs found");
        }
        
        // NEW: Display EC2 Instance costs
        if (costData.EC2Costs && costData.EC2Costs.length > 0) {
            createTable(`${platform} EC2 Instance Costs`, costData.EC2Costs, 
                item => `<td>${item.InstanceId}</td><td>${item.InstanceType}</td><td>${item.Region}</td><td>$${item.HourlyCost.toFixed(4)}</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Instance ID', 'Instance Type', 'Region', 'Hourly Cost', 'Monthly Cost']);
        } else {
            console.log("No EC2 instance costs found");
        }
        
        // NEW: Display S3 Bucket costs
        if (costData.S3Costs && costData.S3Costs.length > 0) {
            console.log(`Found ${costData.S3Costs.length} S3 bucket costs`);
            createTable(`${platform} S3 Bucket Costs`, costData.S3Costs, 
                item => `<td>${item.Name}</td><td>${item.SizeGB} GB</td><td>${item.Region}</td><td>${item.StorageClass}</td><td>$${item.PricePerGB.toFixed(4)}/GB</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Bucket Name', 'Size', 'Region', 'Storage Class', 'Price per GB', 'Monthly Cost']);
        }
        
        // NEW: Display RDS Instance costs
        if (costData.RDSInstanceCosts && costData.RDSInstanceCosts.length > 0) {
            createTable(`${platform} RDS Instance Costs`, costData.RDSInstanceCosts, 
                item => `<td>${item.DBInstanceIdentifier}</td><td>${item.Engine}</td><td>${item.Region}</td><td>${item.DBInstanceClass}</td><td>${item.AllocatedStorage} GB</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Instance ID', 'Engine', 'Region', 'Instance Class', 'Storage', 'Monthly Cost']);
        }
    }
    
    // Display Disk Snapshot costs for Azure
    if (platform === 'Azure') {
        if (costData.DiskSnapshotCosts && costData.DiskSnapshotCosts.length > 0) {
            createTable(`${platform} Disk Snapshot Costs`, costData.DiskSnapshotCosts, 
                item => `<td>${item.Name}</td><td>${item.SizeGB} GB</td><td>${item.Location}</td><td>${item.State || 'N/A'}</td><td>$${item.PricePerGBMonth.toFixed(3)}</td><td>${item.MonthlyCostUSD}</td>`,
                ['Name', 'Size', 'Region', 'State', 'Price per GB/Month', 'Monthly Cost']);
        }
        
        // NEW: Display VM costs
        if (costData.VMCosts && costData.VMCosts.length > 0) {
            createTable(`${platform} Virtual Machine Costs`, costData.VMCosts, 
                item => `<td>${item.Name}</td><td>${item.ResourceGroup}</td><td>${item.Location}</td><td>${item.VMSize}</td><td>$${item.HourlyCost.toFixed(4)}</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Name', 'Resource Group', 'Location', 'VM Size', 'Hourly Cost', 'Monthly Cost']);
        }
        
        // NEW: Display Storage Account costs
        if (costData.StorageAccountCosts && costData.StorageAccountCosts.length > 0) {
            createTable(`${platform} Storage Account Costs`, costData.StorageAccountCosts, 
                item => `<td>${item.Name}</td><td>${item.UsedCapacityGB} GB</td><td>${item.Location}</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Name', 'Used Capacity', 'Location', 'Monthly Cost']);
        }
    }
    
    // Display Disk Snapshot costs for GCP
    if (platform === 'GCP') {
        if (costData.DiskSnapshotCosts && costData.DiskSnapshotCosts.length > 0) {
            createTable(`${platform} Disk Snapshot Costs`, costData.DiskSnapshotCosts, 
                item => `<td>${item.Name}</td><td>${item.SizeGB} GB</td><td>${item.Region}</td><td>${item.Status || 'N/A'}</td><td>$${item.PricePerGBMonth.toFixed(3)}</td><td>${item.MonthlyCostUSD}</td>`,
                ['Name', 'Size', 'Region', 'Status', 'Price per GB/Month', 'Monthly Cost']);
        }
        
        // NEW: Display Compute Instance costs
        if (costData.ComputeCosts && costData.ComputeCosts.length > 0) {
            createTable(`${platform} Compute Instance Costs`, costData.ComputeCosts, 
                item => `<td>${item.Name}</td><td>${item.MachineType}</td><td>${item.Zone}</td><td>$${item.HourlyCost.toFixed(4)}</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Name', 'Machine Type', 'Zone', 'Hourly Cost', 'Monthly Cost']);
        }
        
        // NEW: Display GCS Bucket costs
        if (costData.GCSCosts && costData.GCSCosts.length > 0) {
            createTable(`${platform} Cloud Storage Costs`, costData.GCSCosts, 
                item => `<td>${item.Name}</td><td>${item.SizeGB} GB</td><td>${item.Location}</td><td>${item.StorageClass}</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Bucket Name', 'Size', 'Location', 'Storage Class', 'Monthly Cost']);
        }
        
        // NEW: Display Cloud SQL costs
        if (costData.CloudSQLCosts && costData.CloudSQLCosts.length > 0) {
            createTable(`${platform} Cloud SQL Costs`, costData.CloudSQLCosts, 
                item => `<td>${item.Name}</td><td>${item.DatabaseVersion}</td><td>${item.Region}</td><td>${item.Tier}</td><td>${item.DiskSizeGB} GB</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Name', 'Version', 'Region', 'Tier', 'Disk Size', 'Monthly Cost']);
        }
    }
    
    // Display summary for this platform
    if (costData.Summary) {
        const summaryDiv = document.createElement('div');
        summaryDiv.className = 'cost-summary';
        summaryDiv.innerHTML = `
            <div class="summary-card" style="background-color: var(--card-bg); border-radius: 8px; padding: 15px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                <h3 style="margin-top: 0; border-bottom: 1px solid var(--border-color); padding-bottom: 10px;">${platform} Cost Summary</h3>
                <div style="display: flex; justify-content: space-between; flex-wrap: wrap; margin-top: 10px;">
                    <div class="summary-item" style="margin-bottom: 10px; min-width: 180px;">
                        <div style="font-size: 0.9em; color: var(--secondary-text-color);">Total Snapshot Storage</div>
                        <div style="font-size: 1.5em; font-weight: bold;">${costData.Summary.TotalSnapshotStorage.toFixed(2)} GB</div>
                    </div>
                    ${costData.Summary.TotalComputeCost ? `
                    <div class="summary-item" style="margin-bottom: 10px; min-width: 180px;">
                        <div style="font-size: 0.9em; color: var(--secondary-text-color);">Compute Resources Cost</div>
                        <div style="font-size: 1.5em; font-weight: bold;">$${costData.Summary.TotalComputeCost.toFixed(2)}</div>
                    </div>
                    ` : ''}
                    <div class="summary-item" style="margin-bottom: 10px; min-width: 180px;">
                        <div style="font-size: 0.9em; color: var(--secondary-text-color);">Estimated Total Monthly Cost</div>
                        <div style="font-size: 1.5em; font-weight: bold;">$${costData.Summary.TotalMonthlyCost.toFixed(2)}</div>
                    </div>
                </div>
            </div>
        `;
        document.getElementById('content').appendChild(summaryDiv);
    } else {
        console.log(`No summary found for ${platform}`);
    }
}

function createGlobalSummary(summary) {
    const summaryDiv = document.createElement('div');
    summaryDiv.className = 'cost-global-summary';
    summaryDiv.innerHTML = `
        <div class="summary-card" style="background-color: rgba(74, 144, 226, 0.1); border-radius: 8px; padding: 15px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
            <h3 style="margin-top: 0; border-bottom: 1px solid var(--border-color); padding-bottom: 10px; color: var(--accent-color);">Total Snapshot Cost Across All Platforms</h3>
            <div style="display: flex; justify-content: space-between; margin-top: 10px;">
                <div class="summary-item">
                    <div style="font-size: 0.9em; color: var(--secondary-text-color);">Total Snapshot Storage</div>
                    <div style="font-size: 1.8em; font-weight: bold;">${summary.TotalSnapshotStorage.toFixed(2)} GB</div>
                </div>
                <div class="summary-item">
                    <div style="font-size: 0.9em; color: var(--secondary-text-color);">Estimated Monthly Cost</div>
                    <div style="font-size: 1.8em; font-weight: bold; color: var(--accent-color);">$${summary.TotalMonthlyCost.toFixed(2)}</div>
                </div>
            </div>
        </div>
    `;
    document.getElementById('content').prepend(summaryDiv);
}

function createCostCharts(data) {
    // Create a div for the charts
    const chartsDiv = document.createElement('div');
    chartsDiv.className = 'cost-charts';
    chartsDiv.style.display = 'flex';
    chartsDiv.style.flexWrap = 'wrap';
    chartsDiv.style.justifyContent = 'space-between';
    chartsDiv.style.marginTop = '20px';
    
    // Add chart container divs
    chartsDiv.innerHTML = `
        <div class="chart-wrapper" style="width: 48%; height: 300px; margin-bottom: 20px;">
            <canvas id="storageByPlatformChart"></canvas>
        </div>
        <div class="chart-wrapper" style="width: 48%; height: 300px; margin-bottom: 20px;">
            <canvas id="costByPlatformChart"></canvas>
        </div>
    `;
    
    document.getElementById('content').appendChild(chartsDiv);
    
    // Extract data for charts
    const platforms = [];
    const storageValues = [];
    const costValues = [];
    
    if (data.aws && data.aws.Summary) {
        platforms.push('AWS');
        storageValues.push(data.aws.Summary.TotalSnapshotStorage);
        costValues.push(data.aws.Summary.TotalMonthlyCost);
    }
    
    if (data.azure && data.azure.Summary) {
        platforms.push('Azure');
        storageValues.push(data.azure.Summary.TotalSnapshotStorage);
        costValues.push(data.azure.Summary.TotalMonthlyCost);
    }
    
    if (data.gcp && data.gcp.Summary) {
        platforms.push('GCP');
        storageValues.push(data.gcp.Summary.TotalSnapshotStorage);
        costValues.push(data.gcp.Summary.TotalMonthlyCost);
    }
    
    // Colors for each platform
    const colors = {
        'AWS': '#FF9900',
        'Azure': '#0078D4',
        'GCP': '#4285F4'
    };
    
    const platformColors = platforms.map(platform => colors[platform] || '#777777');
    
    // Create storage chart
    const storageCtx = document.getElementById('storageByPlatformChart').getContext('2d');
    new Chart(storageCtx, {
        type: 'bar',
        data: {
            labels: platforms,
            datasets: [{
                label: 'Snapshot Storage (GB)',
                data: storageValues,
                backgroundColor: platformColors,
                borderColor: platformColors,
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: {
                    display: true,
                    text: 'Snapshot Storage by Platform',
                    font: {
                        size: 16
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: 'Storage (GB)'
                    }
                }
            }
        }
    });
    
    // Create cost chart
    const costCtx = document.getElementById('costByPlatformChart').getContext('2d');
    new Chart(costCtx, {
        type: 'bar',
        data: {
            labels: platforms,
            datasets: [{
                label: 'Monthly Cost (USD)',
                data: costValues,
                backgroundColor: platformColors,
                borderColor: platformColors,
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: {
                    display: true,
                    text: 'Monthly Cost by Platform',
                    font: {
                        size: 16
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: 'Cost (USD)'
                    }
                }
            }
        }
    });
}

function logAllResourceTypes(data) {
    console.log("--- RESOURCE TYPES IN COST DATA ---");
    
    // Check if we have aws/azure/gcp structure or direct structure
    const costsData = data.costs || data;
    
    if (costsData.aws) {
        console.log("AWS resource types:", Object.keys(costsData.aws).filter(k => k !== 'Summary' && k !== 'Message'));
    }
    
    if (costsData.azure) {
        console.log("Azure resource types:", Object.keys(costsData.azure).filter(k => k !== 'Summary' && k !== 'Message'));
    }
    
    if (costsData.gcp) {
        console.log("GCP resource types:", Object.keys(costsData.gcp).filter(k => k !== 'Summary' && k !== 'Message'));
    }
    
    console.log("--- END RESOURCE TYPES ---");
}

function showCostExplorerModal() {
    const modal = document.createElement('div');
    modal.className = 'modal';
    modal.style.position = 'fixed';
    modal.style.top = '0';
    modal.style.left = '0';
    modal.style.width = '100%';
    modal.style.height = '100%';
    modal.style.backgroundColor = 'rgba(0, 0, 0, 0.5)';
    modal.style.display = 'flex';
    modal.style.justifyContent = 'center';
    modal.style.alignItems = 'center';
    modal.style.zIndex = '1000';
    
    const modalContent = document.createElement('div');
    modalContent.className = 'cost-modal';
    modalContent.style.backgroundColor = 'var(--card-bg)';
    modalContent.style.padding = '20px';
    modalContent.style.borderRadius = '8px';
    modalContent.style.maxWidth = '500px';
    modalContent.style.width = '90%';
    modalContent.style.maxHeight = '90vh';
    modalContent.style.overflowY = 'auto';
    modalContent.style.position = 'relative';
    modalContent.style.boxShadow = '0 5px 20px rgba(0,0,0,0.4)';
    modalContent.style.border = '1px solid var(--border-color)';
    
    modalContent.innerHTML = `
        <h3 style="margin-top: 0; color: var(--accent-color); font-size: 1.5em; border-bottom: 1px solid var(--border-color); padding-bottom: 10px;">
            <i class="fas fa-dollar-sign"></i> Cost Explorer
        </h3>
        
        <p style="margin-bottom: 20px;">Select which cloud platform to analyze costs for all resources:</p>
        
        <div class="cost-platform-selection" style="margin: 20px 0;">
            <div style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 10px;">
                <div style="display: flex; align-items: center;">
                    <input type="radio" id="cost-aws" name="cost-platform" value="aws" style="margin-right: 10px;">
                    <label for="cost-aws" style="font-weight: bold; font-size: 1.1em;">
                        <i class="fab fa-aws"></i> AWS
                    </label>
                </div>
                <div style="padding: 5px 0 0 25px; font-size: 0.9em; color: var(--secondary-text-color);">
                    Estimates costs for EC2, S3, RDS, snapshots and more
                </div>
            </div>
            
            <div style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 10px;">
                <div style="display: flex; align-items: center;">
                    <input type="radio" id="cost-azure" name="cost-platform" value="azure" style="margin-right: 10px;">
                    <label for="cost-azure" style="font-weight: bold; font-size: 1.1em;">
                        <i class="fab fa-microsoft"></i> Azure
                    </label>
                </div>
                <div style="padding: 5px 0 0 25px; font-size: 0.9em; color: var(--secondary-text-color);">
                    Estimates costs for VMs, Storage Accounts, SQL Databases, snapshots and more
                </div>
            </div>
            
            <div style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 10px;">
                <div style="display: flex; align-items: center;">
                    <input type="radio" id="cost-gcp" name="cost-platform" value="gcp" style="margin-right: 10px;">
                    <label for="cost-gcp" style="font-weight: bold; font-size: 1.1em;">
                        <i class="fab fa-google"></i> GCP
                    </label>
                </div>
                <div style="padding: 5px 0 0 25px; font-size: 0.9em; color: var(--secondary-text-color);">
                    Estimates costs for Compute Instances, GCS Buckets, Cloud SQL, snapshots and more
                </div>
            </div>
            
            <div style="background: rgba(74, 144, 226, 0.1); border-radius: 6px; padding: 12px; margin-top: 15px; margin-bottom: 10px;">
                <div style="display: flex; align-items: center;">
                    <input type="radio" id="cost-all" name="cost-platform" value="all" style="margin-right: 10px;" checked>
                    <label for="cost-all" style="font-weight: bold; font-size: 1.1em; color: var(--accent-color);">
                        <i class="fas fa-globe"></i> All Connected Platforms
                    </label>
                </div>
                <div style="padding: 5px 0 0 25px; font-size: 0.9em; color: var(--secondary-text-color);">
                    Estimates costs for all cloud resources across all connected platforms
                </div>
            </div>
        </div>

        <div style="margin-top: 15px; margin-bottom: 15px;">
            <label style="display: flex; align-items: center; cursor: pointer;">
                <input type="checkbox" id="use-mock-data" style="margin-right: 10px;">
                <span>Use sample data for testing (no real cloud access needed)</span>
            </label>
        </div>
        
        <div class="modal-buttons" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 25px; border-top: 1px solid var(--border-color); padding-top: 15px;">
            <button id="cost-cancel-btn" class="btn" style="padding: 10px 20px; background-color: var(--button-bg-color); color: var(--button-text-color); border: none; border-radius: 4px; cursor: pointer;">Cancel</button>
            <button id="cost-analyze-btn" class="btn btn-primary" style="padding: 10px 20px; background-color: var(--accent-color); color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">
                <i class="fas fa-calculator"></i> Analyze Costs
            </button>
        </div>
    `;
    
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    document.getElementById('cost-cancel-btn').addEventListener('click', () => {
        console.log("Cancel button clicked");
        modal.remove();
    });
    
    document.getElementById('cost-analyze-btn').addEventListener('click', () => {
        const selectedPlatform = document.querySelector('input[name="cost-platform"]:checked').value;
        const useMock = document.getElementById('use-mock-data').checked;
        
        console.log(`Analyzing costs for platform: ${selectedPlatform}, mock: ${useMock}`);
        fetchCostData(selectedPlatform);
        modal.remove();
    });
    
    checkPlatformAvailability();
    
    function checkPlatformAvailability() {
        fetch('/api/check-credentials')
            .then(response => response.json())
            .then(data => {
                console.log('Connection status for costs:', data);
                
                const awsButton = document.getElementById('cost-aws');
                if (awsButton) awsButton.disabled = !data.aws;
                
                const azureButton = document.getElementById('cost-azure');
                if (azureButton) azureButton.disabled = !data.azure;
                
                const gcpButton = document.getElementById('cost-gcp');
                if (gcpButton) gcpButton.disabled = !data.gcp;
                
                const allButton = document.getElementById('cost-all');
                if (allButton) {
                    const availablePlatforms = [data.aws, data.azure, data.gcp].filter(Boolean).length;
                    allButton.disabled = availablePlatforms === 0;
                
                    if (availablePlatforms === 0) {
                        const modalContent = document.querySelector('.cost-modal');
                        if (modalContent) {
                            const messageDiv = document.createElement('div');
                            messageDiv.style.backgroundColor = 'rgba(255, 165, 0, 0.1)';
                            messageDiv.style.borderLeft = '4px solid #FFA500';
                            messageDiv.style.padding = '10px';
                            messageDiv.style.marginBottom = '20px';
                            
                            messageDiv.innerHTML = `
                                <p style="margin: 0; font-size: 0.9em; color: var(--text-color);">
                                    <i class="fas fa-exclamation-triangle"></i> <strong>No connected platforms detected.</strong> 
                                    <br>Please connect to at least one cloud platform (AWS, Azure, or GCP) before using Cost Explorer.
                                </p>
                            `;
                            
                            modalContent.insertBefore(messageDiv, modalContent.children[1]);
                        }
                    }
                }
            })
            .catch(error => {
                console.error('Error checking platform connections:', error);
            });
    }
}

// Set up event listener for the Cost Explorer button
document.addEventListener('DOMContentLoaded', function() {
    console.log("DOM loaded - Cost Explorer module setting up event listener");
    
    const costButton = document.getElementById('cost-button');
    if (costButton) {
        console.log("Found cost button, setting up handler");
        
        costButton.addEventListener('click', function(event) {
            console.log("Cost button clicked");
            event.preventDefault();
            showCostExplorerModal();
        });
    } else {
        console.error("Could not find cost-button element - you may need to add it to index.html");
    }
});