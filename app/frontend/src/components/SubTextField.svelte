<script lang="ts">
  import AriadneField from "./AriadneField.svelte";
  import type { FieldVariant } from "./AriadneField.svelte";

  export let variant: FieldVariant = "cell";
  export let value: string;
</script>

<AriadneField {variant}>
  <div class="sub-text">{value}</div>
</AriadneField>

<style>
  .sub-text {
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;

    font-size: 9px;
    color: #929baa;
  }
</style>
