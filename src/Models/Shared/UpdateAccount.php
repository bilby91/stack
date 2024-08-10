<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class UpdateAccount
{
    /**
     *
     * @var string $id
     */
    #[\JMS\Serializer\Annotation\SerializedName('id')]
    public string $id;

    /**
     *
     * @var string $ledger
     */
    #[\JMS\Serializer\Annotation\SerializedName('ledger')]
    public string $ledger;

    /**
     * $metadata
     *
     * @var array<string, string> $metadata
     */
    #[\JMS\Serializer\Annotation\SerializedName('metadata')]
    #[\JMS\Serializer\Annotation\Type('array<string, string>')]
    public array $metadata;

    /**
     * @param  ?string  $id
     * @param  ?string  $ledger
     * @param  ?array<string, string>  $metadata
     */
    public function __construct(?string $id = null, ?string $ledger = null, ?array $metadata = null)
    {
        $this->id = $id;
        $this->ledger = $ledger;
        $this->metadata = $metadata;
    }
}