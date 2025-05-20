// vault.js


registerDataHandler('vault', 
    function(data) {
        console.log("Testing if data is Vault data:", data);
        const isVaultData = !!(data.serverInfo || data.authMethods || data.secretEngines || 
                              data.policies || data.auditDevices || data.replicationInfo);
        console.log("Is Vault data?", isVaultData);
        return isVaultData;
    },
    function(data) {
        console.log("Processing Vault data");
        
        if (data.serverInfo) {
            createTable('Server Info', [data.serverInfo], serverInfoRowTemplate, 
                ['Version', 'Cluster Name', 'Cluster ID', 'HA Enabled', 'Storage Type', 'Sealed Status']);
        }
        
        if (data.replicationInfo) {
            createTable('Replication Status', [data.replicationInfo], replicationInfoRowTemplate, 
                ['DR Enabled', 'DR Mode', 'Performance Enabled', 'Performance Mode', 'Primary Cluster']);
        }
        
        if (data.authMethods) {
            createTable('Authentication Methods', data.authMethods, authMethodRowTemplate, 
                ['Path', 'Type', 'Description', 'Accessor', 'Local']);
        }
        
        if (data.secretEngines) {
            createTable('Secret Engines', data.secretEngines, secretEngineRowTemplate, 
                ['Path', 'Type', 'Description', 'Version', 'Local']);
        }
        
        if (data.policies) {
            createTable('Policies', data.policies, policyRowTemplate, 
                ['Name', 'Type', 'Rules']);
        }
        
        if (data.auditDevices) {
            createTable('Audit Devices', data.auditDevices, auditDeviceRowTemplate, 
                ['Path', 'Type', 'Description']);
        }
        
        if (data.namespaces && data.namespaces.length > 0) {
            createTable('Namespaces', data.namespaces, namespaceRowTemplate, 
                ['Path', 'Description']);
        }
        
        setTimeout(() => {
            console.log(`Created Vault tables`);
        }, 100);
    }
);

function serverInfoRowTemplate(item) {
    const sealed = item.sealed ? 'Yes' : 'No';
    const haEnabled = item.haEnabled ? 'Yes' : 'No';
    
    return `
        <td>${item.version || 'N/A'}</td>
        <td>${item.clusterName || 'N/A'}</td>
        <td>${item.clusterId || 'N/A'}</td>
        <td>${haEnabled}</td>
        <td>${item.storageType || 'N/A'}</td>
        <td>${sealed}</td>
    `;
}

function replicationInfoRowTemplate(item) {
    const drEnabled = item.drEnabled ? 'Yes' : 'No';
    const perfEnabled = item.performanceEnabled ? 'Yes' : 'No';
    
    return `
        <td>${drEnabled}</td>
        <td>${item.drMode || 'N/A'}</td>
        <td>${perfEnabled}</td>
        <td>${item.performanceMode || 'N/A'}</td>
        <td>${item.primaryClusterAddr || 'N/A'}</td>
    `;
}

function authMethodRowTemplate(item) {
    return `
        <td>${item.path || 'N/A'}</td>
        <td>${item.type || 'N/A'}</td>
        <td>${item.description || 'N/A'}</td>
        <td>${item.accessor || 'N/A'}</td>
        <td>${item.local ? 'Yes' : 'No'}</td>
    `;
}

function secretEngineRowTemplate(item) {
    return `
        <td>${item.path || 'N/A'}</td>
        <td>${item.type || 'N/A'}</td>
        <td>${item.description || 'N/A'}</td>
        <td>${item.version || 'N/A'}</td>
        <td>${item.local ? 'Yes' : 'No'}</td>
    `;
}

function policyRowTemplate(item) {
    const policyId = `policy-${item.name}`.replace(/[^a-zA-Z0-9-]/g, '-');
    
    return `
        <td>${item.name || 'N/A'}</td>
        <td>${item.type || 'N/A'}</td>
        <td>
            ${item.rules ? `
                <button class="details-button" onclick="toggleVaultDetails('${policyId}')">
                    <i class="fas fa-info-circle"></i> View Policy
                </button>
                <div id="${policyId}" style="display:none;" class="details-panel code-block">
                    <pre><code>${item.rules}</code></pre>
                </div>
            ` : 'No policy content available'}
        </td>
    `;
}

function auditDeviceRowTemplate(item) {
    return `
        <td>${item.path || 'N/A'}</td>
        <td>${item.type || 'N/A'}</td>
        <td>${item.description || 'N/A'}</td>
    `;
}

function namespaceRowTemplate(item) {
    return `
        <td>${item.path || 'N/A'}</td>
        <td>${item.description || 'N/A'}</td>
    `;
}

function toggleVaultDetails(id) {
    const element = document.getElementById(id);
    if (element) {
        if (element.style.display === 'none') {
            element.style.display = 'block';
        } else {
            element.style.display = 'none';
        }
    }
}

document.addEventListener('DOMContentLoaded', function() {
    console.log("DOM loaded - Vault module setting up event listener");
    
    const vaultButton = document.getElementById('vault-button');
    if (vaultButton) {
        console.log("Found Vault button, setting up handler");
        
        const newButton = vaultButton.cloneNode(true);
        if (vaultButton.parentNode) {
            vaultButton.parentNode.replaceChild(newButton, vaultButton);
        }
        
        newButton.addEventListener('click', function(event) {
            console.log("Vault button clicked");
            event.preventDefault();
            
            showVaultConnectionModal();
        });
    } else {
        console.error("Could not find vault-button element");
    }
});

function showVaultConnectionModal() {
    console.log("Creating Vault connection modal");
    const isConnected = document.getElementById('vault-button')?.classList.contains('connected');
    
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
    modalContent.className = 'modal-content vault-modal';
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
                    <i class="fas fa-info-circle"></i> You are already connected to Vault. 
                    You can switch to another server or use different credentials.
                </p>
            </div>
        `;
    }
    
    modalContent.innerHTML = `
        <h3 style="margin-top: 0; color: var(--accent-color); font-size: 1.5em; border-bottom: 1px solid var(--border-color); padding-bottom: 10px;">
            <i class="fas fa-vault"></i> Connect to HashiCorp Vault
        </h3>
        
        ${connectionNote}
        
        <div class="vault-connection-form" style="margin-top: 20px;">
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="vault-token-auth" name="vault-auth-source" value="token" checked>
                <label for="vault-token-auth" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-key"></i> Token Authentication
                </label>
                <div id="vault-token-auth-form" class="source-form" style="margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <div class="form-group">
                        <label for="vault-server" style="font-weight: bold; margin-bottom: 5px;">Vault Server URL:</label>
                        <input type="text" id="vault-server" placeholder="https://vault.example.com:8200" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px; box-sizing: border-box;">
                    </div>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="vault-token" style="font-weight: bold; margin-bottom: 5px;">Vault Token:</label>
                        <input type="password" id="vault-token" placeholder="hvs.xxxxxxxxxxxx" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px; box-sizing: border-box;">
                    </div>
                </div>
            </div>
            
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="vault-userpass-auth" name="vault-auth-source" value="userpass">
                <label for="vault-userpass-auth" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-user"></i> Username & Password Authentication
                </label>
                <div id="vault-userpass-auth-form" class="source-form" style="display: none; margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <div class="form-group">
                        <label for="vault-userpass-server" style="font-weight: bold; margin-bottom: 5px;">Vault Server URL:</label>
                        <input type="text" id="vault-userpass-server" placeholder="https://vault.example.com:8200" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px; box-sizing: border-box;">
                    </div>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="vault-username" style="font-weight: bold; margin-bottom: 5px;">Username:</label>
                        <input type="text" id="vault-username" placeholder="username" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px; box-sizing: border-box;">
                    </div>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="vault-password" style="font-weight: bold; margin-bottom: 5px;">Password:</label>
                        <input type="password" id="vault-password" placeholder="password" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px; box-sizing: border-box;">
                    </div>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="vault-auth-path" style="font-weight: bold; margin-bottom: 5px;">Auth Path:</label>
                        <input type="text" id="vault-auth-path" value="userpass" placeholder="userpass" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px; box-sizing: border-box;">
                    </div>
                </div>
            </div>
            
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="vault-cli-auth" name="vault-auth-source" value="cli">
                <label for="vault-cli-auth" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-terminal"></i> Use Vault CLI Configuration
                </label>
                <div id="vault-cli-auth-form" class="source-form" style="display: none; margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <p style="margin-top: 0;">This option uses credentials from your local Vault CLI configuration.</p>
                    
                    <div id="vault-cli-status" style="margin-top: 10px; padding: 8px; font-size: 0.9em;">
                        Checking Vault CLI status...
                    </div>
                </div>
            </div>
            
            <div class="form-group" style="margin-top: 10px;">
                <div style="display: flex; align-items: center;">
                    <input type="checkbox" id="vault-ignore-ssl" style="margin-right: 8px;">
                    <label for="vault-ignore-ssl">Ignore SSL certificate errors</label>
                </div>
            </div>
            
            <div class="form-group" style="margin-top: 10px;">
                <div style="display: flex; align-items: center;">
                    <input type="checkbox" id="vault-include-policies" style="margin-right: 8px;">
                    <label for="vault-include-policies">Include policy contents (requires additional permissions)</label>
                </div>
            </div>
            
            <div class="modal-buttons" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 25px; border-top: 1px solid var(--border-color); padding-top: 15px;">
                <button id="vault-cancel-btn" class="btn" style="padding: 10px 20px; background-color: var(--button-bg-color); color: var(--button-text-color); border: none; border-radius: 4px; cursor: pointer;">Cancel</button>
                <button id="vault-connect-btn" class="btn btn-primary" style="padding: 10px 20px; background-color: var(--accent-color); color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">
                    <i class="fas fa-plug"></i> Connect
                </button>
            </div>
        </div>
    `;
    
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    console.log("Modal created and added to DOM");

    const sourceForms = document.querySelectorAll('.source-form');
    document.querySelectorAll('input[name="vault-auth-source"]').forEach(radio => {
        radio.addEventListener('change', () => {
            console.log(`Radio changed to: ${radio.value}`);
            sourceForms.forEach(form => form.style.display = 'none');
            const selectedForm = document.getElementById(`vault-${radio.value}-auth-form`);
            if (selectedForm) {
                selectedForm.style.display = 'block';
            }
        });
    });
    
    checkVaultCliStatus();

    document.getElementById('vault-cancel-btn').addEventListener('click', () => {
        console.log("Cancel button clicked");
        modal.remove();
    });

    document.getElementById('vault-connect-btn').addEventListener('click', () => {
        console.log("Connect button clicked");
        
        const authSource = document.querySelector('input[name="vault-auth-source"]:checked').value;
        console.log(`Selected authentication source: ${authSource}`);
        
        const ignoreSSL = document.getElementById('vault-ignore-ssl').checked;
        const includePolicies = document.getElementById('vault-include-policies').checked;
        
        let config = {
            type: authSource,
            insecure: ignoreSSL,
            includePolicies: includePolicies
        };
        
        if (authSource === 'token') {
            const server = document.getElementById('vault-server').value.trim();
            const token = document.getElementById('vault-token').value.trim();
            
            if (!server || !token) {
                alert('Please provide both Server URL and Vault Token');
                return;
            }
            
            config.server = server;
            config.token = token;
        } 
        else if (authSource === 'userpass') {
            const server = document.getElementById('vault-userpass-server').value.trim();
            const username = document.getElementById('vault-username').value.trim();
            const password = document.getElementById('vault-password').value.trim();
            const authPath = document.getElementById('vault-auth-path').value.trim();
            
            if (!server || !username || !password) {
                alert('Please provide Server URL, Username, and Password');
                return;
            }
            
            config.server = server;
            config.username = username;
            config.password = password;
            config.authPath = authPath || 'userpass';
        }
        
        connectToVault(config);
    });
    
    function checkVaultCliStatus() {
        const statusContainer = document.getElementById('vault-cli-status');
        if (!statusContainer) return;
        
        fetch('/api/vault/cli-status')
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                if (data.installed) {
                    let html = `
                        <div style="background-color: rgba(0,255,0,0.1); border-left: 4px solid #4CAF50; padding: 8px;">
                            <p style="margin: 0;"><i class="fas fa-check-circle"></i> <strong>Vault CLI detected!</strong></p>
                            <p style="margin: 5px 0 0 0; font-size: 0.9em;">Version: ${data.version || 'Unknown'}</p>
                    `;
                    
                    if (data.authenticated) {
                        html += `
                            <p style="margin: 5px 0 0 0; font-size: 0.9em;">
                                <i class="fas fa-unlock"></i> Authenticated: Yes
                            </p>
                        `;
                    } else {
                        html += `
                            <p style="margin: 5px 0 0 0; font-size: 0.9em;">
                                <i class="fas fa-lock"></i> Authenticated: No (Login required)
                            </p>
                        `;
                    }
                    
                    html += `</div>`;
                    statusContainer.innerHTML = html;
                } else {
                    statusContainer.innerHTML = `
                        <div style="background-color: rgba(255,0,0,0.1); border-left: 4px solid #FF5252; padding: 8px;">
                            <p style="margin: 0;"><i class="fas fa-times-circle"></i> <strong>Vault CLI not detected</strong></p>
                            <p style="margin: 5px 0 0 0; font-size: 0.9em;">Please install the Vault CLI or use another authentication method.</p>
                        </div>
                    `;
                }
            })
            .catch(error => {
                console.error("Error checking Vault CLI status:", error);
                statusContainer.innerHTML = `
                    <div style="background-color: rgba(255,165,0,0.1); border-left: 4px solid #FFA500; padding: 8px;">
                        <p style="margin: 0;"><i class="fas fa-exclamation-triangle"></i> <strong>Could not verify Vault CLI</strong></p>
                        <p style="margin: 5px 0 0 0; font-size: 0.9em;">Error: ${error.message}</p>
                    </div>
                `;
            });
    }
    
    function connectToVault(config) {
        console.log("Connecting to Vault with config:", config);
        showLoadingIndicator();
        
        fetch('/api/vault/connect', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(config)
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
            console.log("Connection response:", data);
            if (data.status === 'success') {
                const button = document.getElementById('vault-button');
                if (button) {
                    button.classList.add('connected');
                    button.classList.remove('not-connected');
                    
                    const existingBadges = button.querySelectorAll('.connection-badge');
                    existingBadges.forEach(badge => badge.remove());
                    
                    const badge = document.createElement('span');
                    badge.className = 'connection-badge connected';
                    button.appendChild(badge);
                    
                    button.title = 'Vault (Connected)';
                    console.log("Button updated to connected state");
                }
                
                modal.remove();
                console.log("Modal removed, reloading page");
                setTimeout(() => {
                    location.reload();
                }, 300);
            } else {
                throw new Error(data.message || 'Failed to connect to Vault');
            }
        })
        .catch(error => {
            console.error('Vault connection error:', error);
            alert(`Error connecting to Vault: ${error.message}`);
        })
        .finally(() => {
            hideLoadingIndicator();
        });
    }
}