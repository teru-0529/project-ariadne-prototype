<script lang="ts">
  import { onMount } from "svelte";
  import AriadneButton from "./AriadneButton.svelte";

  let dialog: HTMLDialogElement;

  onMount(() => {
    dialog.showModal();
  });

  type Props = {
    title: string;
    message: string;
    confirmLabel: string;
    cancelLabel: string;
    onConfirm?: () => void;
    onCancel?: () => void;
  };

  let { title, message, confirmLabel, cancelLabel, onConfirm, onCancel }: Props = $props();
</script>

<dialog bind:this={dialog} class="modal-dialog" aria-labelledby="modal-title" aria-describedby="modal-message">
  <div class="modal-title" id="modal-title">
    {title}
  </div>

  <div class="modal-message" id="modal-message">
    {message}
  </div>

  <div class="modal-actions">
    <AriadneButton label={cancelLabel} variant="secondary" onclick={onCancel} />
    <AriadneButton label={confirmLabel} variant="primary" onclick={onConfirm} />
  </div>
</dialog>

<style>
  .modal-dialog {
    position: fixed;
    box-sizing: border-box;
    width: 380px;
    max-width: calc(100vw - 32px);
    padding: 20px;
    background: #fff;
    border: 1px solid #bfc5ce;
    border-radius: 8px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.18);
    color: #25324a;
  }

  .modal-title {
    font-size: 16px;
    font-weight: 600;
    color: #25324a;
  }

  .modal-message {
    margin-top: 10px;
    margin-left: 4px;
    font-size: 13px;
    color: #354052;
    white-space: pre-wrap;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    margin-top: 36px;
  }

  .modal-dialog::backdrop {
    background: rgba(0, 0, 0, 0.35);
  }
</style>
