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

document.addEventListener('DOMContentLoaded', function() {
    console.log("DOM loaded - Azure module setting up event listener");
    
    const azureButton = document.getElementById('azure-button');
    if (azureButton) {
        console.log("Found Azure button, setting up handler");
        
        const newButton = azureButton.cloneNode(true);
        if (azureButton.parentNode) {
            azureButton.parentNode.replaceChild(newButton, azureButton);
        }
        
        newButton.addEventListener('click', function(event) {
            console.log("Azure button clicked");
            event.preventDefault();
            
            showAzureCredentialsModal();
        });
    } else {
        console.error("Could not find azure-button element");
    }
});

function showAzureCredentialsModal() {
    console.log("Creating Azure credentials modal");
    const isConnected = document.getElementById('azure-button')?.classList.contains('connected');
    
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
    modalContent.className = 'modal-content azure-modal';
    modalContent.style.backgroundColor = 'var(--card-bg)';
    modalContent.style.color = 'var(--text-color)';
    modalContent.style.padding = '25px';
    modalContent.style.borderRadius = '8px';
    modalContent.style.maxWidth = '500px';
    modalContent.style.width = '90%';
    modalContent.style.boxShadow = '0 5px 20px rgba(0,0,0,0.4)';
    modalContent.style.border = '1px solid var(--border-color)';
    
    let connectionNote = '';
    if (isConnected) {
        connectionNote = `
            <div style="background-color: rgba(0,255,0,0.1); border-left: 4px solid #4CAF50; padding: 8px; margin-bottom: 15px;">
                <p style="margin: 0; color: var(--text-color);">
                    <i class="fas fa-info-circle"></i> You are already connected to Azure. 
                    You can switch to another subscription or use different credentials.
                </p>
            </div>
        `;
    }
    
    modalContent.innerHTML = `
        <h3 style="margin-top: 0; color: var(--accent-color); font-size: 1.5em; border-bottom: 1px solid var(--border-color); padding-bottom: 10px;">
            <i class="fab fa-microsoft"></i> Connect to Azure
        </h3>
        
        ${connectionNote}
        
        <div class="azure-connection-form" style="margin-top: 20px;">
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="azure-default-config" name="azure-config-source" value="default" checked>
                <label for="azure-default-config" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-cog"></i> Use Azure CLI Configuration
                </label>
                <div id="azure-default-config-form" class="source-form" style="margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <p style="margin-top: 0;">This option uses credentials and settings from your Azure CLI.</p>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="azure-subscription-selector" style="font-weight: bold; margin-bottom: 5px;">Select Azure Subscription:</label>
                        <select id="azure-subscription-selector" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;">
                            <option value="">Loading subscriptions...</option>
                        </select>
                    </div>
                </div>
            </div>
            
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="azure-manual-config" name="azure-config-source" value="manual">
                <label for="azure-manual-config" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-key"></i> Enter Service Principal Credentials
                </label>
                <div id="azure-manual-config-form" class="source-form" style="display: none; margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <div class="form-group">
                        <label for="azure-tenant-id" style="font-weight: bold; margin-bottom: 5px;">Tenant ID:</label>
                        <input type="text" id="azure-tenant-id" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;" placeholder="00000000-0000-0000-0000-000000000000">
                    </div>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="azure-client-id" style="font-weight: bold; margin-bottom: 5px;">Client ID:</label>
                        <input type="text" id="azure-client-id" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;" placeholder="00000000-0000-0000-0000-000000000000">
                    </div>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="azure-client-secret" style="font-weight: bold; margin-bottom: 5px;">Client Secret:</label>
                        <input type="password" id="azure-client-secret" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;" placeholder="Your client secret">
                    </div>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="azure-subscription-id" style="font-weight: bold; margin-bottom: 5px;">Subscription ID:</label>
                        <input type="text" id="azure-subscription-id" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;" placeholder="00000000-0000-0000-0000-000000000000">
                    </div>
                    <p class="tip" style="margin-top: 15px; font-size: 0.85em; color: var(--secondary-text-color); font-style: italic;">
                        Note: Credentials are used for the current session only and are never stored.
                    </p>
                </div>
            </div>
            
            <div class="modal-buttons" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 25px; border-top: 1px solid var(--border-color); padding-top: 15px;">
                <button id="azure-cancel-btn" class="btn" style="padding: 10px 20px; background-color: var(--button-bg-color); color: var(--button-text-color); border: none; border-radius: 4px; cursor: pointer;">Cancel</button>
                <button id="azure-connect-btn" class="btn btn-primary" style="padding: 10px 20px; background-color: var(--accent-color); color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">
                    <i class="fas fa-plug"></i> Connect
                </button>
            </div>
        </div>
    `;
    
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    console.log("Modal created and added to DOM");

    const sourceForms = document.querySelectorAll('.source-form');
    document.querySelectorAll('input[name="azure-config-source"]').forEach(radio => {
        radio.addEventListener('change', () => {
            console.log(`Radio changed to: ${radio.value}`);
            sourceForms.forEach(form => form.style.display = 'none');
            const selectedForm = document.getElementById(`azure-${radio.value}-config-form`);
            if (selectedForm) {
                selectedForm.style.display = 'block';
            }
        });
    });

    loadAzureSubscriptions();

    document.getElementById('azure-cancel-btn').addEventListener('click', () => {
        console.log("Cancel button clicked");
        modal.remove();
    });

    document.getElementById('azure-connect-btn').addEventListener('click', () => {
    console.log("Connect button clicked");
    
    const configSource = document.querySelector('input[name="azure-config-source"]:checked').value;
    console.log(`Selected source: ${configSource}`);
    
    if (configSource === 'default') {
        const subscription = document.getElementById('azure-subscription-selector').value;
        console.log(`Using Azure subscription: ${subscription || 'default'}`);
        connectToAzure({ 
            type: 'cli', 
            subscription: subscription 
        });
    } else {
            const tenantId = document.getElementById('azure-tenant-id').value;
            const clientId = document.getElementById('azure-client-id').value;
            const clientSecret = document.getElementById('azure-client-secret').value;
            const subscriptionId = document.getElementById('azure-subscription-id').value;
            
            if (!tenantId || !clientId || !clientSecret || !subscriptionId) {
                alert('Please provide all required Azure service principal credentials');
                return;
            }
            
            console.log(`Using manual Azure service principal credentials`);
            connectToAzure({ 
                type: 'service_principal', 
                tenantId: tenantId,
                clientId: clientId,
                clientSecret: clientSecret,
                subscriptionId: subscriptionId
            });
        }
    });

    function loadAzureSubscriptions() {
    console.log("Loading Azure subscriptions");
    
    fetch('/api/azure/subscriptions')
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => {
                    throw new Error(`HTTP error ${response.status}: ${text}`);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("Subscriptions loaded:", data);
            console.log("Number of subscriptions:", data.subscriptions ? data.subscriptions.length : 0);
            
            const subscriptionSelector = document.getElementById('azure-subscription-selector');
            subscriptionSelector.innerHTML = '';
            
            if (data.subscriptions && data.subscriptions.length > 0) {
                data.subscriptions.forEach(subscription => {
                    console.log("Processing subscription:", subscription);
                    
                    const option = document.createElement('option');
                    option.value = subscription.id;
                    option.textContent = `${subscription.name} (${subscription.id})`;
                    
                    if (subscription.isDefault === "true") {
                        option.selected = true;
                    }
                    
                    subscriptionSelector.appendChild(option);
                });
                
                console.log("Added subscriptions to dropdown");
            } else {
                const option = document.createElement('option');
                option.value = "";
                option.textContent = "No subscriptions found - CLI may not be properly detected";
                subscriptionSelector.appendChild(option);
                
                const noteDiv = document.createElement('div');
                noteDiv.style.marginTop = '10px';
                noteDiv.style.padding = '8px';
                noteDiv.style.backgroundColor = 'rgba(255, 165, 0, 0.1)';
                noteDiv.style.borderLeft = '4px solid #FFA500';
                
                let subscriberCount = data.subscriptionsCount || 'multiple';
                noteDiv.innerHTML = `
                    <p style="margin: 0; font-size: 0.9em; color: var(--text-color);">
                        <i class="fas fa-info-circle"></i> <strong>Note:</strong> Your Azure CLI appears to be configured with ${subscriberCount} subscription(s), but we couldn't display them properly.
                        <br><br>
                        Try one of these options:
                        <ul style="margin-top: 5px; margin-bottom: 0; padding-left: 20px;">
                            <li>Click Connect to try with your default subscription (<code>Michael Cade Sub</code>)</li>
                            <li>Use service principal credentials instead (select option above)</li>
                            <li>If using Connect fails, you can manually set your default subscription in terminal with:<br>
                            <code>az account set --subscription "Subscription Name"</code></li>
                        </ul>
                    </p>
                `;
                
                const formGroup = subscriptionSelector.closest('.form-group');
                formGroup.appendChild(noteDiv);
            }
        })
        .catch(error => {
            console.error("Error loading Azure subscriptions:", error);
        });
}

    function connectToAzure(config) {
    console.log("Connecting to Azure with config:", config);
    showLoadingIndicator();
    
    fetch('/api/azure/connect', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(config)
    })
    .then(response => {
        if (!response.ok) {
            return response.text().then(text => {
                if (response.status === 403) {
                    throw new Error("Permission denied: You don't have sufficient permissions for some Azure resources");
                } else if (response.status === 401) {
                    throw new Error("Authentication failed: The credentials provided are not valid");
                } else {
                    throw new Error(`HTTP error ${response.status}: ${text}`);
                }
            });
        }
        return response.json();
    })
    .then(data => {
        console.log("Connection response:", data);
        if (data.status === 'success') {
            const button = document.getElementById('azure-button');
            if (button) {
                button.classList.add('connected');
                button.classList.remove('not-connected');
                
                const existingBadges = button.querySelectorAll('.connection-badge');
                existingBadges.forEach(badge => badge.remove());
                
                const badge = document.createElement('span');
                badge.className = 'connection-badge connected';
                button.appendChild(badge);
                
                button.title = 'Azure (Connected)';
                console.log("Button updated to connected state");
            }
            
            modal.remove();
            console.log("Modal removed, reloading page");
            setTimeout(() => {
                location.reload();
            }, 300);
        } else {
            throw new Error(data.message || 'Failed to connect to Azure');
        }
    })
    .catch(error => {
        console.error('Azure connection error:', error);
        alert(`Error connecting to Azure: ${error.message}`);
    })
    .finally(() => {
        hideLoadingIndicator();
    });
}
}