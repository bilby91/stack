package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListConnectors(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	backend := NewMockBackend(ctrl)

	api := newAPI(t, backend)
	srv := httptest.NewServer(api.Router())
	t.Cleanup(srv.Close)

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/connectors", nil)
	require.NoError(t, err)

	connectors := []ingester.Connector{
		ingester.NewConnector(ingester.NewConnectorConfiguration("connector1", json.RawMessage(`{}`))),
		ingester.NewConnector(ingester.NewConnectorConfiguration("connector2", json.RawMessage(`{}`))),
	}
	backend.EXPECT().
		ListConnectors(gomock.Any()).
		Return(&bunpaginate.Cursor[ingester.Connector]{
			Data: connectors,
		}, nil)

	rsp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, "application/json", rsp.Header.Get("Content-Type"))
}
