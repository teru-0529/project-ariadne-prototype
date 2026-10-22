<script lang="ts">
  import type { Snippet } from "svelte";
  import AriadneButton from "./AriadneButton.svelte";

  type Props = {
    children: Snippet;
    submitLabel: string;
    onSubmit?: () => void;
  };

  let { children, submitLabel }: Props = $props();

  // TODO:(A07/A08): ValidationとConfirmation Modalの実装時に置き換える
  function handleSubmit() {
    window.alert("Task登録処理は未実装です。");
  }
</script>

<form onsubmit={(event) => event.preventDefault()}>
  {@render children()}

  <AriadneButton label={submitLabel} variant="primary" type="button" onclick={handleSubmit} />
</form>
