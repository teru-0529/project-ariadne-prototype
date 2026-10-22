<script lang="ts">
  import AriadneForm from "../../components/AriadneForm.svelte";
  import TextField from "../../components/TextField.svelte";
  import TextareaField from "../../components/TextareaField.svelte";
  import TaskPriorityField from "./TaskPriorityField.svelte";
  import TaskStatusField from "./TaskStatusField.svelte";
</script>

<div class="entry-layout">
  <AriadneForm submitLabel="Create Task">
    <div class="entry-fields">
      <TextField variant="edit" label="タイトル" required={true} value="" lines={2} maxLength={50} />
      <TextareaField variant="edit" label="説明" value="" lines={3} />
      <TaskPriorityField variant="edit" required={true} value={3} />
      <TaskStatusField variant="edit" required={true} value="TODO" />
    </div>
  </AriadneForm>
</div>

<style>
  .entry-layout {
    width: 100%;
    height: 100%;
    min-width: 0;
    min-height: 0;
  }

  .entry-fields {
    width: 500px;
    padding: 8px 6px;
  }
</style>
