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

## 2. Common Rule — Navigation from Task Create or task Edit

### CR01 未保存の変更がある場合の画面遷移確認

適用Behavior:

- S02: TaskからConfigへ移動する（移動元がCreate画面/Edit画面の場合）
- L01: Side NavigationのListへ移動する（移動元がCreate画面/Edit画面の場合）
- C01: Side NavigationのCreateへ移動する（移動元がCreate画面/Edit画面の場合）
- M02: Editモードから参照モードへ戻る

Rule:

- 未保存の変更（Changed Fields）がある場合、遷移処理の前に確認Modalを表示する
- Cancel: Modalを閉じ、画面・入力内容・Edit状態を維持する。遷移先の処理は実行しない
- Confirm: Modalを閉じ、要求された遷移先のBehaviorを実行する
- 未保存の変更がない場合は確認Modalを表示せず、遷移先のBehaviorを実行する
- 遷移先の取得に失敗した場合は、そのBehaviorのFailure Patternに従う。取得失敗を編集内容の破棄成功として扱わない
- M06の更新成功後の参照モードへの遷移は、編集破棄ではないため本Ruleの対象外

Note:

- 確認Modalの有無は共通ルールで決める。再取得・画面表示・Failureは遷移先のBehaviorが定義する
- 適用Behaviorが増えた場合は、上記の適用一覧にも追加する

---

## 3. Application Shell

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

Note:

- Create画面/Edit画面からの移動には共通ルールCR01を適用する

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

## 4. Task List / Read

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

Note:

- Create画面/Edit画面からの移動には共通ルールCR01を適用する

### L02 Taskを選択する

Start:

- Task List表示中

Action:

- Task ListからTaskを選択する

Result / Pattern:

#### Success

- 選択されたTaskをIDで再取得する
- Task参照モードで取得した内容を表示する
- 取得したTaskを参照モードの保持データとする

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

- Listで取得済みのTaskをそのまま参照画面へ渡さず、GetTask / Pickによって再取得する
- Not FoundとRead FailureはApplication上の意味が異なるため、別Patternとして扱う
- Editモードへの切り替えはM01で扱い、L02には含めない

---

## 5. Task Create

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

Note:

- Create画面/Edit画面からの移動には共通ルールCR01を適用する

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

## 6. Task Edit / Modify

### M01 参照モードからEditモードへ切り替える

Start:

- Task参照モード表示中

Action:

- Editを選択する

Result:

- 参照モードで保持しているTaskを使用し、バックエンドへの再取得は行わない
- Task Editモードを表示する
- Original ValuesとCurrent Valuesを保持Taskから初期化する
- Changed FieldsをClearする
- Single Field ErrorとCorrelation ErrorをClearする
- Correlation Validation Mode = INITIAL
- Modify Button = Disabled

Scope:

- Front

### M02 Editモードから参照モードへ戻る

Start:

- Task Editモード表示中

Action:

- 参照モードへ戻る操作を行う

Process:

- 共通ルールCR01を適用し、必要な場合は画面遷移の意思を確認する
- 遷移が許可された場合、編集中TaskのIDでバックエンドから再取得する

Result / Pattern:

#### Success

- 再取得したTaskを参照モードで表示する
- 参照モードの保持データを再取得結果に更新する
- Edit状態を終了する

#### Cancel（CR01）

- Editモード・入力内容を維持する
- バックエンドの再取得は行わない

#### Not Found

- 対象Taskが存在しないことをError Toastで通知する
- Editモード・入力内容を維持する

#### Read Failure

- Taskを再取得できないことをError Toastで通知する
- Editモード・入力内容を維持する

Scope:

- Front + Back

Note:

- 参照画面に戻す際はEdit開始時の保持データを再利用せず、GetTask / Pickで再取得する
- 確認Modalの共通仕様はCR01を参照する

### M03 入力項目からFocus Outする

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

### M04 Modify Buttonを押す

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

### M05 更新をCancelする

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

### M06 更新をConfirmする

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
- 更新したTaskをIDでバックエンドから再取得する
- 再取得したTaskを参照モードで表示し、保持データを更新する
- Edit状態を終了する

#### Update Success - Reload Failure

- Taskの更新・保存は成功しているが、参照用の再取得に失敗する
- 更新成功と再取得失敗を区別してError Toastで通知する
- 更新が失敗したものとして再実行しない
- 参照モードへの遷移は完了しない（表示状態の扱いは実装前に確定する）

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

- Success時は再取得したTaskを参照モードで表示する。Edit Baseline Resetは行わない
- Not FoundとWrite FailureはApplication上の意味が異なるため別Pattern
- 想定外ErrorはCommon Application Error Handlingで扱う

---

## Behavior Map 項目チェックリスト

Behavior Map作成時は、以下の項目がBehavior / Patternとして
網羅されていることを確認する。

| 領域        | 確認項目                                             | Behavior                |
| :---------- | :--------------------------------------------------- | :---------------------- |
| Application | Application起動                                      | S01                     |
| Navigation  | Task → Config                                        | S02                     |
| Navigation  | Config → Task                                        | S03                     |
| List        | Task一覧取得・表示                                   | L01                     |
| List        | Task 0件                                             | L01                     |
| List        | Read Failure                                         | L01                     |
| List        | Unsupported Version                                  | L01                     |
| List        | Task選択・再取得・参照モード表示                     | L02                     |
| List        | 選択Task Not Found                                   | L02                     |
| List        | 選択Task Read Failure                                | L02                     |
| Create      | Create画面初期表示                                   | C01                     |
| Create      | 初期値設定                                           | C01                     |
| Create      | 初期Create Button Disabled                           | C01                     |
| Validation  | Focus Out時の単項目Validation                        | C02 / M03               |
| Validation  | Single Field Error Set / Clear                       | C02 / M03               |
| Validation  | ACTIVE後の相関Validation再実行                       | C02 / M03               |
| Validation  | INITIALでは相関Validationを遅延                      | C02 / M03               |
| Validation  | Button押下時の全単項目Validation                     | C03 / M04               |
| Validation  | Single Field Error時は相関Validationを実行しない     | C03 / M04               |
| Validation  | 全相関Validation                                     | C03 / M04               |
| Validation  | Correlation Validation Mode INITIAL → ACTIVE         | C03 / M04               |
| Validation  | Validation Priorityに従ったError Messageの選択・表示 | Common                  |
| Correlation | priority = 5 の場合 description必須                  | C02 / C03 / M03 / M04   |
| Create      | Confirmation Modal表示                               | C03                     |
| Create      | Confirmation Modal Cancel                            | C04                     |
| Create      | Confirm / Deposit                                    | C05                     |
| Create      | Write Failure                                        | C05                     |
| Create      | Success時Info Toast                                  | C05                     |
| Create      | Success後Form Reset                                  | C05                     |
| Edit        | 参照モードからEditモードへの切り替え                 | M01                     |
| Edit        | Editモードから参照モードへ再取得して戻る             | M02                     |
| Navigation  | 編集中の画面遷移確認（共通ルール）                   | CR01（S02/L01/C01/M02） |
| Edit        | Original Values / Current Values管理                 | M03                     |
| Edit        | Changed Fields判定                                   | M03                     |
| Edit        | Original Valuesへ戻した場合Not Changed               | M03                     |
| Edit        | 変更なしの場合Modify Button Disabled                 | M03                     |
| Edit        | 変更ありの場合Modify Button Enabled                  | M03                     |
| Edit        | Confirmation Modal表示                               | M04                     |
| Edit        | Confirmation Modal Cancel                            | M05                     |
| Edit        | Changed FieldsのみPropertiesへ連携                   | M06                     |
| Edit        | Confirm / Replace                                    | M06                     |
| Edit        | 更新対象Task Not Found                               | M06                     |
| Edit        | Read Failure                                         | M06                     |
| Edit        | Write Failure                                        | M06                     |
| Edit        | Success時Info Toast                                  | M06                     |
| Edit        | 更新成功・再取得失敗の区別                           | M06                     |
| Edit        | Success後の再取得・参照モード遷移                    | M06                     |
| Error       | Application ErrorをError Toastで表示                 | Common                  |
| Read        | 正常ReadではInfo Toastを表示しない                   | L01 / L02               |
| Write       | 正常WriteではInfo Toastを表示する                    | C05 / M06               |
