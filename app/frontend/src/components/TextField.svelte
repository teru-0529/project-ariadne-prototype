<script lang="ts">
  import AriadneField from "./AriadneField.svelte";
  import type { FieldVariant } from "./AriadneField.svelte";
  import Tooltip from "./AriadneTooltip.svelte";

  type Props = {
    variant: FieldVariant;
    label?: string;
    value: string;
    lines?: number;
  };

  let { variant, label = "", value, lines = 1 }: Props = $props();
</script>

<AriadneField {variant} {label}>
  {#if variant === "reference"}
    <div class="text-value" style:--lines={lines}>
      {value}
    </div>
  {:else if variant === "cell"}
    <Tooltip text={value} />
  {:else}
    <span>{value}</span>
  {/if}
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
</style>
