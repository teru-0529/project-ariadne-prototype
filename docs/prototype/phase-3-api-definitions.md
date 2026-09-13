# Prototype Phase 3：API定義 YAML

**Status:** IN PROGRESS

> 本ドキュメントは Prototype Phase 3
> の途中時点における設計決定事項を記録する。\
> Phase 3 完了までに、Resolved Model、Validation Rule、OAS 3.1
> Generation 等の設計進展に合わせて更新する。

------------------------------------------------------------------------

## 1. Phase 3 の目的

Prototype Phase 3 では、Project ARIADNE が管理する **API 定義 YAML**
のモデルを確定する。

Phase 2 で定義した Element を値定義の正本として利用し、Phase 3
ではそれらを組み合わせて API 上の Resource、Parameter、Request /
Response、Operation 等を表現する。

Phase 3 の主な目的は以下とする。

- API 定義の Source of Truth となる ARIADNE YAML 表記を確定する
- Service / Parameter / Resource / SubResource / Variant / Action /
    Custom の責務を整理する
- Phase 2 Element を参照して API Schema を構成する
- OpenAPI 3.1 を生成するために必要な情報を定義する
- File Validation / API Validation / Service Validation / OpenAPI
    Validation の責務を分離する
- ARIADNE Source から Raw Model / Resolved Model
    を生成するモデル変換を定義する

Phase 3 では **値そのものの型・桁・範囲等を再定義しない**。\
Scalar Value の定義は Phase 2 Element の責務とする。

Prototype では設計を先に確定し、正本 YAML
と中間モデルを手作業で作成して妥当性を確認する。Validator / Generator
の実装は後続工程で行う。

------------------------------------------------------------------------

## 2. Source of Truth と配置

Phase 3 の Source of Truth は `src/api/services/` 配下の ARIADNE YAML
とする。

Service ごとにディレクトリを分け、Service 固有の Parameter
Definition、Resource、Action、Custom、補足ドキュメントをその配下で管理する。

現時点の `order-management` Service の構成例を以下に示す。

```text
src/
└─ api/
   └─ services/
      ├─ order-management/
      │  ├─ service.yaml
      │  ├─ parameters.yaml
      │  ├─ docs/
      │  │  └─ order-create.md
      │  ├─ order/
      │  │  ├─ main.yaml
      │  │  ├─ subresources/
      │  │  │  └─ order-detail.yaml
      │  │  └─ actions/
      │  │     ├─ order-shipment.yaml
      │  │     └─ cancellation.yaml
      │  └─ customs/
      │     └─ export.yaml
      └─ task-management/
         └─ ...
```

### 2.1 配置原則

- Service ディレクトリ名は `service.yaml` の Service ID と一致させる
- `service.yaml` は Service 自身を定義する
- `parameters.yaml` は **Service 単位**の Parameter Definition
    を定義する
- Main Resource は Resource 用ディレクトリの `main.yaml` に定義する
- SubResource / Action は Main Resource 配下で管理する
- Custom は特定 Main Resource の子ではなく **Service 配下**で管理する
- `docs/` は Service 固有の `externalDocs` 用補足文書の配置候補とする
- `externalDocs.url` の相対パスは **Service Root**
    を基準として解決する

`externalDocs` の Markdown は Service Root 配下の `docs/` を正本とする。
OpenAPI / API Document 生成時に HTML へ変換し、 Service 単位の OAS 成果物ディレクトリへ生成する。

------------------------------------------------------------------------

## 3. 基本モデルと責務

### 3.1 Service

Service は API の管理単位を表す。

Service は以下を管理する。

- Service ID
- 表示名
- 説明
- Service 共通 Request Header

Service ID は kebab-case とする。

例:

``` yaml
service: order-management
name: 受注サービス
description: |
  受注および顧客情報を管理するサービス。
requestHeaders:
  - traceId
```

### 3.2 Parameter Definition

Parameter Definition は、Element を API Parameter として利用する際の
Service 固有の定義を表す。

Parameter Definition は以下を管理できる。

- 参照する Element
- description
- example
- 特殊な Header 名を使用する場合の `headerName`

Parameter Definition 自身は Path / Query / Header
のどこで利用されるかを持たない。利用場所は API 側から決定する。

同じ Phase 2 Element を異なる Service が異なる Parameter Definition
として利用してよい。

### 3.3 Resource

Resource は API 上のデータ構造を表す。

Resource は Property を持ち、各 Property
は以下のいずれか一つを参照する。

- Element
- Resource
- Array

Main Resource は `parent` を持たない。

### 3.4 SubResource

SubResource は親 Resource に従属する Resource を表す。

Source では `kind: subresource` と `parent` によって明示する。

例:

``` yaml
kind: subresource
resource: OrderDetail
parent: Order
```

Raw Model では Main Resource / SubResource を同じ Resourceとして扱い、`parentRef` の有無で区別する。

### 3.5 Variant

Variant は Resource の派生表現を定義する。

Variant は Resource を複製して別 Schema
を直接定義するのではなく、以下の差分で表現する。

- `include`
- `exclude`
- `add`
- `overrides`

`include` と `exclude` は同一 Variant で併用しない。

Variant から別 Variant を継承することはしない。

### 3.6 Parent Variant

Parent Variant は、SubResource 側から親 Resource に Variant を追加するための定義である。

ARIADNE YAML では `parentVariants` として記述する。

典型例として、`OrderDetail` から親 `Order` に `WithDetails` Variantを追加する。

``` yaml
parentVariants:
  WithDetails:
    add:
      details:
        array:
          resource: OrderDetail
        name: 受注明細
```

Source / Raw Model では SubResource 側に保持し、親 Resource への統合は Resolve 時に行う。

### 3.7 Action

Action は業務上の操作を、verb Path ではなく
**生成・受付される業務上の名詞 Resource** として表現する。

例:

- `OrderShipment`
- `Cancellation`

Action は必要に応じて Action-local Resource を持つ。

Action-local Resource は Source / Raw Model 上では Action のローカル定義として管理する。

ただし Resource / SubResource / Local Resource の名前は、Resolved Schema 名の一意性を保証するため Service 内で一意とする。

Raw Model 上の Local Resource 参照は API Validation で解決し、Resolved Model では Service 内で一意な Schema として扱う。

Action 名と Action-local Resource 名が同一でも許容する。両者はスコープが異なるため曖昧とはみなさない。

### 3.8 Custom

Custom は通常の Resource / SubResource / Action
のモデルでは自然に表現できない API の escape hatch とする。

Custom では以下を許容する。

- verb 的な Path
- 任意の深さの Path
- Resource 横断的な API
- Custom-local Resource
- 明示的な Success Response Status
- 0件以上の明示的な Error Response
- Response ごとの description
- Response Body の Resource / Resource Array / Body なし

Custom は Service スコープとし、特定 Main Resource には従属させない。

### 3.9 API / Operation

API は logical entry、Path、HTTP Method の階層で記述する。

例:

``` yaml
api:
  member:
    path: /orders/{received_order_no}

    get:
      summary: 受注取得
      response:
        resource: Order

    put:
      summary: 受注更新
      request:
        resource: Order
      response:
        resource: Order
```

logical entry は人間が API を整理するための論理名であり、HTTP 上の一意性は最終的に `(path, method)` で判定する。

Operation は OpenAPI 上の分類を示す `tag` を1件持つ。

`tag` の決定方法は API の定義元により異なる。

- Main Resource の API：その Main Resource 名を使用する
- SubResource の API：親 Main Resource 名を使用する
- Action の API：所属する Main Resource 名を使用する
- Custom の API：Operation に `tag` を明示する

ARIADNE では 1 Operation に対して複数 Tag を許可しない。OpenAPI 生成時は `tag` を `tags` 配列へ変換する。

------------------------------------------------------------------------

## 4. 命名規約

現時点の命名規約を以下とする。

| 対象 | 規約 | 例 |
| --- | --- | --- |
| Service | kebab-case | `order-management` |
| Resource | PascalCase | `Order` |
| SubResource | PascalCase | `OrderDetail` |
| Variant | PascalCase | `Summary` |
| Action | PascalCase | `OrderShipment` |
| Custom | PascalCase | `ExportOrders` |
| Local Resource | PascalCase | `ExportCondition` |
| Property | camelCase | `customerId` |
| Parameter Definition | camelCase | `receivedOrderNo` |
| API logical entry | camelCase | `createOrder` |
| HTTP Path Parameter | snake_case | `received_order_no` |
| HTTP Query Parameter | snake_case | `customer_id` |

Path Token は semantic key の snake_case 表現とし、Raw Model 生成時に
camelCase の Parameter Definition key へ正規化する。

------------------------------------------------------------------------

## 5. 継承・Override の原則

Phase 3 では、Phase 2 Element を値定義の最下層として利用する。

概念的な階層は以下とする。

``` text
Phase 2 Element
    │
    ├─ Resource Property
    │      ↓
    │    Variant
    │      ↓
    │    API Usage
    │
    └─ Parameter Definition
```

Resource Property / Variant / API Usage、および Parameter Definition では、利用コンテキストに応じた情報を追加・Overrideできる。

### 5.1 Override 可能な情報

主に以下をコンテキスト情報として扱う。

- `name`
- `description`
- `example`
- `required`
- `readOnly`
- `writeOnly`
- Array の `minItems`
- Array の `maxItems`

`example` は利用コンテキストに応じて Override できる。ただし、Phase 2 Element が定義する Type / Format / Constraint
に適合しなければならない。

### 5.2 Override できない情報

Phase 2 Element が管理する Scalar Constraint は Phase 3 で再定義しない。

例:

- length / minLength / maxLength
- regex / pattern
- minimum / maximum
- precision / scale
- enum
- format

Phase 3 は **値を定義する場所ではなく、値を組み合わせて API
を構成する場所**とする。

------------------------------------------------------------------------

## 6. Constraint

### 6.1 Resource Property

Resource Property は `element` / `resource` / `array` の
**いずれか一つだけ**を持つ。

Resource Property ではコンテキストに応じて以下を指定できる。

| 属性 | 型 | 用途 |
| --- | --- | --- |
| `required` | Boolean | Resource 上で必須か |
| `readOnly` | Boolean | Response 側を基本とする Property か |
| `writeOnly` | Boolean | Request 側を基本とする Property か |
| `minItems` | Integer | Array Property の最小要素数 |
| `maxItems` | Integer | Array Property の最大要素数 |
| `name` | String | コンテキスト上の表示名 |
| `description` | String | コンテキスト上の説明 |
| `example` | Scalar / Object / Array | コンテキスト上の例 |

`minItems` / `maxItems` は `array` を持つ Property にのみ指定できる。 `array` 配下ではなく、`array` と同階層に指定する。

`example` は Resource Property 利用時のコンテキスト情報として指定できる。

- `element` を参照する Property：Scalar
- `resource` を参照する Property：Object
- `array` を参照する Property：Array

Object / Array 内に含まれる Element 由来の値を含め、`example` は参照先 Element の Type / Format / Constraint
を満たす必要がある。

`readOnly: true` と `writeOnly: true` の同時指定は禁止する。

### 6.2 Array

Array は `element` または `resource` のいずれか一つを参照する。

`array` 配下には、配列要素の参照先のみを指定する。

例:

```yaml
details:
  array:
    resource: OrderDetail
  minItems: 1
  required: true
```

`required` / `minItems` / `maxItems` が Resource / Variant 自体の構造上の制約である場合は、
API Usage ではなく Resource / Variant 側に定義する。

例えば `WithDetails` が「1件以上の明細を持つ受注」を意味する場合、
`details.required: true` および `details.minItems: 1` は `WithDetails` 自体に定義する。

`minItems` / `maxItems` は `Array Property` / `Array API Usage` の属性として、 `array` と同階層に指定する。

`required` も同様に `array` の内部には指定しない。

`minItems` / `maxItems` は 0 以上の整数とし、両方指定する場合は `minItems` <= `maxItems` とする。

### 6.3 Variant

Variant は以下を利用できる。

- `include`
- `exclude`
- `add`
- `overrides`

`include` と `exclude` は併用しない。

Variant の `overrides` では、対象 Property のコンテキスト情報を Override できる。

`example` も Override 対象とし、対象 Property の構造に応じて Scalar / Object / Array を指定できる。

例:

```yaml
variants:
  Summary:
    overrides:
      remainingQuantity:
        name: 未出荷・未キャンセル数
        description: |
          受注数のうち、まだ出荷またはキャンセルされていない数量。
        example: 5
```

`add` で同一 Variant 内に追加した Property を、同じ Variant の`overrides` から再度 Override しない。

### 6.4 Parameter Definition

Parameter Definition は Scalar Constraint を持たない。

指定可能な主な属性は以下とする。

- `element`
- `description`
- `example`
- `headerName`

Parameter Definition の `example` は、参照する Element の`example` を Parameter 利用時のコンテキストで Override できる。

`example` は Scalar とし、参照する Element の Type / Format / Constraint を満たす必要がある。

例:

```yaml
orderPic:
  element: userId
  description: 受注担当者。
  example: U1234
```

### 6.5 API Usage Override

API Usage では Resource / Variant の利用コンテキストに対して Override を指定できる。

例:

```yaml
request:
  resource: Order
  variant: WithDetails
  overrides:
    details:
      example:
        - productNo: P123456
          quantity: 10
          sellingPrice: 1000
        - productNo: P654321
          quantity: 5
          sellingPrice: 1200
```

API Usage の `overrides` では、構造に関する情報と `example` を Override できる。

構造に関する Override により元の Resource / Variant と異なる Schema が必要となる場合は、
Resolve 時に API Usage 専用 Schema を生成する。

一方、`example` のみを Override する場合は Schema の構造を変更しない。
この場合、API Usage 専用 Schema は生成せず、元の Resource / Variant Schema をそのまま利用する。

API Usage の `example` は Request / Response の利用コンテキストにおける差分として扱う。
Resolved Model では `exampleOverride` として保持し、完成した Example は OpenAPI Generation 時に生成する。

Array Property では複数要素を含む Array 全体を `example` として指定できる。

Array Property に対する `minItems` / `maxItems` の Override は、Property の同名属性を Overrideする。

Scalar Constraint の Override は禁止する。

------------------------------------------------------------------------

## 7. API 表現ルール

### 7.1 Path Parameter

Path Parameter は Path の `{token}` から導出する。

Source:

``` yaml
path: /orders/{received_order_no}/details/{detail_no}
```

Raw Model:

``` yaml
path: /orders/{received_order_no}/details/{detail_no}
pathParameters:
  - parameterRef: $receivedOrderNo
  - parameterRef: $detailNo
```

`pathParameters` は Operation ではなく **API logical entry** に属する。

Path Parameter は HTTP 上必須であり、参照する Element は API Identifier
として利用可能でなければならない。

Element の `identifier: true` は、その Element が API
の識別子として利用可能であることを示す。単独主キー・単独一意キーであることは意味せず、複合識別子の一要素として利用してよい。

### 7.2 Query Parameter

Query Parameter は Operation の `queryParameters` で指定する。

``` yaml
get:
  queryParameters:
    - orderPic
    - customerId
```

### 7.3 Request Header

Request Header は Service または Operation から Parameter Definition
を参照する。

Service の `requestHeaders` は Service 共通 Header を表す。

Service と Operation の双方から同一 Header を重複指定することは禁止する。

Service の `requestHeaders` は Resolve 時にユーザー定義 Operation へ展開する。

`/health` / `/version` 等の ARIADNE built-in Operation は Service `requestHeaders` の展開対象外とする。

### 7.4 Pagination

Collection GET では `pagination: true` により Pagination を有効化できる。

Raw Model では `pagination: true` を保持し、Resolve 時に以下へ展開する。

- Query Parameter `limit`
- Query Parameter `offset`
- Response Header `Has-More`

`limit` / `offset` は ARIADNE が管理する built-in Parameter とし、Resolved Model では以下のように表現する。

```yaml
parameters:
  limit:
    builtIn: pagination
    usages:
      - in: query
        name: limit
        required: false

  offset:
    builtIn: pagination
    usages:
      - in: query
        name: offset
        required: false
```

Pagination を利用する Operation では、これらを通常の Parameter と同様に parameterRef から参照する。

```yaml
queryParameters:
- parameterRef: $limit
- parameterRef: $offset
```

Pagination Response では、後続データの有無を Has-More Response Header として返す。

```yaml
headers:
  Has-More:
    builtIn: pagination
```

Has-More は後続データの有無のみを表す。総件数は Pagination の built-in 情報には含めない。
総件数が業務 API として必要な場合は、通常の Response Schema として明示的に定義する。

limit / offset は ARIADNE の予約 Parameter 名とし、 Source の Parameter Definition として定義してはならない。

### 7.5 Location

POST Operation が生成したリソースの URI を `Location` Response Header として返却する場合、 Operation 直下に `location` を指定する。

```yaml
post:
  location:
    example: /orders/ORD-20230827-001
```

`location` は POST Operation にのみ指定できる。

`location` を指定する場合、`example` は必須とする。

`location` を省略した場合、`Location` Response Header は生成しない。

`response` は返却データの定義を表し、 `location` は POST Operation によるリソース生成時の振る舞いを表すため、
`response` 配下ではなく Operation 直下に定義する。

### 7.6 HTTP Method / Success Response

Prototype で利用する HTTP Method は以下とする。

- GET
- POST
- PUT
- PATCH
- DELETE

標準 API の基本的な Success Status は以下とする。

| Method | 基本 Status | 備考 |
| --- | ---: | --- |
| GET | 200 | Response Body あり |
| POST | 201 | Response Body は API 定義に応じる |
| PUT | 200 | 全体更新 / upsert を許容 |
| PATCH | 200 | 部分更新 |
| DELETE | 204 | Response Body なし |

標準 API では Success Status を HTTP Method から Resolve 時に補完する。

Custom では Success Response を明示的に定義できる。

```yaml
response:
  success:
    status: 202
    description: エクスポート正常終了
    resource: ExportResult
```

success.status は必須とする。description は任意とする。

Success Response Body は resource または array を指定できる。
いずれも指定しない場合は Response Body なしとする。

### 7.7 PATCH

PATCH Request では、対象 Schema の Property は Request 上すべて optional
として扱う。

PATCH 用 Variant に明示的に含めた `readOnly` Property
は更新対象として利用できる。

### 7.8 Standard Error

全 Operation に Standard Error の `default` Response を自動付与する。

Raw Model では付与せず、Resolve 時に以下の built-in Response として補完する。

```yaml
default:
  builtIn: error
```

Custom で明示的な Error Response を定義した場合も、 Standard Error の `default` Response は併存する。

Resolved Model では Source / Raw Model の `success` / `errors` 表現は解消し、
HTTP Status を Key とする Response Map へ正規化する。

```yaml
responses:
  "202":
    description: エクスポート正常終了
    schemaRef: $ExportResult

  "400":
    description: 入力チェックエラー
    array:
      schemaRef: $ValidationError

  "409":
    description: 処理コンフリクト

  default:
    builtIn: error
```

Standard Error の具体的な OAS Response / Error Schema は OpenAPI Generation の built-in 定義として生成する。

### 7.9 Built-in API

以下を ARIADNE built-in API として Resolve 時に自動生成する。

- `/health`
- `/version`

Source / Raw Model には存在せず、Resolved Model で通常 Operation と同程度まで展開する。

built-in Operation の `tag` は `System` とする。

#### `/health`

```yaml
/health:
  get:
    builtIn: health
    operationId: health-get
    tag: System
    summary: ヘルスチェック

    responses:
      "200":
        description: 正常終了
      default:
        builtIn: error
```

`/health` は Response Body を持たない。 HTTP 200 が返ること自体を正常状態の表現とする。

#### `/version`

```yaml
/version:
  get:
    builtIn: version
    operationId: version-get
    tag: System
    summary: バージョン取得

    responses:
      "200":
        description: 正常終了
        builtIn: version
      default:
        builtIn: error
```

`builtIn: version` の Response は Application Version のみを返す。

```json
{
  "version": "1.0.0"
}
```

Application Version は Git Tag 等から供給することを想定する。

ユーザー定義 API で `/health` / `/version` を定義することは禁止する。

built-in Operation には Service `requestHeaders` を展開しない。

### 7.10 externalDocs

OpenAPI Schema だけでは表現しにくい業務ルールや処理上の補足には `externalDocs` を利用する。

例:

``` yaml
externalDocs:
  url: ./docs/order-create.md
  description: 受注登録の詳細仕様
```

相対パスは Service Root を基準として解決する。

補足文書は ARIADNE YAML の代替正本ではない。Schema / Constraint
として表現可能な情報の正本は ARIADNE YAML 側とする。

------------------------------------------------------------------------

## 8. YAML ファイル仕様

### 8.1 共通 Header

Phase 3 Source YAML は共通 Header を持つ。

``` yaml
formatVersion: "1.0"
updatedAt: "..."
domain: api
kind: ...
```

- `formatVersion` は Phase 3 内で共通とする
- Phase 2 / Phase 4 とは独立して Version を管理できる
- `updatedAt` は Source YAML に保持する
- `domain` は `api`
- `kind` はファイルの役割を示す

Phase 3 の許可 `kind` は以下とする。

- `service`
- `parameters`
- `resource`
- `subresource`
- `action`
- `custom`

### 8.2 Service YAML

Service 自身を定義する。

主な属性:

| 属性 | 必須 | 型 | 意味 |
| --- | --- | --- | --- |
| `service` | 必須 | String | Service ID |
| `name` | 必須 | String | 表示名 |
| `description` | 任意 | String | 説明 |
| `requestHeaders` | 任意 | Array | Service 共通 Request Header |

### 8.3 Parameter YAML

Service 内で利用する Parameter Definition を定義する。

| 属性 | 必須 | 型 | 意味 |
| --- | --- | --- | --- |
| `parameters` | 必須 | Map | Parameter Definition の集合 |

#### Parameter Definition

| 属性 | 必須 | 型 | 意味 |
| --- | --- | --- | --- |
| `element` | 必須 | String | 参照する Phase 2 Element |
| `description` | 任意 | String | Parameter の説明 |
| `example` | 任意 | Scalar | Parameter の例 |
| `headerName` | 任意 | String | Header利用時の物理Header名 |

`example` の Override 規則は「6.4 Parameter Definition」に従う。

例:

``` yaml
parameters:
  receivedOrderNo:
    element: receivedOrderNo
    description: 受注番号。

  orderPic:
    element: userId
    description: 受注担当者。
    example: U1234

  traceId:
    element: traceId
    headerName: X-Trace-ID
    description: トレースID。
```

未使用 Parameter Definition が存在しても正常とする。Raw Model では全定義を保持し、必要な定義の抽出は Resolve で行う。

### 8.4 Main Resource YAML

Main Resource を定義する。

主に以下を管理する。

| 属性 | 必須 | 型 | 意味 |
| --- | --- | --- | --- |
| `resource` | 必須 | String | Resource 名 |
| `name` | 必須 | String | 表示名 |
| `description` | 任意 | String | 説明 |
| `properties` | 必須 | Map | Property 定義 |
| `variants` | 任意 | Map | Variant 定義 |
| `api` | 任意 | Map | API 定義 |

Main Resource は `parent` を持たない。

### 8.5 SubResource YAML

SubResource を定義する。

Main Resource と同様に以下を管理できる。

| 属性 | 必須 | 型 | 意味 |
| --- | --- | --- | --- |
| `resource` | 必須 | String | Resource 名 |
| `parent` | 必須 | String | 親 Main Resource |
| `name` | 必須 | String | 表示名 |
| `description` | 任意 | String | 説明 |
| `properties` | 必須 | Map | Property 定義 |
| `variants` | 任意 | Map | Variant 定義 |
| `parentVariants` | 任意 | Map | 親 Resource に追加する Variant 定義 |
| `api` | 任意 | Map | API 定義 |

### 8.6 Action YAML

Action を定義する。

| 属性 | 必須 | 型 | 意味 |
| --- | --- | --- | --- |
| `action` | 必須 | String | Action 名 |
| `resources` | 必須 | Map | Action-local Resource 定義 |
| `api` | 必須 | Map | API 定義 |

Action 自身は表示名・説明を持たない。API の表示名・説明は Operation の `summary` / `description` に定義する。

Local Resource は Source / Raw Model 上では Action-local 定義として参照する。
ただし Local Resource 名は Service 内で一意とする。

### 8.7 Custom YAML

Custom API を定義する。

| 属性 | 必須 | 型 | 意味 |
| --- | --- | --- | --- |
| `custom` | 必須 | String | Custom 名 |
| `resources` | 任意 | Map | Custom-local Resource 定義 |
| `api` | 必須 | Map | API 定義 |

Custom 自身は表示名・説明を持たない。API の表示名・説明は Operation の `summary` / `description` に定義する。

Custom は特定 Main Resource に従属しないため、各 Operation に `tag` を必須指定する。

Custom の Response は、`success` と `errors` により定義する。

- `success` は1件定義する
- `errors` は任意で、0件以上定義できる
- 各 Response の `status` は必須
- `description` は任意
- Response Body は `resource` または `array` を指定できる
- `resource` / `array` のいずれも指定しない場合は Body なしとする
- `resource` と `array` は同時指定できない

| 属性 | 必須 | 型 | 意味 |
| --- | --- | --- | --- |
| `response.success` | 必須 | Map | Success Response |
| `response.success.status` | 必須 | Integer | HTTP Status |
| `response.success.description` | 任意 | String | Response説明 |
| `response.success.resource` | 任意 | Resource Reference | 単一Resource Body |
| `response.success.array` | 任意 | Array | Resource Array Body |
| `response.errors` | 任意 | Array | 明示Error Response群 |
| `response.errors[].status` | 必須 | Integer | HTTP Status |
| `response.errors[].description` | 任意 | String | Response説明 |
| `response.errors[].resource` | 任意 | Resource Reference | 単一Resource Body |
| `response.errors[].array` | 任意 | Array | Resource Array Body |

例:

```yaml
response:
  success:
    status: 202
    description: エクスポート正常終了
    resource: ExportResult

  errors:
    - status: 400
      description: 入力チェックエラー
      array:
        resource: ValidationError

    - status: 409
      description: 処理コンフリクト
```

Custom-local Resource も Source / Raw Model 上では Custom-local 定義として管理する。
ただし Local Resource 名は Service 内で一意とする。

------------------------------------------------------------------------

## 9. Raw Model

### 9.1 目的

Raw Model は、Phase 3 ARIADNE Source を **Service 単位の意味モデル**へ正規化した中間モデルである。

Source of Truth ではなく、ARIADNE Source から再生成可能な成果物とする。

Raw Model は YAML として出力し、デバッグや Golden Test に利用できるようにする。

生成先の現時点案:

``` text
dist/api/model/order-management.raw.yaml
```

### 9.2 基本構造

``` yaml
formatVersion: "1.0"

service:
  id: order-management
  name: 受注サービス
  description: |
    受注および顧客情報を管理するサービス。
  requestHeaders:
    - parameterRef: $traceId

parameters:
  ...

resources:
  ...

actions:
  ...

customs:
  ...
```

Raw Model には `updatedAt` を持たせない。生成物の Timestamp による不要な差分を避けるためである。

### 9.3 Source kind の正規化

| Source | Raw Model |
| --- | --- |
| `kind: service` | `service` |
| `kind: parameters` | `parameters` |
| `kind: resource` | `resources` |
| `kind: subresource` | `resources` + `parentRef` |
| `kind: action` | `actions` |
| `kind: custom` | `customs` |

Source の `kind` および `resource` / `action` / `custom` といった識別フィールドは、Raw Model では構造そのものへ吸収する。

### 9.4 Main / SubResource

Raw Model に SubResource 専用型は設けない。

``` yaml
resources:
  Order:
    ...

  OrderDetail:
    parentRef: $Order
    ...
```

Main / Sub の違いは `parentRef` の有無で表現する。

### 9.5 Raw Model で保持するもの

以下は Raw Model では未解決のまま保持する。

- Variant の `include` / `exclude` / `add` / `overrides`
- `parentVariants`
- Resource / Variant 参照
- API Usage Override
- `pagination: true`
- Service / Operation の Parameter Definition 参照
- Action / Custom の Local Resource Scope
- Custom の Response 定義
  - `success`
  - `errors`
  - `status`
  - `description`
  - Response Body
- `externalDocs`

### 9.6 参照表現

Raw Model では、ARIADNE が管理する別定義への参照を `xxxRef` 属性と `$` Prefix によって明示する。

| 参照対象 | Raw Model 表記 | 例 |
| --- | --- | --- |
| Phase 2 Element | `elementRef` | `elementRef: $receivedOrderNo` |
| Resource / Local Resource | `resourceRef` | `resourceRef: $Order` |
| Variant | `variantRef` | `variantRef: $Summary` |
| Parent Resource | `parentRef` | `parentRef: $Order` |
| Parameter Definition | `parameterRef` | `parameterRef: $receivedOrderNo` |

`$` は Raw / Resolved Model において、ARIADNE が管理する別定義への参照値であることを示す。

`$` は Source YAML では使用せず、中間モデル上の参照表現としてのみ使用する。

Property 名など、同一定義内部の構成要素を選択するための名称には `$` を付与しない。
そのため、Variant の `include` / `exclude` / `overrides` の対象 Property 名は参照表現とはしない。

### 9.7 Raw Model 生成時に補完するもの

File Validation の範囲で一意に導出できる情報のみ補完する。

現時点では Path Token から `pathParameters` を生成する。

``` yaml
member:
  path: /orders/{received_order_no}
  pathParameters:
    - parameterRef: $receivedOrderNo

  get:
    ...
  put:
    ...
```

Path Token は semantic key に正規化したうえで、対応する Parameter Definition への `parameterRef` として表現する。

### 9.8 Raw Model で行わないこと

Raw Model 生成時には以下を行わない。

- Phase 3 ファイル間参照の解決
- Phase 2 Element Definition の取り込み
- Element 参照の存在確認
- Variant の展開
- Parent Variant の親 Resource への統合
- API Usage Override の展開
- Parameter Definition の参照解決
- 未使用 Parameter Definition の削除
- Service requestHeaders の Operation への展開
- Pagination の built-in Parameter への展開
- HTTP Success Status の自動決定
- Standard Error Response の付与
- Location Header の付与
- `/health` / `/version` の生成
- OperationId の生成
- OpenAPI Schema への変換

Raw Model の境界は、

> **File Validation の範囲で一意に構造化できるところまで**

とする。

------------------------------------------------------------------------

## 10. Validation Architecture

Phase 3 以降の Validation は以下の4段階に分離する。

``` text
Phase 2 / Phase 3 YAML
        ↓
① File Validation
        ↓
Raw Model
        ↓
② API Validation
        ↓
Resolved Model
        ↓
③ Service Validation
        ↓
OpenAPI 3.1
        ↓
④ OpenAPI Validation
```

### 10.1 File Validation

単一 Phase 3 YAML の構造を検証する。

外部ファイルを参照せず、

> このファイルを Raw Model の一部として一意に解釈できるか

を保証する。

### 10.2 API Validation

Service 内の Phase 3 ファイルを横断して検証する。

主な対象:

- Resource / Variant / Parameter Definition の参照
- SubResource の parent
- Parent Variant
- Resource Reference Cycle
- `(path, method)` の一意性
- Local Resource Scope
- Service 内 Schema 名の衝突
- Phase 3 formatVersion の一致

Phase 2 Element は参照しない。

### 10.3 Service Validation

Phase 2 と Phase 3 を接続して検証する。

主な対象:

- Element 参照の存在
- Path Parameter の Element が `identifier: true` か
- example と Element Type / Constraint の整合性

### 10.4 OpenAPI Validation

生成した OpenAPI 3.1 を標準 OpenAPI Validator で検証する。

------------------------------------------------------------------------

## 11. Validation Rules v0.1

Phase 2 の Validation ID `V-01` ～ `V-20` に続き、Phase 3 は `V-21`
以降を使用する。

**注:** Phase 3 は設計途中のため、以下は現時点で確定している Rule
群の整理である。最終 ID 割当および Rule の追加・統合は Phase 3
完了時に見直す。

### 11.1 File Validation

| ID | Validation Rule | 判定 |
| --- | --- | --- |
| V-21 | 共通 Header の必須項目が存在する | Error |
| V-22 | `formatVersion` が ARIADNE の対応する Phase 3 Format である | Error |
| V-23 | `updatedAt` が ISO 8601 として妥当である | Error |
| V-24 | `domain` が `api` である | Error |
| V-25 | `kind` が Phase 3 で許可された値である | Error |
| V-26 | 各名称が対象ごとの命名規約を満たす | Error |
| V-27 | `kind` ごとの必須 / 許可属性を満たす | Error |
| V-28 | 未定義属性を持たない | Error |
| V-29 | YAML Map に重複 Key が存在しない | Error |
| V-30 | Resource Property は `element` / `resource` / `array` のいずれか一つだけを持つ | Error |
| V-31 | `readOnly: true` と `writeOnly: true` を同時指定しない | Error |
| V-32 | Phase 3 で Scalar Constraint を再定義しない | Error |
| V-33 | Array は `element` / `resource` のいずれか一つだけを持つ | Error |
| V-34 | `minItems` / `maxItems` は 0 以上の整数であり、両方指定時は `minItems <= maxItems` | Error |
| V-35 | 空 Map / 空 Array を明示的に記述しない | Error |
| V-36 | 必須文字列属性に空文字 / 空白のみを指定しない | Error |
| V-37 | `/health` / `/version` をユーザー API として定義しない | Error |
| V-38 | 同一 Path 内で同じ Path Token を複数回使用しない | Error |
| V-39 | Path Token が snake_case の命名規約を満たす | Error |
| V-40 | HTTP Method が `get/post/put/patch/delete` のいずれかである | Error |
| V-41 | GET に Request Body を定義しない | Error |
| V-42 | DELETE に Request / Response Body を定義しない | Error |
| V-43 | `pagination` は GET にのみ指定する | Error |
| V-44 | `externalDocs` を指定する場合 `url` が存在する | Error |
| V-45 | Operation の `tag` は Custom では必須、Resource / SubResource / Action では指定しない | Error |
| V-46 | `limit` / `offset` を Source の Parameter Definition 名として定義しない | Error |
| V-47 | `minItems` / `maxItems` は Array Property、または Array Property に対する Variant / API Usage Override にのみ指定する | Error |
| V-48 | `location` は POST にのみ指定する | Error |
| V-49 | `location` を指定する場合 `example` が存在する | Error |

### 11.2 API Validation

以下は Rule 内容を確定済みとし、最終 ID は Phase 3
完了時に連番整理する。

- Resource 参照先が存在する
- Variant 参照先が存在する
- Parameter Definition 参照先が存在する
- Variant の `include` / `exclude` / `overrides` の対象 Property が妥当である
- Variant の `add` が既存 Property と衝突しない
- 同一 Variant の `add` Property を同じ Variant の `overrides` で指定しない
- SubResource の `parent` が存在し、Main Resource である
- SubResource を親とする多段 SubResource を禁止する
- `parentVariants` の対象と統合内容が妥当である
- Resource Reference Cycle が存在しない
- Local Resource の参照 Scope が妥当である
- Service 内で Resource / Local Resource の名前が衝突しない
- `(path, method)` が Service 内で一意である
- Phase 3 の全 Source が同一 `formatVersion` を使用する
- Service と Operation で同一 Request Header を重複指定しない
- Service YAML が Service 内に1つだけ存在する
- parameters.yaml は Service 内に0または1つとする
- ユーザー定義 Operation が1件以上存在する
- Resolve 後に生成される Schema 名が衝突しない
- Custom の同一 Operation 内で `response.errors[].status` が重複しない
- Variant / API Usage の `minItems` / `maxItems` の対象が Array Property である
- 生成される `operationId` が Service 内で一意である

### 11.3 Service Validation

以下は Phase 2 / Phase 3 の接続ルールとして扱う。

- すべての `element` 参照先が Phase 2 に存在する
- Path Token から正規化した semantic key が Parameter Definition に存在する
- Path Parameter が参照する Element の effective `identifier` が `true` である
- `example` に含まれる Element 由来の値が、対応する Element の Type / Format / Constraint に適合する

------------------------------------------------------------------------

## 12. Resolved Model

Resolved Model は、API Validation を通過した Raw Model から生成する、
API として意味解決済みの Service Model である。

Raw Model が ARIADNE Source の Authoring 構造や差分表現を保持するのに対し、
Resolved Model では Variant、Parent Variant、API Usage Override、
Parameter Usage、built-in 等を解決し、OpenAPI 3.1 へ変換可能な完成した API 意味モデルを表現する。

Resolved Model は Source of Truth ではなく、Raw Model から再生成可能な中間成果物とする。

生成先の現時点案:

```text
dist/api/model/order-management.resolved.yaml
```

Raw Model と同様、Resolved Model にも `updatedAt` は持たせない。

### 12.1 基本構造

Resolved Model は Service 単位で生成する。

基本構造を以下とする。

```yaml
formatVersion: "1.0"

service:
  id: order-management
  name: 受注サービス
  description: |
    受注および顧客情報を管理するサービス。

parameters:
  ...

schemas:
  ...

apis:
  ...

elementRefs:
  ...
```

Resolved Model では、Source / Raw Model の以下の Authoring 上の分類は最終的な API 意味モデルへ統合する。

- Main Resource / SubResource
- Variant
- Parent Variant
- Action-local Resource
- Custom-local Resource
- API logical entry

これらの由来は必要に応じて `origin` に保持する。

### 12.2 Service

Resolved Model は Service の Identity / Metadata を保持する。

```yaml
service:
  id: order-management
  name: 受注サービス
  description: |
    受注および顧客情報を管理するサービス。
```

`id` / `name` / `description` は Raw Model から引き継ぐ。

Raw Model の `service.requestHeaders` は Resolve 時に
ユーザー定義 Operation へ展開するため、
Resolved Model の `service` には保持しない。

Resolved Model の `service` は、
OpenAPI Generation に必要な Service 自身の意味情報を保持する。

### 12.3 Reference

Resolved Model では、参照対象に応じて以下の Reference を使用する。

| 参照対象 | Resolved Model 表記 | 例 |
| --- | --- | --- |
| Phase 2 Element | `elementRef` | `elementRef: $receivedOrderNo` |
| 解決済み Schema | `schemaRef` | `schemaRef: $Order` |
| Parameter | `parameterRef` | `parameterRef: $receivedOrderNo` |

Raw Model の `resourceRef` / `variantRef` は Resolve 時に解決し、 Resolved Model では `schemaRef` に統合する。

例:

```yaml
resourceRef: $Order
variantRef: $Summary
```

は、Resolved Model では以下となる。

```yaml
schemaRef: $Order.Summary
```

Variant を使用しない Resource Reference は以下となる。

```yaml
schemaRef: $Order
```

Phase 2 Element は Resolved Model に定義本体を取り込まず、 `elementRef` のまま保持する。

### 12.4 Schema

Resolved Model の `schemas` には、 Resolve 後の完成した Phase 3 Schema を格納する。

Schema 名は以下を基本とする。

| Source | Resolved Schema 名 |
| --- | --- |
| Main Resource | `Resource` |
| SubResource | `Resource` |
| Variant | `Resource.Variant` |
| Parent Variant | `Resource.Variant` |
| Action-local Resource | Local Resource 名 |
| Custom-local Resource | Local Resource 名 |
| API Usage により生成された Schema | 生成された一意な Schema 名 |

例:

```yaml
schemas:
  Order:
    ...

  Order.Summary:
    ...

  Order.WithDetails:
    ...

  OrderDetail:
    ...

  OrderDetail.Summary:
    ...

  OrderShipment:
    ...

  ExportCondition:
    ...

  PostOrdersRequest:
    ...
```

Main Resource / SubResource の名前は Service 内で一意とする。

Variant 名は Resource 内で一意とし、 Resolved Schema 名では `Resource.Variant` として表現する。

Action / Custom の Local Resource は、 Resolved Model では Action / Custom 名を Namespace とせず、
Service 内の通常 Schema として扱う。そのため Local Resource 名も Service 内で一意でなければならない。

### 12.5 Variant / Parent Variant の Resolve

Variant の以下の差分表現は Resolve 時に適用する。

- `include`
- `exclude`
- `add`
- `overrides`

Resolved Model では Variant の差分を保持せず、最終的な Property を持つ完成 Schema とする。

Parent Variant も同様に、
SubResource の `parentVariants` を親 Main Resource へ統合し、完成した Variant Schema とする。

例として、`OrderDetail` が親 `Order` に定義した `WithDetails` は、 Resolved Model では以下の Schema となる。

```yaml
schemas:
  Order.WithDetails:
    ...
    properties:
      ...
      details:
        array:
          schemaRef: $OrderDetail
```

Resolved Model には `parentVariants` 自体は残さない。

### 12.6 API Usage Override

API Usage の `overrides` は、構造に関する Override と `example` の Override を分けて Resolve する。

#### 構造に関する Override

API Usage によって元の Resource / Variant と異なる Schema 構造が必要となる場合、
Override を適用した API Usage 専用 Schema を生成する。

例:

```yaml
PostOrdersRequest:
  origin:
    kind: apiUsage
    operationId: orders-post
    usage: request
    baseSchemaRef: $Order.WithDetails
  ...
```

API Usage Schema は、API 上で実際に使用される完成した Property 構造を持つ。

構造に関する Override が存在しない場合は API Usage 専用 Schema を生成せず、元の Resource / Variant Schema を直接参照する。

#### Example Override

API Usage で指定された `example` は Schema 構造へ適用しない。

Resolved Model では API Usage の差分 Example として `exampleOverride` に保持する。

例:

```yaml
request:
  schemaRef: $Order.WithDetails
  exampleOverride:
    details:
      - productNo: P123456
        quantity: 10
        sellingPrice: 1000
      - productNo: P654321
        quantity: 5
        sellingPrice: 1200
```

`exampleOverride` は完成した Example ではなく、API Usage で指定された差分を表す。

Request / Response の完成した Example は OpenAPI Generation 時に生成する。

### 12.7 Schema origin

Resolved Schema には、生成元を追跡するため `origin` を保持できる。

`origin` は API Schema の意味そのものではなく、ARIADNE Source からどのように生成されたかを示す provenance 情報とする。

Main Resource:

```yaml
origin:
  kind: resource
  resourceKind: mainResource
  resource: Order
```

Variant:

```yaml
origin:
  kind: variant
  resourceKind: mainResource
  resource: Order
  variant: Summary
```

SubResource:

```yaml
origin:
  kind: resource
  resourceKind: subResource
  resource: OrderDetail
  parent: Order
```

Parent Variant:

```yaml
origin:
  kind: parentVariant
  resourceKind: mainResource
  resource: Order
  variant: WithDetails
  definedBy: OrderDetail
```

Action-local Resource:

```yaml
origin:
  kind: action
  action: OrderShipment
  resource: OrderShipment
  parent: Order
```

Custom-local Resource:

```yaml
origin:
  kind: custom
  custom: ExportOrders
  resource: ExportCondition
```

API Usage Schema:

```yaml
origin:
  kind: apiUsage
  operationId: orders-post
  usage: request
  baseSchemaRef: $Order.WithDetails
```

`resourceKind` は Main Resource / SubResource のトポロジ上の役割を表すために使用し、
Action / Custom / API Usage には使用しない。

### 12.8 Parameter

Raw Model では Service 内の Parameter Definition をすべて保持するが、
Resolved Model では実際に使用される Parameter のみを保持する。

Parameter Definition 自体を各 Operation へインライン展開せず、
Resolved Model の `parameters` に保持し、API 側から `parameterRef` で参照する。

Resolved Parameter は、その Parameter が HTTP 上でどのように利用されるかを
`usages` として保持する。

例:

```yaml
parameters:
  customerId:
    elementRef: $customerId
    description: 顧客ID。
    usages:
      - in: query
        name: customer_id
        required: false
```

同じ Parameter Definition が異なる HTTP Usage を持つ場合は、複数の `usages` を持つことができる。

`usages` は Operation ごとの利用履歴ではなく、OpenAPI Parameter として必要となる物理的な利用形態の集合を表す。

Path Parameter は `required: true` とする。

Query Parameter は明示的な Override がない場合 `required: false` とする。

Service 共通 Request Header は `required: true` とする。

Source の `headerName` は Resolve 時に `usages.name` へ反映する。

### 12.9 Pagination Parameter

`pagination: true` は Resolve 時に ARIADNE built-in Parameter へ展開する。

```yaml
parameters:
  limit:
    builtIn: pagination
    usages:
      - in: query
        name: limit
        required: false

  offset:
    builtIn: pagination
    usages:
      - in: query
        name: offset
        required: false
```

Pagination を利用する Operation は以下を参照する。

```yaml
queryParameters:
  - parameterRef: $limit
  - parameterRef: $offset
```

`limit` / `offset` は Phase 2 Element を参照しないため、 `elementRefs` には含めない。

### 12.10 API / Path / Operation

Resolved Model の `apis` は HTTP Path を直接 Key とする。

Source / Raw Model の API logical entry は Resolved Model では保持しない。

```yaml
apis:
  /orders:
    get:
      ...

    post:
      ...

  /orders/{received_order_no}:
    pathParameters:
      - parameterRef: $receivedOrderNo

    get:
      ...

    put:
      ...
```

Path Parameter は Path 共通情報として Path 階層に保持する。

HTTP Method は Path 配下に配置する。

Resource / Action / Custom といった Source 上の API grouping は Resolved Model の `apis` 構造には残さない。

### 12.11 OperationId

`operationId` は Resolve 時に HTTP Path と Method から生成する。

基本形式は以下とする。

```text
normalized-path-method
```

Path を先、Method を末尾とする。

例:

| Path / Method | operationId |
| --- | --- |
| `GET /orders` | `orders-get` |
| `POST /orders` | `orders-post` |
| `GET /orders/{received_order_no}` | `orders-received-order-no-get` |
| `PATCH /orders/{received_order_no}/date` | `orders-received-order-no-date-patch` |
| `DELETE /orders/{received_order_no}/details/{detail_no}` | `orders-received-order-no-details-detail-no-delete` |
| `GET /health` | `health-get` |
| `GET /version` | `version-get` |

Path Token の `{}` は除去し、snake_case は kebab-case へ正規化する。

生成された `operationId` は Service 内で一意でなければならない。

### 12.12 Operation tag

Resolved Operation は `tag` を1件持つ。

Source の定義元に応じて以下のように決定する。

- Main Resource：Main Resource 名
- SubResource：親 Main Resource 名
- Action：所属する Main Resource 名
- Custom：Source Operation で明示した `tag`
- built-in Operation：`System`

Resolved Model では単数の `tag` を保持し、OpenAPI Generation 時に OAS の `tags` 配列へ変換する。

### 12.13 API origin

Resolved Operation には、API の Source 上の定義元を示す `origin` を保持する。

API `origin` は「その Operation がどこに定義されていたか」を表し、
Request / Response がどの Schema を利用するかとは独立する。

Main Resource:

```yaml
origin:
  kind: resource
  resourceKind: mainResource
  resource: Order
```

SubResource:

```yaml
origin:
  kind: resource
  resourceKind: subResource
  resource: OrderDetail
  parent: Order
```

Action:

```yaml
origin:
  kind: action
  action: OrderShipment
  parent: Order
```

Custom:

```yaml
origin:
  kind: custom
  custom: ExportOrders
```

Variant / Parent Variant は Schema の生成元であり、API の定義元ではないため API `origin.kind` には使用しない。

API の由来、利用 Schema、Schema の由来はそれぞれ独立して扱う。

```text
API はどこに定義されていたか
    → API origin

API がどの Schema を利用するか
    → schemaRef

その Schema がどの Source から生成されたか
    → Schema origin
```

### 12.14 externalDocs

Source / Raw Model の `externalDocs` は、Resolve 後も対象 Operation に保持する。

`externalDocs` は API の補足文書への参照であり、
Variant や Parameter のような意味解決対象ではないため、
Resolve では内容を展開しない。

相対 `url` の基準は Source と同様に Service Root とする。

Resolved Model では Source の Markdown URL をそのまま保持する。
Markdown から HTML への変換および OAS 上の URL 変換は OpenAPI Generation / API Document Generation の責務とする。

### 12.15 Request / Response Schema

Resolved Operation の Request / Response Body は、解決済み Schema を `schemaRef` で参照する。

単一 Resource:

```yaml
request:
  schemaRef: $Order
```

Array:

```yaml
responses:
  "200":
    description: 正常終了
    array:
      schemaRef: $Order.Summary
```

Request Body の必須性など、ARIADNE の API 意味モデルとして確定できる情報は Resolve 時に補完する。

### 12.16 Responses

Resolved Model の `responses` は HTTP Status を Key とする Map とする。

標準 API では HTTP Method から Success Status を Resolve 時に決定する。

例:

```yaml
responses:
  "200":
    description: 正常終了
    schemaRef: $Order

  default:
    builtIn: error
```

Custom の Source / Raw Model で使用する `success` / `errors` は Resolve 時に同じ Status Map へ変換する。

```yaml
responses:
  "202":
    description: エクスポート正常終了
    schemaRef: $ExportResult

  "400":
    description: 入力チェックエラー
    array:
      schemaRef: $ValidationError

  "409":
    description: 処理コンフリクト

  default:
    builtIn: error
```

Response Header は各 Status Response の `headers` に保持する。

Header の物理 HTTP 名を Map Key とする。

Pagination:

```yaml
headers:
  Has-More:
    builtIn: pagination
```

POST Operation に `location` が指定されている場合:

```yaml
headers:
  Location:
    builtIn: location
    example: /orders/ORD-20230827-001
```

Source / Raw Model の `location` は Resolve 時に、成功 Response の `Location` Header へ展開する。

`location` が指定されていない POST Operation には、 `Location` Header を生成しない。

### 12.17 Standard Error

全 Operation に Standard Error の `default` Response を Resolve 時に付与する。

```yaml
default:
  builtIn: error
```

`builtIn: error` は、
ARIADNE が管理する Standard Error Response であることを示す。

Resolved Model では Standard Error の具体的な Schema を展開しない。
具体的な OAS Response / Error Schema は OpenAPI Generation で生成する。

Custom が明示的な Error Status を持つ場合も、 `default` Standard Error は併存する。

### 12.18 builtIn

`builtIn` は、Source で業務 API として明示的に定義されたものではなく、
ARIADNE が規約に基づいて Resolve 時に補完した定義であることを示す。

現時点で以下を使用する。

| builtIn | 対象 | 意味 |
| --- | --- | --- |
| `pagination` | Parameter / Response Header | Pagination により生成 |
| `location` | Response Header | Source / Raw Model の `location` 指定により生成 |
| `error` | default Response | Standard Error |
| `health` | Operation | `/health` built-in API |
| `version` | Operation / Response | `/version` built-in API |

### 12.19 Built-in API

`/health` / `/version` は Source / Raw Model には存在せず、Resolve 時に built-in Operation として追加する。

```yaml
/health:
  get:
    builtIn: health
    operationId: health-get
    tag: System
    summary: ヘルスチェック
    responses:
      "200":
        description: 正常終了
      default:
        builtIn: error

/version:
  get:
    builtIn: version
    operationId: version-get
    tag: System
    summary: バージョン取得
    responses:
      "200":
        description: 正常終了
        builtIn: version
      default:
        builtIn: error
```

`/health` は Response Body を持たない。

`/version` の `builtIn: version` Response は Application Version のみを返す。

built-in Operation には Service `requestHeaders` を展開しない。

### 12.20 elementRefs

Resolved Model は、その Service が利用する Phase 2 Element の依存集合を `elementRefs` として保持する。

```yaml
elementRefs:
  - $receivedOrderNo
  - $orderDate
  - $userId
  - $customerId
  - $orderStatus
  ...
```

`elementRefs` は Phase 2 Element Definition 自体を複製するものではない。

主目的は OpenAPI Generation 時に、必要な Element Schema を一度だけ生成できるようにすることである。

ARIADNE built-in Parameter 等、Phase 2 Element に由来しない定義は `elementRefs` に含めない。

### 12.21 Resolve の境界

Resolve では、ARIADNE の API 意味モデルとして一意に決定できる情報を完成させる。

主な Resolve 対象は以下とする。

- Variant の完全 Schema 化
- Parent Variant の親 Resource への統合
- API Usage Override の適用
- Local Resource の Schema 化
- `resourceRef` / `variantRef` から `schemaRef` への解決
- 使用 Parameter Definition の抽出
- Parameter Usage の解決
- Service `requestHeaders` のユーザー定義 Operation への展開
- Pagination の built-in Parameter / Response Header への展開
- HTTP Success Status の決定
- Standard Error Response の付与
- `location` 指定による Location Header の付与
- OperationId の生成
- Operation tag の確定
- API / Schema origin の付与
- `/health` / `/version` built-in Operation の生成
- `elementRefs` の抽出

一方、以下は Resolve では行わない。

- Phase 2 Element Definition の展開
- OpenAPI `$ref` への変換
- OpenAPI `components` 構造への変換
- OpenAPI 固有 Object への変換
- Standard Error の具体的な OAS Schema / Response 生成
- `/version` の具体的な OAS Response Schema 生成
- OpenAPI `servers` の生成

これらは OpenAPI 3.1 Generation の責務とする。

Resolved Model は、

> **ARIADNE が補完・解決すべき意味をすべて確定し、OpenAPI 固有表現への変換だけを残した Service Model**

と位置付ける。

------------------------------------------------------------------------

## 13. OpenAPI 3.1 Generation

Prototype Phase 3 の生成対象は **OpenAPI 3.1** とする。
OpenAPI 3.0 は対象外とする。

OpenAPI Generation は、Service Validation を通過した Resolved Model を入力とし、
ARIADNE の意味モデルを OpenAPI 3.1 の構造へ変換する。

Resolved Model の時点で API としての意味は確定済みとし、
OpenAPI Generation では新たな業務的意味の解決を行わない。

Prototype では `order-management` Service の OAS 3.1 サンプルを
`dist/api/oas/order-management/openapi.yaml` に手作業で作成し、
Redocly CLI による Validation と ReDoc 表示確認を通じて本仕様の妥当性を検証した。

### 13.1 Generation の責務

OpenAPI Generation の主な責務は以下とする。

- Resolved `service` を OpenAPI `info` へ変換する
- Resolved `schemas` を `components/schemas` へ変換する
- `elementRefs` が参照する Phase 2 Element を `components/schemas` へ変換する
- Resolved `parameters` / Parameter Usage を OpenAPI Parameter Object へ変換する
- Resolved `apis` を OpenAPI `paths` / Operation Object へ変換する
- `schemaRef` / `elementRef` / `parameterRef` を OpenAPI `$ref` へ変換する
- Resolved `responses` を OpenAPI Response Object へ変換する
- ARIADNE built-in を OpenAPI の具体的な定義へ変換する
- `externalDocs` を OpenAPI External Documentation Object へ変換する
- 実行環境未設定を示す固定の `servers` を生成する

Generation は Variant の展開、Parameter Usage の決定、
Success Status の決定、OperationId の生成等を行わない。
これらは Resolve 時点で確定済みとする。

### 13.2 OpenAPI 基本構造

生成する OpenAPI の基本構造は以下とする。

```yaml
openapi: 3.1.0

info:
  title: 受注サービス
  description: |
    受注および顧客情報を管理するサービス。
  version: ...

servers:
  ...

paths:
  ...

components:
  schemas:
    ...
  parameters:
    ...
  responses:
    ...
```

`info.title` / `info.description` は Resolved `service` から生成する。

`info.version` は ARIADNE Source / Resolved Model では管理せず、 Application Version を使用する。

Application Version は、Semantic Version として解釈可能な Git Tag から Build 時に取得する。
Semantic Version として解釈できない開発用 Tag 等は Application Version として使用しない。

`servers` は ARIADNE Source / Resolved Model では管理しない。

OpenAPI Generation 時に、実行環境が未設定であることを示す固定の Server を生成する。

```yaml
servers:
  - url: https://not-configured.invalid
    description: 実行環境で設定
```

この OAS 自体は実行環境の接続先を定義しない。実際の接続先は実行環境側で設定し、OpenAPI Generation の責務とはしない。

### 13.3 Schema Generation

Resolved Model の `schemas` は OpenAPI `components/schemas` へ変換する。

Resolved Schema はすでに Variant / Parent Variant、および API Usage の構造に関する Override が
適用された完成 Schema であるため、OpenAPI Generation では Schema の再構成を行わない。
API Usage の `exampleOverride` は Schema 構造には適用せず、
Request / Response Generation 時に完成 Example の生成へ使用する。

例えば、

```yaml
schemaRef: $Order.Summary
```

は、対応する OpenAPI Schema Component への `$ref` に変換する。

Phase 2 Element を参照する

```yaml
elementRef: $customerId
```

については、Phase 2 Element Definition から対応する OpenAPI Schema を生成し、
`components/schemas` に配置したうえで `$ref` へ変換する。

同一 Element は Service 内で一度だけ Schema Component を生成する。

生成対象となる Element は Resolved Model の `elementRefs` から決定する。

### 13.4 Parameter Generation

Resolved Model の `parameters` は、 `elementRef` と `usages` を組み合わせて OpenAPI Parameter Object へ変換する。

ARIADNE Parameter Definition は semantic な Parameter を表すが、
OpenAPI Parameter Object は `in` / `name` / `required` 等を含む物理的な HTTP Parameter を表す。

そのため、1つの Resolved Parameter が複数の `usages` を持つ場合、
Usage ごとに OpenAPI Parameter Object を生成する。

Operation 側の

```yaml
parameterRef: $customerId
```

は、その Operation に対応する Usage の OpenAPI Parameter Component への `$ref` に変換する。

Pagination により生成された `limit` / `offset` も、通常の OpenAPI Query Parameter として生成する。

### 13.5 Path / Operation Generation

Resolved Model の `apis` は OpenAPI `paths` へ直接変換する。

Resolved:

```yaml
apis:
  /orders/{received_order_no}:
    pathParameters:
      - parameterRef: $receivedOrderNo

    get:
      operationId: orders-received-order-no-get
      tag: Order
      summary: 受注取得
      ...
```

OpenAPI Generation では、

- Path Key を OpenAPI Path Item へ変換する
- `pathParameters` を Path Item の `parameters` へ変換する
- HTTP Method を Operation Object へ変換する
- `operationId` をそのまま使用する
- ARIADNE `tag` を OpenAPI `tags` 配列へ変換する
- `summary` / `description` を対応する Operation 属性へ変換する

`operationId` や `tag` の決定は Generation では行わない。

### 13.6 Request / Response Generation

Resolved Model の `request` / `responses` は、 OpenAPI Request Body / Response Object へ変換する。

`schemaRef` は OpenAPI Schema `$ref` へ変換する。

Array の場合は OpenAPI の

```yaml
type: array
items:
  $ref: ...
```

へ変換し、`minItems` / `maxItems` が存在する場合は Array Schema の同名属性へ変換する。

Resolved `responses` の HTTP Status Key はそのまま OpenAPI `responses` の Status Key として使用する。

Response Header が存在する場合は、OpenAPI Response Object の `headers` へ変換する。

#### Example Generation

Resolved Request / Response に `exampleOverride` が存在する場合、
OpenAPI Generation で Request / Response の利用コンテキストに応じた完成 Example を生成する。

Example の値は、以下の優先順位で解決する。

```text
API Usage exampleOverride
    ↓
Resource / Variant Property example
    ↓
Phase 2 Element example
```

上位で指定された Example は、対応する下位の Example を Override する。

Request Example の生成では `readOnly: true` の Property を送信対象とせず、
Response Example の生成では `writeOnly: true` の Property を応答対象としない。

API Usage の Example は OpenAPI Schema Component の `example` には出力せず、
Request / Response の利用コンテキストに対応する Example として出力する。

Request の場合は、例えば以下へ出力する。

```yaml
requestBody:
  content:
    application/json:
      schema:
        $ref: "#/components/schemas/Order.WithDetails"
      example:
        orderPic: U0672
        customerId: C1234567
        urgent: true
        details:
          - productNo: P123456
            quantity: 10
            sellingPrice: 1000
```

`readOnly: true` かつ `required: true` の Property を含む Schema を Request で利用する場合でも、
Request Example では当該 Property を省略できる。

Resource / Variant Schema は Request / Response の双方で共通利用し、
Request / Response の利用コンテキストのみを理由として専用 Schema を生成しない。

### 13.7 Built-in Generation

Resolved Model の `builtIn` は、ARIADNE が所有する OpenAPI 定義へ変換するための識別子として使用する。

#### Standard Error

```yaml
default:
  builtIn: error
```

は、ARIADNE 標準の Error Response / Error Schema を使用した OpenAPI `default` Response へ変換する。

Standard Error の具体的な Component 名、Schema、 Media Type、Example は OpenAPI Generation の built-in 定義として管理する。

#### Pagination

```yaml
builtIn: pagination
```

を持つ `limit` / `offset` は Query Parameter へ、 `Has-More` は Response Header へ変換する。

#### Location

```yaml
Location:
  builtIn: location
  example: /orders/ORD-20230827-001
```

は以下の OpenAPI Response Header `Location` へ変換する。

```yaml
Location:
  description: 登録されたリソースのURI。
  schema:
    type: string
    format: uri-reference
  example: /orders/ORD-20230827-001
```

`Location` は `components/headers` には生成せず、対象 Response の `headers` へ inline で生成する。

#### Health

`builtIn: health` Operation は、
Resolved Model に確定済みの `/health` Operation を通常の OpenAPI Operation として生成する。

#### Version

`builtIn: version` Response は、 Application Version のみを返す OpenAPI Response Schema へ変換する。

Application Version は、Semantic Version として解釈可能な Git Tag から Build 時に取得する。

### 13.8 externalDocs

Resolved Model に保持された `externalDocs` は、 OpenAPI External Documentation Object へ変換する。

Source / Resolved Model では、補足文書の Source of Truth である Markdown を参照する。

```yaml
externalDocs:
  url: ./docs/order-create.md
  description: 受注登録の詳細仕様
```

OpenAPI Generation では、Markdown を HTML へ変換した成果物を生成し、
OpenAPI の `externalDocs.url` は生成された HTML を参照する。

```yaml
externalDocs:
  url: ./docs/order-create.html
  description: 受注登録の詳細仕様
```

Markdown の相対パスは Service Root を基準として解決する。

補足文書は以下の対応で生成する。

```text
src/api/services/{service-id}/docs/*.md
                    ↓
dist/api/oas/{service-id}/docs/*.html
```

Prototype では Markdown → HTML 変換に `marked` を使用する。

生成 HTML は共通テンプレートを使用し、Markdown 内の最初の H1 を HTML の `title` とする。
H1 が存在しない場合は補足文書の生成を Error とする。

### 13.9 Generation の境界

OpenAPI Generation は、

> **Resolved Model で確定した ARIADNE の API 意味モデルを、
> OpenAPI 3.1 の物理表現へ写像する処理**

と位置付ける。

Generation では以下を行わない。

- Variant / Parent Variant の意味解決
- API Usage Override の適用
- Parameter Usage の決定
- Service `requestHeaders` の展開
- Pagination を使用するかどうかの判断
- Success Status の決定
- Location Header を付与するかどうかの判断
- Standard Error を付与するかどうかの判断
- OperationId の生成
- Operation tag の決定
- `/health` / `/version` を追加するかどうかの判断
- API / Schema origin の決定

これらはすべて Resolve の責務とする。

OpenAPI Generation に残すのは、

- OpenAPI Object への構造変換
- OpenAPI `$ref` の生成
- Component 化
- built-in の具体的な OAS 表現への変換
- Build / Environment 情報の注入

である。

------------------------------------------------------------------------

## 14. API Document

生成した OpenAPI 3.1 から、人間が参照する API Document を生成する。

Prototype では Redocly を利用し、
OpenAPI 3.1 の妥当性確認と API Document の目視確認を行う。

### 14.1 位置付け

API Document は Source of Truth ではない。

成果物の関係は以下とする。

```text
ARIADNE Source
      ↓
Raw Model
      ↓
Resolved Model
      ↓
OpenAPI 3.1
      ↓
API Document
```

ARIADNE Source を正本とし、
Raw Model / Resolved Model / OpenAPI 3.1 / API Document は
すべて再生成可能な成果物とする。

OpenAPI 3.1 は OpenAPI Validator により機械的に検証する。

API Document は、生成された API が人間から見て
意図した構造・表現になっていることを確認するためにも利用する。

### 14.2 Prototype における生成

Prototype では Redocly CLI を利用する。

主な用途は以下とする。

- OpenAPI 3.1 の Validation
- API Document の静的 HTML 生成
- 生成された API Document の目視確認

Prototype では ReDoc 固有の表示カスタマイズを最小限とし、
まず OpenAPI 3.1 の標準的な表現を確認する。

OpenAPI 3.1 の作成中は、

```text
OAS 作成
  ↓
OpenAPI Validation
  ↓
externalDocs HTML 生成
  ↓
ReDoc 生成
  ↓
Browser 確認
```

を繰り返し、OpenAPI Generation の設計を検証する。

Prototype では Service 単位で上記処理を実行する。

OpenAPI Validation が Error となった場合は、古い API Document を誤って参照しないよう、
当該 Service の既存 `redoc.html` および生成済み `docs/` を削除する。

### 14.3 成果物配置

Service 単位の OpenAPI 成果物は以下へ配置する。

```text
dist/
└─ api/
   └─ oas/
      └─ order-management/
         ├─ openapi.yaml
         ├─ redoc.html
         └─ docs/
            └─ order-create.html
```

- `openapi.yaml`
  - 生成された OpenAPI 3.1
- `redoc.html`
  - OpenAPI 3.1 から生成した ReDoc API Document
- `docs/`
  - `externalDocs` から参照する生成済み HTML

Service 単位のディレクトリを、
ローカル参照および Web 公開の双方で利用可能な自己完結した API Document 成果物とする。

### 14.4 externalDocs

`externalDocs` の補足文書の Source of Truth は、Service Root 配下の `docs/` に配置した Markdown とする。

例:

```text
src/api/services/order-management/
└─ docs/
   └─ order-create.md
```

API Document 生成時に Markdown を HTML へ変換し、Service の OAS 成果物ディレクトリへ生成する。

```text
src/api/services/order-management/docs/order-create.md
                         ↓
dist/api/oas/order-management/docs/order-create.html
```

OpenAPI 3.1 の `externalDocs.url` は、生成された HTML を相対 URL で参照する。

```yaml
externalDocs:
  url: ./docs/order-create.html
  description: 受注登録の詳細仕様
```

これにより、Service 単位の OAS 成果物だけで ReDoc と externalDocs を自己完結して参照できる。

### 14.5 Runtime との境界

Prototype では Redocly CLI を
OpenAPI Validation / API Document 生成のための開発ツールとして利用する。

最終的な ARIADNE では、
利用者が Wails Application 上で API を設計し、
生成された OpenAPI から API Document を参照できることを要件とする。

最終利用者に Node.js / npm / Redocly CLI 等の
個別インストールは要求しない。

Wails Application への API Document 表示機能の組み込み方式は、
後続 Phase で決定する。

### 14.6 Development Tools

生成した OpenAPI 3.1 は、API Document の生成だけでなく、Backend / Frontend の開発支援環境からも利用する。

Prototype では以下を想定する。

| Tool | 主な利用者 | 用途 |
| --- | --- | --- |
| ReDoc | API 利用者 / 開発者 | API 仕様の参照 |
| Swagger UI | Backend 実装者 | OAS を利用した API 実装・動作確認 |
| Mock Server | Frontend 実装者 | Backend 実装前の API 呼び出し・画面開発 |

Swagger UI / Mock Server は ARIADNE Source / Raw Model / Resolved Model に専用定義を持たない。

いずれも生成済み OpenAPI 3.1 を入力として利用する。

```text
ARIADNE Source
      ↓
Raw Model
      ↓
Resolved Model
      ↓
OpenAPI 3.1
      ├─ ReDoc
      ├─ Swagger UI
      └─ Mock Server
```

Swagger UI / Mock Server は Docker Container として起動し、 Docker Compose からまとめて起動可能な開発支援環境とする。

Prototype における Mock Server には Prism を利用する。

```text
OpenAPI 3.1
      ↓
Prism Mock Server
      ↓
Frontend Application
```

Mock Server 固有のレスポンス定義等を ARIADNE Source に追加せず、
OpenAPI 3.1 に定義された Schema / Example 等を Mock Server が利用する。

将来的な ARIADNE Development Environment では、 Swagger UI / Mock Server に加えて、
Phase 4 で生成する DDL を利用した PostgreSQL も Docker Compose から起動可能とする。

```text
ARIADNE Development Environment

OpenAPI 3.1
  ├─ Swagger UI
  └─ Mock Server

DDL
  └─ PostgreSQL
```

PostgreSQL の初期化方法および DDL 適用方式は Phase 4 で設計する。
Database Migration は Prototype Phase 4 の必須要件とはしない。

------------------------------------------------------------------------

## 15. Phase 2 / Phase 4 との境界

ARIADNE の責務分離は以下とする。

``` text
Phase 2
Element / Scalar Value Definition
        ↓
Phase 3
Resource / Parameter / API Composition
        ↓
Phase 4
Database Definition
```

### 15.1 Phase 2 → Phase 3

Phase 2 は「値は何者か」を定義する。

例:

- 型
- 桁
- Format
- Range
- Pattern
- Precision / Scale
- Enum
- API Identifier として利用可能か

Phase 3 は Element を参照し、「その値を API
のどこで、どの意味で利用するか」を定義する。

### 15.2 Phase 3 → Phase 4

Phase 3 は API 契約を管理し、Database の PK / Index / Null Constraint
等を定義しない。

Phase 4 は Phase 2 Element を利用して Database Definition
を構成する予定とする。

API Identifier と Database Primary Key は同義ではない。

------------------------------------------------------------------------

## 16. Phase 3 Completion

Prototype Phase 3 は現在進行中である。

### 現時点の進捗

- [x] Phase 3 の基本責務整理
- [x] Service モデル
- [x] Parameter Definition モデル
- [x] Resource / Property モデル
- [x] Variant モデル
- [x] SubResource / parentVariants モデル
- [x] API / Operation 基本表記
- [x] Path / Query / Header Parameter 方針
- [x] Pagination 基本方針
- [x] HTTP Method / Success Status 基本方針
- [x] Action モデル
- [x] Custom モデル
- [x] externalDocs 基本方針
- [x] Source 配置方針
- [x] Raw Model 基本構造
- [x] Source → Raw Model の変換境界
- [x] Validation Architecture
- [x] File Validation 基本 Rule
- [x] API Validation / Service Validation の責務整理
- [x] Resolved Model 詳細構造
- [x] API Usage Override の Resolved Schema 化
- [x] Parameter Usage の Resolved Model
- [x] Pagination の具体的な展開仕様
- [x] OperationId 生成規則
- [x] Resolved Model 上の Path / Operation 構造
- [x] OAS 3.1 Generation 詳細仕様
- [x] 受注サービス OAS 3.1 サンプル作成
- [x] OAS 成果物の配置方針
- [x] ReDoc による API Document 生成方針
- [x] Swagger UI の利用方針
- [x] Mock Server の利用方針
- [x] externalDocs の OAS / ReDoc 成果物への取り込み方針
- [x] Prototype Task API の ARIADNE Source 定義
- [x] Prototype Task API の Raw Model 作成
- [x] Prototype Task API の Resolved Model 作成
- [x] Prototype Task API の OAS 3.1 作成
- [ ] Validation Rule ID の最終整理
- [ ] Phase 3 ドキュメント最終更新
- [ ] Prototype Phase 3 最終レビュー

------------------------------------------------------------------------

**Prototype Phase 3：API定義 YAML** --- IN PROGRESS
