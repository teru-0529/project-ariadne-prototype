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
Wails Generated Binding / Runtime
  ↓
Go Application Service
  ↓
Persistence Port
  ↑
File Persistence Adapter
  ↓
tasks.yaml
```

Phase 5では **Go Service First** を採用する。

Wails v3ではGo Serviceを起点としてFrontend向けBindingを生成する構造を基本とし、
Wails専用のController Interface / Adapter等は追加しない。

Go Application ServiceはWailsから公開されるServiceである一方、
Wailsに依存しない通常のGo Application Serviceとして実装する。

WailsはGo Application ServiceをFrontendへ公開するApplication境界として利用し、
Application Service自体にはWailsとの接続に関する責務を持たせない。

Persistenceについては、
Application ServiceがFile / YAMLを直接認識しないようPersistence Portを設ける。

File Persistence AdapterはPersistence Portを実装し、
Taskの永続化をFile / YAMLとして実現する。

この構造により、
Application ServiceからPersistence技術を分離し、
将来的にFile以外のPersistence方式を利用する場合にも
Application Serviceの変更を最小化できる構造とする。

Wails Generated BindingはFramework管理のgenerated filesとして扱い、
手動編集しない。

FrontendはSvelte + TypeScript、
BackendはGo、
Desktop Application FrameworkはWails v3を使用する。

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
- Frontend Validation
- Validation ErrorのForm表示
- Confirmation Modal
- Notification Toast
- Frontend上の入力状態・画面状態の管理

Frontend ValidationはUXのためのValidationとしてSvelte側で実行する。

Focus Out時の単独項目Validation、
Create / Modify Button押下時の全項目Validationおよび相関ValidationはSvelte側で実行する。

Frontend Validationに失敗した場合はBackendを呼び出さず、対象Form FieldへErrorを表示する。

#### Wails

WailsはSvelteとGo Application Serviceを接続するApplication境界として利用する。

Wails Generated BindingはFramework管理のgenerated filesとして扱い、
原則として手動編集しない。

Application内部から返されたErrorについては、
原因情報を失わない状態でWails Application Boundaryまで伝達する。

Wails Application Boundaryでは、
内部診断に必要な情報とFrontendへ公開する情報を分離する。

```text
Internal Error
    ↓
Wails Application Boundary
    ├─ Diagnostic Information → Log
    └─ Public Error           → Svelte
                                  ↓
                               Error Toast
```

具体的なWails v3でのError Boundary実現方式は、Phase 5実装時に技術検証して決定する。 TODO:

#### Go Application Service

Go Application Serviceは、
Frontendから要求されるApplication Operationを実現する。

Prototypeでは以下を基本的なPublic Serviceとする。

```text
ListTasks()
GetTask(ID)
CreateTask(Properties)
ModifyTask(ID, Properties)
```

Application Serviceは、Taskを成立・変更させる主体とする。

Application ServiceはFile / YAML等のPersistence技術を直接認識せず、
Persistence Portを通じて永続化を要求する。

また、Application ServiceはFrontendの表示方法を認識しない。

PrototypeではBusiness TaskをそのままFrontendへ公開する。

ただし、
Application Serviceが外部へ公開する情報を決定する責務を持つものとし、
Business ModelとFrontendへ公開するModelが常に同一であることをArchitecture上の前提とはしない。

#### Persistence Port

Persistence Portは、
Application ServiceからPersistenceへ要求する操作を定義する。

Prototypeでは以下を基本とする。

```text
Retrieve()
Pick(ID)
Store(Task, StoreIntent)
```

各Operationの意味は以下とする。

- `Retrieve`
  - 格納されているTask群を取得する
- `Pick`
  - IDを指定して1件のTaskを取得する
- `Store`
  - Taskを格納する

`Store` ではApplication ServiceからPersistenceへ
格納Intentを明示する。

```text
Deposit
    新しいTaskを初めて格納する

Replace
    既に格納されているTaskを置き換える
```

Persistence AdapterはTaskの状態からIntentを推測しない。

将来、`Retrieve` に取得対象を指定する必要が生じた場合は、
引数としてCriteriaを導入できる。

Criteriaは概念的に以下を含む。

```text
Criteria
├─ Conditions
├─ Order
└─ Range
```

ただしPrototypeでは必要性がないためCriteriaは実装しない。

#### File Persistence Adapter

File Persistence Adapterは、
Persistence PortをFile / YAMLとして実現する。

主な責務は以下とする。

- `tasks.yaml` Read / Write
- YAML Parse / Serialization
- `format / version` Validation
- Persistence ModelとBusiness Taskの相互変換
- Deposit時のULID生成
- `createdAt / updatedAt` の生成・維持
- Safe Write
  - Temporary Fileへの全量Write
  - Canonical FileへのReplace

Application Serviceは、
ULID生成方式、Timestamp取得、YAML Format、
File I/O等を認識しない。

### 3.3 Frontend / Backend Data Boundary

Prototypeでは、
SvelteとGo Application Service間の基本Operationを以下とする。

```text
Svelte → Application Service

ListTasks()
GetTask(ID)
CreateTask(Properties)
ModifyTask(ID, Properties)
```

Application ServiceからSvelteへの返却は以下を基本とする。

```text
ListTasks()       → []Task
GetTask(ID)       → Task
CreateTask(...)   → storedTask
ModifyTask(...)   → storedTask
```

PrototypeではBusiness TaskをそのままFrontendへ公開する。

ただし、
Business TaskとFrontend公開Modelが常に同一であることを
Architecture上の制約とはしない。

将来、Application Serviceが公開する情報を制御する必要が生じた場合は、
Frontend公開用Modelを分離できるものとする。

#### Properties

Create / Modifyでは、
SvelteからApplication Serviceへ `Properties` を渡す。

```text
CreateTask(Properties)
ModifyTask(ID, Properties)
```

`Properties` は、
Taskへ与えるProperty値の集合を表す共通概念とする。

Create / ModifyというOperationの意味はService Method側が持ち、
Properties自体へOperationの意味を持たせない。

CreateではPropertiesを
新しいTaskを成立させるための材料として使用する。

```text
Properties
    ↓
newTask
```

Modifyでは、
Repositoryから取得した既存Taskへ
指定されたPropertiesを適用する。

```text
ID
 ↓
Pick
 ↓
storedTask
    +
Properties
 ↓
modifiedTask
```

ModifyのPropertiesでは、
Propertyごとに概念上以下の3状態を区別できる必要がある。

```text
未指定
    → 変更しない

指定 + null
    → nullへ変更

指定 + value
    → 指定値へ変更
```

Svelte / Wails / Go間で
この3状態をどのような型として表現するかは、
Phase 5実装時の技術検証課題とする。 TODO:

Step 1では具体的な実装型を固定しない。

### 3.4 Error / Logging

Application内部で発生したErrorは、
発生箇所で握り潰さず、
原因情報を維持したままWails Application Boundaryまで伝達する。

```text
Adapter / Repository / Service
    ↓
Error Chain
    ↓
Wails Application Boundary
```

下位層で発生した具体的なErrorと、
上位層が扱うErrorの意味は同一である必要はない。

例えばFile / YAMLに関する具体的な失敗は、
Persistenceとして意味のあるErrorへ変換してよい。

ただし、
元の原因情報は失わないものとする。

Application Serviceの振る舞いが変わらないErrorについて、
Application Serviceへ不要な詳細分類を公開しない。

例えば以下はPersistence内部では異なる原因であるが、
Application Serviceがそれらを区別して振る舞う必要がなければ
同一のRead Failureとして扱ってよい。

```text
File Read Failure
YAML Parse Failure
Invalid Format
Unsupported Version
```

一方、`Pick(ID)` で対象Taskが存在しないことは、
Persistence Operationとして意味のある結果であるため
Not Foundとして区別する。

#### FrontendへのError公開

Wailsを経由してFrontendへ返却されたErrorは、
Application Errorとして扱う。

```text
Internal Error
    ↓
Wails Application Boundary
    ├─ Diagnostic Information → Log
    └─ Public Message         → Svelte
                                  ↓
                               Error Toast
```

Frontendへ内部Error Chainそのものを公開しない。

Application Error Codeは設けない。

Phase 2〜4で定義したValidation Rule ID等、
Core機能そのものに意味を持つ識別子とは区別する。

#### Logging

本格的なApplication LoggingはPrototype Phase 5の必須実装とはしない。

ただし、
将来CoreでStructured Loggingを導入することを想定し、
PrototypeでもErrorの原因情報を失わない設計とする。

Loggingでは以下を基本的な考え方とする。

```text
Error
    → Error Chainとして上位へ伝達
    → Application Boundaryで診断情報を記録

Information / Debug
    → その情報を認識している場所で記録
```

Errorを各Layerで重複してLogすることは基本としない。

具体的なLogger、出力先、Log Rotation等は
Prototype Phase 5では設計・実装対象としない。

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

`createdAt / updatedAt` はPersistence側で管理し、
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

Frontend Validationは、
以下の2種類に分けて扱う。

```text
Single Field Validation
    単独Fieldで判定できるValidation

Correlation Validation
    複数Fieldの関係によって判定するValidation
```

Single Field Validationは、
原則としてFocus Out時に対象Fieldについて実行する。

Create / Modify Button押下時には、
Focus Outの有無にかかわらず全FieldのSingle Field Validationを再実行する。

Correlation Validationについては、
Formごとに以下のModeを持つ。

```text
INITIAL
    全Correlation Validationをまだ実行していない

ACTIVE
    全Correlation Validationを一度実行済み
```

INITIALでは、
Focus Out時にCorrelation Validationを実行しない。

Create / Modify Button押下時に、
全Single Field ValidationがValidとなった場合のみ
全Correlation Validationを実行する。

全Correlation Validationを実際に実行した時点で、
Correlation Validation ModeをACTIVEとする。

ACTIVEとなった後は、
Focus OutしたFieldに関連するCorrelation Validationを再実行する。

ただし、
Correlation Validationに関連するFieldが
Single Field Validation上Validでない場合、
そのCorrelation Ruleは評価しない。

その場合、
当該Ruleに対応する既存Correlation ErrorはClearする。

```text
Focus Out

Correlation Validation Mode = INITIAL
    ↓
対象FieldのSingle Field Validation
    ↓
終了


Correlation Validation Mode = ACTIVE
    ↓
対象FieldのSingle Field Validation
    ↓
関連FieldがすべてSingle Valid
    ↓
関連するCorrelation Validation
```

Create / Modify Button押下時は以下とする。

```text
Create / Modify
    ↓
全Single Field Validation
    ↓
Single Field Errorあり？
    ├─ YES
    │    ↓
    │  Correlation Validationを実行しない
    │
    └─ NO
         ↓
       全Correlation Validation
         ↓
       Correlation Validation Mode = ACTIVE
         ↓
       Correlation Errorあり？
         ├─ YES → Form FieldへError表示
         └─ NO  → Confirmation Modal
```

1 Fieldに表示するValidation Messageは1つとする。

Single Field ErrorとCorrelation Errorの双方が存在する場合は、
Single Field Errorを優先して表示する。

同一Fieldに複数のCorrelation Errorが存在する場合は、
Validation Priorityにより表示するErrorを決定する。

Prototype Phase 5で使用する具体的なValidation Ruleおよび
Validation Priorityについては、
別途 `validation-priority.md` に定義する。

Frontend ValidationはSvelte側で実行する。

Application Service側でもBusiness / Applicationの整合性を保証するため
Validationを実施する場合があるが、
Frontend Validationを正常に通過した入力が
Application Service側でValidation Errorとなることは異常系として扱う。

その場合、
Validation ErrorをFormへ戻すのではなくApplication Errorとして扱い、
Wails経由でFrontendへ返却しError Toastで通知する。

#### Notification Toast

入力項目の修正では解消できないApplication Errorおよび
Application処理結果の通知には、
画面右下のNotification Toastを使用する。

Prototypeでは以下を基本とする。

| Severity | 用途                              | 消去                        |
| -------- | --------------------------------- | --------------------------- |
| Error    | File I/O失敗等のApplication Error | ユーザーがCloseするまで表示 |
| Info     | 登録・更新成功等の情報通知        | 一定時間後に自動消去        |

Svelte上で検出したErrorのうち、
ユーザーが入力内容を修正することで解消できるValidation Errorは
Form Fieldへ表示する。

入力項目の修正では解消できないFrontend上のApplication Errorは、
Error Toastで通知する。

Wailsを経由してApplication Serviceから返却されたErrorは、
Application ErrorとしてError Toastで通知する。

```text
Svelte
├─ Input Validation Error
│      → Form Field
│
└─ Application Error
       → Error Toast

Wails経由のApplication Error
       → Error Toast
```

Warning等、その他のSeverityは必要性が生じた場合に検討する。

---

## 5. Task Data

Phase 5では、
Business Logicが扱うTaskとPersistence上のTask表現を分離する。

### 5.1 Business Task

Business Taskは以下を基本とする。

```text
Task
├─ id
├─ title
├─ description
├─ priority
└─ status
```

各Propertyは以下とする。

| Property      | 内容                                     |
| ------------- | ---------------------------------------- |
| `id`          | ULID。TaskのIdentityとして利用する       |
| `title`       | 必須。空文字・空白のみ不可。最大50文字   |
| `description` | 任意。複数行可                           |
| `priority`    | 必須。1〜5。値が大きいほどPriorityが高い |
| `status`      | 必須。`TODO / DOING / DONE`              |

Business Taskは `createdAt / updatedAt` を持たない。

`createdAt / updatedAt` はBusiness Logicで利用する情報ではなく、
Persistence上で管理する情報とする。

TaskはIDが未設定の状態でもBusiness Taskとして成立し得る。

Create時にはApplication ServiceがPropertiesから新しいTaskを成立させ、
PersistenceへDepositする時点でFile Persistence AdapterがULIDを生成する。

```text
CreateTask(Properties)
    ↓
newTask
    ↓
Validate
    ↓
validatedTask
    ↓
Store(validatedTask, Deposit)
    ↓
storedTask
```

Application ServiceはTaskのIdentityとしてIDを認識・利用するが、
IDの生成方式は認識しない。

### 5.2 Persistence Task

`tasks.yaml` 上ではBusiness TaskのPropertyに加えて、
以下を保持する。

```text
createdAt
updatedAt
```

Deposit時には以下を行う。

- ULIDを生成して `id` を設定する
- `createdAt` を現在日時で設定する
- `updatedAt` を `createdAt` と同一日時で設定する
- TaskをYAML末尾へ追加する

Replace時には以下を行う。

- `id` により既存Taskを特定する
- `createdAt` を維持する
- `updatedAt` を現在日時へ更新する
- YAML上のTask順序を維持する

TimestampはTimezoneを含むRFC3339 / ISO 8601形式を使用する。

Application ServiceはTimestampの取得・管理を行わない。

### 5.3 tasks.yaml Format

`tasks.yaml` のRootは以下とする。

```yaml
format: ariadne-tasks
version: "2.1"

tasks: []
```

`format` は `ariadne-tasks` 固定とする。

`version` はStringによる `Major.Minor` 形式とする。

Prototype Applicationが対応するVersionを `2.1` とし、
互換性は以下とする。

```text
2.1 → OK
2.0 → OK
1.0 → NG(Majorが異なる)
2.2 → NG(Minorが対応バージョンよりも大きい)
3.0 → NG(Majorが異なる)
```

同一MajorかつFile MinorがApplication対応Minor以下の場合のみ読み込み可能とする。

`format` が異なる場合、
または非対応Versionの場合はRead Errorとする。

Task 0件は以下で表現する。

```yaml
format: ariadne-tasks
version: "2.1"

tasks: []
```

`description` が空の場合は以下とする。

```yaml
description: ""
```

`description` に値がある場合は、
1行の場合を含めLiteral Block `|` を使用する。

複数行のdescriptionはRead / Write後も内容を維持する。

TaskのYAML Sequence順を一覧表示順として扱う。

PrototypeではSortを行わない。

CreateではTaskを末尾へ追加し、
Modifyでは既存Taskの位置を維持する。

### 5.4 File存在条件

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
       → Read Error
```

Project / Workspaceの初期化方式については、
Coreで必要に応じて検討する。

### 5.5 File Write

Task群は1つのYAML Fileで管理し、
全量Read / 全量Writeを基本とする。

Write時はCanonical Fileへ直接書き込まず、
Temporary Fileへ全量Writeした後、
Canonical `tasks.yaml` を置き換える。

```text
Task Data
    ↓
Temporary File
    ↓ Write成功
Canonical tasks.yamlをReplace
```

Temporary FileへのWriteに失敗した場合は、
既存のCanonical Fileを維持する。

Canonical FileのReplaceに失敗した場合はWrite Errorとする。

Backup GenerationおよびCrash Recovery機構はPrototypeでは実装しない。

Windows上でのFile Replaceの具体的な実現方式は、
Phase 5実装時に技術検証する。TODO:

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
ユーザーによる1つのActionをBehaviorの基本単位とする。

同一Actionに対して結果が分岐する場合は、
Behaviorを分割せずResult Patternとして整理する。

基本形式は以下とする。

```text
Behavior
├─ Start    : Action実行前の状態
├─ Action   : ユーザーによる1つのAction
├─ Result
│   └─ Pattern : Resultの分岐
└─ Scope    : Front / Back
```

各Behaviorでは、
Frontend / Backendの関与をScopeとして整理する。

```text
Front
    → Front B

Front + Back
    → Front B
    → Back B

代表的なVertical Slice
    → Integration A
```

すべてのBehaviorについてIntegration Aを作成するのではなく、
Svelte / Wails / Go間の実際の境界接続を保証する価値のある
代表的な経路のみIntegration Aとして確認する。

Phase 5のBehavior Mapは、
別途 `behavior-map.md` に定義する。

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

Step 1では以下を確定した。

- Go Service FirstによるApplication Architecture
- Svelte / Wails / Application Service / Persistence Port / File Persistence Adapterの責務境界
- Application ServiceのPublic Operation
  - `ListTasks`
  - `GetTask`
  - `CreateTask`
  - `ModifyTask`
- Persistence Operation
  - `Retrieve`
  - `Pick`
  - `Store`
  - `Deposit / Replace`
- Business Task Model
- Persistence Task Model
- ULIDおよびTimestampの責務
- `tasks.yaml` Format / Version
- File存在条件
- Empty状態
- Safe Write方針
- Frontend / Backend間で受け渡すData
- Create / Modifyで共通するProperties概念
- Frontend ValidationとApplication Errorの境界
- Error Chain / Wails Application Boundary / Loggingの基本方針

以下は実装時の技術検証課題とする。

- Modify Propertiesの
  `未指定 / null / value`
  をSvelte / Wails / Go間で表現する方式
- Wails v3におけるApplication Error Boundaryの具体的な実現方式
- Windows上でのSafe File Replace方式
- Integration Testの具体的な自動化方式

---

### Step 2：Test / Build Foundation

Phase 5全体で使用するTest / Build基盤を構築する。

機能実装を開始する前に、
Behavior Mapによって保証対象を整理した上で、
Go / FrontendそれぞれのTest Foundationを
最小のDummy Testを実際に動かしながら構築する。

その後、
日常的なTest実行およびBuild Gateを成立させる。

主な実施順序：

```text
2-1 Behavior Map
    ↓
2-2 Go Test Foundation
    ↓
2-3 Frontend Test Foundation
    ↓
2-4 test.sh
    ↓
2-5 build.sh / Build Gate
    ↓
2-6 Behavior / Logic Test 最小技術検証
    ↓
2-7 Integration Test 技術検証
```

#### 2-1 Behavior Map

Step 3以降で実現する主要なBehaviorを整理する。

Behaviorはユーザーによる1つのActionを基本単位とし、
同一Actionに対する結果の分岐はResult Patternとして整理する。

各Behaviorについて、
Front / Backの関与を明確にする。

詳細は `behavior-map.md` に定義する。

#### 2-2 Go Test Foundation

Go Testを実際に動かしながら、
Backend Testの基本構造を成立させる。

最小のDummy Testを作成し、
Go Testを実行できることを確認する。

#### 2-3 Frontend Test Foundation

Frontend Testを実際に動かしながら、
Svelte / TypeScript Testの基本構造を成立させる。

Vitest + Svelte Testing Libraryを基本候補とし、
最小のDummy Testを作成して実行できることを確認する。

#### 2-4 test.sh

Go TestおよびFrontend Testを
`test.sh` から一括実行できるようにする。

TDD中に繰り返し利用するため、
日常的に高速に実行できるTestを対象とする。

#### 2-5 build.sh / Build Gate

`build.sh` からTestを実行し、
Test成功時のみWails ApplicationをBuildする。

意図的にTestを失敗させ、
TestがREDの場合にBuildが実行されないことを確認する。

#### 2-6 Behavior / Logic Test 最小技術検証

B. Behavior Testを中心として、
Phase 5で採用するTest方式を最小単位で技術検証する。

必要に応じて、
C. Logic / Development Step Testについても
小さな実装の足場として利用できることを確認する。

#### 2-7 Integration Test 技術検証

Svelte / Wails / Go間の実際の接続を確認する
Integration Test方式を技術検証する。

自動化方式、実行時間、安定性を確認した上で、
`test.sh` およびRelease Build Gateへの組み込み方を決定する。

詳細なTest Caseおよび
実装途中で必要となるCは、
Step 3以降の各Behavior実装時に具体化する。

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
tasks.yaml
    ↓
File Persistence Adapter
    ↓
Repository
    ↓
ListTasks()
    ↓
Wails
    ↓
Svelte
    ↓
Task List
```

また、Task一覧からTaskを選択した場合は、
一覧表示時に取得したTaskをそのまま編集対象として使用せず、
`GetTask(ID)` により対象Taskを改めて取得する。

```text
Task選択
    ↓
GetTask(ID)
    ↓
Application Service
    ↓
Pick(ID)
    ↓
File Persistence Adapter
    ↓
tasks.yaml
    ↓
Task
    ↓
Svelte
    ↓
Task Edit
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
- Create画面初期表示
  - `title = ""`
  - `description = ""`
  - `priority = 3`
  - `status = TODO`
  - Correlation Validation Mode = INITIAL
  - Create Button = Disabled
- Focus Out時のSingle Field Validation
- Correlation Validation Mode = ACTIVE後の関連Correlation Validation再実行
- Create操作時の全Single Field Validation
- 全Single Field Validation成功後の全Correlation Validation
- Validation ErrorのField表示
- Validation成功時のConfirmation Modal表示
- Cancel時の入力内容維持
- `CreateTask(Properties)` によるTask生成
- Application ServiceでのTask Validation
- `Store(..., Deposit)` による永続化
- File Persistence AdapterによるULID生成
- File Persistence Adapterによる `createdAt / updatedAt` 設定
- Confirm時のTask登録
- Task YAML全量Write
- Write Failureの返却
- Error Toast
- 保存成功時のInfo Toast
- 保存成功後のForm Reset
  - 初期値へ戻す
  - Validation ErrorをClear
  - Correlation Validation Mode = INITIAL
  - Create Button = Disabled

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
- `GetTask(ID)` による編集対象Task取得
- Original Values / Current Valuesの管理
- Original ValuesとCurrent Valuesの比較によるChanged Fields管理
- Original Valuesへ戻したFieldをNot Changedとして扱う
- Changed Fieldsがない場合はModify ButtonをDisabledとする
- Focus Out時のSingle Field Validation
- Correlation Validation Mode = ACTIVE後の関連Correlation Validation再実行
- Modify操作時の全Single Field Validation
- 全Single Field Validation成功後の全Correlation Validation
- Validation ErrorのField表示
- Validation成功時のConfirmation Modal表示
- Cancel時のCurrent Values / Original Values / Changed Fields維持
- Changed FieldsのみをPropertiesとして `ModifyTask(ID, Properties)` へ渡す
- Application ServiceでのTask Validation
- `Store(..., Replace)` による永続化
- File Persistence Adapterによる `createdAt` 維持
- File Persistence Adapterによる `updatedAt` 更新
- Confirm時のTask更新
- Task YAML全量Write
- Not Foundの返却
- Read Failureの返却
- Write Failureの返却
- Error Toast
- 保存成功時のInfo Toast
- 保存成功後のEdit Baseline Reset
  - Current Valuesを新しいOriginal Valuesとする
  - Changed FieldsをClear
  - Validation ErrorをClear
  - Correlation Validation Mode = INITIAL
  - Modify Button = Disabled

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
- Focus Out時にSingle Field Validationを実行できる
- Correlation Validation Mode = INITIAL / ACTIVEに応じて相関Validationを適切なタイミングで実行できる
- Create / Modify Button押下時に全Single Field Validationおよび必要なCorrelation Validationを実行できる
- Validation Errorを対象Form Field付近へ適切な優先順位で表示できる
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
