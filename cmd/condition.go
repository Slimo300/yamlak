package cmd

import yamlak "github.com/Slimo300/yamlak/pkg/yaml"

var conditions []string

func CheckConditions(doc interface{}, conditions []string) bool {
	if len(conditions) == 0 {
		return true
	}
	for _, cond := range conditions {
		cond, err := yamlak.ParseCondition(cond)
		if err != nil {
			return false
		}
		conditionPassed, err := yamlak.CheckCondition(doc, cond)
		if err != nil || !conditionPassed {
			return false
		}
	}

	return true
}
