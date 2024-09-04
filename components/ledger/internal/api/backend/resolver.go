package backend

import (
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/platform/postgres"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"

	"github.com/pkg/errors"

	"github.com/formancehq/ledger/internal/opentelemetry/tracer"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/go-chi/chi/v5"
)

var (
	r  *rand.Rand
	mu sync.Mutex
)

const (
	ErrOutdatedSchema = "OUTDATED_SCHEMA"
)

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomTraceID(n int) string {
	mu.Lock()
	defer mu.Unlock()

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

func LedgerMiddleware(
	resolver Backend,
	excludePathFromSchemaCheck []string,
) func(handler http.Handler) http.Handler {

	mu := sync.RWMutex{}
	ledgers := make(map[string]Ledger, 0)
	upToDateLedgers := collectionutils.Set[string]{}

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			name := chi.URLParam(r, "ledger")
			if name == "" {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			ctx, span := tracer.Start(r.Context(), name)
			defer span.End()

			r = r.WithContext(ctx)

			loggerFields := map[string]any{
				"ledger": name,
			}
			if span.SpanContext().TraceID().IsValid() {
				loggerFields["trace-id"] = span.SpanContext().TraceID().String()
			} else {
				loggerFields["trace-id"] = randomTraceID(10)
			}

			r = r.WithContext(logging.ContextWithFields(r.Context(), loggerFields))

			var (
				l  Ledger
				ok bool
			)

			mu.RLock()
			if l, ok = ledgers[name]; ok {
				mu.RUnlock()
			} else {
				mu.RUnlock()
				mu.Lock()
				if l, ok = ledgers[name]; ok {
					mu.Unlock()
				} else {
					var err error
					l, err = resolver.GetLedgerController(r.Context(), name)
					if err != nil {
						switch {
						case postgres.IsNotFoundError(err):
							sharedapi.WriteErrorResponse(w, http.StatusNotFound, "LEDGER_NOT_FOUND", err)
						default:
							sharedapi.InternalServerError(w, r, err)
						}
						return
					}
					ledgers[name] = l
					mu.Unlock()
				}

				if !upToDateLedgers.Contains(name) {
					pathWithoutLedger := r.URL.Path[1:]
					nextSlash := strings.Index(pathWithoutLedger, "/")
					if nextSlash >= 0 {
						pathWithoutLedger = pathWithoutLedger[nextSlash:]
					} else {
						pathWithoutLedger = ""
					}

					excluded := false
					for _, path := range excludePathFromSchemaCheck {
						if pathWithoutLedger == path {
							excluded = true
							break
						}
					}

					if !excluded {
						isUpToDate, err := l.IsDatabaseUpToDate(ctx)
						if err != nil {
							sharedapi.BadRequest(w, sharedapi.ErrorInternal, err)
							return
						}
						if !isUpToDate {
							sharedapi.BadRequest(w, ErrOutdatedSchema, errors.New("You need to upgrade your ledger schema to the last version"))
							return
						}

						upToDateLedgers.Put(name)
					}
				}
			}

			handler.ServeHTTP(w, r.WithContext(ContextWithLedger(r.Context(), l)))
		})
	}
}
