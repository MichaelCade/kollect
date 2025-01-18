// aws.js

document.getElementById('aws-button').addEventListener('click', () => {
    showLoadingIndicator();
    fetch('/api/switch?type=aws')
        .then(response => response.json())
        .then(data => {
            location.reload();
        })
        .catch(error => console.error('Error switching to AWS:', error))
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
            if (data.EC2Instances) {
                createTable('EC2 Instances', data.EC2Instances, ec2InstanceRowTemplate, ['Name', 'Instance ID', 'Type', 'State', 'Region']);
            }
            if (data.S3Buckets) {
                createTable('S3 Buckets', data.S3Buckets, s3BucketRowTemplate, ['Bucket Name', 'Immutable', 'Region']);
            }
            if (data.RDSInstances) {
                createTable('RDS Instances', data.RDSInstances, rdsInstanceRowTemplate, ['Instance ID', 'Engine', 'Status', 'Region']);
            }
            if (data.DynamoDBTables) {
                createTable('DynamoDB Tables', data.DynamoDBTables, dynamoDBTableRowTemplate, ['Table Name', 'Status', 'Region']);
            }
            if (data.VPCs) {
                createTable('VPCs', data.VPCs, vpcRowTemplate, ['VPC ID', 'State', 'Region']);
            }
        } catch (error) {
            console.error("Error processing data:", error);
        }
    }
});

function ec2InstanceRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.InstanceID}</td><td>${item.Type}</td><td>${item.State}</td><td>${item.Region}</td>`;
}

function s3BucketRowTemplate(item) {
    return `<td>${item.Name}</td><td>${item.Immutable}</td><td>${item.Region}</td>`;
}

function rdsInstanceRowTemplate(item) {
    return `<td>${item.InstanceID}</td><td>${item.Engine}</td><td>${item.Status}</td><td>${item.Region}</td>`;
}

function dynamoDBTableRowTemplate(item) {
    return `<td>${item.TableName}</td><td>${item.Status}</td><td>${item.Region}</td>`;
}

function vpcRowTemplate(item) {
    return `<td>${item.VPCID}</td><td>${item.State}</td><td>${item.Region}</td>`;
}