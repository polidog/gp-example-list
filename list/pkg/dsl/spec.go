package dsl

// Spec represents a parsed DSL specification.
type Spec struct {
	Name        string   `yaml:"name"`
	Language    string   `yaml:"language"`
	ElementType string   `yaml:"element_type"`
	Storage     Storage  `yaml:"storage"`
	Operations  []any    `yaml:"operations"`
}

type Storage struct {
	Type            string `yaml:"type"`
	InitialCapacity int    `yaml:"initial_capacity"`
	GrowthStrategy  string `yaml:"growth_strategy"`
	GrowthIncrement int    `yaml:"growth_increment"`
	Direction       string `yaml:"direction"`
}

// ResolvedSpec is the fully resolved specification after applying defaults and auto-completion.
type ResolvedSpec struct {
	Name        string
	Language    string
	ElementType string
	Storage     ResolvedStorage
	Operations  ResolvedOperations
}

type ResolvedStorage struct {
	Type            string
	InitialCapacity int
	GrowthStrategy  string
	GrowthIncrement int
	Direction       string
}

type ResolvedOperations struct {
	Remove       bool
	Insert       bool
	LinearSearch bool
	Contains     bool
	Sort         bool
	SortAlgorithm string
	Iteration    bool
	IterationDir string
}
