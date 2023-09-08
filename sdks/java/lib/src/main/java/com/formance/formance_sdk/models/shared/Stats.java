/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;


public class Stats {
    @JsonProperty("accounts")
    public Long accounts;

    public Stats withAccounts(Long accounts) {
        this.accounts = accounts;
        return this;
    }
    
    @JsonProperty("transactions")
    public Long transactions;

    public Stats withTransactions(Long transactions) {
        this.transactions = transactions;
        return this;
    }
    
    public Stats(@JsonProperty("accounts") Long accounts, @JsonProperty("transactions") Long transactions) {
        this.accounts = accounts;
        this.transactions = transactions;
  }
}
