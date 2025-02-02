package cve_resolution

import (
	"cve_resolution/version"
	"fmt"
	"sort"
)

type CVEProduct struct {
	Product               string `json:"product"`
	Version               string `json:"version"`
	VersionStartIncluding string `json:"version_start_including"`
	VersionStartExcluding string `json:"version_start_excluding"`
	VersionEndIncluding   string `json:"version_end_including"`
	VersionEndExcluding   string `json:"version_end_excluding"`
}

type Condition struct {
	Operator    string      `json:"operator"`
	Value       string      `json:"value"`
	And         []Condition `json:"and"`
	Or          []Condition `json:"or"`
	OriginIndex int         `json:"origin_index"`
}

func CheckCompareVersions(versionA string, operator string, versionB string) bool {
	vA, err := version.NewVersion(versionA)
	if err != nil {
		// Xử lý lỗi khi phiên bản versionA không hợp lệ
		return false
	}

	vB, err := version.NewVersion(versionB)
	if err != nil {
		// Xử lý lỗi khi phiên bản versionB không hợp lệ
		return false
	}

	switch operator {
	case "<":
		return vA.LessThan(vB)
	case "<=":
		return vA.LessThanOrEqual(vB)
	case ">":
		return vA.GreaterThan(vB)
	case ">=":
		return vA.GreaterThanOrEqual(vB)
	case "==":
		return vA.Equal(vB)
	case "!=":
		return !vA.Equal(vB)
	default:
		// Xử lý lỗi khi toán tử không hợp lệ
		return false
	}
}

func GetConditionEqual(conditions []Condition) (conditionsEqual []Condition) {
	for _, condition := range conditions {
		if condition.Operator == "!=" {
			conditionsEqual = append(conditionsEqual, condition)
		}
	}
	return
}

func GetConditionNotEqual(conditions []Condition) (conditionsEqual []Condition) {
	for _, condition := range conditions {
		if (len(condition.And) > 0 || len(condition.Or) > 0) && condition.Operator != "!=" {
			conditionsEqual = append(conditionsEqual, condition)
		}
	}
	return
}

func GetConditionOne(conditions []Condition) (conditionsOne []Condition) {
	for _, condition := range conditions {
		if len(condition.And) == 0 && len(condition.Or) == 0 && condition.Operator != "!=" {
			conditionsOne = append(conditionsOne, condition)
		}
	}
	return
}

func ConvertToResolution(cve CVEProduct) (count int, resolution Condition) {
	if cve.VersionStartIncluding != "" {
		count++
	}
	if cve.VersionEndIncluding != "" {
		count++
	}
	if cve.VersionStartExcluding != "" {
		count++
	}
	if cve.VersionEndExcluding != "" {
		count++
	}

	if count == 1 {
		if cve.VersionStartIncluding != "" {
			resolution = Condition{Operator: "<", Value: cve.VersionStartIncluding}
		}
		if cve.VersionEndIncluding != "" {
			resolution = Condition{Operator: ">", Value: cve.VersionEndIncluding}
		}
		if cve.VersionStartExcluding != "" {
			resolution = Condition{Operator: "<=", Value: cve.VersionStartExcluding}
		}
		if cve.VersionEndExcluding != "" {
			resolution = Condition{Operator: ">=", Value: cve.VersionEndExcluding}
		}
	}

	if count == 2 {
		if cve.VersionStartIncluding != "" && cve.VersionEndIncluding != "" {
			resolution = Condition{
				Operator: "<",
				Value:    cve.VersionStartIncluding,
				Or:       []Condition{{Operator: ">", Value: cve.VersionEndIncluding}},
			}
		}
		if cve.VersionStartExcluding != "" && cve.VersionEndExcluding != "" {
			resolution = Condition{
				Operator: "<=",
				Value:    cve.VersionStartExcluding,
				Or:       []Condition{{Operator: ">=", Value: cve.VersionEndExcluding}},
			}
		}
		if cve.VersionStartIncluding != "" && cve.VersionEndExcluding != "" {
			resolution = Condition{
				Operator: "<",
				Value:    cve.VersionStartIncluding,
				Or:       []Condition{{Operator: ">=", Value: cve.VersionEndExcluding}},
			}
		}
		if cve.VersionStartExcluding != "" && cve.VersionEndIncluding != "" {
			resolution = Condition{
				Operator: "<=",
				Value:    cve.VersionStartExcluding,
				Or:       []Condition{{Operator: ">", Value: cve.VersionEndIncluding}},
			}
		}

	}
	return
}

func OptimizeTwoVersion(version1, version2 Condition) (versionSelect Condition) {
	version1Value := version1.Value
	operatorVersion1 := version1.Operator

	version2Value := version2.Value
	operatorVersion2 := version2.Operator

	// >, >=, <, <=
	// >, >
	if operatorVersion1 == ">" && operatorVersion2 == ">" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = version2
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = version1
		}
	}

	// >, >=
	if operatorVersion1 == ">" && operatorVersion2 == ">=" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = version2
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = version1
		}
	}

	// >, <
	if operatorVersion1 == ">" && operatorVersion2 == "<" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			version1.And = []Condition{version2}
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
	}

	// >, <=
	if operatorVersion1 == ">" && operatorVersion2 == "<=" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			version1.And = []Condition{version2}
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
	}

	// >=, >=
	if operatorVersion1 == ">=" && operatorVersion2 == ">=" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = version2
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = Condition{Operator: "==", Value: version1Value}
		}
	}

	// >=, >
	if operatorVersion1 == ">=" && operatorVersion2 == ">" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = version2
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = version2
		}
	}

	// >=, <
	if operatorVersion1 == ">=" && operatorVersion2 == "<" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			version1.And = []Condition{version2}
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
	}

	// >=, <=
	if operatorVersion1 == ">=" && operatorVersion2 == "<=" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			version1.And = []Condition{version2}
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = Condition{Operator: "==", Value: version1Value}
		}
	}

	// <, <
	if operatorVersion1 == "<" && operatorVersion2 == "<" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = version2
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = version1
		}
	}

	// <, >
	if operatorVersion1 == "<" && operatorVersion2 == ">" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			version1.And = []Condition{version2}
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
	}

	// <, >=
	if operatorVersion1 == "<" && operatorVersion2 == ">=" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			version1.And = []Condition{version2}
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
	}

	// <, <=
	if operatorVersion1 == "<" && operatorVersion2 == "<=" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = version2
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = version1
		}
	}

	// <=, <=
	if operatorVersion1 == "<=" && operatorVersion2 == "<=" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = version2
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = version1
		}
	}

	// <=, >
	if operatorVersion1 == "<=" && operatorVersion2 == ">" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			version1.And = []Condition{version2}
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
	}

	// <=, >=
	if operatorVersion1 == "<=" && operatorVersion2 == ">=" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = Condition{} //vo nghiem
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			version1.And = []Condition{version2}
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = Condition{Operator: "==", Value: version1Value}
		}
	}

	// <=, <
	if operatorVersion1 == "<=" && operatorVersion2 == "<" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = version2
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = version2
		}
	}

	//!= >
	if operatorVersion1 == "!=" && operatorVersion2 == ">" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = version2
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			version1.And = []Condition{version2}
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = version2
		}
	}

	//!= >=
	if operatorVersion1 == "!=" && operatorVersion2 == ">=" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			versionSelect = version2
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			version1.And = []Condition{version2}
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = Condition{Operator: ">", Value: version1Value}
		}
	}

	//!= <
	if operatorVersion1 == "!=" && operatorVersion2 == "<" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			version1.And = []Condition{version2}
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = version2
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = version2
		}
	}

	//!= <=
	if operatorVersion1 == "!=" && operatorVersion2 == "<=" {
		if CheckCompareVersions(version1Value, "<", version2Value) {
			version1.And = []Condition{version2}
			versionSelect = version1
		}
		if CheckCompareVersions(version1Value, ">", version2Value) {
			versionSelect = version2
		}
		if CheckCompareVersions(version1Value, "==", version2Value) {
			versionSelect = Condition{Operator: "<", Value: version1Value}
		}
	}

	return
}

// VerifyConditionAnd kiem tra condition != co nam trong khoang condition and khong
func VerifyConditionAnd(conditionsAnd []Condition, conditionEqual Condition) (conditions []Condition) {
	flagPass := 0
	for _, and := range conditionsAnd {
		select1 := OptimizeTwoVersion(conditionEqual, and)
		if select1.Operator != "" && select1.Value != "" {
			if len(select1.And) > 0 {
				flagPass += 1
			}
		}
	}
	//ca 2 deu phai ket hop voi dieu kien != thi => condition != nam trong khoang can ket hop
	if flagPass == len(conditionsAnd) {
		conditionsAnd = append(conditionsAnd, conditionEqual)
	}
	for i, and := range conditionsAnd {
		if i == 0 {
			conditions = append(conditions, and)
		} else {
			conditions[0].And = append(conditions[0].And, and)
		}
	}
	return
}

func UnsetItems(myStructs []Condition, indexes []int) []Condition {
	// Sắp xếp các chỉ số theo thứ tự giảm dần để không ảnh hưởng đến việc xóa các phần tử
	sort.Sort(sort.Reverse(sort.IntSlice(indexes)))

	// Xóa các phần tử từ slice ban đầu theo các chỉ số
	for _, index := range indexes {
		myStructs = append(myStructs[:index], myStructs[index+1:]...)
	}

	return myStructs
}

func ConditionExist(conditions []Condition, condition Condition) (status bool) {
	status = false
	for _, c := range conditions {
		if c.Operator == condition.Operator && c.Value == condition.Value {
			status = true
			return
		}
	}
	return
}

func VerifyOptimize(optimize []Condition) (conditions []Condition) {
	//rut gon duoc cac condtion trung lap
	var itemAlone []Condition
	var itemCollapse []Condition
	for i, opt := range optimize {
		opt.OriginIndex = i
		if len(opt.Or) == 0 && len(opt.And) == 0 {
			itemAlone = append(itemAlone, opt)
		} else {
			itemCollapse = append(itemCollapse, opt)
		}
	}

	if len(itemAlone) == 0 {
		conditions = itemCollapse
		return
	}

	var itemIndexRemove []int
	for _, itemA := range itemAlone {
		for _, itemC := range itemCollapse {
			if itemA.Operator == itemC.Operator && itemA.Value == itemC.Value || (len(itemC.Or) > 0 && itemC.Or[0].Operator == itemA.Operator && itemC.Or[0].Value == itemA.Value) || (len(itemC.And) > 0 && itemC.And[0].Operator == itemA.Operator && itemC.And[0].Value == itemA.Value) {
				itemIndexRemove = append(itemIndexRemove, itemA.OriginIndex)
			}
		}
	}

	if len(itemIndexRemove) > 0 {
		optimize = UnsetItems(optimize, itemIndexRemove)
	}
	conditions = optimize
	return
}

func TwoCondition(mathAnd []Condition, next Condition) (optimize []Condition) {
	//(a+b) x c 1 trong 2 phep nhan a x c hoac b x c vo nghiem thi break
	isNoSolution := false
	for _, m := range mathAnd {
		select1 := OptimizeTwoVersion(m, next)
		if select1.Operator != "" && select1.Value != "" {
			if len(select1.And) > 0 {
				if len(mathAnd) == 1 {
					optimize = append(optimize, select1)
				} else {
					optimize = append(optimize, m)
				}
			} else {
				optimize = append(optimize, select1)
			}
		} else {
			//vo nghiem
			if len(mathAnd) > 1 {
				optimize = append(optimize, m)
			}
			isNoSolution = true
			break
		}
	}
	if isNoSolution {
		return []Condition{}
	}
	if len(optimize) == 2 {
		optimize = []Condition{
			{
				Operator: optimize[0].Operator,
				Value:    optimize[0].Value,
				And: []Condition{
					{
						Operator: optimize[1].Operator,
						Value:    optimize[1].Value,
					},
				},
			},
		}
		if len(optimize[0].And) == 0 && len(optimize[0].Or) == 0 && len(optimize[1].And) == 0 && len(optimize[1].Or) == 0 {
			optimize = []Condition{
				{
					Operator: optimize[0].Operator,
					Value:    optimize[0].Value,
					Or: []Condition{
						{
							Operator: optimize[1].Operator,
							Value:    optimize[1].Value,
						},
					},
				},
			}
		}
	}
	return
}

func TwoConditionEqual(conditionsEqual []Condition, optimizeConditionNotEqual []Condition) (optimize []Condition) {
	if len(optimizeConditionNotEqual) > 0 {
		for _, equal := range conditionsEqual {
			condition1 := []Condition{
				{
					Operator: optimizeConditionNotEqual[0].Operator,
					Value:    optimizeConditionNotEqual[0].Value,
				},
			}
			if len(optimizeConditionNotEqual[0].And) > 0 {
				itemConditionAnd := optimizeConditionNotEqual[0].And[0]
				condition1 = append(condition1, Condition{
					Operator: itemConditionAnd.Operator,
					Value:    itemConditionAnd.Value,
				})
			}
			optimize1 := VerifyConditionAnd(condition1, equal)
			var optimize2 []Condition

			if len(optimizeConditionNotEqual[0].Or) > 0 {
				for _, or := range optimizeConditionNotEqual[0].Or {
					condition2 := []Condition{
						{
							Operator: or.Operator,
							Value:    or.Value,
						},
					}
					if len(or.And) > 0 {
						itemConditionAnd := or.And[0]
						condition2 = append(condition2, Condition{
							Operator: itemConditionAnd.Operator,
							Value:    itemConditionAnd.Value,
						})
					}
					opt := VerifyConditionAnd(condition2, equal)
					optimize2 = append(optimize2, opt...)
				}
			}

			optimize = append(optimize1, optimize2...)
		}
	} else {
		optimize = conditionsEqual
	}

	return
}

func ConvertOneConditionOr(optimize []Condition) (oneConditionOr []Condition) {
	for i, opt := range optimize {
		if i == 0 {
			oneConditionOr = append(oneConditionOr, opt)
		} else {
			oneConditionOr[0].Or = append(oneConditionOr[0].Or, opt)
		}
	}
	return
}

func CalculatorProductResolution(cves []CVEProduct) (resolution []Condition, combinedCondition string, err error) {
	for _, cve := range cves {
		if cve.Version != "*" {
			resolution = append(resolution, Condition{Operator: "!=", Value: cve.Version})
		} else {
			_, convert := ConvertToResolution(cve)
			resolution = append(resolution, convert)
		}

	}

	var finalOptimize []Condition

	//lay ra cac condition co 1 dieu kien
	conditionsOne := GetConditionOne(resolution)
	var afterOptimizeConditionOne []Condition

	for {
		if len(conditionsOne) <= 1 {
			afterOptimizeConditionOne = conditionsOne
			break
		}
		current := conditionsOne[0]
		next := conditionsOne[1]
		var optimize []Condition

		if len(current.And) > 0 {
			math := []Condition{
				{
					Operator: current.Operator,
					Value:    current.Value,
				},
				{
					Operator: current.And[0].Operator,
					Value:    current.And[0].Value,
				},
			}
			optimize = TwoCondition(math, next)
			fmt.Println(optimize)
		} else {
			versionSelect := OptimizeTwoVersion(current, next)
			//khong vo nghiem thi append
			if versionSelect.Operator != "" && versionSelect.Value != "" {
				optimize = append(optimize, versionSelect)
			}
		}

		//set lai array resolution moi
		indexes := []int{0, 1} // Các chỉ số phần tử cần xóa
		conditionsOne = UnsetItems(conditionsOne, indexes)
		if len(optimize) > 0 {
			newConditionsOne := []Condition{optimize[0]}
			newConditionsOne = append(newConditionsOne, conditionsOne...)
			conditionsOne = newConditionsOne
		}
	}

	//lay ra danh sach cac condition !=
	conditionsEqual := GetConditionEqual(resolution)

	//lay ra danh sach cac condition khac !=
	conditionsNotEqualSelect := GetConditionNotEqual(resolution)

	//neu condition one vo nghiem thi doi chieu condition != co nam trong khoang vo nghiem day khong
	if len(afterOptimizeConditionOne) == 0 && len(conditionsNotEqualSelect) <= 0 {
		for i, equal := range conditionsEqual {
			conditionsOneOrigin := GetConditionOne(resolution)
			for _, conditionOneOrigin := range conditionsOneOrigin {
				versionSelect := OptimizeTwoVersion(equal, conditionOneOrigin)
				if versionSelect.Operator == equal.Operator && versionSelect.Value == equal.Value {
					conditionsEqual = UnsetItems(conditionsEqual, []int{i})
					break
				}
			}
		}
	}

	//gop conditionsOne vao conditionsNotEqual
	conditionsNotEqual := append(afterOptimizeConditionOne, conditionsNotEqualSelect...)

	flag := 1
	for {
		if len(conditionsNotEqual) <= 1 {
			break
		}
		current := conditionsNotEqual[0]
		next := conditionsNotEqual[1] //luon luon la OR
		var optimize []Condition

		//item dau tien la AND
		if len(current.And) > 0 && len(current.Or) == 0 {
			mathAnd := []Condition{
				{
					Operator: current.Operator,
					Value:    current.Value,
				},
				{
					Operator: current.And[0].Operator,
					Value:    current.And[0].Value,
				},
			}

			mathOr := []Condition{
				{
					Operator: next.Operator,
					Value:    next.Value,
				},
				{
					Operator: next.Or[0].Operator,
					Value:    next.Or[0].Value,
				},
			}

			var optOr []Condition
			for i := range mathOr {
				or := mathOr[i]
				optimize1 := TwoCondition(mathAnd, or)
				optOr = append(optOr, optimize1...)
			}
			optimize = append(optimize, optOr...)
			if len(optimize) >= 2 {
				optimize = ConvertOneConditionOr(optimize)
			}

			fmt.Println(optimize)

		}
		//item dau tien la OR
		if len(current.Or) > 0 && len(current.And) == 0 {
			mathOr := []Condition{
				{
					Operator: current.Operator,
					Value:    current.Value,
				},
				{
					Operator: current.Or[0].Operator,
					Value:    current.Or[0].Value,
				},
			}

			mathOrNext := []Condition{
				{
					Operator: next.Operator,
					Value:    next.Value,
				},
				{
					Operator: next.Or[0].Operator,
					Value:    next.Or[0].Value,
				},
			}

			var opt []Condition
			for _, or := range mathOr {
				for _, orNext := range mathOrNext {
					versionSelect := OptimizeTwoVersion(or, orNext)
					if versionSelect.Operator != "" && versionSelect.Value != "" {
						opt = append(opt, versionSelect)
					}
				}
			}
			if len(opt) >= 2 {
				optimize1 := ConvertOneConditionOr(opt)
				optimize = append(optimize, optimize1...)
			} else {
				optimize = append(optimize, opt...)
			}

			fmt.Println(opt)
		}

		if len(current.And) == 0 && len(current.Or) == 0 {
			mathAnd := []Condition{
				{
					Operator: current.Operator,
					Value:    current.Value,
				},
			}

			mathOr := []Condition{
				{
					Operator: next.Operator,
					Value:    next.Value,
				},
				{
					Operator: next.Or[0].Operator,
					Value:    next.Or[0].Value,
				},
			}

			var optOr []Condition
			for i := range mathOr {
				or := mathOr[i]
				optimize1 := TwoCondition(mathAnd, or)
				optOr = append(optOr, optimize1...)
			}
			optimize = append(optimize, optOr...)
			if len(optimize) >= 2 {
				optimize = UniqueConditions(optimize)
				optimize = ConvertOneConditionOr(optimize)
			}
		}

		if len(current.And) > 0 && len(current.Or) > 0 {
			optimize = FirstExistAndOr(current, next)
		}

		//set lai array resolution moi
		indexes := []int{0, 1} // Các chỉ số phần tử cần xóa
		conditionsNotEqual = UnsetItems(conditionsNotEqual, indexes)
		if len(optimize) > 0 {
			newConditionsNotEqual := []Condition{optimize[0]}
			newConditionsNotEqual = append(newConditionsNotEqual, conditionsNotEqual...)
			conditionsNotEqual = newConditionsNotEqual
		}
		flag++
	}

	//case dac biet
	//if len(conditionsNotEqual) > 0 && len(conditionsNotEqual[0].And) > 0 && len(conditionsNotEqual[0].Or) > 0 {
	// conditionsNotEqual = SpecialOptimize(conditionsNotEqual)
	//}
	//doi chieu voi tung condition != de ra resolution final
	if len(conditionsEqual) > 0 {
		finalOptimize = TwoConditionEqual(conditionsEqual, conditionsNotEqual)
	} else {
		finalOptimize = conditionsNotEqual
	}

	finalOptimize = UniqueConditions(finalOptimize)
	// Kết hợp các điều kiện
	for i, res := range finalOptimize {
		operator := res.Operator
		numValue := res.Value

		// Tạo điều kiện tương ứng
		condition := ""

		switch operator {
		case "<":
			condition = fmt.Sprintf("version < %v", numValue)
		case "<=":
			condition = fmt.Sprintf("version <= %v", numValue)
		case ">":
			condition = fmt.Sprintf("version > %v", numValue)
		case ">=":
			condition = fmt.Sprintf("version >= %v", numValue)
		case "!=":
			condition = fmt.Sprintf("version != %v", numValue)
		case "==":
			condition = fmt.Sprintf("version == %v", numValue)
		default:
			//logger.GetInstance().WithFields(logrus.Fields{"operator": operator}).Info("Operator undefined")
			continue
		}

		if len(res.And) > 0 {
			for _, and := range res.And {
				condition += " && version " + and.Operator + " " + and.Value
			}
		}

		if len(res.Or) > 0 {
			for _, or := range res.Or {
				condition += " || version " + or.Operator + " " + or.Value
				if len(or.And) > 0 {
					for _, and := range or.And {
						condition += " && version " + and.Operator + " " + and.Value
					}
				}
			}
		}

		// Kết hợp điều kiện với điều kiện trước đó (nếu có)
		if i > 0 {
			combinedCondition += " || "
		}
		combinedCondition += condition
	}

	return
}

func FirstExistAndOr(current, next Condition) (optimize []Condition) {
	mathAnd := []Condition{
		{
			Operator: current.Operator,
			Value:    current.Value,
		},
	}
	for _, ma := range current.And {
		mathAnd = append(mathAnd, ma)
	}

	opt1 := TwoCondition(mathAnd, next)
	fmt.Println(opt1)

	var optOr []Condition
	for _, or := range current.Or {
		mathAnd2 := []Condition{
			{
				Operator: or.Operator,
				Value:    or.Value,
			},
		}
		if len(or.And) > 0 {
			for _, o := range or.And {
				mathAnd2 = append(mathAnd2, o)
			}
		}
		opt2 := TwoCondition(mathAnd2, next)
		optOr = append(optOr, opt2...)
	}

	optimize = append(opt1, optOr...)
	optimize = ConvertOneConditionOr(optimize)
	return
}

func SpecialOptimize(conditionsNotEqual []Condition) (optimize []Condition) {
	mathAnd := []Condition{
		{
			Operator: conditionsNotEqual[0].Operator,
			Value:    conditionsNotEqual[0].Value,
		},
	}

	for _, and := range conditionsNotEqual[0].And {
		mathAnd = append(mathAnd, Condition{
			Operator: and.Operator,
			Value:    and.Value,
		})
	}

	mathOr := []Condition{
		{
			Operator: conditionsNotEqual[0].Or[0].Operator,
			Value:    conditionsNotEqual[0].Or[0].Value,
		},
	}

	for _, or := range conditionsNotEqual[0].Or[0].And {
		mathOr = append(mathOr, Condition{
			Operator: or.Operator,
			Value:    or.Value,
		})
	}
	//cac phep OR nhan voi nhau
	for _, a := range mathAnd {
		for _, o := range mathOr {
			selectVersion := OptimizeTwoVersion(a, o)
			if selectVersion.Operator != "" && selectVersion.Value != "" {
				optimize = append(optimize, selectVersion)
			}
		}
	}
	optimize = UniqueConditions(optimize)
	optimize = ConvertOneConditionOr(optimize)
	return
}

func UniqueConditions(conditions []Condition) []Condition {
	uniqueMap := make(map[string]struct{})
	var uniqueConditions []Condition

	for _, condition := range conditions {
		key := fmt.Sprintf("%s:%s", condition.Operator, condition.Value)
		if _, ok := uniqueMap[key]; !ok {
			uniqueMap[key] = struct{}{}
			uniqueConditions = append(uniqueConditions, condition)
		}
	}

	return uniqueConditions
}

func GetCondition(conditions []Condition, conditionSearch Condition) (status bool, index int) {
	for i, r := range conditions {
		if conditionSearch.Operator == r.Operator && conditionSearch.Value == r.Value {
			index = i
			return true, index
		}
	}

	return false, 0
}
