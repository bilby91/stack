// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"time"
)

type TaskModulrDescriptor struct {
	AccountID *string `json:"accountID,omitempty"`
	Key       *string `json:"key,omitempty"`
	Name      *string `json:"name,omitempty"`
}

func (o *TaskModulrDescriptor) GetAccountID() *string {
	if o == nil {
		return nil
	}
	return o.AccountID
}

func (o *TaskModulrDescriptor) GetKey() *string {
	if o == nil {
		return nil
	}
	return o.Key
}

func (o *TaskModulrDescriptor) GetName() *string {
	if o == nil {
		return nil
	}
	return o.Name
}

type TaskModulrState struct {
}

type TaskModulr struct {
	ConnectorID string               `json:"connectorId"`
	CreatedAt   time.Time            `json:"createdAt"`
	Descriptor  TaskModulrDescriptor `json:"descriptor"`
	Error       *string              `json:"error,omitempty"`
	ID          string               `json:"id"`
	State       TaskModulrState      `json:"state"`
	Status      PaymentStatus        `json:"status"`
	UpdatedAt   time.Time            `json:"updatedAt"`
}

func (o *TaskModulr) GetConnectorID() string {
	if o == nil {
		return ""
	}
	return o.ConnectorID
}

func (o *TaskModulr) GetCreatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.CreatedAt
}

func (o *TaskModulr) GetDescriptor() TaskModulrDescriptor {
	if o == nil {
		return TaskModulrDescriptor{}
	}
	return o.Descriptor
}

func (o *TaskModulr) GetError() *string {
	if o == nil {
		return nil
	}
	return o.Error
}

func (o *TaskModulr) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *TaskModulr) GetState() TaskModulrState {
	if o == nil {
		return TaskModulrState{}
	}
	return o.State
}

func (o *TaskModulr) GetStatus() PaymentStatus {
	if o == nil {
		return PaymentStatus("")
	}
	return o.Status
}

func (o *TaskModulr) GetUpdatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.UpdatedAt
}
