// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package types

import (
	"fmt"

	"github.com/ericlagergren/decimal"
)

// MustNewDecimalFromString returns an instance of Decimal from a string
// Avoid using this function in production code.
func MustNewDecimalFromString(s string) *decimal.Big {
	d, ok := new(decimal.Big).SetString(s)
	if !ok {
		panic(fmt.Errorf("failed to parse string as decimal.Big"))
	}

	return d
}
