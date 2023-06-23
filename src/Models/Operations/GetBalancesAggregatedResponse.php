<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;


class GetBalancesAggregatedResponse
{
    /**
     * OK
     * 
     * @var ?\formance\stack\Models\Shared\AggregateBalancesResponse $aggregateBalancesResponse
     */
	
    public ?\formance\stack\Models\Shared\AggregateBalancesResponse $aggregateBalancesResponse = null;
    
	
    public string $contentType;
    
    /**
     * Error
     * 
     * @var ?\formance\stack\Models\Shared\ErrorResponse $errorResponse
     */
	
    public ?\formance\stack\Models\Shared\ErrorResponse $errorResponse = null;
    
	
    public int $statusCode;
    
	
    public ?\Psr\Http\Message\ResponseInterface $rawResponse = null;
    
	public function __construct()
	{
		$this->aggregateBalancesResponse = null;
		$this->contentType = "";
		$this->errorResponse = null;
		$this->statusCode = 0;
		$this->rawResponse = null;
	}
}
