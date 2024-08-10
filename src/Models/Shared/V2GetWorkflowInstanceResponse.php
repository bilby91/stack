<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2GetWorkflowInstanceResponse
{
    /**
     *
     * @var V2WorkflowInstance $data
     */
    #[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Shared\V2WorkflowInstance')]
    public V2WorkflowInstance $data;

    /**
     * @param  ?V2WorkflowInstance  $data
     */
    public function __construct(?V2WorkflowInstance $data = null)
    {
        $this->data = $data;
    }
}