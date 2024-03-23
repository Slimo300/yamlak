package yaml

import (
	"fmt"
	"strings"
)

func DeleteValueByQuery(doc interface{}, query string) error {

	nodes := strings.Split(query, ".")
	current := doc

	var prev interface{}

	for node_iter, node := range nodes {
		name, indexes, err := ParseNode(node)
		if err != nil {
			return err
		}

		currentAsMap, ok := current.(map[string]interface{})
		if !ok {
			return fmt.Errorf("value not found")
		}

		// if this is last node and it has no indexes then set it to value and return
		if node_iter == len(nodes)-1 && len(indexes) == 0 {
			delete(currentAsMap, node)
			return nil
		}
		// else set as current
		prev = current
		current = currentAsMap[name]

		for index_iter, index := range indexes {
			currentAsArray, ok := current.([]interface{})
			if !ok {
				return fmt.Errorf("value not found")
			}
			if len(currentAsArray) <= index {
				return fmt.Errorf("value not found")
			}

			// if this is the last node and last index of a node then delete value of index from array and reassign it to prev
			if node_iter == len(nodes)-1 && index_iter == len(indexes)-1 {

				currentAsArray = append(currentAsArray[:index], currentAsArray[index+1:]...)
				// index_iter being 0 means prev is a map, if its not then its a multi-level slice
				if index_iter == 0 {
					prevAsMap, ok := prev.(map[string]interface{})
					if !ok {
						return fmt.Errorf("internal error")
					}
					prevAsMap[name] = currentAsArray
				} else {
					prevAsArray, ok := prev.([]interface{})
					if !ok {
						return fmt.Errorf("internal error")
					}
					prevAsArray[indexes[index_iter-1]] = currentAsArray
				}
				return nil
			}
			// else set as current
			prev = current
			current = currentAsArray[index]
		}
	}

	return nil
}
