package workflow

import (
	"context"

	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/query"
	"go.temporal.io/sdk/workflow"
)

type TerminateSchedules struct {
	ConnectorID models.ConnectorID
}

func (w Workflow) runTerminateSchedules(
	ctx workflow.Context,
	terminateSchedules TerminateSchedules,
) error {
	query := storage.NewListSchedulesQuery(
		bunpaginate.NewPaginatedQueryOptions(storage.ScheduleQuery{}).
			WithPageSize(100).
			WithQueryBuilder(
				query.Match("connector_id", terminateSchedules.ConnectorID),
			),
	)
	for {
		schedules, err := activities.StorageSchedulesList(infiniteRetryContext(ctx), query)
		if err != nil {
			return err
		}

		for _, schedule := range schedules.Data {
			// TODO(polo): context.Background() ?
			scheduleHandler := w.temporalClient.ScheduleClient().GetHandle(context.Background(), schedule.ID)
			if err := scheduleHandler.Delete(context.Background()); err != nil {
				// TODO(polo): log error but continue
			}
		}

		if !schedules.HasMore {
			break
		}

		err = bunpaginate.UnmarshalCursor(schedules.Next, &query)
		if err != nil {
			return err
		}
	}

	return nil
}

var RunTerminateSchedules any

func init() {
	RunTerminateSchedules = Workflow{}.runTerminateSchedules
}
