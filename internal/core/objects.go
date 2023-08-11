package core

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

var ErrIncorectType = errors.New("incorrect types")

//nolint
func GetObjectFromMap[T any](objectMap any, unmarshalObject T) (T, error) {
	var objectBytes []byte

	var err error

	switch p := objectMap.(type) {
	case map[interface{}]interface{}:
		objectBytes, err = yaml.Marshal(p)
	default:
		return unmarshalObject, fmt.Errorf("%w for : %T", ErrIncorectType, p)
	}

	if err != nil {
		return unmarshalObject, fmt.Errorf("failed to marshal objectMap: %w", err)
	}

	err = yaml.Unmarshal(objectBytes, &unmarshalObject)
	if err != nil {
		return unmarshalObject, fmt.Errorf("failed to unmarshal output object: %w", err)
	}

	return unmarshalObject, nil
}
