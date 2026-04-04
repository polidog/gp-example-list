<?php

declare(strict_types=1);

// Generated from: docs/gp/08-dsl-examples/example-05-minimal-php.yaml
// Configuration: Copy + Monomorphic, No LengthCounter, No Tracing
// Language: PHP

namespace Generated\List;

/**
 * R2: Node構造体 — CopyMonomorphicNode<T>
 * Ownership=Copy: 値埋め込み（PHPは値型を自動コピー���オブジェクトはclone）
 * Morphology=Monomorphic: 型Tのみ許容
 *
 * @template T
 */
final class Node
{
    /**
     * @param T $data
     * @param self<T>|null $next
     */
    public function __construct(
        public mixed $data,
        public ?self $next = null,
    ) {
    }
}
