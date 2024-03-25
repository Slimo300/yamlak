package yaml_test

import (
	"fmt"
	"reflect"
	"testing"

	yamlak "github.com/Slimo300/yamlak/pkg/yaml"
)

func TestSetValueByQueryWithCreate(t *testing.T) {
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
    spec: node
`

	tests := []struct {
		yaml   interface{}
		query  string
		value  string
		result interface{}
	}{
		{yaml1, "spec.template.spec.containers[0].imageData.name", "some_name", `spec:
  template:
    spec:
      containers:
      - image: nginx
        imageData:
          name: some_name
`},
		{yaml1, "spec.template.spec.containers[1].someName", "some_name", `spec:
  template:
    spec:
      containers:
      - image: nginx
      - someName: some_name
`},
		{yaml2, "spec.template.containers[0]", "nginx", `
  spec:
    template:
      containers:
      - nginx
      spec: node
`},
		{yaml2, "spec.template[0]", "nginx", `
  spec:
    template:
    - nginx
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

			if err := yamlak.CreateValueByQuery(test.yaml, test.query, test.value); err != nil {
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
