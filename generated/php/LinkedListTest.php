<?php

declare(strict_types=1);

namespace Generated\List;

require_once __DIR__ . '/Node.php';
require_once __DIR__ . '/LinkedList.php';

use PHPUnit\Framework\TestCase;

final class LinkedListTest extends TestCase
{
    public function testInsertAndLength(): void
    {
        $list = new LinkedList();
        $list->insert(1);
        $list->insert(2);
        $list->insert(3);

        $this->assertSame(3, $list->length());
    }

    public function testRemove(): void
    {
        $list = new LinkedList();
        $list->insert(1);
        $list->insert(2);

        $this->assertTrue($list->remove());
        $this->assertSame(1, $list->length());
    }

    public function testRemoveEmpty(): void
    {
        $list = new LinkedList();
        $this->assertFalse($list->remove());
    }

    public function testFind(): void
    {
        $list = new LinkedList();
        $list->insert(10);
        $list->insert(20);
        $list->insert(30);

        $this->assertSame(20, $list->find(20));
        $this->assertNull($list->find(99));
    }

    public function testTraverse(): void
    {
        $list = new LinkedList();
        $list->insert(3);
        $list->insert(2);
        $list->insert(1);

        $result = [];
        $list->traverse(function (mixed $v) use (&$result): void {
            $result[] = $v;
        });

        $this->assertSame([1, 2, 3], $result);
    }

    public function testDestroy(): void
    {
        $list = new LinkedList();
        $list->insert(1);
        $list->insert(2);
        $list->destroy();

        $this->assertTrue($list->isEmpty());
        $this->assertSame(0, $list->length());
    }

    public function testIsEmpty(): void
    {
        $list = new LinkedList();
        $this->assertTrue($list->isEmpty());

        $list->insert(1);
        $this->assertFalse($list->isEmpty());
    }

    public function testCopyOwnership(): void
    {
        // Ownership=Copy: スカラー値はコピーされ、外部変更の影響を受けない
        $list = new LinkedList();
        $value = 42;
        $list->insert($value);

        $value = 99; // 外部で変更

        $found = $list->find(42);
        $this->assertSame(42, $found); // コピーなので42のまま
    }

    public function testStringElements(): void
    {
        $list = new LinkedList();
        $list->insert('world');
        $list->insert('hello');

        $result = [];
        $list->traverse(function (mixed $v) use (&$result): void {
            $result[] = $v;
        });

        $this->assertSame(['hello', 'world'], $result);
    }
}
