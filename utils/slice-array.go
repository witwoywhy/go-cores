package utils

func SliceArray[T any, Slice ~[]T](items Slice, maxPerSlice int) []Slice {
	var response []Slice
	var ln = len(items)

	if ln == 0 {
		return response
	}

	for start := 0; start < ln; start += maxPerSlice {
		end := start + maxPerSlice
		if end > ln {
			end = ln
		}
		response = append(response, items[start:end])
	}

	return response
}
