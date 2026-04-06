# ドメイン固有記法: リストコンテナ

## 概要

### 記法の目的
リストコンテナの構成（データ構造・操作・パラメータ）を宣言的に記述し、AIによるコード生成の入力仕様とする。

### 利用者
- ジェネレーティブプログラミング学習者（仕様を書いてコード生成を体験する）
- ソフトウェアエンジニア（必要な構成を宣言し、コードを得る）

### フォーマット
YAML

選定理由:
- 人間が読み書きしやすい
- 構造化データを自然に表現できる
- Go の既存バリデータ（`gopkg.in/yaml.v3`）で解析可能

## 構文定義

### 全体構造

```yaml
name: <string>           # リストの名前（生成される型名に使用）
language: <string>        # 生成先の言語
element_type: <string>    # 要素の型

storage:
  type: <string>          # データ構造の種類
  # 以下は type に応じた追加設定

operations:               # 有効にする操作のリスト
  - <operation>
  - <operation>:
      <option>: <value>
```

### 記述項目

| 項目 | 型 | 必須 | 既定値 | 対応するフィーチャー | 説明 |
|------|-----|------|--------|-------------------|------|
| `name` | string | はい | — | — | リストの名前。生成される型名・ファイル名に使用 |
| `language` | string | はい | — | — | コード生成先の言語 |
| `element_type` | string | はい | — | ElementType | 要素の型。言語の型名で指定 |
| `storage` | object | いいえ | `{type: array}` | Storage | 内部データ構造の設定 |
| `storage.type` | string | いいえ | `array` | Storage | データ構造の種類 |
| `storage.initial_capacity` | integer | いいえ | `16` | InitialCapacity | 配列の初期容量（`type: array` 時のみ） |
| `storage.growth_strategy` | string | いいえ | `doubling` | GrowthStrategy | 配列の拡張戦略（`type: array` 時のみ） |
| `storage.growth_increment` | integer | いいえ | `16` | GrowthStrategy.Additive | 加算する固定値（`growth_strategy: additive` 時のみ） |
| `storage.direction` | string | いいえ | `singly` | Direction | 連結リストの方向（`type: linked_list` 時のみ） |
| `operations` | list | いいえ | 下記参照 | OptionalOperations | 有効にするオプション操作のリスト |

**注意**: `add`, `get`, `size` は常に生成される基本操作であり、`operations` に記述する必要はない。

### 選択肢のある項目

| 項目 | 選択肢 | 説明 |
|------|--------|------|
| `language` | `go`, `python`, `typescript`, ... | コード生成先の言語。拡張可能 |
| `storage.type` | `array`, `linked_list` | `array`: 配列ベース実装 / `linked_list`: 連結リスト実装 |
| `storage.growth_strategy` | `doubling`, `additive` | `doubling`: 容量2倍 / `additive`: 固定値加算 |
| `storage.direction` | `singly`, `doubly` | `singly`: 単方向 / `doubly`: 双方向 |
| `operations[].sort.algorithm` | `insertion_sort`, `merge_sort` | ソートアルゴリズムの選択 |
| `operations[].iteration.direction` | `forward`, `reverse` | イテレーション方向の選択 |

### operations の記述形式

操作はシンプル形式（文字列）またはオプション付き形式（オブジェクト）で指定できる。

```yaml
# シンプル形式
operations:
  - remove
  - insert
  - linear_search
  - contains

# オプション付き形式
operations:
  - sort:
      algorithm: merge_sort
  - iteration:
      direction: forward
```

### operations 一覧

| 操作名 | オプション | 説明 |
|--------|-----------|------|
| `remove` | なし | インデックス指定で要素を削除 |
| `insert` | なし | 指定位置に要素を挿入 |
| `linear_search` | なし | 線形検索（インデックスを返す） |
| `contains` | なし | 要素の存在確認（bool を返す） |
| `sort` | `algorithm` | ソート |
| `iteration` | `direction` | イテレータによる走査 |

## 制約ルール

| 制約 | 条件 | 説明 |
|------|------|------|
| storage-array-only | `initial_capacity` と `growth_strategy` は `storage.type: array` のときのみ有効 | 配列ベース固有の設定 |
| storage-linked-only | `direction` は `storage.type: linked_list` のときのみ有効 | 連結リスト固有の設定 |
| growth-additive-only | `growth_increment` は `growth_strategy: additive` のときのみ有効 | 加算戦略の固有パラメータ |
| contains-requires-search | `contains` を指定する場合、`linear_search` も必要（自動補完される） | Contains は LinearSearch に依存 |
| sort-requires-algorithm | `sort` を指定する場合、`algorithm` が必須 | アルゴリズム未指定のソートは生成不可 |
| iteration-requires-direction | `iteration` を指定する場合、`direction` が必須 | 方向未指定のイテレーションは生成不可 |
| initial-capacity-range | `initial_capacity` は 1 以上 65536 以下 | 過小・過大な初期確保の防止 |
| singly-reverse-warning | `storage.direction: singly` かつ `iteration.direction: reverse` は非推奨 | 単方向リストの逆走査は O(n) で非効率 |

## デフォルト動作

`operations` を省略した場合のデフォルト構成:

```yaml
operations:
  - remove
  - insert
  - linear_search
  - contains
  - iteration:
      direction: forward
```

## サンプル仕様書

サンプルは `08-dsl-examples/` ディレクトリを参照。
