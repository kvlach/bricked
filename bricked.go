package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type client struct {
	workspaceURL string
	apiVersion   string
}

func NewClient(workspaceURL, apiVersion string) client {
	return client{
		workspaceURL: workspaceURL,
		apiVersion:   apiVersion,
	}
}

// StrictUnmarshalJSON ensures that the struct and the JSON response are aligned
// by panicking if any one of the following is true:
//   - the JSON contains a key not present in the struct
//   - the struct contains a field not present in the JSON
//   - the struct contains a field without a `json` tag
func StrictUnmarshalJSON(data []byte, obj any) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	expectedFields := make(map[string]struct{})

	rv := reflect.ValueOf(obj).Elem()

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		key := field.Tag.Get("json")
		if key == "" {
			return fmt.Errorf("json tag not set for field '%s'", field.Name)
		}
		if _, ok := raw[key]; !ok {
			return fmt.Errorf("field '%s' is missing from JSON", key)
		}
		expectedFields[key] = struct{}{}
	}

	for key := range raw {
		if _, ok := expectedFields[key]; !ok {
			return fmt.Errorf("unexpected field '%s' found in JSON", key)
		}
	}

	for i := 0; i < rv.NumField(); i++ {
		vf := rv.Field(i)
		tf := rv.Type().Field(i)
		key := tf.Tag.Get("json")

		rawVal, _ := raw[key]

		if vf.Kind() == reflect.Struct {
			// Recursively unmarshal nested structs
			parsedVal := reflect.New(tf.Type).Interface()
			if err := StrictUnmarshalJSON(rawVal, parsedVal); err != nil {
				return err
			}
			vf.Set(reflect.ValueOf(parsedVal).Elem())
		} else {
			// Unmarshal non-struct fields
			if err := json.Unmarshal(rawVal, vf.Addr().Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}

func azCli(arguments ...string) {
	// e.g. azCli("databricks", "workspace", "list")
	panic("TODO: Implement")
}

func azLogin() {
	// Potential helpful error message:
	//   Use `az login` or `az login --tenant <tenant-id>` to login.
	//   Might also need to select the correct subscription using `az account set --subscription <name|id>`
	panic("TODO: Implement")
}

func getEntraToken() {
	// az account get-access-token \
	// 	--resource 2ff814a6-3304-4ab8-85cb-cd0e6f879c1d \
	// 	--query 'accessToken' \
	// 	--output tsv

	// azCli(
	// 	"account", "get-access-token",
	// 	"--resource", "2ff814a6-3304-4ab8-85cb-cd0e6f879c1d",
	// 	"--query", "accessToken",
	// 	"--output", "tsv",
	// )

	panic("TODO: Implement")
}

func transientError(statusCode int) bool {
	switch statusCode {
	// taken from the --retry flag in the 8.8.0 curl manual
	case 408, 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

func (c *client) request(what string, ret any) {
	// GET/POST/etc. will just wrap this method
	panic("TODO: Implement")
}

func (c *client) GET(what string, ret any) {
	url := fmt.Sprintf("%s/%s/%s", c.workspaceURL, c.version, what)
	// url := "https://httpstat.us/random/200,500"

	var err error
	var resp *http.Response
	var body []byte

	for i := 0; i < 5; i++ {
		resp, err = http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		// Read body first for better error messages
		if transientError(resp.StatusCode) {
			fmt.Print("Retrying...")
			continue
		}

		if err := StrictUnmarshalJSON(body, ret); err != nil {
			panic(err)
		}
	}

	if resp.StatusCode >= 400 {
		panic(fmt.Errorf("Request failed, status code: %d, body: %s", resp.StatusCode, body))
	}
	fmt.Printf("Satus code: %d, body: %s\n", resp.StatusCode, body)
}
