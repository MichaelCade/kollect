// common.js

function showLoadingIndicator() {
    document.getElementById('loading-indicator').style.display = 'flex';
}

function hideLoadingIndicator() {
    document.getElementById('loading-indicator').style.display = 'none';
}

document.getElementById('export-button').addEventListener('click', () => {
    fetch('/api/data')
        .then(response => response.json())
        .then(data => {
            const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'data.json';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        })
        .catch(error => console.error('Error exporting data:', error));
});

document.getElementById('import-button').addEventListener('click', () => {
    document.getElementById('import-file').click();
});

document.getElementById('import-file').addEventListener('change', (event) => {
    const file = event.target.files[0];
    if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
            const data = JSON.parse(e.target.result);
            fetch('/api/import', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            })
            .then(response => response.json())
            .then(data => {
                // Refresh the page to display the imported data
                location.reload();
            })
            .catch(error => console.error('Error importing data:', error));
        };
        reader.readAsText(file);
    }
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
            // Add your row templates and createTable calls here
        } catch (error) {
            console.error("Error processing data:", error);
        }
    }
});