package reporter

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Simplifying-Cloud/sonarqube-review/pkg/sonarqube"
)

// MarkdownReporter generates Markdown reports
type MarkdownReporter struct{}

// NewMarkdownReporter creates a new Markdown reporter
func NewMarkdownReporter() *MarkdownReporter {
	return &MarkdownReporter{}
}

// Generate creates a Markdown report
func (r *MarkdownReporter) Generate(issues []sonarqube.Issue, outputFile, projectKey string) (string, error) {
	outputPath := outputFile + ".md"
	
	// Sort issues by severity and type
	sortIssues(issues)
	
	// Group issues by type
	issuesByType := groupIssuesByType(issues)
	
	// Create file
	f, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()
	
	// Write header
	fmt.Fprintf(f, "# SonarQube Issues Report\n\n")
	fmt.Fprintf(f, "**Project:** %s  \n", projectKey)
	fmt.Fprintf(f, "**Generated:** %s  \n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f, "**Total Issues:** %d\n\n", len(issues))
	
	// Write summary
	fmt.Fprintf(f, "## Summary\n\n")
	fmt.Fprintf(f, "| Type | Count |\n")
	fmt.Fprintf(f, "|------|-------|\n")
	
	for _, issueType := range []string{"VULNERABILITY", "SECURITY_HOTSPOT", "BUG", "CODE_SMELL"} {
		if count, ok := issuesByType[issueType]; ok {
			fmt.Fprintf(f, "| %s | %d |\n", formatType(issueType), count)
		}
	}
	fmt.Fprintf(f, "\n")
	
	// Write severity breakdown
	severityCounts := groupIssuesBySeverity(issues)
	fmt.Fprintf(f, "### By Severity\n\n")
	fmt.Fprintf(f, "| Severity | Count |\n")
	fmt.Fprintf(f, "|----------|-------|\n")
	
	for _, severity := range []string{"BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"} {
		if count, ok := severityCounts[severity]; ok {
			fmt.Fprintf(f, "| %s | %d |\n", formatSeverity(severity), count)
		}
	}
	fmt.Fprintf(f, "\n")
	
	// Write detailed issues grouped by type
	for _, issueType := range []string{"VULNERABILITY", "SECURITY_HOTSPOT", "BUG", "CODE_SMELL"} {
		typeIssues := filterIssuesByType(issues, issueType)
		if len(typeIssues) == 0 {
			continue
		}
		
		fmt.Fprintf(f, "## %s (%d)\n\n", formatType(issueType), len(typeIssues))
		
		for i, issue := range typeIssues {
			fmt.Fprintf(f, "### %d. %s\n\n", i+1, issue.Message)
			fmt.Fprintf(f, "- **Severity:** %s\n", formatSeverity(issue.Severity))
			fmt.Fprintf(f, "- **Rule:** `%s`\n", issue.Rule)
			fmt.Fprintf(f, "- **File:** `%s`", issue.Component)
			if issue.Line > 0 {
				fmt.Fprintf(f, " (Line %d)", issue.Line)
			}
			fmt.Fprintf(f, "\n")
			if issue.IssueURL != "" {
				fmt.Fprintf(f, "- **View in SonarQube:** [Open Issue](%s)\n", issue.IssueURL)
			}
			fmt.Fprintf(f, "- **Status:** %s\n", issue.Status)
			if issue.Effort != "" {
				fmt.Fprintf(f, "- **Effort:** %s\n", issue.Effort)
			}
			if len(issue.Tags) > 0 {
				fmt.Fprintf(f, "- **Tags:** %s\n", strings.Join(issue.Tags, ", "))
			}
			
			// Add code snippet if available
			if len(issue.CodeSnippet) > 0 {
				fmt.Fprintf(f, "\n**Code Snippet:**\n\n```\n")
				for _, line := range issue.CodeSnippet {
					fmt.Fprintf(f, "%4d | %s\n", line.Line, line.Code)
				}
				fmt.Fprintf(f, "```\n")
			}
			
			fmt.Fprintf(f, "\n")
		}
	}
	
	return outputPath, nil
}

func formatType(issueType string) string {
	switch issueType {
	case "VULNERABILITY":
		return "Vulnerability"
	case "SECURITY_HOTSPOT":
		return "Security Hotspot"
	case "BUG":
		return "Bug"
	case "CODE_SMELL":
		return "Code Smell"
	default:
		return issueType
	}
}

func formatSeverity(severity string) string {
	emoji := map[string]string{
		"BLOCKER":  "🔴",
		"CRITICAL": "🟠",
		"MAJOR":    "🟡",
		"MINOR":    "🔵",
		"INFO":     "⚪",
	}
	
	if e, ok := emoji[severity]; ok {
		return fmt.Sprintf("%s %s", e, severity)
	}
	return severity
}

func sortIssues(issues []sonarqube.Issue) {
	sort.Slice(issues, func(i, j int) bool {
		// First by type priority
		if sonarqube.GetTypePriority(issues[i].Type) != sonarqube.GetTypePriority(issues[j].Type) {
			return sonarqube.GetTypePriority(issues[i].Type) > sonarqube.GetTypePriority(issues[j].Type)
		}
		// Then by severity
		if sonarqube.GetSeverityLevel(issues[i].Severity) != sonarqube.GetSeverityLevel(issues[j].Severity) {
			return sonarqube.GetSeverityLevel(issues[i].Severity) > sonarqube.GetSeverityLevel(issues[j].Severity)
		}
		// Then by component
		return issues[i].Component < issues[j].Component
	})
}

func groupIssuesByType(issues []sonarqube.Issue) map[string]int {
	counts := make(map[string]int)
	for _, issue := range issues {
		counts[issue.Type]++
	}
	return counts
}

func groupIssuesBySeverity(issues []sonarqube.Issue) map[string]int {
	counts := make(map[string]int)
	for _, issue := range issues {
		counts[issue.Severity]++
	}
	return counts
}

func filterIssuesByType(issues []sonarqube.Issue, issueType string) []sonarqube.Issue {
	var filtered []sonarqube.Issue
	for _, issue := range issues {
		if issue.Type == issueType {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}
