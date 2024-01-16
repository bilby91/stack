<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack;

class SDKConfiguration
{
	public ?\GuzzleHttp\ClientInterface $defaultClient = null;
	public ?\GuzzleHttp\ClientInterface $securityClient = null;
	public ?Models\Shared\Security $security = null;
	public string $serverUrl = '';
	public int $serverIndex = 0;
	public string $language = 'php';
	public string $openapiDocVersion = 'v2.0.0-beta.9';
	public string $sdkVersion = 'v2.0.0-beta.9';
	public string $genVersion = '2.230.1';
	public string $userAgent = 'speakeasy-sdk/php v2.0.0-beta.9 2.230.1 v2.0.0-beta.9 formance-sdk-php';
	

	public function getServerUrl(): string
	{
		
		if ($this->serverUrl !== '') {
			return $this->serverUrl;
		};
		return SDK::SERVERS[$this->serverIndex];
	}
	
}