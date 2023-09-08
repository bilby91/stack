/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.net.http.HttpResponse;


public class WalletsgetServerInfoResponse {
    
    public String contentType;

    public WalletsgetServerInfoResponse withContentType(String contentType) {
        this.contentType = contentType;
        return this;
    }
    
    /**
     * Server information
     */
    
    public com.formance.formance_sdk.models.shared.ServerInfo serverInfo;

    public WalletsgetServerInfoResponse withServerInfo(com.formance.formance_sdk.models.shared.ServerInfo serverInfo) {
        this.serverInfo = serverInfo;
        return this;
    }
    
    
    public Integer statusCode;

    public WalletsgetServerInfoResponse withStatusCode(Integer statusCode) {
        this.statusCode = statusCode;
        return this;
    }
    
    
    public HttpResponse<byte[]> rawResponse;

    public WalletsgetServerInfoResponse withRawResponse(HttpResponse<byte[]> rawResponse) {
        this.rawResponse = rawResponse;
        return this;
    }
    
    /**
     * Error
     */
    
    public com.formance.formance_sdk.models.shared.WalletsErrorResponse walletsErrorResponse;

    public WalletsgetServerInfoResponse withWalletsErrorResponse(com.formance.formance_sdk.models.shared.WalletsErrorResponse walletsErrorResponse) {
        this.walletsErrorResponse = walletsErrorResponse;
        return this;
    }
    
    public WalletsgetServerInfoResponse(@JsonProperty("ContentType") String contentType, @JsonProperty("StatusCode") Integer statusCode) {
        this.contentType = contentType;
        this.statusCode = statusCode;
  }
}
