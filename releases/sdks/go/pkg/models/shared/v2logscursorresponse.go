// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type V2LogsCursorResponseCursor struct {
	Data     []V2Log `json:"data"`
	HasMore  bool    `json:"hasMore"`
	Next     *string `json:"next,omitempty"`
	PageSize int64   `json:"pageSize"`
	Previous *string `json:"previous,omitempty"`
}

func (o *V2LogsCursorResponseCursor) GetData() []V2Log {
	if o == nil {
		return []V2Log{}
	}
	return o.Data
}

func (o *V2LogsCursorResponseCursor) GetHasMore() bool {
	if o == nil {
		return false
	}
	return o.HasMore
}

func (o *V2LogsCursorResponseCursor) GetNext() *string {
	if o == nil {
		return nil
	}
	return o.Next
}

func (o *V2LogsCursorResponseCursor) GetPageSize() int64 {
	if o == nil {
		return 0
	}
	return o.PageSize
}

func (o *V2LogsCursorResponseCursor) GetPrevious() *string {
	if o == nil {
		return nil
	}
	return o.Previous
}

type V2LogsCursorResponse struct {
	Cursor V2LogsCursorResponseCursor `json:"cursor"`
}

func (o *V2LogsCursorResponse) GetCursor() V2LogsCursorResponseCursor {
	if o == nil {
		return V2LogsCursorResponseCursor{}
	}
	return o.Cursor
}
