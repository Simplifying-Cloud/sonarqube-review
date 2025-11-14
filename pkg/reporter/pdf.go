package reporter

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/Simplifying-Cloud/sonarqube-review/pkg/sonarqube"
)

// PDFReporter generates PDF reports
type PDFReporter struct{}

// NewPDFReporter creates a new PDF reporter
func NewPDFReporter() *PDFReporter {
	return &PDFReporter{}
}

// Generate creates a PDF report
func (r *PDFReporter) Generate(issues []sonarqube.Issue, outputFile, projectKey string) (string, error) {
	outputPath := outputFile + ".pdf"
	
	// Sort issues
	sortIssues(issues)
	
	// Group issues
	issuesByType := groupIssuesByType(issues)
	severityCounts := groupIssuesBySeverity(issues)
	
	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 15)
	
	// Add first page
	pdf.AddPage()
	
	// Title
	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(102, 126, 234)
	pdf.CellFormat(0, 15, "SonarQube Issues Report", "", 1, "C", false, 0, "")
	pdf.Ln(5)
	
	// Metadata
	pdf.SetFont("Arial", "", 11)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(0, 6, fmt.Sprintf("Project: %s", projectKey), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Total Issues: %d", len(issues)), "", 1, "L", false, 0, "")
	pdf.Ln(5)
	
	// Summary section
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(51, 51, 51)
	pdf.CellFormat(0, 10, "Summary", "", 1, "L", false, 0, "")
	pdf.Ln(2)
	
	// By Type
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 7, "By Type", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	
	for _, issueType := range []string{"VULNERABILITY", "SECURITY_HOTSPOT", "BUG", "CODE_SMELL"} {
		if count, ok := issuesByType[issueType]; ok {
			pdf.CellFormat(100, 6, "  "+formatType(issueType), "", 0, "L", false, 0, "")
			pdf.CellFormat(0, 6, fmt.Sprintf("%d", count), "", 1, "L", false, 0, "")
		}
	}
	pdf.Ln(3)
	
	// By Severity
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 7, "By Severity", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	
	for _, severity := range []string{"BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"} {
		if count, ok := severityCounts[severity]; ok {
			setSeverityColor(pdf, severity)
			pdf.CellFormat(100, 6, "  "+severity, "", 0, "L", false, 0, "")
			pdf.SetTextColor(51, 51, 51)
			pdf.CellFormat(0, 6, fmt.Sprintf("%d", count), "", 1, "L", false, 0, "")
		}
	}
	pdf.Ln(5)
	
	// Issues by type
	for _, issueType := range []string{"VULNERABILITY", "SECURITY_HOTSPOT", "BUG", "CODE_SMELL"} {
		typeIssues := filterIssuesByType(issues, issueType)
		if len(typeIssues) == 0 {
			continue
		}
		
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		pdf.SetTextColor(51, 51, 51)
		pdf.CellFormat(0, 10, fmt.Sprintf("%s (%d)", formatType(issueType), len(typeIssues)), "", 1, "L", false, 0, "")
		pdf.Ln(3)
		
		for i, issue := range typeIssues {
			// Check if we need a new page
			if pdf.GetY() > 250 {
				pdf.AddPage()
			}
			
			// Issue number and message
			pdf.SetFont("Arial", "B", 11)
			pdf.SetTextColor(51, 51, 51)
			pdf.MultiCell(0, 6, fmt.Sprintf("%d. %s", i+1, issue.Message), "", "L", false)
			pdf.Ln(1)
			
			// Metadata
			pdf.SetFont("Arial", "", 9)
			pdf.SetTextColor(80, 80, 80)
			
			// Severity
			setSeverityColor(pdf, issue.Severity)
			pdf.CellFormat(30, 5, "Severity:", "", 0, "L", false, 0, "")
			pdf.CellFormat(0, 5, issue.Severity, "", 1, "L", false, 0, "")
			
			// Rule
			pdf.SetTextColor(80, 80, 80)
			pdf.CellFormat(30, 5, "Rule:", "", 0, "L", false, 0, "")
			pdf.SetFont("Courier", "", 8)
			pdf.CellFormat(0, 5, issue.Rule, "", 1, "L", false, 0, "")
			
			// File
			pdf.SetFont("Arial", "", 9)
			pdf.CellFormat(30, 5, "File:", "", 0, "L", false, 0, "")
			pdf.SetFont("Courier", "", 7)
			fileText := issue.Component
			if issue.Line > 0 {
				fileText = fmt.Sprintf("%s (Line %d)", issue.Component, issue.Line)
			}
			pdf.MultiCell(0, 5, fileText, "", "L", false)
			
			// Status
			pdf.SetFont("Arial", "", 9)
			pdf.CellFormat(30, 5, "Status:", "", 0, "L", false, 0, "")
			pdf.CellFormat(0, 5, issue.Status, "", 1, "L", false, 0, "")
			
			// Effort
			if issue.Effort != "" {
				pdf.CellFormat(30, 5, "Effort:", "", 0, "L", false, 0, "")
				pdf.CellFormat(0, 5, issue.Effort, "", 1, "L", false, 0, "")
			}
			
			// SonarQube Link
			if issue.IssueURL != "" {
				pdf.Ln(2)
				pdf.SetFont("Arial", "U", 9)
				pdf.SetTextColor(102, 126, 234) // Blue color for link
				pdf.WriteLinkString(5, "View in SonarQube", issue.IssueURL)
				pdf.Ln(1)
				pdf.SetTextColor(80, 80, 80) // Reset color
			}
			
			// Code snippet
			if len(issue.CodeSnippet) > 0 {
				pdf.Ln(2)
				pdf.SetFont("Arial", "B", 9)
				pdf.SetTextColor(51, 51, 51)
				pdf.CellFormat(0, 5, "Code Snippet:", "", 1, "L", false, 0, "")
				pdf.Ln(1)
				
				// Background for code
				pdf.SetFillColor(245, 245, 245)
				pdf.SetTextColor(51, 51, 51)
				pdf.SetFont("Courier", "", 7)
				
				for _, line := range issue.CodeSnippet {
					// Check if we need a new page
					if pdf.GetY() > 270 {
						pdf.AddPage()
					}
					
					lineText := fmt.Sprintf("%4d | %s", line.Line, line.Code)
					// Truncate very long lines
					if len(lineText) > 120 {
						lineText = lineText[:117] + "..."
					}
					pdf.CellFormat(0, 4, lineText, "", 1, "L", true, 0, "")
				}
			}
			
			pdf.Ln(5)
		}
	}
	
	// Save PDF
	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return "", fmt.Errorf("saving PDF: %w", err)
	}
	
	return outputPath, nil
}

func setSeverityColor(pdf *gofpdf.Fpdf, severity string) {
	switch severity {
	case "BLOCKER":
		pdf.SetTextColor(211, 47, 47) // Red
	case "CRITICAL":
		pdf.SetTextColor(245, 124, 0) // Orange
	case "MAJOR":
		pdf.SetTextColor(251, 192, 45) // Yellow
	case "MINOR":
		pdf.SetTextColor(25, 118, 210) // Blue
	case "INFO":
		pdf.SetTextColor(117, 117, 117) // Gray
	default:
		pdf.SetTextColor(51, 51, 51) // Dark gray
	}
}
