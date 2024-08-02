package workflow

import (
	"go.temporal.io/sdk/client"
)

type Workflow struct {
	taskQueue      string
	temporalClient client.Client
}
