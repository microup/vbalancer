package core

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

var ErrIncorectType = errors.New("incorrect types")

//nolint
func YamlToObject[T any](yamlData any, unmarshalObject T) (T, error) {
	var objectBytes []byte

	var err error

	switch p := yamlData.(type) {
	case map[interface{}]interface{}:
		objectBytes, err = yaml.Marshal(p)
		if err != nil {
			return unmarshalObject, fmt.Errorf("failed to marshal yamlData: %w", err)
		}
	default:
		return unmarshalObject, fmt.Errorf("%w for : %T", ErrIncorectType, p)
	}

	err = yaml.Unmarshal(objectBytes, &unmarshalObject)
	if err != nil {
		return unmarshalObject, fmt.Errorf("failed to unmarshal output object: %w", err)
	}

	return unmarshalObject, nil
}
