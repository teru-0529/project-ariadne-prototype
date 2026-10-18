<script lang="ts">
  import AriadneField from "./AriadneField.svelte";
  import type { FieldVariant } from "./AriadneField.svelte";

  export let variant: FieldVariant;
  export let label: string = "";
  export let value: string;
  export let lines: number = 3;
</script>

<AriadneField {variant} {label}>
  {#if variant === "reference"}
    <div class="textarea-value" style:--lines={lines}>
      {value}
    </div>
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

  /* INFO: Referenceのスクロール領域のキーボードFocus */
  .textarea-value:focus-visible {
    outline: 2px solid rgba(59, 111, 189, 0.65);
    outline-offset: -2px;
    border-radius: 3px;
  }
</style>
