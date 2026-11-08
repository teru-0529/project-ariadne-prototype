<script lang="ts">
  import AriadnePage from "../../components/AriadnePage.svelte";
  import { TaskContext, TaskScreen } from "./TaskContext";
  import TaskView from "./TaskView.svelte";
  import TaskEntry from "./TaskEntry.svelte";
  import type { TaskValues } from "./TaskValues";
  import AriadneToggle from "../../components/AriadneToggle.svelte";

  type Props = {
    context: TaskContext;
    onReload: () => void;
  };

  let { context, onReload }: Props = $props();

  // INFO: TaskEntry参照
  let taskEntry = $state<TaskEntry>();

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
    // TODO: Editの場合は判定あり
    callback();
  }

  // FUNCTION: トグル変更時のふるまい
  function handleModeSelect(value: string): void {
    if (value === "edit" && mode === "reference") {
      mode = "edit";
    }

    // TODO: Edit → ReferenceはA07-07で実装
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
      {initialValues}
      submitCondition={hasChanges}
    />
  {/if}
</AriadnePage>
