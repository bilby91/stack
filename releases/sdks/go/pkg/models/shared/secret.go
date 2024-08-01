// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type Secret struct {
	Clear      string         `json:"clear"`
	ID         string         `json:"id"`
	LastDigits string         `json:"lastDigits"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	Name       string         `json:"name"`
}

func (o *Secret) GetClear() string {
	if o == nil {
		return ""
	}
	return o.Clear
}

func (o *Secret) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *Secret) GetLastDigits() string {
	if o == nil {
		return ""
	}
	return o.LastDigits
}

func (o *Secret) GetMetadata() map[string]any {
	if o == nil {
		return nil
	}
	return o.Metadata
}

func (o *Secret) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}
