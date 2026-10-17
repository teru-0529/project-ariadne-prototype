<script lang="ts">
  import { TaskContext, TaskScreen } from "./TaskContext";
  import AriadnePage from "../../components/AriadnePage.svelte";
  import AriadneListTable from "../../components/AriadneListTable.svelte";
  import Tooltip from "../../components/Tooltip.svelte";

  import type { TaskSummary } from "../../../bindings/ariadne/backend/services/publish";

  // INFO: タスクコンテキスト
  export let context: TaskContext;

  // INFO: タスク一覧（表示元データ）
  const tasks = context.getData(TaskScreen.List, "taskList") ?? [];

  // TODO: ダミー実装（行選択時のふるまい）
  function handleRowActivate(task: TaskSummary) {
    window.alert(`Task選択（仮実装）\n${task.Title}\nID: ${task.ID}`);
  }
</script>

{#snippet titleCell(task: TaskSummary)}
  <Tooltip text={task.Title} />
  <div class="task-id">{task.ID}</div>
{/snippet}

<!-- TODO: Field実装時、星のサイズ・行間・縦位置・連続表示を維持すること -->
{#snippet priorityCell(task: TaskSummary)}
  <span class="priority">
    <span class="priority-filled">{"★".repeat(task.Priority)}</span><span class="priority-empty"
      >{"☆".repeat(5 - task.Priority)}</span
    >
  </span>
{/snippet}

{#snippet statusCell(task: TaskSummary)}
  <span class="status-badge" data-status={task.Status}>
    {task.Status}
  </span>
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

<style>
  .task-id {
    margin-top: 2px;
    font-size: 0.7em;
    color: #929baa;
  }

  .priority {
    font-size: 16px;
    line-height: 1;
    transform: translateY(-1px);
  }

  .priority-filled {
    color: #dc8a18;
  }

  .priority-empty {
    color: #aeb8c5;
  }

  .status-badge {
    display: inline-block;
    min-width: 52px;
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 600;
    line-height: 16px;
    text-align: center;
  }

  .status-badge[data-status="TODO"] {
    background: #eee4cc;
    color: #80652d;
  }

  .status-badge[data-status="DOING"] {
    background: #dbeafe;
    color: #2563eb;
  }

  .status-badge[data-status="DONE"] {
    background: #d1d5db;
    color: #4b5563;
  }
</style>
