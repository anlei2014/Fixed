document.addEventListener('DOMContentLoaded', function () {
    // 保持原有 main.js 功能
    const clearBtn = document.getElementById('errorClearButton');
    if (clearBtn) {
        clearBtn.onclick = function() {
            document.getElementById('errorcode').value = '';
            document.getElementById('errorResult').textContent = '';
        };
    }
});
