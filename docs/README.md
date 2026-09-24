# Project ARIADNE — Prototype Roadmap

> Follow the thread from design to code.

Prototypeの目的は、Project ARIADNE Coreへ進む前に、以下の2つを実証すること。

1. Svelte + Wails + Go による Windows Application の技術的成立性
2. Element / API / DDL を対象とした ARIADNE Model の設計的成立性

Windows Application では、Task管理アプリという小さな題材を使い、
画面操作からローカルYAMLの Read / Edit / Write までを一巡する。

ARIADNE Model では、正本・中間モデル・生成成果物の責務と変換構造を手作業で設計・検証する。

Validator / Resolver / Generator 等、ARIADNE Model を処理する機能そのものの実装は Core で行う。

---

Project ARIADNE / Prototype

Phase 0：設計原則
│
│ Prototypeの目的・スコープ
│ Windows Application / ARIADNE Model の検証方針
│ Prototype / Core の責務境界を確定
│
└── ✅ COMPLETE
↓
Phase 1：Repository / Workspace
│
│ Repository構造・Workspace・成果物配置等を確定
│
└── ✅ COMPLETE
↓
Phase 2：Element Definition
│
│ types.yaml / elements.yaml
│ 共通語彙・型・制約・責任範囲を確定
│
└── ✅ COMPLETE
↓
Phase 3：API Definition
│
│ API定義の正本モデル
│ Raw / Resolved Model
│ OASへの変換構造
│ Element Definitionとの責任分界
│
└── ✅ COMPLETE
↓
Phase 4：DDL Definition
│
│ DDL定義の正本モデル
│ Resolved Model
│ Table / Column / Constraint / Index
│ PostgreSQL DDLへの変換構造
│ Element Definitionとの責任分界
│
└── 🔨 NOW
↓
Phase 5：Windows Application
│
│ Task管理アプリを題材として
│ Svelte + Wails + Go による
│ Windows Application の技術的成立性を検証
│
│ Application Shell
│ ├─ Task一覧
│ └─ 選択したTaskの編集領域
│
│ File I/O
│ ├─ YAML Read
│ ├─ YAML Edit
│ ├─ YAML Write
│ └─ Restart / Restore
│
│ Build
│ └─ ariadne-prototype.exe
│
└── ⬜
↓
Phase 6：Prototype Review
│
│ Windows Application 技術検証
│ ARIADNE Model 設計検証
│ Prototype / Core 責務境界
│ 残課題 / 技術リスク
│ Coreへの持越し事項
│
│ Phase 0 Definition of Done を最終確認
│
└── ⬜
↓
🏁 Prototype COMPLETE
↓
Project ARIADNE Core
