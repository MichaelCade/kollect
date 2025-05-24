// docker.js

console.log("Loading Docker module");

registerDataHandler('docker', 
    function(data) {
        return data.containers || data.images || data.volumes || data.networks || data.info;
    },
    function(data) {
        console.log("Processing Docker data");
        
        createContainerHealthSummary(data);
        
        createNetworkTopologyVisualization(data);
        
        if (data.info) {
            const infoCategories = {
                'Server': {
                    'Docker ID': data.info.ID,
                    'Name': data.info.Name, 
                    'Docker Version': data.info.DockerVersion,
                    'OS': data.info.OS,
                    'OS Type': data.info.OSType,
                    'Architecture': data.info.Architecture,
                    'Kernel Version': data.info.KernelVersion
                },
                'Containers': {
                    'Total': data.info.Containers,
                    'Running': data.info.ContainersRunning,
                    'Paused': data.info.ContainersPaused,
                    'Stopped': data.info.ContainersStopped
                },
                'Images': {
                    'Total': data.info.Images
                }
            };

            for (const [category, items] of Object.entries(infoCategories)) {
                createTable(`Docker ${category}`, [items], dockerInfoCategoryRowTemplate, 
                    ['Property', 'Value']);
            }
        }
        
        if (data.containers && data.containers.length > 0) {
            const runningContainers = data.containers.filter(c => c.state === "running");
            const stoppedContainers = data.containers.filter(c => c.state !== "running");
            
            if (runningContainers.length > 0) {
                createTable('Running Containers', runningContainers, containerRowTemplate, 
                    ['Name', 'Image', 'Status', 'Created', 'Ports', 'Networks', 'Actions']);
            }
            
            if (stoppedContainers.length > 0) {
                createTable('Stopped Containers', stoppedContainers, containerRowTemplate, 
                    ['Name', 'Image', 'Status', 'Created', 'Ports', 'Networks', 'Actions']);
            }
        }
        
        if (data.images && data.images.length > 0) {
            createTable('Images', data.images, imageRowTemplate, 
                ['Repository:Tag', 'ID', 'Created', 'Size', 'Actions']);
        }
        
        if (data.volumes && data.volumes.length > 0) {
            createTable('Volumes', data.volumes, volumeRowTemplate, 
                ['Name', 'Driver', 'Mountpoint', 'Scope', 'Actions']);
        }
        
        if (data.networks && data.networks.length > 0) {
            createTable('Networks', data.networks, networkRowTemplate, 
                ['Name', 'ID', 'Driver', 'Scope', 'Subnet', 'Actions']);
        }
        
        createResourceUtilizationCharts(data);
        
        setTimeout(() => {
            console.log(`Created Docker tables`);
        }, 100);
    }
);

function createContainerHealthSummary(data) {
    const containerStats = {
        running: 0,
        exited: 0,
        paused: 0,
        total: data.containers ? data.containers.length : 0
    };
    
    let totalCpuPercent = 0;
    let totalMemoryUsed = 0;
    let totalMemoryLimit = 0;
    let totalNetworkRx = 0;
    let totalNetworkTx = 0;
    
    if (data.containers) {
        data.containers.forEach(container => {
            if (container.state === 'running') containerStats.running++;
            else if (container.state === 'exited') containerStats.exited++;
            else if (container.state === 'paused') containerStats.paused++;
            
            if (data.stats && data.stats[container.id]) {
                const stats = data.stats[container.id];
                totalCpuPercent += stats.cpuPercentage || 0;
                totalMemoryUsed += stats.memoryUsage || 0;
                totalMemoryLimit = Math.max(totalMemoryLimit, stats.memoryLimit || 0);
                totalNetworkRx += stats.networkRx || 0;
                totalNetworkTx += stats.networkTx || 0;
            }
        });
    }
    
    const summaryDiv = document.createElement('div');
    summaryDiv.className = 'docker-summary';
    summaryDiv.innerHTML = `
        <div class="summary-panel">
            <h3>Container Health Status</h3>
            <div class="status-indicators">
                <div class="status-indicator ${containerStats.running > 0 ? 'active' : ''}">
                    <div class="status-count">${containerStats.running}</div>
                    <div class="status-label">Running</div>
                </div>
                <div class="status-indicator ${containerStats.exited > 0 ? 'inactive' : ''}">
                    <div class="status-count">${containerStats.exited}</div>
                    <div class="status-label">Stopped</div>
                </div>
                <div class="status-indicator ${containerStats.paused > 0 ? 'paused' : ''}">
                    <div class="status-count">${containerStats.paused}</div>
                    <div class="status-label">Paused</div>
                </div>
                <div class="status-indicator">
                    <div class="status-count">${containerStats.total}</div>
                    <div class="status-label">Total</div>
                </div>
            </div>
            
            <h3>Resource Utilization</h3>
            <div class="resource-gauges">
                <div class="gauge-container">
                    <div class="gauge-title">CPU Usage</div>
                    <div class="gauge">
                        <div class="gauge-fill" style="width: ${Math.min(totalCpuPercent, 100)}%"></div>
                    </div>
                    <div class="gauge-value">${totalCpuPercent.toFixed(2)}%</div>
                </div>
                <div class="gauge-container">
                    <div class="gauge-title">Memory Usage</div>
                    <div class="gauge">
                        <div class="gauge-fill" style="width: ${(totalMemoryUsed / totalMemoryLimit * 100) || 0}%"></div>
                    </div>
                    <div class="gauge-value">${formatBytes(totalMemoryUsed)} / ${formatBytes(totalMemoryLimit)}</div>
                </div>
                <div class="gauge-container">
                    <div class="gauge-title">Network I/O</div>
                    <div class="network-stats">
                        <div><i class="fas fa-arrow-down"></i> ${formatBytes(totalNetworkRx)}</div>
                        <div><i class="fas fa-arrow-up"></i> ${formatBytes(totalNetworkTx)}</div>
                    </div>
                </div>
            </div>
        </div>
        
        <style>
            .docker-summary {
                margin-bottom: 20px;
            }
            .summary-panel {
                background-color: var(--card-bg);
                border-radius: 8px;
                padding: 15px;
                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            }
            .summary-panel h3 {
                margin-top: 0;
                color: var(--accent-color);
                border-bottom: 1px solid var(--border-color);
                padding-bottom: 10px;
                margin-bottom: 15px;
            }
            .status-indicators {
                display: flex;
                justify-content: space-between;
                margin-bottom: 20px;
            }
            .status-indicator {
                text-align: center;
                padding: 10px;
                border-radius: 5px;
                flex: 1;
                margin: 0 5px;
                background: rgba(0,0,0,0.05);
            }
            .status-indicator.active {
                background: rgba(0,255,0,0.1);
                box-shadow: 0 0 0 1px rgba(0,255,0,0.2);
            }
            .status-indicator.inactive {
                background: rgba(255,0,0,0.1);
                box-shadow: 0 0 0 1px rgba(255,0,0,0.2);
            }
            .status-indicator.paused {
                background: rgba(255,165,0,0.1);
                box-shadow: 0 0 0 1px rgba(255,165,0,0.2);
            }
            .status-count {
                font-size: 24px;
                font-weight: bold;
                margin-bottom: 5px;
            }
            .resource-gauges {
                display: flex;
                flex-wrap: wrap;
                gap: 15px;
            }
            .gauge-container {
                flex: 1;
                min-width: 200px;
            }
            .gauge-title {
                margin-bottom: 5px;
                font-weight: bold;
                color: var(--secondary-text-color);
            }
            .gauge {
                height: 10px;
                background: rgba(0,0,0,0.1);
                border-radius: 5px;
                overflow: hidden;
                margin-bottom: 5px;
            }
            .gauge-fill {
                height: 100%;
                background: var(--accent-color);
                border-radius: 5px;
            }
            .gauge-value {
                font-size: 0.9em;
                color: var(--secondary-text-color);
            }
            .network-stats {
                display: flex;
                justify-content: space-between;
            }
        </style>
    `;
    
    const content = document.getElementById('content');
    if (content.firstChild) {
        content.insertBefore(summaryDiv, content.firstChild);
    } else {
        content.appendChild(summaryDiv);
    }
}

function createNetworkTopologyVisualization(data) {
    if (!data.networks || data.networks.length === 0) return;
    
    const tableContainer = document.createElement('div');
    tableContainer.className = 'collapsible-table';
    tableContainer.id = 'docker-network-topology';
    
    const tableHeader = document.createElement('div');
    tableHeader.className = 'table-header collapsed';
    tableHeader.innerHTML = `
        <span>Container Network Topology</span>
        <span class="icon"><i class="fas fa-chevron-down"></i></span>
    `;
    
    const tableContent = document.createElement('div');
    tableContent.className = 'table-content collapsed';
    
    const topologyContainer = document.createElement('div');
    topologyContainer.id = 'network-topology-viz';
    topologyContainer.style.padding = '15px';
    
    data.networks.forEach(network => {
        const networkDiv = document.createElement('div');
        networkDiv.className = `network ${network.internal ? 'network-internal' : ''} ${network.driver === 'bridge' ? 'network-bridge' : ''}`;
        networkDiv.id = `network-${network.id.substring(0, 8)}`;
        
        const subnet = network.ipam && 
                      network.ipam.config && 
                      network.ipam.config.length > 0 ? 
                      network.ipam.config[0].subnet : 'N/A';
        
        const networkHeaderDiv = document.createElement('div');
        networkHeaderDiv.className = 'network-header';
        networkHeaderDiv.innerHTML = `
            <div class="network-name">${network.name} <span class="network-type">(${network.driver})</span></div>
            <div class="network-attrs">Subnet: ${subnet} | Internal: ${network.internal ? 'Yes' : 'No'} | Attachable: ${network.attachable ? 'Yes' : 'No'}</div>
        `;
        
        networkDiv.appendChild(networkHeaderDiv);
        
        const containerIds = Object.keys(network.containers || {});
        
        if (containerIds.length > 0) {
            const containersDiv = document.createElement('div');
            containersDiv.className = 'containers-grid';
            
            containerIds.forEach(containerId => {
                const containerInfo = network.containers[containerId];
                const container = data.containers ? data.containers.find(c => c.id === containerId) : null;
                
                const status = container ? container.state : 'unknown';
                const statusClass = status === 'running' ? 'running' : 
                                   status === 'exited' ? 'stopped' : 'other';
                
                const containerBox = document.createElement('div');
                containerBox.className = `container-box ${statusClass}`;
                containerBox.innerHTML = `
                    <div class="container-name">${containerInfo.name}</div>
                    <div class="container-ip">${containerInfo.ipv4Address || 'No IP'}</div>
                `;
                
                containersDiv.appendChild(containerBox);
            });
            
            networkDiv.appendChild(containersDiv);
        } else {
            const emptyMsg = document.createElement('p');
            emptyMsg.textContent = 'No connected containers';
            networkDiv.appendChild(emptyMsg);
        }
        
        topologyContainer.appendChild(networkDiv);
    });
    
    const style = document.createElement('style');
    style.textContent = `
        .network {
            margin-bottom: 20px;
            padding: 15px;
            border: 1px solid var(--border-color);
            border-radius: 8px;
            position: relative;
        }
        .network-name {
            font-weight: bold;
            font-size: 16px;
            margin-bottom: 10px;
        }
        .network-attrs {
            font-size: 12px;
            color: var(--secondary-text-color);
            margin-bottom: 10px;
        }
        .network-internal {
            background: rgba(255,0,0,0.1);
        }
        .network-bridge {
            background: rgba(0,0,255,0.1);
        }
        .containers-grid {
            display: flex;
            flex-wrap: wrap;
            gap: 8px;
        }
        .container-box {
            padding: 8px;
            border-radius: 4px;
            background: var(--background-color);
            border: 1px solid var(--border-color);
        }
        .container-box.running {
            border-left: 4px solid #4CAF50;
        }
        .container-box.stopped {
            border-left: 4px solid #F44336;
        }
        .container-name {
            font-weight: bold;
            margin-bottom: 4px;
        }
        .container-ip {
            font-size: 0.8em;
            color: var(--secondary-text-color);
        }
    `;
    
    document.head.appendChild(style);
    
    tableContent.appendChild(topologyContainer);
    
    tableHeader.addEventListener('click', function() {
        const isCollapsed = tableHeader.classList.contains('collapsed');
        if (isCollapsed) {
            tableHeader.classList.remove('collapsed');
            tableContent.classList.remove('collapsed');
        } else {
            tableHeader.classList.add('collapsed');
            tableContent.classList.add('collapsed');
        }
    });
    
    tableContainer.appendChild(tableHeader);
    tableContainer.appendChild(tableContent);
    
    document.getElementById('content').appendChild(tableContainer);
}

function createResourceUtilizationCharts(data) {
    if (!data.containers || data.containers.length === 0 || !data.stats) return;
    
    const chartsDiv = document.createElement('div');
    chartsDiv.className = 'docker-charts';
    chartsDiv.style.display = 'grid';
    chartsDiv.style.gridTemplateColumns = 'repeat(auto-fit, minmax(450px, 1fr))';
    chartsDiv.style.gap = '20px';
    chartsDiv.style.marginTop = '20px';
    
    chartsDiv.innerHTML = `
        <div class="chart-wrapper" style="height: 300px; margin-bottom: 20px;">
            <canvas id="containerCpuChart"></canvas>
        </div>
        <div class="chart-wrapper" style="height: 300px; margin-bottom: 20px;">
            <canvas id="containerMemoryChart"></canvas>
        </div>
        <div class="chart-wrapper" style="height: 300px; margin-bottom: 20px;">
            <canvas id="containerNetworkChart"></canvas>
        </div>
    `;
    
    document.getElementById('content').appendChild(chartsDiv);
    
    const runningContainers = data.containers.filter(c => 
        c.state === 'running' && 
        data.stats && 
        data.stats[c.id]
    );
    
    if (runningContainers.length === 0) return;
    
    const containerNames = runningContainers.map(c => c.name);
    const cpuValues = runningContainers.map(c => data.stats[c.id].cpuPercentage || 0);
    const memoryValues = runningContainers.map(c => {
        const stats = data.stats[c.id];
        return stats.memoryPercentage || 0;
    });
    const networkRxValues = runningContainers.map(c => data.stats[c.id].networkRx || 0);
    const networkTxValues = runningContainers.map(c => data.stats[c.id].networkTx || 0);
    
    const generateColors = (count) => {
        const colors = [];
        for (let i = 0; i < count; i++) {
            const hue = (i * 137.508) % 360; // Use golden angle approximation for distribution
            colors.push(`hsl(${hue}, 75%, 50%)`);
        }
        return colors;
    };
    
    const containerColors = generateColors(containerNames.length);
    
    const cpuCtx = document.getElementById('containerCpuChart');
    if (cpuCtx) {
        new Chart(cpuCtx.getContext('2d'), {
            type: 'bar',
            data: {
                labels: containerNames,
                datasets: [{
                    label: 'CPU Usage (%)',
                    data: cpuValues,
                    backgroundColor: containerColors,
                    borderColor: containerColors,
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Container CPU Usage',
                        font: {
                            size: 16
                        }
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        title: {
                            display: true,
                            text: 'CPU Usage (%)'
                        }
                    }
                }
            }
        });
    }
    
    const memoryCtx = document.getElementById('containerMemoryChart');
    if (memoryCtx) {
        new Chart(memoryCtx.getContext('2d'), {
            type: 'bar',
            data: {
                labels: containerNames,
                datasets: [{
                    label: 'Memory Usage (%)',
                    data: memoryValues,
                    backgroundColor: containerColors,
                    borderColor: containerColors,
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Container Memory Usage',
                        font: {
                            size: 16
                        }
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        title: {
                            display: true,
                            text: 'Memory Usage (%)'
                        }
                    }
                }
            }
        });
    }
    
    const networkCtx = document.getElementById('containerNetworkChart');
    if (networkCtx) {
        new Chart(networkCtx.getContext('2d'), {
            type: 'bar',
            data: {
                labels: containerNames,
                datasets: [
                    {
                        label: 'Network Download',
                        data: networkRxValues,
                        backgroundColor: 'rgba(54, 162, 235, 0.7)',
                        borderColor: 'rgba(54, 162, 235, 1)',
                        borderWidth: 1
                    },
                    {
                        label: 'Network Upload',
                        data: networkTxValues,
                        backgroundColor: 'rgba(255, 99, 132, 0.7)',
                        borderColor: 'rgba(255, 99, 132, 1)',
                        borderWidth: 1
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'Container Network I/O',
                        font: {
                            size: 16
                        }
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return `${context.dataset.label}: ${formatBytes(context.raw)}`;
                            }
                        }
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        title: {
                            display: true,
                            text: 'Data Transfer'
                        },
                        ticks: {
                            callback: function(value) {
                                return formatBytes(value);
                            }
                        }
                    }
                }
            }
        });
    }
}

function dockerInfoCategoryRowTemplate(item) {
    let rows = [];
    for (const [key, value] of Object.entries(item)) {
        rows.push(`<tr><td>${key}</td><td>${value}</td></tr>`);
    }
    return rows.join('');
}

function containerRowTemplate(item) {
    let portMappings = '';
    if (item.ports && item.ports.length > 0) {
        portMappings = item.ports.map(p => {
            if (p.PublicPort) {
                return `${p.PublicPort}:${p.PrivatePort}/${p.Type}`;
            }
            return `${p.PrivatePort}/${p.Type}`;
        }).join(', ');
    }
    
    let networks = item.networks.join(', ') || 'none';
    
    const created = new Date(item.created);
    const now = new Date();
    const diff = (now - created) / 1000;
    
    let createdStr;
    if (diff < 60) {
        createdStr = `${Math.round(diff)} seconds ago`;
    } else if (diff < 3600) {
        createdStr = `${Math.round(diff / 60)} minutes ago`;
    } else if (diff < 86400) {
        createdStr = `${Math.round(diff / 3600)} hours ago`;
    } else {
        createdStr = `${Math.round(diff / 86400)} days ago`;
    }
    
    const containerId = `container-${item.id.substring(0, 12)}`;
    
    return `
        <td>${item.name}</td>
        <td>${item.image}</td>
        <td>${item.status}</td>
        <td>${createdStr}</td>
        <td>${portMappings}</td>
        <td>${networks}</td>
        <td>
            <button class="details-button" onclick="toggleContainerDetails('${containerId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${containerId}" style="display:none;" class="details-panel">
                <h4>Container Details</h4>
                
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>ID:</td><td>${item.id}</td></tr>
                        <tr><td>State:</td><td><span class="badge ${item.state === 'running' ? 'badge-success' : 'badge-secondary'}">${item.state}</span></td></tr>
                        <tr><td>Command:</td><td><code>${item.command}</code></td></tr>
                        <tr><td>Image ID:</td><td>${item.imageId}</td></tr>
                        <tr><td>Created:</td><td>${new Date(item.created).toLocaleString()}</td></tr>
                    </table>
                </div>
                
                <div class="detail-section">
                    <h5>Host Configuration</h5>
                    <table class="nested-table">
                        <tr><td>Network Mode:</td><td>${item.hostConfig.networkMode}</td></tr>
                        <tr><td>Privileged:</td><td>${item.hostConfig.privileged ? '<span class="badge badge-warning">Yes</span>' : '<span class="badge badge-info">No</span>'}</td></tr>
                        <tr><td>Restart Policy:</td><td>${item.hostConfig.restartPolicy}</td></tr>
                    </table>
                </div>
                
                ${item.mounts && item.mounts.length > 0 ? `
                    <div class="detail-section">
                        <h5>Mounts (${item.mounts.length})</h5>
                        <table class="nested-table">
                            <tr>
                                <th>Type</th>
                                <th>Source</th>
                                <th>Destination</th>
                                <th>Mode</th>
                                <th>RW</th>
                            </tr>
                            ${item.mounts.map(m => `
                                <tr>
                                    <td>${m.type}</td>
                                    <td>${m.source}</td>
                                    <td>${m.destination}</td>
                                    <td>${m.mode}</td>
                                    <td>${m.rw ? '<i class="fas fa-check text-success"></i>' : '<i class="fas fa-times text-danger"></i>'}</td>
                                </tr>
                            `).join('')}
                        </table>
                    </div>
                ` : ''}
                
                ${(window.currentData && window.currentData.stats && item.id in window.currentData.stats) ? `
                    <div class="detail-section">
                        <h5>Resource Usage</h5>
                        <div class="resource-metrics">
                            <div class="metric-card">
                                <div class="metric-header">CPU Usage</div>
                                <div class="metric-value">${window.currentData.stats[item.id].cpuPercentage.toFixed(2)}%</div>
                                <div class="metric-gauge">
                                    <div class="gauge-fill" style="width: ${Math.min(window.currentData.stats[item.id].cpuPercentage, 100)}%"></div>
                                </div>
                            </div>
                            
                            <div class="metric-card">
                                <div class="metric-header">Memory Usage</div>
                                <div class="metric-value">
                                    ${formatBytes(window.currentData.stats[item.id].memoryUsage)} / ${formatBytes(window.currentData.stats[item.id].memoryLimit)}
                                    <span class="metric-percentage">(${window.currentData.stats[item.id].memoryPercentage.toFixed(2)}%)</span>
                                </div>
                                <div class="metric-gauge">
                                    <div class="gauge-fill" style="width: ${Math.min(window.currentData.stats[item.id].memoryPercentage, 100)}%"></div>
                                </div>
                            </div>
                            
                            <div class="metric-card">
                                <div class="metric-header">Network I/O</div>
                                <div class="metric-value">
                                    <span class="network-rx"><i class="fas fa-arrow-down"></i> ${formatBytes(window.currentData.stats[item.id].networkRx)}</span>
                                    <span class="network-tx"><i class="fas fa-arrow-up"></i> ${formatBytes(window.currentData.stats[item.id].networkTx)}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                ` : ''}
                
                <style>
                    .detail-section {
                        margin-bottom: 15px;
                        background: var(--background-color);
                        padding: 10px;
                        border-radius: 5px;
                    }
                    .detail-section h5 {
                        margin-top: 0;
                        color: var(--accent-color);
                        font-weight: bold;
                    }
                    .resource-metrics {
                        display: flex;
                        flex-wrap: wrap;
                        gap: 12px;
                    }
                    .metric-card {
                        flex: 1 1 200px;
                        padding: 10px;
                        border: 1px solid var(--border-color);
                        border-radius: 5px;
                    }
                    .metric-header {
                        font-weight: bold;
                        color: var(--secondary-text-color);
                        font-size: 0.9em;
                    }
                    .metric-value {
                        font-size: 1.1em;
                        margin: 8px 0;
                    }
                    .metric-percentage {
                        font-size: 0.9em;
                        color: var(--secondary-text-color);
                    }
                    .network-rx, .network-tx {
                        display: inline-block;
                        margin-right: 10px;
                    }
                    .network-rx i {
                        color: #4CAF50;
                    }
                    .network-tx i {
                        color: #F44336;
                    }
                    .metric-gauge {
                        height: 6px;
                        background: rgba(0,0,0,0.1);
                        border-radius: 3px;
                        overflow: hidden;
                    }
                    .gauge-fill {
                        height: 100%;
                        background: var(--accent-color);
                    }
                    .badge {
                        padding: 3px 6px;
                        border-radius: 3px;
                        font-size: 0.85em;
                        font-weight: normal;
                    }
                    .badge-success {
                        background-color: #4CAF50;
                        color: white;
                    }
                    .badge-secondary {
                        background-color: #607D8B;
                        color: white;
                    }
                    .badge-warning {
                        background-color: #FF9800;
                        color: white;
                    }
                    .badge-info {
                        background-color: #2196F3;
                        color: white;
                    }
                    code {
                        background: rgba(0,0,0,0.1);
                        padding: 2px 4px;
                        border-radius: 3px;
                        font-family: monospace;
                    }
                    .text-success {
                        color: #4CAF50;
                    }
                    .text-danger {
                        color: #F44336;
                    }
                </style>
            </div>
        </td>
    `;
}

function imageRowTemplate(item) {
    let tags = item.repoTags && item.repoTags.length > 0 && item.repoTags[0] !== '<none>:<none>' 
        ? item.repoTags.join(', ') 
        : '<none>';
    
    const created = new Date(item.created);
    const createdStr = `${created.toLocaleDateString()} ${created.toLocaleTimeString()}`;
    const size = formatBytes(item.size);
    const shortId = item.id.substring(7, 19);
    const imageId = `image-${shortId}`;
    
    return `
        <td>${tags}</td>
        <td>${shortId}</td>
        <td>${createdStr}</td>
        <td>${size}</td>
        <td>
            <button class="details-button" onclick="toggleImageDetails('${imageId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${imageId}" style="display:none;" class="details-panel">
                <h4>Image Details</h4>
                <div class="detail-section">
                    <table class="nested-table">
                        <tr><td>Full ID:</td><td>${item.id}</td></tr>
                        <tr><td>Created:</td><td>${new Date(item.created).toLocaleString()}</td></tr>
                        <tr><td>Size:</td><td>${formatBytes(item.size)}</td></tr>
                        ${item.repoDigests && item.repoDigests.length > 0 ? `
                            <tr><td>Digests:</td><td>${item.repoDigests.join('<br>')}</td></tr>
                        ` : ''}
                    </table>
                </div>
                
                ${item.labels && Object.keys(item.labels).length > 0 ? `
                    <div class="detail-section">
                        <h5>Labels</h5>
                        <table class="nested-table">
                            <tr>
                                <th>Key</th>
                                <th>Value</th>
                            </tr>
                            ${Object.entries(item.labels).map(([key, value]) => `
                                <tr>
                                    <td>${key}</td>
                                    <td>${value}</td>
                                </tr>
                            `).join('')}
                        </table>
                    </div>
                ` : ''}
                
                <style>
                    .detail-section {
                        margin-bottom: 15px;
                        background: var(--background-color);
                        padding: 10px;
                        border-radius: 5px;
                    }
                    .detail-section h5 {
                        margin-top: 0;
                        color: var(--accent-color);
                        font-weight: bold;
                    }
                </style>
            </div>
        </td>
    `;
}

function volumeRowTemplate(item) {
    const volumeId = `volume-${item.name.replace(/[^a-zA-Z0-9]/g, '-')}`;
    
    return `
        <td>${item.name}</td>
        <td>${item.driver}</td>
        <td>${item.mountpoint}</td>
        <td>${item.scope}</td>
        <td>
            <button class="details-button" onclick="toggleVolumeDetails('${volumeId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${volumeId}" style="display:none;" class="details-panel">
                <h4>Volume Details</h4>
                <div class="detail-section">
                    ${item.createdAt ? `<p><strong>Created:</strong> ${new Date(item.createdAt).toLocaleString()}</p>` : ''}
                    
                    ${item.labels && Object.keys(item.labels).length > 0 ? `
                        <h5>Labels</h5>
                        <table class="nested-table">
                            <tr>
                                <th>Key</th>
                                <th>Value</th>
                            </tr>
                            ${Object.entries(item.labels).map(([key, value]) => `
                                <tr>
                                    <td>${key}</td>
                                    <td>${value}</td>
                                </tr>
                            `).join('')}
                        </table>
                    ` : ''}
                    
                    ${item.status && Object.keys(item.status).length > 0 ? `
                        <h5>Status</h5>
                        <table class="nested-table">
                            <tr>
                                <th>Key</th>
                                <th>Value</th>
                            </tr>
                            ${Object.entries(item.status).map(([key, value]) => `
                                <tr>
                                    <td>${key}</td>
                                    <td>${value}</td>
                                </tr>
                            `).join('')}
                        </table>
                    ` : ''}
                </div>
                
                <style>
                    .detail-section {
                        margin-bottom: 15px;
                        background: var(--background-color);
                        padding: 10px;
                        border-radius: 5px;
                    }
                    .detail-section h5 {
                        margin-top: 0;
                        color: var(--accent-color);
                        font-weight: bold;
                        margin-bottom: 10px;
                    }
                </style>
            </div>
        </td>
    `;
}

function networkRowTemplate(item) {
    let subnet = '';
    if (item.ipam && item.ipam.config && item.ipam.config.length > 0) {
        subnet = item.ipam.config.map(c => c.subnet).join(', ');
    }
    
    const networkId = `network-${item.id.substring(0, 12)}`;
    
    return `
        <td>${item.name}</td>
        <td>${item.id.substring(0, 12)}</td>
        <td>${item.driver}</td>
        <td>${item.scope}</td>
        <td>${subnet}</td>
        <td>
            <button class="details-button" onclick="toggleNetworkDetails('${networkId}')">
                <i class="fas fa-info-circle"></i> Details
            </button>
            <div id="${networkId}" style="display:none;" class="details-panel">
                <h4>Network Details</h4>
                
                <div class="detail-section">
                    <h5>Basic Information</h5>
                    <table class="nested-table">
                        <tr><td>Full ID:</td><td>${item.id}</td></tr>
                        <tr><td>Internal:</td><td>${item.internal ? 'Yes' : 'No'}</td></tr>
                        <tr><td>Attachable:</td><td>${item.attachable ? 'Yes' : 'No'}</td></tr>
                    </table>
                </div>
                
                ${item.ipam && item.ipam.config && item.ipam.config.length > 0 ? `
                    <div class="detail-section">
                        <h5>IPAM Configuration</h5>
                        <p><strong>Driver:</strong> ${item.ipam.driver || 'default'}</p>
                        
                        <table class="nested-table">
                            <tr>
                                <th>Subnet</th>
                                <th>Gateway</th>
                            </tr>
                            ${item.ipam.config.map(config => `
                                <tr>
                                    <td>${config.subnet || '-'}</td>
                                    <td>${config.gateway || '-'}</td>
                                </tr>
                            `).join('')}
                        </table>
                        
                        ${item.ipam.options && Object.keys(item.ipam.options).length > 0 ? `
                            <h6>IPAM Options</h6>
                            <table class="nested-table">
                                ${Object.entries(item.ipam.options).map(([key, value]) => `
                                    <tr>
                                        <td>${key}</td>
                                        <td>${value}</td>
                                    </tr>
                                `).join('')}
                            </table>
                        ` : ''}
                    </div>
                ` : ''}
                
                ${item.containers && Object.keys(item.containers).length > 0 ? `
                    <div class="detail-section">
                        <h5>Connected Containers (${Object.keys(item.containers).length})</h5>
                        <table class="nested-table">
                            <tr>
                                <th>Name</th>
                                <th>IPv4 Address</th>
                                <th>IPv6 Address</th>
                                <th>MAC Address</th>
                            </tr>
                            ${Object.values(item.containers).map(c => `
                                <tr>
                                    <td>${c.name}</td>
                                    <td>${c.ipv4Address || '-'}</td>
                                    <td>${c.ipv6Address || '-'}</td>
                                    <td>${c.macAddress || '-'}</td>
                                </tr>
                            `).join('')}
                        </table>
                    </div>
                ` : '<div class="detail-section"><p>No connected containers</p></div>'}
                
                ${item.labels && Object.keys(item.labels).length > 0 ? `
                    <div class="detail-section">
                        <h5>Labels</h5>
                        <table class="nested-table">
                            <tr>
                                <th>Key</th>
                                <th>Value</th>
                            </tr>
                            ${Object.entries(item.labels).map(([key, value]) => `
                                <tr>
                                    <td>${key}</td>
                                    <td>${value}</td>
                                </tr>
                            `).join('')}
                        </table>
                    </div>
                ` : ''}
                
                <style>
                    .detail-section {
                        margin-bottom: 15px;
                        background: var(--background-color);
                        padding: 10px;
                        border-radius: 5px;
                    }
                    .detail-section h5 {
                        margin-top: 0;
                        color: var(--accent-color);
                        font-weight: bold;
                        margin-bottom: 10px;
                    }
                    .detail-section h6 {
                        margin-top: 15px;
                        margin-bottom: 5px;
                        color: var(--secondary-text-color);
                    }
                </style>
            </div>
        </td>
    `;
}

function toggleContainerDetails(id) {
    const element = document.getElementById(id);
    if (element) {
        element.style.display = element.style.display === 'none' ? 'block' : 'none';
    }
}

function toggleImageDetails(id) {
    const element = document.getElementById(id);
    if (element) {
        element.style.display = element.style.display === 'none' ? 'block' : 'none';
    }
}

function toggleVolumeDetails(id) {
    const element = document.getElementById(id);
    if (element) {
        element.style.display = element.style.display === 'none' ? 'block' : 'none';
    }
}

function toggleNetworkDetails(id) {
    const element = document.getElementById(id);
    if (element) {
        element.style.display = element.style.display === 'none' ? 'block' : 'none';
    }
}

function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

function showDockerConnectionModal() {
    console.log("Creating Docker connection modal");
    const isConnected = document.getElementById('docker-button')?.classList.contains('connected');
    
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
    modalContent.className = 'modal-content docker-modal';
    modalContent.style.backgroundColor = 'var(--card-bg)';
    modalContent.style.color = 'var(--text-color)';
    modalContent.style.padding = '25px';
    modalContent.style.borderRadius = '8px';
    modalContent.style.maxWidth = '600px';
    modalContent.style.width = '90%';
    modalContent.style.maxHeight = '90vh';
    modalContent.style.overflow = 'auto';
    modalContent.style.boxShadow = '0 5px 20px rgba(0,0,0,0.4)';
    modalContent.style.border = '1px solid var(--border-color)';
    
    let connectionNote = '';
    if (isConnected) {
        connectionNote = `
            <div style="background-color: rgba(0,255,0,0.1); border-left: 4px solid #4CAF50; padding: 8px; margin-bottom: 15px;">
                <p style="margin: 0; color: var(--text-color);">
                    <i class="fas fa-info-circle"></i> You are already connected to Docker. 
                    You can switch to a different Docker daemon if needed.
                </p>
            </div>
        `;
    }
    
    modalContent.innerHTML = `
        <h3 style="margin-top: 0; color: var(--accent-color); font-size: 1.5em; border-bottom: 1px solid var(--border-color); padding-bottom: 10px;">
            <i class="fab fa-docker"></i> Connect to Docker
        </h3>
        
        ${connectionNote}
        
        <div class="docker-connection-form" style="margin-top: 20px;">
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="docker-local" name="docker-source" value="local" checked>
                <label for="docker-local" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-laptop"></i> Local Docker Daemon
                </label>
                <div id="docker-local-form" class="source-form" style="margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <p style="margin-top: 0;">Connect to the Docker daemon running on this machine (default socket).</p>
                </div>
            </div>
            
            <div class="source-option" style="background: var(--background-color); border-radius: 6px; padding: 12px; margin-bottom: 15px;">
                <input type="radio" id="docker-remote" name="docker-source" value="remote">
                <label for="docker-remote" style="font-weight: bold; font-size: 1.1em;">
                    <i class="fas fa-server"></i> Remote Docker Daemon
                </label>
                <div id="docker-remote-form" class="source-form" style="display: none; margin-top: 12px; margin-left: 25px; padding: 10px; background: rgba(255,255,255,0.05); border-radius: 4px;">
                    <div class="form-group">
                        <label for="docker-host" style="font-weight: bold; margin-bottom: 5px; display: block;">Docker Host:</label>
                        <input type="text" id="docker-host" placeholder="tcp://hostname:2375" style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px; box-sizing: border-box;">
                    </div>
                    
            <div class="connection-help" style="margin-top: 15px; padding: 10px; border-radius: 5px; background: rgba(33, 150, 243, 0.1); border-left: 4px solid #2196F3;">
                <p style="margin: 0; font-weight: bold; color: var(--text-color);">Remote Host Setup</p>
                <p style="margin-top: 5px; margin-bottom: 0; font-size: 0.9em; color: var(--secondary-text-color);">
                    For remote connections, you need to configure Docker on the remote host to accept API connections. Run these commands:
                </p>
                <ol style="margin-top: 5px; padding-left: 25px; font-size: 0.9em; color: var(--secondary-text-color);">
                    <li>SSH into your remote server</li>
                    <li>Run these commands to set up Docker for remote connections:
                    <pre style="background: rgba(0,0,0,0.1); padding: 8px; margin: 5px 0; border-radius: 3px; overflow-x: auto; font-size: 12px;"># Create required directories
            sudo mkdir -p /etc/docker
            sudo mkdir -p /etc/systemd/system/docker.service.d/
            
            # Configure Docker to listen on TCP
            sudo bash -c 'cat > /etc/docker/daemon.json << EOF
            {
              "hosts": ["unix:///var/run/docker.sock", "tcp://0.0.0.0:2375"]
            }
            EOF'
            
            # Create systemd override
            sudo bash -c 'cat > /etc/systemd/system/docker.service.d/override.conf << EOF
            [Service]
            ExecStart=
            ExecStart=/usr/bin/dockerd
            EOF'
            
            # Reload systemd and restart Docker
            sudo systemctl daemon-reload
            sudo systemctl restart docker
            
            # Verify Docker is listening on port 2375
            sudo netstat -tuln | grep 2375</pre>
                    </li>
                </ol>
                <p style="margin-top: 5px; font-style: italic; font-size: 0.85em; color: var(--warning-color);">
                    Warning: Exposing Docker API on port 2375 without TLS is insecure. For production, use TLS (port 2376) or SSH tunneling.
                </p>
            </div>
                    
                    <p class="tip" style="margin-top: 15px; font-size: 0.85em; color: var(--secondary-text-color);">
                        Example format: <code>tcp://192.168.1.10:2375</code>
                    </p>
                </div>
            </div>
            
            <div id="connection-status" style="display: none; margin: 15px 0; padding: 10px; border-radius: 5px;"></div>
            
            <div class="modal-buttons" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 25px; border-top: 1px solid var(--border-color); padding-top: 15px;">
                <button id="docker-cancel-btn" class="btn" style="padding: 10px 20px; background-color: var(--button-bg-color); color: var(--button-text-color); border: none; border-radius: 4px; cursor: pointer;">Cancel</button>
                <button id="test-connection-btn" class="btn" style="padding: 10px 20px; background-color: var(--secondary-color); color: white; border: none; border-radius: 4px; cursor: pointer;">
                    <i class="fas fa-heartbeat"></i> Test Connection
                </button>
                <button id="docker-connect-btn" class="btn btn-primary" style="padding: 10px 20px; background-color: var(--accent-color); color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">
                    <i class="fas fa-plug"></i> Connect
                </button>
            </div>
        </div>
    `;
    
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    console.log("Modal created and added to DOM");

    const sourceForms = document.querySelectorAll('.source-form');
    document.querySelectorAll('input[name="docker-source"]').forEach(radio => {
        radio.addEventListener('change', () => {
            console.log(`Radio changed to: ${radio.value}`);
            sourceForms.forEach(form => form.style.display = 'none');
            const selectedForm = document.getElementById(`docker-${radio.value}-form`);
            if (selectedForm) {
                selectedForm.style.display = 'block';
            }
            document.getElementById('connection-status').style.display = 'none';
        });
    });

    document.getElementById('docker-cancel-btn').addEventListener('click', () => {
        console.log("Cancel button clicked");
        modal.remove();
    });
    
    document.getElementById('test-connection-btn').addEventListener('click', () => {
        console.log("Test connection button clicked");
        const source = document.querySelector('input[name="docker-source"]:checked').value;
        if (source === 'remote') {
            const host = document.getElementById('docker-host').value.trim();
            if (!host) {
                showConnectionStatus('error', 'Please enter a Docker host address');
                return;
            }
            testDockerConnection(host);
        } else {
            testDockerConnection('');
        }
    });

    document.getElementById('docker-connect-btn').addEventListener('click', () => {
        console.log("Connect button clicked");
        
        const source = document.querySelector('input[name="docker-source"]:checked').value;
        let dockerHost = '';
        
        if (source === 'remote') {
            dockerHost = document.getElementById('docker-host').value.trim();
            if (!dockerHost) {
                showConnectionStatus('error', 'Please enter a Docker host address');
                return;
            }
        }
        
        connectToDocker(dockerHost);
    });

    function testDockerConnection(host) {
        console.log(`Testing connection to Docker host: ${host || 'local'}`);
        showLoadingIndicator();
        
        showConnectionStatus('pending', 'Testing connection...');
        
        fetch('/api/docker/test-connection', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                host: host
            })
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => {
                    throw new Error(`Error: ${text}`);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("Test connection response:", data);
            if (data.status === 'success') {
                showConnectionStatus('success', `Connection successful! Docker version: ${data.version || 'Unknown'}`);
            } else {
                throw new Error(data.message || 'Failed to connect to Docker');
            }
        })
        .catch(error => {
            console.error('Docker test connection error:', error);
            
            let errorMessage = error.message || 'Unknown error';
            let helpText = '';
            
            if (errorMessage.includes("Cannot connect to the Docker daemon")) {
                helpText = `
                    <div style="margin-top: 10px;">
                        <p><strong>Possible solutions:</strong></p>
                        <ol>
                            <li>Make sure Docker is running on the remote host</li>
                            <li>Check if the Docker daemon is configured to accept remote connections (see setup instructions above)</li>
                            <li>Verify that port 2375 (or the port you specified) is open in the firewall</li>
                            <li>Check network connectivity between this system and the remote host</li>
                        </ol>
                    </div>
                `;
            }
            
            showConnectionStatus('error', `Connection failed: ${errorMessage}${helpText}`);
        })
        .finally(() => {
            hideLoadingIndicator();
        });
    }

    function showConnectionStatus(type, message) {
        const statusDiv = document.getElementById('connection-status');
        statusDiv.style.display = 'block';
        
        let bgColor, iconClass;
        
        switch(type) {
            case 'success':
                bgColor = 'rgba(76, 175, 80, 0.1)';
                iconClass = 'fa-check-circle';
                statusDiv.style.borderLeft = '4px solid #4CAF50';
                break;
            case 'error':
                bgColor = 'rgba(244, 67, 54, 0.1)';
                iconClass = 'fa-exclamation-circle';
                statusDiv.style.borderLeft = '4px solid #F44336';
                break;
            case 'pending':
                bgColor = 'rgba(33, 150, 243, 0.1)';
                iconClass = 'fa-spinner fa-spin';
                statusDiv.style.borderLeft = '4px solid #2196F3';
                break;
            default:
                bgColor = 'rgba(0, 0, 0, 0.1)';
                iconClass = 'fa-info-circle';
                statusDiv.style.borderLeft = '4px solid #9E9E9E';
        }
        
        statusDiv.style.backgroundColor = bgColor;
        statusDiv.innerHTML = `<p style="margin: 0;"><i class="fas ${iconClass}"></i> ${message}</p>`;
    }

    function connectToDocker(host) {
        console.log(`Connecting to Docker host: ${host || 'local'}`);
        showLoadingIndicator();
        
        showConnectionStatus('pending', 'Connecting to Docker...');
        
        fetch('/api/docker/connect', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                host: host
            })
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => {
                    throw new Error(`${text}`);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("Connection response:", data);
            if (data.status === 'success') {
                const button = document.getElementById('docker-button');
                if (button) {
                    button.classList.add('connected');
                    button.classList.remove('not-connected');
                    
                    const existingBadges = button.querySelectorAll('.connection-badge');
                    existingBadges.forEach(badge => badge.remove());
                    
                    const badge = document.createElement('span');
                    badge.className = 'connection-badge connected';
                    button.appendChild(badge);
                    
                    button.title = host ? `Docker (Connected to ${host})` : 'Docker (Connected to local)';
                    console.log("Button updated to connected state");
                }
                
                showConnectionStatus('success', 'Successfully connected! Reloading page...');
                
                setTimeout(() => {
                    modal.remove();
                    console.log("Modal removed, reloading page");
                    location.reload();
                }, 1000);
            } else {
                throw new Error(data.message || 'Failed to connect to Docker');
            }
        })
        .catch(error => {
            console.error('Docker connection error:', error);
            
            let errorMessage = error.message || 'Unknown error';
            let helpText = '';
            
            if (errorMessage.includes("Cannot connect to the Docker daemon")) {
                helpText = `
                    <div style="margin-top: 10px;">
                        <p><strong>Possible solutions:</strong></p>
                        <ol>
                            <li>Make sure Docker is running on the remote host</li>
                            <li>Check if the Docker daemon is configured to accept remote connections (see setup instructions above)</li>
                            <li>Verify that port 2375 (or the port you specified) is open in the firewall</li>
                            <li>Check network connectivity between this system and the remote host</li>
                        </ol>
                    </div>
                `;
            }
            
            showConnectionStatus('error', `Connection failed: ${errorMessage}${helpText}`);
        })
        .finally(() => {
            hideLoadingIndicator();
        });
    }
}

function showLoadingIndicator() {
    let loader = document.getElementById('global-loader');
    if (!loader) {
        loader = document.createElement('div');
        loader.id = 'global-loader';
        loader.style.position = 'fixed';
        loader.style.top = '0';
        loader.style.left = '0';
        loader.style.width = '100%';
        loader.style.height = '100%';
        loader.style.display = 'flex';
        loader.style.justifyContent = 'center';
        loader.style.alignItems = 'center';
        loader.style.backgroundColor = 'rgba(0,0,0,0.5)';
        loader.style.zIndex = '9999';
        
        const spinner = document.createElement('div');
        spinner.className = 'loader-spinner';
        spinner.style.border = '5px solid #f3f3f3';
        spinner.style.borderTop = '5px solid var(--accent-color)';
        spinner.style.borderRadius = '50%';
        spinner.style.width = '50px';
        spinner.style.height = '50px';
        spinner.style.animation = 'spin 2s linear infinite';
        
        const style = document.createElement('style');
        style.textContent = '@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }';
        document.head.appendChild(style);
        
        loader.appendChild(spinner);
        document.body.appendChild(loader);
    } else {
        loader.style.display = 'flex';
    }
}

function hideLoadingIndicator() {
    const loader = document.getElementById('global-loader');
    if (loader) {
        loader.style.display = 'none';
    }
}

document.addEventListener('DOMContentLoaded', function() {
    console.log("DOM loaded - Docker module setting up event listener");
    
    const dockerButton = document.getElementById('docker-button');
    if (dockerButton) {
        console.log("Found docker button, setting up handler");
        
        const newButton = dockerButton.cloneNode(true);
        if (dockerButton.parentNode) {
            dockerButton.parentNode.replaceChild(newButton, dockerButton);
        }
        
        newButton.addEventListener('click', function(event) {
            console.log("Docker button clicked");
            event.preventDefault();
            
            showDockerConnectionModal();
        });
    } else {
        console.error("Could not find docker-button element");
    }
    
    if (!window.currentData) {
        fetch('/api/data')
            .then(response => response.json())
            .then(data => {
                window.currentData = data;
            })
            .catch(error => {
                console.error('Error fetching current data:', error);
            });
    }
});