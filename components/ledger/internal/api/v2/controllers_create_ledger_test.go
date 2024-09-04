package v2_test

import (
	ledger "github.com/formancehq/ledger/internal"
	"net/http"
	"net/http/httptest"
	"testing"

	v2 "github.com/formancehq/ledger/internal/api/v2"
	"github.com/formancehq/ledger/internal/opentelemetry/metrics"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestConfigureLedger(t *testing.T) {
	t.Parallel()

	type testCase struct {
		configuration ledger.Configuration
		name          string
	}

	testCases := []testCase{
		{
			name:          "nominal",
			configuration: ledger.Configuration{},
		},
		{
			name: "with alternative bucket",
			configuration: ledger.Configuration{
				Bucket: "bucket0",
			},
		},
		{
			name: "with metadata",
			configuration: ledger.Configuration{
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			b, _ := newTestingBackend(t, false)
			router := v2.NewRouter(b, nil, metrics.NewNoOpRegistry(), auth.NewNoAuth(), testing.Verbose())

			name := uuid.NewString()
			b.
				EXPECT().
				CreateLedger(gomock.Any(), name, testCase.configuration).
				Return(nil)

			req := httptest.NewRequest(http.MethodPost, "/"+name, api.Buffer(t, testCase.configuration))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
		})
	}
}
