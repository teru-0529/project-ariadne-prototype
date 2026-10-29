<script lang="ts">
  import AriadneField from "./AriadneField.svelte";
  import type { FieldVariant } from "./AriadneField.svelte";
  import type { Validity } from "../lib/validity";

  type Props = {
    variant: FieldVariant;
    label?: string;
    required?: boolean;
    value: string;
    lines?: number;
  };

  let { variant, label = "", required = false, value, lines = 3 }: Props = $props();

  // INFO: 値
  const initialValue = () => value;
  let currentValue = $state(initialValue());

  // INFO: Validity
  const initialValidity = checkRequired(initialValue());
  let validity = $state<Validity>(initialValidity);

  // FUNCTION: 必須チェック
  function checkRequired(value: string): Validity {
    return required && value.trim() === "" ? "INVALID" : "UNKNOWN";
  }

  // FUNCTION: 値変更イベント
  function handleValueChange(nextValue: string) {
    currentValue = nextValue;
    validity = checkRequired(currentValue);
  }

  // TODO: ダミー
  let focusOutCount = $state(0);

  // FUNCTION: フォーカスアウトイベント
  function handleFocusOut() {
    focusOutCount += 1;
  }
</script>

<AriadneField {variant} {label} {required}>
  {#if variant === "reference"}
    <div class="textarea-value" style:--lines={lines}>
      {value}
    </div>
  {:else if variant === "edit"}
    <textarea
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
  <!-- TODO: ダミー -->
  {#if variant === "edit"}
    <div>validity：{validity}</div>
    <div>Focus Out検知：{focusOutCount}</div>
  {/if}
  <!-- TODO: ダミー -->
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

  /* INFO: Referenceのスクロール領域のキーボードFocus */
  .textarea-value:focus-visible {
    outline: 2px solid rgba(59, 111, 189, 0.65);
    outline-offset: -2px;
    border-radius: 3px;
  }
</style>
