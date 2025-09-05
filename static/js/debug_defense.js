// Debug Detection and Defense System
(function() {
    'use strict';
    
    // Configuration
    var config = {
        redirectURL: 'https://www.yuanshen.com',
        enabled: true,
        debugMode: false
    };
    
    // Initialize defense system
    function initDefense(redirectURL) {
        if (redirectURL) {
            config.redirectURL = redirectURL;
        }
        
        if (!config.enabled) return;
        
        // Method 1: Keyboard shortcuts detection
        detectKeyboardShortcuts();
        
        // Method 2: Window size changes detection
        detectWindowSizeChanges();
        
        // Method 3: Console manipulation detection
        detectConsoleManipulation();
        
        // Method 4: DevTools timing detection
        detectDevToolsTiming();
        
        // Method 5: Element inspection detection
        detectElementInspection();
        
        // Method 6: Performance monitoring
        detectPerformanceAnomalies();
    }
    
    // Redirect function
    function triggerRedirect() {
        if (!config.debugMode) {
            config.debugMode = true;
            
            // Send detection notification to server
            fetch('/api/debug-detected', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    timestamp: new Date().toISOString(),
                    userAgent: navigator.userAgent,
                    url: window.location.href
                })
            }).catch(function() {
                // Ignore fetch errors, redirect anyway
            }).finally(function() {
                window.location.href = config.redirectURL;
            });
        }
    }
    
    // Detection Methods
    function detectKeyboardShortcuts() {
        document.addEventListener('keydown', function(e) {
            // F12
            if (e.keyCode === 123) {
                e.preventDefault();
                triggerRedirect();
                return false;
            }
            
            // Ctrl+Shift+I
            if (e.ctrlKey && e.shiftKey && e.keyCode === 73) {
                e.preventDefault();
                triggerRedirect();
                return false;
            }
            
            // Ctrl+Shift+C
            if (e.ctrlKey && e.shiftKey && e.keyCode === 67) {
                e.preventDefault();
                triggerRedirect();
                return false;
            }
            
            // Ctrl+U (View Source)
            if (e.ctrlKey && e.keyCode === 85) {
                e.preventDefault();
                triggerRedirect();
                return false;
            }
            
            // Ctrl+Shift+J (Console)
            if (e.ctrlKey && e.shiftKey && e.keyCode === 74) {
                e.preventDefault();
                triggerRedirect();
                return false;
            }
        });
        
        // Disable right-click context menu
        document.addEventListener('contextmenu', function(e) {
            e.preventDefault();
            return false;
        });
    }
    
    function detectWindowSizeChanges() {
        var threshold = 200;
        var initialOuterHeight = window.outerHeight;
        var initialOuterWidth = window.outerWidth;
        var initialInnerHeight = window.innerHeight;
        var initialInnerWidth = window.innerWidth;
        
        setInterval(function() {
            var heightDiff = Math.abs((window.outerHeight - window.innerHeight) - (initialOuterHeight - initialInnerHeight));
            var widthDiff = Math.abs((window.outerWidth - window.innerWidth) - (initialOuterWidth - initialInnerWidth));
            
            if (heightDiff > threshold || widthDiff > threshold) {
                triggerRedirect();
            }
        }, 500);
    }
    
    function detectConsoleManipulation() {
        if (typeof console !== 'undefined') {
            var originalLog = console.log;
            var originalInfo = console.info;
            var originalWarn = console.warn;
            var originalError = console.error;
            var originalClear = console.clear;
            
            console.log = function() {
                triggerRedirect();
                return originalLog.apply(console, arguments);
            };
            
            console.info = function() {
                triggerRedirect();
                return originalInfo.apply(console, arguments);
            };
            
            console.warn = function() {
                triggerRedirect();
                return originalWarn.apply(console, arguments);
            };
            
            console.error = function() {
                triggerRedirect();
                return originalError.apply(console, arguments);
            };
            
            console.clear = function() {
                triggerRedirect();
                return originalClear.apply(console, arguments);
            };
        }
    }
    
    function detectDevToolsTiming() {
        setInterval(function() {
            var start = performance.now();
            debugger;
            var end = performance.now();
            
            // If execution time is significantly longer, dev tools are likely open
            if (end - start > 100) {
                triggerRedirect();
            }
        }, 1000);
    }
    
    function detectElementInspection() {
        var trapElement = document.createElement('div');
        trapElement.id = '__debug_trap__';
        trapElement.style.display = 'none';
        document.body.appendChild(trapElement);
        
        Object.defineProperty(trapElement, 'id', {
            get: function() {
                triggerRedirect();
                return '__debug_trap__';
            },
            configurable: false
        });
        
        // Additional trap for toString() calls
        trapElement.toString = function() {
            triggerRedirect();
            return '';
        };
    }
    
    function detectPerformanceAnomalies() {
        var checks = 0;
        var maxChecks = 10;
        
        function performanceCheck() {
            checks++;
            var start = performance.now();
            
            setTimeout(function() {
                var end = performance.now();
                var duration = end - start;
                
                // If setTimeout took much longer than expected, tools might be open
                if (duration > 110 && checks < maxChecks) {
                    triggerRedirect();
                }
                
                if (checks < maxChecks) {
                    setTimeout(performanceCheck, 2000);
                }
            }, 100);
        }
        
        setTimeout(performanceCheck, 1000);
    }
    
    // Public API
    window.DebugDefense = {
        init: initDefense,
        enable: function() { config.enabled = true; },
        disable: function() { config.enabled = false; },
        setRedirectURL: function(url) { config.redirectURL = url; }
    };
    
    // Auto-initialize with default settings
    document.addEventListener('DOMContentLoaded', function() {
        initDefense();
    });
    
})();

