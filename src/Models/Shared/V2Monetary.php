<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2Monetary
{
    /**
     * The amount of the monetary value.
     *
     * @var int $amount
     */
    #[\JMS\Serializer\Annotation\SerializedName('amount')]
    public int $amount;

    /**
     * The asset of the monetary value.
     *
     * @var string $asset
     */
    #[\JMS\Serializer\Annotation\SerializedName('asset')]
    public string $asset;

    /**
     * @param  ?int  $amount
     * @param  ?string  $asset
     */
    public function __construct(?int $amount = null, ?string $asset = null)
    {
        $this->amount = $amount;
        $this->asset = $asset;
    }
}