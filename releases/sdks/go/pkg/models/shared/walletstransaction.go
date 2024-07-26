// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"time"
)

type WalletsTransaction struct {
	ID     int64   `json:"id"`
	Ledger *string `json:"ledger,omitempty"`
	// Metadata associated with the wallet.
	Metadata          map[string]string                   `json:"metadata"`
	PostCommitVolumes map[string]map[string]WalletsVolume `json:"postCommitVolumes,omitempty"`
	Postings          []Posting                           `json:"postings"`
	PreCommitVolumes  map[string]map[string]WalletsVolume `json:"preCommitVolumes,omitempty"`
	Reference         *string                             `json:"reference,omitempty"`
	Timestamp         time.Time                           `json:"timestamp"`
}

func (w WalletsTransaction) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(w, "", false)
}

func (w *WalletsTransaction) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &w, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *WalletsTransaction) GetID() int64 {
	if o == nil {
		return 0
	}
	return o.ID
}

func (o *WalletsTransaction) GetLedger() *string {
	if o == nil {
		return nil
	}
	return o.Ledger
}

func (o *WalletsTransaction) GetMetadata() map[string]string {
	if o == nil {
		return map[string]string{}
	}
	return o.Metadata
}

func (o *WalletsTransaction) GetPostCommitVolumes() map[string]map[string]WalletsVolume {
	if o == nil {
		return nil
	}
	return o.PostCommitVolumes
}

func (o *WalletsTransaction) GetPostings() []Posting {
	if o == nil {
		return []Posting{}
	}
	return o.Postings
}

func (o *WalletsTransaction) GetPreCommitVolumes() map[string]map[string]WalletsVolume {
	if o == nil {
		return nil
	}
	return o.PreCommitVolumes
}

func (o *WalletsTransaction) GetReference() *string {
	if o == nil {
		return nil
	}
	return o.Reference
}

func (o *WalletsTransaction) GetTimestamp() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.Timestamp
}
