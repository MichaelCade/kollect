// terraform.js

registerDataHandler('terraform', 
    function(data) {
        return data.Resources || data.Outputs || data.Providers;
    },
    function(data) {
        console.log("Processing Terraform data");
        
        // Create the resources table
        if (data.Resources) {
            createTable('Terraform Resources', data.Resources, resourceRowTemplate, 
                ['Name', 'Type', 'Provider', 'Module', 'Status', 'Details']);
        }
        
        // Create the outputs table
        if (data.Outputs) {
            createTable('Terraform Outputs', data.Outputs, outputRowTemplate, 
                ['Name', 'Value', 'Type']);
        }
        
        // Create the providers table
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
    // Generate a unique ID for the resource
    const resourceId = `tf-resource-${item.Type}-${item.Name}`.replace(/[^a-zA-Z0-9-]/g, '-');
    
    // Format the provider name to be more readable
    const provider = item.Provider.replace(/^provider\[\"/g, '').replace(/\"\]$/g, '');
    
    // Format the module name to be more readable
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
    // For Terraform, we need to get a file from the user since we need a state file
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.tfstate,.json';
    input.onchange = (event) => {
        const file = event.target.files[0];
        if (file) {
            showLoadingIndicator();
            const reader = new FileReader();
            reader.onload = (e) => {
                try {
                    const data = JSON.parse(e.target.result);
                    
                    // Basic validation that this is a Terraform state file
                    if (!data.version || !data.resources) {
                        throw new Error("The selected file does not appear to be a valid Terraform state file");
                    }
                    
                    // Process the Terraform state file
                    const parsedData = parseTerraformState(data);
                    processWithHandler(parsedData);
                } catch (error) {
                    console.error("Error processing Terraform state file:", error);
                    document.getElementById('content').innerHTML = `
                        <div class="error-message">
                            <h2>Error Processing Terraform State</h2>
                            <p>${error.message}</p>
                        </div>
                    `;
                } finally {
                    hideLoadingIndicator();
                }
            };
            reader.readAsText(file);
        }
    };
    input.click();
});

// Parse a Terraform state file into the format expected by our handler
function parseTerraformState(stateData) {
    const result = {
        Resources: [],
        Outputs: [],
        Providers: []
    };
    
    // Parse resources
    if (stateData.resources && Array.isArray(stateData.resources)) {
        stateData.resources.forEach(res => {
            if (res.instances && Array.isArray(res.instances)) {
                res.instances.forEach(inst => {
                    // Extract attributes in a flattened format
                    const attributes = {};
                    if (inst.attributes) {
                        Object.entries(inst.attributes).forEach(([key, value]) => {
                            // Handle simple types only
                            if (typeof value === 'string' || 
                                typeof value === 'number' || 
                                typeof value === 'boolean') {
                                attributes[key] = String(value);
                            } else if (value === null) {
                                attributes[key] = "null";
                            } else {
                                // For complex types, just show the type
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
    
    // Parse outputs
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
    
    // Parse providers
    if (stateData.provider_hash) {
        Object.entries(stateData.provider_hash).forEach(([name, version]) => {
            // Clean up provider names
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