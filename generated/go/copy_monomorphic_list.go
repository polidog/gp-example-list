// Generated from: example-01-minimal (adapted to Go)
// Configuration: Copy + Monomorphic, No LengthCounter, No Tracing
// Language: Go

package list

// --- NodeDefinition: CopyMonomorphicNode ---
// Ownership=Copy: 値埋め込み
// Morphology=Monomorphic: 型Tのみ

// CopyMonoNode はCopy+Monomorphic構成のノード。
type CopyMonoNode[T comparable] struct {
	Data T
	Next *CopyMonoNode[T]
}

// --- ListCore + PublicAPI: CopyMonoList ---

// CopyMonoList はCopy+Monomorphic構成のリスト。
// LengthCounter=無効, Tracing=無効
type CopyMonoList[T comparable] struct {
	head *CopyMonoNode[T]
}

// NewCopyMonoList は新しいリストを生成する。
func NewCopyMonoList[T comparable]() *CopyMonoList[T] {
	return &CopyMonoList[T]{}
}

// R3: Insert — 先頭への挿入
// Ownership=Copy: 値コピー（Goは値型のため自動コピー）
func (l *CopyMonoList[T]) Insert(element T) {
	node := &CopyMonoNode[T]{Data: element, Next: l.head}
	l.head = node
}

// R4: Remove — 先頭要素の削除
// Ownership=Copy: ノード解放のみ（GCが処理）
func (l *CopyMonoList[T]) Remove() bool {
	if l.head == nil {
		return false
	}
	l.head = l.head.Next
	return true
}

// R5: Destroy — リスト破棄
// Ownership=Copy: 参照をnilにしてGCに委ねる
func (l *CopyMonoList[T]) Destroy() {
	l.head = nil
}

// R6: Length — O(n)走査（LengthCounter無効）
func (l *CopyMonoList[T]) Length() int {
	count := 0
	current := l.head
	for current != nil {
		count++
		current = current.Next
	}
	return count
}

// Find — 走査ベースの検索
func (l *CopyMonoList[T]) Find(element T) (T, bool) {
	current := l.head
	for current != nil {
		if current.Data == element {
			return current.Data, true
		}
		current = current.Next
	}
	var zero T
	return zero, false
}

// Traverse — イテレーション
func (l *CopyMonoList[T]) Traverse(fn func(T)) {
	current := l.head
	for current != nil {
		fn(current.Data)
		current = current.Next
	}
}

// IsEmpty はリストが空かどうかを返す。
func (l *CopyMonoList[T]) IsEmpty() bool {
	return l.head == nil
}
