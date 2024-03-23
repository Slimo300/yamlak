package yaml_test

import (
	"fmt"
	"testing"

	"github.com/Slimo300/yamlak/pkg/yaml"
	"github.com/stretchr/testify/assert"
)

func TestIParseNode(t *testing.T) {

	testCases := map[string]struct {
		name    string
		indexes []int
		err     error
	}{
		"containers[0][1]":    {name: "containers", indexes: []int{0, 1}},
		"containers[0][]":     {err: fmt.Errorf("invalid node reference")},
		"containers[][1]":     {err: fmt.Errorf("invalid node reference")},
		"containers[][]":      {err: fmt.Errorf("invalid node reference")},
		"containers[0][1][2]": {name: "containers", indexes: []int{0, 1, 2}},
		"containers[0]":       {name: "containers", indexes: []int{0}},
		"containers[]":        {err: fmt.Errorf("invalid node reference")},
		"containers":          {name: "containers"},
		"contain[0]ers":       {err: fmt.Errorf("invalid node reference")},
		"contain[]ers":        {err: fmt.Errorf("invalid node reference")},
		"contain[ers":         {err: fmt.Errorf("invalid node reference")},
	}

	for node, expected := range testCases {
		t.Run(node, func(t *testing.T) {
			name, indexes, err := yaml.ParseNode(node)

			assert.Equal(t, expected.name, name)
			assert.Equal(t, expected.indexes, indexes)
			assert.Equal(t, expected.err, err)
		})
	}
}
