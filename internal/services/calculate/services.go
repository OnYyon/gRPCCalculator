package services

import (
	"fmt"

	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
)

func ProcessTask(task *proto.Task) (*proto.Task, error) {
	switch task.Operator {
	case "+":
		task.Result = task.Arg1 + task.Arg2
	case "-":
		task.Result = task.Arg1 - task.Arg2
	case "*":
		task.Result = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0.0 {
			return task, fmt.Errorf("zero division error")
		}
		task.Result = task.Arg1 / task.Arg2
	}
	return task, nil
}
