package services

import (
	"context"
	"fmt"

	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
)

func StartResultProcessor(m *manager.Manager) {
	go func() {
		for result := range m.Results {
			expr, exists := m.Expressions[result.ExpressionID]

			if !exists {
				continue

			}

			if expr.Tasks[result.ID].Err {
				m.AddError(result)
				continue
			}
			expr.Completed++
			if expr.Completed == expr.TotalTasks {
				if len(expr.Stack) == 1 {
					expr.FinalResult = result.Result
					m.DB.UpdateExpression(context.TODO(), result.ExpressionID, result.Result)
				} else {
					stack, tasks, err := GenerateTasks(expr.Stack, result.ExpressionID, m)
					if err != nil {
						fmt.Println(err)
					}
					m.AddStack(result.ExpressionID, stack)
					for _, task := range tasks {
						m.AddTask(task)
					}
				}
			}
		}
	}()
}
