package balances

import (
	"math/big"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CreateStore struct {
	BalanceName string `json:"balanceName"`
}
type CreateController struct {
	store *CreateStore
}

const expiresAtFlag = "expires-at"

const priorityFlag = "priority"

var _ fctl.Controller[*CreateStore] = (*CreateController)(nil)

func NewDefaultCreateStore() *CreateStore {
	return &CreateStore{}
}

func NewCreateController() *CreateController {
	return &CreateController{
		store: NewDefaultCreateStore(),
	}
}

func NewCreateCommand() *cobra.Command {
	return fctl.NewCommand("create <balance-name>",
		fctl.WithShortDescription("Create a new balance"),
		fctl.WithAliases("c", "cr"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		internal.WithTargetingWalletByID(),
		internal.WithTargetingWalletByName(),
		fctl.WithStringFlag(expiresAtFlag, "", "Balance expiration date"),
		fctl.WithIntFlag(priorityFlag, 0, "Balance priority"),
		fctl.WithController[*CreateStore](NewCreateController()),
	)
}

func (c *CreateController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	walletID, err := internal.RequireWalletID(cmd, store.Client())
	if err != nil {
		return nil, err
	}

	expiresAt, err := fctl.GetDateTime(cmd, expiresAtFlag)
	if err != nil {
		return nil, err
	}

	var priority *big.Int = nil
	priorityInt := fctl.GetInt(cmd, priorityFlag)
	if priorityInt != 0 {
		priority = big.NewInt(int64(priorityInt))
	}

	request := operations.CreateBalanceRequest{
		ID: walletID,
		CreateBalanceRequest: &shared.CreateBalanceRequest{
			Name:      args[0],
			ExpiresAt: expiresAt,
			Priority:  priority,
		},
	}
	response, err := store.Client().Wallets.V1.CreateBalance(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "creating balance")
	}

	c.store.BalanceName = response.CreateBalanceResponse.Data.Name
	return c, nil
}

func (c *CreateController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln(
		"Balance created successfully with name: %s", c.store.BalanceName)
	return nil

}
