package tree

import (
	"encoding/json"
	"sort"
	"strconv"
)

func Flatten(data []any) []Node {
	if len(data) == 0 {
		return nil
	}

	nodes := make([]Node, 0, len(data)*4)

	if len(data) == 1 {
		flattenValue(&nodes, data[0], 0, ".", "", true)
		return nodes
	}

	nodes = append(nodes, Node{
		Type:        ArrayOpen,
		Depth:       0,
		Path:        ".",
		Collapsible: true,
		ChildCount:  len(data),
		IsLast:      true,
	})

	for i, value := range data {
		flattenValue(&nodes, value, 1, appendArrayPath(".", i), "", i == len(data)-1)
	}

	nodes = append(nodes, Node{
		Type:   ArrayClose,
		Depth:  0,
		Path:   ".",
		IsLast: true,
	})

	return nodes
}

func flattenValue(nodes *[]Node, value any, depth int, path, key string, isLast bool) {
	switch value := value.(type) {
	case map[string]any:
		keys := mapKeys(value)
		*nodes = append(*nodes, Node{
			Type:        ObjectOpen,
			Depth:       depth,
			Key:         key,
			Path:        path,
			Collapsible: true,
			ChildCount:  len(keys),
			IsLast:      isLast,
		})

		for i, objectKey := range keys {
			flattenValue(nodes, value[objectKey], depth+1, appendObjectPath(path, objectKey), objectKey, i == len(keys)-1)
		}

		*nodes = append(*nodes, Node{
			Type:   ObjectClose,
			Depth:  depth,
			Path:   path,
			IsLast: isLast,
		})

	case []any:
		*nodes = append(*nodes, Node{
			Type:        ArrayOpen,
			Depth:       depth,
			Key:         key,
			Path:        path,
			Collapsible: true,
			ChildCount:  len(value),
			IsLast:      isLast,
		})

		for i, item := range value {
			flattenValue(nodes, item, depth+1, appendArrayPath(path, i), "", i == len(value)-1)
		}

		*nodes = append(*nodes, Node{
			Type:   ArrayClose,
			Depth:  depth,
			Path:   path,
			IsLast: isLast,
		})

	default:
		nodeType := ArrayElement
		if key != "" {
			nodeType = KeyValue
		}

		*nodes = append(*nodes, Node{
			Type:   nodeType,
			Depth:  depth,
			Key:    key,
			Value:  value,
			Path:   path,
			IsLast: isLast,
		})
	}
}

func mapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func appendObjectPath(path, key string) string {
	if isIdentifier(key) {
		if path == "." {
			return "." + key
		}
		return path + "." + key
	}

	return path + "[" + jsonString(key) + "]"
}

func appendArrayPath(path string, index int) string {
	return path + "[" + strconv.Itoa(index) + "]"
}

func jsonString(value string) string {
	b, err := json.Marshal(value)
	if err != nil {
		return strconv.Quote(value)
	}

	return string(b)
}

func isIdentifier(value string) bool {
	if value == "" {
		return false
	}

	for i, r := range value {
		isLetter := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
		if i == 0 {
			if r != '_' && !isLetter {
				return false
			}
			continue
		}

		isDigit := r >= '0' && r <= '9'
		if r != '_' && !isLetter && !isDigit {
			return false
		}
	}

	return true
}
