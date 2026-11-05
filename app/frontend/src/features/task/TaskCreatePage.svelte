<script lang="ts">
  import AriadnePage from "../../components/AriadnePage.svelte";
  import TaskEntry from "./TaskEntry.svelte";
  import type { TaskValues } from "./TaskValues";
  import AriadneModal from "../../components/AriadneModal.svelte";

  import { TaskBoundary } from "../../../bindings/ariadne/backend/adapters/wails";
  import type { CreateTaskRequest } from "../../../bindings/ariadne/backend/services/public/models";

  // INFO: 登録時初期値
  const initialValues: TaskValues = {
    title: "",
    description: "",
    priority: 3,
    status: "TODO",
  };

  type Props = {
    onSuccess: (message: string) => void;
    onCreated: () => void;
  };

  let { onSuccess, onCreated }: Props = $props();

  // INFO: taskEntryへの参照
  let taskEntry = $state<TaskEntry>();

  // INFO: 離脱確認待ちのCallback
  let pendingLeaveCallback = $state<(() => void) | undefined>(undefined);

  // FUNCTION: Task登録
  async function handleCreate(): Promise<void> {
    const values = taskEntry!.getValues();

    const request: CreateTaskRequest = {
      Title: values.title,
      Description: values.description,
      Priority: values.priority,
      Status: values.status,
    };

    try {
      await TaskBoundary.CreateTask(request);

      // INFO: 登録成功時の通知・再遷移
      onSuccess("Taskを登録しました。");
      onCreated();
    } catch {
      // TODO: エラー通知はA17-05で実装する
    }
  }

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
  <TaskEntry bind:this={taskEntry} submitLabel="Create Task" {initialValues} onConfirm={handleCreate} />
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
