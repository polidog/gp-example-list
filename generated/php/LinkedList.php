<?php

declare(strict_types=1);

// Generated from: docs/gp/08-dsl-examples/example-05-minimal-php.yaml
// Configuration: Copy + Monomorphic, No LengthCounter, No Tracing
// Language: PHP

namespace Generated\List;

/**
 * R1: List構造体
 * LengthCounter=無効: lengthフィールドなし
 * Tracing=無効: トレースなし
 *
 * @template T
 */
final class LinkedList
{
    /** @var Node<T>|null */
    private ?Node $head = null;

    /**
     * R3: insert — 先頭への挿入
     * Ownership=Copy: 値を直接保持（PHPのスカラー値は自動コピー）
     *
     * @param T $element
     */
    public function insert(mixed $element): void
    {
        $this->head = new Node($element, $this->head);
    }

    /**
     * R4: remove — 先頭要素の削除
     * Ownership=Copy: ノード解放のみ（GCが処理）
     */
    public function remove(): bool
    {
        if ($this->head === null) {
            return false;
        }
        $this->head = $this->head->next;
        return true;
    }

    /**
     * R5: destroy — リスト破棄
     * Ownership=Copy: 参照をnullにしてGCに委ねる
     */
    public function destroy(): void
    {
        $this->head = null;
    }

    /**
     * R6: length — O(n)走査（LengthCounter無効）
     */
    public function length(): int
    {
        $count = 0;
        $current = $this->head;
        while ($current !== null) {
            $count++;
            $current = $current->next;
        }
        return $count;
    }

    /**
     * find — 走査ベースの検索
     *
     * @param T $element
     * @return T|null
     */
    public function find(mixed $element): mixed
    {
        $current = $this->head;
        while ($current !== null) {
            if ($current->data === $element) {
                return $current->data;
            }
            $current = $current->next;
        }
        return null;
    }

    /**
     * traverse — イテレーション
     *
     * @param callable(T): void $fn
     */
    public function traverse(callable $fn): void
    {
        $current = $this->head;
        while ($current !== null) {
            $fn($current->data);
            $current = $current->next;
        }
    }

    public function isEmpty(): bool
    {
        return $this->head === null;
    }
}
