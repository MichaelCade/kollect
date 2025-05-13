// aws.js

registerDataHandler('aws', 
    function(data) {
        return data.EC2Instances || data.S3Buckets || data.RDSInstances || 
               data.DynamoDBTables || data.VPCs;
    },
    function(data) {
        console.log("Processing AWS data");
        
        if (data.EC2Instances) {
            createTable('EC2 Instances', data.EC2Instances, ec2InstanceRowTemplate, 
                ['Name', 'Instance ID', 'Type', 'State', 'Region']);
        }
        
        if (data.S3Buckets) {
            createTable('S3 Buckets', data.S3Buckets, s3BucketRowTemplate, 
                ['Bucket Name', 'Immutable', 'Region']);
        }
        
        if (data.RDSInstances) {
            createTable('RDS Instances', data.RDSInstances, rdsInstanceRowTemplate, 
                ['Instance ID', 'Engine', 'Status', 'Region']);
        }
        
        if (data.DynamoDBTables) {
            createTable('DynamoDB Tables', data.DynamoDBTables, dynamoDBTableRowTemplate, 
                ['Table Name', 'Status', 'Region']);
        }
        
        if (data.VPCs) {
            createTable('VPCs', data.VPCs, vpcRowTemplate, 
                ['VPC ID', 'State', 'Region']);
        }
        
        setTimeout(() => {
            console.log(`Created AWS tables`);
        }, 100);
    }
);

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