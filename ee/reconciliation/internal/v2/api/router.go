package api

import (
	"github.com/formancehq/reconciliation/internal/v2/api/backend"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	b backend.Backend,
	r chi.Router,
) {
	r.Route("/v2", func(r chi.Router) {
		r.Route("/rules", func(r chi.Router) {
			r.Post("/", createRuleHandler(b))
			r.Get("/", listRulesHandler(b))
			r.Route("/{ruleID}", func(r chi.Router) {
				r.Get("/", getRuleHandler(b))
				r.Delete("/", deleteRuleHandler(b))
				// TODO(polo): soft delete rule
			})
		})

		r.Route("/policies", func(r chi.Router) {
			r.Post("/", createPolicyHandler(b))
			r.Get("/", listPoliciesHandler(b))
			r.Route("/{policyID}", func(r chi.Router) {
				r.Get("/", getPolicyHandler(b))
				// TODO(polo): SOFT DELETE POLICY
				r.Delete("/", deletePolicyHandler(b))
				r.Post("/enable", enablePolicyHandler(b))
				r.Post("/disable", disablePolicyHandler(b))
				r.Patch("/rules", updatePolicyRulesHandler(b))
			})
		})

		r.Route("/reconciliations", func(r chi.Router) {
			r.Get("/", listReconciliationsHandler(b))
			r.Post("/", createReconciliationHandler(b))
			r.Route("/{reconciliationID}", func(r chi.Router) {
				r.Get("/", getReconciliationHandler(b))
				r.Get("/details", getReconciliationDetails(b))
			})
		})
	})
}
