package tool

func DuplicateSliceMap(items []map[string]interface{}) []map[string]interface{} {
	newItems := make([]map[string]interface{}, len(items))
	for k, v := range items {
		newItem := make(map[string]interface{}, len(v))
		for k1, v1 := range v {
			newItem[k1] = v1
		}

		newItems[k] = newItem
	}

	return newItems
}
