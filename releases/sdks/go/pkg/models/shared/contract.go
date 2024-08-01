// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type Expr struct {
}

type Contract struct {
	Account *string `json:"account,omitempty"`
	Expr    Expr    `json:"expr"`
}

func (o *Contract) GetAccount() *string {
	if o == nil {
		return nil
	}
	return o.Account
}

func (o *Contract) GetExpr() Expr {
	if o == nil {
		return Expr{}
	}
	return o.Expr
}
