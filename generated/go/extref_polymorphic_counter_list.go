// Generated from: example-04-polymorphic-go.yaml
// Configuration: ExternalReference + Polymorphic, LengthCounter(int), No Tracing
// Language: Go

package list

import "fmt"

// --- NodeDefinition: ExtRefPolymorphicNode ---
// Ownership=ExternalReference: ポインタ保持、要素は外部管理
// Morphology=Polymorphic: インターフェース制約でサブタイプ許容

// Stringer はPolymorphic構成の型制約（インターフェース型）。
// NeedsInterfaceType=true: Go の Polymorphic はインターフェース制約で実現
type Stringer interface {
	comparable
	fmt.Stringer
}

// ExtRefPolyNode はExternalReference+Polymorphic構成のノード。
type ExtRefPolyNode[T Stringer] struct {
	Data *T // ExternalReference: ポインタ保持
	Next *ExtRefPolyNode[T]
}

// --- ListCore + PublicAPI + LengthCounterMixin ---

// ExtRefPolyCounterList はExternalReference+Polymorphic+LengthCounter構成のリスト。
// LengthCounter=有効(int32), Tracing=無効
type ExtRefPolyCounterList[T Stringer] struct {
	head   *ExtRefPolyNode[T]
	length int32 // LengthCounter: CounterType=int → int32
}

// NewExtRefPolyCounterList は新しいリストを生成する。
func NewExtRefPolyCounterList[T Stringer]() *ExtRefPolyCounterList[T] {
	return &ExtRefPolyCounterList[T]{}
}

// R3: Insert — 先頭への挿入
// Ownership=ExternalReference: ポインタを保持（要素の所有権は外部）
// LengthCounter=有効: length++
func (l *ExtRefPolyCounterList[T]) Insert(element *T) {
	node := &ExtRefPolyNode[T]{Data: element, Next: l.head}
	l.head = node
	l.length++ // LengthCounter更新
}

// R4: Remove — 先頭要素の削除
// Ownership=ExternalReference: ポインタのみ破棄（要素は解放しない）
// LengthCounter=有効: length--
func (l *ExtRefPolyCounterList[T]) Remove() bool {
	if l.head == nil {
		return false
	}
	l.head = l.head.Next // ExtRef: 要素は解放しない
	l.length--           // LengthCounter更新
	return true
}

// R5: Destroy — リスト破棄
// Ownership=ExternalReference: ポインタのみクリア（要素は外部管理）
func (l *ExtRefPolyCounterList[T]) Destroy() {
	l.head = nil
	l.length = 0
}

// R6: Length — O(1)返却（LengthCounter有効）
func (l *ExtRefPolyCounterList[T]) Length() int32 {
	return l.length
}

// Find — 走査ベースの検索
func (l *ExtRefPolyCounterList[T]) Find(element *T) (*T, bool) {
	current := l.head
	for current != nil {
		if current.Data != nil && *current.Data == *element {
			return current.Data, true
		}
		current = current.Next
	}
	return nil, false
}

// Traverse — イテレーション
func (l *ExtRefPolyCounterList[T]) Traverse(fn func(*T)) {
	current := l.head
	for current != nil {
		fn(current.Data)
		current = current.Next
	}
}

// IsEmpty はリストが空かどうかを返す。
func (l *ExtRefPolyCounterList[T]) IsEmpty() bool {
	return l.head == nil
}
