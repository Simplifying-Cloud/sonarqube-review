package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Simplifying-Cloud/sonarqube-review/pkg/config"
	"github.com/Simplifying-Cloud/sonarqube-review/pkg/reporter"
	"github.com/Simplifying-Cloud/sonarqube-review/pkg/sonarqube"
)

func main() {
	var (
		sonarURL      = flag.String("url", "", "SonarQube server URL (or set SONAR_URL env var)")
		sonarToken    = flag.String("token", "", "SonarQube authentication token (or set SONAR_TOKEN env var)")
		projectKey    = flag.String("project", "", "SonarQube project key")
		outputFile    = flag.String("output", "sonarqube-report", "Output file name (without extension)")
		format      = flag.String("format", "markdown", "Report format: markdown, html, json, pdf")
		severity      = flag.String("severity", "", "Filter by severity (BLOCKER, CRITICAL, MAJOR, MINOR, INFO)")
		issueType     = flag.String("type", "", "Filter by type (BUG, VULNERABILITY, CODE_SMELL, SECURITY_HOTSPOT)")
		onlyOpen      = flag.Bool("open-only", true, "Only include open issues")
		includeSnippet = flag.Bool("include-snippet", false, "Include code snippets (requires additional API calls)")
	)

	flag.Parse()

	// Load configuration
	cfg := config.New()
	
	// Override with flags if provided
	if *sonarURL != "" {
		cfg.SonarURL = *sonarURL
	}
	if *sonarToken != "" {
		cfg.SonarToken = *sonarToken
	}

	// Validate required parameters
	if cfg.SonarURL == "" {
		fmt.Fprintln(os.Stderr, "Error: SonarQube URL is required (use -url flag or SONAR_URL env var)")
		flag.Usage()
		os.Exit(1)
	}
	if cfg.SonarToken == "" {
		fmt.Fprintln(os.Stderr, "Error: SonarQube token is required (use -token flag or SONAR_TOKEN env var)")
		flag.Usage()
		os.Exit(1)
	}
	if *projectKey == "" {
		fmt.Fprintln(os.Stderr, "Error: Project key is required (use -project flag)")
		flag.Usage()
		os.Exit(1)
	}

	// Create SonarQube client
	client := sonarqube.NewClient(cfg.SonarURL, cfg.SonarToken)

	// Build filter options
	filters := sonarqube.IssueFilters{
		ProjectKey: *projectKey,
		Severity:   *severity,
		Type:       *issueType,
		OnlyOpen:   *onlyOpen,
	}

	// Fetch issues
	fmt.Printf("Fetching issues for project: %s...\n", *projectKey)
	issues, err := client.GetIssues(filters)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching issues: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d issues\n", len(issues))

	// Optionally fetch code snippets
	if *includeSnippet {
		fmt.Println("Fetching code snippets...")
		for i := range issues {
			if issues[i].TextRange != nil {
				snippet, err := client.GetSourceCode(
					issues[i].Component,
					issues[i].TextRange.StartLine,
					issues[i].TextRange.EndLine,
				)
				if err == nil && snippet != nil {
					issues[i].CodeSnippet = snippet
				}
			} else if issues[i].Line > 0 {
				snippet, err := client.GetSourceCode(
					issues[i].Component,
					issues[i].Line,
					issues[i].Line,
				)
				if err == nil && snippet != nil {
					issues[i].CodeSnippet = snippet
				}
			}
		}
	}

	// Generate report
	var rep reporter.Reporter
	switch *format {
	case "markdown", "md":
		rep = reporter.NewMarkdownReporter()
	case "html":
		rep = reporter.NewHTMLReporter()
	case "json":
		rep = reporter.NewJSONReporter()
	case "pdf":
		rep = reporter.NewPDFReporter()
	default:
		fmt.Fprintf(os.Stderr, "Error: Invalid format '%s'. Use: markdown, html, json, or pdf\n", *format)
		os.Exit(1)
	}

	outputPath, err := rep.Generate(issues, *outputFile, *projectKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Report generated successfully: %s\n", outputPath)
}
