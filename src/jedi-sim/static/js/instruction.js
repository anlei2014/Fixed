// Instruction page JavaScript functionality

document.addEventListener('DOMContentLoaded', function() {
    // Initialize the instruction page
    initializeInstructionPage();
});

function initializeInstructionPage() {
    // Add click event listeners to table rows for better UX
    addTableRowInteractions();
    
    // Add copy functionality for CAN format
    addCopyFunctionality();
}

function addTableRowInteractions() {
    const tableRows = document.querySelectorAll('.field-table tbody tr');
    
    tableRows.forEach(row => {
        row.addEventListener('click', function() {
            // Remove highlight from other rows
            tableRows.forEach(r => r.classList.remove('highlighted'));
            // Add highlight to clicked row
            this.classList.add('highlighted');
        });
        
        // Add hover effect
        row.addEventListener('mouseenter', function() {
            this.style.backgroundColor = '#f8f9fa';
        });
        
        row.addEventListener('mouseleave', function() {
            if (!this.classList.contains('highlighted')) {
                this.style.backgroundColor = '';
            }
        });
    });
}

function addCopyFunctionality() {
    const canFormat = document.querySelector('.can-format');
    if (canFormat) {
        canFormat.style.cursor = 'pointer';
        canFormat.title = 'Click to copy CAN format';
        
        canFormat.addEventListener('click', function() {
            const text = this.textContent;
            navigator.clipboard.writeText(text).then(() => {
                showCopyNotification();
            }).catch(err => {
                console.error('Failed to copy text: ', err);
                // Fallback for older browsers
                copyTextFallback(text);
            });
        });
    }
}

function copyTextFallback(text) {
    const textArea = document.createElement('textarea');
    textArea.value = text;
    document.body.appendChild(textArea);
    textArea.select();
    document.execCommand('copy');
    document.body.removeChild(textArea);
    showCopyNotification();
}

function showCopyNotification() {
    // Create notification element
    const notification = document.createElement('div');
    notification.textContent = 'CAN format copied to clipboard!';
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        background: #4caf50;
        color: white;
        padding: 10px 20px;
        border-radius: 4px;
        z-index: 1000;
        font-family: Arial, sans-serif;
        box-shadow: 0 2px 10px rgba(0,0,0,0.2);
    `;
    
    document.body.appendChild(notification);
    
    // Remove notification after 3 seconds
    setTimeout(() => {
        if (notification.parentNode) {
            notification.parentNode.removeChild(notification);
        }
    }, 3000);
}

// Add CSS for highlighted rows
const style = document.createElement('style');
style.textContent = `
    .field-table tbody tr.highlighted {
        background-color: #e3f2fd !important;
        border-left: 4px solid #4f8cff;
    }
    
    .field-table tbody tr {
        transition: background-color 0.2s ease;
    }
`;
document.head.appendChild(style); 