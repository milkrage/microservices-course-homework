package model

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
