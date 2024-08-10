<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class TransactionsCursorResponse
{
    /**
     *
     * @var TransactionsCursorResponseCursor $cursor
     */
    #[\JMS\Serializer\Annotation\SerializedName('cursor')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Shared\TransactionsCursorResponseCursor')]
    public TransactionsCursorResponseCursor $cursor;

    /**
     * @param  ?TransactionsCursorResponseCursor  $cursor
     */
    public function __construct(?TransactionsCursorResponseCursor $cursor = null)
    {
        $this->cursor = $cursor;
    }
}