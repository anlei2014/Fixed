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
    
    // Add error code form handling
    const addErrorCodeForm = document.getElementById('addErrorCodeForm');
    const addErrorCodeResultDiv = document.getElementById('result');
    
    if (addErrorCodeForm) {
        addErrorCodeForm.addEventListener('submit', function (e) {
            e.preventDefault();
            
            // Collect all input values
            const errorCodeData = {
                z0: parseInt(document.getElementById('z0').value) || 0,
                z1: parseInt(document.getElementById('z1').value) || 0,
                z2: parseInt(document.getElementById('z2').value) || 0,
                z3_phase: parseInt(document.getElementById('z3_phase').value) || 0,
                z3_class: parseInt(document.getElementById('z3_class').value) || 0,
                z4z5_errorcode: parseInt(document.getElementById('z4z5_errorcode').value) || 0,
                z6z7_errordata: parseInt(document.getElementById('z6z7_errordata').value) || 0,
                description: document.getElementById('description').value || ''
            };
            
            fetch('/api/add-error-code', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(errorCodeData)
            })
            .then(response => response.json())
            .then(data => {
                if (data.status) {
                    addErrorCodeResultDiv.textContent = 'Succeeded: ' + data.message;
                    addErrorCodeResultDiv.style.color = 'green';
                } else {
                    addErrorCodeResultDiv.textContent = 'Failed: ' + (data.message || 'Unknown error');
                    addErrorCodeResultDiv.style.color = 'red';
                }
            })
            .catch(error => {
                console.error('Error:', error);
                addErrorCodeResultDiv.textContent = 'Request failed';
                addErrorCodeResultDiv.style.color = 'red';
            });
        });
    }
});