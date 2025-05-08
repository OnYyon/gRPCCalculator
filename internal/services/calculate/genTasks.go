package services

import (
	"fmt"
	"strconv"

	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
)

func GenerateTasks(rpn []string, expID string, manager *manager.Manager) ([]string, []*proto.Task, error) {
	// Returns remaining operations and current tasks for execute
	stack := []string{}
	tasks := []*proto.Task{}
	// TODO: cделать релизацую с 0
	for _, v := range rpn {
		if isOperator(v) {
			fmt.Println(stack)
			if len(stack) < 2 {
				stack = append(stack, v)
				continue
			}
			if isNumber(stack[len(stack)-1]) && isNumber(stack[len(stack)-2]) {
				a, err := convertToFloat(stack[len(stack)-2])
				if err != nil {
					return nil, nil, err
				}
				b, err := convertToFloat(stack[len(stack)-1])
				if err != nil {
					return nil, nil, err
				}
				task := &proto.Task{
					ID:           manager.GenerateUUID(),
					Arg1:         a,
					Arg2:         b,
					Operator:     v,
					ExpressionID: expID,
				}
				fmt.Println(task)
				tasks = append(tasks, task)
				stack = stack[:len(stack)-2]
				stack = append(stack, task.ID)
			} else {
				stack = append(stack, v)
			}
		} else if isNumber(v) {
			stack = append(stack, v)
		} else {

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

func convertToFloat(n string) (float64, error) {
	a, err := strconv.ParseFloat(n, 64)
	if err != nil {
		return 0, err
	}
	return a, nil
}
