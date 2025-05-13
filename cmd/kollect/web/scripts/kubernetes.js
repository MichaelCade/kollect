// kubernetes.js

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