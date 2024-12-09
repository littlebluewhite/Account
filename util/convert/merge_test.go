package convert

import (
	"encoding/json"
	"reflect"
	"testing"
)

func jsonEqual(a, b json.RawMessage) bool {
	var objA interface{}
	var objB interface{}

	if err := json.Unmarshal(a, &objA); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &objB); err != nil {
		return false
	}

	return reflect.DeepEqual(objA, objB)
}

func TestMergeJSON(t *testing.T) {
	input1 := []byte(`{"a": 1, "b": "stringB", "c": {"d": 12, "e": "stringE"}}`)
	input2 := []byte(`{"b": "newB", "c": {"e": "newE", "f": "newF"}, "h": 0}`)

	expectedOutput := map[string]interface{}{
		"b": "stringB",
		"c": map[string]interface{}{
			"e": "stringE",
			"f": "newF",
		},
		"h": float64(0), // Changed to float64
	}

	mergedJSON, err := MergeJSON(input1, input2)
	if err != nil {
		t.Fatalf("MergeJSON failed: %v", err)
	}

	var mergedMap map[string]interface{}
	if err := json.Unmarshal(mergedJSON, &mergedMap); err != nil {
		t.Fatalf("Failed to unmarshal merged JSON: %v", err)
	}

	if !reflect.DeepEqual(mergedMap, expectedOutput) {
		t.Errorf("Merged JSON does not match expected.\nExpected: %v\nGot: %v", expectedOutput, mergedMap)
	}
}

func TestMergeJSONExtended(t *testing.T) {
	testCases := []struct {
		name           string
		input1         []byte
		input2         []byte
		expectedOutput map[string]interface{}
	}{
		{
			name:   "Basic Merge",
			input1: []byte(`{"a": 1, "b": "stringB", "c": {"d": 12, "e": "stringE"}}`),
			input2: []byte(`{"b": "newB", "c": {"e": "newE", "f": "newF"}, "h": 0}`),
			expectedOutput: map[string]interface{}{
				"b": "stringB",
				"c": map[string]interface{}{
					"e": "stringE",
					"f": "newF",
				},
				"h": float64(0), // Changed to float64
			},
		},
		{
			name:   "Nested Merge",
			input1: []byte(`{"x": {"y": {"z": 100}}, "a": "alpha"}`),
			input2: []byte(`{"x": {"y": {"z": 200, "w": 300}}, "b": "beta"}`),
			expectedOutput: map[string]interface{}{
				"x": map[string]interface{}{
					"y": map[string]interface{}{
						"z": float64(100), // Changed to float64
						"w": float64(300), // Changed to float64
					},
				},
				"b": "beta",
			},
		},
		{
			name:   "Different Data Types",
			input1: []byte(`{"a": [1, 2, 3], "b": true, "c": "stringC1"}`),
			input2: []byte(`{"a": [4, 5], "b": false, "d": "stringD"}`),
			expectedOutput: map[string]interface{}{
				"a": []interface{}{float64(1), float64(2), float64(3)},
				"b": true,
				"d": "stringD",
			},
		},
		{
			name:   "Empty Input1",
			input1: []byte(`{}`),
			input2: []byte(`{"a": 1, "b": 2}`),
			expectedOutput: map[string]interface{}{
				"a": float64(1), // Changed to float64
				"b": float64(2), // Changed to float64
			},
		},
		{
			name:           "Empty Input2",
			input1:         []byte(`{"a": 10, "b": 20}`),
			input2:         []byte(`{}`),
			expectedOutput: map[string]interface{}{},
		},
		{
			name:   "No Overlap",
			input1: []byte(`{"a": 1}`),
			input2: []byte(`{"b": 2}`),
			expectedOutput: map[string]interface{}{
				"b": float64(2), // Changed to float64
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mergedJSON, err := MergeJSON(tc.input1, tc.input2)
			if err != nil {
				t.Fatalf("MergeJSON failed: %v", err)
			}

			var mergedMap map[string]interface{}
			if err := json.Unmarshal(mergedJSON, &mergedMap); err != nil {
				t.Fatalf("Failed to unmarshal merged JSON: %v", err)
			}

			if !reflect.DeepEqual(mergedMap, tc.expectedOutput) {
				t.Errorf("Merged JSON does not match expected.\nExpected: %v\nGot: %v", tc.expectedOutput, mergedMap)
			}
		})
	}
}
