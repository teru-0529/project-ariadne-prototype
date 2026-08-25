# Prototype Phase 1: Repository / Workspace

**Status:** COMPLETE

---

## 1. Purpose

Project ARIADNE Prototype の開発基盤を構築する。

本Phaseでは、Phase 0で決定した設計原則およびディレクトリ責務を前提として、
Repository / Workspace / Wails / Svelte / Go の開発基盤を確立する。

Phase 0で決定した成果物・正本・実行時データ等の責務については、
Phase 1では変更しない。

---

## 2. Repository

### Repository

- Repository: `project-ariadne-prototype`
- License: MIT
- Main branch: `main`

Prototypeは1 Repositoryで管理する。

### Development Flow

基本的な開発フローは以下とする。

1. `main` から feature branch を作成
2. feature branch 上で開発・commit
3. GitHubへpush
4. Pull Requestを作成
5. Review
6. `main` へmerge
7. remote / local のfeature branchを削除

Phase 1では以下のbranchを使用した。

```text
feature/phase1-repository-workspace
```

Pull Requestによるレビュー・修正を経て `main` へmergeした。

---

## 3. Workspace

ローカルWorkspaceは以下とする。

```text
C:\workspace\project-ariadne-prototype
```

VS CodeではRepository全体をWorkspaceとして扱う。

Workspace定義として `.code-workspace` ファイルを使用する。

---

## 4. Directory Structure

Project ARIADNE Prototypeのディレクトリ構成は、
Phase 0で決定した責務分離を維持する。

概念上の構成は以下とする。

```text
project-ariadne-prototype/
├─ app/
│  ├─ backend/
│  │  └─ services/
│  ├─ frontend/
│  │  ├─ bindings/
│  │  └─ src/
│  ├─ build/
│  ├─ bin/
│  ├─ go.mod
│  ├─ go.sum
│  ├─ main.go
│  └─ Taskfile.yml
│
├─ src/
│  ├─ definitions/
│  ├─ api/
│  └─ database/
│
├─ templates/
├─ dist/
├─ runtime/
├─ docs/
└─ tools/
   ├─ dev.sh
   └─ build.sh
```

### Responsibilities

| Directory | Responsibility |
| --- | --- |
| `app/` | Project ARIADNE PrototypeのWindowsアプリケーション |
| `app/backend/` | Go backend |
| `app/backend/services/` | Wailsから公開するBackend Service |
| `app/frontend/` | Svelte frontend |
| `app/frontend/bindings/` | Wails generated bindings |
| `app/build/` | Wails build configuration |
| `app/bin/` | Application build output |
| `src/definitions/` | 項目定義等、ARIADNE設計情報の正本 |
| `src/api/` | OAS関連の正本・生成処理で使用する資材 |
| `src/database/` | DDL関連の正本・生成処理で使用する資材 |
| `templates/` | Task管理アプリ等、Prototypeで使用するテンプレート／アプリデータ側の資材 |
| `dist/` | Project ARIADNEが生成した成果物 |
| `runtime/` | SQLite等の実行時データ |
| `docs/` | 設計原則・Phaseごとの決定事項等 |
| `tools/` | Repository rootから利用する開発・build用utility script |

### Source of Truth

Phase 0で決定した以下の原則を維持する。

- ARIADNE設計情報の正本は `src/` 配下で管理する
- 項目定義は `src/definitions/` に配置する
- OAS関連は `src/api/` に配置する
- DDL関連は `src/database/` に配置する
- Task管理アプリ等のテンプレート／アプリデータ側の資材は `templates/` として分離する
- 成果物と実行時データを分離する
- `dist/` は生成成果物
- `runtime/` はSQLite等の実行時データ

Phase 1ではこれらの責務を変更しない。

---

## 5. Git Management Policy

`app/bin/` はApplication build outputであるためGit管理対象外とする。

同様に以下はGit管理対象外とする。

```text
app/frontend/node_modules/
app/frontend/dist/
app/bin/
```

一方、Project ARIADNE自身が生成するRepository rootの `dist/` は、
Prototypeでは生成結果を確認可能とするためGit管理対象とする。

したがって、以下の2つの `dist` は責務が異なる。

```text
/dist/                  Project ARIADNEの生成成果物
/app/frontend/dist/     Svelte/Viteのbuild成果物
```

前者はPrototypeではGit管理し、
後者はGit管理しない。

---

## 6. Technology Stack

Prototypeでは以下を使用する。

| Layer | Technology |
| --- | --- |
| Desktop Application | Wails v3 |
| Frontend | Svelte 5 |
| Frontend Language | TypeScript |
| Backend | Go |
| Node.js Management | Volta |

Phase 1完了時点のtoolchainは以下。

```text
Go      1.27.0
Node.js 24.19.0
npm     11.17.0
```

Node.js / npm は `package.json` の `volta` によりProject単位で固定する。

Goは `go.mod` に以下を設定する。

```text
go 1.27.0
```

---

## 7. Wails Application Structure

Wails v3の初期生成コードから、
Prototypeに不要なsample UI / event / serviceを削除する。

一方、Wailsが標準で提供するplatform別build構成等は、
Prototypeで直接使用しないものについても、
Framework標準構成として問題がない限り残置する。

Backend Serviceは以下に配置する。

```text
app/backend/services/
```

Phase 1では `AppService` を作成し、Wails bindingsを介した以下の経路を確認した。

```text
Go Backend Service
        ↓
Wails bindings
        ↓
Svelte / TypeScript
```

Wails bindingsはGo側のService定義から生成されるgenerated filesとして扱い、
原則として手動編集しない。

FrontendからBackendの非同期処理を呼び出す場合は、
可読性を重視し `async / await` を基本形とする。

---

## 8. Application Naming

PrototypeのApplication表示名は以下とする。

```text
Project ARIADNE Prototype
```

以下について同名称を使用する。

- Wails Application Name
- Window Title
- Product Name

Windows実行ファイル名は以下とする。

```text
ariadne-prototype.exe
```

表示名と物理ファイル名は分離する。

---

## 9. Development

Repository rootから開発環境を起動できるようにする。

```bash
./tools/dev.sh
```

`tools/dev.sh` は自身の配置場所を基準として `app/` へ移動し、

```bash
wails3 dev
```

を実行する。

これにより、開発者が毎回 `app/` へ手動で移動する必要をなくす。

---

## 10. Build

Repository rootからApplicationをbuildできるようにする。

```bash
./tools/build.sh
```

`tools/build.sh` は自身の配置場所を基準として `app/` へ移動し、

```bash
wails3 build
```

を実行する。

Windows実行ファイルは以下に生成する。

```text
app/bin/ariadne-prototype.exe
```

Phase 1では以下を確認した。

- build成功
- `ariadne-prototype.exe` 生成
- exe単体起動成功
- Application / Window Title正常
- Go Backend Service呼び出し正常
- Problems 0
- Console errorなし
- Terminal errorなし

一時的に `ariadne-prototype.exe~` が存在したが、
削除後にbuildおよびexeの起動・終了を行っても再現しなかった。

通常のbuild / executionによって生成されるものではないと判断し、
Phase 1では対応不要とする。

---

## 11. Generated Files

Wailsによって生成される資材は、
Frameworkが管理するgenerated filesとして扱う。

代表例は以下。

- Wails bindings
- platform build files
- application assets
- OS-specific build configuration

generated filesについては、
全ファイルをProject ARIADNE固有の設計対象として個別レビューしない。

Pull Request reviewでは、以下を重点的に確認する。

1. Project ARIADNE固有の実装
2. Application entry point
3. Frontend / Backend間の接続
4. Development / Build configuration
5. Dependency / Toolchain configuration
6. Git管理対象・除外対象
7. Wails sample / dummy設定等の不要な残存

---

## 12. Prototype Scope

Phase 0で決定したPrototypeのスコープを維持する。

Prototypeでは以下を対象外とする。

- Dockerによる実行環境構築
- GitHub ActionsによるCI/CD

Wails標準構成としてDockerや他platform向けの設定が存在していても、
Prototypeで利用することを意味しない。

PrototypeのApplication実行対象はWindowsとする。

---

## 13. Phase 1 Verification

Phase 1完了時点で以下を確認した。

- Repository / Workspace構築 OK
- Wails v3 application起動 OK
- Svelte 5 frontend起動 OK
- Go backend起動 OK
- Go → Wails bindings → Svelte 疎通 OK
- `AppService` 呼び出し OK
- dev script OK
- build script OK
- Windows executable単体起動 OK
- Node.js / npm version固定 OK
- Go version設定 OK
- Problems 0
- Console errorなし
- Terminal errorなし
- Pull Request review実施
- Review指摘4ファイル修正
- `main` merge完了
- remote / local feature branch削除完了

---

## 14. Phase 1 Result

Phase 0で決定したProject ARIADNEの設計原則および責務分離を維持したまま、
Prototypeの機能開発を開始するためのRepository / Workspace / Application基盤を構築した。

Phase 1以降はこの基盤を使用し、
Project ARIADNE固有の機能開発へ進む。

---

## Prototype Phase 1: Repository / Workspace — COMPLETE
