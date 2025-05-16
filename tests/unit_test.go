package tests

import (
	"testing"
	"time"

	"github.com/OnYyon/gRPCCalculator/internal/config"
	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestConfig() *config.Config {
	return &config.Config{
		Database: config.DatabaseConfig{
			DBPath:         ":memory:",
			MigrationsPath: "./migrations",
		},
	}
}

func TestNewManager(t *testing.T) {
	cfg := createTestConfig()
	mgr := manager.NewManager(cfg)

	assert.NotNil(t, mgr)
	assert.NotNil(t, mgr.DB)
	assert.NotNil(t, mgr.Queque)
	assert.NotNil(t, mgr.Results)
	assert.NotNil(t, mgr.Expressions)
	assert.Equal(t, cfg, mgr.Cfg)
}

func TestNewExpression(t *testing.T) {
	expr := manager.NewExpression()

	assert.NotNil(t, expr)
	assert.Empty(t, expr.Stack)
	assert.Empty(t, expr.Tasks)
	assert.Equal(t, 0, expr.TotalTasks)
	assert.Equal(t, 0, expr.Completed)
	assert.False(t, expr.AllDone)
}

func TestManager_AddTask(t *testing.T) {
	cfg := createTestConfig()
	mgr := manager.NewManager(cfg)

	task := &proto.Task{
		ID:           "task1",
		ExpressionID: "expr1",
	}

	mgr.AddTask(task)

	expr, exists := mgr.Expressions["expr1"]
	require.True(t, exists)
	require.NotNil(t, expr)

	_, taskExists := expr.Tasks["task1"]
	assert.True(t, taskExists)
	assert.Equal(t, 1, expr.TotalTasks)

	select {
	case queuedTask := <-mgr.Queque:
		assert.Equal(t, task, queuedTask)
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "Task was not queued")
	}

	task2 := &proto.Task{
		ID:           "task2",
		ExpressionID: "expr1",
	}
	mgr.AddTask(task2)

	assert.Equal(t, 2, expr.TotalTasks)
}

func TestManager_AddResult(t *testing.T) {
	cfg := createTestConfig()
	mgr := manager.NewManager(cfg)

	task := &proto.Task{
		ID:           "task1",
		ExpressionID: "expr1",
	}
	mgr.AddTask(task)

	result := &proto.Task{
		ID:           "task1",
		ExpressionID: "expr1",
		Result:       42,
	}

	mgr.AddResult(result)

	expr := mgr.Expressions["expr1"]
	updatedTask := expr.Tasks["task1"]
	assert.Equal(t, float64(42), updatedTask.Result)

	select {
	case res := <-mgr.Results:
		assert.Equal(t, result, res)
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "Result was not sent to channel")
	}
}

func TestManager_AddStack(t *testing.T) {
	cfg := createTestConfig()
	mgr := manager.NewManager(cfg)

	task := &proto.Task{
		ID:           "task1",
		ExpressionID: "expr1",
	}
	mgr.AddTask(task)

	newStack := []string{"1", "2", "+"}
	mgr.AddStack("expr1", newStack)

	expr := mgr.Expressions["expr1"]
	assert.Equal(t, newStack, expr.Stack)
}

func TestManager_GenerateUUID(t *testing.T) {
	cfg := createTestConfig()
	mgr := manager.NewManager(cfg)

	uuid1 := mgr.GenerateUUID()
	uuid2 := mgr.GenerateUUID()

	assert.NotEmpty(t, uuid1)
	assert.NotEmpty(t, uuid2)
	assert.NotEqual(t, uuid1, uuid2)

	_, err := uuid.Parse(uuid1)
	assert.NoError(t, err)
}
