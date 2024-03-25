package yaml_test

import (
	"reflect"
	"testing"

	yamlak "github.com/Slimo300/yamlak/pkg/yaml"
)

func TestDeleteValueByQuery(t *testing.T) {
	// Sample YAML
	yamlData := `
apiVersion: v1
kind: Pod
metadata:
  name: example
spec:
  containers:
    - name: nginx
      image: nginx:1.19.7
  volumes:
    - name: data
      emptyDir: {}
`

	var err error
	var doc interface{}
	doc, err = DecodeYAMLFromString(yamlData)
	if err != nil {
		t.Errorf("Error decoding yaml: %v", err)
	}

	// Testing deleting node existing path
	if err := yamlak.DeleteValueByQuery(doc, "spec.containers[0].image"); err != nil {
		t.Errorf("Error deleting node: %v", err)
	}

	// Checking whether the value was successfully deleted
	expectedYAML := `
apiVersion: v1
kind: Pod
metadata:
  name: example
spec:
  containers:
    - name: nginx
  volumes:
    - name: data
      emptyDir: {}
`
	var expectedDoc interface{}
	expectedDoc, err = DecodeYAMLFromString(expectedYAML)
	if err != nil {
		t.Errorf("Error decoding expected yaml %v", err)
	}

	if !reflect.DeepEqual(doc, expectedDoc) {
		t.Errorf("Result: %v does not match expected result: %v", doc, expectedDoc)
	}

	// Testing deleting for path that does not exist
	if err := yamlak.DeleteValueByQuery(doc, "spec.containers[0].resources.limits.memory"); err == nil {
		t.Error("Expected ValueNotFound error for unexisting path")
	} else if err.Error() != "value not found" {
		t.Errorf("Incorrect error returned, expected 'value not found', received '%s'", err.Error())
	}

	if err := yamlak.DeleteValueByQuery(doc, "spec.volumes[0]"); err != nil {
		t.Errorf("Error deleting member of a table: %v", err)
	}

	// Checking whether the element of array was properly erased
	expectedYAML = `
apiVersion: v1
kind: Pod
metadata:
  name: example
spec:
  containers:
    - name: nginx
  volumes: []
`
	expectedDoc, err = DecodeYAMLFromString(expectedYAML)
	if err != nil {
		t.Fatalf("Error decoding expected yaml: %v", err)
	}

	if !reflect.DeepEqual(doc, expectedDoc) {
		t.Errorf("Result: %v does not match expected result: %v", doc, expectedDoc)
	}
}
