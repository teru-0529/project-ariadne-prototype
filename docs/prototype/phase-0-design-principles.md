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
|------------------------|-------------------------|
| Frontend               | Svelte + TypeScript     |
| Desktop / Backend      | Wails + Go              |
| Prototype 内部 DB      | SQLite                  |
| ファイル正本           | YAML                    |
| API 仕様               | OpenAPI Specification   |
| OAS 検証・加工         | Redocly CLI             |
| Version Control        | Git / GitHub            |
| ローカル Build         | `build.bat`             |
| Prototype 実行ファイル | `ARIADNE-Prototype.exe` |

### 補足

- Wails の Windows アプリとして構築し、Frontend / Backend を別コンテナとして構成しない。
- ARIADNE 本体は Docker 化しない。
- Redocly は ARIADNE 本体へ組み込まず、外部 CLI として扱う。
- Redocly 等の外部ツールはバージョンを固定する。
- GitHub Actions は Prototype の必須範囲外とする。
- Core で Docker を利用する可能性は残す。主用途は、ARIADNE 本体ではなく、PostgreSQL や生成 Backend 等の成果物検証環境を想定する。

---

## 4. Web システムとの境界

Project ARIADNE は Web システムそのものではない。

Web 系技術を利用するが、通常の Web アプリケーションのデプロイ構造をそのまま適用しない。

Prototype の基本構造は以下とする。

``` text
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
|----------------|-------------------------------------------------------------------|
| Prototype 開発 | 使用する。Core の試金石として構造を検証する                       |
| Prototype 利用 | Task 管理が目的。DDL / OAS は利用機能ではない。runtime は利用する |
| Core 開発      | Core 自身の設計・生成・実行に使用する                             |
| Core 利用      | 実プロジェクトの設計情報・生成物・Core 内部データの管理に使用する |

Prototype の DDL / OAS は、**Prototype の利用機能ではなく、Prototype 開発時に Core の構造を検証するためのもの**とする。

---

## 6. `src / dist / runtime` の責務

``` text
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

## 7. Prototype 開発時の概念構成

Phase 0 では詳細なファイル名や YAML の分割単位までは固定しない。役割として以下の構成を採用する。

``` text
src/
├─ definitions/
│   ├─ 項目定義 YAML
│   └─ Task Template YAML
│
├─ database/
│   └─ DB 設計情報 YAML
│       ※ファイル分割単位は Phase 0 では決めない
│
└─ api/
    ├─ root.yaml
    ├─ paths/
    ├─ resources/
    └─ ...

        ↓ 生成・変換

dist/
├─ database/
│   └─ *.sql
│
└─ api/
    ├─ openapi.yaml
    └─ openapi.html

        ↓ 実行

runtime/
└─ task.db
```

項目定義、DB 設計情報、OAS の具体的な構造は、Core の要件および現行 DXSI テンプレートの設計を踏まえて後続 Phase で決定する。

---

## 8. Prototype アプリの最小機能

### 8.1 Task Template

Task Template は **YAML 正本**とする。

最小項目：

``` text
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

### 8.2 Task

Task は **SQLite 正本**とする。

最小項目：

``` text
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

### 8.3 Template → Task

Task 新規作成時には Template 選択を必須とする。

Template の以下の値を Task の初期値としてコピーする。

``` text
defaultPriority
defaultDescription
```

コピー後、Task と Template は独立する。Template の変更・削除が既存 Task に影響してはならない。Task 側の変更も Template に反映しない。

---

## 9. Core の試金石として必須の検証

### 9.1 YAML

アプリケーションから YAML の Read / Edit / Write を一巡する。

Prototype では Task Template を題材とする。

### 9.2 DDL / DB

以下の流れを Prototype 開発時に一度通す。

``` text
src/database/
      ↓
DDL 生成
      ↓
dist/database/*.sql
      ↓
DB 初期化
      ↓
runtime/task.db
```

Prototype では SQLite を利用する。

Core 利用時の設計対象 DB は PostgreSQL 等になり得るが、その実 DB は ARIADNE の外部に存在する。

### 9.3 OAS

OAS の正本を `src/api` で管理し、Redocly を利用して検証・生成する。

``` text
src/api/
   ↓
Redocly lint
   ↓
Redocly bundle
   ↓
dist/api/openapi.yaml
   ↓
HTML 生成
   ↓
dist/api/openapi.html
```

- `openapi.yaml` は後工程・Generate 用
- `openapi.html` は人間によるレビュー・確認用

Prototype アプリ自体から OAS を読み書きする機能は優先度を下げる。YAML Read / Write の技術検証は Task Template で実施できるためである。

---

## 10. OAS の出口検証

OAS のフォルダ構成・内容が後工程で利用可能であることを証明するため、Backend 側の出口検証を Prototype の必須範囲とする。

``` text
dist/api/openapi.yaml
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

## 11. Build / Git

Prototype は Git / GitHub で管理する。

ローカルでは以下の手順で実行ファイルを再生成できることを必須とする。

``` text
build.bat
   ↓
Wails build 等
   ↓
ARIADNE-Prototype.exe
```

「特定の開発環境で偶然動く」状態ではなく、決められた手順で Build 可能な状態を目指す。

GitHub Actions は Prototype の必須範囲外とする。

---

## 12. Prototype でやらないこと

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

## 13. Prototype Definition of Done

Prototype は、以下の2つを満たした時点で Done とする。

### A. Task 管理 Windows アプリとして成立している

- `build.bat` から `ARIADNE-Prototype.exe` を生成できる
- Template を YAML から読み書きできる
- Task を SQLite から読み書きできる
- Template から Task へ初期値をコピーできる
- Template と Task は生成後に同期しない
- 再起動後も各正本から状態を復元できる

### B. Core の試金石として成立している

- `src / dist / runtime` の責務を実際の構成で確認できる
- DB 設計情報から DDL を生成し SQLite を初期化できる
- OAS を分割管理できる
- Redocly で lint / bundle / HTML 生成できる
- bundle 済み OAS から Backend Generate し、API を公開できる
- 各技術要素と責務境界を説明できる

---

## 14. Phase 0 完了判定

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
