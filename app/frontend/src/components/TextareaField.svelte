<script lang="ts">
  import AriadneField from "./AriadneField.svelte";
  import type { FieldVariant } from "./AriadneField.svelte";

  import type { Validation, Validator } from "../lib/validity";
  import { requiredValidator, runValidators } from "../lib/validator";

  type Props = {
    variant: FieldVariant;
    label?: string;
    required?: boolean;
    value: string;
    lines?: number;
    validators?: Validator<string>[];
  };

  let { variant, label = "", required = false, value, lines = 3, validators = [] }: Props = $props();

  // INFO: 値
  const initialValue = () => value;
  let currentValue = $state(initialValue());
  let validation = $state<Validation>(evaluateOnInput(initialValue()));

  // INFO: エラー表示の活性状態
  let errorActivated = $state(false);

  // INFO: 表示用エラーメッセージ
  let displayError = $derived(errorActivated && validation.validity === "INVALID" ? validation.errorMessage : null);

  // FUNCTION: 入力中のValidation評価
  function evaluateOnInput(value: string): Validation {
    let result: Validation = { validity: "UNKNOWN" };

    if (required) {
      result = requiredValidator(value);
    }

    return result;
  }

  // FUNCTION: Validation評価
  function validate() {
    validation = runValidators(currentValue, [...(required ? [requiredValidator] : []), ...validators]);
  }

  // FUNCTION: 値変更イベント
  function handleValueChange(nextValue: string) {
    currentValue = nextValue;
    validation = evaluateOnInput(currentValue);
  }

  // TODO: ダミー
  let focusOutCount = $state(0);

  // FUNCTION: フォーカスアウトイベント
  function handleFocusOut() {
    errorActivated = true;
    validate();
    focusOutCount += 1;
  }
</script>

<!-- TODO: ダミー -->
{#if variant === "edit"}
  <div class="callback-debug">
    validity：{validation.validity}
    {#if validation.validity === "INVALID"}
      / errorMessage：{validation.errorMessage}
    {/if}
    | errorActivated：{errorActivated}
    | Focus Out検知：{focusOutCount}
  </div>
{/if}

<!-- TODO: ダミー -->
<AriadneField {variant} {label} {required} error={displayError}>
  {#if variant === "reference"}
    <div class="textarea-value" style:--lines={lines}>
      {value}
    </div>
  {:else if variant === "edit"}
    <textarea
      class:invalid={!!displayError}
      aria-label={label}
      {value}
      style:--lines={lines}
      oninput={(event) => {
        handleValueChange(event.currentTarget.value);
      }}
      onfocusout={handleFocusOut}></textarea>
  {:else}
    <span>{value}</span>
  {/if}
</AriadneField>

<style>
  .textarea-value {
    font-size: 12px;
    line-height: 1.5;
    height: calc(var(--lines) * 1.5em);

    white-space: pre-wrap;
    overflow-wrap: anywhere;
    overflow-y: scroll;
    overflow-x: hidden;
  }

  textarea {
    display: block;
    width: 100%;
    box-sizing: border-box;
    field-sizing: content;
    resize: none;

    font-size: 12px;
    line-height: 1.5;
    padding: 5px 8px;
    border: 1px solid #bfc5ce;
    border-radius: 4px;
    background: #fff;

    min-height: calc(var(--lines) * 1.5em + 12px);
  }

  textarea:focus {
    outline: 2px solid rgba(59, 111, 189, 0.65);
    outline-offset: 0;
  }

  /* INFO: エラー時：Focusの有無にかかわらず赤枠 */
  textarea.invalid,
  textarea.invalid:focus {
    outline: 2px solid #c62828;
    outline-offset: 0;
  }

  /* INFO: Referenceのスクロール領域のキーボードFocus */
  .textarea-value:focus-visible {
    outline: 2px solid rgba(59, 111, 189, 0.65);
    outline-offset: -2px;
    border-radius: 3px;
  }

  /* TODO: debug */
  .callback-debug {
    padding: 8px 6px;
    font-size: 10px;
    line-height: 1.8;
  }
</style>
