<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class UpdateTransferInitiationStatusRequest
{
    /**
     *
     * @var Status $status
     */
    #[\JMS\Serializer\Annotation\SerializedName('status')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Shared\Status')]
    public Status $status;

    /**
     * @param  ?Status  $status
     */
    public function __construct(?Status $status = null)
    {
        $this->status = $status;
    }
}