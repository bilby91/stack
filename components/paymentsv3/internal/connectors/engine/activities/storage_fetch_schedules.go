package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/paymentsv3/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageFetchSchedules(ctx context.Context, query storage.ListSchedulesQuery) (*bunpaginate.Cursor[models.Schedule], error) {
	return a.storage.ListSchedules(ctx, query)
}

var StorageFetchSchedulesActivity = Activities{}.StorageFetchSchedules

func StorageFetchSchedules(ctx workflow.Context, query storage.ListSchedulesQuery) (*bunpaginate.Cursor[models.Schedule], error) {
	ret := bunpaginate.Cursor[models.Schedule]{}
	if err := executeActivity(ctx, StorageFetchSchedulesActivity, &ret, query); err != nil {
		return nil, err
	}
	return &ret, nil
}
