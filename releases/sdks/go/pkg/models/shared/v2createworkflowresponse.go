// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type V2CreateWorkflowResponse struct {
	Data V2Workflow `json:"data"`
}

func (o *V2CreateWorkflowResponse) GetData() V2Workflow {
	if o == nil {
		return V2Workflow{}
	}
	return o.Data
}
