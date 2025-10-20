package part

import "github.com/milkrage/microservices-course-homework/inventory/internal/repository"

type Service struct {
	partRepository repository.PartRepository
}

func NewPartService(partRepository repository.PartRepository) *Service {
	return &Service{
		partRepository: partRepository,
	}
}
