<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;


class RunWorkflowResponse
{
	
    public string $contentType;
    
    /**
     * General error
     * 
     * @var ?\formance\stack\Models\Shared\Error $error
     */
	
    public ?\formance\stack\Models\Shared\Error $error = null;
    
    /**
     * The workflow instance
     * 
     * @var ?\formance\stack\Models\Shared\RunWorkflowResponse $runWorkflowResponse
     */
	
    public ?\formance\stack\Models\Shared\RunWorkflowResponse $runWorkflowResponse = null;
    
	
    public int $statusCode;
    
	
    public ?\Psr\Http\Message\ResponseInterface $rawResponse = null;
    
	public function __construct()
	{
		$this->contentType = "";
		$this->error = null;
		$this->runWorkflowResponse = null;
		$this->statusCode = 0;
		$this->rawResponse = null;
	}
}
