package bankaccounts

import (
	"fmt"
	"time"

	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	BankAccount *shared.BankAccount `json:"bankAccount"`
}
type ShowController struct {
	PaymentsVersion versions.Version

	store *ShowStore
}

func (c *ShowController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

func NewShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowController() *ShowController {
	return &ShowController{
		store: NewShowStore(),
	}
}

func NewShowCommand() *cobra.Command {
	c := NewShowController()
	return fctl.NewCommand("get <bankAccountID>",
		fctl.WithShortDescription("Get bank account"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("sh", "s"),
		fctl.WithController[*ShowStore](c),
	)
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if err := versions.GetPaymentsVersion(cmd, args, c); err != nil {
		return nil, err
	}

	if c.PaymentsVersion < versions.V1 {
		return nil, fmt.Errorf("bank accounts are only supported in >= v1.0.0")
	}

	response, err := store.Client().Payments.V1.GetBankAccount(cmd.Context(), operations.GetBankAccountRequest{
		BankAccountID: args[0],
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.BankAccount = &response.BankAccountResponse.Data

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) error {
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.BankAccount.ID})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), c.store.BankAccount.Name})
	tableData = append(tableData, []string{pterm.LightCyan("CreatedAt"), c.store.BankAccount.CreatedAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("Country"), c.store.BankAccount.Country})
	if c.store.BankAccount.AccountNumber != nil {
		tableData = append(tableData, []string{pterm.LightCyan("AccountNumber"), *c.store.BankAccount.AccountNumber})
	}
	if c.store.BankAccount.Iban != nil {
		tableData = append(tableData, []string{pterm.LightCyan("Iban"), *c.store.BankAccount.Iban})
	}
	if c.store.BankAccount.SwiftBicCode != nil {
		tableData = append(tableData, []string{pterm.LightCyan("SwiftBicCode"), *c.store.BankAccount.SwiftBicCode})
	}

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("RelatedAccounts")
	tableData = fctl.Map(c.store.BankAccount.RelatedAccounts, func(ba shared.BankAccountRelatedAccounts) []string {
		return []string{
			ba.ID,
			ba.CreatedAt.Format(time.RFC3339),
			ba.ConnectorID,
			ba.Provider,
			ba.AccountID,
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "CreatedAt", "ConnectorID", "Provider", "AccountID"})
	if err := pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return fctl.PrintMetadata(cmd.OutOrStdout(), c.store.BankAccount.Metadata)
}
