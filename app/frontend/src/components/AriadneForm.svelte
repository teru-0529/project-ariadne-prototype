<script lang="ts">
  import type { Snippet } from "svelte";
  import AriadneButton from "./AriadneButton.svelte";

  type Props = {
    children: Snippet;
    submitLabel: string;
    onSubmit?: () => void;
  };

  let { children, submitLabel }: Props = $props();

  function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    // TODO:(A07/A08): ValidationとConfirmation Modalの実装時に置き換える
    window.alert("Task登録処理は未実装です。");
  }
</script>

<form class="ariadne-form" onsubmit={handleSubmit}>
  <div class="form-scroll">
    {@render children()}
  </div>

  <div class="form-actions">
    <AriadneButton label={submitLabel} variant="primary" type="submit" />
  </div>
</form>

<style>
  .ariadne-form {
    position: relative;
    width: 100%;
    height: 100%;
    min-width: 0;
    min-height: 0;
  }

  .form-scroll {
    width: 100%;
    height: 100%;
    min-width: 0;
    min-height: 0;
    overflow: auto;
  }

  .form-actions {
    position: absolute;
    right: 12px;
    bottom: 12px;
    z-index: 1;
  }
</style>
