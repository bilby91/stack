<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class TransactionResponse
{
    /**
     *
     * @var Transaction $data
     */
    #[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Shared\Transaction')]
    public Transaction $data;

    /**
     * @param  ?Transaction  $data
     */
    public function __construct(?Transaction $data = null)
    {
        $this->data = $data;
    }
}