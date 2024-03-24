package yaml

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	indexesPattern = regexp.MustCompile(`\[(\d+)\]`)
	nodeRefPattern = regexp.MustCompile(`^([a-zA-Z]+)((?:\[\d+\])*)$`)
)

// ParseNode is a function that parses a single node in dot-separated node path.
// E.G.: spec.template.spec.containers[0].image (NODE_PATH) has a node containers[0] in it
// ParseNode returns name of node (which is always present) and a slice of integers which are indexes of
// slice referencing
func ParseNode(s string) (string, []int, error) {
	matches := nodeRefPattern.FindStringSubmatch(s)
	if len(matches) < 2 {
		return "", nil, fmt.Errorf("invalid node reference")
	}

	// getting table name
	name := matches[1]

	// getting indexes
	indexesStr := matches[2]

	var indexes []int

	// When we used nodeRefPattern we got division that looks like: "containers[0][1]"" -> "containers", "[0][1]"
	// Here is where parsing of "[0][1]" happens to get all indexes
	for _, match := range indexesPattern.FindAllStringSubmatch(indexesStr, -1) {
		index, err := strconv.Atoi(match[1])
		if err != nil {
			return "", nil, fmt.Errorf("parsing index error: %w", err)
		}
		indexes = append(indexes, index)
	}

	return name, indexes, nil
}
