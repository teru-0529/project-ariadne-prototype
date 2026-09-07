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

現時点の `orders` Service の構成例を以下に示す。

``` text
src/
└─ api/
   └─ services/
      ├─ orders/
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
      └─ tasks/
         └─ ...
```

### 2.1 配置原則

- `service.yaml` は Service 自身を定義する
- `parameters.yaml` は **Service 単位**の Parameter Definition
    を定義する
- Main Resource は Resource 用ディレクトリの `main.yaml` に定義する
- SubResource / Action は Main Resource 配下で管理する
- Custom は特定 Main Resource の子ではなく **Service 配下**で管理する
- `docs/` は Service 固有の `externalDocs` 用補足文書の配置候補とする
- `externalDocs.url` の相対パスは **Service Root**
    を基準として解決する

`externalDocs` の最終的な物理配置および OAS
成果物への取り込み方法は後続工程で確定する。

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
  - requestId
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

Source / Raw Model では SubResource 側に保持し、親 Resource への統合は
Resolve 時に行う。

### 3.7 Action

Action は業務上の操作を、verb Path ではなく
**生成・受付される業務上の名詞 Resource** として表現する。

例:

- `OrderShipment`
- `Cancellation`

Action は必要に応じて Action-local Resource を持つ。

Action 名と Action-local Resource
名が同一でも許容する。両者はスコープが異なるため曖昧とはみなさない。

### 3.8 Custom

Custom は通常の Resource / SubResource / Action
のモデルでは自然に表現できない API の escape hatch とする。

Custom では以下を許容する。

- verb 的な Path
- 任意の深さの Path
- Resource 横断的な API
- Custom-local Resource
- 明示的な Response Status

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

logical entry は人間が API を整理するための論理名であり、HTTP
上の一意性は最終的に `(path, method)` で判定する。

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
| `name` | String | コンテキスト上の表示名 |
| `description` | String | コンテキスト上の説明 |
| `example` | Scalar / Object / Array | コンテキスト上の例 |

`example` は Resource Property 利用時のコンテキスト情報として指定できる。

- `element` を参照する Property：Scalar
- `resource` を参照する Property：Object
- `array` を参照する Property：Array

Object / Array 内に含まれる Element 由来の値を含め、`example` は参照先 Element の Type / Format / Constraint
を満たす必要がある。

`readOnly: true` と `writeOnly: true` の同時指定は禁止する。

### 6.2 Array

Array は `element` または `resource` のいずれか一つを参照する。

Array 自身には以下を指定できる。

- `minItems`
- `maxItems`

`required` は Array の内部ではなく Property / API Usage 側に指定する。

`minItems` / `maxItems` は 0 以上の整数とし、両方指定する場合は
`minItems <= maxItems` とする。

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

API Usage では Resource / Variant の利用コンテキストに対して Override
を指定できる。

例:

```yaml
request:
  resource: Order
  variant: WithDetails
  overrides:
    details:
      required: true
      minItems: 1
      example:
        - productNo: P123456
          quantity: 10
          sellingPrice: 1000
        - productNo: P654321
          quantity: 5
          sellingPrice: 1200
```

API Usage の `overrides` でも `example` を Override できる。

上記のように Array Property では複数要素を含む Array 全体を `example` として指定できる。

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

Service と Operation の双方から同一 Header
を重複指定することは禁止する。

### 7.4 Pagination

Collection GET では `pagination: true` により Pagination
を有効化できる。

Raw Model では `pagination: true` を保持し、built-in の `limit` /
`offset` 等への展開は Resolve で行う。

Pagination Response では、後続データの有無を示す Response Header
を生成する方針とする。具体的な Header 名は後続設計で確定する。

### 7.5 HTTP Method / Success Response

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

Main Resource の新規登録 POST では `Location` Header
を自動生成する方針とする。

SubResource / Action / Custom の POST では自動 Location を付与しない。

### 7.6 PATCH

PATCH Request では、対象 Schema の Property は Request 上すべて optional
として扱う。

PATCH 用 Variant に明示的に含めた `readOnly` Property
は更新対象として利用できる。

### 7.7 Standard Error

全 Operation に Standard Error の default Response
を自動付与する方針とする。

Raw Model では付与せず、後続の Resolve / Generation で補完する。

### 7.8 Built-in API

以下を built-in API として自動生成する方針とする。

- `/health`
- `/version`

ユーザー定義 API でこれらの Path を定義することは禁止する。

### 7.9 externalDocs

OpenAPI Schema だけでは表現しにくい業務ルールや処理上の補足には
`externalDocs` を利用する。

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

  requestId:
    element: requestId
    headerName: X-Request-ID
    description: リクエストを識別するID。
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
| `name` | 必須 | String | 表示名 |
| `description` | 任意 | String | 説明 |
| `resources` | 必須 | Map | Action-local Resource 定義 |
| `api` | 必須 | Map | API 定義 |

Local Resource は Action 内スコープで参照する。

### 8.7 Custom YAML

Custom API を定義する。

| 属性 | 必須 | 型 | 意味 |
| --- | --- | --- | --- |
| `custom` | 必須 | String | Custom 名 |
| `name` | 必須 | String | 表示名 |
| `description` | 任意 | String | 説明 |
| `resources` | 任意 | Map | Custom-local Resource 定義 |
| `api` | 必須 | Map | API 定義 |

Custom では明示的な `response.status` を指定できる。

------------------------------------------------------------------------

## 9. Raw Model

### 9.1 目的

Raw Model は、Phase 3 ARIADNE Source を **Service
単位の意味モデル**へ正規化した中間モデルである。

Source of Truth ではなく、ARIADNE Source から再生成可能な成果物とする。

Raw Model は YAML として出力し、デバッグや Golden Test
に利用できるようにする。

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
  requestHeaders:
    - parameterRef: $requestId

parameters:
  ...

resources:
  ...

actions:
  ...

customs:
  ...
```

Raw Model には `updatedAt` を持たせない。生成物の Timestamp
による不要な差分を避けるためである。

### 9.3 Source kind の正規化

| Source | Raw Model |
| --- | --- |
| `kind: service` | `service` |
| `kind: parameters` | `parameters` |
| `kind: resource` | `resources` |
| `kind: subresource` | `resources` + `parentRef` |
| `kind: action` | `actions` |
| `kind: custom` | `customs` |

Source の `kind` および `resource` / `action` / `custom`
といった識別フィールドは、Raw Model では構造そのものへ吸収する。

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
- Custom の明示 Status
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

### 11.3 Service Validation

以下は Phase 2 / Phase 3 の接続ルールとして扱う。

- すべての `element` 参照先が Phase 2 に存在する
- Path Token から正規化した semantic key が Parameter Definition に存在する
- Path Parameter が参照する Element の effective `identifier` が `true` である
- `example` に含まれる Element 由来の値が、対応する Element の Type / Format / Constraint に適合する

------------------------------------------------------------------------

## 12. Resolved Model

**Status:** DESIGN IN PROGRESS

Resolved Model は API Validation を通過した Raw Model から生成する、API
として解決済みの Service Model とする。

現時点で以下を Resolve 対象とする方針である。

- Variant の完全 Schema 化
- parentVariants の親 Resource への統合
- API Usage Override の反映
- Local Resource の解決
- Parameter Usage の解決
- 使用 Parameter Definition の抽出
- Service 共通 Header の展開
- Pagination の展開
- OperationId の生成
- HTTP Status / Header の補完
- Resource / Variant から生成する Schema 名の確定

一方、Resolved Model の API 構造については設計継続中である。

特に、

- API logical entry / Path の階層をどこまで保持するか
- OperationId 単位へ完全にフラット化するか

は未確定とする。

Path Parameter が Path 共通情報であることを踏まえ、Path 単位の構造を
Resolved Model でも保持する案を検討する。

------------------------------------------------------------------------

## 13. OpenAPI 3.1 Generation

Prototype Phase 3 の生成対象は **OpenAPI 3.1** とする。OpenAPI 3.0
は対象外とする。

現時点の生成方針は以下とする。

- ARIADNE Resource / Variant / Local Resource から
    `components/schemas` を生成する
- Parameter Definition / Usage から OAS Parameter を生成する
- Standard Error Response を全 Operation に生成する
- `/health` / `/version` を built-in Operation として生成する
- Main Resource 新規登録 POST では Location Header を生成する
- OperationId を生成する
- Method に応じた標準 Success Status を補完する
- `servers` は ARIADNE Source では管理せず、環境 / Build
    側から注入する
- `/version` の Version は Application Version とし、Git Tag
    等から供給することを想定する

詳細な Generation Rule は Resolved Model 確定後に更新する。

------------------------------------------------------------------------

## 14. Phase 2 / Phase 4 との境界

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

### 14.1 Phase 2 → Phase 3

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

### 14.2 Phase 3 → Phase 4

Phase 3 は API 契約を管理し、Database の PK / Index / Null Constraint
等を定義しない。

Phase 4 は Phase 2 Element を利用して Database Definition
を構成する予定とする。

API Identifier と Database Primary Key は同義ではない。

------------------------------------------------------------------------

## 15. Phase 3 Completion

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
- [ ] Resolved Model 詳細構造
- [ ] API Usage Override の Resolved Schema 化
- [ ] Parameter Usage の Resolved Model
- [ ] Pagination の具体的な展開仕様
- [ ] OperationId 生成規則
- [ ] Resolved Model 上の Path / Operation 構造
- [ ] OAS 3.1 Generation 詳細仕様
- [ ] Validation Rule ID の最終整理
- [ ] Prototype Phase 3 最終レビュー

------------------------------------------------------------------------

**Prototype Phase 3：API定義 YAML** --- IN PROGRESS
