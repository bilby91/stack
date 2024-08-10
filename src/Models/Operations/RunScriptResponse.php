<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Models\Shared;
class RunScriptResponse
{
    /**
     * HTTP response content type for this operation
     *
     * @var string $contentType
     */
    public string $contentType;

    /**
     * On success, it will return a 200 status code, and the resulting transaction under the `transaction` field.
     *
     *
     * On failure, it will also return a 200 status code, and the following fields:
     *   - `details`: contains a URL. When there is an error parsing Numscript, the result can be difficult to read—the provided URL will render the error in an easy-to-read format.
     *   - `errorCode` and `error_code` (deprecated): contains the string code of the error
     *   - `errorMessage` and `error_message` (deprecated): contains a human-readable indication of what went wrong, for example that an account had insufficient funds, or that there was an error in the provided Numscript.
     *
     *
     * @var ?Shared\ScriptResponse $scriptResponse
     */
    public ?Shared\ScriptResponse $scriptResponse = null;

    /**
     * HTTP response status code for this operation
     *
     * @var int $statusCode
     */
    public int $statusCode;

    /**
     * Raw HTTP response; suitable for custom response parsing
     *
     * @var \Psr\Http\Message\ResponseInterface $rawResponse
     */
    public \Psr\Http\Message\ResponseInterface $rawResponse;

    /**
     * @param  ?string  $contentType
     * @param  ?int  $statusCode
     * @param  ?\Psr\Http\Message\ResponseInterface  $rawResponse
     * @param  ?Shared\ScriptResponse  $scriptResponse
     */
    public function __construct(?string $contentType = null, ?int $statusCode = null, ?\Psr\Http\Message\ResponseInterface $rawResponse = null, ?Shared\ScriptResponse $scriptResponse = null)
    {
        $this->contentType = $contentType;
        $this->statusCode = $statusCode;
        $this->rawResponse = $rawResponse;
        $this->scriptResponse = $scriptResponse;
    }
}