# Validation Priority

## 1. Task Validation Rule

Prototype Phase 5では以下のValidation Ruleを使用する。

| Field       | Single Field Validation   |
| ----------- | ------------------------- |
| title       | Required → Max Length 50  |
| description | Optional                  |
| priority    | Range 1 - 5               |
| status      | Enum: TODO / DOING / DONE |

`title` は以下の順序で評価する。

```text
Required
    ↓ OK
Max Length 50
```

Required Errorの場合は、
Max Length Validationを実行しない。

`description` はOptionalとし、
未入力を正常な値として扱う。

---

## 2. Correlation Validation

Prototype Phase 5では、
相関Validationの技術検証として以下のRuleを使用する。

```text
priority = 5
    ↓
description Required
```

`description` 自体はOptionalである。

`priority = 5` の場合のみ、
Correlation ValidationによってRequiredとなる。

Correlation Validationは、
関連するFieldがSingle Field Validation上Validな場合のみ評価する。

---

## 3. Error Message Priority

1 Fieldに表示するValidation Messageは1つとする。

表示優先順は以下とする。

```text
1. Single Field Error
2. Correlation Error
```

同一Fieldに複数のCorrelation Errorが存在する場合は、
Rule Priorityが高いErrorを表示する。

---

## 4. Focus Out

### Correlation Validation Mode = INITIAL

```text
Focus Out
    ↓
対象FieldのSingle Field Validation
    ↓
終了
```

Correlation Validationは実行しない。

### Correlation Validation Mode = ACTIVE

```text
Focus Out
    ↓
対象FieldのSingle Field Validation
    ↓
対象Fieldに関連するCorrelation Ruleを確認
    ↓
関連FieldがすべてSingle Valid
    ↓
関連するCorrelation Validationを実行
```

Correlation Ruleが評価不能となった場合は、
そのRuleに対応する既存Correlation ErrorをClearする。

---

## 5. Create / Modify Button

```text
Button Click
    ↓
全FieldのSingle Field Validation
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
         ├─ YES → Validation Error表示
         └─ NO  → Confirmation Modal
```

Correlation Validation ModeをACTIVEにするのは、
全Single Field ValidationがValidとなり、
全Correlation Validationを実際に実行した時点とする。

---

## 6. Button Disabled

### Create Button

```text
title未入力
OR
Single Field Errorあり
OR
Correlation Errorあり
```

### Modify Button

```text
Changed Fieldsなし
OR
title未入力
OR
Single Field Errorあり
OR
Correlation Errorあり
```
