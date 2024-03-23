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

// ParseNode parses node reference and extracts node name and possibly slice indexes
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

	// Parsing indexes
	for _, match := range indexesPattern.FindAllStringSubmatch(indexesStr, -1) {
		index, err := strconv.Atoi(match[1])
		if err != nil {
			return "", nil, fmt.Errorf("parsing index error: %w", err)
		}
		indexes = append(indexes, index)
	}

	return name, indexes, nil
}
