# 構成の知識: リストコンテナ

## マッピングルール

### フィーチャー → コンポーネント

| DSL項目 | 選択肢 | 生成されるコンポーネント | 備考 |
|---------|--------|------------------------|------|
| `language` | `go`, `python`, `typescript`, ... | 全コンポーネントの言語固有コード | 生成先言語を決定 |
| `element_type` | 任意の型名 | ListStruct の型パラメータ | 言語のジェネリクス対応に応じて変換 |
| `storage.type` | `array` | ListStruct(配列), Constructor(配列初期化), GrowthLogic | |
| `storage.type` | `linked_list` | ListStruct(ポインタ), NodeStruct, Constructor(nil初期化) | |
| `storage.growth_strategy` | `doubling` | GrowthLogic(2倍拡張) | `storage.type: array` 時のみ |
| `storage.growth_strategy` | `additive` | GrowthLogic(固定値加算) | `storage.type: array` 時のみ |
| `storage.direction` | `singly` | NodeStruct(next), ListStruct(head) | `storage.type: linked_list` 時のみ |
| `storage.direction` | `doubly` | NodeStruct(next+prev), ListStruct(head+tail) | `storage.type: linked_list` 時のみ |
| （常に生成） | — | Add, Get, Size | CoreOperations。DSLに記述不要 |
| `operations[]`: `remove` | 有効/無効 | Remove | |
| `operations[]`: `insert` | 有効/無効 | Insert | |
| `operations[]`: `linear_search` | 有効/無効 | LinearSearch | |
| `operations[]`: `contains` | 有効/無効 | Contains | `linear_search` に依存（自動補完） |
| `operations[]`: `sort` | 有効/無効 | SortImpl | `algorithm` 指定が必須 |
| `sort.algorithm` | `insertion_sort` | SortImpl(挿入ソート) | `sort` 有効時のみ |
| `sort.algorithm` | `merge_sort` | SortImpl(マージソート) | `sort` 有効時のみ |
| `operations[]`: `iteration` | 有効/無効 | Iterator | `direction` 指定が必須 |
| `iteration.direction` | `forward` | Iterator(前方走査) | `iteration` 有効時のみ |
| `iteration.direction` | `reverse` | Iterator(後方走査) | `iteration` 有効時のみ |

### 言語別変換ルール

| DSL項目 | Go | Python | TypeScript |
|---------|-----|--------|------------|
| `element_type` | ジェネリクス `[T any]` または具体型 | 型ヒント `list[int]` 等 | ジェネリクス `<T>` |
| `name` | 型名（PascalCase） | クラス名（PascalCase） | クラス名（PascalCase） |
| ファイル名 | `{snake_case}.go` | `{snake_case}.py` | `{camelCase}.ts` |
| テストFW | `testing/quick` or `pgregory.net/rapid` | `hypothesis` | `fast-check` |

### パラメータ → 設定値

| DSL項目 | 対応する設定 | 変換ルール |
|---------|------------|-----------|
| `element_type` | 構造体の型パラメータまたは具体型 | 言語別変換ルール参照 |
| `storage.initial_capacity` | 配列の初期確保サイズ（定数） | 正の整数をそのまま定数として埋め込む |
| `storage.growth_increment` | 加算拡張時の増分値 | 正の整数を定数として埋め込む（デフォルト: 16） |

### 言語別コントラクトマッピングルール

フィーチャー選択と対象言語の組み合わせから、言語ネイティブのコントラクト（interface / protocol）を**自動的に導出**する。DSLには記述しない。コード生成時に以下のルールを適用する。

#### PHP

| 条件 | コントラクト | 追加コンポーネント | 備考 |
|------|------------|-----------------|------|
| 常時 | `\Countable` | `count(): int`（`size()` に委譲） | `count($list)` で要素数取得可能に |
| Iteration有効 | `\IteratorAggregate` | `getIterator(): \Traversable` | `foreach` で走査可能に |
| 常時 | `{Name}Interface` | 全publicメソッドのinterface（別ファイル） | 差し替え可能性のための契約 |

#### Go

| 条件 | コントラクト | 追加コンポーネント | 備考 |
|------|------------|-----------------|------|
| 常時 | `{Name}er` interface型定義 | 全公開メソッドのinterface型ブロック | 構造的型付けなので`implements`不要 |
| Iteration有効 | `iter.Seq[T]` 互換 | `All() iter.Seq[{type}]` メソッド | Go 1.23+ の range-over-func 対応 |

#### Python

| 条件 | コントラクト | 追加コンポーネント | 備考 |
|------|------------|-----------------|------|
| 常時 | `Sized` | `__len__` メソッド（`size()` に委譲） | `len(list)` で要素数取得可能に |
| Iteration有効 | `Iterable` | `__iter__` + `__next__` | `for x in list` で走査可能に |
| Contains有効 | `Container` | `__contains__` メソッド | `x in list` で存在確認可能に |

#### TypeScript

| 条件 | コントラクト | 追加コンポーネント | 備考 |
|------|------------|-----------------|------|
| Iteration有効 | `Iterable<T>` | `[Symbol.iterator]()` メソッド | `for...of` で走査可能に |

### 組み合わせルール

| 条件 | 追加処理 | 説明 |
|------|---------|------|
| Contains が有効 かつ LinearSearch が無効 | LinearSearch を自動で有効化 | Contains は LinearSearch の結果を利用する |
| Sort が有効 | ElementType に比較可能性が必要 | 比較演算子または比較関数の生成が必要 |
| Search が有効 | ElementType に等価比較が必要 | 等価演算子または等価関数の生成が必要 |
| LinkedList + Singly + Reverse(Iteration) | 逆方向走査のO(n)警告 | 単方向リストでの逆走査は非効率。生成はするが警告コメントを付与 |
| ArrayBased + Insert | GrowthLogic が必要 | 挿入時に容量超過の可能性があるため |

## デフォルト構成

```yaml
name: MyList
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

**注意**: `add`, `get`, `size` は常に生成される基本操作であり、`operations` に記述しない。

このデフォルト構成は以下の制約に違反しない:
- `storage.type: array` → `initial_capacity`, `growth_strategy` が設定済み → OK
- `contains` → `linear_search` が必要 → 設定済み → OK
- `sort` は未選択 → `algorithm` 不要 → OK
- 排他制約: `array` と `linked_list` は同時選択なし → OK

## 変更シナリオ

### シナリオ1: Storage を ArrayBased → LinkedList(Singly) に変更

- **変更するDSL項目**: `storage.type`: `array` → `linked_list`, `storage.direction`: `singly`
- **影響を受けるコンポーネント**: ListStruct, NodeStruct(新規生成), Constructor, Add, Get, Remove, Insert, LinearSearch, Iterator
- **不要になるコンポーネント**: GrowthLogic, InitialCapacity
- **必要な作業**:
  1. ListStruct を配列ベースからポインタベースに差し替え
  2. NodeStruct を新規生成
  3. 全操作の実装フラグメントをLinkedList版に切り替え
  4. GrowthLogic 関連コードを除去

### シナリオ2: Sort を追加（MergeSort）

- **変更するDSL項目**: `operations` に `sort: { algorithm: merge_sort }` を追加
- **影響を受けるコンポーネント**: SortImpl(新規生成)
- **必要な作業**:
  1. ElementType の比較可能性を確認
  2. MergeSort のフラグメントを生成・追加
  3. 比較関数が未定義なら生成

### シナリオ3: ElementType を int → string に変更

- **変更するDSL項目**: `element_type`: `int` → `string`
- **影響を受けるコンポーネント**: ListStruct（型パラメータ変更）、Search（比較ロジック）、Sort（有効なら比較ロジック）
- **必要な作業**:
  1. 型パラメータの置換
  2. Search/Sort の比較ロジックが文字列比較に対応していることを確認
  3. テストデータの更新

## 代数的性質（テスト生成用）

生成されたリストコンテナが満たすべき普遍的な性質。構成に依存せず常に成立する。

### 基本性質

| 性質名 | 記述 | 説明 |
|--------|------|------|
| size-after-add | `list.Add(x); list.Size() == old_size + 1` | 追加後にサイズが1増える |
| get-after-add | `list.Add(x); list.Get(list.Size()-1) == x` | 追加した要素が末尾から取得できる |
| size-empty | `NewList().Size() == 0` | 空リストのサイズは0 |
| get-preserves-order | `list.Add(x); list.Add(y); list.Get(0)==x && list.Get(1)==y` | 追加順序が保持される |

### オプション操作の性質（対応フィーチャーが有効な場合のみ）

| 性質名 | 条件 | 記述 | 説明 |
|--------|------|------|------|
| size-after-remove | Remove有効 | `list.Remove(i); list.Size() == old_size - 1` | 削除後にサイズが1減る |
| insert-get | Insert有効 | `list.Insert(i, x); list.Get(i) == x` | 挿入した位置から取得できる |
| search-after-add | LinearSearch有効 | `list.Add(x); list.LinearSearch(x) >= 0` | 追加した要素は検索で見つかる |
| contains-after-add | Contains有効 | `list.Add(x); list.Contains(x) == true` | 追加した要素はContainsでtrue |
| sort-ordered | Sort有効 | `list.Sort(); ∀i: list.Get(i) <= list.Get(i+1)` | ソート後は昇順 |
| sort-size-preserved | Sort有効 | `list.Sort(); list.Size() == old_size` | ソートでサイズは変わらない |
| iteration-count | Iteration有効 | `count(list.Iterator()) == list.Size()` | イテレーションで全要素を走査 |

### コントラクト適合の性質（言語別・対応コントラクトが導出された場合のみ）

| 性質名 | 条件 | 記述 | 説明 |
|--------|------|------|------|
| php-countable-consistent | PHP + 常時 | `$list->count() === $list->size()` | Countable と size() の一貫性 |
| php-iterator-count | PHP + Iteration有効 | `count(iterator_to_array($list->getIterator())) === $list->size()` | IteratorAggregate の走査が全要素をカバー |
| py-len-consistent | Python + 常時 | `len(list) == list.size()` | `__len__` と size() の一貫性 |
| py-contains-consistent | Python + Contains有効 | `(x in list) == list.contains(x)` | `__contains__` と contains() の一貫性 |
| py-iter-count | Python + Iteration有効 | `len(list(iter(l))) == l.size()` | `__iter__` の走査が全要素をカバー |
| ts-iterable-count | TS + Iteration有効 | `[...list].length === list.size()` | Iterable のスプレッドが全要素をカバー |
