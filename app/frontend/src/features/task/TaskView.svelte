<script lang="ts">
  import type { Task } from "../../../bindings/ariadne/backend/services/publish";
  import TextField from "../../components/TextField.svelte";
  import TextareaField from "../../components/TextareaField.svelte";
  import TaskPriorityField from "./TaskPriorityField.svelte";
  import TaskStatusField from "./TaskStatusField.svelte";

  type Props = {
    task: Task;
  };

  let { task }: Props = $props();
</script>

<div class="detail-layout">
  <TextField variant="reference" label="ID" value={task.ID} />

  <TextField variant="reference" label="タイトル" value={task.Title} lines={2} />

  <TextareaField variant="reference" label="説明" value={task.Description} />

  <div class="priority-layout">
    <TaskPriorityField variant="reference" value={task.Priority} />
  </div>

  <div class="status-layout">
    <TaskStatusField variant="reference" value={task.Status} />
  </div>
</div>

<style>
  .detail-layout {
    padding: 8px 6px;
    width: 500px;
  }

  .priority-layout {
    width: 100px;
  }

  .status-layout {
    width: 70px;
  }
</style>
