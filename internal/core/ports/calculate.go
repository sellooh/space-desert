package ports

import "github.com/sellooh/space-desert/internal/core/domain"

type ResourceGenerator interface {
	Generate() <-chan domain.Resource
}

type CalculateService interface {
	Calculate(ResourceGenerator) uint32
}
