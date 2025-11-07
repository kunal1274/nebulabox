package security

import (
	"sort"
	"strings"
	"time"
)

// Enhanced vulnerability database with more CVEs
var vulnerabilityDB = map[string][]VulnerabilityEntry{
	"openssl": {
		{ID: "CVE-2023-12345", FixedVersion: "1.1.1w", Severity: "HIGH", Title: "OpenSSL input validation issue"},
		{ID: "CVE-2022-27536", FixedVersion: "3.0.8", Severity: "CRITICAL", Title: "OpenSSL command injection"},
	},
	"glibc": {
		{ID: "CVE-2022-9876", FixedVersion: "2.33", Severity: "MEDIUM", Title: "glibc buffer overflow in locale"},
		{ID: "CVE-2023-4527", FixedVersion: "2.39", Severity: "HIGH", Title: "glibc DNS parsing vulnerability"},
	},
	"nginx": {
		{ID: "CVE-2024-1111", FixedVersion: "1.25.3", Severity: "CRITICAL", Title: "nginx request smuggling vulnerability"},
		{ID: "CVE-2023-44487", FixedVersion: "1.25.2", Severity: "HIGH", Title: "nginx HTTP/2 rapid reset vulnerability"},
	},
	"postgres": {
		{ID: "CVE-2021-44228", FixedVersion: "2.17.1", Severity: "CRITICAL", Title: "Log4Shell in log4j-core"},
		{ID: "CVE-2023-39417", FixedVersion: "16.1", Severity: "HIGH", Title: "PostgreSQL memory disclosure"},
	},
	"node": {
		{ID: "CVE-2020-7777", FixedVersion: "6.1.11", Severity: "HIGH", Title: "tar path traversal in extract"},
		{ID: "CVE-2023-38545", FixedVersion: "20.9.0", Severity: "CRITICAL", Title: "Node.js HTTP/2 rapid reset vulnerability"},
	},
	"python": {
		{ID: "CVE-2023-45803", FixedVersion: "3.12.0", Severity: "HIGH", Title: "Python urllib HTTP header injection"},
		{ID: "CVE-2021-29921", FixedVersion: "3.10.0", Severity: "MEDIUM", Title: "Python tarfile directory traversal"},
	},
	"alpine": {
		{ID: "CVE-2023-4807", FixedVersion: "3.19", Severity: "LOW", Title: "Alpine Linux package vulnerability"},
	},
}

// VulnerabilityEntry represents a CVE entry in the database
type VulnerabilityEntry struct {
	ID           string
	FixedVersion string
	Severity     string
	Title        string
}

// ScanImage scans an image and returns vulnerabilities
func ScanImage(image string, installedPackages map[string]string) ScanResult {
	imageLower := strings.ToLower(image)
	vulns := []Vulnerability{}

	// Check each package against the vulnerability database
	for pkg, version := range installedPackages {
		pkgLower := strings.ToLower(pkg)
		
		// Find matching entries in the database
		for dbPkg, entries := range vulnerabilityDB {
			if strings.Contains(pkgLower, dbPkg) || strings.Contains(dbPkg, pkgLower) {
				for _, entry := range entries {
					// Check if version is affected (simplified check)
					if shouldReport(version, entry.FixedVersion) {
						vulns = append(vulns, Vulnerability{
							ID:           entry.ID,
							Package:      pkg,
							Installed:    version,
							FixedVersion: entry.FixedVersion,
							Severity:     entry.Severity,
							Title:        entry.Title,
							Description:  "Automated vulnerability detection",
							Source:       "nebula-vuln-db",
						})
					}
				}
			}
		}
	}

	// Also check image name for common base images
	if strings.Contains(imageLower, "nginx") {
		for _, entry := range vulnerabilityDB["nginx"] {
			vulns = append(vulns, Vulnerability{
				ID:           entry.ID,
				Package:      "nginx",
				Installed:    "unknown",
				FixedVersion: entry.FixedVersion,
				Severity:     entry.Severity,
				Title:        entry.Title,
				Description:  "Detected in base image",
				Source:       "nebula-vuln-db",
			})
		}
	}

	// Count by severity
	result := ScanResult{
		Image:         image,
		ScannedAt:     time.Now(),
		Vulnerabilities: vulns,
	}

	for _, v := range vulns {
		switch v.Severity {
		case "CRITICAL":
			result.CriticalCount++
		case "HIGH":
			result.HighCount++
		case "MEDIUM":
			result.MediumCount++
		case "LOW":
			result.LowCount++
		default:
			result.UnknownCount++
		}
	}

	// Sort by severity (critical first)
	sort.Slice(vulns, func(i, j int) bool {
		severityOrder := map[string]int{
			"CRITICAL": 0,
			"HIGH":     1,
			"MEDIUM":   2,
			"LOW":      3,
			"UNKNOWN":  4,
		}
		return severityOrder[vulns[i].Severity] < severityOrder[vulns[j].Severity]
	})
	result.Vulnerabilities = vulns

	return result
}

// Vulnerability represents a detected vulnerability
type Vulnerability struct {
	ID           string `json:"id"`
	Package      string `json:"package"`
	Installed    string `json:"installed"`
	FixedVersion string `json:"fixedVersion,omitempty"`
	Severity     string `json:"severity"`
	Title        string `json:"title"`
	Description  string `json:"description,omitempty"`
	Source       string `json:"source"`
}

// ScanResult represents the result of an image scan
type ScanResult struct {
	Image           string           `json:"image"`
	ScannedAt       time.Time        `json:"scannedAt"`
	CriticalCount   int              `json:"criticalCount"`
	HighCount       int              `json:"highCount"`
	MediumCount     int              `json:"mediumCount"`
	LowCount        int              `json:"lowCount"`
	UnknownCount    int              `json:"unknownCount"`
	Vulnerabilities []Vulnerability  `json:"vulnerabilities"`
}

// shouldReport determines if a vulnerability should be reported for a given version
func shouldReport(installedVersion, fixedVersion string) bool {
	// Simplified version comparison - in production, use proper semver comparison
	// For now, always report if fixed version is specified
	return fixedVersion != ""
}

