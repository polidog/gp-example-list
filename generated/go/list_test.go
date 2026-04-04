package list

import (
	"fmt"
	"testing"
)

// --- CopyMonoList のテスト ---

func TestCopyMonoList_InsertAndLength(t *testing.T) {
	l := NewCopyMonoList[int]()

	l.Insert(1)
	l.Insert(2)
	l.Insert(3)

	if got := l.Length(); got != 3 {
		t.Errorf("Length() = %d, want 3", got)
	}
}

func TestCopyMonoList_Remove(t *testing.T) {
	l := NewCopyMonoList[int]()
	l.Insert(1)
	l.Insert(2)

	if !l.Remove() {
		t.Error("Remove() returned false on non-empty list")
	}
	if got := l.Length(); got != 1 {
		t.Errorf("Length() after remove = %d, want 1", got)
	}
}

func TestCopyMonoList_RemoveEmpty(t *testing.T) {
	l := NewCopyMonoList[int]()
	if l.Remove() {
		t.Error("Remove() returned true on empty list")
	}
}

func TestCopyMonoList_Find(t *testing.T) {
	l := NewCopyMonoList[int]()
	l.Insert(10)
	l.Insert(20)
	l.Insert(30)

	if v, ok := l.Find(20); !ok || v != 20 {
		t.Errorf("Find(20) = (%v, %v), want (20, true)", v, ok)
	}
	if _, ok := l.Find(99); ok {
		t.Error("Find(99) should return false")
	}
}

func TestCopyMonoList_Traverse(t *testing.T) {
	l := NewCopyMonoList[int]()
	l.Insert(3)
	l.Insert(2)
	l.Insert(1)

	var result []int
	l.Traverse(func(v int) { result = append(result, v) })

	if len(result) != 3 || result[0] != 1 || result[1] != 2 || result[2] != 3 {
		t.Errorf("Traverse result = %v, want [1 2 3]", result)
	}
}

func TestCopyMonoList_Destroy(t *testing.T) {
	l := NewCopyMonoList[int]()
	l.Insert(1)
	l.Insert(2)
	l.Destroy()

	if !l.IsEmpty() {
		t.Error("IsEmpty() should be true after Destroy()")
	}
}

// --- ExtRefPolyCounterList のテスト ---

// テスト用の型（Stringer インターフェースを満たす）
type item struct {
	Name string
}

func (i item) String() string {
	return fmt.Sprintf("item(%s)", i.Name)
}

func TestExtRefPolyCounterList_InsertAndLength(t *testing.T) {
	l := NewExtRefPolyCounterList[item]()

	a := item{Name: "a"}
	b := item{Name: "b"}
	c := item{Name: "c"}

	l.Insert(&a)
	l.Insert(&b)
	l.Insert(&c)

	if got := l.Length(); got != 3 {
		t.Errorf("Length() = %d, want 3", got)
	}
}

func TestExtRefPolyCounterList_RemoveAndLength(t *testing.T) {
	l := NewExtRefPolyCounterList[item]()

	a := item{Name: "a"}
	l.Insert(&a)

	if !l.Remove() {
		t.Error("Remove() returned false on non-empty list")
	}
	if got := l.Length(); got != 0 {
		t.Errorf("Length() after remove = %d, want 0", got)
	}
}

func TestExtRefPolyCounterList_Find(t *testing.T) {
	l := NewExtRefPolyCounterList[item]()

	a := item{Name: "a"}
	b := item{Name: "b"}
	l.Insert(&a)
	l.Insert(&b)

	target := item{Name: "a"}
	if found, ok := l.Find(&target); !ok || found.Name != "a" {
		t.Errorf("Find(a) = (%v, %v), want (a, true)", found, ok)
	}

	missing := item{Name: "z"}
	if _, ok := l.Find(&missing); ok {
		t.Error("Find(z) should return false")
	}
}

func TestExtRefPolyCounterList_ExternalOwnership(t *testing.T) {
	l := NewExtRefPolyCounterList[item]()

	a := item{Name: "original"}
	l.Insert(&a)

	// 外部参照: 外部で値を変更するとリスト内のデータも変わる
	a.Name = "modified"

	var found string
	l.Traverse(func(v *item) { found = v.Name })

	if found != "modified" {
		t.Errorf("ExternalReference should reflect external changes, got %q", found)
	}
}

func TestExtRefPolyCounterList_Destroy(t *testing.T) {
	l := NewExtRefPolyCounterList[item]()

	a := item{Name: "a"}
	l.Insert(&a)
	l.Destroy()

	if !l.IsEmpty() {
		t.Error("IsEmpty() should be true after Destroy()")
	}
	if got := l.Length(); got != 0 {
		t.Errorf("Length() after Destroy = %d, want 0", got)
	}
}
