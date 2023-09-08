<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


/**
 * Error - General error
 * 
 * @package formance\stack\Models\Shared
 * @access public
 */
class Error
{
	#[\JMS\Serializer\Annotation\SerializedName('errorCode')]
    #[\JMS\Serializer\Annotation\Type('enum<formance\stack\Models\Shared\ErrorErrorCode>')]
    public ErrorErrorCode $errorCode;
    
	#[\JMS\Serializer\Annotation\SerializedName('errorMessage')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $errorMessage;
    
	public function __construct()
	{
		$this->errorCode = \formance\stack\Models\Shared\ErrorErrorCode::Validation;
		$this->errorMessage = "";
	}
}
