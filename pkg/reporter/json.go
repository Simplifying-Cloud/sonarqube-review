package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Simplifying-Cloud/sonarqube-review/pkg/sonarqube"
)

// JSONReporter generates JSON reports
type JSONReporter struct{}

// NewJSONReporter creates a new JSON reporter
func NewJSONReporter() *JSONReporter {
	return &JSONReporter{}
}

// JSONReport represents the structure of the JSON report
type JSONReport struct {
	Project      string                 `json:"project"`
	GeneratedAt  string                 `json:"generatedAt"`
	TotalIssues  int                    `json:"totalIssues"`
	Summary      Summary                `json:"summary"`
	Issues       []sonarqube.Issue      `json:"issues"`
}

// Summary contains aggregated statistics
type Summary struct {
	ByType     map[string]int `json:"byType"`
	BySeverity map[string]int `json:"bySeverity"`
}

// Generate creates a JSON report
func (r *JSONReporter) Generate(issues []sonarqube.Issue, outputFile, projectKey string) (string, error) {
	outputPath := outputFile + ".json"
	
	// Sort issues
	sortIssues(issues)
	
	// Create report
	report := JSONReport{
		Project:     projectKey,
		GeneratedAt: time.Now().Format(time.RFC3339),
		TotalIssues: len(issues),
		Summary: Summary{
			ByType:     groupIssuesByType(issues),
			BySeverity: groupIssuesBySeverity(issues),
		},
		Issues: issues,
	}
	
	// Create file
	f, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()
	
	// Write JSON
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return "", fmt.Errorf("encoding JSON: %w", err)
	}
	
	return outputPath, nil
}
