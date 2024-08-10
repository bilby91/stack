<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class TaskModulrDescriptor
{
    /**
     *
     * @var ?string $accountID
     */
    #[\JMS\Serializer\Annotation\SerializedName('accountID')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $accountID = null;

    /**
     *
     * @var ?string $key
     */
    #[\JMS\Serializer\Annotation\SerializedName('key')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $key = null;

    /**
     *
     * @var ?string $name
     */
    #[\JMS\Serializer\Annotation\SerializedName('name')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $name = null;

    /**
     * @param  ?string  $accountID
     * @param  ?string  $key
     * @param  ?string  $name
     */
    public function __construct(?string $accountID = null, ?string $key = null, ?string $name = null)
    {
        $this->accountID = $accountID;
        $this->key = $key;
        $this->name = $name;
    }
}