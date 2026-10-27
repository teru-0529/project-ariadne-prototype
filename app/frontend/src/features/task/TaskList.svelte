<script lang="ts">
  import AriadneListTable from "../../components/AriadneListTable.svelte";
  import TextField from "../../components/TextField.svelte";
  import SubTextField from "../../components/SubTextField.svelte";
  import TaskPriorityField from "./TaskPriorityField.svelte";
  import TaskStatusField from "./TaskStatusField.svelte";

  import type { TaskSummary } from "../../../bindings/ariadne/backend/services/publish";

  type Props = {
    tasks: TaskSummary[];
    onTaskSelect: (id: string) => void;
  };

  let { tasks, onTaskSelect }: Props = $props();

  function handleRowActivate(task: TaskSummary) {
    onTaskSelect(task.ID);
  }
</script>

{#snippet titleCell(task: TaskSummary)}
  <TextField variant="cell" value={task.Title} maxLength={50} />
  <SubTextField variant="cell" value={task.ID} />
{/snippet}

{#snippet priorityCell(task: TaskSummary)}
  <TaskPriorityField variant="cell" value={task.Priority} />
{/snippet}

{#snippet statusCell(task: TaskSummary)}
  <TaskStatusField variant="cell" value={task.Status} />
{/snippet}

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
