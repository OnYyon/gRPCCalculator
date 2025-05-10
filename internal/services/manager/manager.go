package manager

import (
	"fmt"
	"sync"

	"github.com/OnYyon/gRPCCalculator/internal/config"
	"github.com/OnYyon/gRPCCalculator/internal/storage/sqlite"
	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"github.com/google/uuid"
)

type Manager struct {
	Expressions map[string]*Expression
	Queque      chan *proto.Task
	DB          *sqlite.Storage
	mu          sync.Mutex
}

type Expression struct {
	Stack       []string
	Tasks       map[string]*proto.Task
	FinalResult float64
	AllDone     bool
}

func NewManager(cfg *config.Config) *Manager {
	s, err := sqlite.New(cfg.Database.DBPath, cfg.Database.MigrationsPath)
	if err != nil {
		panic(err)
	}
	return &Manager{
		DB:          s,
		Queque:      make(chan *proto.Task, 100),
		Expressions: make(map[string]*Expression),
	}
}

func NewExpression() *Expression {
	return &Expression{
		Tasks:   make(map[string]*proto.Task),
		Stack:   []string{},
		AllDone: false,
	}
}

func (m *Manager) AddTask(task *proto.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Queque <- task
}

func (m *Manager) AddResult(result *proto.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, have := m.Expressions[result.ExpressionID]
	if !have {
		m.Expressions[result.ExpressionID] = NewExpression()
	}
	m.Expressions[result.ExpressionID].Tasks[result.ID] = result
}

func (m *Manager) GetResult() *proto.Task {
	m.mu.Lock()
	defer m.mu.Unlock()
	// TODO:
	return nil
}

func (m *Manager) GenerateUUID() string {
	return fmt.Sprint(uuid.New())
}
