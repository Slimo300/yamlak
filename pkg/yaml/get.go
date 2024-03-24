package yaml

import (
	"strings"
)

// GetValueByQuery takes in doc which should be a yaml document and a node path which it should obtain from it
// It returns a value it found and error information
func GetValueByQuery(doc interface{}, query string) (interface{}, error) {

	nodes := strings.Split(query, ".")

	for _, node := range nodes {
		name, indexes, err := ParseNode(node)
		if err != nil {
			return nil, err
		}

		// Since we expect node to be a string with optional slice referencing doc should be a map
		docAsMap, ok := doc.(map[string]interface{})
		if !ok {
			return nil, ErrValueNotFound
		}
		doc = docAsMap[name] // we extract next level of yaml file and assign it to doc as the rest of a file doesn't matter anymore, with this behavior in the end "doc" will be just the node we are looking for

		for _, index := range indexes {
			// if doc has indexes then current yaml node is an array, so we convert it to slice and cut off higher level as we did before
			docAsArray, ok := doc.([]interface{})
			if !ok {
				return nil, ErrValueNotFound
			}
			if len(docAsArray) <= index {
				return nil, ErrValueNotFound
			}
			doc = docAsArray[index]
		}
	}

	return doc, nil
}
