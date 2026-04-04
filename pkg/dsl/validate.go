package dsl

import (
	"fmt"
	"slices"
)

// ValidationError はバリデーションエラーを表す。
type ValidationError struct {
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Validate は仕様書の制約 C1〜C6 を検証し、違反があればエラーのリストを返す。
func Validate(spec *Spec) []*ValidationError {
	var errs []*ValidationError

	// C4: 必須フィーチャーの検証
	if spec.List.ElementType == "" {
		errs = append(errs, &ValidationError{Code: "E004", Message: "element_type は必須項目です"})
	}
	if spec.List.Ownership == "" {
		errs = append(errs, &ValidationError{Code: "E004", Message: "ownership は必須項目です"})
	}
	if spec.List.Morphology == "" {
		errs = append(errs, &ValidationError{Code: "E004", Message: "morphology は必須項目です"})
	}
	if spec.List.TargetLanguage == "" {
		errs = append(errs, &ValidationError{Code: "E004", Message: "target_language は必須項目です"})
	}

	// C5: Ownership の値検証
	if spec.List.Ownership != "" && !slices.Contains(ValidOwnership, spec.List.Ownership) {
		errs = append(errs, &ValidationError{
			Code:    "E005",
			Message: fmt.Sprintf("ownership の値が不正です: %s。ExternalReference, OwnedReference, Copy のいずれかを指定してください", spec.List.Ownership),
		})
	}

	// C6: Morphology の値検証
	if spec.List.Morphology != "" && !slices.Contains(ValidMorphology, spec.List.Morphology) {
		errs = append(errs, &ValidationError{
			Code:    "E006",
			Message: fmt.Sprintf("morphology の値が不正です: %s。Monomorphic, Polymorphic のいずれかを指定してください", spec.List.Morphology),
		})
	}

	// C1: Polymorphic + Copy 排他
	if spec.List.Morphology == "Polymorphic" && spec.List.Ownership == "Copy" {
		errs = append(errs, &ValidationError{
			Code:    "E001",
			Message: "ownership: Copy と morphology: Polymorphic は同時に指定できません",
		})
	}

	// C2, C3: LengthCounter と CounterType の整合性
	if spec.List.LengthCounter != nil && spec.List.LengthCounter.Enabled {
		if spec.List.LengthCounter.CounterType == "" {
			errs = append(errs, &ValidationError{
				Code:    "E002",
				Message: "length_counter.enabled が true の場合、counter_type は必須です",
			})
		} else if !slices.Contains(ValidCounterType, spec.List.LengthCounter.CounterType) {
			errs = append(errs, &ValidationError{
				Code:    "E005",
				Message: fmt.Sprintf("counter_type の値が不正です: %s。short, int, long のいずれかを指定してください", spec.List.LengthCounter.CounterType),
			})
		}
	}
	if spec.List.LengthCounter != nil && !spec.List.LengthCounter.Enabled && spec.List.LengthCounter.CounterType != "" {
		errs = append(errs, &ValidationError{
			Code:    "E003",
			Message: "length_counter が無効の場合、counter_type は指定できません",
		})
	}

	return errs
}
