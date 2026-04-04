package dsl

import (
	"testing"
)

func TestValidate_MinimalValid(t *testing.T) {
	spec := &Spec{
		List: ListSpec{
			ElementType:    "T",
			Ownership:      "Copy",
			Morphology:     "Monomorphic",
			TargetLanguage: "cpp",
		},
	}
	errs := Validate(spec)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidate_FullValid(t *testing.T) {
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
	errs := Validate(spec)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidate_C1_PolymorphicCopy(t *testing.T) {
	spec := &Spec{
		List: ListSpec{
			ElementType:    "T",
			Ownership:      "Copy",
			Morphology:     "Polymorphic",
			TargetLanguage: "cpp",
		},
	}
	errs := Validate(spec)
	assertHasErrorCode(t, errs, "E001")
}

func TestValidate_C2_CounterTypeRequired(t *testing.T) {
	spec := &Spec{
		List: ListSpec{
			ElementType:    "T",
			Ownership:      "Copy",
			Morphology:     "Monomorphic",
			LengthCounter:  &LengthCounterSpec{Enabled: true},
			TargetLanguage: "cpp",
		},
	}
	errs := Validate(spec)
	assertHasErrorCode(t, errs, "E002")
}

func TestValidate_C3_CounterTypeNotAllowed(t *testing.T) {
	spec := &Spec{
		List: ListSpec{
			ElementType: "T",
			Ownership:   "Copy",
			Morphology:  "Monomorphic",
			LengthCounter: &LengthCounterSpec{
				Enabled:     false,
				CounterType: "int",
			},
			TargetLanguage: "cpp",
		},
	}
	errs := Validate(spec)
	assertHasErrorCode(t, errs, "E003")
}

func TestValidate_C4_MissingRequired(t *testing.T) {
	spec := &Spec{
		List: ListSpec{},
	}
	errs := Validate(spec)
	// element_type, ownership, morphology, target_language の4つ
	count := 0
	for _, e := range errs {
		if e.Code == "E004" {
			count++
		}
	}
	if count != 4 {
		t.Errorf("expected 4 E004 errors, got %d", count)
	}
}

func TestValidate_C5_InvalidOwnership(t *testing.T) {
	spec := &Spec{
		List: ListSpec{
			ElementType:    "T",
			Ownership:      "Invalid",
			Morphology:     "Monomorphic",
			TargetLanguage: "cpp",
		},
	}
	errs := Validate(spec)
	assertHasErrorCode(t, errs, "E005")
}

func TestValidate_C6_InvalidMorphology(t *testing.T) {
	spec := &Spec{
		List: ListSpec{
			ElementType:    "T",
			Ownership:      "Copy",
			Morphology:     "Invalid",
			TargetLanguage: "cpp",
		},
	}
	errs := Validate(spec)
	assertHasErrorCode(t, errs, "E006")
}

func TestValidate_InvalidCounterType(t *testing.T) {
	spec := &Spec{
		List: ListSpec{
			ElementType: "T",
			Ownership:   "Copy",
			Morphology:  "Monomorphic",
			LengthCounter: &LengthCounterSpec{
				Enabled:     true,
				CounterType: "float",
			},
			TargetLanguage: "cpp",
		},
	}
	errs := Validate(spec)
	assertHasErrorCode(t, errs, "E005")
}

func assertHasErrorCode(t *testing.T, errs []*ValidationError, code string) {
	t.Helper()
	for _, e := range errs {
		if e.Code == code {
			return
		}
	}
	t.Errorf("expected error code %s, got %v", code, errs)
}
