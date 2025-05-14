// kubernetes.js

console.log("Loading Kubernetes module");

registerDataHandler('kubernetes', 
    function(data) {
        return data.Nodes || data.Pods || data.Deployments || data.Services ||
               data.PersistentVolumes || data.StorageClasses;
    },
    function(data) {
        console.log("Processing Kubernetes data");
        
        if (data.Nodes) {
            createTable('Nodes', data.Nodes, nodeRowTemplate, 
                ['Name', 'Roles', 'Age', 'Version', 'OS-Image']);
        }
        
        if (data.Namespaces) {
            createTable('Namespaces', data.Namespaces, defaultRowTemplate, 
                ['Namespace']);
        }
        
        if (data.Pods) {
            createTable('Pods', data.Pods, podRowTemplate, 
                ['Pod', 'Namespace', 'Status']);
        }
        
        if (data.Deployments) {
            createTable('Deployments', data.Deployments, deploymentRowTemplate, 
                ['Deployments', 'Namespace', 'Containers', 'Images']);
        }
        
        if (data.StatefulSets) {
            createTable('StatefulSets', data.StatefulSets, stsRowTemplate, 
                ['StatefulSet', 'Namespace', 'Ready Replicas','Image']);
        }
        
        if (data.Services) {
            createTable('Services', data.Services, serviceRowTemplate, 
                ['Service', 'Namespace', 'Type', 'Cluster IP', 'Ports']);
        }
        
        if (data.PersistentVolumes) {
            createTable('PersistentVolumes', data.PersistentVolumes, perVolRowTemplate, 
                ['PersistentVolume', 'Capacity', 'Access Modes', 'Status', 'Claim', 'StorageClass', 'Volume Mode']);
        }
        
        if (data.PersistentVolumeClaims) {
            createTable('PersistentVolumeClaims', data.PersistentVolumeClaims, perVolClaimRowTemplate, 
                ['PersistentVolumeClaim', 'Namespace', 'Status', 'Volume', 'Capacity', 'Access Mode', 'StorageClass']);
        }
        
        if (data.StorageClasses) {
            createTable('StorageClasses', data.StorageClasses, storageClassRowTemplate, 
                ['StorageClass', 'Provisioner', 'Volume Expansion']);
        }
        
        if (data.VolumeSnapshotClasses) {
            createTable('VolumeSnapshotClasses', data.VolumeSnapshotClasses, volSnapshotClassRowTemplate, 
                ['VolumeSnapshotClass', 'Driver']);
        }
        
        if (data.VolumeSnapshots) {
            createTable('VolumeSnapshots', data.VolumeSnapshots, volumeSnapshotRowTemplate, 
                ['Name', 'Namespace', 'Volume', 'CreationTimestamp', 'RestoreSize', 'Status']);
        }
        
        if (data.CustomResourceDefs) {
            createTable('Custom Resource Definitions', data.CustomResourceDefs, crdRowTemplate, 
                ['Name', 'Group', 'Version', 'Kind', 'Scope', 'Age']);
        }
        
        if (data.VirtualMachines) {
            createTable('Virtual Machines', data.VirtualMachines, vmRowTemplate, 
                ['Name', 'Namespace', 'Status', 'Ready', 'Age', 'Run Strategy', 'CPU', 'Memory', 'Data Volumes']);
        }
        
        if (data.DataVolumes) {
            createTable('Data Volumes', data.DataVolumes, dataVolumeRowTemplate, 
                ['Name', 'Namespace', 'Phase', 'Size', 'Source Type', 'Source', 'Age']);
        }
        
        setTimeout(() => {
            console.log(`Created Kubernetes tables`);
        }, 100);
    }
);

function nodeRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Roles}</td><td>${item.Age}</td><td>${item.Version}</td><td>${item.OSImage}</td>`;
}

function defaultRowTemplate(item) {
    return `<td>${item}</td>`;
}

function podRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Namespace}</td><td>${item.Status}</td>`;
}

function deploymentRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Namespace}</td><td>${item.Containers.join(', ')}</td><td>${item.Images.join(', ')}</td>`;
}

function stsRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Namespace}</td><td>${item.ReadyReplicas}</td><td>${item.Image}</td>`;
}

function serviceRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Namespace}</td><td>${item.Type}</td><td>${item.ClusterIP}</td><td>${item.Ports}</td>`;
}

function perVolRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Capacity}</td><td>${item.AccessModes}</td><td>${item.Status}</td><td>${item.AssociatedClaim}</td><td>${item.StorageClass}</td><td>${item.VolumeMode}</td>`;
}

function perVolClaimRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Namespace}</td><td>${item.Status}</td><td>${item.Volume}</td><td>${item.Capacity}</td><td>${item.AccessMode}</td><td>${item.StorageClass}</td>`;
}

function storageClassRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Provisioner}</td><td>${item.VolumeExpansion}</td>`;
}

function volSnapshotClassRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Driver}</td>`;
}

function volumeSnapshotRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Namespace}</td><td>${item.Volume}</td><td>${item.CreationTimestamp}</td><td>${item.RestoreSize}</td><td>${item.Status}</td>`;
}

function crdRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Group}</td><td>${item.Version}</td><td>${item.Kind}</td><td>${item.Scope}</td><td>${item.Age}</td>`;
}

function vmRowTemplate(item) {
    const ready = item.Ready ? 'Yes' : 'No';
    const vmId = `vm-${item.Namespace}-${item.Name}`.replace(/[^a-zA-Z0-9-]/g, '-');
    
    return `
        <td>${item.Name}</td>
        <td>${item.Namespace}</td>
        <td>${item.Status}</td>
        <td>${ready}</td>
        <td>${item.Age}</td>
        <td>${item.RunStrategy}</td>
        <td>${item.CPU || 'N/A'}</td>
        <td>${item.Memory || 'N/A'}</td>
        <td>
            <button class="details-button" onclick="toggleVMDetails('${vmId}')">
                <i class="fas fa-info-circle"></i> Storage
            </button>
            <div id="${vmId}" style="display:none;" class="details-panel">
                <h4>Storage Volumes</h4>
                ${item.Storage && item.Storage.length > 0 ? 
                    `<ul>${item.Storage.map(vol => `<li>${vol}</li>`).join('')}</ul>` : 
                    '<p>No storage volumes found for this VM.</p>'}
                
                <p class="tip">To see detailed PVC information: Check the PersistentVolumeClaims table for volumes in namespace "${item.Namespace}"</p>
            </div>
        </td>
    `;
}

function toggleVMDetails(id) {
    const details = document.getElementById(id);
    if (details) {
        details.style.display = details.style.display === 'none' ? 'block' : 'none';
    }
}

function dataVolumeRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Namespace}</td><td>${item.Phase}</td><td>${item.Size}</td><td>${item.SourceType}</td><td>${item.SourceInfo}</td><td>${item.Age}</td>`;
}

document.addEventListener('DOMContentLoaded', function() {
    console.log("DOM loaded - Kubernetes module setting up event listener");
    
    // Get the kubernetes button
    const kubernetesButton = document.getElementById('kubernetes-button');
    if (kubernetesButton) {
        console.log("Found kubernetes button, setting up handler");
        
        // Clone the button to remove any existing event listeners
        const newButton = kubernetesButton.cloneNode(true);
        if (kubernetesButton.parentNode) {
            kubernetesButton.parentNode.replaceChild(newButton, kubernetesButton);
        }
        
        // Add our click handler that ALWAYS shows the form
        newButton.addEventListener('click', function(event) {
            console.log("Kubernetes button clicked");
            event.preventDefault();
            
            // Always show the connection modal to allow context selection
            showKubernetesConnectionModal();
        });
    } else {
        console.error("Could not find kubernetes-button element");
    }
});

// Function to show the Kubernetes connection modal
function showKubernetesConnectionModal() {
    console.log("Creating Kubernetes connection modal");
    
    // Create a modal dialog for Kubernetes connection
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
    modalContent.className = 'modal-content kubernetes-modal';
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
            <i class="fas fa-dharmachakra"></i> Connect to Kubernetes Cluster
        </h3>
        
        <div class="kubernetes-connection-form" style="margin-top: 20px;">
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="default-kubeconfig" name="kubeconfig-source" value="default" checked>
                <label for="default-kubeconfig" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-home"></i> Default Kubeconfig
                </label>
                <div id="default-kubeconfig-form" class="source-form" style="margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <p style="margin-top: 0;">Using default kubeconfig at: <code>~/.kube/config</code></p>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="default-context-selector" style="font-weight: bold; margin-bottom: 5px;">Select Context:</label>
                        <select id="default-context-selector" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;">
                            <option value="">Loading contexts...</option>
                        </select>
                    </div>
                </div>
            </div>
            
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="custom-kubeconfig" name="kubeconfig-source" value="custom">
                <label for="custom-kubeconfig" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-file"></i> Custom Kubeconfig File
                </label>
                <div id="custom-kubeconfig-form" class="source-form" style="display: none; margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <div class="form-group">
                        <button id="browse-kubeconfig" class="btn" style="background-color: var(--accent-color); color: white; padding: 8px 15px; border-radius: 4px; border: none; cursor: pointer; font-weight: bold;">
                            <i class="fas fa-folder-open"></i> Browse File
                        </button>
                        <span id="selected-kubeconfig-name" style="margin-left: 10px; font-style: italic;"></span>
                    </div>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="custom-context-selector" style="font-weight: bold; margin-bottom: 5px;">Select Context:</label>
                        <select id="custom-context-selector" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;" disabled>
                            <option value="">Select a kubeconfig file first</option>
                        </select>
                    </div>
                </div>
            </div>
            
            <div class="modal-buttons" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 25px; border-top: 1px solid var(--border-color); padding-top: 15px;">
                <button id="kubernetes-cancel-btn" class="btn" style="padding: 10px 20px; background-color: var(--button-bg-color); color: var(--button-text-color); border: none; border-radius: 4px; cursor: pointer;">Cancel</button>
                <button id="kubernetes-connect-btn" class="btn btn-primary" style="padding: 10px 20px; background-color: var(--accent-color); color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">
                    <i class="fas fa-plug"></i> Connect
                </button>
            </div>
        </div>
    `;
    
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    console.log("Modal created and added to DOM");

    // Handle radio button selection
    const sourceForms = document.querySelectorAll('.source-form');
    document.querySelectorAll('input[name="kubeconfig-source"]').forEach(radio => {
        radio.addEventListener('change', () => {
            console.log(`Radio changed to: ${radio.value}`);
            // Hide all forms
            sourceForms.forEach(form => form.style.display = 'none');
            // Show the selected form
            const selectedForm = document.getElementById(`${radio.value}-kubeconfig-form`);
            if (selectedForm) {
                selectedForm.style.display = 'block';
            }
        });
    });

    // Load contexts for default kubeconfig
    console.log("Loading contexts for default kubeconfig");
    loadKubeContexts(null, 'default-context-selector');

    // Handle file selection
    let selectedKubeconfigFile = null;
    document.getElementById('browse-kubeconfig')?.addEventListener('click', () => {
        console.log("Browse kubeconfig button clicked");
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = '.yaml,.yml,.conf,.config';
        input.onchange = (event) => {
            if (event.target.files.length > 0) {
                selectedKubeconfigFile = event.target.files[0];
                document.getElementById('selected-kubeconfig-name').textContent = selectedKubeconfigFile.name;
                console.log(`Selected file: ${selectedKubeconfigFile.name}`);
                
                // Create a FileReader to read the file and get its path
                const reader = new FileReader();
                reader.onload = () => {
                    // Check if this is a valid kubeconfig (basic check for context info)
                    try {
                        const content = reader.result;
                        if (content.includes('contexts:') && content.includes('clusters:') && content.includes('users:')) {
                            console.log("File appears to be a valid kubeconfig");
                            // Looks like a kubeconfig, load contexts
                            const customContextSelector = document.getElementById('custom-context-selector');
                            customContextSelector.disabled = false;
                            
                            // Upload the file and parse it server-side
                            uploadKubeconfigAndGetContexts(selectedKubeconfigFile, 'custom-context-selector');
                        } else {
                            console.error("Invalid kubeconfig file format");
                            document.getElementById('selected-kubeconfig-name').textContent = 'Invalid kubeconfig file selected';
                            document.getElementById('custom-context-selector').disabled = true;
                        }
                    } catch (error) {
                        console.error('Error parsing kubeconfig:', error);
                    }
                };
                reader.readAsText(selectedKubeconfigFile);
            }
        };
        input.click();
    });

    // Handle cancel button
    document.getElementById('kubernetes-cancel-btn').addEventListener('click', () => {
        console.log("Cancel button clicked");
        modal.remove();
    });

    // Handle connect button
    document.getElementById('kubernetes-connect-btn').addEventListener('click', () => {
        console.log("Connect button clicked");
        let kubeconfigPath = '';
        let selectedContext = '';
        
        const kubeSource = document.querySelector('input[name="kubeconfig-source"]:checked').value;
        console.log(`Selected source: ${kubeSource}`);
        
        if (kubeSource === 'default') {
            kubeconfigPath = ''; // Use default
            selectedContext = document.getElementById('default-context-selector').value;
            console.log(`Using default kubeconfig with context: ${selectedContext}`);
            connectToKubernetes(kubeconfigPath, selectedContext);
        } else {
            // For custom kubeconfig, we'd need to upload the file to the server
            if (!selectedKubeconfigFile) {
                alert('Please select a kubeconfig file');
                return;
            }
            
            selectedContext = document.getElementById('custom-context-selector').value;
            console.log(`Using custom kubeconfig with context: ${selectedContext}`);
            
            // Upload the kubeconfig file
            const formData = new FormData();
            formData.append('kubeconfig', selectedKubeconfigFile);
            
            console.log("Uploading kubeconfig file");
            // Send the file to the server
            fetch('/api/kubernetes/upload-kubeconfig', {
                method: 'POST',
                body: formData
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
                    kubeconfigPath = data.path;
                    console.log(`Kubeconfig uploaded to: ${kubeconfigPath}`);
                    connectToKubernetes(kubeconfigPath, selectedContext);
                } else {
                    throw new Error(data.message || 'Failed to upload kubeconfig');
                }
            })
            .catch(error => {
                console.error('Error uploading kubeconfig:', error);
                alert(`Error uploading kubeconfig: ${error.message}`);
            });
        }
    });

    // Helper function to connect to Kubernetes
    function connectToKubernetes(kubeconfigPath, context) {
        console.log(`Connecting to Kubernetes with path: ${kubeconfigPath || 'default'} and context: ${context || 'default'}`);
        showLoadingIndicator();
        
        fetch('/api/kubernetes/connect', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                kubeconfigPath: kubeconfigPath,
                context: context
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
            console.log("Connection response:", data);
            if (data.status === 'success') {
                // Update the button status
                const button = document.getElementById('kubernetes-button');
                if (button) {
                    button.classList.add('connected');
                    button.classList.remove('not-connected');
                    
                    const existingBadges = button.querySelectorAll('.connection-badge');
                    existingBadges.forEach(badge => badge.remove());
                    
                    const badge = document.createElement('span');
                    badge.className = 'connection-badge connected';
                    button.appendChild(badge);
                    
                    button.title = 'Kubernetes (Connected)';
                    console.log("Button updated to connected state");
                }
                
                modal.remove();
                console.log("Modal removed, reloading page");
                setTimeout(() => {
                    location.reload();
                }, 300);
            } else {
                throw new Error(data.message || 'Failed to connect to Kubernetes');
            }
        })
        .catch(error => {
            console.error('Kubernetes connection error:', error);
            alert(`Error connecting to Kubernetes cluster: ${error.message}`);
        })
        .finally(() => {
            hideLoadingIndicator();
        });
    }

    // Function to load kubeconfig contexts from server
    function loadKubeContexts(kubeconfigPath, selectId) {
        const url = kubeconfigPath ? 
            `/api/kubernetes/contexts?path=${encodeURIComponent(kubeconfigPath)}` : 
            '/api/kubernetes/contexts';
        
        console.log(`Loading contexts from: ${url}`);
        
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
                console.log("Contexts loaded:", data);
                const contextSelector = document.getElementById(selectId);
                contextSelector.innerHTML = ''; // Clear existing options
                
                if (data.contexts && data.contexts.length > 0) {
                    data.contexts.forEach(context => {
                        const option = document.createElement('option');
                        option.value = context.name;
                        option.textContent = `${context.name} (${context.cluster})`;
                        
                        if (context.current === "true") {
                            option.textContent += " (current)";
                            option.selected = true;
                        }
                        
                        contextSelector.appendChild(option);
                    });
                    console.log(`Added ${data.contexts.length} contexts to selector`);
                } else {
                    console.log("No contexts found");
                    const option = document.createElement('option');
                    option.value = "";
                    option.textContent = "No contexts found";
                    contextSelector.appendChild(option);
                }
            })
            .catch(error => {
                console.error("Error loading Kubernetes contexts:", error);
                const contextSelector = document.getElementById(selectId);
                contextSelector.innerHTML = '<option value="">Error loading contexts</option>';
            });
    }

    // Function to upload kubeconfig and get contexts
    function uploadKubeconfigAndGetContexts(file, selectId) {
        console.log(`Uploading kubeconfig file to get contexts for selector: ${selectId}`);
        const formData = new FormData();
        formData.append('kubeconfig', file);
        
        fetch('/api/kubernetes/upload-kubeconfig', {
            method: 'POST',
            body: formData
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
                console.log(`Kubeconfig uploaded to: ${data.path}`);
                // Now load the contexts from the uploaded file
                loadKubeContexts(data.path, selectId);
            } else {
                throw new Error(data.message || 'Failed to upload kubeconfig');
            }
        })
        .catch(error => {
            console.error("Error uploading kubeconfig:", error);
            const contextSelector = document.getElementById(selectId);
            contextSelector.innerHTML = '<option value="">Error loading contexts</option>';
        });
    }
}