// terraform.js

registerDataHandler('terraform', 
    function(data) {
        return data.Resources || data.Outputs || data.Providers;
    },
    function(data) {
        console.log("Processing Terraform data");
        
        if (data.Resources) {
            createTable('Terraform Resources', data.Resources, resourceRowTemplate, 
                ['Name', 'Type', 'Provider', 'Module', 'Status', 'Details']);
        }
        
        if (data.Outputs) {
            createTable('Terraform Outputs', data.Outputs, outputRowTemplate, 
                ['Name', 'Value', 'Type']);
        }
        
        if (data.Providers) {
            createTable('Terraform Providers', data.Providers, providerRowTemplate, 
                ['Name', 'Version']);
        }
        
        setTimeout(() => {
            console.log(`Created Terraform tables`);
        }, 100);
    }
);

function resourceRowTemplate(item) {
    const resourceId = `tf-resource-${item.Type}-${item.Name}`.replace(/[^a-zA-Z0-9-]/g, '-');
    const provider = item.Provider.replace(/^provider\[\"/g, '').replace(/\"\]$/g, '');
    
    let module = item.Module || 'root';
    if (module !== 'root') {
        module = module.replace('module.', '');
    }
    
    return `
        <td>${item.Name}</td>
        <td>${item.Type}</td>
        <td>${provider}</td>
        <td>${module}</td>
        <td>${item.Status}</td>
        <td>
            <button class="details-button" onclick="toggleTerraformDetails('${resourceId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${resourceId}" style="display:none;" class="details-panel">
                <h4>Attributes</h4>
                <table class="attributes-table">
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>Value</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${Object.entries(item.Attributes || {}).map(([key, value]) => 
                            `<tr><td>${key}</td><td>${value}</td></tr>`
                        ).join('')}
                    </tbody>
                </table>
                
                ${item.Dependencies && item.Dependencies.length > 0 ? `
                    <h4>Dependencies</h4>
                    <ul>
                        ${item.Dependencies.map(dep => `<li>${dep}</li>`).join('')}
                    </ul>
                ` : ''}
            </div>
        </td>
    `;
}

function outputRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Value}</td><td>${item.Type}</td>`;
}

function providerRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Version}</td>`;
}

function toggleTerraformDetails(id) {
    const details = document.getElementById(id);
    if (details) {
        details.style.display = details.style.display === 'none' ? 'block' : 'none';
    }
}

document.getElementById('terraform-button')?.addEventListener('click', () => {
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
    modalContent.className = 'modal-content terraform-modal';
    modalContent.style.backgroundColor = 'var(--card-bg)';
    modalContent.style.color = 'var(--text-color)';
    modalContent.style.padding = '25px';
    modalContent.style.borderRadius = '8px';
    modalContent.style.maxWidth = '550px';
    modalContent.style.width = '90%';
    modalContent.style.boxShadow = '0 5px 20px rgba(0,0,0,0.4)';
    modalContent.style.border = '1px solid var(--border-color)';
    
    modalContent.innerHTML = `
        <h3 style="margin-top: 0; color: var(--accent-color); font-size: 1.5em; border-bottom: 1px solid var(--border-color); padding-bottom: 10px;">
            <i class="fas fa-server"></i> Select Terraform State Source
        </h3>
        <div style="margin: 20px 0;">
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="local-file" name="state-source" value="file" checked>
                <label for="local-file" style="font-weight: bold; font-size: 1.1em;"><i class="fas fa-file"></i> Local State File</label>
                <div id="file-form" class="source-form" style="margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <button id="browse-file" class="btn" style="background-color: var(--accent-color); color: white; padding: 8px 15px; border-radius: 4px; border: none; cursor: pointer; font-weight: bold;">
                        <i class="fas fa-folder-open"></i> Browse File
                    </button>
                    <span id="selected-file-name" style="margin-left: 10px; font-style: italic;"></span>
                </div>
            </div>
            
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="s3-source" name="state-source" value="s3">
                <label for="s3-source" style="font-weight: bold; font-size: 1.1em;"><i class="fab fa-aws"></i> AWS S3 Bucket</label>
                <div id="s3-form" class="source-form" style="display: none; margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <div class="form-group">
                        <label for="s3-bucket">Bucket:</label>
                        <input type="text" id="s3-bucket" placeholder="my-terraform-states" style="width: 100%; padding: 8px; background: var(--input-bg-color, var(--background-color)); color: var(--text-color); border: 1px solid var(--border-color);">
                    </div>
                    <div class="form-group">
                        <label for="s3-key">Key:</label>
                        <input type="text" id="s3-key" placeholder="env/prod/terraform.tfstate" style="width: 100%; padding: 8px; background: var(--input-bg-color, var(--background-color)); color: var(--text-color); border: 1px solid var(--border-color);">
                    </div>
                    <div class="form-group">
                        <label for="s3-region">Region:</label>
                        <input type="text" id="s3-region" value="us-east-1" style="width: 100%; padding: 8px; background: var(--input-bg-color, var(--background-color)); color: var(--text-color); border: 1px solid var(--border-color);">
                    </div>
                </div>
            </div>
            
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="azure-source" name="state-source" value="azure">
                <label for="azure-source" style="font-weight: bold; font-size: 1.1em;"><i class="fab fa-microsoft"></i> Azure Blob Storage</label>
                <div id="azure-form" class="source-form" style="display: none; margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <div class="form-group">
                        <label for="azure-account">Storage Account:</label>
                        <input type="text" id="azure-account" placeholder="mystorageaccount" style="width: 100%; padding: 8px; background: var(--input-bg-color, var(--background-color)); color: var(--text-color); border: 1px solid var(--border-color);">
                    </div>
                    <div class="form-group">
                        <label for="azure-container">Container:</label>
                        <input type="text" id="azure-container" placeholder="tfstate" style="width: 100%; padding: 8px; background: var(--input-bg-color, var(--background-color)); color: var(--text-color); border: 1px solid var(--border-color);">
                    </div>
                    <div class="form-group">
                        <label for="azure-blob">Blob Name:</label>
                        <input type="text" id="azure-blob" placeholder="prod.terraform.tfstate" style="width: 100%; padding: 8px; background: var(--input-bg-color, var(--background-color)); color: var(--text-color); border: 1px solid var(--border-color);">
                    </div>
                </div>
            </div>
            
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="gcs-source" name="state-source" value="gcs">
                <label for="gcs-source" style="font-weight: bold; font-size: 1.1em;"><i class="fab fa-google"></i> Google Cloud Storage</label>
                <div id="gcs-form" class="source-form" style="display: none; margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <div class="form-group">
                        <label for="gcs-bucket">Bucket:</label>
                        <input type="text" id="gcs-bucket" placeholder="my-tf-states" style="width: 100%; padding: 8px; background: var(--input-bg-color, var(--background-color)); color: var(--text-color); border: 1px solid var(--border-color);">
                    </div>
                    <div class="form-group">
                        <label for="gcs-object">Object:</label>
                        <input type="text" id="gcs-object" placeholder="prod/terraform.tfstate" style="width: 100%; padding: 8px; background: var(--input-bg-color, var(--background-color)); color: var(--text-color); border: 1px solid var(--border-color);">
                    </div>
                </div>
            </div>
        </div>
        <div class="modal-buttons" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 20px; border-top: 1px solid var(--border-color); padding-top: 15px;">
            <button id="cancel-btn" class="btn" style="padding: 10px 20px; background-color: var(--button-bg-color, #444); color: var(--text-color); border: none; border-radius: 4px; cursor: pointer;">Cancel</button>
            <button id="load-state-btn" class="btn btn-primary" style="padding: 10px 20px; background-color: var(--accent-color); color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">
                <i class="fas fa-download"></i> Load State
            </button>
        </div>
    `;
        
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    const sourceForms = document.querySelectorAll('.source-form');
    document.querySelectorAll('input[name="state-source"]').forEach(radio => {
        radio.addEventListener('change', () => {
            sourceForms.forEach(form => form.style.display = 'none');
            const selectedForm = document.getElementById(`${radio.value}-form`);
            if (selectedForm) {
                selectedForm.style.display = 'block';
            }
        });
    });
    
    let selectedFile = null;
    document.getElementById('browse-file')?.addEventListener('click', () => {
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = '.tfstate,.json';
        input.onchange = (event) => {
            if (event.target.files.length > 0) {
                selectedFile = event.target.files[0];
                document.getElementById('selected-file-name').textContent = selectedFile.name;
            }
        };
        input.click();
    });
    
    document.getElementById('cancel-btn').addEventListener('click', () => {
        modal.remove();
    });
    
    document.getElementById('load-state-btn').addEventListener('click', () => {
        const sourceType = document.querySelector('input[name="state-source"]:checked').value;
        showLoadingIndicator();
        
        try {
            switch (sourceType) {
                case 'file':
                    if (!selectedFile) {
                        alert('Please select a file');
                        hideLoadingIndicator();
                        return;
                    }
                    
                    const reader = new FileReader();
                    reader.onload = (e) => {
                        try {
                            const data = JSON.parse(e.target.result);
                            processTerraformState(data);
                            modal.remove();
                        } catch (error) {
                            showError(`Error processing file: ${error.message}`);
                        } finally {
                            hideLoadingIndicator();
                        }
                    };
                    reader.readAsText(selectedFile);
                    break;
                    
                case 's3':
                    const bucket = document.getElementById('s3-bucket').value;
                    const key = document.getElementById('s3-key').value;
                    const region = document.getElementById('s3-region').value;
                    
                    if (!bucket || !key) {
                        alert('Please provide both bucket and key');
                        hideLoadingIndicator();
                        return;
                    }
                    
                    fetch('/api/terraform/s3-state', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            bucket: bucket,
                            key: key,
                            region: region || 'us-east-1'
                        })
                    })
                    .then(response => {
                        if (!response.ok) throw new Error(`HTTP error: ${response.status}`);
                        return response.json();
                    })
                    .then(data => {
                        processWithHandler(data);
                        modal.remove();
                    })
                    .catch(error => showError(`Error fetching from S3: ${error.message}`))
                    .finally(() => hideLoadingIndicator());
                    break;
                    
                case 'azure':
                    const account = document.getElementById('azure-account').value;
                    const container = document.getElementById('azure-container').value;
                    const blob = document.getElementById('azure-blob').value;
                    
                    if (!account || !container || !blob) {
                        alert('Please provide all Azure Blob Storage details');
                        hideLoadingIndicator();
                        return;
                    }
                    
                    fetch('/api/terraform/azure-state', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            account: account,
                            container: container,
                            blob: blob
                        })
                    })
                    .then(response => {
                        if (!response.ok) throw new Error(`HTTP error: ${response.status}`);
                        return response.json();
                    })
                    .then(data => {
                        processWithHandler(data);
                        modal.remove();
                    })
                    .catch(error => showError(`Error fetching from Azure: ${error.message}`))
                    .finally(() => hideLoadingIndicator());
                    break;
                    
                case 'gcs':
                    const gcsBucket = document.getElementById('gcs-bucket').value;
                    const object = document.getElementById('gcs-object').value;
                    
                    if (!gcsBucket || !object) {
                        alert('Please provide both GCS bucket and object');
                        hideLoadingIndicator();
                        return;
                    }
                    
                    fetch('/api/terraform/gcs-state', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            bucket: gcsBucket,
                            object: object
                        })
                    })
                    .then(response => {
                        if (!response.ok) throw new Error(`HTTP error: ${response.status}`);
                        return response.json();
                    })
                    .then(data => {
                        processWithHandler(data);
                        modal.remove();
                    })
                    .catch(error => showError(`Error fetching from GCS: ${error.message}`))
                    .finally(() => hideLoadingIndicator());
                    break;
            }
        } catch (error) {
            hideLoadingIndicator();
            showError(`Error: ${error.message}`);
        }
    });
});

function showError(message) {
    document.getElementById('content').innerHTML = `
        <div class="error-message">
            <h2>Error Processing Terraform State</h2>
            <p>${message}</p>
        </div>
    `;
}

function processTerraformState(stateData) {
    showLoadingIndicator();
    try {
        if (!stateData.version) {
            throw new Error("The selected file does not appear to be a valid Terraform state file");
        }
        const parsedData = parseTerraformState(stateData);
        processWithHandler(parsedData);
    } catch (error) {
        showError(`Failed to process state file: ${error.message}`);
    } finally {
        hideLoadingIndicator();
    }
}

function parseTerraformState(stateData) {
    const result = {
        Resources: [],
        Outputs: [],
        Providers: []
    };
    
    if (stateData.resources && Array.isArray(stateData.resources)) {
        stateData.resources.forEach(res => {
            if (res.instances && Array.isArray(res.instances)) {
                res.instances.forEach(inst => {
                    const attributes = {};
                    if (inst.attributes) {
                        Object.entries(inst.attributes).forEach(([key, value]) => {
                            if (typeof value === 'string' || 
                                typeof value === 'number' || 
                                typeof value === 'boolean') {
                                attributes[key] = String(value);
                            } else if (value === null) {
                                attributes[key] = "null";
                            } else {
                                attributes[key] = `[${Array.isArray(value) ? 'array' : 'object'}]`;
                            }
                        });
                    }
                    
                    result.Resources.push({
                        Name: res.name,
                        Type: res.type,
                        Provider: res.provider || "",
                        Module: res.module || "root",
                        Mode: res.mode || "",
                        Attributes: attributes,
                        Dependencies: inst.dependencies || [],
                        Status: inst.status || "Created"
                    });
                });
            }
        });
    }
    
    if (stateData.outputs) {
        Object.entries(stateData.outputs).forEach(([name, output]) => {
            result.Outputs.push({
                Name: name,
                Value: typeof output.value === 'string' ? output.value : 
                       typeof output.value === 'number' || typeof output.value === 'boolean' ? 
                       String(output.value) : "[complex value]",
                Type: output.type || "unknown"
            });
        });
    }
    
    if (stateData.provider_hash) {
        Object.entries(stateData.provider_hash).forEach(([name, version]) => {
            let cleanName = name;
            if (cleanName.startsWith('provider.')) {
                cleanName = cleanName.substring(9);
            } else if (cleanName.startsWith('provider[')) {
                cleanName = cleanName.substring(9, cleanName.indexOf(']'));
            }
            
            result.Providers.push({
                Name: cleanName,
                Version: version
            });
        });
    }
    
    return result;
}