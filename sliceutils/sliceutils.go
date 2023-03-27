package sliceutils

import (
	"fmt"
	"sort"
)

func SortStringSliceByKey(in []string) []string {
	keys := make([]string, 0, len(in))
	keys = append(keys, in...)

	sort.Strings(keys)

	return keys
}

func SortIntSliceByKey(in []int) []int {
	keys := make([]int, 0, len(in))
	keys = append(keys, in...)

	sort.Ints(keys)

	return keys
}

func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func TabList(items []string) string {
	var list string

	// Sort the list
	sort.Strings(items)

	l := len(items)
	for i, item := range items {
		list += fmt.Sprintf("%v", item)
		if i < l-1 {
			list += "\n\t"
		}
	}

	return list
}

func MinInt(s []int) int {
	if len(s) == 0 {
		return 0
	}

	var r int = s[0]

	for _, a := range s {
		if a < r {
			r = a
		}
	}

	return r
}

func MaxInt(s []int) int {
	if len(s) == 0 {
		return 0
	}

	var r int = s[0]

	for _, a := range s {
		if a > r {
			r = a
		}
	}

	return r
}

func SumInt(s []int) int {
	var r int

	for _, a := range s {
		r += a
	}

	return r
}
