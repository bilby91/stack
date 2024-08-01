package activities

import (
	"errors"

	"github.com/formancehq/paymentsv3/internal/models"
	temporalworker "github.com/formancehq/paymentsv3/internal/temporal/worker"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Activities struct {
	plugin models.Plugin
}

func (a Activities) DefinitionSet() temporalworker.DefinitionSet {
	return temporalworker.NewDefinitionSet().
		Append(temporalworker.Definition{
			Name: "FetchNextAccounts",
			Func: a.FetchNextAccounts,
		}).
		Append(temporalworker.Definition{
			Name: "FetchNextPayments",
			Func: a.FetchNextPayments,
		}).
		Append(temporalworker.Definition{
			Name: "FetchNextOthers",
			Func: a.FetchNextOthers,
		})
}

func New(plugin models.Plugin) Activities {
	return Activities{
		plugin: plugin,
	}
}

func executeActivity(ctx workflow.Context, activity any, ret any, request any) error {
	if err := workflow.ExecuteActivity(ctx, activity, request).Get(ctx, ret); err != nil {
		var timeoutError *temporal.TimeoutError
		if errors.As(err, &timeoutError) {
			return errors.New(timeoutError.Message())
		}
		return err
	}
	return nil
}
