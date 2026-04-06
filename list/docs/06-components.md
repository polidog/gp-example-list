# 実装コンポーネント: リストコンテナ

## コンポーネント一覧

### ListStruct（リスト構造体）
- **責務**: リストの型定義。内部データ構造とメタデータ（サイズ等）を保持する
- **共通/可変**: 可変（Storageに依存）
- **バリアント**:
  - ArrayBased: 内部配列 + サイズ + 容量
  - LinkedList(Singly): headポインタ + サイズ
  - LinkedList(Doubly): head/tailポインタ + サイズ
- **依存先**: なし
- **対応するフィーチャー**: Storage, ElementType

### NodeStruct（ノード構造体）
- **責務**: 連結リストの各要素を表すノードの型定義
- **共通/可変**: 可変（LinkedList選択時のみ生成）
- **バリアント**:
  - Singly: value + next
  - Doubly: value + next + prev
- **依存先**: なし
- **対応するフィーチャー**: LinkedList, Direction

### Constructor（コンストラクタ）
- **責務**: リストの初期化。データ構造に応じた初期状態を設定する
- **共通/可変**: 可変（Storageに依存）
- **バリアント**:
  - ArrayBased: 配列の初期確保（InitialCapacity分）
  - LinkedList: ポインタのnil初期化
- **依存先**: ListStruct
- **対応するフィーチャー**: Storage, InitialCapacity

### Add（末尾追加）
- **責務**: リストの末尾に要素を追加する
- **共通/可変**: 共通（インターフェース固定、実装は可変）
- **バリアント**:
  - ArrayBased: 容量チェック → 拡張 → 末尾に代入
  - LinkedList(Singly): 末尾走査 → ノード追加
  - LinkedList(Doubly): tail参照 → ノード追加
- **依存先**: ListStruct, GrowthStrategy（ArrayBasedの場合）
- **対応するフィーチャー**: CoreOperations.Add

### Get（インデックス取得）
- **責務**: 指定インデックスの要素を返す
- **共通/可変**: 共通（インターフェース固定、実装は可変）
- **バリアント**:
  - ArrayBased: 配列の直接アクセス O(1)
  - LinkedList: 先頭からの走査 O(n)
- **依存先**: ListStruct
- **対応するフィーチャー**: CoreOperations.Get

### Size（サイズ取得）
- **責務**: リストの現在の要素数を返す
- **共通/可変**: 共通（全バリアントで同一実装）
- **依存先**: ListStruct
- **対応するフィーチャー**: CoreOperations.Size

### GrowthLogic（拡張ロジック）
- **責務**: 配列の容量不足時に拡張する
- **共通/可変**: 可変（GrowthStrategyに依存、ArrayBasedのみ）
- **バリアント**:
  - Doubling: 容量を2倍に拡張
  - Additive: 容量を固定値分追加
- **依存先**: ListStruct
- **対応するフィーチャー**: GrowthStrategy

### Remove（削除）
- **責務**: 指定インデックスの要素を削除する
- **共通/可変**: 可変（オプション、Storageに依存）
- **バリアント**:
  - ArrayBased: 要素シフト O(n)
  - LinkedList: ポインタ付け替え O(n)走査 + O(1)削除
- **依存先**: ListStruct
- **対応するフィーチャー**: OptionalOperations.Remove

### Insert（挿入）
- **責務**: 指定位置に要素を挿入する
- **共通/可変**: 可変（オプション、Storageに依存）
- **バリアント**:
  - ArrayBased: 要素シフト → 挿入
  - LinkedList: 走査 → ポインタ付け替え
- **依存先**: ListStruct, GrowthLogic（ArrayBasedの場合）
- **対応するフィーチャー**: OptionalOperations.Insert

### LinearSearch（線形検索）
- **責務**: リスト内の要素を先頭から順に検索し、インデックスを返す
- **共通/可変**: 可変（オプション）
- **依存先**: ListStruct, Get
- **対応するフィーチャー**: OptionalOperations.Search.LinearSearch

### Contains（存在確認）
- **責務**: 指定値がリストに含まれるかを返す
- **共通/可変**: 可変（オプション）
- **依存先**: LinearSearch
- **対応するフィーチャー**: OptionalOperations.Search.Contains

### SortImpl（ソート実装）
- **責務**: リストの要素を昇順にソートする
- **共通/可変**: 可変（オプション、SortAlgorithmに依存）
- **バリアント**:
  - InsertionSort: O(n²) だがシンプルで小規模リスト向き
  - MergeSort: O(n log n) で安定ソート
- **依存先**: ListStruct
- **対応するフィーチャー**: OptionalOperations.Sort

### Iterator（イテレータ）
- **責務**: リストの要素を順番に走査する機構を提供する
- **共通/可変**: 可変（オプション、Directionに依存）
- **バリアント**:
  - Forward: 先頭→末尾の順に走査
  - Reverse: 末尾→先頭の順に走査
- **依存先**: ListStruct
- **対応するフィーチャー**: OptionalOperations.Iteration

### ContractDeclaration（コントラクト宣言）
- **責務**: 言語ネイティブのinterface/protocol宣言と適合メソッドの生成
- **共通/可変**: 可変（言語 + 有効フィーチャーの組み合わせで導出）
- **バリアント**:
  - PHP: `implements \Countable, \IteratorAggregate` + `count()`, `getIterator()` + `{Name}Interface`（別ファイル）
  - Go: `{Name}er` interface型定義 + `All() iter.Seq[T]`（Iteration有効時）
  - Python: `__len__`, `__iter__`, `__contains__` 等のダンダーメソッド
  - TypeScript: `implements Iterable<T>` + `[Symbol.iterator]()`
- **依存先**: ListStruct, 全有効操作
- **対応するフィーチャー**: 導出（フィーチャー選択×言語から自動決定。`09-configuration-knowledge.md` の言語別コントラクトマッピングルール参照）

## コンポーネント依存図

```
                    ┌─────────────┐
                    │ ListStruct  │
                    └──────┬──────┘
           ┌───────────────┼───────────────┐
           │               │               │
     ┌─────┴─────┐   ┌────┴────┐   ┌──────┴──────┐
     │Constructor│   │NodeStruct│   │GrowthLogic  │
     └───────────┘   │(LinkedList)│  │(ArrayBased) │
                     └───────────┘   └─────────────┘
           │
    ┌──────┼──────┬──────────┬──────────┐
    │      │      │          │          │
 ┌──┴─┐ ┌─┴──┐ ┌─┴──┐  ┌───┴───┐ ┌───┴────┐
 │Add │ │Get │ │Size│  │Remove │ │Insert  │
 └────┘ └─┬──┘ └────┘  └───────┘ └────────┘
           │
    ┌──────┼──────────┐
    │      │          │
 ┌──┴──────┴──┐  ┌───┴────┐  ┌────────┐
 │LinearSearch│  │Contains│  │SortImpl│
 └────────────┘  └────────┘  └────────┘
           │
     ┌─────┴─────┐
     │ Iterator  │
     └───────────┘
```

## インターフェース定義

### Storage インターフェース（概念的）
- **目的**: 内部データ構造の差し替えポイント。生成時にフラグメント選択で実現
- **対応する可変ポイント**: Storage（XOR: ArrayBased / LinkedList）
- **現在の実装**: ArrayBased
- **想定される追加実装**: LinkedList(Singly), LinkedList(Doubly)

### Operation インターフェース（概念的）
- **目的**: オプション操作の追加ポイント。生成時にフラグメント追加で実現
- **対応する可変ポイント**: OptionalOperations（OR）
- **現在の実装**: Remove, Insert, Search(LinearSearch, Contains), Iteration(Forward)
- **想定される追加実装**: Sort(InsertionSort / MergeSort), Iteration(Reverse)
