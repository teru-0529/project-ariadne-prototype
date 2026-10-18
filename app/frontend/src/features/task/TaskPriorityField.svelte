<script lang="ts">
  import AriadneField from "../../components/AriadneField.svelte";
  import type { FieldVariant } from "../../components/AriadneField.svelte";

  export let variant: FieldVariant;
  export let label: string = "";
  export let value: number;
</script>

<AriadneField {variant} {label}>
  <span class="priority" role="img" aria-label={`優先度 ${value}`}>
    <span class="priority-filled">{"★".repeat(value)}</span><span class="priority-empty">{"☆".repeat(5 - value)}</span>
  </span>
</AriadneField>

<style>
  .priority {
    font-size: 16px;
    line-height: 1;
    transform: translateY(-1px);
    display: inline-block;
    white-space: nowrap;
  }

  .priority-filled {
    color: #dc8a18;
  }

  .priority-empty {
    color: #aeb8c5;
  }
</style>
