package services

import (
	"context"
	"fmt"

	proto "github.com/OnYyon/gRPCCalculator/proto/gen"
)

func ProcessTask(ctx context.Context, task *proto.Task) (*proto.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
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
