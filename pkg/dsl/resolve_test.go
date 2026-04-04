package dsl

import (
	"testing"
)

func TestResolve_Minimal(t *testing.T) {
	spec := &Spec{
		List: ListSpec{
			ElementType:    "T",
			Ownership:      "Copy",
			Morphology:     "Monomorphic",
			TargetLanguage: "cpp",
		},
	}
	cfg := Resolve(spec)

	assertEqual(t, "NodeVariant", cfg.NodeVariant, "CopyMonomorphicNode")
	assertEqual(t, "MemoryVariant", cfg.MemoryVariant, "CopyMemory")
	assertEqual(t, "LengthCounterVariant", cfg.LengthCounterVariant, "NoCounter")
	assertEqual(t, "TracingVariant", cfg.TracingVariant, "NoTracing")
	assertFalse(t, "TraceCounterUpdates", cfg.TraceCounterUpdates)
	assertFalse(t, "TraceMemoryRelease", cfg.TraceMemoryRelease)
	assertFalse(t, "NeedsVirtualDtor", cfg.NeedsVirtualDtor)
}

func TestResolve_FullFeatures(t *testing.T) {
	spec := &Spec{
		List: ListSpec{
			ElementType: "T",
			Ownership:   "OwnedReference",
			Morphology:  "Polymorphic",
			LengthCounter: &LengthCounterSpec{
				Enabled:     true,
				CounterType: "int",
			},
			Tracing:        &TracingSpec{Enabled: true},
			TargetLanguage: "java",
		},
	}
	cfg := Resolve(spec)

	assertEqual(t, "NodeVariant", cfg.NodeVariant, "OwnedRefPolymorphicNode")
	assertEqual(t, "MemoryVariant", cfg.MemoryVariant, "OwnedRefMemory")
	assertEqual(t, "LengthCounterVariant", cfg.LengthCounterVariant, "IntCounter")
	assertEqual(t, "TracingVariant", cfg.TracingVariant, "StdoutTracing")
	assertTrue(t, "TraceCounterUpdates", cfg.TraceCounterUpdates)
	assertTrue(t, "TraceMemoryRelease", cfg.TraceMemoryRelease)
	assertTrue(t, "NeedsBoundedType", cfg.NeedsBoundedType)
	assertFalse(t, "NeedsVirtualDtor", cfg.NeedsVirtualDtor)
}

func TestResolve_ExtRefPolymorphicCpp(t *testing.T) {
	spec := &Spec{
		List: ListSpec{
			ElementType:    "T",
			Ownership:      "ExternalReference",
			Morphology:     "Polymorphic",
			TargetLanguage: "cpp",
		},
	}
	cfg := Resolve(spec)

	assertEqual(t, "NodeVariant", cfg.NodeVariant, "ExtRefPolymorphicNode")
	assertEqual(t, "MemoryVariant", cfg.MemoryVariant, "ExtRefMemory")
	assertTrue(t, "NeedsVirtualDtor", cfg.NeedsVirtualDtor)
}

func TestResolve_ShortCounter(t *testing.T) {
	spec := &Spec{
		List: ListSpec{
			ElementType: "T",
			Ownership:   "ExternalReference",
			Morphology:  "Monomorphic",
			LengthCounter: &LengthCounterSpec{
				Enabled:     true,
				CounterType: "short",
			},
			TargetLanguage: "cpp",
		},
	}
	cfg := Resolve(spec)

	assertEqual(t, "LengthCounterVariant", cfg.LengthCounterVariant, "ShortCounter")
	assertEqual(t, "CounterType", cfg.CounterType, "short")
}

func TestResolve_GoPolymorphic(t *testing.T) {
	spec := &Spec{
		List: ListSpec{
			ElementType:    "T",
			Ownership:      "ExternalReference",
			Morphology:     "Polymorphic",
			TargetLanguage: "go",
		},
	}
	cfg := Resolve(spec)

	assertTrue(t, "NeedsInterfaceType", cfg.NeedsInterfaceType)
	assertFalse(t, "NeedsVirtualDtor", cfg.NeedsVirtualDtor)
	assertFalse(t, "NeedsBoundedType", cfg.NeedsBoundedType)
}

func assertEqual(t *testing.T, name, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", name, got, want)
	}
}

func assertTrue(t *testing.T, name string, got bool) {
	t.Helper()
	if !got {
		t.Errorf("%s: expected true, got false", name)
	}
}

func assertFalse(t *testing.T, name string, got bool) {
	t.Helper()
	if got {
		t.Errorf("%s: expected false, got true", name)
	}
}
