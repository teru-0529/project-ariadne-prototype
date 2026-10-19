<script lang="ts">
  import AriadneField from "../../components/AriadneField.svelte";
  import type { FieldVariant } from "../../components/AriadneField.svelte";

  type Props = {
    variant: FieldVariant;
    label?: string;
    value: string;
  };

  let { variant, label = "", value }: Props = $props();
</script>

<AriadneField {variant} {label}>
  <span class="status-badge" class:cell={variant === "cell"} data-status={value}>
    {value}
  </span>
</AriadneField>

<style>
  .status-badge {
    display: block;
    width: fit-content;
    min-width: 52px;
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 600;
    line-height: 16px;
    text-align: center;
  }

  .status-badge.cell {
    margin-inline: auto;
  }

  .status-badge[data-status="TODO"] {
    background: #eee4cc;
    color: #80652d;
  }

  .status-badge[data-status="DOING"] {
    background: #dbeafe;
    color: #2563eb;
  }

  .status-badge[data-status="DONE"] {
    background: #d1d5db;
    color: #4b5563;
  }
</style>
