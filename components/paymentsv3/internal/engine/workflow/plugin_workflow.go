package workflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

func (w Workflow) runNextWorkflow(
	ctx workflow.Context,
	config models.Config,
	connectorID models.ConnectorID,
	fromPayload json.RawMessage,
	taskTree []models.TaskTree,
) error {

	var nextWorkflow interface{}
	var request interface{}
	for _, task := range taskTree {
		switch task.TaskType {
		case models.TASK_FETCH_ACCOUNTS:
			req := FetchNextAccounts{
				Config:      config,
				ConnectorID: connectorID,
				FromPayload: fromPayload,
			}

			nextWorkflow = req.GetWorkflow()
			request = req
		case models.TASK_FETCH_EXTERNAL_ACCOUNTS:
			req := FetchNextExternalAccounts{
				Config:      config,
				ConnectorID: connectorID,
				FromPayload: fromPayload,
			}

			nextWorkflow = req.GetWorkflow()
			request = req
		case models.TASK_FETCH_OTHERS:
			req := FetchNextOthers{
				Config:      config,
				ConnectorID: connectorID,
				Name:        task.Name,
				FromPayload: fromPayload,
			}

			nextWorkflow = req.GetWorkflow()
			request = req
		case models.TASK_FETCH_PAYMENTS:
			req := FetchNextPayments{
				Config:      config,
				ConnectorID: connectorID,
				FromPayload: fromPayload,
			}

			nextWorkflow = req.GetWorkflow()
			request = req
		default:
			return fmt.Errorf("unknown task type: %v", task.TaskType)
		}

		scheduleHandle, err := w.temporalClient.ScheduleClient().Create(ctx.(context.Context), client.ScheduleOptions{
			// TODO(polo): id more specific ?
			ID: uuid.New().String(),
			Spec: client.ScheduleSpec{
				Intervals: []client.ScheduleIntervalSpec{
					{
						Every: config.PollingIntervalDuration(),
					},
				},
			},
			Action: &client.ScheduleWorkflowAction{
				Workflow: nextWorkflow,
				Args: []interface{}{
					request,
					task.NextTasks,
				},
				TaskQueue: w.taskQueue,
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
