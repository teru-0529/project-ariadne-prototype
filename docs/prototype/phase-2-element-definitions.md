# Prototype Phase 2：項目定義 YAML

**Status:** COMPLETE

---

## 1. Phase 2 の目的

Prototype Phase 2 では、Project ARIADNE
が管理する「項目定義」のモデルを確定する。

Prototype では項目定義を編集する UI、YAML の読み書き機能、Validator
は実装しない。 この Phase の目的は、以下を設計として確定し、正本 YAML
を手作業で作成して妥当性を確認することである。

- 項目定義で何を管理するか
- 論理 Type と Element の責務
- 正本 YAML の構造（共通ヘッダー `formatVersion` / `updatedAt` / `domain` / `kind` を含む）
- 配置・命名
- YAML 間の整合性ルール
- OAS / DDL への受け渡し方針

読み書き・編集・Validation のツール化は Core で検討する。

## 2. Source of Truth と配置

Phase 0 で決定した責務分離を維持する。

``` text
src/
├─ elements/
│  ├─ types.yaml
│  └─ elements.yaml
├─ api/
└─ ddl/

templates/
dist/
runtime/
```

項目定義の正本は以下とする。

``` text
src/elements/types.yaml
src/elements/elements.yaml
```

各領域の責務は以下のとおり。

| 領域 | 責務 |
| --- | --- |
| `src/elements/` | Type / Element による共通項目定義の正本 |
| `src/api/` | API 定義の ARIADNE Source および API 生成処理で使用する正本資材 |
| `src/ddl/` | Database / DDL 定義の ARIADNE Source および DDL 生成処理で使用する正本資材 |
| `templates/` | Task 管理アプリ等、Prototype で利用するテンプレート／アプリデータ側の資材 |
| `dist/` | ARIADNE が生成した成果物 |
| `runtime/` | SQLite 等の実行時データ |

`types.yaml` と `elements.yaml`
は独立した仕様ではなく、**1つの項目定義モデルを2ファイルに分割したもの**として扱う。

## 3. 基本モデルと責務

### 3.1 Type

`types.yaml` は論理型体系の正本であり、「値は何者か」を定義する。

Type は以下を管理する。

- 論理型の名称・説明
- API 上の識別子として利用可能か
- 利用可能な constraint
- Type 固有の validation
- DDL Default として利用可能な値の種類
- PostgreSQL へのマッピング
- OpenAPI へのマッピング
- 推奨命名

`validation` は、その Type を使用するすべての Element に共通して適用される、Type 固有の不変な形式制約を定義する。

Element ごとに指定可能な `constraints` とは区別する。

### 3.2 Element

`elements.yaml`
は項目定義の正本であり、「業務上どんな意味を持つ項目か」を定義する。

Element は以下を管理する。

- 不変識別子となる Element キー
- 名称
- Type
- 説明
- example
- API 上の識別子として利用可能か（必要な場合のみ）
- Element 固有の constraint 値
- ENUM の値集合

### 3.3 関係

``` text
Type
  1
  │
  N
Element
  │
  └─ type = ENUM の場合のみ
       1..N EnumValue
```

Type / Element / EnumValue は Map を基本構造とする。

## 4. キーの不変性

Element キーおよび EnumValue キーは、不変識別子として扱う。

ここでいう「不変」とは、永久に削除できないという意味ではない。

> 同一概念として存在している間、その識別子を rename しない。

例えば `customerId` を `clientId` に変更したい場合、rename
ではなく以下として扱う。

``` text
DELETE customerId
CREATE clientId
```

削除および新規作成は可能だが、両者は別の Element とする。

EnumValue も同じ原則とする。

## 5. 標準 Type

ARIADNE の標準 Type は以下の14種類とする。

``` text
識別子系
├─ PROVIDED_ID
├─ SEQUENCE_ID
└─ GENERATED_ID

文字列系
├─ FIXED_STRING
├─ STRING
└─ TEXT

数値系
├─ INTEGER
└─ DECIMAL

真偽・選択系
├─ BOOLEAN
├─ ENUM
└─ CODE

日時系
├─ DATE
├─ TIME
└─ DATETIME
```

14種類をアプリケーションにハードコードすることを前提とはしない。 Type は
`types.yaml` の Map として定義し、将来的な追加を可能とする。

### 5.1 識別子系

- `PROVIDED_ID`：外部または人が指定する文字列識別子
- `SEQUENCE_ID`：システム／DB により自動採番される数値識別子
- `GENERATED_ID`：システムが生成する文字列識別子

Type の `identifier: true` は、その Type を使用する Element が
API 上の識別子（例：Path Parameter）として利用可能であることを示す。

また、Type の `identifier` が `false` の場合でも、
個別の Element に `identifier: true` を指定することで、その Element を API 上の識別子として利用可能にできる。

実効的な identifier は以下とする。

effective identifier = Type.identifier OR Element.identifier

Element 側の `identifier` は必要な場合のみ `true` を指定する。 `false` は記述しない。

単独主キー・単独一意キーであることは意味しない。複合識別子の一要素として利用してもよい。

PK / UNIQUE / 複合キー等のデータベース上の制約は DDL 側で扱う。

### 5.2 Default 利用可否

Type は、DDL Column の Default として利用可能な値の種類を `defaultAllowed` により定義する。

`defaultAllowed` は以下の構造とする。

```yaml
defaultAllowed:
  literal: true
  expressions:
    - CURRENT_DATE
```

- `literal`：Literal Default を利用可能かを示す
- `expressions`：その Type で利用可能な ARIADNE BuiltIn Default Expression を列挙する

`defaultAllowed` は Type 固有の性質とし、Element 側では Override できない。

標準 Type の `defaultAllowed` は以下とする。

| Type | literal | expressions |
| --- | --- | --- |
| `PROVIDED_ID` | `false` | `[]` |
| `SEQUENCE_ID` | `false` | `[]` |
| `GENERATED_ID` | `false` | `[]` |
| `FIXED_STRING` | `true` | `[]` |
| `STRING` | `true` | `[]` |
| `TEXT` | `true` | `[]` |
| `INTEGER` | `true` | `[]` |
| `DECIMAL` | `true` | `[]` |
| `BOOLEAN` | `true` | `[]` |
| `ENUM` | `true` | `[]` |
| `CODE` | `false` | `[]` |
| `DATE` | `true` | `CURRENT_DATE` |
| `TIME` | `true` | `CURRENT_TIME` |
| `DATETIME` | `true` | `CURRENT_DATETIME` |

識別子系 Type は値の提供・生成主体を別に持つため、Literal Default を許可しない。

`CODE` は値集合を項目定義の外部で管理するため、Literal Default を許可しない。

`ENUM` は値集合自体が項目定義の仕様として固定されるため、Literal Default を許可する。

`CURRENT_DATE` / `CURRENT_TIME` / `CURRENT_DATETIME` は ARIADNE が定義する BuiltIn Default Expression を設定可能とする。

BuiltIn Default Expression の Database 固有表現は `types.yaml` では定義しない。
PostgreSQL / SQLite 等への変換は DDL 生成側の責務とする。

### 5.3 ENUM と CODE

ENUM と CODE は明確に区別する。

#### **ENUM**

- 有限・離散
- 値集合が仕様として固定
- 値集合を Element 自身が所有
- 値追加・削除は仕様変更
- PostgreSQL enum / OpenAPI enum として利用

#### **CODE**

- 有限・離散
- 値集合が運用中に変化し得る
- 値集合は項目定義の外部で管理
- 主にマスタ等から取得する選択値

原則は以下とする。

> ENUM の値変更 = 仕様変更\
> CODE の値変更 = データ変更

ENUM では `NORMAL = 0` のような別 value を持たせない。 EnumValue
キー自体を実値として扱う。

### 5.4 TIME

`TIME` は、タイムゾーンを持たない時刻値を表す。

形式は `HH:mm:ss` 固定とし、Type 固有の `validation` により形式を定義する。

```text
00:00:00
15:30:00
23:59:59
```

TIME の形式は Element 固有の制約ではないため、各 Element の constraints には定義しない。

### 5.5 DATETIME

`DATETIME` は、UTC オフセットを持つ日時値を表す。

OpenAPI では `string / date-time` として扱う。

例：

```text
2023-10-01T12:00:00+09:00
2023-10-01T03:00:00Z
```

DATETIME は、UTC オフセットを持つ date-time 形式とする。UTC は Z による表現も許容する。
ARIADNE が定義する example では、原則として ±HH:mm 形式の UTC オフセットを使用する。

## 6. Constraint

Type ごとの constraint は以下とする。

| Type | Required | Optional |
| --- | --- | --- |
| `PROVIDED_ID` | `maxLength` | `minLength`, `regex` |
| `SEQUENCE_ID` | なし | なし |
| `GENERATED_ID` | `maxLength` | `minLength`, `regex` |
| `FIXED_STRING` | `length` | `regex` |
| `STRING` | `maxLength` | `minLength`, `regex` |
| `TEXT` | なし | `minLength`, `regex` |
| `INTEGER` | なし | `minimum`, `maximum` |
| `DECIMAL` | `precision`, `scale` | `minimum`, `maximum` |
| `BOOLEAN` | なし | なし |
| `ENUM` | なし | なし |
| `CODE` | `maxLength` | `minLength`, `regex` |
| `DATE` | なし | なし |
| `TIME` | なし | なし |
| `DATETIME` | なし | なし |

`precision / scale` は10進数の桁構造を表す。

`minimum / maximum` は業務上許容する値域を表す。

`SEQUENCE_ID` は自動採番を前提とするため、Element 側の採番方法等の
constraint は持たせない。

## 7. PostgreSQL / OpenAPI マッピング

必須マッピングは PostgreSQL と OpenAPI のみとする。

Java / Go / TypeScript / Python 等の言語型は `types.yaml`
では管理せず、OpenAPI Generator 等の生成側に委ねる。

| Type | PostgreSQL | OpenAPI |
| --- | --- | --- |
| `PROVIDED_ID` | `varchar(n)` | `string` |
| `SEQUENCE_ID` | `bigint` | `integer / int64` |
| `GENERATED_ID` | `varchar(n)` | `string` |
| `FIXED_STRING` | `varchar(n)` + `CHECK (LENGTH(column) = n)` | `string` |
| `STRING` | `varchar(n)` | `string` |
| `TEXT` | `text` | `string` |
| `INTEGER` | `bigint` | `integer / int64` |
| `DECIMAL` | `numeric(p,s)` | `number` |
| `BOOLEAN` | `boolean` | `boolean` |
| `ENUM` | `enum` | `string + enum` |
| `CODE` | `varchar(n)` | `string` |
| `DATE` | `date` | `string / date` |
| `TIME` | `time` | `string` |
| `DATETIME` | `timestamp with time zone` | `string / date-time` |

`SEQUENCE_ID` は値としてのDatabase型を `bigint` とする。

`SEQUENCE_ID` 自体は、そのElementを参照するすべてのColumnが採番主体であることを意味しない。

Database上でそのColumn自身が連番を生成するかどうかはDDL側の責務とし、DDL Columnの `sequence` により定義する。

これにより、同一の `SEQUENCE_ID` Elementを採番元ColumnとForeign Key等の従属Columnの双方で利用可能とする。

`FIXED_STRING` は PostgreSQL の `char(n)` へ直接マッピングせず、
`varchar(n)` と固定長を保証する `CHECK Constraint` の組合せとして生成する。

概念例：

```sql
customer_id varchar(6)
  CHECK (LENGTH(customer_id) = 6)
```

`char(n)` 固有の空白Paddingの挙動を固定長制約の意味に含めず、文字列長が指定値と一致することを明示的に保証する。

`DECIMAL` の OpenAPI `format` に `double` は指定しない。 ARIADNE の
DECIMAL は `precision / scale`
を持つ10進数であり、特定言語の倍精度浮動小数点型へのマッピングは生成側の責務とする。

## 8. `types.yaml`

トップレベルは以下とする。

``` yaml
formatVersion: "1.0"
updatedAt: "2026-08-25T23:00:00+09:00"
domain: elements
kind: types

types:
  ...
```

Type の基本構造は以下とする。

```yaml
TYPE_KEY:
  name: ...
  description: ...
  identifier: false

  constraints:
    required: []
    optional: []

  defaultAllowed:
    literal: true
    expressions: []

  validation:          # 必要な場合のみ
    regex: ...

  postgresql:
    type: ...

  openapi:
    type: ...
    format: ...        # 必要な場合のみ

  recommendations: []
```

`constraints.required / optional`、`defaultAllowed.expressions` および `recommendations` は0件の場合も空リストを明示する。

`defaultAllowed` はすべての Type で必須とする。

`defaultAllowed.literal` は boolean とし、Literal Default の利用可否を示す。

`defaultAllowed.expressions` は、その Type で利用可能な ARIADNE BuiltIn Default Expression の一覧とする。

`validation` は Type 固有の形式制約が存在する場合のみ記述する。Prototype では `regex` を使用する。

### 8.1 Naming Recommendations

標準の推奨命名は以下とする。

``` text
PROVIDED_ID   → *Id, *No
SEQUENCE_ID   → *Id, *No
GENERATED_ID  → *Id, *No

BOOLEAN      → is*, has*, can*, should*
ENUM         → *Type
CODE         → *Code
DATE         → *Date
TIME         → *Time
DATETIME     → *DateTime
```

その他は `recommendations: []` とする。

Naming Recommendation は制約ではなく推奨であり、違反は Error ではなく
Warning とする。

## 9. `elements.yaml`

トップレベルは以下とする。

``` yaml
formatVersion: "1.0"
updatedAt: "2026-08-25T23:00:00+09:00"
domain: elements
kind: elements

elements:
  ...
```

Element の基本構造は以下とする。

``` text
ELEMENT_KEY
├─ name           必須
├─ type           必須
├─ description    必須
├─ example        必須
├─ identifier     必要な場合のみ（true のみ）
├─ constraints    必要な場合のみ
└─ values         ENUM のみ
```

`identifier: true` は、Type 側の `identifier` が `false` であっても、
その Element を API 上の識別子として利用可能にする場合に指定する。

例：

``` yaml
customerId:
  name: 得意先ID
  type: PROVIDED_ID
  description: 得意先を一意に識別するID
  constraints:
    maxLength: 10
    minLength: 1
    regex: "^[A-Z0-9]+$"
  example: "C001"

orderDetailNo:
  name: 明細番号
  type: INTEGER
  identifier: true
  description: 受注内で明細を識別する番号
  constraints:
    minimum: 1
  example: 1
```

constraint が0件の場合、`constraints` 自体を記述しない。

Element 側には `recommendations` を持たせない。

### 9.1 example

`example` は必須とし、文字列化して保持するのではなく、論理 Type
に対応する値型で記述する。

``` yaml
# STRING
example: "ABC"

# INTEGER
example: 100

# DECIMAL
example: 12.5

# BOOLEAN
example: true

# DATE
example: "2026-08-25"
```

`example` は論理 Type の型・形式に適合し、Element に設定された
constraint も満たさなければならない。

### 9.2 ENUM values

ENUM は以下の構造とする。

``` yaml
status:
  name: 状態
  type: ENUM
  description: タスクの状態
  example: "TODO"

  values:
    TODO:
      name: 未着手

    IN_PROGRESS:
      name: 進行中

    DONE:
      name: 完了
```

`values` は Map とし、1件以上必須とする。

EnumValue は以下を持つ。

``` text
ENUM_VALUE_KEY
├─ name           必須
└─ description    任意
```

EnumValue キーは `UPPER_SNAKE_CASE`
の不変識別子であり、そのキー自体を実値として利用する。

## 10. `formatVersion`

`formatVersion` はファイル単体のバージョンではない。

> `types.yaml` と `elements.yaml` で構成される Item Definition Format
> 全体のバージョンである。

したがって両ファイルの `formatVersion` は必ず一致する。

``` text
Item Definition Format 1.0
│
├─ types.yaml
│    formatVersion: "1.0"
│
└─ elements.yaml
     formatVersion: "1.0"
```

`updatedAt`
は各正本ファイルの更新日時であるため、ファイルごとに独立する。

OAS / DDL は別のモデルであり、将来それぞれ `formatVersion`
を持つ場合でも、Item Definition Format
と同じバージョンにする必要はない。

### 10.1 `domain` / `kind`

ARIADNE が管理する正本 YAML は、ファイル名や配置場所だけに依存せず、ファイル自身が「どの設計領域の、何の定義か」を宣言できるようにする。

- `domain`：どの設計領域に属するかを表す
- `kind`：その領域の中で何の定義かを表す

Phase 2 の項目定義では以下とする。

| ファイル | `domain` | `kind` |
| --- | --- | --- |
| `types.yaml` | `elements` | `types` |
| `elements.yaml` | `elements` | `elements` |

Phase 2 の正本 YAML の共通ヘッダーは以下の形となる。

``` yaml
formatVersion: "1.0"
updatedAt: "2026-08-25T23:00:00+09:00"
domain: elements
kind: types     # types.yaml の場合
```

`elements.yaml` の場合は `kind: elements` とする。

`domain` / `kind` は機械向けの自己記述情報であり、ファイル名やディレクトリ構成を置き換えるものではない。配置・命名とあわせて、正本 YAML の種別を明確にするために使用する。

OAS / DDL は別のモデルであり、将来それぞれ `formatVersion`
を持つ場合でも、Item Definition Format
と同じバージョンにする必要はない。

## 11. Validation Rules

Prototype では Validator を実装しないが、正しい状態を仕様として定義する。

| ID | Validation Rule | 判定 |
| --- | --- | --- |
| `V-01` | `types.yaml.formatVersion` と `elements.yaml.formatVersion` が一致する | Error |
| `V-02` | `formatVersion` が ARIADNE の対応する Item Definition Format である | Error |
| `V-03` | Element の `type` が `types.yaml` に存在する | Error |
| `V-04` | Type の `required` constraint を Element がすべて持つ | Error |
| `V-05` | Element は Type の `required / optional` に存在しない constraint を持たない | Error |
| `V-06` | constraint の値そのものが妥当である | Error |
| `V-07` | `example` が論理 Type の型・形式、Type 固有の `validation`、Element の constraint を満たす | Error |
| `V-08` | `type: ENUM` の Element は `values` を持つ | Error |
| `V-09` | ENUM の `values` が1件以上存在する | Error |
| `V-10` | ENUM 以外の Element は `values` を持たない | Error |
| `V-11` | EnumValue キーが `UPPER_SNAKE_CASE` に適合する | Error |
| `V-12` | ENUM の `example` が `values` に存在するキーである | Error |
| `V-13` | Element キーが `camelCase` に適合する | Error |
| `V-14` | Type キーが `UPPER_SNAKE_CASE` に適合する | Error |
| `V-15` | Element キーが Type の `recommendations.naming` に適合する | Warning |
| `V-16` | Element / EnumValue キーを同一概念のまま rename しない | 運用ルール |
| `V-17` | `updatedAt` が各ファイルに存在する | Error |
| `V-18` | `updatedAt` が ISO 8601 形式である | Error |
| `V-19` | `domain` が存在し、値が `elements` である | Error |
| `V-20` | `kind` が存在し、定義種別と一致する（`types` / `elements`） | Error |
| `V-21` | Type の `validation` の定義内容が妥当である | Error |
| `V-22` | Element の `identifier` は、指定する場合 true である | Error |
| `V-23` | Type が `defaultAllowed` を持つ | Error |
| `V-24` | `defaultAllowed.literal` が boolean である | Error |
| `V-25` | `defaultAllowed.expressions` が定義済みの ARIADNE BuiltIn Default Expression のみを持つ | Error |
| `V-26` | Element は `defaultAllowed` を定義しない | Error |

### 11.1 Constraint Validation

`V-06` では少なくとも以下を確認する。

``` text
length > 0

minLength >= 0
maxLength > 0
minLength <= maxLength

precision > 0
scale >= 0
scale <= precision

minimum <= maximum

regex = 有効な正規表現
```

### 11.2 Type Validation

Type に `validation` が定義されている場合、その内容自体を Validation する。

Prototype では以下を確認する。

``` text
validation.regex = 有効な正規表現
```

Type 固有の `validation` は、その Type を使用するすべての Element の Example Validation に適用する。

Element の実効 identifier は、Type.identifier または Element.identifier のいずれかが true の場合に true とする。

### 11.3 Default Allowed Validation

`defaultAllowed` は、DDL Column の Default として利用可能な値の種類を Type ごとに定義する。

``` text
defaultAllowed.literal
  true  → Literal Default を利用可能
  false → Literal Default を利用不可

defaultAllowed.expressions
  → 利用可能な ARIADNE BuiltIn Default Expression
```

`defaultAllowed.expressions` に指定できる値は、ARIADNE が定義する BuiltIn Default Expression に限る。

Prototype では以下を定義する。

``` text
CURRENT_DATE
CURRENT_TIME
CURRENT_DATETIME
```

Element は `defaultAllowed` を Override できない。

DDL Column に指定された Default が実際に利用可能かどうかの Validation は、
参照 Element の Type に定義された `defaultAllowed` を使用して DDL 側で行う。

### 11.4 Example Validation

型の基本対応は以下とする。

``` text
PROVIDED_ID / GENERATED_ID / FIXED_STRING / STRING / TEXT / CODE
  → string

SEQUENCE_ID / INTEGER
  → integer

DECIMAL
  → number

BOOLEAN
  → boolean

ENUM
  → string + values 内に存在

DATE
  → string + date 形式

TIME
  → string + Type 固有 validation（HH:mm:ss）に適合

DATETIME
  → string + date-time 形式

```

さらに Element 固有の
`minimum / maximum / minLength / maxLength / regex` 等も満たすこと。

### 11.5 Error / Warning / Operation Rule

``` text
ERROR
  正本として成立しない。
  将来の生成処理では OAS / DDL 生成を停止する。

WARNING
  正本としては成立する。
  推奨規約から外れていることを通知する。

OPERATION RULE
  単一ファイルの静的 Validation だけでは完全に判定できない。
  Core の操作制御等で担保する。
```

## 12. Core での Validation 実装方針（将来）

Validation は特定画面の機能ではなく、ARIADNE
の共通能力として実装することを想定する。

``` text
validator package
    │
    ├─ アプリ内部
    │    └─ 登録・保存時
    │
    ├─ CLI
    │    └─ ariadne validate
    │
    └─ Git hook / CI
         └─ commit / push / PR
```

Validator 本体はファイル I/O から分離し、Domain Model を入力として
Validation する。

``` text
YAML Loader
    ↓
Domain Model
    ↓
Validator
```

これにより以下で同じ Validator を利用できる。

- UI での登録・保存
- 保存前の未永続データ
- 手作業 YAML の確認
- Git pre-commit / pre-push
- CI / PR
- Go の単体テスト

Go で実装する場合、Validator は共通 package とし、CLI は薄い入口とする。
CLI フレームワークとして Cobra 等を利用することを想定できる。

``` text
Cobra CLI
    ↓
Loader
    ↓
Domain Model
    ↓
Validator Package
```

CLI は将来的に以下のような構成へ拡張可能とする。

``` text
ariadne validate
ariadne generate openapi
ariadne generate ddl
ariadne version
```

Prototype Phase 2 ではこれらを実装しない。

## 13. Prototype 検証用 Element

Prototype では Task を題材に、複数 Type と constraint
を実際に使用してモデルを検証した。

``` text
taskId          → GENERATED_ID
title           → STRING
description     → TEXT
status          → ENUM
priority        → INTEGER
recordDateTime  → DATETIME
```

主な検証内容は以下。

- STRING の `maxLength`
- TEXT の `minLength`
- ENUM の Map / example / values
- INTEGER の `minimum / maximum`
- DATETIME の offset 付き date-time
- `example` と constraint の整合性

`recordDateTime` は、レコードの生成・更新時にシステムが管理する共通日時
Element とする。

Phase 4 では同一 Element を DDL 上で `created_at` / `updated_at`
等の複数カラムとして利用する設計を検証する。

Prototype の手作業 Validation 結果は以下。

``` text
ERROR   : 0
WARNING : 1

WARNING:
status → ENUM naming recommendation "*Type"

RESULT  : VALID
```

`status` は業務上自然な名称であり、Naming Recommendation
は制約ではないため変更しない。

## 14. OAS / DDL への受け渡し

### 14.1 Phase 3：API Definition / OAS

Phase 3 の API Definition は、Element Definition から共通意味を継承する。

主な入力は以下。

- Element キー
- `name`
- `type`
- `description`
- `example`
- constraints
- identifier
- Type 固有 validation
- ENUM values
- OpenAPI mapping

Type 固有の `validation` は、OpenAPI で表現可能な場合、OpenAPI Schema の対応する制約へ変換する。

Phase 3 では Element を参照し、
Resource / Parameter / Request / Response / Operation 等のAPI 利用文脈を定義する。

Resource Property / Variant / API Usage / Parameter Definition 等では、
それぞれの責務に応じて `name` / `description` / `example` 等のコンテキスト情報を追加・Overrideできる。

ARIADNE 標準として、

- Request / Response Body の Property は camelCase
- Path / Query Parameter の物理名は snake_case

とする。

項目定義から OpenAPI 完成形を直接生成するのではない。

```text
Phase 2 Element Definition
        ↓
Phase 3 ARIADNE API Source
        ↓
Raw Model
        ↓
Resolved Model
        ↓
OpenAPI 3.1
```

Phase 2 は共通の値定義を提供し、Phase 3 が API 文脈を付与したうえで OpenAPI 3.1 を生成する。

### 14.2 Phase 4：DDL

DDL 側も項目定義から共通意味を継承する。

主な入力は以下。

- Element キー
- `name`
- `type`
- `description`
- constraints
- ENUM values
- `defaultAllowed`
- PostgreSQL mapping

DDL 固有の以下の情報は DDL 側で管理する。

- schema / table
- column 名
- NOT NULL
- DEFAULT
- PK
- FK
- UNIQUE
- INDEX
- テーブル固有 CHECK
- 同一 Element を複数カラムとして利用する場合の役割・別名

DDL Column の `default` は DDL 側で定義する。

ただし、その Column で利用可能な Default の種類は、参照 Element の Type に定義された `defaultAllowed` に従う。

Literal Default を指定する場合は `defaultAllowed.literal`、
BuiltIn Default Expression を指定する場合は `defaultAllowed.expressions` を参照する。

BuiltIn Default Expression の PostgreSQL / SQLite 等への具体的な変換は DDL 生成側の責務とする。

基本原則は以下とする。

> Phase 2 は共通意味を定義する。\
> Phase 3 / Phase 4 は、その意味を利用文脈に合わせて拡張する。

``` text
types.yaml + elements.yaml
        │
        ├──→ OAS 正本
        │      + API 文脈
        │
        └──→ DDL 正本
               + DB 文脈
```

Phase 3 / Phase 4 で項目定義そのものを再定義しない。

## 15. Phase 2 Completion

Prototype Phase 2 では以下を完了した。

- [x] 項目定義の責務整理
- [x] 概念モデル・Type 体系
- [x] `types.yaml` 詳細構造
- [x] `elements.yaml` 詳細構造
- [x] 整合性ルール
- [x] 配置・ファイル名
- [x] Prototype 用 YAML の手作成
- [x] 手作業 Validation
- [x] Phase 3 / Phase 4 への受け渡し整理

---

## Prototype Phase 2: 項目定義 YAML — COMPLETE
