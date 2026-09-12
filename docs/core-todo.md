# Core TODO

## OAS Generation

- [ ] Git TagからApplication Versionを取得する
- [ ] Semantic Versionとして無効な場合は `0.0.0` にフォールバックする
- [ ] `/version` のレスポンスに同じApplication Versionを使用する

## ReDoc / Swagger UI

- [ ] Redocly CLIの実行環境をCoreに組み込む方法を決める
- [ ] Redocly CLIのバージョンをARIADNEリリース単位で固定する
- [ ] 通常利用時はRedocly Update Noticeを抑止する
- [ ] Telemetryを無効化する
- [ ] Backend開発用Swagger UI環境を用意する
- [ ] Swagger UI用OASにServer Variablesを設定する
- [ ] Mock / 開発中API / 開発済みDockerイメージを切り替えられる構成を検討する
- [ ] Frontendからも同様の接続先切替が可能か検討する

## Validation

- [ ] File Validation実装
- [ ] API Validation実装
- [ ] Service Validation実装
- [ ] OpenAPI Validation実装

## Generator

- [ ] Raw Model生成
- [ ] Resolved Model生成
- [ ] Element → OAS Schema変換
- [ ] Standard Error生成
- [ ] built-in `/health` 生成
- [ ] built-in `/version` 生成
