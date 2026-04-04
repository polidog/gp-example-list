# gp-example-list

ジェネレーティブプログラミング（DEMRAL）のプロセスを現代のAI技術で実践する実験的プロジェクト。

Czarneckiらの「Generative Programming」で紹介されるリストコンテナの例題をベースに、フィーチャーモデルの定義からDSLによる仕様記述、AIによるコード生成までの一連のGPプロセスを体験できる。

## フィーチャーモデル

```
LIST
├── ElementType    [必須] ── ジェネリック型パラメータ <T>
├── Ownership      [必須] ── XOR: ExternalReference | OwnedReference | Copy
├── Morphology     [必須] ── XOR: Monomorphic | Polymorphic
├── LengthCounter  [選択] ── CounterType: short | int | long
└── Tracing        [選択] ── 操作ログの有無
```

**制約**: `Polymorphic + Copy` は不可（値コピーではサブタイプの多態性を実現できない）

有効な構成は **40通り**。

## 使い方

### 1. DSL仕様書を書く

```yaml
# my-list.yaml
list:
  element_type: T
  ownership: Copy
  morphology: Monomorphic
  length_counter:
    enabled: true
    counter_type: int
  target_language: cpp
```

### 2. バリデーションと構成解決

```bash
go run ./cmd/dsl-validate/ my-list.yaml --resolve
```

### 3. コード生成

構成の知識（`docs/gp/09-configuration-knowledge.md`）のルール R1〜R6 に基づき、Claude Code にコード生成を依頼する。

## プロジェクト構造

```
docs/gp/                  # DEMRALプロセスの成果物
├── 01-domain-scope.md        # ドメインスコープ
├── 02-feature-model.md       # フィーチャーモデル
├── 03-concept-model.md       # 概念モデル
├── 04-constraints.md         # フィーチャー間制約
├── 05-architecture.md        # アーキテクチャ設計
├── 06-components.md          # 実装コンポーネント定義
├── 07-dsl-notation.md        # DSL記法定義
├── 08-dsl-examples/          # DSLサンプル仕様書（5件）
└── 09-configuration-knowledge.md  # 構成の知識（生成ルール R1〜R6）

pkg/dsl/                  # DSLバリデータ・構成解決（Go）
cmd/dsl-validate/         # CLIツール

generated/                # 生成されたリストコンテナ
├── cpp/                      # C++（Copy+Mono, コンパイル確認済み）
├── java/                     # Java（OwnedRef+Poly+Counter+Tracing）
├── go/                       # Go（2バリアント, テスト付き）
└── php/                      # PHP（Copy+Mono, テスト付き）

.claude/skills/           # Claude Code スキル（DEMRALの各活動）
```

## DEMRALプロセス

本プロジェクトではDEMRALの9つの開発活動をClaude Codeのスキルとして実装している。

| 活動 | スキル | 成果物 |
|------|--------|--------|
| ドメインスコーピング | `/domain-scoping` | `01-domain-scope.md` |
| フィーチャーモデリング | `/feature-modeling` | `02〜04` |
| ADT仕様定義 | `/adt-specification` | `03a-adt-specification.md` |
| アーキテクチャ設計 | `/architecture-design` | `05〜06` |
| DSL定義 | `/dsl-definition` | `07〜08` |
| 構成の知識定義 | `/configuration-knowledge` | `09` |
| DSL実装 | `/dsl-implementation` | `pkg/dsl/` |
| コンポーネント実装 | `/component-implementation` | `generated/` |
| ジェネレータ実装 | `/generator-implementation` | 構成の知識 + AI |

これらは厳密な順序ではなく、任意の順番で何度でも繰り返せる。

## 開発

```bash
# テスト実行
go test ./...

# DSL仕様書のバリデーション
go run ./cmd/dsl-validate/ docs/gp/08-dsl-examples/example-01-minimal.yaml

# 構成解決結果の確認
go run ./cmd/dsl-validate/ docs/gp/08-dsl-examples/example-02-full-features.yaml --resolve
```

## 設計思想

- **変化への対応 > 再利用性**: 単一アプリケーションが時間とともにどう進化するかを構造化する
- **共通性と可変性の対比分析**: 各要素を「共通（変わらない部分）」と「可変（変わりうる部分）」に分類する
- **AIによるコード生成の前提**: フィーチャーモデルとDSLを宣言的な仕様として扱い、AIが構成の知識に従ってコードを生成する
