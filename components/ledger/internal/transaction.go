package ledger

import (
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/time"
	"slices"
	"sort"

	"github.com/formancehq/stack/libs/go-libs/pointer"

	"github.com/pkg/errors"

	"github.com/formancehq/stack/libs/go-libs/metadata"
)

var (
	ErrNoPostings = errors.New("invalid payload: should contain either postings or script")
)

type Transactions struct {
	Transactions []TransactionData `json:"transactions"`
}

type TransactionData struct {
	Postings   Postings          `json:"postings"`
	Metadata   metadata.Metadata `json:"metadata"`
	Timestamp  time.Time         `json:"timestamp"`
	Reference  string            `json:"reference,omitempty"`
	InsertedAt time.Time         `json:"insertedAt,omitempty"`
}

func (data TransactionData) WithPostings(postings ...Posting) TransactionData {
	data.Postings = append(data.Postings, postings...)
	return data
}

func NewTransactionData() TransactionData {
	return TransactionData{
		Metadata: metadata.Metadata{},
	}
}

func (data TransactionData) Reverse(atEffectiveDate bool) TransactionData {
	ret := NewTransactionData().WithPostings(data.Postings.Reverse()...)
	if atEffectiveDate {
		ret = ret.WithDate(data.Timestamp)
	}

	return ret
}

func (data TransactionData) WithDate(now time.Time) TransactionData {
	data.Timestamp = now

	return data
}

func (data TransactionData) WithReference(ref string) TransactionData {
	data.Reference = ref

	return data
}

func (data TransactionData) WithInsertedAt(date time.Time) TransactionData {
	data.InsertedAt = date

	return data
}

func (data TransactionData) WithMetadata(m metadata.Metadata) TransactionData {
	data.Metadata = m

	return data
}

type Transaction struct {
	TransactionData
	ID       int  `json:"id"`
	Reverted bool `json:"reverted"`
	Seq      int  `json:"-"`
}

func (t Transaction) WithPostings(postings ...Posting) Transaction {
	t.TransactionData = t.TransactionData.WithPostings(postings...)
	return t
}

func (t Transaction) WithReference(ref string) Transaction {
	t.Reference = ref
	return t
}

func (t Transaction) WithDate(ts time.Time) Transaction {
	t.Timestamp = ts
	return t
}

func (t Transaction) WithID(id int) Transaction {
	t.ID = id
	return t
}

func (t Transaction) WithMetadata(m metadata.Metadata) Transaction {
	t.Metadata = m
	return t
}

func (t Transaction) GetMoves() Moves {
	ret := make([]Move, 0)
	for _, p := range t.Postings {
		ret = append(ret, []Move{
			{
				IsSource:       true,
				Account:        p.Source,
				Amount:         p.Amount,
				Asset:          p.Asset,
				InsertedAt:     t.InsertedAt,
				EffectiveDate:  t.Timestamp,
				TransactionSeq: t.Seq,
			},
			{
				IsSource:       false,
				Account:        p.Destination,
				Amount:         p.Amount,
				Asset:          p.Asset,
				InsertedAt:     t.InsertedAt,
				EffectiveDate:  t.Timestamp,
				TransactionSeq: t.Seq,
			},
		}...)
	}
	return ret
}

func NewTransaction() Transaction {
	return Transaction{
		TransactionData: NewTransactionData().
			WithDate(time.Now()),
	}
}

type ExpandedTransaction struct {
	Transaction
	PreCommitVolumes           AccountsAssetsVolumes `json:"preCommitVolumes,omitempty"`
	PostCommitVolumes          AccountsAssetsVolumes `json:"postCommitVolumes,omitempty"`
	PreCommitEffectiveVolumes  AccountsAssetsVolumes `json:"preCommitEffectiveVolumes,omitempty"`
	PostCommitEffectiveVolumes AccountsAssetsVolumes `json:"postCommitEffectiveVolumes,omitempty"`
}

func (t ExpandedTransaction) Base() Transaction {
	return t.Transaction
}

func (t ExpandedTransaction) AppendPosting(p Posting) {
	t.Postings = append(t.Postings, p)
}

func ExpandTransaction(tx *Transaction, preCommitVolumes AccountsAssetsVolumes) ExpandedTransaction {
	postCommitVolumes := preCommitVolumes.Copy()
	for _, posting := range tx.Postings {
		preCommitVolumes.AddInput(posting.Destination, posting.Asset, Zero)
		preCommitVolumes.AddOutput(posting.Source, posting.Asset, Zero)
		postCommitVolumes.AddOutput(posting.Source, posting.Asset, posting.Amount)
		postCommitVolumes.AddInput(posting.Destination, posting.Asset, posting.Amount)
	}
	return ExpandedTransaction{
		Transaction:       *tx,
		PreCommitVolumes:  preCommitVolumes,
		PostCommitVolumes: postCommitVolumes,
	}
}

type TransactionRequest struct {
	Postings  Postings          `json:"postings"`
	Script    ScriptV1          `json:"script"`
	Timestamp time.Time         `json:"timestamp"`
	Reference string            `json:"reference"`
	Metadata  metadata.Metadata `json:"metadata" swaggertype:"object"`
}

func (req *TransactionRequest) ToRunScript() *RunScript {

	if len(req.Postings) > 0 {
		txData := TransactionData{
			Postings:  req.Postings,
			Timestamp: req.Timestamp,
			Reference: req.Reference,
			Metadata:  req.Metadata,
		}

		return pointer.For(TxToScriptData(txData, false))
	}

	return &RunScript{
		Script:    req.Script.ToCore(),
		Timestamp: req.Timestamp,
		Reference: req.Reference,
		Metadata:  req.Metadata,
	}
}

type Moves []Move

func (m Moves) InvolvedAccounts() []string {
	accounts := collectionutils.Map(m, func(from Move) string {
		return from.Account
	})
	sort.Strings(accounts)
	slices.Compact(accounts)

	return accounts
}