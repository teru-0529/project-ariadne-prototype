<script lang="ts">
  import AriadnePage from "../../components/AriadnePage.svelte";
  import { TaskContext, TaskScreen } from "./TaskContext";
  import TaskDetailView from "./TaskDetailView.svelte";
  import AriadneToggle from "../../components/AriadneToggle.svelte";

  // INFO: Ariadne Toggle に渡すリスト
  const modeOptions = [
    { value: "reference", label: "Reference" },
    { value: "edit", label: "Edit" },
  ];

  // INFO: タスクコンテキスト
  export let context: TaskContext;

  // INFO: タスク（表示元データ）
  const task = context.getData(TaskScreen.Reference, "taskDetail");

  // INFO: トグル変更時のふるまい
  function handleModeSelect(value: string) {
    if (value === "edit") {
      // TODO: 将来はTaskManagementへ画面遷移を依頼する
      window.alert("Edit機能は未実装です。");
    }
  }
</script>

{#snippet modeToggle()}
  <!-- TODO: Edit実装時に現在の表示モードと連動させる -->
  <AriadneToggle options={modeOptions} selected="reference" onSelect={handleModeSelect} />
{/snippet}

<AriadnePage title="Task 詳細" headerAction={modeToggle}>
  <TaskDetailView {task} />
</AriadnePage>
