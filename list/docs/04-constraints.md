# フィーチャー間の制約と影響関係: リストコンテナ

## requires（依存）

| フィーチャーA | → | フィーチャーB | 理由 |
|-------------|---|-------------|------|
| InitialCapacity | requires | ArrayBased | 初期容量は配列ベース実装にのみ適用される |
| GrowthStrategy | requires | ArrayBased | 拡張戦略は配列ベース実装にのみ適用される |
| Direction（LinkedList） | requires | LinkedList | 方向の選択は連結リスト実装にのみ適用される |
| Doubly | requires | LinkedList | 双方向は連結リストの属性 |
| SortAlgorithm | requires | Sort | ソートアルゴリズムはSort機能が有効な場合のみ |
| Reverse（Iteration） | requires | Iteration | 逆方向走査はイテレーション機能が有効な場合のみ |

## excludes（排他）

| フィーチャーA | ⊗ | フィーチャーB | 理由 |
|-------------|---|-------------|------|
| ArrayBased | excludes | LinkedList | 内部データ構造は1つだけ選択（XOR） |
| Singly | excludes | Doubly | 連結リストの方向は1つだけ選択（XOR） |
| Doubling | excludes | Additive | 拡張戦略は1つだけ選択（XOR） |
| InsertionSort | excludes | MergeSort | ソートアルゴリズムは1つだけ選択（XOR） |
| Forward（Iteration） | excludes | Reverse（Iteration） | イテレーション方向は1つだけ選択（XOR） |

## impacts（影響）

| フィーチャーA | → | フィーチャーB | 影響の内容 |
|-------------|---|-------------|-----------|
| Storage の変更 | impacts | すべての操作 | データ構造が変わると全操作の実装が変わる |
| LinkedList + Singly | impacts | Reverse（Iteration） | 単方向連結リストでの逆方向走査はO(n)のコスト |
| ElementType の変更 | impacts | Search, Sort | 型が変わると比較ロジックが変わる |

## recommends（推奨）

| フィーチャーA | → | フィーチャーB | 理由 |
|-------------|---|-------------|------|
| Sort | recommends | Search | ソートされたリストでは二分探索などの効率的な検索が可能になる |
| LinkedList + Doubly | recommends | Reverse（Iteration） | 双方向連結リストでは逆方向走査がO(1)で可能 |

## parameter-constraint（パラメータ制約）

| パラメータ | 制約 | 理由 |
|-----------|------|------|
| InitialCapacity | > 0 | 容量は正の整数でなければならない |
| InitialCapacity | ≤ 65536 | 学習目的のため過大な初期割当てを制限 |

## change-chain（変更波及チェーン）

### Storage の差し替え
```
Storage変更 → 全操作の実装変更 → テストの更新
```
データ構造の変更は最も影響範囲が大きい。ArrayBased ↔ LinkedList の切り替えにより、すべての操作の内部実装とテストが影響を受ける。

### ElementType の変更
```
ElementType変更 → Search/Sort の比較ロジック変更 → テストのテストデータ変更
```
要素型の変更は、比較を伴う操作（Search, Sort）に波及する。

### Sort の追加
```
Sort追加 → SortAlgorithm選択 → （推奨）Search の効率化検討
```
ソート機能の追加時に、検索機能との連携を検討する機会となる。
