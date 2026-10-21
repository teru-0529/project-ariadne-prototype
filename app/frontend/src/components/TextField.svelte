<script lang="ts">
  import AriadneField from "./AriadneField.svelte";
  import type { FieldVariant } from "./AriadneField.svelte";
  import Tooltip from "./AriadneTooltip.svelte";

  type Props = {
    variant: FieldVariant;
    label?: string;
    required?: boolean;
    value: string;
    lines?: number;
  };

  let { variant, label = "", required = false, value, lines = 1 }: Props = $props();
</script>

<AriadneField {variant} {label} {required}>
  {#if variant === "reference"}
    <div class="text-value" style:--lines={lines}>
      {value}
    </div>
  {:else if variant === "edit"}
    <input type="text" aria-label={label} {value} />
  {:else if variant === "cell"}
    <Tooltip text={value} />
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
