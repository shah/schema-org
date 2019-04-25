package schemamd

import "fmt"

const (
	UnableToCreateHTTPRequest        string = "SCHEMA_ORG_E-0100"
	UnableToExecuteHTTPGETRequest    string = "SCHEMA_ORG_E-0200"
	InvalidAPIRespHTTPStatusCode     string = "SCHEMA_ORG_E-0300"
	UnableToReadBodyFromHTTPResponse string = "SCHEMA_ORG_E-0400"

	APIErrorResponseFound       string = "SCHEMA_ORG_E-0500"
	NoAPIKeyProvidedInCodeOrEnv string = "SCHEMA_ORG_E-0600"

	UnknownGraphNodeType           string = "SCHEMA_ORG_E-1000"
	GraphNodeIDNotFound            string = "SCHEMA_ORG_E-1001"
	GraphNodeTypeNotFound          string = "SCHEMA_ORG_E-1002"
	GraphNodeIDNotHandled          string = "SCHEMA_ORG_E-1003"
	GraphPropertyNodeKeyNotHandled string = "SCHEMA_ORG_E-1004"
	GraphClassNodeKeyNotHandled    string = "SCHEMA_ORG_E-1005"
)

// Issue is a structured problem identification with context information
type Issue interface {
	IssueContext() interface{} // this will be the schema org object plus location (item index, etc.), it's kept generic so it doesn't require package dependency
	IssueCode() string         // useful to uniquely identify a particular code
	Issue() string             // the

	IsError() bool   // this issue is an error
	IsWarning() bool // this issue is a warning
}

// Issues packages multiple issues into a container
type Issues interface {
	ErrorsAndWarnings() []Issue
	IssueCounts() (uint, uint, uint)
	HandleIssues(errorHandler func(Issue), warningHandler func(Issue))
}

type issue struct {
	APIEndpoint    string `json:"context"`
	Code           string `json:"code"`
	Message        string `json:"message"`
	IsIssueAnError bool   `json:"isError"`
}

func newIssue(apiEndpoint string, code string, message string, isError bool) Issue {
	result := new(issue)
	result.APIEndpoint = apiEndpoint
	result.Code = code
	result.Message = message
	result.IsIssueAnError = isError
	return result
}

func newHTTPResponseIssue(apiEndpoint string, httpRespStatusCode int, message string, isError bool) Issue {
	result := new(issue)
	result.APIEndpoint = apiEndpoint
	result.Code = fmt.Sprintf("%s-HTTP-%d", InvalidAPIRespHTTPStatusCode, httpRespStatusCode)
	result.Message = message
	result.IsIssueAnError = isError
	return result
}

func (i issue) IssueContext() interface{} {
	return i.APIEndpoint
}

func (i issue) IssueCode() string {
	return i.Code
}

func (i issue) Issue() string {
	return i.Message
}

func (i issue) IsError() bool {
	return i.IsIssueAnError
}

func (i issue) IsWarning() bool {
	return !i.IsIssueAnError
}

// Error satisfies the Go error contract
func (i issue) Error() string {
	return i.Message
}
