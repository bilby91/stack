package workflow

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

const (
	// TODO(polo): add config
	defaultPolling = 2 * time.Minute
)

func (w Workflow) runNextWorkflow(ctx workflow.Context, fromPayload json.RawMessage, pageSize int, taskTree []*models.TaskTree) error {
	var nextWorkflow interface{}

	for _, task := range taskTree {
		switch task.TaskType {
		case models.TASK_FETCH_ACCOUNTS:
			req := FetchNextAccounts{
				FromPayload: fromPayload,
				PageSize:    pageSize,
			}

			nextWorkflow = req.GetWorkflow()
		case models.TASK_FETCH_EXTERNAL_ACCOUNTS:
			req := FetchNextExternalAccounts{
				FromPayload: fromPayload,
				PageSize:    pageSize,
			}

			nextWorkflow = req.GetWorkflow()
		case models.TASK_FETCH_OTHERS:
			req := FetchNextOthers{
				Name:        task.Name,
				FromPayload: fromPayload,
				PageSize:    pageSize,
			}

			nextWorkflow = req.GetWorkflow()
		case models.TASK_FETCH_PAYMENTS:
			req := FetchNextPayments{
				FromPayload: fromPayload,
				PageSize:    pageSize,
			}

			nextWorkflow = req.GetWorkflow()
		}

		scheduleHandle, err := w.temporalClient.ScheduleClient().Create(ctx.(context.Context), client.ScheduleOptions{
			// TODO(polo): id more specific
			ID: uuid.New().String(),
			Spec: client.ScheduleSpec{
				Intervals: []client.ScheduleIntervalSpec{
					{
						// TODO(polo): add config
						Every: defaultPolling,
					},
				},
			},
			Action: &client.ScheduleWorkflowAction{
				Workflow: nextWorkflow,
				// TODO(polo): add more args
				Args:      []interface{}{},
				TaskQueue: w.taskQueue,
				// TODO(polo): add retry policy
			},
			Overlap:            enums.SCHEDULE_OVERLAP_POLICY_SKIP,
			TriggerImmediately: true,
		})
		if err != nil {
			return err
		}

		// TODO(polo): store schedule handle
		_ = scheduleHandle
	}
	return nil
}
