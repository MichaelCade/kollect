// kubernetes.js

document.getElementById('kubernetes-button').addEventListener('click', () => {
    showLoadingIndicator();
    fetch('/api/switch?type=kubernetes')
        .then(response => response.json())
        .then(data => {
            location.reload();
        })
        .catch(error => console.error('Error switching to Kubernetes:', error))
        .finally(() => hideLoadingIndicator());
});

document.addEventListener('htmx:afterSwap', (event) => {
    if (event.detail.target.id === 'hidden-content') {
        try {
            const data = JSON.parse(event.detail.xhr.responseText);
            console.log("Fetched Data:", data); // Log fetched data
            const content = document.getElementById('content');
            const template = document.getElementById('table-template').content;
            function createTable(headerText, data, rowTemplate, headers) {
                if (!data || data.length === 0) return; // Ensure data is not null or empty
                const table = template.cloneNode(true);
                table.querySelector('th').textContent = headerText;
                const thead = table.querySelector('thead');
                const headerRow = document.createElement('tr');
                headers.forEach(header => {
                    const th = document.createElement('th');
                    th.textContent = header;
                    headerRow.appendChild(th);
                });
                thead.appendChild(headerRow);
                const tbody = table.querySelector('tbody');
                data.forEach(item => {
                    const row = document.createElement('tr');
                    row.innerHTML = rowTemplate(item);
                    tbody.appendChild(row);
                });
                content.appendChild(table);
            }
            if (data.Nodes) {
                createTable('Nodes', data.Nodes, nodeRowTemplate, ['Name', 'Roles', 'Age', 'Version', 'OS-Image']);
            }
            if (data.Namespaces) {
                createTable('Namespaces', data.Namespaces, defaultRowTemplate, ['Namespace']);
            }
            if (data.Pods) {
                createTable('Pods', data.Pods, podRowTemplate, ['Pod', 'Namespace', 'Status']);
            }
            if (data.Deployments) {
                createTable('Deployments', data.Deployments, deploymentRowTemplate, ['Deployments', 'Namespace', 'Containers', 'Images']);
            }
            if (data.StatefulSets) {
                createTable('StatefulSets', data.StatefulSets, stsRowTemplate, ['StatefulSet', 'Namespace', 'Ready Replicas','Image']);
            }
            if (data.Services) {
                createTable('Services', data.Services, serviceRowTemplate, ['Service', 'Namespace', 'Type', 'Cluster IP', 'Ports']);
            }
            if (data.PersistentVolumes) {
                createTable('PersistentVolumes', data.PersistentVolumes, perVolRowTemplate, ['PersistentVolume', 'Capacity', 'Access Modes', 'Status', 'Claim', 'StorageClass', 'Volume Mode']);
            }
            if (data.PersistentVolumeClaims) {
                createTable('PersistentVolumeClaims', data.PersistentVolumeClaims, perVolClaimRowTemplate, ['PersistentVolumeClaim', 'Namespace', 'Status', 'Volume', 'Capacity', 'Access Mode', 'StorageClass']);
            }
            if (data.StorageClasses) {
                createTable('StorageClasses', data.StorageClasses, storageClassRowTemplate, ['StorageClass', 'Provisioner', 'Volume Expansion']);
            }
            if (data.VolumeSnapshotClasses) {
                createTable('VolumeSnapshotClasses', data.VolumeSnapshotClasses, volSnapshotClassRowTemplate, ['VolumeSnapshotClass', 'Driver']);
            }
            if (data.VolumeSnapshots) {
                createTable('VolumeSnapshots', data.VolumeSnapshots, volumeSnapshotRowTemplate, ['Name', 'Namespace', 'Volume', 'CreationTimestamp', 'RestoreSize', 'Status']);
            }
        } catch (error) {
            console.error("Error processing data:", error);
        }
    }
});

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