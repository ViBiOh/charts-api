package readingtag

// IncludesUint64 checks in an array includes given uint64
func IncludesUint64(array []uint64, lookup uint64) bool {
	for _, item := range array {
		if item == lookup {
			return true
		}
	}

	return false
}
