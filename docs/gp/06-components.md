# 実装コンポーネント: リストコンテナ

## コンポーネント一覧

### ListCore（リストコア）
- **責務**: リストの基本構造（head参照、ノード連結）と共通操作（挿入・削除・走査・検索）のアルゴリズムを提供する
- **共通/可変**: 共通（すべての構成で存在する骨格部分）
- **依存先**: NodeDefinition
- **対応するフィーチャー**: LIST（ルート）

### NodeDefinition（ノード定義）
- **責務**: リストの内部ノード構造体を定義する。要素の保持方法（値/ポインタ）と型制約を決定する
- **共通/可変**: 可変
- **可変の場合**: Ownership × Morphology の5つの有効な組み合わせに対応する5バリアントを持つ。生成時に1つが選択される
- **依存先**: なし
- **対応するフィーチャー**: Ownership, Morphology

### MemoryManagement（メモリ管理）
- **責務**: ノードの生成・破棄時のメモリ確保/解放ロジックを提供する。Ownershipに応じて解放責任が異なる
- **共通/可変**: 可変
- **可変の場合**: Ownershipの3バリアントに対応
  - Copy: ノード破棄時に値も消滅（特別な解放不要）
  - OwnedReference: ノード破棄時に参照先の要素も解放
  - ExternalReference: ノード破棄時にポインタのみ破棄（要素は解放しない）
- **依存先**: NodeDefinition
- **対応するフィーチャー**: Ownership

### LengthCounterMixin（長さカウンタミックスイン）
- **責務**: List構造体にカウンタフィールドを追加し、挿入・削除操作にカウンタ更新コードを織り込む
- **共通/可変**: 可変（オプショナル — LengthCounter有効時のみ生成）
- **可変の場合**: 有効/無効の切り替え + CounterType（short/int/long）の選択
- **依存先**: ListCore
- **対応するフィーチャー**: LengthCounter, CounterType

### TracingMixin（トレーシングミックスイン）
- **責務**: 各操作の前後にトレースログ出力コードを織り込む
- **共通/可変**: 可変（オプショナル — Tracing有効時のみ生成）
- **可変の場合**: 有効/無効の切り替え
- **依存先**: ListCore
- **対応するフィーチャー**: Tracing

### PublicAPI（公開API）
- **責務**: 利用者に公開するメソッド/関数シグネチャを定義する。生成対象言語のイディオムに沿ったインターフェースを提供する
- **共通/可変**: 共通（シグネチャの骨格は共通。LengthCounter有効時にlength()メソッドが追加される）
- **依存先**: ListCore, LengthCounterMixin（有効時）
- **対応するフィーチャー**: LIST, LengthCounter

## コンポーネント依存図

```
┌────────────┐
│ PublicAPI   │ ← 利用者に公開されるインターフェース
└─────┬──────┘
      │ 依存
      v
┌────────────┐     ┌─────────────────────┐
│  ListCore  │<────│ LengthCounterMixin  │ [オプショナル]
│  (共通)    │<────│ TracingMixin        │ [オプショナル]
└─────┬──────┘     └─────────────────────┘
      │ 依存
      v
┌────────────────┐
│ NodeDefinition │ ← Ownership × Morphology で決定
└─────┬──────────┘
      │ 依存
      v
┌──────────────────┐
│ MemoryManagement │ ← Ownership で決定
└──────────────────┘
```

## インターフェース定義

### INodeDefinition（ノード定義インターフェース）
- **目的**: 要素の保持方法と型制約をノード構造体として表現する
- **対応する可変ポイント**: Ownership × Morphology
- **バリアント**:
  - `CopyMonomorphicNode<T>` — 値埋め込み、単一型
  - `ExtRefMonomorphicNode<T>` — 外部参照、単一型
  - `OwnedRefMonomorphicNode<T>` — 所有参照、単一型
  - `ExtRefPolymorphicNode<T>` — 外部参照、多態型
  - `OwnedRefPolymorphicNode<T>` — 所有参照、多態型

### IMemoryManagement（メモリ管理インターフェース）
- **目的**: ノードと要素の生成・解放方式を抽象化する
- **対応する可変ポイント**: Ownership
- **バリアント**:
  - `CopyMemory` — 値のコピーで生成、ノード解放のみ
  - `OwnedRefMemory` — ポインタで生成、ノード解放時に要素も解放
  - `ExtRefMemory` — ポインタで生成、ノード解放時に要素は解放しない

### ILengthCounter（長さカウンタインターフェース）
- **目的**: カウンタの有無とカウンタ型を抽象化する
- **対応する可変ポイント**: LengthCounter, CounterType
- **バリアント**:
  - `NoCounter` — カウンタなし（要素数取得はO(n)走査）
  - `ShortCounter` — short型カウンタ
  - `IntCounter` — int型カウンタ
  - `LongCounter` — long型カウンタ

### ITracing（トレーシングインターフェース）
- **目的**: トレース出力の有無を抽象化する
- **対応する可変ポイント**: Tracing
- **バリアント**:
  - `NoTracing` — トレースなし
  - `StdoutTracing` — 標準出力へのトレース

## フィーチャーからコンポーネントへの対応表

| フィーチャー | 影響するコンポーネント | 影響の内容 |
|-------------|---------------------|-----------|
| ElementType | NodeDefinition, PublicAPI | 型パラメータTの指定 |
| Ownership (ExternalReference) | NodeDefinition, MemoryManagement | 外部参照ノード + 要素非解放 |
| Ownership (OwnedReference) | NodeDefinition, MemoryManagement | 所有参照ノード + 要素解放 |
| Ownership (Copy) | NodeDefinition, MemoryManagement | 値埋め込みノード + コピー生成 |
| Morphology (Monomorphic) | NodeDefinition | 具象型Tのみ許容 |
| Morphology (Polymorphic) | NodeDefinition | Tのサブタイプも許容 |
| LengthCounter (有効) | LengthCounterMixin, ListCore, PublicAPI | カウンタフィールド追加 + 更新ロジック織り込み + length()メソッド追加 |
| CounterType | LengthCounterMixin | カウンタフィールドの型決定 |
| Tracing (有効) | TracingMixin, ListCore | 全操作へのログ出力フック織り込み |
