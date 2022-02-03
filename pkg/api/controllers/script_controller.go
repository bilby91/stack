package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/numary/ledger/pkg/core"
	"github.com/numary/ledger/pkg/ledger"
)

type ScriptResponse struct {
	ErrorResponse
	Link string                           `json:"details,omitempty"`
	Txs  []ledger.CommitTransactionResult `json:"txs,omitempty"`
}

func EncodeLink(err error) string {
	errStr := err.Error()
	errStr = strings.ReplaceAll(errStr, "\n", "\r\n")
	payload, err := json.Marshal(gin.H{
		"error": errStr,
	})
	if err != nil {
		panic(err)
	}
	payloadB64 := base64.StdEncoding.EncodeToString(payload)
	return fmt.Sprintf("https://play.numscript.org/?payload=%v", payloadB64)
}

// ScriptController -
type ScriptController struct {
	BaseController
}

// NewScriptController -
func NewScriptController() ScriptController {
	return ScriptController{}
}

func (ctl *ScriptController) PostScript(c *gin.Context) {
	l, _ := c.Get("ledger")

	var script core.Script
	c.ShouldBind(&script)

	txs, err := l.(*ledger.Ledger).Execute(c.Request.Context(), script)

	res := ScriptResponse{
		Txs: txs,
	}

	if err != nil {
		res.ErrorResponse = ErrorResponse{
			ErrorCode:    ErrInternal,
			ErrorMessage: err.Error(),
		}
		res.Link = EncodeLink(err)
	}

	c.JSON(http.StatusOK, res)
}
