<script lang="ts">
  import AriadnePage from "../../components/AriadnePage.svelte";
  import { TaskContext, TaskScreen } from "./TaskContext";
  import TaskView from "./TaskView.svelte";
  import TaskEntry from "./TaskEntry.svelte";
  import type { TaskValues } from "./TaskValues";
  import AriadneToggle from "../../components/AriadneToggle.svelte";
  import AriadneModal from "../../components/AriadneModal.svelte";

  import { toErrorMessage } from "../../lib/errors/ProcessedError";
  import { TaskBoundary } from "../../../bindings/ariadne/backend/adapters/wails";
  import type { EditTaskRequest } from "../../../bindings/ariadne/backend/services/public/models";

  type Props = {
    context: TaskContext;
    onSuccess: (message: string) => void;
    onError: (message: string) => void;
    onReload: () => void;
  };

  let { context, onSuccess, onError, onReload }: Props = $props();

  // INFO: TaskEntry参照
  let taskEntry = $state<TaskEntry>();

  // INFO: 離脱確認待ちのCallback
  let pendingLeaveCallback = $state<(() => void) | undefined>(undefined);

  // INFO: タスク（表示元データ）
  let task = $derived(context.getData(TaskScreen.Detail, "taskDetail"));

  // INFO: 表示モード
  let mode = $state<"reference" | "edit">("reference");

  // INFO: Ariadne Toggle に渡すリスト
  const modeOptions = [
    { value: "reference", label: "Reference" },
    { value: "edit", label: "Edit" },
  ];

  // INFO: Edit用の初期入力値
  let initialValues = $derived<TaskValues>({
    title: task.Title,
    description: task.Description,
    priority: task.Priority,
    status: task.Status,
  });

  // FUNCTION: 差分生成
  function calculateChanges(currentValues: TaskValues): Partial<TaskValues> {
    const changes: Partial<TaskValues> = {};

    if (currentValues.title !== initialValues.title) {
      changes.title = currentValues.title;
    }
    if (currentValues.description !== initialValues.description) {
      changes.description = currentValues.description;
    }
    if (currentValues.priority !== initialValues.priority) {
      changes.priority = currentValues.priority;
    }
    if (currentValues.status !== initialValues.status) {
      changes.status = currentValues.status;
    }

    return changes;
  }

  // FUNCTION: 更新用差分取得
  function getChanges(): Partial<TaskValues> {
    return calculateChanges(taskEntry!.getValues());
  }

  // FUNCTION: Task更新
  async function handleEdit(): Promise<void> {
    const changes = getChanges();

    const request: EditTaskRequest = {
      Changes: Object.entries(changes).map(([type, value]) => ({
        Type: type,
        Value: value,
      })),
    };

    try {
      await TaskBoundary.ModifyTask(task.ID, request);

      // INFO: 更新成功時の通知・再遷移
      onSuccess("Taskを更新しました。");
      onReload();
    } catch (error) {
      onError(toErrorMessage(error));
    }
  }

  // FUNCTION: 未保存変更判定
  function hasUnsavedChanges(): boolean {
    return Object.keys(getChanges()).length > 0;
  }

  // FUNCTION: 変更有無判定（Submit活性用）
  function hasChanges(currentValues: TaskValues): boolean {
    return Object.keys(calculateChanges(currentValues)).length > 0;
  }

  // FUNCTION: Page離脱
  export function leave(callback: () => void): void {
    if (mode === "reference" || !hasUnsavedChanges()) {
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

  // FUNCTION: トグル変更時のふるまい
  function handleModeSelect(value: string): void {
    if (value === "edit" && mode === "reference") {
      mode = "edit";
    } else if (value === "reference" && mode === "edit") {
      leave(onReload);
    }
  }
</script>

{#snippet modeToggle()}
  <AriadneToggle options={modeOptions} selected={mode} onSelect={handleModeSelect} />
{/snippet}

<AriadnePage title={mode === "reference" ? "Task 詳細（参照）" : "Task 詳細（更新）"} headerAction={modeToggle}>
  {#if mode === "reference"}
    <TaskView {task} />
  {:else}
    <TaskEntry
      bind:this={taskEntry}
      showId={true}
      id={task.ID}
      submitLabel="Modify Task"
      modalTitle="Task Modification"
      modalMessage="Taskを変更しますか？"
      {initialValues}
      submitCondition={hasChanges}
      onConfirm={handleEdit}
    />
  {/if}
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
