<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2Workflow
{
    /**
     *
     * @var V2WorkflowConfig $config
     */
    #[\JMS\Serializer\Annotation\SerializedName('config')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Shared\V2WorkflowConfig')]
    public V2WorkflowConfig $config;

    /**
     *
     * @var \DateTime $createdAt
     */
    #[\JMS\Serializer\Annotation\SerializedName('createdAt')]
    public \DateTime $createdAt;

    /**
     *
     * @var string $id
     */
    #[\JMS\Serializer\Annotation\SerializedName('id')]
    public string $id;

    /**
     *
     * @var \DateTime $updatedAt
     */
    #[\JMS\Serializer\Annotation\SerializedName('updatedAt')]
    public \DateTime $updatedAt;

    /**
     * @param  ?V2WorkflowConfig  $config
     * @param  ?\DateTime  $createdAt
     * @param  ?string  $id
     * @param  ?\DateTime  $updatedAt
     */
    public function __construct(?V2WorkflowConfig $config = null, ?\DateTime $createdAt = null, ?string $id = null, ?\DateTime $updatedAt = null)
    {
        $this->config = $config;
        $this->createdAt = $createdAt;
        $this->id = $id;
        $this->updatedAt = $updatedAt;
    }
}