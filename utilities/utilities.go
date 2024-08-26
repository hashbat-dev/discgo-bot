package utilities

// Contains a bunch of funcs etc that aren't in go but we use a lot

func StringSliceContainsElement(slice []string, search string) bool {
	for _, val := range slice {
		if val == search {
			return true
		}
	}
	return false
}
