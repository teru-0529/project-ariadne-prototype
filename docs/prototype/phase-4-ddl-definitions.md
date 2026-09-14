# Project ARIADNE — Prototype Phase 4：DDL YAML

**Status:** IN PROGRESS

---

## 1. 目的

Phase 4では、Project ARIADNEにおけるデータベース定義の正本となるDDL YAMLの設計、およびDDL生成・参照・実行検証の基本方針を定める。

DDL YAMLは、Phase 2で定義したTypes / Elementsを利用し、データベース上のSchema、Table、Column、Constraint、Index等の配置情報を定義する。

項目の意味・型・桁等をDDL YAMLで再定義せず、Types / Elementsを正本として参照する。

Phase 4 Prototypeでは、実Generator / Validatorそのものの実装は対象としない。

以下をPrototypeの対象とする。

- DDL YAMLの仕様設計
- Validation Ruleの定義
- Resolved DDL Modelの仕様設計
- PostgreSQL向けDDL生成規則の定義
- SQLite向け互換DDL生成規則の定義
- DDL Referenceの仕様設計
- Sample DDLの作成
- PostgreSQL / SQLiteでの実行検証
- ARIADNE Prototype自身が利用するSQLite Databaseの検証

実際のDDL Generator / Validatorの実装はCoreで行う。

---

## 2. 基本原則

### 2.1 Types / Elementsを項目定義の正本とする

Columnは原則としてElementsに登録されたElementを参照して定義する。

DDL YAML上でデータ型、桁数等を直接定義しない。

```text
types.yaml
    ↓
elements.yaml
    ↓
DDL YAML
    ↓
Resolved DDL Model
    ↓
Database DDL
```

DDL生成時には、参照されたElementおよびTypeからDB型・桁・値制約等を導出する。

---

### 2.2 DDL YAMLは「DB上での配置」を定義する

Elementsが「項目そのもの」を定義するのに対し、DDL YAMLでは、その項目をDatabase上でどのように配置するかを定義する。

DDL YAMLの主な責務は以下とする。

- Schema
- Table
- Column
- Primary Key
- Unique Constraint
- Foreign Key
- Index
- Column単位の初期値

Database Function、Stored Procedure、Trigger等のDBMS固有機能は、原則としてDDL YAMLによる構造化管理の対象外とする。

ただし、ARIADNEがBuiltInとして提供する機能については、自動生成対象とすることができる。

---

### 2.3 ColumnとElementを分離する

同一Elementを、一つのTable内または複数Table内で複数回利用できるものとする。

そのため、Column名とElement名は分離して管理する。

例：

```yaml
columns:
  customerId:
    element: customerId

  billingCustomerId:
    element: customerId

  shippingCustomerId:
    element: customerId
```

上記はいずれも同一の`customerId` Elementを利用するが、Database上ではそれぞれ別Columnとして生成される。

Column名はDDL YAML側で新たに定義可能とする。

---

### 2.4 PostgreSQLを基準Databaseとする

ARIADNE DDLの基準DatabaseはPostgreSQLとする。

DDL YAMLでは可能な限りDBMS固有表現を持たないが、機能および意味論の基準はPostgreSQLとする。

```text
DDL YAML
   ↓
Resolved DDL Model
   ↓
PostgreSQL Semantics
   ├─ PostgreSQL
   └─ SQLite Compatibility
```

SQLiteはPostgreSQLと同列の対応Databaseとはせず、ARIADNE Prototype自身のRuntime Databaseとして必要な範囲で互換対応する。

PostgreSQLとの機能差については、SQLite側で以下を許容する。

- 代替実装
- 簡略化
- 制約付き対応
- 非対応

SQLite対応のためにPostgreSQL側の表現力を制限することはしない。

---

## 3. Schema

DDL定義はSchema単位で管理可能とする。

SchemaごとにTableを登録し、複数Schemaを一つのProject内で扱えるものとする。

SchemaをまたぐForeign KeyはPrototype Phase 4では許可しない。

SQLiteにはPostgreSQLと同等のSchema概念が存在しないため、SQLiteでの扱いは互換生成規則として別途定義する。

---

## 4. Table

DDL YAMLではTable単位の定義を行う。

Tableには以下を定義可能とする。

- Table名
- Column
- Primary Key
- Unique Constraint
- Foreign Key
- Index

TableおよびColumnの物理名はProject ARIADNEのNaming Ruleに従い生成する。

生成SQLをTable単位でファイル分割するかどうかはDDLモデルの仕様には含めない。

生成DDLのファイル分割単位および実行順序は、Generatorの出力仕様としてPhase 4で検討する。

---

## 5. Column

ColumnはElementsに登録されたElementを参照して定義する。

DDL YAMLでは原則として型・桁を直接指定しない。

以下の情報はTypes / Elementsから導出する。

- DBデータ型
- 最大長
- Precision
- Scale
- Minimum
- Maximum
- Enum等の値制約

### 5.1 Element由来の制約

ElementにMinimum / Maximum等が定義されている場合、Databaseで表現可能なものはConstraintとして生成する。

例：

```yaml
minimum: 0
```

から、PostgreSQLでは概念的に以下を生成する。

```sql
CHECK (quantity >= 0)
```

具体的な変換規則はPhase 4で決定する。

---

### 5.2 Table上で定義する情報

以下の情報は、そのColumnがTable内でどのように利用されるかに依存するため、DDL YAML側で定義する。

- NOT NULL
- Default
- Primary Key
- Unique Constraint
- Foreign Key
- Index

---

### 5.3 Default

Columnの初期値をDDL YAMLで定義可能とする。

例：

```yaml
columns:
  productType:
    element: productType
    default: NORMAL
```

同一Elementを利用するColumnであっても、TableまたはColumnごとに異なるDefaultを定義可能とする。

Prototypeではまず固定値を対象とする。

CURRENT_TIMESTAMP等のDatabase式を一般的なDefaultとして許可するかについてはPhase 4内で検討する。

ARIADNE BuiltIn Columnについては、一般Defaultとは別にARIADNEが生成規則を管理する。

---

## 6. Column名の変更

Element名とは異なるColumn名を定義可能とする。

これにより、一つのElementを同一Table内で複数回利用できる。

例：

```yaml
billingCustomerId:
  element: customerId
```

この場合、

```text
Element       : customerId
Column論理名  : billingCustomerId
Column物理名  : billing_customer_id
```

として扱う。

生成DDLには、Columnの説明と参照元Elementを追跡可能な情報をCommentとして出力する。

Commentの内容はElementsのdescription等から生成し、DDL YAML内への重複記載は原則行わない。

Commentの具体的な生成形式はPhase 4内で決定する。

---

## 7. BuiltIn Columns / Audit Columns

### 7.1 createdAt / updatedAt

すべてのTableに以下をARIADNE BuiltIn Columnとして自動付与する。

- createdAt
- updatedAt

各TableのDDL YAMLへ個別に記載する必要はない。

#### createdAt

Record作成日時を保持する。

INSERT時に現在日時を自動設定する。

#### updatedAt

Record更新日時を保持する。

INSERT時に現在日時を自動設定する。

UPDATE時にはARIADNEが生成するBuiltIn Function / Trigger等により現在日時へ自動更新する。

PostgreSQLでの具体的な実現方式はPhase 4で決定する。

SQLiteでは必要に応じて代替方式を採用する。

---

### 7.2 Audit Trace Element

作成・更新を追跡するための識別情報は、ARIADNE固定の型・桁とはしない。

Projectごとに、追跡情報として利用するElementを指定可能とする。

例：

```text
Audit Trace Element = traceId
```

`traceId`の型・桁・制約等はTypes / Elementsから導出する。

指定されたElementを利用して、作成時・更新時の追跡Columnを自動付与する。

概念例：

```text
createdAt
createdTraceId  → Element: traceId

updatedAt
updatedTraceId  → Element: traceId
```

作成用・更新用Columnは同一Elementを参照するが、Column名を分離する。

具体的なColumn名およびDDL YAML上の指定形式はPhase 4で決定する。

Audit Trace Elementが指定されていない場合は、`createdAt` / `updatedAt`のみをBuiltIn Columnとして付与する。

---

### 7.3 Audit Trace値の設定

Database自身は、アプリケーション上の利用者やTrace IDを直接知ることができない。

そのため、Audit Trace値をDatabaseへどのように伝搬するかについてはPhase 4で方式を整理する。

PostgreSQLでは、Transaction単位のSession Parameter等を利用し、Trigger / Functionから取得する方式も候補とする。

例：

```sql
SET LOCAL app.trace_id = '...';
```

ただし、接続PoolやTransaction管理との関係があるため、Prototype Phase 4では仕様整理と実証までを対象とし、汎用実装はCoreで行う。

---

## 8. Primary Key

Primary KeyはTable側の責務としてDDL YAMLに定義する。

Elementsの`identifier: true`はAPI上の識別子として利用可能であることを示すものであり、Database上のPrimary Keyを意味しない。

単一Primary Keyおよび複合Primary Keyを許可する。

Primary KeyのConstraint名はARIADNEがNaming Ruleに従って自動生成する。

DDL YAMLでは名称を指定しない。

---

## 9. Unique Constraint

データモデル上、一意でなければならないColumnまたはColumn組合せをUnique Constraintとして定義する。

単一Columnおよび複合ColumnのUnique Constraintを許可する。

例：

```text
orderNo must be unique
```

のように、データそのものの成立条件としての一意性を表現する。

Unique Constraint名はARIADNEがNaming Ruleに従って自動生成する。

DDL YAMLでは名称を指定しない。

### 9.1 Unique Index

Prototype Phase 4ではUnique IndexをDDL YAMLの独立した定義対象とはしない。

一意性はUnique Constraintとして表現する。

PostgreSQL固有のPartial Unique Index、Expression Unique Index等が必要な場合はCustom SQLを利用する。

---

## 10. Foreign Key

Foreign KeyをDDL YAMLで管理可能とする。

Foreign Keyでは以下を定義する。

- 参照元Column
- 参照先Table
- 参照先Column

単一Foreign Keyおよび複合Foreign Keyを許可する。

Foreign Keyで対応する参照元Columnと参照先Columnは、同一Elementを参照していることを必須とする。

複合Foreign Keyでは、対応する各ColumnについてElementの一致をValidationする。

SchemaをまたぐForeign KeyはPrototype Phase 4では許可しない。

Foreign Key名はARIADNEがNaming Ruleに従って自動生成する。

DDL YAMLでは名称を指定しない。

Table生成後にForeign Keyを適用する等、依存関係を考慮したDDL生成順序はGenerator側で管理する。

ファイル名やファイル分割方法はDDL YAMLの概念モデルには含めない。

---

## 11. Index

検索性能等を目的として、通常IndexをDDL YAMLで定義可能とする。

以下を対象とする。

- 単一Column Index
- 複合Column Index

Unique Indexは独立した定義対象としない。

Index名はARIADNEがNaming Ruleに従って自動生成する。

DDL YAMLでは名称を指定しない。

---

### 11.1 Unique Constraintとの重複

Unique Constraintにより一意性を保証するためのIndexがDatabase側で生成されることを前提とする。

そのため、Unique Constraintと完全に同一のColumn・同一順序を持つ通常Indexは定義不可とする。

例：

```text
UNIQUE (customerId, orderNo)
INDEX  (customerId, orderNo)
```

上記はValidation Errorとする。

一方、Column順序が異なる場合は検索特性が異なるため許可する。

```text
UNIQUE (customerId, orderNo)
INDEX  (orderNo, customerId)
```

また、Unique Constraintの一部Columnのみを対象とするIndexについては許可する。

例：

```text
UNIQUE (customerId, orderNo)
INDEX  (customerId)
```

ただし、Unique Constraintによって作成されるIndexで代替可能な場合があるため、Validation Warningの対象とすることを検討する。

---

## 12. Constraint / Index Naming

以下の名称はARIADNEがNaming Ruleに従って自動生成する。

- Primary Key
- Unique Constraint
- Foreign Key
- Index

利用者がDDL YAML上で名称を指定することは原則として行わない。

具体的なNaming RuleはPhase 4内で決定する。

---

## 13. BuiltIn Database Behavior

ARIADNEでは、各Projectが毎回個別に定義する必要のない共通Database BehaviorをBuiltInとして提供可能とする。

Prototype Phase 4では少なくとも以下を対象とする。

### 13.1 createdAt自動設定

INSERT時に現在日時を設定する。

### 13.2 updatedAt自動設定

INSERT時に現在日時を設定する。

### 13.3 updatedAt自動更新

UPDATE時に`updatedAt`を現在日時へ変更する。

PostgreSQLでは共通FunctionおよびTableごとのTriggerを自動生成する方式を候補とする。

SQLiteでは必要に応じて代替方式を利用する。

利用者がこれらのBuiltIn Function / TriggerをDDL YAMLへ直接記述する必要はない。

---

## 14. Custom SQL

Database Function、Stored Procedure、Trigger等のDBMS固有機能は、原則としてDDL YAMLによる構造化管理の対象外とする。

一方、実システムではこれらが必要になる場合があるため、任意のSQLを追加登録可能とする。

例：

- Database Function
- Stored Procedure
- Trigger
- Event Trigger
- Extension設定
- DBMS固有DDL
- 初期化処理

Custom SQLはDBMS固有であることを許容する。

Custom SQLとARIADNE生成DDLとの実行順序を指定または管理できる方式をPhase 4で検討する。

ファイル名による実行順制御はARIADNEの概念モデルには含めない。

---

## 15. Generated ID等のDatabase Function

業務上の採番値等をDatabase Functionで生成する方式を許容する。

例：

```text
ORD_20260914_002
```

このような値の生成は、SEQUENCE_ID等のARIADNE標準型とは異なる業務固有ルールとなり得るため、Prototype Phase 4では原則としてCustom SQLによる実装対象とする。

代表的な利用例：

- Prefix + 連番
- 日付 + 連番
- 年月 + 連番
- Schema単位の採番
- Table単位の採番
- 業務単位の採番

UUID / ULID等の標準的なID生成については、Types / Elements側との責務も含め別途整理する。

Phase 4では業務固有採番を汎用YAMLモデルとして抽象化すること自体を目的としない。

---

## 16. Resolved DDL Model

SourceとなるDDL YAMLから、DDL生成に必要な情報をすべて解決したResolved DDL Modelを定義する。

```text
Types
   ↓
Elements
   ↓
DDL YAML
   ↓
Resolved DDL Model
   ├─ PostgreSQL DDL
   ├─ SQLite Compatibility DDL
   └─ DDL Reference
```

Resolved DDL Modelでは、少なくとも以下を解決済みとする。

- Schema
- Table
- Column
- Element
- Column Alias
- 論理型
- 桁
- Precision / Scale
- Minimum / Maximum
- Element由来Constraint
- NULL制約
- Default
- Primary Key
- Unique Constraint
- Foreign Key
- Index
- BuiltIn Columns
- Audit Trace Columns
- Comment生成情報

DBMS固有SQLへの変換はResolved後の処理とする。

Prototype Phase 4ではResolved DDL Modelの仕様およびSampleを定義し、Generator実装は行わない。

---

## 17. PostgreSQL DDL

PostgreSQLをARIADNE DDLの基準Databaseとする。

Phase 4では以下の生成規則を定義する。

- Schema
- Table
- Column
- Data Type
- Default
- Check Constraint
- Primary Key
- Unique Constraint
- Foreign Key
- Index
- Comment
- BuiltIn Function
- BuiltIn Trigger

実Generatorの実装はCoreで行う。

Phase 4ではSample SQLを作成し、実際のPostgreSQLへ適用可能であることを検証する。

---

## 18. SQLite Compatibility DDL

SQLiteはARIADNE Prototype自身が利用するRuntime Databaseとして必要な範囲で対応する。

PostgreSQLを基準とし、SQLiteで表現できない機能については以下を許容する。

- 代替
- 簡略化
- 制限
- 非対応

Phase 4では、少なくともTask ManagementをSQLite上で実行可能とするために必要な互換規則を整理する。

SQLiteの制約に合わせて、ARIADNE DDLまたはPostgreSQL向け仕様そのものを弱めることはしない。

---

## 19. DDL Output

生成DDLのファイル構成はARIADNEの論理モデルとは分離する。

Phase 4では以下を検討する。

- 1ファイルへ集約する方式
- Schema単位
- Table単位
- DDL種別単位
- 上記の組合せ

重要なのはファイル単位ではなく、依存関係を満たす実行順序である。

概念的な実行Phaseは以下を想定する。

```text
Schema
  ↓
Table
  ↓
Constraint / Index
  ↓
Foreign Key
  ↓
BuiltIn / Custom SQL
```

具体的な順序は、BuiltIn FunctionやCustom SQLの依存関係を含めPhase 4で決定する。

ファイル名の数字等による順序制御は、必要であればGeneratorの出力仕様として採用できるが、DDL YAMLの概念モデルには含めない。

---

## 20. DDL Reference

作成されたDatabase定義をWeb Browserから参照可能とする。

DDL Referenceでは、少なくとも以下を確認可能とする。

- Schema
- Table
- Column
- Element
- Data Type
- Default
- NULL制約
- Element由来Constraint
- Primary Key
- Unique Constraint
- Foreign Key
- Index
- BuiltIn Column
- Audit Trace Column
- Comment
- 生成想定SQL

Phase 3におけるRedocと同様に、生成成果物を人間が容易に確認できることを目的とする。

適切な既存Viewerが存在する場合は利用を検討する。

適切なものがない場合、MarkdownからHTMLを生成する方式を採用する。

---

## 21. Runtime

Phase 4では、設計したDDLが実Database上で成立することを検証する。

### 21.1 PostgreSQL

PostgreSQLをDDL設計の基準Databaseとして利用する。

DockerベースでPostgreSQLを起動する。

Database参照用としてpgwebを併せて起動し、BrowserからDatabaseの状態を確認可能とする。

Order ManagementをPostgreSQLでの主要検証Modelとして利用する。

### 21.2 SQLite

ARIADNE Prototype自身が利用するDatabaseとしてSQLiteを使用する。

Runtime上にSQLite Database Fileを配置し、ARIADNE Prototypeから利用可能な状態とする。

Task ManagementをARIADNE Prototype自身が利用するDatabase Modelとして使用する。

DockerはARIADNE Prototypeアプリケーション自体の実行要件とはしない。

---

## 22. Sample / Verification

### 22.1 Task Management

主用途：

- ARIADNE Prototype自身が利用するSQLite Database
- SQLite互換DDLの実証
- Runtime SQLite Fileの検証

必要に応じて同一ModelからPostgreSQL向けDDLも作成し、互換性を確認する。

### 22.2 Order Management

主用途：

- ARIADNEが設計対象とする業務システムのDatabase Model
- PostgreSQL DDL仕様の実証
- Schema / FK / Index / Constraint / Function等を含む実践的検証

---

## 23. Validation Rule

Prototype Phase 4ではValidatorそのものは実装しない。

将来のValidatorが検証すべきRuleを定義する。

少なくとも以下を対象とする。

### Element / Column

- 参照Elementが存在すること
- Schema内でTable名が重複しないこと
- Table内でColumn名が重複しないこと
- DefaultがElementの型・制約と矛盾しないこと

### Primary Key / Unique Constraint

- 対象Columnが存在すること
- 複合定義内でColumnが重複しないこと

### Foreign Key

- 参照元Columnが存在すること
- 参照先Schema / Table / Columnが存在すること
- Schemaをまたいでいないこと
- 参照元Columnと参照先Columnが同一Elementであること
- 複合FKのColumn数が一致すること
- 複合FKの対応する各ColumnのElementが一致すること

### Index

- Index対象Columnが存在すること
- 同一Index内でColumnが重複しないこと
- 同一Column・同一順序のIndexが重複しないこと
- Unique Constraintと完全に同一Column・同一順序のIndexを定義していないこと

Unique Constraintの先頭Columnと重複するIndex等についてはWarning候補とする。

### Naming

- 自動生成されるConstraint / Index名が衝突しないこと

Validation RuleはPhase 4で追加・精緻化する。

---

## 24. Phase 4で決定する事項

Phase 4では、以下を確定する。

- DDL YAMLのファイル構成
- Schema定義形式
- Table定義形式
- Column定義形式
- Default定義形式
- Primary Key定義形式
- Unique Constraint定義形式
- Foreign Key定義形式
- Index定義形式
- Constraint / Index Naming Rule
- Column Comment生成規則
- Minimum / Maximum等からのConstraint変換規則
- createdAt / updatedAtの詳細仕様
- Audit Trace Elementの指定方式
- Audit Trace Column名
- Audit Trace値のDatabaseへの伝搬方式
- PostgreSQL型変換規則
- SQLite互換型変換規則
- SQLiteにおけるSchemaの扱い
- BuiltIn SQLの仕様
- Custom SQLの登録方法
- Custom SQLと生成DDLの実行順序
- DDL出力ファイル構成
- DDL Reference生成方式
- Validation Rule
- Runtime構成

---

## Prototype Phase 4：作業チェックリスト

### Step 1：DDL基本仕様

- [ ] DDL YAML全体構造を確定する
- [ ] Schema定義形式を確定する
- [ ] Table定義形式を確定する
- [ ] Column定義形式を確定する
- [ ] Element参照方式を確定する
- [ ] Column別名定義方式を確定する
- [ ] Default定義方式を確定する
- [ ] Element由来Constraintの導出規則を確定する

### Step 2：BuiltIn / Audit

- [ ] createdAtの仕様を確定する
- [ ] updatedAtの仕様を確定する
- [ ] Audit Trace Elementの指定方式を確定する
- [ ] created側Audit Trace Columnの仕様を確定する
- [ ] updated側Audit Trace Columnの仕様を確定する
- [ ] Audit Trace値の伝搬方式を検証する

### Step 3：Constraint / Index

- [ ] Primary Key定義形式を確定する
- [ ] Unique Constraint定義形式を確定する
- [ ] Foreign Key定義形式を確定する
- [ ] 複合Foreign Key定義形式を確定する
- [ ] Index定義形式を確定する
- [ ] Unique Constraint / Index重複Ruleを確定する
- [ ] Constraint / Index Naming Ruleを確定する

### Step 4：Validation Rule

- [ ] Element / Column Validation Ruleを確定する
- [ ] Default Validation Ruleを確定する
- [ ] PK / Unique Validation Ruleを確定する
- [ ] FK Validation Ruleを確定する
- [ ] FK Element一致Ruleを確定する
- [ ] 複合FK Validation Ruleを確定する
- [ ] Index Validation Ruleを確定する
- [ ] Naming Validation Ruleを確定する

### Step 5：Resolved DDL Model

- [ ] Resolved DDL Modelの形式を確定する
- [ ] Types / Elementsの解決結果をSample化する
- [ ] BuiltIn Columnsの展開結果をSample化する
- [ ] Audit Trace Columnsの展開結果をSample化する
- [ ] Constraint / Indexの解決結果をSample化する
- [ ] Comment情報の解決結果をSample化する

### Step 6：PostgreSQL DDL仕様

- [ ] PostgreSQL型変換規則を確定する
- [ ] Schema SQLを定義する
- [ ] Table SQLを定義する
- [ ] Default / Check Constraint変換を定義する
- [ ] Column Comment生成規則を定義する
- [ ] PK / Unique Constraint生成規則を定義する
- [ ] Index生成規則を定義する
- [ ] FK生成規則を定義する
- [ ] BuiltIn updatedAt Function / Triggerを定義する
- [ ] Sample PostgreSQL DDLを作成する
- [ ] Sample DDLをPostgreSQLで実行確認する

### Step 7：SQLite互換仕様

- [ ] SQLite型変換規則を確定する
- [ ] Schemaの扱いを確定する
- [ ] PostgreSQLとの差異を整理する
- [ ] BuiltIn updatedAtの代替方式を確定する
- [ ] Task Management用SQLite DDLを作成する
- [ ] SQLiteで実行確認する

### Step 8：Custom SQL

- [ ] Custom SQLの登録方式を確定する
- [ ] DBMS別Custom SQLの扱いを確定する
- [ ] Custom SQLと生成DDLの実行順序を確定する
- [ ] GENERATED_ID等の業務固有FunctionをSampleとして整理する

### Step 9：DDL Output

- [ ] DDLのファイル分割方針を検討する
- [ ] DDLの実行Phaseを確定する
- [ ] ファイル名による順序制御を利用するか決定する
- [ ] Sample出力構成を作成する

### Step 10：DDL Reference

- [ ] DDL Referenceの表示内容を確定する
- [ ] Schema / Table / Columnの表示方式を確定する
- [ ] Constraint / Indexの表示方式を確定する
- [ ] Elementとの対応表示を確定する
- [ ] SQLの参照方式を確定する
- [ ] HTML生成方式を確定する
- [ ] Browserから参照できるSampleを作成する

### Step 11：PostgreSQL Runtime

- [ ] PostgreSQL Docker環境を作成する
- [ ] pgwebを追加する
- [ ] Order Management Sample DDLを適用する
- [ ] Table / Constraint / Indexを確認する
- [ ] BuiltIn Triggerを確認する
- [ ] Custom SQLを確認する
- [ ] pgwebからDatabaseを確認する

### Step 12：SQLite Runtime

- [ ] SQLite Database Fileの配置方式を確定する
- [ ] Task Management DDLを適用する
- [ ] Runtime上にSQLite Fileを生成する
- [ ] ARIADNE Prototypeから接続可能な状態を確認する
- [ ] createdAt / updatedAt等のBuiltIn動作を確認する

### Step 13：Phase 4 Completion

- [ ] Sample DDL YAMLを整理する
- [ ] Resolved DDL Sampleを整理する
- [ ] PostgreSQL Sample SQLを整理する
- [ ] SQLite Sample SQLを整理する
- [ ] Runtime検証を完了する
- [ ] Phase 4ドキュメントを実証結果に合わせて更新する
- [ ] Coreへ引き継ぐGenerator / Validator仕様を整理する
- [ ] 未決事項が残っていないことを確認する
- [ ] Phase 4 COMPLETE

---

## Prototype Phase 4: DDL YAML — IN PROGRESS
