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

document.addEventListener('DOMContentLoaded', function() {
    console.log("DOM loaded - AWS module setting up event listener");
    
    const awsButton = document.getElementById('aws-button');
    if (awsButton) {
        console.log("Found AWS button, setting up handler");
        
        const newButton = awsButton.cloneNode(true);
        if (awsButton.parentNode) {
            awsButton.parentNode.replaceChild(newButton, awsButton);
        }
        
        newButton.addEventListener('click', function(event) {
            console.log("AWS button clicked");
            event.preventDefault();
            
            showAWSCredentialsModal();
        });
    } else {
        console.error("Could not find aws-button element");
    }
});

function showAWSCredentialsModal() {
    console.log("Creating AWS credentials modal");
    const isConnected = document.getElementById('aws-button')?.classList.contains('connected');

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
    modalContent.className = 'modal-content aws-modal';
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
                    <i class="fas fa-info-circle"></i> You are already connected to AWS. 
                    You can switch to another account or use different credentials.
                </p>
            </div>
        `;
    }

    modalContent.innerHTML = `
        <h3 style="margin-top: 0; color: var(--accent-color); font-size: 1.5em; border-bottom: 1px solid var(--border-color); padding-bottom: 10px;">
            <i class="fab fa-aws"></i> Connect to AWS
        </h3>
        
        ${connectionNote}
        
        <div class="aws-connection-form" style="margin-top: 20px;">
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="aws-default-config" name="aws-config-source" value="default" checked>
                <label for="aws-default-config" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-cog"></i> Use AWS CLI Configuration
                </label>
                <div id="aws-default-config-form" class="source-form" style="margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <p style="margin-top: 0;">This option uses credentials from your AWS CLI configuration at <code>~/.aws/credentials</code> and <code>~/.aws/config</code>.</p>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="aws-profile-selector" style="font-weight: bold; margin-bottom: 5px;">Select AWS Profile:</label>
                        <select id="aws-profile-selector" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;">
                            <option value="default">default</option>
                        </select>
                    </div>
                </div>
            </div>
            
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="aws-manual-config" name="aws-config-source" value="manual">
                <label for="aws-manual-config" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-key"></i> Enter Credentials Manually
                </label>
                <div id="aws-manual-config-form" class="source-form" style="display: none; margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <div class="form-group">
                        <label for="aws-access-key" style="font-weight: bold; margin-bottom: 5px;">Access Key ID:</label>
                        <input type="text" id="aws-access-key" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;" placeholder="AKIAIOSFODNN7EXAMPLE">
                    </div>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="aws-secret-key" style="font-weight: bold; margin-bottom: 5px;">Secret Access Key:</label>
                        <input type="password" id="aws-secret-key" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;" placeholder="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY">
                    </div>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="aws-region" style="font-weight: bold; margin-bottom: 5px;">AWS Region:</label>
                        <input type="text" id="aws-region" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;" value="us-east-1" placeholder="us-east-1">
                    </div>
                    <p class="tip" style="margin-top: 15px; font-size: 0.85em; color: var(--secondary-text-color); font-style: italic;">
                    </p>
                </div>
            </div>
            
            <div class="modal-buttons" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 25px; border-top: 1px solid var(--border-color); padding-top: 15px;">
                <button id="aws-cancel-btn" class="btn" style="padding: 10px 20px; background-color: var(--button-bg-color); color: var(--button-text-color); border: none; border-radius: 4px; cursor: pointer;">Cancel</button>
                <button id="aws-connect-btn" class="btn btn-primary" style="padding: 10px 20px; background-color: var(--accent-color); color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">
                    <i class="fas fa-plug"></i> Connect
                </button>
            </div>
        </div>
    `;
    
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    console.log("Modal created and added to DOM");

    const sourceForms = document.querySelectorAll('.source-form');
    document.querySelectorAll('input[name="aws-config-source"]').forEach(radio => {
        radio.addEventListener('change', () => {
            console.log(`Radio changed to: ${radio.value}`);
            sourceForms.forEach(form => form.style.display = 'none');
            const selectedForm = document.getElementById(`aws-${radio.value}-config-form`);
            if (selectedForm) {
                selectedForm.style.display = 'block';
            }
        });
    });

    loadAWSProfiles();

    document.getElementById('aws-cancel-btn').addEventListener('click', () => {
        console.log("Cancel button clicked");
        modal.remove();
    });

    document.getElementById('aws-connect-btn').addEventListener('click', () => {
        console.log("Connect button clicked");
        
        const configSource = document.querySelector('input[name="aws-config-source"]:checked').value;
        console.log(`Selected source: ${configSource}`);
        
        if (configSource === 'default') {
            const profile = document.getElementById('aws-profile-selector').value;
            console.log(`Using AWS profile: ${profile}`);
            connectToAWS({ type: 'profile', profile: profile });
        } else {
            const accessKey = document.getElementById('aws-access-key').value;
            const secretKey = document.getElementById('aws-secret-key').value;
            const region = document.getElementById('aws-region').value;
            
            if (!accessKey || !secretKey) {
                alert('Please provide both Access Key and Secret Key');
                return;
            }
            
            console.log(`Using manual AWS credentials with region: ${region}`);
            connectToAWS({ 
                type: 'credentials', 
                accessKey: accessKey, 
                secretKey: secretKey, 
                region: region || 'us-east-1' 
            });
        }
    });

    function loadAWSProfiles() {
        console.log("Loading AWS profiles");
        
        fetch('/api/aws/profiles')
            .then(response => {
                if (!response.ok) {
                    return response.text().then(text => {
                        throw new Error(`HTTP error ${response.status}: ${text}`);
                    });
                }
                return response.json();
            })
            .then(data => {
                console.log("Profiles loaded:", data);
                const profileSelector = document.getElementById('aws-profile-selector');
                profileSelector.innerHTML = ''; 
                
                if (data.profiles && data.profiles.length > 0) {
                    data.profiles.forEach(profile => {
                        const option = document.createElement('option');
                        option.value = profile;
                        option.textContent = profile;
                        
                        if (profile === 'default') {
                            option.selected = true;
                        }
                        
                        profileSelector.appendChild(option);
                    });
                } else {
                    const option = document.createElement('option');
                    option.value = "default";
                    option.textContent = "default";
                    profileSelector.appendChild(option);
                }
            })
            .catch(error => {
                console.error("Error loading AWS profiles:", error);
                const profileSelector = document.getElementById('aws-profile-selector');
                profileSelector.innerHTML = '<option value="default">default</option>';
            });
    }

    function connectToAWS(config) {
        console.log("Connecting to AWS with config:", config);
        showLoadingIndicator();
        
        fetch('/api/aws/connect', {
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
                const button = document.getElementById('aws-button');
                if (button) {
                    button.classList.add('connected');
                    button.classList.remove('not-connected');
                    
                    const existingBadges = button.querySelectorAll('.connection-badge');
                    existingBadges.forEach(badge => badge.remove());
                    
                    const badge = document.createElement('span');
                    badge.className = 'connection-badge connected';
                    button.appendChild(badge);
                    
                    button.title = 'AWS (Connected)';
                    console.log("Button updated to connected state");
                }
                
                modal.remove();
                console.log("Modal removed, reloading page");
                setTimeout(() => {
                    location.reload();
                }, 300);
            } else {
                throw new Error(data.message || 'Failed to connect to AWS');
            }
        })
        .catch(error => {
            console.error('AWS connection error:', error);
            alert(`Error connecting to AWS: ${error.message}`);
        })
        .finally(() => {
            hideLoadingIndicator();
        });
    }
}