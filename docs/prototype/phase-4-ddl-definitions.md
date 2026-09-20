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

## 3. Service / Schema

DDL SourceはServiceに属するDatabase Schema定義として管理する。

ServiceはAPI / DDL等の設計成果物を束ねるProject ARIADNE上の上位定義単位とする。

1つのDDL Source YAMLには1つのServiceと1つのSchemaを定義し、以下に配置する。

```text
src/
└─ ddl/
   └─ schemas/
      ├─ order-management.yaml
      └─ task-management.yaml
```

DDL Source YAMLのファイル名はService IDと一致させる。

DDL Source YAMLでは、ServiceとSchemaを以下の形式で定義する。

```yaml
service:
  id: order-management

schema: received_order
```

`service.id` はARIADNE上のService識別子とし、必須とする。

Service IDは `kebab-case` とする。

`schema` はARIADNE上のSchema識別子とし、必須とする。

Schema識別子は `snake_case` とする。

Service IDとSchema識別子は、それぞれ異なる責務を持ち、一致する必要はない。

例えば `order-management` Serviceでは、Database Schemaとして `received_order` を利用する。

DDL Source YAMLには物理Schema名を別途保持せず、Database上の物理Schema名はNaming Ruleに従って生成する。

複数Service / Schemaを一つのProject内で扱えるものとする。

SchemaをまたぐForeign KeyはPrototype Phase 4では許可しない。

SQLiteにはPostgreSQLと同等のSchema概念が存在しないため、SQLiteでの扱いは互換生成規則として別途定義する。

---

## 4. Table

DDL Source YAMLでは、Schema配下のTableをMapとして定義する。

基本形式を以下とする。

```yaml
service:
  id: order-management

schema: received_order

tables:
  orders:
    name: 受注
    columns: ...
```

TableのMap KeyはARIADNE上のTable識別子とする。

Table識別子は `snake_case` とし、原則として複数形の名詞を使用する。

例：

```text
orders
order_details
shipping_instructions
cancel_instructions
tasks
statuses
```

Tableには以下を定義する。

- Table識別子
- `name`
- Column
- Primary Key
- Unique Constraint
- Foreign Key
- Index

`name` はTableの論理名・表示名とし、必須とする。

例：

```yaml
tables:
  orders:
    name: 受注

  order_details:
    name: 受注明細
```

`name` はDDL Reference等での表示、およびDatabase Table Commentの生成元として利用する。

DDL Source YAMLには物理Table名を別途保持しない。

Database上の物理Table名はTable識別子からNaming Ruleに従って生成する。

生成SQLをTable単位でファイル分割するかどうかはDDLモデルの仕様には含めない。

生成DDLのファイル分割単位および実行順序は、Generatorの出力仕様としてPhase 4で検討する。

---

## 5. Column

ColumnはElementsに登録されたElementを参照して定義する。

基本形式を以下とする。

```yaml
columns:
  orderQuantity:
    element: quantity
    name: 受注数
    notNull: true

  shippingQuantity:
    element: quantity
    name: 出荷数
    notNull: true
    default: 0
```

ColumnのMap KeyはARIADNE上のColumn識別子とする。

Column識別子は `lowerCamelCase` とする。

Database上の物理Column名はNaming Ruleに従って `snake_case` へ変換する。

例：

```text
orderQuantity
    ↓
order_quantity
```

ColumnとElementは別の識別子として扱う。

同一Elementを、一つのTable内または複数Table内で複数回利用できるものとする。

例：

```yaml
orderQuantity:
  element: quantity
  name: 受注数

shippingQuantity:
  element: quantity
  name: 出荷数

cancelQuantity:
  element: quantity
  name: キャンセル数
```

上記はすべて同一の `quantity` Elementを参照するが、Database上ではそれぞれ別Columnとして生成される。

### 5.1 Columnで定義可能な情報

Columnでは以下を定義可能とする。

- `element`
- `name`
- `notNull`
- `default`
- `sequence`
- `constraints`

`element` は必須とする。

`name` はColumn固有の論理名・表示名とし、Optionalとする。

`name` が指定されていない場合は、参照Elementの `name` を利用する。

同一Elementを異なる業務上の意味で利用する場合、Column側の `name` により表示名をOverrideできる。

例：

```yaml
personInCharge:
  element: userId
  name: 出荷担当者ID
```

### 5.2 NOT NULL

ColumnごとのNULL許容性はDDL Source YAMLで定義する。

```yaml
customerId:
  element: customerId
  notNull: true
```

`notNull: true` の場合、Database DDLに `NOT NULL` を生成する。

`notNull` が指定されていない場合はNULL許容とする。

API上の `required` とDatabase上の `NOT NULL` は別の責務として扱い、自動的な相互導出は行わない。

### 5.3 Element由来情報のOverride禁止

DDL Source YAMLでは、参照ElementまたはTypeが管理する以下の情報を再定義またはOverrideしてはならない。

- Type
- length
- minLength
- maxLength
- regex
- minimum
- maximum
- precision
- scale
- ENUM values
- format

これらはTypes / Elementsを正本とし、DDL生成時に自動的に導出する。

ColumnはElement / Typeを再定義またはOverrideする場所ではない。

ただし、そのElementを特定のColumnとして利用する際に、Element由来Constraintよりも強い追加Constraintを課すことはできる。

Column固有ConstraintはElement定義そのものを変更するものではなく、そのColumnにおける追加条件として扱う。

`sequence` はElement / Typeの値定義をOverrideするものではなく、
そのColumn自身がDatabase上の採番主体となるかを指定するColumn固有の配置・利用情報として扱う。

### 5.4 Element由来Constraint

Types / Elementsに定義された制約は、対象Databaseで表現可能な限りDatabase DDLへ反映する。

Database型そのもので保証できる制約は型として反映し、型のみでは保証できない制約はCHECK Constraint等として生成する。

同一の制約をDatabase型とCHECK Constraintで重複して生成することは原則として行わない。

PostgreSQLでは、Prototype Phase 4において以下を基本変換規則とする。

| Type / Constraint     | PostgreSQLへの反映                               |
| --------------------- | ------------------------------------------------ |
| `FIXED_STRING.length` | `varchar(n)` + `CHECK (LENGTH(column) = n)`      |
| `maxLength`           | `varchar(n)` 等の型として反映                    |
| `minLength`           | `CHECK (LENGTH(column) >= n)`                    |
| `precision / scale`   | `numeric(p,s)`                                   |
| `minimum`             | `CHECK (column >= value)`                        |
| `maximum`             | `CHECK (column <= value)`                        |
| `regex`               | PostgreSQL正規表現による `CHECK`                 |
| `ENUM`                | PostgreSQL ENUM型                                |
| `CODE`                | 値集合に対するConstraintは生成しない             |
| `format`              | 一般的なCHECK Constraintの自動生成対象とはしない |

`FIXED_STRING` は `char(n)` を使用しない。

固定長は `varchar(n)` と `CHECK (LENGTH(column) = n)` の組合せによって保証する。

例：

```sql
customer_id varchar(6)
  CHECK (LENGTH(customer_id) = 6)
```

`maxLength` は `varchar(n)` 等のDatabase型自体によって保証できるため、同一内容の追加CHECK Constraintは生成しない。

`minLength` は型のみでは保証できないためCHECK Constraintとして生成する。

例：

```sql
CHECK (LENGTH(description) >= 10)
```

`regex` はPostgreSQLの正規表現演算子を利用したCHECK Constraintとして生成する。

概念例：

```sql
CHECK (billing_id ~ '^BL-[0-9]{7}$')
```

`precision / scale` はPostgreSQLの `numeric(p,s)` として型に反映する。

例：

```text
precision: 5
scale: 2
```

は概念的に以下へ変換する。

```sql
numeric(5,2)
```

`minimum / maximum` は値域を保証するCHECK Constraintとして生成する。

例：

```yaml
minimum: 0
maximum: 100
```

は概念的に以下へ変換する。

```sql
CHECK (profit_rate >= 0)
CHECK (profit_rate <= 100)
```

ENUMは `varchar + CHECK` ではなく、PostgreSQLのENUM型として生成する。

CODEは値集合が運用中に変化し得るため、Element定義から値集合Constraintを生成しない。

`format` は一般的なCHECK Constraintの自動生成対象とはしない。Database型への変換に影響する場合は、PostgreSQL型変換規則側で扱う。

Element由来ConstraintをDDL Source YAML側からOverride、緩和、無効化することはできない。

### 5.5 Column固有Constraint

Columnには、参照Element由来Constraintに加えて、そのColumn固有の追加Constraintを定義可能とする。

基本形式を以下とする。

```yaml
orderQuantity:
  element: quantity
  notNull: true
  constraints:
    minimum: 1
  name: 受注数
```

Column固有ConstraintはElement由来ConstraintをOverrideするものではない。

Element由来ConstraintとColumn固有Constraintの双方を満たす値を、そのColumnの有効な値域とする。

概念的には以下として扱う。

```text
Effective Constraint
  = Element Constraint
    ∩ Column Constraint
```

例えば、Element `quantity` に以下が定義されている場合、

```yaml
minimum: 0
```

Column側で、

```yaml
constraints:
  minimum: 1
```

を指定することは許可する。

この場合、そのColumnの実効Constraintは `minimum: 1` となる。

一方、Column固有ConstraintによってElement由来Constraintを緩和することはできない。

例えば、Elementが `minimum: 0` の場合、Columnに以下を指定してはならない。

```yaml
constraints:
  minimum: -1
```

また、Element由来Constraintと同一内容のConstraintをColumn側へ重複して指定してはならない。

例えば、Elementが `minimum: 0` の場合、以下もValidation Errorとする。

```yaml
constraints:
  minimum: 0
```

Column固有Constraintは、そのColumnに追加の制約を課す必要がある場合にのみ利用する。

Prototype Phase 4では、少なくとも以下をColumn固有Constraintとして利用可能とする。

- `minLength`
- `maxLength`
- `minimum`
- `maximum`
- `regex`

比較可能なConstraintについては、Element由来Constraintよりも厳しい値のみ指定可能とする。

`regex` については一般的な包含関係の解析を行わない。

Element由来Regexと異なるRegexがColumn側に指定された場合は、双方を満たす追加Constraintとして扱う。

Element由来Regexと完全に同一のRegexをColumn側へ指定した場合は、意味のない重複としてValidation Errorとする。

Column固有ConstraintはElement由来のDatabase型を変更しない。

例えばElementの `maxLength: 100` がPostgreSQLで `varchar(100)` として表現され、Column側に以下が指定された場合、

```yaml
constraints:
  maxLength: 50
```

PostgreSQL型を `varchar(50)` へ変更せず、概念的に以下として生成する。

```sql
description varchar(100),
CHECK (LENGTH(description) <= 50)
```

Database型はElementの物理表現とし、Column固有Constraintはその利用箇所に対する追加制約として表現する。

比較可能なElement由来ConstraintとColumn固有Constraintが存在する場合、
Generatorは実効的に最も強いConstraintのみを生成し、意味的に冗長なCHECK Constraintは生成しない。

例えばElementが `minimum: 0`、Columnが `minimum: 1` の場合、

```sql
CHECK (order_quantity >= 1)
```

のみを生成し、以下のような重複生成は行わない。

```sql
CHECK (order_quantity >= 0)
CHECK (order_quantity >= 1)
```

複数Column間の関係等、単一ColumnのConstraintとして表現できない業務制約はColumn固有Constraintの対象とはせず、Custom SQLを利用する。

### 5.6 Default

Columnの初期値をDDL Source YAMLで定義可能とする。

Defaultには以下の2種類を許可する。

1. Literal
2. ARIADNE BuiltIn Expression

Literalは値を直接指定する。

例：

```yaml
status:
  element: orderStatus
  notNull: true
  default: PREPARING
```

```yaml
shippingQuantity:
  element: quantity
  notNull: true
  default: 0
```

BuiltIn Expressionは以下の形式で指定する。

```yaml
orderDate:
  element: orderDate
  notNull: true
  default:
    expression: CURRENT_DATE
```

Prototype Phase 4では、少なくとも以下のBuiltIn Expressionを定義する。

```text
CURRENT_DATE
CURRENT_TIME
CURRENT_DATETIME
```

`expression` に任意のDatabase SQLを直接記述することは許可しない。

BuiltIn ExpressionはARIADNE上のDBMS非依存な意味として扱い、各Database向けGeneratorが対応するDatabase表現へ変換する。

DefaultのLiteralおよびBuiltIn Expressionは、参照ElementのTypeおよびConstraintと整合していなければならない。

同一Elementを利用するColumnであっても、Columnごとに異なるDefaultを定義可能とする。

Database Function等による業務ロジックとDefaultは別の責務として扱う。

DefaultはINSERT時の初期値を定義するものであり、業務ロジックによる値生成・更新はFunction / Trigger / Custom SQL等の責務とする。

ARIADNE BuiltIn Columnについては、一般ColumnのDefaultとは別にARIADNEが生成規則を管理する。

### 5.7 Sequence

`SEQUENCE_ID` Elementを参照するColumnでは、そのColumn自身をDatabase上の連番採番主体とする場合、`sequence: true` を指定する。

例：

```yaml
orderShipmentId:
  element: orderShipmentId
  notNull: true
  sequence: true
```

`sequence: true` は、そのColumnの値をDatabaseの連番採番機構によって生成することを表す。

`sequence` はPrimary Key、Unique Constraint、またはARIADNEの `identifier` を意味しない。

`SEQUENCE_ID` は値としてのDatabase型を定義し、`sequence` はそのColumnが採番主体であるかを定義する。両者は別の責務として扱う。

同一の `SEQUENCE_ID` Elementを、採番元ColumnとForeign Key等の従属Columnの双方で利用可能とする。

例：

```yaml
# 採番元
orderShipmentId:
  element: orderShipmentId
  notNull: true
  sequence: true
```

```yaml
# Foreign Key等の従属Column
orderShipmentId:
  element: orderShipmentId
  notNull: true
```

PostgreSQLでは、`sequence: true` が指定された `SEQUENCE_ID` Columnを以下の形式へ変換する。

```sql
order_shipment_id bigint GENERATED BY DEFAULT AS IDENTITY
```

`sequence` が指定されていない `SEQUENCE_ID` Columnは、通常の `bigint` として生成する。

```sql
order_shipment_id bigint
```

ARIADNE Sourceの `sequence` は論理的な連番採番を表し、PostgreSQLの物理的な `CREATE SEQUENCE` の使用を直接指定するものではない。

Prototype Phase 4では以下をValidation Ruleとする。

- `SEQUENCE_ID` + `sequence: true`：許可
- `SEQUENCE_ID` + `sequence` 未指定：許可
- `SEQUENCE_ID` 以外 + `sequence: true`：Error

`sequence` が未指定の場合は `false` と同等に扱う。

`sequence: true` のColumnがPrimary KeyまたはUnique Constraintであることは必須としない。
それらのDatabase ConstraintはDDL側で独立して定義する。

---

## 6. Column名と論理名

Element名とは異なるColumn識別子を定義可能とする。

例：

```yaml
billingCustomerId:
  element: customerId
```

この場合、

```text
Element              : customerId
Column識別子         : billingCustomerId
Column物理名         : billing_customer_id
```

として扱う。

さらにColumn固有の論理名が必要な場合は `name` を指定する。

```yaml
orderQuantity:
  element: quantity
  name: 受注数
```

`name` が指定されていない場合は参照Elementの `name` を利用する。

生成DDLのColumn Commentでは、Columnの実効論理名および参照元Elementを追跡可能な情報を利用する。

Commentの具体的な生成形式はPhase 4内で決定する。

---

## 7. BuiltIn Columns / Audit Columns

### 7.1 createdAt / updatedAt

すべてのTableに、以下のColumnをARIADNE BuiltIn Columnとして自動付与する。

- `createdAt`
- `updatedAt`

これらはDDL Source YAMLには記載しない。

DDL Source YAMLに利用者が `createdAt` または `updatedAt` を定義した場合はValidation Errorとする。

#### createdAt

`createdAt` はRecord作成日時を保持するBuiltIn Columnとする。

以下の仕様を持つ。

```text
Column識別子 : createdAt
物理Column名 : created_at
論理名       : 作成日時
Type         : DATETIME
NOT NULL     : true
Default      : CURRENT_DATETIME
```

INSERT時に現在日時を設定する。

UPDATE時には変更しない。

#### updatedAt

`updatedAt` はRecord最終更新日時を保持するBuiltIn Columnとする。

以下の仕様を持つ。

```text
Column識別子 : updatedAt
物理Column名 : updated_at
論理名       : 更新日時
Type         : DATETIME
NOT NULL     : true
Default      : CURRENT_DATETIME
```

INSERT時に現在日時を設定する。

UPDATE時にはARIADNEが生成するBuiltIn Function / Triggerにより現在日時へ自動更新する。

PostgreSQLでは、Schema単位のBuiltIn FunctionとTable単位のTriggerによって実現する。

SQLiteでは必要に応じて代替方式を採用する。

---

### 7.2 Audit

作成者・更新者等のAudit情報を保持する場合、Schema単位でAuditを有効化する。

AuditはOptionalとし、以下の形式で定義する。

```yaml
schema: received_order

audit:
  element: traceId
```

`audit.element` には、Audit Contextの値表現として利用するElementを指定する。

例：

```yaml
audit:
  element: userId
```

Auditが定義された場合、そのSchemaに属するすべてのTableへ以下のARIADNE BuiltIn Columnを自動付与する。

- `createdBy`
- `updatedBy`

Auditが定義されていない場合、`createdBy` / `updatedBy` は付与しない。

`createdAt` / `updatedAt` はAuditの有無にかかわらず、すべてのTableへ付与する。

---

### 7.3 createdBy / updatedBy

`createdBy` / `updatedBy` はARIADNE BuiltIn Columnとし、DDL Source YAMLには記載しない。

DDL Source YAMLに利用者が `createdBy` または `updatedBy` を定義した場合はValidation Errorとする。

#### createdBy

Record作成時のAudit Contextを保持する。

以下の仕様を持つ。

```text
Column識別子 : createdBy
物理Column名 : created_by
論理名       : 作成者
NOT NULL     : true
値表現       : audit.elementから導出
```

INSERT時にAudit Contextを設定する。

UPDATE時には変更しない。

#### updatedBy

Record最終更新時のAudit Contextを保持する。

以下の仕様を持つ。

```text
Column識別子 : updatedBy
物理Column名 : updated_by
論理名       : 更新者
NOT NULL     : true
値表現       : audit.elementから導出
```

INSERT時にAudit Contextを設定する。

UPDATE時には現在のAudit Contextへ更新する。

---

### 7.4 audit.elementから継承する情報

`audit.element` は、`createdBy` / `updatedBy` の値をDatabase上で表現するために利用する。

Types / Elementsから、少なくとも以下の情報を導出する。

- Database型へ変換するためのType情報
- length / minLength / maxLength
- precision / scale
- minimum / maximum
- regex
- その他Database上の値表現に必要なConstraint

一方、以下の情報は継承しない。

- Elementの `name`
- `identifier`
- Element固有の業務上の意味

`createdBy` / `updatedBy` は、参照Elementそのものを表すColumnではなく、ARIADNEが定義するAudit用BuiltIn Columnである。

---

### 7.5 Audit Context

Auditが有効なSchemaでは、アプリケーションからDatabaseへAudit Contextを伝搬する。

ARIADNEは「誰をAudit Contextとするか」というアプリケーション上の判断は行わない。

アプリケーションがTransaction開始後にAudit Contextを設定し、Database側のBuiltIn Function / Triggerがその値を取得する。

PostgreSQLでは、Transaction単位の設定値を利用する。

概念例：

```sql
SET LOCAL ariadne.audit_by = '...';
```

BuiltIn Functionでは、Transactionに設定された値を取得して `createdBy` / `updatedBy` へ設定する。

Connection Pool利用時に別RequestのAudit Contextが残存することを防ぐため、Session全体ではなくTransaction単位で管理する。

Auditが有効であるにもかかわらずAudit Contextが設定されていない場合はDatabase Errorとする。

`SET LOCAL` 等の具体的な実現方式はPostgreSQL向け実装仕様とし、ARIADNEのDBMS非依存なAudit概念そのものには含めない。

---

## 8. Primary Key

Primary KeyはTable側の責務としてDDL YAMLに定義する。

Elementsの `identifier: true` はAPI上の識別子として利用可能であることを示すものであり、Database上のPrimary Keyを意味しない。

Primary KeyはすべてのTableで必須とする。

TableごとにPrimary Keyを1つ定義しなければならない。

Primary Keyは単一Columnまたは複数Columnで構成できる。単一Column、複合Columnのいずれの場合も配列形式で定義する。

```yaml
primaryKey:
  - orderNo
```

複合Primary Keyは以下の形式とする。

```yaml
primaryKey:
  - orderNo
  - detailNo
```

配列の順序は物理Primary KeyのColumn順序として保持する。

Primary Keyに指定するColumnは、DDL Source YAML上で明示的に `notNull: true` が指定されていなければならない。

Primary Key指定によって暗黙的に `notNull: true` へ昇格させることはしない。

Primary Keyに指定されたColumnが `notNull: true` でない場合はValidation Errorとする。

`identifier`、`sequence`、`notNull`、Primary Keyはそれぞれ独立した意味として扱う。

Primary KeyのConstraint名はARIADNEがNaming Ruleに従って自動生成する。

DDL YAMLではConstraint名を指定しない。

---

## 9. Unique Constraint

データモデル上、一意でなければならないColumnまたはColumn組合せをUnique Constraintとして定義する。

Unique ConstraintはTableごとに複数定義可能とする。

基本形式を以下とする。

```yaml
uniqueConstraints:
  - columns:
      - operationDate
      - orderNo
      - detailNo
```

単一ColumnのUnique Constraintも同じ形式で定義する。

```yaml
uniqueConstraints:
  - columns:
      - customerId
```

各Unique Constraintには1つ以上のColumnを指定する。

配列の順序は物理Unique ConstraintのColumn順序として保持する。

同一Unique Constraint内で同じColumnを複数回指定してはならない。

同一のColumn集合を持つUnique Constraintを複数定義してはならない。

この重複判定ではColumn順序を区別しない。

したがって、以下は同一の一意性を表すためValidation Errorとする。

```text
UNIQUE (A, B)
UNIQUE (B, A)
```

Column順序は物理Constraint生成時には保持するが、順序の違いによって別の一意性とはみなさない。

Unique Constraint対象Columnに `notNull: true` を必須とはしない。

NULLを許容するUnique Constraintについては、基準DatabaseであるPostgreSQLの通常のUnique Constraint Semanticsに従う。

Unique Constraint名はARIADNEがNaming Ruleに従って自動生成する。

DDL YAMLではConstraint名を指定しない。

### 9.1 Unique Index

Prototype Phase 4ではUnique IndexをDDL YAMLの独立した定義対象とはしない。

一意性はUnique Constraintとして表現する。

PostgreSQL固有のPartial Unique Index、Expression Unique Index等が必要な場合はCustom SQLを利用する。

---

## 10. Foreign Key

Foreign KeyをDDL YAMLで管理可能とする。

Foreign KeyはTableごとに複数定義可能とする。

単一Foreign Keyの基本形式を以下とする。

```yaml
foreignKeys:
  - columns:
      - orderNo
    reference:
      table: orders
      columns:
        - orderNo
```

複合Foreign Keyは以下の形式とする。

```yaml
foreignKeys:
  - columns:
      - orderNo
      - detailNo
    reference:
      table: order_details
      columns:
        - orderNo
        - detailNo
```

`columns` は参照元Columnを表す。

`reference.table` は参照先Tableを表す。

`reference.columns` は参照先Columnを表す。

複合Foreign Keyでは、参照元 `columns` と参照先 `reference.columns` を配列順に対応付ける。

参照元と参照先のColumn数が一致しない場合はValidation Errorとする。

参照先Columnは、参照先TableのPrimary KeyまたはUnique Constraintとして一意性が保証されていなければならない。

Foreign Keyで対応する参照元Columnと参照先Columnは、同一Elementを参照していることを必須とする。

複合Foreign Keyでは、対応するすべてのColumnについてElementの一致をValidationする。

これはDatabase上の型互換性だけではなく、ARIADNE上で同一の意味を持つElement同士が参照関係を構成することを保証するためのRuleである。

SchemaをまたぐForeign KeyはPrototype Phase 4では許可しない。

Schemaをまたぐ参照が必要な場合はCustom SQLとして利用者の責任で定義する。

同一Tableを参照先とする自己参照Foreign Keyは許可する。

Foreign Keyを理由として、参照元ColumnにIndexを自動生成しない。

また、参照元ColumnがIndexを持つことをARIADNEのValidation条件とはしない。

Indexが必要な場合は、利用者が検索・更新特性に基づいて明示的に定義する。

### 10.1 Foreign Key Action

Foreign Keyには、参照先RecordのDELETE / UPDATEに対するActionを指定可能とする。

```yaml
foreignKeys:
  - columns:
      - orderNo
    reference:
      table: orders
      columns:
        - orderNo
    onDelete: CASCADE
```

Prototype Phase 4では以下を利用可能とする。

```text
NO_ACTION
CASCADE
SET_NULL
SET_DEFAULT
```

`onDelete` / `onUpdate` はOptionalとする。

未指定の場合は `NO_ACTION` として扱う。

`SET_NULL` を指定する場合、対象となる参照元ColumnはNULL許容でなければならない。

`SET_DEFAULT` を指定する場合、対象となる参照元ColumnにはDefaultが定義されていなければならない。

Foreign Key名はARIADNEがNaming Ruleに従って自動生成する。

DDL YAMLではConstraint名を指定しない。

Table生成後にForeign Keyを適用する等、依存関係を考慮したDDL生成順序はGenerator側で管理する。

ファイル名やファイル分割方法はDDL YAMLの概念モデルには含めない。

---

## 11. Index

検索性能等を目的として、通常IndexをDDL YAMLで定義可能とする。

IndexはTableごとに複数定義可能とする。

基本形式を以下とする。

```yaml
indexes:
  - columns:
      - column: orderNo
      - column: detailNo
      - column: operationDate
        order: DESC
```

単一Column Indexおよび複合Column Indexを許可する。

各Indexには1つ以上のColumnを指定する。

Column順序はIndex定義の一部として扱う。

各Columnには `order` を指定可能とする。

Prototype Phase 4では以下を許可する。

```text
ASC
DESC
```

`order` は大文字で指定する。

`order` が未指定の場合は `ASC` として扱う。

各Index Columnには、NULL値の並び順として `nulls` を指定可能とする。

```yaml
indexes:
  - columns:
      - column: operationDate
        order: DESC
        nulls: LAST
```

`nulls` には以下を指定可能とする。

```text
FIRST
LAST
```

`nulls` はOptionalとする。

`nulls` は、対象ColumnがNULL許容の場合のみ指定可能とする。

対象Columnに `notNull: true` が指定されている場合、`nulls` を指定してはならない。

NOT NULL Columnに `nulls` が指定されている場合はValidation Errorとする。

未指定の場合、ARIADNEはNULLの並び順を明示せず、対象DatabaseのDefault Semanticsに従う。

同一Index内で同じColumnを複数回指定してはならない。

Column順序、実効的なASC / DESC指定、および `nulls` 指定が完全に同一のIndexを複数定義してはならない。

以下は異なるIndexとして許可する。

```text
INDEX (A, B)
INDEX (B, A)
```

```text
INDEX (A ASC, B ASC)
INDEX (A DESC, B ASC)
```

Prototype Phase 4では、以下は構造化されたIndex定義の対象外とする。

- INCLUDE
- Partial Index / WHERE
- Expression Index
- Unique Index

これらが必要な場合はCustom SQLを利用する。

### 11.1 Primary Key / Unique Constraintとの重複

Primary KeyおよびUnique Constraintでは、一意性を保証するためのIndexがDatabase側で生成されることを前提とする。

そのため、Primary KeyまたはUnique Constraintによって生成されるIndexと完全に同一の
Column・同一順序・同一の実効Orderを持ち、`nulls` が指定されていない通常Indexは定義不可とする。

Primary KeyおよびUnique Constraintによって生成されるIndexの実効Orderは `ASC` として比較する。

Unique Constraintの対象ColumnがNULL許容であり、通常Index側に `nulls` が明示されている場合は、
NULLの並び順に対する設計意図が異なるため、別のIndexとして許可する。

例：

```text
UNIQUE (customerId, orderNo)
INDEX  (customerId ASC, orderNo ASC)
```

上記はValidation Errorとする。

一方、Column順序が異なる場合は許可する。

```text
UNIQUE (customerId, orderNo)
INDEX  (orderNo, customerId)
```

同じColumn順序であってもOrderが異なる場合は許可する。

```text
UNIQUE (customerId, orderNo)
INDEX  (customerId DESC, orderNo ASC)
```

Column数が異なる場合も許可する。

```text
UNIQUE (customerId, orderNo)
INDEX  (customerId)
```

```text
UNIQUE (customerId, orderNo)
INDEX  (customerId, orderNo, operationDate)
```

NULL許容Columnを含むUnique Constraintに対して、`nulls` が明示されたIndexは別のIndexとして許可する。

```text
UNIQUE (customerId, orderNo)
INDEX  (customerId ASC NULLS LAST, orderNo ASC)
```

ただし、`nulls` はNULL許容Columnにのみ指定可能である。

Primary KeyのColumnは必ずNOT NULLであるため、Primary Key対象Columnに `nulls` を指定したIndexは、
重複判定以前にValidation Errorとなる。

これらのIndexが実際に必要かどうかは設計者の判断とし、ARIADNEはValidation ErrorまたはWarningとはしない。

Foreign Keyの存在を理由としてIndexを自動生成したり、Indexの存在を要求したりしない。

Index名はARIADNEがNaming Ruleに従って自動生成する。

DDL YAMLではIndex名を指定しない。

---

## 12. Constraint / Index Naming

以下の物理名はARIADNEがNaming Ruleに従って自動生成する。

- Primary Key
- Unique Constraint
- Foreign Key
- Index

利用者はDDL YAML上でこれらの名称を指定しない。

Naming Ruleでは、Source YAMLに人工的なConstraint / Index識別子を持たせず、定義内容から決定的かつ安定した物理名を生成する。

### 12.1 Naming形式

物理名は以下の形式とする。

```text
Primary Key
  pk_<table>

Unique Constraint
  uq_<table>_<hash8>

Foreign Key
  fk_<table>_<hash8>

Index
  idx_<table>_<hash8>
```

Primary KeyはTableごとに1つであるためHashを付与しない。

Unique Constraint、Foreign Key、Indexには、定義内容から生成した8桁のHashを付与する。

### 12.2 Hash生成

HashにはSHA-256を使用する。

正規化したCanonical StringをUTF-8としてSHA-256でHash化し、lowercase hexadecimal表現の先頭8文字を利用する。

SchemaはHash計算に含めない。

Constraint / Indexの一意性はSchema内で管理する。

### 12.3 Unique Constraint

Unique ConstraintのHash入力には以下を利用する。

```text
table
columns[順序維持]
```

Canonical Stringの概念例：

```text
unique|shipping_instructions|operation_date,order_no,detail_no
```

Column順序はHash生成時に保持する。

### 12.4 Foreign Key

Foreign KeyのHash入力には以下を利用する。

```text
table
columns[順序維持]
reference.table
```

Canonical Stringの概念例：

```text
foreign_key|shipping_instructions|order_no,detail_no|order_details
```

`reference.columns`、`onDelete`、`onUpdate` はHash入力に含めない。

Foreign Key ActionはForeign Keyの属性であり、物理名上のIdentityとはしない。

### 12.5 Index

IndexのHash入力には以下を利用する。

```text
table
columns[順序維持 + effective order + nulls]
```

`nulls` が未指定の場合は、Canonical String上では `DEFAULT` へ正規化してHashを生成する。

`order` が未指定の場合は `ASC` へ正規化してからHashを生成する。

したがって、以下は同一のCanonical表現として扱う。

```yaml
- column: orderNo
```

```yaml
- column: orderNo
  order: ASC
```

Canonical Stringの概念例：

```text
index|shipping_instructions|order_no:ASC:DEFAULT,detail_no:ASC:DEFAULT,operation_date:DESC:LAST
```

これにより、Source上の省略表現の違いではHashを変化させず、意味が同一の定義から同一の物理名を生成する。

物理識別子長がDatabaseの上限を超える場合のTable名等の短縮Ruleは、Physical Naming Ruleとして別途定義する。

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

### 13.4 createdBy自動設定

Auditが有効なSchemaでは、INSERT時に現在のAudit Contextを `createdBy` へ設定する。

### 13.5 updatedBy自動設定・更新

Auditが有効なSchemaでは、INSERT時に現在のAudit Contextを `updatedBy` へ設定する。

UPDATE時には、現在のAudit Contextを `updatedBy` へ設定する。

### 13.6 PostgreSQLでのBuiltIn実装方針

PostgreSQLでは、`updatedAt` およびAudit Columnの自動更新を、Schema単位のBuiltIn FunctionとTable単位のTriggerによって実現する。

同一のBuiltIn FunctionでTimestamp更新およびAudit Context反映を扱うことを基本方針とする。

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

`GENERATED_ID` は値の型・意味をARIADNE Types / Elementsで定義する。

一方、実際の値をどの方式・どの主体で生成するかはDDL Sourceの責務とはしない。

Database側で値を生成する場合は、Prototype Phase 4では原則としてCustom SQLによる実装対象とする。

UUID / ULID等の一般化可能な生成方式については、Custom SQLの仕様検討時にARIADNE BuiltIn Functionとしてパターン化する価値があるかを検討する。

Phase 4では、業務固有採番を汎用YAMLモデルとして抽象化すること自体を目的としない。

---

## 16. Resolved DDL Model

SourceとなるDDL YAMLから、参照・省略表現・Naming Rule等を解決し、後続処理が一意に解釈可能なResolved DDL Modelを定義する。

Types / Elementsを正本とする情報はResolved DDL Modelへ複製せず、 `elementRef` により参照関係を保持する。

Database DDL生成時には、Resolved DDL ModelとTypes / Elementsを組み合わせて利用する。

```text
DDL Source YAML
      ↓
Resolved DDL Model ─────┐
                        ├─ PostgreSQL DDL
Types / Elements ───────┤
                        ├─ SQLite Compatibility DDL
                        └─ DDL Reference
```

Resolved DDL Modelでは、以下を解決済みとする。

- Service識別子
- Schema識別子
- 利用するENUM ElementとそのPhysical Name
- Table / ColumnのPhysical Name
- Element参照（`elementRef`）
- Column固有論理名のOverride
- Primary Key
- Unique Constraint
- Foreign Key
- Foreign Key ActionのDefault
- Index
- Constraint / Index Physical Name
- Constraint / Index Naming用Hash
- Index Columnの実効Order
- BuiltIn Columns
- Audit Columns
- Default
- Sequence / Generation Strategy
- Comment生成に必要な参照情報

Types / Elementsを正本とする以下の情報は、Resolved DDL Modelへ複製しない。

- 論理Type
- length / minLength / maxLength
- precision / scale
- minimum / maximum
- regex
- ENUM values
- format
- Element由来Constraint

これらはDatabase DDL生成時に `elementRef` からTypes / Elementsを参照して導出する。

DBMS固有SQLへの変換はResolved後の処理とする。

Prototype Phase 4ではResolved DDL Modelの仕様およびSampleを定義し、Generator実装は行わない。

### 16.1 Resolved DDL Model形式

Resolved DDL ModelはDDL Source単位で生成し、Service IDとSchema IDを保持する。

```yaml
formatVersion: "1.0"

service:
  id: order-management

schema:
  id: received_order

enums:
  orderStatus:
    elementRef: $orderStatus
    physicalName: order_status

tables:
  orders:
    name: 受注
    physicalName: orders

    columns:
      orderNo:
        elementRef: $receivedOrderNo
        physicalName: order_no
        notNull: true

      personInCharge:
        elementRef: $userId
        physicalName: person_in_charge
        override:
          name: 出荷担当者ID
```

Source YAMLの `service.id` は、Resolved DDL ModelでもService IDとして保持する。

Source YAMLの `schema` は、Resolved DDL Modelでは `schema.id` として保持する。

Service IDとSchema IDは独立した識別子として扱い、一致する必要はない。

Source YAML上のColumn `element` は、Resolved DDL Modelでは `elementRef` として保持する。

Table / ColumnのDatabase上の物理名はNaming Ruleを適用し、 `physicalName` として解決する。

Column固有の `name` が指定されている場合は、Elementの論理名に対するOverrideとして `override.name` に保持する。

### 16.2 BuiltIn / Audit Columns

DDL Source YAMLに記載されないBuiltIn Columnは、Resolved DDL Modelで明示的に展開する。

```yaml
createdAt:
  builtIn: system
  physicalName: created_at
  notNull: true
  default:
    expression: CURRENT_DATETIME

updatedAt:
  builtIn: system
  physicalName: updated_at
  notNull: true
  default:
    expression: CURRENT_DATETIME

createdBy:
  builtIn: audit
  elementRef: $traceId
  physicalName: created_by
  notNull: true

updatedBy:
  builtIn: audit
  elementRef: $traceId
  physicalName: updated_by
  notNull: true
```

`createdAt` / `updatedAt` は `builtIn: system` とする。

Auditが有効な場合に展開される `createdBy` / `updatedBy` は
`builtIn: audit` とし、`audit.element` を `elementRef` として保持する。

### 16.3 Constraint / Index

Primary Key、Unique Constraint、Foreign Key、Indexでは、
対象となるTable / Columnの物理名およびNaming Ruleによる物理名を解決する。

Primary Keyの例：

```yaml
primaryKey:
  physicalName: pk_order_details
  columns:
    - order_no
    - detail_no
```

Unique Constraintの例：

```yaml
uniqueConstraints:
  - physicalName: uq_shipping_instructions_ba368286
    columns:
      - operation_date
      - order_no
      - detail_no
    hash: ba3682869c048f7b175132052d27e475d10ebc570a0bf20dad5f06467631a7f3
```

Foreign Keyの例：

```yaml
foreignKeys:
  - physicalName: fk_shipping_instructions_3c47e712
    columns:
      - order_no
      - detail_no
    reference:
      table: order_details
      columns:
        - order_no
        - detail_no
    onUpdate: NO_ACTION
    onDelete: NO_ACTION
    hash: 3c47...
```

Foreign Keyの `onUpdate` / `onDelete` がSourceで省略されている場合、
Resolved DDL Modelでは実効値である `NO_ACTION` を明示する。

Indexの例：

```yaml
indexes:
  - physicalName: idx_tasks_fb4a92c9
    columns:
      - column: status_code
        order: ASC
      - column: priority
        order: DESC
      - column: task_pic
        order: ASC
        nulls: FIRST
      - column: task_id
        order: ASC
    hash: fb4...
```

Index Columnの `order` がSourceで省略されている場合、Resolved DDL Modelでは実効値である `ASC` を明示する。

`nulls` がSourceで省略されている場合はResolved DDL Modelにも出力せず、対象DatabaseのDefault Semanticsに従う。

Unique Constraint、Foreign Key、Indexでは、Naming Ruleによって算出したSHA-256 Hashを `hash` として保持する。

### 16.4 Default / Sequence

Literal DefaultはResolved DDL ModelでもLiteralとして保持する。

```yaml
status:
  elementRef: $orderStatus
  physicalName: status
  default: PREPARING
```

ARIADNE BuiltIn ExpressionもDBMS固有表現へ変換せず保持する。

```yaml
orderDate:
  elementRef: $orderDate
  physicalName: order_date
  default:
    expression: CURRENT_DATE
```

Source YAMLの `sequence: true` は、Resolved DDL Modelでは値生成方式として明示的に解決する。

```yaml
orderShipmentId:
  elementRef: $orderShipmentId
  physicalName: order_shipment_id
  notNull: true
  generation:
    strategy: SEQUENCE
```

`generation.strategy` はARIADNE上の値生成方式を表し、
PostgreSQLの `IDENTITY` 等のDBMS固有表現への変換はGeneratorの責務とする。

### 16.5 Comment生成情報

Database CommentそのものはResolved DDL Modelへ保持しない。

Table CommentはTableの `name` を生成元とする。

Column Commentは、Columnに `override.name` が存在する場合はその値を利用し、
存在しない場合は `elementRef` が参照するElementの `name` を利用する。

したがってComment生成に必要な情報はResolved DDL ModelとTypes / Elementsの参照関係によって保持され、
Comment文字列そのものをResolved DDL Modelへ複製しない。

### 16.6 Element由来Constraint

Types / Elementsに定義された型・桁・値域等のConstraintは、Resolved DDL Modelへ複製しない。

Resolved DDL Modelは `elementRef` を保持し、
Database DDL生成時に参照先のElement / Typeから必要なConstraintを導出する。

したがって、Element由来ConstraintのDatabase表現への変換はGeneratorの責務とする。

一方、DDL Source YAMLに明示されたColumn固有Constraintは、そのColumnに固有の定義であるためResolved DDL Modelへ保持する。

例えばSourceが以下の場合、

```yaml
orderQuantity:
  element: quantity
  notNull: true
  constraints:
    minimum: 1
```

Resolved DDL Modelでは以下として保持する。

```yaml
orderQuantity:
  elementRef: $quantity
  physicalName: order_quantity
  notNull: true
  constraints:
    minimum: 1
```

Element由来ConstraintとColumn固有ConstraintをResolved DDL Model上で統合しない。

Element由来Constraintは `elementRef` から参照し、Column固有Constraintは `constraints` として保持することで、両者の責務を分離する。

実効Constraintの決定、およびDatabase表現への変換はGeneratorの責務とする。

### 16.7 ENUM Dependency

Resolved DDL Modelでは、そのDDL Source内で利用するENUM Elementを依存情報として保持する。

Resolverは各Columnの `elementRef` から参照Element / Typeを解決し、Typeが `ENUM` であるElementを抽出する。

同一ENUM Elementが複数Columnから参照されている場合は、Element単位で重複を除去する。

例えば `orderStatus` Elementが利用されている場合、以下として保持する。

```yaml
enums:
  orderStatus:
    elementRef: $orderStatus
    physicalName: order_status
```

ENUMのMap KeyはElement IDとする。

`physicalName` はElement IDをPhysical Naming Ruleに従って `snake_case` へ変換した値とする。

ENUM valuesはTypes / Elementsを正本とし、Resolved DDL Modelへ複製しない。

Database DDL生成時には `elementRef` から参照Element / Typeを取得し、ENUM valuesを導出する。

同一ENUM Elementを複数Columnで利用する場合も、Database上では同一のENUM型を共有する。

利用するENUM Elementが存在しない場合は、Resolved DDL Modelでは以下として表現する。

```yaml
enums: []
```

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

Development EnvironmentとしてDocker Composeを使用し、PostgreSQLとpgwebを起動する。

Docker Compose定義は以下に配置する。

```text
tools/postgresql/compose.yaml
```

Repository Rootの `.env` をDocker Composeから明示的に読み込む。

Repository Rootから以下のWrapper Scriptにより起動・停止する。

```text
./tools/db-up.sh
./tools/db-down.sh
```

Database参照用としてpgwebを利用し、BrowserからDatabaseの状態を確認可能とする。

生成したPostgreSQL DDLは `dist/ddl/postgresql/` に配置し、PostgreSQL Containerの初期化SQLとして適用可能な構成とする。

PostgreSQL Databaseは生成DDLの実行検証を目的とした使い捨て環境とし、永続Volumeは使用しない。

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

Phase 4のValidationは、DDL Source YAMLの構造・記述形式を検証する
File Validationと、DDL Modelとしての意味的な整合性を検証する各種Validationに分離する。

Phase 2 / Phase 3からValidation IDを継続し、Phase 4では `V-083` 以降を使用する。

### 23.1 File Validation

File Validationでは、DDL Source YAMLの構造・記述形式・Source上の命名規約を検証する。

Phase 2 Element / Typeの参照解決や、Constraint間の意味的な整合性は後続Validationで検証する。

| ID    | Rule                                                                                                           | Level |
| ----- | -------------------------------------------------------------------------------------------------------------- | ----- |
| V-083 | 共通Headerの必須項目が存在する                                                                                 | Error |
| V-084 | `formatVersion` がARIADNEの対応するPhase 4 Formatである                                                        | Error |
| V-085 | `updatedAt` がISO 8601として妥当である                                                                         | Error |
| V-086 | `domain` が `ddl` である                                                                                       | Error |
| V-087 | `kind` が `database` である                                                                                    | Error |
| V-088 | 未定義属性を持たない                                                                                           | Error |
| V-089 | YAML Mapに重複Keyが存在しない                                                                                  | Error |
| V-090 | 空Map / 空Arrayを明示的に記述しない                                                                            | Error |
| V-091 | 必須文字列属性に空文字 / 空白のみを指定しない                                                                  | Error |
| V-092 | `schema` が存在する                                                                                            | Error |
| V-093 | `tables` が存在し、1件以上のTableを持つ                                                                        | Error |
| V-094 | Schema識別子が `snake_case` の命名規約を満たす                                                                 | Error |
| V-095 | Table識別子が `snake_case` の命名規約を満たす                                                                  | Error |
| V-096 | 各Tableが `name` を持つ                                                                                        | Error |
| V-097 | 各Tableが `columns` を持ち、1件以上のSource Columnを持つ                                                       | Error |
| V-098 | Column識別子が `lowerCamelCase` の命名規約を満たす                                                             | Error |
| V-099 | 各Columnが `element` を持つ                                                                                    | Error |
| V-100 | `notNull` を指定する場合は `true` のみを許可する                                                               | Error |
| V-101 | `sequence` を指定する場合は `true` のみを許可する                                                              | Error |
| V-102 | `audit` を指定する場合はMapであり、`element` を持つ                                                            | Error |
| V-103 | `default` はScalarまたはExpression Mapである                                                                   | Error |
| V-104 | Expression Map形式の `default` は `expression` を持つ                                                          | Error |
| V-105 | `default` を指定する場合、値は `null` ではない                                                                 | Error |
| V-106 | `primaryKey` はMapであり、`columns` を持つ                                                                     | Error |
| V-107 | `primaryKey.columns` は1件以上のArrayである                                                                    | Error |
| V-108 | `uniqueConstraints` を指定する場合、Unique Constraint Mapを要素とする1件以上のArrayである                      | Error |
| V-109 | 各Unique Constraintは `columns` を持ち、1件以上のArrayである                                                   | Error |
| V-110 | `foreignKeys` を指定する場合、Foreign Key Mapを要素とする1件以上のArrayである                                  | Error |
| V-111 | 各Foreign Keyは `columns` を持ち、1件以上のArrayである                                                         | Error |
| V-112 | 各Foreign Keyは `reference` Mapを持つ                                                                          | Error |
| V-113 | Foreign Keyの `reference` は `table` を持つ                                                                    | Error |
| V-114 | Foreign Keyの `reference` は `columns` を持ち、1件以上のArrayである                                            | Error |
| V-115 | `onDelete` / `onUpdate` を指定する場合は `NO_ACTION` / `CASCADE` / `SET_NULL` / `SET_DEFAULT` のいずれかである | Error |
| V-116 | `indexes` を指定する場合、Index Mapを要素とする1件以上のArrayである                                            | Error |
| V-117 | 各Indexは `columns` を持ち、Index Column Mapを要素とする1件以上のArrayである                                   | Error |
| V-118 | 各Index Columnが `column` を持つ                                                                               | Error |
| V-119 | Index Columnの `order` を指定する場合は `ASC` / `DESC` のいずれかである                                        | Error |
| V-120 | Index Columnの `nulls` を指定する場合は `FIRST` / `LAST` のいずれかである                                      | Error |
| V-121 | `service` が存在し、`id` を持つ                                                                                | Error |
| V-122 | Service IDが `kebab-case` の命名規約を満たす                                                                   | Error |
| V-123 | DDL Source YAMLのファイル名がService IDと一致する                                                              | Error |

### 23.2 Element / Column Validation

Element / Column Validationでは、DDL ColumnとPhase 2 Element / Typeとの意味的な整合性を検証する。

<!-- prettier-ignore-start -->

| ID | Rule | Level |
| ----- | --- | ----- |
| V-124 | Columnの `element` がPhase 2で定義されたElementとして存在する | Error |
| V-125 | ColumnからElement / Type由来の `type` / `length` / `minLength` / `maxLength` / `regex` / `minimum` / `maximum` / `precision` / `scale` / `enum` / `format` を再定義またはOverrideしない | Error |
| V-126 | `sequence: true` を指定したColumnの参照Element Typeが `SEQUENCE_ID` である | Error |
| V-127 | `createdAt` / `updatedAt` / `createdBy` / `updatedBy` をユーザー定義Columnとして定義しない | Error |
| V-128 | Columnの `constraints` に指定されたConstraintが参照Element / Typeで利用可能なConstraintである | Error |
| V-129 | Column固有ConstraintがElement由来Constraintと同一内容ではない | Error |
| V-130 | 比較可能なColumn固有ConstraintがElement由来Constraintを緩和しない | Error |
| V-131 | Column固有 `regex` がElement由来 `regex` と完全に同一ではない | Error |

<!-- prettier-ignore-end -->

Column固有ConstraintはElement由来ConstraintのOverrideではなく、追加Constraintとして扱う。

`minimum` / `maximum` / `minLength` / `maxLength` 等、
Constraint間の強弱を機械的に比較可能な場合は、Element由来Constraintと同一または緩いColumn固有ConstraintをValidation Errorとする。

`regex` については一般的な包含関係・強弱関係の解析を行わない。

Element由来Regexと異なるRegexは追加Constraintとして許可し、双方を満たすものとして扱う。

Element由来Regexと完全に同一のRegexは、意味のない重複としてValidation Errorとする。

`SEQUENCE_ID` Elementを参照するColumnであっても、 `sequence: true` の指定は必須ではない。

### 23.3 Default Validation

Default Validationでは、Column Defaultと参照Element / Typeとの整合性を検証する。

| ID    | Rule                                                                                                         | Level |
| ----- | ------------------------------------------------------------------------------------------------------------ | ----- |
| V-132 | Literal Defaultを指定した場合、参照Element Typeの `defaultAllowed.literal` が `true` である                  | Error |
| V-133 | Literal Defaultが参照Elementの論理Type / Format / Type固有Validation / Element Constraintを満たす            | Error |
| V-134 | Expression Defaultを指定した場合、そのExpressionが参照Element Typeの `defaultAllowed.expressions` に含まれる | Error |

Literal Defaultでは暗黙の型変換を行わない。

例えばINTEGERに対する `"0"` と `0`、BOOLEANに対する `"true"` と `true` は異なる値型として扱う。

任意のDatabase SQLをExpressionとして指定することはできない。
利用可能なExpressionは `defaultAllowed.expressions` により決定する。

### 23.4 Primary Key / Unique Constraint Validation

#### Primary Key

| ID    | Rule                                                                           | Level |
| ----- | ------------------------------------------------------------------------------ | ----- |
| V-135 | すべてのTableがPrimary Keyをちょうど1つ持つ                                    | Error |
| V-136 | `primaryKey.columns` に指定されたすべてのColumnが対象Tableに存在する           | Error |
| V-137 | 同一Primary Key内で同じColumnを重複指定しない                                  | Error |
| V-138 | Primary Keyを構成するすべてのColumnに明示的な `notNull: true` が指定されている | Error |

Primary KeyのColumn順序はSourceのArray記載順を保持する。

#### Unique Constraint

| ID    | Rule                                                                          | Level |
| ----- | ----------------------------------------------------------------------------- | ----- |
| V-139 | `uniqueConstraints[].columns` に指定されたすべてのColumnが対象Tableに存在する | Error |
| V-140 | 同一Unique Constraint内で同じColumnを重複指定しない                           | Error |
| V-141 | 同一Table内に同じColumn集合を持つUnique Constraintを複数定義しない            | Error |

Unique ConstraintのColumn順序はDDL生成時に保持する。

ただし、重複判定ではColumn順序を区別しない。

例えば `[a, b]` と `[b, a]` は、同一Column集合を持つUnique Constraintとして重複Errorとする。

### 23.5 Foreign Key Validation

Foreign Key Validationでは、Foreign Keyの参照関係およびReferential Actionの意味的整合性を検証する。

| ID    | Rule                                                                                                         | Level |
| ----- | ------------------------------------------------------------------------------------------------------------ | ----- |
| V-142 | Foreign KeyのLocal側 `columns` に指定されたすべてのColumnが対象Tableに存在する                               | Error |
| V-143 | `reference.table` に指定されたTableが同一Schema内に存在する                                                  | Error |
| V-144 | `reference.columns` に指定されたすべてのColumnが参照先Tableに存在する                                        | Error |
| V-145 | `reference.columns` が参照先TableのPrimary Keyまたは1つのUnique Constraintに対応する                         | Error |
| V-146 | `onDelete` / `onUpdate` が `SET_NULL` の場合、対象となるすべてのLocal ColumnがNULL許容である                 | Error |
| V-147 | `onDelete` / `onUpdate` が `SET_DEFAULT` の場合、対象となるすべてのLocal Columnに `default` が定義されている | Error |
| V-148 | 同一Foreign KeyのLocal側 `columns` に同じColumnを重複指定しない                                              | Error |
| V-149 | 同一Foreign Keyの `reference.columns` に同じColumnを重複指定しない                                           | Error |
| V-150 | 同一Table内に、Local Columnsの並びとReference Tableが同一であるForeign Keyを複数定義しない                   | Error |

Foreign KeyのIdentityは以下により決定する。

```text
Table
＋ Local Columns（記載順）
＋ Reference Table
```

`reference.columns`、`onDelete`、`onUpdate` は Foreign KeyのIdentityには含めない。

したがって、同一Tableから同一Local Columnsを使用して
同一Reference Tableへ複数のForeign Keyを定義することは許可しない。

Foreign Keyは同一Schema内のみを参照可能とする。

Self Referenceは許可する。

### 23.6 Foreign Key Element一致 Validation

| ID    | Rule                                                                                     | Level |
| ----- | ---------------------------------------------------------------------------------------- | ----- |
| V-151 | Foreign KeyのLocal Columnと対応するReference Columnが同一のPhase 2 Elementを参照している | Error |

Database上の型が互換であっても、異なるElement間にForeign Keyを定義することは許可しない。

Foreign KeyにおけるElement一致は、
Database型の一致ではなく、同一の業務上の値を参照していることを保証するためのRuleとする。

### 23.7 Composite Foreign Key Validation

| ID    | Rule                                                                           | Level |
| ----- | ------------------------------------------------------------------------------ | ----- |
| V-152 | 複合Foreign KeyではLocal側 `columns` と `reference.columns` の要素数が一致する | Error |

複合Foreign KeyのColumn対応はArrayの記載順とする。

第n Local Columnは第n Reference Columnに対応する。

各対応Columnには Foreign Key Element一致Validationを適用する。

### 23.8 Index Validation

| ID    | Rule                                                                            | Level |
| ----- | ------------------------------------------------------------------------------- | ----- |
| V-153 | `indexes[].columns[].column` に指定されたすべてのColumnが対象Tableに存在する    | Error |
| V-154 | 同一Index内で同じColumnを複数回指定しない                                       | Error |
| V-155 | 同一Table内に同一Canonical Identityを持つIndexを複数定義しない                  | Error |
| V-156 | `notNull: true` のColumnに `nulls` を指定しない                                 | Error |
| V-157 | Primary Key / Unique ConstraintのBacking Indexと同一定義の通常Indexを定義しない | Error |

IndexのCanonical Identityは以下により決定する。

```text
Table
＋ Columns（記載順）
    ＋ Effective Order
    ＋ Nulls
```

Canonical Identity上では以下として扱う。

```text
order省略 → ASC
nulls省略 → DEFAULT
```

したがって、

```yaml
- column: orderNo
```

と、

```yaml
- column: orderNo
  order: ASC
```

は同一定義として扱う。

一方、`nulls` の明示指定は設計意図としてCanonical Identityに含める。

Primary Key / Unique Constraintとの重複判定では、
Primary Key / Unique Constraint側のEffective Orderを `ASC`、 `nulls` を `DEFAULT` として比較する。

Column順序、Order、または有効な `nulls` 指定が異なるIndexは別のIndexとして許可する。

Prefix Indexも別定義として許可する。

例：

```text
Unique Constraint : (a, b)
Index             : (a)
```

上記Indexは重複とはみなさない。

### 23.9 Naming Ruleとの関係

Constraint / Indexの物理名はARIADNEがNaming Ruleに従って自動生成するため、Naming自体を独立したValidation Ruleとはしない。

生成名は以下の形式とする。

```text
Primary Key       : pk_<table>
Unique Constraint : uq_<table>_<hash8>
Foreign Key       : fk_<table>_<hash8>
Index              : idx_<table>_<hash8>
```

`hash8` は対象定義のCanonical InputをUTF-8で表現し、SHA-256を計算した結果の先頭8桁をlowercaseで使用する。

Unique ConstraintのCanonical Inputは以下とする。

```text
Table
＋ Columns（記載順）
```

Foreign KeyのCanonical Inputは以下とする。

```text
Table
＋ Local Columns（記載順）
＋ Reference Table
```

IndexのCanonical Inputは以下とする。

```text
Table
＋ Columns（記載順）
    ＋ Effective Order
    ＋ Nulls
```

ARIADNEが異なる定義から同一物理名を生成するHash Collision等が発生した場合は、
Source ValidationではなくGeneratorの生成不能Errorとして扱う。

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

- [x] DDL YAML全体構造を確定する
- [x] Schema定義形式を確定する
- [x] Table定義形式を確定する
- [x] Column定義形式を確定する
- [x] Element参照方式を確定する
- [x] Column別名定義方式を確定する
- [x] Default定義方式を確定する
- [x] Element由来Constraintの導出規則を確定する

### Step 2：BuiltIn / Audit

- [x] createdAtの仕様を確定する
- [x] updatedAtの仕様を確定する
- [x] Audit Trace Elementの指定方式を確定する
- [x] created側Audit Trace Columnの仕様を確定する
- [x] updated側Audit Trace Columnの仕様を確定する
- [x] Audit Trace値の伝搬方式を検証する

### Step 3：Constraint / Index

- [x] Primary Key定義形式を確定する
- [x] Unique Constraint定義形式を確定する
- [x] Foreign Key定義形式を確定する
- [x] 複合Foreign Key定義形式を確定する
- [x] Index定義形式を確定する
- [x] Unique Constraint / Index重複Ruleを確定する
- [x] Constraint / Index Naming Ruleを確定する
- [x] Foreign Key Action仕様を確定する
- [x] Foreign Key参照範囲・自己参照Ruleを確定する

### Step 4：Validation Rule

- [x] Element / Column Validation Ruleを確定する
- [x] Default Validation Ruleを確定する
- [x] PK / Unique Validation Ruleを確定する
- [x] FK Validation Ruleを確定する
- [x] FK Element一致Ruleを確定する
- [x] 複合FK Validation Ruleを確定する
- [x] Index Validation Ruleを確定する
- [x] Naming Validation Ruleを確定する

### Step 5：Resolved DDL Model

- [x] Resolved DDL Modelの形式を確定する
- [x] Types / Elementsの解決結果をSample化する
- [x] BuiltIn Columnsの展開結果をSample化する
- [x] Audit Trace Columnsの展開結果をSample化する
- [x] Constraint / Indexの解決結果をSample化する
- [x] Comment情報の解決結果をSample化する
- [x] Physical Namingの解決結果をSample化する
- [x] Element由来Constraintの解決結果をSample化する
- [x] Default / Sequenceの解決結果をSample化する

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

- [x] PostgreSQL Docker環境を作成する
- [x] pgwebを追加する
- [ ] Order Management Sample DDLを適用する
- [ ] Table / Constraint / Indexを確認する
- [ ] BuiltIn Triggerを確認する
- [ ] Custom SQLを確認する
- [x] pgwebからDatabaseを確認する

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
