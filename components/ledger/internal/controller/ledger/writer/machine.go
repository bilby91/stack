package writer

import (
	"context"
	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/machine/vm"
	"github.com/formancehq/ledger/internal/machine/vm/program"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/pkg/errors"
)

type MachineResult struct {
	Postings        ledger.Postings   `json:"postings"`
	Metadata        metadata.Metadata `json:"metadata"`
	AccountMetadata map[string]metadata.Metadata
}

type MachineInit struct {
	BoundedSources []string
	UnboundedSourcesAndDestinations []string
}

//go:generate mockgen -source machine.go -destination machine_generated.go -package writer . Machine
type Machine interface {
	// Init return all used accounts in the transaction
	Init(context.Context, map[string]string) (*MachineInit, error)
	Execute(context.Context) (*MachineResult, error)
}

type DefaultMachineAdapter struct {
	program program.Program
	machine *vm.Machine
	store   vm.Store
}

func (d *DefaultMachineAdapter) Init(ctx context.Context, vars map[string]string) (*MachineInit, error) {
	d.machine = vm.NewMachine(d.program)
	if err := d.machine.SetVarsFromJSON(vars); err != nil {
		return nil, errors.Wrap(err, "failed to set vars from JSON")
	}
	readLockAccounts, writeLockAccounts, err := d.machine.ResolveResources(ctx, d.store)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve resources")
	}

	return &MachineInit{
		BoundedSources:                  writeLockAccounts,
		UnboundedSourcesAndDestinations: readLockAccounts,
	}, nil
}

func (d *DefaultMachineAdapter) Execute(ctx context.Context) (*MachineResult, error) {
	if err := d.machine.ResolveBalances(ctx, d.store); err != nil {
		return nil, errors.Wrap(err, "failed to resolve balances")
	}

	if err := d.machine.Execute(); err != nil {
		return nil, errors.Wrap(err, "failed to execute machine")
	}

	return &MachineResult{
		Postings: collectionutils.Map(d.machine.Postings, func(from vm.Posting) ledger.Posting {
			return ledger.Posting{
				Source:      from.Source,
				Destination: from.Destination,
				Amount:      from.Amount.ToBigInt(),
				Asset:       from.Asset,
			}
		}),
		Metadata:        d.machine.GetTxMetaJSON(),
		AccountMetadata: d.machine.GetAccountsMetaJSON(),
	}, nil
}

func NewDefaultMachine(p program.Program, store vm.Store) *DefaultMachineAdapter {
	return &DefaultMachineAdapter{
		program: p,
		store:   store,
	}
}

var _ Machine = (*DefaultMachineAdapter)(nil)
