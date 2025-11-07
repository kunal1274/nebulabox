package api

import (
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/nebulabox/nebulabox/internal/security"
)

// ImageScanRequest represents a scan request
type ImageScanRequest struct {
    Image string `json:"image" binding:"required"`
}

// Vulnerability represents a single finding
type Vulnerability struct {
    ID           string `json:"id"`
    Package      string `json:"package"`
    Installed    string `json:"installed"`
    FixedVersion string `json:"fixedVersion,omitempty"`
    Severity     string `json:"severity"` // CRITICAL|HIGH|MEDIUM|LOW|UNKNOWN
    Title        string `json:"title"`
    Description  string `json:"description,omitempty"`
    Source       string `json:"source"` // database/source
}

// ImageScanResult is the scan response
type ImageScanResult struct {
    Image           string           `json:"image"`
    ScannedAt       time.Time        `json:"scannedAt"`
    CriticalCount   int              `json:"criticalCount"`
    HighCount       int              `json:"highCount"`
    MediumCount     int              `json:"mediumCount"`
    LowCount        int              `json:"lowCount"`
    UnknownCount    int              `json:"unknownCount"`
    Vulnerabilities []Vulnerability  `json:"vulnerabilities"`
}

// postScanImage handles POST /api/images/scan
// NOTE: Stub/mock implementation that fabricates findings for demo purposes
func (s *Server) postScanImage(c *gin.Context) {
    var req ImageScanRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body", "details": err.Error()})
        return
    }

    image := strings.TrimSpace(req.Image)
    if image == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
        return
    }

    // Use enhanced scanner
    // Simulate installed packages based on image name
    installedPackages := make(map[string]string)
    base := strings.ToLower(image)
    
    // Add common packages based on image type
    installedPackages["openssl"] = "1.1.1f"
    installedPackages["glibc"] = "2.31"
    
    if strings.Contains(base, "nginx") {
        installedPackages["nginx"] = "1.21.0"
    }
    if strings.Contains(base, "postgres") || strings.Contains(base, "postgresql") {
        installedPackages["postgresql"] = "13.0"
    }
    if strings.Contains(base, "node") {
        installedPackages["nodejs"] = "18.0.0"
        installedPackages["npm"] = "8.0.0"
    }
    if strings.Contains(base, "python") {
        installedPackages["python"] = "3.9.0"
    }
    if strings.Contains(base, "alpine") {
        installedPackages["alpine"] = "3.18"
    }

    // Perform scan
    scanResult := security.ScanImage(image, installedPackages)
    
    // Convert to API format
    res := ImageScanResult{
        Image:           scanResult.Image,
        ScannedAt:       scanResult.ScannedAt,
        CriticalCount:   scanResult.CriticalCount,
        HighCount:       scanResult.HighCount,
        MediumCount:     scanResult.MediumCount,
        LowCount:        scanResult.LowCount,
        UnknownCount:    scanResult.UnknownCount,
        Vulnerabilities: make([]Vulnerability, len(scanResult.Vulnerabilities)),
    }
    
    for i, v := range scanResult.Vulnerabilities {
        res.Vulnerabilities[i] = Vulnerability{
            ID:           v.ID,
            Package:      v.Package,
            Installed:    v.Installed,
            FixedVersion: v.FixedVersion,
            Severity:     v.Severity,
            Title:        v.Title,
            Description:  v.Description,
            Source:       v.Source,
        }
    }

    c.JSON(http.StatusOK, res)
}


