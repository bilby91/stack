<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class TransferRequest
{
    /**
     *
     * @var int $amount
     */
    #[\JMS\Serializer\Annotation\SerializedName('amount')]
    public int $amount;

    /**
     *
     * @var string $asset
     */
    #[\JMS\Serializer\Annotation\SerializedName('asset')]
    public string $asset;

    /**
     *
     * @var string $destination
     */
    #[\JMS\Serializer\Annotation\SerializedName('destination')]
    public string $destination;

    /**
     *
     * @var ?string $source
     */
    #[\JMS\Serializer\Annotation\SerializedName('source')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $source = null;

    /**
     * @param  ?int  $amount
     * @param  ?string  $asset
     * @param  ?string  $destination
     * @param  ?string  $source
     */
    public function __construct(?int $amount = null, ?string $asset = null, ?string $destination = null, ?string $source = null)
    {
        $this->amount = $amount;
        $this->asset = $asset;
        $this->destination = $destination;
        $this->source = $source;
    }
}