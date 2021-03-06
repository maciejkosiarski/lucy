<?php

declare(strict_types=1);

namespace App\Domain;

class Url
{
    private string $url;

    public function __construct(string $url)
    {
        $this->url = $url;
    }

    public function getUrl(): string
    {
        return $this->url;
    }

    public function getHost(): string
    {
        return parse_url($this->url, PHP_URL_HOST);
    }
}
