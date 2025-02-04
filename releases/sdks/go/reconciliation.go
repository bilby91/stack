// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package v2

type Reconciliation struct {
	V1 *FormanceReconciliationV1

	sdkConfiguration sdkConfiguration
}

func newReconciliation(sdkConfig sdkConfiguration) *Reconciliation {
	return &Reconciliation{
		sdkConfiguration: sdkConfig,
		V1:               newFormanceReconciliationV1(sdkConfig),
	}
}
