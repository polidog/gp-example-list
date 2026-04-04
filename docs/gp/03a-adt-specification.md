# 抽象データ型仕様: リストコンテナ

## 型定義

### ソート（型の集合）

| ソート | 説明 | フィーチャーによる変化 |
|--------|------|---------------------|
| `List<T>` | 順序付き要素の集まりを管理するコンテナ（主要型） | LengthCounter によりカウンタフィールドの有無が変化、Tracing によりログフィールドの有無が変化 |
| `T` | リストが保持する要素の型（要素型） | ElementType で指定。Morphology=Monomorphic では正確に T、Polymorphic では T のサブタイプも許容 |
| `Node<T>` | 要素を保持しリストの連結構造を形成する内部型（補助型） | Ownership により data フィールドの保持方法が変化（値 / ポインタ） |
| `Nat` | 自然数（0以上の整数）。要素数やインデックスに使用（補助型） | — |
| `Bool` | 真偽値。操作の成否に使用（補助型） | — |
| `CounterInt` | カウンタの整数型（補助型） | CounterType により short / int / long のいずれか |

### 型パラメータ

| パラメータ | 対応フィーチャー | 制約 |
|-----------|----------------|------|
| `T` | ElementType | Monomorphic: 等値比較可能（`T` は `==` をサポート）。Polymorphic: `T` のサブタイプ `S <: T` も格納可能 |

## 操作

### new — リストの生成

- **シグネチャ**: `new() → List<T>`
- **カテゴリ**: コンストラクタ
- **共通/可変**: 共通
- **事前条件**: なし
- **事後条件**:
  - `isEmpty(new()) = true`
  - `length(new()) = 0`
- **フィーチャーによる変化**:
  - LengthCounter = 有効: カウンタが 0 に初期化される
  - Tracing = 有効: トレースログが空で初期化される

### insert — 要素の挿入

- **シグネチャ**: `insert(l: List<T>, e: T) → List<T>`
- **カテゴリ**: ミューテータ
- **共通/可変**: 共通
- **事前条件**:
  - Monomorphic: `e` の型が正確に `T` であること
  - Polymorphic: `e` の型が `T` またはそのサブタイプであること
- **事後条件**:
  - `isEmpty(insert(l, e)) = false`
  - `length(insert(l, e)) = length(l) + 1`
  - `find(insert(l, e), e) ≠ null`（挿入した要素は検索で見つかる）
- **フィーチャーによる変化**:
  - Ownership = Copy: `e` のコピーがリストに格納される。元の `e` への変更はリスト内の値に影響しない
  - Ownership = ExternalReference: `e` への参照が格納される。元の `e` への変更はリスト経由でも観測される
  - Ownership = OwnedReference: `e` への参照が格納され、所有権がリストに移転する
  - LengthCounter = 有効: 内部カウンタが 1 増加する
  - Tracing = 有効: トレースログに insert 操作が記録される

### remove — 先頭要素の削除

- **シグネチャ**: `remove(l: List<T>) → (List<T>, Bool)`
- **カテゴリ**: ミューテータ
- **共通/可変**: 共通
- **事前条件**: なし（空リストの場合は false を返す）
- **事後条件**:
  - `isEmpty(l) = true` の場合: リストは変化せず、`false` を返す
  - `isEmpty(l) = false` の場合: `length(remove(l)) = length(l) - 1`、`true` を返す
- **フィーチャーによる変化**:
  - Ownership = Copy: 削除されたノードとそのコピー値が解放される
  - Ownership = ExternalReference: ノードのみ解放される。要素は外部が管理し続ける
  - Ownership = OwnedReference: ノードとともに要素も解放される
  - LengthCounter = 有効: 成功時、内部カウンタが 1 減少する
  - Tracing = 有効: トレースログに remove 操作が記録される

### find — 要素の検索

- **シグネチャ**: `find(l: List<T>, e: T) → T | null`
- **カテゴリ**: オブザーバ
- **共通/可変**: 共通
- **事前条件**: なし
- **事後条件**:
  - リスト内に `e` と等しい要素が存在する場合: その要素を返す
  - 存在しない場合: `null` を返す
  - リストの状態は変化しない
- **フィーチャーによる変化**:
  - Tracing = 有効: トレースログに find 操作と結果（found / not found）が記録される

### length — 要素数の取得

- **シグネチャ**: `length(l: List<T>) → Nat`
- **カテゴリ**: オブザーバ
- **共通/可変**: 共通（存在はすべての構成で共通だが、計算量が変化する）
- **事前条件**: なし
- **事後条件**:
  - リスト内の要素数を返す
  - リストの状態は変化しない
- **フィーチャーによる変化**:
  - LengthCounter = 有効: O(1) で内部カウンタの値を返す。返り値の型は `CounterInt`
  - LengthCounter = 無効: O(n) で全ノードを走査してカウントする。返り値の型は `Nat`

### traverse — 要素の走査

- **シグネチャ**: `traverse(l: List<T>, f: (T) → void) → void`
- **カテゴリ**: オブザーバ
- **共通/可変**: 共通
- **事前条件**: なし
- **事後条件**:
  - リスト内の各要素に対して先頭から順に関数 `f` が適用される
  - リストの状態は変化しない
- **フィーチャーによる変化**: なし

### isEmpty — 空判定

- **シグネチャ**: `isEmpty(l: List<T>) → Bool`
- **カテゴリ**: オブザーバ
- **共通/可変**: 共通
- **事前条件**: なし
- **事後条件**:
  - `length(l) = 0` のとき `true`、それ以外のとき `false`
  - リストの状態は変化しない
- **フィーチャーによる変化**: なし

### destroy — リストの破棄

- **シグネチャ**: `destroy(l: List<T>) → void`
- **カテゴリ**: デストラクタ
- **共通/可変**: 共通（存在はすべての構成で共通だが、解放戦略が変化する）
- **事前条件**: なし
- **事後条件**:
  - リストの全ノードが解放される
  - 操作後、リストは使用不可（二重解放の防止は利用者の責任）
- **フィーチャーによる変化**:
  - Ownership = Copy: 各ノードを解放する（値は自動的に破棄される）
  - Ownership = ExternalReference: 各ノードを解放する（要素は解放しない、外部が管理し続ける）
  - Ownership = OwnedReference: 各ノードについて要素を解放し、その後ノードを解放する
  - Tracing = 有効: トレースログに destroy 操作と解放された要素の情報が記録される

## 不変条件

### 共通不変条件（すべての構成で成立）

| # | 不変条件 | 説明 |
|---|---------|------|
| I1 | `length(l) ≥ 0` | リストの要素数は常に0以上 |
| I2 | `isEmpty(l) ⟺ length(l) = 0` | 空判定と要素数の一貫性 |
| I3 | `length(insert(l, e)) = length(l) + 1` | 挿入は必ず要素数を1増やす |
| I4 | `¬isEmpty(l) ⟹ length(remove(l)) = length(l) - 1` | 非空リストからの削除は要素数を1減らす |
| I5 | `find(insert(l, e), e) ≠ null` | 挿入直後の要素は必ず検索で見つかる |
| I6 | 各ノードの `next` 参照は循環しない | リストの連結構造は非循環 |
| I7 | リスト内のすべてのノードで所有権モデルは統一されている | 要素ごとに異なる所有権モデルが混在しない |

### 条件付き不変条件（特定のフィーチャー選択時に成立）

| # | 条件 | 不変条件 | 説明 |
|---|------|---------|------|
| CI1 | Morphology = Monomorphic | `∀e ∈ l: type(e) = T` | リスト内の全要素の型は正確に T |
| CI2 | Morphology = Polymorphic | `∀e ∈ l: type(e) <: T` | リスト内の全要素の型は T またはそのサブタイプ |
| CI3 | LengthCounter = 有効 | `counter = |{n : n はリスト内のノード}|` | 内部カウンタは常にリスト内の実際のノード数と一致 |
| CI4 | Ownership = Copy | 挿入後、`e` への外部からの変更はリスト内の値に影響しない | コピーの独立性 |
| CI5 | Ownership = ExternalReference | リストはいかなる要素の生存期間も管理しない | 外部管理の保証 |
| CI6 | Ownership = OwnedReference | `destroy(l)` は全要素の `delete` を保証する | 所有権による解放責任 |
| CI7 | LengthCounter = 有効 | `0 ≤ counter ≤ MAX(CounterType)` | カウンタ値はCounterTypeの表現範囲内 |

## フィーチャーによるパラメタライズ

| フィーチャー | 型への影響 | 操作への影響 | 不変条件への影響 |
|-------------|-----------|-------------|----------------|
| ElementType | 型パラメータ `T` を決定する | 全操作のシグネチャに `T` として現れる | — |
| Ownership = Copy | Node.data が値型 `T` になる | insert: 要素をコピーして格納。remove/destroy: ノード解放のみ | CI4（コピーの独立性）が成立 |
| Ownership = ExternalReference | Node.data がポインタ型 `*T` になる | insert: 参照を格納。remove/destroy: ノードのみ解放 | CI5（外部管理の保証）が成立 |
| Ownership = OwnedReference | Node.data がポインタ型 `*T` になる | insert: 参照を格納し所有権移転。remove/destroy: 要素も解放 | CI6（所有権による解放責任）が成立 |
| Morphology = Monomorphic | `T` は正確な型一致を要求 | insert の事前条件: `type(e) = T` | CI1（型の厳密一致）が成立 |
| Morphology = Polymorphic | `T` はサブタイプも許容 | insert の事前条件: `type(e) <: T` | CI2（サブタイプ許容）が成立 |
| LengthCounter = 有効 | List に `counter: CounterInt` フィールド追加 | insert/remove にカウンタ更新が追加。length が O(1) になる | CI3（カウンタ整合性）、CI7（範囲制約）が成立 |
| LengthCounter = 無効 | counter フィールドなし | length は O(n) 走査 | CI3, CI7 は適用外 |
| CounterType | counter の具体的な整数型を決定 | length の返り値型に影響 | CI7 の MAX 値が変化（short: 32767, int: ~2.1×10⁹, long: ~9.2×10¹⁸） |
| Tracing = 有効 | List に TraceLog フィールド追加 | 全ミューテータ/デストラクタ操作にログ出力が追加。find にも結果ログ追加 | — |
| Tracing = 無効 | TraceLog フィールドなし | ログ出力なし | — |

## 操作の代数的性質

以下はフィーチャー選択に依存しない、ADT全体の代数的性質:

```
-- 挿入と長さ
length(new()) = 0
length(insert(l, e)) = length(l) + 1

-- 挿入と検索
find(insert(l, e), e) ≠ null

-- 挿入と空判定
isEmpty(new()) = true
isEmpty(insert(l, e)) = false

-- 削除と長さ（非空リスト）
¬isEmpty(l) ⟹ length(remove(l)) = length(l) - 1

-- 削除と空リスト
isEmpty(l) ⟹ remove(l) = (l, false)
```
