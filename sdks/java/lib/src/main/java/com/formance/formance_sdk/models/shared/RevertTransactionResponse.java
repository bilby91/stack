/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * RevertTransactionResponse - OK
 */

public class RevertTransactionResponse {
    @JsonProperty("data")
    public Transaction data;

    public RevertTransactionResponse withData(Transaction data) {
        this.data = data;
        return this;
    }
    
    public RevertTransactionResponse(@JsonProperty("data") Transaction data) {
        this.data = data;
  }
}
