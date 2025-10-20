package model

import (
	"errors"
	"time"
)

var ErrPartNotFound = errors.New("part not found")

type Part struct {
	ID            string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      PartCategory
	Dimensions    *PartDimensions
	Manufacturer  *PartManufacturer
	Tags          []string
	Metadata      map[string]PartMetadataValue
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PartCategory struct {
	number int32
	name   string
}

func NewPartCategory(number int32, name string) PartCategory {
	return PartCategory{number: number, name: name}
}

func (c PartCategory) Number() int32 {
	return c.number
}

func (c PartCategory) Name() string {
	return c.name
}

type PartDimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type PartManufacturer struct {
	Name    string
	Country string
	Website string
}

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

type PartFilter struct {
	IDs                   []string
	Names                 []string
	Categories            []PartCategory
	ManufacturerCountries []string
	Tags                  []string
}
