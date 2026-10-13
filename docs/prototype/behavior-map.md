# Phase 5 Behavior Map

## 1. Purpose

Behavior Mapは、 Prototype Phase
5においてTDDで保証するBehaviorを整理するためのMapとする。

Behaviorは原則として、 1つのユーザーActionを基本単位とする。

同じActionに対して結果が分岐する場合は、 Behaviorを分割せずResult
Patternとして整理する。

基本形式：

```text
Behavior
├─ Start    : Action実行前の状態
├─ Action   : ユーザーによる1つのAction
├─ Result
│   └─ Pattern : Resultの分岐
└─ Scope    : Front / Back
```

---

## 2. Application Shell

### S01 Applicationを起動する

Start:

- Applicationが起動していない

Action:

- Applicationを起動する

Result:

- Task Pageを表示する
- Task Side Navigationを表示する
- Main Areaに `Select an operation.` を表示する
- Footerを表示する
- Task Listは取得しない

Scope:

- Front

### S02 Configへ切り替える

Start:

- Task Page

Action:

- Primary NavigationのConfigを選択する

Result:

- Config Dummy Pageを表示する
- Task Side Navigationを表示しない
- Main Areaに `Configuration is not available in this prototype.`
  を表示する
- Task Featureの画面状態・取得データを破棄する

Scope:

- Front

### S03 Taskへ切り替える

Start:

- Config Dummy Page

Action:

- Primary NavigationのTaskを選択する

Result:

- Task Pageを表示する
- Task Side Navigationを表示する
- Main Areaに `Select an operation.` を表示する
- Task Listは取得しない

Scope:

- Front

Note:

- Task Featureは初期状態から開始する
- 前回のTask Featureの画面状態・取得データは引き継がない

---

## 3. Task List / Read

### L01 Task Listを選択する

Start:

- Task Page
- Main Areaは初期状態

Action:

- Side NavigationのListを選択する

Result / Pattern:

#### Success - Tasks Exist

- Task Listを取得する
- Task Listを表示する

#### Success - Empty

- Task Listを取得する
- 0件のTask Listを正常状態として表示する

#### Read Failure

- Task Listを取得できない
- Main Areaは遷移前の画面を維持する
- Side Navigationの選択状態を維持する
- Read FailureのError Toastを表示する

#### Unsupported Version

- Task Listを取得できない
- Main Areaは遷移前の画面を維持する
- Side Navigationの選択状態を維持する
- 非対応VersionであることをError Toastで通知する

Scope:

- Front + Back

Note:

- 正常なReadではInfo Toastを表示しない
- File Read / YAML Parse / Format / Invalid Task Data等、
  Application上の振る舞いが同じ原因はRead Failureとして扱う
- Unsupported Versionは運用上発生し得る状態としてRead
  Failureと区別する

### L02 Taskを選択する

Start:

- Task List表示中

Action:

- Task ListからTaskを選択する

Result / Pattern:

#### Success

- 選択されたTaskをIDで再取得する
- Task Edit画面を表示する

#### Not Found

- 選択されたTaskが存在しない
- Task List画面を維持する
- Side NavigationのList選択状態を維持する
- Error Toastを表示する

#### Read Failure

- Taskを取得できない
- Task List画面を維持する
- Side NavigationのList選択状態を維持する
- Error Toastを表示する

Scope:

- Front + Back

Note:

- Listで取得済みのTaskをそのままEditへ渡さず、 GetTask /
  Pickによって再取得する
- Not FoundとRead FailureはApplication上の意味が異なるため、
  別Patternとして扱う

---

## 4. Task Create

### C01 Createを選択する

Start:

- Task Page

Action:

- Side NavigationのCreateを選択する

Result:

- Task Create画面を表示する

Initial Values:

- title = ""
- description = ""
- priority = 3
- status = TODO

Validation State:

- Single Field Error = none
- Correlation Error = none
- Correlation Validation Mode = INITIAL

Button:

- titleが未入力のためCreate ButtonはDisabled

Scope:

- Front

Note:

- 初期表示時にはValidation Errorを表示しない
- Invalidな初期値と、表示中のValidation Errorは別概念とする

### C02 入力項目からFocus Outする

Start:

- Task Create画面で入力中

Validation State:

- Single Field Error
- Correlation Error
- Correlation Validation Mode = INITIAL / ACTIVE

Action:

- 入力項目からFocus Outする

Process:

1. 対象Fieldの単項目Validationを実行する
   - Ruleを優先順に評価する
   - 最初のErrorで対象Fieldの評価を終了する
   - ValidならSingle Field ErrorをClear
   - Invalidなら最初のErrorをSet
2. Correlation Validation ModeがACTIVEの場合
   - 対象Fieldに関連する相関Validationを再実行する
   - 関連Fieldがすべて単項目ValidなRuleのみ評価する
   - 評価可能な関連Ruleはすべて評価する
   - ValidならCorrelation ErrorをClear
   - InvalidならCorrelation ErrorをSet
   - 評価不能なら対応するCorrelation ErrorをClear
3. Create Buttonの状態を更新する

Disabled:

- 条件なし単項目必須Fieldが未入力
- OR Single Field Errorあり
- OR Correlation Errorあり

Enabled:

- 上記に該当しない

Scope:

- Front

Note:

- INITIALではFocus Out時に相関Validationを行わない
- Validation Errorの表示は、`phase-5-windows-application.md` Section 4.6 のError Message Priorityに従う

### C03 Create Buttonを押す

Start:

- Task Create画面
- Create ButtonがEnabled

Action:

- Create Buttonを押す

Process:

1. 全Fieldの単項目Validationを実行する
   - Fieldごとに最初のErrorまで評価
   - Single Field ErrorをSet / Clear
2. Single Field Errorが存在する場合
   - 相関Validationは実行しない
   - Correlation Validation Modeは変更しない
3. Single Field Errorが存在しない場合
   - 全相関Validationを実行する
   - Correlation ErrorをSet / Clearする
   - Correlation Validation ModeをACTIVEにする

Result / Pattern:

#### Valid

- Single Field Errorなし
- Correlation Errorなし
- Confirmation Modalを表示する

#### Invalid - Single Field

- Single Field Errorあり
- 相関Validationは実行しない
- Validation Errorを表示する
- Confirmation Modalを表示しない
- Create ButtonをDisabledにする
- Correlation Validation Modeは変更しない

#### Invalid - Correlation

- Single Field Errorなし
- Correlation Errorあり
- Validation Errorを表示する
- Confirmation Modalを表示しない
- Create ButtonをDisabledにする
- Correlation Validation Mode = ACTIVE

Scope:

- Front

Note:

- INITIAL → ACTIVEは、 全単項目ValidationがValidとなり、
  全相関Validationを実際に実行した時点で行う
- Buttonを押しただけではACTIVEにしない
- 相関Validationは最初のErrorで終了せず全Ruleを評価する

### C04 登録をCancelする

Start:

- Confirmation Modal表示中

Action:

- Cancelを選択する

Result:

- Taskを登録しない
- Backendを呼び出さない
- Confirmation Modalを閉じる
- Task Create画面へ戻る
- 入力内容を維持する
- Correlation Validation Mode = ACTIVEを維持する

Scope:

- Front

### C05 登録をConfirmする

Start:

- Confirmation Modal表示中

Action:

- Confirmを選択する

Result / Pattern:

#### Success

- Taskを登録する
- tasks.yamlへ保存する
- Confirmation Modalを閉じる
- Info Toastを表示する
- Task Create画面に留まる

Form Reset:

- title = ""
- description = ""
- priority = 3
- status = TODO
- Single Field ErrorをClear
- Correlation ErrorをClear
- Correlation Validation Mode = INITIAL
- Create Button = Disabled

#### Write Failure

- Taskを保存できない
- Confirmation Modalを閉じる
- Error Toastを表示する
- 入力内容を維持する
- Correlation Validation Mode = ACTIVEを維持する

Scope:

- Front + Back

Note:

- Service側Validation Error等の想定外Errorは Common Application Error
  Handlingで扱う

---

## 5. Task Edit / Modify

### M01 入力項目からFocus Outする

Start:

- Task Edit画面で編集中

Edit State:

- Original Values
- Current Values
- Changed Fields

Validation State:

- Single Field Error
- Correlation Error
- Correlation Validation Mode = INITIAL / ACTIVE

Action:

- 入力項目からFocus Outする

Process:

1. 対象Fieldの単項目Validationを実行する
   - Ruleを優先順に評価
   - 最初のErrorで終了
   - ValidならSingle Field ErrorをClear
   - InvalidならSingle Field ErrorをSet
2. Correlation Validation ModeがACTIVEの場合
   - 対象Fieldに関連する相関Validationを再実行
   - 関連Fieldがすべて単項目ValidなRuleのみ評価
   - 評価可能な関連Ruleはすべて評価
   - ValidならCorrelation ErrorをClear
   - InvalidならCorrelation ErrorをSet
   - 評価不能ならCorrelation ErrorをClear
3. Original ValuesとCurrent Valuesを比較し、 Changed Fieldsを更新する
   - Originalと異なるField → Changed
   - Originalへ戻したField → Not Changed
4. Modify Buttonの状態を更新する

Disabled:

- Changed Fieldsなし
- OR 条件なし単項目必須Fieldが未入力
- OR Single Field Errorあり
- OR Correlation Errorあり

Enabled:

- Changed Fieldsあり
- AND Validation上の問題なし

Result / Pattern:

#### Valid - Changed

- Validation Errorなし
- Changed Fieldsあり
- 他のDisabled条件なし
- Modify ButtonをEnabledにする

#### Valid - Not Changed

- Validation Errorなし
- Changed Fieldsなし
- Modify ButtonをDisabledにする

#### Invalid

- Validation Errorを表示する
- Modify ButtonをDisabledにする

Scope:

- Front

Note:

- Changed FieldsとValidation Stateは別管理とする
- INITIALではFocus Out時に相関Validationを実行しない

### M02 Modify Buttonを押す

Start:

- Task Edit画面
- Changed Fieldsあり
- Modify ButtonがEnabled

Action:

- Modify Buttonを押す

Process:

1. 全Fieldの単項目Validationを実行する
   - Fieldごとに最初のErrorまで評価
   - Single Field ErrorをSet / Clear
2. Single Field Errorが存在する場合
   - 相関Validationは実行しない
   - Correlation Validation Modeは変更しない
3. Single Field Errorが存在しない場合
   - 全相関Validationを実行する
   - Correlation ErrorをSet / Clear
   - Correlation Validation ModeをACTIVEにする

Result / Pattern:

#### Valid

- Single Field Errorなし
- Correlation Errorなし
- Confirmation Modalを表示する

#### Invalid - Single Field

- Single Field Errorあり
- 相関Validationは実行しない
- Validation Errorを表示する
- Confirmation Modalを表示しない
- Modify ButtonをDisabledにする
- Correlation Validation Modeは変更しない

#### Invalid - Correlation

- Single Field Errorなし
- Correlation Errorあり
- Validation Errorを表示する
- Confirmation Modalを表示しない
- Modify ButtonをDisabledにする
- Correlation Validation Mode = ACTIVE

Scope:

- Front

### M03 更新をCancelする

Start:

- Confirmation Modal表示中

Action:

- Cancelを選択する

Result:

- Taskを更新しない
- Backendを呼び出さない
- Confirmation Modalを閉じる
- Task Edit画面へ戻る
- Current Valuesを維持する
- Original Valuesを維持する
- Changed Fieldsを維持する
- Correlation Validation Mode = ACTIVEを維持する

Scope:

- Front

### M04 更新をConfirmする

Start:

- Confirmation Modal表示中

Action:

- Confirmを選択する

Process:

- Changed FieldsだけをPropertiesとして ModifyTask(ID,
  Properties)へ渡す
- Backendは既存Taskへ指定Propertiesを反映する

Result / Pattern:

#### Success

- Taskを更新する
- tasks.yamlへ保存する
- Confirmation Modalを閉じる
- Info Toastを表示する
- Task Edit画面に留まる

Edit Baseline Reset:

- Current Valuesを新しいOriginal Valuesとする
- Changed FieldsをClear
- Single Field ErrorをClear
- Correlation ErrorをClear
- Correlation Validation Mode = INITIAL
- Modify Button = Disabled

#### Not Found

- 更新対象Taskが存在しない
- Confirmation Modalを閉じる
- Error Toastを表示する
- Current Valuesを維持する
- Original Valuesを維持する
- Changed Fieldsを維持する
- Correlation Validation Mode = ACTIVEを維持する

#### Read Failure

- 更新対象Taskを取得できない
- Confirmation Modalを閉じる
- Error Toastを表示する
- Current Valuesを維持する
- Original Valuesを維持する
- Changed Fieldsを維持する
- Correlation Validation Mode = ACTIVEを維持する

#### Write Failure

- Taskを保存できない
- Confirmation Modalを閉じる
- Error Toastを表示する
- Current Valuesを維持する
- Original Valuesを維持する
- Changed Fieldsを維持する
- Correlation Validation Mode = ACTIVEを維持する

Scope:

- Front + Back

Note:

- SuccessはForm ResetではなくEdit Baseline Reset
- Not FoundとWrite FailureはApplication上の意味が異なるため別Pattern
- 想定外ErrorはCommon Application Error Handlingで扱う

---

## Behavior Map 項目チェックリスト

Behavior Map作成時は、以下の項目がBehavior / Patternとして
網羅されていることを確認する。

| 領域        | 確認項目                                             | Behavior              |
| :---------- | :--------------------------------------------------- | :-------------------- |
| Application | Application起動                                      | S01                   |
| Navigation  | Task → Config                                        | S02                   |
| Navigation  | Config → Task                                        | S03                   |
| List        | Task一覧取得・表示                                   | L01                   |
| List        | Task 0件                                             | L01                   |
| List        | Read Failure                                         | L01                   |
| List        | Unsupported Version                                  | L01                   |
| List        | Task選択・再取得                                     | L02                   |
| List        | 選択Task Not Found                                   | L02                   |
| List        | 選択Task Read Failure                                | L02                   |
| Create      | Create画面初期表示                                   | C01                   |
| Create      | 初期値設定                                           | C01                   |
| Create      | 初期Create Button Disabled                           | C01                   |
| Validation  | Focus Out時の単項目Validation                        | C02 / M01             |
| Validation  | Single Field Error Set / Clear                       | C02 / M01             |
| Validation  | ACTIVE後の相関Validation再実行                       | C02 / M01             |
| Validation  | INITIALでは相関Validationを遅延                      | C02 / M01             |
| Validation  | Button押下時の全単項目Validation                     | C03 / M02             |
| Validation  | Single Field Error時は相関Validationを実行しない     | C03 / M02             |
| Validation  | 全相関Validation                                     | C03 / M02             |
| Validation  | Correlation Validation Mode INITIAL → ACTIVE         | C03 / M02             |
| Validation  | Validation Priorityに従ったError Messageの選択・表示 | Common                |
| Correlation | priority = 5 の場合 description必須                  | C02 / C03 / M01 / M02 |
| Create      | Confirmation Modal表示                               | C03                   |
| Create      | Confirmation Modal Cancel                            | C04                   |
| Create      | Confirm / Deposit                                    | C05                   |
| Create      | Write Failure                                        | C05                   |
| Create      | Success時Info Toast                                  | C05                   |
| Create      | Success後Form Reset                                  | C05                   |
| Edit        | Original Values / Current Values管理                 | M01                   |
| Edit        | Changed Fields判定                                   | M01                   |
| Edit        | Original Valuesへ戻した場合Not Changed               | M01                   |
| Edit        | 変更なしの場合Modify Button Disabled                 | M01                   |
| Edit        | 変更ありの場合Modify Button Enabled                  | M01                   |
| Edit        | Confirmation Modal表示                               | M02                   |
| Edit        | Confirmation Modal Cancel                            | M03                   |
| Edit        | Changed FieldsのみPropertiesへ連携                   | M04                   |
| Edit        | Confirm / Replace                                    | M04                   |
| Edit        | 更新対象Task Not Found                               | M04                   |
| Edit        | Read Failure                                         | M04                   |
| Edit        | Write Failure                                        | M04                   |
| Edit        | Success時Info Toast                                  | M04                   |
| Edit        | Success後Edit Baseline Reset                         | M04                   |
| Error       | Application ErrorをError Toastで表示                 | Common                |
| Read        | 正常ReadではInfo Toastを表示しない                   | L01 / L02             |
| Write       | 正常WriteではInfo Toastを表示する                    | C05 / M04             |
