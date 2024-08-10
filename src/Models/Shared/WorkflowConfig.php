<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class WorkflowConfig
{
    /**
     *
     * @var ?string $name
     */
    #[\JMS\Serializer\Annotation\SerializedName('name')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $name = null;

    /**
     * $stages
     *
     * @var array<array<string, mixed>> $stages
     */
    #[\JMS\Serializer\Annotation\SerializedName('stages')]
    #[\JMS\Serializer\Annotation\Type('array<array<string, mixed>>')]
    public array $stages;

    /**
     * @param  ?array<array<string, mixed>>  $stages
     * @param  ?string  $name
     */
    public function __construct(?array $stages = null, ?string $name = null)
    {
        $this->stages = $stages;
        $this->name = $name;
    }
}