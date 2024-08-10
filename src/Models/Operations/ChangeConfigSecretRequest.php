<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Models\Shared;
use formance\stack\Utils\SpeakeasyMetadata;
class ChangeConfigSecretRequest
{
    /**
     *
     * @var ?Shared\ConfigChangeSecret $configChangeSecret
     */
    #[SpeakeasyMetadata('request:mediaType=application/json')]
    public ?Shared\ConfigChangeSecret $configChangeSecret = null;

    /**
     * Config ID
     *
     * @var string $id
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=id')]
    public string $id;

    /**
     * @param  ?string  $id
     * @param  ?Shared\ConfigChangeSecret  $configChangeSecret
     */
    public function __construct(?string $id = null, ?Shared\ConfigChangeSecret $configChangeSecret = null)
    {
        $this->id = $id;
        $this->configChangeSecret = $configChangeSecret;
    }
}