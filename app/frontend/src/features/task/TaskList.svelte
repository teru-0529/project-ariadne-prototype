<script lang="ts">
  import { TaskContext, TaskScreen } from "./TaskContext";

  import Tooltip from "../../components/Tooltip.svelte";

  // INFO: タスクコンテキスト
  export let context: TaskContext;

  // INFO: タスク一覧（表示元データ）
  const tasks = context.getData(TaskScreen.List, "taskList") ?? [];
</script>

<div class="task-list">
  <div class="task-list-header">
    <h2>Task 一覧</h2>
    <span class="task-count">全 {tasks.length} 件</span>
  </div>

  <div class="task-table-scroll">
    <table>
      <colgroup>
        <col />
        <col style="width: 150px;" />
        <col style="width: 130px;" />
      </colgroup>

      <thead>
        <tr>
          <th>タイトル</th>
          <th>優先度</th>
          <th>ステータス</th>
        </tr>
      </thead>

      <tbody>
        {#each tasks as task}
          <tr data-status={task.Status}>
            <td>
              <Tooltip text={task.Title} />
              <div class="task-id">{task.ID}</div>
            </td>
            <td class="priority">
              <span class="priority-filled">{"★".repeat(task.Priority)}</span><span class="priority-empty"
                >{"☆".repeat(5 - task.Priority)}</span
              >
            </td>
            <td>
              <span class="status-badge" data-status={task.Status}>
                {task.Status}
              </span>
            </td>
          </tr>
        {:else}
          <tr class="empty-row">
            <td colspan="3">No tasks.</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

<style>
  .task-list {
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
    min-height: 0;
  }

  .task-list-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 6px;
  }

  .task-list-header h2 {
    margin: 0;
    font-size: 1em;
    font-weight: 600;
  }

  .task-count {
    font-size: 0.75em;
    color: var(--color-text-muted);
    white-space: nowrap;
  }

  .task-table-scroll {
    flex: 1;
    min-height: 0;
    overflow: auto;
  }

  table {
    width: 100%;
    min-width: 430px;
    table-layout: fixed;
    border-collapse: collapse;
    border: 1px solid #d8dde5;
  }

  th,
  td {
    padding: 4px 10px;
    border-bottom: 1px solid #d8dde5;
    font-size: 12px;
    vertical-align: middle;
  }

  th {
    position: sticky;
    top: 0;
    z-index: 1;
    background: #f4f6f8;
    text-align: center;
  }

  tr[data-status="DONE"] {
    background: #e5e7eb;
  }

  td:first-child {
    text-align: left;
  }

  td:nth-child(2),
  td:nth-child(3) {
    text-align: center;
  }

  .empty-row td {
    height: 120px;
    text-align: center;
    color: var(--color-text-muted);
  }

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

  tbody tr:hover {
    outline: 1px solid #4f7fd3;
    outline-offset: -1px;
  }
</style>
