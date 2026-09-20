# Project ARIADNE Prototype --- Phase 0 設計原則

**Status:** COMPLETE

---

## 1. Phase 0 の目的

Prototype は、Task 管理アプリを題材として、Project ARIADNE Core を構築するために必要な技術・設計・責務境界を先行検証する。

Prototype の成功は、Task 管理アプリそのものの完成度ではなく、以下を説明できる状態になることとする。

- Windows デスクトップアプリを構築できること
- Svelte / Go / YAML / SQLite / OAS / 外部ツールの責務境界を理解できること
- ファイル正本と DB 正本を適切に使い分けられること
- Core で利用する `src / dist / runtime` の構造を実際に経験できること
- Core が生成・管理する成果物が後工程で利用可能であることを確認できること

Prototype は「Web システムを作る Prototype」ではなく、**Web システムの設計成果物を作る Windows ツールの試金石**と位置付ける。

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
| Prototype 内部 DB      | SQLite                  |
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

`src / dist / runtime` の意味を考える際は、Prototype / Core と、開発 / 利用を区別する。

| 局面           | `src / dist / runtime` の扱い                                     |
| -------------- | ----------------------------------------------------------------- |
| Prototype 開発 | 使用する。Core の試金石として構造を検証する                       |
| Prototype 利用 | Task 管理が目的。DDL / OAS は利用機能ではない。runtime は利用する |
| Core 開発      | Core 自身の設計・生成・実行に使用する                             |
| Core 利用      | 実プロジェクトの設計情報・生成物・Core 内部データの管理に使用する |

Prototype の DDL / OAS は、**Prototype の利用機能ではなく、Prototype 開発時に Core の構造を検証するためのもの**とする。

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

たとえば Prototype の Task 実データは SQLite が正本であり、`runtime/task.db` に存在する。一方、Task DB を構築するための設計情報は `src` に置き、そこから DDL を `dist` に生成する。

Core 利用時に ARIADNE が設計対象とする PostgreSQL 等の実 DB は、ARIADNE の外側に存在する。ARIADNE が管理するのは、その DB を構築するための設計情報と生成 DDL である。

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
│       └─ {service-id}.yaml
│
├─ api/
│   └─ services/
│       └─ {service-id}/
│           └─ API 定義 YAML
│
└─ templates/
    └─ Task 管理アプリ等の Prototype 用データ

        ↓ 生成・変換

dist/
├─ database/
│   └─ *.sql
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

        ↓ 実行

runtime/
└─ task.db
```

API定義とDDL定義は物理的には異なる領域で管理するが、同一のService IDによってARIADNE上の同一Serviceに関連付ける。

Serviceは物理ディレクトリ階層そのものを表す概念ではない。

項目定義、DB 設計情報、OAS の具体的な構造は、Core の要件および現行 DXSI テンプレートの設計を踏まえて後続 Phase で決定する。

---

## 9. Prototype アプリの最小機能

### 9.1 Task Template

Task Template は **YAML 正本**とする。

最小項目：

```text
id
name
defaultPriority
defaultDescription
createdAt
updatedAt
```

要件：

- Template 一覧を表示できる
- Template を新規作成できる
- Template を編集できる
- YAML へ保存できる
- アプリ再起動後に YAML から復元できる
- Template 群は 1 ファイルで管理し、全量 Read / 全量 Write とする
- `createdAt / updatedAt` はアプリが管理し、UI から編集させない
- 削除は Prototype の必須機能としない

### 9.2 Task

Task は **SQLite 正本**とする。

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
- SQLite へ保存できる
- アプリ再起動後に SQLite から復元できる
- `createdAt / updatedAt` はアプリが管理し、UI から編集させない
- `CreateUser / UpdateUser` は持たない

### 9.3 Template → Task

Task 新規作成時には Template 選択を必須とする。

Template の以下の値を Task の初期値としてコピーする。

```text
defaultPriority
defaultDescription
```

コピー後、Task と Template は独立する。Template の変更・削除が既存 Task に影響してはならない。Task 側の変更も Template に反映しない。

---

## 10. Core の試金石として必須の検証

### 10.1 YAML

アプリケーションから YAML の Read / Edit / Write を一巡する。

Prototype では Task Template を題材とする。

### 10.2 DDL / DB

DDL定義はARIADNE SourceからResolved DDL Modelを生成し、対象Database向けDDLへ変換する。

```text
ARIADNE DDL Source
      ↓
Resolved DDL Model
      ↓
Database DDL
      ├─ PostgreSQL DDL
      └─ SQLite Compatibility DDL
```

PostgreSQLをDDL設計の基準Databaseとする。

SQLiteはPostgreSQLと同列の対応Databaseとはせず、ARIADNE Prototype自身のRuntime Databaseとして必要な範囲で互換対応する。

Core利用時にARIADNEが設計対象とするPostgreSQL等の実DatabaseはARIADNEの外側に存在する。
ARIADNEが管理するのは、そのDatabaseを構築するための設計情報と生成DDLである。

### 10.3 API / OAS

API 定義の Source of Truth は `src/api/services/` 配下の ARIADNE YAML とする。

Prototype では以下の流れを検証する。

```text
ARIADNE Source
      ↓
Raw Model
      ↓
Resolved Model
      ↓
OpenAPI 3.1 Generation
      ↓
dist/api/oas/{service-id}/openapi.yaml
      ↓
OpenAPI Validation
      ↓
API Document
```

OpenAPI 3.1 自体は Source of Truth ではなく、ARIADNE Source から生成される成果物とする。

Raw Model / Resolved Model は Source of Truth ではなく、ARIADNE Source から再生成可能な中間成果物とする。

生成した OpenAPI 3.1 は以下で利用する。

- ReDoc による API Document
- Swagger UI
- Mock Server
- Backend / Frontend 等の後工程

---

## 11. OAS の出口検証

OAS のフォルダ構成・内容が後工程で利用可能であることを証明するため、Backend 側の出口検証を Prototype の必須範囲とする。

```text
dist/api/oas/{service-id}/openapi.yaml
        ↓
Backend Generate
        ↓
Go Backend
        ↓
API 公開
```

目的は Backend アプリケーションの開発ではなく、**ARIADNE が管理した OAS が実際の後工程で利用可能であることの確認**である。

- Backend Generate：対象
- Frontend Generate：対象外
- Backend Generate に利用する具体的なツールは Phase 0 では固定しない

---

## 12. Build / Git

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

---

## 13. Prototype でやらないこと

以下は Prototype の完成条件に含めない。

- Excel ライクな高機能グリッド UI
- Frontend Generate
- 認証・権限管理
- `CreateUser / UpdateUser`
- 高度な検索・ソート・フィルタ
- 履歴・コメント・タグ等の Task 管理拡張
- 過剰な UI デザイン
- 将来を見越した汎用化・抽象化
- ARIADNE 本体の Docker 化
- インストーラー
- 自動更新
- GitHub Actions

必要性が明確になったものは後続 Phase / Core の Backlog とする。

---

## 14. Prototype Definition of Done

Prototype は、以下の2つを満たした時点で Done とする。

### A. Task 管理 Windows アプリとして成立している

- `tools/build.sh` から `ariadne-prototype.exe` を生成できる
- Template を YAML から読み書きできる
- Task を SQLite から読み書きできる
- Template から Task へ初期値をコピーできる
- Template と Task は生成後に同期しない
- 再起動後も各正本から状態を復元できる

### B. Core の試金石として成立している

- `src / dist / runtime` の責務を実際の構成で確認できる
- DB 設計情報から DDL を生成し SQLite を初期化できる
- API 定義を ARIADNE YAML として分割管理できる
- ARIADNE Source から Raw Model / Resolved Model / OpenAPI 3.1 を生成できる
- Redocly で OpenAPI Validation / API Document 生成ができる
- 生成済み OpenAPI 3.1 を後工程で利用できる
- 各技術要素と責務境界を説明できる

---

## 15. Phase 0 完了判定

Phase 0 では以下を確定した。

1. Prototype で何を作るか
2. Core へ進むために何を検証するか
3. 正本・生成物・実行時データの責務
4. Prototype と Core、および開発時と利用時の違い
5. 採用する技術スタック
6. Prototype の最小機能
7. Prototype でやらないこと
8. Prototype の Definition of Done

---

## ★Phase 0：設計原則 --- COMPLETE
