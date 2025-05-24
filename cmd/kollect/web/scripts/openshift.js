// openshift.js

console.log("Loading OpenShift module");

registerDataHandler('openshift', 
    function(data) {
        return data.projects || data.routes || data.clusterInfo || 
               data.buildConfigs || data.builds || data.imageStreams ||
               data.imageStreamTags || data.imageStreamImports || data.templates || 
               data.securityContextConstraints || data.clusterOperators || 
               data.clusterVersions || data.machineConfigs || data.machineConfigPools || 
               data.operatorGroups || data.subscriptions || data.installPlans || 
               data.consoleLinks || data.consoleNotifications || data.consoleCLIDownloads || 
               data.oAuthClients || data.oAuthAccessTokens || data.oAuthAuthorizeTokens ||
               data.ingressControllers || data.dnses;
    },
    function(data) {
        console.log("Processing OpenShift data");
        
        // Create the OpenShift cluster info section
        if (data.clusterInfo) {
            createOpenShiftClusterInfo(data.clusterInfo);
        }
        
        // Create Project overview
        if (data.projects && data.projects.length > 0) {
            createTable('Projects', data.projects, projectRowTemplate, 
                ['Name', 'Display Name', 'Status', 'Created', 'Actions']);
        }
        
        // Create Routes table
        if (data.routes && data.routes.length > 0) {
            createTable('Routes', data.routes, routeRowTemplate, 
                ['Name', 'Namespace', 'Host', 'Path', 'Service', 'TLS', 'Actions']);
        }
        
        // Create BuildConfig table
        if (data.buildConfigs && data.buildConfigs.length > 0) {
            createTable('Build Configs', data.buildConfigs, buildConfigRowTemplate, 
                ['Name', 'Namespace', 'Type', 'Source', 'Output', 'Triggers', 'Actions']);
        }
        
        // Create Builds table
        if (data.builds && data.builds.length > 0) {
            createTable('Builds', data.builds, buildRowTemplate, 
                ['Name', 'Namespace', 'Status', 'Phase', 'Started', 'Duration', 'Actions']);
        }
        
        // Create Image Streams table
        if (data.imageStreams && data.imageStreams.length > 0) {
            createTable('Image Streams', data.imageStreams, imageStreamRowTemplate, 
                ['Name', 'Namespace', 'Tags', 'Updated', 'Actions']);
        }
        
        // Create Image Stream Tags table
        if (data.imageStreamTags && data.imageStreamTags.length > 0) {
            createTable('Image Stream Tags', data.imageStreamTags, imageStreamTagRowTemplate, 
                ['Name', 'Namespace', 'From', 'Created', 'Actions']);
        }
        
        // Create Image Stream Imports table
        if (data.imageStreamImports && data.imageStreamImports.length > 0) {
            createTable('Image Stream Imports', data.imageStreamImports, imageStreamImportRowTemplate, 
                ['Name', 'Namespace', 'Status', 'Images', 'Actions']);
        }
        
        // Create Templates table
        if (data.templates && data.templates.length > 0) {
            createTable('Templates', data.templates, templateRowTemplate, 
                ['Name', 'Namespace', 'Objects', 'Parameters', 'Actions']);
        }
        
        // Create Security Context Constraints table
        if (data.securityContextConstraints && data.securityContextConstraints.length > 0) {
            createTable('Security Context Constraints', data.securityContextConstraints, sccRowTemplate, 
                ['Name', 'Priority', 'Allowed Capabilities', 'Users', 'Groups', 'Actions']);
        }
        
        // Create Cluster Operators table
        if (data.clusterOperators && data.clusterOperators.length > 0) {
            createTable('Cluster Operators', data.clusterOperators, clusterOperatorRowTemplate, 
                ['Name', 'Version', 'Status', 'Message', 'Actions']);
        }
        
        // Create Cluster Versions table
        if (data.clusterVersions && data.clusterVersions.length > 0) {
            createTable('Cluster Versions', data.clusterVersions, clusterVersionRowTemplate, 
                ['Version', 'Available Updates', 'Channel', 'Status', 'Actions']);
        }
        
        // Create Machine Configs table
        if (data.machineConfigs && data.machineConfigs.length > 0) {
            createTable('Machine Configs', data.machineConfigs, machineConfigRowTemplate, 
                ['Name', 'Created', 'Generation', 'Actions']);
        }
        
        // Create Machine Config Pools table
        if (data.machineConfigPools && data.machineConfigPools.length > 0) {
            createTable('Machine Config Pools', data.machineConfigPools, machineConfigPoolRowTemplate, 
                ['Name', 'Configuration', 'Updated', 'Ready/Total', 'Actions']);
        }
        
        // Create Operator Groups table
        if (data.operatorGroups && data.operatorGroups.length > 0) {
            createTable('Operator Groups', data.operatorGroups, operatorGroupRowTemplate, 
                ['Name', 'Namespace', 'Target Namespaces', 'Created', 'Actions']);
        }
        
        // Create Subscriptions table
        if (data.subscriptions && data.subscriptions.length > 0) {
            createTable('Subscriptions', data.subscriptions, subscriptionRowTemplate, 
                ['Name', 'Namespace', 'Package', 'Channel', 'Source', 'Status', 'Actions']);
        }
        
        // Create InstallPlans table
        if (data.installPlans && data.installPlans.length > 0) {
            createTable('Install Plans', data.installPlans, installPlanRowTemplate, 
                ['Name', 'Namespace', 'Phase', 'Components', 'Approved', 'Actions']);
        }
        
        // Create Console Links table
        if (data.consoleLinks && data.consoleLinks.length > 0) {
            createTable('Console Links', data.consoleLinks, consoleLinkRowTemplate, 
                ['Name', 'Text', 'URL', 'Location', 'Actions']);
        }
        
        // Create Console Notifications table
        if (data.consoleNotifications && data.consoleNotifications.length > 0) {
            createTable('Console Notifications', data.consoleNotifications, consoleNotificationRowTemplate, 
                ['Name', 'Text', 'Location', 'Color', 'Actions']);
        }
        
        // Create Console CLI Downloads table
        if (data.consoleCLIDownloads && data.consoleCLIDownloads.length > 0) {
            createTable('Console CLI Downloads', data.consoleCLIDownloads, consoleCLIDownloadRowTemplate, 
                ['Name', 'Display Name', 'Description', 'Actions']);
        }
        
        // Create OAuth Clients table
        if (data.oAuthClients && data.oAuthClients.length > 0) {
            createTable('OAuth Clients', data.oAuthClients, oauthClientRowTemplate, 
                ['Name', 'Secret', 'Redirect URIs', 'Actions']);
        }
        
        // Create OAuth Access Tokens table
        if (data.oAuthAccessTokens && data.oAuthAccessTokens.length > 0) {
            createTable('OAuth Access Tokens', data.oAuthAccessTokens, oauthAccessTokenRowTemplate, 
                ['Name', 'User', 'Client', 'Expires', 'Actions']);
        }
        
        // Create OAuth Authorize Tokens table
        if (data.oAuthAuthorizeTokens && data.oAuthAuthorizeTokens.length > 0) {
            createTable('OAuth Authorize Tokens', data.oAuthAuthorizeTokens, oauthAuthorizeTokenRowTemplate, 
                ['Name', 'User', 'Client', 'Expires', 'Actions']);
        }
        
        // Create Ingress Controllers table
        if (data.ingressControllers && data.ingressControllers.length > 0) {
            createTable('Ingress Controllers', data.ingressControllers, ingressControllerRowTemplate, 
                ['Name', 'Namespace', 'Domain', 'Status', 'Actions']);
        }
        
        // Create DNS table
        if (data.dnses && data.dnses.length > 0) {
            createTable('DNS', data.dnses, dnsRowTemplate, 
                ['Name', 'Base Domain', 'Private Zone', 'Actions']);
        }
    }
);

// Row template functions for each resource type
function projectRowTemplate(project) {
    const name = project.metadata.name;
    const displayName = project.metadata.annotations && project.metadata.annotations['openshift.io/display-name'] || name;
    const status = project.status && project.status.phase || 'Unknown';
    const creationTimestamp = new Date(project.metadata.creationTimestamp).toLocaleString();
    
    const projectId = `project-${name}`;
    
    return `
        <td>${name}</td>
        <td>${displayName}</td>
        <td>${status}</td>
        <td>${creationTimestamp}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${projectId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${projectId}" style="display:none;" class="details-panel">
                <h4>Project Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Display Name:</td><td>${displayName}</td></tr>
                        <tr><td>Status:</td><td>${status}</td></tr>
                        <tr><td>Created:</td><td>${creationTimestamp}</td></tr>
                        <tr><td>Description:</td><td>${project.metadata.annotations && project.metadata.annotations['openshift.io/description'] || 'N/A'}</td></tr>
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Annotations</h5>
                    <table class="nested-table">
                        ${renderAnnotations(project.metadata.annotations)}
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(project.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function routeRowTemplate(route) {
    const name = route.metadata.name;
    const namespace = route.metadata.namespace;
    const host = route.spec.host || 'N/A';
    const path = route.spec.path || '/';
    const service = route.spec.to && route.spec.to.name || 'N/A';
    const tls = route.spec.tls ? 'Yes' : 'No';
    
    const routeId = `route-${namespace}-${name}`;
    
    return `
        <td>${name}</td>
        <td>${namespace}</td>
        <td>${host}</td>
        <td>${path}</td>
        <td>${service}</td>
        <td>${tls}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${routeId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${routeId}" style="display:none;" class="details-panel">
                <h4>Route Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Namespace:</td><td>${namespace}</td></tr>
                        <tr><td>Host:</td><td>${host}</td></tr>
                        <tr><td>Path:</td><td>${path}</td></tr>
                        <tr><td>Service:</td><td>${service}</td></tr>
                        <tr><td>TLS:</td><td>${tls}</td></tr>
                        <tr><td>Created:</td><td>${new Date(route.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                ${route.spec.tls ? `
                <div class="detail-section">
                    <h5>TLS Configuration</h5>
                    <table class="nested-table">
                        <tr><td>Termination:</td><td>${route.spec.tls.termination || 'N/A'}</td></tr>
                        <tr><td>Certificate:</td><td>${route.spec.tls.certificate ? 'Present' : 'None'}</td></tr>
                        <tr><td>Key:</td><td>${route.spec.tls.key ? 'Present' : 'None'}</td></tr>
                        <tr><td>CA Certificate:</td><td>${route.spec.tls.caCertificate ? 'Present' : 'None'}</td></tr>
                        <tr><td>Destination CA Certificate:</td><td>${route.spec.tls.destinationCACertificate ? 'Present' : 'None'}</td></tr>
                        <tr><td>Insecure Edge Termination Policy:</td><td>${route.spec.tls.insecureEdgeTerminationPolicy || 'N/A'}</td></tr>
                    </table>
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(route.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function buildConfigRowTemplate(buildConfig) {
    const name = buildConfig.metadata.name;
    const namespace = buildConfig.metadata.namespace;
    const sourceType = buildConfig.spec && buildConfig.spec.source ? Object.keys(buildConfig.spec.source)[0] || 'N/A' : 'N/A';
    const sourceValue = buildConfig.spec && buildConfig.spec.source && buildConfig.spec.source[sourceType] ? 
                        buildConfig.spec.source[sourceType].uri || 
                        buildConfig.spec.source[sourceType].dockerfile ||
                        'N/A' : 'N/A';
    const output = buildConfig.spec && buildConfig.spec.output && buildConfig.spec.output.to ? 
                  `${buildConfig.spec.output.to.kind}/${buildConfig.spec.output.to.name}` : 'N/A';
    const triggers = buildConfig.spec && buildConfig.spec.triggers ? buildConfig.spec.triggers.map(t => t.type).join(', ') : 'None';
    
    const buildConfigId = `buildconfig-${namespace}-${name}`;
    
    return `
        <td>${name}</td>
        <td>${namespace}</td>
        <td>${sourceType}</td>
        <td>${sourceValue}</td>
        <td>${output}</td>
        <td>${triggers}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${buildConfigId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${buildConfigId}" style="display:none;" class="details-panel">
                <h4>BuildConfig Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Namespace:</td><td>${namespace}</td></tr>
                        <tr><td>Created:</td><td>${new Date(buildConfig.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Source</h5>
                    <table class="nested-table">
                        ${renderBuildSource(buildConfig.spec.source)}
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Strategy</h5>
                    <table class="nested-table">
                        <tr><td>Type:</td><td>${buildConfig.spec && buildConfig.spec.strategy ? Object.keys(buildConfig.spec.strategy)[0] || 'N/A' : 'N/A'}</td></tr>
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Output</h5>
                    <table class="nested-table">
                        ${renderBuildOutput(buildConfig.spec.output)}
                    </table>
                </div>
                
                ${buildConfig.spec && buildConfig.spec.triggers && buildConfig.spec.triggers.length > 0 ? `
                <div class="detail-section">
                    <h5>Triggers</h5>
                    <table class="nested-table">
                        ${renderBuildTriggers(buildConfig.spec.triggers)}
                    </table>
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(buildConfig.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function buildRowTemplate(build) {
    const name = build.metadata.name;
    const namespace = build.metadata.namespace;
    const status = build.status && build.status.phase || 'Unknown';
    
    // Calculating duration
    let duration = 'N/A';
    if (build.status && build.status.startTimestamp) {
        const start = new Date(build.status.startTimestamp);
        let end;
        
        if (build.status.completionTimestamp) {
            end = new Date(build.status.completionTimestamp);
        } else {
            end = new Date();
        }
        
        const durationMs = end - start;
        duration = formatDuration(durationMs);
    }
    
    const buildId = `build-${namespace}-${name}`;
    
    return `
        <td>${name}</td>
        <td>${namespace}</td>
        <td>${status}</td>
        <td>${build.status && build.status.phase || 'Unknown'}</td>
        <td>${build.status && build.status.startTimestamp ? new Date(build.status.startTimestamp).toLocaleString() : 'N/A'}</td>
        <td>${duration}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${buildId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${buildId}" style="display:none;" class="details-panel">
                <h4>Build Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Namespace:</td><td>${namespace}</td></tr>
                        <tr><td>Status:</td><td>${status}</td></tr>
                        <tr><td>Created:</td><td>${new Date(build.metadata.creationTimestamp).toLocaleString()}</td></tr>
                        <tr><td>Started:</td><td>${build.status && build.status.startTimestamp ? new Date(build.status.startTimestamp).toLocaleString() : 'N/A'}</td></tr>
                        <tr><td>Completed:</td><td>${build.status && build.status.completionTimestamp ? new Date(build.status.completionTimestamp).toLocaleString() : 'N/A'}</td></tr>
                        <tr><td>Duration:</td><td>${duration}</td></tr>
                    </table>
                </div>
                
                ${build.spec && build.spec.source ? `
                <div class="detail-section">
                    <h5>Source</h5>
                    <table class="nested-table">
                        ${renderBuildSource(build.spec.source)}
                    </table>
                </div>
                ` : ''}
                
                ${build.status && build.status.logSnippet ? `
                <div class="detail-section">
                    <h5>Log Snippet</h5>
                    <pre>${build.status.logSnippet}</pre>
                </div>
                ` : ''}
                
                ${build.status && build.status.message ? `
                <div class="detail-section">
                    <h5>Status Message</h5>
                    <p>${build.status.message}</p>
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(build.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function imageStreamRowTemplate(imageStream) {
    const name = imageStream.metadata.name;
    const namespace = imageStream.metadata.namespace;
    const tagCount = imageStream.status && imageStream.status.tags ? imageStream.status.tags.length : 0;
    
    // Find latest update time
    let latestUpdate = null;
    if (imageStream.status && imageStream.status.tags) {
        imageStream.status.tags.forEach(tag => {
            if (tag.items && tag.items.length > 0) {
                const tagTime = new Date(tag.items[0].created);
                if (!latestUpdate || tagTime > latestUpdate) {
                    latestUpdate = tagTime;
                }
            }
        });
    }
    
    const latestUpdateStr = latestUpdate ? latestUpdate.toLocaleString() : 'N/A';
    
    const imageStreamId = `imagestream-${namespace}-${name}`;
    
    return `
        <td>${name}</td>
        <td>${namespace}</td>
        <td>${tagCount}</td>
        <td>${latestUpdateStr}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${imageStreamId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${imageStreamId}" style="display:none;" class="details-panel">
                <h4>ImageStream Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Namespace:</td><td>${namespace}</td></tr>
                        <tr><td>Created:</td><td>${new Date(imageStream.metadata.creationTimestamp).toLocaleString()}</td></tr>
                        <tr><td>Docker Repository:</td><td>${imageStream.status && imageStream.status.dockerImageRepository || 'N/A'}</td></tr>
                    </table>
                </div>
                
                ${imageStream.status && imageStream.status.tags && imageStream.status.tags.length > 0 ? `
                <div class="detail-section">
                    <h5>Tags (${imageStream.status.tags.length})</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Tag</th>
                            <th>Image</th>
                            <th>Created</th>
                        </tr>
                        ${renderImageStreamTags(imageStream.status.tags)}
                    </table>
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(imageStream.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function imageStreamTagRowTemplate(imageStreamTag) {
    const name = imageStreamTag.metadata.name;
    const namespace = imageStreamTag.metadata.namespace;
    const parts = name.split(':');
    const imageStreamName = parts[0];
    const tag = parts[1] || 'latest';
    
    let fromInfo = 'N/A';
    let createdTime = 'N/A';
    
    if (imageStreamTag.image && imageStreamTag.image.dockerImageMetadata) {
        if (imageStreamTag.image.dockerImageMetadata.Created) {
            createdTime = new Date(imageStreamTag.image.dockerImageMetadata.Created).toLocaleString();
        }
        
        if (imageStreamTag.image.dockerImageMetadata.Config && imageStreamTag.image.dockerImageMetadata.Config.Labels) {
            const labels = imageStreamTag.image.dockerImageMetadata.Config.Labels;
            if (labels['io.openshift.build.commit.id'] && labels['io.openshift.build.source-location']) {
                fromInfo = `Git ${labels['io.openshift.build.commit.id'].substring(0, 7)}`;
            }
        }
    } else if (imageStreamTag.tag && imageStreamTag.tag.from) {
        fromInfo = `${imageStreamTag.tag.from.kind}/${imageStreamTag.tag.from.name}`;
    }
    
    const imageStreamTagId = `imagestreamtag-${namespace}-${name.replace(':', '-')}`;
    
    return `
        <td>${name}</td>
        <td>${namespace}</td>
        <td>${fromInfo}</td>
        <td>${createdTime}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${imageStreamTagId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${imageStreamTagId}" style="display:none;" class="details-panel">
                <h4>ImageStreamTag Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Namespace:</td><td>${namespace}</td></tr>
                        <tr><td>Created:</td><td>${createdTime}</td></tr>
                        <tr><td>Image Stream:</td><td>${imageStreamName}</td></tr>
                        <tr><td>Tag:</td><td>${tag}</td></tr>
                    </table>
                </div>
                
                ${imageStreamTag.image ? `
                <div class="detail-section">
                    <h5>Image Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${imageStreamTag.image.name || 'N/A'}</td></tr>
                        <tr><td>Docker Image Reference:</td><td>${imageStreamTag.image.dockerImageReference || 'N/A'}</td></tr>
                        ${imageStreamTag.image.dockerImageMetadata ? `
                            <tr><td>Architecture:</td><td>${imageStreamTag.image.dockerImageMetadata.Architecture || 'N/A'}</td></tr>
                            <tr><td>Size:</td><td>${formatBytes(imageStreamTag.image.dockerImageMetadata.Size || 0)}</td></tr>
                        ` : ''}
                    </table>
                </div>
                ` : ''}
                
                ${imageStreamTag.image && imageStreamTag.image.dockerImageMetadata && imageStreamTag.image.dockerImageMetadata.Config ? `
                <div class="detail-section">
                    <h5>Image Config</h5>
                    <table class="nested-table">
                        <tr><td>User:</td><td>${imageStreamTag.image.dockerImageMetadata.Config.User || 'N/A'}</td></tr>
                        <tr><td>Working Dir:</td><td>${imageStreamTag.image.dockerImageMetadata.Config.WorkingDir || 'N/A'}</td></tr>
                        <tr><td>Exposed Ports:</td><td>${imageStreamTag.image.dockerImageMetadata.Config.ExposedPorts ? Object.keys(imageStreamTag.image.dockerImageMetadata.Config.ExposedPorts).join(', ') : 'None'}</td></tr>
                    </table>
                </div>
                ` : ''}
                
                ${imageStreamTag.image && imageStreamTag.image.dockerImageMetadata && imageStreamTag.image.dockerImageMetadata.Config && imageStreamTag.image.dockerImageMetadata.Config.Labels ? `
                <div class="detail-section">
                    <h5>Image Labels</h5>
                    <table class="nested-table">
                        ${renderDockerLabels(imageStreamTag.image.dockerImageMetadata.Config.Labels)}
                    </table>
                </div>
                ` : ''}
            </div>
        </td>
    `;
}

function imageStreamImportRowTemplate(imageStreamImport) {
    const name = imageStreamImport.metadata.name;
    const namespace = imageStreamImport.metadata.namespace;
    
    // Determine status
    let status = 'Complete';
    let imageCount = 0;
    
    if (imageStreamImport.status && imageStreamImport.status.images) {
        imageCount = imageStreamImport.status.images.length;
        
        // Check if any image failed to import
        const failedImages = imageStreamImport.status.images.filter(img => img.status && img.status.status === 'Failure');
        if (failedImages.length > 0) {
            status = 'Failed';
        }
    }
    
    const imageStreamImportId = `imagestreamimport-${namespace}-${name}`;
    
    return `
        <td>${name}</td>
        <td>${namespace}</td>
        <td>${status}</td>
        <td>${imageCount}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${imageStreamImportId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${imageStreamImportId}" style="display:none;" class="details-panel">
                <h4>ImageStreamImport Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Namespace:</td><td>${namespace}</td></tr>
                        <tr><td>Status:</td><td>${status}</td></tr>
                        <tr><td>Created:</td><td>${new Date(imageStreamImport.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                ${imageStreamImport.status && imageStreamImport.status.images ? `
                <div class="detail-section">
                    <h5>Images</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Name</th>
                            <th>Status</th>
                            <th>Reason</th>
                        </tr>
                        ${imageStreamImport.status.images.map((img, idx) => `
                            <tr>
                                <td>${imageStreamImport.spec.images && imageStreamImport.spec.images[idx] 
                                      ? imageStreamImport.spec.images[idx].from.name : `Image ${idx}`}</td>
                                <td>${img.status ? img.status.status : 'N/A'}</td>
                                <td>${img.status && img.status.message ? img.status.message : 'N/A'}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(imageStreamImport.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function templateRowTemplate(template) {
    const name = template.metadata.name;
    const namespace = template.metadata.namespace;
    const objectCount = template.objects ? template.objects.length : 0;
    const parameterCount = template.parameters ? template.parameters.length : 0;
    
    const templateId = `template-${namespace}-${name}`;
    
    return `
        <td>${name}</td>
        <td>${namespace}</td>
        <td>${objectCount}</td>
        <td>${parameterCount}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${templateId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${templateId}" style="display:none;" class="details-panel">
                <h4>Template Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Namespace:</td><td>${namespace}</td></tr>
                        <tr><td>Created:</td><td>${new Date(template.metadata.creationTimestamp).toLocaleString()}</td></tr>
                        <tr><td>Description:</td><td>${template.metadata.annotations ? template.metadata.annotations.description || 'N/A' : 'N/A'}</td></tr>
                    </table>
                </div>
                
                ${template.objects && template.objects.length > 0 ? `
                <div class="detail-section">
                    <h5>Objects (${template.objects.length})</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Kind</th>
                            <th>Name</th>
                        </tr>
                        ${template.objects.map(obj => `
                            <tr>
                                <td>${obj.kind}</td>
                                <td>${obj.metadata.name}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
                
                ${template.parameters && template.parameters.length > 0 ? `
                <div class="detail-section">
                    <h5>Parameters (${template.parameters.length})</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Name</th>
                            <th>Description</th>
                            <th>Default Value</th>
                        </tr>
                        ${template.parameters.map(param => `
                            <tr>
                                <td>${param.name}</td>
                                <td>${param.description || 'N/A'}</td>
                                <td>${param.value || 'N/A'}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(template.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function sccRowTemplate(scc) {
    const name = scc.metadata.name;
    const priority = scc.priority || 'N/A';
    const allowedCapabilities = scc.allowedCapabilities ? scc.allowedCapabilities.join(', ') : 'None';
    
    // Users and Groups
    const users = scc.users ? scc.users.join(', ') : 'None';
    const groups = scc.groups ? scc.groups.join(', ') : 'None';
    
    const sccId = `scc-${name}`;
    
    return `
        <td>${name}</td>
        <td>${priority}</td>
        <td>${allowedCapabilities}</td>
        <td>${users}</td>
        <td>${groups}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${sccId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${sccId}" style="display:none;" class="details-panel">
                <h4>Security Context Constraints Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Priority:</td><td>${priority}</td></tr>
                        <tr><td>Allow Privileged:</td><td>${scc.allowPrivilegedContainer === true ? 'Yes' : 'No'}</td></tr>
                        <tr><td>Default Add Capabilities:</td><td>${scc.defaultAddCapabilities ? scc.defaultAddCapabilities.join(', ') : 'None'}</td></tr>
                        <tr><td>Required Drop Capabilities:</td><td>${scc.requiredDropCapabilities ? scc.requiredDropCapabilities.join(', ') : 'None'}</td></tr>
                        <tr><td>Allowed Capabilities:</td><td>${allowedCapabilities}</td></tr>
                        <tr><td>Allow Host Network:</td><td>${scc.allowHostNetwork === true ? 'Yes' : 'No'}</td></tr>
                        <tr><td>Allow Host Ports:</td><td>${scc.allowHostPorts === true ? 'Yes' : 'No'}</td></tr>
                        <tr><td>Allow Host PID:</td><td>${scc.allowHostPID === true ? 'Yes' : 'No'}</td></tr>
                        <tr><td>Allow Host IPC:</td><td>${scc.allowHostIPC === true ? 'Yes' : 'No'}</td></tr>
                        <tr><td>Read Only Root Filesystem:</td><td>${scc.readOnlyRootFilesystem === true ? 'Yes' : 'No'}</td></tr>
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Run As User Strategy</h5>
                    <table class="nested-table">
                        <tr><td>Type:</td><td>${scc.runAsUser ? scc.runAsUser.type : 'N/A'}</td></tr>
                        ${scc.runAsUser && scc.runAsUser.type === 'MustRunAsRange' ? `
                            <tr><td>UID Range Min:</td><td>${scc.runAsUser.uidRangeMin || 'N/A'}</td></tr>
                            <tr><td>UID Range Max:</td><td>${scc.runAsUser.uidRangeMax || 'N/A'}</td></tr>
                        ` : ''}
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>SELinux Context Strategy</h5>
                    <table class="nested-table">
                        <tr><td>Type:</td><td>${scc.seLinuxContext ? scc.seLinuxContext.type : 'N/A'}</td></tr>
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Allowed Volume Types</h5>
                    <p>${scc.volumes ? scc.volumes.join(', ') : 'None'}</p>
                </div>
                
                <div class="detail-section">
                    <h5>Users and Groups</h5>
                    <table class="nested-table">
                        <tr><td>Users:</td><td>${users}</td></tr>
                        <tr><td>Groups:</td><td>${groups}</td></tr>
                    </table>
                </div>
            </div>
        </td>
    `;
}

function clusterOperatorRowTemplate(operator) {
    const name = operator.metadata.name;
    const version = operator.status && operator.status.versions ? 
                  operator.status.versions.find(v => v.name === 'operator') ?
                  operator.status.versions.find(v => v.name === 'operator').version :
                  'N/A' : 'N/A';
    
    // Determine status based on conditions
    let status = 'Unknown';
    let statusMessage = '';
    
    if (operator.status && operator.status.conditions) {
        // First check for Available=True
        const availableCond = operator.status.conditions.find(c => c.type === 'Available');
        const degradedCond = operator.status.conditions.find(c => c.type === 'Degraded');
        const progressingCond = operator.status.conditions.find(c => c.type === 'Progressing');
        
        if (availableCond && availableCond.status === 'True') {
            status = 'Available';
            statusMessage = availableCond.message || '';
        } 
        
        if (degradedCond && degradedCond.status === 'True') {
            status = 'Degraded';
            statusMessage = degradedCond.message || '';
        } else if (progressingCond && progressingCond.status === 'True') {
            status = 'Progressing';
            statusMessage = progressingCond.message || '';
        }
    }
    
    let statusClass = '';
    switch (status.toLowerCase()) {
        case 'available':
            statusClass = 'status-available';
            break;
        case 'progressing':
            statusClass = 'status-progressing';
            break;
        case 'degraded':
            statusClass = 'status-degraded';
            break;
        default:
            statusClass = 'status-unknown';
    }
    
    const operatorId = `operator-${name}`;
    
    return `
        <td>${name}</td>
        <td>${version}</td>
        <td><span class="status-badge ${statusClass}">${status}</span></td>
        <td>${statusMessage || 'N/A'}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${operatorId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${operatorId}" style="display:none;" class="details-panel">
                <h4>Cluster Operator Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Status:</td><td><span class="status-badge ${statusClass}">${status}</span></td></tr>
                        <tr><td>Message:</td><td>${statusMessage || 'N/A'}</td></tr>
                    </table>
                </div>
                
                ${operator.status && operator.status.versions ? `
                <div class="detail-section">
                    <h5>Versions</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Component</th>
                            <th>Version</th>
                        </tr>
                        ${operator.status.versions.map(v => `
                            <tr>
                                <td>${v.name}</td>
                                <td>${v.version}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
                
                ${operator.status && operator.status.relatedObjects ? `
                <div class="detail-section">
                    <h5>Related Objects</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Group</th>
                            <th>Resource</th>
                            <th>Name</th>
                            <th>Namespace</th>
                        </tr>
                        ${operator.status.relatedObjects.map(obj => `
                            <tr>
                                <td>${obj.group || 'core'}</td>
                                <td>${obj.resource}</td>
                                <td>${obj.name}</td>
                                <td>${obj.namespace || '-'}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
                
                ${operator.status && operator.status.conditions ? `
                <div class="detail-section">
                    <h5>Conditions</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Type</th>
                            <th>Status</th>
                            <th>Last Transition</th>
                            <th>Message</th>
                        </tr>
                        ${operator.status.conditions.map(cond => `
                            <tr>
                                <td>${cond.type}</td>
                                <td>${cond.status}</td>
                                <td>${new Date(cond.lastTransitionTime).toLocaleString()}</td>
                                <td>${cond.message}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
            </div>
        </td>
    `;
}

function clusterVersionRowTemplate(version) {
    // Extract core info
    const versionNumber = version.status && version.status.desired ? version.status.desired.version : 'Unknown';
    const channel = version.spec && version.spec.channel || 'N/A';
    
    // Determine available updates
    let availableUpdates = 'None';
    if (version.status && version.status.availableUpdates && version.status.availableUpdates.length > 0) {
        availableUpdates = version.status.availableUpdates.map(u => u.version).join(', ');
    }
    
    // Determine status
    let status = 'Unknown';
    if (version.status && version.status.conditions) {
        const progressingCond = version.status.conditions.find(c => c.type === 'Progressing');
        const availableCond = version.status.conditions.find(c => c.type === 'Available');
        
        if (progressingCond && progressingCond.status === 'True') {
            status = 'Updating';
        } else if (availableCond && availableCond.status === 'True') {
            status = 'Available';
        }
    }
    
    const versionId = `clusterversion-${version.metadata.name}`;
    
    return `
        <td>${versionNumber}</td>
        <td>${availableUpdates}</td>
        <td>${channel}</td>
        <td>${status}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${versionId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${versionId}" style="display:none;" class="details-panel">
                <h4>Cluster Version Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${version.metadata.name}</td></tr>
                        <tr><td>Version:</td><td>${versionNumber}</td></tr>
                        <tr><td>Channel:</td><td>${channel}</td></tr>
                        <tr><td>Status:</td><td>${status}</td></tr>
                        <tr><td>Created:</td><td>${new Date(version.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                ${version.status && version.status.history && version.status.history.length > 0 ? `
                <div class="detail-section">
                    <h5>Update History</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Version</th>
                            <th>State</th>
                            <th>Started</th>
                            <th>Completed</th>
                        </tr>
                        ${version.status.history.map(h => `
                            <tr>
                                <td>${h.version}</td>
                                <td>${h.state}</td>
                                <td>${h.startedTime ? new Date(h.startedTime).toLocaleString() : 'N/A'}</td>
                                <td>${h.completionTime ? new Date(h.completionTime).toLocaleString() : 'N/A'}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
                
                ${version.status && version.status.conditions ? `
                <div class="detail-section">
                    <h5>Conditions</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Type</th>
                            <th>Status</th>
                            <th>Message</th>
                            <th>Last Updated</th>
                        </tr>
                        ${version.status.conditions.map(cond => `
                            <tr>
                                <td>${cond.type}</td>
                                <td>${cond.status}</td>
                                <td>${cond.message || '-'}</td>
                                <td>${cond.lastTransitionTime ? new Date(cond.lastTransitionTime).toLocaleString() : '-'}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
            </div>
        </td>
    `;
}

function machineConfigRowTemplate(machineConfig) {
    const name = machineConfig.metadata.name;
    const createdAt = new Date(machineConfig.metadata.creationTimestamp).toLocaleString();
    const generation = machineConfig.metadata.generation || 'N/A';
    
    const configId = `machineconfig-${name}`;
    
    return `
        <td>${name}</td>
        <td>${createdAt}</td>
        <td>${generation}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${configId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${configId}" style="display:none;" class="details-panel">
                <h4>Machine Config Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Created:</td><td>${createdAt}</td></tr>
                        <tr><td>Generation:</td><td>${generation}</td></tr>
                        <tr><td>OSIMAGEURL:</td><td>${machineConfig.spec && machineConfig.spec.osImageURL || 'N/A'}</td></tr>
                    </table>
                </div>
                
                ${machineConfig.spec && machineConfig.spec.config ? `
                <div class="detail-section">
                    <h5>Ignition Config</h5>
                    <p>Version: ${machineConfig.spec.config.ignition ? machineConfig.spec.config.ignition.version : 'N/A'}</p>
                    
                    ${machineConfig.spec.config.storage && machineConfig.spec.config.storage.files ? `
                        <h6>Files (${machineConfig.spec.config.storage.files.length})</h6>
                        <table class="nested-table">
                            <tr>
                                <th>Path</th>
                                <th>Mode</th>
                            </tr>
                            ${machineConfig.spec.config.storage.files.map(file => `
                                <tr>
                                    <td>${file.path}</td>
                                    <td>${file.mode || '-'}</td>
                                </tr>
                            `).join('')}
                        </table>
                    ` : ''}
                    
                    ${machineConfig.spec.config.systemd && machineConfig.spec.config.systemd.units ? `
                        <h6>Systemd Units (${machineConfig.spec.config.systemd.units.length})</h6>
                        <table class="nested-table">
                            <tr>
                                <th>Name</th>
                                <th>Enabled</th>
                            </tr>
                            ${machineConfig.spec.config.systemd.units.map(unit => `
                                <tr>
                                    <td>${unit.name}</td>
                                    <td>${unit.enabled === true ? 'Yes' : 'No'}</td>
                                </tr>
                            `).join('')}
                        </table>
                    ` : ''}
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(machineConfig.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function machineConfigPoolRowTemplate(pool) {
    const name = pool.metadata.name;
    const configuration = pool.status && pool.status.configuration && pool.status.configuration.name || 'N/A';
    
    let updated = 'No';
    const conditions = pool.status && pool.status.conditions || [];
    for (const condition of conditions) {
        if (condition.type === 'Updated' && condition.status === 'True') {
            updated = 'Yes';
            break;
        }
    }
    
    const readyCount = pool.status && pool.status.readyMachineCount || 0;
    const totalCount = pool.status && pool.status.machineCount || 0;
    const readyTotal = `${readyCount}/${totalCount}`;
    
    const poolId = `machineconfigpool-${name}`;
    
    return `
        <td>${name}</td>
        <td>${configuration}</td>
        <td>${updated}</td>
        <td>${readyTotal}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${poolId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${poolId}" style="display:none;" class="details-panel">
                <h4>Machine Config Pool Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Created:</td><td>${new Date(pool.metadata.creationTimestamp).toLocaleString()}</td></tr>
                        <tr><td>Configuration:</td><td>${configuration}</td></tr>
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Status</h5>
                    <table class="nested-table">
                        <tr><td>Machine Count:</td><td>${pool.status ? pool.status.machineCount : 'N/A'}</td></tr>
                        <tr><td>Updated Machine Count:</td><td>${pool.status ? pool.status.updatedMachineCount : 'N/A'}</td></tr>
                        <tr><td>Ready Machine Count:</td><td>${pool.status ? pool.status.readyMachineCount : 'N/A'}</td></tr>
                        <tr><td>Unavailable Machine Count:</td><td>${pool.status ? pool.status.unavailableMachineCount : 'N/A'}</td></tr>
                        <tr><td>Degraded Machine Count:</td><td>${pool.status ? pool.status.degradedMachineCount : 'N/A'}</td></tr>
                    </table>
                </div>
                
                ${pool.spec && pool.spec.machineConfigSelector ? `
                <div class="detail-section">
                    <h5>Machine Config Selector</h5>
                    <table class="nested-table">
                        ${renderSelectorLabels(pool.spec.machineConfigSelector.matchLabels)}
                    </table>
                </div>
                ` : ''}
                
                ${pool.spec && pool.spec.nodeSelector ? `
                <div class="detail-section">
                    <h5>Node Selector</h5>
                    <table class="nested-table">
                        ${renderSelectorLabels(pool.spec.nodeSelector.matchLabels)}
                    </table>
                </div>
                ` : ''}
                
                ${conditions.length > 0 ? `
                <div class="detail-section">
                    <h5>Conditions</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Type</th>
                            <th>Status</th>
                            <th>Last Transition</th>
                            <th>Message</th>
                        </tr>
                        ${conditions.map(cond => `
                            <tr>
                                <td>${cond.type}</td>
                                <td>${cond.status}</td>
                                <td>${new Date(cond.lastTransitionTime).toLocaleString()}</td>
                                <td>${cond.message || '-'}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(pool.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function operatorGroupRowTemplate(operatorGroup) {
    const name = operatorGroup.metadata.name;
    const namespace = operatorGroup.metadata.namespace;
    const targetNamespaces = operatorGroup.status && operatorGroup.status.namespaces ? 
                            operatorGroup.status.namespaces.join(', ') : 
                            (operatorGroup.spec && operatorGroup.spec.targetNamespaces ? operatorGroup.spec.targetNamespaces.join(', ') : 'All Namespaces');
    
    const operatorGroupId = `operatorgroup-${namespace}-${name}`;
    
    return `
        <td>${name}</td>
        <td>${namespace}</td>
        <td>${targetNamespaces}</td>
        <td>${new Date(operatorGroup.metadata.creationTimestamp).toLocaleString()}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${operatorGroupId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${operatorGroupId}" style="display:none;" class="details-panel">
                <h4>Operator Group Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Namespace:</td><td>${namespace}</td></tr>
                        <tr><td>Target Namespaces:</td><td>${targetNamespaces}</td></tr>
                        <tr><td>Created:</td><td>${new Date(operatorGroup.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(operatorGroup.metadata.labels)}
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Annotations</h5>
                    <table class="nested-table">
                        ${renderAnnotations(operatorGroup.metadata.annotations)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function subscriptionRowTemplate(subscription) {
    const name = subscription.metadata.name;
    const namespace = subscription.metadata.namespace;
    const packageName = subscription.spec.name || 'N/A';
    const channel = subscription.spec.channel || 'N/A';
    const source = subscription.spec.source || 'N/A';
    
    let status = 'N/A';
    if (subscription.status) {
        if (subscription.status.installedCSV) {
            status = 'Installed';
        } else if (subscription.status.state) {
            status = subscription.status.state;
        }
    }
    
    const subscriptionId = `subscription-${namespace}-${name}`;
    
    return `
        <td>${name}</td>
        <td>${namespace}</td>
        <td>${packageName}</td>
        <td>${channel}</td>
        <td>${source}</td>
        <td>${status}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${subscriptionId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${subscriptionId}" style="display:none;" class="details-panel">
                <h4>Subscription Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Namespace:</td><td>${namespace}</td></tr>
                        <tr><td>Package:</td><td>${packageName}</td></tr>
                        <tr><td>Channel:</td><td>${channel}</td></tr>
                        <tr><td>Source:</td><td>${source}</td></tr>
                        <tr><td>Status:</td><td>${status}</td></tr>
                        <tr><td>Created:</td><td>${new Date(subscription.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                ${subscription.status && subscription.status.installedCSV ? `
                <div class="detail-section">
                    <h5>Installed CSV</h5>
                    <p>${subscription.status.installedCSV}</p>
                </div>
                ` : ''}
                
                ${subscription.status && subscription.status.conditions && subscription.status.conditions.length > 0 ? `
                <div class="detail-section">
                    <h5>Conditions</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Type</th>
                            <th>Status</th>
                            <th>Message</th>
                            <th>Last Transition</th>
                        </tr>
                        ${subscription.status.conditions.map(cond => `
                            <tr>
                                <td>${cond.type}</td>
                                <td>${cond.status}</td>
                                <td>${cond.message || '-'}</td>
                                <td>${cond.lastTransitionTime ? new Date(cond.lastTransitionTime).toLocaleString() : '-'}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(subscription.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function installPlanRowTemplate(installPlan) {
    const name = installPlan.metadata.name;
    const namespace = installPlan.metadata.namespace;
    const phase = installPlan.status ? installPlan.status.phase : 'Unknown';
    
    let componentsCount = 0;
    if (installPlan.spec && installPlan.spec.clusterServiceVersionNames) {
        componentsCount = installPlan.spec.clusterServiceVersionNames.length;
    }
    
    const approved = installPlan.spec && installPlan.spec.approved === true ? 'Yes' : 'No';
    
    const installPlanId = `installplan-${namespace}-${name}`;
    
    return `
        <td>${name}</td>
        <td>${namespace}</td>
        <td>${phase}</td>
        <td>${componentsCount}</td>
        <td>${approved}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${installPlanId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${installPlanId}" style="display:none;" class="details-panel">
                <h4>Install Plan Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Namespace:</td><td>${namespace}</td></tr>
                        <tr><td>Phase:</td><td>${phase}</td></tr>
                        <tr><td>Approved:</td><td>${approved}</td></tr>
                        <tr><td>Created:</td><td>${new Date(installPlan.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                ${installPlan.spec && installPlan.spec.clusterServiceVersionNames && installPlan.spec.clusterServiceVersionNames.length > 0 ? `
                <div class="detail-section">
                    <h5>Cluster Service Versions</h5>
                    <ul>
                        ${installPlan.spec.clusterServiceVersionNames.map(csv => `<li>${csv}</li>`).join('')}
                    </ul>
                </div>
                ` : ''}
                
                ${installPlan.status && installPlan.status.plan && installPlan.status.plan.length > 0 ? `
                <div class="detail-section">
                    <h5>Plan Steps</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Resource</th>
                            <th>Kind</th>
                            <th>Status</th>
                        </tr>
                        ${installPlan.status.plan.map(step => `
                            <tr>
                                <td>${step.resource.name}</td>
                                <td>${step.resource.kind}</td>
                                <td>${step.status}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(installPlan.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function consoleLinkRowTemplate(link) {
    const name = link.metadata.name;
    const text = link.spec.text || 'N/A';
    const url = link.spec.href || 'N/A';
    const location = link.spec.location || 'N/A';
    
    const linkId = `consolelink-${name}`;
    
    return `
        <td>${name}</td>
        <td>${text}</td>
        <td>${url}</td>
        <td>${location}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${linkId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${linkId}" style="display:none;" class="details-panel">
                <h4>Console Link Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Text:</td><td>${text}</td></tr>
                        <tr><td>URL:</td><td><a href="${url}" target="_blank">${url}</a></td></tr>
                        <tr><td>Location:</td><td>${location}</td></tr>
                        <tr><td>Created:</td><td>${new Date(link.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(link.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function consoleNotificationRowTemplate(notification) {
    const name = notification.metadata.name;
    const text = notification.spec.text || 'N/A';
    const location = notification.spec.location || 'N/A';
    const color = notification.spec.color || 'N/A';
    
    const notificationId = `consolenotification-${name}`;
    
    return `
        <td>${name}</td>
        <td>${text}</td>
        <td>${location}</td>
        <td>${color}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${notificationId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${notificationId}" style="display:none;" class="details-panel">
                <h4>Console Notification Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Text:</td><td>${text}</td></tr>
                        <tr><td>Location:</td><td>${location}</td></tr>
                        <tr><td>Color:</td><td>${color}</td></tr>
                        <tr><td>Created:</td><td>${new Date(notification.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(notification.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function consoleCLIDownloadRowTemplate(download) {
    const name = download.metadata.name;
    const displayName = download.spec.displayName || name;
    const description = download.spec.description || 'N/A';
    
    const downloadId = `consoleclidownload-${name}`;
    
    return `
        <td>${name}</td>
        <td>${displayName}</td>
        <td>${description}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${downloadId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${downloadId}" style="display:none;" class="details-panel">
                <h4>Console CLI Download Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Display Name:</td><td>${displayName}</td></tr>
                        <tr><td>Description:</td><td>${description}</td></tr>
                        <tr><td>Created:</td><td>${new Date(download.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                ${download.spec.links && download.spec.links.length > 0 ? `
                <div class="detail-section">
                    <h5>Download Links</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Text</th>
                            <th>URL</th>
                        </tr>
                        ${download.spec.links.map(link => `
                            <tr>
                                <td>${link.text}</td>
                                <td><a href="${link.href}" target="_blank">${link.href}</a></td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(download.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function oauthClientRowTemplate(client) {
    const name = client.metadata.name;
    const secret = client.secret || 'Hidden';
    const redirectURIs = client.redirectURIs ? client.redirectURIs.join(', ') : 'N/A';
    
    const clientId = `oauthclient-${name}`;
    
    return `
        <td>${name}</td>
        <td>${secret}</td>
        <td>${redirectURIs}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${clientId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${clientId}" style="display:none;" class="details-panel">
                <h4>OAuth Client Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Secret:</td><td>${secret}</td></tr>
                        <tr><td>Created:</td><td>${new Date(client.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                ${client.redirectURIs && client.redirectURIs.length > 0 ? `
                <div class="detail-section">
                    <h5>Redirect URIs</h5>
                    <ul>
                        ${client.redirectURIs.map(uri => `<li>${uri}</li>`).join('')}
                    </ul>
                </div>
                ` : ''}
                
                ${client.scopeRestrictions && client.scopeRestrictions.length > 0 ? `
                <div class="detail-section">
                    <h5>Scope Restrictions</h5>
                    <table class="nested-table">
                        ${renderScopeRestrictions(client.scopeRestrictions)}
                    </table>
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(client.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

function oauthAccessTokenRowTemplate(token) {
    const name = token.metadata.name;
    const userName = token.userName || 'N/A';
    const clientName = token.clientName || 'N/A';
    
    // Format expiration time
    let expiresTime = 'N/A';
    if (token.expiresIn) {
        const expiresDate = new Date(token.metadata.creationTimestamp);
        expiresDate.setSeconds(expiresDate.getSeconds() + token.expiresIn);
        expiresTime = expiresDate.toLocaleString();
    }
    
    const tokenId = `oauthaccesstoken-${name}`;
    
    return `
        <td>${name}</td>
        <td>${userName}</td>
        <td>${clientName}</td>
        <td>${expiresTime}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${tokenId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${tokenId}" style="display:none;" class="details-panel">
                <h4>OAuth Access Token Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>User Name:</td><td>${userName}</td></tr>
                        <tr><td>Client Name:</td><td>${clientName}</td></tr>
                        <tr><td>Expires:</td><td>${expiresTime}</td></tr>
                        <tr><td>Created:</td><td>${new Date(token.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                ${token.scopes && token.scopes.length > 0 ? `
                <div class="detail-section">
                    <h5>Scopes</h5>
                    <ul>
                        ${token.scopes.map(scope => `<li>${scope}</li>`).join('')}
                    </ul>
                </div>
                ` : ''}
            </div>
        </td>
    `;
}

function oauthAuthorizeTokenRowTemplate(token) {
    const name = token.metadata.name;
    const userName = token.userName || 'N/A';
    const clientName = token.clientName || 'N/A';
    
    // Format expiration time
    let expiresTime = 'N/A';
    if (token.expiresIn) {
        const expiresDate = new Date(token.metadata.creationTimestamp);
        expiresDate.setSeconds(expiresDate.getSeconds() + token.expiresIn);
        expiresTime = expiresDate.toLocaleString();
    }
    
    const tokenId = `oauthauthorizetoken-${name}`;
    
    return `
        <td>${name}</td>
        <td>${userName}</td>
        <td>${clientName}</td>
        <td>${expiresTime}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${tokenId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${tokenId}" style="display:none;" class="details-panel">
                <h4>OAuth Authorize Token Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>User Name:</td><td>${userName}</td></tr>
                        <tr><td>Client Name:</td><td>${clientName}</td></tr>
                        <tr><td>Expires:</td><td>${expiresTime}</td></tr>
                        <tr><td>Created:</td><td>${new Date(token.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                ${token.scopes && token.scopes.length > 0 ? `
                <div class="detail-section">
                    <h5>Scopes</h5>
                    <ul>
                        ${token.scopes.map(scope => `<li>${scope}</li>`).join('')}
                    </ul>
                </div>
                ` : ''}
                
                ${token.redirectURI ? `
                <div class="detail-section">
                    <h5>Redirect URI</h5>
                    <p>${token.redirectURI}</p>
                </div>
                ` : ''}
            </div>
        </td>
    `;
}

function ingressControllerRowTemplate(controller) {
    const name = controller.metadata.name;
    const namespace = controller.metadata.namespace;
    const domain = controller.spec && controller.spec.domain || 'N/A';
    
    // Determine status
    let status = 'Unknown';
    if (controller.status && controller.status.conditions) {
        const availableCondition = controller.status.conditions.find(c => c.type === 'Available');
        if (availableCondition && availableCondition.status === 'True') {
            status = 'Available';
        } else {
            status = 'Not Available';
        }
    }
    
    const controllerId = `ingresscontroller-${namespace}-${name}`;
    
    return `
        <td>${name}</td>
        <td>${namespace}</td>
        <td>${domain}</td>
        <td>${status}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${controllerId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${controllerId}" style="display:none;" class="details-panel">
                <h4>Ingress Controller Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Namespace:</td><td>${namespace}</td></tr>
                        <tr><td>Domain:</td><td>${domain}</td></tr>
                        <tr><td>Status:</td><td>${status}</td></tr>
                        <tr><td>Created:</td><td>${new Date(controller.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                ${controller.spec && controller.spec.endpointPublishingStrategy ? `
                <div class="detail-section">
                    <h5>Endpoint Publishing Strategy</h5>
                    <table class="nested-table">
                        <tr><td>Type:</td><td>${controller.spec.endpointPublishingStrategy.type || 'N/A'}</td></tr>
                    </table>
                </div>
                ` : ''}
                
                ${controller.status && controller.status.conditions ? `
                <div class="detail-section">
                    <h5>Conditions</h5>
                    <table class="nested-table">
                        <tr>
                            <th>Type</th>
                            <th>Status</th>
                            <th>Message</th>
                            <th>Last Transition</th>
                        </tr>
                        ${controller.status.conditions.map(cond => `
                            <tr>
                                <td>${cond.type}</td>
                                <td>${cond.status}</td>
                                <td>${cond.message || '-'}</td>
                                <td>${cond.lastTransitionTime ? new Date(cond.lastTransitionTime).toLocaleString() : '-'}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
            </div>
        </td>
    `;
}

function dnsRowTemplate(dns) {
    const name = dns.metadata.name;
    const baseDomain = dns.spec && dns.spec.baseDomain || 'N/A';
    const privateZone = dns.spec && dns.spec.privateZone === true ? 'Yes' : 'No';
    
    const dnsId = `dns-${name}`;
    
    return `
        <td>${name}</td>
        <td>${baseDomain}</td>
        <td>${privateZone}</td>
        <td>
            <button class="details-button" onclick="toggleDetails('${dnsId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${dnsId}" style="display:none;" class="details-panel">
                <h4>DNS Details</h4>
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Name:</td><td>${name}</td></tr>
                        <tr><td>Base Domain:</td><td>${baseDomain}</td></tr>
                        <tr><td>Private Zone:</td><td>${privateZone}</td></tr>
                        <tr><td>Created:</td><td>${new Date(dns.metadata.creationTimestamp).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                ${dns.spec && dns.spec.publicZones ? `
                <div class="detail-section">
                    <h5>Public Zones</h5>
                    <table class="nested-table">
                        <tr>
                            <th>ID</th>
                            <th>Region</th>
                        </tr>
                        ${dns.spec.publicZones.map(zone => `
                            <tr>
                                <td>${zone.id || 'N/A'}</td>
                                <td>${zone.region || 'N/A'}</td>
                            </tr>
                        `).join('')}
                    </table>
                </div>
                ` : ''}
                
                <div class="detail-section">
                    <h5>Labels</h5>
                    <table class="nested-table">
                        ${renderLabels(dns.metadata.labels)}
                    </table>
                </div>
            </div>
        </td>
    `;
}

// Helper functions for rendering resources
function renderLabels(labels) {
    if (!labels || Object.keys(labels).length === 0) {
        return '<tr><td colspan="2">No labels</td></tr>';
    }
    
    return Object.entries(labels).map(([key, value]) => 
        `<tr><td>${key}</td><td>${value}</td></tr>`
    ).join('');
}

function renderAnnotations(annotations) {
    if (!annotations || Object.keys(annotations).length === 0) {
        return '<tr><td colspan="2">No annotations</td></tr>';
    }
    
    return Object.entries(annotations).map(([key, value]) => 
        `<tr><td>${key}</td><td>${value}</td></tr>`
    ).join('');
}

function renderBuildSource(source) {
    if (!source) return '<tr><td colspan="2">No source information</td></tr>';
    
    let result = [];
    
    if (source.git) {
        result.push(`<tr><td>Type:</td><td>Git</td></tr>`);
        result.push(`<tr><td>URL:</td><td>${source.git.uri}</td></tr>`);
        if (source.git.ref) result.push(`<tr><td>Reference:</td><td>${source.git.ref}</td></tr>`);
    } else if (source.binary) {
        result.push(`<tr><td>Type:</td><td>Binary</td></tr>`);
    } else if (source.dockerfile) {
        result.push(`<tr><td>Type:</td><td>Dockerfile</td></tr>`);
        result.push(`<tr><td>Content:</td><td><pre>${source.dockerfile}</pre></td></tr>`);
    }
    
    return result.join('') || '<tr><td colspan="2">No source information</td></tr>';
}

function renderBuildOutput(output) {
    if (!output) return '<tr><td colspan="2">No output information</td></tr>';
    
    let result = [];
    
    if (output.to) {
        result.push(`<tr><td>Type:</td><td>${output.to.kind}</td></tr>`);
        result.push(`<tr><td>Name:</td><td>${output.to.name}</td></tr>`);
        if (output.to.namespace) result.push(`<tr><td>Namespace:</td><td>${output.to.namespace}</td></tr>`);
    }
    
    return result.join('') || '<tr><td colspan="2">No output information</td></tr>';
}

function renderBuildTriggers(triggers) {
    if (!triggers || triggers.length === 0) return '<tr><td colspan="2">No triggers</td></tr>';
    
    return triggers.map(trigger => 
        `<tr><td>Type:</td><td>${trigger.type}</td></tr>`
    ).join('');
}

function renderImageStreamTags(tags) {
    if (!tags || tags.length === 0) return '<tr><td colspan="3">No tags</td></tr>';
    
    return tags.map(tag => {
        let created = 'N/A';
        if (tag.items && tag.items.length > 0) {
            created = new Date(tag.items[0].created).toLocaleString();
        }
        
        let image = 'N/A';
        if (tag.items && tag.items.length > 0) {
            image = tag.items[0].dockerImageReference || tag.items[0].image || 'N/A';
        }
        
        return `<tr><td>${tag.tag}</td><td>${image}</td><td>${created}</td></tr>`;
    }).join('');
}

function renderSelectorLabels(labels) {
    if (!labels || Object.keys(labels).length === 0) {
        return '<tr><td colspan="2">No selectors</td></tr>';
    }
    
    return Object.entries(labels).map(([key, value]) => 
        `<tr><td>${key}</td><td>${value}</td></tr>`
    ).join('');
}

function renderDockerLabels(labels) {
    if (!labels || Object.keys(labels).length === 0) {
        return '<tr><td colspan="2">No labels</td></tr>';
    }
    
    return Object.entries(labels).map(([key, value]) => 
        `<tr><td>${key}</td><td>${value}</td></tr>`
    ).join('');
}

function renderScopeRestrictions(restrictions) {
    if (!restrictions || restrictions.length === 0) {
        return '<tr><td colspan="2">No scope restrictions</td></tr>';
    }
    
    return restrictions.map(restriction => {
        if (restriction.literals) {
            return `<tr><td>Literals:</td><td>${restriction.literals.join(', ')}</td></tr>`;
        } else {
            return '<tr><td>Restriction:</td><td>Custom</td></tr>';
        }
    }).join('');
}

function formatDuration(milliseconds) {
    if (milliseconds < 1000) {
        return `${milliseconds}ms`;
    } else if (milliseconds < 60000) {
        return `${Math.floor(milliseconds / 1000)}s`;
    } else if (milliseconds < 3600000) {
        const minutes = Math.floor(milliseconds / 60000);
        const seconds = Math.floor((milliseconds % 60000) / 1000);
        return `${minutes}m ${seconds}s`;
    } else {
        const hours = Math.floor(milliseconds / 3600000);
        const minutes = Math.floor((milliseconds % 3600000) / 60000);
        return `${hours}h ${minutes}m`;
    }
}

function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

// Create OpenShift cluster info section
function createOpenShiftClusterInfo(clusterInfo) {
    const section = document.createElement('div');
    section.className = 'cluster-info-section';
    
    section.innerHTML = `
        <style>
            .cluster-info-section {
                background: var(--card-bg);
                padding: 20px;
                border-radius: 8px;
                margin-bottom: 20px;
                box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            }
            
            .cluster-header {
                display: flex;
                align-items: center;
                margin-bottom: 20px;
            }
            
            .cluster-header h2 {
                margin: 0;
                font-size: 1.8em;
                color: var(--accent-color);
            }
            
            .cluster-details {
                display: grid;
                grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
                gap: 20px;
            }
            
            .detail-card {
                background: var(--background-color);
                border: 1px solid var(--border-color);
                border-radius: 6px;
                padding: 15px;
            }
            
            .detail-card h3 {
                margin-top: 0;
                font-size: 1.2em;
                color: var(--secondary-text-color);
                border-bottom: 1px solid var(--border-color);
                padding-bottom: 8px;
                margin-bottom: 10px;
            }
            
            .detail-card .detail-row {
                display: flex;
                justify-content: space-between;
                margin-bottom: 8px;
            }
            
            .detail-card .label {
                font-weight: bold;
            }
            
            .status-available {
                color: #4CAF50;
            }
            
            .status-progressing {
                color: #2196F3;
            }
            
            .status-degraded {
                color: #F44336;
            }
        </style>
        
        <div class="cluster-header">
            <h2>OpenShift Cluster</h2>
        </div>
        
        <div class="cluster-details">
            <div class="detail-card">
                <h3>Cluster Overview</h3>
                <div class="detail-row">
                    <span class="label">Version:</span>
                    <span>${clusterInfo.version || 'Not Available'}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Cluster ID:</span>
                    <span>${clusterInfo.clusterId || 'Not Available'}</span>
                </div>
                <div class="detail-row">
                    <span class="label">API Server:</span>
                    <span>${clusterInfo.apiServerURL || 'Not Available'}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Platform:</span>
                    <span>${clusterInfo.platform || 'Not Available'}</span>
                </div>
            </div>
            
            <div class="detail-card">
                <h3>Resources</h3>
                <div class="detail-row">
                    <span class="label">Projects:</span>
                    <span>${clusterInfo.projectCount || 0}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Nodes:</span>
                    <span>${clusterInfo.nodeCount || 0}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Pods:</span>
                    <span>${clusterInfo.podCount || 0}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Storage Classes:</span>
                    <span>${clusterInfo.storageClassCount || 0}</span>
                </div>
            </div>
            
            ${clusterInfo.status ? `
            <div class="detail-card">
                <h3>Status</h3>
                <div class="detail-row">
                    <span class="label">Cluster Status:</span>
                    <span class="status-${clusterInfo.status.toLowerCase()}">${clusterInfo.status}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Available Operators:</span>
                    <span>${clusterInfo.availableOperators || 0}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Degraded Operators:</span>
                    <span>${clusterInfo.degradedOperators || 0}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Progressing Operators:</span>
                    <span>${clusterInfo.progressingOperators || 0}</span>
                </div>
            </div>
            ` : ''}
        </div>
    `;
    
    document.getElementById('content').prepend(section);
}

// Generic toggle function for all details panels
function toggleDetails(id) {
    const element = document.getElementById(id);
    if (element) {
        element.style.display = element.style.display === 'none' ? 'block' : 'none';
    }
}

// Add DOM event listener for the OpenShift connection button
document.addEventListener('DOMContentLoaded', function() {
    console.log("DOM loaded - OpenShift module setting up event listener");
    
    const openshiftButton = document.getElementById('openshift-button');
    if (openshiftButton) {
        console.log("Found OpenShift button, setting up handler");
        
        const newButton = openshiftButton.cloneNode(true);
        if (openshiftButton.parentNode) {
            openshiftButton.parentNode.replaceChild(newButton, openshiftButton);
        }
        
        newButton.addEventListener('click', function(event) {
            console.log("OpenShift button clicked");
            event.preventDefault();
            
            showOpenShiftConnectionModal();
        });
    } else {
        console.error("Could not find openshift-button element");
    }
});

function showOpenShiftConnectionModal() {
    console.log("Creating OpenShift connection modal");
    
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
    modalContent.className = 'modal-content openshift-modal';
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
            <i class="fas fa-dharmachakra"></i> Connect to OpenShift Cluster
        </h3>
        
        <div class="openshift-connection-form" style="margin-top: 20px;">
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="default-kubeconfig" name="openshift-source" value="default" checked>
                <label for="default-kubeconfig" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-home"></i> Default Kubeconfig
                </label>
                <div id="default-kubeconfig-form" class="source-form" style="margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <p style="margin-top: 0;">Using default kubeconfig at: <code>~/.kube/config</code></p>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="openshift-default-context" style="font-weight: bold; margin-bottom: 5px;">Select Context:</label>
                        <select id="openshift-default-context" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px; box-sizing: border-box;">
                            <option value="">Loading contexts...</option>
                        </select>
                    </div>
                </div>
            </div>
            
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="custom-kubeconfig" name="openshift-source" value="custom">
                <label for="custom-kubeconfig" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-file"></i> Custom Kubeconfig File
                </label>
                <div id="custom-kubeconfig-form" class="source-form" style="display: none; margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <div class="form-group">
                        <button id="browse-openshift-kubeconfig" class="btn" style="background-color: var(--accent-color); color: white; padding: 8px 15px; border-radius: 4px; box-sizing: border-box; border: none; cursor: pointer; font-weight: bold;">
                            <i class="fas fa-folder-open"></i> Browse File
                        </button>
                        <span id="selected-openshift-kubeconfig-name" style="margin-left: 10px; font-style: italic;"></span>
                    </div>
                    <div class="form-group" style="margin-top: 15px;">
                        <label for="openshift-custom-context" style="font-weight: bold; margin-bottom: 5px;">Select Context:</label>
                        <select id="openshift-custom-context" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px; box-sizing: border-box;" disabled>
                            <option value="">Select a kubeconfig file first</option>
                        </select>
                    </div>
                </div>
            </div>
            
            <div class="modal-buttons" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 25px; border-top: 1px solid var(--border-color); padding-top: 15px;">
                <button id="openshift-cancel-btn" class="btn" style="padding: 10px 20px; background-color: var(--button-bg-color); color: var(--button-text-color); border: none; border-radius: 4px; cursor: pointer;">Cancel</button>
                <button id="openshift-connect-btn" class="btn btn-primary" style="padding: 10px 20px; background-color: var(--accent-color); color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">
                    <i class="fas fa-plug"></i> Connect
                </button>
            </div>
        </div>
    `;
    
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    console.log("Modal created and added to DOM");

    const sourceForms = document.querySelectorAll('.source-form');
    document.querySelectorAll('input[name="openshift-source"]').forEach(radio => {
        radio.addEventListener('change', () => {
            console.log(`Radio changed to: ${radio.value}`);
            sourceForms.forEach(form => form.style.display = 'none');
            const selectedForm = document.getElementById(`${radio.value}-kubeconfig-form`);
            if (selectedForm) {
                selectedForm.style.display = 'block';
            }
        });
    });

    console.log("Loading contexts for default kubeconfig");
    loadKubeContexts(null, 'openshift-default-context');

    let selectedKubeconfigFile = null;
    document.getElementById('browse-openshift-kubeconfig')?.addEventListener('click', () => {
        console.log("Browse kubeconfig button clicked");
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = '.yaml,.yml,.conf,.config';
        input.onchange = (event) => {
            if (event.target.files.length > 0) {
                selectedKubeconfigFile = event.target.files[0];
                document.getElementById('selected-openshift-kubeconfig-name').textContent = selectedKubeconfigFile.name;
                console.log(`Selected file: ${selectedKubeconfigFile.name}`);
                
                const reader = new FileReader();
                reader.onload = () => {
                    try {
                        const content = reader.result;
                        if (content.includes('contexts:') && content.includes('clusters:') && content.includes('users:')) {
                            console.log("File appears to be a valid kubeconfig");
                            const customContextSelector = document.getElementById('openshift-custom-context');
                            customContextSelector.disabled = false;
                            
                            uploadKubeconfigAndGetContexts(selectedKubeconfigFile, 'openshift-custom-context');
                        } else {
                            console.error("Invalid kubeconfig file format");
                            document.getElementById('selected-openshift-kubeconfig-name').textContent = 'Invalid kubeconfig file selected';
                            document.getElementById('openshift-custom-context').disabled = true;
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

    document.getElementById('openshift-cancel-btn').addEventListener('click', () => {
        console.log("Cancel button clicked");
        modal.remove();
    });

    document.getElementById('openshift-connect-btn').addEventListener('click', () => {
        console.log("Connect button clicked");
        let kubeconfigPath = '';
        let selectedContext = '';
        
        const kubeSource = document.querySelector('input[name="openshift-source"]:checked').value;
        console.log(`Selected source: ${kubeSource}`);
        
        if (kubeSource === 'default') {
            kubeconfigPath = '';  
            selectedContext = document.getElementById('openshift-default-context').value;
            console.log(`Using default kubeconfig with context: ${selectedContext}`);
            connectToOpenShift(kubeconfigPath, selectedContext);
        } else {
            if (!selectedKubeconfigFile) {
                alert('Please select a kubeconfig file');
                return;
            }
            
            selectedContext = document.getElementById('openshift-custom-context').value;
            console.log(`Using custom kubeconfig with context: ${selectedContext}`);
            
            const formData = new FormData();
            formData.append('kubeconfig', selectedKubeconfigFile);
            
            console.log("Uploading kubeconfig file");
            fetch('/api/openshift/upload-kubeconfig', {
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
                    connectToOpenShift(kubeconfigPath, selectedContext);
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

    function connectToOpenShift(kubeconfigPath, context) {
        console.log(`Connecting to OpenShift with config path: ${kubeconfigPath} and context: ${context}`);
        showLoadingIndicator();
        
        fetch('/api/openshift/connect', {
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
                const button = document.getElementById('openshift-button');
                if (button) {
                    button.classList.add('connected');
                    button.classList.remove('not-connected');
                    
                    const existingBadges = button.querySelectorAll('.connection-badge');
                    existingBadges.forEach(badge => badge.remove());
                    
                    const badge = document.createElement('span');
                    badge.className = 'connection-badge connected';
                    button.appendChild(badge);
                    
                    button.title = 'OpenShift (Connected)';
                }
                
                modal.remove();
                console.log("Modal removed, reloading page");
                setTimeout(() => {
                    location.reload();
                }, 300);
            } else {
                throw new Error(data.message || 'Failed to connect to OpenShift');
            }
        })
        .catch(error => {
            console.error('OpenShift connection error:', error);
            alert(`Error connecting to OpenShift cluster: ${error.message}`);
        })
        .finally(() => {
            hideLoadingIndicator();
        });
    }

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
                contextSelector.innerHTML = ''; 
                
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