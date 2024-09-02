package api

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/api"
)

func (a *API) listPipelines(w http.ResponseWriter, r *http.Request) {
	pipelines, err := a.backend.ListPipelines(r.Context())
	if err != nil {
		api.InternalServerError(w, r, err)
		return
	}

	api.RenderCursor(w, *pipelines)
}
