<script lang="ts">
  import { TaskContext, TaskScreen } from "./TaskContext";
  import AriadnePage from "../../components/AriadnePage.svelte";
  import AriadneListTable from "../../components/AriadneListTable.svelte";
  import TextField from "../../components/TextField.svelte";
  import SubTextField from "../../components/SubTextField.svelte";
  import TaskPriorityField from "./TaskPriorityField.svelte";
  import TaskStatusField from "./TaskStatusField.svelte";

  import type { TaskSummary } from "../../../bindings/ariadne/backend/services/publish";

  export let onTaskSelect: (id: string) => void;

  // INFO: タスクコンテキスト
  export let context: TaskContext;

  // INFO: タスク一覧（表示元データ）
  const tasks = context.getData(TaskScreen.List, "taskList");

  // INFO: 行選択時のふるまい
  function handleRowActivate(task: TaskSummary) {
    onTaskSelect(task.ID);
  }
</script>

{#snippet titleCell(task: TaskSummary)}
  <TextField variant="cell" value={task.Title} />
  <SubTextField variant="cell" value={task.ID} />
{/snippet}

{#snippet priorityCell(task: TaskSummary)}
  <TaskPriorityField variant="cell" value={task.Priority} />
{/snippet}

{#snippet statusCell(task: TaskSummary)}
  <TaskStatusField variant="cell" value={task.Status} />
{/snippet}

<AriadnePage title="Task 一覧">
  <AriadneListTable
    minWidth="430px"
    data={tasks}
    getKey={(task: TaskSummary) => task.ID}
    columns={[
      {
        key: "title",
        header: "タイトル",
        cell: titleCell,
      },
      {
        key: "priority",
        header: "優先度",
        cell: priorityCell,
        width: "150px",
        align: "center",
      },
      {
        key: "status",
        header: "ステータス",
        cell: statusCell,
        width: "130px",
        align: "center",
      },
    ]}
    getRowVariant={(task: TaskSummary) => (task.Status === "DONE" ? "muted" : undefined)}
    onRowActivate={handleRowActivate}
  />
</AriadnePage>
