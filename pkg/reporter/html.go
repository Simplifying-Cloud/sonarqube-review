package reporter

import (
	"fmt"
	"html"
	"os"
	"time"

	"github.com/Simplifying-Cloud/sonarqube-review/pkg/sonarqube"
)

// HTMLReporter generates HTML reports
type HTMLReporter struct{}

// NewHTMLReporter creates a new HTML reporter
func NewHTMLReporter() *HTMLReporter {
	return &HTMLReporter{}
}

// Generate creates an HTML report
func (r *HTMLReporter) Generate(issues []sonarqube.Issue, outputFile, projectKey string) (string, error) {
	outputPath := outputFile + ".html"
	
	// Sort issues
	sortIssues(issues)
	
	// Group issues
	issuesByType := groupIssuesByType(issues)
	severityCounts := groupIssuesBySeverity(issues)
	
	// Create file
	f, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()
	
	// Write HTML
	fmt.Fprint(f, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SonarQube Issues Report</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            border-radius: 10px;
            margin-bottom: 30px;
        }
        h1 {
            margin: 0 0 10px 0;
        }
        .meta {
            opacity: 0.9;
            font-size: 14px;
        }
        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .card {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .card h3 {
            margin: 0 0 15px 0;
            color: #333;
            font-size: 16px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .stat {
            display: flex;
            justify-content: space-between;
            padding: 8px 0;
            border-bottom: 1px solid #eee;
        }
        .stat:last-child {
            border-bottom: none;
        }
        .issue-section {
            background: white;
            padding: 25px;
            border-radius: 8px;
            margin-bottom: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .issue-section h2 {
            margin-top: 0;
            color: #333;
            border-bottom: 3px solid #667eea;
            padding-bottom: 10px;
        }
        .issue {
            border-left: 4px solid #ddd;
            padding: 15px;
            margin: 15px 0;
            background: #fafafa;
            border-radius: 4px;
        }
        .issue.BLOCKER { border-left-color: #d32f2f; }
        .issue.CRITICAL { border-left-color: #f57c00; }
        .issue.MAJOR { border-left-color: #fbc02d; }
        .issue.MINOR { border-left-color: #1976d2; }
        .issue.INFO { border-left-color: #757575; }
        .issue-title {
            font-weight: 600;
            font-size: 16px;
            color: #333;
            margin-bottom: 10px;
        }
        .issue-meta {
            display: flex;
            flex-wrap: wrap;
            gap: 15px;
            font-size: 13px;
            color: #666;
        }
        .badge {
            display: inline-block;
            padding: 3px 8px;
            border-radius: 3px;
            font-size: 12px;
            font-weight: 600;
            text-transform: uppercase;
        }
        .badge.BLOCKER { background: #ffebee; color: #d32f2f; }
        .badge.CRITICAL { background: #fff3e0; color: #f57c00; }
        .badge.MAJOR { background: #fffde7; color: #f57f17; }
        .badge.MINOR { background: #e3f2fd; color: #1976d2; }
        .badge.INFO { background: #f5f5f5; color: #757575; }
        .badge.VULNERABILITY { background: #fce4ec; color: #c2185b; }
        .badge.BUG { background: #ffebee; color: #d32f2f; }
        .badge.CODE_SMELL { background: #e8f5e9; color: #388e3c; }
        .badge.SECURITY_HOTSPOT { background: #fff3e0; color: #f57c00; }
        code {
            background: #f5f5f5;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
            font-size: 12px;
        }
        .links {
            margin-top: 10px;
            display: flex;
            gap: 10px;
        }
        .links a {
            color: #667eea;
            text-decoration: none;
            font-size: 13px;
        }
        .links a:hover {
            text-decoration: underline;
        }
        .code-snippet {
            margin-top: 15px;
            background: #1e1e1e;
            color: #d4d4d4;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
            font-family: 'Courier New', monospace;
            font-size: 13px;
            line-height: 1.5;
        }
        .code-snippet-title {
            font-weight: 600;
            margin-bottom: 10px;
            color: #333;
        }
        .code-line {
            white-space: pre;
        }
        .line-number {
            color: #858585;
            margin-right: 15px;
            user-select: none;
        }
    </style>
</head>
<body>
`)
	
	// Header
	fmt.Fprintf(f, `    <div class="header">
        <h1>SonarQube Issues Report</h1>
        <div class="meta">
            <strong>Project:</strong> %s<br>
            <strong>Generated:</strong> %s<br>
            <strong>Total Issues:</strong> %d
        </div>
    </div>
`, html.EscapeString(projectKey), time.Now().Format("2006-01-02 15:04:05"), len(issues))
	
	// Summary cards
	fmt.Fprint(f, `    <div class="summary">
        <div class="card">
            <h3>By Type</h3>
`)
	for _, issueType := range []string{"VULNERABILITY", "SECURITY_HOTSPOT", "BUG", "CODE_SMELL"} {
		if count, ok := issuesByType[issueType]; ok {
			fmt.Fprintf(f, `            <div class="stat"><span>%s</span><strong>%d</strong></div>
`, formatType(issueType), count)
		}
	}
	fmt.Fprint(f, `        </div>
        <div class="card">
            <h3>By Severity</h3>
`)
	for _, severity := range []string{"BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"} {
		if count, ok := severityCounts[severity]; ok {
			fmt.Fprintf(f, `            <div class="stat"><span>%s</span><strong>%d</strong></div>
`, severity, count)
		}
	}
	fmt.Fprint(f, `        </div>
    </div>
`)
	
	// Issues by type
	for _, issueType := range []string{"VULNERABILITY", "SECURITY_HOTSPOT", "BUG", "CODE_SMELL"} {
		typeIssues := filterIssuesByType(issues, issueType)
		if len(typeIssues) == 0 {
			continue
		}
		
		fmt.Fprintf(f, `    <div class="issue-section">
        <h2>%s (%d)</h2>
`, formatType(issueType), len(typeIssues))
		
		for _, issue := range typeIssues {
			fmt.Fprintf(f, `        <div class="issue %s">
            <div class="issue-title">%s</div>
            <div class="issue-meta">
                <span><span class="badge %s">%s</span></span>
                <span><span class="badge %s">%s</span></span>
                <span><strong>Rule:</strong> <code>%s</code></span>
                <span><strong>File:</strong> <code>%s</code>`,
				issue.Severity,
				html.EscapeString(issue.Message),
				issue.Severity, issue.Severity,
				issue.Type, formatType(issue.Type),
				html.EscapeString(issue.Rule),
				html.EscapeString(issue.Component))
			
			if issue.Line > 0 {
				fmt.Fprintf(f, ` (Line %d)`, issue.Line)
			}
			fmt.Fprint(f, `</span>
`)
			
			if issue.Effort != "" {
				fmt.Fprintf(f, `                <span><strong>Effort:</strong> %s</span>
`, html.EscapeString(issue.Effort))
			}
			
			fmt.Fprint(f, `            </div>
`)
			
			// Add links
			if issue.IssueURL != "" {
				fmt.Fprint(f, `            <div class="links">
`)
				fmt.Fprintf(f, `                <a href="%s" target="_blank">View in SonarQube →</a>
`, html.EscapeString(issue.IssueURL))
				fmt.Fprint(f, `            </div>
`)
			}
			
			// Add code snippet
			if len(issue.CodeSnippet) > 0 {
				fmt.Fprint(f, `            <div class="code-snippet-title">Code Snippet:</div>
            <div class="code-snippet">
`)
				for _, line := range issue.CodeSnippet {
					fmt.Fprintf(f, `                <div class="code-line"><span class="line-number">%4d</span>%s</div>
`, line.Line, html.EscapeString(line.Code))
				}
				fmt.Fprint(f, `            </div>
`)
			}
			
			fmt.Fprint(f, `        </div>
`)
		}
		
		fmt.Fprint(f, `    </div>
`)
	}
	
	fmt.Fprint(f, `</body>
</html>
`)
	
	return outputPath, nil
}
