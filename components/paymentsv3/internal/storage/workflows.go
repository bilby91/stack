package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type workflow struct {
	bun.BaseModel `bun:"table:workflows"`

	// Mandatory fields
	ID          string             `bun:"id,pk,type:text,notnull"`
	ConnectorID models.ConnectorID `bun:"connector_id,type:character varying,notnull"`
	CreatedAt   time.Time          `bun:"created_at,type:timestamp without time zone,notnull"`
	Capability  models.Capability  `bun:"capability,type:text,notnull"`

	// Optional fields with default
	// c.f. https://bun.uptrace.dev/guide/models.html#default
	Metadata map[string]string `bun:"metadata,type:jsonb,nullzero,notnull,default:'{}'"`
}

func (s *store) UpsertWorkflow(ctx context.Context, workflow models.Workflow) error {
	toInsert := fromWorkflowModel(workflow)

	_, err := s.db.NewInsert().
		Model(&toInsert).
		On("CONFLICT (id) DO NOTHING").
		Exec(ctx)

	return e("failed to insert workflow", err)
}

func (s *store) GetWorflow(ctx context.Context, id string) (*models.Workflow, error) {
	var w workflow

	err := s.db.NewSelect().
		Model(&w).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to fetch workflow", err)
	}

	ret := toWorkflowModel(w)
	return &ret, nil
}

func (s *store) DeleteWorkflowsFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error {
	_, err := s.db.NewDelete().
		Model((*workflow)(nil)).
		Where("connector_id = ?", connectorID).
		Exec(ctx)

	return e("failed to delete workflow", err)
}

type WorkflowQuery struct{}

type ListWorkflowsQuery bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[WorkflowQuery]]

func NewListWorkflowsQuery(opts bunpaginate.PaginatedQueryOptions[WorkflowQuery]) ListWorkflowsQuery {
	return ListWorkflowsQuery{
		Order:    bunpaginate.OrderAsc,
		PageSize: opts.PageSize,
		Options:  opts,
	}
}

func (s *store) workflowsQueryContext(qb query.Builder) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "connector_id",
			key == "capability":
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrValidation, fmt.Sprintf("'%s' column can only be used with $match", key))
			}
			return fmt.Sprintf("%s = ?", key), []any{value}, nil
		case metadataRegex.Match([]byte(key)):
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrValidation, "'metadata' column can only be used with $match")
			}
			match := metadataRegex.FindAllStringSubmatch(key, 3)

			key := "metadata"
			return key + " @> ?", []any{map[string]any{
				match[0][1]: value,
			}}, nil
		default:
			return "", nil, errors.Wrap(ErrValidation, fmt.Sprintf("unknown key '%s' when building query", key))
		}
	}))
}

func (s *store) ListWorkflows(ctx context.Context, q ListWorkflowsQuery) (*bunpaginate.Cursor[models.Workflow], error) {
	var (
		where string
		args  []any
		err   error
	)
	if q.Options.QueryBuilder != nil {
		where, args, err = s.workflowsQueryContext(q.Options.QueryBuilder)
		if err != nil {
			return nil, err
		}
	}

	cursor, err := paginateWithOffset[bunpaginate.PaginatedQueryOptions[WorkflowQuery], workflow](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[WorkflowQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			if where != "" {
				query = query.Where(where, args...)
			}

			query = query.Order("created_at DESC")

			return query
		},
	)
	if err != nil {
		return nil, e("failed to fetch workflows", err)
	}

	workflows := make([]models.Workflow, 0, len(cursor.Data))
	for _, w := range cursor.Data {
		workflows = append(workflows, toWorkflowModel(w))
	}

	return &bunpaginate.Cursor[models.Workflow]{
		PageSize: cursor.PageSize,
		HasMore:  cursor.HasMore,
		Previous: cursor.Previous,
		Next:     cursor.Next,
		Data:     workflows,
	}, nil
}

func fromWorkflowModel(from models.Workflow) workflow {
	return workflow{
		ID:          from.ID,
		ConnectorID: from.ConnectorID,
		CreatedAt:   from.CreatedAt,
		Capability:  from.Capability,
		Metadata:    from.Metadata,
	}
}

func toWorkflowModel(to workflow) models.Workflow {
	return models.Workflow{
		ID:          to.ID,
		ConnectorID: to.ConnectorID,
		CreatedAt:   to.CreatedAt,
		Capability:  to.Capability,
		Metadata:    to.Metadata,
	}
}
