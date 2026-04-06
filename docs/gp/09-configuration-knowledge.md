# 構成の知識: リストコンテナ

## マッピングルール

### フィーチャー → コンポーネント

| フィーチャー（DSL項目） | 選択肢 | 生成/選択されるコンポーネント | 生成内容 |
|----------------------|--------|---------------------------|---------|
| `element_type` | 任意の型名 T | NodeDefinition, PublicAPI | 型パラメータ T をノード定義と公開APIのジェネリクス引数として埋め込む |
| `ownership` | `Copy` | NodeDefinition → `CopyMonomorphicNode<T>` or (Poly不可), MemoryManagement → `CopyMemory` | Node の data フィールドを値埋め込み（`T data`）で生成。破棄時は要素の特別な解放なし |
| `ownership` | `ExternalReference` | NodeDefinition → `ExtRef*Node<T>`, MemoryManagement → `ExtRefMemory` | Node の data フィールドをポインタ（`T* data`）で生成。破棄時はポインタのみ破棄、要素は解放しない |
| `ownership` | `OwnedReference` | NodeDefinition → `OwnedRef*Node<T>`, MemoryManagement → `OwnedRefMemory` | Node の data フィールドをポインタ（`T* data`）で生成。破棄時に要素も `delete` / `free` で解放 |
| `morphology` | `Monomorphic` | NodeDefinition（Ownership との組み合わせで決定） | 型制約: 要素は正確に型 T のみ。コンパイル時/静的に型チェック |
| `morphology` | `Polymorphic` | NodeDefinition（Ownership との組み合わせで決定） | 型制約: T のサブタイプも許容。ポインタ/参照経由で格納 |
| `length_counter.enabled` | `true` | LengthCounterMixin, ListCore（拡張）, PublicAPI（拡張） | List 構造体にカウンタフィールド追加、挿入・削除にカウンタ更新を織り込み、`length()` メソッド追加 |
| `length_counter.enabled` | `false` / 省略 | （生成なし） | カウンタ関連コードを一切生成しない。要素数取得は O(n) 走査で実装 |
| `length_counter.counter_type` | `short` / `int` / `long` | LengthCounterMixin | カウンタフィールドの型を指定された整数型で生成 |
| `tracing.enabled` | `true` | TracingMixin, ListCore（拡張） | 各操作（insert, remove, find 等）の前後にトレースログ出力コードを織り込む |
| `tracing.enabled` | `false` / 省略 | （生成なし） | トレース関連コードを一切生成しない |
| `target_language` | `cpp` / `java` / `go` 等 | 全コンポーネント | 指定言語のイディオム・構文・型システムに沿ったコードを生成 |

### NodeDefinition の組み合わせ解決表

Ownership と Morphology の組み合わせで NodeDefinition のバリアントが一意に決まる:

| Ownership \ Morphology | Monomorphic | Polymorphic |
|------------------------|-------------|-------------|
| **Copy** | `CopyMonomorphicNode<T>` — `T data; Node* next;` | **不可**（C1制約） |
| **ExternalReference** | `ExtRefMonomorphicNode<T>` — `T* data; Node* next;` | `ExtRefPolymorphicNode<T>` — `T* data; Node* next;`（サブタイプ許容） |
| **OwnedReference** | `OwnedRefMonomorphicNode<T>` — `T* data; Node* next;` | `OwnedRefPolymorphicNode<T>` — `T* data; Node* next;`（サブタイプ許容、破棄時に要素解放） |

### パラメータ → 設定値

| パラメータ（DSL項目） | 対応する生成コード上の設定 | 説明 |
|---------------------|------------------------|------|
| `element_type` | ジェネリクス型パラメータ `T` | `List<T>`, `template<typename T>` 等、言語に応じた型パラメータ構文 |
| `counter_type: short` | カウンタフィールド型: `short` / `int16_t` | 言語に応じた短整数型に変換 |
| `counter_type: int` | カウンタフィールド型: `int` / `int32_t` | 言語に応じた標準整数型に変換 |
| `counter_type: long` | カウンタフィールド型: `long` / `int64_t` / `long long` | 言語に応じた長整数型に変換 |

### 言語別型マッピング

| DSL の `counter_type` | C++ | Java | Go | Python | Rust |
|----------------------|-----|------|----|--------|------|
| `short` | `int16_t` | `short` | `int16` | `int`（型ヒント注釈のみ） | `i16` |
| `int` | `int32_t` | `int` | `int32` | `int` | `i32` |
| `long` | `int64_t` | `long` | `int64` | `int` | `i64` |

### 組み合わせルール

| 条件 | 追加処理 | 説明 |
|------|---------|------|
| `length_counter.enabled: true` + `tracing.enabled: true` | カウンタ更新操作もトレース対象に含める | トレースログに要素数の変化も記録される |
| `ownership: OwnedReference` + `tracing.enabled: true` | メモリ解放操作もトレース対象に含める | 要素の解放タイミングをトレースログで追跡可能にする |
| `morphology: Polymorphic` + `target_language: cpp` | `virtual` デストラクタの生成を保証 | C++で多態削除を安全に行うため |
| `morphology: Polymorphic` + `target_language: java` | 型パラメータの上限境界を設定（`<T extends Base>`） | Java のジェネリクスでサブタイプ関係を表現 |
| `morphology: Polymorphic` + `target_language: go` | インターフェース型制約を使用 | Go の type parameter でサブタイプ関係を表現 |

## 生成ルールの詳細

### ルール R1: List 構造体の生成

```
入力: element_type, length_counter, tracing
出力: List<T> 構造体の定義

生成フィールド:
  - head: Node<T>*          ← 常に生成
  - length: <counter_type>   ← length_counter.enabled == true の場合のみ
  - trace_log: TraceLog*     ← tracing.enabled == true の場合のみ
```

### ルール R2: Node 構造体の生成

```
入力: element_type, ownership, morphology
出力: Node<T> 構造体の定義

生成フィールド:
  - data:
      Copy         → T data          （値埋め込み）
      ExtRef       → T* data         （外部ポインタ）
      OwnedRef     → T* data         （所有ポインタ）
  - next: Node<T>*                    ← 常に生成
  
型制約:
  Monomorphic  → data の型は正確に T
  Polymorphic  → data の型は T またはそのサブタイプ
```

### ルール R3: insert 操作の生成

```
入力: ownership, length_counter, tracing
出力: insert メソッド/関数

処理の構成:
  1. [tracing == true]  → trace("insert:begin", element)
  2. node = create_node(element)
     - Copy       → node.data = copy(element)
     - ExtRef     → node.data = &element
     - OwnedRef   → node.data = new T(element)  または  node.data = &element（所有権移転）
  3. node.next = head; head = node
  4. [length_counter == true]  → length++
  5. [tracing == true]  → trace("insert:end", length?)
```

### ルール R4: remove 操作の生成

```
入力: ownership, length_counter, tracing
出力: remove メソッド/関数

処理の構成:
  1. [tracing == true]  → trace("remove:begin", position)
  2. ノード連結の解除（共通ロジック）
  3. メモリ解放:
     - Copy       → ノード解放のみ（値は自動破棄）
     - ExtRef     → ノード解放のみ（要素は外部管理）
     - OwnedRef   → delete node.data; ノード解放
  4. [length_counter == true]  → length--
  5. [tracing == true]  → trace("remove:end", length?)
```

### ルール R5: destroy（リスト破棄）操作の生成

```
入力: ownership, tracing
出力: destroy / デストラクタ

処理の構成:
  1. [tracing == true]  → trace("destroy:begin")
  2. 全ノードを走査して解放:
     - Copy       → 各ノードを free（値は自動破棄）
     - ExtRef     → 各ノードを free（要素は解放しない）
     - OwnedRef   → 各ノードについて delete node.data、その後ノードを free
  3. [tracing == true]  → trace("destroy:end")
```

### ルール R6: length 操作の生成

```
入力: length_counter
出力: length メソッド/関数

分岐:
  - length_counter.enabled == true  → return this.length （O(1)）
  - length_counter.enabled == false → ノードを走査してカウント（O(n)）
```

## デフォルト構成

```yaml
list:
  element_type: T
  ownership: Copy
  morphology: Monomorphic
  target_language: cpp
```

この構成は:
- 最もシンプルな値コピー・単相リスト
- カウンタなし、トレースなし
- 制約違反なし（Copy + Monomorphic は有効な組み合わせ）

## 変更シナリオ

### シナリオ1: 所有権モデルの変更（Copy → OwnedReference）

- **変更するDSL項目**: `ownership`: `Copy` → `OwnedReference`
- **影響を受けるコンポーネント**:
  - NodeDefinition: `CopyMonomorphicNode<T>` → `OwnedRefMonomorphicNode<T>`（data フィールドが値からポインタに変更）
  - MemoryManagement: `CopyMemory` → `OwnedRefMemory`（破棄時に要素解放ロジック追加）
  - ListCore: insert の create_node ロジック変更、remove と destroy に要素解放追加
- **必要な作業**: Node 構造体、insert、remove、destroy の再生成

### シナリオ2: Morphology の変更（Monomorphic → Polymorphic）

- **変更するDSL項目**: `morphology`: `Monomorphic` → `Polymorphic`
- **前提条件**: `ownership` が `Copy` でないことを確認（C1制約）。Copy の場合はエラー E001
- **影響を受けるコンポーネント**:
  - NodeDefinition: 型制約の変更（サブタイプ許容）
  - 言語固有の追加処理（C++: virtual デストラクタ、Java: 上限境界、Go: インターフェース制約）
- **必要な作業**: Node 構造体の再生成、言語固有の型制約コードの追加

### シナリオ3: LengthCounter の追加

- **変更するDSL項目**: `length_counter` を追加（`enabled: true`, `counter_type: int`）
- **影響を受けるコンポーネント**:
  - LengthCounterMixin: 新規生成（int型カウンタフィールド）
  - ListCore: insert と remove にカウンタ更新コード織り込み
  - PublicAPI: `length()` メソッド追加（O(n)走査 → O(1)返却に変更）
- **必要な作業**: List 構造体にフィールド追加、insert/remove の再生成、length メソッドの再生成

### シナリオ4: CounterType の変更（int → short）

- **変更するDSL項目**: `length_counter.counter_type`: `int` → `short`
- **影響を受けるコンポーネント**:
  - LengthCounterMixin: カウンタフィールドの型変更のみ
- **必要な作業**: List 構造体のカウンタフィールド型の変更（影響は局所的）

### シナリオ5: Tracing の追加（OwnedReference 構成に追加）

- **変更するDSL項目**: `tracing` を追加（`enabled: true`）、既存の `ownership: OwnedReference`
- **影響を受けるコンポーネント**:
  - TracingMixin: 新規生成
  - ListCore: 全操作にトレースフック織り込み
  - MemoryManagement への影響: メモリ解放操作もトレース対象（組み合わせルール）
- **必要な作業**: 全操作（insert, remove, find, destroy）の再生成にトレースコード追加

### シナリオ6: 言語の変更（cpp → java）

- **変更するDSL項目**: `target_language`: `cpp` → `java`
- **影響を受けるコンポーネント**: 全コンポーネント（構文・型システム・メモリ管理イディオムが変わる）
- **必要な作業**: 全コードの再生成
  - テンプレート構文 → ジェネリクス構文
  - ポインタ → 参照
  - 手動メモリ管理 → GC前提（OwnedReference の場合は参照のnull化）
  - counter_type のマッピング変更（`int32_t` → `int` 等）

## 構成解決アルゴリズム

DSL仕様書からコード生成を行う際の処理手順:

```
1. YAML パース
2. バリデーション（制約 C1〜C6 を検証）
   → 違反時はエラーコード E001〜E006 を返却
3. NodeDefinition バリアント決定
   → Ownership × Morphology の組み合わせ解決表から選択
4. MemoryManagement バリアント決定
   → Ownership から選択
5. LengthCounterMixin 決定
   → enabled == true なら CounterType に応じたバリアント選択
6. TracingMixin 決定
   → enabled == true なら StdoutTracing を選択
7. 組み合わせルール適用
   → 追加処理が必要な組み合わせを検出
8. target_language に応じた言語固有変換
   → 型マッピング、構文変換、イディオム適用
9. コード生成（ルール R1〜R6 に従い各コンポーネントを合成）
```

## テスト性質

生成コードが満たすべき性質。`/test` スキルでテストコードとして生成される。

### 共通性質（すべての構成で成立）

| # | 性質 | 検証内容 |
|---|------|---------|
| P1 | `length(new()) = 0` | 新規リストの要素数は0 |
| P2 | `isEmpty(new()) = true` | 新規リストは空 |
| P3 | `length(insert(l, e)) = length(l) + 1` | 挿入は要素数を1増やす |
| P4 | `isEmpty(insert(l, e)) = false` | 挿入後は空ではない |
| P5 | `find(insert(l, e), e) ≠ null` | 挿入した要素は検索で見つかる |
| P6 | `¬isEmpty(l) ⟹ length(remove(l)) = length(l) - 1` | 非空リストからの削除は要素数を1減らす |
| P7 | `isEmpty(l) ⟹ remove(l) returns false` | 空リストからの削除は失敗する |

### 構成固有の性質

| # | 条件 | 性質 | 検証内容 |
|---|------|------|---------|
| PC1 | Ownership = Copy | 挿入後に元の値を変更しても、リスト内の値は変化しない | コピーの独立性 |
| PE1 | Ownership = ExternalReference | 挿入後に元の値を変更すると、リスト経由でも変更が観測される | 参照の共有性 |
| PO1 | Ownership = OwnedReference | destroy後、全要素が解放される | 所有権による解放責任 |
| PL1 | LengthCounter = 有効 | 複数回のinsert/removeの後もlength()は正確 | カウンタ整合性 |
| PT1 | Tracing = 有効 | insert/remove/find/destroyの各操作がトレースログに記録される | トレース網羅性 |
