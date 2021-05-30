package util

import "testing"

func TestTrimPrefixFromString1(t *testing.T) {
	prefix := "hoge"
	str := "hoge/fuga"
	if *TrimPrefixFromString(&str, &prefix) != "fuga" {
		t.Fatal("failed test")
	}
}
func TestTrimPrefixFromStringToBlank(t *testing.T) {
	prefix := "hoge"
	str := "hoge/"
	if *TrimPrefixFromString(&str, &prefix) != "" {
		t.Fatal("failed test")
	}
}
func TestTrimPrefixFromStringFail(t *testing.T) {
	prefix := "hoge"
	str := "hoge"
	if *TrimPrefixFromString(&str, &prefix) != str {
		t.Fatal("failed test")
	}
}
