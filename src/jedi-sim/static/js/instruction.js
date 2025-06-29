// Instruction Page JavaScript
document.addEventListener('DOMContentLoaded', function() {
    initInstructionPage();
});

function initInstructionPage() {
    // Add smooth scrolling behavior
    document.documentElement.style.scrollBehavior = 'smooth';
    
    // Add table row hover effects
    addTableHoverEffects();
    
    // Add back button functionality
    setupBackButton();
    
    // Add page load animations
    addPageAnimations();
    
    // Add responsive behavior
    setupResponsiveBehavior();
    
    // Add keyboard shortcuts
    addKeyboardShortcuts();
}

function addTableHoverEffects() {
    const tableRows = document.querySelectorAll('.field-table tbody tr');
    
    tableRows.forEach((row, index) => {
        // Add staggered animation delay
        row.style.animationDelay = `${(index + 1) * 0.1}s`;
        
        // Add hover effect
        row.addEventListener('mouseenter', function() {
            this.style.transform = 'scale(1.01)';
            this.style.boxShadow = '0 4px 20px rgba(79, 140, 255, 0.15)';
        });
        
        row.addEventListener('mouseleave', function() {
            this.style.transform = 'scale(1)';
            this.style.boxShadow = 'none';
        });
        
        // Add click effect for better UX
        row.addEventListener('click', function() {
            this.style.transform = 'scale(0.98)';
            setTimeout(() => {
                this.style.transform = 'scale(1)';
            }, 150);
        });
    });
}

function setupBackButton() {
    const backBtn = document.querySelector('.back-btn');
    
    if (backBtn) {
        // Add ripple effect
        backBtn.addEventListener('click', function(e) {
            createRippleEffect(e, this);
        });
        
        // Add keyboard navigation
        backBtn.addEventListener('keydown', function(e) {
            if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                window.location.href = 'index.html';
            }
        });
        
        // Add focus styles
        backBtn.addEventListener('focus', function() {
            this.style.outline = '2px solid #4f8cff';
            this.style.outlineOffset = '2px';
        });
        
        backBtn.addEventListener('blur', function() {
            this.style.outline = 'none';
        });
    }
}

function createRippleEffect(event, element) {
    const ripple = document.createElement('span');
    const rect = element.getBoundingClientRect();
    const size = Math.max(rect.width, rect.height);
    const x = event.clientX - rect.left - size / 2;
    const y = event.clientY - rect.top - size / 2;
    
    ripple.style.width = ripple.style.height = size + 'px';
    ripple.style.left = x + 'px';
    ripple.style.top = y + 'px';
    ripple.classList.add('ripple');
    
    element.appendChild(ripple);
    
    setTimeout(() => {
        ripple.remove();
    }, 600);
}

function addPageAnimations() {
    // Add CSS for ripple effect
    const style = document.createElement('style');
    style.textContent = `
        .ripple {
            position: absolute;
            border-radius: 50%;
            background: rgba(255, 255, 255, 0.3);
            transform: scale(0);
            animation: ripple-animation 0.6s linear;
            pointer-events: none;
        }
        
        @keyframes ripple-animation {
            to {
                transform: scale(4);
                opacity: 0;
            }
        }
        
        .field-table tr {
            transition: all 0.3s ease;
        }
        
        .back-btn {
            position: relative;
            overflow: hidden;
        }
    `;
    document.head.appendChild(style);
    
    // Add entrance animation for the container
    const container = document.querySelector('.container');
    if (container) {
        container.style.opacity = '0';
        container.style.transform = 'translateY(30px)';
        
        setTimeout(() => {
            container.style.transition = 'all 0.8s ease';
            container.style.opacity = '1';
            container.style.transform = 'translateY(0)';
        }, 100);
    }
}

function setupResponsiveBehavior() {
    // Handle window resize
    window.addEventListener('resize', debounce(function() {
        adjustTableForScreenSize();
    }, 250));
    
    // Initial adjustment
    adjustTableForScreenSize();
}

function adjustTableForScreenSize() {
    const table = document.querySelector('.field-table');
    const isMobile = window.innerWidth <= 768;
    
    if (table && isMobile) {
        // Add horizontal scroll for mobile
        table.style.minWidth = '600px';
        
        // Wrap table in scrollable container
        if (!table.parentElement.classList.contains('table-scroll-container')) {
            const wrapper = document.createElement('div');
            wrapper.className = 'table-scroll-container';
            wrapper.style.overflowX = 'auto';
            wrapper.style.borderRadius = '12px';
            wrapper.style.boxShadow = '0 4px 20px rgba(0, 0, 0, 0.08)';
            
            table.parentNode.insertBefore(wrapper, table);
            wrapper.appendChild(table);
        }
    } else if (table) {
        // Reset for desktop
        table.style.minWidth = 'auto';
        const wrapper = table.parentElement;
        if (wrapper && wrapper.classList.contains('table-scroll-container')) {
            wrapper.parentNode.insertBefore(table, wrapper);
            wrapper.remove();
        }
    }
}

function addKeyboardShortcuts() {
    document.addEventListener('keydown', function(e) {
        // Escape key to go back
        if (e.key === 'Escape') {
            window.location.href = 'index.html';
        }
        
        // Ctrl/Cmd + Home to scroll to top
        if ((e.ctrlKey || e.metaKey) && e.key === 'Home') {
            e.preventDefault();
            scrollToTop();
        }
    });
}

// Utility functions
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

function scrollToTop() {
    window.scrollTo({
        top: 0,
        behavior: 'smooth'
    });
}

// Performance optimization
function optimizePerformance() {
    // Use Intersection Observer for better performance
    if ('IntersectionObserver' in window) {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.style.opacity = '1';
                    entry.target.style.transform = 'translateY(0)';
                }
            });
        });
        
        // Observe table rows for better performance
        document.querySelectorAll('.field-table tr').forEach(row => {
            observer.observe(row);
        });
    }
}

// Initialize performance optimizations
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', optimizePerformance);
} else {
    optimizePerformance();
} 