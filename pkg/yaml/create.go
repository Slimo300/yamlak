package yaml

import "strings"

func CreateValueByQuery(doc interface{}, query, value string) error {
	nodes := strings.Split(query, ".")
	current := doc // current is a value to use for extracting deeper levels of yaml body

	var prev interface{}

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
		prev = current
		current, ok = currentAsMap[name]
		if !ok {
			if len(indexes) == 0 {
				currentAsMap[name] = map[string]interface{}{}
			} else {
				currentAsMap[name] = []interface{}{}
			}
			current = currentAsMap[name]
		}

		for index_iter, index := range indexes {
			currentAsSlice, ok := current.([]interface{})
			if !ok {
				if index_iter == 0 {
					prevAsMap, ok := prev.(map[string]interface{})
					if !ok {
						return ErrValueNotFound
					}
					prevAsMap[name] = []interface{}{}
				} else {
					prevAsSlice, ok := prev.([]interface{})
					if !ok {
						return ErrValueNotFound
					}
					prevAsSlice[indexes[index_iter-1]] = []interface{}{}
				}
			}

			// else set as current
			if index > len(currentAsSlice) {
				return ErrValueNotFound
			}
			if index == len(currentAsSlice) {
				if index_iter == len(indexes)-1 {
					currentAsSlice = append(currentAsSlice, map[string]interface{}{})
				} else {
					currentAsSlice = append(currentAsSlice, []interface{}{})
				}

				// this lines are obligatory because when we assign result of append to currentAsSlice it doesn't necessarly has to be
				// the same slice in the same address, so doc value will not be changed further without a reassignment
				if index_iter == 0 {
					prevAsMap, ok := prev.(map[string]interface{})
					if !ok {
						return ErrValueNotFound
					}
					prevAsMap[name] = currentAsSlice
				} else {
					prevAsSlice, ok := prev.([]interface{})
					if !ok {
						return ErrValueNotFound
					}
					prevAsSlice[indexes[index_iter-1]] = currentAsSlice
				}
			}

			// if this is the last node and last index of a node then set it to value and return
			if node_iter == len(nodes)-1 && index_iter == len(indexes)-1 {
				currentAsSlice[index] = value
				return nil
			}

			prev = current
			current = currentAsSlice[index]
		}
	}

	return nil
}
