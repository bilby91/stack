<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2BulkElementRevertTransactionData
{
    /**
     *
     * @var ?bool $atEffectiveDate
     */
    #[\JMS\Serializer\Annotation\SerializedName('atEffectiveDate')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?bool $atEffectiveDate = null;

    /**
     *
     * @var ?bool $force
     */
    #[\JMS\Serializer\Annotation\SerializedName('force')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?bool $force = null;

    /**
     *
     * @var int $id
     */
    #[\JMS\Serializer\Annotation\SerializedName('id')]
    public int $id;

    /**
     * @param  ?int  $id
     * @param  ?bool  $atEffectiveDate
     * @param  ?bool  $force
     */
    public function __construct(?int $id = null, ?bool $atEffectiveDate = null, ?bool $force = null)
    {
        $this->id = $id;
        $this->atEffectiveDate = $atEffectiveDate;
        $this->force = $force;
    }
}