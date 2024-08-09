package storage

import (
	"context"
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/uptrace/bun"
)

type state struct {
	bun.BaseModel `bun:"table:states"`

	ID          models.StateID     `bun:"id,pk,type:character varying,notnull"`
	ConnectorID models.ConnectorID `bun:"connector_id,type:character varying,notnull"`
	State       json.RawMessage    `bun:"state,type:json,notnull"`
}

func (s *store) UpsertState(ctx context.Context, state models.State) error {
	toInsert := fromStateModels(state)

	_, err := s.db.NewInsert().
		Model(&toInsert).
		On("CONFLICT (id) DO UPDATE").
		Set("state = EXCLUDED.state").
		Exec(ctx)
	return e("failed to upsert state", err)
}

func (s *store) GetState(ctx context.Context, id models.StateID) (models.State, error) {
	var state state

	err := s.db.NewSelect().
		Model(&state).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return models.State{}, e("failed to get state", err)
	}

	res := toStateModels(state)
	return res, nil
}

func fromStateModels(from models.State) state {
	return state{
		ID:          from.ID,
		ConnectorID: from.ConnectorID,
		State:       from.State,
	}
}

func toStateModels(from state) models.State {
	return models.State{
		ID:          from.ID,
		ConnectorID: from.ConnectorID,
		State:       from.State,
	}
}
