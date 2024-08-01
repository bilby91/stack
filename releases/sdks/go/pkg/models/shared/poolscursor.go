// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type PoolsCursorCursor struct {
	Data     []Pool  `json:"data"`
	HasMore  bool    `json:"hasMore"`
	Next     *string `json:"next,omitempty"`
	PageSize int64   `json:"pageSize"`
	Previous *string `json:"previous,omitempty"`
}

func (o *PoolsCursorCursor) GetData() []Pool {
	if o == nil {
		return []Pool{}
	}
	return o.Data
}

func (o *PoolsCursorCursor) GetHasMore() bool {
	if o == nil {
		return false
	}
	return o.HasMore
}

func (o *PoolsCursorCursor) GetNext() *string {
	if o == nil {
		return nil
	}
	return o.Next
}

func (o *PoolsCursorCursor) GetPageSize() int64 {
	if o == nil {
		return 0
	}
	return o.PageSize
}

func (o *PoolsCursorCursor) GetPrevious() *string {
	if o == nil {
		return nil
	}
	return o.Previous
}

type PoolsCursor struct {
	Cursor PoolsCursorCursor `json:"cursor"`
}

func (o *PoolsCursor) GetCursor() PoolsCursorCursor {
	if o == nil {
		return PoolsCursorCursor{}
	}
	return o.Cursor
}
