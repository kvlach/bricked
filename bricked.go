package bricked

import (
	"fmt"
	"io"
	"net/http"
)

type client struct {
	workspaceURL string
	apiVersion   string
}

func NewClient(workspaceName, apiVersion string) client {
	workspaceURL := findWorkspaceUrl(workspaceName)

	return client{
		workspaceURL: workspaceURL,
		apiVersion:   apiVersion,
	}
}

func (c *client) assertVersion(expected string) {
	if c.apiVersion != expected {
		panic(fmt.Errorf("Expected API version '%s', got '%s'", c.apiVersion, expected))
	}
}

func (c *client) request(what string, ret any) {
	// GET/POST/etc. will just wrap this method
	panic("TODO: Implement")
}

func (c *client) GET(what string, ret any) {
	url := fmt.Sprintf("%s/%s/%s", c.workspaceURL, c.apiVersion, what)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode >= 400 {
		panic(fmt.Errorf("Request failed, status code: %d, body: %s", resp.StatusCode, body))
	}
	fmt.Printf("Satus code: %d, body: %s\n", resp.StatusCode, body)

	if err := StrictUnmarshalJSON(body, ret); err != nil {
		panic(err)
	}
}

// func transientError(statusCode int) bool {
// 	switch statusCode {
// 	// taken from the --retry flag in the 8.8.0 curl manual
// 	case 408, 429, 500, 502, 503, 504:
// 		return true
// 	default:
// 		return false
// 	}
// }

// func (c *client) GET(what string, ret any) {
// 	url := fmt.Sprintf("%s/%s/%s", c.workspaceURL, c.apiVersion, what)
// 	// url := "https://httpstat.us/random/200,500"

// 	var err error
// 	var resp *http.Response
// 	var body []byte

// 	for i := 0; i < 5; i++ {
// 		resp, err = http.Get(url)
// 		if err != nil {
// 			panic(err)
// 		}
// 		defer resp.Body.Close()

// 		body, err = io.ReadAll(resp.Body)
// 		if err != nil {
// 			panic(err)
// 		}

// 		// Read body first for better error messages
// 		if transientError(resp.StatusCode) {
// 			fmt.Print("Retrying...")
// 			continue
// 		}

// 		if err := StrictUnmarshalJSON(body, ret); err != nil {
// 			panic(err)
// 		}
// 	}

// 	if resp.StatusCode >= 400 {
// 		panic(fmt.Errorf("Request failed, status code: %d, body: %s", resp.StatusCode, body))
// 	}
// 	fmt.Printf("Satus code: %d, body: %s\n", resp.StatusCode, body)
// }
