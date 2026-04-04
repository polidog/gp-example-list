// Generated from: docs/gp/08-dsl-examples/example-02-full-features.yaml
// Configuration: OwnedReference + Polymorphic, LengthCounter(int), Tracing
// Language: Java

import java.util.function.Consumer;

// R1: List構造体
// Morphology=Polymorphic + Java → 上限境界 <T>
// LengthCounter=有効(int): lengthフィールドあり
// Tracing=有効: トレースログあり
public class List<T> {

    // R2: Node構造体 — OwnedRefPolymorphicNode<T>
    // Ownership=OwnedReference: 参照保持、破棄時にnull化
    // Morphology=Polymorphic: Tのサブタイプも許容
    private static class Node<T> {
        T data;
        Node<T> next;

        Node(T data, Node<T> next) {
            this.data = data;
            this.next = next;
        }
    }

    private Node<T> head;
    private int length;  // LengthCounter: CounterType=int

    public List() {
        this.head = null;
        this.length = 0;
    }

    // R3: insert — 先頭への挿入
    // Ownership=OwnedReference: 参照を保持（所有権をリストが取得）
    // Tracing=有効: 操作前後にトレース出力
    // LengthCounter=有効: length++
    public void insert(T element) {
        trace("insert:begin", element);
        Node<T> node = new Node<>(element, head);
        head = node;
        length++;  // LengthCounter更新
        trace("insert:end", "length=" + length);  // 組み合わせルール: TraceCounterUpdates
    }

    // R4: remove — 先頭要素の削除
    // Ownership=OwnedReference: 参照をnull化（GCによる解放を促進）
    // Tracing=有効: 操作前後にトレース + メモリ解放トレース
    // LengthCounter=有効: length--
    public boolean remove() {
        trace("remove:begin", null);
        if (head == null) {
            trace("remove:end", "empty list");
            return false;
        }
        Node<T> oldHead = head;
        head = head.next;
        // Ownership=OwnedReference: 参照null化で解放を促進
        trace("memory:release", oldHead.data);  // 組み合わせルール: TraceMemoryRelease
        oldHead.data = null;
        oldHead.next = null;
        length--;  // LengthCounter更新
        trace("remove:end", "length=" + length);
        return true;
    }

    // R5: destroy — リスト破棄
    // Ownership=OwnedReference: 全要素の参照をnull化
    // Tracing=有効: 破棄トレース
    public void destroy() {
        trace("destroy:begin", null);
        Node<T> current = head;
        while (current != null) {
            Node<T> next = current.next;
            trace("memory:release", current.data);  // TraceMemoryRelease
            current.data = null;
            current.next = null;
            current = next;
        }
        head = null;
        length = 0;
        trace("destroy:end", null);
    }

    // R6: length — O(1)返却（LengthCounter有効）
    public int length() {
        return length;
    }

    // find — 走査ベースの検索
    // Tracing=有効: 検索トレース
    public T find(T element) {
        trace("find:begin", element);
        Node<T> current = head;
        while (current != null) {
            if (current.data != null && current.data.equals(element)) {
                trace("find:end", "found");
                return current.data;
            }
            current = current.next;
        }
        trace("find:end", "not found");
        return null;
    }

    // traverse — イテレーション
    public void traverse(Consumer<T> fn) {
        Node<T> current = head;
        while (current != null) {
            fn.accept(current.data);
            current = current.next;
        }
    }

    public boolean isEmpty() {
        return head == null;
    }

    // Tracing実装: 標準出力へのトレースログ
    private void trace(String operation, Object detail) {
        System.out.println("[TRACE] " + operation + (detail != null ? " | " + detail : ""));
    }
}
