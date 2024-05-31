window.onload = function() {
    fetch('/api/data')
        .then(response => response.json())
        .then(data => {
            document.getElementById('data').textContent = JSON.stringify(data, null, 2);
        })
        .catch(error => console.error('Error:', error));
};