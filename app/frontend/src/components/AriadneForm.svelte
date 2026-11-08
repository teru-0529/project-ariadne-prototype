<script lang="ts">
  import type { Snippet } from "svelte";
  import AriadneButton from "./AriadneButton.svelte";
  import AriadneModal from "./AriadneModal.svelte";

  type Props = {
    children: Snippet;
    submitLabel: string;
    canSubmit: boolean;
    modalTitle: string;
    modalMessage: string;
    onSubmit?: () => boolean;
    onConfirm?: () => Promise<void>;
  };

  let { children, submitLabel, canSubmit, modalTitle, modalMessage, onSubmit, onConfirm }: Props = $props();

  // INFO: Confirmation Modal表示状態
  let showConfirmation = $state(false);

  // FUNCTION: ボタン押下イベント
  function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    if (showConfirmation) {
      return;
    }

    if (onSubmit?.() !== true) {
      return;
    }

    showConfirmation = true;
  }

  // FUNCTION: Confirmation ModalのCancel処理
  function handleCancel(): void {
    showConfirmation = false;
  }

  // FUNCTION: Confirmation ModalのConfirm処理
  async function handleConfirm(): Promise<void> {
    await onConfirm?.();
    showConfirmation = false;
  }
</script>

<form class="ariadne-form" onsubmit={handleSubmit}>
  <div class="form-scroll">
    {@render children()}
  </div>

  <div class="form-actions">
    <AriadneButton label={submitLabel} variant="primary" type="submit" disabled={!canSubmit} />
  </div>
</form>

{#if showConfirmation}
  <AriadneModal
    title={modalTitle}
    message={modalMessage}
    confirmLabel="Confirm"
    cancelLabel="Cancel"
    onCancel={handleCancel}
    onConfirm={handleConfirm}
  />
{/if}

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
