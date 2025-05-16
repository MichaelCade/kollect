package cost

import (
	"fmt"
)

// convertToSnapshotList converts various snapshot data representations to a list of map[string]string
func convertToSnapshotList(data interface{}) ([]map[string]string, bool) {
	// Check if already in the right format
	if snapshots, ok := data.([]map[string]string); ok {
		return snapshots, true
	}

	// Check if it's a slice of interfaces (common when unmarshaling JSON)
	if items, ok := data.([]interface{}); ok {
		result := make([]map[string]string, 0, len(items))
		for _, item := range items {
			if itemMap, ok := item.(map[string]interface{}); ok {
				snapshot := make(map[string]string)
				for k, v := range itemMap {
					if strVal, ok := v.(string); ok {
						snapshot[k] = strVal
					} else {
						// Handle non-string values by converting them to strings
						snapshot[k] = fmt.Sprintf("%v", v)
					}
				}
				result = append(result, snapshot)
			}
		}
		return result, len(result) > 0
	}

	return nil, false
}

// stringOrDefault returns the string value or a default if pointer is nil
func stringOrDefault(ptr *string, defaultValue string) string {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}
