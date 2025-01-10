// chatbot.js

document.addEventListener('DOMContentLoaded', () => {
    const chatbotInput = document.getElementById('chatbot-input');
    const chatbotOutput = document.getElementById('chatbot-output');

    chatbotInput.addEventListener('keypress', (event) => {
        if (event.key === 'Enter') {
            const query = chatbotInput.value;
            fetch('/api/chatbot', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ query })
            })
            .then(response => response.json())
            .then(data => {
                chatbotOutput.innerHTML += `<div class="chatbot-response">${data.response}</div>`;
                chatbotInput.value = '';
            })
            .catch(error => console.error('Error:', error));
        }
    });
});