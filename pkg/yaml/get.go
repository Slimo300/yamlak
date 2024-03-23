package yaml

import (
	"strings"
)

func GetValueByQuery(doc interface{}, query string) (interface{}, error) {

	nodes := strings.Split(query, ".")

	for _, node := range nodes {
		name, indexes, err := ParseNode(node)
		if err != nil {
			return nil, err
		}

		docAsMap, ok := doc.(map[string]interface{})
		if !ok {
			return nil, ErrValueNotFound
		}
		doc = docAsMap[name]

		for _, index := range indexes {
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
