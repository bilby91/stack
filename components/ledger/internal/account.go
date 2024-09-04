package ledger

import (
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/time"
	"github.com/uptrace/bun"
)

const (
	WORLD = "world"
)

type Account struct {
	// todo: clean models by keeping all related bun specificity in the storage package
	bun.BaseModel `bun:"table:accounts,alias:accounts"`

	Address    string            `json:"address,type:varchar"`
	Metadata   metadata.Metadata `json:"metadata,type:jsonb"`
	FirstUsage time.Time         `json:"-" bun:"first_usage,type:timestamp"`
	InsertionDate time.Time         `bun:"insertion_date,type:timestamp"`
	UpdatedAt     time.Time         `bun:"updated_at,type:timestamp"`
}

func (a Account) copy() Account {
	a.Metadata = a.Metadata.Copy()
	return a
}

func NewAccount(address string) Account {
	return Account{
		Address:  address,
		Metadata: metadata.Metadata{},
	}
}

type ExpandedAccount struct {
	Account          `bun:",extend"`
	Volumes          VolumesByAssets `json:"volumes,omitempty" bun:"volumes,type:jsonb"`
	EffectiveVolumes VolumesByAssets `json:"effectiveVolumes,omitempty" bun:"effective_volumes,type:jsonb"`
}

func NewExpandedAccount(address string) ExpandedAccount {
	return ExpandedAccount{
		Account: Account{
			Address:  address,
			Metadata: metadata.Metadata{},
		},
		Volumes: map[string]*Volumes{},
	}
}

func (v ExpandedAccount) Copy() ExpandedAccount {
	v.Account = v.Account.copy()
	v.Volumes = v.Volumes.copy()
	return v
}
