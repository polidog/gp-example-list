package dsl

// Spec はリストコンテナDSL仕様書のルート構造。
type Spec struct {
	List ListSpec `yaml:"list"`
}

// ListSpec はリストコンテナのフィーチャー選択を表す。
type ListSpec struct {
	ElementType   string             `yaml:"element_type"`
	Ownership     string             `yaml:"ownership"`
	Morphology    string             `yaml:"morphology"`
	LengthCounter *LengthCounterSpec `yaml:"length_counter,omitempty"`
	Tracing       *TracingSpec       `yaml:"tracing,omitempty"`
	TargetLanguage string            `yaml:"target_language"`
}

// LengthCounterSpec は長さカウンタの設定。
type LengthCounterSpec struct {
	Enabled     bool   `yaml:"enabled"`
	CounterType string `yaml:"counter_type,omitempty"`
}

// TracingSpec はトレーシングの設定。
type TracingSpec struct {
	Enabled bool `yaml:"enabled"`
}

// 有効な選択肢の定数
var (
	ValidOwnership  = []string{"ExternalReference", "OwnedReference", "Copy"}
	ValidMorphology = []string{"Monomorphic", "Polymorphic"}
	ValidCounterType = []string{"short", "int", "long"}
)
