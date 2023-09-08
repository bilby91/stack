/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.net.http.HttpResponse;


public class DebitWalletResponse {
    
    public String contentType;

    public DebitWalletResponse withContentType(String contentType) {
        this.contentType = contentType;
        return this;
    }
    
    /**
     * Wallet successfully debited as a pending hold
     */
    
    public com.formance.formance_sdk.models.shared.DebitWalletResponse debitWalletResponse;

    public DebitWalletResponse withDebitWalletResponse(com.formance.formance_sdk.models.shared.DebitWalletResponse debitWalletResponse) {
        this.debitWalletResponse = debitWalletResponse;
        return this;
    }
    
    
    public Integer statusCode;

    public DebitWalletResponse withStatusCode(Integer statusCode) {
        this.statusCode = statusCode;
        return this;
    }
    
    
    public HttpResponse<byte[]> rawResponse;

    public DebitWalletResponse withRawResponse(HttpResponse<byte[]> rawResponse) {
        this.rawResponse = rawResponse;
        return this;
    }
    
    /**
     * Error
     */
    
    public com.formance.formance_sdk.models.shared.WalletsErrorResponse walletsErrorResponse;

    public DebitWalletResponse withWalletsErrorResponse(com.formance.formance_sdk.models.shared.WalletsErrorResponse walletsErrorResponse) {
        this.walletsErrorResponse = walletsErrorResponse;
        return this;
    }
    
    public DebitWalletResponse(@JsonProperty("ContentType") String contentType, @JsonProperty("StatusCode") Integer statusCode) {
        this.contentType = contentType;
        this.statusCode = statusCode;
  }
}
