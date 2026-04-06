# gp-example-list

ジェネレーティブプログラミング（DEMRAL）のプロセスを現代のAI技術で実践する実験的プロジェクト。

Czarneckiらの「Generative Programming」で紹介されるリストコンテナの例題をベースに、フィーチャーモデルの定義からDSLによる仕様記述、AIによるコード生成までの一連のGPプロセスを体験できる。

## フィーチャーモデル

```
ListContainer
├── ElementType [パラメータ型] ── 要素の型（int, string 等）
├── Storage [XOR]
│   ├── ArrayBased
│   │   ├── InitialCapacity [パラメータ] = 16
│   │   └── GrowthStrategy [XOR]: Doubling | Additive
│   └── LinkedList
│       └── Direction [XOR]: Singly | Doubly
├── CoreOperations [固定] ── Add, Get, Size
└── OptionalOperations [OR]
    ├── Remove
    ├── Insert
    ├── Search (LinearSearch, Contains)
    ├── Sort (InsertionSort | MergeSort)
    └── Iteration (Forward | Reverse)
```

言語ネイティブのコントラクト（PHP interface, Go interface, Python protocol 等）はフィーチャー選択と対象言語から自動導出される。

## 使い方

### 1. DSL仕様書を書く

```yaml
# my-list.yaml
name: IntList
language: go
element_type: int

storage:
  type: array
  initial_capacity: 16
  growth_strategy: doubling

operations:
  - remove
  - insert
  - linear_search
  - contains
  - iteration:
      direction: forward
```

### 2. バリデーションと構成解決

```bash
go run ./list/cmd/dsl-validate/ my-list.yaml --resolve
```

### 3. コード生成

構成の知識（`list/docs/09-configuration-knowledge.md`）のマッピングルールに基づき、Claude Code にコード生成を依頼する。

```bash
# Claude Code のスキルを使う場合
/implement my-list.yaml
```

## プロジェクト構造

```
list/
├── docs/                         # DEMRALプロセスの成果物
│   ├── 01-domain-scope.md            # ドメインスコープ
│   ├── 02-feature-model.md           # フィーチャーモデル
│   ├── 04-constraints.md             # フィーチャー間制約
│   ├── 05-architecture.md            # アーキテクチャ設計
│   ├── 06-components.md              # 実装コンポーネント定義
│   ├── 07-dsl-notation.md            # DSL記法定義
│   ├── 08-dsl-examples/              # DSLサンプル仕様書
│   └── 09-configuration-knowledge.md # 構成の知識（マッピングルール）
├── pkg/dsl/                      # DSLバリデータ・構成解決（Go）
├── cmd/dsl-validate/             # CLIツール
└── generated/                    # 生成されたリストコンテナ
    ├── go/                           # Go（配列ベース + iter.Seq対応）
    └── php/                          # PHP（Countable, IteratorAggregate対応）

.claude/skills/                   # Claude Code スキル（DEMRALの各活動）
```

## DEMRALプロセス

本プロジェクトではDEMRALの開発活動をClaude Codeのスキルとして実装している。

| 活動 | スキル | 成果物 |
|------|--------|--------|
| ドメインスコーピング | `/domain-scoping` | `01-domain-scope.md` |
| フィーチャーモデリング | `/feature-modeling` | `02`, `04` |
| ドメイン設計 | `/domain-design` | `05`, `06`, `09` |
| DSL定義 | `/dsl-definition` | `07`, `08` |
| 実装 | `/implement` | `generated/`, `list/pkg/dsl/` |
| テスト | `/test` | テストコード |

これらは厳密な順序ではなく、任意の順番で何度でも繰り返せる。

## 開発

```bash
# テスト実行
go test ./...

# DSL仕様書のバリデーション
go run ./list/cmd/dsl-validate/ list/docs/08-dsl-examples/basic-array-list.yaml

# 構成解決結果の確認
go run ./list/cmd/dsl-validate/ list/docs/08-dsl-examples/basic-array-list.yaml --resolve
```

## 設計思想

- **変化への対応 > 再利用性**: 単一アプリケーションが時間とともにどう進化するかを構造化する
- **共通性と可変性の対比分析**: 各要素を「共通（変わらない部分）」と「可変（変わりうる部分）」に分類する
- **AIによるコード生成の前提**: フィーチャーモデルとDSLを宣言的な仕様として扱い、AIが構成の知識に従ってコードを生成する
- **言語非依存のDSL**: 仕様は言語に依存せず、コントラクト（interface/protocol）は構成の知識により自動導出される
