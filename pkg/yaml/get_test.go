package yaml_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	yamlak "github.com/Slimo300/yamlak/pkg/yaml"
	"gopkg.in/yaml.v3"
)

func TestGetValueByQuery(t *testing.T) {
	// Poprawne pliki YAML
	yaml1 := `
spec:
  template:
    spec:
      containers:
      - image: nginx
`
	yaml2 := `
spec:
  template:
    spec:
      containers:
      - - image: nginx
`
	yaml3 := `
spec:
  template:
    spec:
      containers:
      - image: nginx
        ports:
        - containerPort: 80
`
	yaml4 := `
spec:
  template:
    spec:
      containers:
      - - image: nginx
          ports:
          - containerPort: 80
`
	yaml5 := `
spec:
  template:
    spec:
      containers:
      - - - image: nginx
`
	yaml6 := `
spec:
  template:
    spec:
      containers:
      - - - image: nginx
            ports:
            - containerPort: 80
`

	tests := []struct {
		yaml   interface{}
		query  string
		result interface{}
		err    error
	}{
		// Testy dla plików YAML z pojedynczymi tablicami
		{yaml1, "spec.template.spec.containers[0].image", "nginx", nil},
		{yaml2, "spec.template.spec.containers[0][0].image", "nginx", nil},
		{yaml3, "spec.template.spec.containers[0].ports[0].containerPort", 80, nil},
		{yaml4, "spec.template.spec.containers[0][0].ports[0].containerPort", 80, nil},
		{yaml5, "spec.template.spec.containers[0][0][0].image", "nginx", nil},
		{yaml6, "spec.template.spec.containers[0][0][0].ports[0].containerPort", 80, nil},

		// Test dla nieznalezionej ścieżki
		{yaml1, "spec.template.spec.containers[1].image", nil, yamlak.ErrValueNotFound},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test %d", i+1), func(t *testing.T) {

			docAsString, ok := test.yaml.(string)
			if !ok {
				t.Error("couldn't convert to byte")
			}

			if err := yaml.NewDecoder(strings.NewReader(docAsString)).Decode(&test.yaml); err != nil {
				t.Errorf("Decoding yaml error: %v", err)
			}

			result, err := yamlak.GetValueByQuery(test.yaml, test.query)
			if err != test.err {
				t.Errorf("Expected error: %v, Received error: %v", test.err, err)
			}

			if !reflect.DeepEqual(result, test.result) {
				t.Errorf("Expected result: %v, Received result: %v", test.result, result)
			}
		})
	}
}
