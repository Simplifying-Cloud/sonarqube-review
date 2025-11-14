# SonarQube Review

A Go application that generates developer-friendly reports from SonarQube findings that need to be fixed.

## Features

- Fetches issues from SonarQube API
- Generates reports in multiple formats:
  - **Markdown** - Easy to read and share
  - **HTML** - Beautiful, styled reports that can be opened in a browser
  - **JSON** - Machine-readable format for integration with other tools
  - **PDF** - Professional PDF reports with color-coded severity levels
- **Direct links to SonarQube** - Each issue includes a link to view it in SonarQube
- **Optional code snippets** - Include the problematic code with context
- Filter issues by severity, type, and status
- Automatic pagination for large issue sets
- Sorts issues by priority (type and severity)
- Clean, organized output with summary statistics

## Installation

### Prerequisites

- Go 1.21 or later
- Access to a SonarQube server
- SonarQube authentication token

### Build from source

```bash
# Clone the repository
git clone https://github.com/Simplifying-Cloud/sonarqube-review.git
cd sonarqube-review

# Build the application
go build -o sonarqube-review

# Optional: Install to GOPATH/bin
go install
```

## Usage

### Basic Usage

```bash
# Using environment variables
export SONAR_URL="https://sonarqube.example.com"
export SONAR_TOKEN="your-token-here"

./sonarqube-review -project my-project-key
```

### Using Command-Line Flags

```bash
./sonarqube-review \
  -url https://sonarqube.example.com \
  -token your-token-here \
  -project my-project-key \
  -format markdown \
  -output my-report
```

### Command-Line Options

| Flag | Description | Default | Required |
|------|-------------|---------|----------|
| `-url` | SonarQube server URL (or use `SONAR_URL` env var) | - | Yes |
| `-token` | SonarQube authentication token (or use `SONAR_TOKEN` env var) | - | Yes |
| `-project` | SonarQube project key | - | Yes |
| `-format` | Report format: `markdown`, `html`, `json`, or `pdf` | `markdown` | No |
| `-output` | Output file name (without extension) | `sonarqube-report` | No |
| `-severity` | Filter by severity: `BLOCKER`, `CRITICAL`, `MAJOR`, `MINOR`, `INFO` | (all) | No |
| `-type` | Filter by type: `BUG`, `VULNERABILITY`, `CODE_SMELL`, `SECURITY_HOTSPOT` | (all) | No |
| `-open-only` | Only include open issues | `true` | No |
| `-include-snippet` | Include code snippets in the report (requires additional API calls) | `false` | No |

### Examples

**Generate HTML report for critical vulnerabilities:**
```bash
./sonarqube-review \
  -project my-app \
  -format html \
  -severity CRITICAL \
  -type VULNERABILITY \
  -output critical-vulnerabilities
```

**Generate JSON report for all open bugs:**
```bash
./sonarqube-review \
  -project my-app \
  -format json \
  -type BUG \
  -output bugs-report
```

**Generate Markdown report for blocker and critical issues:**
```bash
./sonarqube-review \
  -project my-app \
  -severity BLOCKER,CRITICAL \
  -output high-priority-issues
```

**Generate report with code snippets and links:**
```bash
./sonarqube-review \
  -project my-app \
  -format html \
  -include-snippet \
  -output detailed-report
```

**Generate PDF report:**
```bash
./sonarqube-review \
  -project my-app \
  -format pdf \
  -include-snippet \
  -output security-audit
```

## Getting a SonarQube Token

1. Log in to your SonarQube server
2. Go to **My Account** > **Security**
3. Generate a new token
4. Copy the token and store it securely

For more information, see the [SonarQube documentation](https://docs.sonarqube.org/latest/user-guide/user-token/).

## Report Formats

### Markdown
- Easy to read in any text editor
- Perfect for including in pull requests or documentation
- Contains summary statistics and detailed issue information

### HTML
- Beautiful, styled reports with color-coded severity levels
- Can be opened directly in any web browser
- Includes visual statistics and organized sections

### JSON
- Machine-readable format
- Includes all raw data from SonarQube
- Perfect for integration with CI/CD pipelines or other tools

### PDF
- Professional, print-ready reports
- Color-coded severity levels (red, orange, yellow, blue, gray)
- Includes summary statistics and detailed issue breakdowns
- Perfect for sharing with stakeholders or archiving

## Project Structure

```
.
├── main.go                 # Application entry point and CLI
├── pkg/
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── sonarqube/         # SonarQube API client
│   │   ├── client.go      # API client implementation
│   │   └── types.go       # Data models
│   └── reporter/          # Report generators
│       ├── reporter.go    # Reporter interface
│       ├── markdown.go    # Markdown reporter
│       ├── html.go        # HTML reporter
│       └── json.go        # JSON reporter
├── go.mod
└── README.md
```

## Development

### Running Tests

```bash
go test ./...
```

### Adding a New Report Format

1. Create a new reporter in `pkg/reporter/`
2. Implement the `Reporter` interface
3. Add the format option in `main.go`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Support

For issues and questions, please open an issue on GitHub.
