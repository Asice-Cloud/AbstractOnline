package server

import (
	"net/http"
)

// DefenseConfig holds configuration for debug defense
type DefenseConfig struct {
	RedirectURL string
	Enabled     bool
}

// DefaultDefenseConfig returns default configuration
func DefaultDefenseConfig() *DefenseConfig {
	return &DefenseConfig{
		RedirectURL: "https://www.yuanshen.com",
		Enabled:     true,
	}
}

// HandleDebugDetection handles requests when F12 is detected
func HandleDebugDetection(w http.ResponseWriter, r *http.Request, config *DefenseConfig) {
	if !config.Enabled {
		return
	}

	// Log the detection attempt (optional)
	// logger.Info("Debug tools detected from IP: " + r.RemoteAddr)

	// Redirect to specified URL
	http.Redirect(w, r, config.RedirectURL, http.StatusTemporaryRedirect)
}

// GenerateDebugDetectionScript generates JavaScript code for F12 detection
func GenerateDebugDetectionScript(redirectURL string) string {
	return `
<script>
(function() {
    'use strict';
    
    var redirectURL = '` + redirectURL + `';
    var isDebugMode = false;
    
    // Method 1: Detect F12 key press and other debug shortcuts
    var keyDetection = function() {
        document.addEventListener('keydown', function(e) {
            // F12 key
            if (e.keyCode === 123) {
                e.preventDefault();
                if (!isDebugMode) {
                    isDebugMode = true;
                    window.location.href = redirectURL;
                }
                return false;
            }
            
            // Ctrl+Shift+I (Chrome dev tools)
            if (e.ctrlKey && e.shiftKey && e.keyCode === 73) {
                e.preventDefault();
                if (!isDebugMode) {
                    isDebugMode = true;
                    window.location.href = redirectURL;
                }
                return false;
            }
            
            // Ctrl+Shift+C (Chrome inspect element)
            if (e.ctrlKey && e.shiftKey && e.keyCode === 67) {
                e.preventDefault();
                if (!isDebugMode) {
                    isDebugMode = true;
                    window.location.href = redirectURL;
                }
                return false;
            }
            
            // Ctrl+U (View source)
            if (e.ctrlKey && e.keyCode === 85) {
                e.preventDefault();
                if (!isDebugMode) {
                    isDebugMode = true;
                    window.location.href = redirectURL;
                }
                return false;
            }
        });
        
        // Disable right-click context menu
        document.addEventListener('contextmenu', function(e) {
            e.preventDefault();
            return false;
        });
    };
    
    // Method 2: Improved window size detection with better logic
    var windowSizeDetection = function() {
        var threshold = 300; // Increased threshold to reduce false positives
        var initialOuterHeight = window.outerHeight;
        var initialInnerHeight = window.innerHeight;
        var detectionCount = 0;
        var maxDetections = 3; // Require multiple detections before triggering
        
        setInterval(function() {
            var currentOuterHeight = window.outerHeight;
            var currentInnerHeight = window.innerHeight;
            
            // Calculate the difference in the available viewport
            var initialViewport = initialInnerHeight;
            var currentViewport = currentInnerHeight;
            var viewportDiff = Math.abs(initialViewport - currentViewport);
            
            // Also check if the height difference between outer and inner changed significantly
            var initialHeightDiff = initialOuterHeight - initialInnerHeight;
            var currentHeightDiff = currentOuterHeight - currentInnerHeight;
            var heightDiffChange = Math.abs(currentHeightDiff - initialHeightDiff);
            
            // Only trigger if both conditions are met and viewport actually shrunk
            if (heightDiffChange > threshold && currentViewport < initialViewport) {
                detectionCount++;
                if (detectionCount >= maxDetections && !isDebugMode) {
                    isDebugMode = true;
                    window.location.href = redirectURL;
                }
            } else {
                // Reset detection count if conditions aren't met
                detectionCount = Math.max(0, detectionCount - 1);
            }
        }, 1000); // Increased interval to reduce CPU usage
    };
    
    // Method 3: Detect debugger statements (but with timeout protection)
    var debuggerDetection = function() {
        var detectionAttempts = 0;
        var maxAttempts = 5;
        
        var debugCheck = function() {
            if (detectionAttempts >= maxAttempts) return; // Stop after max attempts
            
            detectionAttempts++;
            var start = performance.now();
            
            try {
                debugger;
            } catch(e) {
                // Ignore any errors
            }
            
            var end = performance.now();
            if (end - start > 100) { // If debugger took too long, dev tools are open
                if (!isDebugMode) {
                    isDebugMode = true;
                    window.location.href = redirectURL;
                }
            }
        };
        
        // Only run debugger detection periodically and with limits
        var interval = setInterval(function() {
            debugCheck();
            if (detectionAttempts >= maxAttempts) {
                clearInterval(interval);
            }
        }, 5000); // Less frequent checks
    };
    
    // Initialize detection methods
    keyDetection();
    windowSizeDetection();
    debuggerDetection();
    
})();
</script>`
}
