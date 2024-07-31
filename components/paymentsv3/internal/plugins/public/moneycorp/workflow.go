package moneycorp

import "github.com/formancehq/paymentsv3/internal/plugins/models"

func workflow() models.Workflow {
	return []models.TaskTree{
		{
			TaskType: models.TASK_FETCH_ACCOUNTS,
			Name:     "fetch_accounts",
			NextTasks: []models.TaskTree{
				{
					TaskType:  models.TASK_FETCH_PAYMENTS,
					Name:      "fetch_payments",
					NextTasks: []models.TaskTree{},
				},
			},
		},
	}
}
