package part

import (
	"sync"

	"github.com/milkrage/microservices-course-homework/inventory/internal/repository/model"
)

type Repository struct {
	storage map[string]model.Part
	mu      sync.RWMutex
}

func NewPartRepository() *Repository {
	repo := &Repository{
		storage: make(map[string]model.Part),
	}

	repo.init()

	return repo
}
