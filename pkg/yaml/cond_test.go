package yaml_test

import (
	"reflect"
	"testing"

	yamlak "github.com/Slimo300/yamlak/pkg/yaml"
)

func TestParseCondition(t *testing.T) {
	tests := []struct {
		expr          string
		expectedCond  *yamlak.Cond
		expectedError error
	}{
		{"metadata.name == emails", &yamlak.Cond{"metadata.name", "==", "emails"}, nil},
		{"value != 10", &yamlak.Cond{"value", "!=", "10"}, nil},
		{"timestamp > 2024-01-01", &yamlak.Cond{"timestamp", ">", "2024-01-01"}, nil},
		{"value <= 100", &yamlak.Cond{"value", "<=", "100"}, nil},
		{"invalid expression", nil, yamlak.ErrInvalidCondition},
	}

	for _, test := range tests {
		cond, err := yamlak.ParseCondition(test.expr)
		if err != test.expectedError {
			t.Errorf("Expected error: %v, but got: %v", test.expectedError, err)
		}

		if !reflect.DeepEqual(cond, test.expectedCond) {
			t.Errorf("For expression: %s, expected condition: %+v, but got: %+v", test.expr, test.expectedCond, cond)
		}
	}
}
