package activities

import (
	"errors"

	"github.com/formancehq/paymentsv3/internal/connectors/engine/plugins"
	"github.com/formancehq/paymentsv3/internal/storage"
	temporalworker "github.com/formancehq/stack/libs/go-libs/temporal"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Activities struct {
	storage storage.Storage
}

func (a Activities) DefinitionSet() temporalworker.DefinitionSet {
	return temporalworker.NewDefinitionSet().
		Append(temporalworker.Definition{
			Name: "PluginInstallConnector",
			Func: a.PluginInstallConnector,
		}).
		Append(temporalworker.Definition{
			Name: "PluginFetchNextAccounts",
			Func: a.PluginFetchNextAccounts,
		}).
		Append(temporalworker.Definition{
			Name: "PluginFetchNextExternalAccounts",
			Func: a.PluginFetchNextExternalAccounts,
		}).
		Append(temporalworker.Definition{
			Name: "PluginFetchNextPayments",
			Func: a.PluginFetchNextPayments,
		}).
		Append(temporalworker.Definition{
			Name: "PluginFetchNextOthers",
			Func: a.PluginFetchNextOthers,
		}).
		Append(temporalworker.Definition{
			Name: "StorageStoreAccounts",
			Func: a.StorageStoreAccounts,
		}).
		Append(temporalworker.Definition{
			Name: "StorageDeleteAccounts",
			Func: a.StorageDeleteAccounts,
		}).
		Append(temporalworker.Definition{
			Name: "StorageStorePayments",
			Func: a.StorageStorePayments,
		}).
		Append(temporalworker.Definition{
			Name: "StorageDeletePayments",
			Func: a.StorageDeletePayments,
		}).
		Append(temporalworker.Definition{
			Name: "StorageFetchState",
			Func: a.StorageFetchState,
		}).
		Append(temporalworker.Definition{
			Name: "StorageStoreState",
			Func: a.StorageStoreState,
		}).
		Append(temporalworker.Definition{
			Name: "StorageDeleteStates",
			Func: a.StorageDeleteStates,
		}).
		Append(temporalworker.Definition{
			Name: "StorageStoreTasksTree",
			Func: a.StorageStoreTasksTree,
		}).
		Append(temporalworker.Definition{
			Name: "StorageDeleteTasksTree",
			Func: a.StorageDeleteTasksTree,
		}).
		Append(temporalworker.Definition{
			Name: "StorageStoreConnector",
			Func: a.StorageStoreConnector,
		}).
		Append(temporalworker.Definition{
			Name: "StorageDeleteConnector",
			Func: a.StorageDeleteConnector,
		}).
		Append(temporalworker.Definition{
			Name: "StorageStoreSchedule",
			Func: a.StorageStoreSchedule,
		}).
		Append(temporalworker.Definition{
			Name: "StorageFetchSchedules",
			Func: a.StorageFetchSchedules,
		}).
		Append(temporalworker.Definition{
			Name: "StorageDeleteSchedules",
			Func: a.StorageDeleteSchedules,
		}).
		Append(temporalworker.Definition{
			Name: "StorageStoreInstance",
			Func: a.StorageStoreInstance,
		}).
		Append(temporalworker.Definition{
			Name: "StorageUpdateInstance",
			Func: a.StorageUpdateInstance,
		}).
		Append(temporalworker.Definition{
			Name: "StorageDeleteInstances",
			Func: a.StorageDeleteInstances,
		}).
		Append(temporalworker.Definition{
			Name: "StorageStoreWorkflow",
			Func: a.StorageStoreWorkflow,
		}).
		Append(temporalworker.Definition{
			Name: "StorageDeleteWorkflow",
			Func: a.StorageDeleteWorkflow,
		})
}

func New(storage storage.Storage, plugins plugins.Plugins) Activities {
	return Activities{
		storage: storage,
	}
}

func executeActivity(ctx workflow.Context, activity any, ret any, args ...any) error {
	if err := workflow.ExecuteActivity(ctx, activity, args...).Get(ctx, ret); err != nil {
		var timeoutError *temporal.TimeoutError
		if errors.As(err, &timeoutError) {
			return errors.New(timeoutError.Message())
		}
		return err
	}
	return nil
}
