// snapshot.js


registerDataHandler('snapshot', 
    function(data) {
        return data.kubernetes || data.aws || data.azure || data.gcp;
    },
    function(data) {
        console.log("Processing snapshot data:", data);
        
        if (data.kubernetes && data.kubernetes.VolumeSnapshots && data.kubernetes.VolumeSnapshots.length > 0) {
            createTable('Kubernetes Volume Snapshots', data.kubernetes.VolumeSnapshots, kubernetesVolumeSnapshotRowTemplate, 
                ['Name', 'Namespace', 'Volume', 'Creation Time', 'Restore Size', 'Status']);
        }
        
        if (data.kubernetes && data.kubernetes.VolumeSnapshotContents && data.kubernetes.VolumeSnapshotContents.length > 0) {
            createTable('Kubernetes Volume Snapshot Contents', data.kubernetes.VolumeSnapshotContents, kubernetesVolumeSnapshotContentRowTemplate, 
                ['Name', 'Driver', 'Volume Handle', 'Snapshot Handle', 'Restore Size']);
        }
        
        if (data.aws && data.aws.EBSSnapshots && data.aws.EBSSnapshots.length > 0) {
            createTable('AWS EBS Snapshots', data.aws.EBSSnapshots, awsEbsSnapshotRowTemplate, 
                ['Snapshot ID', 'Volume ID', 'Size', 'State', 'Creation Time', 'Description', 'Encrypted']);
        }
        
        if (data.aws && data.aws.RDSSnapshots && data.aws.RDSSnapshots.length > 0) {
            createTable('AWS RDS Snapshots', data.aws.RDSSnapshots, awsRdsSnapshotRowTemplate, 
                ['Snapshot ID', 'DB Instance', 'Type', 'Status', 'Engine', 'Size', 'Creation Time', 'Encrypted']);
        }
        
        if (data.azure) {
            console.log("Azure data found:", data.azure);
            if (data.azure.DiskSnapshots) {
                console.log("Azure DiskSnapshots found:", data.azure.DiskSnapshots);
                console.log("DiskSnapshots length:", data.azure.DiskSnapshots.length);
                if (data.azure.DiskSnapshots.length > 0) {
                    console.log("First snapshot item:", data.azure.DiskSnapshots[0]);
                }
            } else {
                console.log("No Azure DiskSnapshots found in data");
            }
        } else {
            console.log("No Azure data found");
        }

        if (data.azure && data.azure.DiskSnapshots && data.azure.DiskSnapshots.length > 0) {
            createTable('Azure Disk Snapshots', data.azure.DiskSnapshots, azureDiskSnapshotRowTemplate, 
                ['Name', 'Location', 'Size', 'State', 'Creation Time']);
         }
        
        if (data.gcp && data.gcp.DiskSnapshots && data.gcp.DiskSnapshots.length > 0) {
            createTable('GCP Disk Snapshots', data.gcp.DiskSnapshots, gcpDiskSnapshotRowTemplate, 
                ['Name', 'Source Disk', 'Size', 'Status', 'Creation Time']);
        }
        
        if (!data.kubernetes?.VolumeSnapshots?.length && 
            !data.kubernetes?.VolumeSnapshotContents?.length && 
            !data.aws?.EBSSnapshots?.length && 
            !data.aws?.RDSSnapshots?.length && 
            !data.azure?.DiskSnapshots?.length && 
            !data.gcp?.DiskSnapshots?.length) {
            document.getElementById('content').innerHTML = `
                <div class="empty-state">
                    <h3><i class="fas fa-camera"></i> No snapshots found</h3>
                    <p>No snapshots were found in the selected platforms. This could mean:</p>
                    <ul>
                        <li>You don't have any snapshots in these platforms</li>
                        <li>You don't have permission to access the snapshots</li>
                        <li>There was an error retrieving the snapshots</li>
                    </ul>
                </div>
            `;
        }
        
        setTimeout(() => {
            console.log(`Created Snapshot tables`);
        }, 100);
    }
);

function kubernetesVolumeSnapshotRowTemplate(item) {
    let statusClasses = {
        "Ready": "badge-success",
        "Creating": "badge-info",
        "Pending": "badge-warning",
        "Error": "badge-danger"
    };
    
    let stateClass = statusClasses[item.State] || "badge-secondary";
    let statusBadge = `<span class="badge ${stateClass}">${item.State}</span>`;
    
    return `<td>${item.Name}</td>
            <td>${item.Namespace || "-"}</td>
            <td>${item.Volume || "-"}</td>
            <td>${item.CreationTimestamp || "-"}</td>
            <td>${item.RestoreSize || "-"}</td>
            <td>${statusBadge}</td>`;
}

function kubernetesVolumeSnapshotContentRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Driver || "-"}</td><td>${item.VolumeHandle || "-"}</td><td>${item.SnapshotHandle || "-"}</td><td>${item.RestoreSize || "-"}</td>`;
}

function awsEbsSnapshotRowTemplate(item) {
    let encrypted = item.Encrypted === "true" ? '<i class="fas fa-lock" title="Encrypted"></i>' : '<i class="fas fa-unlock" title="Not encrypted"></i>';
    return `<td>${item.SnapshotId}</td><td>${item.VolumeId}</td><td>${item.VolumeSize}</td><td>${item.State}</td><td>${item.StartTime}</td><td>${item.Description || "-"}</td><td>${encrypted}</td>`;
}

function awsRdsSnapshotRowTemplate(item) {
    let encrypted = item.Encrypted === "true" ? '<i class="fas fa-lock" title="Encrypted"></i>' : '<i class="fas fa-unlock" title="Not encrypted"></i>';
    return `<td>${item.SnapshotId}</td><td>${item.DBInstanceId}</td><td>${item.SnapshotType}</td><td>${item.Status}</td><td>${item.Engine}</td><td>${item.AllocatedStorage}</td><td>${item.CreationTime}</td><td>${encrypted}</td>`;
}

function azureDiskSnapshotRowTemplate(item) {
    console.log("Azure snapshot item:", item);
    
    const name = item.Name || "Unknown";
    const location = item.Location || "-";
    
    let sizeDisplay = "-";
    if (item.SizeGB) {
        sizeDisplay = `${item.SizeGB} GB`;
    } else if (item.DiskSizeGB) {
        sizeDisplay = `${item.DiskSizeGB} GB`;
    }
    
    let state = "-";
    if (item.State) {
        state = item.State;
    } else if (item.ProvisioningState) {
        state = item.ProvisioningState;
    } else if (item.Status) {
        state = item.Status;
    }
    
    const creationTime = item.CreationTime || item.TimeCreated || "-";
    
    return `<td>${name}</td><td>${location}</td><td>${sizeDisplay}</td><td>${state}</td><td>${creationTime}</td>`;
}

function gcpDiskSnapshotRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.SourceDiskName || "-"}</td><td>${item.DiskSizeGB} GB</td><td>${item.Status}</td><td>${item.CreationTime}</td>`;
}

function showSnapshotHunterModal() {
    console.log("Creating Snapshot Hunter modal");
    
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
    modalContent.className = 'modal-content snapshot-modal';
    modalContent.style.backgroundColor = 'var(--card-bg)';
    modalContent.style.color = 'var(--text-color)';
    modalContent.style.padding = '25px';
    modalContent.style.borderRadius = '8px';
    modalContent.style.maxWidth = '500px';
    modalContent.style.width = '90%';
    modalContent.style.boxShadow = '0 5px 20px rgba(0,0,0,0.4)';
    modalContent.style.border = '1px solid var(--border-color)';
    
    modalContent.innerHTML = `
        <h3 style="margin-top: 0; color: var(--accent-color); font-size: 1.5em; border-bottom: 1px solid var(--border-color); padding-bottom: 10px;">
            <i class="fas fa-camera"></i> Snapshot Hunter
        </h3>
        
        <p style="margin-bottom: 20px;">Select which cloud platforms to gather snapshot inventory from:</p>
        
        <div class="snapshot-platform-selection" style="margin: 20px 0;">
            <div style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 10px;">
                <div style="display: flex; align-items: center;">
                    <input type="checkbox" id="snapshot-k8s" value="kubernetes" style="margin-right: 10px;">
                    <label for="snapshot-k8s" style="font-weight: bold; font-size: 1.1em;">
                        <i class="fas fa-dharmachakra"></i> Kubernetes Volume Snapshots
                    </label>
                </div>
                <div style="padding: 5px 0 0 25px; font-size: 0.9em; color: var(--secondary-text-color);">
                    Collects Volume Snapshots and Snapshot Contents
                </div>
            </div>
            
            <div style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 10px;">
                <div style="display: flex; align-items: center;">
                    <input type="checkbox" id="snapshot-aws" value="aws" style="margin-right: 10px;">
                    <label for="snapshot-aws" style="font-weight: bold; font-size: 1.1em;">
                        <i class="fab fa-aws"></i> AWS Snapshots
                    </label>
                </div>
                <div style="padding: 5px 0 0 25px; font-size: 0.9em; color: var(--secondary-text-color);">
                    Collects EBS and RDS snapshots
                </div>
            </div>
            
            <div style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 10px;">
                <div style="display: flex; align-items: center;">
                    <input type="checkbox" id="snapshot-azure" value="azure" style="margin-right: 10px;">
                    <label for="snapshot-azure" style="font-weight: bold; font-size: 1.1em;">
                        <i class="fab fa-microsoft"></i> Azure Snapshots
                    </label>
                </div>
                <div style="padding: 5px 0 0 25px; font-size: 0.9em; color: var(--secondary-text-color);">
                    Collects disk snapshots
                </div>
            </div>
            
            <div style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 10px;">
                <div style="display: flex; align-items: center;">
                    <input type="checkbox" id="snapshot-gcp" value="gcp" style="margin-right: 10px;">
                    <label for="snapshot-gcp" style="font-weight: bold; font-size: 1.1em;">
                        <i class="fab fa-google"></i> GCP Snapshots
                    </label>
                </div>
                <div style="padding: 5px 0 0 25px; font-size: 0.9em; color: var(--secondary-text-color);">
                    Collects disk snapshots
                </div>
            </div>
            
            <div style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 10px;">
                <div style="display: flex; align-items: center;">
                    <input type="checkbox" id="snapshot-all" value="all" style="margin-right: 10px;">
                    <label for="snapshot-all" style="font-weight: bold; font-size: 1.1em;">
                        <i class="fas fa-check-double"></i> All Platforms
                    </label>
                </div>
                <div style="padding: 5px 0 0 25px; font-size: 0.9em; color: var(--secondary-text-color);">
                    Collects snapshots from all connected platforms
                </div>
            </div>
        </div>
        
        <div class="modal-buttons" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 25px; border-top: 1px solid var(--border-color); padding-top: 15px;">
            <button id="snapshot-cancel-btn" class="btn" style="padding: 10px 20px; background-color: var(--button-bg-color); color: var(--button-text-color); border: none; border-radius: 4px; cursor: pointer;">Cancel</button>
            <button id="snapshot-collect-btn" class="btn btn-primary" style="padding: 10px 20px; background-color: var(--accent-color); color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">
                <i class="fas fa-search"></i> Find Snapshots
            </button>
        </div>
    `;
    
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    console.log("Modal created and added to DOM");
    

document.getElementById('snapshot-all').addEventListener('change', function() {
    const isChecked = this.checked;
    
    document.querySelectorAll('.snapshot-platform-selection input[type="checkbox"]:not(#snapshot-all)').forEach(checkbox => {
        checkbox.checked = false;
        checkbox.disabled = isChecked;
    });
});

const platformCheckboxes = [
    document.getElementById('snapshot-k8s'),
    document.getElementById('snapshot-aws'),
    document.getElementById('snapshot-azure'),
    document.getElementById('snapshot-gcp')
];

platformCheckboxes.forEach(checkbox => {
    checkbox.addEventListener('change', function() {
        if (this.checked) {
            platformCheckboxes.forEach(box => {
                if (box !== this) {
                    box.checked = false;
                }
            });
            
            document.getElementById('snapshot-all').checked = false;
        }
    });
});

document.getElementById('snapshot-cancel-btn').addEventListener('click', () => {
    console.log("Cancel button clicked");
    modal.remove();
});

document.getElementById('snapshot-collect-btn').addEventListener('click', () => {
    console.log("Collect button clicked");
    
    const selectedPlatforms = [];
    
    if (document.getElementById('snapshot-all').checked) {
        selectedPlatforms.push('all');
    } else {
        if (document.getElementById('snapshot-k8s').checked) selectedPlatforms.push('kubernetes');
        if (document.getElementById('snapshot-aws').checked) selectedPlatforms.push('aws');
        if (document.getElementById('snapshot-azure').checked) selectedPlatforms.push('azure');
        if (document.getElementById('snapshot-gcp').checked) selectedPlatforms.push('gcp');
    }
    
    if (selectedPlatforms.length === 0) {
        alert('Please select at least one platform');
        return;
    }
    
    collectSnapshots(selectedPlatforms);
    modal.remove();
});
    
    checkPlatformAvailability();
}

function checkPlatformAvailability() {
    fetch('/api/check-credentials')
        .then(response => response.json())
        .then(data => {
            console.log('Connection status for snapshots:', data);
            
            document.getElementById('snapshot-k8s').disabled = !data.kubernetes;
            
            document.getElementById('snapshot-aws').disabled = !data.aws;
            
            document.getElementById('snapshot-azure').disabled = !data.azure;
            
            document.getElementById('snapshot-gcp').disabled = !data.gcp;
            
            const availablePlatforms = [data.kubernetes, data.aws, data.azure, data.gcp].filter(Boolean).length;
            
            if (availablePlatforms === 0) {
                const modalContent = document.querySelector('.snapshot-modal');
                
                const messageDiv = document.createElement('div');
                messageDiv.style.backgroundColor = 'rgba(255, 165, 0, 0.1)';
                messageDiv.style.borderLeft = '4px solid #FFA500';
                messageDiv.style.padding = '10px';
                messageDiv.style.marginBottom = '20px';
                
                messageDiv.innerHTML = `
                    <p style="margin: 0; font-size: 0.9em; color: var(--text-color);">
                        <i class="fas fa-exclamation-triangle"></i> <strong>No connected platforms detected.</strong> 
                        <br>Please connect to at least one platform (Kubernetes, AWS, Azure, or GCP) before using Snapshot Hunter.
                    </p>
                `;
                
                modalContent.insertBefore(messageDiv, modalContent.children[1]);
                
                document.getElementById('snapshot-collect-btn').disabled = true;
                document.getElementById('snapshot-all').disabled = true;
            }
        })
        .catch(error => {
            console.error('Error checking platform connections:', error);
        });
}

function collectSnapshots(platforms) {
    console.log(`Collecting snapshots from: ${platforms.join(', ')}`);
    showLoadingIndicator();
    
    let targetPlatform = platforms[0];
    if (platforms.includes('all')) {
        targetPlatform = 'all';
    }
    
    console.log(`Making API request to: /api/snapshots?platform=${targetPlatform}`);
    
    fetch(`/api/snapshots?platform=${targetPlatform}`)
        .then(response => {
            console.log(`Response status: ${response.status}`);
            if (!response.ok) {
                return response.text().then(text => {
                    console.error(`Error response body: ${text}`);
                    throw new Error(`HTTP error ${response.status}: ${text}`);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("Raw snapshot data received:", data);
            console.log("Data keys:", Object.keys(data));
            
            document.getElementById('content').innerHTML = '';
            
            if (targetPlatform === 'all') {
                console.log("Processing multi-platform data");
                processWithHandler(data);
            } else {
                console.log(`Processing single platform data for: ${targetPlatform}`);
                const wrappedData = {};
                wrappedData[targetPlatform] = data;
                console.log("Wrapped data:", wrappedData);
                processWithHandler(wrappedData);
            }
            
            updateResourceNav();
        })
        .catch(error => {
            console.error("Error collecting snapshots:", error);
            alert(`Error collecting snapshots: ${error.message}`);
            document.getElementById('content').innerHTML = `
                <div class="error-message">
                    <h2>Error Collecting Snapshots</h2>
                    <p>${error.message}</p>
                </div>
            `;
        })
        .finally(() => {
            hideLoadingIndicator();
        });
}

document.addEventListener('DOMContentLoaded', function() {
    console.log("DOM loaded - Snapshot module setting up event listener");
    
    const snapshotButton = document.getElementById('snapshot-button');
    if (snapshotButton) {
        console.log("Found snapshot button, setting up handler");
        
        const newButton = snapshotButton.cloneNode(true);
        if (snapshotButton.parentNode) {
            snapshotButton.parentNode.replaceChild(newButton, snapshotButton);
        }
        
        newButton.addEventListener('click', function(event) {
            console.log("Snapshot button clicked");
            event.preventDefault();
            
            showSnapshotHunterModal();
        });
    } else {
        console.error("Could not find snapshot-button element");
    }
});