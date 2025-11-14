package sonarqube

import (
	"encoding/json"
	"time"
)

// SonarTime is a custom time type that handles SonarQube's timestamp format
type SonarTime struct {
	time.Time
}

// UnmarshalJSON handles parsing SonarQube timestamps
func (st *SonarTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	
	if s == "" {
		return nil
	}
	
	// Try multiple formats that SonarQube might use
	formats := []string{
		"2006-01-02T15:04:05-0700",  // SonarQube format: 2025-11-08T02:43:49+0000
		time.RFC3339,                 // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05Z",      // ISO 8601 without timezone
	}
	
	var t time.Time
	var err error
	
	for _, format := range formats {
		t, err = time.Parse(format, s)
		if err == nil {
			st.Time = t
			return nil
		}
	}
	
	return err
}

// MarshalJSON handles encoding SonarTime to JSON
func (st SonarTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(st.Time)
}

// IssuesResponse represents the response from the issues API
type IssuesResponse struct {
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
}

// Issue represents a SonarQube issue/finding
type Issue struct {
	Key          string     `json:"key"`
	Rule         string     `json:"rule"`
	Severity     string     `json:"severity"`
	Component    string     `json:"component"`
	Project      string     `json:"project"`
	Line         int        `json:"line"`
	Message      string     `json:"message"`
	Type         string     `json:"type"`
	Status       string     `json:"status"`
	CreationDate SonarTime  `json:"creationDate"`
	UpdateDate   SonarTime  `json:"updateDate"`
	Author       string     `json:"author"`
	Effort       string     `json:"effort"`
	Tags         []string   `json:"tags"`
	TextRange    *TextRange `json:"textRange,omitempty"`
	
	// Enriched fields (not from API)
	IssueURL     string     `json:"issueUrl,omitempty"`
	RuleURL      string     `json:"ruleUrl,omitempty"`
	CodeSnippet  []CodeLine `json:"codeSnippet,omitempty"`
}

// TextRange represents the location of an issue in the source code
type TextRange struct {
	StartLine   int `json:"startLine"`
	EndLine     int `json:"endLine"`
	StartOffset int `json:"startOffset"`
	EndOffset   int `json:"endOffset"`
}

// CodeLine represents a line of source code
type CodeLine struct {
	Line int    `json:"line"`
	Code string `json:"code"`
}

// Project represents a SonarQube project
type Project struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Qualifier   string `json:"qualifier"`
}

// GetSeverityLevel returns a numeric level for severity (higher = more severe)
func GetSeverityLevel(severity string) int {
	switch severity {
	case "BLOCKER":
		return 5
	case "CRITICAL":
		return 4
	case "MAJOR":
		return 3
	case "MINOR":
		return 2
	case "INFO":
		return 1
	default:
		return 0
	}
}

// GetTypePriority returns a numeric priority for issue types
func GetTypePriority(issueType string) int {
	switch issueType {
	case "VULNERABILITY":
		return 4
	case "SECURITY_HOTSPOT":
		return 3
	case "BUG":
		return 2
	case "CODE_SMELL":
		return 1
	default:
		return 0
	}
}
