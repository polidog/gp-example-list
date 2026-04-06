package dsl

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadSpec(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading spec file: %w", err)
	}
	var spec Spec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}
	return &spec, nil
}

func Validate(spec *Spec) []string {
	var errs []string

	// Required fields
	if spec.Name == "" {
		errs = append(errs, "name is required")
	}
	if spec.Language == "" {
		errs = append(errs, "language is required")
	}
	if spec.ElementType == "" {
		errs = append(errs, "element_type is required")
	}

	// Storage validation
	st := spec.Storage.Type
	if st == "" {
		st = "array"
	}
	switch st {
	case "array":
		if spec.Storage.Direction != "" {
			errs = append(errs, "storage.direction is only valid for linked_list")
		}
		if spec.Storage.InitialCapacity != 0 {
			if spec.Storage.InitialCapacity < 1 || spec.Storage.InitialCapacity > 65536 {
				errs = append(errs, "storage.initial_capacity must be between 1 and 65536")
			}
		}
		gs := spec.Storage.GrowthStrategy
		if gs != "" && gs != "doubling" && gs != "additive" {
			errs = append(errs, fmt.Sprintf("storage.growth_strategy must be 'doubling' or 'additive', got '%s'", gs))
		}
		if gs != "additive" && spec.Storage.GrowthIncrement != 0 {
			errs = append(errs, "storage.growth_increment is only valid when growth_strategy is 'additive'")
		}
	case "linked_list":
		if spec.Storage.InitialCapacity != 0 {
			errs = append(errs, "storage.initial_capacity is only valid for array")
		}
		if spec.Storage.GrowthStrategy != "" {
			errs = append(errs, "storage.growth_strategy is only valid for array")
		}
		if spec.Storage.GrowthIncrement != 0 {
			errs = append(errs, "storage.growth_increment is only valid for array")
		}
		dir := spec.Storage.Direction
		if dir != "" && dir != "singly" && dir != "doubly" {
			errs = append(errs, fmt.Sprintf("storage.direction must be 'singly' or 'doubly', got '%s'", dir))
		}
	default:
		errs = append(errs, fmt.Sprintf("storage.type must be 'array' or 'linked_list', got '%s'", st))
	}

	// Operations validation
	ops := parseOperations(spec.Operations)
	if ops.Sort && ops.SortAlgorithm == "" {
		errs = append(errs, "sort requires algorithm to be specified")
	}
	if ops.SortAlgorithm != "" && ops.SortAlgorithm != "insertion_sort" && ops.SortAlgorithm != "merge_sort" {
		errs = append(errs, fmt.Sprintf("sort.algorithm must be 'insertion_sort' or 'merge_sort', got '%s'", ops.SortAlgorithm))
	}
	if ops.Iteration && ops.IterationDir == "" {
		errs = append(errs, "iteration requires direction to be specified")
	}
	if ops.IterationDir != "" && ops.IterationDir != "forward" && ops.IterationDir != "reverse" {
		errs = append(errs, fmt.Sprintf("iteration.direction must be 'forward' or 'reverse', got '%s'", ops.IterationDir))
	}

	// Warning: singly + reverse
	if (st == "linked_list") && (spec.Storage.Direction == "" || spec.Storage.Direction == "singly") && ops.IterationDir == "reverse" {
		errs = append(errs, "[warning] singly linked list with reverse iteration is O(n) per step — consider doubly")
	}

	return errs
}

func Resolve(spec *Spec) *ResolvedSpec {
	r := &ResolvedSpec{
		Name:        spec.Name,
		Language:    spec.Language,
		ElementType: spec.ElementType,
	}

	// Storage defaults
	r.Storage.Type = spec.Storage.Type
	if r.Storage.Type == "" {
		r.Storage.Type = "array"
	}
	if r.Storage.Type == "array" {
		r.Storage.InitialCapacity = spec.Storage.InitialCapacity
		if r.Storage.InitialCapacity == 0 {
			r.Storage.InitialCapacity = 16
		}
		r.Storage.GrowthStrategy = spec.Storage.GrowthStrategy
		if r.Storage.GrowthStrategy == "" {
			r.Storage.GrowthStrategy = "doubling"
		}
		if r.Storage.GrowthStrategy == "additive" {
			r.Storage.GrowthIncrement = spec.Storage.GrowthIncrement
			if r.Storage.GrowthIncrement == 0 {
				r.Storage.GrowthIncrement = 16
			}
		}
	}
	if r.Storage.Type == "linked_list" {
		r.Storage.Direction = spec.Storage.Direction
		if r.Storage.Direction == "" {
			r.Storage.Direction = "singly"
		}
	}

	// Operations
	ops := parseOperations(spec.Operations)

	// Auto-completion: contains requires linear_search
	if ops.Contains && !ops.LinearSearch {
		ops.LinearSearch = true
	}

	r.Operations = ops
	return r
}

func parseOperations(raw []any) ResolvedOperations {
	var ops ResolvedOperations
	for _, item := range raw {
		switch v := item.(type) {
		case string:
			switch v {
			case "remove":
				ops.Remove = true
			case "insert":
				ops.Insert = true
			case "linear_search":
				ops.LinearSearch = true
			case "contains":
				ops.Contains = true
			}
		case map[string]any:
			for key, val := range v {
				switch key {
				case "sort":
					ops.Sort = true
					if m, ok := val.(map[string]any); ok {
						if alg, ok := m["algorithm"].(string); ok {
							ops.SortAlgorithm = alg
						}
					}
				case "iteration":
					ops.Iteration = true
					if m, ok := val.(map[string]any); ok {
						if dir, ok := m["direction"].(string); ok {
							ops.IterationDir = dir
						}
					}
				}
			}
		}
	}
	return ops
}

func FormatResolved(r *ResolvedSpec) string {
	var b strings.Builder
	fmt.Fprintf(&b, "name: %s\n", r.Name)
	fmt.Fprintf(&b, "language: %s\n", r.Language)
	fmt.Fprintf(&b, "element_type: %s\n", r.ElementType)
	fmt.Fprintf(&b, "storage:\n")
	fmt.Fprintf(&b, "  type: %s\n", r.Storage.Type)
	if r.Storage.Type == "array" {
		fmt.Fprintf(&b, "  initial_capacity: %d\n", r.Storage.InitialCapacity)
		fmt.Fprintf(&b, "  growth_strategy: %s\n", r.Storage.GrowthStrategy)
		if r.Storage.GrowthStrategy == "additive" {
			fmt.Fprintf(&b, "  growth_increment: %d\n", r.Storage.GrowthIncrement)
		}
	}
	if r.Storage.Type == "linked_list" {
		fmt.Fprintf(&b, "  direction: %s\n", r.Storage.Direction)
	}
	fmt.Fprintf(&b, "operations:\n")
	fmt.Fprintf(&b, "  add: true (fixed)\n")
	fmt.Fprintf(&b, "  get: true (fixed)\n")
	fmt.Fprintf(&b, "  size: true (fixed)\n")
	fmt.Fprintf(&b, "  remove: %t\n", r.Operations.Remove)
	fmt.Fprintf(&b, "  insert: %t\n", r.Operations.Insert)
	fmt.Fprintf(&b, "  linear_search: %t\n", r.Operations.LinearSearch)
	fmt.Fprintf(&b, "  contains: %t\n", r.Operations.Contains)
	fmt.Fprintf(&b, "  sort: %t\n", r.Operations.Sort)
	if r.Operations.Sort {
		fmt.Fprintf(&b, "    algorithm: %s\n", r.Operations.SortAlgorithm)
	}
	fmt.Fprintf(&b, "  iteration: %t\n", r.Operations.Iteration)
	if r.Operations.Iteration {
		fmt.Fprintf(&b, "    direction: %s\n", r.Operations.IterationDir)
	}
	return b.String()
}
