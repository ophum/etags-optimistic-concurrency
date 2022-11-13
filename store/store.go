package store

import (
	"errors"
	"sort"
	"sync"

	"github.com/ophum/etags-optimistic-concurrency/entities"
)

type Store struct {
	pets map[string]*entities.Pet
	mu   sync.Mutex
}

func New() *Store {
	return &Store{
		pets: map[string]*entities.Pet{},
		mu:   sync.Mutex{},
	}
}

var ErrNotFound = errors.New("not found")

func (s *Store) GetAll() ([]*entities.Pet, error) {
	r := make([]*entities.Pet, 0, len(s.pets))
	for _, v := range s.pets {
		r = append(r, v.DeepCopy())
	}

	// order by create_at asc
	sort.Slice(r, func(i, j int) bool {
		return r[i].CreatedAt.Before(r[j].CreatedAt)
	})
	return r, nil
}

func (s *Store) Get(id string) (*entities.Pet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.pets[id]; !exists {
		return nil, ErrNotFound
	}
	return s.pets[id].DeepCopy(), nil
}

func (s *Store) Put(pet *entities.Pet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pets[pet.ID] = pet.DeepCopy()
	return nil
}

func (s *Store) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.pets, id)
}
