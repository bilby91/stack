package webhooks

import (
	"strings"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListWebhookStore struct {
	Webhooks []shared.WebhooksConfig `json:"webhooks"`
}
type ListWebhookController struct {
	store *ListWebhookStore
}

var _ fctl.Controller[*ListWebhookStore] = (*ListWebhookController)(nil)

func NewDefaultListWebhookStore() *ListWebhookStore {
	return &ListWebhookStore{
		Webhooks: []shared.WebhooksConfig{},
	}
}

func NewListWebhookController() *ListWebhookController {
	return &ListWebhookController{
		store: NewDefaultListWebhookStore(),
	}
}
func (c *ListWebhookController) GetStore() *ListWebhookStore {
	return c.store
}

func (c *ListWebhookController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())
	request := operations.GetManyConfigsRequest{}
	response, err := store.Client().Webhooks.V1.GetManyConfigs(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "listing all config")
	}

	c.store.Webhooks = response.ConfigsResponse.Cursor.Data

	return c, nil
}

func (c *ListWebhookController) Render(cmd *cobra.Command, args []string) error {
	// TODO: WebhooksConfig is missing ?
	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(cmd.OutOrStdout()).
		WithData(
			fctl.Prepend(
				fctl.Map(c.store.Webhooks,
					func(src shared.WebhooksConfig) []string {
						return []string{
							src.ID,
							src.CreatedAt.Format(time.RFC3339),
							src.Secret,
							src.Endpoint,
							fctl.BoolToString(src.Active),
							strings.Join(src.EventTypes, ","),
						}
					}),
				[]string{"ID", "Created at", "Secret", "Endpoint", "Active", "Event types"},
			),
		).Render(); err != nil {
		return errors.Wrap(err, "rendering table")
	}
	return nil
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithShortDescription("List all configs"),
		fctl.WithAliases("ls", "l"),
		fctl.WithController[*ListWebhookStore](NewListWebhookController()),
	)
}
