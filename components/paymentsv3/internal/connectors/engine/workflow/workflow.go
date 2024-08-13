package workflow

import (
	temporalworker "github.com/formancehq/stack/libs/go-libs/temporal"
	"go.temporal.io/sdk/client"
)

const SearchAttributeWorkflowID = "PaymentWorkflowID"

type Workflow struct {
	temporalClient client.Client
}

func New(temporalClient client.Client) Workflow {
	return Workflow{
		temporalClient: temporalClient,
	}
}

func (w Workflow) DefinitionSet() temporalworker.DefinitionSet {
	return temporalworker.NewDefinitionSet().
		Append(temporalworker.Definition{
			Name: "FetchAccounts",
			Func: w.runFetchNextAccounts,
		}).
		Append(temporalworker.Definition{
			Name: "FetchExternalAccounts",
			Func: w.runFetchNextExternalAccounts,
		}).
		Append(temporalworker.Definition{
			Name: "FetchOthers",
			Func: w.runFetchNextOthers,
		}).
		Append(temporalworker.Definition{
			Name: "FetchPayments",
			Func: w.runFetchNextPayments,
		}).
		Append(temporalworker.Definition{
			Name: "TerminateSchedules",
			Func: w.runTerminateSchedules,
		}).
		Append(temporalworker.Definition{
			Name: "InstallConnector",
			Func: w.runInstallConnector,
		}).
		Append(temporalworker.Definition{
			Name: "UninstallConnector",
			Func: w.runUninstallConnector,
		}).
		Append(temporalworker.Definition{
			Name: "Run",
			Func: w.run,
		})
}
