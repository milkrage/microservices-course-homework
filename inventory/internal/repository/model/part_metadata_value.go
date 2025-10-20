package model

import "errors"

type PartMetadataValue struct {
	value any
}

func NewPartMetadataValue(value any) (PartMetadataValue, error) {
	switch value.(type) {
	case string, int64, float64, bool:
		return PartMetadataValue{value: value}, nil
	default:
		return PartMetadataValue{}, errors.New("invalid value type: must be one of string, int64, float64, or bool")
	}
}

func (v PartMetadataValue) Get() any {
	return v.value
}
