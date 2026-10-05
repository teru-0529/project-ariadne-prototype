# Behavior Test Specification

## 1. Purpose

本ドキュメントは、Prototype Phase 5のBehaviorに対して、どのTest境界から何を確認することでBehaviorを保証するかを定義する。

Behaviorそのものの定義は `behavior-map.md` を正本とする。

本ドキュメントでは、Behaviorごとに必要となる以下のTestを整理する。

- Frontend Module Test
- Backend Module Test
- Integration Test

Test Codeそのものの詳細設計を目的とはせず、Behaviorに対して何をTestするかを明確にすることを目的とする。

### Test Environment

Frontend Module Testでは、Backendとの境界にTest Doubleを使用する。

Backend Module Testで `tasks.yaml` を使用する場合は、
`t.TempDir()` によって作成した一時Directory上に `tasks.yaml` を配置し、
実際のFile Adapter / File I/Oを使用する。

### Testにおける機能依存

TestのGivenを構築するために、先行するBehaviorで実装された機能を利用する場合がある。

この場合、対象BehaviorのTestは、依存する機能を保証するTestがGREENであることを前提とする。
ただし、先行BehaviorのExpectedを重複して検証することはせず、
対象BehaviorのGivenを構築するための操作として利用する。

Behavior間にこの依存がある場合は、各BehaviorのTest Specificationに依存先を明記する。

---

## 2. S01 Applicationを起動する

### Frontend Module Test

#### When

- Applicationを立ち上げる

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

## 3. S02 Configを選択する

### Frontend Module Test

#### Given

- Primary Navigationで `Task` が選択されている

#### When

- Primary Navigationの `Config` を選択する

#### Expected

- `Task` の選択状態が解除される
- `Config` が選択状態になる
- Config Dummy Pageが表示される
- Task Side Navigationが表示されない
- Main Areaに `Configuration is not available in this prototype.` が表示される

### Backend Module Test

N/A

### Integration Test

N/A

---

## 4. S03 Taskを選択する

### Frontend Module Test

#### Dependency

- S02 `Configを選択する`

#### Given

- Primary Navigationで `Config` が選択されている

#### When

- Primary Navigationの `Task` を選択する

#### Expected

- `Config` の選択状態が解除される
- `Task` が選択状態になる
- Task Pageが表示される
- Task Side Navigationが表示される
- Task Side Navigationに `List` が表示される
- Task Side Navigationに `Create` が表示される
- `List` / `Create` のどちらも選択されていない
- Main Areaに `Select an operation.` が表示される
- Task List取得処理は呼び出されない

### Backend Module Test

N/A

### Integration Test

N/A

---

## 5. L01 Task Listを表示する

### Result Pattern: Success - Tasks Exist

#### Frontend Module Test

##### Given

- Application起動直後で、Primary Navigationの `Task` が選択されている
- Backendから複数のTaskSummaryを取得できる状態である
- TaskSummaryは意図的に自然順ではない順序で取得できる状態とする
- 画面内に収まらない件数のTaskSummaryを取得できる状態とする

##### When

- Task Side Navigationの `List` を選択する

##### Expected

- Main AreaがTask List表示に切り替わる
- Backendから取得したすべてのTaskが表示される
- Backendから取得した順序でTaskが表示される

#### Backend Module Test

##### Given

- `tasks.yaml` に複数のTaskが保存されている
- Taskは意図的に自然順ではない順序で保存されている

##### When

- `ListTasks` を呼び出す

##### Expected

- 保存されているすべてのTaskがTaskSummaryとして返される
- `description` はTaskSummaryに含まれない
- `tasks.yaml` に保存されている順序でTaskSummaryが返される

#### Integration Test

対象とする。

### Result Pattern: Success - Empty

#### Frontend Module Test

##### Given

- Application起動直後で、Primary Navigationの `Task` が選択されている
- Backendから空のTaskSummary一覧を取得できる状態である

##### When

- Task Side Navigationの `List` を選択する

##### Expected

- Main AreaがTask List表示に切り替わる
- Taskは表示されない
- `No tasks.` が表示される
- Error Toastは表示されない

#### Backend Module Test

##### Given

- `tasks.yaml` のTaskが0件である

##### When

- `ListTasks` を呼び出す

##### Expected

- 空のTaskSummary一覧が返される
- Errorは返されない

#### Integration Test

対象とする。

### Result Pattern: Read Failure

#### Frontend Module Test

##### Given

- Application起動直後で、Primary Navigationの `Task` が選択されている
- BackendからRead FailureのPublic Errorが返る状態である

##### When

- Task Side Navigationの `List` を選択する

##### Expected

- Main Areaに `Select an operation.` が表示されたままである
- Task Side Navigationは`List` / `Create` のどちらも選択されていない
- Read FailureのError Toastが表示される

#### Backend Module Test

##### Given

以下のRead Failureが発生する状態をTestごとに構築する。

- File Read Failure
- YAML Parse Failure
- Invalid Format
- Unsupported Minor Version
- Invalid Task Data

Invalid Formatは、YAMLとしては成立するが、
ARIADNEが定義するYAML Formatの構造、項目または型に適合しない状態とする。

Invalid Task Dataは、YAML Formatには適合するが、
Task Validationを満たさない状態とする。

Invalid Task Dataは代表1ケースをBackend Module Testで確認する。
Task Validationの各Validation PatternはLogic Test（D）で保証する。

##### When

- `ListTasks` を呼び出す

##### Expected

- Read Failureが返される

#### Integration Test

対象とする。

### Result Pattern: Unsupported Major Version

#### Frontend Module Test

##### Given

- Application起動直後で、Primary Navigationの `Task` が選択されている
- BackendからUnsupported Major VersionのPublic Errorが返る状態である

##### When

- Task Side Navigationの `List` を選択する

##### Expected

- Main Areaに `Select an operation.` が表示されたままである
- Task Side Navigationは`List` / `Create` のどちらも選択されていない
- Unsupported Major VersionのError Toastが表示される

#### Backend Module Test

##### Given

- Applicationが対応しないMajor Versionを持つ `tasks.yaml` が保存されている
- Task本文は現在のMajor VersionのFormatでは解釈できない構造である

##### When

- `ListTasks` を呼び出す

##### Expected

- YAML本文のParse / Format Failureではなく、Unsupported Major Versionが返される

#### Integration Test

対象とする。

---

## 6. Note

CSS、Layout、Color、Spacing、Size等のVisual Designは実装対象とするが、自動Module Testによる保証対象とはしない。

Visual DesignはBuild後の実機確認によって確認し、必要に応じて調整する。
