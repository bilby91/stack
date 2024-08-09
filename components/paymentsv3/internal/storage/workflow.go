package storage

import (
	"context"
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type workflow struct {
	bun.BaseModel `bun:"table:workflows"`

	ConnectorID models.ConnectorID `bun:"connector_id,pk,type:character varying,notnull"`
	Workflow    json.RawMessage    `bun:"workflow,type:json,notnull"`
}

func (s *store) UpsertWorkflow(ctx context.Context, connectorID models.ConnectorID, tasks models.Workflow) error {
	payload, err := json.Marshal(&tasks)
	if err != nil {
		return errors.Wrap(err, "failed to marshal workflow")
	}

	workflow := workflow{
		ConnectorID: connectorID,
		Workflow:    payload,
	}

	_, err = s.db.NewInsert().
		Model(&workflow).
		On("CONFLICT (connector_id) DO UPDATE").
		Set("workflow = EXCLUDED.workflow").
		Exec(ctx)
	return e("failed to insert workflow", err)
}

func (s *store) GetWorkflow(ctx context.Context, connectorID models.ConnectorID) (*models.Workflow, error) {
	var workflow workflow

	err := s.db.NewSelect().
		Model(&workflow).
		Where("connector_id = ?", connectorID).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to fetch workflow", err)
	}

	var tasks models.Workflow
	if err := json.Unmarshal(workflow.Workflow, &tasks); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal workflow")
	}

	return &tasks, nil
}
