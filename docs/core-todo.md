# Core TODO

## Source / Model

### Element Definition

- [ ] Phase 2 Element Definition Loaderを実装する

### API Definition

- [ ] Phase 3 API Definition Loaderを実装する
- [ ] API Source → Raw Model変換を実装する
- [ ] API Raw Model → Resolved Model変換を実装する

### DDL Definition

- [ ] Phase 4 DDL Definition Loaderを実装する
- [ ] DDL Source → Resolved DDL Model変換を実装する

### 共通

- [ ] Raw / Resolved Modelのデバッグ出力方式を決める

## Validation

### Phase 2

- [ ] Phase 2 Validation Rule `V-001` ～ `V-028` を実装する

### Phase 3

- [ ] File Validationを実装する
- [ ] API Validationを実装する
- [ ] Service Validationを実装する
- [ ] OpenAPI Validationを実装する
- [ ] Phase 3 Validation Rule `V-029` ～ `V-084` を実装する

### Phase 4

- [ ] DDL Validationを実装する
- [ ] Phase 4 Validation Rule `V-085` 以降を実装する

### 共通

- [ ] ValidatorをFile I/Oから分離した共通Packageとして実装する
- [ ] UI / CLI / Testから同一Validatorを利用できる構造にする

## OpenAPI 3.1 Generation

- [ ] Resolved Service → OpenAPI `info` 変換を実装する
- [ ] Element → OpenAPI Schema変換を実装する
- [ ] Resolved Schema → OpenAPI Schema変換を実装する
- [ ] Parameter / Usage → OpenAPI Parameter Component変換を実装する
- [ ] Path / Operation変換を実装する
- [ ] Request / Response変換を実装する
- [ ] Example生成を実装する
- [ ] Pagination built-in生成を実装する
- [ ] Location Header生成を実装する
- [ ] Standard Error生成を実装する
- [ ] built-in `/health` を生成する
- [ ] built-in `/version` を生成する
- [ ] `externalDocs` を変換する
- [ ] `servers` の生成方式を実装する

## PostgreSQL DDL Generation

- [ ] Resolved DDL Model + Types / ElementsからPostgreSQL DDLを生成する
- [ ] Schema / Table / Columnを生成する
- [ ] Primary Key / Unique Constraintを生成する
- [ ] Foreign Keyを生成する
- [ ] Indexを生成する
- [ ] Defaultを生成する
- [ ] Database / BuiltIn / Custom Generationを生成する
- [ ] BuiltIn Column / Audit Columnを生成する
- [ ] BuiltIn Function / Triggerを生成する
- [ ] Commentを生成する
- [ ] Custom SQLを出力へ組み込む
- [ ] DDLの依存関係と実行順序を管理する
- [ ] DDL成果物の物理的なファイル構成を決定する

## DDL Document

- [ ] DDL Definitionから人間向けDDL Documentを生成する
- [ ] Table / Column / Constraint / Indexを参照できるTable Definitionを生成する
- [ ] Element由来の型・桁・説明等をDDL Documentへ反映する
- [ ] Markdown / HTML等、DDL Documentの成果物形式を決定する
- [ ] Wails ApplicationからDDL Documentを参照する方式を決める

## Application Version

- [ ] Git TagからApplication Versionを取得する
- [ ] Semantic Versionとして利用可能なTagのみ採用する
- [ ] 有効なVersionを取得できない場合の開発時Versionを決定する
- [ ] OpenAPI `info.version` にApplication Versionを使用する
- [ ] `/version` Responseに同じApplication Versionを使用する

## API Document

- [ ] Redocly CLIの実行環境をCoreへ組み込む方式を決める
- [ ] 最終利用者にNode.js / npm / Redocly CLIの個別インストールを要求しない構成にする
- [ ] Redocly CLIのバージョンをARIADNEリリース単位で固定する
- [ ] 通常利用時はRedocly Update Noticeを抑止する
- [ ] Telemetryを無効化する
- [ ] OpenAPIからReDoc静的HTMLを生成する
- [ ] Wails ApplicationからAPI Documentを参照する方式を決める

## externalDocs

- [ ] CoreにおけるMarkdown → HTML変換方式を決定する
- [ ] Markdown内の最初のH1をHTML titleとして使用する
- [ ] H1が存在しない場合をErrorとする
- [ ] OAS成果物ディレクトリへ補足HTMLを生成する
- [ ] OpenAPI `externalDocs.url` を生成HTMLへの相対URLへ変換する

## Windows Application

- [ ] Element Definitionの編集UIを設計・実装する
- [ ] API Definitionの編集UIを設計・実装する
- [ ] DDL Definitionの編集UIを設計・実装する
- [ ] ARIADNE SourceのRead / Edit / Writeを実装する
- [ ] Validation結果をApplication上で表示する
- [ ] GenerationをApplicationから実行できるようにする
- [ ] Coreの内部実行時データおよび永続化方式の必要性を判断する
- [ ] SQLite等の内部Databaseを採用するか判断する

## Git Integration

- [ ] ARIADNE ApplicationからGit操作を行う必要性を判断する
- [ ] 必要な場合、status / diff / commit / push等の責務と実装方式を決定する

## Future Feature

### AI-assisted Editing

- [ ] ARIADNE Sourceの登録・編集をAIで支援するAI Modeを検討する
- [ ] 自然言語からElement / API / DDL Definitionの登録候補を生成する方式を検討する
- [ ] AIの生成結果を直接正本へ反映せず、利用者が確認・修正してから確定する操作モデルを設計する
- [ ] AI Provider / Model / API接続情報をARIADNEでどのように管理するか検討する
- [ ] AI利用に必要な契約・認証・利用制約を整理する
- [ ] GitHub Copilot等、既存の企業向けAI契約を利用可能か調査する

## Development Environment

- [ ] Swagger UIを利用したAPI参照・実行環境を構築する
- [ ] Prism Mock Server環境を構築する
- [ ] Development EnvironmentにおけるDocker / Docker Composeの利用方式を決定する
- [ ] PostgreSQLを含むDevelopment Environmentを構築する
- [ ] Mock / 開発中API / 開発済みBackendの接続先切替方式を決定する
- [ ] Frontendの接続先切替方式を決定する

## CLI

- [ ] Core CLIの構成を決める
- [ ] `ariadne validate` を実装する
- [ ] `ariadne generate openapi` を実装する
- [ ] `ariadne generate ddl` を実装する
- [ ] `ariadne version` を実装する
