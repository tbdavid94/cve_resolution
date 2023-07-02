package cve_resolution

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertToResolution(t *testing.T) {
	cve := CVEProduct{
		Product: "a",
		Version: "2.3",
	}
	cve1 := CVEProduct{
		Product:               "a",
		Version:               "2.3",
		VersionStartExcluding: "3.3",
	}
	cve2 := CVEProduct{
		Product:               "a",
		Version:               "2.3",
		VersionStartExcluding: "3.3",
		VersionEndExcluding:   "3.4",
	}

	count, _ := ConvertToResolution(cve)
	count1, _ := ConvertToResolution(cve1)
	count2, _ := ConvertToResolution(cve2)
	assert.Equal(t, 0, count)
	assert.Equal(t, 1, count1)
	assert.Equal(t, 2, count2)
}

func TestProductCase1(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
	}
	t.Log(cve)
	_, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(combinedCondition)
	assert.Equal(t, "version != 2.3", combinedCondition)
}

func TestProductCase2(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.2", // < 2.2
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version < 2.2", combinedCondition)
}

func TestProductCase3(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.3",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version > 2.3", combinedCondition)
}

func TestProductCase4(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.4",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version <= 2.4", combinedCondition)
}

func TestProductCase5(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.5",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version >= 2.5", combinedCondition)
}

func TestProductCase6(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.5",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version != 2.3 && version < 2.5", combinedCondition)
}

func TestProductCase7(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.5",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version > 2.5", combinedCondition)
}

func TestProductCase8(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.5",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version != 2.3 && version <= 2.5", combinedCondition)
}

func TestProductCase9(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.5",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version >= 2.5", combinedCondition)
}

func TestProductCase10(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.5",
			VersionEndIncluding:   "2.7",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version < 2.5 || version > 2.7 && version != 2.3", combinedCondition)
}

func TestProductCase11(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.5",
			VersionEndExcluding:   "2.7",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version <= 2.5 || version >= 2.7 && version != 2.3", combinedCondition)
}

func TestProductCase12(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.5",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.7",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version < 2.5 || version >= 2.7 && version != 2.3", combinedCondition)
}

func TestProductCase13(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.7",
			VersionStartExcluding: "2.5",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version <= 2.5 || version > 2.7 && version != 2.3", combinedCondition)
}

func TestProductCase14(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.5",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.7",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "", combinedCondition)
}

func TestProductCase15(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.5",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.7",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "", combinedCondition)
}

func TestProductCase16(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.5",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.7",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "", combinedCondition)
}

func TestProductCase17(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.5",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.7",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "", combinedCondition)
}

func TestProductCase18(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.5",
			VersionEndIncluding:   "2.7",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.9",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version > 2.7 && version != 2.3 && version <= 2.9", combinedCondition)
}

func TestProductCase19(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.5",
			VersionEndIncluding:   "2.7",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.9",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version >= 2.9 || version < 2.5", combinedCondition)
}

func TestProductCase20(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.5",
			VersionEndExcluding:   "2.7",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.9",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version >= 2.7 && version != 2.3 && version < 2.9", combinedCondition)
}

func TestProductCase21(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.5",
			VersionEndExcluding:   "2.7",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.9",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version > 2.9", combinedCondition)
}

func TestProductCase22(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.5",
			VersionEndIncluding:   "2.7",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.9",
			VersionEndExcluding:   "2.11",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version >= 2.11 || version <= 2.9", combinedCondition)
}

func TestProductCase23(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.5",
			VersionEndIncluding:   "2.7",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.9",
			VersionEndExcluding:   "2.11",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version >= 2.11 || version <= 2.9", combinedCondition)
}

func TestOptimizeCondition1(t *testing.T) {
	conditionCase1 := []Condition{
		{
			Operator: ">",
			Value:    "2",
		},
		{
			Operator: ">=",
			Value:    "7",
		},
	}

	conditionCase2 := []Condition{
		{
			Operator: ">",
			Value:    "7",
		},
		{
			Operator: ">=",
			Value:    "2",
		},
	}

	conditionCase3 := []Condition{
		{
			Operator: ">",
			Value:    "7",
		},
		{
			Operator: ">=",
			Value:    "7",
		},
	}

	optimize1, _ := OptimizeCondition1(conditionCase1)
	optimize2, _ := OptimizeCondition1(conditionCase2)
	optimize3, _ := OptimizeCondition1(conditionCase3)
	t.Log(optimize1)
	t.Log(optimize2)
	t.Log(optimize3)
	assert.Equal(t, Condition{Operator: ">=", Value: "7"}, optimize1)
	assert.Equal(t, Condition{Operator: ">", Value: "7", Or: []Condition{
		{Operator: ">=", Value: "2"},
		{Operator: "<", Value: "7"},
	}}, optimize2)
	assert.Equal(t, Condition{Operator: ">=", Value: "7"}, optimize3)
}

func TestOptimizeCondition2(t *testing.T) {
	conditionCase1 := []Condition{
		{
			Operator: "<",
			Value:    "2",
		},
		{
			Operator: "<=",
			Value:    "7",
		},
	}

	conditionCase2 := []Condition{
		{
			Operator: "<",
			Value:    "7",
		},
		{
			Operator: "<=",
			Value:    "2",
		},
	}

	conditionCase3 := []Condition{
		{
			Operator: "<",
			Value:    "7",
		},
		{
			Operator: "<=",
			Value:    "7",
		},
	}

	optimize1, _ := OptimizeCondition2(conditionCase1)
	optimize2, _ := OptimizeCondition2(conditionCase2)
	optimize3, _ := OptimizeCondition2(conditionCase3)
	t.Log(optimize1)
	t.Log(optimize2)
	t.Log(optimize3)

	assert.Equal(t, Condition{Operator: "<=", Value: "7", And: []Condition{
		{Operator: "!=", Value: "2"},
	}}, optimize1)
	assert.Equal(t, Condition{Operator: "<", Value: "7"}, optimize2)
	assert.Equal(t, Condition{Operator: "<=", Value: "7"}, optimize3)
}

// 1
// output 1 optimize condition 1 vs output 1 optimize condition 2
func TestOptimizeCondition1Vs2CaseCd11vsCaseCd21(t *testing.T) {
	optimize1ConditionCase1 := []Condition{
		{
			Operator: ">",
			Value:    "2",
		},
		{
			Operator: ">=",
			Value:    "7",
		},
	} //ra >=7

	optimize2ConditionCase1 := []Condition{
		{
			Operator: "<",
			Value:    "2",
		},
		{
			Operator: "<=",
			Value:    "7",
		},
	}

	optimize2ConditionCase2 := []Condition{
		{
			Operator: "<",
			Value:    "2",
		},
		{
			Operator: "<=",
			Value:    "9",
		},
	}

	optimize2ConditionCase3 := []Condition{
		{
			Operator: "<",
			Value:    "2",
		},
		{
			Operator: "<=",
			Value:    "4",
		},
	}

	optimize1, _ := OptimizeCondition1(optimize1ConditionCase1)
	optimize21, _ := OptimizeCondition2(optimize2ConditionCase1)
	optimize22, _ := OptimizeCondition2(optimize2ConditionCase2)
	optimize23, _ := OptimizeCondition2(optimize2ConditionCase3)
	optimize1vs2Case1, _ := OptimizeCondition1Vs2(optimize1, optimize21)
	optimize1vs2Case2, _ := OptimizeCondition1Vs2(optimize1, optimize22)
	optimize1vs2Case3, _ := OptimizeCondition1Vs2(optimize1, optimize23)
	assert.Equal(t, []Condition{
		{
			Operator: "=",
			Value:    "7",
			Or: []Condition{
				{
					Operator: "!=",
					Value:    "2",
				},
			},
		},
	}, optimize1vs2Case1)

	assert.Equal(t, []Condition{
		{
			Operator: ">=",
			Value:    "7",
		},
		{
			Operator: "<=",
			Value:    "9",
			And: []Condition{
				{
					Operator: "!=",
					Value:    "2",
				},
			},
		},
	}, optimize1vs2Case2)

	assert.Equal(t, []Condition{
		{
			Operator: ">=",
			Value:    "7",
			And: []Condition{
				{
					Operator: "!=",
					Value:    "2",
				},
			},
			Or: []Condition{
				{
					Operator: "<=",
					Value:    "4",
				},
			},
		},
	}, optimize1vs2Case3)
}

// output 1 optimize condition 1 vs output 2 optimize condition 2
func TestOptimizeCondition1Vs2CaseCd11vsCaseCd22(t *testing.T) {
	optimizeCondition1Case1 := []Condition{
		{
			Operator: ">",
			Value:    "2",
		},
		{
			Operator: ">=",
			Value:    "7",
		},
	}

	optimizeCondition2Case2 := []Condition{
		{
			Operator: "<",
			Value:    "7",
		},
		{
			Operator: "<=",
			Value:    "2",
		},
	}

	optimizeCondition1, _ := OptimizeCondition1(optimizeCondition1Case1)
	t.Log(optimizeCondition1)
	optimizeCondition2, _ := OptimizeCondition2(optimizeCondition2Case2)
	optimizeCondition1Vs2, _ := OptimizeCondition1Vs2(optimizeCondition1, optimizeCondition2)
	t.Log(optimizeCondition1Vs2)
	assert.Equal(t, []Condition{
		{
			Operator: ">=",
			Value:    "7",
		},
	}, optimizeCondition1Vs2)
}

// output 1 optimize condition 1 vs output 3 optimize condition 2
func TestOptimizeCondition1Vs2CaseCd11vsCaseCd23(t *testing.T) {
	optimizeCondition1Case1 := []Condition{
		{
			Operator: ">",
			Value:    "2",
		},
		{
			Operator: ">=",
			Value:    "7",
		},
	}

	optimizeCondition2Case3 := []Condition{
		{
			Operator: "<",
			Value:    "7",
		},
		{
			Operator: "<=",
			Value:    "7",
		},
	}

	optimizeCondition1, _ := OptimizeCondition1(optimizeCondition1Case1)
	t.Log(optimizeCondition1)
	optimizeCondition2, _ := OptimizeCondition2(optimizeCondition2Case3)
	t.Log(optimizeCondition2)
	optimizeCondition1Vs2, _ := OptimizeCondition1Vs2(optimizeCondition1, optimizeCondition2)
	t.Log(optimizeCondition1Vs2)
	assert.Equal(t, []Condition{
		{Operator: "=", Value: "7"},
	}, optimizeCondition1Vs2)
}

func TestLessThanCompareVersions(t *testing.T) {
	versionA := "1.2.3"
	versionB := "1.2.4"

	result, err := CompareVersions(versionA, versionB)
	if err != nil {
		fmt.Println(err)
		return
	}

	ex := ""
	if result == -1 {
		fmt.Printf("Phiên bản A (%s) nhỏ hơn phiên bản B (%s)\n", versionA, versionB)
		ex = "A < B"
	} else if result == 1 {
		fmt.Printf("Phiên bản A (%s) lớn hơn phiên bản B (%s)\n", versionA, versionB)
		ex = "A > B"
	} else {
		fmt.Printf("Phiên bản A (%s) và phiên bản B (%s) bằng nhau\n", versionA, versionB)
		ex = "A = B"
	}
	assert.Equal(t, "A < B", ex)
}

func TestGreaterThanCompareVersions(t *testing.T) {
	versionA := "1.2.4"
	versionB := "1.2.2"

	result, err := CompareVersions(versionA, versionB)
	if err != nil {
		fmt.Println(err)
		return
	}

	ex := ""
	if result == -1 {
		fmt.Printf("Phiên bản A (%s) nhỏ hơn phiên bản B (%s)\n", versionA, versionB)
		ex = "A < B"
	} else if result == 1 {
		fmt.Printf("Phiên bản A (%s) lớn hơn phiên bản B (%s)\n", versionA, versionB)
		ex = "A > B"
	} else {
		fmt.Printf("Phiên bản A (%s) và phiên bản B (%s) bằng nhau\n", versionA, versionB)
		ex = "A = B"
	}
	assert.Equal(t, "A > B", ex)
}

func TestEqualToCompareVersions(t *testing.T) {
	versionA := "1.2.4"
	versionB := "1.2.4"

	result, err := CompareVersions(versionA, versionB)
	if err != nil {
		fmt.Println(err)
		return
	}

	ex := ""
	if result == -1 {
		fmt.Printf("Phiên bản A (%s) nhỏ hơn phiên bản B (%s)\n", versionA, versionB)
		ex = "A < B"
	} else if result == 1 {
		fmt.Printf("Phiên bản A (%s) lớn hơn phiên bản B (%s)\n", versionA, versionB)
		ex = "A > B"
	} else {
		fmt.Printf("Phiên bản A (%s) và phiên bản B (%s) bằng nhau\n", versionA, versionB)
		ex = "A = B"
	}
	assert.Equal(t, "A = B", ex)
}
