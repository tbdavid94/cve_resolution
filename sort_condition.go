package cve_resolution

import (
	"sort"
)

type Conditions []Condition

func (c Conditions) Len() int {
	return len(c)
}

func (c Conditions) Less(i, j int) bool {
	operators := map[string]int{
		">":  5,
		">=": 4,
		"<":  3,
		"<=": 2,
		"!=": 1,
	}

	opI := operators[c[i].Operator]
	opJ := operators[c[j].Operator]

	return opI > opJ
}

func (c Conditions) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// SortConditions Sắp xếp mảng Conditions theo Operator từ >, >=, <, <=, !=
func SortConditions(conditions []Condition) {
	sort.Sort(Conditions(conditions))
}
