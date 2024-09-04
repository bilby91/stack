package writer

import (
	"context"
	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/time"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)

	store := NewMockStore(ctrl)
	machine := NewMockMachine(ctrl)
	machineFactory := NewMockMachineFactory(ctrl)
	sqlTX := NewMockTX(ctrl)

	l := New(store, machineFactory)

	runScript := ledger.RunScript{}

	machineFactory.EXPECT().
		Make(runScript.Plain).
		Return(machine, nil)

	machine.EXPECT().
		Init(gomock.Any(), runScript.Vars).
		Return([]string{"b", "a", "c"}, nil)

	store.EXPECT().
		BeginTX(gomock.Any()).
		Return(sqlTX, nil)

	sqlTX.EXPECT().
		LockAccounts(gomock.Any(), "a", "b", "c").
		Return(nil)

	machine.EXPECT().
		Execute(gomock.Any()).
		Return(&MachineResult{}, nil)

	tx := ledger.NewTransaction()
	sqlTX.EXPECT().
		InsertTransaction(gomock.Any(), ledger.NewTransactionData()).
		Return(tx, nil)

	sqlTX.EXPECT().
		InsertLog(gomock.Any(), ledger.NewTransactionLogWithDate(tx, nil, time.Time{})).
		Return(pointer.For(ledger.NewTransactionLog(tx, nil).ChainLog(nil)), nil)

	sqlTX.EXPECT().
		Commit(gomock.Any()).
		Return(nil)
	sqlTX.EXPECT().
		Rollback(gomock.Any()).
		Return(errors.New("already commited"))

	_, err := l.CreateTransaction(context.Background(), Parameters{}, runScript)
	require.NoError(t, err)
}
