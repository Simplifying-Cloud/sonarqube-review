package reporter

import "github.com/Simplifying-Cloud/sonarqube-review/pkg/sonarqube"

// Reporter defines the interface for generating reports
type Reporter interface {
	Generate(issues []sonarqube.Issue, outputFile, projectKey string) (string, error)
}
