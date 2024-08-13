package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageStoreSchedule(ctx context.Context, schedule models.Schedule) error {
	return a.storage.UpsertSchedule(ctx, schedule)
}

var StorageStoreScheduleActivity = Activities{}.StorageStoreSchedule

func StorageStoreSchedule(ctx workflow.Context, schedule models.Schedule) error {
	return executeActivity(ctx, StorageStoreScheduleActivity, nil, schedule)
}
