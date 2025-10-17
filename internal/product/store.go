package product

import (
	"sync"

	"github.com/google/uuid"
)

type Store struct {
	mu   sync.RWMutex
	data []Product
}

func NewStore() *Store { return &Store{data: make([]Product, 0)} }

func (s *Store) List() []Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Product, len(s.data))
	copy(out, s.data)
	return out
}

func (s *Store) Create(name string, price float64) Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := Product{ID: uuid.NewString(), Name: name, Price: price}
	s.data = append(s.data, p)
	return p
}
