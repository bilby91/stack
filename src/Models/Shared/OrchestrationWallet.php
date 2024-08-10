<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class OrchestrationWallet
{
    /**
     *
     * @var \DateTime $createdAt
     */
    #[\JMS\Serializer\Annotation\SerializedName('createdAt')]
    public \DateTime $createdAt;

    /**
     * The unique ID of the wallet.
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
     * Metadata associated with the wallet.
     *
     * @var array<string, string> $metadata
     */
    #[\JMS\Serializer\Annotation\SerializedName('metadata')]
    #[\JMS\Serializer\Annotation\Type('array<string, string>')]
    public array $metadata;

    /**
     *
     * @var string $name
     */
    #[\JMS\Serializer\Annotation\SerializedName('name')]
    public string $name;

    /**
     * @param  ?\DateTime  $createdAt
     * @param  ?string  $id
     * @param  ?string  $ledger
     * @param  ?array<string, string>  $metadata
     * @param  ?string  $name
     */
    public function __construct(?\DateTime $createdAt = null, ?string $id = null, ?string $ledger = null, ?array $metadata = null, ?string $name = null)
    {
        $this->createdAt = $createdAt;
        $this->id = $id;
        $this->ledger = $ledger;
        $this->metadata = $metadata;
        $this->name = $name;
    }
}