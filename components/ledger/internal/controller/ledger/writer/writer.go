package writer

import (
	"context"
	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/opentelemetry/tracer"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/time"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"slices"
)

type Writer struct {
	store          Store
	machineFactory MachineFactory
}

func (w *Writer) withTX(ctx context.Context, parameters Parameters, fn func(sqlTX TX) (*ledger.Log, error)) (*ledger.ChainedLog, error) {
	tx, err := w.store.BeginTX(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		// Ignore error, will be a noop if the transaction is already commited
		_ = tx.Rollback(ctx)
	}()

	log, err := fn(tx)
	if err != nil {
		return nil, err
	}
	log.IdempotencyKey = parameters.IdempotencyKey

	chainedLog, latency, err := tracer.TraceWithLatency(ctx, "InsertLog", func(ctx context.Context) (*ledger.ChainedLog, error) {
		return tx.InsertLog(ctx, *log)
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert log")
	}
	logging.FromContext(ctx).
		WithField("latency", latency.String()).
		Debugf("log inserted with id %d", chainedLog.ID)

	if parameters.DryRun {
		return chainedLog, errors.Wrap(tx.Rollback(ctx), "failed to rollback transaction")
	}

	// TODO: check errors for conflict on IK
	// if so, we can read from the database and respond with the already written log
	ret, latency, err := tracer.TraceWithLatency(ctx, "CommitTransaction", func(ctx context.Context) (*ledger.ChainedLog, error) {
		return chainedLog, errors.Wrap(tx.Commit(ctx), "failed to commit transaction")
	})
	if err != nil {
		return nil, err
	}

	logging.FromContext(ctx).
		WithField("latency", latency.String()).
		Debugf("store transaction commited")

	return ret, nil
}

// todo: handle deadlocks
func (w *Writer) CreateTransaction(ctx context.Context, parameters Parameters, runScript ledger.RunScript) (*ledger.Transaction, error) {
	logger := logging.FromContext(ctx).
		WithField("req", uuid.NewString()[:8])
	ctx = logging.ContextWithLogger(ctx, logger)

	machine, err := w.machineFactory.Make(runScript.Plain)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile script")
	}

	machineInit, err := machine.Init(ctx, runScript.Vars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init program")
	}

	logger.WithFields(map[string]any{
		"boundedSources": machineInit.BoundedSources,
		"otherAccounts":  machineInit.UnboundedSourcesAndDestinations,
	}).Debugf("creating new machine")

	slices.Sort(machineInit.BoundedSources)
	slices.Sort(machineInit.UnboundedSourcesAndDestinations)

	// todo: add logs
	log, err := w.withTX(ctx, parameters, func(sqlTX TX) (*ledger.Log, error) {

		if len(machineInit.BoundedSources) > 0 {
			_, latency, err := tracer.TraceWithLatency(ctx, "LockAccounts", func(ctx context.Context) (*struct{}, error) {
				return nil, sqlTX.LockAccounts(ctx, machineInit.BoundedSources...)
			}, func(ctx context.Context, _ *struct{}) {
				trace.SpanFromContext(ctx).SetAttributes(
					attribute.StringSlice("accounts", machineInit.BoundedSources),
				)
			})
			if err != nil {
				return nil, errors.Wrap(err, "failed to acquire accounts locks")
			}
			logger.WithFields(map[string]any{
				"latency":  latency.String(),
				"accounts": machineInit.BoundedSources,
			}).Debugf("locked accounts")
		} else {
			logger.Debugf("no bounded sources to lock")
		}

		result, latency, err := tracer.TraceWithLatency(ctx, "ExecuteMachine", func(ctx context.Context) (*MachineResult, error) {
			return machine.Execute(ctx)
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to execute program")
		}

		logger.WithFields(map[string]any{
			"latency": latency.String(),
		}).Debugf("machine executed")

		transaction, latency, err := tracer.TraceWithLatency(ctx, "InsertTransaction", func(ctx context.Context) (*ledger.Transaction, error) {
			return sqlTX.InsertTransaction(ctx,
				ledger.NewTransactionData().
					WithPostings(result.Postings...).
					WithDate(runScript.Timestamp). // If empty will be filled by the database
					WithReference(runScript.Reference),
			)
		}, func(ctx context.Context, tx *ledger.Transaction) {
			trace.SpanFromContext(ctx).SetAttributes(
				attribute.Int("id", tx.ID),
				attribute.String("timestamp", tx.Timestamp.Format(time.RFC3339Nano)),
			)
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert transaction")
		}

		logger.WithFields(map[string]any{
			"latency": latency.String(),
			"txID":    transaction.ID,
		}).Debugf("transaction inserted")

		for _, address := range transaction.GetMoves().InvolvedAccounts() {
			upserted, latency, err := tracer.TraceWithLatency(ctx, "UpsertAccount", func(ctx context.Context) (bool, error) {
				return sqlTX.UpsertAccount(ctx, ledger.Account{
					Address:    address,
					Metadata:   result.AccountMetadata[address],
					FirstUsage: transaction.Timestamp,
				})
			}, func(ctx context.Context, upserted bool) {
				trace.SpanFromContext(ctx).SetAttributes(
					attribute.String("address", address),
					attribute.Bool("upserted", upserted),
				)
			})
			if err != nil {
				return nil, errors.Wrap(err, "failed to upsert account")
			} else if upserted {
				logger.WithField("latency", latency.String()).Debugf("account upserted")
			} else {
				logger.WithField("latency", latency.String()).Debugf("account not modified")
			}
		}

		for _, account := range transaction.GetMoves().InvolvedAccounts() {
			_, latency, err = tracer.TraceWithLatency(ctx, "LockAccounts", func(ctx context.Context) (struct{}, error) {
				return struct{}{}, sqlTX.LockAccounts(ctx, account)
			}, func(ctx context.Context, _ struct{}) {
				trace.SpanFromContext(ctx).SetAttributes(
					attribute.StringSlice("accounts", transaction.GetMoves().InvolvedAccounts()),
				)
			})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to acquire account lock on %s", account)
			}
			logger.WithField("latency", latency.String()).Debugf("account locked: %s", account)
		}

		_, latency, err = tracer.TraceWithLatency(ctx, "InsertMoves", func(ctx context.Context) (struct{}, error) {
			return struct{}{}, sqlTX.InsertMoves(ctx, transaction.GetMoves()...)
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert moves")
		}

		logger.WithField("latency", latency.String()).Debugf("moves inserted")

		// notes(gfyrag): force date to be zero to let postgres fill it
		// todo: clean that
		return pointer.For(ledger.NewTransactionLogWithDate(*transaction, result.AccountMetadata, time.Time{})), err
	})
	if err != nil {
		logger.Errorf("failed to create transaction: %s", err)
		return nil, err
	}

	return pointer.For(log.Data.(ledger.NewTransactionLogPayload).Transaction), nil
}

func (w *Writer) RevertTransaction(ctx context.Context, parameters Parameters, id int, force bool, atEffectiveDate bool) (*ledger.Transaction, error) {
	log, err := w.withTX(ctx, parameters, func(sqlTX TX) (*ledger.Log, error) {
		// todo: check if account has enough funds
		// no need to use numscript for that, we just n
		// todo reimplement force
		originalTransaction, hasBeenReverted, err := sqlTX.RevertTransaction(ctx, id)
		if err != nil {
			return nil, err
		}
		if !hasBeenReverted {
			return nil, errors.New("transaction already reverted")
		}

		transaction, err := sqlTX.InsertTransaction(ctx, originalTransaction.Reverse(atEffectiveDate))
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert transaction")
		}

		return ledger.NewRevertedTransactionLog(time.Time{}, id, transaction), nil
	})
	if err != nil {
		return nil, err
	}

	return log.Data.(ledger.RevertedTransactionLogPayload).RevertTransaction, nil
}

func (w *Writer) SaveMeta(ctx context.Context, parameters Parameters, targetType string, targetID any, m metadata.Metadata) error {

	_, err := w.withTX(ctx, parameters, func(sqlTX TX) (*ledger.Log, error) {
		switch targetType {
		case ledger.MetaTargetTypeTransaction:
			_, err := sqlTX.UpdateTransactionMetadata(ctx, targetID.(int), m)
			if err != nil {
				return nil, err
			}
		case ledger.MetaTargetTypeAccount:
			if err := sqlTX.UpdateAccountMetadata(ctx, targetID.(string), m); err != nil {
				return nil, err
			}
		default:
			panic(errors.Errorf("unknown target type '%s'", targetType))
		}

		return ledger.NewSetMetadataLog(time.Now(), ledger.SetMetadataLogPayload{
			TargetType: targetType,
			TargetID:   targetID,
			Metadata:   m,
		}), nil
	})
	return err
}

func (w *Writer) DeleteMetadata(ctx context.Context, parameters Parameters, targetType string, targetID any, key string) error {
	_, err := w.withTX(ctx, parameters, func(sqlTX TX) (*ledger.Log, error) {
		switch targetType {
		case ledger.MetaTargetTypeTransaction:
			_, err := sqlTX.DeleteTransactionMetadata(ctx, targetID.(int), key)
			if err != nil {
				return nil, err
			}
		case ledger.MetaTargetTypeAccount:
			if err := sqlTX.DeleteAccountMetadata(ctx, targetID.(string), key); err != nil {
				return nil, err
			}
		default:
			panic(errors.Errorf("unknown target type '%s'", targetType))
		}

		return ledger.NewDeleteMetadataLog(time.Now(), ledger.DeleteMetadataLogPayload{
			TargetType: targetType,
			TargetID:   targetID,
			Key:        key,
		}), nil
	})
	return err
}

func New(store Store, machineFactory MachineFactory) *Writer {
	return &Writer{
		store:          store,
		machineFactory: machineFactory,
	}
}
