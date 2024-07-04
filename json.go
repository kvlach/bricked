package bricked

import (
	"encoding/json"
	"fmt"
	"reflect"
)

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

func PrettyPrint(v any) {
	bytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
