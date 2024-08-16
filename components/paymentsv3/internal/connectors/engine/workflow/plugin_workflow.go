package workflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/formancehq/paymentsv3/internal/connectors/engine/activities"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

func (w Workflow) run(
	ctx workflow.Context,
	plugin models.Plugin,
	config models.Config,
	connectorID models.ConnectorID,
	fromPayload json.RawMessage,
	taskTree []models.TaskTree,
) error {
	var nextWorkflow interface{}
	var request interface{}
	var capability models.Capability
	metadata := make(map[string]string)
	for _, task := range taskTree {
		switch task.TaskType {
		case models.TASK_FETCH_ACCOUNTS:
			req := FetchNextAccounts{
				Config:      config,
				ConnectorID: connectorID,
				FromPayload: fromPayload,
			}

			nextWorkflow = RunFetchNextAccounts
			request = req
			capability = models.CAPABILITY_FETCH_ACCOUNTS
		case models.TASK_FETCH_EXTERNAL_ACCOUNTS:
			req := FetchNextExternalAccounts{
				Config:      config,
				ConnectorID: connectorID,
				FromPayload: fromPayload,
			}

			nextWorkflow = RunFetchNextExternalAccounts
			request = req
			capability = models.CAPABILITY_FETCH_EXTERNAL_ACCOUNTS
		case models.TASK_FETCH_OTHERS:
			req := FetchNextOthers{
				Config:      config,
				ConnectorID: connectorID,
				Name:        task.Name,
				FromPayload: fromPayload,
			}

			nextWorkflow = RunFetchNextOthers
			request = req
			capability = models.CAPABILITY_FETCH_OTHERS
			metadata["name"] = task.Name
		case models.TASK_FETCH_PAYMENTS:
			req := FetchNextPayments{
				Config:      config,
				ConnectorID: connectorID,
				FromPayload: fromPayload,
			}

			nextWorkflow = RunFetchNextPayments
			request = req
			capability = models.CAPABILITY_FETCH_PAYMENTS
		default:
			return fmt.Errorf("unknown task type: %v", task.TaskType)
		}

		// Create next wk in database
		wk := models.Workflow{
			ID:          uuid.New().String(),
			ConnectorID: connectorID,
			CreatedAt:   workflow.Now(ctx).UTC(),
			Capability:  capability,
			Metadata:    metadata,
		}
		err := activities.StorageWorkflowsStore(
			infiniteRetryContext(ctx),
			wk,
		)
		if err != nil {
			return err
		}

		// Schedule next workflow every polling duration
		scheduleHandle, err := w.temporalClient.ScheduleClient().Create(ctx.(context.Context), client.ScheduleOptions{
			ID: uuid.New().String(),
			Spec: client.ScheduleSpec{
				Intervals: []client.ScheduleIntervalSpec{
					{
						Every: config.PollingDuration,
					},
				},
			},
			Action: &client.ScheduleWorkflowAction{
				Workflow: nextWorkflow,
				Args: []interface{}{
					plugin,
					request,
					task.NextTasks,
				},
				TaskQueue: connectorID.Reference,
				// Search attributes are used to query workflows
				SearchAttributes: map[string]any{
					SearchAttributeWorkflowID: wk.ID,
				},
			},
			Overlap:            enums.SCHEDULE_OVERLAP_POLICY_SKIP,
			TriggerImmediately: true,
		})
		if err != nil {
			return err
		}

		err = activities.StorageSchedulesStore(ctx, models.Schedule{
			ID:          scheduleHandle.GetID(),
			ConnectorID: connectorID,
			CreatedAt:   workflow.Now(ctx).UTC(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

var Run any

func init() {
	Run = Workflow{}.run
}
