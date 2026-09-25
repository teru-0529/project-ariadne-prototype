# Project ARIADNE — Prototype Phase 5：Windows Application

**Status:** IN PROGRESS

---

## 1. 目的

Prototype Phase 5では、Svelte + Wails v3 + GoによるWindows Applicationを実装し、
Project ARIADNEのWindows Applicationとしての技術的成立性を検証する。

Phase 0で定義したとおり、Prototype ApplicationではTask管理アプリを題材として、
以下の一連の動作を実際に実装・確認する。

- ローカルYAMLからTaskを読み込む
- Task一覧を画面へ表示する
- Taskを新規作成する
- Taskを編集する
- TaskをローカルYAMLへ保存する
- Application再起動後に保存内容を復元する
- Svelte / Wails / Go間で正常系・異常系を適切に伝達する
- Windows実行ファイルとしてBuildし、単体で起動・操作できることを確認する

Phase 5は、Phase 2〜4で設計したElements / API / DDLをApplicationから操作するPhaseではない。

Elements / API / DDLについてはPrototypeではモデルおよび変換構造の設計検証までとし、
それらを扱うApplication機能、Validator / Resolver / Generator等の実装はCoreで行う。

---

## 2. Phase 5の基本方針

Phase 5では、完成形を事前に詳細設計してから実装するのではなく、
小さな機能単位で実装・テスト・Build・実機確認を繰り返しながら進める。

基本方針を以下とする。

- 小さな単位で機能追加しながら進める
- 詳細な設計ドキュメントの作成は必須としない
- 実装を通じて確定した重要な設計判断や、後から参照する価値のある事項のみ本ドキュメントへ記載する
- コードから明らかな内容をドキュメントへ重複記載しない
- Prototypeでは内部Databaseを使用せず、Fileを使用する
- 機能追加の節目ごとにBuildを行い、リリースモジュールとして実際に起動・操作できることを確認する
- テストを通過していない状態ではBuildしない
- `test.sh` / `build.sh` をPhase 5初期に整備し、その後の開発サイクルの基盤とする
- 未知の技術要素については、一度に完成形まで実装せず、小さな単位で動作を確認しながら段階的に進める

Phase 5では、動くApplicationを作ることだけでなく、
Svelte / Wails / Go / Fileの責務を説明できる状態で実装することを重視する。

---

## 3. Application Architecture

### 3.1 基本構造

Phase 5 Applicationの基本構造は以下とする。

```text
Svelte
  ↓
Wails generated binding / runtime
  ↓
Go
  ↓
File
```

FrontendはSvelte + TypeScript、
Desktop / BackendはWails v3 + Goを使用する。

TaskデータはローカルYAMLをSource of Truthとする。

---

### 3.2 責務境界

Phase 5では、以下の責務を基本として実装を進める。

#### Svelte

SvelteはユーザーとのInteractionおよび画面状態を担当する。

主な責務は以下とする。

- 画面表示
- ユーザー入力
- Navigation
- Task一覧表示
- Task新規登録画面
- Task編集画面
- 入力Validationの表示
- Confirmation Modal
- Notification Toast
- Frontend上の入力状態・画面状態の管理

#### Wails

WailsはSvelteとGoを接続するApplication境界として利用する。

Wails generated bindingsはFramework管理のgenerated filesとして扱い、
原則として手動編集しない。

#### Go

GoはFrontendから利用されるApplication ServiceおよびFile操作を担当する。

主な責務は以下とする。

- TaskデータのRead
- TaskデータのWrite
- Task登録・更新に必要なBackend処理
- File I/O Error等のBackend Errorの返却

GoはFrontendの表示方法を認識しない。

例えばFile Write Errorが発生した場合、
GoはErrorをFrontendへ返すところまでを責務とし、
Toast等による表示方法はSvelte側の責務とする。

#### File

PrototypeではTaskデータをFileで管理する。

Task群は1つのYAML Fileで管理し、
全量Read / 全量Writeを基本とする。

---

## 4. Application Structure

Phase 5では、以下を基本的なApplication構造とする。

```text
Application
│
├─ Header
│   └─ PROJECT ARIADNE
│
├─ Primary Navigation
│   ├─ Task
│   └─ Config
│
├─ Task Page
│   ├─ Side Navigation
│   │   ├─ List
│   │   └─ Create
│   │
│   └─ Main Area
│       ├─ Task List
│       ├─ Task Create
│       └─ Task Edit
│
├─ Config Dummy Page
│
├─ Confirmation Modal
│
└─ Notification / Validation Presentation
    ├─ Field Validation
    └─ Notification Toast
```

### 4.1 Header

HeaderはApplication自体を識別する領域とする。

Prototypeでは `PROJECT ARIADNE` を表示する。

---

### 4.2 Primary Navigation

Applicationの上位機能を切り替えるNavigationとする。

Prototypeでは以下を配置する。

- Task
- Config

Primary Navigationの選択によりMain Areaの表示対象を切り替える。

Task選択時はTask Pageを表示し、Task用のSide Navigationを表示する。

Config選択時はConfig Dummy Pageを表示し、
Task用のSide Navigationは表示しない。

---

### 4.3 Task

TaskはPrototype Applicationの実機能とする。

#### List

Task YAMLから読み込んだTaskを一覧表示する。

一覧からTaskを選択することで、対象Taskの編集画面へ遷移する。

Prototypeでは以下を対象外とする。

- Sort
- Filter
- Search
- Pagination

Taskは全件を一覧表示する。

表示件数が一覧表示領域を超える場合は、
Main Area内の一覧領域にスクロールバーを表示する。

#### Create

Taskの新規登録画面を表示する。

#### Edit

既存Taskの編集画面を表示する。

`createdAt / updatedAt` はApplicationが管理し、
ユーザーから直接編集させない。

---

### 4.4 Config

Configは将来的なApplication設定領域を想定する。

例えばTask YAMLのPath等を設定する可能性があるが、
Prototype Phase 5では設定機能そのものは実装しない。

ただし、Primary Navigationによる画面切替を検証するため、
Dummy Config Pageは実装する。

---

### 4.5 Confirmation Modal

Taskの登録・更新等、
ユーザーによる最終確認が必要な操作ではConfirmation Modalを使用する。

ModalはApplication処理の結果通知ではなく、
ユーザーに意思決定を求めるUIとして扱う。

登録・更新操作時には、
Confirmation Modalを表示する前に必要なValidationを実行する。

Validationに失敗した場合はConfirmation Modalを表示せず、
対象Form FieldへValidation Errorを表示して入力内容の修正を求める。

Validationに成功した場合のみConfirmation Modalを表示する。

---

### 4.6 Validation / Notification

ErrorおよびApplicationからの通知は、
内容に応じて表示方法を使い分ける。

基本的な考え方は以下とする。

```text
Validation Error
      ↓
Form Field

Application Error
      ↓
Error Toast

Application Info
      ↓
Info Toast
```

#### Field Validation

ユーザーが入力内容を修正することで解消できるErrorは、
対象項目付近へ表示する。

基本的な表示方法は以下とする。

- 対象項目の直下にError Messageを表示する
- 対象項目を視覚的に強調する
- Error状態が認識できる程度の背景色・Border等を使用する

単独項目で判定できるValidationは、
原則としてFocus Out時に実行する。

例：

- 必須入力
- 文字数
- 単独項目内で完結する形式チェック

登録・更新操作時には、
Focus Outの有無にかかわらず全単独項目Validationを再実行する。

複数項目の関係によって判定する相関Validationは、
原則として登録・更新操作時に実行する。

相関Validation Errorについても、
ユーザーが修正すべき項目を特定できる場合は
Toastではなく関連するForm Field付近へ表示する。

```text
Focus Out
    ↓
対象項目の単独Validation


Create / Update
    ↓
全単独項目Validation
    +
相関Validation
    ↓
NG
    → Form FieldへError表示

OK
    → Confirmation Modal
```

Prototypeで実際に必要となるValidation Ruleは、
Task Dataおよび各画面の実装時に決定する。

存在しないValidation Ruleを技術検証のためだけに追加することはしない。

#### Notification Toast

入力項目の修正では解消できないApplication Errorおよび
Application処理結果の通知には、
画面右下のNotification Toastを使用する。

Prototypeでは以下を基本とする。

| Severity | 用途                              | 消去                        |
| -------- | --------------------------------- | --------------------------- |
| Error    | File I/O失敗等のApplication Error | ユーザーがCloseするまで表示 |
| Info     | 登録・更新成功等の情報通知        | 一定時間後に自動消去        |

Frontend / Backendのどちらで発生したかではなく、
Errorの意味に応じてPresentation方法を決定する。

Warning等、その他のSeverityは必要性が生じた場合に検討する。

---

## 5. Task Data

Phase 0で決定したTaskを使用する。

最小項目は以下とする。

```text
id
title
description
priority
status
createdAt
updatedAt
```

Task群は1つのYAML Fileで管理する。

具体的なTask Modelおよび `tasks.yaml` のFormatは、
Step 1で決定する。

Step 1では少なくとも以下を確定する。

- YAML全体構造
- 各項目の表現
- IDの形式
- `priority` の値
- `status` の値
- `createdAt / updatedAt` の形式
- Empty状態の表現

### 5.1 File存在条件

Prototypeでは `tasks.yaml` が事前に存在することを前提とする。

Applicationによる自動生成および初期化処理は実装しない。

`tasks.yaml` が存在しない場合は、
空のTask一覧として扱わずRead Errorとする。

```text
tasks.yamlあり
  │
  ├─ Taskあり
  │    → 正常
  │
  └─ Task 0件
       → 正常

tasks.yamlなし
       → Error
```

Project / Workspaceの初期化方式については、
Coreで必要に応じて検討する。

Phase 5では以下を実装する。

- Task一覧表示
- Task新規作成
- Task編集
- YAMLへの保存
- Application再起動後のYAMLからの復元

削除はPrototypeの必須機能としない。

---

## 6. Test Strategy

Phase 5では、テストを技術要素ごとではなく、
保証する範囲によってA / B / Cに分類する。

基本戦略は **A + B** とする。

CはBを適切な粒度に保つため、
または実装を小さな段階に分けて確認するための補助的なテストと位置付ける。

---

### 6.1 A. Integration Test

Svelte / Wails / Go間の実際の接続を保証する。

```text
Svelte
  ↓
Wails generated binding / runtime
  ↓
Go
```

- Mockを使用せず、実際のWails経由で確認する
- 全機能・全パターンを確認するものではない
- 代表的なVertical Sliceを少数通し、Svelte ↔ Goの境界そのものが成立していることを保証する
- Phase 5ではFileまで含めた実経路を確認対象としてよい

例：

```text
test data
   ↓
File
   ↓
Go
   ↓
Wails
   ↓
Svelte
   ↓
画面表示
```

具体的な自動化方法はPhase 5実装時にWails v3環境で検証して決定する。

---

### 6.2 B. Behavior Test

公開された振る舞いを単位としてテストする。

Phase 5におけるTDDは、
原則としてこのBを中心に行う。

Use CaseがFrontendのみで完結する場合は、
Front Bのみを作成する。

FrontendとBackendの双方を利用するUse Caseでは、
原則としてFront BとBack Bに分けてBehaviorを保証する。

```text
FrontendのみのUse Case

Front B
  Svelte
    ↓
Behavior


Frontend + BackendのUse Case

Front B
  Svelte
    ↓
Go Mock

Back B
  Go Public API
    ↓
File

代表経路
  ↓
Integration A
```

すべてのBehaviorについてIntegration Testを作成するのではなく、
実際の境界接続を保証する価値のある代表的な経路をAとして確認する。

#### Front B

Front Bでは、
ユーザー操作を起点としたFrontendの振る舞いを確認する。

Backendを利用するBehaviorではGo側をMockし、
ユーザー操作に対するBackend Responseを再現する。

例：

- Task / Configを選択すると対象Pageへ切り替わる
- Listを選択するとTask一覧取得を要求し、取得したTaskを一覧表示する
- Taskを選択すると編集画面が表示される
- Task一覧取得に失敗した場合にError Toastが表示される
- Focus Outにより対象項目のValidationが実行され、Errorがある場合はFieldへ表示される
- 登録操作を行いValidationに成功するとConfirmation Modalが表示される

デザインそのもの
（色、余白、サイズ等）は原則として自動テスト対象とせず、
実機確認する。

Svelte側のテスト基盤は、
Vitest + Svelte Testing Libraryを基本候補とする。

#### Back B

Svelteから呼ばれるPublic API / Service等、
外部から見た公開単位をテストする。

内部関数の構成は原則として意識せず、

```text
入力
 ↓
Public API / Service
 ↓
期待する結果
```

を確認する。

PrototypeではFile I/Oも基本的にBehavior Testへ含める。

```text
test data / temporary file
       ↓
Go Public API
       ↓
expected result
```

File System自体をMockすることを基本とはせず、
`t.TempDir()` 等を利用して実際のFile I/Oを行う。

---

### 6.3 C. Logic / Development Step Test

Cは必要な場合のみ使用する。

Cを利用する主な目的は以下の2つとする。

#### 1. Logicの独立保証

BのBehavior Testだけではパターン数が過大になる場合、
複雑な内部Logicを切り出して直接テストする。

GoではTable-Driven Testを基本候補とする。

Frontend側でも、
複雑なTypeScript Logicが存在する場合のみ純粋Logicとして切り出し、
Vitestで直接テストする。

File I/OについてもBを基本とするが、
パス解決、File判定等で多数のパターンが発生する局所Logicについては
Cとして切り出してよい。

#### 2. 段階的な実装確認

未知の技術要素を含むBehaviorについて、
最終的なBまで一度に実装せず、
小さな技術単位で動作を確認しながら段階的に進めたい場合にもCを利用してよい。

例えばTask YAML Loadでは、
最終的なBack Bへ到達するまでに以下のような段階を踏むことができる。

```text
[C] Fileを読み込める
       ↓
[C] YAMLをParseできる
       ↓
[C] Task Modelへ変換できる
       ↓
[B] LoadTasks()からTask一覧を取得できる
```

この用途のCは、
完成後の仕様保証だけを目的とするものではなく、
実装中に「ここまで動いた」ことを確認して次へ進むための足場として利用する。

上位のBが成立した後、
そのCを継続して保持する価値がない場合は削除してよい。

一方、
内部Logicを直接保証する価値がある場合は継続して保持する。

Cを利用する場合でも、
実装内部を網羅的にTestすることを目的とはしない。

---

### 6.4 A / B / Cの関係

```text
A. Integration Test
────────────────────────
Svelte → Wails → Go
境界そのものを保証
代表経路のみ


B. Behavior Test
────────────────────────
       ★ TDDの主戦場 ★

Front B               Back B
Svelte                 Go
   │                    │
Go Mock             Public API
                        │
                       File

公開された振る舞いを保証


C. Logic / Development Step Test
────────────────────────
       必要な場合のみ

① 複雑なLogicを独立して保証する

② Bへ到達するまで
   小さな段階で動作を確認する

上位B成立後、
保証価値のない足場のCは削除してよい
```

---

## 7. TDD Development Cycle

Phase 5では、
実装開始前にPhase全体の主要なUse Case / Behaviorを
Behavior Mapとして整理する。

Behavior MapはPhase 5で実現する振る舞いの全体像を把握するためのものであり、
詳細なTest Caseを事前にすべて定義するものではない。

Behavior Mapでは、
各Use CaseについてFrontend / Backendの関与を整理する。

例：

```text
Use Case
│
├─ Frontendのみ
│    └─ Front B
│
└─ Frontend + Backend
     ├─ Front B
     ├─ Back B
     └─ Integration A
          ※ 代表経路のみ
```

各StepではBehavior Mapを基に、
実装対象となるBehaviorのTest Caseを実装直前に具体化する。

基本的なDevelopment Cycleは以下とする。

```text
Behaviorを選択
       ↓
Front / Backの関与を確認
       ↓
必要なBを決定
       ↓
B：RED
       ↓
実装
       │
       ├─ そのまま実装可能
       │      ↓
       │   B：GREEN
       │
       └─ 小さな段階で確認したい
              ↓
           C：RED
              ↓
           最小実装
              ↓
           C：GREEN
              ↓
           次の段階
              ↓
           Bへ戻る
              ↓
           B：GREEN
       ↓
必要ならRefactor
       ↓
全テスト実行
       ↓
代表経路の場合
       ↓
Integration A
```

Cは、
複雑なLogicを切り出す場合だけでなく、
未知の技術要素を一段ずつ確認しながら実装するためにも利用できる。

Phase 5では、

「ここまで動いた。じゃあ次。」

と確認できる単位で実装を進めることを重視する。

---

## 8. test.sh / build.sh

Phase 5の機能実装開始前に、
Test / Buildの実行基盤を整備する。

### 8.1 test.sh

基本イメージを以下とする。

```text
test.sh
├─ Go tests
└─ Svelte / TypeScript tests
```

TDD中に頻繁に実行するため、
高速に完了することを重視する。

Integration Testについては、
実際のWails v3での自動化方式、実行時間、安定性を確認した上で、
`test.sh`への組み込み方を決定する。

---

### 8.2 build.sh

`build.sh` はTestをBuild Gateとして実行した上で、
Wails ApplicationをBuildする。

基本イメージを以下とする。

```text
build.sh
├─ Test
│   ├─ Behavior / Logic Test
│   └─ Integration Test
│
└─ Wails Application Build
       ↓
   ariadne-prototype.exe
```

`test.sh` にIntegration Testを含めるかどうかは、
Integration Testの実装方式、実行時間、安定性を確認して決定する。

ただし `build.sh` によるRelease Buildの前には、
Phase 5で必要なTestをすべて実行し、
Integration Testの成功もBuild条件とする。

Testを通過していない状態ではApplication Buildを行わない。

Build後は生成されたWindows実行ファイルを実際に起動し、
Applicationとして動作することを節目ごとに確認する。

---

## 9. Implementation Steps

Phase 5は以下のStepで進める。

各Stepでは、
事前に詳細設計を完成させることよりも、
必要な設計判断を行いながら小さく実装・検証することを優先する。

---

### Step 1：Application Architecture / Task Data

Svelte / Wails / Go / Fileの責務境界と、
Applicationが扱うTask Dataの基本仕様を整理する。

主な確認事項：

- Svelteの責務
- Wailsの位置付け
- Go Backend Serviceの責務
- File Accessの責務
- Frontend / Backend間のError伝達
- Public API / Serviceの単位
- Task Model
- `tasks.yaml` Format
- `tasks.yaml` のFile存在条件
- Empty状態の表現
- Frontend / Backend間で受け渡すData

---

### Step 2：Test / Build Foundation

Phase 5全体で使用するTest / Build基盤を構築する。

主な実施事項：

- Phase 5 Behavior Mapの整理
- 各Use CaseのFrontend / Backend関与の整理
- Go Test環境確認
- Svelte / TypeScript Test環境構築
- Behavior Testの最小技術検証
- Logic / Development Step Testの最小技術検証
- Integration Test方式の技術検証
- `test.sh` 整備
- `build.sh` 整備
- Build Gateの確立

Behavior Mapでは、
Step 3以降で実現する主要なUse Case / Behaviorを整理する。

詳細なTest Caseおよび実装途中で必要となるCは、
各Stepの実装時に具体化する。

---

### Step 3：Application Shell

Applicationの基本的な画面構造を構築する。

主な実施事項：

- Header
- Primary Navigation
  - Task
  - Config
- Config Dummy Page
- Side Navigation
  - List
  - Create
- Main Area
- Confirmation Modalの基本構造
- Notification Toastの基本構造
- Field Validation表示の基本構造
- Prototype用の暫定Visual Design

このStepでは各業務機能を完成させることを目的とせず、
後続Stepで機能を追加できるApplication Shellを成立させる。

Primary Navigationについては、
Task / Configの選択によるPage切替をBehavior Testで確認する。

---

### Step 4：Task List / Read

Step開始時にTask List画面の簡易画面設計を行う。

主に以下を確認する。

- 表示項目
- 一覧領域
- Task選択方法
- Editへの遷移
- スクロール方法
- Error表示

詳細な画面設計書は作成せず、
実装に必要なUI構造を決定する。

Task YAMLを読み込み、
Task一覧として画面へ表示するVertical Sliceを構築する。

```text
Task YAML
    ↓
Go
    ↓
Wails
    ↓
Svelte
    ↓
Task List
```

主な実施事項：

- Task YAML Read
- Go Backend Service
- Wails経由でのFrontend呼び出し
- Task List表示
- Task選択
- Task Edit画面への遷移
- Task 0件時の正常動作
- `tasks.yaml` 不存在時のRead Error
- Backend ErrorのFrontendへの伝達
- Error Toast
- 一覧領域内のスクロール

PrototypeではSort / Filter / Search / Paginationは実装しない。

`tasks.yaml` が存在しない場合は自動生成せず、
Read Errorとして扱う。

このError Caseにより、
File → Go → Wails → Svelte → Notificationまでの
Error伝達経路も検証する。

Task YAML Loadについては、
必要に応じてCを足場として利用し、
File Read → YAML Parse → Task Model変換等を
小さな単位で確認しながら実装する。

---

### Step 5：Task Create

Step開始時にTask Create画面の簡易画面設計を行う。

主に以下を確認する。

- 項目配置
- 入力Component
- Button配置
- Validation Error表示
- Confirmation Modal
- Notification

Taskの新規作成を、
入力からFileへの保存まで一つのUse Caseとして実装する。

```text
Create
  ↓
入力
  ↓
Validation
  ↓
Confirmation Modal
  ↓
Confirm
  ↓
Svelte
  ↓
Wails
  ↓
Go
  ↓
Task YAML Write
  ↓
Info Toast
```

主な実施事項：

- Task Create画面
- Task入力
- Focus Out時の単独項目Validation
- 登録時の全項目Validation
- 必要な場合の相関Validation
- Validation ErrorのField表示
- `createdAt / updatedAt` のApplication管理
- Confirmation Modal
- Confirm時のTask登録
- Cancel時に登録しないこと
- Task YAML全量Write
- Write Errorの返却
- Error Toast
- 保存成功時のInfo Toast

正常系と、
Create Use Case内で発生する異常系を同じStepで検証する。

---

### Step 6：Task Edit

Step開始時にTask Edit画面の簡易画面設計を行う。

主に以下を確認する。

- 項目配置
- 入力Component
- Button配置
- Validation Error表示
- Confirmation Modal
- Notification

既存Taskの編集を、
Task選択からFileへの保存まで一つのUse Caseとして実装する。

```text
Task List
  ↓
Task選択
  ↓
Edit
  ↓
変更
  ↓
Validation
  ↓
Confirmation Modal
  ↓
Confirm
  ↓
Svelte
  ↓
Wails
  ↓
Go
  ↓
Task YAML Write
  ↓
Info Toast
```

主な実施事項：

- 既存TaskのEdit画面表示
- Task編集
- Focus Out時の単独項目Validation
- 更新時の全項目Validation
- 必要な場合の相関Validation
- Validation ErrorのField表示
- `createdAt` の維持
- `updatedAt` のApplication管理
- Confirmation Modal
- Confirm時のTask更新
- Cancel時に更新しないこと
- Task YAML全量Write
- Write Errorの返却
- Error Toast
- 保存成功時のInfo Toast

正常系と、
Edit Use Case内で発生する異常系を同じStepで検証する。

---

### Step 7：Windows Application Verification

Phase 5で構築したApplicationをWindows実行ファイルとしてBuildし、
Prototypeとしての成立性を最終確認する。

Release Build前には、
Behavior / Logic TestおよびIntegration Testを実行し、
すべて成功していることを確認する。

最終的なVertical Sliceは以下とする。

```text
Task YAML
    ↓ Read
Go
    ↓
Wails
    ↓
Svelte
    ↓
表示
    ↓
Create / Edit
    ↓
Validation
    ↓
Confirmation
    ↓
Svelte
    ↓
Wails
    ↓
Go
    ↓ Write
Task YAML
    ↓
Application再起動
    ↓
Read
    ↓
保存内容復元
```

主な確認事項：

- 全Behavior / Logic Test成功
- Integration Test成功
- Release Build成功
- `ariadne-prototype.exe` 単体起動
- Task List表示
- Task Create
- Task Edit
- Task保存
- Application再起動
- 保存内容の復元
- 主要Error Caseの実機確認

---

## 10. Visual Design

Visual DesignそのものはPrototype Phase 5の主目的としない。

PrototypeではApplicationとして操作可能であり、
各UI要素の役割を識別できる程度のデザインを行う。

基本的な方向性として以下を想定する。

- 白〜薄いグレーを基調とする
- Blue / Purple系をAccent Colorとする
- 装飾を抑えたシンプルなUIとする
- 情報そのものを主役とし、UIの視覚的主張を強くしすぎない
- Validation Errorは対象項目を識別できる程度に強調する
- Application Error等、意味を持つColorは必要に応じて使い分ける

Application全体の暫定Visual DesignはStep 3で適用する。

Task List / Create / Editの具体的な画面構成は、
Step 4〜6の各Step開始時に簡易設計する。

色、余白、Typography、Component Design等の
詳細なDesign SystemはCoreで必要に応じて検討する。

---

## 11. Future Consideration

将来的に、ApplicationからAIを利用した登録・編集支援を行う可能性がある。

Application Shellとしては、
Task / Config等の機能Navigationとは別に、
AI Modeを起動する入口を配置できる構造を妨げないものとする。

AIは「何を扱うか」を表すApplication機能ではなく、
既存機能を別の操作方法から利用する支援機能として位置付けることを想定する。

ただし、AI Mode、Prompt UI、AI API連携等は
Prototype Phase 5では一切実装しない。

---

## 12. Prototype Phase 5 対象外

Phase 5では以下を対象外とする。

- Elements編集UI
- API編集UI
- DDL編集UI
- Validator / Resolver / Generator実装
- OAS / DDL生成機能
- 内部Database
- Project / Workspace初期化機能
- `tasks.yaml` の自動生成
- DB Access Test
- git等の外部Command実行に関するTest
- Configの実機能
- Task ListのSort / Filter / Search / Pagination
- Task削除機能
- AI Mode
- AI API連携
- 詳細なDesign System
- GitHub ActionsによるCI/CD

Phase 5では、
Svelte + Wails + Go + FileによるWindows Applicationとしての成立性確認に集中する。

---

## 13. Phase 5 Completion Criteria

以下をすべて満たした時点で、
Prototype Phase 5をCOMPLETEとする。

- Svelte / Wails / Go / Fileの責務境界を説明できる
- Task Modelおよび `tasks.yaml` Formatが確定している
- `tasks.yaml` が存在しない場合をRead Errorとして扱える
- Task YAMLを読み込める
- Task一覧を表示できる
- Task 0件を正常状態として扱える
- Taskを新規作成できる
- 既存Taskを編集できる
- TaskをYAMLへ保存できる
- Application再起動後に保存内容を復元できる
- Primary NavigationによりTask / Config Dummy Pageを切り替えられる
- 単独項目Validationを適切なタイミングで実行できる
- 登録・更新時に必要なValidationを実行できる
- Validation Errorを対象Form Field付近へ表示できる
- Application ErrorをError Toastとして通知できる
- Info Toastが動作する
- Confirmation Modalが動作する
- Phase 5の主要BehaviorがBehavior Mapとして整理されている
- Use CaseごとにFront / Backの関与を整理できる
- Behavior Testを中心とした自動テストが成立する
- 必要に応じてCをLogic保証または段階的な実装確認に利用できる
- 代表的なSvelte ↔ Wails ↔ GoのIntegrationを自動テストで確認できる
- `test.sh` により日常的なTestを実行できる
- Release Build前にIntegration Testを含む必要なTestがすべて実行される
- Test成功後に `build.sh` からWindows ApplicationをBuildできる
- `ariadne-prototype.exe` を単体起動し、一連のTask操作を実行できる
