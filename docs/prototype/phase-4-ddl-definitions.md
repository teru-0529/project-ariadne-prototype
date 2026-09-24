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
- Custom SQLのモデル上の扱いの定義
- DDL Outputの基本構造および論理的な実行順序の定義
- Sample DDLの作成
- PostgreSQLでの実行検証

実際のDDL Generator / Validatorの実装はCoreで行う。

SQLiteを含むARIADNE自身の内部DatabaseについてはPrototypeでは扱わず、必要性および採用方式をCoreで改めて判断する。

DDL Referenceの生成およびBrowserからの参照機能についてもPrototypeでは実装せず、Coreで扱う。

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

DDL YAMLでは可能な限りDBMS固有表現を持たないが、
機能および意味論の基準はPostgreSQLとする。

```text
DDL YAML
   ↓
Resolved DDL Model
   ↓
PostgreSQL Semantics
   ↓
PostgreSQL DDL
```

Prototype Phase 4では、PostgreSQL向けDDLの設計および実行検証を対象とする。

他Databaseへの対応については、PostgreSQL向け仕様そのものを制限せず、必要性が生じた時点でCoreにおいて検討する。

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

`name` はTableの論理名・表示名、およびDatabase Table Commentの生成元として利用する。

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
- `generation`
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

`generation` はElement / Typeの値定義をOverrideするものではなく、
そのColumn配置において値を生成するか、およびその生成主体・生成方式を指定するColumn固有の配置・利用情報として扱う。

値生成が可能かどうかは、参照ElementのTypeが持つ `capability.generation` によって決定する。
Typeは値そのものの意味を定義し、実際にそのColumnで値生成を行うかどうかは `generation` によって定義する。

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

DefaultとGenerationは別の責務として扱う。

DefaultはINSERT時に値が明示されなかった場合の初期値を定義する。

同一Columnに `default` と `generation` を同時に指定することはできない。

値そのものを生成する場合はDefaultではなく、`generation` を利用する。

### 5.7 Generation

Columnでは、そのColumn配置において値を生成する場合、`generation` を指定する。

Generationは、Element / Typeが表す値そのものの意味とは分離して扱う。

Typeは値の意味および値生成可能性を定義し、
Columnの `generation` は、そのColumn配置において実際に値を生成するか、およびその生成主体・生成方式を定義する。

`generation` を利用できるのは、参照ElementのTypeが以下を満たす場合のみとする。

```yaml
capability:
  generation: true
```

Prototype Phase 4では、Generationを以下の3種類に分類する。

```text
Generation
├─ DATABASE
├─ BUILTIN
│  └─ ULID
└─ CUSTOM
   └─ Custom Generation Function
```

#### DATABASE Generation

Database自身の標準的な値生成機構を利用する場合、以下の形式で指定する。

```yaml
orderShipmentId:
  element: orderShipmentId
  notNull: true
  generation: DATABASE
```

`DATABASE` はDBMS非依存なARIADNE上の意味として扱う。

具体的なDatabase表現への変換は各Database Generatorの責務とする。

例えば、`SEQUENCE` Typeを参照するColumnに対する `generation: DATABASE` は、
PostgreSQL Generatorでは以下のようなDatabaseネイティブの連番生成へ変換する。

```sql
order_shipment_id bigint GENERATED BY DEFAULT AS IDENTITY
```

`SEQUENCE` Typeそのものは「順序を持つ連番整数値」であることを表す。

したがって、同じ `SEQUENCE` Elementを参照するColumnであっても、
そのColumn自身が値を生成しない場合は `generation` を指定しない。

```yaml
detailNo:
  element: orderDetailNo
  notNull: true
```

この場合、PostgreSQLでは通常の `bigint` Columnとして生成する。

#### BUILTIN Generation

ARIADNEが標準提供する値生成方式を利用する場合、
BuiltIn Generation名を `generation` に指定する。

Prototype Phase 4では、少なくとも以下を定義する。

```text
ULID
```

例：

```yaml
taskId:
  element: taskId
  notNull: true
  generation: ULID
```

`ULID` はARIADNE上のDBMS非依存な値生成方式であり、具体的なFunction / Trigger等への変換は各Database Generatorの責務とする。

#### CUSTOM Generation

Service固有の値生成ロジックを利用する場合、Custom Generation Functionを以下の形式で指定する。

```yaml
orderNo:
  element: receivedOrderNo
  notNull: true
  generation:
    custom: generateOrderNo
```

Custom Generation Functionは対象TableのRowを入力として受け取り、
対象Columnへ設定する単一の値を返す。

```text
Custom Generation Function

input  : target Table Row
output : generated value
```

Custom Generation Function自身は対象Columnを書き換えない。

生成された値を対象Columnへ設定する責務はARIADNEが生成するGeneration Triggerが持つ。

例えば概念的には以下の処理となる。

```sql
NEW.order_no :=
  received_order.generate_order_no(NEW);
```

CUSTOM Generationでは、同一Row内の他のGeneration Columnが生成した値への依存を許可しない。

Generation Column間の実行順序は保証しない。

#### GenerationとDatabase Constraint

`generation` はPrimary Key、Unique Constraint、Foreign Key等のDatabase Constraintを意味しない。

値生成とDatabase Constraintは独立した責務として定義する。

また、`generation` が指定されていないことは、そのElement / Typeが値生成不可能であることを意味しない。

同一Elementを参照する複数Columnのうち、一部のColumnのみでGenerationを行うこともできる。

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

Element / Typeが表す値の意味やType Capabilityは、Database上のPrimary Keyを意味しない。

Primary Keyは、対象TableにおけるDatabase ConstraintとしてDDL Source YAMLで明示的に定義する。

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

Columnの `generation`、`notNull`、Primary Keyはそれぞれ独立した意味として扱う。

値生成を行うColumnがPrimary Keyである必要はなく、Primary Keyを構成するColumnが値生成を行う必要もない。

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

## 15. Custom Generation Function

Service固有の値生成ロジックをDatabase Functionとして実装する場合、
DDL Source YAMLではCUSTOM Generationとして定義する。

代表的な利用例：

- Prefix + 連番
- 日付 + 連番
- 年月 + 連番
- Schema単位の採番
- Table単位の採番
- 業務単位の採番

例えば、以下のような受注番号を生成する場合、

```text
ORD_20260914_002
```

DDL Source YAMLでは対象ColumnにCustom Generation Functionを指定する。

```yaml
orderNo:
  element: receivedOrderNo
  notNull: true
  generation:
    custom: generateOrderNo
```

Custom Generation Functionは、対象TableのRowを入力として受け取り、
Generation対象Columnへ設定する単一の値を返すDatabase Functionとする。

```text
Custom Generation Function

input  : target Table Row
output : generated value
```

Custom Generation Function自身はGeneration対象Columnを書き換えない。

対象Columnへの値設定は、ARIADNEが生成するGeneration Triggerの責務とする。

例えばPostgreSQLでは、概念的に以下のようなFunction呼び出しをGeneration Triggerから行う。

```sql
NEW.order_no :=
  received_order.generate_order_no(NEW);
```

Custom Generation Functionは、対象Serviceに属するCustom Functionとして管理する。

Functionの実装自体はDatabase固有SQLとして記述することを許容する。

ただし、DDL Source YAML上のCUSTOM Generation定義によって、
そのFunctionが単なるCustom SQLではなく、Columnの値生成に利用されるFunctionであることをARIADNEが認識する。

Custom Generation Functionは以下の契約を満たさなければならない。

- 対象TableのRowを入力として受け取る
- Generation対象Columnへ設定する単一値を返す
- 戻り値型がGeneration対象ColumnのDatabase物理型と一致する
- Generation対象Column自身を書き換えない
- 同一Row内の他のGeneration Columnが生成した値へ依存しない

Prototype Phase 4ではCustom Generation Functionの戻り値に対する暗黙の型変換を許可しない。

複数のGeneration Columnが存在する場合、それらのGeneration実行順序は保証しない。

Custom Generation Functionと、一般的なCustom Function / Custom Triggerは責務を分離する。

```text
Custom Generation Function
  → RowからGeneration対象Columnの値を生成する

Custom Function / Custom Trigger
  → Generation以外の業務ロジックを実装する
```

DATABASE GenerationおよびARIADNE BuiltIn GenerationはCustom Generation Functionとは別のGeneration方式として扱う。

ULID等、ARIADNEが標準提供する値生成方式はBuiltIn Generationとして定義し、
Service固有のCustom Generation Functionを利用者が実装する必要はない。

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
Types / Elements ───────┘
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
- Generation
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

### 16.4 Default

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

### 16.5 Generation

Source YAMLの `generation` は、Resolved DDL ModelではGeneration種別を明示した共通形式へ正規化する。

#### DATABASE Generation

Source：

```yaml
generation: DATABASE
```

Resolved：

```yaml
generation:
  kind: DATABASE
```

`DATABASE` はDatabase自身の標準的な値生成機構を利用することを表す。

具体的なDatabase表現はResolved DDL Modelでは保持しない。

例えば `SEQUENCE` Typeを参照するColumnの場合、
PostgreSQLの `IDENTITY` 等への変換はPostgreSQL Generatorの責務とする。

#### BUILTIN Generation

Source：

```yaml
generation: ULID
```

Resolved：

```yaml
generation:
  kind: BUILTIN
  strategy: ULID
```

`kind: BUILTIN` はARIADNEが標準提供するGenerationであることを表す。

`strategy` にはARIADNE上のBuiltIn Generation名を保持する。

具体的なDatabase Function / Trigger等への変換は各Database Generatorの責務とする。

#### CUSTOM Generation

Source：

```yaml
generation:
  custom: generateOrderNo
```

Resolved：

```yaml
generation:
  kind: CUSTOM
  functionRef: $generateOrderNo
```

`kind: CUSTOM` はService固有のCustom Generation Functionを利用することを表す。

Sourceの `custom` に指定されたFunction識別子は、
Resolved DDL Modelでは `functionRef` として参照関係を明示する。

Custom Generation Functionの実装内容や戻り値型はResolved DDL Modelへ複製しない。

Generation対象Column、`functionRef` が参照するFunction、およびTypes / Elementsを組み合わせて後続Generatorが利用する。

Generationを実現するために必要となるGeneration Triggerは、Resolved DDL Modelへ展開しない。

Resolved DDL ModelではColumnの `generation` を値生成の正本として保持し、
各Database GeneratorがGeneration方式に応じて必要なDatabase Function / Trigger等を生成する。

したがって、GenerationのためにARIADNEが自動生成するTriggerは `customTriggers` には含めない。

`customTriggers` はDDL Source YAMLで利用者が明示的に定義したCustom Triggerのみを保持する。

```text
Resolved DDL Model
├─ Column.generation
│    └─ Generatorが必要なGeneration Triggerを生成
│
└─ customTriggers
     └─ 利用者定義のCustom Triggerのみ
```

Generation TriggerとCustom Triggerは、Resolved DDL Model上でも別の責務として扱う。

### 16.6 Comment生成情報

Database CommentそのものはResolved DDL Modelへ保持しない。

Table CommentはTableの `name` を生成元とする。

Column Commentは、Columnに `override.name` が存在する場合はその値を利用し、
存在しない場合は `elementRef` が参照するElementの `name` を利用する。

したがってComment生成に必要な情報はResolved DDL ModelとTypes / Elementsの参照関係によって保持され、
Comment文字列そのものをResolved DDL Modelへ複製しない。

### 16.7 Element由来Constraint

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

### 16.8 ENUM Dependency

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

PostgreSQL Generatorは、Resolved DDL ModelとTypes / Elementsを入力として、PostgreSQLで実行可能なDDLを生成する。

Phase 4では以下を生成対象とする。

- Schema
- ENUM
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
- DATABASE Generation
- BUILTIN Generation
- CUSTOM Generation
- BuiltIn Function
- BuiltIn Trigger
- Generation Trigger
- Custom Function
- Custom Trigger

### 17.1 DDL生成順序

PostgreSQL DDLは、Database Object間の依存関係を満たす順序で生成する。

基本的な生成順序を以下とする。

```text
Schema
  ↓
ENUM
  ↓
Table / Column / Default / Check Constraint
  ↓
Primary Key / Unique Constraint
  ↓
Index
  ↓
Foreign Key
  ↓
BuiltIn Function / Custom Function
  ↓
BuiltIn Trigger / Generation Trigger / Custom Trigger
```

Functionを利用するTriggerは、参照するFunctionの生成後に生成する。

Custom Function / Custom Triggerについても、ARIADNE生成Objectとの依存関係を満たす順序で配置する。

生成DDLを単一ファイルとするか複数ファイルへ分割するかは論理モデルには含めず、
Generatorの出力仕様として扱う。

### 17.2 Schema / ENUM / Table

SchemaはResolved DDL Modelの `schema.id` から生成する。

Prototype Phase 4のSample DDLでは、再実行可能な検証環境を構築するため、
Schemaを再作成する形式を利用する。

```sql
DROP SCHEMA IF EXISTS received_order CASCADE;
CREATE SCHEMA received_order;
```

ENUMは、Resolved DDL Modelの `enums` とTypes / Elementsから値集合を導出し、
PostgreSQL ENUM型としてSchema内に生成する。

Table / ColumnはResolved DDL ModelのPhysical Nameを利用して生成する。

ColumnのDatabase型、桁、値制約等は `elementRef` からTypes / Elementsを参照して導出する。

Element由来ConstraintおよびColumn固有Constraintのうち、
Database型のみでは保証できないものはCHECK Constraintとして生成する。

### 17.3 Constraint / Index / Comment

Primary Key、Unique Constraint、Foreign Key、Indexは、
Resolved DDL Modelで解決済みのPhysical NameおよびPhysical Column Nameを利用して生成する。

Foreign Keyは参照先Table / Constraintが生成された後に適用する。

Table CommentにはTableの論理名を利用する。

Column CommentにはColumnの実効論理名を利用し、参照Elementを追跡可能な情報を含める。

### 17.4 DATABASE Generation

`generation.kind: DATABASE` は、対象Typeに対応するPostgreSQLネイティブの値生成機構へ変換する。

Prototype Phase 4では、`SEQUENCE` TypeのDATABASE Generationを PostgreSQL Identity Columnとして生成する。

```sql
order_shipment_id bigint
  GENERATED BY DEFAULT AS IDENTITY
```

DATABASE GenerationではARIADNEによるGeneration Triggerを生成しない。

値生成はPostgreSQL自身へ委譲する。

### 17.5 BUILTIN Generation

`generation.kind: BUILTIN` は、ARIADNEが提供するPostgreSQL向けBuiltIn Generationへ変換する。

Prototype Phase 4では `ULID` を対象とする。

ULID Generationでは、Schema単位のBuiltIn Functionを生成する。

概念例：

```sql
CREATE FUNCTION task.ariadne_builtin_ulid()
RETURNS varchar
...
```

BuiltIn Generation対象Columnを持つTableには、
ARIADNEがTable単位のGeneration Trigger FunctionおよびGeneration Triggerを生成する。

概念例：

```sql
NEW.task_id := task.ariadne_builtin_ulid();
```

BuiltIn ULID Functionは対象Rowを引数として受け取らない。

Generation TriggerがBuiltIn Functionを呼び出し、返された値をGeneration対象Columnへ設定する。

### 17.6 CUSTOM Generation

`generation.kind: CUSTOM` は、
Service固有のCustom Generation FunctionとARIADNEが生成するGeneration Triggerを組み合わせて実現する。

Custom Generation Functionは対象TableのRowを入力として受け取り、
Generation対象Columnへ設定する単一値を返す。

```text
Custom Generation Function

input  : target Table Row
output : generated value
```

Custom Generation Function自身はGeneration対象Columnを書き換えない。

対象Columnへの代入はGeneration Triggerが行う。

概念例：

```sql
NEW.order_no :=
  received_order.generate_order_no(NEW);
```

同一Tableに複数のGeneration Columnが存在する場合、ARIADNEはそれらのGeneration処理順序を仕様として保証しない。

そのためCustom Generation Functionは、同一Row内の他のGeneration Columnが生成した値へ依存してはならない。

### 17.7 BuiltIn Function / BuiltIn Trigger

`createdAt` / `updatedAt` およびAudit Columnの自動設定・更新は、
ARIADNEが生成するBuiltIn Function / BuiltIn Triggerによって実現する。

PostgreSQLではSchema単位のBuiltIn Functionと、Table単位のRow Triggerを生成する。

INSERT時には以下を設定する。

- `createdAt`
- `updatedAt`
- Audit有効時の `createdBy`
- Audit有効時の `updatedBy`

UPDATE時には以下を更新する。

- `updatedAt`
- Audit有効時の `updatedBy`

Audit ContextはTransaction単位で設定された値から取得する。

Auditが有効であるにもかかわらずAudit Contextが存在しない場合はDatabase Errorとする。

### 17.8 Trigger

ARIADNEが扱うTriggerは責務によって以下の3種類に分離する。

```text
BuiltIn Trigger
  → createdAt / updatedAt / Audit等のARIADNE共通Behavior

Generation Trigger
  → Generation対象Columnへの値設定

Custom Trigger
  → Service固有の業務・イベント処理
```

BuiltIn TriggerおよびGeneration TriggerはARIADNEが生成する。

Custom Triggerは利用者が定義したDatabase固有SQLを利用する。

PostgreSQLでは、同一Table・同一Eventに複数Triggerが存在する場合、Trigger Nameのアルファベット順で実行される。

また、BEFORE Triggerが返した変更後のRowは、後続するBEFORE Triggerの入力となる。

ただし、ARIADNEの論理仕様ではTrigger Nameによる実行順序を Generation間の依存関係解決手段として利用しない。

各GenerationおよびCustom処理は、暗黙のTrigger実行順序へ依存しない設計を基本とする。

### 17.9 PostgreSQL実行検証

Phase 4ではOrder ManagementおよびTask ManagementのSample DDLを生成し、
PostgreSQL Runtimeへ適用して以下を実証する。

- Schema / ENUM / Table生成
- Data Type / Default / Check Constraint
- Primary Key / Unique Constraint / Foreign Key
- Index
- Comment
- DATABASE Generation
- BUILTIN GenerationによるULID生成
- CUSTOM Generation
- Generation Trigger
- BuiltIn Function / BuiltIn Trigger
- Custom Function / Custom Trigger
- Audit Context
- INSERT / UPDATE

Sample DDLは使い捨てのPostgreSQL Runtimeへ先頭から再適用し、
初期状態から正常に構築・実行できることを確認する。

実Generator / Validatorの実装はCoreで行う。

---

## 18. DDL Output

生成DDLのファイル構成はARIADNEの論理モデルとは分離する。

Prototype Phase 4では、生成DDLの具体的なファイル分割方式を論理モデルには含めない。

Generatorは、生成DDLの物理的なファイル構成とは独立して、依存関係を満たす論理的な実行順序を管理する。

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

## 19. Runtime

Phase 4では、設計したPostgreSQL DDLが実Database上で成立することを検証する。

### 19.1 PostgreSQL

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

---

## 20. Sample / Verification

### 20.1 Task Management

主用途：

- 小規模なDDL Modelによる基本仕様の確認
- Element / Column / Constraint / BuiltIn等の基本変換規則の確認

### 20.2 Order Management

主用途：

- ARIADNEが設計対象とする業務システムのDatabase Model
- PostgreSQL DDL仕様の実証
- Schema / FK / Index / Constraint / Function等を含む実践的検証

---

## 21. Validation Rule

Prototype Phase 4ではValidatorそのものは実装しない。

将来のValidatorが検証すべきRuleを定義する。

Phase 4のValidationは、DDL Source YAMLの構造・記述形式を検証する
File Validationと、DDL Modelとしての意味的な整合性を検証する各種Validationに分離する。

Phase 2 / Phase 3からValidation IDを継続し、Phase 4では `V-085` 以降を使用する。

### 21.1 File Validation

File Validationでは、DDL Source YAMLの構造・記述形式・Source上の命名規約を検証する。

Phase 2 Element / Typeの参照解決や、Constraint間の意味的な整合性は後続Validationで検証する。

| ID    | Rule                                                                                                           | Level |
| ----- | -------------------------------------------------------------------------------------------------------------- | ----- |
| V-085 | 共通Headerの必須項目が存在する                                                                                 | Error |
| V-086 | `formatVersion` がARIADNEの対応するPhase 4 Formatである                                                        | Error |
| V-087 | `updatedAt` がISO 8601として妥当である                                                                         | Error |
| V-088 | `domain` が `ddl` である                                                                                       | Error |
| V-089 | `kind` が `database` である                                                                                    | Error |
| V-090 | 未定義属性を持たない                                                                                           | Error |
| V-091 | YAML Mapに重複Keyが存在しない                                                                                  | Error |
| V-092 | 空Map / 空Arrayを明示的に記述しない                                                                            | Error |
| V-093 | 必須文字列属性に空文字 / 空白のみを指定しない                                                                  | Error |
| V-094 | `schema` が存在する                                                                                            | Error |
| V-095 | `tables` が存在し、1件以上のTableを持つ                                                                        | Error |
| V-096 | Schema識別子が `snake_case` の命名規約を満たす                                                                 | Error |
| V-097 | Table識別子が `snake_case` の命名規約を満たす                                                                  | Error |
| V-098 | 各Tableが `name` を持つ                                                                                        | Error |
| V-099 | 各Tableが `columns` を持ち、1件以上のSource Columnを持つ                                                       | Error |
| V-100 | Column識別子が `lowerCamelCase` の命名規約を満たす                                                             | Error |
| V-101 | 各Columnが `element` を持つ                                                                                    | Error |
| V-102 | `notNull` を指定する場合は `true` のみを許可する                                                               | Error |
| V-103 | `generation` は `DATABASE` / BuiltIn Generation名、または `custom` を持つMapのいずれかである                   | Error |
| V-104 | `audit` を指定する場合はMapであり、`element` を持つ                                                            | Error |
| V-105 | `default` はScalarまたはExpression Mapである                                                                   | Error |
| V-106 | Expression Map形式の `default` は `expression` を持つ                                                          | Error |
| V-107 | `default` を指定する場合、値は `null` ではない                                                                 | Error |
| V-108 | `primaryKey` はMapであり、`columns` を持つ                                                                     | Error |
| V-109 | `primaryKey.columns` は1件以上のArrayである                                                                    | Error |
| V-110 | `uniqueConstraints` を指定する場合、Unique Constraint Mapを要素とする1件以上のArrayである                      | Error |
| V-111 | 各Unique Constraintは `columns` を持ち、1件以上のArrayである                                                   | Error |
| V-112 | `foreignKeys` を指定する場合、Foreign Key Mapを要素とする1件以上のArrayである                                  | Error |
| V-113 | 各Foreign Keyは `columns` を持ち、1件以上のArrayである                                                         | Error |
| V-114 | 各Foreign Keyは `reference` Mapを持つ                                                                          | Error |
| V-115 | Foreign Keyの `reference` は `table` を持つ                                                                    | Error |
| V-116 | Foreign Keyの `reference` は `columns` を持ち、1件以上のArrayである                                            | Error |
| V-117 | `onDelete` / `onUpdate` を指定する場合は `NO_ACTION` / `CASCADE` / `SET_NULL` / `SET_DEFAULT` のいずれかである | Error |
| V-118 | `indexes` を指定する場合、Index Mapを要素とする1件以上のArrayである                                            | Error |
| V-119 | 各Indexは `columns` を持ち、Index Column Mapを要素とする1件以上のArrayである                                   | Error |
| V-120 | 各Index Columnが `column` を持つ                                                                               | Error |
| V-121 | Index Columnの `order` を指定する場合は `ASC` / `DESC` のいずれかである                                        | Error |
| V-122 | Index Columnの `nulls` を指定する場合は `FIRST` / `LAST` のいずれかである                                      | Error |
| V-123 | `service` が存在し、`id` を持つ                                                                                | Error |
| V-124 | Service IDが `kebab-case` の命名規約を満たす                                                                   | Error |
| V-125 | DDL Source YAMLのファイル名がService IDと一致する                                                              | Error |

### 21.2 Element / Column Validation

Element / Column Validationでは、DDL ColumnとPhase 2 Element / Typeとの意味的な整合性を検証する。

<!-- prettier-ignore-start -->

| ID | Rule | Level |
| ----- | --- | ----- |
| V-126 | Columnの `element` がPhase 2で定義されたElementとして存在する | Error |
| V-127 | ColumnからElement / Type由来の `type` / `length` / `minLength` / `maxLength` / `regex` / `minimum` / `maximum` / `precision` / `scale` / `enum` / `format` を再定義またはOverrideしない | Error |
| V-128 | `createdAt` / `updatedAt` / `createdBy` / `updatedBy` をユーザー定義Columnとして定義しない | Error |
| V-129 | Columnの `constraints` に指定されたConstraintが参照Element / Typeで利用可能なConstraintである | Error |
| V-130 | Column固有ConstraintがElement由来Constraintと同一内容ではない | Error |
| V-131 | 比較可能なColumn固有ConstraintがElement由来Constraintを緩和しない | Error |
| V-132 | Column固有 `regex` がElement由来 `regex` と完全に同一ではない | Error |

<!-- prettier-ignore-end -->

Column固有ConstraintはElement由来ConstraintのOverrideではなく、追加Constraintとして扱う。

`minimum` / `maximum` / `minLength` / `maxLength` 等、
Constraint間の強弱を機械的に比較可能な場合は、Element由来Constraintと同一または緩いColumn固有ConstraintをValidation Errorとする。

`regex` については一般的な包含関係・強弱関係の解析を行わない。

Element由来Regexと異なるRegexは追加Constraintとして許可し、双方を満たすものとして扱う。

Element由来Regexと完全に同一のRegexは、意味のない重複としてValidation Errorとする。

参照ElementのTypeが `capability.generation: true` であっても、`generation` の指定は必須ではない。

Type Capabilityは、そのTypeを参照するColumnで値生成を定義可能であることを表す。
実際にそのColumnで値生成を行うかどうかは、DDL Sourceの `generation` によって決定する。

### 21.3 Default Validation

Default Validationでは、Column Defaultと参照Element / Typeとの整合性を検証する。

| ID    | Rule                                                                                                         | Level |
| ----- | ------------------------------------------------------------------------------------------------------------ | ----- |
| V-133 | Literal Defaultを指定した場合、参照Element Typeの `defaultAllowed.literal` が `true` である                  | Error |
| V-134 | Literal Defaultが参照Elementの論理Type / Format / Type固有Validation / Element Constraintを満たす            | Error |
| V-135 | Expression Defaultを指定した場合、そのExpressionが参照Element Typeの `defaultAllowed.expressions` に含まれる | Error |
| V-136 | 同一Columnに `default` と `generation` が同時に指定されていない                                              | Error |

Literal Defaultでは暗黙の型変換を行わない。

例えばINTEGERに対する `"0"` と `0`、BOOLEANに対する `"true"` と `true` は異なる値型として扱う。

任意のDatabase SQLをExpressionとして指定することはできない。
利用可能なExpressionは `defaultAllowed.expressions` により決定する。

### 21.4 Generation Validation

Generation Validationでは、Columnに指定されたGeneration方式と参照Element / Type、
およびCustom Generation Functionとの意味的な整合性を検証する。

<!-- prettier-ignore-start -->

| ID | Rule | Level |
| ----- | --- | ----- |
| V-137 | `generation` を指定したColumnの参照Element Typeが `capability.generation: true` である | Error |
| V-138 | `DATABASE` Generationが、参照Element Typeに対して対象Databaseで利用可能なDatabase Generationへ変換可能である | Error |
| V-139 | BuiltIn Generation名がARIADNEで定義されたBuiltIn Generationとして存在する | Error |
| V-140 | BuiltIn Generationが参照Element Typeに適用可能である | Error |
| V-141 | CUSTOM Generationの `custom` に指定されたCustom Generation Functionが対象Serviceに存在する | Error |
| V-142 | Custom Generation Functionが対象TableのRowを入力として受け取る | Error |
| V-143 | Custom Generation Functionの戻り値型がGeneration対象ColumnのDatabase物理型と一致する | Error |

<!-- prettier-ignore-end -->

Generation方式ごとの検証責務は以下とする。

```text
DATABASE
  → 対象DatabaseでそのTypeをDatabase Generationへ変換可能か

BUILTIN
  → BuiltIn Generationが存在するか
  → そのTypeに適用可能か

CUSTOM
  → Functionが存在するか
  → 対象Table Rowを入力とするか
  → 戻り値型が対象ColumnのDatabase物理型と一致するか
```

Custom Generation Functionの戻り値について、Prototype Phase 4では暗黙の型変換を許可しない。

例えば、PostgreSQL上で対象Columnが `bigint` として生成される場合、Custom Generation Functionの戻り値も `bigint` でなければならない。

Custom Generation Functionは対象Columnへ値を設定せず、生成値のみを返す。

対象Columnへの代入はARIADNEが生成するGeneration Triggerの責務とする。

また、Custom Generation Functionから同一Row内の他のGeneration Columnが生成した値へ依存することを許可しない。

Generation Column間の実行順序は保証しない。

### 21.5 Primary Key / Unique Constraint Validation

#### Primary Key

| ID    | Rule                                                                           | Level |
| ----- | ------------------------------------------------------------------------------ | ----- |
| V-144 | すべてのTableがPrimary Keyをちょうど1つ持つ                                    | Error |
| V-145 | `primaryKey.columns` に指定されたすべてのColumnが対象Tableに存在する           | Error |
| V-146 | 同一Primary Key内で同じColumnを重複指定しない                                  | Error |
| V-147 | Primary Keyを構成するすべてのColumnに明示的な `notNull: true` が指定されている | Error |

Primary KeyのColumn順序はSourceのArray記載順を保持する。

#### Unique Constraint

| ID    | Rule                                                                          | Level |
| ----- | ----------------------------------------------------------------------------- | ----- |
| V-148 | `uniqueConstraints[].columns` に指定されたすべてのColumnが対象Tableに存在する | Error |
| V-149 | 同一Unique Constraint内で同じColumnを重複指定しない                           | Error |
| V-150 | 同一Table内に同じColumn集合を持つUnique Constraintを複数定義しない            | Error |

Unique ConstraintのColumn順序はDDL生成時に保持する。

ただし、重複判定ではColumn順序を区別しない。

例えば `[a, b]` と `[b, a]` は、同一Column集合を持つUnique Constraintとして重複Errorとする。

### 21.6 Foreign Key Validation

Foreign Key Validationでは、Foreign Keyの参照関係およびReferential Actionの意味的整合性を検証する。

| ID    | Rule                                                                                                         | Level |
| ----- | ------------------------------------------------------------------------------------------------------------ | ----- |
| V-151 | Foreign KeyのLocal側 `columns` に指定されたすべてのColumnが対象Tableに存在する                               | Error |
| V-152 | `reference.table` に指定されたTableが同一Schema内に存在する                                                  | Error |
| V-153 | `reference.columns` に指定されたすべてのColumnが参照先Tableに存在する                                        | Error |
| V-154 | `reference.columns` が参照先TableのPrimary Keyまたは1つのUnique Constraintに対応する                         | Error |
| V-155 | `onDelete` / `onUpdate` が `SET_NULL` の場合、対象となるすべてのLocal ColumnがNULL許容である                 | Error |
| V-156 | `onDelete` / `onUpdate` が `SET_DEFAULT` の場合、対象となるすべてのLocal Columnに `default` が定義されている | Error |
| V-157 | 同一Foreign KeyのLocal側 `columns` に同じColumnを重複指定しない                                              | Error |
| V-158 | 同一Foreign Keyの `reference.columns` に同じColumnを重複指定しない                                           | Error |
| V-159 | 同一Table内に、Local Columnsの並びとReference Tableが同一であるForeign Keyを複数定義しない                   | Error |

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

### 21.7 Foreign Key Element一致 Validation

| ID    | Rule                                                                                     | Level |
| ----- | ---------------------------------------------------------------------------------------- | ----- |
| V-160 | Foreign KeyのLocal Columnと対応するReference Columnが同一のPhase 2 Elementを参照している | Error |

Database上の型が互換であっても、異なるElement間にForeign Keyを定義することは許可しない。

Foreign KeyにおけるElement一致は、
Database型の一致ではなく、同一の業務上の値を参照していることを保証するためのRuleとする。

### 21.8 Composite Foreign Key Validation

| ID    | Rule                                                                           | Level |
| ----- | ------------------------------------------------------------------------------ | ----- |
| V-161 | 複合Foreign KeyではLocal側 `columns` と `reference.columns` の要素数が一致する | Error |

複合Foreign KeyのColumn対応はArrayの記載順とする。

第n Local Columnは第n Reference Columnに対応する。

各対応Columnには Foreign Key Element一致Validationを適用する。

### 21.9 Index Validation

| ID    | Rule                                                                            | Level |
| ----- | ------------------------------------------------------------------------------- | ----- |
| V-162 | `indexes[].columns[].column` に指定されたすべてのColumnが対象Tableに存在する    | Error |
| V-163 | 同一Index内で同じColumnを複数回指定しない                                       | Error |
| V-164 | 同一Table内に同一Canonical Identityを持つIndexを複数定義しない                  | Error |
| V-165 | `notNull: true` のColumnに `nulls` を指定しない                                 | Error |
| V-166 | Primary Key / Unique ConstraintのBacking Indexと同一定義の通常Indexを定義しない | Error |

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

### 21.10 Naming Ruleとの関係

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

## 22. Phase 4で決定する事項

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
- BuiltIn SQLの仕様
- Custom SQLの登録方法
- Custom SQLと生成DDLの実行順序
- DDL出力ファイル構成
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
- [x] Default / Generationの解決結果をSample化する

### Step 6：PostgreSQL DDL仕様

- [x] PostgreSQL型変換規則を確定する
- [x] Schema SQLを定義する
- [x] Table / Column SQLを定義する
- [x] Default / Check Constraint変換を定義する
- [x] Column Comment生成規則を定義する
- [x] PK / Unique Constraint生成規則を定義する
- [x] Index生成規則を定義する
- [x] FK生成規則を定義する
- [x] DATABASE Generation生成規則を定義する
- [x] BUILTIN Generation（ULID）の生成規則を定義する
- [x] CUSTOM Generation / Generation Trigger生成規則を定義する
- [x] BuiltIn updatedAt / Audit Function / Triggerを定義する
- [x] Custom Function / Custom Trigger生成規則を定義する
- [x] Sample PostgreSQL DDLを作成する
- [x] Sample DDLをPostgreSQLで実行確認する

### Step 7：Custom SQL

- [ ] Custom SQLの登録方式を確定する
- [ ] DBMS別Custom SQLの扱いを確定する
- [ ] Custom SQLと生成DDLの実行順序を確定する

### Step 8：DDL Output

- [ ] DDLのファイル分割方針を検討する
- [ ] DDLの実行Phaseを確定する
- [ ] ファイル名による順序制御を利用するか決定する
- [ ] Sample出力構成を作成する

### Step 9：PostgreSQL Runtime

- [x] PostgreSQL Docker環境を作成する
- [x] pgwebを追加する
- [x] Order Management Sample DDLを適用する
- [x] Table / Constraint / Indexを確認する
- [x] BuiltIn Triggerを確認する
- [x] Custom SQLを確認する
- [x] pgwebからDatabaseを確認する

### Step 10：Phase 4 Completion

- [ ] Sample DDL YAMLを整理する
- [ ] Resolved DDL Sampleを整理する
- [ ] PostgreSQL Sample SQLを整理する
- [x] PostgreSQL Runtime検証を完了する
- [ ] Phase 4ドキュメントを実証結果に合わせて更新する
- [ ] Coreへ引き継ぐGenerator / Validator仕様を整理する
- [ ] 未決事項が残っていないことを確認する
- [ ] Phase 4 COMPLETE

---

## Prototype Phase 4: DDL YAML — IN PROGRESS
