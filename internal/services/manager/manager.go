package manager

import (
	"sync"

	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
)

type Manager struct {
	Tasks   chan *proto.Task
	Results chan *proto.Task
	mu      sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		Tasks:   make(chan *proto.Task, 100),
		Results: make(chan *proto.Task, 100),
	}
}

func (m *Manager) AddTask(task *proto.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Tasks <- task
}

func (m *Manager) AddResult(result *proto.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Results <- result
}

func (m *Manager) GetResult() *proto.Task {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := <-m.Results
	return result
}
