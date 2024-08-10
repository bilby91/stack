<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Utils\SpeakeasyMetadata;
class GetTransferInitiationRequest
{
    /**
     * The transfer ID.
     *
     * @var string $transferId
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=transferId')]
    public string $transferId;

    /**
     * @param  ?string  $transferId
     */
    public function __construct(?string $transferId = null)
    {
        $this->transferId = $transferId;
    }
}