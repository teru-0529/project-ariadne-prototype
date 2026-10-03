# Behavior Test Specification

## 1. Purpose

本ドキュメントは、Prototype Phase 5のBehaviorに対して、どのTest境界から何を確認することでBehaviorを保証するかを定義する。

Behaviorそのものの定義は `behavior-map.md` を正本とする。

本ドキュメントでは、Behaviorごとに必要となる以下のTestを整理する。

- Frontend Module Test
- Backend Module Test
- Integration Test

Test Codeそのものの詳細設計を目的とはせず、Behaviorに対して何をTestするかを明確にすることを目的とする。

---

## 2. S01 Applicationを起動する

### Frontend Module Test

#### Given

- Applicationを初期状態でRenderする

#### Expected

- Headerに `PROJECT ARIADNE` が表示される
- Primary Navigationに `Task` が表示される
- Primary Navigationに `Config` が表示される
- `Task` が選択状態である
- Task Pageが表示される
- Task Side Navigationに `List` が表示される
- Task Side Navigationに `Create` が表示される
- `List` / `Create` のどちらも選択されていない
- Main Areaに `Select an operation.` が表示される
- Task List取得処理は呼び出されない
- Footerが表示される

### Backend Module Test

N/A

### Integration Test

N/A

---

## 3. Note

CSS、Layout、Color、Spacing、Size等のVisual Designは実装対象とするが、自動Module Testによる保証対象とはしない。

Visual DesignはBuild後の実機確認によって確認し、必要に応じて調整する。
