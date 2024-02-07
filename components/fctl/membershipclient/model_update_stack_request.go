/*
Membership API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package membershipclient

import (
	"encoding/json"
)

// checks if the UpdateStackRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UpdateStackRequest{}

// UpdateStackRequest struct for UpdateStackRequest
type UpdateStackRequest struct {
	// Stack name
	Name string `json:"name"`
	Metadata map[string]string `json:"metadata"`
}

// NewUpdateStackRequest instantiates a new UpdateStackRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpdateStackRequest(name string, metadata map[string]string) *UpdateStackRequest {
	this := UpdateStackRequest{}
	this.Name = name
	this.Metadata = metadata
	return &this
}

// NewUpdateStackRequestWithDefaults instantiates a new UpdateStackRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpdateStackRequestWithDefaults() *UpdateStackRequest {
	this := UpdateStackRequest{}
	return &this
}

// GetName returns the Name field value
func (o *UpdateStackRequest) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *UpdateStackRequest) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *UpdateStackRequest) SetName(v string) {
	o.Name = v
}

// GetMetadata returns the Metadata field value
func (o *UpdateStackRequest) GetMetadata() map[string]string {
	if o == nil {
		var ret map[string]string
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *UpdateStackRequest) GetMetadataOk() (*map[string]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *UpdateStackRequest) SetMetadata(v map[string]string) {
	o.Metadata = v
}

func (o UpdateStackRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UpdateStackRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["metadata"] = o.Metadata
	return toSerialize, nil
}

type NullableUpdateStackRequest struct {
	value *UpdateStackRequest
	isSet bool
}

func (v NullableUpdateStackRequest) Get() *UpdateStackRequest {
	return v.value
}

func (v *NullableUpdateStackRequest) Set(val *UpdateStackRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableUpdateStackRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableUpdateStackRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUpdateStackRequest(val *UpdateStackRequest) *NullableUpdateStackRequest {
	return &NullableUpdateStackRequest{value: val, isSet: true}
}

func (v NullableUpdateStackRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUpdateStackRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


