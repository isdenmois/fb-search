package mocks

import (
	"sync"

	"fb-search/domain"
)

type MockInpParser struct {
	RebuildCalled bool
	Progress      *domain.ParseProgress
	Mu            sync.Mutex
}

func NewMockInpParser() *MockInpParser {
	return &MockInpParser{
		Progress: &domain.ParseProgress{},
	}
}

func (m *MockInpParser) RebuildDb(progress *domain.ParseProgress) {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	m.RebuildCalled = true
	progress.Files = 3
	progress.Books = 100
	progress.Time = 5000
}

func (m *MockInpParser) WasRebuildCalled() bool {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	return m.RebuildCalled
}
