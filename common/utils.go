package common

// ContainsStringArray return true if the dest is a subset of src, otherwise return false and unmatched string
func ContainsStringArray(src []string, dest []string) (string, bool) {
	stringSet := make(map[string]struct{})
	for _, str := range src {
		stringSet[str] = struct{}{}
	}

	for _, str := range dest {
		if _, ok := stringSet[str]; !ok {
			return str, false
		}
	}

	return "", true
}
