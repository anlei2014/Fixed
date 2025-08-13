// This file contains the JavaScript code that handles the functionality of the buttons.

document.addEventListener('DOMContentLoaded', function () {
    // Error code report form handling
    const form = document.getElementById('reportForm');
    const resultDiv = document.getElementById('errorResult');

    if (form) {
        form.addEventListener('submit', function (e) {
            e.preventDefault();
            const errorcode = document.getElementById('errorcode').value;

            fetch('/api/report', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                },
                body: 'errorcode=' + encodeURIComponent(errorcode)
            })
            .then(response => response.json())
            .then(data => {
                if (data.status) {
                    resultDiv.textContent = 'Succeeded: ' + data.message;
                    resultDiv.style.color = 'green';
                } else {
                    resultDiv.textContent = 'Failed: ' + (data.message || 'Unknown error');
                    resultDiv.style.color = 'red';
                }
            })
            .catch(() => {
                resultDiv.textContent = 'Request failed';
                resultDiv.style.color = 'red';
            });
        });
    }
});