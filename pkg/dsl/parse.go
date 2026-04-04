package dsl

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ParseFile はYAMLファイルを読み込み、Spec に変換する。
func ParseFile(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗しました: %w", err)
	}
	return Parse(data)
}

// Parse はYAMLバイト列を Spec に変換する。
func Parse(data []byte) (*Spec, error) {
	var spec Spec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("YAMLのパースに失敗しました: %w", err)
	}
	return &spec, nil
}
