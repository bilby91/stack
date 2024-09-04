package ledger

import (
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/time"
)

type Configuration struct {
	Bucket   string            `json:"bucket"`
	Metadata metadata.Metadata `json:"metadata"`
}

type Ledger struct {
	Configuration
	Name string `json:"name"`
	AddedAt  time.Time         `json:"addedAt"`
	State    string            `json:"-"`
}
