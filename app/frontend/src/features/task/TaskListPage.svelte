<script lang="ts">
  import { TaskContext, TaskScreen } from "./TaskContext";
  import AriadnePage from "../../components/AriadnePage.svelte";
  import TaskList from "./TaskList.svelte";

  type Props = {
    onTaskSelect: (id: string) => void;
    context: TaskContext;
  };

  let { onTaskSelect, context }: Props = $props();

  // INFO: タスク一覧（表示元データ）
  const tasks = $derived(context.getData(TaskScreen.List, "taskList"));
</script>

<AriadnePage title="Task 一覧">
  <TaskList {tasks} {onTaskSelect} />
</AriadnePage>
