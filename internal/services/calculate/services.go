package services

import proto "github.com/OnYyon/gRPCCalculator/proto/gen"

func ProcessTask(task *proto.Task) *proto.Task {
	task.ID = "modifed"
	return task
}
