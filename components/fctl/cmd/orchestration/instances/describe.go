package instances

import (
	"fmt"
	"io"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	formance "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type InstancesDescribeStore struct {
	WorkflowInstancesHistory []shared.WorkflowInstanceHistory `json:"workflowInstanceHistory"`
}
type InstancesDescribeController struct {
	store *InstancesDescribeStore
}

var _ fctl.Controller[*InstancesDescribeStore] = (*InstancesDescribeController)(nil)

func NewDefaultInstancesDescribeStore() *InstancesDescribeStore {
	return &InstancesDescribeStore{}
}

func NewInstancesDescribeController() *InstancesDescribeController {
	return &InstancesDescribeController{
		store: NewDefaultInstancesDescribeStore(),
	}
}

func NewDescribeCommand() *cobra.Command {
	c := NewInstancesDescribeController()
	return fctl.NewCommand("describe <instance-id>",
		fctl.WithShortDescription("Describe a specific workflow instance"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*InstancesDescribeStore](c),
	)
}

func (c *InstancesDescribeController) GetStore() *InstancesDescribeStore {
	return c.store
}

func (c *InstancesDescribeController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	response, err := store.Client().Orchestration.V1.GetInstanceHistory(cmd.Context(), operations.GetInstanceHistoryRequest{
		InstanceID: args[0],
	})
	if err != nil {
		return nil, err
	}

	c.store.WorkflowInstancesHistory = response.GetWorkflowInstanceHistoryResponse.Data

	return c, nil
}

func (c *InstancesDescribeController) Render(cmd *cobra.Command, args []string) error {
	store := fctl.GetStackStore(cmd.Context())
	for i, history := range c.store.WorkflowInstancesHistory {
		if err := printStage(cmd, i, store.Client(), args[0], history); err != nil {
			return err
		}
	}

	return nil
}

func printHistoryBaseInfo(out io.Writer, name string, ind int, history shared.WorkflowInstanceHistory) {
	fctl.Section.WithWriter(out).Printf("Stage %d : %s\n", ind, name)
	fctl.BasicText.WithWriter(out).Printfln("Started at: %s", history.StartedAt.Format(time.RFC3339))
	if history.Terminated {
		fctl.BasicText.WithWriter(out).Printfln("Terminated at: %s", history.StartedAt.Format(time.RFC3339))
	}
}

func stageSourceName(src *shared.StageSendSource) string {
	switch {
	case src.Wallet != nil:
		return fmt.Sprintf("wallet '%s' (balance: %s)", src.Wallet.ID, *src.Wallet.Balance)
	case src.Account != nil:
		return fmt.Sprintf("account '%s' (ledger: %s)", src.Account.ID, *src.Account.Ledger)
	case src.Payment != nil:
		return fmt.Sprintf("payment '%s'", src.Payment.ID)
	default:
		return "unknown_source_type"
	}
}

func stageDestinationName(dst *shared.StageSendDestination) string {
	switch {
	case dst.Wallet != nil:
		return fmt.Sprintf("wallet '%s' (balance: %s)", dst.Wallet.ID, *dst.Wallet.Balance)
	case dst.Account != nil:
		return fmt.Sprintf("account '%s' (ledger: %s)", dst.Account.ID, *dst.Account.Ledger)
	case dst.Payment != nil:
		return dst.Payment.Psp
	default:
		return "unknown_source_type"
	}
}

func subjectName(src shared.Subject) string {
	switch {
	case src.WalletSubject != nil:
		return fmt.Sprintf("wallet %s (balance: %s)", src.WalletSubject.Identifier, *src.WalletSubject.Balance)
	case src.LedgerAccountSubject != nil:
		return fmt.Sprintf("account %s", src.LedgerAccountSubject.Identifier)
	default:
		return "unknown_subject_type"
	}
}

func printMetadata(metadata map[string]string) []pterm.BulletListItem {
	ret := make([]pterm.BulletListItem, 0)
	ret = append(ret, historyItemDetails("Added metadata:"))
	for k, v := range metadata {
		ret = append(ret, pterm.BulletListItem{
			Level: 2,
			Text:  fmt.Sprintf("%s: %s", k, v),
		})
	}
	return ret
}

func printStage(cmd *cobra.Command, i int, client *formance.Formance, id string, history shared.WorkflowInstanceHistory) error {
	cyanWriter := fctl.BasicTextCyan
	defaultWriter := fctl.BasicText

	listItems := make([]pterm.BulletListItem, 0)

	switch history.Input.Type {
	case shared.StageTypeStageSend:
		printHistoryBaseInfo(cmd.OutOrStdout(), "send", i, history)
		cyanWriter.Printfln("Send %v %s from %s to %s", history.Input.StageSend.Amount.Amount,
			history.Input.StageSend.Amount.Asset, stageSourceName(history.Input.StageSend.Source),
			stageDestinationName(history.Input.StageSend.Destination))
		fctl.Println()

		stageResponse, err := client.Orchestration.V1.GetInstanceStageHistory(cmd.Context(), operations.GetInstanceStageHistoryRequest{
			InstanceID: id,
			Number:     int64(i),
		})
		if err != nil {
			return err
		}

		for _, historyStage := range stageResponse.GetWorkflowInstanceHistoryStageResponse.Data {
			switch {
			case historyStage.Input.StripeTransfer != nil:
				listItems = append(listItems, historyItemTitle("Send %v %s to Stripe connected account: %s",
					*historyStage.Input.StripeTransfer.Amount,
					*historyStage.Input.StripeTransfer.Asset,
					*historyStage.Input.StripeTransfer.Destination,
				))
			case historyStage.Input.CreateTransaction != nil:
				listItems = append(listItems, historyItemTitle("Send %v %s from account %s to account %s (ledger %s)",
					historyStage.Input.CreateTransaction.Data.Postings[0].Amount,
					historyStage.Input.CreateTransaction.Data.Postings[0].Asset,
					historyStage.Input.CreateTransaction.Data.Postings[0].Source,
					historyStage.Input.CreateTransaction.Data.Postings[0].Destination,
					*historyStage.Input.CreateTransaction.Ledger,
				))
				if historyStage.Error == nil && historyStage.LastFailure == nil && historyStage.Terminated {
					listItems = append(listItems, historyItemDetails("Created transaction: %d", historyStage.Output.CreateTransaction.Data.ID))
					if historyStage.Input.CreateTransaction.Data.Reference != nil {
						listItems = append(listItems, historyItemDetails("Reference: %s", *historyStage.Output.CreateTransaction.Data.Reference))
					}
					if len(historyStage.Input.CreateTransaction.Data.Metadata) > 0 {
						listItems = append(listItems, printMetadata(historyStage.Input.CreateTransaction.Data.Metadata)...)
					}
				}
			case historyStage.Input.ConfirmHold != nil:
				listItems = append(listItems, historyItemTitle("Confirm debit hold %s", historyStage.Input.ConfirmHold.ID))
			case historyStage.Input.CreditWallet != nil:
				listItems = append(listItems, historyItemTitle("Credit wallet %s (balance: %s) of %v %s from %s",
					*historyStage.Input.CreditWallet.ID,
					*historyStage.Input.CreditWallet.Data.Balance,
					historyStage.Input.CreditWallet.Data.Amount.Amount,
					historyStage.Input.CreditWallet.Data.Amount.Asset,
					subjectName(historyStage.Input.CreditWallet.Data.Sources[0]),
				))
				if historyStage.Error == nil && historyStage.LastFailure == nil && historyStage.Terminated {
					if len(historyStage.Input.CreditWallet.Data.Metadata) > 0 {
						listItems = append(listItems, printMetadata(historyStage.Input.CreditWallet.Data.Metadata)...)
					}
				}
			case historyStage.Input.DebitWallet != nil:
				destination := "@world"
				if historyStage.Input.DebitWallet.Data.Destination != nil {
					destination = subjectName(*historyStage.Input.DebitWallet.Data.Destination)
				}

				listItems = append(listItems, historyItemTitle("Debit wallet %s (balance: %s) of %v %s to %s",
					*historyStage.Input.DebitWallet.ID,
					historyStage.Input.DebitWallet.Data.Balances[0],
					historyStage.Input.DebitWallet.Data.Amount.Amount,
					historyStage.Input.DebitWallet.Data.Amount.Asset,
					destination,
				))
				if historyStage.Error == nil && historyStage.LastFailure == nil && historyStage.Terminated {
					if len(historyStage.Input.DebitWallet.Data.Metadata) > 0 {
						listItems = append(listItems, printMetadata(historyStage.Input.DebitWallet.Data.Metadata)...)
					}
				}
			case historyStage.Input.GetAccount != nil:
				listItems = append(listItems, historyItemTitle("Read account %s of ledger %s",
					historyStage.Input.GetAccount.ID,
					historyStage.Input.GetAccount.Ledger,
				))
			case historyStage.Input.GetPayment != nil:
				listItems = append(listItems, historyItemTitle("Read payment %s",
					historyStage.Input.GetPayment.ID))
			case historyStage.Input.GetWallet != nil:
				listItems = append(listItems, historyItemTitle("Read wallet '%s'", historyStage.Input.GetWallet.ID))
			case historyStage.Input.RevertTransaction != nil:
				listItems = append(listItems, historyItemTitle("Revert transaction %s", historyStage.Input.RevertTransaction.ID))
				if historyStage.Error == nil {
					listItems = append(listItems, historyItemTitle("Created transaction: %d", historyStage.Output.RevertTransaction.Data.ID))
				}
			case historyStage.Input.VoidHold != nil:
				listItems = append(listItems, historyItemTitle("Cancel debit hold %s", historyStage.Input.VoidHold.ID))
			case historyStage.Input.ListWallets != nil:
				listItems = append(listItems, historyItemTitle("List wallets"))
			}
			if historyStage.LastFailure != nil {
				listItems = append(listItems, historyItemError(*historyStage.LastFailure))
				if historyStage.NextExecution != nil {
					listItems = append(listItems, historyItemError("Next try: %s", historyStage.NextExecution.Format(time.RFC3339)))
					listItems = append(listItems, historyItemError("Attempt: %d", historyStage.Attempt))
				}
			}
			if historyStage.Error != nil {
				listItems = append(listItems, historyItemError(*historyStage.Error))
			}
		}
	case shared.StageTypeStageDelay:
		printHistoryBaseInfo(cmd.OutOrStdout(), "delay", i, history)
		switch {
		case history.Input.StageDelay.Duration != nil:
			listItems = append(listItems, historyItemTitle("Pause workflow for a delay of %s", *history.Input.StageDelay.Duration))
		case history.Input.StageDelay.Until != nil:
			listItems = append(listItems, historyItemTitle("Pause workflow until %s", *history.Input.StageDelay.Until))
		}
	case shared.StageTypeStageWaitEvent:
		printHistoryBaseInfo(cmd.OutOrStdout(), "wait_event", i, history)
		listItems = append(listItems, historyItemTitle("Waiting event '%s'", history.Input.StageWaitEvent.Event))
		if history.Error == nil {
			if history.Terminated {
				listItems = append(listItems, historyItemDetails("Event received!"))
			} else {
				listItems = append(listItems, historyItemDetails("Still waiting event..."))
			}
		}
	case shared.StageTypeUpdate:
		printHistoryBaseInfo(cmd.OutOrStdout(), "update", i, history)
		switch {
		case history.Input.Update.Account != nil:
			account := history.Input.Update.Account
			listItems = append(listItems, historyItemTitle("Update account '%s' of ledger '%s'", account.ID, account.Ledger))
			listItems = append(listItems, printMetadata(account.Metadata)...)
		}
	default:
		// Display error?
	}
	if history.Error != nil {
		fctl.BasicTextRed.WithWriter(cmd.OutOrStdout()).Printfln("Stage terminated with error: %s", *history.Error)
	}

	if len(listItems) > 0 {
		defaultWriter.Print("History :\n")
		return pterm.DefaultBulletList.WithWriter(cmd.OutOrStdout()).WithItems(listItems).Render()
	}
	return nil
}

func historyItemTitle(format string, args ...any) pterm.BulletListItem {
	return pterm.BulletListItem{
		Level:     0,
		TextStyle: fctl.StyleGreen,
		Text:      fmt.Sprintf(format, args...),
	}
}

func historyItemDetails(format string, args ...any) pterm.BulletListItem {
	return pterm.BulletListItem{
		Level: 1,
		Text:  fmt.Sprintf(format, args...),
	}
}

func historyItemError(format string, args ...any) pterm.BulletListItem {
	return pterm.BulletListItem{
		Level:     1,
		TextStyle: fctl.StyleRed,
		Text:      fmt.Sprintf(format, args...),
	}
}
