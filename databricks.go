package bricked

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

type databricks struct {
	workspaceURL string
	apiVersion   string
	token        string
}

func (c *databricks) assertVersion(expected string) {
	if c.apiVersion != expected {
		panic(fmt.Errorf("Expected API version '%s', got '%s'", c.apiVersion, expected))
	}
}

func (c *databricks) request(what string, ret any) {
	// GET/POST/etc. will just wrap this method
	panic("TODO: Implement")
}

func (c *databricks) GET(what string, ret any) {
	url := fmt.Sprintf("%s/api/%s/%s", c.workspaceURL, c.apiVersion, what)
	slog.Debug("GET", "url", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)

	// resp, err := http.Get(url)
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
	slog.Debug("request", "status code", resp.StatusCode, "body", body)

	if err := StrictUnmarshalJSON(body, ret); err != nil {
		panic(err)
	}
	PrettyPrint(ret)
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
