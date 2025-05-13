// gcp.js

document.getElementById('google-button').addEventListener('click', () => {
    showLoadingIndicator();
    fetch('/api/switch?type=gcp')  // Change 'google' to 'gcp' to match the Go code
        .then(response => response.json())
        .then(data => {
            location.reload();
        })
        .catch(error => {
            console.error('Error switching to GCP:', error);
            const content = document.getElementById('content');
            content.innerHTML = '<div class="error-message">Error loading GCP data. Please check the console for details.</div>';
        })
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
            if (data.ComputeInstances) {
                createTable('Compute Instances', data.ComputeInstances, computeInstanceRowTemplate, ['Name', 'Zone', 'Machine Type', 'Status', 'Project']);
            }
            if (data.GCSBuckets) {
                createTable('Cloud Storage Buckets', data.GCSBuckets, gcsBucketRowTemplate, ['Name', 'Location', 'Storage Class', 'Retention Policy', 'Retention Duration', 'Project']);
            }
            if (data.CloudSQLInstances) {
                createTable('Cloud SQL Instances', data.CloudSQLInstances, cloudSQLInstanceRowTemplate, ['Name', 'Database Version', 'Region', 'Tier', 'Status', 'Project']);
            }
            if (data.CloudRunServices) {
                createTable('Cloud Run Services', data.CloudRunServices, cloudRunServiceRowTemplate, ['Name', 'Region', 'URL', 'Replicas', 'Container', 'Project']);
            }
            if (data.CloudFunctions) {
                createTable('Cloud Functions', data.CloudFunctions, cloudFunctionRowTemplate, ['Name', 'Region', 'Runtime', 'Status', 'Entry Point', 'Available Memory', 'Project']);
            }
        } catch (error) {
            console.error("Error processing data:", error);
        }
    }
});

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