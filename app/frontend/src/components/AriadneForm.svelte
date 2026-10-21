<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = {
    children: Snippet;
    submitLabel: string;
  };

  let { children, submitLabel }: Props = $props();
</script>

<form>
  {@render children()}

  <button type="button" disabled>
    {submitLabel}
  </button>
</form>
