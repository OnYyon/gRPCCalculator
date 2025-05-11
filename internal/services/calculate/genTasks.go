package services

import (
	"fmt"
	"strconv"

	"github.com/OnYyon/gRPCCalculator/internal/config"
	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
)

func GenerateTasks(
	rpn []string,
	expID string,
	manager *manager.Manager,
) ([]string, []*proto.Task, error) {
	// Returns remaining operations and current tasks for execute
	stack := []string{}
	tasks := []*proto.Task{}
	// TODO: cделать релизацую с 0
	for _, v := range rpn {
		// ----- NOTE: for tests -----
		// fmt.Println(stack)
		// ---------------------------
		if isNumber(v) {
			stack = append(stack, v)
		} else if isOperator(v) {
			if len(stack) < 2 {
				stack = append(stack, v)
			} else if isNumber(stack[len(stack)-1]) && isNumber(stack[len(stack)-2]) {
				operand2, err := strconv.ParseFloat(stack[len(stack)-1], 64)
				if err != nil {
					return nil, nil, err
				}
				operand1, err := strconv.ParseFloat(stack[len(stack)-2], 64)
				if err != nil {
					return nil, nil, err
				}

				task := &proto.Task{
					ID:           manager.GenerateUUID(),
					Arg1:         operand1,
					Arg2:         operand2,
					ExpressionID: expID,
					Operator:     v,
					Timeout:      getTimeout(v, manager.Cfg),
					Err:          false,
				}
				stack = stack[:len(stack)-2]

				tasks = append(tasks, task)
				stack = append(stack, task.ID)
			} else {
				stack = append(stack, v)
			}
		} else {
			stack = append(stack, fmt.Sprint(manager.Expressions[expID].Tasks[v].Result))
		}
	}
	return stack, tasks, nil
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
