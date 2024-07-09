package core

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

var ErrIncorectType = errors.New("incorrect types")

//nolint:ireturn //this is why generic type
func YamlToObject[T any](yamlData any, unmarshalObject T) (T, error) {
	var results []byte

	switch p := yamlData.(type) {
	case map[interface{}]interface{}:
		var err error
		if results, err = yaml.Marshal(p); err != nil {
			return unmarshalObject, fmt.Errorf("failed to marshal yamlData: %w", err)
		}
	default:
		return unmarshalObject, fmt.Errorf("%w for : %T", ErrIncorectType, p)
	}

	err := yaml.Unmarshal(results, &unmarshalObject)
	if err != nil {
		return unmarshalObject, fmt.Errorf("failed to unmarshal output object: %w", err)
	}

	return unmarshalObject, nil
}
