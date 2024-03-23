package yaml

import (
	"strings"
)

func SetValueByQuery(doc interface{}, query, value string) error {
	nodes := strings.Split(query, ".")
	current := doc // current is a value to use for extracting deeper levels of yaml body

	for node_iter, node := range nodes {
		name, indexes, err := ParseNode(node)
		if err != nil {
			return err
		}

		currentAsMap, ok := current.(map[string]interface{})
		if !ok {
			return ErrValueNotFound
		}

		// if this is last node and it has no indexes then set it to value and return
		if node_iter == len(nodes)-1 && len(indexes) == 0 {
			currentAsMap[name] = value
			return nil
		}
		// else set as current
		current = currentAsMap[name]

		for index_iter, index := range indexes {
			currentAsArray, ok := current.([]interface{})
			if !ok {
				return ErrValueNotFound
			}
			if len(currentAsArray) <= index {
				return ErrValueNotFound
			}

			// if this is the last node and last index of a node then set it to value and return
			if node_iter == len(nodes)-1 && index_iter == len(indexes)-1 {
				currentAsArray[index] = value
				return nil
			}
			// else set as current
			current = currentAsArray[index]
		}
	}

	return nil
}
