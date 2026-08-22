# Project ARIADNE Prototype

Project ARIADNE のアーキテクチャ及び技術要素を検証するための Prototype。

本Repositoryは検証を目的としており、
製品版（Project ARIADNE）は別Repositoryとして構築する。

## Repository Structure

- `app/`
  - Prototypeアプリケーション
- `src/`
  - ARIADNEが扱う設計ソース
  - 項目定義、OAS、Database定義などを配置する
- `dist/`
  - ARIADNEによって生成・統合された最終成果物
  - Prototypeでは生成結果の検証のためGit管理対象とする
- `tools/`
  - 定義変換、OAS統合、DDL生成などのツール・スクリプト
- `templates/`
  - Prototypeアプリケーションで使用するテンプレート
- `runtime/`
  - アプリケーション実行時のデータ
  - 原則としてGit管理対処外
- `docs/`
  - 設計資料・ドキュメント

## Prototype Policy

Prototypeでは、個々の技術要素だけでなく、
設計ソースから成果物およびアプリケーションコードまでの
一連のフローが成立することを検証する。
