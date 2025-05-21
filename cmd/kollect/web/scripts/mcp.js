// MCP (Model Context Protocol) script for AI Assistant features

// Check if MCP is active/enabled
function checkMCPStatus() {
    fetch('/api/mcp/status')
        .then(response => response.json())
        .then(data => {
            if (data.enabled) {
                console.log('MCP is active, showing AI Assistant interface');
                showMCPCard();
            } else {
                console.log('MCP is not active');
            }
        })
        .catch(error => {
            console.error('Error checking MCP status:', error);
        });
}

// Create the MCP UI card
function showMCPCard() {
    const content = document.getElementById('content');
    
    // Check if the card already exists
    if (document.getElementById('mcp-card')) {
        return;
    }
    
    // Create the AI Assistant card
    const mcpCard = document.createElement('div');
    mcpCard.id = 'mcp-card';
    mcpCard.className = 'mcp-card';
    mcpCard.innerHTML = `
        <div class="card" style="margin: 20px 0; padding: 20px; border-radius: 8px; background-color: var(--card-bg); box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
            <h3><i class="fas fa-robot"></i> AI Assistant Access</h3>
            <p>Your infrastructure data is available to AI assistants via the Model Context Protocol.</p>
            
            <div style="margin-top: 15px;">
                <div class="form-group">
                    <label for="mcp-endpoint" style="font-weight: bold;">MCP Endpoint:</label>
                    <input type="text" id="mcp-endpoint" value="http://${window.location.host}/api/mcp/retrieve" readonly 
                           style="width: 100%; padding: 8px; margin-top: 5px; background: var(--secondary-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;">
                </div>
                
                <div style="margin-top: 15px;">
                    <button id="copy-mcp-endpoint" class="btn btn-primary">
                        <i class="fas fa-copy"></i> Copy Endpoint
                    </button>
                    <button id="test-mcp-query" class="btn">
                        <i class="fas fa-search"></i> Test Query
                    </button>
                </div>
            </div>
            
            <div id="mcp-test-container" style="display: none; margin-top: 15px;">
                <div class="form-group">
                    <label for="mcp-test-query" style="font-weight: bold;">Test Query:</label>
                    <input type="text" id="mcp-test-query" placeholder="Enter a query about your infrastructure..." 
                           style="width: 100%; padding: 8px; margin-top: 5px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;">
                </div>
                <button id="run-mcp-test" class="btn btn-primary" style="margin-top: 10px;">
                    <i class="fas fa-play"></i> Run Query
                </button>
                
                <div id="mcp-results" style="margin-top: 15px; max-height: 300px; overflow-y: auto; display: none;">
                </div>
            </div>
        </div>
    `;
    
    // Insert at the top of the content
    if (content.firstChild) {
        content.insertBefore(mcpCard, content.firstChild);
    } else {
        content.appendChild(mcpCard);
        checkMCPIndexStatus();
    }
    
    setupMCPEventHandlers();
}

function setupMCPEventHandlers() {
    // Event handlers for the MCP card
    document.getElementById('copy-mcp-endpoint')?.addEventListener('click', () => {
        const endpoint = document.getElementById('mcp-endpoint').value;
        navigator.clipboard.writeText(endpoint).then(() => {
            showToast('Endpoint copied to clipboard');
        }).catch(err => {
            console.error('Could not copy text: ', err);
        });
    });
    
    document.getElementById('test-mcp-query')?.addEventListener('click', () => {
        const container = document.getElementById('mcp-test-container');
        container.style.display = container.style.display === 'none' ? 'block' : 'none';
    });
    
    document.getElementById('run-mcp-test')?.addEventListener('click', () => {
        const query = document.getElementById('mcp-test-query').value;
        if (!query) {
            showToast('Please enter a query', 'warning');
            return;
        }
        
        document.getElementById('mcp-results').style.display = 'none';
        showLoadingIndicator();
        checkMCPIndexStatus();
        
        fetch('/api/mcp/retrieve', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                query: query,
                limit: 5
            })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            const resultsContainer = document.getElementById('mcp-results');
            resultsContainer.innerHTML = '';
            
            if (data.documents && data.documents.length > 0) {
                const resultList = document.createElement('div');
                resultList.className = 'mcp-results-list';
                resultList.style.display = 'flex';
                resultList.style.flexDirection = 'column';
                resultList.style.gap = '10px';
                
                data.documents.forEach(doc => {
                    const docElement = document.createElement('div');
                    docElement.className = 'mcp-result-item';
                    docElement.style.padding = '10px';
                    docElement.style.margin = '5px 0';
                    docElement.style.borderLeft = '3px solid var(--accent-color)';
                    docElement.style.backgroundColor = 'var(--secondary-bg-color)';
                    
                    docElement.innerHTML = `
                        <div><strong>Source:</strong> ${doc.source} (${doc.source_type})</div>
                        <div style="margin-top: 5px; white-space: pre-wrap;">${doc.content}</div>
                        ${doc.score ? `<div style="margin-top: 5px; font-size: 0.9em; color: var(--secondary-text-color);">Relevance: ${(doc.score * 100).toFixed(1)}%</div>` : ''}
                    `;
                    
                    resultList.appendChild(docElement);
                });
                
                resultsContainer.appendChild(resultList);
            } else {
                resultsContainer.innerHTML = '<div style="padding: 10px; color: var(--secondary-text-color);">No results found for your query.</div>';
            }
            
            resultsContainer.style.display = 'block';
            hideLoadingIndicator();
        })
        .catch(error => {
            console.error('Error testing MCP query:', error);
            hideLoadingIndicator();
            document.getElementById('mcp-results').innerHTML = `<div style="padding: 10px; color: #FF5252;">Error: ${error.message}</div>`;
            document.getElementById('mcp-results').style.display = 'block';
        });
    });
}

// Helper functions
function showToast(message, type = 'info') {
    // Check if toast container exists
    let toastContainer = document.getElementById('toast-container');
    
    if (!toastContainer) {
        toastContainer = document.createElement('div');
        toastContainer.id = 'toast-container';
        toastContainer.style.position = 'fixed';
        toastContainer.style.bottom = '20px';
        toastContainer.style.right = '20px';
        toastContainer.style.zIndex = '9999';
        document.body.appendChild(toastContainer);
    }
    
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    toast.style.backgroundColor = type === 'warning' ? '#ffdd57' : 'var(--accent-color)';
    toast.style.color = type === 'warning' ? '#000' : '#fff';
    toast.style.padding = '12px 20px';
    toast.style.borderRadius = '4px';
    toast.style.marginTop = '10px';
    toast.style.boxShadow = '0 2px 10px rgba(0,0,0,0.2)';
    toast.style.minWidth = '200px';
    toast.style.textAlign = 'center';
    toast.textContent = message;
    
    toastContainer.appendChild(toast);
    
    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transition = 'opacity 0.5s ease';
        setTimeout(() => {
            toast.remove();
            if (toastContainer.childNodes.length === 0) {
                toastContainer.remove();
            }
        }, 500);
    }, 3000);
}

// Watch for changes in the content
function setupMCPObserver() {
    // Create a MutationObserver to watch for changes in content
    const observer = new MutationObserver((mutations) => {
        for (const mutation of mutations) {
            if (mutation.type === 'childList' && !document.getElementById('mcp-card')) {
                showMCPCard();
            }
        }
    });
    
    // Start observing the content with config
    const content = document.getElementById('content');
    if (content) {
        observer.observe(content, { childList: true, subtree: true });
    }
}

// Add this function to check the MCP index status
function checkMCPIndexStatus() {
    console.log("Checking MCP index status...");
    fetch('/api/mcp/status-info')
        .then(response => response.json())
        .then(data => {
            console.log("MCP indexed documents info:", data);
            
            // Add a debug section to the MCP card
            const mcpCard = document.getElementById('mcp-card');
            if (mcpCard) {
                let debugInfo = document.getElementById('mcp-debug-info');
                if (!debugInfo) {
                    debugInfo = document.createElement('div');
                    debugInfo.id = 'mcp-debug-info';
                    debugInfo.style.marginTop = '15px';
                    debugInfo.style.padding = '10px';
                    debugInfo.style.backgroundColor = 'rgba(0,0,0,0.05)';
                    debugInfo.style.borderRadius = '4px';
                    debugInfo.style.fontSize = '0.9em';
                    debugInfo.style.fontFamily = 'monospace';
                    mcpCard.querySelector('.card').appendChild(debugInfo);
                }
                
                debugInfo.innerHTML = `
                    <div style="font-weight: bold; margin-bottom: 5px;">MCP Debug Info:</div>
                    <div>Enabled: ${data.enabled}</div>
                    <div>Document Count: ${data.docCount || 0}</div>
                    <div>Document Types: ${JSON.stringify(data.docTypes || [])}</div>
                    <div>Query Engine: ${data.engineType || 'Unknown'}</div>
                    <button id="mcp-refresh-debug" style="margin-top: 10px;" class="btn btn-sm">Refresh Debug Info</button>
                `;
                
                document.getElementById('mcp-refresh-debug')?.addEventListener('click', () => {
                    checkMCPIndexStatus();
                });
            }
            
            // If no documents, show a warning
            if (!data.docCount || data.docCount === 0) {
                showToast("No documents have been indexed by MCP yet!", "warning");
            }
        })
        .catch(error => {
            console.error("Error checking MCP index status:", error);
        });
}

// Initialize when page loads
document.addEventListener('DOMContentLoaded', function() {
    console.log('MCP script loaded');
    checkMCPStatus();
    setupMCPObserver();
    
    // Also listen for data load events
    window.addEventListener('dataLoaded', function() {
        console.log('Data loaded event received by MCP');
        showMCPCard();
    });
    
    // Re-add the card when someone imports data
    window.addEventListener('dataImported', function() {
        console.log('Data imported event received by MCP');
        showMCPCard();
    });
});

// Add a dedicated method for showing the card programmatically 
// (to be called from other scripts)
window.showMCPInterface = function() {
    showMCPCard();
};