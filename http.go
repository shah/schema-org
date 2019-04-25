package schemamd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPUserAgent may be passed into getHTTPResult as the default HTTP User-Agent header parameter
const HTTPUserAgent = "github.com/lectio/score"

// HTTPTimeout may be passed into getHTTPResult function as the default HTTP timeout parameter
const HTTPTimeout = time.Second * 90

// HTTPResult encapsulates an API call
type httpResult struct {
	apiEndpoint string
	body        *[]byte
}

// GetHTTPResult runs the apiEndpoint and returns the body of the HTTP result
// TODO: Consider using [HTTP Cache](https://github.com/gregjones/httpcache)
func getHTTPResult(apiEndpoint string, userAgent string, timeout time.Duration) (*httpResult, Issue) {
	result := new(httpResult)
	result.apiEndpoint = apiEndpoint

	httpClient := http.Client{
		Timeout: timeout,
	}
	req, reqErr := http.NewRequest(http.MethodGet, apiEndpoint, nil)
	if reqErr != nil {
		return nil, newIssue(apiEndpoint, UnableToCreateHTTPRequest, fmt.Sprintf("Unable to create HTTP request: %v", reqErr), true)
	}
	req.Header.Set("User-Agent", userAgent)
	resp, getErr := httpClient.Do(req)
	if getErr != nil {
		return nil, newIssue(apiEndpoint, UnableToExecuteHTTPGETRequest, fmt.Sprintf("Unable to execute HTTP GET request: %v", getErr), true)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, newHTTPResponseIssue(apiEndpoint, resp.StatusCode, fmt.Sprintf("HTTP response status is not 200: %v", resp.StatusCode), true)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, newIssue(apiEndpoint, UnableToReadBodyFromHTTPResponse, fmt.Sprintf("Unable to read body from HTTP response: %v", readErr), true)
	}

	result.body = &body
	return result, nil
}
