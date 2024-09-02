package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListPipelines(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	backend := NewMockBackend(ctrl)

	api := newAPI(t, backend)
	srv := httptest.NewServer(api.Router())
	t.Cleanup(srv.Close)

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/pipelines", nil)
	require.NoError(t, err)

	pipelines := []ingester.Pipeline{
		ingester.NewPipeline(ingester.NewPipelineConfiguration("module1", "connector1"), ingester.NewReadyState()),
		ingester.NewPipeline(ingester.NewPipelineConfiguration("module2", "connector2"), ingester.NewReadyState()),
	}
	backend.EXPECT().
		ListPipelines(gomock.Any()).
		Return(&bunpaginate.Cursor[ingester.Pipeline]{
			Data: pipelines,
		}, nil)

	rsp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, "application/json", rsp.Header.Get("Content-Type"))
}
