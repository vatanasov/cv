package populator

import (
	"autobiography/internal/models"
	"iter"
	"maps"
)

type Settable[K comparable, V any] interface {
	Key() K
	Value() V
}

type Set[K comparable, V any] struct {
	values map[K]V
}

func NewSet[T comparable, K any]() Set[T, K] {
	return Set[T, K]{values: map[T]K{}}
}

func (s *Set[T, K]) Add(value Settable[T, K]) {
	s.values[value.Key()] = value.Value()
}

func (s *Set[T, K]) Remove(value Settable[T, K]) {
	delete(s.values, value.Key())
}

func (s *Set[T, K]) GetAll() iter.Seq[K] {
	return maps.Values(s.values)
}

// In the function below the external module was not allowed
type internalTechnology = models.Technology

func insertTechnologies(models models.Models, technologies Set[string, internalTechnology]) error {
	for technology := range technologies.GetAll() {
		err := models.Technologies.Insert(technology)
		if err != nil {
			return err
		}
	}

	return nil
}
