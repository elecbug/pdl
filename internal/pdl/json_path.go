package pdl

import (
	"fmt"
	"strconv"
	"strings"
)

type pathPart struct {
	Key     string
	Index   *int
	IsArray bool
}

func parseJSONPath(path string) ([]pathPart, error) {
	if path == "" {
		return nil, fmt.Errorf("empty json path")
	}

	rawParts := strings.Split(path, ".")
	parts := make([]pathPart, 0, len(rawParts))

	for _, raw := range rawParts {
		if raw == "" {
			return nil, fmt.Errorf("invalid empty path part in %q", path)
		}

		part, err := parsePathPart(raw)
		if err != nil {
			return nil, err
		}

		parts = append(parts, part)
	}

	return parts, nil
}

func parsePathPart(raw string) (pathPart, error) {
	// [0]
	if strings.HasPrefix(raw, "[") {
		if !strings.HasSuffix(raw, "]") {
			return pathPart{}, fmt.Errorf("invalid array path part %q", raw)
		}

		idxStr := strings.TrimSuffix(strings.TrimPrefix(raw, "["), "]")
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			return pathPart{}, fmt.Errorf("invalid array index %q", idxStr)
		}

		return pathPart{
			IsArray: true,
			Index:   &idx,
		}, nil
	}

	// key[0]
	if i := strings.Index(raw, "["); i >= 0 {
		if !strings.HasSuffix(raw, "]") {
			return pathPart{}, fmt.Errorf("invalid indexed path part %q", raw)
		}

		key := raw[:i]
		idxStr := raw[i+1 : len(raw)-1]

		if key == "" {
			return pathPart{}, fmt.Errorf("empty key in %q", raw)
		}

		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			return pathPart{}, fmt.Errorf("invalid array index %q", idxStr)
		}

		return pathPart{
			Key:     key,
			IsArray: true,
			Index:   &idx,
		}, nil
	}

	return pathPart{Key: raw}, nil
}

func setJSONPath(root *map[string]any, path string, value any) error {
	parts, err := parseJSONPath(path)
	if err != nil {
		return err
	}

	// Special case: root array path like [0]
	if parts[0].IsArray && parts[0].Key == "" {
		return fmt.Errorf("root array path is not supported yet: %q", path)
	}

	var cur any = *root

	for i, part := range parts {
		last := i == len(parts)-1

		obj, ok := cur.(map[string]any)
		if !ok {
			return fmt.Errorf("path conflict at %q", path)
		}

		if part.IsArray {
			arrAny, exists := obj[part.Key]
			if !exists {
				arrAny = []any{}
			}

			arr, ok := arrAny.([]any)
			if !ok {
				return fmt.Errorf("path conflict: %q is not array", part.Key)
			}

			arr = ensureArrayLen(arr, *part.Index+1)

			if last {
				arr[*part.Index] = value
				obj[part.Key] = arr
				return nil
			}

			if arr[*part.Index] == nil {
				arr[*part.Index] = map[string]any{}
			}

			nextObj, ok := arr[*part.Index].(map[string]any)
			if !ok {
				return fmt.Errorf("path conflict at array %q[%d]", part.Key, *part.Index)
			}

			obj[part.Key] = arr
			cur = nextObj
			continue
		}

		if last {
			if _, exists := obj[part.Key]; exists {
				return fmt.Errorf("duplicate json path %q", path)
			}
			obj[part.Key] = value
			return nil
		}

		next, exists := obj[part.Key]
		if !exists {
			next = map[string]any{}
			obj[part.Key] = next
		}

		nextObj, ok := next.(map[string]any)
		if !ok {
			return fmt.Errorf("path conflict: %q is not object", part.Key)
		}

		cur = nextObj
	}

	return nil
}

func ensureArrayLen(arr []any, n int) []any {
	for len(arr) < n {
		arr = append(arr, nil)
	}
	return arr
}
