<script lang="ts">
  import AriadneField from "./AriadneField.svelte";
  import type { FieldVariant } from "./AriadneField.svelte";
  import Tooltip from "./AriadneTooltip.svelte";

  type Props = {
    variant: FieldVariant;
    label?: string;
    required?: boolean;
    value: string;
    maxLength: number;
    lines?: number;
  };

  let { variant, label = "", required = false, value, maxLength, lines = 1 }: Props = $props();

  let composing = false;

  const initialValue = () => value;
  let currentValue = $state(initialValue());

  // FUNCTION: 値変更イベント
  function handleValueChange(nextValue: string) {
    currentValue = nextValue;
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
    <div class="text-value" style:--lines={lines}>
      {value}
    </div>
  {:else if variant === "edit"}
    <textarea
      aria-label={label}
      {value}
      rows={lines}
      maxlength={maxLength}
      onkeydown={(event) => {
        if (event.key === "Enter" && !event.isComposing) {
          event.preventDefault();
        }
      }}
      oncompositionstart={() => {
        composing = true;
      }}
      oncompositionend={(event) => {
        composing = false;
        event.currentTarget.value = event.currentTarget.value.replace(/[\r\n]/g, "");
        handleValueChange(event.currentTarget.value);
      }}
      oninput={(event) => {
        if (composing) return;

        const element = event.currentTarget;
        element.value = element.value.replace(/[\r\n]/g, "");
        handleValueChange(element.value);
      }}
      onfocusout={handleFocusOut}></textarea>
  {:else if variant === "cell"}
    <Tooltip text={value} />
  {/if}
  <!-- TODO: ダミー -->
  {#if variant === "edit"}
    <div>変更検知：{currentValue}</div>
    <div>Focus Out検知：{focusOutCount}</div>
  {/if}
  <!-- TODO: ダミー -->
</AriadneField>

<style>
  .text-value {
    font-size: 12px;
    line-height: 1.5;
    height: calc(var(--lines) * 1.5em);

    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: var(--lines);
    line-clamp: var(--lines);
    overflow: hidden;

    overflow-wrap: anywhere;
  }

  textarea {
    display: block;
    width: 100%;
    box-sizing: border-box;

    font-size: 12px;
    line-height: 1.5;
    padding: 5px 8px;

    border: 1px solid #bfc5ce;
    border-radius: 4px;
    background: #fff;

    resize: none;
    overflow-y: auto;
    overflow-x: hidden;
  }

  textarea:focus {
    outline: 2px solid rgba(59, 111, 189, 0.65);
    outline-offset: 0;
  }
</style>
