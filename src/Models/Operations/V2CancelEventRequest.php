<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Utils\SpeakeasyMetadata;
class V2CancelEventRequest
{
    /**
     * The instance id
     *
     * @var string $instanceID
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=instanceID')]
    public string $instanceID;

    /**
     * @param  ?string  $instanceID
     */
    public function __construct(?string $instanceID = null)
    {
        $this->instanceID = $instanceID;
    }
}