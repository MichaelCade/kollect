package cost

import (
	"fmt"
)

func convertToSnapshotList(data interface{}) ([]map[string]string, bool) {
	if snapshots, ok := data.([]map[string]string); ok {
		return snapshots, true
	}

	if items, ok := data.([]interface{}); ok {
		result := make([]map[string]string, 0, len(items))
		for _, item := range items {
			if itemMap, ok := item.(map[string]interface{}); ok {
				snapshot := make(map[string]string)
				for k, v := range itemMap {
					if strVal, ok := v.(string); ok {
						snapshot[k] = strVal
					} else {
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

func stringOrDefault(ptr *string, defaultValue string) string {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}
