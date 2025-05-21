// MCP (Model Context Protocol) script for AI Assistant features

function checkMCPStatus() {
    fetch('/api/mcp/status')
        .then(response => response.json())
        .then(data => {
            if (data.enabled) {
                console.log('MCP is active, adding AI Assistant button');
                addMCPButton();
            } else {
                console.log('MCP is not active');
                // Make sure to hide the button if it exists and MCP is disabled
                const existingButton = document.getElementById('ai-assistant-button');
                if (existingButton) {
                    existingButton.style.display = 'none';
                }
            }
        })
        .catch(error => {
            console.error('Error checking MCP status:', error);
        });
}

// Add the AI Assistant button next to the cost explorer button
function addMCPButton() {
    let aiButton = document.getElementById('ai-assistant-button');
    if (!aiButton) {
        const detailsButtons = document.querySelector('.details-buttons');
        if (detailsButtons) {
            aiButton = document.createElement('button');
            aiButton.id = 'ai-assistant-button';
            aiButton.className = 'utility-button';
            aiButton.title = 'AI Assistant - Query your infrastructure data';
            aiButton.innerHTML = '<i class="fas fa-robot" style="font-size: 24px;"></i>';
            aiButton.addEventListener('click', function(event) {
                event.preventDefault();
                showMCPModal();
            });
            detailsButtons.appendChild(aiButton);
        }
    } else {
        aiButton.style.display = 'block';
    }
}

// Create the MCP modal dialog
function showMCPModal() {
    // First check if MCP is enabled
    fetch('/api/mcp/status')
        .then(response => response.json())
        .then(data => {
            if (!data.enabled) {
                showToast('AI Assistant is not enabled. Run with --mcp flag to enable it.', 'warning');
                return;
            }
            
            const modal = document.createElement('div');
            modal.className = 'modal';
            modal.style.position = 'fixed';
            modal.style.top = '0';
            modal.style.left = '0';
            modal.style.width = '100%';
            modal.style.height = '100%';
            modal.style.backgroundColor = 'rgba(0, 0, 0, 0.5)';
            modal.style.display = 'flex';
            modal.style.justifyContent = 'center';
            modal.style.alignItems = 'center';
            modal.style.zIndex = '1000';
            
            const modalContent = document.createElement('div');
            modalContent.className = 'mcp-modal';
            modalContent.style.backgroundColor = 'var(--card-bg)';
            modalContent.style.padding = '25px';
            modalContent.style.borderRadius = '8px';
            modalContent.style.maxWidth = '700px';
            modalContent.style.width = '90%';
            modalContent.style.maxHeight = '90vh';
            modalContent.style.overflowY = 'auto';
            modalContent.style.position = 'relative';
            modalContent.style.boxShadow = '0 5px 20px rgba(0,0,0,0.4)';
            modalContent.style.border = '1px solid var(--border-color)';
            
            modalContent.innerHTML = `
                <h3 style="margin-top: 0; color: var(--accent-color); font-size: 1.5em; border-bottom: 1px solid var(--border-color); padding-bottom: 10px;">
                    <i class="fas fa-robot"></i> AI Assistant
                </h3>
                
                <p style="margin-bottom: 20px;">Your infrastructure data is available to AI assistants via the Model Context Protocol.</p>
                
                <div class="form-group" style="margin-bottom: 20px;">
                    <label for="mcp-endpoint-modal" style="font-weight: bold; display: block; margin-bottom: 5px;">MCP Endpoint:</label>
                    <div style="display: flex; gap: 10px;">
                        <input type="text" id="mcp-endpoint-modal" value="http://${window.location.host}/api/mcp/retrieve" readonly 
                            style="flex-grow: 1; padding: 8px; background: var(--secondary-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;">
                        <button id="copy-mcp-endpoint-modal" class="btn btn-primary" style="white-space: nowrap;">
                            <i class="fas fa-copy"></i> Copy
                        </button>
                    </div>
                </div>
                
                <div style="margin-bottom: 20px;">
                    <h4 style="margin-top: 0; margin-bottom: 10px; border-bottom: 1px solid var(--border-color); padding-bottom: 5px;">
                        <i class="fas fa-search"></i> Test Query
                    </h4>
                    <div class="form-group" style="margin-bottom: 10px;">
                        <input type="text" id="mcp-test-query-modal" placeholder="Enter a query about your infrastructure..." 
                            style="width: 100%; padding: 8px; background: var(--input-bg-color); color: var(--text-color); border: 1px solid var(--border-color); border-radius: 4px;">
                    </div>
                    <button id="run-mcp-test-modal" class="btn btn-primary">
                        <i class="fas fa-play"></i> Run Query
                    </button>
                </div>
                
                <div id="mcp-results-modal" style="margin-top: 15px; max-height: 300px; overflow-y: auto; display: none;">
                </div>
                
                <div id="mcp-debug-info-modal" style="margin-top: 15px; padding: 10px; background-color: rgba(0,0,0,0.05); border-radius: 4px; font-size: 0.9em; font-family: monospace;">
                    <div style="font-weight: bold; margin-bottom: 5px;">MCP Debug Info:</div>
                    <div>Loading...</div>
                </div>
                
                <div class="modal-buttons" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 25px; border-top: 1px solid var(--border-color); padding-top: 15px;">
                    <button id="mcp-close-btn" class="btn" style="padding: 10px 20px; background-color: var(--button-bg-color); color: var(--button-text-color); border: none; border-radius: 4px; cursor: pointer;">Close</button>
                </div>
            `;
            
            modal.appendChild(modalContent);
            document.body.appendChild(modal);
            
            // Add event handlers
            document.getElementById('copy-mcp-endpoint-modal')?.addEventListener('click', () => {
                const endpoint = document.getElementById('mcp-endpoint-modal').value;
                navigator.clipboard.writeText(endpoint).then(() => {
                    showToast('Endpoint copied to clipboard');
                }).catch(err => {
                    console.error('Could not copy text: ', err);
                });
            });
            
            document.getElementById('run-mcp-test-modal')?.addEventListener('click', () => {
                const query = document.getElementById('mcp-test-query-modal').value;
                if (!query) {
                    showToast('Please enter a query', 'warning');
                    return;
                }
                
                const resultsContainer = document.getElementById('mcp-results-modal');
                resultsContainer.style.display = 'none';
                showMCPLoadingIndicator();
                
                fetch('/api/mcp/retrieve', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        query: query,
                        limit: 10
                    })
                })
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`HTTP error: ${response.status}`);
                    }
                    return response.json();
                })
                .then(data => {
                    console.log("MCP query response:", data);
                    
                    resultsContainer.innerHTML = '';
                    
                    if (data && data.documents && data.documents.length > 0) {
                        // Results rendering
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
                                <div><strong>Source:</strong> ${doc.source || 'Unknown'} (${doc.source_type || 'Unknown'})</div>
                                <div style="margin-top: 5px; white-space: pre-wrap;">${doc.content || 'No content'}</div>
                                ${doc.score ? `<div style="margin-top: 5px; font-size: 0.9em; color: var(--secondary-text-color);">Relevance: ${(doc.score * 100).toFixed(1)}%</div>` : ''}
                            `;
                            
                            resultList.appendChild(docElement);
                        });
                        
                        resultsContainer.appendChild(resultList);
                    } else if (Array.isArray(data) && data.length > 0) {
                        // Alternative format handling
                        const resultList = document.createElement('div');
                        resultList.className = 'mcp-results-list';
                        resultList.style.display = 'flex';
                        resultList.style.flexDirection = 'column';
                        resultList.style.gap = '10px';
                        
                        data.forEach(doc => {
                            const docElement = document.createElement('div');
                            docElement.className = 'mcp-result-item';
                            docElement.style.padding = '10px';
                            docElement.style.margin = '5px 0';
                            docElement.style.borderLeft = '3px solid var(--accent-color)';
                            docElement.style.backgroundColor = 'var(--secondary-bg-color)';
                            
                            docElement.innerHTML = `
                                <div><strong>Source:</strong> ${doc.source || 'Unknown'} (${doc.source_type || 'Unknown'})</div>
                                <div style="margin-top: 5px; white-space: pre-wrap;">${doc.content || 'No content'}</div>
                                ${doc.score ? `<div style="margin-top: 5px; font-size: 0.9em; color: var(--secondary-text-color);">Relevance: ${(doc.score * 100).toFixed(1)}%</div>` : ''}
                            `;
                            
                            resultList.appendChild(docElement);
                        });
                        
                        resultsContainer.appendChild(resultList);
                    } else {
                        resultsContainer.innerHTML = `
                            <div style="padding: 10px; color: var(--secondary-text-color);">
                                <p>No results found for your query: "${query}"</p>
                                <p>Try using one of these terms instead:</p>
                                <ul style="margin-top: 5px; margin-left: 20px;">
                                    <li>kubernetes</li>
                                    <li>pod</li>
                                    <li>nodes</li>
                                    <li>storage</li>
                                    <li>aws</li>
                                    <li>azure</li>
                                    <li>gcp</li>
                                </ul>
                            </div>`;
                    }
                    
                    resultsContainer.style.display = 'block';
                    hideMCPLoadingIndicator();
                })
                .catch(error => {
                    console.error('Error testing MCP query:', error);
                    hideMCPLoadingIndicator();
                    resultsContainer.innerHTML = `
                        <div style="padding: 10px; color: #FF5252;">
                            <p>Error: ${error.message}</p>
                            <p style="margin-top: 10px;">This may happen if:</p>
                            <ul style="margin-top: 5px; margin-left: 20px;">
                                <li>The MCP service is not initialized correctly</li>
                                <li>There are no documents indexed yet</li>
                                <li>The query format is incorrect</li>
                            </ul>
                        </div>`;
                    resultsContainer.style.display = 'block';
                });
            });
            
            document.getElementById('mcp-close-btn')?.addEventListener('click', () => {
                modal.remove();
            });
            
            // Load MCP debug info
            updateMCPModalDebugInfo();
        })
        .catch(error => {
            console.error('Error checking MCP status:', error);
            showToast('Error checking AI Assistant status', 'warning');
        });
}

// Helper function to update the MCP debug info in the modal
function updateMCPModalDebugInfo() {
    const debugInfoContainer = document.getElementById('mcp-debug-info-modal');
    if (!debugInfoContainer) return;
    
    fetch('/api/mcp/status-info')
        .then(response => response.json())
        .then(data => {
            // Format document types to display in a more compact way
            let docTypes = '';
            if (data.docTypes && data.docTypes.length > 0) {
                // Group by main type 
                const typeGroups = {};
                data.docTypes.forEach(type => {
                    const mainType = type.split(':')[0];
                    if (!typeGroups[mainType]) {
                        typeGroups[mainType] = [];
                    }
                    typeGroups[mainType].push(type.split(':')[1]);
                });
                
                // Format as a list of main types with subtypes
                docTypes = Object.entries(typeGroups)
                    .map(([mainType, subtypes]) => 
                        `<div style="margin-bottom:5px;">
                            <strong>${mainType}:</strong> 
                            <span style="word-break:break-all;">${subtypes.join(', ')}</span>
                        </div>`
                    ).join('');
            } else {
                docTypes = "None";
            }
            
            debugInfoContainer.innerHTML = `
                <div style="font-weight: bold; margin-bottom: 5px;">MCP Debug Info:</div>
                <div>Enabled: ${data.enabled}</div>
                <div>Document Count: ${data.docCount || 0}</div>
                <div style="margin-top: 5px;">Document Types:</div>
                <div style="max-width: 100%; margin-left: 10px;">${docTypes}</div>
                <div>Query Engine: ${data.engineType || 'Unknown'}</div>
                <button id="mcp-refresh-debug-modal" style="margin-top: 10px;" class="btn btn-sm">Refresh Debug Info</button>
            `;
            
            document.getElementById('mcp-refresh-debug-modal')?.addEventListener('click', () => {
                updateMCPModalDebugInfo();
            });
            
            // If no documents, show a warning
            if (!data.docCount || data.docCount === 0) {
                showToast("No documents have been indexed by MCP yet!", "warning");
            }
        })
        .catch(error => {
            console.error('Error fetching MCP debug info:', error);
            debugInfoContainer.innerHTML = `
                <div style="font-weight: bold; margin-bottom: 5px;">MCP Debug Info:</div>
                <div style="color: #FF5252;">Error fetching information: ${error.message}</div>
                <button id="mcp-refresh-debug-modal" style="margin-top: 10px;" class="btn btn-sm">Retry</button>
            `;
            
            document.getElementById('mcp-refresh-debug-modal')?.addEventListener('click', () => {
                updateMCPModalDebugInfo();
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

// Loading indicator functions specific to MCP
function showMCPLoadingIndicator() {
    // For MCP UI
    const resultsContainer = document.getElementById('mcp-results-modal');
    if (resultsContainer) {
        // We're in the MCP UI
        let loader = document.getElementById('mcp-loader');
        if (!loader) {
            loader = document.createElement('div');
            loader.id = 'mcp-loader';
            loader.style.position = 'relative';
            loader.style.padding = '20px';
            loader.style.textAlign = 'center';
            loader.innerHTML = '<div class="spinner" style="display: inline-block; width: 30px; height: 30px; border: 3px solid rgba(0,0,0,0.1); border-radius: 50%; border-top-color: var(--accent-color); animation: spin 1s ease-in-out infinite;"></div>';
            loader.innerHTML += '<style>@keyframes spin { to { transform: rotate(360deg); } }</style>';
            
            resultsContainer.innerHTML = '';
            resultsContainer.appendChild(loader);
            resultsContainer.style.display = 'block';
        }
    } else {
        // Use the common loading indicator
        if (typeof window.showLoadingIndicator === 'function') {
            window.showLoadingIndicator();
        } else {
            // Fallback
            let globalLoader = document.getElementById('loading-indicator');
            if (globalLoader) {
                globalLoader.style.display = 'flex';
            } else {
                // Create one if needed
                globalLoader = document.createElement('div');
                globalLoader.id = 'loading-indicator';
                globalLoader.style.position = 'fixed';
                globalLoader.style.top = '0';
                globalLoader.style.left = '0';
                globalLoader.style.width = '100%';
                globalLoader.style.height = '100%';
                globalLoader.style.backgroundColor = 'rgba(0,0,0,0.5)';
                globalLoader.style.display = 'flex';
                globalLoader.style.justifyContent = 'center';
                globalLoader.style.alignItems = 'center';
                globalLoader.style.zIndex = '9999';
                
                const spinner = document.createElement('div');
                spinner.style.width = '50px';
                spinner.style.height = '50px';
                spinner.style.border = '5px solid rgba(255,255,255,0.3)';
                spinner.style.borderRadius = '50%';
                spinner.style.borderTop = '5px solid var(--accent-color, #4CAF50)';
                spinner.style.animation = 'spin 1s linear infinite';
                
                globalLoader.appendChild(spinner);
                globalLoader.innerHTML += '<style>@keyframes spin { to { transform: rotate(360deg); } }</style>';
                
                document.body.appendChild(globalLoader);
            }
            
            globalLoader.style.display = 'flex';
        }
    }
}

function hideMCPLoadingIndicator() {
    // Hide MCP-specific loader
    const loader = document.getElementById('mcp-loader');
    if (loader) {
        loader.remove();
    }
    
    // Use the common hide function
    if (typeof window.hideLoadingIndicator === 'function') {
        window.hideLoadingIndicator();
    } else {
        // Fallback
        const globalLoader = document.getElementById('loading-indicator');
        if (globalLoader) {
            globalLoader.style.display = 'none';
        }
    }
}

// Initialize when page loads
document.addEventListener('DOMContentLoaded', function() {
    console.log('MCP script loaded');
    
    // Remove any existing MCP cards from previous version
    const existingCards = document.querySelectorAll('#mcp-card');
    existingCards.forEach(card => card.remove());
    
    // Check MCP status to add button if enabled
    checkMCPStatus();
    
    // Listen for data load events - check MCP status to ensure button is shown/hidden correctly
    window.addEventListener('dataLoaded', function() {
        console.log('Data loaded event received by MCP');
        checkMCPStatus();
    });
    
    // Recheck MCP status when someone imports data
    window.addEventListener('dataImported', function() {
        console.log('Data imported event received by MCP');
        checkMCPStatus();
    });
});

// Export a function for showing the AI Assistant from other scripts
window.showMCPInterface = function() {
    showMCPModal();
};