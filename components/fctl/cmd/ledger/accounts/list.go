package accounts

import (
	"github.com/formancehq/stack/libs/go-libs/collectionutils"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	Accounts []shared.Account `json:"accounts"`
}
type ListController struct {
	store        *ListStore
	metadataFlag string
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}

func NewListController() *ListController {
	return &ListController{
		store:        NewDefaultListStore(),
		metadataFlag: "metadata",
	}
}

func NewListCommand() *cobra.Command {
	c := NewListController()
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List accounts"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithStringSliceFlag(c.metadataFlag, []string{}, "Filter accounts with metadata"),
		fctl.WithController[*ListStore](c),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetStackStore(cmd.Context())

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	body := make([]map[string]map[string]any, 0)
	for key, value := range metadata {
		body = append(body, map[string]map[string]any{
			"$match": {
				"metadata[" + key + "]": value,
			},
		})
	}

	request := operations.ListAccountsRequest{
		Ledger:   fctl.GetString(cmd, internal.LedgerFlag),
		Metadata: collectionutils.ConvertMap(metadata, collectionutils.ToAny[string]),
	}
	rsp, err := store.Client().Ledger.V1.ListAccounts(cmd.Context(), request)
	if err != nil {
		return nil, err
	}

	c.store.Accounts = rsp.AccountsCursorResponse.Cursor.Data

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {

	tableData := fctl.Map(c.store.Accounts, func(account shared.Account) []string {
		return []string{
			account.Address,
			fctl.MetadataAsShortString(account.Metadata),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"Address", "Metadata"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
