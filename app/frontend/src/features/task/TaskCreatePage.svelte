<script lang="ts">
  import AriadnePage from "../../components/AriadnePage.svelte";
  import TaskEntry from "./TaskEntry.svelte";
  import type { TaskValues } from "./TaskValues";
  import AriadneModal from "../../components/AriadneModal.svelte";

  // INFO: 登録時初期値
  const initialValues: TaskValues = {
    title: "",
    description: "",
    priority: 3,
    status: "TODO",
  };

  // INFO: taskEntryへの参照
  let taskEntry = $state<TaskEntry>();

  // INFO: 離脱確認待ちのCallback
  let pendingLeaveCallback = $state<(() => void) | undefined>(undefined);

  // FUNCTION: Page離脱
  export function leave(callback: () => void): void {
    // INFO: 未保存変更の有無を確認
    if (!hasUnsavedChanges()) {
      callback();
      return;
    }

    pendingLeaveCallback = callback;
  }

  // FUNCTION: 離脱確認ModalのCancel処理
  function handleLeaveCancel(): void {
    pendingLeaveCallback = undefined;
  }

  // FUNCTION: 離脱確認ModalのConfirm処理
  function handleLeaveConfirm(): void {
    const callback = pendingLeaveCallback;
    pendingLeaveCallback = undefined;
    callback?.();
  }

  // FUNCTION: 未保存変更の有無
  function hasUnsavedChanges(): boolean {
    const currentValues = taskEntry!.getValues();

    return (
      currentValues.title !== initialValues.title ||
      currentValues.description !== initialValues.description ||
      currentValues.priority !== initialValues.priority ||
      currentValues.status !== initialValues.status
    );
  }
</script>

<AriadnePage title="Task 作成">
  <TaskEntry bind:this={taskEntry} submitLabel="Create Task" {initialValues} />
</AriadnePage>

{#if pendingLeaveCallback}
  <AriadneModal
    title="Unsaved Changes"
    message="入力内容が保存されていません。画面を離れますか？"
    confirmLabel="Leave"
    cancelLabel="Cancel"
    onConfirm={handleLeaveConfirm}
    onCancel={handleLeaveCancel}
  />
{/if}
