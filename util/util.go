package util

import "unicode/utf8"

func GetTwoSliceDiff(a []string, b []string) []string {
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

func GetUniqueSlice(a []string) []string {
	amap := make(map[string]bool)
	var uniqueSlice []string
	for _, e := range a {
		if !amap[e] {
			uniqueSlice = append(uniqueSlice, e)
			amap[e] = true
		}
	}
	return uniqueSlice
}

// get XXX from prefix/XXX
func TrimPrefixFromString(str *string, prefix *string) *string {
	if len(*prefix) >= len(*str) {
		return str // cannot trim prefix
	}
	prefixCount := utf8.RuneCountInString(*prefix) + 1 // userName/fileName => fileName
	fileName := string([]rune(*str)[prefixCount:])
	return &fileName
}
