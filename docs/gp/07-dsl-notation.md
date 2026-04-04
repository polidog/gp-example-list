# ドメイン固有記法: リストコンテナ

## 概要

### 記法の目的

リストコンテナのフィーチャー選択とパラメータ指定を宣言的に記述し、AIジェネレータに対してコード生成を「発注」するための記法。この仕様書1つで、どのような構成のリストコンテナを、どの言語で生成するかが完全に決定される。

### 利用者

- **ライブラリ開発者**: 必要な構成のリストコンテナを仕様書として記述し、コード生成を依頼する
- **GP学習者**: フィーチャー選択を変えながら生成結果の違いを観察する

### フォーマット

**YAML** — 宣言的で人間が読みやすく、AIが解析しやすい。コメントも記述可能。

## 構文定義

### 全体構造

```yaml
# リストコンテナ仕様書
list:
  # 必須フィーチャー
  element_type: <型名>        # 要素の型（ジェネリクス型パラメータ）
  ownership: <所有権モデル>    # ExternalReference | OwnedReference | Copy
  morphology: <形態>          # Monomorphic | Polymorphic

  # 選択フィーチャー
  length_counter:              # 省略時: 無効
    enabled: true
    counter_type: <カウンタ型>  # short | int | long（enabled: true 時に必須）
  tracing:                     # 省略時: 無効
    enabled: true

  # 生成オプション
  target_language: <言語名>    # 生成対象言語
```

### 記述項目

| 項目 | 型 | 必須 | 既定値 | 対応するフィーチャー | 説明 |
|------|-----|------|--------|-------------------|------|
| `list.element_type` | string | はい | — | ElementType | リストが保持する要素の型パラメータ名。ジェネリクスの型引数として使用される |
| `list.ownership` | enum | はい | — | Ownership | 要素の所有権モデル |
| `list.morphology` | enum | はい | — | Morphology | 要素の型の扱い方（単相/多相） |
| `list.length_counter` | object | いいえ | 無効 | LengthCounter | 長さカウンタの設定。省略またはキー自体がない場合は無効 |
| `list.length_counter.enabled` | boolean | いいえ | `false` | LengthCounter | カウンタの有効/無効 |
| `list.length_counter.counter_type` | enum | 条件付き | — | CounterType | カウンタの整数型。`enabled: true` の場合に必須 |
| `list.tracing` | object | いいえ | 無効 | Tracing | トレーシングの設定。省略またはキー自体がない場合は無効 |
| `list.tracing.enabled` | boolean | いいえ | `false` | Tracing | トレーシングの有効/無効 |
| `list.target_language` | string | はい | — | （生成オプション） | 生成対象のプログラミング言語 |

### 選択肢のある項目

| 項目 | 選択肢 | 説明 |
|------|--------|------|
| `ownership` | `ExternalReference` | リストは要素への参照を保持するが、要素の生存期間は管理しない |
| | `OwnedReference` | リストは要素への参照を保持し、リスト破棄時に要素も解放する |
| | `Copy` | リストは要素のコピーを値として保持する |
| `morphology` | `Monomorphic` | すべての要素が同一の型Tでなければならない |
| | `Polymorphic` | Tのサブタイプも格納可能 |
| `counter_type` | `short` | 小さな整数型（最大約32K〜65K要素） |
| | `int` | 標準整数型（最大約21億要素） |
| | `long` | 大きな整数型（事実上無制限） |
| `target_language` | `cpp` | C++（テンプレートを使用） |
| | `java` | Java（ジェネリクスを使用） |
| | `go` | Go（type parameterを使用） |
| | `python` | Python（型ヒントを使用） |
| | `rust` | Rust（ジェネリクス + 所有権システムを活用） |
| | その他 | 上記以外の言語も指定可能 |

### 簡略記法

選択フィーチャーが無効の場合、キー自体を省略できる:

```yaml
# 完全記法
list:
  element_type: T
  ownership: Copy
  morphology: Monomorphic
  length_counter:
    enabled: false
  tracing:
    enabled: false
  target_language: cpp

# 簡略記法（上と同義）
list:
  element_type: T
  ownership: Copy
  morphology: Monomorphic
  target_language: cpp
```

## 制約ルール

仕様書は以下の制約を満たさなければならない。バリデーション時にこれらを検証する。

| # | 制約 | 条件 | 説明 |
|---|------|------|------|
| C1 | Polymorphic + Copy 排他 | `morphology: Polymorphic` かつ `ownership: Copy` は不可 | 値コピーではサブタイプのスライシングが発生し多態性を実現できない |
| C2 | CounterType 必須条件 | `length_counter.enabled: true` の場合、`counter_type` の指定が必須 | カウンタ有効時にカウンタ型が未指定だと生成できない |
| C3 | CounterType 無効条件 | `length_counter.enabled: false`（または省略）の場合、`counter_type` は指定不可 | 無効なカウンタに型を指定しても意味がない |
| C4 | 必須フィーチャー | `element_type`, `ownership`, `morphology`, `target_language` は省略不可 | これらがないと構成を決定できない |
| C5 | Ownership XOR | `ownership` は3つの選択肢から正確に1つ | 複数の所有権モデルは混在不可 |
| C6 | Morphology XOR | `morphology` は2つの選択肢から正確に1つ | 単相と多相は排他 |

## バリデーションエラーメッセージ

| エラーコード | メッセージ | 発生条件 |
|------------|-----------|---------|
| `E001` | `ownership: Copy と morphology: Polymorphic は同時に指定できません` | C1違反 |
| `E002` | `length_counter.enabled が true の場合、counter_type は必須です` | C2違反 |
| `E003` | `length_counter が無効の場合、counter_type は指定できません` | C3違反 |
| `E004` | `{項目名} は必須項目です` | C4違反 |
| `E005` | `ownership の値が不正です: {値}。ExternalReference, OwnedReference, Copy のいずれかを指定してください` | C5違反 |
| `E006` | `morphology の値が不正です: {値}。Monomorphic, Polymorphic のいずれかを指定してください` | C6違反 |

## サンプル仕様書

サンプル仕様書は `docs/gp/08-dsl-examples/` に配置する。
