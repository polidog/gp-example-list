# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

ジェネレーティブプログラミングのプロセスを現代のAI技術で実践する実験的プロジェクト。従来のジェネレーティブプログラミングがドメインの再利用性を中心に据えていたのに対し、本プロジェクトではWebアプリ・スマホアプリにおける**アプリケーションの変化（進化）への対応**を重視する。

## ジェネレーティブプログラミングのプロセス（DEMRAL準拠）

本プロジェクトはDEMRAL（Domain Engineering Method for Reusable Algorithmic Libraries）のプロセスに準拠する。以下の6つの活動はClaude Codeのスキル（スラッシュコマンド）として実装されている。

**重要: これらは厳密な順序のフェーズではなく、任意の順番でスケジュールし、何度でも繰り返してよい活動である。**

### ドメイン分析
- **ドメインスコーピング** (`/domain-scoping`) — システムの目的・利用者・範囲を明確にし、ドメインの境界を定義する。成果物: `docs/gp/01-domain-scope.md`
- **フィーチャーモデリング** (`/feature-modeling`) — 重要概念の洗い出しと共通性・可変性分析を行い、フィーチャーツリーを構造化する。成果物: `docs/gp/02-feature-model.md`, `docs/gp/04-constraints.md`

### ドメイン設計
- **ドメイン設計** (`/domain-design`) — アーキテクチャ設計と構成の知識（フィーチャー選択→実装マッピング）を定義する。成果物: `docs/gp/05-architecture.md`, `docs/gp/06-components.md`, `docs/gp/09-configuration-knowledge.md`
- **ドメイン固有記法の定義** (`/dsl-definition`) — システムを「発注」するための宣言的な記法を定義する。成果物: `docs/gp/07-dsl-notation.md`, `docs/gp/08-dsl-examples/`

### ドメイン実装
- **実装** (`/implement`) — DSL仕様書と構成の知識に基づきコードを生成する。DSLバリデータの実装もここで行う。
- **テスト** (`/test`) — 生成コードの正しさをプロパティベーステストで検証する。

### 運用サイクル
活動が一通り完了した後は「仕様書を書いてシステムを発注 → フィーチャー変更を反映して進化」のサイクルで開発を進める。このサイクルの中でも、必要に応じて上記の活動に戻り成果物を更新する。

## 設計思想

- **変化への対応 > 再利用性**: 従来のジェネレーティブプログラミングが「製品ライン」の再利用を目指したのに対し、本プロジェクトでは単一アプリケーションが時間とともにどう進化するかを構造化する
- **共通性と可変性の対比分析**: ドメインの各要素を「共通（変わらない部分）」と「可変（変わりうる部分）」に分類し、洗い出しと判定を同時に行う。可変要素には変化のタイプ（差し替え・追加・パラメータ・構造）を付与する
- **AIによるコード生成の前提**: フィーチャーモデルとDSLを宣言的な仕様として扱い、AIが構成の知識に従ってコードを生成・変更する流れを想定する

## ドキュメント構成

すべての成果物は `docs/gp/` に番号付きで配置する。

- `docs/gp/01-domain-scope.md` — ドメインスコープ
- `docs/gp/02-feature-model.md` — フィーチャーモデル
- `docs/gp/04-constraints.md` — フィーチャー間の制約と影響関係
- `docs/gp/05-architecture.md` — アーキテクチャ設計
- `docs/gp/06-components.md` — 実装コンポーネント一覧
- `docs/gp/07-dsl-notation.md` — ドメイン固有記法の定義
- `docs/gp/08-dsl-examples/` — DSLのサンプル仕様書
- `docs/gp/09-configuration-knowledge.md` — 構成の知識（生成ルール + テスト性質）
