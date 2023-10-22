package ports

import "github.com/sellooh/space-desert/internal/core/domain"

type ResourceGenerator interface {
	Generate(chan<- error) <-chan domain.Resource
}
