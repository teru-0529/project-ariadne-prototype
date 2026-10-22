<script lang="ts">
  import AriadneField from "./AriadneField.svelte";
  import type { FieldVariant } from "./AriadneField.svelte";

  type Props = {
    variant: FieldVariant;
    label?: string;
    required?: boolean;
    value: string;
    lines?: number;
  };

  let { variant, label = "", required = false, value, lines = 3 }: Props = $props();
</script>

<AriadneField {variant} {label} {required}>
  {#if variant === "reference"}
    <div class="textarea-value" style:--lines={lines}>
      {value}
    </div>
  {:else if variant === "edit"}
    <textarea aria-label={label} {value} style:--lines={lines}></textarea>
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

  /* INFO: Referenceのスクロール領域のキーボードFocus */
  .textarea-value:focus-visible {
    outline: 2px solid rgba(59, 111, 189, 0.65);
    outline-offset: -2px;
    border-radius: 3px;
  }
</style>
