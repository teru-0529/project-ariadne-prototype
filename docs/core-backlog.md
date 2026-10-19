# Core Backlog

Prototypeの開発を通じて見つかった、Coreで実現・検討したいアイデアを記録する。

本Backlogへの記載は、Coreでの実施を確定するものではない。

Core開始時に棚卸しを行い、実施すると決定した項目は `core-todo.md` へ移す。

## UI / UX

- [ ] Side Navigationを折りたたみ可能にし、Main Areaを広く利用できるようにする
- [ ] Breadcrumb Navigationを導入し、Application内の現在位置を認識しやすくする
- [ ] Focus時のVisual Feedbackを整備し、Keyboard操作を含めて現在FocusしているElementを認識しやすくする
- [ ] Footerを設ける
  - Application Versionを表示する
  - Copyrightを表示する
- [ ] Git Integrationを利用する場合、ログイン中のGit AccountをHeaderまたはFooterに表示する
- [ ] 項目数が多い場合に備え、Side Navigation内をスクロール可能にする
- [ ] Side Navigationの項目を絞り込むFilter機能を検討する

## AI-assisted Editing

- [ ] ARIADNE Sourceの登録・編集をAIで支援するAI Modeを検討する
- [ ] 自然言語からElement / API / DDL Definitionの登録候補を生成する方式を検討する
- [ ] AIの生成結果を直接正本へ反映せず、利用者が確認・修正してから確定する操作モデルを設計する
- [ ] AI Provider / Model / API接続情報をARIADNEでどのように管理するか検討する
- [ ] AI利用に必要な契約・認証・利用制約を整理する
- [ ] GitHub Copilot等、既存の企業向けAI契約を利用可能か調査する

## Git Integration

- [ ] ARIADNE ApplicationからGit操作を行う必要性を検討する
- [ ] Git Integrationを採用する場合、status / diff / commit / push等の責務と実装方式を検討する
- [ ] ConfigからGit関連設定を行えるようにする
  - Git Account / 認証等の設定
  - 具体的な設定項目・認証方式はCoreで検討

