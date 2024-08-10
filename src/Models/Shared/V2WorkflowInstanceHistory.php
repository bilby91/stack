<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2WorkflowInstanceHistory
{
    /**
     *
     * @var ?string $error
     */
    #[\JMS\Serializer\Annotation\SerializedName('error')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $error = null;

    /**
     *
     * @var V2StageSend|V2StageDelay|V2StageWaitEvent|V2Update $input
     */
    #[\JMS\Serializer\Annotation\SerializedName('input')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Shared\V2StageSend|\formance\stack\Models\Shared\V2StageDelay|\formance\stack\Models\Shared\V2StageWaitEvent|\formance\stack\Models\Shared\V2Update')]
    public V2StageSend|V2StageDelay|V2StageWaitEvent|V2Update $input;

    /**
     *
     * @var string $name
     */
    #[\JMS\Serializer\Annotation\SerializedName('name')]
    public string $name;

    /**
     *
     * @var \DateTime $startedAt
     */
    #[\JMS\Serializer\Annotation\SerializedName('startedAt')]
    public \DateTime $startedAt;

    /**
     *
     * @var bool $terminated
     */
    #[\JMS\Serializer\Annotation\SerializedName('terminated')]
    public bool $terminated;

    /**
     *
     * @var ?\DateTime $terminatedAt
     */
    #[\JMS\Serializer\Annotation\SerializedName('terminatedAt')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?\DateTime $terminatedAt = null;

    /**
     * @param  V2StageSend|V2StageDelay|V2StageWaitEvent|V2Update|null  $input
     * @param  ?string  $name
     * @param  ?\DateTime  $startedAt
     * @param  ?bool  $terminated
     * @param  ?string  $error
     * @param  ?\DateTime  $terminatedAt
     */
    public function __construct(V2StageSend|V2StageDelay|V2StageWaitEvent|V2Update|null $input = null, ?string $name = null, ?\DateTime $startedAt = null, ?bool $terminated = null, ?string $error = null, ?\DateTime $terminatedAt = null)
    {
        $this->input = $input;
        $this->name = $name;
        $this->startedAt = $startedAt;
        $this->terminated = $terminated;
        $this->error = $error;
        $this->terminatedAt = $terminatedAt;
    }
}