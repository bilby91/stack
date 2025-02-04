package install

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type PaymentsConnectorsAtlarStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
	ConnectorID   string `json:"connectorId"`
}

type PaymentsConnectorsAtlarController struct {
	store *PaymentsConnectorsAtlarStore
}

var _ fctl.Controller[*PaymentsConnectorsAtlarStore] = (*PaymentsConnectorsAtlarController)(nil)

func NewDefaultPaymentsConnectorsAtlarStore() *PaymentsConnectorsAtlarStore {
	return &PaymentsConnectorsAtlarStore{
		Success: false,
	}
}

func NewPaymentsConnectorsAtlarController() *PaymentsConnectorsAtlarController {
	return &PaymentsConnectorsAtlarController{
		store: NewDefaultPaymentsConnectorsAtlarStore(),
	}
}

func NewAtlarCommand() *cobra.Command {
	c := NewPaymentsConnectorsAtlarController()
	return fctl.NewCommand(internal.AtlarConnector+" <file>|-",
		fctl.WithShortDescription("Install an atlar connector"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*PaymentsConnectorsAtlarStore](c),
	)
}

func (c *PaymentsConnectorsAtlarController) GetStore() *PaymentsConnectorsAtlarStore {
	return c.store
}

func (c *PaymentsConnectorsAtlarController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())
	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to install connector '%s'", internal.AtlarConnector) {
		return nil, fctl.ErrMissingApproval
	}
	script, err := fctl.ReadFile(cmd, store.Stack(), args[0])
	if err != nil {
		return nil, err
	}

	var config shared.AtlarConfig
	if err := json.Unmarshal([]byte(script), &config); err != nil {
		return nil, err
	}

	response, err := store.Client().Payments.V1.InstallConnector(cmd.Context(), operations.InstallConnectorRequest{
		ConnectorConfig: shared.ConnectorConfig{
			AtlarConfig: &config,
		},
		Connector: shared.ConnectorAtlar,
	})
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorName = internal.AtlarConnector

	if response.ConnectorResponse != nil {
		c.store.ConnectorID = response.ConnectorResponse.Data.ConnectorID
	}

	return c, nil
}

func (c *PaymentsConnectorsAtlarController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ConnectorID == "" {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector installed!", c.store.ConnectorName)
	} else {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector '%s' installed!", c.store.ConnectorName, c.store.ConnectorID)
	}

	return nil
}
