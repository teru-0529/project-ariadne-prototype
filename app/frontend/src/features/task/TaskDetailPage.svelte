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

<AriadnePage title="Task 詳細" headerAction={modeToggle}>
  {#if mode === "reference"}
    <TaskView {task} />
  {:else}
    <TaskEntry showId={true} id={task.ID} submitLabel="Edit Task" {initialValues} />
  {/if}
</AriadnePage>
