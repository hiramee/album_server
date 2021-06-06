package test

import (
	"album-server/util"
	"reflect"
	"testing"
)

// Test struct for TrimPrefixFromString.
type TrimPrefixFromStringTest struct {
	str      string
	prefix   string
	expected string
}

var trimPrefixFromString = []TrimPrefixFromStringTest{
	{"hoge/fuga", "hoge", "fuga"},
	{"hoge/", "hoge", ""},
	{"hoge", "hoge", "hoge"},
}

func TestTrimPrefixFromString(t *testing.T) {
	for i, test := range trimPrefixFromString {
		actual := *util.TrimPrefixFromString(&test.str, &test.prefix)
		if actual != test.expected {
			t.Errorf("#%d got: %#v want: %#v", i, actual, test.expected)
		}
	}
}

// Test struct for GenerateTwoSliceDiff.
type GenerateTwoSliceDiffTest struct {
	amap     []string
	bmap     []string
	expected []string
}

var generateTwoSliceDiffSlice = []GenerateTwoSliceDiffTest{
	{[]string{"a"}, []string{"a"}, nil},
	{[]string{"a", "b"}, []string{"b"}, []string{"a"}},
	{nil, []string{"b"}, nil},
	{[]string{"a"}, nil, []string{"a"}},
	{nil, nil, nil},
}

func TestGenerateTwoSliceDiff(t *testing.T) {
	for i, test := range generateTwoSliceDiffSlice {
		actual := util.GenerateTwoSliceDiff(test.amap, test.bmap)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("#%d got: %#v want: %#v", i, actual, test.expected)
		}
	}
}

// Test struct for GenerateTwoSliceDiff.
type GenerateUniqueSliceTest struct {
	org      []string
	expected []string
}

var generateUniqueSlice = []GenerateUniqueSliceTest{
	{[]string{"a"}, []string{"a"}},
	{[]string{"a", "a"}, []string{"a"}},
	{[]string{"a", "b", "b", "a"}, []string{"a", "b"}},
	{nil, nil},
}

func TestGenerateUniqueSlice(t *testing.T) {
	for i, test := range generateUniqueSlice {
		actual := util.GenerateUniqueSlice(test.org)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("#%d got: %#v want: %#v", i, actual, test.expected)
		}
	}
}
