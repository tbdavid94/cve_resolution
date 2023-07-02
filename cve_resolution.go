package cve_resolution

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
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
	Operator string      `json:"operator"`
	Value    string      `json:"value"`
	And      []Condition `json:"and"`
	Or       []Condition `json:"or"`
}

func GetOppositeOperator(operator string) string {
	switch operator {
	case ">":
		return "<"
	case ">=":
		return "<="
	case "<":
		return ">"
	case "<=":
		return ">="
	default:
		return operator
	}
}

func CheckCompareVersions(versionA string, operator string, versionB string) (bool, error) {
	partsA := strings.Split(versionA, ".")
	partsB := strings.Split(versionB, ".")

	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}

	for i := 0; i < maxLen; i++ {
		var partA, partB int

		if i < len(partsA) {
			numA, err := strconv.Atoi(partsA[i])
			if err != nil {
				return false, fmt.Errorf("Invalid version format for versionA: %s", versionA)
			}
			partA = numA
		}

		if i < len(partsB) {
			numB, err := strconv.Atoi(partsB[i])
			if err != nil {
				return false, fmt.Errorf("Invalid version format for versionB: %s", versionB)
			}
			partB = numB
		}

		switch {
		case partA < partB:
			return operator == "<", nil
		case partA > partB:
			return operator == ">", nil
		}
	}

	switch operator {
	case "==":
		return true, nil
	case "!=":
		return false, nil
	case "<", "<=":
		return true, nil
	case ">", ">=":
		return false, nil
	default:
		return false, fmt.Errorf("Invalid operator: %s", operator)
	}
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

func CalculatorProductResolution(cves []CVEProduct) (resolution []Condition, combinedCondition string, err error) {
	for _, cve := range cves {
		if cve.Version != "*" {
			resolution = append(resolution, Condition{Operator: "!=", Value: cve.Version})
		} else {
			_, convert := ConvertToResolution(cve)
			resolution = append(resolution, convert)
		}

	}

	SortConditions(resolution)

	//Kiem tra dieu kien co nam trong khoang bi vuln cua cac version khong
	var checkResolution []Condition
	for i, r := range resolution {
		versionA := r.Value
		operatorVersionA := r.Operator
		areaVulnVersionA := GetOppositeOperator(operatorVersionA) //khoang bi vuln cua version nay
		//check xem version con lai co nam trong khoang bi vuln nay khong
		for j, nr := range resolution {
			if i == j {
				continue
			}
			checkCompare, _ := CheckCompareVersions(nr.Value, areaVulnVersionA, versionA)
			if checkCompare {
				checkResolution = append(checkResolution, nr)
			}
		}
	}

	if len(checkResolution) > 0 {
		var newResolution []Condition
		for _, r := range resolution {
			search := GetCondition(checkResolution, r)
			if !search {
				newResolution = append(newResolution, r)
			}
		}

		resolution = newResolution
	}

	if len(resolution) == 0 {
		return resolution, "", nil
	}

	//optimize 1 la lay ra nhung operator > va >=
	var conditions1 []Condition
	var conditions2 []Condition
	for _, c := range resolution {
		if c.Operator == ">" || c.Operator == ">=" {
			conditions1 = append(conditions1, c)
		}
		if c.Operator == "<" || c.Operator == "<=" {
			conditions2 = append(conditions2, c)
		}
	}
	optimize1, _ := OptimizeCondition1(conditions1)
	optimize2, _ := OptimizeCondition2(conditions2)
	optimize1vs2, _ := OptimizeCondition1Vs2(optimize1, optimize2)

	//la lay ra operate !=
	conditionNotEqualTo := GetConditionNotEqualTo(resolution, "!=")
	conditionFirst := Condition{}
	if len(conditionNotEqualTo) > 0 {
		conditionFirst = conditionNotEqualTo[0]
	}

	if len(optimize1vs2) > 0 {
		var resolutionFinal []Condition
		for _, opt := range optimize1vs2 {
			simplifyCondition := SimplifyCondition(conditionFirst, opt)
			for _, sc := range simplifyCondition {
				resolutionFinal = append(resolutionFinal, sc)
			}
		}

		resolution = resolutionFinal
	} else {
		if len(conditionNotEqualTo) == 0 {
			resolution = optimize1vs2
		}
	}

	// Kết hợp các điều kiện
	for i, res := range resolution {
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
		default:
			//logger.GetInstance().WithFields(logrus.Fields{"operator": operator}).Info("Operator undefined")
			continue
		}

		if len(res.Or) > 0 {
			for _, or := range res.Or {
				//result, err := CompareVersions(numValue, or.Value)
				//if err != nil {
				// logger.GetInstance().WithFields(logrus.Fields{"err": err}).Info("Error compare 2 versions")
				// return []Condition{}, "", err
				//}
				//if (operator == ">=" || operator == ">") && (or.Operator == "<=" || or.Operator == "<") && result == 1 {
				// continue
				//}
				condition += " || version " + or.Operator + " " + or.Value
			}
		}

		// Kết hợp điều kiện với điều kiện trước đó (nếu có)
		if i > 0 {
			combinedCondition += " && "
		}
		combinedCondition += condition
	}

	return
}

func GetCondition(conditions []Condition, conditionSearch Condition) (status bool) {
	for _, r := range conditions {
		if conditionSearch.Operator == r.Operator && conditionSearch.Value == r.Value {
			return true
		}
	}

	return false
}

func UnsetResolution(arr []Condition, index int) []Condition {
	if index < 0 || index >= len(arr) {
		return arr
	}

	// Sử dụng phép cắt mảng để tạo một mảng mới bỏ qua phần tử ở vị trí index
	return append(arr[:index], arr[index+1:]...)
}

func GetConditionNotEqualTo(conditions []Condition, value string) (uniqueConditions []Condition) {
	var cdt []Condition
	for _, r := range conditions {
		if r.Operator == "!=" {
			cdt = append(cdt, r)
		}
	}
	uniqueConditions = RemoveDuplicates(cdt)
	sort.Slice(uniqueConditions, func(i, j int) bool {
		compare, _ := CompareVersions(uniqueConditions[i].Value, uniqueConditions[j].Value)
		return compare > 0
	})

	return
}

func RemoveDuplicates(conditions []Condition) []Condition {
	// Tạo một map để lưu trữ các giá trị đã xuất hiện
	seen := make(map[string]bool)
	result := []Condition{}

	for _, cond := range conditions {
		// Kiểm tra nếu giá trị đã xuất hiện, bỏ qua nếu đã tồn tại trong map
		if seen[cond.Value] {
			continue
		}

		// Thêm giá trị vào map và vào mảng kết quả
		seen[cond.Value] = true
		result = append(result, cond)
	}

	return result
}

func SimplifyCondition(conditionFirst, conditionOptimize Condition) (condition []Condition) {
	versionA := conditionFirst.Value
	versionB := conditionOptimize.Value

	if versionA == "" {
		return []Condition{
			conditionOptimize,
		}
	}

	operator := conditionOptimize.Operator

	result, err := CompareVersions(versionA, versionB)
	if err != nil {
		fmt.Println(err)
		return
	}

	condition = []Condition{
		conditionOptimize,
	}

	//if (operator == ">" || operator == ">=") && result == -1 {
	// condition = []Condition{
	// conditionOptimize,
	// }
	//}

	if (operator == "<" || operator == "<=") && result == -1 {
		condition = []Condition{
			conditionFirst,
			conditionOptimize,
		}
	}

	if (operator == ">" || operator == ">=") && result == 1 {
		condition = []Condition{
			conditionFirst,
			conditionOptimize,
		}
	}

	//if (operator == "<" || operator == "<=") && result == 1 {
	// condition = []Condition{
	// conditionOptimize,
	// }
	//}

	if len(conditionOptimize.Or) > 0 {
		versionC := conditionOptimize.Or[0].Value
		result2, err := CompareVersions(versionA, versionC)
		if err != nil {
			fmt.Println(err)
			return
		}
		if result2 == -1 {
			condition = []Condition{
				conditionOptimize,
				conditionFirst,
			}
		}

	}

	return
}

func OptimizeCondition1(conditions []Condition) (optimize Condition, err error) {

	if len(conditions) == 1 {
		return conditions[0], nil
	}

	if len(conditions) == 2 {
		//lay ra > va >=
		//lon hon
		greaterThan := conditions[0].Value
		//lon hon hoac bang
		greaterThanOrEqualTo := conditions[1].Value

		result, err := CompareVersions(greaterThan, greaterThanOrEqualTo)
		if err != nil {
			//logger.GetInstance().WithFields(logrus.Fields{"err": err}).Info("Error compare 2 versions")
			return Condition{}, err
		}

		if result == -1 {
			optimize = conditions[1]
		} else if result == 1 {
			optimize = Condition{
				Operator: conditions[0].Operator,
				Value:    conditions[0].Value,
				Or: []Condition{
					{
						Operator: conditions[1].Operator,
						Value:    conditions[1].Value,
					},
					{
						Operator: "<",
						Value:    conditions[0].Value,
					},
				},
			}
		} else {
			optimize = conditions[1]
		}

		//if greaterThan < greaterThanOrEqualTo || greaterThan == greaterThanOrEqualTo {
		// optimize = conditions[1]
		//} else if greaterThan > greaterThanOrEqualTo {
		// optimize = Condition{
		// Operator: conditions[0].Operator,
		// Value: conditions[0].Value,
		// Or: []Condition{
		// {
		// Operator: conditions[1].Operator,
		// Value: conditions[1].Value,
		// },
		// {
		// Operator: "<",
		// Value: conditions[0].Value,
		// },
		// },
		// }
		//}
	}

	if len(conditions) > 2 {
		//lay value > max
		sort.Slice(conditions, func(i, j int) bool {
			compare, _ := CompareVersions(conditions[i].Value, conditions[j].Value)
			return compare > 0
		})
		optimize = conditions[0]
	}
	return
}

func OptimizeCondition2(conditions []Condition) (optimize Condition, err error) {

	if len(conditions) == 1 {
		return conditions[0], nil
	}

	if len(conditions) == 2 {
		//lay ra < va <=
		//nho hon
		lessThan := conditions[0].Value
		//nho hon hoac bang
		lessThanOrEqualTo := conditions[1].Value

		result, err := CompareVersions(lessThan, lessThanOrEqualTo)
		if err != nil {
			//logger.GetInstance().WithFields(logrus.Fields{"err": err}).Info("Error compare 2 versions")
			return Condition{}, err
		}

		if result == -1 {
			optimize = Condition{
				Operator: conditions[0].Operator,
				Value:    conditions[0].Value,
				//And: []Condition{
				// {
				// Operator: "!=",
				// Value: conditions[0].Value,
				// },
				//},
			}
		} else if result == 1 {
			optimize = conditions[0]
		} else {
			optimize = conditions[1]
		}
	}

	if len(conditions) > 2 {
		//lay value < min
		sort.Slice(conditions, func(i, j int) bool {
			compare, _ := CompareVersions(conditions[i].Value, conditions[j].Value)
			return compare < 0
		})
		optimize = conditions[0]
	}

	return
}

func OptimizeCondition1Vs2(optimizeCondition1, optimizeCondition2 Condition) (optimize []Condition, err error) {

	if optimizeCondition1.Operator == "" && optimizeCondition2.Operator != "" {
		return []Condition{
			optimizeCondition2,
		}, nil
	}

	if optimizeCondition1.Operator != "" && optimizeCondition2.Operator == "" {
		return []Condition{
			optimizeCondition1,
		}, nil
	}

	versionA := optimizeCondition1.Value
	versionB := optimizeCondition2.Value

	result, err := CompareVersions(versionA, versionB)
	if err != nil {
		//logger.GetInstance().WithFields(logrus.Fields{"err": err}).Info("Error compare 2 versions")
		return
	}

	optimize = []Condition{
		optimizeCondition1,
		optimizeCondition2,
	}
	if optimizeCondition1.Operator == ">=" && optimizeCondition2.Operator == "<=" {
		if result == -1 {
			optimize = []Condition{
				optimizeCondition1,
				optimizeCondition2,
			}
		} else if result == 1 {
			optimize = []Condition{
				{
					Operator: optimizeCondition1.Operator,
					Value:    optimizeCondition1.Value,
					Or: []Condition{
						{
							Operator: optimizeCondition2.Operator,
							Value:    optimizeCondition2.Value,
						},
					},
					And: optimizeCondition2.And,
				},
			}
		} else {
			optimize = []Condition{
				{
					Operator: "=",
					Value:    optimizeCondition1.Value,
					Or:       optimizeCondition2.And,
				},
			}
		}
	} else if (optimizeCondition1.Operator == ">=" || optimizeCondition1.Operator == ">") && (optimizeCondition2.Operator == "<" || optimizeCondition2.Operator == "<=") {
		if result == -1 {
			optimize = []Condition{
				optimizeCondition1,
				optimizeCondition2,
			}
		} else if result == 1 {
			optimize = []Condition{
				{
					Operator: optimizeCondition1.Operator,
					Value:    optimizeCondition1.Value,
					Or: []Condition{
						optimizeCondition2,
					},
				},
			}
		} else {
			optimize = []Condition{
				optimizeCondition1,
			}
		}
	} else if optimizeCondition1.Operator == ">=" && optimizeCondition2.Operator == "<=" {
		if result == -1 {
			optimize = []Condition{
				optimizeCondition1,
				optimizeCondition2,
			}
		} else if result == 1 {
			optimize = []Condition{
				{
					Operator: optimizeCondition1.Operator,
					Value:    optimizeCondition1.Value,
					Or: []Condition{
						optimizeCondition2,
					},
				},
			}
		} else {
			optimize = []Condition{
				{
					Operator: "=",
					Value:    optimizeCondition1.Value,
				},
			}
		}
	} else if optimizeCondition1.Operator == "<" && optimizeCondition2.Operator == ">" {
		if result == -1 {
			optimize = []Condition{
				{
					Operator: optimizeCondition1.Operator,
					Value:    optimizeCondition1.Value,
					Or: []Condition{
						{
							Operator: optimizeCondition2.Operator,
							Value:    optimizeCondition2.Value,
						},
					},
				},
			}
		} else if result == 1 {
			optimize = []Condition{
				optimizeCondition1,
				optimizeCondition2,
			}
		} else {
			optimize = []Condition{}
		}
	}

	return
}

func CompareVersions(versionA, versionB string) (int, error) {
	aSegments := strings.Split(versionA, ".")
	bSegments := strings.Split(versionB, ".")

	maxLen := len(aSegments)
	if len(bSegments) > maxLen {
		maxLen = len(bSegments)
	}

	for i := 0; i < maxLen; i++ {
		aVal := 0
		if i < len(aSegments) {
			var err error
			aVal, err = strconv.Atoi(aSegments[i])
			if err != nil {
				return 0, fmt.Errorf("version A undefined")
			}
		}

		bVal := 0
		if i < len(bSegments) {
			var err error
			bVal, err = strconv.Atoi(bSegments[i])
			if err != nil {
				return 0, fmt.Errorf("version B undefined")
			}
		}

		if aVal < bVal {
			return -1, nil
		} else if aVal > bVal {
			return 1, nil
		}
	}

	return 0, nil
}
