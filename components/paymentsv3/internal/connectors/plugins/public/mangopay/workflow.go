package mangopay

import "github.com/formancehq/paymentsv3/internal/models"

func workflow() models.Tasks {
	return []models.TaskTree{
		{
			TaskType: models.TASK_FETCH_OTHERS,
			Name:     usersTaskName,
			NextTasks: []models.TaskTree{
				{
					TaskType:  models.TASK_FETCH_EXTERNAL_ACCOUNTS,
					Name:      "fetch_external_accounts",
					NextTasks: []models.TaskTree{},
				},
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
			},
		},
	}
}
