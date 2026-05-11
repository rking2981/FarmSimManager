package server

import (
	"sync"

	"github.com/farmsimcompanymanager/companion/internal/parser"
)

type Store struct {
	mu        sync.RWMutex
	companies []parser.Company
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) SetCompanies(companies []parser.Company) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.companies = companies
}

func (s *Store) GetCompanies() []parser.Company {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.companies
}
