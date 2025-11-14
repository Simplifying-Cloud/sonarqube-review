package sonarqube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Client represents a SonarQube API client
type Client struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewClient creates a new SonarQube client
func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		client:  &http.Client{},
	}
}

// IssueFilters contains filter options for fetching issues
type IssueFilters struct {
	ProjectKey string
	Severity   string
	Type       string
	OnlyOpen   bool
}

// GetIssues fetches issues from SonarQube
func (c *Client) GetIssues(filters IssueFilters) ([]Issue, error) {
	params := url.Values{}
	params.Add("componentKeys", filters.ProjectKey)
	
	if filters.Severity != "" {
		params.Add("severities", filters.Severity)
	}
	
	if filters.Type != "" {
		params.Add("types", filters.Type)
	}
	
	if filters.OnlyOpen {
		params.Add("resolved", "false")
	}
	
	params.Add("ps", "500") // page size
	
	var allIssues []Issue
	page := 1
	
	for {
		params.Set("p", fmt.Sprintf("%d", page))
		
		apiURL := fmt.Sprintf("%s/api/issues/search?%s", c.baseURL, params.Encode())
		
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}
		
		req.Header.Set("Authorization", "Bearer "+c.token)
		
		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("making request: %w", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
		}
		
		var result IssuesResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("decoding response: %w", err)
		}
		
		allIssues = append(allIssues, result.Issues...)
		
		// Check if we need to fetch more pages
		if len(allIssues) >= result.Total {
			break
		}
		page++
	}
	
	// Enrich issues with URLs
	for i := range allIssues {
		allIssues[i].IssueURL = c.GetIssueURL(allIssues[i].Key, allIssues[i].Project)
		allIssues[i].RuleURL = c.GetRuleURL(allIssues[i].Rule)
	}
	
	return allIssues, nil
}

// GetIssueURL returns the URL to view an issue in SonarQube
func (c *Client) GetIssueURL(issueKey, projectKey string) string {
	// For SonarCloud
	if strings.Contains(c.baseURL, "sonarcloud.io") {
		return fmt.Sprintf("%s/project/issues?id=%s&open=%s", 
			c.baseURL, url.QueryEscape(projectKey), url.QueryEscape(issueKey))
	}
	// For SonarQube server
	return fmt.Sprintf("%s/project/issues?id=%s&open=%s", 
		c.baseURL, url.QueryEscape(projectKey), url.QueryEscape(issueKey))
}

// GetRuleURL returns the URL to view rule details in SonarQube
func (c *Client) GetRuleURL(ruleKey string) string {
	// For SonarCloud - simpler format that works
	if strings.Contains(c.baseURL, "sonarcloud.io") {
		return fmt.Sprintf("%s/coding_rules?open=%s", 
			c.baseURL, url.QueryEscape(ruleKey))
	}
	// For SonarQube server
	return fmt.Sprintf("%s/coding_rules?open=%s", 
		c.baseURL, url.QueryEscape(ruleKey))
}

// GetSourceCode fetches the source code snippet for an issue
func (c *Client) GetSourceCode(componentKey string, startLine, endLine int) ([]CodeLine, error) {
	if startLine <= 0 {
		startLine = 1
	}
	if endLine <= 0 {
		endLine = startLine + 5 // Show 5 lines if no range specified
	}
	
	// Add some context lines
	contextStart := startLine - 2
	if contextStart < 1 {
		contextStart = 1
	}
	contextEnd := endLine + 2
	
	params := url.Values{}
	params.Add("key", componentKey)
	params.Add("from", fmt.Sprintf("%d", contextStart))
	params.Add("to", fmt.Sprintf("%d", contextEnd))
	
	apiURL := fmt.Sprintf("%s/api/sources/lines?%s", c.baseURL, params.Encode())
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+c.token)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		// Source code might not be available, return empty
		return nil, nil
	}
	
	var result struct {
		Sources []CodeLine `json:"sources"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	
	// Strip HTML from code
	for i := range result.Sources {
		result.Sources[i].Code = stripHTML(result.Sources[i].Code)
	}
	
	return result.Sources, nil
}

// stripHTML removes HTML tags from a string
func stripHTML(s string) string {
	// Remove all HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	result := re.ReplaceAllString(s, "")
	
	// Decode common HTML entities
	result = strings.ReplaceAll(result, "&lt;", "<")
	result = strings.ReplaceAll(result, "&gt;", ">")
	result = strings.ReplaceAll(result, "&amp;", "&")
	result = strings.ReplaceAll(result, "&quot;", "\"")
	result = strings.ReplaceAll(result, "&#39;", "'")
	result = strings.ReplaceAll(result, "&nbsp;", " ")
	
	return result
}

// GetProjectInfo fetches project information
func (c *Client) GetProjectInfo(projectKey string) (*Project, error) {
	apiURL := fmt.Sprintf("%s/api/components/show?component=%s", c.baseURL, url.QueryEscape(projectKey))
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+c.token)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var result struct {
		Component Project `json:"component"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	
	return &result.Component, nil
}
