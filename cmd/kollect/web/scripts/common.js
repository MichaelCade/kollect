// common.js

function showLoadingIndicator() {
    document.getElementById('loading-indicator').style.display = 'flex';
}

function hideLoadingIndicator() {
    document.getElementById('loading-indicator').style.display = 'none';
}

function createTable(headerText, data, rowTemplate, headers) {
    if (!data || data.length === 0) return;
    const tableId = headerText.replace(/\s+/g, '-').toLowerCase();
    const tableContainer = document.createElement('div');
    tableContainer.className = 'collapsible-table';
    tableContainer.id = `table-container-${tableId}`;
    const tableHeader = document.createElement('div');
    tableHeader.className = 'table-header collapsed';
    tableHeader.innerHTML = `
        <span>${headerText}</span>
        <div>
            <span class="table-counter">${data.length}</span>
            <span class="icon">▼</span>
        </div>
    `;

    const tableContent = document.createElement('div');
    tableContent.className = 'table-content collapsed';
    
    const table = document.createElement('table');
    const thead = document.createElement('thead');
    const headerRow = document.createElement('tr');
    
    headers.forEach(header => {
        const th = document.createElement('th');
        th.textContent = header;
        headerRow.appendChild(th);
    });
    
    thead.appendChild(headerRow);
    table.appendChild(thead);

    const tbody = document.createElement('tbody');
    data.forEach(item => {
        const row = document.createElement('tr');
        row.innerHTML = rowTemplate(item);
        tbody.appendChild(row);
    });
    table.appendChild(tbody);
    
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
    
    tableContent.appendChild(table);
    tableContainer.appendChild(tableHeader);
    tableContainer.appendChild(tableContent);
    
    document.getElementById('content').appendChild(tableContainer);
    
    return tableContainer;
}

function toggleAllTables(collapse) {
    console.log("Toggling all tables:", collapse ? "Collapse" : "Expand");
    const tables = document.querySelectorAll('.collapsible-table');
    console.log(`Found ${tables.length} tables`);
    
    tables.forEach(table => {
        const header = table.querySelector('.table-header');
        const content = table.querySelector('.table-content');
        
        if (collapse) {
            header.classList.add('collapsed');
            content.classList.add('collapsed');
        } else {
            header.classList.remove('collapsed');
            content.classList.remove('collapsed');
        }
    });
    
    const toggleButton = document.getElementById('toggle-tables');
    if (toggleButton) {
        if (collapse) {
            toggleButton.innerHTML = '<i class="fas fa-expand-alt"></i>';
            toggleButton.setAttribute('data-collapsed', 'true');
            toggleButton.title = 'Expand All Tables';
        } else {
            toggleButton.innerHTML = '<i class="fas fa-compress-alt"></i>';
            toggleButton.setAttribute('data-collapsed', 'false');
            toggleButton.title = 'Collapse All Tables';
        }
    }
}

function updateResourceNav() {
    const resourceNav = document.getElementById('resource-nav');
    if (!resourceNav) return;
    
    resourceNav.innerHTML = '';
    
    const tables = document.querySelectorAll('.collapsible-table');
    
    if (tables.length === 0) {
        resourceNav.style.display = 'none';
        const navToggle = document.getElementById('resource-nav-toggle');
        if (navToggle) navToggle.style.display = 'none';
        return;
    } else {
        const navToggle = document.getElementById('resource-nav-toggle');
        if (navToggle) navToggle.style.display = 'block';
    }

    const heading = document.createElement('h4');
    heading.textContent = 'Resources';
    heading.style.margin = '0 0 10px 0';
    heading.style.padding = '0';
    heading.style.textAlign = 'center';
    resourceNav.appendChild(heading);
    
    tables.forEach(table => {
        const headerEl = table.querySelector('.table-header span');
        if (!headerEl) return;
        
        const header = headerEl.textContent;
        const id = table.id;
        
        const link = document.createElement('a');
        link.textContent = header;
        link.href = `#${id}`;
        link.onclick = function(e) {
            e.preventDefault();
            table.scrollIntoView({ behavior: 'smooth' });
            
            const content = table.querySelector('.table-content');
            const headerElement = table.querySelector('.table-header');
            headerElement.classList.remove('collapsed');
            content.classList.remove('collapsed');
        };
        
        resourceNav.appendChild(link);
    });
    
    console.log("Navigation updated");
}

const navToggle = document.getElementById('resource-nav-toggle');
if (navToggle) {
    console.log("Found nav toggle button");
    navToggle.addEventListener('click', () => {
        const nav = document.getElementById('resource-nav');
        if (nav) {
            if (nav.style.display === 'none' || nav.style.display === '') {
                nav.style.display = 'block';
            } else {
                nav.style.display = 'none';
            }
        }
    });
}

window.dataHandlers = {};

function registerDataHandler(identifier, testFn, handlerFn) {
    window.dataHandlers[identifier] = {
        test: testFn,     
        handler: handlerFn 
    };
}

function processWithHandler(data) {
    document.getElementById('content').innerHTML = '';
    const chartsContainer = document.getElementById('charts-container');
    if (chartsContainer) {
        chartsContainer.style.display = 'none';
    }
    
    let handlerFound = false;

    for (const [id, handler] of Object.entries(window.dataHandlers)) {
        if (handler.test(data)) {
            console.log(`Processing data with ${id} handler`);
            handler.handler(data);
            handlerFound = true;
            break;
        }
    }
    
    if (!handlerFound) {
        console.log("No handler found for this data format");
        displayUnknownDataFormat(data);
    }
    
    setTimeout(() => {
        const tables = document.querySelectorAll('.collapsible-table');
        console.log(`Found ${tables.length} tables`);
        updateResourceNav();
    }, 200);
}

function displayUnknownDataFormat(data) {
    const content = document.getElementById('content');
    content.innerHTML = `
        <div class="unknown-data">
            <h2>No Platform Selected</h2>
            <p>Select available platforms, or import your JSON</p>
            <p>Available data keys, if any: ${Object.keys(data).join(', ')}</p>
        </div>
    `;
}

function loadTestJson(platform) {
    showLoadingIndicator();
    
    const fileMap = {
        'kubernetes': 'k8s',
        'aws': 'aws',
        'azure': 'azure',
        'gcp': 'gcp',
        'veeam': 'veeam'
    };
    
    const filePrefix = fileMap[platform] || platform;
    
    fetch(`/test/${filePrefix}.json`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log(`Loaded test data from ${filePrefix}.json`);
            processWithHandler(data);
        })
        .catch(error => {
            console.error(`Error loading test data (${filePrefix}.json):`, error);
            document.getElementById('content').innerHTML = `
                <div class="error-message">
                    <h2>Error Loading Test Data</h2>
                    <p>${error.message}</p>
                </div>
            `;
        })
        .finally(() => {
            hideLoadingIndicator();
        });
}

function loadTerraformState(source) {
    showLoadingIndicator();
    
    if (source.startsWith('s3://')) {
        fetch('/api/terraform/s3-state', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ path: source })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            processWithHandler(data);
        })
        .catch(error => {
            console.error(`Error loading Terraform state from S3: ${error}`);
            document.getElementById('content').innerHTML = `
                <div class="error-message">
                    <h2>Error Loading Terraform State</h2>
                    <p>${error.message}</p>
                </div>
            `;
        })
        .finally(() => hideLoadingIndicator());
    }
}

function enhanceConnectionButtons() {
    console.log("Enhancing connection buttons with loading indicators");
    
    const awsConnectBtn = document.getElementById('aws-connect-btn');
    if (awsConnectBtn) {
        const originalClickHandler = awsConnectBtn.onclick;
        awsConnectBtn.onclick = function(event) {
            showLoadingIndicator();
            if (originalClickHandler) {
                originalClickHandler.call(this, event);
            }
        };
    }
    
    const azureConnectBtn = document.getElementById('azure-connect-btn');
    if (azureConnectBtn) {
        const originalClickHandler = azureConnectBtn.onclick;
        azureConnectBtn.onclick = function(event) {
            showLoadingIndicator();
            if (originalClickHandler) {
                originalClickHandler.call(this, event);
            }
        };
    }
    
    const gcpConnectBtn = document.getElementById('gcp-connect-btn');
    if (gcpConnectBtn) {
        const originalClickHandler = gcpConnectBtn.onclick;
        gcpConnectBtn.onclick = function(event) {
            showLoadingIndicator();
            if (originalClickHandler) {
                originalClickHandler.call(this, event);
            }
        };
    }
    
    const k8sConnectBtn = document.getElementById('kubernetes-connect-btn');
    if (k8sConnectBtn) {
        const originalClickHandler = k8sConnectBtn.onclick;
        k8sConnectBtn.onclick = function(event) {
            showLoadingIndicator();
            if (originalClickHandler) {
                originalClickHandler.call(this, event);
            }
        };
    }
    
    const tfConnectBtns = document.querySelectorAll('[id$="-terraform-connect-btn"]');
    tfConnectBtns.forEach(btn => {
        const originalClickHandler = btn.onclick;
        btn.onclick = function(event) {
            showLoadingIndicator();
            if (originalClickHandler) {
                originalClickHandler.call(this, event);
            }
        };
    });
}

document.addEventListener('htmx:afterSwap', (event) => {
    if (event.detail.target.id === 'hidden-content') {
        try {
            const data = JSON.parse(event.detail.xhr.responseText);
            console.log("Fetched Data:", Object.keys(data));
            processWithHandler(data);
        } catch (error) {
            console.error("Error processing data:", error);
        } finally {
            hideLoadingIndicator();
        }
    }
});

document.addEventListener('DOMContentLoaded', () => {
    console.log("DOM Loaded - Setting up event listeners");
    
    const toggleTablesButton = document.getElementById('toggle-tables');
    if (toggleTablesButton) {
        console.log("Found toggle-tables button");
        toggleTablesButton.addEventListener('click', () => {
            const isCollapsed = toggleTablesButton.getAttribute('data-collapsed') === 'true';
            console.log(`Toggle tables clicked, current state: ${isCollapsed ? 'Collapsed' : 'Expanded'}`);
            
            if (isCollapsed) {
                toggleAllTables(false);
            } else {
                toggleAllTables(true);
            }
        });
    }
    
    const navToggle = document.getElementById('resource-nav-toggle');
    if (navToggle) {
        console.log("Found nav toggle button");
        navToggle.addEventListener('click', () => {
            const nav = document.getElementById('resource-nav');
            if (nav) {
                nav.style.display = nav.style.display === 'none' ? 'block' : 'none';
            }
        });
    }
    
    document.getElementById('veeam-button')?.addEventListener('click', () => {
        showLoadingIndicator();
        fetch('/api/switch?type=veeam')
            .then(response => response.json())
            .then(data => {
                location.reload();
            })
            .catch(error => console.error('Error switching to Veeam:', error))
            .finally(() => hideLoadingIndicator());
    });
    
    const bodyObserver = new MutationObserver(function(mutations) {
        mutations.forEach(function(mutation) {
            if (mutation.addedNodes && mutation.addedNodes.length > 0) {
                for (let i = 0; i < mutation.addedNodes.length; i++) {
                    const node = mutation.addedNodes[i];
                    if (node.id && node.id.endsWith('-connect-btn')) {
                        console.log(`Detected new connection button: ${node.id}`);
                        enhanceConnectionButtons();
                        break;
                    }
                }
            }
        });
    });
    
    bodyObserver.observe(document.body, {
        childList: true,
        subtree: true
    });
    
    enhanceConnectionButtons();
});

document.getElementById('export-button')?.addEventListener('click', () => {
    showLoadingIndicator();
    fetch('/api/data')
        .then(response => response.json())
        .then(data => {
            const dataStr = JSON.stringify(data, null, 2);
            const dataUri = 'data:application/json;charset=utf-8,'+ encodeURIComponent(dataStr);
            
            const exportFileDefaultName = 'kollect_data.json';
            
            const linkElement = document.createElement('a');
            linkElement.setAttribute('href', dataUri);
            linkElement.setAttribute('download', exportFileDefaultName);
            linkElement.click();
        })
        .catch(error => console.error('Error exporting data:', error))
        .finally(() => hideLoadingIndicator());
});

document.getElementById('import-button')?.addEventListener('click', () => {
    document.getElementById('import-file').click();
});

document.getElementById('import-file')?.addEventListener('change', (event) => {
    const file = event.target.files[0];
    if (file) {
        showLoadingIndicator();
        const reader = new FileReader();
        reader.onload = (e) => {
            try {
                const data = JSON.parse(e.target.result);
                console.log("Imported data:", Object.keys(data));
                processWithHandler(data);
                hideLoadingIndicator();
            } catch (error) {
                console.error("Error parsing imported data:", error);
                document.getElementById('content').innerHTML = `
                    <div class="error-message">
                        <h2>Error Importing Data</h2>
                        <p>${error.message}</p>
                    </div>
                `;
                hideLoadingIndicator();
            }
        };
        reader.readAsText(file);
    }
});