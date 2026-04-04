# フィーチャー間の制約と影響関係: リストコンテナ

## 制約一覧

### requires（依存）

| フィーチャーA | 関係 | フィーチャーB | 理由 |
|-------------|------|-------------|------|
| Polymorphic | requires | ExternalReference or OwnedReference | 多相リストでは要素をポインタ経由で扱う必要がある（値埋め込みではサブタイプの実体を格納できない） |
| CounterType | requires | LengthCounter | CounterType は LengthCounter が有効な場合にのみ意味を持つ。LengthCounter 有効時は CounterType の指定が必須 |

### excludes（排他）

| フィーチャーA | 関係 | フィーチャーB | 理由 |
|-------------|------|-------------|------|
| Polymorphic | excludes | Copy | コピーセマンティクスでは要素を値として格納するため、サブタイプのスライシングが発生し多態性を実現できない |

### impacts（影響）

| フィーチャーA | 関係 | フィーチャーB | 影響内容 |
|-------------|------|-------------|---------|
| Ownership | impacts | リストの破棄処理 | 所有権モデルにより、リスト破棄時の要素解放ロジックが変わる |
| Morphology | impacts | 要素の挿入処理 | Polymorphic時は型チェックがサブタイプを考慮する必要がある |
| LengthCounter | impacts | 挿入・削除処理 | 有効時、挿入・削除のたびにカウンタの更新処理が追加される |
| CounterType | impacts | LengthCounter | カウンタ型により、リストの最大要素数とメモリ使用量が変わる |
| Tracing | impacts | 全操作 | 有効時、各操作にログ出力のフックが追加される |

### recommends（推奨）

| フィーチャーA | 関係 | フィーチャーB | 理由 |
|-------------|------|-------------|------|
| Tracing | recommends | LengthCounter | トレーシングでリスト状態を記録する際、要素数が即座に参照できると有用 |

## 変更波及パス（change-chain）

### パス1: Ownership変更時
```
Ownership変更
  → Node のデータ保持方法の変更
  → リスト破棄処理の変更
  → 挿入・削除時のメモリ管理ロジックの変更
```

### パス2: Morphology変更時
```
Morphology変更（Monomorphic ↔ Polymorphic）
  → Ownership の選択肢が制約される（Polymorphic時、Copyは不可）
  → Node の型制約の変更
  → 挿入時の型チェックロジックの変更
```

### パス3: LengthCounter の有効/無効切り替え時
```
LengthCounter切り替え
  → List の属性変更（length フィールドの追加/削除）
  → 挿入・削除処理へのカウンタ更新ロジックの追加/削除
  → 要素数取得の計算量変更（O(1) ↔ O(n)）
```

### パス4: Tracing の有効/無効切り替え時
```
Tracing切り替え
  → TraceLog エンティティの追加/削除
  → 全操作へのログ出力フックの追加/削除
```

## 有効な構成の検証ルール

1. **Ownership は XOR**: ExternalReference / OwnedReference / Copy のうち正確に1つを選択
2. **Morphology は XOR**: Monomorphic / Polymorphic のうち正確に1つを選択
3. **Polymorphic + Copy は不可**: Morphology = Polymorphic かつ Ownership = Copy の組み合わせは禁止
4. **ElementType は必須**: 型パラメータ T は必ず指定する
5. **LengthCounter と Tracing は独立**: それぞれ独立に有効/無効を設定可能
6. **CounterType は XOR**: LengthCounter 有効時、short / int / long のうち正確に1つを選択
7. **CounterType は LengthCounter に従属**: LengthCounter 無効時、CounterType は指定不可

## 有効な構成の組み合わせ

Ownership × Morphology の有効な組み合わせ:

| | Monomorphic | Polymorphic |
|---|---|---|
| **ExternalReference** | OK | OK |
| **OwnedReference** | OK | OK |
| **Copy** | OK | **不可** |

LengthCounter の選択肢: 無効（1通り）+ 有効時の CounterType（short / int / long = 3通り）= 4通り

上記の有効な組み合わせ（5通り）に対して、LengthCounter（4通り）× Tracing（2通り）= 8通りを掛け合わせ、合計 **40通り** の有効な構成が存在する。
