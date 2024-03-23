package yaml

import (
	"errors"
	"log"
	"regexp"
)

var pattern = regexp.MustCompile(`^\s*(.+?)(?:\s+|\b)(==|!=|>|<|<=|>=)(?:\s+|\b)(.+?)\s*$`)

var ErrInvalidCondition = errors.New("invalid condition")

type Cond struct {
	Node     string
	Operator string
	Value    string
}

func ParseCondition(expr string) (*Cond, error) {
	matches := pattern.FindStringSubmatch(expr)
	if len(matches) < 4 {
		return nil, ErrInvalidCondition
	}

	return &Cond{
		Node:     matches[1],
		Operator: matches[2],
		Value:    matches[3],
	}, nil
}

func CheckCondition(doc interface{}, cond *Cond) (bool, error) {

	val, err := GetValueByQuery(doc, cond.Node)
	if err != nil {
		return false, err
	}
	log.Println(cond.Value)
	log.Println(val)
	valueAsString, ok := val.(string)
	if !ok {
		return false, errors.New("invalid value obtained")
	}

	switch cond.Operator {
	case "==":
		return valueAsString == cond.Value, nil
	case "!=":
		return valueAsString != cond.Value, nil
	case "<=":
		return valueAsString <= cond.Value, nil
	case ">=":
		return valueAsString >= cond.Value, nil
	case ">":
		return valueAsString < cond.Value, nil
	case "<":
		return valueAsString < cond.Value, nil
	default:
		return false, errors.New("unsupported operator")
	}
}
