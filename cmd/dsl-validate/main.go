package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/polidog/gp-example-list/pkg/dsl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "使い方: dsl-validate <仕様書.yaml> [--resolve]")
		os.Exit(1)
	}

	path := os.Args[1]
	showResolve := len(os.Args) >= 3 && os.Args[2] == "--resolve"

	spec, err := dsl.ParseFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	errs := dsl.Validate(spec)
	if len(errs) > 0 {
		fmt.Fprintln(os.Stderr, "バリデーションエラー:")
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  %s\n", e.Error())
		}
		os.Exit(1)
	}

	fmt.Println("OK: 仕様書は有効です")

	if showResolve {
		cfg := dsl.Resolve(spec)
		data, _ := json.MarshalIndent(cfg, "", "  ")
		fmt.Println("\n構成解決結果:")
		fmt.Println(string(data))
	}
}
