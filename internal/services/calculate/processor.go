package services

import (
	"fmt"

	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
)

func StartResultProcessor(m *manager.Manager) {
	fmt.Println("start processor")
	go func() {
		for result := range m.Results {
			fmt.Println(result)
			expr, exists := m.Expressions[result.ExpressionID]

			if !exists {
				continue

			}

			expr.Completed++
			if expr.Completed == expr.TotalTasks {
				if len(expr.Stack) == 1 {
					expr.FinalResult = result.Result
					fmt.Println(expr.FinalResult)
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
