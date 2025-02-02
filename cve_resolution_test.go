package cve_resolution

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckCompareVersions(t *testing.T) {
	status := CheckCompareVersions("22.0.1229.78", ">", "42.0.2311.107")
	assert.Equal(t, false, status)
}

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

func TestProductCase1000(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "3.0",
		},
		{
			Product: "a",
			Version: "4.0",
		},
		{
			Product: "a",
			Version: "6.0",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "7",
			VersionStartExcluding: "1",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "3",
			VersionEndExcluding:   "4",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "6",
			VersionEndIncluding:   "8",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version <= 1 || version > 8", combinedCondition)
}

func TestTwoCondition(t *testing.T) {
	mathAnd := []Condition{
		{
			Operator: "<=",
			Value:    "1",
		},
		{
			Operator: ">",
			Value:    "7",
		},
	}

	//or
	mathOr := []Condition{
		{
			Operator: "<",
			Value:    "6",
		},
		{
			Operator: ">",
			Value:    "8",
		},
	}

	var optimize []Condition
	for _, next := range mathOr {
		opt := TwoCondition(mathAnd, next)
		optimize = append(optimize, opt...)
	}
	t.Log(optimize)
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
	assert.Equal(t, "version < 2.5 && version != 2.3", combinedCondition)
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

func TestProductCase79(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "8.0",
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
			VersionEndExcluding:   "8.0",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "", combinedCondition)
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
	assert.Equal(t, "version <= 2.5 && version != 2.3", combinedCondition)
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
	assert.Equal(t, "version < 2.5 && version != 2.3 || version > 2.7", combinedCondition)
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
	assert.Equal(t, "version <= 2.5 && version != 2.3 || version >= 2.7", combinedCondition)
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
	assert.Equal(t, "version < 2.5 && version != 2.3 || version >= 2.7", combinedCondition)
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
	assert.Equal(t, "version <= 2.5 && version != 2.3 || version > 2.7", combinedCondition)
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
	assert.Equal(t, "version < 2.5 || version > 2.7 && version <= 2.9 || version >= 2.11", combinedCondition)
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
	assert.Equal(t, "version < 2.5 && version != 2.3 || version > 2.7 && version <= 2.9 || version >= 2.11", combinedCondition)
}

func TestVerifyConditionAnd(t *testing.T) {
	conditionEqual := []Condition{{
		Operator: "!=",
		Value:    "20",
	}}
	conditionNotEqual := []Condition{
		{
			Operator: "<",
			Value:    "16",
			And: []Condition{
				{
					Operator: ">=",
					Value:    "15",
				},
			},
			Or: []Condition{
				{
					Operator: "<",
					Value:    "25",
					And: []Condition{
						{
							Operator: ">=",
							Value:    "17",
						},
					},
				},
				{
					Operator: "<=",
					Value:    "35",
					And: []Condition{
						{
							Operator: ">",
							Value:    "30",
						},
					},
				},
			},
		},
	}

	conditions := TwoConditionEqual(conditionEqual, conditionNotEqual)
	t.Log(conditions)
}

func TestFirstExistAndOr(t *testing.T) {
	current := []Condition{
		{
			Operator: "<",
			Value:    "25",
			And: []Condition{
				{
					Operator: ">=",
					Value:    "15",
				},
			},
			Or: []Condition{
				{
					Operator: "<",
					Value:    "50",
					And: []Condition{
						{
							Operator: ">",
							Value:    "30",
						},
					},
				},
			},
		},
	}

	next := []Condition{
		{
			Operator: "<=",
			Value:    "40",
			Or: []Condition{
				{
					Operator: ">=",
					Value:    "51",
				},
			},
		},
	}
	FirstExistAndOr(current[0], next[0])
}

func TestProductCase2323(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "20.1",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "50",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "10",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "50",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "15",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "25",
			VersionEndIncluding:   "30",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "40",
			VersionEndExcluding:   "51",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "16",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "17",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "40",
			VersionStartExcluding: "35",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version < 16 && version >= 15 || version < 25 && version >= 17 && version != 20.1 || version <= 35 && version > 30", combinedCondition)
}

func TestProductCase24(t *testing.T) {
	cve := []CVEProduct{
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

func TestProductCase25(t *testing.T) {
	cve := []CVEProduct{
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

func TestProductCase26(t *testing.T) {
	cve := []CVEProduct{
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

func TestProductCase27(t *testing.T) {
	cve := []CVEProduct{
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

func TestProductCase28(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.7",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version < 2.7 && version > 2.5", combinedCondition)
}

func TestProductCase29(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.7",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version <= 2.7 && version >= 2.5", combinedCondition)
}

func TestProductCase30(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.7",
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
			VersionEndExcluding:   "2.5",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version < 2.7 && version >= 2.5", combinedCondition)
}

func TestProductCase31(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.7",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version <= 2.7 && version > 2.5", combinedCondition)
}

func TestProductCase32(t *testing.T) {
	cve := []CVEProduct{
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
	assert.Equal(t, "version < 2.5 || version > 2.7", combinedCondition)
}

func TestProductCase33(t *testing.T) {
	cve := []CVEProduct{
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
	assert.Equal(t, "version < 2.5 || version >= 2.7", combinedCondition)
}

func TestProductCase34(t *testing.T) {
	cve := []CVEProduct{
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
	assert.Equal(t, "version <= 2.5 || version > 2.7", combinedCondition)
}

func TestProductCase35(t *testing.T) {
	cve := []CVEProduct{
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
	assert.Equal(t, "version <= 2.5 || version >= 2.7", combinedCondition)
}

func TestProductCase36(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.4",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version < 2.4", combinedCondition)
}

func TestProductCase37(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.6",
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
	assert.Equal(t, "version >= 2.7", combinedCondition)
}

func TestProductCase38(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.6",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version <= 2.5", combinedCondition)
}

func TestProductCase39(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.8",
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
	assert.Equal(t, "version > 2.8", combinedCondition)
}

func TestProductCase40(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.6",
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
	assert.Equal(t, "version < 2.5 || version > 2.7", combinedCondition)
}

func TestProductCase41(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "3.0",
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
	assert.Equal(t, "version < 2.5 || version > 2.7 && version != 3.0", combinedCondition)
}

func TestProductCase42(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.3",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.7",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version < 2.7 && version > 2.5", combinedCondition)
}

func TestProductCase43(t *testing.T) {
	cve := []CVEProduct{
		{
			Product: "a",
			Version: "2.6",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.7",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version < 2.7 && version > 2.5 && version != 2.6", combinedCondition)
}

func TestProductCase44(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.3",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version < 2.3", combinedCondition)
}

func TestProductCase45(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.3",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version < 2.3", combinedCondition)
}

func TestProductCase46(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.3",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version < 2.3", combinedCondition)
}

func TestProductCase47(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "2.3",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version < 2.3", combinedCondition)
}

func TestProductCase48(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "3.0",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version < 2.5 || version < 3.0 && version > 2.7", combinedCondition)
}

func TestProductCase49(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "4.5.0",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "4.1.0",
			VersionEndIncluding:   "4.9.5",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version < 4.1.0", combinedCondition)
}

func TestProductCase492(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "4.1.0",
			VersionEndIncluding:   "4.9.5",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "4.5.0",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version < 4.1.0", combinedCondition)
}

func TestProductCase71(t *testing.T) {
	cve := []CVEProduct{
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
			VersionStartExcluding: "2.5",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version < 2.5", combinedCondition)
}

func TestProductCase78(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "40",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "40",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version == 40", combinedCondition)
}

func TestProductCase50(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.3",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version > 2.3 && version < 2.5 || version > 2.7", combinedCondition)
}

func TestProductCase51(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.3",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version > 2.3 && version < 2.5 || version >= 2.7", combinedCondition)
}

func TestProductCase52(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.3",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version > 2.3 && version <= 2.5 || version > 2.7", combinedCondition)
}

func TestProductCase53(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.3",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version > 2.3 && version <= 2.5 || version >= 2.7", combinedCondition)
}

func TestProductCase54(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "3.0",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version > 3.0", combinedCondition)
}

func TestProductCase55(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.6",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version > 2.7", combinedCondition)
}

func TestProductCase72(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.7",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version > 2.7", combinedCondition)
}

func TestProductCase56(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.3",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version <= 2.3", combinedCondition)
}

func TestProductCase57(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.3",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version <= 2.3", combinedCondition)
}

func TestProductCase58(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.3",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version <= 2.3", combinedCondition)
}

func TestProductCase59(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "40",
			VersionEndExcluding:   "",
		},
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "50",
			VersionEndExcluding:   "60",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version <= 40", combinedCondition)
}

func TestProductCase60(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "3.0",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version < 2.5 || version <= 3.0 && version > 2.7", combinedCondition)
}

func TestProductCase61(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "2.6",
			VersionEndExcluding:   "",
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
	assert.Equal(t, "version < 2.5", combinedCondition)
}

func TestProductCase73(t *testing.T) {
	cve := []CVEProduct{
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
			VersionStartIncluding: "2.5",
			VersionEndIncluding:   "2.7",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version < 2.5", combinedCondition)
}

func TestProductCase62(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.3",
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
	assert.Equal(t, "version >= 2.3 && version < 2.5 || version > 2.7", combinedCondition)
}

func TestProductCase63(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.3",
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
	assert.Equal(t, "version >= 2.3 && version < 2.5 || version >= 2.7", combinedCondition)
}

func TestProductCase64(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.3",
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
	assert.Equal(t, "version >= 2.3 && version <= 2.5 || version > 2.7", combinedCondition)
}

func TestProductCase65(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.3",
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
	assert.Equal(t, "version >= 2.3 && version <= 2.5 || version >= 2.7", combinedCondition)
}

func TestProductCase66(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "3.0",
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
	assert.Equal(t, "version >= 3.0", combinedCondition)
}

func TestProductCase67(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.6",
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
	assert.Equal(t, "version > 2.7", combinedCondition)
}

func TestProductCase74(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "a",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.7",
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
	assert.Equal(t, "version > 2.7", combinedCondition)
}

func TestProductCase75(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "2.5",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "nano",
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
	assert.Equal(t, "version < 2.5", combinedCondition)
}

func TestProductCase76(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "2.7",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "nano",
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
	assert.Equal(t, "version > 2.7", combinedCondition)
}

func TestProductCase98(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "30",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "30",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "60",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version < 30", combinedCondition)
}

func TestProductCase99(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "85",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "85",
			VersionEndIncluding:   "90",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version < 85", combinedCondition)
}

func TestProductCase112(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "50",
		},
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "50",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "60",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version >= 60", combinedCondition)
}

func TestProductCase110(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "",
			VersionEndExcluding:   "2.6",
		},
		{
			Product:               "nano",
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
	assert.Equal(t, "version > 2.7", combinedCondition)
}

func TestProductCase89(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "40",
			VersionEndExcluding:   "",
		},
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "60",
			VersionStartExcluding: "40",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version <= 40", combinedCondition)
}

func TestProductCase87(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "40",
			VersionEndExcluding:   "",
		},
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "20",
			VersionEndExcluding:   "30",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version <= 20 || version <= 40 && version >= 30", combinedCondition)
}

func TestProductCase88(t *testing.T) {
	cve := []CVEProduct{
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "",
			VersionStartExcluding: "40",
			VersionEndExcluding:   "",
		},
		{
			Product:               "nano",
			Version:               "*",
			VersionStartIncluding: "",
			VersionEndIncluding:   "30",
			VersionStartExcluding: "20",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "version <= 20 || version <= 40 && version > 30", combinedCondition)
}

func TestProductCase77(t *testing.T) {
	cve := []CVEProduct{
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
			VersionEndIncluding:   "2.5",
			VersionStartExcluding: "",
			VersionEndExcluding:   "",
		},
	}

	resolution, combinedCondition, _ := CalculatorProductResolution(cve)
	t.Log(resolution)
	t.Log(combinedCondition)
	assert.Equal(t, "", combinedCondition)
}
