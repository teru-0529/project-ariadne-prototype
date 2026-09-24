# Project ARIADNE Prototype --- Phase 0 設計原則

**Status:** COMPLETE

---

## 1. Phase 0 の目的

Prototype は、現行の Excel ベースの運用を Windows デスクトップアプリとして再構築するために、
Svelte + Wails + Go を用いたアプリケーション構成の技術的成立性を先行検証する。

Prototype の主な目的は、Task 管理アプリという小さな題材を使いながら、以下を実際に確認することである。

- Svelte + Wails + Go により Windows デスクトップアプリを構築できること
- Svelte / Wails / Go の責務境界を理解できること
- Windows アプリケーションからローカルファイルを Read / Edit / Write できること
- Git で差分管理可能な可読形式を Source of Truth とする構成が技術的に成立すること
- 決められた手順から Windows 実行ファイルを Build できること

また Prototype では、Project ARIADNE Core が将来扱う設計情報について、
Element / API / DDL のモデルおよび変換構造を手作業で設計・検証する。

Validator / Resolver / Generator 等、ARIADNE Model を処理する機能そのものの実装は Core で行う。

Prototype は「Web システムを作る Prototype」ではなく、
**Web システムの設計成果物を作る Windows ツールの技術および設計の試金石**と位置付ける。

---

## 2. 設計原則

### 第1条：すべてのデータには、ひとつの正本がある

UI・DB・ファイルなど複数箇所に同じ情報が存在しても、「どれが正しい情報か」を一意に説明できること。同じ情報を複数箇所で正本としない。

### 第2条：動くことより、誰の仕事かが分かること

Svelte / Go / DB / ファイル / 外部ツールの責務を意識する。実装について「なぜこの処理をここで行うのか」を説明できることを重視する。

### 第3条：ファイルを正本とする情報は、オープンで可読な形式を選ぶ

Git で差分を確認でき、専用アプリケーションなしでも人間が理解できる形式を基本とする。Prototype / Core では YAML を基本形式とする。

### 第4条：正本の形を、UI の都合で決めない

正本は情報本来の意味と構造から設計する。UI に必要な形への加工・変換は許容するが、UI の都合を正本へ持ち込まない。

### 第5条：仕組みで守る。ただし、人間の判断を奪わない

機械的に保証すべき制約は Error として防止する。利用者の判断余地があるものは Warning / Info 等を使い分け、使い勝手とのバランスを取る。

### 第6条：分からないものを、先回りして作らない

「将来使いそう」という理由だけで抽象化・汎用化・機能追加を行わない。必要性が確定していないものは Backlog / Not now とする。

### 第7条：自分たちが理解できるものだけを積み上げる

コードが動くだけでは完成としない。技術・構造・コードについて「なぜそうなっているか」を説明できる状態で進める。

---

## 3. 技術スタック

| 領域                   | 採用技術 / 方針         |
| ---------------------- | ----------------------- |
| Frontend               | Svelte + TypeScript     |
| Desktop / Backend      | Wails + Go              |
| ファイル正本           | YAML                    |
| API 仕様               | OpenAPI 3.1             |
| OAS 検証・加工         | Redocly CLI             |
| Version Control        | Git / GitHub            |
| ローカル Development   | `tools/dev.sh`          |
| ローカル Build         | `tools/build.sh`        |
| Prototype 実行ファイル | `ariadne-prototype.exe` |

### 補足

- Wails の Windows アプリとして構築し、Frontend / Backend を別コンテナとして構成しない。
- ARIADNE 本体は Docker 化しない。
- Redocly CLI は外部ツールとして扱い、ARIADNE のドメインロジックへ組み込まない。
- Core では、最終利用者に Node.js / npm / Redocly CLI の個別インストールを要求しない実行方式を採用する。
- Redocly 等の外部ツールは ARIADNE リリース単位でバージョンを固定する。
- GitHub Actions は Prototype の必須範囲外とする。
- Core で Docker を利用する可能性は残す。主用途は、ARIADNE 本体ではなく、PostgreSQL や生成 Backend 等の成果物検証環境を想定する。
- Prototype では内部 Database を使用しない。
- SQLite を含む内部 Database の必要性は Core の設計時に改めて判断する。

---

## 4. Web システムとの境界

Project ARIADNE は Web システムそのものではない。

Web 系技術を利用するが、通常の Web アプリケーションのデプロイ構造をそのまま適用しない。

Prototype の基本構造は以下とする。

```text
Svelte
  ↓
Wails Binding
  ↓
Go
```

Svelte と Go は責務を分離するが、最終的には Windows デスクトップアプリとして構築する。

Core 利用時には、ARIADNE が管理する設計領域は実 Web 開発プロジェクト全体の一部分となる。

---

## 5. 「作る」と「使う」の区別

| 局面           | 扱い                                                                       |
| -------------- | -------------------------------------------------------------------------- |
| Prototype 開発 | Windows アプリ技術と ARIADNE Model の設計を検証する                        |
| Prototype 利用 | Task 管理アプリとして、ローカル YAML の Read / Edit / Write を行う         |
| Core 開発      | Prototype で検証した技術・設計をもとに ARIADNE 本体を実装する              |
| Core 利用      | 実プロジェクトの設計情報を管理し、Validation / Resolve / Generation を行う |

Prototype の DDL / OAS は Prototype の利用機能ではない。

Prototype では、ARIADNE Model から DDL / OAS へ至るモデルおよび変換構造を設計・検証する。
Validator / Resolver / Generator 等、その変換を実行する機能は Core で実装する。

---

## 6. `src / dist / runtime` の責務

```text
src
  = 設計情報の正本

dist
  = src から生成される成果物

runtime
  = ARIADNE 自身が実行時に利用するデータ
```

`src` は「あらゆる正本」を意味しない。

Prototype の Task データは YAML を Source of Truth とする。

Prototype において `runtime` を必須の永続化領域とはしない。
Core における実行時データや一時データの必要性、およびその保存方式は Core の設計時に改めて判断する。

Core 利用時に ARIADNE が設計対象とする PostgreSQL 等の実 Database は、ARIADNE の外側に存在する。
ARIADNE が管理するのは、その Database を構築するための設計情報と生成 DDL である。

---

## 7. Service

Project ARIADNEでは、API / DDL等の設計成果物を束ねる上位の定義単位を `Service` とする。

Serviceは特定の技術表現そのものではなく、ARIADNE上で同一の業務・機能領域に属する設計情報を関連付けるための単位である。

例えば `order-management` Serviceでは、APIとDDLで以下の異なる識別子を利用できる。

```text
Service: order-management
├─ API
│  └─ Resource: Order
└─ DDL
   └─ Schema: received_order
```

Service ID、API Resource識別子、Database Schema識別子は、それぞれ異なる責務を持ち、一致する必要はない。

---

## 8. Prototype 開発時の概念構成

Phase 0 では詳細なファイル名や YAML の分割単位までは固定しない。役割として以下の構成を採用する。

```text
src/
├─ elements/
│   ├─ types.yaml
│   └─ elements.yaml
│
├─ ddl/
│   └─ schemas/
│       └─ {service-id}/
│           └─ DDL 定義 YAML
│
└─ api/
    └─ services/
        └─ {service-id}/
            └─ API 定義 YAML

        ↓ 生成・変換

dist/
├─ ddl/
│   ├─ model/
│   │   └─ {service-id}.resolved.yaml
│   │
│   └─ postgresql/
│       └─ {service-id}.sql
│
└─ api/
    ├─ model/
    │   ├─ {service-id}.raw.yaml
    │   └─ {service-id}.resolved.yaml
    │
    └─ oas/
        └─ {service-id}/
            ├─ openapi.yaml
            ├─ redoc.html
            └─ docs/

prototype-data/
 └─ tasks.yaml
```

Phase 0 では、Task YAML の具体的なファイル名や配置方法は固定しない。

`prototype-data` は Prototype アプリケーションの技術検証に使用する専用データ領域であり、
ARIADNE Model の設計情報を表す `src` および生成成果物を表す `dist` とは責務を分離する。

ARIADNE Model であるAPI定義とDDL定義は物理的には異なる領域で管理するが、同一のService IDによってARIADNE上の同一Serviceに関連付ける。

Serviceは物理ディレクトリ階層そのものを表す概念ではない。

---

## 9. Prototype アプリの最小機能

### 9.1 Task

Task は **YAML 正本**とする。

最小項目：

```text
id
title
description
priority
status
createdAt
updatedAt
```

要件：

- Task 一覧を表示できる
- Task を新規作成できる
- Task を編集できる
- YAML へ保存できる
- アプリ再起動後に YAML から復元できる
- Task 群は 1 ファイルで管理し、全量 Read / 全量 Write とする
- `createdAt / updatedAt` はアプリが管理し、UI から編集させない
- `CreateUser / UpdateUser` は持たない
- 削除は Prototype の必須機能としない

---

## 10. Prototype で必須の検証

Prototype では、以下の2つを異なる目的の検証として扱う。

1. Windows Application の技術検証
2. ARIADNE Model の設計検証

### 10.1 Windows Application の技術検証

Svelte + Wails + Go により、Windows デスクトップアプリケーションが成立することを確認する。

基本構造は以下とする。

```text
Svelte
  ↓
Wails Binding
  ↓
Go
```

Task 管理アプリを題材として、画面表示・操作から Go の処理までを一巡できることを確認する。

また、一覧から対象を選択し、編集領域へ遷移して操作できる Application Shell が成立すること。

#### File I/O

Task YAML を利用し、Windows アプリケーションからローカルファイルの Read / Edit / Write を一巡する。

```text
Task YAML
    ↓ Read
   Go
    ↓
Wails Binding
    ↓
 Svelte
    ↓ Edit
Wails Binding
    ↓
   Go
    ↓ Write
Task YAML
```

アプリケーション終了後に再起動し、YAML から保存済み状態を復元できることを確認する。

Git 操作そのものは Prototype では実装しない。
保存された YAML が人間に可読であり、Git による差分管理が可能な形式として維持されることを確認できればよい。

### 10.2 ARIADNE Model の設計検証

Prototype では、Core が将来扱う ARIADNE Model について、以下の設計を手作業で検証する。

- Element Definition
- API Definition
- DDL Definition
- Source of Truth の責務
- Validation の責務
- Raw / Resolved Model の役割
- OAS / PostgreSQL DDL への変換構造

Prototype では、正本 YAML、中間モデル、生成結果のサンプルを手作業で作成し、モデルおよび変換規則が成立することを確認する。

Validator / Resolver / Generator は実装しない。これらの実装は Project ARIADNE Core で行う。

---

## 11. Build / Git

Prototype は Git / GitHub で管理する。

ローカルでは以下の手順で実行ファイルを再生成できることを必須とする。

```text
tools/build.sh
   ↓
Wails build 等
   ↓
ariadne-prototype.exe
```

「特定の開発環境で偶然動く」状態ではなく、決められた手順で Build 可能な状態を目指す。

GitHub Actions は Prototype の必須範囲外とする。

Prototype アプリケーションから Git の操作は行わない。
Git status / diff / commit / push 等を ARIADNE から操作する機能については、Core で必要性を判断する。

---

## 12. Prototype でやらないこと

以下は Prototype の完成条件に含めない。

- Excel ライクな高機能グリッド UI
- 認証・権限管理
- 高度な検索・ソート・フィルタ
- 履歴・コメント・タグ等の Task 管理拡張
- 過剰な UI デザイン
- 将来を見越した汎用化・抽象化
- ARIADNE 本体の Docker 化
- インストーラー
- 自動更新
- GitHub Actions
- SQLite 等を利用した Prototype 内部 Database
- SQLite Compatibility DDL
- Prototype アプリケーションからの Git 操作
- ARIADNE Model の Validator / Resolver / Generator の実装
- ARIADNE Model から生成した OAS / DDL を利用した Application Integration
- Backend / Frontend Generate
- Core の最終的なファイル分割方式・Working Model・内部永続化方式の確定

必要性が明確になったものは後続 Phase / Core の Backlog とする。

---

## 13. Prototype Definition of Done

Prototype は、以下の2つを満たした時点で Done とする。

### A. Windows Application の技術検証が成立している

- `tools/build.sh` から `ariadne-prototype.exe` を生成できる
- Svelte + Wails + Go による Windows デスクトップアプリケーションとして動作する
- Task を YAML から読み込める
- Task を画面から新規作成・編集できる
- Task を YAML へ保存できる
- アプリ再起動後に YAML から状態を復元できる
- YAML が人間に可読で、Git による差分管理が可能な形式として維持される
- 一覧から対象を選択し、編集領域へ遷移して操作できる Application Shell が成立する

### B. ARIADNE Model の設計検証が成立している

- Element Definition の正本モデルと責務を説明できる
- API Definition の正本モデルと責務を説明できる
- DDL Definition の正本モデルと責務を説明できる
- Source / Raw / Resolved / Generated Artifact の責務境界を説明できる
- OAS へ至る変換構造および規則を説明できる
- PostgreSQL DDL へ至る変換構造および規則を説明できる
- Prototype で手作業とした領域と、Core で実装する領域を説明できる

---

## 14. Phase 0 完了判定

Phase 0 では以下を確定した。

1. Prototype の目的
2. Windows Application として検証する技術要素
3. ARIADNE Model として設計検証する領域
4. Source of Truth と成果物の責務
5. Prototype と Core の責務境界
6. 採用する技術スタック
7. Prototype アプリの最小機能
8. Prototype でやらないこと
9. Prototype の Definition of Done

---

## ★Phase 0：設計原則 --- COMPLETE
