package manager

import (
	"context"
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
	Results     chan *proto.Task
	DB          *sqlite.Storage
	mu          sync.Mutex
	Cfg         *config.Config
}

type Expression struct {
	Stack       []string
	Tasks       map[string]*proto.Task
	TotalTasks  int
	Completed   int
	FinalResult float64
	AllDone     bool
	Err         bool
	// mu          sync.Mutex
}

func NewManager(cfg *config.Config) *Manager {
	s, err := sqlite.New(cfg.Database.DBPath, cfg.Database.MigrationsPath)
	if err != nil {
		panic(err)
	}
	return &Manager{
		DB:          s,
		Queque:      make(chan *proto.Task, 100),
		Results:     make(chan *proto.Task, 100),
		Expressions: make(map[string]*Expression),
		Cfg:         cfg,
	}
}

func NewExpression() *Expression {
	return &Expression{
		Tasks:      make(map[string]*proto.Task),
		Stack:      []string{},
		TotalTasks: 0,
		Completed:  0,
		AllDone:    false,
		Err:        false,
	}
}

func (m *Manager) AddTask(task *proto.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, have := m.Expressions[task.ExpressionID]
	if !have {
		m.Expressions[task.ExpressionID] = NewExpression()
	}
	expr := m.Expressions[task.ExpressionID]
	expr.Tasks[task.ID] = task
	expr.TotalTasks++
	m.Queque <- task
}

func (m *Manager) AddResult(result *proto.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Expressions[result.ExpressionID].Tasks[result.ID] = result
	m.Results <- result
}

func (m *Manager) AddError(result *proto.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Expressions[result.ExpressionID].Err = true
	err := m.DB.AddError(context.TODO(), result.ExpressionID)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *Manager) AddStack(expID string, stack []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Expressions[expID].Stack = stack
}

// NOTE: нужна или нет вот в чем вопрос
func (m *Manager) GetResult() *proto.Task {
	m.mu.Lock()
	defer m.mu.Unlock()
	// TODO:
	return nil
}

func (m *Manager) GenerateUUID() string {
	return fmt.Sprint(uuid.New())
}
