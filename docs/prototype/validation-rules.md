# Validation Rule

## 1. Single Field Validation Rule

Prototype では以下のValidation Ruleを使用する。

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

Required判定では前後の空白を除いた値を使用する。
空白のみの場合はRequired Errorとする。

ValidationのためにTrimした値で、
入力値そのものを置き換えない。

Required Errorの場合は、
Max Length Validationを実行しない。

`description` はOptionalとし、
未入力を正常な値として扱う。

---

## 2. Correlation Validation Rule

Prototype では、
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
