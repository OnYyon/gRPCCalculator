package tests

import (
	"strconv"
	"testing"

	"github.com/OnYyon/gRPCCalculator/internal/config"
	services "github.com/OnYyon/gRPCCalculator/internal/services/calculate"
	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
	"github.com/stretchr/testify/assert"
)

type MockManager struct {
	Cfg         *config.Config
	Expressions map[string]*proto.Expression
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

func getTimeout(operator string, cfg *config.Config) int64 {
	switch operator {
	case "+":
		return cfg.Server.TimeAdditionMS
	case "-":
		return cfg.Server.TimeSubtractionMS
	case "*":
		return cfg.Server.TimeMultiplicationMS
	case "/":
		return cfg.Server.TimeDivisionMS
	}
	return 0
}

func (m *MockManager) GenerateUUID() string {
	return "mock-uuid"
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"-123.45", true},
		{"abc", false},
		{"", false},
		{"12a3", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, isNumber(tt.input))
		})
	}
}

func TestIsOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"+", true},
		{"-", true},
		{"*", true},
		{"/", true},
		{"%", false},
		{"abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, isOperator(tt.input))
		})
	}
}

func TestGetTimeout(t *testing.T) {
	cfg, err := config.Load("./internal/config/config.yaml")
	if err != nil {
		return
	}

	tests := []struct {
		operator string
		expected int64
	}{
		{"+", 100},
		{"-", 100},
		{"*", 300},
		{"/", 400},
		{"%", 0},
	}

	for _, tt := range tests {
		t.Run(tt.operator, func(t *testing.T) {
			assert.Equal(t, tt.expected, getTimeout(tt.operator, cfg))
		})
	}
}

func TestGenerateTasks(t *testing.T) {
	mockCfg, err := config.Load("./internal/config/config.yaml")
	if err != nil {
		return
	}

	mockManager := &manager.Manager{
		Cfg:    mockCfg,
		Queque: make(chan *proto.Task),
		Expressions: map[string]*manager.Expression{
			"task1": {
				Tasks: map[string]*proto.Task{
					"t1": {Result: 42},
				},
			},
		},
	}

	tests := []struct {
		name          string
		rpn           []string
		expID         string
		expectedStack []string
		expectedTasks []*proto.Task
		expectError   bool
	}{
		{
			name:          "Simple addition",
			rpn:           []string{"2", "3", "+"},
			expID:         "exp1",
			expectedStack: []string{},
			expectedTasks: []*proto.Task{
				{
					ID:           "mock-uuid",
					Arg1:         2,
					Arg2:         3,
					ExpressionID: "exp1",
					Operator:     "+",
					Timeout:      100,
					Err:          false,
				},
			},
			expectError: false,
		},
		{
			name:          "Multiple operations",
			rpn:           []string{"2", "3", "+", "5", "*"},
			expID:         "exp1",
			expectedStack: []string{"mock-uuid", "5", "*"},
			expectedTasks: []*proto.Task{
				{
					ID:           "mock-uuid",
					Arg1:         2,
					Arg2:         3,
					ExpressionID: "exp1",
					Operator:     "+",
					Timeout:      100,
					Err:          false,
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, tasks, err := services.GenerateTasks(tt.rpn, tt.expID, mockManager)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if len(tt.expectedTasks) > 0 {
				assert.Equal(t, len(tt.expectedTasks), len(tasks))
				for i, expectedTask := range tt.expectedTasks {
					assert.Equal(t, expectedTask.Arg1, tasks[i].Arg1)
					assert.Equal(t, expectedTask.Arg2, tasks[i].Arg2)
					assert.Equal(t, expectedTask.ExpressionID, tasks[i].ExpressionID)
					assert.Equal(t, expectedTask.Operator, tasks[i].Operator)
					assert.Equal(t, expectedTask.Timeout, tasks[i].Timeout)
					assert.Equal(t, expectedTask.Err, tasks[i].Err)
				}
			} else {
				assert.Empty(t, tasks)
			}
		})
	}
}
