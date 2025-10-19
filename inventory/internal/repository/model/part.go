package model

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
}

type PartFilter struct {
	IDs                   []string
	Names                 []string
	Categories            []PartCategory
	ManufacturerCountries []string
	Tags                  []string
}
