package dsl

// ResolvedConfig は構成解決後の結果を表す。
// フィーチャー選択からコンポーネントバリアントが一意に決定された状態。
type ResolvedConfig struct {
	// フィーチャー選択（入力の正規化）
	ElementType    string
	Ownership      string
	Morphology     string
	LengthCounter  bool
	CounterType    string // LengthCounter == true の場合のみ有効
	Tracing        bool
	TargetLanguage string

	// 解決されたコンポーネントバリアント
	NodeVariant          string // "CopyMonomorphicNode", "ExtRefPolymorphicNode" 等
	MemoryVariant        string // "CopyMemory", "OwnedRefMemory", "ExtRefMemory"
	LengthCounterVariant string // "NoCounter", "ShortCounter", "IntCounter", "LongCounter"
	TracingVariant       string // "NoTracing", "StdoutTracing"

	// 組み合わせルールによる追加処理フラグ
	TraceCounterUpdates bool // LengthCounter + Tracing
	TraceMemoryRelease  bool // OwnedReference + Tracing
	NeedsVirtualDtor    bool // Polymorphic + cpp
	NeedsBoundedType    bool // Polymorphic + java
	NeedsInterfaceType  bool // Polymorphic + go
}

// Resolve は検証済みの仕様書から構成を解決し、コンポーネントバリアントを決定する。
// Validate を先に呼び出してエラーがないことを確認してから呼ぶこと。
func Resolve(spec *Spec) *ResolvedConfig {
	cfg := &ResolvedConfig{
		ElementType:    spec.List.ElementType,
		Ownership:      spec.List.Ownership,
		Morphology:     spec.List.Morphology,
		TargetLanguage: spec.List.TargetLanguage,
	}

	// LengthCounter
	if spec.List.LengthCounter != nil && spec.List.LengthCounter.Enabled {
		cfg.LengthCounter = true
		cfg.CounterType = spec.List.LengthCounter.CounterType
	}

	// Tracing
	if spec.List.Tracing != nil && spec.List.Tracing.Enabled {
		cfg.Tracing = true
	}

	// NodeDefinition バリアント解決 (Ownership × Morphology)
	cfg.NodeVariant = resolveNodeVariant(cfg.Ownership, cfg.Morphology)

	// MemoryManagement バリアント解決
	cfg.MemoryVariant = resolveMemoryVariant(cfg.Ownership)

	// LengthCounter バリアント解決
	cfg.LengthCounterVariant = resolveCounterVariant(cfg.LengthCounter, cfg.CounterType)

	// Tracing バリアント解決
	if cfg.Tracing {
		cfg.TracingVariant = "StdoutTracing"
	} else {
		cfg.TracingVariant = "NoTracing"
	}

	// 組み合わせルール適用
	cfg.TraceCounterUpdates = cfg.LengthCounter && cfg.Tracing
	cfg.TraceMemoryRelease = cfg.Ownership == "OwnedReference" && cfg.Tracing
	cfg.NeedsVirtualDtor = cfg.Morphology == "Polymorphic" && cfg.TargetLanguage == "cpp"
	cfg.NeedsBoundedType = cfg.Morphology == "Polymorphic" && cfg.TargetLanguage == "java"
	cfg.NeedsInterfaceType = cfg.Morphology == "Polymorphic" && cfg.TargetLanguage == "go"

	return cfg
}

func resolveNodeVariant(ownership, morphology string) string {
	switch ownership {
	case "Copy":
		return "CopyMonomorphicNode" // Copy + Polymorphic は制約で排除済み
	case "ExternalReference":
		if morphology == "Polymorphic" {
			return "ExtRefPolymorphicNode"
		}
		return "ExtRefMonomorphicNode"
	case "OwnedReference":
		if morphology == "Polymorphic" {
			return "OwnedRefPolymorphicNode"
		}
		return "OwnedRefMonomorphicNode"
	default:
		return ""
	}
}

func resolveMemoryVariant(ownership string) string {
	switch ownership {
	case "Copy":
		return "CopyMemory"
	case "ExternalReference":
		return "ExtRefMemory"
	case "OwnedReference":
		return "OwnedRefMemory"
	default:
		return ""
	}
}

func resolveCounterVariant(enabled bool, counterType string) string {
	if !enabled {
		return "NoCounter"
	}
	switch counterType {
	case "short":
		return "ShortCounter"
	case "int":
		return "IntCounter"
	case "long":
		return "LongCounter"
	default:
		return ""
	}
}
