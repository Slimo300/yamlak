package yaml_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	yamlak "github.com/Slimo300/yamlak/pkg/yaml"
	"gopkg.in/yaml.v3"
)

func TestSetValueByQuery(t *testing.T) {
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
		value  string
		result interface{}
	}{
		// Testy dla plików YAML z pojedynczymi tablicami
		{yaml1, "spec.template.spec.containers[0].image", "updated_image", `spec:
  template:
    spec:
      containers:
      - image: updated_image
`},
		{yaml2, "spec.template.spec.containers[0][0].image", "updated_image", `
spec:
  template:
    spec:
      containers:
      - - image: updated_image
`},
		{yaml3, "spec.template.spec.containers[0].ports[0].containerPort", "8080", `
spec:
  template:
    spec:
      containers:
      - image: nginx
        ports:
        - containerPort: "8080"
`},
		{yaml4, "spec.template.spec.containers[0][0].ports[0].containerPort", "8080", `
spec:
  template:
    spec:
      containers:
      - - image: nginx
          ports:
          - containerPort: "8080"
`},
		{yaml5, "spec.template.spec.containers[0][0][0].image", "updated_image", `
spec:
  template:
    spec:
      containers:
      - - - image: updated_image
`},
		{yaml6, "spec.template.spec.containers[0][0][0].ports[0].containerPort", "8080", `
spec:
  template:
    spec:
      containers:
      - - - image: nginx
            ports:
            - containerPort: "8080"
`},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("Test %d", i+1), func(t *testing.T) {
			var err error
			test.yaml, err = DecodeYAMLFromString(test.yaml)
			if err != nil {
				t.Errorf("Error decoding input yaml: %v", err)
			}
			test.result, err = DecodeYAMLFromString(test.result)
			if err != nil {
				t.Errorf("Error decoding expected result yaml: %v", err)
			}

			if err := yamlak.SetValueByQuery(test.yaml, test.query, test.value); err != nil {
				t.Errorf("Error setting value: %v", err)
			}

			received, err := StringifyYAML(test.yaml)
			if err != nil {
				t.Errorf("Error stringifying received response: %v", err)
			}
			expected, err := StringifyYAML(test.result)
			if err != nil {
				t.Errorf("Error stringifying expected response: %v", err)
			}

			if !reflect.DeepEqual(received, expected) {
				t.Errorf("Expected output: %v, Received output: %v", expected, received)
			}
		})
	}
}

// StringifyYAML is a helper function
func StringifyYAML(doc interface{}) (string, error) {
	builder := &strings.Builder{}

	if err := yaml.NewEncoder(builder).Encode(&doc); err != nil {
		return "", err
	}

	return builder.String(), nil
}

func DecodeYAMLFromString(doc interface{}) (interface{}, error) {

	docAsString, ok := doc.(string)
	if !ok {
		return nil, errors.New("couldn't convert to string")
	}
	if err := yaml.NewDecoder(strings.NewReader(docAsString)).Decode(&doc); err != nil {
		return nil, err
	}

	return doc, nil

}
