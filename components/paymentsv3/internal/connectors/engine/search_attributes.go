package engine

import (
	"github.com/formancehq/paymentsv3/internal/connectors/engine/workflow"
	"go.temporal.io/api/enums/v1"
)

var (
	SearchAttributes = map[string]enums.IndexedValueType{
		workflow.SearchAttributeWorkflowID: enums.INDEXED_VALUE_TYPE_TEXT,
	}
)
