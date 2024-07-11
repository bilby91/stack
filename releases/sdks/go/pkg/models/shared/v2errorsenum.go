// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type V2ErrorsEnum string

const (
	V2ErrorsEnumInternal          V2ErrorsEnum = "INTERNAL"
	V2ErrorsEnumInsufficientFund  V2ErrorsEnum = "INSUFFICIENT_FUND"
	V2ErrorsEnumValidation        V2ErrorsEnum = "VALIDATION"
	V2ErrorsEnumConflict          V2ErrorsEnum = "CONFLICT"
	V2ErrorsEnumCompilationFailed V2ErrorsEnum = "COMPILATION_FAILED"
	V2ErrorsEnumMetadataOverride  V2ErrorsEnum = "METADATA_OVERRIDE"
	V2ErrorsEnumNotFound          V2ErrorsEnum = "NOT_FOUND"
	V2ErrorsEnumRevertOccurring   V2ErrorsEnum = "REVERT_OCCURRING"
	V2ErrorsEnumAlreadyRevert     V2ErrorsEnum = "ALREADY_REVERT"
	V2ErrorsEnumNoPostings        V2ErrorsEnum = "NO_POSTINGS"
	V2ErrorsEnumLedgerNotFound    V2ErrorsEnum = "LEDGER_NOT_FOUND"
	V2ErrorsEnumImport            V2ErrorsEnum = "IMPORT"
)

func (e V2ErrorsEnum) ToPointer() *V2ErrorsEnum {
	return &e
}
func (e *V2ErrorsEnum) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "INTERNAL":
		fallthrough
	case "INSUFFICIENT_FUND":
		fallthrough
	case "VALIDATION":
		fallthrough
	case "CONFLICT":
		fallthrough
	case "COMPILATION_FAILED":
		fallthrough
	case "METADATA_OVERRIDE":
		fallthrough
	case "NOT_FOUND":
		fallthrough
	case "REVERT_OCCURRING":
		fallthrough
	case "ALREADY_REVERT":
		fallthrough
	case "NO_POSTINGS":
		fallthrough
	case "LEDGER_NOT_FOUND":
		fallthrough
	case "IMPORT":
		*e = V2ErrorsEnum(v)
		return nil
	default:
		return fmt.Errorf("invalid value for V2ErrorsEnum: %v", v)
	}
}
