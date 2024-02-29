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
	public string $openapiDocVersion = 'v2.0.0-rc.14';
	public string $sdkVersion = 'v2.0.0-rc.14';
	public string $genVersion = '2.237.2';
	public string $userAgent = 'speakeasy-sdk/php v2.0.0-rc.14 2.237.2 v2.0.0-rc.14 formance-sdk-php';
	

	public function getServerUrl(): string
	{
		
		if ($this->serverUrl !== '') {
			return $this->serverUrl;
		};
		return SDK::SERVERS[$this->serverIndex];
	}
	
}