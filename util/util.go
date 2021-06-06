package util

// GenerateTwoSliceDiff generates a slice of elements which the first slice contains and the second do not contains.
// If there is no element to return, this function returns a nil slice.
func GenerateTwoSliceDiff(a []string, b []string) []string {
	bmap := make(map[string]bool)
	for _, e := range b {
		if !bmap[e] {
			bmap[e] = true
		}
	}
	var diff []string
	for _, e := range a {
		if !bmap[e] {
			diff = append(diff, e)
		}
	}
	return diff
}

// GenerateUniqueSlice generates a unique string slice from argument string slice.
// If there is no element to return, this function returns a nil slice.
func GenerateUniqueSlice(org []string) []string {
	amap := make(map[string]bool)
	var uniqueSlice []string
	for _, e := range org {
		if !amap[e] {
			uniqueSlice = append(uniqueSlice, e)
			amap[e] = true
		}
	}
	return uniqueSlice
}
