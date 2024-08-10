<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class Script
{
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
     * @var string $plain
     */
    #[\JMS\Serializer\Annotation\SerializedName('plain')]
    public string $plain;

    /**
     * Reference to attach to the generated transaction
     *
     * @var ?string $reference
     */
    #[\JMS\Serializer\Annotation\SerializedName('reference')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $reference = null;

    /**
     * $vars
     *
     * @var ?array<string, mixed> $vars
     */
    #[\JMS\Serializer\Annotation\SerializedName('vars')]
    #[\JMS\Serializer\Annotation\Type('array<string, mixed>')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?array $vars = null;

    /**
     * @param  ?string  $plain
     * @param  ?array<string, mixed>  $metadata
     * @param  ?string  $reference
     * @param  ?array<string, mixed>  $vars
     */
    public function __construct(?string $plain = null, ?array $metadata = null, ?string $reference = null, ?array $vars = null)
    {
        $this->plain = $plain;
        $this->metadata = $metadata;
        $this->reference = $reference;
        $this->vars = $vars;
    }
}