# Core TODO

## Source / Model

- [ ] Phase 2 Element Definition Loaderを実装する
- [ ] Phase 3 API Definition Loaderを実装する
- [ ] Source → Raw Model生成を実装する
- [ ] Raw Model → Resolved Model生成を実装する
- [ ] Raw / Resolved Modelのデバッグ出力を実装する

## Validation

### Phase 2

- [ ] Phase 2 Validation Rule `V-01` ～ `V-22` を実装する
- [ ] DECIMAL の値が `precision / scale` に適合することをValidationする
- [ ] Naming Recommendation Warningを実装する

### Phase 3

- [ ] File Validationを実装する
- [ ] API Validationを実装する
- [ ] Service Validationを実装する
- [ ] OpenAPI Validationを実装する
- [ ] Phase 3 Validation Rule `V-23` ～ `V-78` を実装する

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
- [ ] built-in `/health` 生成を実装する
- [ ] built-in `/version` 生成を実装する
- [ ] `externalDocs` 変換を実装する
- [ ] 未設定環境用の固定 `servers` を生成する

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

- [ ] Markdown → HTML変換を実装する
- [ ] Prototype / Coreで利用する`marked`の実行方式を決める
- [ ] Markdown内の最初のH1をHTML titleとして使用する
- [ ] H1が存在しない場合をErrorとする
- [ ] OAS成果物ディレクトリへ補足HTMLを生成する
- [ ] OpenAPI `externalDocs.url` を生成HTMLへの相対URLへ変換する

## Development Environment

- [ ] Backend開発用Swagger UI環境を用意する
- [ ] Prism Mock Server環境を用意する
- [ ] Swagger UI / Mock ServerをDocker Composeから起動できるようにする
- [ ] Swagger UI実行時に接続先Serverを設定・切替する方式を決める
- [ ] Mock / 開発中API / 開発済みDocker Imageを切り替えられる構成を検討する
- [ ] Frontendからも同様の接続先切替が可能か検討する
- [ ] Phase 4 DDL完成後、PostgreSQLを含むDevelopment Environmentを構築する
- [ ] Task Applicationを用いたEnd-to-End検証を行う

## CLI

- [ ] Core CLIの構成を決める
- [ ] `ariadne validate` を実装する
- [ ] `ariadne generate openapi` を実装する
- [ ] `ariadne generate ddl` をPhase 4で検討する
- [ ] `ariadne version` を実装する
