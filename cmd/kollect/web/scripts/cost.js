// cost.js - Cloud Cost Explorer

registerDataHandler('cost', 
    function(data) {
        return data.costs || data.aws || data.azure || data.gcp;
    },
    function(data) {
        console.log("Processing cost data:", data);
        logAllResourceTypes(data);

        window.currentCostData = data.costs || data;
        
        const disclaimer = data.disclaimer || "Cost estimates are approximations based on publicly available pricing information. Actual costs may vary based on your specific agreements, reserved capacity, and other factors.";
                
        const disclaimerDiv = document.createElement('div');
        disclaimerDiv.className = 'cost-disclaimer';
        disclaimerDiv.innerHTML = `
            <div style="background-color: rgba(255, 165, 0, 0.1); border-left: 4px solid #FFA500; padding: 10px; margin-bottom: 20px;">
                <p style="margin: 0; font-size: 0.9em; color: var(--text-color);"><i class="fas fa-info-circle"></i> <strong>Disclaimer:</strong> ${disclaimer}</p>
            </div>
        `;
        document.getElementById('content').prepend(disclaimerDiv);
        
        let costsData = data.costs || data;
        
        if (costsData.aws) {
            processPlatformCosts('AWS', costsData.aws);
        }
        
        if (costsData.azure) {
            processPlatformCosts('Azure', costsData.azure);
        }
        
        if (costsData.gcp) {
            processPlatformCosts('GCP', costsData.gcp);
        }
        
        if (costsData.GlobalSummary) {
            createGlobalSummary(costsData.GlobalSummary);
        }
        
        createCostCharts(costsData);
    }
);

function debugCostData(data) {
    console.log("--- COST DATA DEBUG ---");
    console.log("Raw data:", data);
    console.log("Data structure:", JSON.stringify(data, null, 2));
    
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
    
    if (platform === 'AWS') {
        console.log("AWS cost data details:", {
            hasEBSSnapshots: costData.EBSSnapshotCosts ? costData.EBSSnapshotCosts.length : 'none',
            hasRDSSnapshots: costData.RDSSnapshotCosts ? costData.RDSSnapshotCosts.length : 'none',
            hasEC2Instances: costData.EC2Costs ? costData.EC2Costs.length : 'none',
            summary: costData.Summary
        });
    }

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
    
    let storageCost = 0;
    let dataServicesCost = 0;
    let snapshotCost = 0;
    let computeCost = 0;
        
    if (platform === 'AWS') {
        if (costData.EBSSnapshotCosts && costData.EBSSnapshotCosts.length > 0) {
            costData.EBSSnapshotCosts.forEach(item => {
                if (item.MonthlyCost !== undefined) {
                    snapshotCost += parseFloat(item.MonthlyCost) || 0;
                } else if (item.MonthlyCostUSD !== undefined) {
                    snapshotCost += parseFloat(item.MonthlyCostUSD.replace('$', '')) || 0;
                }
            });
            
            createTable(`${platform} EBS Snapshot Costs`, costData.EBSSnapshotCosts, 
                item => `<td>${item.SnapshotId}</td><td>${item.VolumeId || 'N/A'}</td><td>${item.SizeGB} GB</td><td>${item.Region}</td><td>$${item.PricePerGBMonth.toFixed(3)}</td><td>${item.MonthlyCostUSD}</td>`,
                ['Snapshot ID', 'Volume ID', 'Size', 'Region', 'Price per GB/Month', 'Monthly Cost']);
        } else {
            console.log("No EBS snapshot costs found");
        }
        
        if (costData.RDSSnapshotCosts && costData.RDSSnapshotCosts.length > 0) {
            costData.RDSSnapshotCosts.forEach(item => {
                if (item.MonthlyCost !== undefined) {
                    snapshotCost += parseFloat(item.MonthlyCost) || 0;
                } else if (item.MonthlyCostUSD !== undefined) {
                    snapshotCost += parseFloat(item.MonthlyCostUSD.replace('$', '')) || 0;
                }
            });
            
            createTable(`${platform} RDS Snapshot Costs`, costData.RDSSnapshotCosts, 
                item => `<td>${item.SnapshotId}</td><td>${item.Engine || 'N/A'}</td><td>${item.SizeGB} GB</td><td>${item.Region}</td><td>$${item.PricePerGBMonth.toFixed(3)}</td><td>${item.MonthlyCostUSD}</td>`,
                ['Snapshot ID', 'Engine', 'Size', 'Region', 'Price per GB/Month', 'Monthly Cost']);
        } else {
            console.log("No RDS snapshot costs found");
        }
        
        if (costData.EC2Costs && costData.EC2Costs.length > 0) {
            costData.EC2Costs.forEach(item => {
                computeCost += item.MonthlyCost || 0;
            });
            
            createTable(`${platform} EC2 Instance Costs`, costData.EC2Costs, 
                item => `<td>${item.InstanceId}</td><td>${item.InstanceType}</td><td>${item.Region}</td><td>$${item.HourlyCost.toFixed(4)}</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Instance ID', 'Instance Type', 'Region', 'Hourly Cost', 'Monthly Cost']);
        } else {
            console.log("No EC2 instance costs found");
        }
        
        if (costData.S3Costs && costData.S3Costs.length > 0) {
            costData.S3Costs.forEach(item => {
                storageCost += item.MonthlyCost || 0;
            });
            
            console.log(`Found ${costData.S3Costs.length} S3 bucket costs`);
            createTable(`${platform} S3 Bucket Costs`, costData.S3Costs, 
                item => `<td>${item.Name}</td><td>${item.SizeGB} GB</td><td>${item.Region}</td><td>${item.StorageClass}</td><td>$${item.PricePerGB.toFixed(4)}/GB</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Bucket Name', 'Size', 'Region', 'Storage Class', 'Price per GB', 'Monthly Cost']);
        }
        
        if (costData.RDSInstanceCosts && costData.RDSInstanceCosts.length > 0) {
            costData.RDSInstanceCosts.forEach(item => {
                dataServicesCost += item.MonthlyCost || 0;
            });
            
            createTable(`${platform} RDS Instance Costs`, costData.RDSInstanceCosts, 
                item => `<td>${item.DBInstanceIdentifier}</td><td>${item.Engine}</td><td>${item.Region}</td><td>${item.DBInstanceClass}</td><td>${item.AllocatedStorage} GB</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Instance ID', 'Engine', 'Region', 'Instance Class', 'Storage', 'Monthly Cost']);
        }
    }
    
    if (platform === 'Azure') {
        if (costData.DiskSnapshotCosts && costData.DiskSnapshotCosts.length > 0) {
            costData.DiskSnapshotCosts.forEach(item => {
                if (item.MonthlyCost !== undefined) {
                    snapshotCost += parseFloat(item.MonthlyCost) || 0;
                } else if (item.MonthlyCostUSD !== undefined) {
                    snapshotCost += parseFloat(item.MonthlyCostUSD.replace('$', '')) || 0;
                }
            });
            
            createTable(`${platform} Disk Snapshot Costs`, costData.DiskSnapshotCosts, 
                item => `<td>${item.Name}</td><td>${item.SizeGB} GB</td><td>${item.Location}</td><td>${item.State || 'N/A'}</td><td>$${item.PricePerGBMonth.toFixed(3)}</td><td>${item.MonthlyCostUSD}</td>`,
                ['Name', 'Size', 'Region', 'State', 'Price per GB/Month', 'Monthly Cost']);
        }
        
        if (costData.VMCosts && costData.VMCosts.length > 0) {
            costData.VMCosts.forEach(item => {
                computeCost += item.MonthlyCost || 0;
            });
            
            createTable(`${platform} Virtual Machine Costs`, costData.VMCosts, 
                item => `<td>${item.Name}</td><td>${item.ResourceGroup}</td><td>${item.Location}</td><td>${item.VMSize}</td><td>$${item.HourlyCost.toFixed(4)}</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Name', 'Resource Group', 'Location', 'VM Size', 'Hourly Cost', 'Monthly Cost']);
        }
        
        if (costData.StorageAccountCosts && costData.StorageAccountCosts.length > 0) {
            costData.StorageAccountCosts.forEach(item => {
                storageCost += item.MonthlyCost || 0;
            });
            
            createTable(`${platform} Storage Account Costs`, costData.StorageAccountCosts, 
                item => `<td>${item.Name}</td><td>${item.UsedCapacityGB} GB</td><td>${item.Location}</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Name', 'Used Capacity', 'Location', 'Monthly Cost']);
        }
    }
    
    if (platform === 'GCP') {
        if (costData.DiskSnapshotCosts && costData.DiskSnapshotCosts.length > 0) {
            costData.DiskSnapshotCosts.forEach(item => {
                if (item.MonthlyCost !== undefined) {
                    snapshotCost += parseFloat(item.MonthlyCost) || 0;
                } else if (item.MonthlyCostUSD !== undefined) {
                    snapshotCost += parseFloat(item.MonthlyCostUSD.replace('$', '')) || 0;
                }
            });
            
            createTable(`${platform} Disk Snapshot Costs`, costData.DiskSnapshotCosts, 
                item => `<td>${item.Name}</td><td>${item.SizeGB} GB</td><td>${item.Region}</td><td>${item.Status || 'N/A'}</td><td>$${item.PricePerGBMonth.toFixed(3)}</td><td>${item.MonthlyCostUSD}</td>`,
                ['Name', 'Size', 'Region', 'Status', 'Price per GB/Month', 'Monthly Cost']);
        }
        
        if (costData.ComputeCosts && costData.ComputeCosts.length > 0) {
            costData.ComputeCosts.forEach(item => {
                computeCost += item.MonthlyCost || 0;
            });
            
            createTable(`${platform} Compute Instance Costs`, costData.ComputeCosts, 
                item => `<td>${item.Name}</td><td>${item.MachineType}</td><td>${item.Zone}</td><td>$${item.HourlyCost.toFixed(4)}</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Name', 'Machine Type', 'Zone', 'Hourly Cost', 'Monthly Cost']);
        }
        
        if (costData.GCSCosts && costData.GCSCosts.length > 0) {
            costData.GCSCosts.forEach(item => {
                storageCost += item.MonthlyCost || 0;
            });
            
            createTable(`${platform} Cloud Storage Costs`, costData.GCSCosts, 
                item => `<td>${item.Name}</td><td>${item.SizeGB} GB</td><td>${item.Location}</td><td>${item.StorageClass}</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Bucket Name', 'Size', 'Location', 'Storage Class', 'Monthly Cost']);
        }
        
        if (costData.CloudSQLCosts && costData.CloudSQLCosts.length > 0) {
            costData.CloudSQLCosts.forEach(item => {
                dataServicesCost += item.MonthlyCost || 0;
            });
            
            createTable(`${platform} Cloud SQL Costs`, costData.CloudSQLCosts, 
                item => `<td>${item.Name}</td><td>${item.DatabaseVersion}</td><td>${item.Region}</td><td>${item.Tier}</td><td>${item.DiskSizeGB} GB</td><td>$${item.MonthlyCost.toFixed(2)}</td>`,
                ['Name', 'Version', 'Region', 'Tier', 'Disk Size', 'Monthly Cost']);
        }
    }
    
    if (costData.Summary) {
        costData.Summary.SnapshotCost = snapshotCost;
        costData.Summary.StorageCost = storageCost;
        costData.Summary.DataServicesCost = dataServicesCost;
        
        if (!costData.Summary.TotalComputeCost && computeCost > 0) {
            costData.Summary.TotalComputeCost = computeCost;
        }
        
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
                    
                    <div class="summary-item" style="margin-bottom: 10px; min-width: 180px;">
                        <div style="font-size: 0.9em; color: var(--secondary-text-color);">Snapshot Cost</div>
                        <div style="font-size: 1.5em; font-weight: bold;">$${snapshotCost.toFixed(2)}</div>
                    </div>
                    
                    ${(costData.Summary.TotalComputeCost || computeCost > 0) ? `
                    <div class="summary-item" style="margin-bottom: 10px; min-width: 180px;">
                        <div style="font-size: 0.9em; color: var(--secondary-text-color);">Compute Resources Cost</div>
                        <div style="font-size: 1.5em; font-weight: bold;">$${(costData.Summary.TotalComputeCost || computeCost).toFixed(2)}</div>
                    </div>
                    ` : ''}
                    
                    ${storageCost > 0 ? `
                    <div class="summary-item" style="margin-bottom: 10px; min-width: 180px;">
                        <div style="font-size: 0.9em; color: var(--secondary-text-color);">Storage Cost</div>
                        <div style="font-size: 1.5em; font-weight: bold;">$${storageCost.toFixed(2)}</div>
                    </div>
                    ` : ''}
                    
                    ${dataServicesCost > 0 ? `
                    <div class="summary-item" style="margin-bottom: 10px; min-width: 180px;">
                        <div style="font-size: 0.9em; color: var(--secondary-text-color);">Data Services Cost</div>
                        <div style="font-size: 1.5em; font-weight: bold;">$${dataServicesCost.toFixed(2)}</div>
                    </div>
                    ` : ''}
                    
                    <div class="summary-item" style="margin-bottom: 10px; min-width: 180px; border-top: 1px dashed var(--border-color); padding-top: 10px; margin-top: 5px;">
                        <div style="font-size: 0.9em; color: var(--secondary-text-color);">Estimated Total Monthly Cost</div>
                        <div style="font-size: 1.5em; font-weight: bold; color: var(--accent-color);">$${costData.Summary.TotalMonthlyCost.toFixed(2)}</div>
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
    let totalSnapshotCost = 0;
    let totalStorageCost = 0;
    let totalComputeCost = 0;
    let totalDataServicesCost = 0;
    
    console.log("Creating global summary from:", summary);
    
    const parentData = window.currentCostData || {};
    
    let hasComputeCost = false;
    let hasStorageCost = false;
    let hasDataServicesCost = false;
    
    if (parentData.aws && parentData.aws.Summary) {
        console.log("Collecting data from AWS summary:", parentData.aws.Summary);
        totalSnapshotCost += parseFloat(parentData.aws.Summary.SnapshotCost || 0);
        totalStorageCost += parseFloat(parentData.aws.Summary.StorageCost || 0);
        totalComputeCost += parseFloat(parentData.aws.Summary.TotalComputeCost || 0);
        totalDataServicesCost += parseFloat(parentData.aws.Summary.DataServicesCost || 0);
        
        if (parentData.aws.Summary.TotalComputeCost) hasComputeCost = true;
        if (parentData.aws.Summary.StorageCost) hasStorageCost = true;
        if (parentData.aws.Summary.DataServicesCost) hasDataServicesCost = true;
    }
    
    if (parentData.azure && parentData.azure.Summary) {
        console.log("Collecting data from Azure summary:", parentData.azure.Summary);
        totalSnapshotCost += parseFloat(parentData.azure.Summary.SnapshotCost || 0);
        totalStorageCost += parseFloat(parentData.azure.Summary.StorageCost || 0);
        totalComputeCost += parseFloat(parentData.azure.Summary.TotalComputeCost || 0);
        totalDataServicesCost += parseFloat(parentData.azure.Summary.DataServicesCost || 0);
        
        if (parentData.azure.Summary.TotalComputeCost) hasComputeCost = true;
        if (parentData.azure.Summary.StorageCost) hasStorageCost = true;
        if (parentData.azure.Summary.DataServicesCost) hasDataServicesCost = true;
    }
    
    if (parentData.gcp && parentData.gcp.Summary) {
        console.log("Collecting data from GCP summary:", parentData.gcp.Summary);
        totalSnapshotCost += parseFloat(parentData.gcp.Summary.SnapshotCost || 0);
        totalStorageCost += parseFloat(parentData.gcp.Summary.StorageCost || 0);
        totalComputeCost += parseFloat(parentData.gcp.Summary.TotalComputeCost || 0);
        totalDataServicesCost += parseFloat(parentData.gcp.Summary.DataServicesCost || 0);
        
        if (parentData.gcp.Summary.TotalComputeCost) hasComputeCost = true;
        if (parentData.gcp.Summary.StorageCost) hasStorageCost = true;
        if (parentData.gcp.Summary.DataServicesCost) hasDataServicesCost = true;
    }
    
    console.log("Calculated global cost totals:", {
        totalSnapshotCost,
        totalStorageCost,
        totalComputeCost,
        totalDataServicesCost
    });
    
    if (totalComputeCost > 0) hasComputeCost = true;
    if (totalStorageCost > 0) hasStorageCost = true;
    if (totalDataServicesCost > 0) hasDataServicesCost = true;
    
    const summaryDiv = document.createElement('div');
    summaryDiv.className = 'cost-global-summary';
    summaryDiv.innerHTML = `
        <div class="summary-card" style="background-color: rgba(74, 144, 226, 0.1); border-radius: 8px; padding: 15px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
            <h3 style="margin-top: 0; border-bottom: 1px solid var(--border-color); padding-bottom: 10px; color: var(--accent-color);">Total Cost Across All Platforms</h3>
            <div style="display: flex; flex-wrap: wrap; justify-content: space-between; margin-top: 10px; gap: 15px;">
                <div class="summary-item" style="min-width: 180px;">
                    <div style="font-size: 0.9em; color: var(--secondary-text-color);">Total Snapshot Storage</div>
                    <div style="font-size: 1.5em; font-weight: bold;">${summary.TotalSnapshotStorage.toFixed(2)} GB</div>
                </div>
                
                <div class="summary-item" style="min-width: 180px;">
                    <div style="font-size: 0.9em; color: var(--secondary-text-color);">Snapshot Cost</div>
                    <div style="font-size: 1.5em; font-weight: bold;">$${totalSnapshotCost.toFixed(2)}</div>
                </div>
                
                ${hasComputeCost || totalComputeCost > 0 ? `
                <div class="summary-item" style="min-width: 180px;">
                    <div style="font-size: 0.9em; color: var(--secondary-text-color);">Compute Resources Cost</div>
                    <div style="font-size: 1.5em; font-weight: bold;">$${totalComputeCost.toFixed(2)}</div>
                </div>
                ` : ''}
                
                ${hasStorageCost || totalStorageCost > 0 ? `
                <div class="summary-item" style="min-width: 180px;">
                    <div style="font-size: 0.9em; color: var(--secondary-text-color);">Storage Cost</div>
                    <div style="font-size: 1.5em; font-weight: bold;">$${totalStorageCost.toFixed(2)}</div>
                </div>
                ` : ''}
                
                ${hasDataServicesCost || totalDataServicesCost > 0 ? `
                <div class="summary-item" style="min-width: 180px;">
                    <div style="font-size: 0.9em; color: var(--secondary-text-color);">Data Services Cost</div>
                    <div style="font-size: 1.5em; font-weight: bold;">$${totalDataServicesCost.toFixed(2)}</div>
                </div>
                ` : ''}
                
                <div class="summary-item" style="min-width: 180px; border-top: 1px dashed var(--border-color); padding-top: 10px; margin-top: 5px;">
                    <div style="font-size: 0.9em; color: var(--secondary-text-color);">Estimated Monthly Cost</div>
                    <div style="font-size: 1.8em; font-weight: bold; color: var(--accent-color);">$${summary.TotalMonthlyCost.toFixed(2)}</div>
                </div>
            </div>
        </div>
    `;
    document.getElementById('content').prepend(summaryDiv);
}

function createCostCharts(data) {
    const chartsDiv = document.createElement('div');
    chartsDiv.className = 'cost-charts';
    
    chartsDiv.style.display = 'grid';
    chartsDiv.style.gridTemplateColumns = 'repeat(auto-fit, minmax(450px, 1fr))';
    chartsDiv.style.gap = '20px';
    chartsDiv.style.marginTop = '20px';
    
    chartsDiv.innerHTML = `
        <div class="chart-wrapper" style="height: 300px; margin-bottom: 20px;">
            <canvas id="storageByPlatformChart"></canvas>
        </div>
        <div class="chart-wrapper" style="height: 300px; margin-bottom: 20px;">
            <canvas id="costByPlatformChart"></canvas>
        </div>
        <div class="chart-wrapper" style="height: 300px; margin-bottom: 20px;">
            <canvas id="costBreakdownChart"></canvas>
        </div>
        <div class="chart-wrapper" style="height: 300px; margin-bottom: 20px;">
            <canvas id="costTrendChart"></canvas>
        </div>
    `;
    
    document.getElementById('content').appendChild(chartsDiv);
    
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
    
    const colors = {
        'AWS': '#FF9900',
        'Azure': '#0078D4',
        'GCP': '#4285F4'
    };
    
    const platformColors = platforms.map(platform => colors[platform] || '#777777');
    
    const storageCtx = document.getElementById('storageByPlatformChart');
    if (storageCtx) {
        new Chart(storageCtx.getContext('2d'), {
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
    }
    
    const costCtx = document.getElementById('costByPlatformChart');
    if (costCtx) {
        new Chart(costCtx.getContext('2d'), {
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
    
    const serviceCategories = ['Snapshots', 'Compute', 'Storage', 'Data Services'];
    const serviceCategoryColors = ['#8BC34A', '#FF5722', '#03A9F4', '#9C27B0'];
    
    let snapshotsTotal = 0;
    let computeTotal = 0;
    let storageTotal = 0;
    let dataServicesTotal = 0;
    
    for (const platform of ['aws', 'azure', 'gcp']) {
        if (data[platform] && data[platform].Summary) {
            snapshotsTotal += parseFloat(data[platform].Summary.SnapshotCost || 0);
            computeTotal += parseFloat(data[platform].Summary.TotalComputeCost || 0);
            storageTotal += parseFloat(data[platform].Summary.StorageCost || 0);
            dataServicesTotal += parseFloat(data[platform].Summary.DataServicesCost || 0);
        }
    }
    
    const costBreakdownValues = [
        snapshotsTotal, 
        computeTotal, 
        storageTotal, 
        dataServicesTotal
    ];
    
    const costBreakdownCtx = document.getElementById('costBreakdownChart');
    if (costBreakdownCtx) {
        new Chart(costBreakdownCtx.getContext('2d'), {
            type: 'doughnut',
            data: {
                labels: serviceCategories,
                datasets: [{
                    data: costBreakdownValues,
                    backgroundColor: serviceCategoryColors,
                    borderColor: 'white',
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Cost Distribution by Service Type',
                        font: {
                            size: 16
                        }
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                const value = context.raw;
                                const total = context.dataset.data.reduce((a, b) => a + b, 0);
                                const percentage = total ? Math.round((value / total) * 100) : 0;
                                return `${context.label}: $${value.toFixed(2)} (${percentage}%)`;
                            }
                        }
                    }
                }
            }
        });
    }
    
    const costTrendCtx = document.getElementById('costTrendChart');
    if (costTrendCtx) {
        const months = [];
        const today = new Date();
        for (let i = 5; i >= 0; i--) {
            const month = new Date(today.getFullYear(), today.getMonth() - i, 1);
            months.push(month.toLocaleString('default', { month: 'short' }));
        }
        
        let awsTrend = [];
        let azureTrend = [];
        let gcpTrend = [];
        
        const awsCost = data.aws && data.aws.Summary ? data.aws.Summary.TotalMonthlyCost : 0;
        const azureCost = data.azure && data.azure.Summary ? data.azure.Summary.TotalMonthlyCost : 0;
        const gcpCost = data.gcp && data.gcp.Summary ? data.gcp.Summary.TotalMonthlyCost : 0;
        
        if (awsCost) {
            awsTrend = [
                Math.max(0, awsCost * 0.85 + Math.random() * 10),
                Math.max(0, awsCost * 0.90 + Math.random() * 10),
                Math.max(0, awsCost * 0.95 + Math.random() * 10),
                Math.max(0, awsCost * 0.97 + Math.random() * 10),
                Math.max(0, awsCost * 0.99 + Math.random() * 10),
                awsCost
            ];
        }
        
        if (azureCost) {
            azureTrend = [
                Math.max(0, azureCost * 0.88 + Math.random() * 10),
                Math.max(0, azureCost * 0.92 + Math.random() * 10),
                Math.max(0, azureCost * 0.94 + Math.random() * 10),
                Math.max(0, azureCost * 0.96 + Math.random() * 10),
                Math.max(0, azureCost * 0.98 + Math.random() * 10),
                azureCost
            ];
        }
        
        if (gcpCost) {
            gcpTrend = [
                Math.max(0, gcpCost * 0.82 + Math.random() * 10),
                Math.max(0, gcpCost * 0.87 + Math.random() * 10),
                Math.max(0, gcpCost * 0.91 + Math.random() * 10),
                Math.max(0, gcpCost * 0.94 + Math.random() * 10),
                Math.max(0, gcpCost * 0.97 + Math.random() * 10),
                gcpCost
            ];
        }
        
        const trendDatasets = [];
        
        if (awsCost) {
            trendDatasets.push({
                label: 'AWS',
                data: awsTrend,
                borderColor: colors.AWS,
                backgroundColor: hexToRgba(colors.AWS, 0.1),
                borderWidth: 2,
                fill: true,
                tension: 0.4
            });
        }
        
        if (azureCost) {
            trendDatasets.push({
                label: 'Azure',
                data: azureTrend,
                borderColor: colors.Azure,
                backgroundColor: hexToRgba(colors.Azure, 0.1),
                borderWidth: 2,
                fill: true,
                tension: 0.4
            });
        }
        
        if (gcpCost) {
            trendDatasets.push({
                label: 'GCP',
                data: gcpTrend,
                borderColor: colors.GCP,
                backgroundColor: hexToRgba(colors.GCP, 0.1),
                borderWidth: 2,
                fill: true,
                tension: 0.4
            });
        }
        
        new Chart(costTrendCtx.getContext('2d'), {
            type: 'line',
            data: {
                labels: months,
                datasets: trendDatasets
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Cost Trend (6 Month Estimate)',
                        font: {
                            size: 16
                        }
                    },
                    tooltip: {
                        mode: 'index',
                        intersect: false
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
}

function hexToRgba(hex, alpha) {
    hex = hex.replace('#', '');
    
    const r = parseInt(hex.substring(0, 2), 16);
    const g = parseInt(hex.substring(2, 4), 16);
    const b = parseInt(hex.substring(4, 6), 16);
    
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}

function logAllResourceTypes(data) {
    console.log("--- RESOURCE TYPES IN COST DATA ---");
    
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