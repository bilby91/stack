package ledger

import (
	"context"
	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/time"
	"github.com/uptrace/bun"
	"strings"
)

type Move struct {
	bun.BaseModel `bun:"table:moves"`

	Ledger              string              `bun:"ledger,type:varchar"`
	IsSource            bool                `bun:"is_source,type:bool"`
	Account             string              `bun:"account_address,type:varchar"`
	AccountAddressArray []string            `bun:"account_address_array,type:jsonb"`
	Amount              *bunpaginate.BigInt `bun:"amount,type:numeric"`
	Asset               string              `bun:"asset,type:varchar"`
	TransactionSeq      int                 `bun:"transactions_seq,type:bigint"`
	AccountSeq          int                 `bun:"accounts_seq,type:bigint,scanonly"`
	InsertionDate       time.Time           `bun:"insertion_date,type:timestamp"`
	EffectiveDate       time.Time           `bun:"effective_date,type:timestamp"`
}

func (s *Store) InsertMoves(ctx context.Context, moves ...ledger.Move) error {

	mappedMoves := pointer.For(Map(moves, func(from ledger.Move) Move {
		return Move{
			Ledger:              s.ledgerName,
			IsSource:            from.IsSource,
			Account:             from.Account,
			AccountAddressArray: strings.Split(from.Account, ":"),
			Amount:              (*bunpaginate.BigInt)(from.Amount),
			Asset:               from.Asset,
			InsertionDate:       from.InsertedAt,
			TransactionSeq:      from.TransactionSeq,
			EffectiveDate:       from.EffectiveDate,
		}
	}))

	// todo: rewrite, we just need to basically insert moves
	_, err := s.db.NewInsert().
		With("_rows", s.db.NewValues(mappedMoves)).
		//todo: we should already have the sequence when using UpsertAccount
		With("_account_sequences",
			s.db.NewSelect().
				Column("seq", "address").
				ModelTableExpr(s.PrefixWithBucketUsingModel(Account{})).
				Join("join _rows on _rows.account_address = address and _rows.ledger = accounts.ledger"),
		).
		With("_computed_rows",
			s.db.NewSelect().
				ColumnExpr("_rows.*").
				ColumnExpr("_account_sequences.seq as accounts_seq").
				//ColumnExpr("("+
				//	"coalesce(((last_move_by_seq.post_commit_volumes).inputs), 0) + case when is_source then 0 else amount end, "+
				//	"coalesce(((last_move_by_seq.post_commit_volumes).outputs), 0) + case when is_source then amount else 0 end"+
				//	")::"+s.PrefixWithBucket("volumes")+" as post_commit_volumes").
				Join("join _account_sequences on _account_sequences.address = address").
				//Join("left join lateral ("+
				//	s.db.NewSelect().
				//		ColumnExpr("last_move.post_commit_volumes").
				//		ModelTableExpr(s.PrefixWithBucketUsingModel(Move{})+" as last_move").
				//		Where("_rows.account_address = last_move.account_address").
				//		Where("_rows.asset = last_move.asset").
				//		Where("_rows.ledger = last_move.ledger").
				//		Order("seq desc").
				//		Limit(1).
				//		String()+
				//	") last_move_by_seq on true").
				Table("_rows"),
		).
		Model(&Move{}).
		Column(
			"ledger",
			"is_source",
			"account_address",
			"account_address_array",
			"amount",
			"asset",
			"transactions_seq",
			"insertion_date",
			"effective_date",
			"accounts_seq",
			//"post_commit_volumes",
		).
		ModelTableExpr(s.PrefixWithBucketUsingModel(Move{})).
		Table("_computed_rows").
		Exec(ctx)
	return err
}
