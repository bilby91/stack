package ledger

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrBucketAlreadyExists = errors.New("bucket already exists")
	ErrStoreAlreadyExists  = errors.New("store already exists")
	ErrStoreNotFound       = errors.New("store not found")
)

type ErrInvalidQuery struct {
	msg string
}

func (e *ErrInvalidQuery) Error() string {
	return e.msg
}

func (e *ErrInvalidQuery) Is(err error) bool {
	_, ok := err.(*ErrInvalidQuery)
	return ok
}

func newErrInvalidQuery(msg string, args ...any) *ErrInvalidQuery {
	return &ErrInvalidQuery{
		msg: fmt.Sprintf(msg, args...),
	}
}

func IsErrInvalidQuery(err error) bool {
	return errors.Is(err, &ErrInvalidQuery{})
}
