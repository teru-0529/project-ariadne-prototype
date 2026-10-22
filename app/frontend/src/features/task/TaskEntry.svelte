<script lang="ts">
  import AriadneForm from "../../components/AriadneForm.svelte";
  import TextField from "../../components/TextField.svelte";
  import TextareaField from "../../components/TextareaField.svelte";
  import TaskPriorityField from "./TaskPriorityField.svelte";
  import TaskStatusField from "./TaskStatusField.svelte";
</script>

<div class="entry-layout">
  <AriadneForm submitLabel="Create Task">
    <TextField variant="edit" label="タイトル" required={true} value="" lines={1} maxLength={100} />
    <TextareaField variant="edit" label="説明" value="" lines={5} />
    <TaskPriorityField variant="edit" required={true} value={3} />
    <TaskStatusField variant="edit" required={true} value="TODO" />
  </AriadneForm>
</div>

<style>
  .entry-layout {
    padding: 8px 6px;
    width: 500px;
  }
</style>
