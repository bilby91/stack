<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class Account
{
    /**
     *
     * @var string $address
     */
    #[\JMS\Serializer\Annotation\SerializedName('address')]
    public string $address;

    /**
     * $metadata
     *
     * @var ?array<string, mixed> $metadata
     */
    #[\JMS\Serializer\Annotation\SerializedName('metadata')]
    #[\JMS\Serializer\Annotation\Type('array<string, mixed>')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?array $metadata = null;

    /**
     *
     * @var ?string $type
     */
    #[\JMS\Serializer\Annotation\SerializedName('type')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $type = null;

    /**
     * @param  ?string  $address
     * @param  ?array<string, mixed>  $metadata
     * @param  ?string  $type
     */
    public function __construct(?string $address = null, ?array $metadata = null, ?string $type = null)
    {
        $this->address = $address;
        $this->metadata = $metadata;
        $this->type = $type;
    }
}