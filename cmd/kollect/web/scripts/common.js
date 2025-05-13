// common.js

function showLoadingIndicator() {
    document.getElementById('loading-indicator').style.display = 'flex';
}

function hideLoadingIndicator() {
    document.getElementById('loading-indicator').style.display = 'none';
}

// Function to create collapsible tables
function createTable(headerText, data, rowTemplate, headers) {
    if (!data || data.length === 0) return;
    
    // Create a unique ID for the table based on the header text
    const tableId = headerText.replace(/\s+/g, '-').toLowerCase();
    
    // Create the table container
    const tableContainer = document.createElement('div');
    tableContainer.className = 'collapsible-table';
    tableContainer.id = `table-container-${tableId}`;
    
    // Create the table header with collapse functionality
    const tableHeader = document.createElement('div');
    tableHeader.className = 'table-header collapsed'; // Start collapsed
    tableHeader.innerHTML = `
        <span>${headerText}</span>
        <div>
            <span class="table-counter">${data.length}</span>
            <span class="icon">▼</span>
        </div>
    `;
    
    // Create the table content area
    const tableContent = document.createElement('div');
    tableContent.className = 'table-content collapsed'; // Start collapsed
    
    // Create table
    const table = document.createElement('table');
    
    // Create header
    const thead = document.createElement('thead');
    const headerRow = document.createElement('tr');
    
    headers.forEach(header => {
        const th = document.createElement('th');
        th.textContent = header;
        headerRow.appendChild(th);
    });
    
    thead.appendChild(headerRow);
    table.appendChild(thead);
    
    // Create body
    const tbody = document.createElement('tbody');
    data.forEach(item => {
        const row = document.createElement('tr');
        row.innerHTML = rowTemplate(item);
        tbody.appendChild(row);
    });
    table.appendChild(tbody);
    
    // Add click handler for toggling
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

// Function to toggle all tables
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
}

// Function to update the resource navigation panel
function updateResourceNav() {
    const resourceNav = document.getElementById('resource-nav');
    if (!resourceNav) return;
    
    resourceNav.innerHTML = '';
    
    // Get all collapsible tables
    const tables = document.querySelectorAll('.collapsible-table');
    
    // Only show the navigation panel if we have tables
    if (tables.length === 0) {
        resourceNav.style.display = 'none';
        const navToggle = document.getElementById('resource-nav-toggle');
        if (navToggle) navToggle.style.display = 'none';
        return;
    } else {
        const navToggle = document.getElementById('resource-nav-toggle');
        if (navToggle) navToggle.style.display = 'block';
    }
    
    // Add a heading
    const heading = document.createElement('h4');
    heading.textContent = 'Resources';
    heading.style.margin = '0 0 10px 0';
    heading.style.padding = '0';
    heading.style.textAlign = 'center';
    resourceNav.appendChild(heading);
    
    // Add a link for each table
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
            
            // Expand the clicked table
            const content = table.querySelector('.table-content');
            const headerElement = table.querySelector('.table-header');
            headerElement.classList.remove('collapsed');
            content.classList.remove('collapsed');
        };
        
        resourceNav.appendChild(link);
    });
    
    console.log("Navigation updated");
}

// Then improve the toggle button functionality in the DOMContentLoaded event:

// Setup event listener for resource navigation toggle
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

// Create a global map to register data handlers
window.dataHandlers = {};

// Register a handler for a specific data format
function registerDataHandler(identifier, testFn, handlerFn) {
    window.dataHandlers[identifier] = {
        test: testFn,     // Function that tests if data matches this handler
        handler: handlerFn // Function that processes the data
    };
}

// Process data with the appropriate handler
function processWithHandler(data) {
    // First clear previous content
    document.getElementById('content').innerHTML = '';
    
    // Hide charts container by default
    const chartsContainer = document.getElementById('charts-container');
    if (chartsContainer) {
        chartsContainer.style.display = 'none';
    }
    
    // Find a handler that recognizes this data
    let handlerFound = false;
    
    // Try each registered handler
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
    
    // Update navigation regardless of which handler was used
    setTimeout(() => {
        const tables = document.querySelectorAll('.collapsible-table');
        console.log(`Found ${tables.length} tables`);
        updateResourceNav();
    }, 200);
}

// Display a message for unknown data formats
function displayUnknownDataFormat(data) {
    const content = document.getElementById('content');
    content.innerHTML = `
        <div class="unknown-data">
            <h2>Unknown Data Format</h2>
            <p>The imported data doesn't match any known platform format.</p>
            <p>Available data keys: ${Object.keys(data).join(', ')}</p>
        </div>
    `;
}

// Function to load JSON test data for development purposes
function loadTestJson(platform) {
    showLoadingIndicator();
    
    // Map platform names to file prefixes
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

// Update HTMX event listener to use the handler system
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

// Add this to the DOMContentLoaded event
document.addEventListener('DOMContentLoaded', () => {
    console.log("DOM Loaded - Setting up event listeners");
    
    // Setup event listeners for expand/collapse all buttons
    const expandAllButton = document.getElementById('expand-all');
    if (expandAllButton) {
        console.log("Found expand-all button");
        expandAllButton.addEventListener('click', () => {
            console.log("Expand all clicked");
            toggleAllTables(false);
        });
    }
    
    const collapseAllButton = document.getElementById('collapse-all');
    if (collapseAllButton) {
        console.log("Found collapse-all button");
        collapseAllButton.addEventListener('click', () => {
            console.log("Collapse all clicked");
            toggleAllTables(true);
        });
    }
    
    // Setup event listener for resource navigation toggle
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
    
    // Setup event listeners for platform buttons to call the API endpoints
    document.getElementById('aws-button')?.addEventListener('click', () => {
        showLoadingIndicator();
        fetch('/api/switch?type=aws')
            .then(response => response.json())
            .then(data => {
                location.reload();
            })
            .catch(error => console.error('Error switching to AWS:', error))
            .finally(() => hideLoadingIndicator());
    });
    
    document.getElementById('kubernetes-button')?.addEventListener('click', () => {
        showLoadingIndicator();
        fetch('/api/switch?type=kubernetes')
            .then(response => response.json())
            .then(data => {
                location.reload();
            })
            .catch(error => console.error('Error switching to Kubernetes:', error))
            .finally(() => hideLoadingIndicator());
    });
    
    document.getElementById('azure-button')?.addEventListener('click', () => {
        showLoadingIndicator();
        fetch('/api/switch?type=azure')
            .then(response => response.json())
            .then(data => {
                location.reload();
            })
            .catch(error => console.error('Error switching to Azure:', error))
            .finally(() => hideLoadingIndicator());
    });
    
    document.getElementById('google-button')?.addEventListener('click', () => {
        showLoadingIndicator();
        fetch('/api/switch?type=gcp')
            .then(response => response.json())
            .then(data => {
                location.reload();
            })
            .catch(error => console.error('Error switching to GCP:', error))
            .finally(() => hideLoadingIndicator());
    });
    
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
});

// Handle import/export
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
                // Process the data directly
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