package model

import "github.com/go-faster/errors"

var (
	ErrInternal         = errors.New("internal server error")
	ErrNotFoundAllParts = errors.New("not found all parts")
)

type Part struct {
	ID    string
	Price float64
}

type ListParts struct {
	Parts []Part
}

type PartFilter struct {
	IDs []string
}
