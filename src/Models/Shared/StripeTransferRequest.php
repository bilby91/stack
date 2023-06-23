<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class StripeTransferRequest
{
	#[\JMS\Serializer\Annotation\SerializedName('amount')]
    #[\JMS\Serializer\Annotation\Type('int')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?int $amount = null;
    
	#[\JMS\Serializer\Annotation\SerializedName('asset')]
    #[\JMS\Serializer\Annotation\Type('string')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $asset = null;
    
	#[\JMS\Serializer\Annotation\SerializedName('destination')]
    #[\JMS\Serializer\Annotation\Type('string')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $destination = null;
    
    /**
     * A set of key/value pairs that you can attach to a transfer object.
     * 
     * It can be useful for storing additional information about the transfer in a structured format.
     * 
     * 
     * @var ?\formance\stack\Models\Shared\StripeTransferRequestMetadata $metadata
     */
	#[\JMS\Serializer\Annotation\SerializedName('metadata')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\StripeTransferRequestMetadata')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?StripeTransferRequestMetadata $metadata = null;
    
	public function __construct()
	{
		$this->amount = null;
		$this->asset = null;
		$this->destination = null;
		$this->metadata = null;
	}
}
