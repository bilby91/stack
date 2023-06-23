<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack;

class Wallets 
{

	private SDKConfiguration $sdkConfiguration;

	/**
	 * @param SDKConfiguration $sdkConfig
	 */
	public function __construct(SDKConfiguration $sdkConfig)
	{
		$this->sdkConfiguration = $sdkConfig;
	}
	
    /**
     * Confirm a hold
     * 
     * @param \formance\stack\Models\Operations\ConfirmHoldRequest $request
     * @return \formance\stack\Models\Operations\ConfirmHoldResponse
     */
	public function confirmHold(
        \formance\stack\Models\Operations\ConfirmHoldRequest $request,
    ): \formance\stack\Models\Operations\ConfirmHoldResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/holds/{hold_id}/confirm', \formance\stack\Models\Operations\ConfirmHoldRequest::class, $request);
        
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, "confirmHoldRequest", "json");
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('POST', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\ConfirmHoldResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 204) {
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Create a balance
     * 
     * @param \formance\stack\Models\Operations\CreateBalanceRequest $request
     * @return \formance\stack\Models\Operations\CreateBalanceResponse
     */
	public function createBalance(
        \formance\stack\Models\Operations\CreateBalanceRequest $request,
    ): \formance\stack\Models\Operations\CreateBalanceResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/wallets/{id}/balances', \formance\stack\Models\Operations\CreateBalanceRequest::class, $request);
        
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, "createBalanceRequest", "json");
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = 'application/json;q=1, application/json;q=0';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('POST', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\CreateBalanceResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 201) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->createBalanceResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\CreateBalanceResponse', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Create a new wallet
     * 
     * @param \formance\stack\Models\Shared\CreateWalletRequest $request
     * @return \formance\stack\Models\Operations\CreateWalletResponse
     */
	public function createWallet(
        \formance\stack\Models\Shared\CreateWalletRequest $request,
    ): \formance\stack\Models\Operations\CreateWalletResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/wallets');
        
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, "request", "json");
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = 'application/json;q=1, application/json;q=0';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('POST', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\CreateWalletResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 201) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->createWalletResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\CreateWalletResponse', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Credit a wallet
     * 
     * @param \formance\stack\Models\Operations\CreditWalletRequest $request
     * @return \formance\stack\Models\Operations\CreditWalletResponse
     */
	public function creditWallet(
        \formance\stack\Models\Operations\CreditWalletRequest $request,
    ): \formance\stack\Models\Operations\CreditWalletResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/wallets/{id}/credit', \formance\stack\Models\Operations\CreditWalletRequest::class, $request);
        
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, "creditWalletRequest", "json");
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('POST', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\CreditWalletResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 204) {
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Debit a wallet
     * 
     * @param \formance\stack\Models\Operations\DebitWalletRequest $request
     * @return \formance\stack\Models\Operations\DebitWalletResponse
     */
	public function debitWallet(
        \formance\stack\Models\Operations\DebitWalletRequest $request,
    ): \formance\stack\Models\Operations\DebitWalletResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/wallets/{id}/debit', \formance\stack\Models\Operations\DebitWalletRequest::class, $request);
        
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, "debitWalletRequest", "json");
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = 'application/json;q=1, application/json;q=0';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('POST', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\DebitWalletResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 201) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->debitWalletResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\DebitWalletResponse', 'json');
            }
        }
        else if ($httpResponse->getStatusCode() === 204) {
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Get detailed balance
     * 
     * @param \formance\stack\Models\Operations\GetBalanceRequest $request
     * @return \formance\stack\Models\Operations\GetBalanceResponse
     */
	public function getBalance(
        \formance\stack\Models\Operations\GetBalanceRequest $request,
    ): \formance\stack\Models\Operations\GetBalanceResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/wallets/{id}/balances/{balanceName}', \formance\stack\Models\Operations\GetBalanceRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json;q=1, application/json;q=0';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\GetBalanceResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->getBalanceResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\GetBalanceResponse', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Get a hold
     * 
     * @param \formance\stack\Models\Operations\GetHoldRequest $request
     * @return \formance\stack\Models\Operations\GetHoldResponse
     */
	public function getHold(
        \formance\stack\Models\Operations\GetHoldRequest $request,
    ): \formance\stack\Models\Operations\GetHoldResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/holds/{holdID}', \formance\stack\Models\Operations\GetHoldRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json;q=1, application/json;q=0';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\GetHoldResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->getHoldResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\GetHoldResponse', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Get all holds for a wallet
     * 
     * @param \formance\stack\Models\Operations\GetHoldsRequest $request
     * @return \formance\stack\Models\Operations\GetHoldsResponse
     */
	public function getHolds(
        \formance\stack\Models\Operations\GetHoldsRequest $request,
    ): \formance\stack\Models\Operations\GetHoldsResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/holds');
        
        $options = ['http_errors' => false];
        $options = array_merge_recursive($options, Utils\Utils::getQueryParams(\formance\stack\Models\Operations\GetHoldsRequest::class, $request, null));
        $options['headers']['Accept'] = 'application/json;q=1, application/json;q=0';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\GetHoldsResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->getHoldsResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\GetHoldsResponse', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * getTransactions
     * 
     * @param \formance\stack\Models\Operations\GetTransactionsRequest $request
     * @return \formance\stack\Models\Operations\GetTransactionsResponse
     */
	public function getTransactions(
        \formance\stack\Models\Operations\GetTransactionsRequest $request,
    ): \formance\stack\Models\Operations\GetTransactionsResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/transactions');
        
        $options = ['http_errors' => false];
        $options = array_merge_recursive($options, Utils\Utils::getQueryParams(\formance\stack\Models\Operations\GetTransactionsRequest::class, $request, null));
        $options['headers']['Accept'] = 'application/json;q=1, application/json;q=0';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\GetTransactionsResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->getTransactionsResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\GetTransactionsResponse', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Get a wallet
     * 
     * @param \formance\stack\Models\Operations\GetWalletRequest $request
     * @return \formance\stack\Models\Operations\GetWalletResponse
     */
	public function getWallet(
        \formance\stack\Models\Operations\GetWalletRequest $request,
    ): \formance\stack\Models\Operations\GetWalletResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/wallets/{id}', \formance\stack\Models\Operations\GetWalletRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json;q=1, application/json;q=0';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\GetWalletResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->getWalletResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\GetWalletResponse', 'json');
            }
        }
        else if ($httpResponse->getStatusCode() === 404) {
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Get wallet summary
     * 
     * @param \formance\stack\Models\Operations\GetWalletSummaryRequest $request
     * @return \formance\stack\Models\Operations\GetWalletSummaryResponse
     */
	public function getWalletSummary(
        \formance\stack\Models\Operations\GetWalletSummaryRequest $request,
    ): \formance\stack\Models\Operations\GetWalletSummaryResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/wallets/{id}/summary', \formance\stack\Models\Operations\GetWalletSummaryRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json;q=1, application/json;q=0';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\GetWalletSummaryResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->getWalletSummaryResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\GetWalletSummaryResponse', 'json');
            }
        }
        else if ($httpResponse->getStatusCode() === 404) {
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * List balances of a wallet
     * 
     * @param \formance\stack\Models\Operations\ListBalancesRequest $request
     * @return \formance\stack\Models\Operations\ListBalancesResponse
     */
	public function listBalances(
        \formance\stack\Models\Operations\ListBalancesRequest $request,
    ): \formance\stack\Models\Operations\ListBalancesResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/wallets/{id}/balances', \formance\stack\Models\Operations\ListBalancesRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\ListBalancesResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->listBalancesResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ListBalancesResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * List all wallets
     * 
     * @param \formance\stack\Models\Operations\ListWalletsRequest $request
     * @return \formance\stack\Models\Operations\ListWalletsResponse
     */
	public function listWallets(
        \formance\stack\Models\Operations\ListWalletsRequest $request,
    ): \formance\stack\Models\Operations\ListWalletsResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/wallets');
        
        $options = ['http_errors' => false];
        $options = array_merge_recursive($options, Utils\Utils::getQueryParams(\formance\stack\Models\Operations\ListWalletsRequest::class, $request, null));
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\ListWalletsResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->listWalletsResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ListWalletsResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Update a wallet
     * 
     * @param \formance\stack\Models\Operations\UpdateWalletRequest $request
     * @return \formance\stack\Models\Operations\UpdateWalletResponse
     */
	public function updateWallet(
        \formance\stack\Models\Operations\UpdateWalletRequest $request,
    ): \formance\stack\Models\Operations\UpdateWalletResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/wallets/{id}', \formance\stack\Models\Operations\UpdateWalletRequest::class, $request);
        
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, "requestBody", "json");
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('PATCH', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\UpdateWalletResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 204) {
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Cancel a hold
     * 
     * @param \formance\stack\Models\Operations\VoidHoldRequest $request
     * @return \formance\stack\Models\Operations\VoidHoldResponse
     */
	public function voidHold(
        \formance\stack\Models\Operations\VoidHoldRequest $request,
    ): \formance\stack\Models\Operations\VoidHoldResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/holds/{hold_id}/void', \formance\stack\Models\Operations\VoidHoldRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('POST', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\VoidHoldResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 204) {
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Get server info
     * 
     * @return \formance\stack\Models\Operations\WalletsgetServerInfoResponse
     */
	public function walletsgetServerInfo(
    ): \formance\stack\Models\Operations\WalletsgetServerInfoResponse
    {
        $baseUrl = Utils\Utils::templateUrl($this->sdkConfiguration->getServerUrl(), $this->sdkConfiguration->getServerDefaults());
        $url = Utils\Utils::generateUrl($baseUrl, '/api/wallets/_info');
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json;q=1, application/json;q=0';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s %s', $this->sdkConfiguration->language, $this->sdkConfiguration->sdkVersion, $this->sdkConfiguration->genVersion, $this->sdkConfiguration->openapiDocVersion);
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\WalletsgetServerInfoResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->serverInfo = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ServerInfo', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->walletsErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\WalletsErrorResponse', 'json');
            }
        }

        return $response;
    }
}