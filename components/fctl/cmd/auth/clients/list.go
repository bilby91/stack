package clients

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type Client struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IsPublic    string   `json:"isPublic"`
	Scopes      []string `json:"scopes"`
}

type ListStore struct {
	Clients []Client `json:"clients"`
}
type ListController struct {
	store *ListStore
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}

func NewListController() *ListController {
	return &ListController{
		store: NewDefaultListStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithShortDescription("List clients"),
		fctl.WithController[*ListStore](NewListController()),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	clients, err := store.Client().Auth.V1.ListClients(cmd.Context())
	if err != nil {
		return nil, err
	}

	if clients.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", clients.StatusCode)
	}

	c.store.Clients = fctl.Map(clients.ListClientsResponse.Data, func(o shared.Client) Client {
		return Client{
			ID:   o.ID,
			Name: o.Name,
			Description: func() string {
				if o.Description == nil {
					return ""
				}
				return ""
			}(),
			Scopes:   o.Scopes,
			IsPublic: fctl.BoolPointerToString(o.Public),
		}
	})

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Clients, func(o Client) []string {
		return []string{
			o.ID,
			o.Name,
			o.Description,
			o.IsPublic,
			fmt.Sprintf("%d", len(o.Scopes)),
		}
	})

	tableData = fctl.Prepend(tableData, []string{"ID", "Name", "Description", "Public", "Permissions"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
