package convert

import (
	"encoding/json"
	"fmt"
)

// MergeJSON merges two JSON objects based on the keys of input2.
// For each key in input2, if input1 has a value for that key, use input1's value.
// If the value is a nested object, perform the merge recursively.
// Returns the merged JSON as a byte slice.
func MergeJSON(input1, input2 []byte) ([]byte, error) {
	var map1 map[string]interface{}
	var map2 map[string]interface{}

	// Unmarshal input1
	if err := json.Unmarshal(input1, &map1); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input1: %w", err)
	}

	// Unmarshal input2
	if err := json.Unmarshal(input2, &map2); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input2: %w", err)
	}

	// Merge map1 into map2 based on keys in map2
	mergedMap := mergeMaps(map1, map2)

	// Marshal the merged map back to JSON
	mergedJSON, err := json.Marshal(mergedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merged JSON: %w", err)
	}

	return mergedJSON, nil
}

// mergeMaps merges map1 into map2 based on the keys in map2.
// It returns a new map containing the merged data.
func mergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	for key, val2 := range map2 {
		if val1, exists := map1[key]; exists {
			// If both values are maps, merge recursively
			mapVal1, ok1 := val1.(map[string]interface{})
			mapVal2, ok2 := val2.(map[string]interface{})
			if ok1 && ok2 {
				merged[key] = mergeMaps(mapVal1, mapVal2)
				continue
			}
			// Else, use val1 from map1
			merged[key] = val1
			continue
		}
		// If key not in map1, use val2 from map2
		merged[key] = val2
	}

	return merged
}
