// gcp.js

registerDataHandler('gcp', 
    function(data) {
        return data.ComputeInstances || data.GCSBuckets || data.CloudSQLInstances ||
               data.CloudRunServices || data.CloudFunctions;
    },
    function(data) {
        console.log("Processing GCP data");
        
        if (data.ComputeInstances) {
            createTable('Compute Instances', data.ComputeInstances, computeInstanceRowTemplate, 
                ['Name', 'Zone', 'Machine Type', 'Status', 'Project']);
        }
        
        if (data.GCSBuckets) {
            createTable('Cloud Storage Buckets', data.GCSBuckets, gcsBucketRowTemplate, 
                ['Name', 'Location', 'Storage Class', 'Retention Policy', 'Retention Duration', 'Project']);
        }
        
        if (data.CloudSQLInstances) {
            createTable('Cloud SQL Instances', data.CloudSQLInstances, cloudSQLInstanceRowTemplate, 
                ['Name', 'Database Version', 'Region', 'Tier', 'Status', 'Project']);
        }
        
        if (data.CloudRunServices) {
            createTable('Cloud Run Services', data.CloudRunServices, cloudRunServiceRowTemplate, 
                ['Name', 'Region', 'URL', 'Replicas', 'Container', 'Project']);
        }
        
        if (data.CloudFunctions) {
            createTable('Cloud Functions', data.CloudFunctions, cloudFunctionRowTemplate, 
                ['Name', 'Region', 'Runtime', 'Status', 'Entry Point', 'Available Memory', 'Project']);
        }
        
        setTimeout(() => {
            console.log(`Created GCP tables`);
        }, 100);
    }
);

function computeInstanceRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Zone}</td><td>${item.MachineType}</td><td>${item.Status}</td><td>${item.Project}</td>`;
}

function gcsBucketRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Location}</td><td>${item.StorageClass}</td><td>${item.RetentionPolicy}</td><td>${item.RetentionDuration}</td><td>${item.Project}</td>`;
}

function cloudSQLInstanceRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.DatabaseVersion}</td><td>${item.Region}</td><td>${item.Tier}</td><td>${item.Status}</td><td>${item.Project}</td>`;
}

function cloudRunServiceRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Region}</td><td><a href="${item.URL}" target="_blank">${item.URL}</a></td><td>${item.Replicas}</td><td>${item.Container}</td><td>${item.Project}</td>`;
}

function cloudFunctionRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Region}</td><td>${item.Runtime}</td><td>${item.Status}</td><td>${item.EntryPoint}</td><td>${item.AvailableMemory}</td><td>${item.Project}</td>`;
}

// Update the DOMContentLoaded event handler
document.addEventListener('DOMContentLoaded', function() {
    console.log("DOM loaded - GCP module setting up event listener");
    
    const gcpButton = document.getElementById('google-button');
    if (gcpButton) {
        console.log("Found GCP button, setting up handler");
        
        // Clone the button to remove any existing event listeners
        const newButton = gcpButton.cloneNode(true);
        if (gcpButton.parentNode) {
            gcpButton.parentNode.replaceChild(newButton, gcpButton);
        }
        
        // Add our click handler that ALWAYS shows the form
        newButton.addEventListener('click', function(event) {
            console.log("GCP button clicked");
            event.preventDefault();
            
            // Always show the connection form, regardless of current connection state
            showGCPCredentialsModal();
        });
    } else {
        console.error("Could not find google-button element");
    }
});

function showGCPCredentialsModal() {
    console.log("Creating GCP credentials modal");
    const isConnected = document.getElementById('google-button')?.classList.contains('connected');
    
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
    modalContent.className = 'modal-content gcp-modal';
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
                    <i class="fas fa-info-circle"></i> You are already connected to Google Cloud. 
                    You can switch to another project or use different credentials.
                </p>
            </div>
        `;
    }
    
    modalContent.innerHTML = `
        <h3 style="margin-top: 0; color: var(--accent-color); font-size: 1.5em; border-bottom: 1px solid var(--border-color); padding-bottom: 10px;">
            <i class="fab fa-google"></i> Connect to Google Cloud
        </h3>
        
        ${connectionNote}
        
        <div class="gcp-connection-form" style="margin-top: 20px;">
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="gcp-default-config" name="gcp-config-source" value="default" checked>
                <label for="gcp-default-config" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-cog"></i> Use gcloud CLI Configuration
                </label>
                <div id="gcp-default-config-form" class="source-form" style="margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <p style="margin-top: 0;">This option uses credentials and settings from your gcloud CLI configuration.</p>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="gcp-project-selector" style="font-weight: bold; margin-bottom: 5px;">Select GCP Project:</label>
                        <select id="gcp-project-selector" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;">
                            <option value="">Loading projects...</option>
                        </select>
                    </div>
                </div>
            </div>
            
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="gcp-service-account" name="gcp-config-source" value="service_account">
                <label for="gcp-service-account" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-key"></i> Use Service Account Key File
                </label>
                <div id="gcp-service_account-config-form" class="source-form" style="display: none; margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <div class="form-group">
                        <button id="browse-gcp-key" class="btn" style="background-color: var(--accent-color); color: white; padding: 8px 15px; border-radius: 4px; border: none; cursor: pointer; font-weight: bold;">
                            <i class="fas fa-folder-open"></i> Browse Key File
                        </button>
                        <span id="selected-key-name" style="margin-left: 10px; font-style: italic;"></span>
                    </div>
                    <p class="tip" style="margin-top: 15px; font-size: 0.85em; color: var(--secondary-text-color); font-style: italic;">
                        Select a JSON key file downloaded from the GCP console. The key file is only used for authentication and is not stored.
                    </p>
                </div>
            </div>
            
            <div class="modal-buttons" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 25px; border-top: 1px solid var(--border-color); padding-top: 15px;">
                <button id="gcp-cancel-btn" class="btn" style="padding: 10px 20px; background-color: var(--button-bg-color); color: var(--button-text-color); border: none; border-radius: 4px; cursor: pointer;">Cancel</button>
                <button id="gcp-connect-btn" class="btn btn-primary" style="padding: 10px 20px; background-color: var(--accent-color); color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">
                    <i class="fas fa-plug"></i> Connect
                </button>
            </div>
        </div>
    `;
    
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    console.log("Modal created and added to DOM");

    const sourceForms = document.querySelectorAll('.source-form');
    document.querySelectorAll('input[name="gcp-config-source"]').forEach(radio => {
        radio.addEventListener('change', () => {
            console.log(`Radio changed to: ${radio.value}`);
            sourceForms.forEach(form => form.style.display = 'none');
            const selectedForm = document.getElementById(`gcp-${radio.value}-config-form`);
            if (selectedForm) {
                selectedForm.style.display = 'block';
            }
        });
    });

    loadGCPProjects();

    let selectedKeyFile = null;
    document.getElementById('browse-gcp-key')?.addEventListener('click', () => {
        console.log("Browse GCP key file button clicked");
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = '.json';
        input.onchange = (event) => {
            if (event.target.files.length > 0) {
                selectedKeyFile = event.target.files[0];
                document.getElementById('selected-key-name').textContent = selectedKeyFile.name;
                console.log(`Selected key file: ${selectedKeyFile.name}`);
            }
        };
        input.click();
    });

    document.getElementById('gcp-cancel-btn').addEventListener('click', () => {
        console.log("Cancel button clicked");
        modal.remove();
    });

    document.getElementById('gcp-connect-btn').addEventListener('click', () => {
        console.log("Connect button clicked");
        
        const configSource = document.querySelector('input[name="gcp-config-source"]:checked').value;
        console.log(`Selected source: ${configSource}`);
        
        if (configSource === 'default') {
            const project = document.getElementById('gcp-project-selector').value;
            if (!project) {
                alert('Please select a GCP project');
                return;
            }
            console.log(`Using GCP project: ${project}`);
            connectToGCP({ type: 'gcloud', project: project });
        } else {
            if (!selectedKeyFile) {
                alert('Please select a service account key file');
                return;
            }
            
            const reader = new FileReader();
            reader.onload = (e) => {
                try {
                    const keyData = JSON.parse(e.target.result);
                    console.log("Using service account key");
                    connectToGCP({ 
                        type: 'service_account', 
                        keyData: keyData
                    });
                } catch (error) {
                    console.error('Error parsing key file:', error);
                    alert('Invalid service account key file. Please make sure it is a valid JSON file.');
                }
            };
            reader.readAsText(selectedKeyFile);
        }
    });

    function loadGCPProjects() {
        console.log("Loading GCP projects");
        
        fetch('/api/gcp/projects')
            .then(response => {
                if (!response.ok) {
                    return response.text().then(text => {
                        throw new Error(`HTTP error ${response.status}: ${text}`);
                    });
                }
                return response.json();
            })
            .then(data => {
                console.log("Projects loaded:", data);
                const projectSelector = document.getElementById('gcp-project-selector');
                projectSelector.innerHTML = ''; 
                
                if (data.projects && data.projects.length > 0) {
                    data.projects.forEach(project => {
                        const option = document.createElement('option');
                        option.value = project.id;
                        option.textContent = `${project.name} (${project.id})`;
                        
                        if (project.isDefault) {
                            option.selected = true;
                        }
                        
                        projectSelector.appendChild(option);
                    });
                } else {
                    const option = document.createElement('option');
                    option.value = "";
                    option.textContent = "No projects found";
                    projectSelector.appendChild(option);
                }
            })
            .catch(error => {
                console.error("Error loading GCP projects:", error);
                const projectSelector = document.getElementById('gcp-project-selector');
                projectSelector.innerHTML = '<option value="">Error loading projects</option>';
            });
    }

    function connectToGCP(config) {
    console.log("Connecting to GCP with config:", config);
    showLoadingIndicator();
    
    fetch('/api/gcp/connect', {
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
            const button = document.getElementById('google-button');
            if (button) {
                button.classList.add('connected');
                button.classList.remove('not-connected');
                
                const existingBadges = button.querySelectorAll('.connection-badge');
                existingBadges.forEach(badge => badge.remove());
                
                const badge = document.createElement('span');
                badge.className = 'connection-badge connected';
                button.appendChild(badge);
                
                button.title = 'GCP (Connected)';
                console.log("Button updated to connected state");
            }
            
            modal.remove();
            console.log("Modal removed, reloading page");
            setTimeout(() => {
                location.reload();
            }, 300);
        } else {
            throw new Error(data.message || 'Failed to connect to GCP');
        }
    })
    .catch(error => {
        console.error('GCP connection error:', error);
        alert(`Error connecting to GCP: ${error.message}`);
    })
    .finally(() => {
        hideLoadingIndicator();
    });
}
}